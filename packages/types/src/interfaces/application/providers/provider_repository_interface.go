package interfaces

import (
	"context"
)

// ProviderRepositoryInterface defines the contract for the provider repository
// without creating circular dependencies with the application package.
//
// This interface provides methods for:
// - Provider registration and loading
// - Provider lifecycle management
// - Provider introspection and debugging
// - Manifest management
type ProviderRepositoryInterface interface {
	// RegisterProvider registers a service provider instance with the repository.
	//
	// Parameters:
	//   provider: The service provider instance to register
	//
	// Returns:
	//   error: Any error that occurred during registration
	RegisterProvider(provider ServiceProviderInterface) error

	// LoadProviders registers the application service providers from the given list.
	//
	// Parameters:
	//   providers: List of provider instances to load
	//
	// Returns:
	//   error: Any error that occurred during provider loading
	LoadProviders(providers []ServiceProviderInterface) error

	// LoadDeferredProvider loads a deferred provider when its service is requested.
	//
	// Parameters:
	//   service: The service name that triggered the provider loading
	//
	// Returns:
	//   error: Any error that occurred during loading
	LoadDeferredProvider(service string) error

	// BootProviders boots all loaded providers.
	//
	// Returns:
	//   error: Any error that occurred during booting
	BootProviders() error

	// GetLoadedProviders returns the list of loaded provider types.
	//
	// Returns:
	//   []string: List of loaded provider type names
	GetLoadedProviders() []string

	// GetBootedProviders returns the list of booted provider types.
	//
	// Returns:
	//   []string: List of booted provider type names
	GetBootedProviders() []string

	// GetRegisteredProviders returns all registered provider instances.
	//
	// Returns:
	//   []ServiceProviderInterface: List of all registered providers
	GetRegisteredProviders() []ServiceProviderInterface

	// IsProviderLoaded checks if a provider has been loaded.
	//
	// Parameters:
	//   providerType: The provider type name to check
	//
	// Returns:
	//   bool: true if the provider is loaded
	IsProviderLoaded(providerType string) bool

	// IsProviderBooted checks if a provider has been booted.
	//
	// Parameters:
	//   providerType: The provider type name to check
	//
	// Returns:
	//   bool: true if the provider is booted
	IsProviderBooted(providerType string) bool

	// IsProviderRegistered checks if a provider has been registered.
	//
	// Parameters:
	//   providerType: The provider type name to check
	//
	// Returns:
	//   bool: true if the provider is registered
	IsProviderRegistered(providerType string) bool

	// GetProviderInstance returns a registered provider instance by type.
	//
	// Parameters:
	//   providerType: The provider type name
	//
	// Returns:
	//   ServiceProviderInterface: The provider instance
	//   bool: true if the provider was found
	GetProviderInstance(providerType string) (ServiceProviderInterface, bool)

	// RecompileManifest forces recompilation of the provider manifest.
	//
	// Returns:
	//   error: Any error that occurred during recompilation
	RecompileManifest() error

	// GetTerminatableProviders returns all registered terminatable providers.
	//
	// Returns:
	//   []TerminatableProvider: List of terminatable providers sorted by priority
	GetTerminatableProviders() []TerminatableProvider

	// TerminateProviders terminates all registered terminatable providers.
	//
	// Parameters:
	//   ctx: Context for the termination process
	//
	// Returns:
	//   []error: List of errors that occurred during termination
	TerminateProviders(ctx context.Context) []error
}
