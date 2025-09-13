package tests

import (
	"testing"

	"govel/packages/application"
)

// TestApplicationInfo tests comprehensive application information
func TestApplicationInfo(t *testing.T) {
	app := application.New()
	app.SetName("Test App")
	app.SetVersion("1.2.3")
	
	info := app.GetApplicationInfo()
	
	// Check basic info
	if info["name"] != "Test App" {
		t.Errorf("Expected name 'Test App', got %v", info["name"])
	}
	
	if info["version"] != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got %v", info["version"])
	}
	
	// Check that trait info is included
	if _, ok := info["environment"]; !ok {
		t.Error("Expected environment info to be included")
	}
	
	if _, ok := info["directories"]; !ok {
		t.Error("Expected directories info to be included")
	}
	
	if _, ok := info["locale"]; !ok {
		t.Error("Expected locale info to be included")
	}
}
