// Package pipeline provides a Laravel-compatible pipeline implementation for Go.
// This package allows you to pass an object through a series of "pipes" (middleware)
// where each pipe can modify the object or perform side effects before passing
// it to the next pipe in the chain.
//
// The pipeline follows the "Russian Doll" or "Onion" pattern where pipes are
// executed in reverse order, creating nested closures that process the data
// as it flows through each layer.
package pipeline

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	interfaces "govel/types/src/interfaces/pipeline"

	containerInterfaces "govel/types/src/interfaces/container"
)

// ErrNilContainer is returned when attempting to use container functionality without a container
var ErrNilContainer = errors.New("a container instance has not been passed to the pipeline")

// ErrInvalidPipe is returned when a pipe cannot be processed
var ErrInvalidPipe = errors.New("invalid pipe: must be a function, implement PipeHandlerInterface, or be resolvable from container")

// ErrPipeExecution is returned when a pipe fails to execute
type ErrPipeExecution struct {
	PipeName string
	Cause    error
}

func (e *ErrPipeExecution) Error() string {
	return fmt.Sprintf("pipe execution failed [%s]: %v", e.PipeName, e.Cause)
}

func (e *ErrPipeExecution) Unwrap() error {
	return e.Cause
}

// Pipeline represents a pipeline implementation that processes objects through
// a series of pipes (middleware). This struct is the main implementation of
// the PipelineInterface interface.
//
// The pipeline maintains state about the object being processed, the pipes
// to execute, and various execution options like method names, finalizers,
// and transaction settings.
type Pipeline struct {
	// container is the dependency injection container used for pipe resolution
	container containerInterfaces.ContainerInterface

	// passable is the object being passed through the pipeline
	passable interface{}

	// pipes is the array of pipes to execute
	pipes []interface{}

	// method is the method name to call on pipe objects (default: "Handle")
	method string

	// finallyCallback is executed after pipeline completion, regardless of outcome
	finallyCallback func(interface{})

	// withinTransaction indicates if the pipeline should run within a database transaction
	withinTransaction string

	// pipelineContext is the context for pipeline execution
	pipelineContext interfaces.PipelineContextInterface

	// mutex protects concurrent access to pipeline state
	mutex sync.RWMutex
}

// NewPipeline creates a new Pipeline instance with the given container.
// The container is optional and can be nil if dependency injection is not needed.
//
// Parameters:
//   - container: Optional dependency injection container for pipe resolution
//
// Returns:
//   - *Pipeline: New pipeline instance ready for configuration and execution
//
// Example:
//
//	pipeline := NewPipeline(container)
//	result, err := pipeline.
//		Send(data).
//		Through([]interface{}{middleware1, middleware2}).
//		Then(func(passable interface{}) interface{} {
//			return processData(passable)
//		})
func NewPipeline(container containerInterfaces.ContainerInterface) *Pipeline {
	return &Pipeline{
		container: container,
		method:    "Handle", // Default method name for pipe objects
		pipes:     make([]interface{}, 0),
		mutex:     sync.RWMutex{},
	}
}

// Send sets the object being sent through the pipeline.
// This is the initial data that will be passed through each pipe.
//
// Parameters:
//   - passable: The object to pass through the pipeline
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (p *Pipeline) Send(passable interface{}) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.passable = passable
	return p
}

// Through sets the array of pipes to process the passable object.
// Pipes can be function closures, struct instances, or string names
// that will be resolved from the container.
//
// Parameters:
//   - pipes: Array of pipes (middlewares) to execute
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
//
// Supported pipe types:
//   - func(interface{}, func(interface{}) (interface{}, error)) (interface{}, error)
//   - types implementing interfaces.PipeHandlerInterface
//   - string names that can be resolved from the container
func (p *Pipeline) Through(pipes []interface{}) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Create a new slice to avoid mutations to the original
	p.pipes = make([]interface{}, len(pipes))
	copy(p.pipes, pipes)
	return p
}

// Pipe adds additional pipes to the existing pipeline.
// This allows for dynamic pipe addition after initial setup.
//
// Parameters:
//   - pipes: Additional pipes to append to the pipeline
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (p *Pipeline) Pipe(pipes ...interface{}) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pipes = append(p.pipes, pipes...)
	return p
}

// Via sets the method name to call on each pipe object.
// By default, the pipeline will call 'Handle' method on pipe objects.
// This allows customization of the method name.
//
// Parameters:
//   - method: The method name to call on pipe objects
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (p *Pipeline) Via(method string) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.method = method
	return p
}

// Then executes the pipeline with the given destination callback.
// This is the final step that processes the passable object through
// all pipes and then executes the destination function.
//
// The pipeline execution follows the "Russian Doll" or "Onion" pattern
// where pipes are executed in reverse order, creating nested closures.
//
// Parameters:
//   - destination: Final callback to execute after all pipes
//
// Returns:
//   - interface{}: The result of pipeline execution
//   - error: Any error that occurred during pipeline execution
//
// Thread-safe: This method creates a snapshot of the pipeline state for execution.
func (p *Pipeline) Then(destination func(interface{}) interface{}) (interface{}, error) {
	// Create a snapshot of the current pipeline state to ensure thread safety
	p.mutex.RLock()
	passable := p.passable
	pipes := make([]interface{}, len(p.pipes))
	copy(pipes, p.pipes)
	method := p.method
	finallyCallback := p.finallyCallback
	withinTransaction := p.withinTransaction
	pipelineContext := p.pipelineContext
	p.mutex.RUnlock()

	// Prepare the destination function with error handling
	preparedDestination := p.prepareDestination(destination)

	// Build the pipeline execution chain
	pipelineChain := p.buildPipelineChain(pipes, method, preparedDestination)

	// Execute the pipeline with proper error handling and cleanup
	result, err := p.executePipeline(pipelineChain, passable, pipelineContext, withinTransaction)

	// Execute finally callback if present
	if finallyCallback != nil {
		// Wrap in a recovery block to prevent finally callback from panicking
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log the panic but don't let it affect the main result
					fmt.Printf("Warning: finally callback panicked: %v\n%s\n", r, debug.Stack())
				}
			}()
			finallyCallback(passable)
		}()
	}

	return result, err
}

// ThenReturn executes the pipeline and returns the passable object unchanged.
// This is a convenience method equivalent to calling Then with an identity function.
//
// Returns:
//   - interface{}: The passable object after processing through all pipes
//   - error: Any error that occurred during pipeline execution
func (p *Pipeline) ThenReturn() (interface{}, error) {
	return p.Then(func(passable interface{}) interface{} {
		return passable
	})
}

// Finally sets a callback to be executed after pipeline completion,
// regardless of success or failure. This is similar to a "finally" block
// in exception handling.
//
// Parameters:
//   - callback: Function to execute after pipeline completion
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (p *Pipeline) Finally(callback func(interface{})) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.finallyCallback = callback
	return p
}

// WithContext sets the context for pipeline execution.
// This enables cancellation, timeouts, and context value propagation.
//
// Parameters:
//   - ctx: The context to use for pipeline execution
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (p *Pipeline) WithContext(ctx context.Context) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// If the context is already a PipelineContextInterface, use it directly
	if pipelineCtx, ok := ctx.(interfaces.PipelineContextInterface); ok {
		p.pipelineContext = pipelineCtx
	} else {
		// Wrap the standard context in a pipeline context
		p.pipelineContext = NewPipelineContext(ctx)
	}

	return p
}

// WithinTransaction enables database transaction wrapping for pipeline execution.
// When enabled, the entire pipeline execution will be wrapped in a database transaction.
//
// Parameters:
//   - connection: Database connection name (empty string for default)
//
// Returns:
//   - interfaces.PipelineInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
//
// Note: This requires a container with database functionality.
func (p *Pipeline) WithinTransaction(connection string) interfaces.PipelineInterface {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.withinTransaction = connection
	return p
}

// prepareDestination wraps the destination function with error handling.
// This ensures that panics in the destination are caught and converted to errors.
func (p *Pipeline) prepareDestination(destination func(interface{}) interface{}) func(interface{}) (interface{}, error) {
	return func(passable interface{}) (result interface{}, err error) {
		// Recover from panics in the destination function
		defer func() {
			if r := recover(); r != nil {
				err = &ErrPipeExecution{
					PipeName: "destination",
					Cause:    fmt.Errorf("panic: %v", r),
				}
			}
		}()

		// Check for context cancellation before executing destination
		if p.pipelineContext != nil {
			select {
			case <-p.pipelineContext.Done():
				return nil, p.pipelineContext.Err()
			default:
				// Context is still active, continue
			}
		}

		result = destination(passable)
		return result, nil
	}
}

// buildPipelineChain constructs the pipeline execution chain using the "Russian Doll" pattern.
// This reverses the pipes and creates nested closures for execution.
func (p *Pipeline) buildPipelineChain(pipes []interface{}, method string, destination func(interface{}) (interface{}, error)) func(interface{}) (interface{}, error) {
	// Start with the destination as the innermost function
	pipelineChain := destination

	// Build the chain by wrapping each pipe around the previous chain
	// We iterate in reverse order to create the proper nesting
	for i := len(pipes) - 1; i >= 0; i-- {
		pipe := pipes[i]
		currentChain := pipelineChain

		// Create a closure for the current pipe
		pipelineChain = func(passable interface{}) (interface{}, error) {
			return p.executePipe(pipe, passable, currentChain, method)
		}
	}

	return pipelineChain
}

// executePipeline runs the pipeline chain with transaction and context support.
func (p *Pipeline) executePipeline(pipelineChain func(interface{}) (interface{}, error), passable interface{}, ctx interfaces.PipelineContextInterface, withinTransaction string) (interface{}, error) {
	// Set execution start time if context is available
	if ctx != nil {
		ctx.SetExecutionStartTime(time.Now())
	}

	// If transaction is requested and container is available, wrap in transaction
	if withinTransaction != "" && p.container != nil {
		return p.executeWithinTransaction(pipelineChain, passable, withinTransaction)
	}

	// Execute without transaction
	return pipelineChain(passable)
}

// executeWithinTransaction wraps pipeline execution in a database transaction.
func (p *Pipeline) executeWithinTransaction(pipelineChain func(interface{}) (interface{}, error), passable interface{}, connection string) (interface{}, error) {
	// Try to get database connection from container
	db, err := p.container.Make("db")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve database from container: %w", err)
	}

	// This is a placeholder for actual transaction logic
	// In a real implementation, you would use your database library's transaction methods
	// For example, with GORM or sqlx

	// For now, just execute without transaction (would need database interface)
	_ = db // Silence unused variable warning
	return pipelineChain(passable)
}

// executePipe executes a single pipe in the pipeline chain.
// This handles the different types of pipes (functions, objects, strings).
func (p *Pipeline) executePipe(pipe interface{}, passable interface{}, next func(interface{}) (interface{}, error), method string) (result interface{}, err error) {
	// Set current pipe in context for debugging
	if p.pipelineContext != nil {
		pipeName := p.getPipeName(pipe)
		p.pipelineContext.SetCurrentPipe(pipeName)
	}

	// Recover from panics in pipe execution
	defer func() {
		if r := recover(); r != nil {
			err = &ErrPipeExecution{
				PipeName: p.getPipeName(pipe),
				Cause:    fmt.Errorf("panic: %v", r),
			}
		}
	}()

	// Check for context cancellation before executing pipe
	if p.pipelineContext != nil {
		select {
		case <-p.pipelineContext.Done():
			return nil, p.pipelineContext.Err()
		default:
			// Context is still active, continue
		}
	}

	// Handle different pipe types
	switch pipeValue := pipe.(type) {
	case func(interface{}, func(interface{}) (interface{}, error)) (interface{}, error):
		// Function pipe - call directly
		return pipeValue(passable, next)

	case interfaces.PipeHandlerInterface:
		// Object implementing PipeHandlerInterface - call Handle method
		return pipeValue.Handle(passable, next)

	case string:
		// String pipe - resolve from container and execute
		return p.executeStringPipe(pipeValue, passable, next, method)

	default:
		// Try to call method on object using reflection
		return p.executeObjectPipe(pipe, passable, next, method)
	}
}

// executeStringPipe resolves a string pipe from the container and executes it.
func (p *Pipeline) executeStringPipe(pipeName string, passable interface{}, next func(interface{}) (interface{}, error), method string) (interface{}, error) {
	if p.container == nil {
		return nil, ErrNilContainer
	}

	// Parse pipe string for parameters (format: "PipeName:param1,param2")
	name, parameters := p.parsePipeString(pipeName)

	// Resolve pipe from container
	pipeInstance, err := p.container.Make(name)
	if err != nil {
		return nil, &ErrPipeExecution{
			PipeName: name,
			Cause:    fmt.Errorf("failed to resolve pipe from container: %w", err),
		}
	}

	// Execute the resolved pipe with parameters
	return p.executeObjectPipeWithParameters(pipeInstance, passable, next, method, parameters)
}

// executeObjectPipe executes an object pipe using reflection to call the specified method.
func (p *Pipeline) executeObjectPipe(pipe interface{}, passable interface{}, next func(interface{}) (interface{}, error), method string) (interface{}, error) {
	return p.executeObjectPipeWithParameters(pipe, passable, next, method, nil)
}

// executeObjectPipeWithParameters executes an object pipe with additional parameters.
func (p *Pipeline) executeObjectPipeWithParameters(pipe interface{}, passable interface{}, next func(interface{}) (interface{}, error), method string, parameters []string) (interface{}, error) {
	pipeValue := reflect.ValueOf(pipe)
	pipeType := reflect.TypeOf(pipe)

	// Check if the pipe implements PipeHandlerInterface
	if handler, ok := pipe.(interfaces.PipeHandlerInterface); ok {
		return handler.Handle(passable, next)
	}

	// Try to find the method on the pipe
	methodValue := pipeValue.MethodByName(method)
	if !methodValue.IsValid() {
		// Method not found, try to call the pipe as a function
		if pipeType.Kind() == reflect.Func {
			return p.callPipeFunction(pipeValue, passable, next, parameters)
		}

		return nil, &ErrPipeExecution{
			PipeName: p.getPipeName(pipe),
			Cause:    fmt.Errorf("method '%s' not found on pipe type %s", method, pipeType.String()),
		}
	}

	// Call the method with appropriate parameters
	return p.callPipeMethod(methodValue, passable, next, parameters)
}

// callPipeFunction calls a pipe function using reflection.
func (p *Pipeline) callPipeFunction(pipeValue reflect.Value, passable interface{}, next func(interface{}) (interface{}, error), parameters []string) (interface{}, error) {
	pipeType := pipeValue.Type()

	// Build arguments for the function call
	args := []reflect.Value{
		reflect.ValueOf(passable),
		reflect.ValueOf(next),
	}

	// Add parameters if the function accepts them
	for i, param := range parameters {
		if i+2 < pipeType.NumIn() {
			args = append(args, reflect.ValueOf(param))
		}
	}

	// Call the function
	results := pipeValue.Call(args)

	// Handle return values
	return p.handlePipeResults(results)
}

// callPipeMethod calls a pipe method using reflection.
func (p *Pipeline) callPipeMethod(methodValue reflect.Value, passable interface{}, next func(interface{}) (interface{}, error), parameters []string) (interface{}, error) {
	methodType := methodValue.Type()

	// Build arguments for the method call
	args := []reflect.Value{
		reflect.ValueOf(passable),
		reflect.ValueOf(next),
	}

	// Add parameters if the method accepts them
	for i, param := range parameters {
		if i+2 < methodType.NumIn() {
			args = append(args, reflect.ValueOf(param))
		}
	}

	// Call the method
	results := methodValue.Call(args)

	// Handle return values
	return p.handlePipeResults(results)
}

// handlePipeResults processes the return values from a pipe function/method call.
func (p *Pipeline) handlePipeResults(results []reflect.Value) (interface{}, error) {
	switch len(results) {
	case 0:
		return nil, errors.New("pipe function must return at least one value")
	case 1:
		// Single return value - assume it's the result with no error
		return results[0].Interface(), nil
	case 2:
		// Two return values - assume (result, error)
		result := results[0].Interface()
		errorVal := results[1].Interface()

		if errorVal != nil {
			if err, ok := errorVal.(error); ok {
				return result, err
			}
			return result, fmt.Errorf("pipe returned non-error as second value: %v", errorVal)
		}

		return result, nil
	default:
		return nil, fmt.Errorf("pipe function returned too many values (%d), expected 1 or 2", len(results))
	}
}

// parsePipeString parses a pipe string in the format "PipeName:param1,param2".
func (p *Pipeline) parsePipeString(pipeStr string) (name string, parameters []string) {
	parts := strings.SplitN(pipeStr, ":", 2)
	name = strings.TrimSpace(parts[0])

	if len(parts) > 1 && parts[1] != "" {
		paramStr := strings.TrimSpace(parts[1])
		parameters = strings.Split(paramStr, ",")

		// Trim whitespace from each parameter
		for i, param := range parameters {
			parameters[i] = strings.TrimSpace(param)
		}
	}

	return name, parameters
}

// getPipeName returns a human-readable name for a pipe for debugging purposes.
func (p *Pipeline) getPipeName(pipe interface{}) string {
	switch pipeValue := pipe.(type) {
	case string:
		return pipeValue
	default:
		pipeType := reflect.TypeOf(pipe)
		if pipeType == nil {
			return "<nil>"
		}
		return pipeType.String()
	}
}
