package traits

import (
	"sync"

	"govel/packages/application/helpers"
	"govel/packages/config"
	interfaces "govel/packages/types/src/interfaces/config"
)

/**
 * Configurable provides configuration management functionality in a thread-safe manner.
 * This trait follows the embedding pattern where it embeds the config instance directly,
 * providing both trait-specific management methods and transparent access to all
 * config interface methods through embedding.
 */
type Configurable struct {
	/**
	 * mutex provides thread-safe access to config properties
	 */
	mutex sync.RWMutex

	/**
	 * Config instance embedded directly to provide transparent delegation
	 * All ConfigInterface methods are automatically available
	 */
	*config.Config
}

/**
 * NewConfigurable creates a new Configurable instance with a config.
 *
 * @param configInstance *config.Config The config instance to use
 * @return *Configurable The newly created trait instance
 */
func NewConfigurable(configInstance *config.Config) *Configurable {
	if configInstance == nil {
		configInstance = config.New()
	}

	return &Configurable{
		Config: configInstance,
	}
}

/**
 * NewConfigurableWithEnvironment creates a new Configurable instance with environment-specific config.
 * If no environment is provided, it will be read from environment variables.
 *
 * @param environment string Optional environment to configure for (variadic, first value used if provided)
 * @return *Configurable The newly created trait instance
 *
 * Example:
 *   // Using environment from env vars
 *   config := NewConfigurableWithEnvironment()
 *   // Providing explicit environment
 *   config := NewConfigurableWithEnvironment("production")
 */
func NewConfigurableWithEnvironment(environment ...string) *Configurable {
	// Get helper instance
	envHelper := helpers.NewEnvHelper()

	// Use provided environment or fallback to environment helper
	appEnv := envHelper.GetAppEnvironment() // Default from environment
	if len(environment) > 0 && environment[0] != "" {
		appEnv = environment[0]
	}

	return &Configurable{
		Config: config.NewWithEnvironment(appEnv),
	}
}

/**
 * GetConfig returns the config instance.
 *
 * @return interfaces.ConfigInterface The config instance
 */
func (t *Configurable) GetConfig() interfaces.ConfigInterface {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Config
}

/**
 * SetConfig sets the config instance.
 *
 * @param configInstance interface{} The config instance to set (using interface{} to avoid circular import)
 */
func (t *Configurable) SetConfig(configInstance interfaces.ConfigInterface) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if configPtr, ok := configInstance.(*config.Config); ok {
		if configPtr == nil {
			configPtr = config.New()
		}
		t.Config = configPtr
	} else {
		// Fallback to default config if invalid type
		t.Config = config.New()
	}
}

/**
 * HasConfig returns whether a config instance is set.
 *
 * @return bool true if a config is set
 */
func (t *Configurable) HasConfig() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Config != nil
}

/**
 * GetConfigInfo returns information about the config.
 *
 * @return map[string]interface{} Config information
 */
func (t *Configurable) GetConfigInfo() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	info := map[string]interface{}{
		"has_config": t.Config != nil,
	}

	if t.Config != nil {
		// Add config-specific information if available
		// This would depend on the actual config implementation
		info["config_type"] = "default"

		// You could add more specific info like:
		info["config_keys_count"] = len(t.Config.AllConfig())
		info["environment"] = t.Config.GetEnvironment()
	}

	return info
}

// Compile-time interface compliance check
var _ interfaces.ConfigurableInterface = (*Configurable)(nil)
