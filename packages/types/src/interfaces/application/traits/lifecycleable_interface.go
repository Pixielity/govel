package interfaces

import "context"

// LifecycleableInterface defines the contract for application lifecycle management functionality.
type LifecycleableInterface interface {
	// Boot initializes the application and its components
	Boot(ctx context.Context) error
	
	// Booting registers a callback to be executed before providers are booted
	Booting(callback func(interface{}))
	
	// IsBooted returns whether the application has been booted
	IsBooted() bool
	
	// SetBooted sets the booted state of the application
	SetBooted(booted bool)
	
	// Booted registers a callback to be executed after providers have been booted
	Booted(callback func(interface{}))
	
	// Starting registers a callback to be executed before application starts
	Starting(callback func(interface{}))
	
	// Start starts the application after booting
	Start(ctx context.Context) error
	
	// Started registers a callback to be executed after application has started
	Started(callback func(interface{}))
	
	// IsStarted returns whether the application has been started
	IsStarted() bool
	
	// SetStarted sets the started state of the application
	SetStarted(started bool)
	
	// Restart restarts the application (stop then start)
	Restart(ctx context.Context) error
	
	// Stopping registers a callback to be executed before application stops
	Stopping(callback func(interface{}))
	
	// Stop stops the application gracefully
	Stop(ctx context.Context) error
	
	// Stopped registers a callback to be executed after application has stopped
	Stopped(callback func(interface{}))
	
	// IsStopped returns whether the application has been stopped
	IsStopped() bool
	
	// SetStopped sets the stopped state of the application
	SetStopped(stopped bool)
	
	// Terminating registers a callback to be executed during application termination
	Terminating(callback func(interface{})) interface{}
	
	// Terminate terminates the application completely
	Terminate(ctx context.Context) error
	
	// IsTerminated returns whether the application has been terminated
	IsTerminated() bool
	
	// SetTerminated sets the terminated state of the application
	SetTerminated(terminated bool)
	
	// Terminated registers a callback to be executed after application has terminated
	Terminated(callback func(interface{}))
	
	// GetState returns the current lifecycle state of the application
	GetState() string
	
	// IsState checks if the application is in the specified state
	IsState(state string) bool
	
	// GetLifecycleInfo returns comprehensive lifecycle information
	GetLifecycleInfo() map[string]interface{}
	
	// SetState sets the current lifecycle state
	SetState(state string)
}
