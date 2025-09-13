package tests

import (
	"errors"
	"testing"

	"govel/container"
)

// TestContainerErrorHandling tests error scenarios
func TestContainerErrorHandling(t *testing.T) {
	c := container.New()
	
	// Test resolving non-existent service
	_, err := c.Make("non.existent")
	if err == nil {
		t.Error("Expected error when resolving non-existent service")
	}
	
	// Test binding with invalid key
	err = c.Bind("", "empty_key")
	if err == nil {
		t.Error("Expected error when binding with empty key")
	}
	
	// Test factory that returns error
	c.Bind("error.factory", func() interface{} {
		return errors.New("factory error")
	})
	
	resolved, err := c.Make("error.factory")
	if err != nil {
		// If the container handles factory errors, this is expected
		if resolved != nil {
			t.Error("Expected resolved value to be nil when factory fails")
		}
	} else {
		// If the container doesn't handle factory errors, the error should be in resolved
		if resolvedErr, ok := resolved.(error); ok {
			if resolvedErr.Error() != "factory error" {
				t.Errorf("Expected 'factory error', got %s", resolvedErr.Error())
			}
		}
	}
}
