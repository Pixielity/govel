package tests

import (
	"testing"
	"time"

	"govel/application/builders"
)

// TestAppBuilderConfiguration tests all configuration methods
func TestAppBuilderConfiguration(t *testing.T) {
	testName := "Test Application"
	testVersion := "2.1.0"
	testEnv := "production"
	testLocale := "fr"
	testFallbackLocale := "en"
	testTimezone := "Europe/Paris"
	testBasePath := "/opt/myapp"
	testTimeout := 45 * time.Second
	
	app := builders.NewApp().
		WithName(testName).
		WithVersion(testVersion).
		WithEnvironment(testEnv).
		WithLocale(testLocale).
		WithFallbackLocale(testFallbackLocale).
		WithTimezone(testTimezone).
		WithBasePath(testBasePath).
		WithShutdownTimeout(testTimeout).
		WithDebug(false).
		Build()
	
	// Test that all configurations were applied
	if app.GetName() != testName {
		t.Errorf("Expected name %s, got %s", testName, app.GetName())
	}
	
	if app.GetVersion() != testVersion {
		t.Errorf("Expected version %s, got %s", testVersion, app.GetVersion())
	}
	
	if app.GetEnvironment() != testEnv {
		t.Errorf("Expected environment %s, got %s", testEnv, app.GetEnvironment())
	}
	
	if app.GetLocale() != testLocale {
		t.Errorf("Expected locale %s, got %s", testLocale, app.GetLocale())
	}
	
	if app.GetFallbackLocale() != testFallbackLocale {
		t.Errorf("Expected fallback locale %s, got %s", testFallbackLocale, app.GetFallbackLocale())
	}
	
	if app.GetTimezone() != testTimezone {
		t.Errorf("Expected timezone %s, got %s", testTimezone, app.GetTimezone())
	}
	
	if app.IsDebug() != false {
		t.Error("Expected debug to be false")
	}
	
	if app.GetShutdownTimeout() != testTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", testTimeout, app.GetShutdownTimeout())
	}
}

// TestAppBuilderRuntimeModes tests InConsole and InTesting methods
func TestAppBuilderRuntimeModes(t *testing.T) {
	// Test InConsole
	consoleApp := builders.NewApp().
		InConsole().
		Build()
	
	if !consoleApp.IsRunningInConsole() {
		t.Error("Expected application to be running in console mode")
	}
	
	// Test InTesting
	testApp := builders.NewApp().
		InTesting().
		Build()
	
	if !testApp.IsRunningUnitTests() {
		t.Error("Expected application to be running in testing mode")
	}
	
	// Test both together
	bothApp := builders.NewApp().
		InConsole().
		InTesting().
		Build()
	
	if !bothApp.IsRunningInConsole() {
		t.Error("Expected application to be running in console mode")
	}
	
	if !bothApp.IsRunningUnitTests() {
		t.Error("Expected application to be running in testing mode")
	}
}

// TestAppBuilderDebugAutomaticDisable tests that debug is automatically disabled in production
func TestAppBuilderDebugAutomaticDisable(t *testing.T) {
	// First set debug to true, then set environment to production
	app := builders.NewApp().
		WithDebug(true).
		WithEnvironment("production").
		Build()
	
	// Debug should be automatically disabled in production
	if app.IsDebug() {
		t.Error("Expected debug to be automatically disabled in production environment")
	}
}

// TestAppBuilderWithContainer tests container injection
func TestAppBuilderWithContainer(t *testing.T) {
	// This test might need adjustment based on actual container integration
	// For now, we'll just test that the method exists and doesn't panic
	app := builders.NewApp().
		WithContainer(nil). // Passing nil to test the method exists
		Build()
	
	if app == nil {
		t.Error("Expected Build() to succeed even with nil container")
	}
}
