package tests

import (
	"testing"

	"govel/container"
)

// TestContainerSingleton tests singleton service registration and resolution
func TestContainerSingleton(t *testing.T) {
	c := container.New()
	
	// Register a singleton
	err := c.Singleton("singleton.service", func() interface{} {
		return &struct {
			ID    string
			Value int
		}{
			ID:    "unique_id_123",
			Value: 999,
		}
	})
	
	if err != nil {
		t.Errorf("Expected no error when registering singleton, got %v", err)
	}
	
	// Resolve singleton multiple times
	instance1, err := c.Make("singleton.service")
	if err != nil {
		t.Errorf("Expected no error when resolving singleton, got %v", err)
	}
	
	instance2, err := c.Make("singleton.service")
	if err != nil {
		t.Errorf("Expected no error when resolving singleton second time, got %v", err)
	}
	
	// They should be the exact same instance (same pointer)
	if instance1 != instance2 {
		t.Error("Expected singleton instances to be the same")
	}
	
	// Verify the content is correct
	struct1, ok := instance1.(*struct {
		ID    string
		Value int
	})
	if !ok {
		t.Error("Expected instance to be the correct struct type")
	} else {
		if struct1.ID != "unique_id_123" {
			t.Errorf("Expected ID 'unique_id_123', got %s", struct1.ID)
		}
		if struct1.Value != 999 {
			t.Errorf("Expected Value 999, got %d", struct1.Value)
		}
	}
}
