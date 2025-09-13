package pipeline_test

import (
	"testing"

	containerMocks "govel/container/mocks"
	"govel/new/pipeline/src"
)

// TestNewPipeline tests the NewPipeline constructor function
func TestNewPipeline(t *testing.T) {
	t.Run("NewPipeline with nil container", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(nil)
		
		if pipeline == nil {
			t.Fatal("NewPipeline should not return nil")
		}
		
		// Test that pipeline has default method "Handle"
		// We can't access private fields directly, but we can test behavior
		// This is a structural test to ensure the constructor works
	})

	t.Run("NewPipeline with mock container", func(t *testing.T) {
		container := containerMocks.NewMockContainer()
		pipeline := pipeline.NewPipeline(container)
		
		if pipeline == nil {
			t.Fatal("NewPipeline should not return nil")
		}
		
		// Test that the pipeline is properly initialized
		// We'll test this indirectly through method chaining
		result := pipeline.Send("test data")
		
		if result == nil {
			t.Error("Pipeline should support method chaining")
		}
	})
}

// TestNewPipelineInitialState tests that a new pipeline has the correct initial state
func TestNewPipelineInitialState(t *testing.T) {
	container := containerMocks.NewMockContainer()
	pipeline := pipeline.NewPipeline(container)
	
	// Test that we can immediately call methods without errors
	// This indirectly tests the initialization
	result, err := pipeline.ThenReturn()
	
	// Should return nil since no data was sent
	if result != nil {
		t.Errorf("Expected nil result for empty pipeline, got %v", result)
	}
	
	if err != nil {
		t.Errorf("Expected no error for empty pipeline, got %v", err)
	}
}

// TestNewPipelineThreadSafety tests basic thread safety of pipeline creation
func TestNewPipelineThreadSafety(t *testing.T) {
	container := containerMocks.NewMockContainer()
	
	// Create multiple pipelines concurrently
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			pipeline := pipeline.NewPipeline(container)
			if pipeline == nil {
				t.Error("NewPipeline returned nil")
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
