package tests

import (
	"testing"

	"govel/application"
)

// TestApplicationContainer tests container functionality via delegation
func TestApplicationContainer(t *testing.T) {
	app := application.New()
	
	// Test binding a service
	err := app.Bind("test_service", func() interface{} {
		return "Hello from service"
	})
	
	if err != nil {
		t.Errorf("Expected no error when binding service, got %v", err)
	}
	
	// Test checking if bound
	if !app.IsBound("test_service") {
		t.Error("Expected service to be bound")
	}
	
	// Test resolving the service
	service, err := app.Make("test_service")
	if err != nil {
		t.Errorf("Expected no error when resolving service, got %v", err)
	}
	
	if service != "Hello from service" {
		t.Errorf("Expected 'Hello from service', got %v", service)
	}
	
	// Test singleton
	err = app.Singleton("singleton_service", func() interface{} {
		return &struct{ Value int }{Value: 123}
	})
	
	if err != nil {
		t.Errorf("Expected no error when binding singleton, got %v", err)
	}
	
	// Resolve singleton twice and ensure it's the same instance
	instance1, _ := app.Make("singleton_service")
	instance2, _ := app.Make("singleton_service")
	
	// They should be the same pointer for singletons
	if instance1 != instance2 {
		t.Error("Expected singleton instances to be the same")
	}
}
