// Package support provides environment variable utilities for Go applications.
//
// This module provides Laravel-style environment variable handling with
// type-safe conversions, default values, and validation. It supports all
// common data types and provides convenient methods for accessing configuration
// from environment variables.
//
// Key Features:
//   - Type-safe environment variable access
//   - Default value support for all types
//   - Boolean parsing with multiple formats
//   - Array/slice parsing with custom separators
//   - URL and file path validation
//   - Environment variable existence checking
//   - Cached environment loading for performance
//
// Usage Example:
//
//	// Basic usage
//	dbHost := Get("DB_HOST", "localhost")
//	dbPort := GetInt("DB_PORT", 5432)
//	debug := GetBool("APP_DEBUG", false)
//
//	// Advanced usage
//	allowedHosts := GetArray("ALLOWED_HOSTS", []string{"localhost"}, ",")
//	dbUrl := GetURL("DATABASE_URL", nil)
package support

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// envCache caches environment variables to improve performance
	envCache     = make(map[string]string)
	envCacheMux  sync.RWMutex
	envCacheInit sync.Once
)

// initEnvCache initializes the environment cache with all current environment variables
func initEnvCache() {
	envCacheInit.Do(func() {
		envCacheMux.Lock()
		defer envCacheMux.Unlock()

		for _, env := range os.Environ() {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) == 2 {
				envCache[pair[0]] = pair[1]
			}
		}
	})
}

// Get retrieves an environment variable as a string with optional default value.
//
// This is the base method for environment variable access, similar to Laravel's env() helper.
// It provides a clean interface for accessing string environment variables with fallback values.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or empty
//
// Returns:
//   - string: The environment variable value or default value
//
// Example:
//
//	appName := Get("APP_NAME", "MyApp")
//	secret := Get("APP_SECRET") // No default, returns empty string if not set
func Get(key string, defaultValue ...string) string {
	initEnvCache()

	envCacheMux.RLock()
	value, exists := envCache[key]
	envCacheMux.RUnlock()

	// If variable exists and is not empty, return it
	if exists && value != "" {
		return value
	}

	// Return default value if provided
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	// Return empty string if no default provided
	return ""
}

// GetString is an alias for Get for explicit string type indication.
//
// This provides a more explicit method name when you want to emphasize that
// you're retrieving a string environment variable.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or empty
//
// Returns:
//   - string: The environment variable value or default value
//
// Example:
//
//	appName := GetString("APP_NAME", "MyApp")
//	secret := GetString("APP_SECRET")
func GetString(key string, defaultValue ...string) string {
	return Get(key, defaultValue...)
}

// GetInt retrieves an environment variable as an integer with optional default value.
//
// Converts the environment variable value to an integer using strconv.Atoi.
// If conversion fails, returns the default value or 0.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or conversion fails
//
// Returns:
//   - int: The environment variable value as integer or default value
//
// Example:
//
//	port := GetInt("PORT", 8080)
//	workers := GetInt("WORKERS", 4)
func GetInt(key string, defaultValue ...int) int {
	value := Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetInt64 retrieves an environment variable as an int64 with optional default value.
//
// Similar to GetInt but returns int64 for larger integer values.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or conversion fails
//
// Returns:
//   - int64: The environment variable value as int64 or default value
//
// Example:
//
//	maxSize := GetInt64("MAX_FILE_SIZE", 1073741824) // 1GB default
func GetInt64(key string, defaultValue ...int64) int64 {
	value := Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intValue
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetFloat retrieves an environment variable as a float64 with optional default value.
//
// Converts the environment variable value to a float64 using strconv.ParseFloat.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or conversion fails
//
// Returns:
//   - float64: The environment variable value as float64 or default value
//
// Example:
//
//	rate := GetFloat("RATE_LIMIT", 10.5)
//	timeout := GetFloat("TIMEOUT_SECONDS", 30.0)
func GetFloat(key string, defaultValue ...float64) float64 {
	value := Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}

	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0.0
}

// GetBool retrieves an environment variable as a boolean with optional default value.
//
// Supports multiple boolean formats for flexibility:
//   - true, TRUE, True, 1, yes, YES, Yes, on, ON, On
//   - false, FALSE, False, 0, no, NO, No, off, OFF, Off
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or invalid format
//
// Returns:
//   - bool: The environment variable value as boolean or default value
//
// Example:
//
//	debug := GetBool("DEBUG", false)
//	ssl := GetBool("SSL_ENABLED", true)
func GetBool(key string, defaultValue ...bool) bool {
	value := strings.ToLower(strings.TrimSpace(Get(key)))
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	// Check for truthy values
	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}

	// Try standard bool parsing as fallback
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetDuration retrieves an environment variable as a time.Duration with optional default value.
//
// Supports duration formats like "5s", "10m", "1h", "24h", etc.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default value if variable is not set or parsing fails
//
// Returns:
//   - time.Duration: The environment variable value as Duration or default value
//
// Example:
//
//	timeout := GetDuration("REQUEST_TIMEOUT", 30*time.Second)
//	interval := GetDuration("POLLING_INTERVAL", 5*time.Minute)
func GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	value := Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetArray retrieves an environment variable as a string slice with optional default and separator.
//
// Splits the environment variable value by the specified separator (default: comma).
// Trims whitespace from each element and filters out empty strings.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default slice if variable is not set
//   - separator: Optional separator (default: ",")
//
// Returns:
//   - []string: The environment variable value as string slice or default value
//
// Example:
//
//	hosts := GetArray("ALLOWED_HOSTS", []string{"localhost"}, ",")
//	tags := GetArray("TAGS", []string{}, ";")
func GetArray(key string, defaultValue []string, separator ...string) []string {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	sep := ","
	if len(separator) > 0 {
		sep = separator[0]
	}

	parts := strings.Split(value, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// GetURL retrieves an environment variable as a parsed URL with optional default.
//
// Parses the environment variable value as a URL and validates its format.
// Returns nil if parsing fails and no default is provided.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default URL if variable is not set or invalid
//
// Returns:
//   - *url.URL: The parsed URL or default value
//
// Example:
//
//	dbURL := GetURL("DATABASE_URL", nil)
//	if dbURL != nil {
//	    fmt.Printf("Host: %s, Port: %s", dbURL.Hostname(), dbURL.Port())
//	}
func GetURL(key string, defaultValue *url.URL) *url.URL {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	if parsedURL, err := url.Parse(value); err == nil {
		return parsedURL
	}

	return defaultValue
}

// GetPath retrieves an environment variable as a validated file path with optional default.
//
// Validates that the path exists and optionally checks if it's a file or directory.
// Returns the default value if the path doesn't exist or validation fails.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: Optional default path if variable is not set or invalid
//   - mustExist: Optional flag to require path existence (default: false)
//
// Returns:
//   - string: The validated file path or default value
//
// Example:
//
//	configPath := GetPath("CONFIG_PATH", "/etc/app/config.yaml", true)
//	logDir := GetPath("LOG_DIR", "/var/log/app", false)
func GetPath(key string, defaultValue string, mustExist ...bool) string {
	value := Get(key)
	if value == "" {
		return defaultValue
	}

	// Clean the path
	cleanPath := filepath.Clean(value)

	// Check existence if required
	if len(mustExist) > 0 && mustExist[0] {
		if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
			return defaultValue
		}
	}

	return cleanPath
}

// Exists checks if an environment variable is set (even if empty).
//
// This is useful for distinguishing between unset variables and empty string values.
//
// Parameters:
//   - key: The environment variable name
//
// Returns:
//   - bool: True if the variable is set, false otherwise
//
// Example:
//
//	if Exists("OPTIONAL_FEATURE") {
//	    // Variable is set, use its value (even if empty)
//	    value := Get("OPTIONAL_FEATURE")
//	} else {
//	    // Variable is not set, use different logic
//	}
func Exists(key string) bool {
	initEnvCache()

	envCacheMux.RLock()
	_, exists := envCache[key]
	envCacheMux.RUnlock()

	return exists
}

// SetEnv sets an environment variable for the current process and updates the cache.
//
// This is useful for testing or runtime configuration changes. Note that this
// only affects the current process and its children.
//
// Parameters:
//   - key: The environment variable name
//   - value: The value to set
//
// Returns:
//   - error: Any error that occurred during the set operation
//
// Example:
//
//	err := SetEnv("TEST_VAR", "test_value")
//	if err != nil {
//	    log.Printf("Failed to set environment variable: %v", err)
//	}
func SetEnv(key, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}

	// Update cache
	envCacheMux.Lock()
	envCache[key] = value
	envCacheMux.Unlock()

	return nil
}

// UnsetEnv removes an environment variable from the current process and cache.
//
// Parameters:
//   - key: The environment variable name to remove
//
// Returns:
//   - error: Any error that occurred during the unset operation
//
// Example:
//
//	err := UnsetEnv("TEMP_VAR")
//	if err != nil {
//	    log.Printf("Failed to unset environment variable: %v", err)
//	}
func UnsetEnv(key string) error {
	err := os.Unsetenv(key)
	if err != nil {
		return err
	}

	// Update cache
	envCacheMux.Lock()
	delete(envCache, key)
	envCacheMux.Unlock()

	return nil
}

// GetAllEnv returns a copy of all environment variables as a map.
//
// This is useful for debugging or passing environment context to other systems.
// Returns a copy to prevent external modifications to the cache.
//
// Returns:
//   - map[string]string: A copy of all environment variables
//
// Example:
//
//	allVars := GetAllEnv()
//	for key, value := range allVars {
//	    fmt.Printf("%s=%s\n", key, value)
//	}
func GetAllEnv() map[string]string {
	initEnvCache()

	envCacheMux.RLock()
	defer envCacheMux.RUnlock()

	// Return a copy to prevent external modifications
	result := make(map[string]string, len(envCache))
	for k, v := range envCache {
		result[k] = v
	}

	return result
}

// GetWithPrefix returns all environment variables with a specific prefix as a map.
//
// This is useful for grouping related configuration variables, similar to how
// Laravel groups config by prefixes.
//
// Parameters:
//   - prefix: The prefix to filter by (case-sensitive)
//   - stripPrefix: Optional flag to remove prefix from keys (default: false)
//
// Returns:
//   - map[string]string: Map of matching environment variables
//
// Example:
//
//	dbVars := GetWithPrefix("DB_", true)
//	// Returns {"HOST": "localhost", "PORT": "5432"} instead of
//	// {"DB_HOST": "localhost", "DB_PORT": "5432"}
func GetWithPrefix(prefix string, stripPrefix ...bool) map[string]string {
	initEnvCache()

	envCacheMux.RLock()
	defer envCacheMux.RUnlock()

	result := make(map[string]string)
	strip := len(stripPrefix) > 0 && stripPrefix[0]

	for key, value := range envCache {
		if strings.HasPrefix(key, prefix) {
			if strip {
				result[strings.TrimPrefix(key, prefix)] = value
			} else {
				result[key] = value
			}
		}
	}

	return result
}

// GetRequired retrieves a required environment variable, panicking if not set.
//
// This is useful for critical configuration values that must be present.
// Similar to Laravel's env() function when a variable is absolutely required.
//
// Parameters:
//   - key: The environment variable name
//   - message: Optional custom error message
//
// Returns:
//   - string: The environment variable value
//
// Panics:
//   - If the environment variable is not set or is empty
//
// Example:
//
//	secret := GetRequired("APP_SECRET", "Application secret is required")
//	dbHost := GetRequired("DB_HOST") // Uses default error message
func GetRequired(key string, message ...string) string {
	value := Get(key)
	if value == "" {
		msg := fmt.Sprintf("Required environment variable '%s' is not set", key)
		if len(message) > 0 {
			msg = message[0]
		}
		panic(msg)
	}
	return value
}

// LoadEnvFile loads environment variables from a .env file.
//
// This provides basic .env file support similar to Laravel's vlucas/phpdotenv.
// Supports simple KEY=VALUE format with basic comment support.
//
// Parameters:
//   - filepath: Path to the .env file
//   - override: Whether to override existing environment variables
//
// Returns:
//   - error: Any error that occurred during file loading
//
// Example:
//
//	err := LoadEnvFile(".env", false)
//	if err != nil {
//	    log.Printf("Failed to load .env file: %v", err)
//	}
func LoadEnvFile(filepath string, override bool) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read env file %s: %w", filepath, err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first = sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format at line %d: %s", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
				(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if override is true or variable doesn't exist
		if override || !Exists(key) {
			if err := SetEnv(key, value); err != nil {
				return fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	return nil
}
