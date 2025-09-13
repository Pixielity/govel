package pipeline_test

import (
	"testing"

	containerMocks "govel/container/mocks"
	"govel/new/pipeline/src"
	"govel/new/pipeline/src/interfaces"
)

// TestNewHub tests the NewHub constructor function
func TestNewHub(t *testing.T) {
	t.Run("NewHub with nil container", func(t *testing.T) {
		hub := pipeline.NewHub(nil)
		
		if hub == nil {
			t.Fatal("NewHub should not return nil")
		}
		
		// Test that we can define a default pipeline
		hub.Defaults(func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return passable
		})
		
		// Test that we can execute the pipeline
		result, err := hub.Pipe("test")
		if err != nil {
			t.Errorf("Hub execution failed: %v", err)
		}
		
		if result != "test" {
			t.Errorf("Expected 'test', got %v", result)
		}
	})

	t.Run("NewHub with mock container", func(t *testing.T) {
		container := containerMocks.NewMockContainer()
		hub := pipeline.NewHub(container)
		
		if hub == nil {
			t.Fatal("NewHub should not return nil")
		}
		
		// Test that the hub uses the provided container
		hub.Defaults(func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			// This tests that the pipeline created by the hub has access to the container
			result, _ := pipeline.Send(passable).ThenReturn()
			return result
		})
		
		result, err := hub.Pipe("test")
		if err != nil {
			t.Errorf("Hub execution failed: %v", err)
		}
		
		// The result should be a tuple (interface{}, error) from ThenReturn
		// But since we're in the callback, we need to handle it properly
		if resultVal, ok := result.(string); !ok || resultVal != "test" {
			t.Errorf("Expected 'test', got %v", result)
		}
	})
}

// TestNewHubInitialState tests that a new hub has the correct initial state
func TestNewHubInitialState(t *testing.T) {
	container := containerMocks.NewMockContainer()
	hub := pipeline.NewHub(container)
	
	// Test that calling Pipe without defining pipelines returns error
	result, err := hub.Pipe("test")
	
	if err == nil {
		t.Error("Expected error when calling Pipe without defining pipelines")
	}
	
	if result != nil {
		t.Errorf("Expected nil result when no pipelines defined, got %v", result)
	}
}

// TestNewHubThreadSafety tests basic thread safety of hub creation
func TestNewHubThreadSafety(t *testing.T) {
	container := containerMocks.NewMockContainer()
	
	// Create multiple hubs concurrently
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			hub := pipeline.NewHub(container)
			if hub == nil {
				t.Error("NewHub returned nil")
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestNewHubIsolation tests that different hub instances are isolated
func TestNewHubIsolation(t *testing.T) {
	container := containerMocks.NewMockContainer()
	
	hub1 := pipeline.NewHub(container)
	hub2 := pipeline.NewHub(container)
	
	// Define different pipelines on each hub
	hub1.Pipeline("test", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
		return "from-hub1"
	})
	
	hub2.Pipeline("test", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
		return "from-hub2"
	})
	
	// Test that each hub executes its own pipeline
	result1, err1 := hub1.Pipe("input", "test")
	result2, err2 := hub2.Pipe("input", "test")
	
	if err1 != nil || err2 != nil {
		t.Errorf("Hub execution failed: %v, %v", err1, err2)
	}
	
	if result1 != "from-hub1" {
		t.Errorf("Expected 'from-hub1', got %v", result1)
	}
	
	if result2 != "from-hub2" {
		t.Errorf("Expected 'from-hub2', got %v", result2)
	}
}
