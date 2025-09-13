package facades

import (
	facade "govel/packages/support/src"
	applicationInterfaces "govel/packages/types/src/interfaces/application/base"
	applicationTokens "govel/packages/types/src/interfaces/application"
)

// App provides a clean, static-like interface to the application's core service.
//
// This facade implements the facade pattern, providing global access to the application
// service configured in the dependency injection container. It offers a Laravel-style
// API for application lifecycle management, environment detection, service resolution,
// configuration access, and core application functionality.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved application service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent application operations across goroutines
//   - Central hub for application state and lifecycle management
//   - Built-in service container and dependency injection integration
//
// Behavior:
//   - First call: Resolves app service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if app service cannot be resolved (fail-fast behavior)
//   - Automatically handles application state, environment detection, and service management
//
// Returns:
//   - ApplicationInterface: The application's core service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "app" service is not registered in the container
//   - If the resolved service doesn't implement ApplicationInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call App() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Application state access is thread-safe and consistent
//
// Usage Examples:
//
//	// Environment detection
//	if facades.App().Environment() == "production" {
//	    fmt.Println("Running in production mode")
//	}
//
//	if facades.App().IsProduction() {
//	    enableProductionLogging()
//	}
//
//	if facades.App().IsLocal() {
//	    enableDebugMode()
//	}
//
//	if facades.App().IsTesting() {
//	    setupTestDatabase()
//	}
//
//	// Application information
//	appName := facades.App().Name()
//	version := facades.App().Version()
//	buildTime := facades.App().BuildTime()
//
//	fmt.Printf("Application: %s v%s (built %s)\n", appName, version, buildTime)
//
//	// Service resolution through app
//	logger := facades.App().Make("logger")
//	cache := facades.App().Make("cache")
//	db := facades.App().Make("database")
//
//	// Type-safe service resolution
//	loggerService := facades.App().MakeWithType[LoggerInterface]("logger")
//	cacheService := facades.App().MakeWithType[CacheInterface]("cache")
//
//	// Check if services are bound
//	if facades.App().Bound("redis") {
//	    redisClient := facades.App().Make("redis")
//	    // Use Redis client
//	}
//
//	if facades.App().Bound("elasticsearch") {
//	    esClient := facades.App().Make("elasticsearch")
//	    // Use Elasticsearch client
//	}
//
//	// Application paths
//	basePath := facades.App().BasePath()
//	configPath := facades.App().ConfigPath()
//	storagePath := facades.App().StoragePath()
//	publicPath := facades.App().PublicPath()
//
//	fmt.Printf("Base path: %s\n", basePath)
//	fmt.Printf("Config path: %s\n", configPath)
//	fmt.Printf("Storage path: %s\n", storagePath)
//
//	// Construct paths relative to application directories
//	logPath := facades.App().StoragePath("logs", "app.log")
//	configFile := facades.App().ConfigPath("database.yaml")
//	publicAsset := facades.App().PublicPath("assets", "style.css")
//
//	// Application lifecycle
//	facades.App().Booting(func() {
//	    fmt.Println("Application is booting...")
//	    // Initialize critical services
//	})
//
//	facades.App().Booted(func() {
//	    fmt.Println("Application has booted")
//	    // Perform post-boot initialization
//	})
//
//	facades.App().Terminating(func() {
//	    fmt.Println("Application is terminating...")
//	    // Cleanup resources
//	})
//
//	// Application state management
//	facades.App().SetMaintenance(true)
//	if facades.App().IsDownForMaintenance() {
//	    http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
//	    return
//	}
//
//	// Exit maintenance mode
//	facades.App().SetMaintenance(false)
//
//	// Service provider registration
//	facades.App().Register(&DatabaseServiceProvider{})
//	facades.App().Register(&CacheServiceProvider{})
//	facades.App().Register(&LoggingServiceProvider{})
//
//	// Bind services at runtime
//	facades.App().Bind("custom_service", func() interface{} {
//	    return &CustomService{
//	        config: facades.Config().Get("custom"),
//	    }
//	})
//
//	// Singleton services
//	facades.App().Singleton("file_manager", func() interface{} {
//	    return &FileManager{
//	        basePath: facades.App().StoragePath(),
//	    }
//	})
//
//	// Instance binding
//	customInstance := &CustomService{}
//	facades.App().Instance("custom_instance", customInstance)
//
// Advanced Application Patterns:
//
//	// Feature flags and environment-specific behavior
//	func setupFeatures() {
//	    if facades.App().Environment() == "development" {
//	        facades.App().Bind("profiler", func() interface{} {
//	            return &DevelopmentProfiler{}
//	        })
//	    }
//
//	    if facades.App().IsProduction() {
//	        facades.App().Bind("monitor", func() interface{} {
//	            return &ProductionMonitor{}
//	        })
//	    }
//	}
//
//	// Application bootstrapping
//	func BootstrapApplication() error {
//	    // Load configuration
//	    configPath := facades.App().ConfigPath("app.yaml")
//	    if err := loadConfiguration(configPath); err != nil {
//	        return fmt.Errorf("failed to load config: %w", err)
//	    }
//
//	    // Register service providers
//	    providers := []ServiceProvider{
//	        &DatabaseServiceProvider{},
//	        &CacheServiceProvider{},
//	        &LoggingServiceProvider{},
//	    }
//
//	    for _, provider := range providers {
//	        facades.App().Register(provider)
//	    }
//
//	    // Boot all providers
//	    return facades.App().Boot()
//	}
//
//	// Graceful shutdown handling
//	func HandleShutdown() {
//	    c := make(chan os.Signal, 1)
//	    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//
//	    go func() {
//	        <-c
//	        fmt.Println("Shutting down gracefully...")
//
//	        // Trigger termination callbacks
//	        facades.App().Terminate()
//
//	        os.Exit(0)
//	    }()
//	}
//
//	// Health checks
//	func HealthCheck() map[string]interface{} {
//	    status := map[string]interface{}{
//	        "app": map[string]interface{}{
//	            "name":        facades.App().Name(),
//	            "version":     facades.App().Version(),
//	            "environment": facades.App().Environment(),
//	            "uptime":      facades.App().Uptime(),
//	        },
//	        "services": map[string]interface{}{},
//	    }
//
//	    // Check critical services
//	    services := []string{"database", "cache", "logger"}
//	    for _, service := range services {
//	        if facades.App().Bound(service) {
//	            status["services"].(map[string]interface{})[service] = "available"
//	        } else {
//	            status["services"].(map[string]interface{})[service] = "unavailable"
//	        }
//	    }
//
//	    return status
//	}
//
//	// Environment-specific service binding
//	func bindEnvironmentServices() {
//	    switch facades.App().Environment() {
//	    case "local":
//	        facades.App().Bind("mail", func() interface{} {
//	            return &LogMailer{} // Log emails instead of sending
//	        })
//
//	    case "testing":
//	        facades.App().Bind("database", func() interface{} {
//	            return &InMemoryDatabase{} // Use in-memory database
//	        })
//
//	    case "production":
//	        facades.App().Bind("mail", func() interface{} {
//	            return &SMTPMailer{} // Use real SMTP
//	        })
//	    }
//	}
//
// Best Practices:
//   - Use the App facade for environment detection and path resolution
//   - Register service providers early in application lifecycle
//   - Use singleton binding for stateful services
//   - Implement proper shutdown handling with termination callbacks
//   - Use feature flags for environment-specific functionality
//   - Monitor application health through bound services
//   - Keep service bindings organized by feature or domain
//   - Use type-safe service resolution when possible
//
// Application Lifecycle:
//  1. Application creation and container setup
//  2. Service provider registration
//  3. Configuration loading
//  4. Service provider booting
//  5. Request handling
//  6. Graceful termination and cleanup
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume app service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	app, err := facade.Resolve[ApplicationInterface]("app")
//	if err != nil {
//	    // Handle app service unavailability gracefully
//	    return defaultEnvironment, fmt.Errorf("app unavailable: %w", err)
//	}
//	environment := app.Environment()
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestApplicationBehavior(t *testing.T) {
//	    // Create a test application instance
//	    testApp := &TestApplication{
//	        environment: "testing",
//	        name:        "TestApp",
//	        version:     "1.0.0-test",
//	        bindings:    make(map[string]interface{}),
//	    }
//
//	    // Swap the real app with test app
//	    restore := support.SwapService("app", testApp)
//	    defer restore() // Always restore after test
//
//	    // Now facades.App() returns testApp
//	    assert.Equal(t, "testing", facades.App().Environment())
//	    assert.True(t, facades.App().IsTesting())
//	    assert.Equal(t, "TestApp", facades.App().Name())
//
//	    // Test service binding
//	    facades.App().Bind("test_service", func() interface{} {
//	        return &TestService{}
//	    })
//
//	    assert.True(t, facades.App().Bound("test_service"))
//	    service := facades.App().Make("test_service")
//	    assert.NotNil(t, service)
//	}
//
// Container Configuration:
// Ensure the application service is properly configured in your container:
//
//	// Example app registration
//	container.Singleton("app", func() interface{} {
//	    config := application.Config{
//	        // Basic application information
//	        Name:        "MyApplication",
//	        Version:     "1.0.0",
//	        Environment: os.Getenv("APP_ENV"), // local, testing, staging, production
//	        Debug:       os.Getenv("APP_DEBUG") == "true",
//
//	        // Application paths
//	        BasePath:    "/app",
//	        ConfigPath:  "/app/config",
//	        StoragePath: "/app/storage",
//	        PublicPath:  "/app/public",
//
//	        // Application settings
//	        Timezone:    "UTC",
//	        Locale:      "en",
//
//	        // Service container
//	        Container: container,
//
//	        // Service providers to register
//	        Providers: []ServiceProvider{
//	            &DatabaseServiceProvider{},
//	            &CacheServiceProvider{},
//	            &LoggingServiceProvider{},
//	            &AuthServiceProvider{},
//	            &MailServiceProvider{},
//	        },
//
//	        // Lifecycle callbacks
//	        BootingCallbacks:     []func(){},
//	        BootedCallbacks:      []func(){},
//	        TerminatingCallbacks: []func(){},
//	    }
//
//	    app, err := application.NewApplication(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create application: %v", err)
//	    }
//
//	    return app
//	})
func App() applicationInterfaces.ApplicationInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "app" service from the dependency injection container
	// - Performs type assertion to ApplicationInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[applicationInterfaces.ApplicationInterface](applicationTokens.APPLICATION_TOKEN)
}

// AppWithError provides error-safe access to the application service.
//
// This function offers the same functionality as App() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle application service unavailability gracefully.
//
// This is a convenience wrapper around facade.Resolve() that provides
// the same caching and performance benefits as App() but with error handling.
//
// Returns:
//   - ApplicationInterface: The resolved app instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ApplicationInterface
//
// Usage Examples:
//
//	// Basic error-safe application access
//	app, err := facades.AppWithError()
//	if err != nil {
//	    log.Printf("App service unavailable: %v", err)
//	    return defaultConfig, fmt.Errorf("application not available")
//	}
//	environment := app.Environment()
//
//	// Conditional application operations
//	if app, err := facades.AppWithError(); err == nil {
//	    if app.IsProduction() {
//	        // Perform production-specific operations
//	        enableProductionMode()
//	    }
//	}
//
//	// Health check pattern
//	func CheckAppHealth() error {
//	    app, err := facades.AppWithError()
//	    if err != nil {
//	        return fmt.Errorf("application service unavailable: %w", err)
//	    }
//
//	    // Test basic app functionality
//	    if app.Name() == "" {
//	        return fmt.Errorf("application not properly configured")
//	    }
//
//	    // Check critical services
//	    if !app.Bound("logger") {
//	        return fmt.Errorf("logger service not bound")
//	    }
//
//	    return nil
//	}
func AppWithError() (applicationInterfaces.ApplicationInterface, error) {
	// Use facade.Resolve() for error-return behavior:
	// - Resolves "app" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[applicationInterfaces.ApplicationInterface](applicationTokens.APPLICATION_TOKEN)
}
