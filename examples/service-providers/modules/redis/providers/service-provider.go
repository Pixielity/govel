package redis

import (
	"context"
	"fmt"
	"time"

	"../../internal/container"
	"../base"
)

/**
 * RedisServiceProvider provides Redis caching services
 * 
 * This service provider demonstrates:
 * - Multiple connection management
 * - Configuration-based initialization
 * - Support for different Redis modes (standalone, sentinel, cluster)
 * - Health monitoring and statistics
 * - Proper resource cleanup
 * - Deferred loading for optimal performance
 */
type RedisServiceProvider struct {
	*base.BaseProvider
	connectionManager *ConnectionManager
	defaultConfig     ConnectionConfig
	initialized       bool
}

// NewRedisServiceProvider creates a new Redis service provider
func NewRedisServiceProvider() *RedisServiceProvider {
	return &RedisServiceProvider{
		BaseProvider: base.NewBaseProvider("redis", true), // Deferred loading enabled
		defaultConfig: ConnectionConfig{
			Host:         "localhost",
			Port:         6379,
			Database:     0,
			PoolSize:     10,
			MinIdleConns: 2,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		},
	}
}

// Provides returns the services this provider offers
func (p *RedisServiceProvider) Provides() []string {
	return []string{
		"redis.cache",
		"redis.connection_manager",
		"redis.default_connection",
		"cache", // Alias for default connection
	}
}

// Register registers the Redis services with the container
func (p *RedisServiceProvider) Register(c container.ContainerInterface) error {
	p.SetContainer(c)

	// Register the connection manager as a singleton
	err := c.Singleton("redis.connection_manager", func(c container.ContainerInterface) (interface{}, error) {
		if p.connectionManager == nil {
			p.connectionManager = NewConnectionManager()
		}
		return p.connectionManager, nil
	})
	if err != nil {
		return fmt.Errorf("failed to register Redis connection manager: %w", err)
	}

	// Register cache factory
	err = c.Bind("redis.cache", func(c container.ContainerInterface) (interface{}, error) {
		// Get connection manager
		mgr, err := c.Make("redis.connection_manager")
		if err != nil {
			return nil, fmt.Errorf("failed to resolve connection manager: %w", err)
		}

		connectionManager := mgr.(*ConnectionManager)

		// Load configuration from container if available
		config := p.defaultConfig
		if configService, err := c.Make("config"); err == nil {
			if cfg, ok := configService.(ConfigInterface); ok {
				config = p.loadRedisConfig(cfg)
			}
		}

		// Create and return cache connection
		return connectionManager.CreateConnection("default", config)
	})
	if err != nil {
		return fmt.Errorf("failed to register Redis cache factory: %w", err)
	}

	// Register default connection
	err = c.Singleton("redis.default_connection", func(c container.ContainerInterface) (interface{}, error) {
		cache, err := c.Make("redis.cache")
		if err != nil {
			return nil, fmt.Errorf("failed to create default Redis cache: %w", err)
		}

		return cache, nil
	})
	if err != nil {
		return fmt.Errorf("failed to register default Redis connection: %w", err)
	}

	// Register cache alias
	err = c.Singleton("cache", func(c container.ContainerInterface) (interface{}, error) {
		return c.Make("redis.default_connection")
	})
	if err != nil {
		return fmt.Errorf("failed to register cache alias: %w", err)
	}

	p.SetRegistered(true)
	return nil
}

// Boot initializes the Redis services after all providers are registered
func (p *RedisServiceProvider) Boot(c container.ContainerInterface) error {
	if p.IsBooted() {
		return nil
	}

	// Resolve connection manager to ensure it's initialized
	mgr, err := c.Make("redis.connection_manager")
	if err != nil {
		return fmt.Errorf("failed to resolve Redis connection manager during boot: %w", err)
	}

	p.connectionManager = mgr.(*ConnectionManager)
	p.initialized = true
	p.SetBooted(true)

	// Log successful boot
	if logger, err := c.Make("logger"); err == nil {
		if log, ok := logger.(LoggerInterface); ok {
			log.Info("Redis service provider booted successfully", map[string]interface{}{
				"provider": "redis",
				"deferred": p.IsDeferred(),
			})
		}
	}

	return nil
}

// Terminate cleans up Redis resources
func (p *RedisServiceProvider) Terminate(c container.ContainerInterface) error {
	if p.connectionManager != nil {
		// Close all Redis connections
		if err := p.connectionManager.CloseAllConnections(); err != nil {
			// Log error but don't fail termination
			if logger, err := c.Make("logger"); err == nil {
				if log, ok := logger.(LoggerInterface); ok {
					log.Error("Error closing Redis connections during termination", map[string]interface{}{
						"error":    err.Error(),
						"provider": "redis",
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

// HealthCheck performs health checks on Redis connections
func (p *RedisServiceProvider) HealthCheck(ctx context.Context) map[string]interface{} {
	if !p.initialized || p.connectionManager == nil {
		return map[string]interface{}{
			"status":  "not_initialized",
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
		"status":      "initialized",
		"healthy":     allHealthy,
		"connections": results,
		"stats":       p.connectionManager.GetAllConnectionStats(),
	}
}

// GetConnectionManager returns the connection manager (if initialized)
func (p *RedisServiceProvider) GetConnectionManager() *ConnectionManager {
	return p.connectionManager
}

// CreateNamedConnection creates a named Redis connection with specific configuration
func (p *RedisServiceProvider) CreateNamedConnection(name string, config ConnectionConfig) (CacheInterface, error) {
	if p.connectionManager == nil {
		return nil, fmt.Errorf("Redis provider not initialized")
	}

	return p.connectionManager.CreateConnection(name, config)
}

// loadRedisConfig loads Redis configuration from the config service
func (p *RedisServiceProvider) loadRedisConfig(config ConfigInterface) ConnectionConfig {
	redisConfig := p.defaultConfig

	// Basic connection settings
	if host := config.GetString("redis.host", redisConfig.Host); host != "" {
		redisConfig.Host = host
	}
	if port := config.GetInt("redis.port", redisConfig.Port); port != 0 {
		redisConfig.Port = port
	}
	if password := config.GetString("redis.password", ""); password != "" {
		redisConfig.Password = password
	}
	if database := config.GetInt("redis.database", redisConfig.Database); database >= 0 {
		redisConfig.Database = database
	}

	// Connection pool settings
	if poolSize := config.GetInt("redis.pool_size", redisConfig.PoolSize); poolSize > 0 {
		redisConfig.PoolSize = poolSize
	}
	if minIdle := config.GetInt("redis.min_idle_connections", redisConfig.MinIdleConns); minIdle >= 0 {
		redisConfig.MinIdleConns = minIdle
	}
	if maxRetries := config.GetInt("redis.max_retries", redisConfig.MaxRetries); maxRetries >= 0 {
		redisConfig.MaxRetries = maxRetries
	}

	// Timeout settings
	if dialTimeout := config.GetDuration("redis.dial_timeout", redisConfig.DialTimeout); dialTimeout > 0 {
		redisConfig.DialTimeout = dialTimeout
	}
	if readTimeout := config.GetDuration("redis.read_timeout", redisConfig.ReadTimeout); readTimeout > 0 {
		redisConfig.ReadTimeout = readTimeout
	}
	if writeTimeout := config.GetDuration("redis.write_timeout", redisConfig.WriteTimeout); writeTimeout > 0 {
		redisConfig.WriteTimeout = writeTimeout
	}
	if poolTimeout := config.GetDuration("redis.pool_timeout", redisConfig.PoolTimeout); poolTimeout > 0 {
		redisConfig.PoolTimeout = poolTimeout
	}

	// TLS configuration
	if config.GetBool("redis.tls.enabled", false) {
		redisConfig.TLSConfig = &TLSConfig{
			Enabled:            true,
			InsecureSkipVerify: config.GetBool("redis.tls.insecure_skip_verify", false),
			CertFile:           config.GetString("redis.tls.cert_file", ""),
			KeyFile:            config.GetString("redis.tls.key_file", ""),
			CAFile:             config.GetString("redis.tls.ca_file", ""),
		}
	}

	// Sentinel configuration
	if config.GetBool("redis.sentinel.enabled", false) {
		sentinelAddrs := config.GetStringSlice("redis.sentinel.addresses", []string{})
		if len(sentinelAddrs) > 0 {
			redisConfig.SentinelConfig = &SentinelConfig{
				Enabled:          true,
				MasterName:       config.GetString("redis.sentinel.master_name", "mymaster"),
				SentinelAddrs:    sentinelAddrs,
				SentinelPassword: config.GetString("redis.sentinel.password", ""),
			}
		}
	}

	// Cluster configuration
	if config.GetBool("redis.cluster.enabled", false) {
		clusterAddrs := config.GetStringSlice("redis.cluster.addresses", []string{})
		if len(clusterAddrs) > 0 {
			redisConfig.ClusterConfig = &ClusterConfig{
				Enabled:      true,
				Addrs:        clusterAddrs,
				MaxRedirects: config.GetInt("redis.cluster.max_redirects", 3),
				ReadOnly:     config.GetBool("redis.cluster.read_only", false),
			}
		}
	}

	return redisConfig
}

// Supporting interfaces that this provider expects from other services
type ConfigInterface interface {
	GetString(key, defaultValue string) string
	GetInt(key string, defaultValue int) int
	GetBool(key string, defaultValue bool) bool
	GetDuration(key string, defaultValue time.Duration) time.Duration
	GetStringSlice(key string, defaultValue []string) []string
}

type LoggerInterface interface {
	Info(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
}
