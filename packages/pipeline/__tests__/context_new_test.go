package tests

import (
	"context"
	"testing"
	"time"

	"govel/packages/new/pipeline/src"
)

// TestNewPipelineContext tests the NewPipelineContext constructor function
func TestNewPipelineContext(t *testing.T) {
	t.Run("NewPipelineContext with Background context", func(t *testing.T) {
		ctx := context.Background()
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		if pipelineCtx == nil {
			t.Fatal("NewPipelineContext should not return nil")
		}
		
		// Test that it implements context.Context interface
		if pipelineCtx.Err() != nil {
			t.Error("Background context should not have error")
		}
		
		select {
		case <-pipelineCtx.Done():
			t.Error("Background context should not be done")
		default:
			// Expected: context is not done
		}
	})

	t.Run("NewPipelineContext with nil context", func(t *testing.T) {
		pipelineCtx := pipeline.NewPipelineContext(nil)
		
		if pipelineCtx == nil {
			t.Fatal("NewPipelineContext should not return nil even with nil input")
		}
		
		// Should default to background context
		if pipelineCtx.Err() != nil {
			t.Error("Default context should not have error")
		}
	})

	t.Run("NewPipelineContext with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately
		
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		if pipelineCtx == nil {
			t.Fatal("NewPipelineContext should not return nil")
		}
		
		// Should inherit cancellation
		if pipelineCtx.Err() == nil {
			t.Error("Pipeline context should inherit cancellation")
		}
		
		select {
		case <-pipelineCtx.Done():
			// Expected: context is done
		default:
			t.Error("Pipeline context should be done")
		}
	})

	t.Run("NewPipelineContext with timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		if pipelineCtx == nil {
			t.Fatal("NewPipelineContext should not return nil")
		}
		
		// Should inherit deadline
		deadline, ok := pipelineCtx.Deadline()
		if !ok {
			t.Error("Pipeline context should have deadline")
		}
		
		// Deadline should be in the future but soon
		if time.Until(deadline) > 100*time.Millisecond {
			t.Error("Deadline should be within timeout period")
		}
	})
}

// TestNewPipelineContextInitialState tests the initial state of a new pipeline context
func TestNewPipelineContextInitialState(t *testing.T) {
	pipelineCtx := pipeline.NewPipelineContext(context.Background())
	
	// Test that metadata starts empty
	_, exists := pipelineCtx.GetMetadata("nonexistent")
	if exists {
		t.Error("New pipeline context should have empty metadata")
	}
	
	// Test that we can set metadata immediately
	pipelineCtx.SetMetadata("test", "value")
	
	value, exists := pipelineCtx.GetMetadata("test")
	if !exists || value != "value" {
		t.Error("Should be able to set and get metadata immediately")
	}
}

// TestNewPipelineContextValues tests context value inheritance
func TestNewPipelineContextValues(t *testing.T) {
	// Create context with values
	parentCtx := context.WithValue(context.Background(), "key1", "parent-value")
	parentCtx = context.WithValue(parentCtx, "key2", 123)
	
	pipelineCtx := pipeline.NewPipelineContext(parentCtx)
	
	// Test that parent values are accessible
	if value := pipelineCtx.Value("key1"); value != "parent-value" {
		t.Errorf("Expected 'parent-value', got %v", value)
	}
	
	if value := pipelineCtx.Value("key2"); value != 123 {
		t.Errorf("Expected 123, got %v", value)
	}
	
	// Test that metadata values take precedence over parent context values
	pipelineCtx.SetMetadata("key1", "metadata-value")
	
	if value := pipelineCtx.Value("key1"); value != "metadata-value" {
		t.Errorf("Expected 'metadata-value', got %v", value)
	}
	
	// Parent value should still be accessible for key2
	if value := pipelineCtx.Value("key2"); value != 123 {
		t.Errorf("Expected 123, got %v", value)
	}
}

// TestNewPipelineContextThreadSafety tests thread safety of context creation
func TestNewPipelineContextThreadSafety(t *testing.T) {
	parentCtx := context.Background()
	
	// Create multiple pipeline contexts concurrently
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			pipelineCtx := pipeline.NewPipelineContext(parentCtx)
			
			if pipelineCtx == nil {
				t.Error("NewPipelineContext returned nil")
			}
			
			// Test basic functionality
			pipelineCtx.SetMetadata("id", id)
			if value, exists := pipelineCtx.GetMetadata("id"); !exists || value != id {
				t.Errorf("Metadata not set correctly: expected %d, got %v", id, value)
			}
			
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestNewPipelineContextWithDifferentContextTypes tests with various context types
func TestNewPipelineContextWithDifferentContextTypes(t *testing.T) {
	t.Run("With value context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "test", "value")
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		if pipelineCtx.Value("test") != "value" {
			t.Error("Value context not properly inherited")
		}
	})

	t.Run("With deadline context", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()
		
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		ctxDeadline, ok := pipelineCtx.Deadline()
		if !ok || !ctxDeadline.Equal(deadline) {
			t.Error("Deadline context not properly inherited")
		}
	})

	t.Run("With timeout context", func(t *testing.T) {
		timeout := time.Hour
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		
		pipelineCtx := pipeline.NewPipelineContext(ctx)
		
		_, ok := pipelineCtx.Deadline()
		if !ok {
			t.Error("Timeout context should have deadline")
		}
	})
}
