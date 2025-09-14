package tests

import (
	"reflect"
	"sync"
	"testing"

	containerMocks "govel/container/mocks"
	"govel/new/pipeline"
)

// TestPipelineSend tests the Send method functionality
func TestPipelineSend(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Send with string data", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		testData := "hello world"
		
		result := pipeline.Send(testData)
		
		// Test method chaining - Send should return PipelineInterface
		if result == nil {
			t.Error("Send should return PipelineInterface for method chaining")
		}
		
		// Test that the data was stored by executing pipeline
		finalResult, err := result.ThenReturn()
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if finalResult != testData {
			t.Errorf("Expected %v, got %v", testData, finalResult)
		}
	})

	t.Run("Send with struct data", func(t *testing.T) {
		type TestStruct struct {
			ID   int
			Name string
		}
		
		pipeline := pipeline.NewPipeline(container)
		testData := TestStruct{ID: 1, Name: "test"}
		
		result := pipeline.Send(testData)
		finalResult, err := result.ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if finalResult != testData {
			t.Errorf("Expected %+v, got %+v", testData, finalResult)
		}
	})

	t.Run("Send with nil data", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		result := pipeline.Send(nil)
		finalResult, err := result.ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if finalResult != nil {
			t.Errorf("Expected nil, got %v", finalResult)
		}
	})

	t.Run("Send overwrites previous data", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		// Send first data
		pipeline.Send("first")
		
		// Send second data (should overwrite)
		result := pipeline.Send("second")
		finalResult, err := result.ThenReturn()
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if finalResult != "second" {
			t.Errorf("Expected 'second', got %v", finalResult)
		}
	})
}

// TestPipelineSendThreadSafety tests thread safety of the Send method
func TestPipelineSendThreadSafety(t *testing.T) {
	container := containerMocks.NewMockContainer()
	pipeline := pipeline.NewPipeline(container)
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	// Test concurrent Send operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			
			// Each goroutine sends different data
			testData := map[string]int{"id": id}
			result := pipeline.Send(testData)
			
			if result == nil {
				t.Error("Send should return PipelineInterface")
			}
		}(i)
	}
	
	wg.Wait()
	
	// After all goroutines complete, the pipeline should have the data
	// from one of the goroutines (the last one to complete)
	finalResult, err := pipeline.ThenReturn()
	if err != nil {
		t.Errorf("Pipeline execution failed: %v", err)
	}
	
	// Result should be one of the sent values
	if resultMap, ok := finalResult.(map[string]int); ok {
		if id, exists := resultMap["id"]; !exists || id < 0 || id >= numGoroutines {
			t.Errorf("Unexpected result: %v", finalResult)
		}
	} else {
		t.Errorf("Expected map[string]int, got %T", finalResult)
	}
}

// TestPipelineSendChaining tests method chaining after Send
func TestPipelineSendChaining(t *testing.T) {
	container := containerMocks.NewMockContainer()
	pipeline := pipeline.NewPipeline(container)
	
	// Test chaining Send -> Through -> Then
	result, err := pipeline.
		Send("test").
		Through([]interface{}{}).
		Then(func(passable interface{}) interface{} {
			return passable
		})
	
	if err != nil {
		t.Errorf("Chained pipeline execution failed: %v", err)
	}
	
	if result != "test" {
		t.Errorf("Expected 'test', got %v", result)
	}
}

// TestPipelineSendWithComplexData tests Send with complex data structures
func TestPipelineSendWithComplexData(t *testing.T) {
	container := containerMocks.NewMockContainer()
	pipeline := pipeline.NewPipeline(container)
	
	// Test with nested map
	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    123,
			"name":  "John Doe",
			"roles": []string{"admin", "user"},
		},
		"request": map[string]interface{}{
			"path":   "/api/users",
			"method": "GET",
		},
	}
	
	result, err := pipeline.Send(complexData).ThenReturn()
	
	if err != nil {
		t.Errorf("Pipeline execution failed: %v", err)
	}
	
	// Deep compare complex structures
	if !compareComplexData(result, complexData) {
		t.Errorf("Complex data not preserved through pipeline")
	}
}

// Helper function to compare complex data structures
func compareComplexData(a, b interface{}) bool {
	// Use reflect.DeepEqual for proper comparison of complex structures
	// This handles slices, maps, and other complex types correctly
	return reflect.DeepEqual(a, b)
}
