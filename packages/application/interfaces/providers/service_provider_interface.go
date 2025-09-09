package interfaces

import (
	applicationInterfaces "govel/packages/application/interfaces/application"
)

// ServiceProviderInterface defines the contract that all service providers must implement.
// This interface follows Laravel's service provider pattern, providing methods for
// service registration and bootstrapping.
//
// Service providers should implement this interface to participate in the application
// bootstrap process. The Register method is called first for all providers, followed
// by the Boot method after all providers have been registered.
type ServiceProviderInterface interface {
	// Register is called during the service registration phase.
	// Within this method, you should only bind things into the service container.
	// You should never attempt to register event listeners, routes, or any other
	// piece of functionality within the register method.
	//
	// Parameters:
	//   application: The application instance for service container access
	//
	// Returns:
	//   error: Any error that occurred during registration, nil if successful
	Register(application applicationInterfaces.ApplicationInterface) error

	// Boot is called after all service providers have been registered.
	// This method is called with all other services registered by the framework,
	// meaning you have access to all other services that have been registered.
	//
	// Use this method to:
	// - Register event listeners
	// - Register routes
	// - Register view composers
	// - Configure other services that depend on registered services
	//
	// Parameters:
	//   application: The application instance
	//
	// Returns:
	//   error: Any error that occurred during boot, nil if successful
	Boot(application applicationInterfaces.ApplicationInterface) error

	// Priority returns the registration priority for this service provider.
	// Lower values are registered first. This allows for proper dependency ordering.
	//
	// Priority levels (suggested convention):
	// - 0-99: Core infrastructure providers
	// - 100-199: Framework service providers
	// - 200-299: Application service providers
	// - 300+: Optional/plugin providers
	//
	// Returns:
	//   int: Priority level for registration ordering
	Priority() int
}
