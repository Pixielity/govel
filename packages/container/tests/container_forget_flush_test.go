package tests

import (
	"testing"

	"govel/container"
)

// TestContainerForget tests the Forget functionality
func TestContainerForget(t *testing.T) {
	c := container.New()
	
	// Bind a service
	err := c.Bind("test.service", "test_value")
	if err != nil {
		t.Errorf("Expected no error when binding service, got %v", err)
	}
	
	// Verify it's bound
	if !c.IsBound("test.service") {
		t.Error("Expected service to be bound")
	}
	
	// Forget the service
	c.Forget("test.service")
	
	// Verify it's no longer bound
	if c.IsBound("test.service") {
		t.Error("Expected service to be forgotten")
	}
	
	// Trying to resolve should fail
	_, err = c.Make("test.service")
	if err == nil {
		t.Error("Expected error when resolving forgotten service")
	}
}

// TestContainerFlush tests the FlushContainer functionality
func TestContainerFlush(t *testing.T) {
	c := container.New()
	
	// Bind multiple services
	c.Bind("service1", "value1")
	c.Bind("service2", "value2")
	c.Singleton("singleton1", func() interface{} { return "singleton_value" })
	
	// Verify they're all bound
	if !c.IsBound("service1") || !c.IsBound("service2") || !c.IsBound("singleton1") {
		t.Error("Expected all services to be bound")
	}
	
	// Flush the container
	c.FlushContainer()
	
	// Verify all services are gone
	if c.IsBound("service1") || c.IsBound("service2") || c.IsBound("singleton1") {
		t.Error("Expected all services to be flushed")
	}
	
	// Trying to resolve any should fail
	_, err := c.Make("service1")
	if err == nil {
		t.Error("Expected error when resolving from flushed container")
	}
}
