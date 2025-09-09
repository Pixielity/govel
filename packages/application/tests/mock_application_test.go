package tests

import (
	"testing"

	"govel/packages/application/mocks"
)

// TestMockApplication tests our MockApplication implementation
func TestMockApplication(t *testing.T) {
	mockApp := mocks.NewMockApplication()

	// Test basic properties
	if mockApp.GetName() != "MockApp" {
		t.Errorf("Expected mock name 'MockApp', got %s", mockApp.GetName())
	}

	if mockApp.GetVersion() != "1.0.0-mock" {
		t.Errorf("Expected mock version '1.0.0-mock', got %s", mockApp.GetVersion())
	}

	// Test trait capabilities
	if !mockApp.HasConfig() {
		t.Error("Expected mock to be configurable")
	}

	if !mockApp.HasContainer() {
		t.Error("Expected mock to be containable")
	}

	if !mockApp.HasLogger() {
		t.Error("Expected mock to be loggable")
	}

	// Test configuration through mock
	mockApp.Set("mock.test", "mock_value")
	value := mockApp.GetString("mock.test", "default")

	if value != "mock_value" {
		t.Errorf("Expected 'mock_value', got '%s'", value)
	}

	// Test container through mock
	err := mockApp.Bind("mock_service", "mock_implementation")
	if err != nil {
		t.Errorf("Expected no error binding to mock container, got %v", err)
	}

	service, err := mockApp.Make("mock_service")
	if err != nil {
		t.Errorf("Expected no error resolving from mock container, got %v", err)
	}

	if service != "mock_implementation" {
		t.Errorf("Expected 'mock_implementation', got %v", service)
	}

	// Test logger through mock
	mockApp.Info("Test log message")
	mockApp.Debug("Debug message with args: %s", "test")

	// Get the embedded mock logger to verify messages
	mockLogger := mockApp.GetMockLogger()
	if mockLogger == nil {
		t.Fatal("Expected to get mock logger instance")
	}

	messages := mockLogger.GetMessages()
	if len(messages) < 2 {
		t.Errorf("Expected at least 2 log messages, got %d", len(messages))
	}
}

// TestMockApplicationFailureMode tests failure simulation
func TestMockApplicationFailureMode(t *testing.T) {
	mockApp := mocks.NewMockApplication()

	// Enable failure mode for container operations
	mockApp.SetFailureMode(true, true, false, false)

	// Test that container operations fail
	mockContainer := mockApp.GetMockContainer()
	if mockContainer == nil {
		t.Fatal("Expected to get mock container instance")
	}

	mockContainer.SetFailureMode(true, true, false)

	err := mockApp.Bind("failing_service", "implementation")
	if err == nil {
		t.Error("Expected bind operation to fail in failure mode")
	}

	_, err = mockApp.Make("some_service")
	if err == nil {
		t.Error("Expected make operation to fail in failure mode")
	}
}
