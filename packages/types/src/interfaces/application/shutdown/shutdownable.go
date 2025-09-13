package interfaces

import (
	"context"
	"os"
	"time"

	shared "govel/packages/types/src/shared"
)

// ShutdownableInterface defines the contract for shutdownable functionality.
// This interface provides comprehensive shutdown management capabilities including
// graceful shutdown, force shutdown, callback management, and signal handling.
type ShutdownableInterface interface {
	// RegisterShutdownCallback registers a callback to be executed during shutdown.
	//
	// Parameters:
	//   name: The name of the callback for identification
	//   callback: The callback function to register
	RegisterShutdownCallback(name string, callback shared.ShutdownCallback)

	// UnregisterShutdownCallback removes a registered shutdown callback.
	//
	// Parameters:
	//   name: The name of the callback to unregister
	//
	// Returns:
	//   bool: true if the callback was found and removed
	UnregisterShutdownCallback(name string) bool

	// GetShutdownCallbacks returns all registered shutdown callbacks.
	//
	// Returns:
	//   map[string]shared.ShutdownCallback: Map of callback names to callbacks
	GetShutdownCallbacks() map[string]shared.ShutdownCallback

	// Shutdown initiates the application shutdown process.
	//
	// Parameters:
	//   ctx: Context for controlling the shutdown process
	//
	// Returns:
	//   error: Any error that occurred during shutdown
	Shutdown(ctx context.Context) error

	// GracefulShutdown performs a graceful shutdown with timeout.
	//
	// Parameters:
	//   timeout: The maximum time to wait for shutdown
	//
	// Returns:
	//   error: Any error that occurred during shutdown
	GracefulShutdown(timeout time.Duration) error

	// ForceShutdown performs an immediate shutdown without waiting for callbacks.
	ForceShutdown()

	// IsShuttingDown returns whether the application is currently shutting down.
	//
	// Returns:
	//   bool: true if shutdown is in progress, false otherwise
	IsShuttingDown() bool

	// SetShuttingDown sets the shutting down state.
	//
	// Parameters:
	//   shutting: Whether the application is shutting down
	SetShuttingDown(shutting bool)

	// IsShutdown returns whether the application has completed shutdown.
	//
	// Returns:
	//   bool: true if shutdown is complete, false otherwise
	IsShutdown() bool

	// SetShutdown sets the shutdown completion state.
	//
	// Parameters:
	//   shutdown: Whether shutdown is complete
	SetShutdown(shutdown bool)

	// HandleSignals sets up signal handling for graceful shutdown.
	//
	// Parameters:
	//   signals: The signals to handle
	//   timeout: The timeout for graceful shutdown
	HandleSignals(signals []os.Signal, timeout time.Duration)

	// GetShutdownTimeout returns the configured shutdown timeout.
	//
	// Returns:
	//   time.Duration: The shutdown timeout
	GetShutdownTimeout() time.Duration

	// SetShutdownTimeout sets the shutdown timeout.
	//
	// Parameters:
	//   timeout: The shutdown timeout duration
	SetShutdownTimeout(timeout time.Duration)

	// GetShutdownInfo returns comprehensive shutdown information.
	//
	// Returns:
	//   map[string]interface{}: Shutdown status and details
	GetShutdownInfo() map[string]interface{}
}
