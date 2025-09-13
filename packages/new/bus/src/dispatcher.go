package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"govel/bus/src/interfaces"
)

// Dispatcher is the main command dispatcher implementation
type Dispatcher struct {
	mu           sync.RWMutex
	handlers     map[string]interface{}
	middleware   []interfaces.MiddlewarePipe
	queueResolver func() interface{} // Queue resolver function
}

// NewDispatcher creates a new dispatcher instance
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:   make(map[string]interface{}),
		middleware: make([]interfaces.MiddlewarePipe, 0),
	}
}

// Dispatch dispatches a command to its appropriate handler
func (d *Dispatcher) Dispatch(ctx context.Context, command interface{}) (interface{}, error) {
	// Check if command should be queued
	if d.queueResolver != nil && d.commandShouldBeQueued(command) {
		return nil, d.dispatchToQueue(ctx, command)
	}

	return d.DispatchNow(ctx, command)
}

// DispatchNow dispatches a command immediately in the current process
func (d *Dispatcher) DispatchNow(ctx context.Context, command interface{}) (interface{}, error) {
	// Check if command handles itself
	if selfHandling, ok := command.(interfaces.SelfHandling); ok {
		return d.executeWithPipeline(ctx, command, func(ctx context.Context, cmd interface{}) (interface{}, error) {
			return selfHandling.Handle(ctx)
		})
	}

	// Get handler for command
	handler, err := d.GetCommandHandler(command)
	if err != nil {
		return nil, err
	}

	// Execute with pipeline
	return d.executeWithPipeline(ctx, command, func(ctx context.Context, cmd interface{}) (interface{}, error) {
		return d.callHandler(ctx, handler, cmd)
	})
}

// DispatchSync dispatches a command synchronously
func (d *Dispatcher) DispatchSync(ctx context.Context, command interface{}) (interface{}, error) {
	// For now, DispatchSync is the same as DispatchNow
	// In a full implementation, this would handle sync queue connections
	return d.DispatchNow(ctx, command)
}

// HasCommandHandler checks if a handler exists for the given command
func (d *Dispatcher) HasCommandHandler(command interface{}) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	commandType := d.getCommandType(command)
	_, exists := d.handlers[commandType]
	return exists
}

// GetCommandHandler retrieves the handler for a command
func (d *Dispatcher) GetCommandHandler(command interface{}) (interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	commandType := d.getCommandType(command)
	handler, exists := d.handlers[commandType]
	if !exists {
		return nil, fmt.Errorf("no handler found for command type: %s", commandType)
	}

	return handler, nil
}

// Map registers command to handler mappings
func (d *Dispatcher) Map(commands map[string]string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for commandType, handlerType := range commands {
		d.handlers[commandType] = handlerType
	}
}

// MapHandler registers a specific handler instance for a command type
func (d *Dispatcher) MapHandler(commandType string, handler interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers[commandType] = handler
}

// PipeThrough adds middleware to the pipeline
func (d *Dispatcher) PipeThrough(pipes []interfaces.MiddlewarePipe) interfaces.Dispatcher {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.middleware = append(d.middleware, pipes...)
	return d
}

// SetQueueResolver sets the queue resolver function
func (d *Dispatcher) SetQueueResolver(resolver func() interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.queueResolver = resolver
}

// executeWithPipeline executes a command through the middleware pipeline
func (d *Dispatcher) executeWithPipeline(ctx context.Context, command interface{}, handler func(context.Context, interface{}) (interface{}, error)) (interface{}, error) {
	d.mu.RLock()
	middleware := make([]interfaces.MiddlewarePipe, len(d.middleware))
	copy(middleware, d.middleware)
	d.mu.RUnlock()

	// If no middleware, execute handler directly
	if len(middleware) == 0 {
		return handler(ctx, command)
	}

	// Build the pipeline
	next := handler
	for i := len(middleware) - 1; i >= 0; i-- {
		pipe := middleware[i]
		currentNext := next
		next = func(ctx context.Context, cmd interface{}) (interface{}, error) {
			return pipe.Handle(ctx, cmd, currentNext)
		}
	}

	return next(ctx, command)
}

// callHandler calls the appropriate method on a handler
func (d *Dispatcher) callHandler(ctx context.Context, handler interface{}, command interface{}) (interface{}, error) {
	// Check if handler implements Handler interface
	if h, ok := handler.(interfaces.Handler); ok {
		return h.Handle(ctx, command)
	}

	// Use reflection to call handler methods
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	// Try to find Handle method
	handleMethod, exists := handlerType.MethodByName("Handle")
	if exists {
		// Call Handle method with context and command
		args := []reflect.Value{
			handlerValue,
			reflect.ValueOf(ctx),
			reflect.ValueOf(command),
		}

		results := handleMethod.Func.Call(args)
		if len(results) >= 2 {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				return nil, err
			}
			return results[0].Interface(), nil
		}
		return results[0].Interface(), nil
	}

	return nil, fmt.Errorf("handler does not implement Handle method")
}

// getCommandType returns the type name of a command
func (d *Dispatcher) getCommandType(command interface{}) string {
	t := reflect.TypeOf(command)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

// commandShouldBeQueued determines if a command should be queued
func (d *Dispatcher) commandShouldBeQueued(command interface{}) bool {
	// Check if command implements a Queueable interface
	// This would be defined in traits
	if _, ok := command.(interface{ ShouldQueue() bool }); ok {
		return true
	}
	return false
}

// dispatchToQueue dispatches a command to a queue
func (d *Dispatcher) dispatchToQueue(ctx context.Context, command interface{}) error {
	if d.queueResolver == nil {
		return fmt.Errorf("no queue resolver configured")
	}

	// This would integrate with an actual queue system
	// For now, we'll return an error indicating it's not implemented
	return fmt.Errorf("queue dispatching not implemented")
}

// Ensure Dispatcher implements the Dispatcher interface
var _ interfaces.Dispatcher = (*Dispatcher)(nil)
