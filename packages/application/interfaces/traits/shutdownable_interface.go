package interfaces

import (
	"context"
	"govel/packages/application/types"
	"os"
	"time"
)

/**
 * ShutdownableInterface defines the contract for components that provide
 * application shutdown management functionality. This interface follows the
 * Interface Segregation Principle by focusing solely on shutdown operations.
 */
type ShutdownableInterface interface {
	/**
	 * RegisterShutdownCallback registers a callback to be executed during shutdown.
	 *
	 * @param name string The name of the callback for identification
	 * @param callback ShutdownCallback The callback function to register
	 */
	RegisterShutdownCallback(name string, callback types.ShutdownCallback)

	/**
	 * UnregisterShutdownCallback removes a registered shutdown callback.
	 *
	 * @param name string The name of the callback to unregister
	 * @return bool true if the callback was found and removed
	 */
	UnregisterShutdownCallback(name string) bool

	/**
	 * GetShutdownCallbacks returns all registered shutdown callbacks.
	 *
	 * @return map[string]ShutdownCallback Map of callback names to callbacks
	 */
	GetShutdownCallbacks() map[string]types.ShutdownCallback

	/**
	 * Shutdown initiates the application shutdown process.
	 *
	 * @param ctx context.Context The context for the shutdown operation
	 * @return error Any error that occurred during shutdown
	 */
	Shutdown(ctx context.Context) error

	/**
	 * GracefulShutdown performs a graceful shutdown with timeout.
	 *
	 * @param timeout time.Duration The maximum time to wait for shutdown
	 * @return error Any error that occurred during shutdown
	 */
	GracefulShutdown(timeout time.Duration) error

	/**
	 * ForceShutdown performs an immediate shutdown without waiting for callbacks.
	 */
	ForceShutdown()

	/**
	 * IsShuttingDown returns whether the application is currently shutting down.
	 *
	 * @return bool true if shutdown is in progress
	 */
	IsShuttingDown() bool

	/**
	 * SetShuttingDown sets the shutting down state.
	 *
	 * @param shutting bool Whether the application is shutting down
	 */
	SetShuttingDown(shutting bool)

	/**
	 * IsShutdown returns whether the application has completed shutdown.
	 *
	 * @return bool true if shutdown is complete
	 */
	IsShutdown() bool

	/**
	 * SetShutdown sets the shutdown completion state.
	 *
	 * @param shutdown bool Whether shutdown is complete
	 */
	SetShutdown(shutdown bool)

	/**
	 * HandleSignals sets up signal handling for graceful shutdown.
	 *
	 * @param signals []os.Signal The signals to handle
	 * @param timeout time.Duration The timeout for graceful shutdown
	 */
	HandleSignals(signals []os.Signal, timeout time.Duration)

	/**
	 * GetShutdownTimeout returns the configured shutdown timeout.
	 *
	 * @return time.Duration The shutdown timeout
	 */
	GetShutdownTimeout() time.Duration

	/**
	 * SetShutdownTimeout sets the shutdown timeout.
	 *
	 * @param timeout time.Duration The timeout to set
	 */
	SetShutdownTimeout(timeout time.Duration)

	/**
	 * GetShutdownInfo returns comprehensive shutdown information.
	 *
	 * @return map[string]interface{} Shutdown details
	 */
	GetShutdownInfo() map[string]interface{}
}
