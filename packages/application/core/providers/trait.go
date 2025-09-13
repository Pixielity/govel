package providers

import (
	"context"
	"fmt"
	"sync"

	applicationInterfaces "govel/types/src/interfaces/application"
	providerInterfaces "govel/types/src/interfaces/application/providers"
)

// Serviceable trait provides service provider management functionality.
// It handles registration, booting, and termination of service providers
// following Laravel-like patterns.
type Serviceable struct {
	// Service Provider Management

	/**
	 * repository manages service provider registration and lifecycle
	 */
	repository *ProviderRepository

	/**
	 * mutex provides thread-safe access to trait fields
	 */
	mutex sync.RWMutex
}

// NewServiceable creates a new serviceable trait with the provided repository.
//
// Parameters:
//
//	app: The application instance
//	manifestPath: Path to the provider manifest file
//
// Returns:
//
//	*Serviceable: A new serviceable trait instance
//
// Example:
//
//	serviceable := NewServiceable(app, "/path/to/manifest")
func NewServiceable(app applicationInterfaces.ApplicationInterface, manifestPath string) *Serviceable {
	return &Serviceable{
		repository: NewProviderRepository(app, manifestPath),
	}
}

// RegisterProvider registers a single service provider.
//
// Parameters:
//
//	provider: The service provider instance to register (should implement ServiceProviderInterface)
//
// Returns:
//
//	error: Any error that occurred during registration
//
// Example:
//
//	postgresProvider := modules.NewPostgreSQLServiceProvider(app)
//	if err := app.RegisterProvider(postgresProvider); err != nil {
//	    return fmt.Errorf("failed to register PostgreSQL provider: %w", err)
//	}
func (sp *Serviceable) RegisterProvider(provider interface{}) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	// Cast to ServiceProviderInterface
	serviceProvider, ok := provider.(applicationInterfaces.ServiceProviderInterface)
	if !ok {
		return fmt.Errorf("provider must implement ServiceProviderInterface, got %T", provider)
	}

	if err := sp.repository.RegisterProvider(serviceProvider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	return nil
}

// RegisterProviders registers multiple service providers.
// This is a convenience method for registering multiple providers at once.
//
// Parameters:
//
//	providers: A slice of service provider instances to register (should implement ServiceProviderInterface)
//
// Returns:
//
//	error: Any error that occurred during registration
//
// Example:
//
//	providers := []interface{}{
//	    modules.NewPostgreSQLServiceProvider(app),
//	    modules.NewRedisServiceProvider(app),
//	}
//	if err := app.RegisterProviders(providers); err != nil {
//	    return fmt.Errorf("failed to register providers: %w", err)
//	}
func (sp *Serviceable) RegisterProviders(providers []interface{}) error {
	for _, provider := range providers {
		if err := sp.RegisterProvider(provider); err != nil {
			return fmt.Errorf("failed to register provider %T: %w", provider, err)
		}
	}
	return nil
}

// BootProviders boots all registered service providers.
// This method handles the provider lifecycle by loading providers
// (respecting eager/deferred loading) and then booting all loaded providers.
//
// Parameters:
//
//	ctx: Context for the boot process
//
// Returns:
//
//	error: Any error that occurred during booting
//
// Example:
//
//	ctx := context.Background()
//	if err := app.BootProviders(ctx); err != nil {
//	    return fmt.Errorf("failed to boot providers: %w", err)
//	}
func (sp *Serviceable) BootProviders(ctx context.Context) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	// Load providers using the repository (this handles eager/deferred loading)
	registeredProviders := sp.repository.GetRegisteredProviders()
	if err := sp.repository.LoadProviders(registeredProviders); err != nil {
		return fmt.Errorf("failed to load providers: %w", err)
	}

	// Boot all loaded providers
	if err := sp.repository.BootProviders(); err != nil {
		return fmt.Errorf("failed to boot providers: %w", err)
	}

	return nil
}

// TerminateProviders gracefully terminates all terminatable
// This method is called during application shutdown to allow providers
// to clean up resources, close connections, and perform graceful shutdown.
//
// Parameters:
//
//	ctx: Context for the termination process
//
// Returns:
//
//	[]error: List of errors that occurred during termination
//
// Example:
//
//	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	errors := app.TerminateProviders(shutdownCtx, app)
//	if len(errors) > 0 {
//		for _, err := range errors {
//			log.Printf("Termination error: %v", err)
//		}
//	}
func (sp *Serviceable) TerminateProviders(ctx context.Context) []error {
	return sp.repository.TerminateProviders(ctx)
}

// GetTerminatableProviders returns all registered terminatable providers.
// This method delegates to the repository.
//
// Returns:
//
//	[]TerminatableProvider: List of terminatable providers
func (sp *Serviceable) GetTerminatableProviders() []providerInterfaces.TerminatableProvider {
	return sp.repository.GetTerminatableProviders()
}

// LoadDeferredProvider loads a deferred provider when its service is requested.
// This method is typically called by the container when a deferred service is requested.
//
// Parameters:
//
//	service: The service name that triggered the provider loading
//
// Returns:
//
//	error: Any error that occurred during loading
//
// Example:
//
//	// This would typically be called by the container automatically
//	if err := app.LoadDeferredProvider("redis"); err != nil {
//	    return fmt.Errorf("failed to load Redis provider: %w", err)
//	}
func (sp *Serviceable) LoadDeferredProvider(service string) error {
	return sp.repository.LoadDeferredProvider(service)
}

// GetProviderRepository returns the provider repository for advanced operations.
// This allows access to provider introspection and management functionality.
//
// Returns:
//
//	interface{}: The provider repository instance (cast to *ProviderRepository when needed)
//
// Example:
//
//	repo := app.GetProviderRepository().(*ProviderRepository)
//	if repo.IsProviderLoaded("*modules.PostgreSQLServiceProvider") {
//	    fmt.Println("PostgreSQL provider is loaded")
//	}
func (sp *Serviceable) GetProviderRepository() interface{} {
	return sp.repository
}

// GetRegisteredProviders returns all registered service provider instances.
// This is useful for introspection and debugging.
// IsDownForMaintenance determines if the application is currently down for maintenance.
// This is compatible with Laravel's isDownForMaintenance method.
//
// Returns:
//
//	[]interface{}: List of registered providers (cast to []ServiceProviderInterface when needed)
//	bool: true if the application is in maintenance mode, false otherwise
//
// Example:
//
//	providers := app.GetRegisteredProviders()
//	for _, provider := range providers {
//	    fmt.Printf("Registered provider: %T\n", provider)
//	if app.IsDownForMaintenance() {
//		fmt.Println("Application is currently down for maintenance")
//		return
//	}
func (sp *Serviceable) GetRegisteredProviders() []interface{} {
	registeredProviders := sp.repository.GetRegisteredProviders()
	result := make([]interface{}, len(registeredProviders))
	for i, provider := range registeredProviders {
		result[i] = provider
	}
	return result
}

// GetLoadedProviders returns the list of loaded provider type names.
// This is useful for debugging and monitoring.
//
// Returns:
//
//	[]string: List of loaded provider type names
//
// Example:
//
//	loadedProviders := app.GetLoadedProviders()
//	fmt.Printf("Loaded providers: %v\n", loadedProviders)
func (sp *Serviceable) GetLoadedProviders() []string {
	return sp.repository.GetLoadedProviders()
}

// GetBootedProviders returns the list of booted provider type names.
// This is useful for debugging and monitoring.
//
// Returns:
//
//	[]string: List of booted provider type names
//
// Example:
//
//	bootedProviders := app.GetBootedProviders()
//	fmt.Printf("Booted providers: %v\n", bootedProviders)
func (sp *Serviceable) GetBootedProviders() []string {
	return sp.repository.GetBootedProviders()
}

// Compile-time interface compliance check
var _ applicationInterfaces.ServiceableInterface = (*Serviceable)(nil)
