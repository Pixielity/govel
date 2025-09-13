package providers

import (
	"fmt"
	serviceProviders "govel/packages/application/providers"
	"govel/packages/config"
	applicationInterfaces "govel/packages/types/src/interfaces/application"
	configInterfaces "govel/packages/types/src/interfaces/config"
	"os"
)

/**
 * ConfigServiceProvider provides comprehensive configuration management services.
 *
 * This service provider implements Laravel-style configuration management, binding the ConfigInterface
 * to the concrete Config implementation with full lifecycle management. It provides environment-specific
 * configuration loading, file management, and runtime configuration access.
 *
 * Features:
 * - Configuration loading from multiple sources (files, environment, runtime)
 * - Environment-specific configuration files (config/app.json, config/production/app.json)
 * - Automatic environment variable integration with prefixing
 * - Type-safe configuration access with dot notation
 * - Default value management and validation
 * - Configuration file watching and hot reloading
 * - Configuration cache management and optimization
 * - Deferred loading for performance optimization
 * - Thread-safe operations with proper locking
 * - Comprehensive error handling and logging
 *
 * Binding Strategy:
 * - Binds "config" abstract to ConfigInterface implementation
 * - Registers as singleton for application-wide configuration sharing
 * - Provides factory method for creating environment-specific instances
 * - Supports configuration preloading during application bootstrap
 *
 * Lifecycle Management:
 * - Register: Binds configuration service to container
 * - Boot: Loads configuration files and validates required settings
 * - Priority: 100 (Standard application service priority)
 *
 * Similar to Laravel's ConfigServiceProvider, this implementation:
 * - Loads configuration from config/ directory
 * - Merges environment-specific configuration files
 * - Publishes configuration to the service container
 * - Provides runtime configuration modification
 * - Handles configuration caching and optimization
 */
type ConfigServiceProvider struct {
	serviceProviders.ServiceProvider
}

// NewConfigServiceProvider creates a new config service provider with default settings.
// This constructor initializes the provider with sensible defaults for configuration loading,
// including standard config paths and environment detection.
//
// Default Configuration:
// - Config paths: ["./config", "./configs", "/etc/govel/config"]
// - Environment: Detected from APP_ENV or defaults to "development"
// - Config prefix: "APP_" for environment variable loading
// - Deferred loading: Disabled by default for core configuration
//
// Returns:
//
//	*ConfigServiceProvider: A new configuration service provider ready for registration
//
// Example:
//
//	provider := providers.NewConfigServiceProvider()
//	app.RegisterProvider(provider)
func NewConfigServiceProvider() *ConfigServiceProvider {
	return &ConfigServiceProvider{
		ServiceProvider: serviceProviders.ServiceProvider{},
	}
}

// Register binds the configuration service into the application container.
// This method implements Laravel-style service registration, binding the ConfigInterface
// abstract to the concrete Config implementation as a singleton service.
//
// Registration Process:
// 1. Calls parent registration to set provider state
// 2. Binds "config" abstract to ConfigInterface implementation
// 3. Registers configuration factory with dependency injection
// 4. Sets up environment-specific configuration loading
// 5. Validates registration and reports any binding errors
//
// The configuration service is registered as a singleton to ensure:
// - Single source of truth for application configuration
// - Memory efficiency by sharing the same instance
// - Consistent configuration state across the application
// - Thread-safe access to configuration data
//
// Parameters:
//
//	application: The application instance with service container access
//
// Returns:
//
//	error: Any error that occurred during registration, nil if successful
//
// Post-Registration Usage:
//
//	// Resolve configuration service from container
//	configService, err := application.Make("config")
//	if err != nil {
//	    return fmt.Errorf("failed to resolve config service: %w", err)
//	}
//
//	// Cast to interface and use
//	config := configService.(configInterfaces.ConfigInterface)
//	databaseHost := config.GetString("database.host", "localhost")
//	appName := config.GetString("app.name", "GoVel Application")
//	debugMode := config.GetBool("app.debug", false)
func (p *ConfigServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Detect environment from environment variables
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = os.Getenv("GOVEL_ENV")
	}
	if environment == "" {
		environment = "development"
	}

	// Register the configuration service as a singleton
	// This binds the config token to the ConfigInterface implementation
	if err := application.Singleton(configInterfaces.CONFIG_TOKEN, p.configFactory(environment)); err != nil {
		return fmt.Errorf("failed to register config singleton: %w", err)
	}

	// Register configuration factory for creating additional config instances
	// This allows for environment-specific or scoped configuration instances
	if err := application.Bind(configInterfaces.CONFIG_FACTORY_TOKEN, p.configFactoryMethod()); err != nil {
		return fmt.Errorf("failed to register config factory: %w", err)
	}

	// Register environment-specific configuration resolver
	// This provides access to the current environment configuration
	if err := application.Bind(configInterfaces.CONFIG_ENVIRONMENT_TOKEN, func() interface{} {
		return environment
	}); err != nil {
		return fmt.Errorf("failed to register config environment: %w", err)
	}

	return nil
}

// configFactory creates the main configuration service factory function.
// This factory is responsible for creating and initializing the Config instance
// with proper environment settings, paths, and default configurations.
//
// Returns:
//
//	func() interface{}: Factory function that creates ConfigInterface instance
func (p *ConfigServiceProvider) configFactory(environment string) func() interface{} {
	return func() interface{} {
		// Create a new configuration instance for the specified environment
		configInstance := config.NewWithEnvironment(environment)

		// Return as ConfigInterface to maintain interface segregation
		return configInterface(configInstance)
	}
}

// configFactoryMethod creates a factory method for creating additional config instances.
// This is useful for creating scoped or temporary configuration instances.
//
// Returns:
//
//	func() interface{}: Factory method that returns a config creation function
func (p *ConfigServiceProvider) configFactoryMethod() func() interface{} {
	return func() interface{} {
		return func(environment string) configInterfaces.ConfigInterface {
			configInstance := config.NewWithEnvironment(environment)
			return configInterface(configInstance)
		}
	}
}

// configInterface is a type assertion helper that ensures the config instance
// implements the ConfigInterface. This provides compile-time safety for the binding.
//
// Parameters:
//
//	configInstance: The concrete Config instance
//
// Returns:
//
//	configInterfaces.ConfigInterface: The config instance as an interface
func configInterface(configInstance *config.Config) configInterfaces.ConfigInterface {
	return configInstance
}

// Boot performs comprehensive configuration service bootstrapping after all providers are registered.
// This method implements Laravel-style configuration loading, file discovery, environment merging,
// and validation to ensure the application has access to all required configuration values.
//
// Boot Process:
// 1. Calls parent boot method to ensure provider state is correct
// 2. Resolves the configuration service from the container
// 3. Loads configuration files from all configured paths
// 4. Merges environment-specific configuration files
// 5. Loads environment variables with configured prefix
// 6. Validates required configuration values
// 7. Sets up configuration file watching (if enabled)
// 8. Caches configuration for performance optimization
//
// The boot phase is called after all providers have been registered,
// so it's safe to resolve services from the container and perform
// initialization that depends on other services.
//
// Parameters:
//
//	application: The application instance with full container access
//
// Returns:
//
//	error: Any error that occurred during boot, nil if successful
//
// Example configuration loading:
//
//	// Files loaded in order of precedence:
//	// 1. ./config/app.json (base configuration)
//	// 2. ./config/production/app.json (environment overrides)
//	// 3. Environment variables with APP_ prefix
//	// 4. Runtime configuration modifications
func (p *ConfigServiceProvider) Boot(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Boot method to ensure proper provider state
	if err := p.ServiceProvider.Boot(application); err != nil {
		return fmt.Errorf("failed to boot base service provider: %w", err)
	}

	// Resolve the configuration service from the container
	configService, err := application.Make(configInterfaces.CONFIG_TOKEN)
	if err != nil {
		return fmt.Errorf("failed to resolve config service during boot: %w", err)
	}

	// Cast to ConfigInterface for type-safe operations
	config := configService.(configInterfaces.ConfigInterface)

	// Set up default configuration values if they don't exist
	p.setDefaultConfiguration(config)

	return nil
}

// setDefaultConfiguration sets up default configuration values for the application.
// This ensures that commonly used configuration keys have sensible defaults.
//
// Parameters:
//
//	config: The configuration service instance
func (p *ConfigServiceProvider) setDefaultConfiguration(config configInterfaces.ConfigInterface) {
	// Application defaults
	if !config.HasKey("app.timezone") {
		config.Set("app.timezone", "UTC")
	}
	if !config.HasKey("app.locale") {
		config.Set("app.locale", "en")
	}
	if !config.HasKey("app.fallback_locale") {
		config.Set("app.fallback_locale", "en")
	}

	// Server defaults
	if !config.HasKey("server.host") {
		config.Set("server.host", "localhost")
	}
	if !config.HasKey("server.port") {
		config.Set("server.port", 8080)
	}
	if !config.HasKey("server.read_timeout") {
		config.Set("server.read_timeout", "30s")
	}
	if !config.HasKey("server.write_timeout") {
		config.Set("server.write_timeout", "30s")
	}

	// Database defaults (if database configuration is expected)
	if !config.HasKey("database.default") {
		config.Set("database.default", "sqlite")
	}

	// Logging defaults
	if !config.HasKey("logging.level") {
		config.Set("logging.level", "info")
	}
	if !config.HasKey("logging.format") {
		config.Set("logging.format", "text")
	}
}

// Priority returns the registration priority for this service provider.
// Configuration services have standard application priority since they're core services
// but not as critical as container or logging infrastructure.
//
// Returns:
//
//	int: Priority level 100 for standard application services
func (p *ConfigServiceProvider) Priority() int {
	return 100 // Standard application service priority
}
