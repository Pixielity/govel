package facades

import (
	configInterfaces "govel/types/interfaces/config"
	facade "govel/support"
)

// Config provides a clean, static-like interface to the application's configuration service.
//
// This facade implements the facade pattern, providing global access to the configuration
// service configured in the dependency injection container. It offers a Laravel-style
// API for configuration management with automatic service resolution, type safety, and
// support for multiple configuration sources and formats.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved config service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent access across goroutines
//   - Supports multiple config sources (files, environment, remote, etc.)
//   - Built-in support for configuration hot-reloading and change detection
//
// Behavior:
//   - First call: Resolves config service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if config service cannot be resolved (fail-fast behavior)
//   - Automatically handles configuration loading, parsing, and validation
//
// Returns:
//   - ConfigInterface: The application's configuration service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "config" service is not registered in the container
//   - If the resolved service doesn't implement ConfigInterface
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
//   - Multiple goroutines can call Config() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Configuration access is thread-safe and consistent
//
// Usage Examples:
//
//	// Basic configuration access
//	appName := facades.Config().GetString("app.name")
//	debugMode := facades.Config().GetBool("app.debug")
//	port := facades.Config().GetInt("server.port")
//	timeout := facades.Config().GetDuration("server.timeout")
//
//	// Configuration with default values
//	theme := facades.Config().GetStringWithDefault("ui.theme", "dark")
//	maxRetries := facades.Config().GetIntWithDefault("api.max_retries", 3)
//	enabled := facades.Config().GetBoolWithDefault("feature.new_ui", false)
//
//	// Nested configuration access
//	dbHost := facades.Config().GetString("database.connections.mysql.host")
//	dbPort := facades.Config().GetInt("database.connections.mysql.port")
//	dbName := facades.Config().GetString("database.connections.mysql.database")
//
//	// Array/slice configuration
//	allowedHosts := facades.Config().GetStringSlice("security.allowed_hosts")
//	corsOrigins := facades.Config().GetStringSlice("cors.allowed_origins")
//	adminEmails := facades.Config().GetStringSlice("notifications.admin_emails")
//
//	// Map configuration
//	logLevels := facades.Config().GetStringMap("logging.levels")
//	featureFlags := facades.Config().GetBoolMap("features.flags")
//	environmentVars := facades.Config().GetStringMap("environment.variables")
//
//	// Complex configuration structures
//	type DatabaseConfig struct {
//	    Host     string `mapstructure:"host"`
//	    Port     int    `mapstructure:"port"`
//	    Database string `mapstructure:"database"`
//	    Username string `mapstructure:"username"`
//	    Password string `mapstructure:"password"`
//	}
//
//	var dbConfig DatabaseConfig
//	facades.Config().UnmarshalKey("database.connections.mysql", &dbConfig)
//
//	// Environment variable integration
//	// Automatically reads from ENV variables with fallbacks
//	appEnv := facades.Config().GetString("app.environment") // APP_ENVIRONMENT
//	secretKey := facades.Config().GetString("app.secret_key") // APP_SECRET_KEY
//	maxWorkers := facades.Config().GetInt("workers.max_count") // WORKERS_MAX_COUNT
//
//	// Configuration validation and type safety
//	if !facades.Config().IsSet("database.host") {
//	    log.Fatal("Database host configuration is required")
//	}
//
//	if facades.Config().GetString("app.environment") == "" {
//	    log.Fatal("Application environment must be specified")
//	}
//
//	// Dynamic configuration updates (if supported)
//	facades.Config().Set("runtime.last_updated", time.Now())
//	facades.Config().Set("feature.experimental", true)
//
//	// Configuration watching and hot-reloading
//	facades.Config().WatchConfig()
//	facades.Config().OnConfigChange(func(in fsnotify.Event) {
//	    log.Printf("Configuration file changed: %s", in.Name)
//	    // Handle configuration reload
//	})
//
//	// Scoped configuration access
//	loggerConfig := facades.Config().Sub("logging")
//	level := loggerConfig.GetString("level") // equivalent to facades.Config().GetString("logging.level")
//	format := loggerConfig.GetString("format")
//
// Configuration File Examples:
//
//	// config/app.yaml
//	app:
//	  name: "MyApplication"
//	  version: "1.0.0"
//	  environment: "${APP_ENVIRONMENT:development}"
//	  debug: ${APP_DEBUG:false}
//	  secret_key: "${APP_SECRET_KEY}"
//
//	server:
//	  host: "${SERVER_HOST:localhost}"
//	  port: ${SERVER_PORT:8080}
//	  timeout: "${SERVER_TIMEOUT:30s}"
//
//	database:
//	  connections:
//	    mysql:
//	      host: "${DB_HOST:localhost}"
//	      port: ${DB_PORT:3306}
//	      database: "${DB_DATABASE:myapp}"
//	      username: "${DB_USERNAME:root}"
//	      password: "${DB_PASSWORD}"
//
//	// config/features.json
//	{
//	  "features": {
//	    "flags": {
//	      "new_ui": false,
//	      "beta_features": true,
//	      "experimental_api": false
//	    }
//	  },
//	  "security": {
//	    "allowed_hosts": ["localhost", "127.0.0.1", "*.example.com"],
//	    "cors": {
//	      "allowed_origins": ["http://localhost:3000", "https://app.example.com"]
//	    }
//	  }
//	}
//
// Best Practices:
//   - Use hierarchical configuration keys ("app.name", "database.host")
//   - Provide sensible defaults for non-critical configuration
//   - Use environment variables for environment-specific values
//   - Validate critical configuration at application startup
//   - Use type-safe getters (GetString, GetInt, GetBool, etc.)
//   - Group related configuration under common prefixes
//   - Use configuration structs for complex nested config
//
// Environment Variable Patterns:
//   - Use UPPER_CASE with underscores for env var names
//   - Mirror configuration hierarchy: "app.secret_key" → "APP_SECRET_KEY"
//   - Provide defaults in configuration files when possible
//   - Use ${VAR:default} syntax for environment substitution
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume configuration always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	config, err := facade.TryResolve[ConfigInterface]("config")
//	if err != nil {
//	    // Handle config unavailability gracefully
//	    useDefaultConfiguration()
//	    return
//	}
//	appName := config.GetString("app.name")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserService(t *testing.T) {
//	    // Create a test configuration
//	    testConfig := &TestConfig{
//	        values: map[string]interface{}{
//	            "app.name":    "TestApp",
//	            "app.debug":   true,
//	            "server.port": 8888,
//	        },
//	    }
//
//	    // Swap the real config with test config
//	    restore := support.SwapService("config", testConfig)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Config() returns testConfig
//	    userService := NewUserService()
//
//	    // Test with controlled configuration
//	    assert.Equal(t, "TestApp", facades.Config().GetString("app.name"))
//	    assert.True(t, facades.Config().GetBool("app.debug"))
//	    assert.Equal(t, 8888, facades.Config().GetInt("server.port"))
//	}
//
// Container Configuration:
// Ensure the config service is properly configured in your container:
//
//	// Example config registration
//	container.Singleton("config", func() interface{} {
//	    config := viper.New()
//
//	    // Set configuration file locations
//	    config.SetConfigName("app")
//	    config.SetConfigType("yaml")
//	    config.AddConfigPath("/etc/myapp/")
//	    config.AddConfigPath("$HOME/.myapp")
//	    config.AddConfigPath(".")
//
//	    // Enable environment variable support
//	    config.AutomaticEnv()
//	    config.SetEnvPrefix("MYAPP") // MYAPP_APP_NAME → app.name
//	    config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
//
//	    // Set defaults
//	    config.SetDefault("app.environment", "development")
//	    config.SetDefault("server.port", 8080)
//	    config.SetDefault("app.debug", false)
//
//	    // Read configuration
//	    if err := config.ReadInConfig(); err != nil {
//	        log.Printf("No config file found: %v", err)
//	    }
//
//	    return config
//	})
func Config() configInterfaces.ConfigInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves config service using type-safe token from the dependency injection container
	// - Performs type assertion to ConfigInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[configInterfaces.ConfigInterface](configInterfaces.CONFIG_TOKEN)
}

// ConfigWithError provides error-safe access to the configuration service.
//
// This function offers the same functionality as Config() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle configuration unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Config() but with error handling.
//
// Returns:
//   - ConfigInterface: The resolved config instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ConfigInterface
//
// Usage Examples:
//
//	// Basic error-safe configuration access
//	config, err := facades.ConfigWithError()
//	if err != nil {
//	    log.Printf("Configuration unavailable: %v", err)
//	    useDefaultSettings() // Fallback to hardcoded defaults
//	    return
//	}
//	appName := config.GetString("app.name")
//
//	// Conditional configuration loading
//	if config, err := facades.ConfigWithError(); err == nil {
//	    if config.GetBool("features.advanced_logging") {
//	        setupAdvancedLogging()
//	    }
//	}
func ConfigWithError() (configInterfaces.ConfigInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves config service using type-safe token from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[configInterfaces.ConfigInterface](configInterfaces.CONFIG_TOKEN)
}
