package tests

import (
	"testing"

	"govel/application"
)

// TestApplicationConfiguration tests configuration functionality via delegation
func TestApplicationConfiguration(t *testing.T) {
	app := application.New()
	
	// Test setting and getting string values
	app.Set("test.string", "hello world")
	value := app.GetString("test.string", "default")
	
	if value != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", value)
	}
	
	// Test setting and getting int values
	app.Set("test.number", 42)
	intValue := app.GetInt("test.number", 0)
	
	if intValue != 42 {
		t.Errorf("Expected 42, got %d", intValue)
	}
	
	// Test setting and getting bool values
	app.Set("test.flag", true)
	boolValue := app.GetBool("test.flag", false)
	
	if !boolValue {
		t.Error("Expected true, got false")
	}
	
	// Test HasKey
	if !app.HasKey("test.string") {
		t.Error("Expected HasKey to return true for existing key")
	}
	
	if app.HasKey("nonexistent.key") {
		t.Error("Expected HasKey to return false for non-existent key")
	}
	
	// Test AllConfig - check that we can retrieve what we set
	allConfig := app.AllConfig()
	t.Logf("AllConfig returned %d entries: %v", len(allConfig), allConfig)
	
	// At minimum, we should have at least 1 entry if the configuration works
	if len(allConfig) == 0 {
		t.Error("Expected at least some config entries")
	}
}
