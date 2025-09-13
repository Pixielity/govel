// Package pipeline provides a Laravel-compatible pipeline hub implementation for Go.
package pipeline

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	containerInterfaces "govel/types/src/interfaces/container"

	interfaces "govel/types/src/interfaces/pipeline"
	"govel/types/src/types/pipeline"
)

// ErrPipelineNotFound is returned when trying to execute a pipeline that doesn't exist
var ErrPipelineNotFound = errors.New("pipeline not found")

// Hub implements the HubInterface interface and manages named pipelines.
// It provides a centralized way to define and execute different pipeline
// configurations under specific names.
//
// The hub maintains a registry of named pipeline configurations and
// provides methods to define, execute, and manage these pipelines.
// It supports a default pipeline for cases where no specific name is provided.
//
// This implementation is thread-safe and can be used across multiple goroutines.
type Hub struct {
	// container is the dependency injection container used for pipeline creation
	container containerInterfaces.ContainerInterface

	// pipelines stores the registered pipeline configurations by name
	pipelines map[string]types.PipelineCallback

	// mutex protects concurrent access to the pipelines map
	mutex sync.RWMutex
}

// NewHub creates a new Hub instance with the given container.
// The container is optional and can be nil if dependency injection is not needed.
//
// Parameters:
//   - container: Optional dependency injection container for pipeline creation
//
// Returns:
//   - *Hub: New hub instance ready for pipeline registration and execution
//
// Example:
//
//	hub := NewHub(container)
//
//	// Define a default pipeline
//	hub.Defaults(func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
//		result, _ := pipeline.Send(passable).Through(defaultMiddlewares).ThenReturn()
//		return result
//	})
func NewHub(container containerInterfaces.ContainerInterface) *Hub {
	return &Hub{
		container: container,
		pipelines: make(map[string]types.PipelineCallback),
		mutex:     sync.RWMutex{},
	}
}

// Defaults defines the default pipeline configuration.
// This pipeline will be used when no specific pipeline name is provided
// to the Pipe method.
//
// Parameters:
//   - callback: Function that defines the default pipeline configuration
//     The callback receives a PipelineInterface instance and the passable object
//
// Thread-safe: This method is safe for concurrent use.
//
// Example:
//
//	hub.Defaults(func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
//		// Configure the pipeline with default middleware
//		result, err := pipeline.
//			Send(passable).
//			Through([]interface{}{AuthMiddleware{}, LoggingMiddleware{}}).
//			Then(func(p interface{}) interface{} {
//				// Process the request
//				return processRequest(p)
//			})
//
//		if err != nil {
//			// Handle error or return error response
//			return ErrorResponse{Error: err}
//		}
//
//		return result
//	})
func (h *Hub) Defaults(callback func(interfaces.PipelineInterface, interface{}) interface{}) {
	h.Pipeline("default", callback)
}

// Pipeline defines a named pipeline configuration.
// This allows for different pipeline setups to be registered under
// specific names for later execution.
//
// Parameters:
//   - name: Unique name for the pipeline configuration
//   - callback: Function that defines the pipeline configuration
//     The callback receives a PipelineInterface instance and the passable object
//
// Thread-safe: This method is safe for concurrent use.
//
// Example:
//
//	// Define an API pipeline with API-specific middleware
//	hub.Pipeline("api", func(pipeline interfaces.PipelineInterface, passable interface{}) interface{} {
//		result, err := pipeline.
//			Send(passable).
//			Through([]interface{}{"RateLimitMiddleware", "ApiAuthMiddleware"}).
//			Via("Handle").
//			Then(func(p interface{}) interface{} {
//				return processApiRequest(p)
//			})
//
//		if err != nil {
//			return ApiErrorResponse{Error: err}
//		}
//
//		return result
//	})
func (h *Hub) Pipeline(name string, callback func(interfaces.PipelineInterface, interface{}) interface{}) {
	if name == "" {
		panic("pipeline name cannot be empty")
	}

	if callback == nil {
		panic("pipeline callback cannot be nil")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.pipelines[name] = callback
}

// Pipe sends an object through one of the available pipelines.
// If no pipeline name is specified, the default pipeline will be used.
//
// Parameters:
//   - object: The object to pass through the pipeline
//   - pipelineName: Optional name of the pipeline to use (uses default if empty)
//
// Returns:
//   - interface{}: The result of pipeline execution
//   - error: Any error that occurred during pipeline execution
//
// Thread-safe: This method is safe for concurrent use.
//
// Example:
//
//	// Execute with default pipeline
//	result, err := hub.Pipe(request)
//	if err != nil {
//		log.Printf("Pipeline execution failed: %v", err)
//		return nil, err
//	}
//
//	// Execute with named pipeline
//	result, err = hub.Pipe(apiRequest, "api")
//	if err != nil {
//		log.Printf("API pipeline execution failed: %v", err)
//		return nil, err
//	}
func (h *Hub) Pipe(object interface{}, pipelineName ...string) (interface{}, error) {
	// Determine which pipeline to use
	name := "default"
	if len(pipelineName) > 0 && pipelineName[0] != "" {
		name = pipelineName[0]
	}

	// Get the pipeline callback
	h.mutex.RLock()
	callback, exists := h.pipelines[name]
	h.mutex.RUnlock()

	if !exists {
		return nil, &PipelineExecutionError{
			PipelineName: name,
			Cause:        ErrPipelineNotFound,
			Message:      fmt.Sprintf("pipeline '%s' is not registered", name),
		}
	}

	// Create a new pipeline instance
	pipeline := NewPipeline(h.container)

	// Execute the pipeline callback with panic recovery
	var result interface{}
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = &PipelineExecutionError{
					PipelineName: name,
					Cause:        fmt.Errorf("panic: %v", r),
					Message:      fmt.Sprintf("pipeline '%s' panicked during execution", name),
				}
			}
		}()

		result = callback(pipeline, object)
	}()

	return result, err
}

// GetContainer returns the container instance used by the hub.
// This provides access to the dependency injection container for
// pipe resolution and other container services.
//
// Returns:
//   - interfaces.Container: The container instance (may be nil)
func (h *Hub) GetContainer() containerInterfaces.ContainerInterface {
	return h.container
}

// SetContainer sets the container instance used by the hub.
// This allows for dependency injection container replacement
// or initialization after hub creation.
//
// Parameters:
//   - container: The container instance to use
//
// Returns:
//   - interfaces.HubInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (h *Hub) SetContainer(container containerInterfaces.ContainerInterface) interfaces.HubInterface {
	h.container = container
	return h
}

// HasPipeline checks if a named pipeline exists in the hub.
// This is useful for conditional pipeline execution or validation.
//
// Parameters:
//   - name: Name of the pipeline to check
//
// Returns:
//   - bool: True if the pipeline exists, false otherwise
//
// Thread-safe: This method is safe for concurrent use.
func (h *Hub) HasPipeline(name string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.pipelines[name]
	return exists
}

// GetPipelineNames returns all registered pipeline names.
// This is useful for debugging or administrative purposes.
//
// Returns:
//   - []string: Slice of all registered pipeline names (sorted alphabetically)
//
// Thread-safe: This method returns a copy of the pipeline names to prevent race conditions.
func (h *Hub) GetPipelineNames() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	names := make([]string, 0, len(h.pipelines))
	for name := range h.pipelines {
		names = append(names, name)
	}

	// Sort names for consistent output
	sort.Strings(names)
	return names
}

// RemovePipeline removes a named pipeline from the hub.
// This allows for dynamic pipeline management.
//
// Parameters:
//   - name: Name of the pipeline to remove
//
// Returns:
//   - bool: True if the pipeline was removed, false if it didn't exist
//
// Thread-safe: This method is safe for concurrent use.
//
// Note: Removing the "default" pipeline will prevent Pipe() calls
// without a specific pipeline name from working.
func (h *Hub) RemovePipeline(name string) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.pipelines[name]; exists {
		delete(h.pipelines, name)
		return true
	}

	return false
}

// GetPipelineCount returns the number of registered pipelines.
// This is useful for monitoring and debugging purposes.
//
// Returns:
//   - int: Number of registered pipelines
//
// Thread-safe: This method is safe for concurrent use.
func (h *Hub) GetPipelineCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.pipelines)
}

// ClearPipelines removes all registered pipelines from the hub.
// This is useful for testing or resetting the hub state.
//
// Thread-safe: This method is safe for concurrent use.
//
// Warning: This will remove all pipelines including the default one.
func (h *Hub) ClearPipelines() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.pipelines = make(map[string]types.PipelineCallback)
}

// PipelineExecutionError represents an error that occurred during pipeline execution.
// This provides detailed information about which pipeline failed and why.
type PipelineExecutionError struct {
	// PipelineName is the name of the pipeline that failed
	PipelineName string

	// Cause is the underlying error that caused the failure
	Cause error

	// Message is a human-readable description of the error
	Message string
}

// Error returns a string representation of the pipeline execution error.
func (e *PipelineExecutionError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("pipeline '%s' execution failed: %v", e.PipelineName, e.Cause)
}

// Unwrap returns the underlying cause of the error.
// This allows errors.Is and errors.As to work with the wrapped error.
func (e *PipelineExecutionError) Unwrap() error {
	return e.Cause
}
