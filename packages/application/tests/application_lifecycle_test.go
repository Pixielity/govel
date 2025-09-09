package tests

import (
	"context"
	"testing"

	"govel/packages/application/mocks"
)

// TestApplicationLifecycle tests application lifecycle methods
func TestApplicationLifecycle(t *testing.T) {
	mockApp := mocks.NewMockApplication()
	ctx := context.Background()
	
	// Test boot process
	err := mockApp.Boot(ctx)
	if err != nil {
		t.Errorf("Expected no error during boot, got %v", err)
	}
	
	if !mockApp.IsBooted() {
		t.Error("Expected application to be booted")
	}
	
	// Test start process
	err = mockApp.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error during start, got %v", err)
	}
	
	if !mockApp.IsStarted() {
		t.Error("Expected application to be started")
	}
	
	// Test state
	state := mockApp.GetState()
	if state != "running" {
		t.Errorf("Expected state 'running', got %s", state)
	}
	
	if !mockApp.IsState("running") {
		t.Error("Expected IsState to return true for 'running'")
	}
}
