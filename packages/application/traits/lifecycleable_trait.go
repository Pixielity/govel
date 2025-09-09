package traits

import (
	"context"
	"sync"

	"govel/packages/application/constants"
	traitInterfaces "govel/packages/application/interfaces/traits"
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
	 * state stores the current lifecycle state
	 */
	state string
}

/**
 * NewLifecycleable creates a new Lifecycleable instance with default values.
 *
 * @return *Lifecycleable The newly created trait instance
 */
func NewLifecycleable() *Lifecycleable {
	return &Lifecycleable{
		booted:  false,
		started: false,
		stopped: false,
		state:   constants.StateInitializing,
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

	// Boot logic would go here
	t.booted = true
	t.state = constants.StateBooted
	return nil
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
		t.state = constants.StateBooted
	} else {
		t.state = constants.StateInitializing
	}
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
		t.booted = true
	}

	if t.started {
		return nil // Already started
	}

	// Start logic would go here
	t.started = true
	t.stopped = false
	t.state = constants.StateRunning
	return nil
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
		t.state = constants.StateRunning
		t.stopped = false
	}
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

	// Stop logic would go here
	t.started = false
	t.stopped = true
	t.state = constants.StateStopped
	return nil
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
		t.state = constants.StateStopped
		t.started = false
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
 * GetState returns the current lifecycle state of the application.
 *
 * @return string The current lifecycle state
 */
func (t *Lifecycleable) GetState() string {
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
		"state":      t.state,
		"is_running": t.started && !t.stopped,
	}
}

/**
 * SetState sets the current lifecycle state.
 *
 * @param state string The state to set
 */
func (t *Lifecycleable) SetState(state string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.state = state

	// Update other states based on the new state
	switch state {
	case constants.StateInitializing:
		t.booted = false
		t.started = false
		t.stopped = false
	case constants.StateBooted:
		t.booted = true
		t.started = false
		t.stopped = false
	case constants.StateRunning:
		t.booted = true
		t.started = true
		t.stopped = false
	case constants.StateStopped:
		t.started = false
		t.stopped = true
	}
}

// Compile-time interface compliance check
var _ traitInterfaces.LifecycleableInterface = (*Lifecycleable)(nil)
