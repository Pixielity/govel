package tests

import (
	"testing"
	"time"

	"govel/packages/container"
)

// TestContainerResolution tests service resolution
func TestContainerResolution(t *testing.T) {
	c := container.New()
	
	// Bind a simple value
	c.Bind("test.value", "hello world")
	
	// Resolve the value
	resolved, err := c.Make("test.value")
	if err != nil {
		t.Errorf("Expected no error when resolving service, got %v", err)
	}
	
	if resolved != "hello world" {
		t.Errorf("Expected 'hello world', got %v", resolved)
	}
	
	// Test factory resolution
	c.Bind("test.factory", func() interface{} {
		return map[string]interface{}{
			"created_at": time.Now().Unix(),
			"value":      42,
		}
	})
	
	resolved, err = c.Make("test.factory")
	if err != nil {
		t.Errorf("Expected no error when resolving factory, got %v", err)
	}
	
	if resolved == nil {
		t.Error("Expected resolved factory to not be nil")
	}
	
	// Verify the resolved value is a map
	resolvedMap, ok := resolved.(map[string]interface{})
	if !ok {
		t.Error("Expected resolved value to be a map")
	} else {
		if resolvedMap["value"] != 42 {
			t.Errorf("Expected value 42, got %v", resolvedMap["value"])
		}
	}
}
