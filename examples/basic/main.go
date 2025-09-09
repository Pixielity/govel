package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"govel/packages/application/builders"
)

// Example service interfaces and implementations
type DatabaseService interface {
	Connect() error
	Query(sql string) ([]map[string]interface{}, error)
	Close() error
}

type PostgreSQLService struct {
	Host      string
	Port      int
	Database  string
	connected bool
}

func (p *PostgreSQLService) Connect() error {
	fmt.Printf("🔗 Connecting to PostgreSQL at %s:%d/%s\n", p.Host, p.Port, p.Database)
	p.connected = true
	return nil
}

func (p *PostgreSQLService) Query(sql string) ([]map[string]interface{}, error) {
	if !p.connected {
		return nil, fmt.Errorf("database not connected")
	}
	fmt.Printf("📊 Executing SQL: %s\n", sql)
	// Mock result
	return []map[string]interface{}{
		{"id": 1, "name": "John Doe", "email": "john@example.com"},
		{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
	}, nil
}

func (p *PostgreSQLService) Close() error {
	fmt.Println("🔌 Closing database connection")
	p.connected = false
	return nil
}

// Cache service
type CacheService interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
}

type RedisCache struct {
	Host string
	Port int
	data map[string]interface{}
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	if r.data == nil {
		r.data = make(map[string]interface{})
	}
	fmt.Printf("💾 Cache SET: %s = %v (TTL: %v)\n", key, value, ttl)
	r.data[key] = value
	return nil
}

func (r *RedisCache) Get(key string) (interface{}, bool) {
	if r.data == nil {
		return nil, false
	}
	value, exists := r.data[key]
	fmt.Printf("🔍 Cache GET: %s = %v (exists: %t)\n", key, value, exists)
	return value, exists
}

func (r *RedisCache) Delete(key string) error {
	if r.data != nil {
		delete(r.data, key)
		fmt.Printf("🗑️  Cache DELETE: %s\n", key)
	}
	return nil
}

// HTTP service
type HTTPService struct {
	Host string
	Port int
}

func (h *HTTPService) Start() error {
	fmt.Printf("🚀 Starting HTTP server on %s:%d\n", h.Host, h.Port)
	return nil
}

func (h *HTTPService) Stop() error {
	fmt.Println("🛑 Stopping HTTP server")
	return nil
}

func main() {
	fmt.Println("🎯 GoVel Framework - Basic Example")
	fmt.Println("==================================")

	// Create a comprehensive application using the AppBuilder with all features
	app := builders.NewApp().
		WithName("GoVel Basic Example").
		WithVersion("1.0.0").
		WithEnvironment("development").
		WithDebug(true).
		WithLocale("en").
		WithFallbackLocale("en").
		WithTimezone("UTC").
		WithBasePath(".").
		WithShutdownTimeout(10 * time.Second).
		InConsole().
		Build()

	fmt.Printf("📋 Application Information:\n")
	fmt.Printf("   Name: %s\n", app.GetName())
	fmt.Printf("   Version: %s\n", app.GetVersion())
	fmt.Printf("   Environment: %s\n", app.GetEnvironment())
	fmt.Printf("   Debug Mode: %t\n", app.IsDebug())
	fmt.Printf("   Locale: %s\n", app.GetLocale())
	fmt.Printf("   Fallback Locale: %s\n", app.GetFallbackLocale())
	fmt.Printf("   Timezone: %s\n", app.GetTimezone())
	fmt.Printf("   Running in Console: %t\n", app.IsRunningInConsole())
	fmt.Printf("   Shutdown Timeout: %v\n", app.GetShutdownTimeout())
	fmt.Println()

	// Test Configuration functionality
	fmt.Println("⚙️  Testing Configuration:")

	// Set various configuration values
	app.Set("database.host", "localhost")
	app.Set("database.port", 5432)
	app.Set("database.name", "govel_example")
	app.Set("cache.enabled", true)
	app.Set("cache.ttl", 300)
	app.Set("server.host", "0.0.0.0")
	app.Set("server.port", 8080)
	app.Set("features", []string{"auth", "api", "websockets"})

	// Retrieve and display configuration values
	fmt.Printf("   Database Host: %s\n", app.GetString("database.host", "default"))
	fmt.Printf("   Database Port: %d\n", app.GetInt("database.port", 5432))
	fmt.Printf("   Database Name: %s\n", app.GetString("database.name", ""))
	fmt.Printf("   Cache Enabled: %t\n", app.GetBool("cache.enabled", false))
	fmt.Printf("   Cache TTL: %d seconds\n", app.GetInt("cache.ttl", 60))
	fmt.Printf("   Server Host: %s\n", app.GetString("server.host", "localhost"))
	fmt.Printf("   Server Port: %d\n", app.GetInt("server.port", 3000))

	// Test slice configuration
	features := app.GetStringSlice("features", []string{})
	fmt.Printf("   Features: %v\n", features)
	fmt.Println()

	// Test Container functionality
	fmt.Println("📦 Testing Container (Dependency Injection):")

	// Register database service
	err := app.Singleton("database", func() interface{} {
		return &PostgreSQLService{
			Host:     app.GetString("database.host", "localhost"),
			Port:     app.GetInt("database.port", 5432),
			Database: app.GetString("database.name", "govel"),
		}
	})
	if err != nil {
		fmt.Printf("❌ Failed to register database service: %v\n", err)
		return
	}

	// Register cache service
	err = app.Singleton("cache", func() interface{} {
		return &RedisCache{
			Host: app.GetString("cache.host", "localhost"),
			Port: app.GetInt("cache.port", 6379),
		}
	})
	if err != nil {
		fmt.Printf("❌ Failed to register cache service: %v\n", err)
		return
	}

	// Register HTTP service
	err = app.Bind("http", func() interface{} {
		return &HTTPService{
			Host: app.GetString("server.host", "localhost"),
			Port: app.GetInt("server.port", 8080),
		}
	})
	if err != nil {
		fmt.Printf("❌ Failed to register HTTP service: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Registered services: database, cache, http\n")
	fmt.Printf("   Database bound: %t\n", app.IsBound("database"))
	fmt.Printf("   Cache bound: %t\n", app.IsBound("cache"))
	fmt.Printf("   HTTP bound: %t\n", app.IsBound("http"))
	fmt.Println()

	// Test Logger functionality
	fmt.Println("📝 Testing Logger:")

	app.Info("Application started successfully")
	app.Debug("Debug mode is enabled - detailed logging active")
	app.Warn("This is a warning message")

	// Test structured logging with fields
	app.WithField("component", "database").Info("Database service registered")
	app.WithFields(map[string]interface{}{
		"service": "cache",
		"host":    "localhost",
		"port":    6379,
	}).Info("Cache service configured")

	app.WithField("user_id", 12345).WithField("action", "login").Info("User authentication successful")
	fmt.Println()

	// Demonstrate service resolution and usage
	fmt.Println("🔧 Testing Service Resolution and Usage:")

	// Resolve and use database service
	dbService, err := app.Make("database")
	if err != nil {
		fmt.Printf("❌ Failed to resolve database service: %v\n", err)
		return
	}

	db, ok := dbService.(DatabaseService)
	if !ok {
		fmt.Println("❌ Database service type assertion failed")
		return
	}

	// Use database service
	err = db.Connect()
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		return
	}

	results, err := db.Query("SELECT * FROM users LIMIT 2")
	if err != nil {
		fmt.Printf("❌ Failed to execute query: %v\n", err)
		return
	}

	fmt.Printf("   📊 Query results: %v\n", results)

	// Resolve and use cache service
	cacheService, err := app.Make("cache")
	if err != nil {
		fmt.Printf("❌ Failed to resolve cache service: %v\n", err)
		return
	}

	cache, ok := cacheService.(CacheService)
	if !ok {
		fmt.Println("❌ Cache service type assertion failed")
		return
	}

	// Use cache service
	cache.Set("user:12345", map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}, 5*time.Minute)

	cachedUser, found := cache.Get("user:12345")
	if found {
		fmt.Printf("   💾 Cached user data: %v\n", cachedUser)
	}

	// Test singleton behavior - resolve database service again
	dbService2, err := app.Make("database")
	if err == nil {
		fmt.Printf("   🔄 Singleton test: Same instance? %t\n", dbService == dbService2)
	}

	// Resolve HTTP service
	httpService, err := app.Make("http")
	if err != nil {
		fmt.Printf("❌ Failed to resolve HTTP service: %v\n", err)
		return
	}

	http, ok := httpService.(*HTTPService)
	if !ok {
		fmt.Println("❌ HTTP service type assertion failed")
		return
	}

	http.Start()
	fmt.Println()

	// Display comprehensive application information
	fmt.Println("📊 Comprehensive Application Information:")
	
	// Use defer to catch any panics in GetApplicationInfo
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("   ⚠️  GetApplicationInfo encountered an error: %v\n", r)
				fmt.Printf("   📋 Basic info - Name: %s, Version: %s, Environment: %s\n", 
					app.GetName(), app.GetVersion(), app.GetEnvironment())
			}
		}()
		
		appInfo := app.GetApplicationInfo()
		for key, value := range appInfo {
			fmt.Printf("   %s: %v\n", key, value)
		}
	}()
	fmt.Println()

	// Test application timing
	fmt.Println("⏰ Testing Application Timing:")
	startTime := time.Now()
	app.SetStartTime(startTime)

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("   Start Time: %v\n", app.GetStartTime().Format(time.RFC3339))
	fmt.Printf("   Uptime: %v\n", app.GetUptime())
	fmt.Println()

	// Test different AppBuilder configurations
	fmt.Println("🏗️  Testing Different AppBuilder Configurations:")

	// Production configuration
	prodApp := builders.NewApp().
		ForProduction().
		WithName("Production App").
		Build()

	fmt.Printf("   Production App - Debug: %t, Environment: %s, Timeout: %v\n",
		prodApp.IsDebug(), prodApp.GetEnvironment(), prodApp.GetShutdownTimeout())

	// Testing configuration
	testApp := builders.NewApp().
		ForTesting().
		WithName("Test App").
		Build()

	fmt.Printf("   Test App - Debug: %t, Environment: %s, Testing: %t, Timeout: %v\n",
		testApp.IsDebug(), testApp.GetEnvironment(), testApp.IsRunningUnitTests(), testApp.GetShutdownTimeout())

	// Custom configuration with method chaining
	customApp := builders.NewApp().
		WithName("Custom App").
		WithVersion("2.5.1").
		ForDevelopment().
		WithLocale("fr").
		WithFallbackLocale("en").
		WithTimezone("Europe/Paris").
		InConsole().
		WithShutdownTimeout(15 * time.Second).
		Build()

	fmt.Printf("   Custom App - Name: %s, Locale: %s, Timezone: %s\n",
		customApp.GetName(), customApp.GetLocale(), customApp.GetTimezone())
	fmt.Println()

	// Set up graceful shutdown
	fmt.Println("🛡️  Setting up graceful shutdown...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown
	go func() {
		<-quit
		fmt.Println("\n🛑 Shutdown signal received...")

		app.WithField("component", "shutdown").Info("Starting graceful shutdown")

		// Cleanup services
		app.Info("Closing database connection...")
		db.Close()

		app.Info("Stopping HTTP server...")
		http.Stop()

		app.Info("Clearing cache...")
		cache.Delete("user:12345")

		app.Info("Graceful shutdown completed")
		os.Exit(0)
	}()

	fmt.Println("🎉 GoVel Framework Example Running Successfully!")
	fmt.Println("📊 All features tested:")
	fmt.Println("   ✅ AppBuilder with fluent interface")
	fmt.Println("   ✅ Application configuration and identity")
	fmt.Println("   ✅ Configuration management (strings, ints, bools, slices)")
	fmt.Println("   ✅ Dependency injection container (singletons and bindings)")
	fmt.Println("   ✅ Structured logging with fields")
	fmt.Println("   ✅ Service resolution and usage")
	fmt.Println("   ✅ Application timing and uptime")
	fmt.Println("   ✅ Environment-specific configurations")
	fmt.Println("   ✅ Graceful shutdown handling")
	fmt.Println()
	fmt.Println("🚀 Press Ctrl+C to test graceful shutdown...")

	// Keep the application running
	select {}
}
