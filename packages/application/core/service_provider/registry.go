package service_provider

import (
	providerInterfaces "govel/packages/application/interfaces/providers"
)

// ServiceProviderRegistry manages the collection of registered service providers.
// This struct handles provider registration, sorting by priority, and lifecycle management.
//
// The registry maintains separate collections for different types of providers:
// - Regular service providers for normal registration
// - Deferred providers for on-demand loading
// - Booted providers tracking for lifecycle management
type ServiceProviderRegistry struct {
	// providers holds the registered service providers
	providers []providerInterfaces.ServiceProviderInterface

	// deferredProviders holds providers that should be loaded on-demand
	deferredProviders map[string]providerInterfaces.ServiceProviderInterface

	// bootedProviders tracks which providers have been booted
	bootedProviders map[providerInterfaces.ServiceProviderInterface]bool
}

// NewServiceProviderRegistry creates a new service provider registry.
//
// Returns:
//
//	*ServiceProviderRegistry: A new registry instance ready for provider registration
//
// Example:
//
//	registry := service_providers.NewServiceProviderRegistry()
//	registry.Register(&MyServiceProvider{})
func NewServiceProviderRegistry() *ServiceProviderRegistry {
	return &ServiceProviderRegistry{
		providers:         make([]providerInterfaces.ServiceProviderInterface, 0),
		deferredProviders: make(map[string]providerInterfaces.ServiceProviderInterface),
		bootedProviders:   make(map[providerInterfaces.ServiceProviderInterface]bool),
	}
}

// Register adds a service provider to the registry.
// The provider will be registered during the application boot process.
//
// Parameters:
//
//	provider: The service provider to register
//
// Example:
//
//	provider := &DatabaseServiceProvider{}
//	registry.Register(provider)
func (r *ServiceProviderRegistry) Register(provider providerInterfaces.ServiceProviderInterface) {
	r.providers = append(r.providers, provider)
}

// Providers returns all registered service providers sorted by priority.
// Providers with lower priority values are returned first, ensuring proper
// initialization order based on dependencies.
//
// Returns:
//
//	[]providerInterfaces.ServiceProviderInterface: All registered providers in priority order
//
// Example:
//
//	providers := registry.Providers()
//	for _, provider := range providers {
//	    provider.Register(application)
//	}
func (r *ServiceProviderRegistry) Providers() []providerInterfaces.ServiceProviderInterface {
	// Sort providers by priority (lower values first)
	sorted := make([]providerInterfaces.ServiceProviderInterface, len(r.providers))
	copy(sorted, r.providers)

	// Simple bubble sort by priority
	// TODO: Consider using sort.Slice for better performance with many providers
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Priority() > sorted[j+1].Priority() {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// AllServiceProviders returns all registered service providers sorted by priority.
// This is an alias for Providers() method for API consistency.
//
// Returns:
//
//	[]providerInterfaces.ServiceProviderInterface: All registered providers in priority order
//
// Example:
//
//	providers := registry.AllServiceProviders()
//	for _, provider := range providers {
//	    provider.Boot(application)
//	}
func (r *ServiceProviderRegistry) AllServiceProviders() []providerInterfaces.ServiceProviderInterface {
	return r.Providers()
}

// Count returns the number of registered service providers.
//
// Returns:
//
//	int: Total number of registered providers
//
// Example:
//
//	if registry.Count() == 0 {
//	    log.Println("No service providers registered")
//	}
func (r *ServiceProviderRegistry) Count() int {
	return len(r.providers)
}

// IsBooted returns whether a specific provider has been booted.
// This is useful for debugging and ensuring proper provider lifecycle.
//
// Parameters:
//
//	provider: The provider to check
//
// Returns:
//
//	bool: true if the provider has been booted, false otherwise
//
// Example:
//
//	if !registry.IsBooted(provider) {
//	    log.Printf("Provider %T not yet booted", provider)
//	}
func (r *ServiceProviderRegistry) IsBooted(provider providerInterfaces.ServiceProviderInterface) bool {
	return r.bootedProviders[provider]
}

// MarkAsBooted marks a provider as booted.
// This should be called after successfully booting a provider.
//
// Parameters:
//
//	provider: The provider to mark as booted
//
// Example:
//
//	err := provider.Boot(application)
//	if err == nil {
//	    registry.MarkAsBooted(provider)
//	}
func (r *ServiceProviderRegistry) MarkAsBooted(provider providerInterfaces.ServiceProviderInterface) {
	r.bootedProviders[provider] = true
}

// BootedCount returns the number of providers that have been booted.
//
// Returns:
//
//	int: Number of booted providers
func (r *ServiceProviderRegistry) BootedCount() int {
	return len(r.bootedProviders)
}

// UnbootedProviders returns providers that haven't been booted yet.
//
// Returns:
//
//	[]providerInterfaces.ServiceProviderInterface: Unbooted providers
func (r *ServiceProviderRegistry) UnbootedProviders() []providerInterfaces.ServiceProviderInterface {
	var unbooted []providerInterfaces.ServiceProviderInterface

	for _, provider := range r.providers {
		if !r.IsBooted(provider) {
			unbooted = append(unbooted, provider)
		}
	}

	return unbooted
}

// Clear removes all registered providers and resets the registry state.
// This is primarily useful for testing scenarios.
//
// Example:
//
//	// In test cleanup
//	defer registry.Clear()
func (r *ServiceProviderRegistry) Clear() {
	r.providers = make([]providerInterfaces.ServiceProviderInterface, 0)
	r.deferredProviders = make(map[string]providerInterfaces.ServiceProviderInterface)
	r.bootedProviders = make(map[providerInterfaces.ServiceProviderInterface]bool)
}
