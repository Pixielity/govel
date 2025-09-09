package interfaces

import "context"

/**
 * LifecycleableInterface defines the contract for components that provide
 * application lifecycle management functionality. This interface follows the
 * Interface Segregation Principle by focusing solely on lifecycle operations.
 */
type LifecycleableInterface interface {
	/**
	 * Boot initializes the application and its components.
	 *
	 * @param ctx context.Context The context for the boot operation
	 * @return error Any error that occurred during boot
	 */
	Boot(ctx context.Context) error

	/**
	 * IsBooted returns whether the application has been booted.
	 *
	 * @return bool true if the application is booted
	 */
	IsBooted() bool

	/**
	 * SetBooted sets the booted state of the application.
	 *
	 * @param booted bool Whether the application is booted
	 */
	SetBooted(booted bool)

	/**
	 * Start starts the application after booting.
	 *
	 * @param ctx context.Context The context for the start operation
	 * @return error Any error that occurred during start
	 */
	Start(ctx context.Context) error

	/**
	 * IsStarted returns whether the application has been started.
	 *
	 * @return bool true if the application is started
	 */
	IsStarted() bool

	/**
	 * SetStarted sets the started state of the application.
	 *
	 * @param started bool Whether the application is started
	 */
	SetStarted(started bool)

	/**
	 * Stop stops the application gracefully.
	 *
	 * @param ctx context.Context The context for the stop operation
	 * @return error Any error that occurred during stop
	 */
	Stop(ctx context.Context) error

	/**
	 * IsStopped returns whether the application has been stopped.
	 *
	 * @return bool true if the application is stopped
	 */
	IsStopped() bool

	/**
	 * SetStopped sets the stopped state of the application.
	 *
	 * @param stopped bool Whether the application is stopped
	 */
	SetStopped(stopped bool)

	/**
	 * Restart restarts the application (stop then start).
	 *
	 * @param ctx context.Context The context for the restart operation
	 * @return error Any error that occurred during restart
	 */
	Restart(ctx context.Context) error

	/**
	 * GetState returns the current lifecycle state of the application.
	 *
	 * @return string The current lifecycle state
	 */
	GetState() string

	/**
	 * IsState checks if the application is in the specified state.
	 *
	 * @param state string The state to check against
	 * @return bool true if the application is in the specified state
	 */
	IsState(state string) bool

	/**
	 * GetLifecycleInfo returns comprehensive lifecycle information.
	 *
	 * @return map[string]interface{} Lifecycle details
	 */
	GetLifecycleInfo() map[string]interface{}
}
