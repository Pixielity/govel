package interfaces

import (
	"context"
)

// Dispatcher defines the contract for command/job dispatching
type Dispatcher interface {
	// Dispatch dispatches a command to its appropriate handler
	// Returns the result of the command execution or an error
	Dispatch(ctx context.Context, command interface{}) (interface{}, error)

	// DispatchNow dispatches a command immediately in the current process
	// Bypasses any queue system and executes synchronously
	DispatchNow(ctx context.Context, command interface{}) (interface{}, error)

	// DispatchSync dispatches a command synchronously
	// Similar to DispatchNow but respects queue configuration for sync execution
	DispatchSync(ctx context.Context, command interface{}) (interface{}, error)

	// HasCommandHandler checks if a handler exists for the given command
	HasCommandHandler(command interface{}) bool

	// GetCommandHandler retrieves the handler for a command
	GetCommandHandler(command interface{}) (interface{}, error)

	// Map registers a command to handler mapping
	Map(commands map[string]string)

	// PipeThrough adds middleware to the pipeline
	PipeThrough(pipes []MiddlewarePipe) Dispatcher
}

// MiddlewarePipe defines the contract for middleware in the dispatcher pipeline
type MiddlewarePipe interface {
	// Handle processes the command through middleware
	Handle(ctx context.Context, command interface{}, next func(ctx context.Context, command interface{}) (interface{}, error)) (interface{}, error)
}

// Handler defines the contract for command handlers
type Handler interface {
	// Handle processes the command and returns a result
	Handle(ctx context.Context, command interface{}) (interface{}, error)
}

// SelfHandling defines the contract for commands that handle themselves
type SelfHandling interface {
	// Handle processes the command
	Handle(ctx context.Context) (interface{}, error)
}
