package drivers

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"

	configInterfaces "govel/types/interfaces/config"
)

// EnvDriver implements config.Driver interface for environment variable configuration.
// It uses Viper internally for environment variable operations.
type EnvDriver struct {
	viper     *viper.Viper
	prefix    string
	replacer  *strings.Replacer
	mu        sync.RWMutex
	watchFunc func()
	envKeyMap map[string]string // Maps config keys to env var names
}

// EnvDriverOptions contains configuration options for EnvDriver.
type EnvDriverOptions struct {
	Prefix       string            // Environment variable prefix (e.g., "APP_")
	Replacer     *strings.Replacer // Custom replacer for key transformation
	EnvKeyMap    map[string]string // Maps config keys to specific env var names
	AutomaticEnv bool              // Enable automatic environment variable reading
}

// NewEnvDriver creates a new environment variable configuration driver.
//
// Example:
//
//	driver := drivers.NewEnvDriver(&drivers.EnvDriverOptions{
//	    Prefix:       "MYAPP_",
//	    AutomaticEnv: true,
//	    Replacer:     strings.NewReplacer(".", "_", "-", "_"),
//	})
func NewEnvDriver(options *EnvDriverOptions) *EnvDriver {
	v := viper.New()

	// Set defaults
	if options == nil {
		options = &EnvDriverOptions{
			Prefix:       "",
			AutomaticEnv: true,
		}
	}

	// Configure viper for environment variables
	if options.AutomaticEnv {
		v.AutomaticEnv()
	}

	if options.Prefix != "" {
		v.SetEnvPrefix(options.Prefix)
	}

	// Set custom key replacer (e.g., "." -> "_")
	if options.Replacer != nil {
		v.SetEnvKeyReplacer(options.Replacer)
	} else {
		// Default replacer: dots become underscores
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	}

	driver := &EnvDriver{
		viper:     v,
		prefix:    options.Prefix,
		replacer:  options.Replacer,
		envKeyMap: options.EnvKeyMap,
	}

	// Bind specific environment variables if provided
	if options.EnvKeyMap != nil {
		for configKey, envKey := range options.EnvKeyMap {
			v.BindEnv(configKey, envKey)
		}
	}

	return driver
}

// Get retrieves a configuration value by key.
func (e *EnvDriver) Get(key string) (interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Check if the key is bound to a specific environment variable
	if envKey, exists := e.envKeyMap[key]; exists {
		if value := os.Getenv(envKey); value != "" {
			return value, nil
		}
	}

	// Use viper's automatic environment variable resolution
	if value := e.viper.Get(key); value != nil {
		return value, nil
	}

	return nil, fmt.Errorf("key '%s' not found", key)
}

// Set stores a configuration value by setting an environment variable.
func (e *EnvDriver) Set(key string, value interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Convert value to string
	strValue := fmt.Sprintf("%v", value)

	// Check if there's a specific env var mapping
	if envKey, exists := e.envKeyMap[key]; exists {
		return os.Setenv(envKey, strValue)
	}

	// Generate environment variable name
	envKey := e.generateEnvKey(key)
	return os.Setenv(envKey, strValue)
}

// Load loads configuration from environment variables.
func (e *EnvDriver) Load() error {
	// Environment variables are automatically loaded by viper
	// when AutomaticEnv is enabled, so this is a no-op
	return nil
}

// Watch starts watching for environment variable changes.
func (e *EnvDriver) Watch(callback func()) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.watchFunc != nil {
		return fmt.Errorf("already watching for config changes")
	}

	e.watchFunc = callback

	// Note: Environment variables don't have built-in file system watching
	// This would require a more sophisticated implementation with polling
	// or external tools to detect environment changes
	return nil
}

// Unwatch stops watching for environment variable changes.
func (e *EnvDriver) Unwatch() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.watchFunc = nil
	return nil
}

// Invalidate clears the configuration cache.
func (e *EnvDriver) Invalidate() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Create a new viper instance to clear cache
	v := viper.New()

	if e.viper != nil {
		v.AutomaticEnv()
		if e.prefix != "" {
			v.SetEnvPrefix(e.prefix)
		}
		if e.replacer != nil {
			v.SetEnvKeyReplacer(e.replacer)
		} else {
			v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		}
	}

	e.viper = v

	// Re-bind specific environment variables if provided
	if e.envKeyMap != nil {
		for configKey, envKey := range e.envKeyMap {
			e.viper.BindEnv(configKey, envKey)
		}
	}

	return nil
}

// GetAll returns all configuration data from environment variables.
func (e *EnvDriver) GetAll() (map[string]interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make(map[string]interface{})

	// Get all environment variables
	env := os.Environ()
	for _, envVar := range env {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envKey := parts[0]
		envValue := parts[1]

		// Check if this env var matches our prefix
		if e.prefix != "" && !strings.HasPrefix(envKey, e.prefix) {
			continue
		}

		// Convert environment variable name to config key
		configKey := e.envKeyToConfigKey(envKey)
		result[configKey] = envValue
	}

	// Add specifically mapped environment variables
	if e.envKeyMap != nil {
		for configKey, envKey := range e.envKeyMap {
			if value := os.Getenv(envKey); value != "" {
				result[configKey] = value
			}
		}
	}

	return result, nil
}

// Has checks if a configuration key exists in environment variables.
func (e *EnvDriver) Has(key string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Check specific mapping first
	if envKey, exists := e.envKeyMap[key]; exists {
		_, exists := os.LookupEnv(envKey)
		return exists
	}

	// Use viper's IsSet method
	return e.viper.IsSet(key)
}

// Delete removes an environment variable.
func (e *EnvDriver) Delete(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if there's a specific env var mapping
	if envKey, exists := e.envKeyMap[key]; exists {
		return os.Unsetenv(envKey)
	}

	// Generate environment variable name
	envKey := e.generateEnvKey(key)
	return os.Unsetenv(envKey)
}

// BindEnv binds a configuration key to a specific environment variable.
func (e *EnvDriver) BindEnv(configKey string, envKeys ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.envKeyMap == nil {
		e.envKeyMap = make(map[string]string)
	}

	if len(envKeys) > 0 {
		e.envKeyMap[configKey] = envKeys[0]
		args := append([]string{configKey}, envKeys...)
		return e.viper.BindEnv(args...)
	}

	return e.viper.BindEnv(configKey)
}

// SetEnvPrefix sets the environment variable prefix.
func (e *EnvDriver) SetEnvPrefix(prefix string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.prefix = prefix
	e.viper.SetEnvPrefix(prefix)
}

// SetEnvKeyReplacer sets the strings.Replacer for env key transformation.
func (e *EnvDriver) SetEnvKeyReplacer(r *strings.Replacer) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.replacer = r
	e.viper.SetEnvKeyReplacer(r)
}

// Helper methods

// generateEnvKey converts a config key to an environment variable name.
func (e *EnvDriver) generateEnvKey(configKey string) string {
	envKey := configKey

	// Apply replacer transformations
	if e.replacer != nil {
		envKey = e.replacer.Replace(envKey)
	} else {
		// Default transformation: dots and dashes to underscores
		envKey = strings.ReplaceAll(envKey, ".", "_")
		envKey = strings.ReplaceAll(envKey, "-", "_")
	}

	// Convert to uppercase
	envKey = strings.ToUpper(envKey)

	// Add prefix if specified
	if e.prefix != "" {
		envKey = e.prefix + envKey
	}

	return envKey
}

// envKeyToConfigKey converts an environment variable name to a config key.
func (e *EnvDriver) envKeyToConfigKey(envKey string) string {
	configKey := envKey

	// Remove prefix if present
	if e.prefix != "" && strings.HasPrefix(configKey, e.prefix) {
		configKey = strings.TrimPrefix(configKey, e.prefix)
	}

	// Convert to lowercase
	configKey = strings.ToLower(configKey)

	// Apply reverse replacer transformations (underscore to dot)
	if e.replacer != nil {
		// This is a simplified reverse transformation
		// A full implementation would need a reverse replacer
		configKey = strings.ReplaceAll(configKey, "_", ".")
	} else {
		// Default reverse transformation: underscores to dots
		configKey = strings.ReplaceAll(configKey, "_", ".")
	}

	return configKey
}

// Ensure EnvDriver implements the Driver interface
var _ configInterfaces.DriverInterface = (*EnvDriver)(nil)
