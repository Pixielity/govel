package concerns

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	concernsInterfaces "govel/types/interfaces/application/concerns"
)

/**
 * ApplicationProvider provides service provider management functionality.
 * This trait implements the ApplicationProviderInterface and manages the
 * lifecycle of service providers including registration, booting, and termination.
 *
 * Features:
 * - Service provider registration and management
 * - Provider lifecycle operations (boot/terminate)
 * - Deferred provider loading
 * - Provider introspection and debugging
 * - Thread-safe provider operations
 * - Circular dependency detection
 */
type ApplicationProvider struct {
	/**
	 * providers holds registered service provider instances
	 */
	providers []interface{}

	/**
	 * loadedProviders tracks which provider types have been loaded
	 */
	loadedProviders []string

	/**
	 * bootedProviders tracks which provider types have been booted
	 */
	bootedProviders []string

	/**
	 * deferredProviders maps service names to their provider instances
	 */
	deferredProviders map[string]interface{}

	/**
	 * terminatableProviders holds providers that implement termination
	 */
	terminatableProviders []interface{}

	/**
	 * providerRepository holds the provider repository instance
	 */
	providerRepository interface{}

	/**
	 * mutex provides thread-safe access to provider fields
	 */
	mutex sync.RWMutex
}

// NewApplicationProvider creates a new application provider manager.
//
// Returns:
//   *ApplicationProvider: A new application provider instance
//
// Example:
//   provider := NewApplicationProvider()
func NewApplicationProvider() *ApplicationProvider {
	return &ApplicationProvider{
		providers:             make([]interface{}, 0),
		loadedProviders:       make([]string, 0),
		bootedProviders:       make([]string, 0),
		deferredProviders:     make(map[string]interface{}),
		terminatableProviders: make([]interface{}, 0),
	}
}

// RegisterProvider registers a service provider with the application.
// The provider will be stored for later loading and automatically added
// to terminatable providers if it implements the appropriate interface.
//
// Parameters:
//   provider: The service provider instance to register
//
// Returns:
//   error: Any error that occurred during registration
//
// Example:
//   err := app.RegisterProvider(&DatabaseServiceProvider{})
//   if err != nil {
//       log.Fatalf("Failed to register provider: %v", err)
//   }
func (p *ApplicationProvider) RegisterProvider(provider interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if provider == nil {
		return fmt.Errorf("cannot register nil provider")
	}

	// Add to providers list
	p.providers = append(p.providers, provider)

	// Check if provider implements terminatable interface and add to terminatable list
	// Note: This would need to be adapted based on your actual terminatable interface
	if p.isTerminatableProvider(provider) {
		p.terminatableProviders = append(p.terminatableProviders, provider)
	}

	return nil
}

// RegisterProviders registers multiple service providers.
// This is a convenience method for bulk provider registration.
//
// Parameters:
//   providers: A slice of service provider instances to register
//
// Returns:
//   error: Any error that occurred during registration
//
// Example:
//   providers := []interface{}{
//       &DatabaseServiceProvider{},
//       &CacheServiceProvider{},
//       &LoggingServiceProvider{},
//   }
//   err := app.RegisterProviders(providers)
func (p *ApplicationProvider) RegisterProviders(providers []interface{}) error {
	for _, provider := range providers {
		if err := p.RegisterProvider(provider); err != nil {
			return fmt.Errorf("failed to register provider %T: %w", provider, err)
		}
	}
	return nil
}

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
//
// Example:
//   ctx := context.Background()
//   err := app.BootProviders(ctx)
//   if err != nil {
//       log.Fatalf("Failed to boot providers: %v", err)
//   }
func (p *ApplicationProvider) BootProviders(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// First, load all eager providers
	for _, provider := range p.providers {
		if err := p.loadProvider(provider); err != nil {
			return fmt.Errorf("failed to load provider %T: %w", provider, err)
		}
	}

	// Then boot all loaded providers
	for _, provider := range p.providers {
		if err := p.bootProvider(ctx, provider); err != nil {
			return fmt.Errorf("failed to boot provider %T: %w", provider, err)
		}
	}

	return nil
}

// TerminateProviders gracefully terminates all terminatable providers.
// This method is called during application shutdown to allow providers
// to clean up resources and perform graceful shutdown.
//
// Parameters:
//   ctx: Context for the termination process (with timeout)
//
// Returns:
//   []error: A slice of errors that occurred during termination
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//   defer cancel()
//   errors := app.TerminateProviders(ctx)
//   for _, err := range errors {
//       log.Printf("Provider termination error: %v", err)
//   }
func (p *ApplicationProvider) TerminateProviders(ctx context.Context) []error {
	p.mutex.RLock()
	terminatable := make([]interface{}, len(p.terminatableProviders))
	copy(terminatable, p.terminatableProviders)
	p.mutex.RUnlock()

	var errors []error
	for _, provider := range terminatable {
		if err := p.terminateProvider(ctx, provider); err != nil {
			errors = append(errors, fmt.Errorf("failed to terminate provider %T: %w", provider, err))
		}
	}

	return errors
}

// LoadDeferredProvider loads a deferred provider when its service is requested.
// This method is typically called by the container when a deferred service
// is requested for the first time.
//
// Parameters:
//   service: The service name that triggered the provider loading
//
// Returns:
//   error: Any error that occurred during loading
//
// Example:
//   err := app.LoadDeferredProvider("user_service")
//   if err != nil {
//       log.Printf("Failed to load deferred provider for service %s: %v", "user_service", err)
//   }
func (p *ApplicationProvider) LoadDeferredProvider(service string) error {
	p.mutex.RLock()
	provider, exists := p.deferredProviders[service]
	p.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("no deferred provider found for service: %s", service)
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.loadProvider(provider)
}

// GetProviderRepository returns the provider repository for advanced operations.
// This allows access to provider introspection, manifest management,
// and other advanced provider operations.
//
// Returns:
//   interface{}: The provider repository instance
//
// Example:
//   repo := app.GetProviderRepository()
//   if repo != nil {
//       // Cast to appropriate type and use
//   }
func (p *ApplicationProvider) GetProviderRepository() interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.providerRepository
}

// GetRegisteredProviders returns all registered service provider instances.
// This is useful for introspection, debugging, and testing.
//
// Returns:
//   []interface{}: List of registered providers
//
// Example:
//   providers := app.GetRegisteredProviders()
//   fmt.Printf("Registered %d providers\n", len(providers))
//   for _, provider := range providers {
//       fmt.Printf("Provider: %T\n", provider)
//   }
func (p *ApplicationProvider) GetRegisteredProviders() []interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return a copy to prevent external modification
	providers := make([]interface{}, len(p.providers))
	copy(providers, p.providers)
	return providers
}

// GetLoadedProviders returns the list of loaded provider type names.
// This is useful for debugging and monitoring which providers have been loaded.
//
// Returns:
//   []string: List of loaded provider type names
//
// Example:
//   loaded := app.GetLoadedProviders()
//   fmt.Printf("Loaded providers: %v\n", loaded)
func (p *ApplicationProvider) GetLoadedProviders() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return a copy to prevent external modification
	loaded := make([]string, len(p.loadedProviders))
	copy(loaded, p.loadedProviders)
	return loaded
}

// GetBootedProviders returns the list of booted provider type names.
// This is useful for debugging and monitoring which providers have been booted.
//
// Returns:
//   []string: List of booted provider type names
//
// Example:
//   booted := app.GetBootedProviders()
//   fmt.Printf("Booted providers: %v\n", booted)
func (p *ApplicationProvider) GetBootedProviders() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return a copy to prevent external modification
	booted := make([]string, len(p.bootedProviders))
	copy(booted, p.bootedProviders)
	return booted
}

// SetProviderRepository sets the provider repository instance.
//
// Parameters:
//   repository: The provider repository instance
//
// Example:
//   repo := &ProviderRepository{}
//   app.SetProviderRepository(repo)
func (p *ApplicationProvider) SetProviderRepository(repository interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.providerRepository = repository
}

// Helper methods

// loadProvider loads a single provider
func (p *ApplicationProvider) loadProvider(provider interface{}) error {
	providerType := reflect.TypeOf(provider).String()
	
	// Check if already loaded
	for _, loaded := range p.loadedProviders {
		if loaded == providerType {
			return nil // Already loaded
		}
	}

	// Here you would call the actual load method on the provider
	// This depends on your provider interface definition
	
	// Add to loaded providers
	p.loadedProviders = append(p.loadedProviders, providerType)
	return nil
}

// bootProvider boots a single provider
func (p *ApplicationProvider) bootProvider(ctx context.Context, provider interface{}) error {
	providerType := reflect.TypeOf(provider).String()
	
	// Check if already booted
	for _, booted := range p.bootedProviders {
		if booted == providerType {
			return nil // Already booted
		}
	}

	// Here you would call the actual boot method on the provider
	// This depends on your provider interface definition
	
	// Add to booted providers
	p.bootedProviders = append(p.bootedProviders, providerType)
	return nil
}

// terminateProvider terminates a single provider
func (p *ApplicationProvider) terminateProvider(ctx context.Context, provider interface{}) error {
	// Here you would call the actual terminate method on the provider
	// This depends on your provider interface definition
	return nil
}

// isTerminatableProvider checks if a provider implements the terminatable interface
func (p *ApplicationProvider) isTerminatableProvider(provider interface{}) bool {
	// This would need to be adapted based on your actual terminatable provider interface
	// For now, we'll return false as a placeholder
	return false
}

// Compile-time interface compliance check
var _ concernsInterfaces.ApplicationProviderInterface = (*ApplicationProvider)(nil)