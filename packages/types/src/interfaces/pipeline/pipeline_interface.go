package interfaces

import (
	"context"
)

// PipelineInterface defines the contract for pipeline operations.
// This interface provides a Laravel-compatible pipeline implementation that allows
// passing objects through a series of "pipes" (middleware) where each pipe can
// modify the object or perform side effects.
//
// Key features:
//   - Method chaining for fluent API
//   - Support for multiple pipe types (functions, objects, strings)
//   - Context support for cancellation and timeouts
//   - Transaction support for database operations
//   - Finally callbacks for cleanup operations
//   - Thread-safe pipeline state management
type PipelineInterface interface {
	// Send sets the object being sent through the pipeline.
	// This is the initial data that will be passed through each pipe.
	//
	// Parameters:
	//   - passable: The object to pass through the pipeline
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	Send(passable interface{}) PipelineInterface

	// Through sets the array of pipes to process the passable object.
	// Pipes can be function closures, struct instances, or string names.
	//
	// Parameters:
	//   - pipes: Array of pipes (middlewares) to execute
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	Through(pipes []interface{}) PipelineInterface

	// Pipe adds additional pipes to the existing pipeline.
	// This allows for dynamic pipe addition after initial setup.
	//
	// Parameters:
	//   - pipes: Additional pipes to append to the pipeline
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	Pipe(pipes ...interface{}) PipelineInterface

	// Via sets the method name to call on each pipe object.
	// By default, the pipeline will call 'Handle' method on pipe objects.
	//
	// Parameters:
	//   - method: The method name to call on pipe objects
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	Via(method string) PipelineInterface

	// Then executes the pipeline with the given destination callback.
	// This is the final step that processes the passable object through
	// all pipes and then executes the destination function.
	//
	// Parameters:
	//   - destination: Final callback to execute after all pipes
	//
	// Returns:
	//   - interface{}: The result of pipeline execution
	//   - error: Any error that occurred during pipeline execution
	Then(destination func(interface{}) interface{}) (interface{}, error)

	// ThenReturn executes the pipeline and returns the passable object unchanged.
	// This is a convenience method equivalent to calling Then with an identity function.
	//
	// Returns:
	//   - interface{}: The passable object after processing through all pipes
	//   - error: Any error that occurred during pipeline execution
	ThenReturn() (interface{}, error)

	// Finally sets a callback to be executed after pipeline completion,
	// regardless of success or failure. This is similar to a "finally" block.
	//
	// Parameters:
	//   - callback: Function to execute after pipeline completion
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	Finally(callback func(interface{})) PipelineInterface

	// WithContext sets the context for pipeline execution.
	// This enables cancellation, timeouts, and context value propagation.
	//
	// Parameters:
	//   - ctx: The context to use for pipeline execution
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	WithContext(ctx context.Context) PipelineInterface

	// WithinTransaction enables database transaction wrapping for pipeline execution.
	// When enabled, the entire pipeline execution will be wrapped in a database transaction.
	//
	// Parameters:
	//   - connection: Database connection name (empty string for default)
	//
	// Returns:
	//   - PipelineInterface: Returns self for method chaining
	WithinTransaction(connection string) PipelineInterface
}
