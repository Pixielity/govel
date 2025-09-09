package tests

import (
	"testing"

	"govel/packages/container"
)

// TestContainerBinding tests basic service binding
func TestContainerBinding(t *testing.T) {
	c := container.New()
	
	// Test simple string binding
	err := c.Bind("simple.service", "simple_value")
	if err != nil {
		t.Errorf("Expected no error when binding simple service, got %v", err)
	}
	
	// Test that service is bound
	if !c.IsBound("simple.service") {
		t.Error("Expected service to be bound")
	}
	
	// Test factory function binding
	err = c.Bind("factory.service", func() interface{} {
		return "factory_result"
	})
	
	if err != nil {
		t.Errorf("Expected no error when binding factory service, got %v", err)
	}
	
	if !c.IsBound("factory.service") {
		t.Error("Expected factory service to be bound")
	}
}
