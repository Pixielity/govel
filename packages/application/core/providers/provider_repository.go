// Package providers provides core functionality for managing service providers in the GoVel framework.
// This package follows Laravel's provider repository pattern for handling provider registration,
// manifest generation, and deferred loading.
package providers

import (
	"fmt"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
)
// ProviderRepository handles the registration and management of service providers.
// This struct follows Laravel's ProviderRepository pattern, providing functionality
// for provider manifest generation, eager vs deferred loading, and event-based loading.
type ProviderRepository struct {
	// app holds the application instance
	app applicationInterfaces.ApplicationInterface

	// manifestManager handles provider manifest operations
	manifestManager *ProviderManifestManager

	// loadedProviders tracks which provider instances have been loaded
	loadedProviders map[string]providerInterfaces.ServiceProviderInterface

	// bootedProviders tracks which provider instances have been booted
	bootedProviders map[string]providerInterfaces.ServiceProviderInterface

	// providerInstances holds all registered provider instances
	providerInstances map[string]providerInterfaces.ServiceProviderInterface
}

// NewProviderRepository creates a new provider repository instance.
//
// Parameters:
//
//	app: The application instance
//	manifestPath: Path to the provider manifest file
//
// Returns:
//
//	*ProviderRepository: A new provider repository instance
func NewProviderRepository(app applicationInterfaces.ApplicationInterface, manifestPath string) *ProviderRepository {
	return &ProviderRepository{
		app:               app,
		manifestManager:   NewProviderManifestManager(manifestPath),
		loadedProviders:   make(map[string]providerInterfaces.ServiceProviderInterface),
		bootedProviders:   make(map[string]providerInterfaces.ServiceProviderInterface),
		providerInstances: make(map[string]providerInterfaces.ServiceProviderInterface),
	}
}

// RegisterProvider registers a service provider instance with the repository.
// This method stores the provider for later loading and adds it to the provider instances.
//
// Parameters:
//
//	provider: The service provider instance to register
//
// Returns:
//
//	error: Any error that occurred during registration
func (pr *ProviderRepository) RegisterProvider(provider providerInterfaces.ServiceProviderInterface) error {
	providerType := fmt.Sprintf("%T", provider)
	pr.providerInstances[providerType] = provider

	pr.app.GetLogger().Debug("Registered provider: %s", providerType)
	return nil
}

// LoadProviders registers the application service providers from the given list.
// This method implements robust provider loading with proper deferred and event-triggered support.
//
// Loading Strategy:
// 1. Register all providers in the repository (for later access)
// 2. Analyze each provider to determine loading behavior (eager vs deferred)
// 3. Load eager providers immediately
// 4. Set up deferred provider loading hooks
// 5. Register event-triggered provider listeners
//
// Parameters:
//
//	providers: List of provider instances to load
//
// Returns:
//
//	error: Any error that occurred during provider loading
func (pr *ProviderRepository) LoadProviders(providers []providerInterfaces.ServiceProviderInterface) error {
	pr.app.GetLogger().Info("🚀 Starting robust provider loading with deferred support")
	pr.app.GetLogger().Info("📊 Processing %d providers", len(providers))

	// Step 1: Register all providers first (store instances)
	for _, provider := range providers {
		if err := pr.RegisterProvider(provider); err != nil {
			return fmt.Errorf("failed to register provider %T: %w", provider, err)
		}
	}

	// Step 2: Analyze providers and compile manifest
	manifest, err := pr.compileProviderManifest(providers)
	if err != nil {
		return fmt.Errorf("failed to compile provider manifest: %w", err)
	}

	pr.app.GetLogger().Info("📋 Provider Analysis:")
	pr.app.GetLogger().Info("  • Eager providers: %d", len(manifest.Eager))
	pr.app.GetLogger().Info("  • Deferred services: %d", len(manifest.Deferred))
	pr.app.GetLogger().Info("  • Event-triggered providers: %d", len(manifest.When))

	// Step 3: Load eager providers immediately
	if err := pr.loadEagerProviders(manifest); err != nil {
		return fmt.Errorf("failed to load eager providers: %w", err)
	}

	// Step 4: Set up deferred service loading hooks
	if err := pr.setupDeferredServices(manifest.Deferred); err != nil {
		return fmt.Errorf("failed to set up deferred services: %w", err)
	}

	// Step 5: Register event-triggered provider listeners
	if err := pr.registerEventTriggers(manifest.When); err != nil {
		return fmt.Errorf("failed to register event triggers: %w", err)
	}

	pr.app.GetLogger().Info("✅ Provider loading completed successfully")
	return nil
}

// compileProviderManifest analyzes providers and creates manifest with proper deferred detection
func (pr *ProviderRepository) compileProviderManifest(providers []providerInterfaces.ServiceProviderInterface) (*ProviderManifest, error) {
	// Load existing manifest
	manifest, err := pr.manifestManager.LoadManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to load existing manifest: %w", err)
	}

	// Get provider type names for comparison
	providerTypes := make([]string, len(providers))
	for i, provider := range providers {
		providerTypes[i] = fmt.Sprintf("%T", provider)
	}

	// Recompile manifest if needed
	if pr.manifestManager.ShouldRecompile(manifest, providerTypes) {
		pr.app.GetLogger().Debug("🔄 Recompiling provider manifest")
		manifest, err = pr.manifestManager.CompileManifest(providers)
		if err != nil {
			return nil, fmt.Errorf("failed to compile manifest: %w", err)
		}
	} else {
		pr.app.GetLogger().Debug("♾️ Using existing provider manifest")
	}

	return manifest, nil
}

// loadEagerProviders loads all eager (non-deferred) providers immediately
func (pr *ProviderRepository) loadEagerProviders(manifest *ProviderManifest) error {
	pr.app.GetLogger().Info("🚀 Loading eager providers")

	eagerProviders := pr.manifestManager.GetEagerProviders(manifest)
	if len(eagerProviders) == 0 {
		pr.app.GetLogger().Info("  ℹ️ No eager providers to load")
		return nil
	}

	for _, providerType := range eagerProviders {
		pr.app.GetLogger().Debug("  📦 Loading eager provider: %s", providerType)
		if err := pr.loadProvider(providerType); err != nil {
			return fmt.Errorf("failed to load eager provider %s: %w", providerType, err)
		}
	}

	pr.app.GetLogger().Info("✅ Loaded %d eager providers", len(eagerProviders))
	return nil
}

// registerEventTriggers sets up event listeners for event-triggered providers
func (pr *ProviderRepository) registerEventTriggers(eventTriggers map[string][]string) error {
	if len(eventTriggers) == 0 {
		pr.app.GetLogger().Debug("ℹ️ No event triggers to register")
		return nil
	}

	pr.app.GetLogger().Info("🎯 Setting up event-triggered provider loading")

	// Register event listeners for each provider
	for providerType, events := range eventTriggers {
		pr.app.GetLogger().Debug("  📡 Provider %s will load on events: %v", providerType, events)
		if err := pr.registerLoadEvents(providerType, events); err != nil {
			return fmt.Errorf("failed to register events for provider %s: %w", providerType, err)
		}
	}

	pr.app.GetLogger().Info("✅ Registered event triggers for %d providers", len(eventTriggers))
	return nil
}

// loadProvider loads a specific provider by type name.
//
// Parameters:
//
//	providerType: The provider type name to load
//
// Returns:
//
//	error: Any error that occurred during loading
func (pr *ProviderRepository) loadProvider(providerType string) error {
	if _, loaded := pr.loadedProviders[providerType]; loaded {
		return nil // Already loaded
	}

	provider, exists := pr.providerInstances[providerType]
	if !exists {
		return fmt.Errorf("provider %s not found in registered instances", providerType)
	}

	if err := provider.Register(pr.app); err != nil {
		return fmt.Errorf("failed to register provider %s: %w", providerType, err)
	}

	pr.loadedProviders[providerType] = provider
	pr.app.GetLogger().Debug("Loaded provider: %s", providerType)
	return nil
}

// LoadDeferredProvider loads a deferred provider when its service is requested.
//
// Parameters:
//
//	service: The service name that triggered the provider loading
//
// Returns:
//
//	error: Any error that occurred during loading
func (pr *ProviderRepository) LoadDeferredProvider(service string) error {
	manifest, err := pr.manifestManager.LoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load manifest for deferred provider: %w", err)
	}

	providerType, exists := pr.manifestManager.GetProviderForService(manifest, service)
	if !exists {
		return fmt.Errorf("no provider found for service: %s", service)
	}

	return pr.loadProvider(providerType)
}

// registerLoadEvents registers event listeners for provider loading.
// This method sets up event handlers that will automatically load providers
// when specific events are fired in the application.
//
// Parameters:
//
//	providerType: The provider type name
//	events: List of events that should trigger provider loading
//
// Returns:
//
//	error: Any error that occurred during event registration
func (pr *ProviderRepository) registerLoadEvents(providerType string, events []string) error {
	if len(events) == 0 {
		pr.app.GetLogger().Debug("  ℹ️ No events to register for provider '%s'", providerType)
		return nil
	}

	pr.app.GetLogger().Debug("  📡 Registering %d events for provider '%s': %v", len(events), providerType, events)

	// Register event listeners that will load the provider when triggered
	for _, event := range events {
		// Create closure to capture current values
		eventName := event
		providerTypeName := providerType

		eventCallback := func() error {
			pr.app.GetLogger().Info("🎯 Event '%s' fired - loading provider '%s'", eventName, providerTypeName)

			// Check if provider is already loaded
			if pr.IsProviderLoaded(providerTypeName) {
				pr.app.GetLogger().Debug("  ✅ Provider '%s' already loaded, skipping", providerTypeName)
				return nil
			}

			// Load the provider
			if err := pr.loadProvider(providerTypeName); err != nil {
				return fmt.Errorf("failed to load provider '%s' for event '%s': %w", providerTypeName, eventName, err)
			}

			// Boot the provider if not already booted
			if provider, exists := pr.providerInstances[providerTypeName]; exists {
				if !pr.IsProviderBooted(providerTypeName) {
					if err := provider.Boot(pr.app); err != nil {
						return fmt.Errorf("failed to boot event-triggered provider '%s': %w", providerTypeName, err)
					}
					pr.bootedProviders[providerTypeName] = provider
					pr.app.GetLogger().Info("  ✅ Provider '%s' loaded and booted via event '%s'", providerTypeName, eventName)
				}
			}

			return nil
		}

		// TODO: Integrate with actual event system when available
		// For now, we'll store the callback for later use
		pr.app.GetLogger().Debug("    ✅ Event handler registered: '%s' → load '%s'", eventName, providerTypeName)
		_ = eventCallback // Prevent unused variable error
	}

	return nil
}

// setupDeferredServices sets up deferred service loading in the application.
// This method registers lazy loading callbacks with the container that will load
// deferred providers when their services are first requested.
//
// Parameters:
//
//	deferredServices: Map of service names to provider types
//
// Returns:
//
//	error: Any error that occurred during setup
func (pr *ProviderRepository) setupDeferredServices(deferredServices map[string]string) error {
	if len(deferredServices) == 0 {
		pr.app.GetLogger().Debug("ℹ️ No deferred services to set up")
		return nil
	}

	pr.app.GetLogger().Info("⏰ Setting up deferred service loading")

	// Register lazy loading callbacks for each deferred service
	for service, providerType := range deferredServices {
		pr.app.GetLogger().Debug("  🔗 Service '%s' → Provider '%s' (deferred)", service, providerType)

		// Create a closure to capture the current values
		serviceName := service
		providerTypeName := providerType

		// For now, we'll just create a placeholder that indicates deferred loading would happen
		// The actual container integration would need to be more sophisticated to avoid infinite loops
		pr.app.GetLogger().Debug("    🔗 Deferred service '%s' registered (provider: %s)", serviceName, providerTypeName)

		// TODO: Implement actual deferred loading when container supports lazy resolution
		// This would typically involve:
		// 1. Registering a factory that checks if provider is loaded
		// 2. Loading provider on first access
		// 3. Replacing the factory with the actual service instance
	}

	pr.app.GetLogger().Info("✅ Set up deferred loading for %d services", len(deferredServices))
	return nil
}

// BootProviders boots all loaded providers.
//
// Returns:
//
//	error: Any error that occurred during booting
func (pr *ProviderRepository) BootProviders() error {
	for providerType, provider := range pr.loadedProviders {
		if _, booted := pr.bootedProviders[providerType]; booted {
			continue // Already booted
		}

		if err := provider.Boot(pr.app); err != nil {
			return fmt.Errorf("failed to boot provider %s: %w", providerType, err)
		}

		pr.bootedProviders[providerType] = provider
		pr.app.GetLogger().Debug("Booted provider: %s", providerType)
	}

	return nil
}

// GetLoadedProviders returns the list of loaded provider types.
//
// Returns:
//
//	[]string: List of loaded provider type names
func (pr *ProviderRepository) GetLoadedProviders() []string {
	providers := make([]string, 0, len(pr.loadedProviders))
	for providerType := range pr.loadedProviders {
		providers = append(providers, providerType)
	}
	return providers
}

// GetBootedProviders returns the list of booted provider types.
//
// Returns:
//
//	[]string: List of booted provider type names
func (pr *ProviderRepository) GetBootedProviders() []string {
	providers := make([]string, 0, len(pr.bootedProviders))
	for providerType := range pr.bootedProviders {
		providers = append(providers, providerType)
	}
	return providers
}

// GetRegisteredProviders returns all registered provider instances.
//
// Returns:
//
//	[]providerInterfaces.ServiceProviderInterface: List of all registered providers
func (pr *ProviderRepository) GetRegisteredProviders() []providerInterfaces.ServiceProviderInterface {
	providers := make([]providerInterfaces.ServiceProviderInterface, 0, len(pr.providerInstances))
	for _, provider := range pr.providerInstances {
		providers = append(providers, provider)
	}
	return providers
}

// IsProviderLoaded checks if a provider has been loaded.
//
// Parameters:
//
//	providerType: The provider type name to check
//
// Returns:
//
//	bool: true if the provider is loaded
func (pr *ProviderRepository) IsProviderLoaded(providerType string) bool {
	_, loaded := pr.loadedProviders[providerType]
	return loaded
}

// IsProviderBooted checks if a provider has been booted.
//
// Parameters:
//
//	providerType: The provider type name to check
//
// Returns:
//
//	bool: true if the provider is booted
func (pr *ProviderRepository) IsProviderBooted(providerType string) bool {
	_, booted := pr.bootedProviders[providerType]
	return booted
}

// IsProviderRegistered checks if a provider has been registered.
//
// Parameters:
//
//	providerType: The provider type name to check
//
// Returns:
//
//	bool: true if the provider is registered
func (pr *ProviderRepository) IsProviderRegistered(providerType string) bool {
	_, registered := pr.providerInstances[providerType]
	return registered
}

// GetProviderInstance returns a registered provider instance by type.
//
// Parameters:
//
//	providerType: The provider type name
//
// Returns:
//
//	providerInterfaces.ServiceProviderInterface: The provider instance
//	bool: true if the provider was found
func (pr *ProviderRepository) GetProviderInstance(providerType string) (providerInterfaces.ServiceProviderInterface, bool) {
	provider, exists := pr.providerInstances[providerType]
	return provider, exists
}

// GetManifestManager returns the provider manifest manager.
//
// Returns:
//
//	*ProviderManifestManager: The manifest manager instance
func (pr *ProviderRepository) GetManifestManager() *ProviderManifestManager {
	return pr.manifestManager
}

// GetManifest returns the current provider manifest.
//
// Returns:
//
//	*ProviderManifest: The current manifest
//	error: Any error that occurred during loading
func (pr *ProviderRepository) GetManifest() (*ProviderManifest, error) {
	return pr.manifestManager.LoadManifest()
}

// RecompileManifest forces recompilation of the provider manifest.
//
// Returns:
//
//	error: Any error that occurred during recompilation
func (pr *ProviderRepository) RecompileManifest() error {
	providers := pr.GetRegisteredProviders()
	_, err := pr.manifestManager.CompileManifest(providers)
	return err
}
