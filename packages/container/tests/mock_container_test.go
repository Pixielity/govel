package tests

import (
	"testing"

	"govel/packages/container/mocks"
)

// TestMockContainer tests the mock container implementation
func TestMockContainer(t *testing.T) {
	mockContainer := mocks.NewMockContainer()
	
	// Test basic operations
	err := mockContainer.Bind("mock.service", "mock_value")
	if err != nil {
		t.Errorf("Expected no error binding to mock container, got %v", err)
	}
	
	// Test that service is bound
	if !mockContainer.IsBound("mock.service") {
		t.Error("Expected mock service to be bound")
	}
	
	// Test resolution
	resolved, err := mockContainer.Make("mock.service")
	if err != nil {
		t.Errorf("Expected no error resolving from mock container, got %v", err)
	}
	
	if resolved != "mock_value" {
		t.Errorf("Expected 'mock_value', got %v", resolved)
	}
	
	// Test singleton
	err = mockContainer.Singleton("mock.singleton", func() interface{} {
		return "singleton_value"
	})
	
	if err != nil {
		t.Errorf("Expected no error registering mock singleton, got %v", err)
	}
	
	// Resolve singleton twice
	instance1, _ := mockContainer.Make("mock.singleton")
	instance2, _ := mockContainer.Make("mock.singleton")
	
	if instance1 != instance2 {
		t.Error("Expected mock singleton instances to be the same")
	}
}

// TestMockContainerFailureMode tests failure simulation in mock container
func TestMockContainerFailureMode(t *testing.T) {
	mockContainer := mocks.NewMockContainer()
	
	// Enable failure modes
	mockContainer.SetFailureMode(true, true, false) // bind fails, make fails, forget succeeds
	
	// Test that bind fails
	err := mockContainer.Bind("failing.service", "value")
	if err == nil {
		t.Error("Expected bind to fail in failure mode")
	}
	
	// Test that make fails
	_, err = mockContainer.Make("some.service")
	if err == nil {
		t.Error("Expected make to fail in failure mode")
	}
	
	// Test that forget still works (not in failure mode)
	mockContainer.Forget("any.service") // Should not panic
}
