package tests

import (
	"testing"

	"govel/config"
)

// TestConfigCreation tests basic config creation
func TestConfigCreation(t *testing.T) {
	cfg := config.New()
	
	if cfg == nil {
		t.Fatal("Expected config to be created, got nil")
	}
	
	// Test that new config is empty initially
	allConfig := cfg.AllConfig()
	if len(allConfig) != 0 {
		t.Errorf("Expected empty config initially, got %d items", len(allConfig))
	}
}

// TestConfigBasicOperations tests basic config operations
func TestConfigBasicOperations(t *testing.T) {
	cfg := config.New()
	
	// Test setting and getting string values
	cfg.Set("app.name", "Test Application")
	name := cfg.GetString("app.name", "Default App")
	
	if name != "Test Application" {
		t.Errorf("Expected 'Test Application', got '%s'", name)
	}
	
	// Test getting non-existent key with default
	nonExistent := cfg.GetString("non.existent", "default_value")
	if nonExistent != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", nonExistent)
	}
	
	// Test HasKey
	if !cfg.HasKey("app.name") {
		t.Error("Expected HasKey to return true for existing key")
	}
	
	if cfg.HasKey("non.existent") {
		t.Error("Expected HasKey to return false for non-existent key")
	}
}
