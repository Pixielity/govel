package tests

import (
	"testing"
	"time"

	"govel/application/builders"
)

// TestAppBuilderForProduction tests the ForProduction convenience method
func TestAppBuilderForProduction(t *testing.T) {
	app := builders.NewApp().
		ForProduction().
		Build()

	// Check that production settings are applied
	if app.GetEnvironment() != "production" {
		t.Errorf("Expected environment 'production', got '%s'", app.GetEnvironment())
	}

	if app.IsDebug() {
		t.Error("Expected debug to be false in production")
	}

	// Check that shutdown timeout is set (60 seconds by default for production)
	expectedTimeout := 60 * time.Second
	if app.GetShutdownTimeout() != expectedTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", expectedTimeout, app.GetShutdownTimeout())
	}
}

// TestAppBuilderForDevelopment tests the ForDevelopment convenience method
func TestAppBuilderForDevelopment(t *testing.T) {
	app := builders.NewApp().
		ForDevelopment().
		Build()

	// Check that development settings are applied
	if app.GetEnvironment() != "development" {
		t.Errorf("Expected environment 'development', got '%s'", app.GetEnvironment())
	}

	if !app.IsDebug() {
		t.Error("Expected debug to be true in development")
	}

	// Check that shutdown timeout is set (10 seconds by default for development)
	expectedTimeout := 10 * time.Second
	if app.GetShutdownTimeout() != expectedTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", expectedTimeout, app.GetShutdownTimeout())
	}
}

// TestAppBuilderForTesting tests the ForTesting convenience method
func TestAppBuilderForTesting(t *testing.T) {
	app := builders.NewApp().
		ForTesting().
		Build()

	// Check that testing settings are applied
	if app.GetEnvironment() != "testing" {
		t.Errorf("Expected environment 'testing', got '%s'", app.GetEnvironment())
	}

	if !app.IsDebug() {
		t.Error("Expected debug to be true in testing")
	}

	if !app.IsRunningUnitTests() {
		t.Error("Expected application to be running in unit test mode")
	}

	// Check that shutdown timeout is set (5 seconds by default for testing)
	expectedTimeout := 5 * time.Second
	if app.GetShutdownTimeout() != expectedTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", expectedTimeout, app.GetShutdownTimeout())
	}
}

// TestAppBuilderConvenienceOverride tests that convenience methods can be overridden
func TestAppBuilderConvenienceOverride(t *testing.T) {
	// Test that we can override production settings
	app := builders.NewApp().
		ForProduction().
		WithDebug(true).                        // Override debug setting
		WithShutdownTimeout(120 * time.Second). // Override timeout
		Build()

	if app.GetEnvironment() != "production" {
		t.Errorf("Expected environment 'production', got '%s'", app.GetEnvironment())
	}

	// Debug should be true because we overrode it
	if !app.IsDebug() {
		t.Error("Expected debug to be true after override")
	}

	// Timeout should be 120 seconds because we overrode it
	expectedTimeout := 120 * time.Second
	if app.GetShutdownTimeout() != expectedTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", expectedTimeout, app.GetShutdownTimeout())
	}
}

// TestAppBuilderMethodChaining tests complex method chaining
func TestAppBuilderMethodChaining(t *testing.T) {
	// Test a complex chain of method calls
	app := builders.NewApp().
		WithName("Complex App").
		WithVersion("3.2.1").
		ForDevelopment().
		WithLocale("es").
		WithFallbackLocale("en").
		WithTimezone("Europe/Madrid").
		InConsole().
		WithShutdownTimeout(30 * time.Second).
		Build()

	// Verify all settings
	if app.GetName() != "Complex App" {
		t.Errorf("Expected name 'Complex App', got '%s'", app.GetName())
	}

	if app.GetVersion() != "3.2.1" {
		t.Errorf("Expected version '3.2.1', got '%s'", app.GetVersion())
	}

	if app.GetEnvironment() != "development" {
		t.Errorf("Expected environment 'development', got '%s'", app.GetEnvironment())
	}

	if !app.IsDebug() {
		t.Error("Expected debug to be true (from ForDevelopment)")
	}

	if app.GetLocale() != "es" {
		t.Errorf("Expected locale 'es', got '%s'", app.GetLocale())
	}

	if app.GetFallbackLocale() != "en" {
		t.Errorf("Expected fallback locale 'en', got '%s'", app.GetFallbackLocale())
	}

	if app.GetTimezone() != "Europe/Madrid" {
		t.Errorf("Expected timezone 'Europe/Madrid', got '%s'", app.GetTimezone())
	}

	if !app.IsRunningInConsole() {
		t.Error("Expected application to be running in console mode")
	}

	expectedTimeout := 30 * time.Second
	if app.GetShutdownTimeout() != expectedTimeout {
		t.Errorf("Expected shutdown timeout %v, got %v", expectedTimeout, app.GetShutdownTimeout())
	}
}
