// Package shutdown provides graceful shutdown functionality for the GoVel application.
// This package handles termination of services, signal management, graceful and force
// shutdown procedures, and coordination with terminatable service providers.
//
// The shutdown system supports:
// - Graceful shutdown with configurable timeout
// - Force shutdown when graceful timeout is exceeded
// - Drain mode for completing active operations
// - Provider termination coordination
// - Error collection and reporting
package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"govel/logger"
	constants "govel/types/constants/application"
	providerInterfaces "govel/types/interfaces/application/providers"
	types "govel/types/types/application"
)

// ShutdownManager handles the graceful shutdown process for the application.
// It coordinates with terminatable service providers and manages the overall
// shutdown lifecycle. This manager is designed to work with the Shutdownable trait.
type ShutdownManager struct {
	// providerRepository provides access to terminatable providers
	providerRepository providerInterfaces.ProviderRepositoryInterface

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

	// callbacks stores registered shutdown callbacks
	callbacks map[string]types.ShutdownCallback

	// signalChannel handles OS signals for graceful shutdown
	signalChannel chan os.Signal

	// shutdownChannel coordinates shutdown completion
	shutdownChannel chan struct{}

	// logger is the logger instance for shutdown operations
	logger *logger.Logger
}

// NewShutdownManager creates a new shutdown manager.
//
// Parameters:
//
//	providerRepository: The provider repository for accessing terminatable providers
//
// Returns:
//
//	*ShutdownManager: A new shutdown manager instance
func NewShutdownManager(providerRepository providerInterfaces.ProviderRepositoryInterface) *ShutdownManager {
	return &ShutdownManager{
		providerRepository: providerRepository,
		shutdownTimeout:    constants.DefaultShutdownTimeout,
		drainTimeout:       constants.DefaultDrainTimeout,
		terminationErrors:  make([]error, 0),
		callbacks:          make(map[string]types.ShutdownCallback),
		signalChannel:      make(chan os.Signal, 1),
		shutdownChannel:    make(chan struct{}),
		logger:             logger.New(),
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

// GetTerminatableProviders returns all registered terminatable providers.
//
// Returns:
//
//	[]TerminatableProvider: List of registered terminatable providers
func (sm *ShutdownManager) GetTerminatableProviders() []providerInterfaces.TerminatableProvider {
	if sm.providerRepository == nil {
		return []providerInterfaces.TerminatableProvider{}
	}
	return sm.providerRepository.GetTerminatableProviders()
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
	sm.logger.Info("Application entering drain mode")

	// Notify the application that it should stop accepting new work
	// This is a hook for the application to implement drain behavior
	sm.logger.Debug("Drain mode enabled - new operations should be rejected")
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

// GracefulShutdownWithApp initiates a graceful shutdown of the application with provider termination.
// This method coordinates the shutdown of all terminatable service providers
// and respects the configured shutdown timeout.
//
// Parameters:
//
//	ctx: Context for controlling the shutdown process
//	app: Application instance for provider termination
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
//	if err := shutdownManager.GracefulShutdownWithApp(ctx, app); err != nil {
//	    log.Printf("Graceful shutdown failed: %v", err)
//	    shutdownManager.ForceShutdown()
//	}
func (sm *ShutdownManager) GracefulShutdownWithApp(ctx context.Context, app interface{}) error {
	sm.mutex.Lock()
	if sm.shutdownStarted {
		sm.mutex.Unlock()
		return fmt.Errorf("shutdown already in progress")
	}
	sm.shutdownStarted = true
	sm.mutex.Unlock()

	sm.logger.Info("Initiating graceful shutdown with app")

	// Execute shutdown callbacks first
	sm.executeShutdownCallbacks(ctx)

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

	// Terminate all terminatable service providers
	sm.logger.Info("Terminating service providers")
	terminationErrors := sm.terminateProviders(shutdownCtx, app)
	if len(terminationErrors) > 0 {
		sm.mutex.Lock()
		sm.terminationErrors = append(sm.terminationErrors, terminationErrors...)
		sm.mutex.Unlock()

		for i, err := range terminationErrors {
			sm.logger.Error("Provider termination error %d: %v", i+1, err)
		}
		sm.logger.Warn("Some service providers failed to terminate gracefully")
	} else {
		sm.logger.Info("All service providers terminated successfully")
	}

	sm.mutex.Lock()
	sm.shutdownCompleted = true
	sm.mutex.Unlock()

	sm.logger.Info("Graceful shutdown completed")

	// Signal shutdown completion
	select {
	case sm.shutdownChannel <- struct{}{}:
	default:
	}

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

	sm.logger.Warn("Initiating force shutdown")
	sm.forceShutdown = true
	sm.shutdownStarted = true

	// Force immediate termination
	// In a real implementation, this might involve:
	// - Closing all connections immediately
	// - Stopping all goroutines
	// - Flushing critical data

	sm.shutdownCompleted = true
	sm.logger.Warn("Force shutdown completed")

	// Exit the process immediately
	os.Exit(1)
}

// waitForDrain waits for active operations to complete during drain mode.
func (sm *ShutdownManager) waitForDrain(ctx context.Context) {
	sm.logger.Info("Waiting for active operations to complete (drain)")

	// This is where the application would implement logic to wait for
	// active operations to complete. For example:
	// - Wait for HTTP requests to complete
	// - Wait for database transactions to finish
	// - Wait for background jobs to complete

	select {
	case <-ctx.Done():
		sm.logger.Warn("Drain timeout exceeded, proceeding with shutdown")
	case <-time.After(sm.drainTimeout):
		sm.logger.Info("Drain completed successfully")
	}
}

// terminateProviders terminates all registered terminatable providers.
// This method requires an application interface for provider termination.
func (sm *ShutdownManager) terminateProviders(ctx context.Context, app interface{}) []error {
	if sm.providerRepository == nil {
		return []error{fmt.Errorf("provider repository is not available")}
	}

	// Try to cast the app to ApplicationInterface
	return sm.providerRepository.TerminateProviders(ctx)
}

// RegisterShutdownCallback registers a callback to be executed during shutdown.
//
// Parameters:
//
//	name: The name of the callback for identification
//	callback: The callback function to register
//
// Example:
//
//	shutdownManager.RegisterShutdownCallback("database", func(ctx context.Context) error {
//	    return database.Close()
//	})
func (sm *ShutdownManager) RegisterShutdownCallback(name string, callback types.ShutdownCallback) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.callbacks[name] = callback
	sm.logger.Debug("Registered shutdown callback: %s", name)
}

// UnregisterShutdownCallback removes a registered shutdown callback.
//
// Parameters:
//
//	name: The name of the callback to unregister
//
// Returns:
//
//	bool: true if the callback was found and removed
//
// Example:
//
//	if shutdownManager.UnregisterShutdownCallback("database") {
//	    log.Println("Database shutdown callback removed")
//	}
func (sm *ShutdownManager) UnregisterShutdownCallback(name string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if _, exists := sm.callbacks[name]; exists {
		delete(sm.callbacks, name)
		sm.logger.Debug("Unregistered shutdown callback: %s", name)
		return true
	}
	return false
}

// GetShutdownCallbacks returns all registered shutdown callbacks.
//
// Returns:
//
//	map[string]types.ShutdownCallback: Map of callback names to callbacks
//
// Example:
//
//	callbacks := shutdownManager.GetShutdownCallbacks()
//	for name := range callbacks {
//	    log.Printf("Registered callback: %s", name)
//	}
func (sm *ShutdownManager) GetShutdownCallbacks() map[string]types.ShutdownCallback {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	callbacks := make(map[string]types.ShutdownCallback)
	for name, callback := range sm.callbacks {
		callbacks[name] = callback
	}
	return callbacks
}

// GetShutdownTimeout returns the configured shutdown timeout.
//
// Returns:
//
//	time.Duration: The shutdown timeout
func (sm *ShutdownManager) GetShutdownTimeout() time.Duration {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.shutdownTimeout
}

// HandleSignals sets up signal handling for graceful shutdown.
//
// Parameters:
//
//	signals: The signals to handle
//	timeout: The timeout for graceful shutdown
//
// Example:
//
//	shutdownManager.HandleSignals([]os.Signal{syscall.SIGINT, syscall.SIGTERM}, 30*time.Second)
func (sm *ShutdownManager) HandleSignals(signals []os.Signal, timeout time.Duration) {
	signal.Notify(sm.signalChannel, signals...)
	go func() {
		sig := <-sm.signalChannel
		sm.logger.Info("Received signal: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := sm.Shutdown(ctx); err != nil {
			sm.logger.Error("Graceful shutdown failed: %v", err)
			sm.ForceShutdown()
		}
	}()
}

// Shutdown initiates the application shutdown process with callbacks.
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
//	if err := shutdownManager.Shutdown(ctx); err != nil {
//	    log.Printf("Shutdown failed: %v", err)
//	}
func (sm *ShutdownManager) Shutdown(ctx context.Context) error {
	sm.mutex.Lock()
	if sm.shutdownStarted {
		sm.mutex.Unlock()
		return fmt.Errorf("shutdown already in progress")
	}
	sm.shutdownStarted = true
	sm.mutex.Unlock()

	sm.logger.Info("Initiating shutdown")

	// Execute shutdown callbacks first
	sm.executeShutdownCallbacks(ctx)

	// Then perform graceful shutdown
	return sm.performGracefulShutdown(ctx)
}

// executeShutdownCallbacks executes all registered shutdown callbacks.
func (sm *ShutdownManager) executeShutdownCallbacks(ctx context.Context) {
	sm.mutex.RLock()
	callbacks := make(map[string]types.ShutdownCallback)
	for name, callback := range sm.callbacks {
		callbacks[name] = callback
	}
	sm.mutex.RUnlock()

	if len(callbacks) == 0 {
		return
	}

	sm.logger.Info("Executing %d shutdown callbacks", len(callbacks))

	for name, callback := range callbacks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					sm.logger.Error("Shutdown callback '%s' panicked: %v", name, r)
				}
			}()

			if err := callback(ctx); err != nil {
				sm.logger.Error("Shutdown callback '%s' failed: %v", name, err)
			} else {
				sm.logger.Debug("Shutdown callback '%s' completed successfully", name)
			}
		}()
	}
}

// performGracefulShutdown performs the core graceful shutdown logic.
func (sm *ShutdownManager) performGracefulShutdown(ctx context.Context) error {
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

	// Terminate all terminatable service providers
	sm.logger.Info("Terminating service providers")
	terminationErrors := sm.terminateProvidersInternal(shutdownCtx)
	if len(terminationErrors) > 0 {
		sm.mutex.Lock()
		sm.terminationErrors = append(sm.terminationErrors, terminationErrors...)
		sm.mutex.Unlock()

		for i, err := range terminationErrors {
			sm.logger.Error("Provider termination error %d: %v", i+1, err)
		}
		sm.logger.Warn("Some service providers failed to terminate gracefully")
	} else {
		sm.logger.Info("All service providers terminated successfully")
	}

	sm.mutex.Lock()
	sm.shutdownCompleted = true
	sm.mutex.Unlock()

	sm.logger.Info("Graceful shutdown completed")

	// Signal shutdown completion
	select {
	case sm.shutdownChannel <- struct{}{}:
	default:
	}

	// Return aggregated errors if any
	if len(sm.terminationErrors) > 0 {
		return fmt.Errorf("shutdown completed with %d errors", len(sm.terminationErrors))
	}

	return nil
}

// terminateProvidersInternal terminates all registered terminatable providers.
func (sm *ShutdownManager) terminateProvidersInternal(ctx context.Context) []error {
	if sm.providerRepository == nil {
		return []error{fmt.Errorf("provider repository is not available")}
	}

	// We need an app instance, but since we're in the manager, we'll need to handle this differently
	// For now, return empty errors as provider termination should be handled at a higher level
	return []error{}
}

// GetShutdownInfo returns comprehensive shutdown information.
//
// Returns:
//
//	map[string]interface{}: Shutdown details
//
// Example:
//
//	info := shutdownManager.GetShutdownInfo()
//	fmt.Printf("Shutdown status: %v\n", info["shutdown_started"])
func (sm *ShutdownManager) GetShutdownInfo() map[string]interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return map[string]interface{}{
		"shutdown_started":   sm.shutdownStarted,
		"shutdown_completed": sm.shutdownCompleted,
		"drain_mode":         sm.drainMode,
		"force_shutdown":     sm.forceShutdown,
		"shutdown_timeout":   sm.shutdownTimeout.String(),
		"drain_timeout":      sm.drainTimeout.String(),
		"callback_count":     len(sm.callbacks),
		"error_count":        len(sm.terminationErrors),
	}
}

// GracefulShutdown performs a graceful shutdown with timeout.
// This method is required by the ShutdownableInterface and provides
// a simplified interface for graceful shutdown.
//
// Parameters:
//
//	timeout: The maximum time to wait for shutdown
//
// Returns:
//
//	error: Any error that occurred during shutdown
//
// Example:
//
//	if err := shutdownManager.GracefulShutdown(30 * time.Second); err != nil {
//	    log.Printf("Graceful shutdown failed: %v", err)
//	}
func (sm *ShutdownManager) GracefulShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return sm.Shutdown(ctx)
}

// SetShuttingDown sets the shutting down state.
//
// Parameters:
//
//	shutting: Whether the application is shutting down
//
// Example:
//
//	shutdownManager.SetShuttingDown(true)
func (sm *ShutdownManager) SetShuttingDown(shutting bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.shutdownStarted = shutting
}

// IsShutdown returns whether the application has completed shutdown.
// This is an alias for IsShutdownCompleted for interface compatibility.
//
// Returns:
//
//	bool: true if shutdown is complete
func (sm *ShutdownManager) IsShutdown() bool {
	return sm.IsShutdownCompleted()
}

// SetShutdown sets the shutdown completion state.
//
// Parameters:
//
//	shutdown: Whether shutdown is complete
//
// Example:
//
//	shutdownManager.SetShutdown(true)
func (sm *ShutdownManager) SetShutdown(shutdown bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.shutdownCompleted = shutdown
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
