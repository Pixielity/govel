package traits

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"govel/packages/application/constants"
	traitInterfaces "govel/packages/application/interfaces/traits"
	"govel/packages/application/types"
)

/**
 * Shutdownable provides application shutdown management functionality in a thread-safe manner.
 * This trait follows the self-contained pattern with dependency injection through interfaces.
 */
type Shutdownable struct {
	/**
	 * mutex provides thread-safe access to shutdown properties
	 */
	mutex sync.RWMutex

	/**
	 * callbacks stores registered shutdown callbacks
	 */
	callbacks map[string]types.ShutdownCallback

	/**
	 * isShuttingDown indicates whether shutdown is in progress
	 */
	isShuttingDown bool

	/**
	 * isShutdown indicates whether shutdown is complete
	 */
	isShutdown bool

	/**
	 * shutdownTimeout stores the timeout for shutdown operations
	 */
	shutdownTimeout time.Duration

	/**
	 * signalChannel handles OS signals for graceful shutdown
	 */
	signalChannel chan os.Signal
}

/**
 * NewShutdownable creates a new Shutdownable instance with default values.
 *
 * @return *Shutdownable The newly created trait instance
 */
func NewShutdownable() *Shutdownable {
	return &Shutdownable{
		callbacks:       make(map[string]types.ShutdownCallback),
		isShuttingDown:  false,
		isShutdown:      false,
		shutdownTimeout: constants.DefaultShutdownTimeout,
		signalChannel:   make(chan os.Signal, 1),
	}
}

/**
 * RegisterShutdownCallback registers a callback to be executed during shutdown.
 *
 * @param name string The name of the callback for identification
 * @param callback ShutdownCallback The callback function to register
 */
func (t *Shutdownable) RegisterShutdownCallback(name string, callback types.ShutdownCallback) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.callbacks[name] = callback
}

/**
 * UnregisterShutdownCallback removes a registered shutdown callback.
 *
 * @param name string The name of the callback to unregister
 * @return bool true if the callback was found and removed
 */
func (t *Shutdownable) UnregisterShutdownCallback(name string) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.callbacks[name]; exists {
		delete(t.callbacks, name)
		return true
	}
	return false
}

/**
 * GetShutdownCallbacks returns all registered shutdown callbacks.
 *
 * @return map[string]ShutdownCallback Map of callback names to callbacks
 */
func (t *Shutdownable) GetShutdownCallbacks() map[string]types.ShutdownCallback {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Return a copy to prevent external modification
	callbacks := make(map[string]types.ShutdownCallback)
	for name, callback := range t.callbacks {
		callbacks[name] = callback
	}
	return callbacks
}

/**
 * Shutdown initiates the application shutdown process.
 *
 * @param ctx context.Context The context for the shutdown operation
 * @return error Any error that occurred during shutdown
 */
func (t *Shutdownable) Shutdown(ctx context.Context) error {
	t.mutex.Lock()

	if t.isShutdown || t.isShuttingDown {
		t.mutex.Unlock()
		return nil // Already shutdown or shutting down
	}

	t.isShuttingDown = true
	callbacks := make(map[string]types.ShutdownCallback)
	for name, callback := range t.callbacks {
		callbacks[name] = callback
	}
	t.mutex.Unlock()

	// Execute shutdown callbacks
	for _, callback := range callbacks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := callback(ctx); err != nil {
				// Log error but continue with shutdown
				// In a real implementation, you might want to log this
				continue
			}
		}
	}

	t.mutex.Lock()
	t.isShuttingDown = false
	t.isShutdown = true
	t.mutex.Unlock()

	return nil
}

/**
 * GracefulShutdown performs a graceful shutdown with timeout.
 *
 * @param timeout time.Duration The maximum time to wait for shutdown
 * @return error Any error that occurred during shutdown
 */
func (t *Shutdownable) GracefulShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return t.Shutdown(ctx)
}

/**
 * ForceShutdown performs an immediate shutdown without waiting for callbacks.
 */
func (t *Shutdownable) ForceShutdown() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.isShuttingDown = false
	t.isShutdown = true
}

/**
 * IsShuttingDown returns whether the application is currently shutting down.
 *
 * @return bool true if shutdown is in progress
 */
func (t *Shutdownable) IsShuttingDown() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.isShuttingDown
}

/**
 * SetShuttingDown sets the shutting down state.
 *
 * @param shutting bool Whether the application is shutting down
 */
func (t *Shutdownable) SetShuttingDown(shutting bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.isShuttingDown = shutting
}

/**
 * IsShutdown returns whether the application has completed shutdown.
 *
 * @return bool true if shutdown is complete
 */
func (t *Shutdownable) IsShutdown() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.isShutdown
}

/**
 * SetShutdown sets the shutdown completion state.
 *
 * @param shutdown bool Whether shutdown is complete
 */
func (t *Shutdownable) SetShutdown(shutdown bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.isShutdown = shutdown
	if shutdown {
		t.isShuttingDown = false
	}
}

/**
 * HandleSignals sets up signal handling for graceful shutdown.
 *
 * @param signals []os.Signal The signals to handle
 * @param timeout time.Duration The timeout for graceful shutdown
 */
func (t *Shutdownable) HandleSignals(signals []os.Signal, timeout time.Duration) {
	t.mutex.Lock()
	t.shutdownTimeout = timeout
	t.mutex.Unlock()

	signal.Notify(t.signalChannel, signals...)

	go func() {
		<-t.signalChannel
		t.GracefulShutdown(timeout)
	}()
}

/**
 * GetShutdownTimeout returns the configured shutdown timeout.
 *
 * @return time.Duration The shutdown timeout
 */
func (t *Shutdownable) GetShutdownTimeout() time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.shutdownTimeout
}

/**
 * SetShutdownTimeout sets the shutdown timeout.
 *
 * @param timeout time.Duration The timeout to set
 */
func (t *Shutdownable) SetShutdownTimeout(timeout time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.shutdownTimeout = timeout
}

/**
 * GetShutdownInfo returns comprehensive shutdown information.
 *
 * @return map[string]interface{} Shutdown details
 */
func (t *Shutdownable) GetShutdownInfo() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	callbackNames := make([]string, 0, len(t.callbacks))
	for name := range t.callbacks {
		callbackNames = append(callbackNames, name)
	}

	return map[string]interface{}{
		"is_shutting_down": t.isShuttingDown,
		"is_shutdown":      t.isShutdown,
		"shutdown_timeout": t.shutdownTimeout,
		"callback_count":   len(t.callbacks),
		"callback_names":   callbackNames,
	}
}

/**
 * ClearCallbacks removes all registered shutdown callbacks.
 */
func (t *Shutdownable) ClearCallbacks() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.callbacks = make(map[string]types.ShutdownCallback)
}

/**
 * GetCallbackCount returns the number of registered shutdown callbacks.
 *
 * @return int The number of registered callbacks
 */
func (t *Shutdownable) GetCallbackCount() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return len(t.callbacks)
}

// Compile-time interface compliance check
var _ traitInterfaces.ShutdownableInterface = (*Shutdownable)(nil)
