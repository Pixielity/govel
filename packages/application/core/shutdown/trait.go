// Package shutdown provides graceful shutdown functionality for the GoVel application.
package shutdown

import (
	"context"
	"os"
	"time"

	"govel/types/src/types/application"
	applicationInterfaces "govel/types/src/interfaces/application"
	providerInterfaces "govel/types/src/interfaces/application/providers"
)

// Shutdownable provides shutdown functionality using the ShutdownManager.
// This trait implements the ShutdownableInterface and serves as a lightweight
// wrapper around the ShutdownManager, delegating all operations to the manager.
// All shutdown functionality is handled by the manager to maintain a single
// source of truth for shutdown state and operations.
type Shutdownable struct {
	// manager handles all shutdown orchestration and functionality
	manager *ShutdownManager
}

// NewShutdownable creates a new shutdownable trait with the provided manager.
//
// Parameters:
//
//	providerRepository: The provider repository for accessing terminatable providers
//
// Returns:
//
//	*Shutdownable: A new shutdownable trait instance
//
// Example:
//
//	shutdownable := NewShutdownable(providerRepository)
func NewShutdownable(providerRepository providerInterfaces.ProviderRepositoryInterface) *Shutdownable {
	return &Shutdownable{
		manager: NewShutdownManager(providerRepository),
	}
}

// RegisterShutdownCallback registers a callback to be executed during shutdown.
// Delegates to the manager.
func (s *Shutdownable) RegisterShutdownCallback(name string, callback types.ShutdownCallback) {
	s.manager.RegisterShutdownCallback(name, callback)
}

// UnregisterShutdownCallback removes a registered shutdown callback.
// Delegates to the manager.
func (s *Shutdownable) UnregisterShutdownCallback(name string) bool {
	return s.manager.UnregisterShutdownCallback(name)
}

// GetShutdownCallbacks returns all registered shutdown callbacks.
// Delegates to the manager.
func (s *Shutdownable) GetShutdownCallbacks() map[string]types.ShutdownCallback {
	return s.manager.GetShutdownCallbacks()
}

// Shutdown initiates the application shutdown process.
// Delegates to the manager.
func (s *Shutdownable) Shutdown(ctx context.Context) error {
	return s.manager.Shutdown(ctx)
}

// GracefulShutdown performs a graceful shutdown with timeout.
// Delegates to the manager.
func (s *Shutdownable) GracefulShutdown(timeout time.Duration) error {
	return s.manager.GracefulShutdown(timeout)
}

// ForceShutdown performs an immediate shutdown without waiting for callbacks.
// Delegates to the manager.
func (s *Shutdownable) ForceShutdown() {
	s.manager.ForceShutdown()
}

// IsShuttingDown returns whether the application is currently shutting down.
// Delegates to the manager.
func (s *Shutdownable) IsShuttingDown() bool {
	return s.manager.IsShuttingDown()
}

// SetShuttingDown sets the shutting down state.
// Delegates to the manager.
func (s *Shutdownable) SetShuttingDown(shutting bool) {
	s.manager.SetShuttingDown(shutting)
}

// IsShutdown returns whether the application has completed shutdown.
// Delegates to the manager.
func (s *Shutdownable) IsShutdown() bool {
	return s.manager.IsShutdown()
}

// SetShutdown sets the shutdown completion state.
// Delegates to the manager.
func (s *Shutdownable) SetShutdown(shutdown bool) {
	s.manager.SetShutdown(shutdown)
}

// HandleSignals sets up signal handling for graceful shutdown.
// Delegates to the manager.
func (s *Shutdownable) HandleSignals(signals []os.Signal, timeout time.Duration) {
	s.manager.HandleSignals(signals, timeout)
}

// GetShutdownTimeout returns the configured shutdown timeout.
// Delegates to the manager.
func (s *Shutdownable) GetShutdownTimeout() time.Duration {
	return s.manager.GetShutdownTimeout()
}

// SetShutdownTimeout sets the shutdown timeout.
// Delegates to the manager.
func (s *Shutdownable) SetShutdownTimeout(timeout time.Duration) {
	s.manager.SetShutdownTimeout(timeout)
}

// GetShutdownInfo returns comprehensive shutdown information.
// Delegates to the manager.
func (s *Shutdownable) GetShutdownInfo() map[string]interface{} {
	return s.manager.GetShutdownInfo()
}

// Compile-time interface compliance check
var _ applicationInterfaces.ShutdownableInterface = (*Shutdownable)(nil)
