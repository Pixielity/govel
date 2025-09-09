package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"../../internal/container"
	"../base"
)

/**
 * PostgreSQL Service Provider Example
 * 
 * This example demonstrates how to use the PostgreSQL service provider
 * in a complete application setup with proper configuration and error handling.
 */

// ExampleUsage demonstrates how to use the PostgreSQL service provider
func ExampleUsage() {
	fmt.Println("=== PostgreSQL Service Provider Example ===")

	// Create container
	c := container.NewContainer()

	// Create and register a mock config service for this example
	mockConfig := &MockConfig{
		values: map[string]interface{}{
			"database.host":               "localhost",
			"database.port":               5432,
			"database.database":           "example_db",
			"database.username":           "postgres",
			"database.password":           "password",
			"database.ssl_mode":           "disable",
			"database.max_open_connections": 20,
			"database.max_idle_connections": 10,
			"database.connection_lifetime":  "5m",
			"database.connect_timeout":      "30s",
			"database.query_timeout":        "1m",
		},
	}

	err := c.Singleton("config", func(c container.ContainerInterface) (interface{}, error) {
		return mockConfig, nil
	})
	if err != nil {
		log.Fatalf("Failed to register config service: %v", err)
	}

	// Create and register a mock logger service
	mockLogger := &MockLogger{}
	err = c.Singleton("logger", func(c container.ContainerInterface) (interface{}, error) {
		return mockLogger, nil
	})
	if err != nil {
		log.Fatalf("Failed to register logger service: %v", err)
	}

	// Create PostgreSQL service provider
	provider := NewPostgreSQLServiceProvider()

	// Register the provider
	if err := provider.Register(c); err != nil {
		log.Fatalf("Failed to register PostgreSQL provider: %v", err)
	}

	// Boot the provider
	if err := provider.Boot(c); err != nil {
		log.Fatalf("Failed to boot PostgreSQL provider: %v", err)
	}

	// Demonstrate service usage
	demonstrateServices(c)

	// Demonstrate health checks
	demonstrateHealthCheck(provider)

	// Demonstrate connection manager
	demonstrateConnectionManager(provider)

	// Clean up
	if err := provider.Terminate(c); err != nil {
		log.Printf("Error during termination: %v", err)
	}

	fmt.Println("=== Example completed successfully ===")
}

// demonstrateServices shows how to use the registered services
func demonstrateServices(c container.ContainerInterface) {
	fmt.Println("\\n--- Demonstrating Service Usage ---")

	// Get the default database connection
	dbService, err := c.Make("database")
	if err != nil {
		log.Printf("Failed to resolve database service: %v", err)
		return
	}

	// Type assert to DatabaseInterface
	db, ok := dbService.(DatabaseInterface)
	if !ok {
		log.Printf("Database service does not implement DatabaseInterface")
		return
	}

	fmt.Printf("Database connection status: %t\\n", db.IsConnected())
	fmt.Printf("Database statistics: %+v\\n", db.Stats())

	// Get connection manager
	mgService, err := c.Make("postgres.connection_manager")
	if err != nil {
		log.Printf("Failed to resolve connection manager: %v", err)
		return
	}

	mgr, ok := mgService.(*ConnectionManager)
	if !ok {
		log.Printf("Connection manager service has wrong type")
		return
	}

	connections := mgr.ListConnections()
	fmt.Printf("Managed connections: %v\\n", connections)

	// Demonstrate connection creation with custom config
	customConfig := DatabaseConfig{
		Host:            "localhost",
		Port:            5433, // Different port
		Database:        "custom_db",
		Username:        "custom_user",
		Password:        "custom_pass",
		SSLMode:         "require",
		MaxOpenConns:    15,
		MaxIdleConns:    3,
		ConnMaxLifetime: 3 * time.Minute,
		ConnectTimeout:  20 * time.Second,
		QueryTimeout:    45 * time.Second,
	}

	customDB, err := mgr.CreateConnection(customConfig)
	if err != nil {
		log.Printf("Failed to create custom connection: %v", err)
	} else {
		fmt.Printf("Custom connection created successfully: %t\\n", customDB != nil)
	}
}

// demonstrateHealthCheck shows the health check functionality
func demonstrateHealthCheck(provider *PostgreSQLServiceProvider) {
	fmt.Println("\\n--- Demonstrating Health Check ---")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	healthData := provider.HealthCheck(ctx)
	fmt.Printf("Health check results: %+v\\n", healthData)
}

// demonstrateConnectionManager shows connection management features
func demonstrateConnectionManager(provider *PostgreSQLServiceProvider) {
	fmt.Println("\\n--- Demonstrating Connection Manager ---")

	mgr := provider.GetConnectionManager()
	if mgr == nil {
		log.Printf("Connection manager not available")
		return
	}

	// List all connections
	connections := mgr.ListConnections()
	fmt.Printf("All connections: %v\\n", connections)

	// Get statistics for all connections
	stats := mgr.GetAllConnectionStats()
	for name, stat := range stats {
		fmt.Printf("Connection %s stats: Open=%d, InUse=%d, Idle=%d\\n",
			name, stat.OpenConnections, stat.InUseConnections, stat.IdleConnections)
	}

	// Health check all connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthResults := mgr.HealthCheck(ctx)
	for name, healthy := range healthResults {
		status := "unhealthy"
		if healthy {
			status = "healthy"
		}
		fmt.Printf("Connection %s: %s\\n", name, status)
	}
}

// Mock implementations for the example

// MockConfig implements the ConfigInterface for testing
type MockConfig struct {
	values map[string]interface{}
}

func (m *MockConfig) GetString(key, defaultValue string) string {
	if val, ok := m.values[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (m *MockConfig) GetInt(key string, defaultValue int) int {
	if val, ok := m.values[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return defaultValue
}

func (m *MockConfig) GetDuration(key string, defaultValue time.Duration) time.Duration {
	if val, ok := m.values[key]; ok {
		if str, ok := val.(string); ok {
			if duration, err := time.ParseDuration(str); err == nil {
				return duration
			}
		}
		if dur, ok := val.(time.Duration); ok {
			return dur
		}
	}
	return defaultValue
}

// MockLogger implements the LoggerInterface for testing
type MockLogger struct{}

func (m *MockLogger) Info(message string, fields map[string]interface{}) {
	fmt.Printf("[INFO] %s %+v\\n", message, fields)
}

func (m *MockLogger) Error(message string, fields map[string]interface{}) {
	fmt.Printf("[ERROR] %s %+v\\n", message, fields)
}

// ExampleWithRealDatabase shows how to use with a real PostgreSQL database
// Note: This requires an actual PostgreSQL instance running
func ExampleWithRealDatabase() {
	fmt.Println("\\n=== Real Database Example ===")
	
	// This example would work with a real database
	// You would need to:
	// 1. Have PostgreSQL running
	// 2. Create a database and user
	// 3. Update the configuration values

	config := DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		Database:        "testdb",
		Username:        "testuser", 
		Password:        "testpass",
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnectTimeout:  30 * time.Second,
		QueryTimeout:    30 * time.Second,
	}

	// Create database instance
	db := NewDatabase(config)
	
	// Connect (this would fail without a real database)
	ctx := context.Background()
	if err := db.Connect(ctx); err != nil {
		fmt.Printf("Connection failed (expected without real DB): %v\\n", err)
		return
	}
	
	// Example queries (would work with real database)
	rows, err := db.Query(ctx, "SELECT version()")
	if err != nil {
		fmt.Printf("Query failed: %v\\n", err)
	} else {
		fmt.Printf("Query results: %+v\\n", rows)
	}
	
	// Transaction example
	err = db.Transaction(ctx, func(tx TransactionInterface) error {
		// Execute statements within transaction
		if err := tx.Execute("CREATE TABLE IF NOT EXISTS test (id SERIAL PRIMARY KEY, name TEXT)"); err != nil {
			return err
		}
		
		if err := tx.Execute("INSERT INTO test (name) VALUES ($1)", "Test Name"); err != nil {
			return err
		}
		
		results, err := tx.Query("SELECT id, name FROM test LIMIT 5")
		if err != nil {
			return err
		}
		
		fmt.Printf("Transaction results: %+v\\n", results)
		return nil
	})
	
	if err != nil {
		fmt.Printf("Transaction failed: %v\\n", err)
	}
	
	// Clean up
	db.Disconnect()
	fmt.Println("Real database example completed")
}
