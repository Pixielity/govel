package tests

import (
	"testing"

	"govel/application"
)

// TestApplicationCreation tests basic application creation and setup
func TestApplicationCreation(t *testing.T) {
	app := application.New()
	
	if app == nil {
		t.Fatal("Expected application to be created, got nil")
	}
	
	// Test default values
	if app.GetName() == "" {
		t.Error("Expected application to have a default name")
	}
	
	if app.GetVersion() == "" {
		t.Error("Expected application to have a default version")
	}
	
	// Test that start time is zero initially (not started yet)
	if !app.GetStartTime().IsZero() {
		t.Error("Expected start time to be zero initially")
	}
}
