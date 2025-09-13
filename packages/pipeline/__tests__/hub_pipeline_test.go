package pipeline_test

import (
	"fmt"
	"sync"
	"testing"

	containerMocks "govel/packages/container/mocks"
	"govel/packages/new/pipeline/src"
	"govel/packages/new/pipeline/src/interfaces"
)

// TestHubPipeline tests the Pipeline method functionality
func TestHubPipeline(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Pipeline with valid name and callback", func(t *testing.T) {
		hub := pipeline.NewHub(container)
		
		// Define a named pipeline
		hub.Pipeline("api", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "api-" + passable.(string)
		})
		
		// Execute the pipeline
		result, err := hub.Pipe("test", "api")
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		expected := "api-test"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Pipeline overwrites existing pipeline with same name", func(t *testing.T) {
		hub := pipeline.NewHub(container)
		
		// Define first pipeline
		hub.Pipeline("test", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "first"
		})
		
		// Overwrite with second pipeline
		hub.Pipeline("test", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "second"
		})
		
		// Should execute the second pipeline
		result, err := hub.Pipe("input", "test")
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if result != "second" {
			t.Errorf("Expected 'second', got %v", result)
		}
	})

	t.Run("Pipeline with multiple named pipelines", func(t *testing.T) {
		hub := pipeline.NewHub(container)
		
		// Define multiple pipelines
		hub.Pipeline("api", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "api-result"
		})
		
		hub.Pipeline("web", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "web-result"
		})
		
		hub.Pipeline("console", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return "console-result"
		})
		
		// Test each pipeline
		apiResult, err1 := hub.Pipe("input", "api")
		webResult, err2 := hub.Pipe("input", "web")
		consoleResult, err3 := hub.Pipe("input", "console")
		
		if err1 != nil || err2 != nil || err3 != nil {
			t.Errorf("Pipeline execution failed: %v, %v, %v", err1, err2, err3)
		}
		
		if apiResult != "api-result" {
			t.Errorf("Expected 'api-result', got %v", apiResult)
		}
		
		if webResult != "web-result" {
			t.Errorf("Expected 'web-result', got %v", webResult)
		}
		
		if consoleResult != "console-result" {
			t.Errorf("Expected 'console-result', got %v", consoleResult)
		}
	})
}

// TestHubPipelinePanic tests panic conditions in Pipeline method
func TestHubPipelinePanic(t *testing.T) {
	container := containerMocks.NewMockContainer()
	hub := pipeline.NewHub(container)

	t.Run("Pipeline panics with empty name", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty pipeline name")
			}
		}()
		
		hub.Pipeline("", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
			return passable
		})
	})

	t.Run("Pipeline panics with nil callback", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil callback")
			}
		}()
		
		hub.Pipeline("test", nil)
	})
}

// TestHubPipelineThreadSafety tests thread safety of the Pipeline method
func TestHubPipelineThreadSafety(t *testing.T) {
	container := containerMocks.NewMockContainer()
	hub := pipeline.NewHub(container)
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	// Register pipelines concurrently
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			
			pipelineName := fmt.Sprintf("pipeline-%d", id)
			expectedResult := fmt.Sprintf("result-%d", id)
			
			hub.Pipeline(pipelineName, func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
				return expectedResult
			})
		}(i)
	}
	
	wg.Wait()
	
	// Test that all pipelines were registered correctly
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			
			pipelineName := fmt.Sprintf("pipeline-%d", id)
			expectedResult := fmt.Sprintf("result-%d", id)
			
			result, err := hub.Pipe("test", pipelineName)
			
			if err != nil {
				t.Errorf("Pipeline execution failed for %s: %v", pipelineName, err)
				return
			}
			
			if result != expectedResult {
				t.Errorf("Expected '%s', got %v for pipeline %s", expectedResult, result, pipelineName)
			}
		}(i)
	}
	
	wg.Wait()
}

// TestHubPipelineWithComplexCallback tests Pipeline method with complex callbacks
func TestHubPipelineWithComplexCallback(t *testing.T) {
	container := containerMocks.NewMockContainer()
	hub := pipeline.NewHub(container)
	
	// Define middleware for the pipeline callback
	authMiddleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
		// Simulate authentication
		if request, ok := passable.(map[string]interface{}); ok {
			request["authenticated"] = true
		}
		return next(passable)
	}
	
	loggerMiddleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
		// Simulate logging
		if request, ok := passable.(map[string]interface{}); ok {
			if logs, ok := request["logs"].([]string); ok {
				request["logs"] = append(logs, "logged")
			} else {
				request["logs"] = []string{"logged"}
			}
		}
		return next(passable)
	}
	
	// Register a complex pipeline
	hub.Pipeline("complex", func(pipelineInstance interfaces.PipelineInterface, passable interface{}) interface{} {
		result, err := pipelineInstance.
			Send(passable).
			Through([]interface{}{authMiddleware, loggerMiddleware}).
			Then(func(passable interface{}) interface{} {
				if request, ok := passable.(map[string]interface{}); ok {
					request["processed"] = true
				}
				return passable
			})
		
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}
		
		return result
	})
	
	// Test the complex pipeline
	input := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}
	
	result, err := hub.Pipe(input, "complex")
	
	if err != nil {
		t.Errorf("Complex pipeline execution failed: %v", err)
	}
	
	if resultMap, ok := result.(map[string]interface{}); ok {
		if !resultMap["authenticated"].(bool) {
			t.Error("Request should be authenticated")
		}
		
		if !resultMap["processed"].(bool) {
			t.Error("Request should be processed")
		}
		
		if logs, ok := resultMap["logs"].([]string); !ok || len(logs) != 1 || logs[0] != "logged" {
			t.Errorf("Expected logs to contain 'logged', got %v", logs)
		}
	} else {
		t.Errorf("Expected map[string]interface{}, got %T", result)
	}
}
