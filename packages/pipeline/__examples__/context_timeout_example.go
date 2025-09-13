package main

import (
	"context"
	"fmt"
	"time"

	"govel/new/pipeline/src"
)

// Task represents a task to be processed
type Task struct {
	ID       string
	Type     string
	Data     interface{}
	Priority int
}

// SlowMiddleware simulates a slow middleware operation
func SlowMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	task := passable.(*Task)
	fmt.Printf("[SlowMiddleware] Processing task %s (sleeping for 2 seconds)\n", task.ID)
	
	// Simulate slow operation
	time.Sleep(2 * time.Second)
	
	task.Data = "processed by slow middleware"
	fmt.Printf("[SlowMiddleware] Completed processing task %s\n", task.ID)
	
	return next(passable)
}

// FastMiddleware simulates a fast middleware operation
func FastMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	task := passable.(*Task)
	fmt.Printf("[FastMiddleware] Processing task %s\n", task.ID)
	
	// Add metadata to track processing
	task.Priority += 1
	fmt.Printf("[FastMiddleware] Increased priority of task %s to %d\n", task.ID, task.Priority)
	
	return next(passable)
}

// ExampleContextWithTimeout demonstrates using pipeline context with timeout
func ExampleContextWithTimeout() {
	fmt.Println("=== Context Timeout Example ===")
	
	// Create pipeline
	p := pipeline.NewPipeline(nil)
	
	// Create base context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	// Create pipeline context
	pipelineCtx := pipeline.NewPipelineContext(ctx)
	
	// Add metadata
	pipelineCtx.SetMetadata("request_id", "REQ-001")
	pipelineCtx.SetMetadata("user_id", 123)
	pipelineCtx.SetMetadata("start_time", time.Now())
	
	task := &Task{
		ID:       "TASK-001",
		Type:     "data_processing",
		Data:     "initial data",
		Priority: 1,
	}
	
	fmt.Println("\n--- Testing with timeout (will fail) ---")
	
	// This will timeout because SlowMiddleware takes 2 seconds but timeout is 1 second
	result, err := p.
		Send(task).
		WithContext(pipelineCtx).
		Through([]interface{}{FastMiddleware, SlowMiddleware}).
		Then(func(passable interface{}) interface{} {
			task := passable.(*Task)
			fmt.Printf("[Final] Task %s completed successfully\n", task.ID)
			return task
		})
	
	if err != nil {
		fmt.Printf("Expected timeout error: %v\n", err)
	} else {
		fmt.Printf("Unexpected success: %+v\n", result)
	}
}

// ExampleContextWithoutTimeout demonstrates successful processing without timeout
func ExampleContextWithoutTimeout() {
	fmt.Println("\n=== Context Without Timeout Example ===")
	
	// Create pipeline
	p := pipeline.NewPipeline(nil)
	
	// Create base context with longer timeout (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create pipeline context
	pipelineCtx := pipeline.NewPipelineContext(ctx)
	
	// Add metadata
	pipelineCtx.SetMetadata("request_id", "REQ-002")
	pipelineCtx.SetMetadata("user_id", 456)
	pipelineCtx.SetMetadata("start_time", time.Now())
	
	task := &Task{
		ID:       "TASK-002",
		Type:     "data_processing",
		Data:     "initial data",
		Priority: 1,
	}
	
	fmt.Println("\n--- Testing with longer timeout (will succeed) ---")
	
	// This will succeed because timeout is 5 seconds
	result, err := p.
		Send(task).
		WithContext(pipelineCtx).
		Through([]interface{}{FastMiddleware, SlowMiddleware}).
		Then(func(passable interface{}) interface{} {
			task := passable.(*Task)
			
			// Access context metadata in the final handler
			if ctx, ok := passable.(*Task); ok {
				_ = ctx // Use the task, but we want to access pipeline context
			}
			
			fmt.Printf("[Final] Task %s completed successfully\n", task.ID)
			return task
		})
	
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		if task, ok := result.(*Task); ok {
			fmt.Printf("Success: Task %s processed (Priority: %d, Data: %s)\n", 
				task.ID, task.Priority, task.Data)
		}
	}
}

// ExampleContextMetadata demonstrates context metadata usage
func ExampleContextMetadata() {
	fmt.Println("\n=== Context Metadata Example ===")
	
	// Create pipeline
	p := pipeline.NewPipeline(nil)
	
	// Create pipeline context
	pipelineCtx := pipeline.NewPipelineContext(context.Background())
	
	// Add various types of metadata
	pipelineCtx.SetMetadata("request_id", "REQ-003")
	pipelineCtx.SetMetadata("user_id", 789)
	pipelineCtx.SetMetadata("permissions", []string{"read", "write"})
	pipelineCtx.SetMetadata("config", map[string]interface{}{
		"debug":   true,
		"timeout": 30,
	})
	
	// Middleware that uses context metadata
	metadataMiddleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
		task := passable.(*Task)
		
		fmt.Printf("[MetadataMiddleware] Processing task %s\n", task.ID)
		
		// Access metadata through context (we'd need to pass context to middleware)
		// For this example, we'll simulate accessing metadata
		
		// In a real implementation, middleware would have access to the context
		// through the pipeline's context propagation mechanism
		
		fmt.Printf("[MetadataMiddleware] Task metadata would be accessible here\n")
		
		return next(passable)
	}
	
	task := &Task{
		ID:       "TASK-003",
		Type:     "metadata_test",
		Data:     "metadata example",
		Priority: 5,
	}
	
	result, err := p.
		Send(task).
		WithContext(pipelineCtx).
		Through([]interface{}{metadataMiddleware}).
		Then(func(passable interface{}) interface{} {
			task := passable.(*Task)
			fmt.Printf("[Final] Task %s completed with metadata support\n", task.ID)
			return task
		})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %+v\n", result)
	}
	
	// Demonstrate metadata retrieval
	fmt.Println("\n--- Metadata Retrieval ---")
	
	if requestID, exists := pipelineCtx.GetMetadata("request_id"); exists {
		fmt.Printf("Request ID: %v\n", requestID)
	}
	
	if userID, exists := pipelineCtx.GetMetadata("user_id"); exists {
		fmt.Printf("User ID: %v\n", userID)
	}
	
	if permissions, exists := pipelineCtx.GetMetadata("permissions"); exists {
		fmt.Printf("Permissions: %v\n", permissions)
	}
	
	if config, exists := pipelineCtx.GetMetadata("config"); exists {
		fmt.Printf("Config: %v\n", config)
	}
	
	// Try to get non-existent metadata
	if _, exists := pipelineCtx.GetMetadata("nonexistent"); !exists {
		fmt.Println("Non-existent metadata correctly returns false")
	}
}

// ExampleContextCancellation demonstrates context cancellation
func ExampleContextCancellation() {
	fmt.Println("\n=== Context Cancellation Example ===")
	
	// Create pipeline
	p := pipeline.NewPipeline(nil)
	
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create pipeline context
	pipelineCtx := pipeline.NewPipelineContext(ctx)
	pipelineCtx.SetMetadata("request_id", "REQ-004")
	
	task := &Task{
		ID:       "TASK-004",
		Type:     "cancellation_test",
		Data:     "will be cancelled",
		Priority: 1,
	}
	
	// Start processing in a goroutine
	resultChan := make(chan interface{})
	errorChan := make(chan error)
	
	go func() {
		result, err := p.
			Send(task).
			WithContext(pipelineCtx).
			Through([]interface{}{SlowMiddleware}).
			Then(func(passable interface{}) interface{} {
				task := passable.(*Task)
				fmt.Printf("[Final] This should not be reached due to cancellation\n")
				return task
			})
		
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()
	
	// Cancel after 500ms (before SlowMiddleware completes)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Main] Cancelling context...")
	cancel()
	
	// Wait for result
	select {
	case result := <-resultChan:
		fmt.Printf("Unexpected success: %+v\n", result)
	case err := <-errorChan:
		fmt.Printf("Expected cancellation error: %v\n", err)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout waiting for result")
	}
}

func main() {
	ExampleContextWithTimeout()
	ExampleContextWithoutTimeout()
	ExampleContextMetadata()
	ExampleContextCancellation()
}
