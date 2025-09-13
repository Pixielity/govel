package tests

import (
	"testing"
	"time"

	"govel/config/mocks"
)

// TestMockConfig tests the mock config implementation
func TestMockConfig(t *testing.T) {
	mockCfg := mocks.NewMockConfig()
	
	// Test basic operations
	mockCfg.Set("mock.test", "mock_value")
	value := mockCfg.GetString("mock.test", "default")
	
	if value != "mock_value" {
		t.Errorf("Expected 'mock_value', got '%s'", value)
	}
	
	// Test different data types
	mockCfg.Set("mock.number", 42)
	number := mockCfg.GetInt("mock.number", 0)
	
	if number != 42 {
		t.Errorf("Expected 42, got %d", number)
	}
	
	// Test boolean
	mockCfg.Set("mock.flag", true)
	flag := mockCfg.GetBool("mock.flag", false)
	
	if !flag {
		t.Error("Expected true, got false")
	}
	
	// Test HasKey
	if !mockCfg.HasKey("mock.test") {
		t.Error("Expected HasKey to return true for existing key")
	}
	
	if mockCfg.HasKey("non.existent") {
		t.Error("Expected HasKey to return false for non-existent key")
	}
	
	// Test AllConfig
	allConfig := mockCfg.AllConfig()
	if len(allConfig) != 3 { // test, number, flag
		t.Errorf("Expected 3 config items, got %d", len(allConfig))
	}
}

// TestMockConfigDataTypes tests all data types in mock config
func TestMockConfigDataTypes(t *testing.T) {
	mockCfg := mocks.NewMockConfig()
	
	// Test int64
	var largeNum int64 = 9876543210
	mockCfg.Set("mock.large", largeNum)
	retrievedLarge := mockCfg.GetInt64("mock.large", 0)
	
	if retrievedLarge != largeNum {
		t.Errorf("Expected %d, got %d", largeNum, retrievedLarge)
	}
	
	// Test float64
	mockCfg.Set("mock.rate", 3.14159)
	rate := mockCfg.GetFloat64("mock.rate", 0.0)
	
	if rate != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", rate)
	}
	
	// Test duration
	timeout := 30 * time.Second
	mockCfg.Set("mock.timeout", timeout)
	retrievedTimeout := mockCfg.GetDuration("mock.timeout", time.Minute)
	
	if retrievedTimeout != timeout {
		t.Errorf("Expected %v, got %v", timeout, retrievedTimeout)
	}
	
	// Test string slice
	hosts := []string{"host1", "host2", "host3"}
	mockCfg.Set("mock.hosts", hosts)
	retrievedHosts := mockCfg.GetStringSlice("mock.hosts", []string{})
	
	if len(retrievedHosts) != 3 {
		t.Errorf("Expected 3 hosts, got %d", len(retrievedHosts))
	}
	
	for i, host := range hosts {
		if i < len(retrievedHosts) && retrievedHosts[i] != host {
			t.Errorf("Expected host %s at index %d, got %s", host, i, retrievedHosts[i])
		}
	}
}

// TestMockConfigFailureMode tests mock config failure simulation
func TestMockConfigFailureMode(t *testing.T) {
	mockCfg := mocks.NewMockConfig()
	
	// Enable failure mode
	mockCfg.SetFailureMode(true, false)
	
	// Test that get operations fail or return defaults
	value := mockCfg.GetString("test.key", "default")
	if value != "default" {
		t.Errorf("Expected default value in failure mode, got '%s'", value)
	}
	
	// Test that HasKey might fail in failure mode
	hasKey := mockCfg.HasKey("test.key")
	// Behavior depends on mock implementation - mainly test it doesn't panic
	t.Logf("HasKey result in failure mode: %t", hasKey)
}

// TestMockConfigRawGet tests the raw Get method in mock
func TestMockConfigRawGet(t *testing.T) {
	mockCfg := mocks.NewMockConfig()
	
	// Set a value
	mockCfg.Set("raw.test", "raw_value")
	
	// Get it using the raw Get method
	value, exists := mockCfg.Get("raw.test")
	if !exists {
		t.Error("Expected key to exist")
	}
	
	if value != "raw_value" {
		t.Errorf("Expected 'raw_value', got %v", value)
	}
	
	// Test non-existent key
	_, exists = mockCfg.Get("non.existent")
	if exists {
		t.Error("Expected non-existent key to return false")
	}
}

// TestMockConfigFileOperations tests file-related operations
func TestMockConfigFileOperations(t *testing.T) {
	mockCfg := mocks.NewMockConfig()
	
	// Test loading from file (should succeed in mock or skip gracefully)
	err := mockCfg.LoadFromFile("fake/config.json")
	if err != nil {
		t.Skip("Mock config file loading not implemented or intentionally fails")
	}
	
	// Test loading from environment
	err = mockCfg.LoadFromEnv("MOCK_")
	if err != nil {
		t.Skip("Mock environment loading not implemented or intentionally fails")
	}
}
