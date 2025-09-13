package interfaces

import applicationInterfaces "govel/packages/types/src/interfaces/application"

// ServiceProviderInterface defines the contract that all service providers must implement.
// This interface follows Laravel's ServiceProvider pattern, providing methods for
// registration and booting of services within the application.
//
// Service providers are the central place for all application bootstrapping.
// Your own application, as well as all core framework services, are bootstrapped
// via service providers.
type ServiceProviderInterface interface {
	// Register any application services.
	// This method is called during the service registration phase and should be used
	// to bind services into the service container. You should never attempt to use
	// any services in the register method since the service you are trying to use
	// may not have been registered yet.
	//
	// Parameters:
	//   app: The application instance for service registration
	//
	// Returns:
	//   error: Any error that occurred during registration
	Register(app applicationInterfaces.ApplicationInterface) error

	// Boot any application services.
	// This method is called after all other service providers have been registered,
	// meaning you have access to all other services that have been registered.
	// This is where you should place bootstrap logic that depends on other services.
	//
	// Parameters:
	//   app: The application instance for service booting
	//
	// Returns:
	//   error: Any error that occurred during booting
	Boot(app applicationInterfaces.ApplicationInterface) error

	// GetProvides returns the services provided by the provider.
	// This method is used to determine which services this provider offers,
	// particularly useful for deferred service providers.
	//
	// Returns:
	//   []string: A slice of service identifiers provided by this provider
	GetProvides() []string

	Priority() int
}
