package tests

import (
	"reflect"
	"testing"
	"time"

	"govel/packages/config"
)

// TestConfigOptionalDefaultValues tests that all Get* methods work correctly
// with optional default values using variadic parameters
func TestConfigOptionalDefaultValues(t *testing.T) {
	cfg := config.New()

	// Set some values for testing
	cfg.Set("app.name", "GoVel Framework")
	cfg.Set("database.port", 5432)
	cfg.Set("feature.enabled", true)
	cfg.Set("server.timeout", "30s")
	cfg.Set("rate.limit", 95.5)
	cfg.Set("max.size", int64(1048576))
	cfg.Set("allowed.hosts", []string{"localhost", "127.0.0.1"})

	t.Run("GetString with optional defaults", func(t *testing.T) {
		// Test existing value without default
		result := cfg.GetString("app.name")
		if result != "GoVel Framework" {
			t.Errorf("Expected 'GoVel Framework', got '%s'", result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetString("app.name", "Default App")
		if result != "GoVel Framework" {
			t.Errorf("Expected 'GoVel Framework', got '%s'", result)
		}

		// Test missing value without default (should return empty string)
		result = cfg.GetString("missing.key")
		if result != "" {
			t.Errorf("Expected empty string, got '%s'", result)
		}

		// Test missing value with default (should return default)
		result = cfg.GetString("missing.key", "Fallback Value")
		if result != "Fallback Value" {
			t.Errorf("Expected 'Fallback Value', got '%s'", result)
		}
	})

	t.Run("GetInt with optional defaults", func(t *testing.T) {
		// Test existing value without default
		result := cfg.GetInt("database.port")
		if result != 5432 {
			t.Errorf("Expected 5432, got %d", result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetInt("database.port", 3306)
		if result != 5432 {
			t.Errorf("Expected 5432, got %d", result)
		}

		// Test missing value without default (should return 0)
		result = cfg.GetInt("missing.port")
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}

		// Test missing value with default (should return default)
		result = cfg.GetInt("missing.port", 8080)
		if result != 8080 {
			t.Errorf("Expected 8080, got %d", result)
		}
	})

	t.Run("GetBool with optional defaults", func(t *testing.T) {
		// Test existing value without default
		result := cfg.GetBool("feature.enabled")
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetBool("feature.enabled", false)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}

		// Test missing value without default (should return false)
		result = cfg.GetBool("missing.flag")
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}

		// Test missing value with default (should return default)
		result = cfg.GetBool("missing.flag", true)
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("GetDuration with optional defaults", func(t *testing.T) {
		expectedDuration := 30 * time.Second

		// Test existing value without default
		result := cfg.GetDuration("server.timeout")
		if result != expectedDuration {
			t.Errorf("Expected %v, got %v", expectedDuration, result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetDuration("server.timeout", 60*time.Second)
		if result != expectedDuration {
			t.Errorf("Expected %v, got %v", expectedDuration, result)
		}

		// Test missing value without default (should return 0)
		result = cfg.GetDuration("missing.duration")
		if result != 0 {
			t.Errorf("Expected 0s, got %v", result)
		}

		// Test missing value with default (should return default)
		defaultDuration := 45 * time.Second
		result = cfg.GetDuration("missing.duration", defaultDuration)
		if result != defaultDuration {
			t.Errorf("Expected %v, got %v", defaultDuration, result)
		}
	})

	t.Run("GetFloat64 with optional defaults", func(t *testing.T) {
		// Test existing value without default
		result := cfg.GetFloat64("rate.limit")
		if result != 95.5 {
			t.Errorf("Expected 95.5, got %f", result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetFloat64("rate.limit", 80.0)
		if result != 95.5 {
			t.Errorf("Expected 95.5, got %f", result)
		}

		// Test missing value without default (should return 0.0)
		result = cfg.GetFloat64("missing.rate")
		if result != 0.0 {
			t.Errorf("Expected 0.0, got %f", result)
		}

		// Test missing value with default (should return default)
		result = cfg.GetFloat64("missing.rate", 75.5)
		if result != 75.5 {
			t.Errorf("Expected 75.5, got %f", result)
		}
	})

	t.Run("GetInt64 with optional defaults", func(t *testing.T) {
		// Test existing value without default
		result := cfg.GetInt64("max.size")
		if result != 1048576 {
			t.Errorf("Expected 1048576, got %d", result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetInt64("max.size", 2097152)
		if result != 1048576 {
			t.Errorf("Expected 1048576, got %d", result)
		}

		// Test missing value without default (should return 0)
		result = cfg.GetInt64("missing.size")
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}

		// Test missing value with default (should return default)
		result = cfg.GetInt64("missing.size", 4194304)
		if result != 4194304 {
			t.Errorf("Expected 4194304, got %d", result)
		}
	})

	t.Run("GetStringSlice with optional defaults", func(t *testing.T) {
		expectedSlice := []string{"localhost", "127.0.0.1"}

		// Test existing value without default
		result := cfg.GetStringSlice("allowed.hosts")
		if !reflect.DeepEqual(result, expectedSlice) {
			t.Errorf("Expected %v, got %v", expectedSlice, result)
		}

		// Test existing value with default (should return existing value)
		result = cfg.GetStringSlice("allowed.hosts", []string{"fallback.com"})
		if !reflect.DeepEqual(result, expectedSlice) {
			t.Errorf("Expected %v, got %v", expectedSlice, result)
		}

		// Test missing value without default (should return empty slice)
		result = cfg.GetStringSlice("missing.slice")
		expectedEmpty := []string{}
		if !reflect.DeepEqual(result, expectedEmpty) {
			t.Errorf("Expected %v, got %v", expectedEmpty, result)
		}

		// Test missing value with default (should return default)
		defaultSlice := []string{"default1", "default2"}
		result = cfg.GetStringSlice("missing.slice", defaultSlice)
		if !reflect.DeepEqual(result, defaultSlice) {
			t.Errorf("Expected %v, got %v", defaultSlice, result)
		}
	})
}

// TestConfigOptionalDefaultsRealisticUsage tests realistic usage patterns
// combining optional defaults with application-level fallbacks
func TestConfigOptionalDefaultsRealisticUsage(t *testing.T) {
	cfg := config.New()

	// Set some values to simulate partial configuration
	cfg.Set("database.port", 5432)

	t.Run("Mixed usage patterns", func(t *testing.T) {
		// Application-level fallback pattern
		dbHost := cfg.GetString("database.host") // Will return "" if not set
		if dbHost == "" {
			dbHost = "localhost" // Application-level fallback
		}

		// Config-level default pattern
		dbPort := cfg.GetInt("database.port", 3306) // With config-level default

		// No default, use zero value pattern
		debug := cfg.GetBool("app.debug") // No default, false if not set

		// Default with meaningful fallback
		timeout := cfg.GetDuration("api.timeout", 30*time.Second) // With default

		// Validate results
		if dbHost != "localhost" {
			t.Errorf("Expected 'localhost' for dbHost, got '%s'", dbHost)
		}

		if dbPort != 5432 {
			t.Errorf("Expected 5432 for dbPort, got %d", dbPort)
		}

		if debug != false {
			t.Errorf("Expected false for debug, got %v", debug)
		}

		if timeout != 30*time.Second {
			t.Errorf("Expected 30s for timeout, got %v", timeout)
		}
	})
}

// TestConfigOptionalDefaultsEdgeCases tests edge cases and boundary conditions
func TestConfigOptionalDefaultsEdgeCases(t *testing.T) {
	cfg := config.New()

	t.Run("Multiple default values (should use first)", func(t *testing.T) {
		// Only the first default value should be used
		result := cfg.GetString("missing.key", "first", "second", "third")
		if result != "first" {
			t.Errorf("Expected 'first', got '%s'", result)
		}

		intResult := cfg.GetInt("missing.int", 100, 200, 300)
		if intResult != 100 {
			t.Errorf("Expected 100, got %d", intResult)
		}
	})

	t.Run("Empty slice as default", func(t *testing.T) {
		result := cfg.GetStringSlice("missing.slice", []string{})
		expected := []string{}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Zero duration as default", func(t *testing.T) {
		result := cfg.GetDuration("missing.duration", 0)
		if result != 0 {
			t.Errorf("Expected 0s, got %v", result)
		}
	})

	t.Run("Negative values as defaults", func(t *testing.T) {
		intResult := cfg.GetInt("missing.int", -1)
		if intResult != -1 {
			t.Errorf("Expected -1, got %d", intResult)
		}

		floatResult := cfg.GetFloat64("missing.float", -1.5)
		if floatResult != -1.5 {
			t.Errorf("Expected -1.5, got %f", floatResult)
		}

		int64Result := cfg.GetInt64("missing.int64", -1000)
		if int64Result != -1000 {
			t.Errorf("Expected -1000, got %d", int64Result)
		}
	})
}
