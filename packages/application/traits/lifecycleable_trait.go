package traits

import (
	"context"
	"sync"

	"govel/types/src/enums/application"
	traitInterfaces "govel/types/src/interfaces/application/traits"
)

/**
 * Lifecycleable provides application lifecycle management functionality in a thread-safe manner.
 * This trait follows the self-contained pattern with dependency injection through interfaces.
 */
type Lifecycleable struct {
	/**
	 * mutex provides thread-safe access to lifecycle properties
	 */
	mutex sync.RWMutex

	/**
	 * booted indicates whether the application has been booted
	 */
	booted bool

	/**
	 * started indicates whether the application has been started
	 */
	started bool

	/**
	 * stopped indicates whether the application has been stopped
	 */
	stopped bool

	/**
	 * terminated indicates whether the application has been terminated
	 */
	terminated bool

	/**
	 * state stores the current lifecycle state
	 */
	state enums.LifecycleState

	/**
	 * lifecycle callbacks
	 */
	bootingCallbacks     []func(interface{})
	bootedCallbacks      []func(interface{})
	startingCallbacks    []func(interface{})
	startedCallbacks     []func(interface{})
	stoppingCallbacks    []func(interface{})
	stoppedCallbacks     []func(interface{})
	terminatingCallbacks []func(interface{})
	terminatedCallbacks  []func(interface{})
}

/**
 * NewLifecycleable creates a new Lifecycleable instance with default values.
 *
 * @return *Lifecycleable The newly created trait instance
 */
func NewLifecycleable() *Lifecycleable {
	return &Lifecycleable{
		booted:     false,
		started:    false,
		stopped:    false,
		terminated: false,
		state:      enums.StateInitializing,
	}
}

/**
 * Boot initializes the application and its components.
 *
 * @param ctx context.Context The context for the boot operation
 * @return error Any error that occurred during boot
 */
func (t *Lifecycleable) Boot(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.booted {
		return nil // Already booted
	}

	// Set booting state
	t.state = enums.StateBooting

	// Execute booting callbacks
	t.executeCallbacks(t.bootingCallbacks, t)

	// Boot logic would go here
	t.booted = true
	t.state = enums.StateBooted

	// Execute booted callbacks
	t.executeCallbacks(t.bootedCallbacks, t)

	return nil
}

/**
 * Booting registers a callback to be executed before providers are booted.
 *
 * @param callback func(interface{}) The function to execute before booting
 */
func (t *Lifecycleable) Booting(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.bootingCallbacks = append(t.bootingCallbacks, callback)
}

/**
 * IsBooted returns whether the application has been booted.
 *
 * @return bool true if the application is booted
 */
func (t *Lifecycleable) IsBooted() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.booted
}

/**
 * SetBooted sets the booted state of the application.
 *
 * @param booted bool Whether the application is booted
 */
func (t *Lifecycleable) SetBooted(booted bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.booted = booted
	if booted {
		t.state = enums.StateBooted
	} else {
		t.state = enums.StateInitializing
	}
}

/**
 * Booted registers a callback to be executed after providers have been booted.
 *
 * @param callback func(interface{}) The function to execute after booting
 */
func (t *Lifecycleable) Booted(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.bootedCallbacks = append(t.bootedCallbacks, callback)
}

/**
 * Starting registers a callback to be executed before application starts.
 *
 * @param callback func(interface{}) The function to execute before starting
 */
func (t *Lifecycleable) Starting(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.startingCallbacks = append(t.startingCallbacks, callback)
}

/**
 * Start starts the application after booting.
 *
 * @param ctx context.Context The context for the start operation
 * @return error Any error that occurred during start
 */
func (t *Lifecycleable) Start(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.booted {
		// Auto-boot if not already booted
		t.state = enums.StateBooting
		t.executeCallbacks(t.bootingCallbacks, t)
		t.booted = true
		t.state = enums.StateBooted
		t.executeCallbacks(t.bootedCallbacks, t)
	}

	if t.started {
		return nil // Already started
	}

	// Set starting state
	t.state = enums.StateStarting

	// Execute starting callbacks
	t.executeCallbacks(t.startingCallbacks, t)

	// Start logic would go here
	t.started = true
	t.stopped = false
	t.state = enums.StateRunning

	// Execute started callbacks
	t.executeCallbacks(t.startedCallbacks, t)

	return nil
}

/**
 * Started registers a callback to be executed after application has started.
 *
 * @param callback func(interface{}) The function to execute after starting
 */
func (t *Lifecycleable) Started(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.startedCallbacks = append(t.startedCallbacks, callback)
}

/**
 * IsStarted returns whether the application has been started.
 *
 * @return bool true if the application is started
 */
func (t *Lifecycleable) IsStarted() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.started
}

/**
 * SetStarted sets the started state of the application.
 *
 * @param started bool Whether the application is started
 */
func (t *Lifecycleable) SetStarted(started bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.started = started
	if started {
		t.state = enums.StateRunning
		t.stopped = false
	}
}

/**
 * Restart restarts the application (stop then start).
 *
 * @param ctx context.Context The context for the restart operation
 * @return error Any error that occurred during restart
 */
func (t *Lifecycleable) Restart(ctx context.Context) error {
	// Stop first
	if err := t.Stop(ctx); err != nil {
		return err
	}

	// Then start
	return t.Start(ctx)
}

/**
 * Stopping registers a callback to be executed before application stops.
 *
 * @param callback func(interface{}) The function to execute before stopping
 */
func (t *Lifecycleable) Stopping(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.stoppingCallbacks = append(t.stoppingCallbacks, callback)
}

/**
 * Stop stops the application gracefully.
 *
 * @param ctx context.Context The context for the stop operation
 * @return error Any error that occurred during stop
 */
func (t *Lifecycleable) Stop(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.stopped {
		return nil // Already stopped
	}

	// Set stopping state
	t.state = enums.StateStopping

	// Execute stopping callbacks
	t.executeCallbacks(t.stoppingCallbacks, t)

	// Stop logic would go here
	t.started = false
	t.stopped = true
	t.state = enums.StateStopped

	// Execute stopped callbacks
	t.executeCallbacks(t.stoppedCallbacks, t)

	return nil
}

/**
 * Stopped registers a callback to be executed after application has stopped.
 *
 * @param callback func(interface{}) The function to execute after stopping
 */
func (t *Lifecycleable) Stopped(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.stoppedCallbacks = append(t.stoppedCallbacks, callback)
}

/**
 * IsStopped returns whether the application has been stopped.
 *
 * @return bool true if the application is stopped
 */
func (t *Lifecycleable) IsStopped() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.stopped
}

/**
 * SetStopped sets the stopped state of the application.
 *
 * @param stopped bool Whether the application is stopped
 */
func (t *Lifecycleable) SetStopped(stopped bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.stopped = stopped
	if stopped {
		t.state = enums.StateStopped
		t.started = false
	}
}

/**
 * Terminating registers a callback to be executed during application termination.
 *
 * @param callback func(interface{}) The function to execute during termination
 * @return interface{} Returns the trait instance for method chaining
 */
func (t *Lifecycleable) Terminating(callback func(interface{})) interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.terminatingCallbacks = append(t.terminatingCallbacks, callback)
	return t
}

/**
 * Terminate terminates the application completely.
 *
 * @param ctx context.Context The context for the terminate operation
 * @return error Any error that occurred during termination
 */
func (t *Lifecycleable) Terminate(ctx context.Context) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.terminated {
		return nil // Already terminated
	}

	// Set terminating state
	t.state = enums.StateTerminating

	// Execute terminating callbacks
	t.executeCallbacks(t.terminatingCallbacks, t)

	// Terminate logic would go here
	t.terminated = true
	t.started = false
	t.stopped = true
	t.state = enums.StateTerminated

	// Execute terminated callbacks
	t.executeCallbacks(t.terminatedCallbacks, t)

	return nil
}

/**
 * IsTerminated returns whether the application has been terminated.
 *
 * @return bool true if the application is terminated
 */
func (t *Lifecycleable) IsTerminated() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.terminated
}

/**
 * SetTerminated sets the terminated state of the application.
 *
 * @param terminated bool Whether the application is terminated
 */
func (t *Lifecycleable) SetTerminated(terminated bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.terminated = terminated
	if terminated {
		t.state = enums.StateTerminated
		t.started = false
		t.stopped = true
	}
}

/**
 * Terminated registers a callback to be executed after application has terminated.
 *
 * @param callback func(interface{}) The function to execute after termination
 */
func (t *Lifecycleable) Terminated(callback func(interface{})) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.terminatedCallbacks = append(t.terminatedCallbacks, callback)
}

/**
 * GetState returns the current lifecycle state of the application.
 *
 * @return string The current lifecycle state
 */
func (t *Lifecycleable) GetState() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.state.String()
}

/**
 * GetStateEnum returns the current lifecycle state as an enum.
 *
 * @return enums.LifecycleState The current lifecycle state enum
 */
func (t *Lifecycleable) GetStateEnum() enums.LifecycleState {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.state
}

/**
 * IsState checks if the application is in the specified state.
 *
 * @param state string The state to check against
 * @return bool true if the application is in the specified state
 */
func (t *Lifecycleable) IsState(state string) bool {
	return t.GetState() == state
}

/**
 * GetLifecycleInfo returns comprehensive lifecycle information.
 *
 * @return map[string]interface{} Lifecycle details
 */
func (t *Lifecycleable) GetLifecycleInfo() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return map[string]interface{}{
		"booted":     t.booted,
		"started":    t.started,
		"stopped":    t.stopped,
		"terminated": t.terminated,
		"state":      t.state,
		"is_running": t.started && !t.stopped && !t.terminated,
	}
}

/**
 * SetState sets the current lifecycle state.
 *
 * @param state string The state to set
 */
func (t *Lifecycleable) SetState(state string) {
	// Convert string to enum and call SetStateEnum
	stateEnum, valid := enums.FromString(state)
	if !valid {
		// Default to error state if invalid
		stateEnum = enums.StateError
	}
	t.SetStateEnum(stateEnum)
}

/**
 * SetStateEnum sets the current lifecycle state using an enum.
 *
 * @param state enums.LifecycleState The state enum to set
 */
func (t *Lifecycleable) SetStateEnum(state enums.LifecycleState) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !state.IsValid() {
		// Default to error state if invalid
		state = enums.StateError
	}

	t.state = state

	// Update other states based on the new state
	switch state {
	case enums.StateInitializing:
		t.booted = false
		t.started = false
		t.stopped = false
		t.terminated = false
	case enums.StateBooting:
		// During boot process - no changes to boolean flags yet
	case enums.StateBooted:
		t.booted = true
		t.started = false
		t.stopped = false
	case enums.StateStarting:
		// During start process - booted should already be true
		t.booted = true
	case enums.StateRunning:
		t.booted = true
		t.started = true
		t.stopped = false
	case enums.StateStopping:
		// During stop process - no changes to boolean flags yet
	case enums.StateStopped:
		t.started = false
		t.stopped = true
	case enums.StateTerminating:
		// During terminate process - no changes to boolean flags yet
	case enums.StateTerminated:
		t.terminated = true
		t.started = false
		t.stopped = true
	// No explicit cases for StateMaintenance, StateError, StateShuttingDown
	// as they don't directly map to boolean state changes
	}
}

/**
 * executeCallbacks safely executes a slice of callbacks with panic recovery.
 *
 * @param callbacks []func(interface{}) The callbacks to execute
 * @param app interface{} The application instance to pass to callbacks
 */
func (t *Lifecycleable) executeCallbacks(callbacks []func(interface{}), app interface{}) {
	for _, callback := range callbacks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but don't crash the application
					// In a real implementation, you might want to use a proper logger
				}
			}()
			callback(app)
		}()
	}
}

// Compile-time interface compliance check
var _ traitInterfaces.LifecycleableInterface = (*Lifecycleable)(nil)
