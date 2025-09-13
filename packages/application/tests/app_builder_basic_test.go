package tests

import (
	"testing"

	"govel/application/builders"
)

// TestAppBuilderCreation tests basic AppBuilder creation
func TestAppBuilderCreation(t *testing.T) {
	builder := builders.NewApp()
	
	if builder == nil {
		t.Fatal("Expected AppBuilder to be created, got nil")
	}
	
	// Test that we can build an application
	app := builder.Build()
	if app == nil {
		t.Fatal("Expected Build() to return an application, got nil")
	}
	
	// Test that the built application has basic properties
	if app.GetName() == "" {
		t.Error("Expected built application to have a name")
	}
	
	if app.GetVersion() == "" {
		t.Error("Expected built application to have a version")
	}
}

// TestAppBuilderFluentInterface tests that all builder methods return the builder for chaining
func TestAppBuilderFluentInterface(t *testing.T) {
	builder := builders.NewApp()
	
	// Test that each method returns the builder instance for chaining
	result := builder.WithName("Test App")
	if result != builder {
		t.Error("Expected WithName to return the same builder instance")
	}
	
	result = builder.WithVersion("1.0.0")
	if result != builder {
		t.Error("Expected WithVersion to return the same builder instance")
	}
	
	result = builder.WithEnvironment("testing")
	if result != builder {
		t.Error("Expected WithEnvironment to return the same builder instance")
	}
	
	result = builder.WithDebug(true)
	if result != builder {
		t.Error("Expected WithDebug to return the same builder instance")
	}
}
