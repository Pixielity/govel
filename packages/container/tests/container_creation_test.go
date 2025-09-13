package tests

import (
	"testing"

	"govel/packages/container"
)

// TestContainerCreation tests basic container creation
func TestContainerCreation(t *testing.T) {
	c := container.New()
	
	if c == nil {
		t.Fatal("Expected container to be created, got nil")
	}
	
	// Test that new container is empty initially
	if c.IsBound("non.existent") {
		t.Error("Expected new container to not have any bindings")
	}
}
