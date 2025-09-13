package interfaces

// ConfigurableInterface defines the contract for components that can be configured.
// This interface allows components to receive and work with configuration.
type ConfigurableInterface interface {
	ConfigInterface

	// GetConfig returns the current configuration.
	GetConfig() ConfigInterface

	// SetConfig sets the configuration for the component.
	SetConfig(configInstance ConfigInterface)

	// HasConfig checks if configuration is available.
	HasConfig() bool

	// GetConfigInfo returns information about the config.
	GetConfigInfo() map[string]interface{}
}
