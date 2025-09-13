package interfaces

import "time"

// ConfigInterface defines the contract for configuration management.
// This interface provides a standardized API for accessing configuration
// values with dot notation support, environment variable integration,
// and type-safe retrieval methods.
type ConfigInterface interface {
	// GetString retrieves a string configuration value.
	GetString(key string, defaultValue ...string) string

	// GetInt retrieves an integer configuration value.
	GetInt(key string, defaultValue ...int) int

	// GetInt64 retrieves a 64-bit integer configuration value.
	GetInt64(key string, defaultValue ...int64) int64

	// GetFloat64 retrieves a floating-point configuration value.
	GetFloat64(key string, defaultValue ...float64) float64

	// GetBool retrieves a boolean configuration value.
	GetBool(key string, defaultValue ...bool) bool

	// GetDuration retrieves a duration configuration value.
	GetDuration(key string, defaultValue ...time.Duration) time.Duration

	// GetStringSlice retrieves a string slice configuration value.
	GetStringSlice(key string, defaultValue ...[]string) []string

	// Get retrieves a raw configuration value as interface{}.
	Get(key string) (interface{}, bool)

	// Set stores a configuration value.
	Set(key string, value interface{})

	// HasKey checks if a configuration key exists.
	HasKey(key string) bool

	// AllConfig returns all configuration values as a map.
	AllConfig() map[string]interface{}

	// LoadFromFile loads configuration from a file.
	LoadFromFile(filePath string) error

	// LoadFromEnv loads configuration from environment variables.
	LoadFromEnv(prefix string) error
}