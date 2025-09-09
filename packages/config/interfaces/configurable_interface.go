package interfaces

/**
 * ConfigurableInterface defines the contract for components that provide
 * configuration functionality. This interface follows the Interface Segregation
 * Principle by focusing solely on configuration operations.
 *
 * By embedding ConfigInterface, all config methods are automatically
 * available through this interface, providing transparent delegation.
 */
type ConfigurableInterface interface {
	/**
	 * Embed ConfigInterface to provide transparent access to all config methods
	 */
	ConfigInterface

	/**
	 * GetConfig returns the config instance.
	 *
	 * @return ConfigInterface The config instance
	 */
	GetConfig() ConfigInterface

	/**
	 * SetConfig sets the config instance.
	 *
	 * @param config interface{} The config instance to set (using interface{} to avoid circular import)
	 */
	SetConfig(config interface{})

	/**
	 * HasConfig returns whether a config instance is set.
	 *
	 * @return bool true if a config is set
	 */
	HasConfig() bool

	/**
	 * GetConfigInfo returns information about the config.
	 *
	 * @return map[string]interface{} Config information
	 */
	GetConfigInfo() map[string]interface{}
}
