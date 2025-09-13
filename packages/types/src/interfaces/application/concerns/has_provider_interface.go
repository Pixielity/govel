package interfaces

import (
	"context"
)

// ApplicationProviderInterface defines the contract for service provider management
// within the application. This interface follows the Interface Segregation Principle
// by focusing solely on provider lifecycle operations.
//
// Note: This interface uses interface{} for provider parameters to avoid circular
// imports with the provider interfaces package. The actual types should be cast
// to the appropriate ServiceProviderInterface when implementing.
//
// This interface provides methods for:
// - Provider registration and management
// - Provider lifecycle operations (boot/terminate)
// - Provider introspection and debugging
// - Deferred provider loading
type ApplicationProviderInterface interface {
	// RegisterProvider registers a service provider with the application.
	// The provider will be stored for later loading and automatically added
	// to terminatable providers if it implements TerminatableProvider.
	//
	// Parameters:
	//   provider: The service provider instance to register (should implement ServiceProviderInterface)
	//
	// Returns:
	//   error: Any error that occurred during registration
	RegisterProvider(provider interface{}) error

	// RegisterProviders registers multiple service providers.
	// This is a convenience method for bulk provider registration.
	//
	// Parameters:
	//   providers: A slice of service provider instances to register (should implement ServiceProviderInterface)
	//
	// Returns:
	//   error: Any error that occurred during registration
	RegisterProviders(providers []interface{}) error

	// BootProviders boots all registered service providers.
	// This method handles the complete provider lifecycle:
	// 1. Load providers (respecting eager/deferred loading)
	// 2. Boot all loaded providers
	//
	// Parameters:
	//   ctx: Context for the boot process
	//
	// Returns:
	//   error: Any error that occurred during booting
	BootProviders(ctx context.Context) error

	// TerminateProviders gracefully terminates all terminatable providers.
	// This method is called during application shutdown to allow providers
	// to clean up resources and perform graceful shutdown.
	//
	// Parameters:
	//   ctx: Context for the termination process (with timeout)
	//
	// Returns:
	//   []error: A slice of errors that occurred during termination
	TerminateProviders(ctx context.Context) []error

	// LoadDeferredProvider loads a deferred provider when its service is requested.
	// This method is typically called by the container when a deferred service
	// is requested for the first time.
	//
	// Parameters:
	//   service: The service name that triggered the provider loading
	//
	// Returns:
	//   error: Any error that occurred during loading
	LoadDeferredProvider(service string) error

	// GetProviderRepository returns the provider repository for advanced operations.
	// This allows access to provider introspection, manifest management,
	// and other advanced provider operations.
	//
	// Returns:
	//   interface{}: The provider repository instance (should be cast to ProviderRepositoryInterface)
	GetProviderRepository() interface{}

	// GetRegisteredProviders returns all registered service provider instances.
	// This is useful for introspection, debugging, and testing.
	//
	// Returns:
	//   []interface{}: List of registered providers (should be cast to []ServiceProviderInterface)
	GetRegisteredProviders() []interface{}

	// GetLoadedProviders returns the list of loaded provider type names.
	// This is useful for debugging and monitoring which providers have been loaded.
	//
	// Returns:
	//   []string: List of loaded provider type names
	GetLoadedProviders() []string

	// GetBootedProviders returns the list of booted provider type names.
	// This is useful for debugging and monitoring which providers have been booted.
	//
	// Returns:
	//   []string: List of booted provider type names
	GetBootedProviders() []string
}
