// Package config provides configuration management functionality for the GoVel framework.
// This package implements a flexible configuration system that supports multiple
// sources, environments, and data types with dot notation access.
//
// The config system supports:
// - Environment-based configurations
// - Multiple file formats (JSON, YAML, TOML planned)
// - Dot notation for nested values (e.g., "database.host")
// - Type-safe value retrieval
// - Thread-safe operations
// - Default value fallbacks
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"govel/types/src/interfaces/config"
)

// Config represents a configuration manager that handles loading,
// storing, and retrieving configuration values from various sources.
//
// The configuration manager supports hierarchical key access using
// dot notation and provides type-safe retrieval methods.
//
// Example usage:
//
//	config := config.New()
//	config.Set("database.host", "localhost")
//	config.Set("database.port", 5432)
//
//	host := config.GetString("database.host", "127.0.0.1")
//	port := config.GetInt("database.port", 3306)
type Config struct {
	// data holds the configuration data in a nested map structure
	data map[string]interface{}

	// mutex provides thread-safe access to configuration data
	mutex sync.RWMutex

	// environment holds the current environment name
	environment string

	// environmentPrefix is the prefix for environment variable lookups
	environmentPrefix string
}

// New creates a new configuration manager with default settings.
//
// Returns:
//
//	*Config: A new configuration manager ready for use
//
// Example:
//
//	config := config.New()
//	config.Set("application.name", "GoVel Application")
func New() *Config {
	return &Config{
		data:              make(map[string]interface{}),
		environment:       "development",
		environmentPrefix: "",
	}
}

// NewWithEnvironment creates a new configuration manager for a specific environment.
//
// Parameters:
//
//	environment: The environment name (e.g., "development", "production")
//
// Returns:
//
//	*Config: A new configuration manager configured for the environment
//
// Example:
//
//	config := config.NewWithEnvironment("production")
func NewWithEnvironment(environment string) *Config {
	return &Config{
		data:              make(map[string]interface{}),
		environment:       environment,
		environmentPrefix: "",
	}
}

// SetEnvironmentPrefix sets the prefix for environment variable lookups.
// When set, environment variables will be looked up with this prefix.
//
// Parameters:
//
//	prefix: The prefix to use for environment variables
//
// Example:
//
//	config.SetEnvironmentPrefix("GOVEL_")
//	// Now config.GetString("database.host", "") will look for GOVEL_DATABASE_HOST
func (c *Config) SetEnvironmentPrefix(prefix string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.environmentPrefix = prefix
}

// LoadFromFile loads configuration data from a JSON file.
// The file contents will be merged with existing configuration data.
//
// Parameters:
//
//	filename: Path to the JSON configuration file
//
// Returns:
//
//	error: Any error that occurred during file loading, nil if successful
//
// Example:
//
//	err := config.LoadFromFile("config/application.json")
//	if err != nil {
//		log.Fatal("Failed to load config:", err)
//	}
func (c *Config) LoadFromFile(filename string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", filename)
	}

	// Read file contents
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON data
	var fileData map[string]interface{}
	if err := json.Unmarshal(data, &fileData); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	// Merge with existing data
	c.mergeData(fileData)
	return nil
}

// LoadFromEnvironment loads configuration from environment variables.
// Environment variables are converted to lowercase with dots for nesting.
// For example: DATABASE_HOST becomes "database.host"
//
// Returns:
//
//	error: Any error that occurred during environment loading
//
// Example:
//
//	// With environment variables: DATABASE_HOST=localhost, DATABASE_PORT=5432
//	config.LoadFromEnvironment()
//	host := config.GetString("database.host", "") // Returns "localhost"
func (c *Config) LoadFromEnvironment() error {
	return c.LoadFromEnv("")
}

// LoadFromEnv loads configuration from environment variables with a prefix (interface-compatible method).
//
// Parameters:
//
//	prefix: Environment variable prefix (e.g., "APP_")
//
// Returns:
//
//	error: Any error that occurred during environment loading
func (c *Config) LoadFromEnv(prefix string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	envVars := os.Environ()
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Skip if doesn't match prefix
		if prefix != "" {
			if !strings.HasPrefix(key, prefix) {
				continue
			}
			key = strings.TrimPrefix(key, prefix)
		}

		// Convert to config key format
		configKey := strings.ToLower(strings.ReplaceAll(key, "_", "."))

		// Try to parse as different types
		parsedValue := c.parseValue(value)
		c.setNestedValue(configKey, parsedValue)
	}

	return nil
}

// Set sets a configuration value using dot notation.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	value: The value to set
//
// Example:
//
//	config.Set("database.host", "localhost")
//	config.Set("database.connections.max", 100)
func (c *Config) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.setNestedValue(key, value)
}

// GetWithDefault retrieves a configuration value using dot notation.
// Returns the default value if the key is not found.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: The default value to return if key is not found
//
// Returns:
//
//	interface{}: The configuration value or default value
//
// Example:
//
//	value := config.GetWithDefault("database.host", "localhost")
func (c *Config) GetWithDefault(key string, defaultValue interface{}) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// First try to get from configuration data
	if value := c.getNestedValue(key); value != nil {
		return value
	}

	// Try environment variable as fallback
	if envValue := c.getEnvironmentValue(key); envValue != "" {
		return c.parseValue(envValue)
	}

	return defaultValue
}

// Get retrieves a configuration value (interface-compatible method).
//
// Parameters:
//
//	key: The configuration key using dot notation
//
// Returns:
//
//	interface{}: The raw configuration value or nil if not found
//	bool: true if the key exists, false otherwise
func (c *Config) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// First try to get from configuration data
	if value := c.getNestedValue(key); value != nil {
		return value, true
	}

	// Try environment variable as fallback
	if envValue := c.getEnvironmentValue(key); envValue != "" {
		return c.parseValue(envValue), true
	}

	return nil, false
}

// GetString retrieves a configuration value as a string.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default string value to return if key is not found (empty string if not provided)
//
// Returns:
//
//	string: The configuration value as a string or default value
//
// Example:
//
//	host := config.GetString("database.host", "localhost") // With default
//	name := config.GetString("app.name")                   // Without default
func (c *Config) GetString(key string, defaultValue ...string) string {
	default_ := ""
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// GetInt retrieves a configuration value as an integer.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default integer value to return if key is not found (0 if not provided)
//
// Returns:
//
//	int: The configuration value as an integer or default value
//
// Example:
//
//	port := config.GetInt("database.port", 3306) // With default
//	workers := config.GetInt("server.workers")   // Without default
func (c *Config) GetInt(key string, defaultValue ...int) int {
	default_ := 0
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	return default_
}

// GetBool retrieves a configuration value as a boolean.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default boolean value to return if key is not found (false if not provided)
//
// Returns:
//
//	bool: The configuration value as a boolean or default value
//
// Example:
//
//	debug := config.GetBool("application.debug", false) // With default
//	enabled := config.GetBool("feature.enabled")       // Without default
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	default_ := false
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}

	return default_
}

// GetInt64 retrieves a configuration value as an int64.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default int64 value to return if key is not found (0 if not provided)
//
// Returns:
//
//	int64: The configuration value as an int64 or default value
func (c *Config) GetInt64(key string, defaultValue ...int64) int64 {
	default_ := int64(0)
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}

	return default_
}

// GetFloat64 retrieves a configuration value as a float64.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default float64 value to return if key is not found (0.0 if not provided)
//
// Returns:
//
//	float64: The configuration value as a float64 or default value
//
// Example:
//
//	timeout := config.GetFloat64("server.timeout", 30.0) // With default
//	rate := config.GetFloat64("cache.hit_rate")         // Without default
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	default_ := 0.0
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	return default_
}

// GetDuration retrieves a configuration value as a time.Duration.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default duration to return if key is not found (0 if not provided)
//
// Returns:
//
//	time.Duration: The configuration duration or default value
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	default_ := time.Duration(0)
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case time.Duration:
		return v
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	case int:
		// Assume seconds if integer
		return time.Duration(v) * time.Second
	case int64:
		// Assume seconds if integer
		return time.Duration(v) * time.Second
	case float64:
		// Assume seconds if float
		return time.Duration(v * float64(time.Second))
	}

	return default_
}

// GetStringSlice retrieves a configuration value as a string slice.
//
// Parameters:
//
//	key: The configuration key using dot notation
//	defaultValue: Optional default slice to return if key is not found (empty slice if not provided)
//
// Returns:
//
//	[]string: The configuration slice or default value
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	default_ := []string{}
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	value := c.GetWithDefault(key, default_)

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		// Convert interface slice to string slice
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case string:
		// Split comma-separated string
		if v == "" {
			return default_
		}
		parts := strings.Split(v, ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return parts
	}

	return default_
}

// HasKey checks if a configuration key exists.
//
// Parameters:
//
//	key: The configuration key using dot notation
//
// Returns:
//
//	bool: true if the key exists, false otherwise
//
// Example:
//
//	if config.HasKey("database.host") {
//		host := config.GetString("database.host", "")
//	}
func (c *Config) HasKey(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.getNestedValue(key) != nil || c.getEnvironmentValue(key) != ""
}

// GetEnvironment returns the current environment name.
//
// Returns:
//
//	string: The current environment name
//
// Example:
//
//	env := config.GetEnvironment()
//	if env == "production" {
//		// Production-specific logic
//	}
func (c *Config) GetEnvironment() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.environment
}

// SetEnvironment sets the current environment name.
//
// Parameters:
//
//	environment: The environment name to set
//
// Example:
//
//	config.SetEnvironment("production")
func (c *Config) SetEnvironment(environment string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.environment = environment
}

// GetAll returns a copy of all configuration data.
//
// Returns:
//
//	map[string]interface{}: A copy of all configuration data
//
// Example:
//
//	allConfig := config.GetAll()
//	for key, value := range allConfig {
//		fmt.Printf("%s: %v\n", key, value)
//	}
func (c *Config) GetAll() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Create a deep copy of the data
	result := make(map[string]interface{})
	c.copyMap(c.data, result)
	return result
}

// AllConfig returns all configuration values (interface-compatible method).
//
// Returns:
//
//	map[string]interface{}: All configuration key-value pairs
func (c *Config) AllConfig() map[string]interface{} {
	return c.GetAll()
}

// setNestedValue sets a value in the nested map structure using dot notation.
func (c *Config) setNestedValue(key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := c.data

	// Navigate to the parent of the target key
	for _, k := range keys[:len(keys)-1] {
		if current[k] == nil {
			current[k] = make(map[string]interface{})
		}

		// Ensure we have a map at this level
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			// Overwrite non-map value with a new map
			newMap := make(map[string]interface{})
			current[k] = newMap
			current = newMap
		}
	}

	// Set the final value
	current[keys[len(keys)-1]] = value
}

// getNestedValue retrieves a value from the nested map structure using dot notation.
func (c *Config) getNestedValue(key string) interface{} {
	keys := strings.Split(key, ".")
	current := c.data

	for _, k := range keys[:len(keys)-1] {
		if nested, ok := current[k].(map[string]interface{}); ok {
			current = nested
		} else {
			return nil
		}
	}

	return current[keys[len(keys)-1]]
}

// getEnvironmentValue retrieves a value from environment variables.
func (c *Config) getEnvironmentValue(key string) string {
	// Convert config key to environment variable format
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if c.environmentPrefix != "" {
		envKey = c.environmentPrefix + envKey
	}
	return os.Getenv(envKey)
}

// parseValue attempts to parse a string value to its appropriate type.
func (c *Config) parseValue(value string) interface{} {
	// Try boolean
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	// Try integer
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	// Try float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Return as string
	return value
}

// mergeData merges new data into existing configuration data.
func (c *Config) mergeData(newData map[string]interface{}) {
	for key, value := range newData {
		if existingValue, exists := c.data[key]; exists {
			// If both values are maps, merge recursively
			if existingMap, ok := existingValue.(map[string]interface{}); ok {
				if newMap, ok := value.(map[string]interface{}); ok {
					c.mergeMap(existingMap, newMap)
					continue
				}
			}
		}
		c.data[key] = value
	}
}

// mergeMap recursively merges two maps.
func (c *Config) mergeMap(existing, new map[string]interface{}) {
	for key, value := range new {
		if existingValue, exists := existing[key]; exists {
			if existingMap, ok := existingValue.(map[string]interface{}); ok {
				if newMap, ok := value.(map[string]interface{}); ok {
					c.mergeMap(existingMap, newMap)
					continue
				}
			}
		}
		existing[key] = value
	}
}

// copyMap creates a deep copy of a map.
func (c *Config) copyMap(src, dst map[string]interface{}) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]interface{}); ok {
			dstMap := make(map[string]interface{})
			c.copyMap(srcMap, dstMap)
			dst[key] = dstMap
		} else {
			dst[key] = value
		}
	}
}

// Compile-time interface compliance check
// This ensures that Config implements the ConfigInterface
var _ interfaces.ConfigInterface = (*Config)(nil)
