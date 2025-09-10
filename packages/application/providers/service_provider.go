// Package providers provides the base service provider functionality for the GoVel framework.
// This package follows Laravel's service provider pattern, providing a foundation for
// all application service registration and bootstrapping.
package providers

import (
	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	"govel/packages/application/types"
)

// ServiceProvider provides a base implementation for service providers.
// This struct follows Laravel's abstract ServiceProvider pattern, providing
// common functionality that can be embedded in concrete service providers.
//
// Service providers are the central place for all application bootstrapping.
// Your own application, as well as all core framework services, are bootstrapped
// via service providers.
type ServiceProvider struct {
	// app holds the application instance
	app applicationInterfaces.ApplicationInterface

	// bootingCallbacks holds callbacks to run before the Boot method
	bootingCallbacks []types.ProviderCallback

	// bootedCallbacks holds callbacks to run after the Boot method
	bootedCallbacks []types.ProviderCallback

	// provides holds the services provided by this provider (for deferred providers)
	provides []string

	// deferred indicates whether this provider is deferred
	deferred bool
}

// NewServiceProvider creates a new base service provider instance.
//
// Parameters:
//
//	app: The application instance
//
// Returns:
//
//	*ServiceProvider: A new service provider instance
func NewServiceProvider(app applicationInterfaces.ApplicationInterface) *ServiceProvider {
	return &ServiceProvider{
		app:              app,
		bootingCallbacks: make([]types.ProviderCallback, 0),
		bootedCallbacks:  make([]types.ProviderCallback, 0),
		provides:         make([]string, 0),
		deferred:         false,
	}
}

// GetApp returns the application instance.
//
// Returns:
//
//	ApplicationInterface: The application instance
func (sp *ServiceProvider) GetApp() applicationInterfaces.ApplicationInterface {
	return sp.app
}

// Register any application services.
// This is the base implementation and should be overridden by concrete providers.
// The register method is used to bind services into the service container.
//
// Parameters:
//
//	app: The application instance for service registration
//
// Returns:
//
//	error: Any error that occurred during registration (base implementation returns nil)
func (sp *ServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	// Base implementation does nothing - override in concrete providers
	return nil
}

// Boot any application services.
// This is the base implementation and should be overridden by concrete providers.
// The boot method is called after all providers have been registered.
//
// Parameters:
//
//	app: The application instance for service booting
//
// Returns:
//
//	error: Any error that occurred during booting (base implementation returns nil)
func (sp *ServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	// Base implementation does nothing - override in concrete providers
	return nil
}

// GetProvides returns the services provided by the provider.
//
// Returns:
//
//	[]string: A slice of service identifiers provided by this provider
func (sp *ServiceProvider) GetProvides() []string {
	return sp.provides
}

// SetProvides sets the services provided by the provider.
// This is typically used in the constructor of concrete providers.
//
// Parameters:
//
//	provides: A slice of service identifiers that this provider can resolve
func (sp *ServiceProvider) SetProvides(provides []string) {
	sp.provides = provides
}

// IsProviderDeferred is a helper function to determine if any provider is deferred.
// This checks if the provider implements the DeferrableProvider interface.
//
// Parameters:
//
//	provider: The provider instance to check
//
// Returns:
//
//	bool: true if the provider implements DeferrableProvider and is deferred, false otherwise
func IsProviderDeferred(provider interface{}) bool {
	// Check if the provider implements DeferrableProvider interface
	if deferrable, ok := provider.(providerInterfaces.DeferrableProvider); ok {
		return deferrable.IsDeferred()
	}
	return false
}

// SetDeferred sets whether the provider should be deferred.
//
// Parameters:
//
//	deferred: true if the provider should be deferred, false otherwise
func (sp *ServiceProvider) SetDeferred(deferred bool) {
	sp.deferred = deferred
}

// Booting registers a callback to be run before the Boot method is called.
// This allows providers to register pre-boot logic.
//
// Parameters:
//
//	callback: The callback function to run before booting
func (sp *ServiceProvider) Booting(callback types.ProviderCallback) {
	sp.bootingCallbacks = append(sp.bootingCallbacks, callback)
}

// Booted registers a callback to be run after the Boot method is called.
// This allows providers to register post-boot logic.
//
// Parameters:
//
//	callback: The callback function to run after booting
func (sp *ServiceProvider) Booted(callback types.ProviderCallback) {
	sp.bootedCallbacks = append(sp.bootedCallbacks, callback)
}

// CallBootingCallbacks executes all registered booting callbacks.
// This should be called before the Boot method.
//
// Returns:
//
//	error: Any error that occurred during callback execution
func (sp *ServiceProvider) CallBootingCallbacks() error {
	for _, callback := range sp.bootingCallbacks {
		// Cast app back to interface{} for the callback
		if err := callback(interface{}(sp.app)); err != nil {
			return err
		}
	}
	return nil
}

// CallBootedCallbacks executes all registered booted callbacks.
// This should be called after the Boot method.
//
// Returns:
//
//	error: Any error that occurred during callback execution
func (sp *ServiceProvider) CallBootedCallbacks() error {
	for _, callback := range sp.bootedCallbacks {
		// Cast app back to interface{} for the callback
		if err := callback(interface{}(sp.app)); err != nil {
			return err
		}
	}
	return nil
}

// Compile-time interface compliance check
var _ providerInterfaces.ServiceProviderInterface = (*ServiceProvider)(nil)
