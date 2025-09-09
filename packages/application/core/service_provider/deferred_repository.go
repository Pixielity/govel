package service_provider

import (
	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
)

// Type aliases for cleaner usage
type DeferredServiceProvider = providerInterfaces.DeferredServiceProvider

// DeferredProviderRepository manages deferred service providers and handles
// on-demand loading when services are requested. This component maintains
// a mapping between service identifiers and their corresponding providers.
type DeferredProviderRepository struct {
	// serviceToProvider maps service identifiers to their providers
	serviceToProvider map[string]DeferredServiceProvider

	// loadedProviders tracks which deferred providers have been loaded
	loadedProviders map[DeferredServiceProvider]bool

	// application holds a reference to the application for provider registration
	application applicationInterfaces.ApplicationInterface
}

// NewDeferredProviderRepository creates a new deferred provider repository.
//
// Parameters:
//
//	application: The application instance for provider registration
//
// Returns:
//
//	*DeferredProviderRepository: A new repository instance
func NewDeferredProviderRepository(application applicationInterfaces.ApplicationInterface) *DeferredProviderRepository {
	return &DeferredProviderRepository{
		serviceToProvider: make(map[string]DeferredServiceProvider),
		loadedProviders:   make(map[DeferredServiceProvider]bool),
		application:       application,
	}
}

// Register adds a deferred service provider to the repository.
// The provider is not immediately registered with the application,
// but is stored for on-demand loading.
//
// Parameters:
//
//	provider: The deferred service provider to register
//
// Example:
//
//	emailProvider := &EmailServiceProvider{}
//	repository.Register(emailProvider)
func (r *DeferredProviderRepository) Register(provider DeferredServiceProvider) {
	// Map each provided service to this provider
	for _, service := range provider.Provides() {
		r.serviceToProvider[service] = provider
	}
}

// LoadFor loads the deferred provider that provides the specified service.
// If the provider has already been loaded, this method does nothing.
// If no provider provides the service, this method returns false.
//
// Parameters:
//
//	service: The service identifier to load a provider for
//
// Returns:
//
//	bool: true if a provider was found and loaded, false otherwise
//	error: Any error that occurred during provider loading
//
// Example:
//
//	loaded, err := repository.LoadFor("mailer")
//	if err != nil {
//	    return fmt.Errorf("failed to load mailer provider: %w", err)
//	}
//	if !loaded {
//	    return fmt.Errorf("no provider found for mailer service")
//	}
func (r *DeferredProviderRepository) LoadFor(service string) (bool, error) {
	provider, exists := r.serviceToProvider[service]
	if !exists {
		return false, nil // No provider for this service
	}

	// Check if already loaded
	if r.loadedProviders[provider] {
		return true, nil // Already loaded
	}

	// Load the provider
	if err := r.loadProvider(provider); err != nil {
		return false, err
	}

	// Mark as loaded
	r.loadedProviders[provider] = true
	return true, nil
}

// Provides returns all service identifiers that have deferred providers.
// This is useful for debugging and service discovery.
//
// Returns:
//
//	[]string: All service identifiers with deferred providers
//
// Example:
//
//	services := repository.Provides()
//	for _, service := range services {
//	    fmt.Printf("Deferred service: %s\n", service)
//	}
func (r *DeferredProviderRepository) Provides() []string {
	services := make([]string, 0, len(r.serviceToProvider))
	for service := range r.serviceToProvider {
		services = append(services, service)
	}
	return services
}

// IsLoaded returns whether a deferred provider has been loaded.
//
// Parameters:
//
//	provider: The provider to check
//
// Returns:
//
//	bool: true if the provider has been loaded, false otherwise
func (r *DeferredProviderRepository) IsLoaded(provider DeferredServiceProvider) bool {
	return r.loadedProviders[provider]
}

// LoadedCount returns the number of deferred providers that have been loaded.
//
// Returns:
//
//	int: Number of loaded providers
func (r *DeferredProviderRepository) LoadedCount() int {
	return len(r.loadedProviders)
}

// TotalCount returns the total number of registered deferred providers.
//
// Returns:
//
//	int: Total number of deferred providers
func (r *DeferredProviderRepository) TotalCount() int {
	providerSet := make(map[DeferredServiceProvider]bool)
	for _, provider := range r.serviceToProvider {
		providerSet[provider] = true
	}
	return len(providerSet)
}

// loadProvider performs the actual loading of a deferred provider.
// This includes registering and booting the provider.
//
// Parameters:
//
//	provider: The provider to load
//
// Returns:
//
//	error: Any error that occurred during loading
func (r *DeferredProviderRepository) loadProvider(provider DeferredServiceProvider) error {
	// Register the provider
	if err := provider.Register(r.application); err != nil {
		return err
	}

	// Boot the provider
	if err := provider.Boot(r.application); err != nil {
		return err
	}

	return nil
}
