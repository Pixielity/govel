package interfaces

import (
	"context"
	"os"
	"time"

	shared "govel/packages/types/src/shared"
)

// ShutdownableInterface defines the contract for graceful shutdown functionality.
type ShutdownableInterface interface {
	// RegisterShutdownCallback registers a callback to be executed during shutdown
	RegisterShutdownCallback(name string, callback shared.ShutdownCallback)
	
	// UnregisterShutdownCallback removes a registered shutdown callback
	UnregisterShutdownCallback(name string) bool
	
	// GetShutdownCallbacks returns all registered shutdown callbacks
	GetShutdownCallbacks() map[string]shared.ShutdownCallback
	
	// Shutdown initiates the application shutdown process
	Shutdown(ctx context.Context) error
	
	// GracefulShutdown performs a graceful shutdown with timeout
	GracefulShutdown(timeout time.Duration) error
	
	// ForceShutdown performs an immediate shutdown without waiting for callbacks
	ForceShutdown()
	
	// IsShuttingDown returns whether the application is currently shutting down
	IsShuttingDown() bool
	
	// SetShuttingDown sets the shutting down state
	SetShuttingDown(shutting bool)
	
	// IsShutdown returns whether the application has completed shutdown
	IsShutdown() bool
	
	// SetShutdown sets the shutdown completion state
	SetShutdown(shutdown bool)
	
	// HandleSignals sets up signal handling for graceful shutdown
	HandleSignals(signals []os.Signal, timeout time.Duration)
	
	// GetShutdownTimeout returns the configured shutdown timeout
	GetShutdownTimeout() time.Duration
	
	// SetShutdownTimeout sets the shutdown timeout
	SetShutdownTimeout(timeout time.Duration)
	
	// GetShutdownInfo returns comprehensive shutdown information
	GetShutdownInfo() map[string]interface{}
}
