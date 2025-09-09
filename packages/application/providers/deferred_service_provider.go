package providers

/**
 * DeferredServiceProvider provides deferred loading functionality for service providers.
 *
 * This service provider binds the DeferredProviderInterface to implementations that support
 * deferred loading, enabling lazy initialization and improved application startup performance.
 *
 * Features:
 * - Lazy service initialization
 * - On-demand loading
 * - Performance optimization
 * - Memory efficiency
 * - Conditional service loading
 */
type DeferredServiceProvider struct {
	name     string
	deferred bool

	// Deferred loading state
	loadedServices  map[string]bool
	loaderFunctions map[string]func() (interface{}, error)
}

// NewDeferredServiceProvider creates a new deferred service provider
func NewDeferredServiceProvider() *DeferredServiceProvider {
	return &DeferredServiceProvider{
		name:            "deferred_service_provider",
		deferred:        true, // This provider itself uses deferred loading
		loadedServices:  make(map[string]bool),
		loaderFunctions: make(map[string]func() (interface{}, error)),
	}
}

// Provides returns the services this provider offers
func (p *DeferredServiceProvider) Provides() []string {
	return []string{
		"deferred_provider",
		"lazy_loader",
		"deferred_service_manager",
	}
}
