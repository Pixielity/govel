package tests

import (
	"errors"
	"fmt"
	"testing"

	containerMocks "govel/packages/container/mocks"
	"govel/packages/new/pipeline/src"
)

// MockMiddleware simulates a middleware for testing
type MockMiddleware struct {
	name      string
	shouldErr bool
}

func (m *MockMiddleware) Handle(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	if m.shouldErr {
		return nil, errors.New("middleware error")
	}
	
	// Modify the passable data
	if str, ok := passable.(string); ok {
		passable = str + "->" + m.name
	}
	
	return next(passable)
}

// TestPipelineThrough tests the Through method functionality
func TestPipelineThrough(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Through with empty pipes array", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		result := pipeline.Send("test").Through([]interface{}{})
		
		// Should return pipeline interface for chaining
		if result == nil {
			t.Error("Through should return PipelineInterface for method chaining")
		}
		
		// Execute to verify empty pipes don't affect data
		finalResult, err := result.ThenReturn()
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if finalResult != "test" {
			t.Errorf("Expected 'test', got %v", finalResult)
		}
	})

	t.Run("Through with function pipes", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		// Define middleware functions
		middleware1 := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			if str, ok := passable.(string); ok {
				passable = str + "->func1"
			}
			return next(passable)
		}
		
		middleware2 := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			if str, ok := passable.(string); ok {
				passable = str + "->func2"
			}
			return next(passable)
		}
		
		pipes := []interface{}{middleware1, middleware2}
		result, err := pipeline.Send("test").Through(pipes).ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		// Pipes execute in order: test -> func1 -> func2
		expected := "test->func1->func2"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Through with object pipes", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		middleware1 := &MockMiddleware{name: "obj1", shouldErr: false}
		middleware2 := &MockMiddleware{name: "obj2", shouldErr: false}
		
		pipes := []interface{}{middleware1, middleware2}
		result, err := pipeline.Send("test").Through(pipes).ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		expected := "test->obj1->obj2"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Through with string pipes", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		// Register middleware in container
		container.Bind("middleware1", &MockMiddleware{name: "str1", shouldErr: false})
		container.Bind("middleware2", &MockMiddleware{name: "str2", shouldErr: false})
		
		pipes := []interface{}{"middleware1", "middleware2"}
		result, err := pipeline.Send("test").Through(pipes).ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		expected := "test->str1->str2"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Through overwrites previous pipes", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		middleware1 := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			if str, ok := passable.(string); ok {
				passable = str + "->first"
			}
			return next(passable)
		}
		
		middleware2 := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			if str, ok := passable.(string); ok {
				passable = str + "->second"
			}
			return next(passable)
		}
		
		// Set first pipes
		pipeline.Send("test").Through([]interface{}{middleware1})
		
		// Overwrite with second pipes
		result, err := pipeline.Through([]interface{}{middleware2}).ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		// Should only have second middleware
		expected := "test->second"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})
}

// TestPipelineThroughErrorHandling tests error handling in Through method
func TestPipelineThroughErrorHandling(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Through with middleware that errors", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		middleware1 := &MockMiddleware{name: "good", shouldErr: false}
		middleware2 := &MockMiddleware{name: "bad", shouldErr: true}
		middleware3 := &MockMiddleware{name: "never_reached", shouldErr: false}
		
		pipes := []interface{}{middleware1, middleware2, middleware3}
		result, err := pipeline.Send("test").Through(pipes).ThenReturn()
		
		// Should get error from middleware2
		if err == nil {
			t.Error("Expected error from middleware, got nil")
		}
		
		if !errors.Is(err, errors.New("middleware error")) && err.Error() != "middleware error" {
			t.Errorf("Expected middleware error, got: %v", err)
		}
		
		// Result should be nil when error occurs
		if result != nil {
			t.Errorf("Expected nil result on error, got: %v", result)
		}
	})

	t.Run("Through with unresolvable string pipe", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		// Don't bind "nonexistent" to container
		pipes := []interface{}{"nonexistent"}
		result, err := pipeline.Send("test").Through(pipes).ThenReturn()
		
		// Should handle unresolvable pipes gracefully
		if err == nil {
			t.Error("Expected error for unresolvable pipe")
		}
		
		if result != nil {
			t.Errorf("Expected nil result on error, got: %v", result)
		}
	})
}

// TestPipelineThroughThreadSafety tests thread safety of the Through method
func TestPipelineThroughThreadSafety(t *testing.T) {
	container := containerMocks.NewMockContainer()
	
	// Create multiple pipelines and run them concurrently
	numGoroutines := 50
	results := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			pipeline := pipeline.NewPipeline(container)
			
			middleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
				if str, ok := passable.(string); ok {
					passable = str + "->goroutine"
				}
				return next(passable)
			}
			
			result, err := pipeline.
				Send("test").
				Through([]interface{}{middleware}).
				ThenReturn()
			
			if err != nil {
				errors <- err
				return
			}
			
			if str, ok := result.(string); ok {
				results <- str
			} else {
				errors <- fmt.Errorf("invalid result type")
			}
		}(i)
	}
	
	// Collect results
	for i := 0; i < numGoroutines; i++ {
		select {
		case result := <-results:
			expected := "test->goroutine"
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
		case err := <-errors:
			t.Errorf("Goroutine failed: %v", err)
		}
	}
}

// TestPipelineThroughPipeOrdering tests that pipes execute in the correct order
func TestPipelineThroughPipeOrdering(t *testing.T) {
	container := containerMocks.NewMockContainer()
	pipeline := pipeline.NewPipeline(container)
	
	// Create middleware that tracks execution order
	var executionOrder []string
	
	createOrderingMiddleware := func(name string) func(interface{}, func(interface{}) (interface{}, error)) (interface{}, error) {
		return func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			executionOrder = append(executionOrder, name)
			return next(passable)
		}
	}
	
	pipes := []interface{}{
		createOrderingMiddleware("first"),
		createOrderingMiddleware("second"),
		createOrderingMiddleware("third"),
	}
	
	_, err := pipeline.Send("test").Through(pipes).ThenReturn()
	
	if err != nil {
		t.Errorf("Pipeline execution failed: %v", err)
	}
	
	// Verify execution order
	expectedOrder := []string{"first", "second", "third"}
	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d middleware executions, got %d", len(expectedOrder), len(executionOrder))
	}
	
	for i, expected := range expectedOrder {
		if i >= len(executionOrder) || executionOrder[i] != expected {
			t.Errorf("Expected middleware %s at position %d, got %v", expected, i, executionOrder)
		}
	}
}
