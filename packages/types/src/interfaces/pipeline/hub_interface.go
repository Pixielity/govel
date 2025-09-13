package interfaces

import (
	containerInterfaces "govel/packages/types/src/interfaces/container"
)

// HubInterface defines the contract for pipeline hub operations.
// This interface provides a centralized way to manage and execute named pipelines.
//
// Key features:
//   - Named pipeline registration and management
//   - Default pipeline configuration
//   - Pipeline execution with error handling
//   - Container integration for dependency injection
//   - Thread-safe operations
type HubInterface interface {
	// Defaults defines the default pipeline configuration.
	// This pipeline will be used when no specific pipeline name is provided.
	//
	// Parameters:
	//   - callback: Function that defines the default pipeline configuration
	Defaults(callback func(PipelineInterface, interface{}) interface{})

	// Pipeline defines a named pipeline configuration.
	// This allows for different pipeline setups to be registered under specific names.
	//
	// Parameters:
	//   - name: Unique name for the pipeline configuration
	//   - callback: Function that defines the pipeline configuration
	Pipeline(name string, callback func(PipelineInterface, interface{}) interface{})

	// Pipe sends an object through one of the available pipelines.
	// If no pipeline name is specified, the default pipeline will be used.
	//
	// Parameters:
	//   - object: The object to pass through the pipeline
	//   - pipelineName: Optional name of the pipeline to use
	//
	// Returns:
	//   - interface{}: The result of pipeline execution
	//   - error: Any error that occurred during pipeline execution
	Pipe(object interface{}, pipelineName ...string) (interface{}, error)

	// GetContainer returns the container instance used by the hub.
	// This provides access to the dependency injection container.
	//
	// Returns:
	//   - ContainerInterface: The container instance (may be nil)
	GetContainer() containerInterfaces.ContainerInterface

	// SetContainer sets the container instance used by the hub.
	// This allows for dependency injection container replacement.
	//
	// Parameters:
	//   - container: The container instance to use
	//
	// Returns:
	//   - HubInterface: Returns self for method chaining
	SetContainer(container containerInterfaces.ContainerInterface) HubInterface

	// HasPipeline checks if a named pipeline exists in the hub.
	// This is useful for conditional pipeline execution or validation.
	//
	// Parameters:
	//   - name: Name of the pipeline to check
	//
	// Returns:
	//   - bool: True if the pipeline exists, false otherwise
	HasPipeline(name string) bool

	// GetPipelineNames returns all registered pipeline names.
	// This is useful for debugging and introspection.
	//
	// Returns:
	//   - []string: List of all registered pipeline names
	GetPipelineNames() []string

	// RemovePipeline removes a named pipeline from the hub.
	// This is useful for dynamic pipeline management.
	//
	// Parameters:
	//   - name: Name of the pipeline to remove
	//
	// Returns:
	//   - bool: True if the pipeline was removed, false if it didn't exist
	RemovePipeline(name string) bool

	// ClearPipelines removes all registered pipelines from the hub.
	// This is useful for testing or resetting the hub state.
	ClearPipelines()

	// GetPipelineCount returns the number of registered pipelines.
	// This is useful for monitoring and debugging purposes.
	//
	// Returns:
	//   - int: Number of registered pipelines
	GetPipelineCount() int
}
