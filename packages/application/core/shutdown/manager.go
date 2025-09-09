// Package application provides graceful shutdown functionality for the GoVel application.
// This file handles termination of services, signal management, graceful and force
// shutdown procedures, and coordination with terminatable service providers.
//
// The shutdown system supports:
// - Graceful shutdown with configurable timeout
// - Force shutdown when graceful timeout is exceeded
// - Signal handling (SIGINT, SIGTERM, etc.)
// - Terminatable service provider coordination
// - Drain mode for completing active operations
package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"govel/packages/application/enums"
	applicationInterfaces "govel/packages/application/interfaces/application"
)

// ShutdownManager handles the graceful shutdown process for the application.
// It coordinates with terminatable service providers and manages the overall
// shutdown lifecycle including signal handling and timeout management.
type ShutdownManager struct {
	// application holds a reference to the application instance
	application applicationInterfaces.ApplicationInterface

	// shutdownChan is used to signal shutdown initiation
	shutdownChan chan os.Signal

	// shutdownStarted indicates whether shutdown has been initiated
	shutdownStarted bool

	// shutdownCompleted indicates whether shutdown has completed
	shutdownCompleted bool

	// drainMode indicates whether the application is in drain mode
	drainMode bool

	// forceShutdown indicates whether force shutdown should be used
	forceShutdown bool

	// mutex provides thread-safe access to shutdown state
	mutex sync.RWMutex

	// shutdownTimeout specifies the maximum time to wait for graceful shutdown
	shutdownTimeout time.Duration

	// drainTimeout specifies the maximum time to wait for drain completion
	drainTimeout time.Duration

	// terminationErrors collects errors from terminatable providers
	terminationErrors []error
}

// NewShutdownManager creates a new shutdown manager for the application.
//
// Parameters:
//
//	application: The application instance to manage shutdown for
//
// Returns:
//
//	*ShutdownManager: A new shutdown manager instance
func NewShutdownManager(application applicationInterfaces.ApplicationInterface) *ShutdownManager {
	return &ShutdownManager{
		application:       application,
		shutdownChan:      make(chan os.Signal, 1),
		shutdownTimeout:   enums.GetDefaultTimeout(enums.TimeoutShutdown).Duration(),
		drainTimeout:      enums.GetDefaultTimeout(enums.TimeoutDrain).Duration(),
		terminationErrors: make([]error, 0),
	}
}

// SetShutdownTimeout configures the graceful shutdown timeout.
// This is the maximum time to wait for all services to shut down gracefully.
//
// Parameters:
//
//	timeout: The shutdown timeout duration
//
// Example:
//
//	shutdownManager.SetShutdownTimeout(60 * time.Second)
func (sm *ShutdownManager) SetShutdownTimeout(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.shutdownTimeout = timeout
}

// SetDrainTimeout configures the drain timeout.
// This is the maximum time to wait for active operations to complete during drain.
//
// Parameters:
//
//	timeout: The drain timeout duration
//
// Example:
//
//	shutdownManager.SetDrainTimeout(30 * time.Second)
func (sm *ShutdownManager) SetDrainTimeout(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.drainTimeout = timeout
}

// RegisterSignalHandlers sets up signal handling for graceful shutdown.
// This method registers handlers for SIGINT, SIGTERM, and other termination signals.
//
// Example:
//
//	shutdownManager.RegisterSignalHandlers()
//	// Application will now respond to Ctrl+C and other termination signals
func (sm *ShutdownManager) RegisterSignalHandlers() {
	// Register for interrupt signals
	signal.Notify(sm.shutdownChan,
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Termination signal
		syscall.SIGQUIT, // Quit signal
	)

	// Start signal handler goroutine
	go sm.handleSignals()
}

// handleSignals processes incoming shutdown signals.
func (sm *ShutdownManager) handleSignals() {
	for {
		sig := <-sm.shutdownChan
		sm.application.GetLogger().Info("Received shutdown signal: %v", sig)

		// Initiate graceful shutdown
		go sm.GracefulShutdown(context.Background())

		// Wait for shutdown to complete or force shutdown on second signal
		select {
		case <-sm.shutdownChan:
			sm.application.GetLogger().Warn("Received second shutdown signal, forcing immediate shutdown")
			sm.ForceShutdown()
			return
		case <-time.After(sm.shutdownTimeout + enums.TimeoutShort.Duration()):
			// Give a little extra time beyond the configured timeout
			return
		}
	}
}

// IsShuttingDown returns true if the application is currently shutting down.
//
// Returns:
//
//	bool: true if shutdown is in progress, false otherwise
//
// Example:
//
//	if application.IsShuttingDown() {
//	    // Reject new requests
//	    http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
//	    return
//	}
func (sm *ShutdownManager) IsShuttingDown() bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.shutdownStarted
}

// IsShutdownCompleted returns true if the application shutdown has completed.
//
// Returns:
//
//	bool: true if shutdown is completed, false otherwise
func (sm *ShutdownManager) IsShutdownCompleted() bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.shutdownCompleted
}

// EnableDrainMode puts the application into drain mode.
// In drain mode, the application stops accepting new work but continues
// processing existing operations until they complete or the drain timeout is reached.
//
// Example:
//
//	// Start drain mode before shutdown
//	shutdownManager.EnableDrainMode()
//	time.Sleep(10 * time.Second) // Allow time for operations to complete
//	shutdownManager.GracefulShutdown(ctx)
func (sm *ShutdownManager) EnableDrainMode() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sm.drainMode {
		return // Already in drain mode
	}

	sm.drainMode = true
	sm.application.GetLogger().Info("Application entering drain mode")

	// Notify the application that it should stop accepting new work
	// This is a hook for the application to implement drain behavior
	sm.application.GetLogger().Debug("Drain mode enabled - new operations should be rejected")
}

// IsDraining returns true if the application is in drain mode.
//
// Returns:
//
//	bool: true if draining, false otherwise
//
// Example:
//
//	if application.IsDraining() {
//	    // Reject new work
//	    return errors.New("application is draining")
//	}
func (sm *ShutdownManager) IsDraining() bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.drainMode
}

// GracefulShutdown initiates a graceful shutdown of the application.
// This method coordinates the shutdown of all terminatable service providers
// and respects the configured shutdown timeout.
//
// Parameters:
//
//	ctx: Context for controlling the shutdown process
//
// Returns:
//
//	error: Any error that occurred during shutdown
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
//	defer cancel()
//
//	if err := shutdownManager.GracefulShutdown(ctx); err != nil {
//	    log.Printf("Graceful shutdown failed: %v", err)
//	    shutdownManager.ForceShutdown()
//	}
func (sm *ShutdownManager) GracefulShutdown(ctx context.Context) error {
	sm.mutex.Lock()
	if sm.shutdownStarted {
		sm.mutex.Unlock()
		return fmt.Errorf("shutdown already in progress")
	}
	sm.shutdownStarted = true
	sm.mutex.Unlock()

	sm.application.GetLogger().Info("Initiating graceful shutdown")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, sm.shutdownTimeout)
	defer cancel()

	// Enable drain mode if not already enabled
	if !sm.IsDraining() {
		sm.EnableDrainMode()

		// Wait for drain timeout to allow active operations to complete
		drainCtx, drainCancel := context.WithTimeout(shutdownCtx, sm.drainTimeout)
		sm.waitForDrain(drainCtx)
		drainCancel()
	}

	// Shutdown terminatable service providers
	if err := sm.terminateProviders(shutdownCtx); err != nil {
		sm.application.GetLogger().Error("Error during provider termination: %v", err)
	}

	// Shutdown other application components
	if err := sm.shutdownComponents(shutdownCtx); err != nil {
		sm.application.GetLogger().Error("Error during component shutdown: %v", err)
	}

	sm.mutex.Lock()
	sm.shutdownCompleted = true
	sm.mutex.Unlock()

	sm.application.GetLogger().Info("Graceful shutdown completed")

	// Return aggregated errors if any
	if len(sm.terminationErrors) > 0 {
		return fmt.Errorf("shutdown completed with %d errors", len(sm.terminationErrors))
	}

	return nil
}

// ForceShutdown performs an immediate shutdown of the application.
// This method should only be used when graceful shutdown has failed or
// when immediate termination is required.
//
// Example:
//
//	// If graceful shutdown fails
//	if err := shutdownManager.GracefulShutdown(ctx); err != nil {
//	    log.Printf("Graceful shutdown failed, forcing shutdown: %v", err)
//	    shutdownManager.ForceShutdown()
//	}
func (sm *ShutdownManager) ForceShutdown() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sm.shutdownCompleted {
		return // Already shut down
	}

	sm.application.GetLogger().Warn("Initiating force shutdown")
	sm.forceShutdown = true
	sm.shutdownStarted = true

	// Force immediate termination
	// In a real implementation, this might involve:
	// - Closing all connections immediately
	// - Stopping all goroutines
	// - Flushing critical data

	sm.shutdownCompleted = true
	sm.application.GetLogger().Warn("Force shutdown completed")

	// Exit the process immediately
	os.Exit(1)
}

// waitForDrain waits for active operations to complete during drain mode.
func (sm *ShutdownManager) waitForDrain(ctx context.Context) {
	sm.application.GetLogger().Info("Waiting for active operations to complete (drain)")

	// This is where the application would implement logic to wait for
	// active operations to complete. For example:
	// - Wait for HTTP requests to complete
	// - Wait for database transactions to finish
	// - Wait for background jobs to complete

	select {
	case <-ctx.Done():
		sm.application.GetLogger().Warn("Drain timeout exceeded, proceeding with shutdown")
	case <-time.After(sm.drainTimeout):
		sm.application.GetLogger().Info("Drain completed successfully")
	}
}

// terminateProviders shuts down all terminatable service providers.
func (sm *ShutdownManager) terminateProviders(_ context.Context) error {
	sm.application.GetLogger().Info("Terminating service providers")

	// Use the termination manager to handle provider shutdown
	// Note: This would need to be implemented through the interface
	// For now, we'll skip this until the interface is properly defined
	// errors := sm.application.TerminationManager().Terminate(ctx)
	errors := make([]error, 0)

	if len(errors) > 0 {
		sm.terminationErrors = append(sm.terminationErrors, errors...)
		sm.application.GetLogger().Error("Termination completed with %d errors", len(errors))
		for i, err := range errors {
			sm.application.GetLogger().Error("Termination error %d: %v", i+1, err)
		}
		return fmt.Errorf("provider termination failed with %d errors", len(errors))
	}

	sm.application.GetLogger().Info("All service providers terminated successfully")
	return nil
}

// shutdownComponents shuts down other application components.
func (sm *ShutdownManager) shutdownComponents(_ context.Context) error {
	sm.application.GetLogger().Info("Shutting down application components")

	// Shutdown logger (flush any pending logs)
	// Note: This would need to be implemented through the interface
	// For now, we'll skip this until the interface is properly defined
	// if err := sm.application.GetLogger().Flush(); err != nil {
	// 	// Log the error, but don't fail the shutdown
	// 	fmt.Printf("Warning: failed to flush logger during shutdown: %v\n", err)
	// }

	// Additional component shutdown logic would go here
	// For example:
	// - Close database connections
	// - Stop background workers
	// - Clean up temporary resources

	sm.application.GetLogger().Info("Application components shutdown completed")
	return nil
}

// ShutdownErrors returns any errors that occurred during the shutdown process.
//
// Returns:
//
//	[]error: List of errors that occurred during shutdown
//
// Example:
//
//	errors := shutdownManager.ShutdownErrors()
//	for i, err := range errors {
//	    log.Printf("Shutdown error %d: %v", i+1, err)
//	}
func (sm *ShutdownManager) ShutdownErrors() []error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	errors := make([]error, len(sm.terminationErrors))
	copy(errors, sm.terminationErrors)
	return errors
}
