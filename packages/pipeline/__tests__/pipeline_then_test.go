package tests

import (
	"errors"
	"testing"

	containerMocks "govel/packages/container/mocks"
	"govel/packages/new/pipeline/src"
)

// TestPipelineThen tests the Then method functionality
func TestPipelineThen(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Then with simple destination", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		destination := func(passable interface{}) interface{} {
			if str, ok := passable.(string); ok {
				return str + "->destination"
			}
			return passable
		}
		
		result, err := pipeline.Send("test").Then(destination)
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		expected := "test->destination"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Then with middleware and destination", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		middleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			if str, ok := passable.(string); ok {
				passable = str + "->middleware"
			}
			return next(passable)
		}
		
		destination := func(passable interface{}) interface{} {
			if str, ok := passable.(string); ok {
				return str + "->destination"
			}
			return passable
		}
		
		result, err := pipeline.
			Send("test").
			Through([]interface{}{middleware}).
			Then(destination)
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		expected := "test->middleware->destination"
		if result != expected {
			t.Errorf("Expected '%s', got %v", expected, result)
		}
	})

	t.Run("Then with complex data transformation", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		type User struct {
			Name string
			Age  int
		}
		
		destination := func(passable interface{}) interface{} {
			if user, ok := passable.(User); ok {
				return map[string]interface{}{
					"name":      user.Name,
					"age":       user.Age,
					"processed": true,
				}
			}
			return passable
		}
		
		user := User{Name: "John", Age: 30}
		result, err := pipeline.Send(user).Then(destination)
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if resultMap, ok := result.(map[string]interface{}); ok {
			if resultMap["name"] != "John" || resultMap["age"] != 30 || resultMap["processed"] != true {
				t.Errorf("Unexpected result: %+v", result)
			}
		} else {
			t.Errorf("Expected map[string]interface{}, got %T", result)
		}
	})
}

// TestPipelineThenErrorHandling tests error handling in Then method
func TestPipelineThenErrorHandling(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Then with panicking destination", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		destination := func(passable interface{}) interface{} {
			panic("destination panic")
		}
		
		result, err := pipeline.Send("test").Then(destination)
		
		// Should recover from panic and return error
		if err == nil {
			t.Error("Expected error from panicking destination")
		}
		
		if result != nil {
			t.Errorf("Expected nil result on panic, got: %v", result)
		}
	})

	t.Run("Then with middleware error", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		middleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
			return nil, errors.New("middleware error")
		}
		
		destination := func(passable interface{}) interface{} {
			return passable
		}
		
		result, err := pipeline.
			Send("test").
			Through([]interface{}{middleware}).
			Then(destination)
		
		if err == nil {
			t.Error("Expected error from middleware")
		}
		
		if result != nil {
			t.Errorf("Expected nil result on error, got: %v", result)
		}
	})
}

// TestPipelineThenWithFinally tests Then method with Finally callback
func TestPipelineThenWithFinally(t *testing.T) {
	container := containerMocks.NewMockContainer()

	t.Run("Then with Finally on success", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		finallyCalled := false
		var finallyData interface{}
		
		destination := func(passable interface{}) interface{} {
			return passable.(string) + "->destination"
		}
		
		finallyCallback := func(data interface{}) {
			finallyCalled = true
			finallyData = data
		}
		
		result, err := pipeline.
			Send("test").
			Finally(finallyCallback).
			Then(destination)
		
		if err != nil {
			t.Errorf("Pipeline execution failed: %v", err)
		}
		
		if !finallyCalled {
			t.Error("Finally callback was not called")
		}
		
		if finallyData != "test" {
			t.Errorf("Finally callback received wrong data: %v", finallyData)
		}
		
		if result != "test->destination" {
			t.Errorf("Expected 'test->destination', got %v", result)
		}
	})

	t.Run("Then with Finally on error", func(t *testing.T) {
		pipeline := pipeline.NewPipeline(container)
		
		finallyCalled := false
		
		destination := func(passable interface{}) interface{} {
			panic("destination error")
		}
		
		finallyCallback := func(data interface{}) {
			finallyCalled = true
		}
		
		result, err := pipeline.
			Send("test").
			Finally(finallyCallback).
			Then(destination)
		
		if err == nil {
			t.Error("Expected error from panicking destination")
		}
		
		if !finallyCalled {
			t.Error("Finally callback should be called even on error")
		}
		
		if result != nil {
			t.Errorf("Expected nil result on error, got: %v", result)
		}
	})
}
