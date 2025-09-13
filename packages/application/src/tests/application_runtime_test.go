package tests

import (
	"testing"

	"govel/application"
)

// TestApplicationRuntime tests application runtime state methods
func TestApplicationRuntime(t *testing.T) {
	app := application.New()
	
	// Test console mode
	app.SetRunningInConsole(true)
	if !app.IsRunningInConsole() {
		t.Error("Expected application to be running in console mode")
	}
	
	app.SetRunningInConsole(false)
	if app.IsRunningInConsole() {
		t.Error("Expected application to not be running in console mode")
	}
	
	// Test unit testing mode
	app.SetRunningUnitTests(true)
	if !app.IsRunningUnitTests() {
		t.Error("Expected application to be running unit tests")
	}
	
	app.SetRunningUnitTests(false)
	if app.IsRunningUnitTests() {
		t.Error("Expected application to not be running unit tests")
	}
}
