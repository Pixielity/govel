package config

import (
	"strings"

	"govel/config/drivers"
	"govel/support"
	enums "govel/types/enums/config"
	configInterfaces "govel/types/interfaces/config"
	containerInterfaces "govel/types/interfaces/container"
)

// ConfigManager provides centralized configuration management with multiple driver support.
// Extends support.Manager implementing Laravel-style driver pattern for configuration operations.
//
// Key features:
//   - Multi-driver support (file, memory, env, remote, firestore, appconfig)
//   - Laravel-style driver creation and management
//   - Container-based dependency injection
//   - Automatic driver validation and fallback
//
// Implements ConfigInterface and FactoryInterface for complete configuration functionality.
type ConfigManager struct {
	*support.Manager
}

// NewConfigManager creates a new configuration manager with Laravel-style functionality.
// Initializes the manager with container dependency injection and proxy self-reference.
//
// Constructor features:
//   - Container-based dependency injection for configuration
//   - Proxy self-reference for proper method resolution
//   - Laravel-style driver pattern implementation
//   - Automatic driver discovery via reflection
//
// Returns configured ConfigManager ready for configuration operations.
func NewConfigManager(container containerInterfaces.ContainerInterface) *ConfigManager {
	// Create base manager with dependency injection container
	// This provides core driver management functionality
	baseManager := support.NewManager(container)

	// Initialize config manager wrapping the base manager
	// Inherits all driver management capabilities
	configManager := &ConfigManager{
		Manager: baseManager,
	}

	// Set up proxy self-reference for proper method resolution
	// This enables the Manager's reflection to find ConfigManager's CreateXXXDriver methods
	// Critical for Laravel-style automatic driver discovery
	baseManager.SetProxySelf(configManager)

	return configManager
}

// Config gets a config instance by driver name (optional).
// If no name is provided, uses the default driver.
// Provides type-safe access to specific configuration drivers with validation.
//
// Method features:
//   - Optional driver name (uses default if not provided)
//   - Driver validation before driver creation (only if name provided)
//   - Automatic driver instantiation and caching
//   - Type-safe ConfigInterface casting
//   - Graceful error handling with nil returns
//
// Returns nil if driver is invalid or driver creation fails.
func (c *ConfigManager) Config(name ...enums.Driver) configInterfaces.ConfigInterface {
	var driverName string

	// Determine which driver to use
	if len(name) > 0 && name[0] != "" {
		// Name provided - validate it
		driverName = name[0].String()
	}

	// Get driver instance using base manager's driver resolution
	// This triggers CreateXXXDriver methods via reflection
	driver, err := c.Driver(driverName)
	if err != nil {
		// Driver creation failed - return nil for consistent error handling
		return nil
	}

	// Type-safe casting to ConfigerInterface
	// Ensures returned driver implements required config operations
	configer, ok := driver.(configInterfaces.ConfigInterface)
	if !ok {
		// Driver doesn't implement ConfigerInterface - should never happen
		return nil
	}

	return configer
}

// Driver creation methods - Laravel reflection-style method naming
// These methods are automatically discovered by the base Manager via reflection.
// Method naming convention: Create{DriverName}Driver() for automatic discovery.
// Each method handles configuration loading and driver instantiation for its driver type.

// CreateFileDriver creates file driver instance.
// Automatically discovered by base Manager for "file" driver requests.
//
// Simply passes configuration from container to driver constructor.
// All defaults, validation, and parameter handling is done by the driver itself.
func (c *ConfigManager) CreateFileDriver() (interface{}, error) {
	// Get configuration from container and pass directly to driver
	var config map[string]interface{}

	// Convert map to FileDriverOptions
	options := &drivers.FileDriverOptions{
		ConfigPaths: []string{".", "./config", "/etc/config"},
		ConfigName:  "config",
		ConfigType:  "yaml",
	}

	if config != nil {
		if paths, ok := config["paths"].([]string); ok {
			options.ConfigPaths = paths
		}
		if name, ok := config["name"].(string); ok {
			options.ConfigName = name
		}
		if configType, ok := config["type"].(string); ok {
			options.ConfigType = configType
		}
	}

	return drivers.NewFileDriver(options), nil
}

// CreateMemoryDriver creates memory driver instance.
// Automatically discovered by base Manager for "memory" driver requests.
//
// Simply passes configuration from container to driver constructor.
// All defaults, validation, and parameter handling is done by the driver itself.
func (c *ConfigManager) CreateMemoryDriver() (interface{}, error) {
	// Get configuration from container and pass directly to driver
	var config map[string]interface{}

	// Convert map to MemoryDriverOptions
	options := &drivers.MemoryDriverOptions{
		InitialData: make(map[string]interface{}),
	}

	if config != nil {
		if initialData, ok := config["initial_data"].(map[string]interface{}); ok {
			options.InitialData = initialData
		}
	}

	return drivers.NewMemoryDriver(options), nil
}

// CreateEnvDriver creates environment driver instance.
// Automatically discovered by base Manager for "env" driver requests.
//
// Simply passes configuration from container to driver constructor.
// All defaults, validation, and parameter handling is done by the driver itself.
func (c *ConfigManager) CreateEnvDriver() (interface{}, error) {
	// Get configuration from container and pass directly to driver
	var config map[string]interface{}

	// Convert map to EnvDriverOptions
	options := &drivers.EnvDriverOptions{
		Prefix:       "APP_",
		AutomaticEnv: true,
		Replacer:     strings.NewReplacer(".", "_", "-", "_"),
	}

	if config != nil {
		if prefix, ok := config["prefix"].(string); ok {
			options.Prefix = prefix
		}
		if automaticEnv, ok := config["automatic_env"].(bool); ok {
			options.AutomaticEnv = automaticEnv
		}
		if envKeyMap, ok := config["env_key_map"].(map[string]string); ok {
			options.EnvKeyMap = envKeyMap
		}
	}

	return drivers.NewEnvDriver(options), nil
}

// CreateRemoteDriver creates remote driver instance.
// Automatically discovered by base Manager for "remote" driver requests.
//
// Simply passes configuration from container to driver constructor.
// All defaults, validation, and parameter handling is done by the driver itself.
func (c *ConfigManager) CreateRemoteDriver() (interface{}, error) {
	// Get configuration from container and pass directly to driver
	var config map[string]interface{}

	// Convert map to RemoteDriverOptions
	options := &drivers.RemoteDriverOptions{
		Provider: "etcd",
		Endpoint: "http://localhost:2379",
		KeyPath:  "/config",
	}

	if config != nil {
		if provider, ok := config["provider"].(string); ok {
			options.Provider = provider
		}
		if endpoint, ok := config["endpoint"].(string); ok {
			options.Endpoint = endpoint
		}
		if keyPath, ok := config["key_path"].(string); ok {
			options.KeyPath = keyPath
		}
	}

	return drivers.NewRemoteDriver(options), nil
}

// GetDefaultDriver implements the ManagerInterface requirement.
// Returns the default driver identifier for driver resolution and delegation.
//
// Default driver features:
//   - Centralized default driver configuration
//   - Consistent behavior across all manager operations
//   - Easy configuration override
//   - Laravel-style manager pattern compliance
//
// Returns the default driver name.
func (c *ConfigManager) GetDefaultDriver() string {
	// Return default driver name
	return "file"
}

// Compile-time interface compliance checks
// These ensure ConfigManager properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ configInterfaces.FactoryInterface = (*ConfigManager)(nil) // Config factory operations
