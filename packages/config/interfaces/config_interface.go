package interfaces

import "time"

// ConfigInterface defines the contract for configuration management.
// This interface provides a standardized API for accessing configuration
// values with dot notation support, environment variable integration,
// and type-safe retrieval methods.
//
// The configuration system supports multiple sources including configuration
// files, environment variables, command-line arguments, and programmatic
// configuration. It provides type-safe methods for retrieving common data
// types with optional default values.
//
// Example usage:
//
//	config := &Config{}
//	host := config.GetString("database.host", "localhost")
//	port := config.GetInt("database.port", 3306)
//	debug := config.GetBool("application.debug", false)
//	timeout := config.GetDuration("server.timeout", 30*time.Second)
//
//	// Set configuration values
//	config.Set("application.name", "My Application")
//	config.Set("database.host", "remote-db.example.com")
//
// The interface promotes:
// - Centralized configuration management
// - Type-safe configuration access
// - Environment variable integration
// - Easy testing with mock configurations
// - Flexible configuration sources
type ConfigInterface interface {
	// GetString retrieves a string configuration value.
	// Uses dot notation for nested configuration keys.
	//
	// Parameters:
	//   key: The configuration key using dot notation (e.g., "database.host")
	//   defaultValue: The default value to return if key is not found
	//
	// Returns:
	//   string: The configuration value or default if not found
	//
	// Example:
	//   host := config.GetString("database.host", "localhost")
	//   appName := config.GetString("application.name", "Default Application")
	GetString(key, defaultValue string) string

	// GetInt retrieves an integer configuration value.
	// Automatically converts string values to integers when possible.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default value to return if key is not found or invalid
	//
	// Returns:
	//   int: The configuration value or default if not found/invalid
	//
	// Example:
	//   port := config.GetInt("database.port", 3306)
	//   workers := config.GetInt("server.workers", 4)
	GetInt(key string, defaultValue int) int

	// GetInt64 retrieves a 64-bit integer configuration value.
	// Useful for large numbers, timestamps, or IDs.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default value to return if key is not found or invalid
	//
	// Returns:
	//   int64: The configuration value or default if not found/invalid
	//
	// Example:
	//   maxSize := config.GetInt64("upload.max_size", 10485760) // 10MB default
	GetInt64(key string, defaultValue int64) int64

	// GetFloat64 retrieves a floating-point configuration value.
	// Useful for percentages, ratios, or decimal values.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default value to return if key is not found or invalid
	//
	// Returns:
	//   float64: The configuration value or default if not found/invalid
	//
	// Example:
	//   rate := config.GetFloat64("cache.hit_rate", 0.95)
	//   timeout := config.GetFloat64("request.timeout_seconds", 30.0)
	GetFloat64(key string, defaultValue float64) float64

	// GetBool retrieves a boolean configuration value.
	// Recognizes common boolean representations: true/false, yes/no, 1/0.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default value to return if key is not found or invalid
	//
	// Returns:
	//   bool: The configuration value or default if not found/invalid
	//
	// Example:
	//   debug := config.GetBool("application.debug", false)
	//   enabled := config.GetBool("feature.new_ui", true)
	GetBool(key string, defaultValue bool) bool

	// GetDuration retrieves a duration configuration value.
	// Supports various duration formats: "30s", "5m", "2h", "24h".
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default duration to return if key is not found or invalid
	//
	// Returns:
	//   time.Duration: The configuration duration or default if not found/invalid
	//
	// Example:
	//   timeout := config.GetDuration("server.timeout", 30*time.Second)
	//   interval := config.GetDuration("cleanup.interval", 24*time.Hour)
	GetDuration(key string, defaultValue time.Duration) time.Duration

	// GetStringSlice retrieves a string slice configuration value.
	// Supports comma-separated values or array formats depending on source.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   defaultValue: The default slice to return if key is not found
	//
	// Returns:
	//   []string: The configuration slice or default if not found
	//
	// Example:
	//   hosts := config.GetStringSlice("database.hosts", []string{"localhost"})
	//   features := config.GetStringSlice("application.enabled_features", []string{})
	GetStringSlice(key string, defaultValue []string) []string

	// Get retrieves a raw configuration value as interface{}.
	// Useful when the exact type is unknown or for custom type handling.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//
	// Returns:
	//   interface{}: The raw configuration value or nil if not found
	//   bool: true if the key exists, false otherwise
	//
	// Example:
	//   value, exists := config.Get("custom.setting")
	//   if exists {
	//       // Process the raw value
	//   }
	Get(key string) (interface{}, bool)

	// Set stores a configuration value.
	// Useful for programmatic configuration or runtime updates.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//   value: The value to store
	//
	// Example:
	//   config.Set("application.name", "My Application")
	//   config.Set("database.port", 5432)
	Set(key string, value interface{})

	// HasKey checks if a configuration key exists.
	//
	// Parameters:
	//   key: The configuration key using dot notation
	//
	// Returns:
	//   bool: true if the key exists, false otherwise
	//
	// Example:
	//   if config.HasKey("database.host") {
	//       host := config.GetString("database.host", "")
	//   }
	HasKey(key string) bool

	// AllConfig returns all configuration values as a map.
	// Useful for debugging, serialization, or bulk operations.
	//
	// Returns:
	//   map[string]interface{}: All configuration key-value pairs
	//
	// Example:
	//   allConfig := config.AllConfig()
	//   for key, value := range allConfig {
	//       fmt.Printf("%s: %v\n", key, value)
	//   }
	AllConfig() map[string]interface{}

	// LoadFromFile loads configuration from a file.
	// Supports various formats: JSON, YAML, TOML, etc.
	//
	// Parameters:
	//   filePath: Path to the configuration file
	//
	// Returns:
	//   error: Any error that occurred during loading
	//
	// Example:
	//   err := config.LoadFromFile("config/application.yaml")
	//   if err != nil {
	//       log.Fatal("Failed to load config:", err)
	//   }
	LoadFromFile(filePath string) error

	// LoadFromEnv loads configuration from environment variables.
	// Typically uses a prefix to filter relevant environment variables.
	//
	// Parameters:
	//   prefix: Environment variable prefix (e.g., "APP_")
	//
	// Returns:
	//   error: Any error that occurred during loading
	//
	// Example:
	//   err := config.LoadFromEnv("MYAPP_")
	//   // This would load MYAPP_DATABASE_HOST as database.host
	LoadFromEnv(prefix string) error
}
