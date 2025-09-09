package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./internal/application"
	"./internal/container"
	"./internal/config"
	"./internal/logger"
	
	"./modules/base"
	"./modules/postgres"
	"./modules/redis"
)

/**
 * Complete Service Provider Example Application
 * 
 * This application demonstrates the full service provider pattern with:
 * - Multiple service providers working together
 * - Proper bootstrapping and lifecycle management
 * - Configuration management
 * - Dependency injection
 * - Health monitoring
 * - Graceful shutdown
 */

func main() {
	fmt.Println("=== Complete Service Provider Example ===")
	
	// Create application instance
	app := application.NewApplication("service-provider-example", "1.0.0")
	
	// Set up configuration
	if err := setupConfiguration(app); err != nil {
		log.Fatalf("Failed to setup configuration: %v", err)
	}
	
	// Register core service providers
	if err := registerCoreProviders(app); err != nil {
		log.Fatalf("Failed to register core providers: %v", err)
	}
	
	// Register feature service providers
	if err := registerFeatureProviders(app); err != nil {
		log.Fatalf("Failed to register feature providers: %v", err)
	}
	
	// Boot the application
	if err := app.Boot(); err != nil {
		log.Fatalf("Failed to boot application: %v", err)
	}
	
	// Demonstrate service usage
	if err := demonstrateServices(app); err != nil {
		log.Printf("Error demonstrating services: %v", err)
	}
	
	// Start health monitoring
	go startHealthMonitoring(app)
	
	// Wait for shutdown signal
	waitForShutdown(app)
	
	fmt.Println("=== Application shutdown complete ===")
}

// setupConfiguration sets up the application configuration
func setupConfiguration(app application.ApplicationInterface) error {
	fmt.Println("\\n--- Setting up Configuration ---")
	
	// Create configuration with default values and environment variable loading
	cfg := config.NewConfig()
	
	// Set default configuration values
	defaults := map[string]interface{}{
		"app.name":        "Service Provider Example",
		"app.version":     "1.0.0",
		"app.environment": "development",
		"app.debug":       true,
		
		// Database configuration
		"database.host":                  "localhost",
		"database.port":                  5432,
		"database.database":              "example_db",
		"database.username":              "postgres",
		"database.password":              "password",
		"database.ssl_mode":              "disable",
		"database.max_open_connections":  25,
		"database.max_idle_connections":  5,
		"database.connection_lifetime":   "5m",
		"database.connect_timeout":       "30s",
		"database.query_timeout":         "30s",
		
		// Redis configuration
		"redis.host":                "localhost",
		"redis.port":                6379,
		"redis.password":            "",
		"redis.database":            0,
		"redis.pool_size":           10,
		"redis.min_idle_connections": 2,
		"redis.max_retries":         3,
		"redis.dial_timeout":        "5s",
		"redis.read_timeout":        "3s",
		"redis.write_timeout":       "3s",
		"redis.pool_timeout":        "4s",
		
		// Logging configuration
		"logging.level":   "info",
		"logging.format":  "json",
		"logging.output":  "stdout",
	}
	
	for key, value := range defaults {
		cfg.Set(key, value)
	}
	
	// Load from environment variables
	cfg.LoadFromEnv("APP_")
	
	// Register configuration service
	container := app.GetContainer()
	err := container.Singleton("config", func(c container.ContainerInterface) (interface{}, error) {
		return cfg, nil
	})
	if err != nil {
		return fmt.Errorf("failed to register config service: %w", err)
	}
	
	fmt.Printf("Configuration loaded successfully\\n")
	return nil
}

// registerCoreProviders registers core service providers
func registerCoreProviders(app application.ApplicationInterface) error {
	fmt.Println("\\n--- Registering Core Service Providers ---")
	
	// Create and register logger provider
	loggerProvider := NewLoggerServiceProvider()
	if err := app.RegisterProvider(loggerProvider); err != nil {
		return fmt.Errorf("failed to register logger provider: %w", err)
	}
	
	fmt.Printf("Core providers registered successfully\\n")
	return nil
}

// registerFeatureProviders registers feature service providers
func registerFeatureProviders(app application.ApplicationInterface) error {
	fmt.Println("\\n--- Registering Feature Service Providers ---")
	
	// Create and register PostgreSQL provider
	postgresProvider := postgres.NewPostgreSQLServiceProvider()
	if err := app.RegisterProvider(postgresProvider); err != nil {
		return fmt.Errorf("failed to register PostgreSQL provider: %w", err)
	}
	
	// Create and register Redis provider
	redisProvider := redis.NewRedisServiceProvider()
	if err := app.RegisterProvider(redisProvider); err != nil {
		return fmt.Errorf("failed to register Redis provider: %w", err)
	}
	
	fmt.Printf("Feature providers registered successfully\\n")
	return nil
}

// demonstrateServices demonstrates using the registered services
func demonstrateServices(app application.ApplicationInterface) error {
	fmt.Println("\\n--- Demonstrating Service Usage ---")
	
	container := app.GetContainer()
	ctx := context.Background()
	
	// Demonstrate configuration service
	if err := demonstrateConfig(container); err != nil {
		return fmt.Errorf("config demo failed: %w", err)
	}
	
	// Demonstrate logging service
	if err := demonstrateLogging(container); err != nil {
		return fmt.Errorf("logging demo failed: %w", err)
	}
	
	// Demonstrate database service
	if err := demonstrateDatabase(ctx, container); err != nil {
		log.Printf("Database demo failed (expected if no real DB): %v", err)
	}
	
	// Demonstrate cache service
	if err := demonstrateCache(ctx, container); err != nil {
		log.Printf("Cache demo failed (expected if no Redis): %v", err)
	}
	
	fmt.Println("Service demonstrations completed")
	return nil
}

// demonstrateConfig shows configuration service usage
func demonstrateConfig(c container.ContainerInterface) error {
	fmt.Println("\\n-- Configuration Service Demo --")
	
	configService, err := c.Make("config")
	if err != nil {
		return err
	}
	
	cfg, ok := configService.(ConfigInterface)
	if !ok {
		return fmt.Errorf("config service does not implement ConfigInterface")
	}
	
	appName := cfg.GetString("app.name", "Unknown")
	appVersion := cfg.GetString("app.version", "0.0.0")
	debug := cfg.GetBool("app.debug", false)
	
	fmt.Printf("App: %s v%s (Debug: %t)\\n", appName, appVersion, debug)
	
	return nil
}

// demonstrateLogging shows logging service usage  
func demonstrateLogging(c container.ContainerInterface) error {
	fmt.Println("\\n-- Logging Service Demo --")
	
	loggerService, err := c.Make("logger")
	if err != nil {
		return err
	}
	
	log, ok := loggerService.(LoggerInterface)
	if !ok {
		return fmt.Errorf("logger service does not implement LoggerInterface")
	}
	
	log.Info("Logging service is working", map[string]interface{}{
		"component": "demo",
		"action":    "test",
	})
	
	log.Error("This is a test error message", map[string]interface{}{
		"component": "demo",
		"error":     "simulated error",
	})
	
	fmt.Println("Logging demonstration completed")
	return nil
}

// demonstrateDatabase shows database service usage
func demonstrateDatabase(ctx context.Context, c container.ContainerInterface) error {
	fmt.Println("\\n-- Database Service Demo --")
	
	dbService, err := c.Make("database")
	if err != nil {
		return fmt.Errorf("failed to get database service: %w", err)
	}
	
	db, ok := dbService.(postgres.DatabaseInterface)
	if !ok {
		return fmt.Errorf("database service does not implement DatabaseInterface")
	}
	
	fmt.Printf("Database connection status: %t\\n", db.IsConnected())
	
	// Try to connect (will fail without real database)
	if err := db.Connect(ctx); err != nil {
		fmt.Printf("Database connection failed (expected): %v\\n", err)
		return nil
	}
	
	// If connected, demonstrate some operations
	stats := db.Stats()
	fmt.Printf("Database stats: %+v\\n", stats)
	
	// Test query (example)
	results, err := db.Query(ctx, "SELECT version()")
	if err != nil {
		fmt.Printf("Query failed: %v\\n", err)
	} else {
		fmt.Printf("Query results: %+v\\n", results)
	}
	
	return nil
}

// demonstrateCache shows cache service usage
func demonstrateCache(ctx context.Context, c container.ContainerInterface) error {
	fmt.Println("\\n-- Cache Service Demo --")
	
	cacheService, err := c.Make("cache")
	if err != nil {
		return fmt.Errorf("failed to get cache service: %w", err)
	}
	
	cache, ok := cacheService.(redis.CacheInterface)
	if !ok {
		return fmt.Errorf("cache service does not implement CacheInterface")
	}
	
	// Test ping (will fail without Redis)
	if err := cache.Ping(ctx); err != nil {
		fmt.Printf("Cache ping failed (expected): %v\\n", err)
		return nil
	}
	
	// If connected, demonstrate cache operations
	fmt.Println("Cache is connected, demonstrating operations...")
	
	// Set a value
	err = cache.Set(ctx, "demo:key", "Hello, World!", 5*time.Minute)
	if err != nil {
		fmt.Printf("Cache set failed: %v\\n", err)
		return err
	}
	
	// Get the value
	value, err := cache.Get(ctx, "demo:key")
	if err != nil {
		fmt.Printf("Cache get failed: %v\\n", err)
	} else {
		fmt.Printf("Cache value: %s\\n", value)
	}
	
	// Check if exists
	exists, err := cache.Exists(ctx, "demo:key")
	if err != nil {
		fmt.Printf("Cache exists check failed: %v\\n", err)
	} else {
		fmt.Printf("Key exists: %t\\n", exists)
	}
	
	// Get stats
	stats := cache.Stats()
	fmt.Printf("Cache stats: %+v\\n", stats)
	
	return nil
}

// startHealthMonitoring starts background health monitoring
func startHealthMonitoring(app application.ApplicationInterface) {
	fmt.Println("\\n--- Starting Health Monitoring ---")
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			performHealthCheck(app)
		}
	}
}

// performHealthCheck performs health checks on all providers
func performHealthCheck(app application.ApplicationInterface) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	providers := app.GetProviders()
	healthResults := make(map[string]interface{})
	
	for _, provider := range providers {
		if healthProvider, ok := provider.(HealthCheckProvider); ok {
			health := healthProvider.HealthCheck(ctx)
			healthResults[provider.Name()] = health
		}
	}
	
	// Log overall health status
	container := app.GetContainer()
	if loggerService, err := container.Make("logger"); err == nil {
		if log, ok := loggerService.(LoggerInterface); ok {
			log.Info("Health check completed", map[string]interface{}{
				"results": healthResults,
			})
		}
	}
}

// waitForShutdown waits for shutdown signals and gracefully shuts down
func waitForShutdown(app application.ApplicationInterface) {
	fmt.Println("\\n--- Application Running (Press Ctrl+C to exit) ---")
	
	// Create channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Block until signal received
	<-sigChan
	
	fmt.Println("\\n--- Shutting down application ---")
	
	// Perform graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := app.Terminate(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

// Supporting interfaces and types

type ConfigInterface interface {
	GetString(key, defaultValue string) string
	GetInt(key string, defaultValue int) int
	GetBool(key string, defaultValue bool) bool
	GetDuration(key string, defaultValue time.Duration) time.Duration
}

type LoggerInterface interface {
	Info(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Debug(message string, fields map[string]interface{})
	Warning(message string, fields map[string]interface{})
}

type HealthCheckProvider interface {
	HealthCheck(ctx context.Context) map[string]interface{}
}

// Simple logger service provider for the example
type LoggerServiceProvider struct {
	*base.BaseProvider
}

func NewLoggerServiceProvider() *LoggerServiceProvider {
	return &LoggerServiceProvider{
		BaseProvider: base.NewBaseProvider("logger", false),
	}
}

func (p *LoggerServiceProvider) Provides() []string {
	return []string{"logger"}
}

func (p *LoggerServiceProvider) Register(c container.ContainerInterface) error {
	p.SetContainer(c)
	
	err := c.Singleton("logger", func(c container.ContainerInterface) (interface{}, error) {
		// Get config if available
		var logLevel string = "info"
		if configService, err := c.Make("config"); err == nil {
			if cfg, ok := configService.(ConfigInterface); ok {
				logLevel = cfg.GetString("logging.level", "info")
			}
		}
		
		logger := logger.NewLogger(logger.Config{
			Level:  logLevel,
			Format: "text",
			Output: "stdout",
		})
		
		return logger, nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to register logger service: %w", err)
	}
	
	p.SetRegistered(true)
	return nil
}

func (p *LoggerServiceProvider) Boot(c container.ContainerInterface) error {
	p.SetBooted(true)
	return nil
}

func (p *LoggerServiceProvider) Terminate(c container.ContainerInterface) error {
	p.SetBooted(false)
	return nil
}
