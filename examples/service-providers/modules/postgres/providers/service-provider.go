package postgres

import (
	"context"
	"fmt"
	"time"

	"../../internal/container"
	"../base"
)

/**
 * PostgreSQLServiceProvider provides database services and connection management.
 * 
 * This service provider demonstrates:
 * - Deferred service loading
 * - Configuration-based initialization
 * - Proper resource cleanup
 * - Health monitoring
 * - Connection pooling
 */
type PostgreSQLServiceProvider struct {
	*base.BaseProvider
	connectionManager *ConnectionManager
	defaultConfig     DatabaseConfig
	initialized       bool
}

// NewPostgreSQLServiceProvider creates a new PostgreSQL service provider
func NewPostgreSQLServiceProvider() *PostgreSQLServiceProvider {
	return &PostgreSQLServiceProvider{
		BaseProvider: base.NewBaseProvider("postgresql", true), // Deferred loading enabled
		defaultConfig: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			SSLMode:         "prefer",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			ConnectTimeout:  30 * time.Second,
			QueryTimeout:    30 * time.Second,
		},
	}
}

// Provides returns the services this provider offers
func (p *PostgreSQLServiceProvider) Provides() []string {
	return []string{
		"postgres.database",
		"postgres.connection_manager", 
		"postgres.default_connection",
		"database", // Alias for default connection
	}
}

// Register registers the PostgreSQL services with the container
func (p *PostgreSQLServiceProvider) Register(c container.ContainerInterface) error {
	p.SetContainer(c)

	// Register the connection manager as a singleton
	err := c.Singleton("postgres.connection_manager", func(c container.ContainerInterface) (interface{}, error) {
		if p.connectionManager == nil {
			p.connectionManager = NewConnectionManager()
		}
		return p.connectionManager, nil
	})
	if err != nil {
		return fmt.Errorf("failed to register connection manager: %w", err)
	}

	// Register database factory
	err = c.Bind("postgres.database", func(c container.ContainerInterface) (interface{}, error) {
		// Get connection manager
		mgr, err := c.Make("postgres.connection_manager")
		if err != nil {
			return nil, fmt.Errorf("failed to resolve connection manager: %w", err)
		}

		connectionManager := mgr.(*ConnectionManager)
		
		// Load configuration from container if available
		config := p.defaultConfig
		if configService, err := c.Make("config"); err == nil {
			if cfg, ok := configService.(ConfigInterface); ok {
				config = p.loadDatabaseConfig(cfg)
			}
		}

		// Create and return database connection
		return connectionManager.CreateConnection(config)
	})
	if err != nil {
		return fmt.Errorf("failed to register database factory: %w", err)
	}

	// Register default connection
	err = c.Singleton("postgres.default_connection", func(c container.ContainerInterface) (interface{}, error) {
		db, err := c.Make("postgres.database")
		if err != nil {
			return nil, fmt.Errorf("failed to create default database connection: %w", err)
		}
		
		// Connect to database
		database := db.(DatabaseInterface)
		if err := database.Connect(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		
		return database, nil
	})
	if err != nil {
		return fmt.Errorf("failed to register default connection: %w", err)
	}

	// Register database alias
	err = c.Singleton("database", func(c container.ContainerInterface) (interface{}, error) {
		return c.Make("postgres.default_connection")
	})
	if err != nil {
		return fmt.Errorf("failed to register database alias: %w", err)
	}

	p.SetRegistered(true)
	return nil
}

// Boot initializes the PostgreSQL services after all providers are registered
func (p *PostgreSQLServiceProvider) Boot(c container.ContainerInterface) error {
	if p.IsBooted() {
		return nil
	}

	// Resolve connection manager to ensure it's initialized
	mgr, err := c.Make("postgres.connection_manager")
	if err != nil {
		return fmt.Errorf("failed to resolve connection manager during boot: %w", err)
	}

	p.connectionManager = mgr.(*ConnectionManager)
	p.initialized = true
	p.SetBooted(true)

	// Log successful boot
	if logger, err := c.Make("logger"); err == nil {
		if log, ok := logger.(LoggerInterface); ok {
			log.Info("PostgreSQL service provider booted successfully", map[string]interface{}{
				"provider": "postgresql",
				"deferred": p.IsDeferred(),
			})
		}
	}

	return nil
}

// Terminate cleans up PostgreSQL resources
func (p *PostgreSQLServiceProvider) Terminate(c container.ContainerInterface) error {
	if p.connectionManager != nil {
		// Close all database connections
		if err := p.connectionManager.CloseAllConnections(); err != nil {
			// Log error but don't fail termination
			if logger, err := c.Make("logger"); err == nil {
				if log, ok := logger.(LoggerInterface); ok {
					log.Error("Error closing database connections during termination", map[string]interface{}{
						"error": err.Error(),
						"provider": "postgresql",
					})
				}
			}
		}
		p.connectionManager = nil
	}

	p.initialized = false
	p.SetBooted(false)
	return nil
}

// HealthCheck performs health checks on database connections
func (p *PostgreSQLServiceProvider) HealthCheck(ctx context.Context) map[string]interface{} {
	if !p.initialized || p.connectionManager == nil {
		return map[string]interface{}{
			"status": "not_initialized",
			"healthy": false,
		}
	}

	results := p.connectionManager.HealthCheck(ctx)
	allHealthy := true
	for _, healthy := range results {
		if !healthy {
			allHealthy = false
			break
		}
	}

	return map[string]interface{}{
		"status": "initialized",
		"healthy": allHealthy,
		"connections": results,
		"stats": p.connectionManager.GetAllConnectionStats(),
	}
}

// GetConnectionManager returns the connection manager (if initialized)
func (p *PostgreSQLServiceProvider) GetConnectionManager() *ConnectionManager {
	return p.connectionManager
}

// loadDatabaseConfig loads database configuration from the config service
func (p *PostgreSQLServiceProvider) loadDatabaseConfig(config ConfigInterface) DatabaseConfig {
	dbConfig := p.defaultConfig

	// Load database configuration with fallbacks
	if host := config.GetString("database.host", dbConfig.Host); host != "" {
		dbConfig.Host = host
	}
	if port := config.GetInt("database.port", dbConfig.Port); port != 0 {
		dbConfig.Port = port
	}
	if database := config.GetString("database.database", ""); database != "" {
		dbConfig.Database = database
	}
	if username := config.GetString("database.username", ""); username != "" {
		dbConfig.Username = username
	}
	if password := config.GetString("database.password", ""); password != "" {
		dbConfig.Password = password
	}
	if sslMode := config.GetString("database.ssl_mode", dbConfig.SSLMode); sslMode != "" {
		dbConfig.SSLMode = sslMode
	}

	// Connection pool settings
	if maxOpen := config.GetInt("database.max_open_connections", dbConfig.MaxOpenConns); maxOpen > 0 {
		dbConfig.MaxOpenConns = maxOpen
	}
	if maxIdle := config.GetInt("database.max_idle_connections", dbConfig.MaxIdleConns); maxIdle >= 0 {
		dbConfig.MaxIdleConns = maxIdle
	}

	// Timeout settings
	if lifetime := config.GetDuration("database.connection_lifetime", dbConfig.ConnMaxLifetime); lifetime > 0 {
		dbConfig.ConnMaxLifetime = lifetime
	}
	if connectTimeout := config.GetDuration("database.connect_timeout", dbConfig.ConnectTimeout); connectTimeout > 0 {
		dbConfig.ConnectTimeout = connectTimeout
	}
	if queryTimeout := config.GetDuration("database.query_timeout", dbConfig.QueryTimeout); queryTimeout > 0 {
		dbConfig.QueryTimeout = queryTimeout
	}

	return dbConfig
}

// Supporting interfaces that this provider expects from other services
type ConfigInterface interface {
	GetString(key, defaultValue string) string
	GetInt(key string, defaultValue int) int
	GetDuration(key string, defaultValue time.Duration) time.Duration
}

type LoggerInterface interface {
	Info(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
}
