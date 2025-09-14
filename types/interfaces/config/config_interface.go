package config

// DriverInterface defines the interface that all configuration drivers must implement.
type DriverInterface interface {
	// Get retrieves a configuration value by key.
	Get(key string) (interface{}, error)

	// Set stores a configuration value by key.
	Set(key string, value interface{}) error

	// Load loads the configuration data.
	Load() error

	// Watch starts watching for configuration changes.
	Watch(callback func()) error

	// Unwatch stops watching for configuration changes.
	Unwatch() error

	// Invalidate clears the configuration cache.
	Invalidate() error

	// GetAll returns all configuration data as a map.
	GetAll() (map[string]interface{}, error)

	// Has checks if a configuration key exists.
	Has(key string) bool

	// Delete removes a configuration key.
	Delete(key string) error
}

// ConfigInterface defines the interface for the main configuration manager.
type ConfigInterface interface {
	// GetString returns a string value for the given key with optional default.
	GetString(key string, defaultValue ...string) string

	// GetInt returns an integer value for the given key with optional default.
	GetInt(key string, defaultValue ...int) int

	// GetBool returns a boolean value for the given key with optional default.
	GetBool(key string, defaultValue ...bool) bool

	// GetFloat64 returns a float64 value for the given key with optional default.
	GetFloat64(key string, defaultValue ...float64) float64

	// Set stores a configuration value by key.
	Set(key string, value interface{})

	// Has checks if a configuration key exists.
	Has(key string) bool

	// GetAll returns all configuration data.
	GetAll() (map[string]interface{}, error)

	// Load loads the configuration.
	Load() error

	// Watch starts watching for configuration changes.
	Watch(callback func()) error

	// Unwatch stops watching for configuration changes.
	Unwatch() error
}