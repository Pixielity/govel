package tests

import (
	"testing"
	"time"

	"govel/config"
)

// TestConfigIntegerValues tests integer configuration values
func TestConfigIntegerValues(t *testing.T) {
	cfg := config.New()
	
	// Set integer value
	cfg.Set("app.port", 8080)
	port := cfg.GetInt("app.port", 3000)
	
	if port != 8080 {
		t.Errorf("Expected 8080, got %d", port)
	}
	
	// Test default for non-existent int
	defaultPort := cfg.GetInt("app.default_port", 3000)
	if defaultPort != 3000 {
		t.Errorf("Expected 3000, got %d", defaultPort)
	}
	
	// Test negative integers
	cfg.Set("app.offset", -10)
	offset := cfg.GetInt("app.offset", 0)
	
	if offset != -10 {
		t.Errorf("Expected -10, got %d", offset)
	}
}

// TestConfigBooleanValues tests boolean configuration values
func TestConfigBooleanValues(t *testing.T) {
	cfg := config.New()
	
	// Set boolean value
	cfg.Set("app.debug", true)
	debug := cfg.GetBool("app.debug", false)
	
	if !debug {
		t.Error("Expected true, got false")
	}
	
	// Test default for non-existent bool
	production := cfg.GetBool("app.production", false)
	if production {
		t.Error("Expected false, got true")
	}
	
	// Test setting false explicitly
	cfg.Set("app.maintenance", false)
	maintenance := cfg.GetBool("app.maintenance", true)
	
	if maintenance {
		t.Error("Expected false, got true")
	}
}

// TestConfigFloatValues tests float configuration values
func TestConfigFloatValues(t *testing.T) {
	cfg := config.New()
	
	// Set float value
	cfg.Set("app.rate", 3.14159)
	rate := cfg.GetFloat64("app.rate", 0.0)
	
	if rate != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", rate)
	}
	
	// Test default for non-existent float
	defaultRate := cfg.GetFloat64("app.default_rate", 1.0)
	if defaultRate != 1.0 {
		t.Errorf("Expected 1.0, got %f", defaultRate)
	}
}

// TestConfigInt64Values tests int64 configuration values
func TestConfigInt64Values(t *testing.T) {
	cfg := config.New()
	
	// Set int64 value
	var largeNumber int64 = 9223372036854775807 // max int64
	cfg.Set("app.large_number", largeNumber)
	retrievedNumber := cfg.GetInt64("app.large_number", 0)
	
	if retrievedNumber != largeNumber {
		t.Errorf("Expected %d, got %d", largeNumber, retrievedNumber)
	}
	
	// Test default for non-existent int64
	defaultLarge := cfg.GetInt64("app.default_large", 1000000)
	if defaultLarge != 1000000 {
		t.Errorf("Expected 1000000, got %d", defaultLarge)
	}
}

// TestConfigDurationValues tests duration configuration values
func TestConfigDurationValues(t *testing.T) {
	cfg := config.New()
	
	// Set duration value
	timeout := 30 * time.Second
	cfg.Set("app.timeout", timeout)
	retrievedTimeout := cfg.GetDuration("app.timeout", time.Minute)
	
	if retrievedTimeout != timeout {
		t.Errorf("Expected %v, got %v", timeout, retrievedTimeout)
	}
	
	// Test default for non-existent duration
	defaultTimeout := cfg.GetDuration("app.default_timeout", 5*time.Minute)
	if defaultTimeout != 5*time.Minute {
		t.Errorf("Expected %v, got %v", 5*time.Minute, defaultTimeout)
	}
}

// TestConfigStringSliceValues tests string slice configuration values
func TestConfigStringSliceValues(t *testing.T) {
	cfg := config.New()
	
	// Set slice value
	hosts := []string{"localhost", "127.0.0.1", "::1"}
	cfg.Set("database.hosts", hosts)
	
	retrievedHosts := cfg.GetStringSlice("database.hosts", []string{})
	if len(retrievedHosts) != 3 {
		t.Errorf("Expected 3 hosts, got %d", len(retrievedHosts))
	}
	
	// Verify the actual values
	for i, expectedHost := range hosts {
		if i < len(retrievedHosts) && retrievedHosts[i] != expectedHost {
			t.Errorf("Expected host %s at index %d, got %s", expectedHost, i, retrievedHosts[i])
		}
	}
	
	// Test default for non-existent slice
	defaultHosts := cfg.GetStringSlice("database.default_hosts", []string{"fallback"})
	if len(defaultHosts) != 1 || defaultHosts[0] != "fallback" {
		t.Errorf("Expected default slice with 'fallback', got %v", defaultHosts)
	}
}

// TestConfigRawGet tests the raw Get method
func TestConfigRawGet(t *testing.T) {
	cfg := config.New()
	
	// Set various types
	cfg.Set("string.value", "test_string")
	cfg.Set("int.value", 42)
	cfg.Set("bool.value", true)
	
	// Test getting existing values
	stringVal, exists := cfg.Get("string.value")
	if !exists {
		t.Error("Expected string.value to exist")
	}
	if stringVal != "test_string" {
		t.Errorf("Expected 'test_string', got %v", stringVal)
	}
	
	intVal, exists := cfg.Get("int.value")
	if !exists {
		t.Error("Expected int.value to exist")
	}
	if intVal != 42 {
		t.Errorf("Expected 42, got %v", intVal)
	}
	
	// Test getting non-existent value
	_, exists = cfg.Get("non.existent")
	if exists {
		t.Error("Expected non.existent to not exist")
	}
}
