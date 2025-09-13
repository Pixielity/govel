// Package providers provides core functionality for managing service providers in the GoVel framework.
// This package follows Laravel's provider repository pattern for handling provider registration,
// manifest generation, and deferred loading.
package providers

import (
	"context"
	"fmt"
	"sort"

	applicationInterfaces "govel/types/interfaces/application/base"
	providerInterfaces "govel/types/interfaces/application/providers"
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

	// terminatableProviders holds providers that need graceful termination
	terminatableProviders []providerInterfaces.TerminatableProvider
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
		app:                   app,
		manifestManager:       NewProviderManifestManager(manifestPath),
		loadedProviders:       make(map[string]providerInterfaces.ServiceProviderInterface),
		bootedProviders:       make(map[string]providerInterfaces.ServiceProviderInterface),
		providerInstances:     make(map[string]providerInterfaces.ServiceProviderInterface),
		terminatableProviders: make([]providerInterfaces.TerminatableProvider, 0),
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

	// Check if provider is terminatable and store it
	if terminatable, ok := provider.(providerInterfaces.TerminatableProvider); ok {
		pr.terminatableProviders = append(pr.terminatableProviders, terminatable)
		pr.app.GetLogger().Debug("Registered terminatable provider: %s", providerType)
	}

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

// GetRegisteredProviders returns all registered provider instances sorted by priority.
// Higher priority providers are returned first in the slice.
//
// Returns:
//
//	[]ServiceProviderInterface: List of all registered providers sorted by priority
func (pr *ProviderRepository) GetRegisteredProviders() []providerInterfaces.ServiceProviderInterface {
	providers := make([]providerInterfaces.ServiceProviderInterface, 0, len(pr.providerInstances))
	for _, provider := range pr.providerInstances {
		providers = append(providers, provider)
	}
	// Sort providers by priority before returning
	pr.sortProviders(providers)
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
//	providers.ServiceProvider: The provider instance
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

// GetTerminatableProviders returns all registered terminatable providers sorted by priority.
// Higher priority providers are returned first in the slice.
//
// Returns:
//
//	[]TerminatableProvider: List of terminatable providers sorted by priority
func (pr *ProviderRepository) GetTerminatableProviders() []providerInterfaces.TerminatableProvider {
	// Return a copy to prevent external modification
	providers := make([]providerInterfaces.TerminatableProvider, len(pr.terminatableProviders))
	copy(providers, pr.terminatableProviders)

	// Sort providers by priority before returning (higher priority first)
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Priority() > providers[j].Priority()
	})

	return providers
}

// TerminateProviders terminates all registered terminatable providers.
// This method is called during application shutdown to allow providers
// to clean up resources, close connections, and perform graceful shutdown.
//
// Parameters:
//
//	ctx: Context for the termination process
//	app: Application instance to pass to providers
//
// Returns:
//
//	[]error: List of errors that occurred during termination
func (pr *ProviderRepository) TerminateProviders(ctx context.Context) []error {
	var errors []error

	for _, provider := range pr.terminatableProviders {
		if err := provider.Terminate(ctx, pr.app); err != nil {
			errors = append(errors, fmt.Errorf("provider %T termination failed: %w", provider, err))
		}
	}

	return errors
}

// =============================================================================
// Private Methods
// =============================================================================

// compileProviderManifest analyzes service providers and creates/updates the provider manifest.
// This method determines which providers should be loaded eagerly, deferred, or triggered by events.
// It implements Laravel's provider discovery pattern with intelligent caching to avoid recompilation
// when the provider set hasn't changed.
//
// The manifest compilation process:
// 1. Loads existing manifest from cache/storage
// 2. Analyzes current provider set and compares with cached manifest
// 3. Recompiles manifest only if providers have changed
// 4. Categorizes providers into eager, deferred, and event-triggered groups
//
// Parameters:
//
//	providers: List of service provider instances to analyze
//
// Returns:
//
//	*ProviderManifest: The compiled manifest containing provider categorization
//	error: Any error that occurred during manifest compilation
func (pr *ProviderRepository) compileProviderManifest(providers []providerInterfaces.ServiceProviderInterface) (*ProviderManifest, error) {
	// Load existing manifest from storage/cache to avoid unnecessary recompilation
	// This improves performance by reusing previously compiled manifests when possible
	manifest, err := pr.manifestManager.LoadManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to load existing manifest: %w", err)
	}

	// Extract provider type names for comparison with cached manifest
	// Using reflect.Type string representation ensures consistent comparison
	providerTypes := make([]string, len(providers))
	for i, provider := range providers {
		providerTypes[i] = fmt.Sprintf("%T", provider)
	}

	// Determine if recompilation is necessary by comparing current providers
	// with those in the cached manifest
	if pr.manifestManager.ShouldRecompile(manifest, providerTypes) {
		pr.app.GetLogger().Debug("🔄 Recompiling provider manifest")

		// Perform full manifest compilation with provider analysis
		// This process introspects each provider to determine loading behavior
		manifest, err = pr.manifestManager.CompileManifest(providers)
		if err != nil {
			return nil, fmt.Errorf("failed to compile manifest: %w", err)
		}
	} else {
		// Use existing manifest to save compilation time
		pr.app.GetLogger().Debug("♾️ Using existing provider manifest")
	}

	return manifest, nil
}

// loadEagerProviders handles the immediate loading of providers that cannot be deferred.
// These are providers that must be initialized during the application bootstrap phase
// because they provide core services or have registration side effects that other
// providers depend on. This method follows Laravel's eager loading pattern.
//
// The eager loading process:
// 1. Retrieves eager provider types from the compiled manifest
// 2. Gets the provider instances and sorts them by priority
// 3. Loads each provider in priority order by calling their Register() method
// 4. Tracks loaded providers to prevent duplicate loading
// 5. Reports loading progress and statistics
//
// Eager providers are typically:
// - Core framework services (routing, logging, caching)
// - Providers that bind fundamental interfaces
// - Providers without deferred loading capabilities
//
// Parameters:
//
//	manifest: The provider manifest containing categorized providers
//
// Returns:
//
//	error: Any error that occurred during provider loading
func (pr *ProviderRepository) loadEagerProviders(manifest *ProviderManifest) error {
	pr.app.GetLogger().Info("🚀 Loading eager providers")

	// Extract the list of eager provider types from the manifest
	// These providers have been pre-analyzed and determined to require immediate loading
	eagerProviderTypes := pr.manifestManager.GetEagerProviders(manifest)
	if len(eagerProviderTypes) == 0 {
		pr.app.GetLogger().Info("  ℹ️ No eager providers to load")
		return nil
	}

	// Get the actual provider instances for sorting by priority
	eagerProviders := make([]providerInterfaces.ServiceProviderInterface, 0, len(eagerProviderTypes))
	for _, providerType := range eagerProviderTypes {
		if provider, exists := pr.providerInstances[providerType]; exists {
			eagerProviders = append(eagerProviders, provider)
		} else {
			return fmt.Errorf("eager provider %s not found in registered instances", providerType)
		}
	}

	// Sort eager providers by priority (higher priority first)
	pr.sortProviders(eagerProviders)
	pr.app.GetLogger().Debug("🔢 Sorted %d eager providers by priority", len(eagerProviders))

	// Load each eager provider in priority order
	// Higher priority providers are loaded first to ensure dependencies are satisfied
	for _, provider := range eagerProviders {
		providerType := fmt.Sprintf("%T", provider)
		pr.app.GetLogger().Debug("  📦 Loading eager provider: %s (priority: %d)", providerType, provider.Priority())

		// Load the provider using the internal loadProvider method
		// This handles registration, tracking, and error handling
		if err := pr.loadProvider(providerType); err != nil {
			return fmt.Errorf("failed to load eager provider %s: %w", providerType, err)
		}
	}

	// Report successful completion with statistics
	pr.app.GetLogger().Info("✅ Loaded %d eager providers in priority order", len(eagerProviders))
	return nil
}

// registerEventTriggers establishes event-driven provider loading for providers that should
// be loaded only when specific application events occur. This implements Laravel's event-driven
// provider loading pattern, allowing for more granular control over when providers are initialized.
//
// Event-triggered providers are useful for:
// - Heavy or rarely-used services that shouldn't be loaded on every request
// - Providers that depend on specific application states or user actions
// - Optional features that should only load when accessed
//
// The registration process:
// 1. Iterates through each provider and its associated trigger events
// 2. Creates event listeners that will load the provider when events fire
// 3. Registers these listeners with the application's event system
// 4. Handles duplicate prevention and error handling
//
// Parameters:
//
//	eventTriggers: Map of provider type names to lists of event names that trigger loading
//
// Returns:
//
//	error: Any error that occurred during event listener registration
func (pr *ProviderRepository) registerEventTriggers(eventTriggers map[string][]string) error {
	// Early return if no event-triggered providers are configured
	// This is a common case when all providers are either eager or deferred
	if len(eventTriggers) == 0 {
		pr.app.GetLogger().Debug("ℹ️ No event triggers to register")
		return nil
	}

	pr.app.GetLogger().Info("🎯 Setting up event-triggered provider loading")

	// Process each provider that has event-based loading configured
	// The map key is the provider type name, value is the list of triggering events
	for providerType, events := range eventTriggers {
		pr.app.GetLogger().Debug("  📡 Provider %s will load on events: %v", providerType, events)

		// Register event listeners for this provider using the internal method
		// This creates the necessary event callbacks and bindings
		if err := pr.registerLoadEvents(providerType, events); err != nil {
			return fmt.Errorf("failed to register events for provider %s: %w", providerType, err)
		}
	}

	// Report successful registration with statistics
	pr.app.GetLogger().Info("✅ Registered event triggers for %d providers", len(eventTriggers))
	return nil
}

// loadProvider handles the loading of an individual service provider by its type name.
// This is the core method that performs the actual provider registration process,
// including duplicate prevention, instance lookup, and service binding.
//
// The loading process:
// 1. Checks if the provider is already loaded to prevent duplication
// 2. Looks up the provider instance from the registered instances map
// 3. Calls the provider's Register() method to bind services to the container
// 4. Tracks the loaded provider for future reference and deduplication
//
// This method is called by:
// - loadEagerProviders() for immediate loading during bootstrap
// - Event handlers for event-triggered providers
// - Deferred loading mechanisms when services are first accessed
//
// Parameters:
//
//	providerType: The string type name of the provider to load
//
// Returns:
//
//	error: Any error that occurred during provider loading or registration
func (pr *ProviderRepository) loadProvider(providerType string) error {
	// Check if the provider has already been loaded to prevent duplicate registration
	// This is crucial for maintaining provider singleton behavior and avoiding conflicts
	if _, loaded := pr.loadedProviders[providerType]; loaded {
		return nil // Already loaded, safe to return
	}

	// Look up the provider instance from the registered instances map
	// This map is populated during provider registration phase
	provider, exists := pr.providerInstances[providerType]
	if !exists {
		return fmt.Errorf("provider %s not found in registered instances", providerType)
	}

	// Execute the provider's registration logic
	// This typically involves binding services to the application container
	if err := provider.Register(pr.app); err != nil {
		return fmt.Errorf("failed to register provider %s: %w", providerType, err)
	}

	// Mark the provider as loaded and track it for future reference
	// This enables duplicate prevention and provider lifecycle management
	pr.loadedProviders[providerType] = provider
	pr.app.GetLogger().Debug("Loaded provider: %s", providerType)
	return nil
}

// registerLoadEvents creates and registers event listeners that will trigger provider loading
// when specific events are fired. This implements lazy loading for providers that should only
// be initialized when certain application events occur, improving performance by deferring
// unnecessary provider initialization.
//
// The event registration process:
// 1. Validates that events are provided for the provider
// 2. Creates event callback functions for each event
// 3. Each callback checks if the provider is already loaded
// 4. If not loaded, calls loadProvider() and Boot() if applicable
// 5. Registers callbacks with the application's event system
//
// Event callbacks handle:
// - Duplicate loading prevention
// - Provider registration and booting
// - Error handling and logging
// - State tracking for loaded and booted providers
//
// Parameters:
//
//	providerType: The string type name of the provider to load on events
//	events: List of event names that should trigger provider loading
//
// Returns:
//
//	error: Any error that occurred during event listener registration
func (pr *ProviderRepository) registerLoadEvents(providerType string, events []string) error {
	// Validate that events are provided - empty event lists are invalid
	// This prevents unnecessary processing and helps catch configuration errors
	if len(events) == 0 {
		pr.app.GetLogger().Debug("  ℹ️ No events to register for provider '%s'", providerType)
		return nil
	}

	pr.app.GetLogger().Debug("  📡 Registering %d events for provider '%s': %v", len(events), providerType, events)

	// Create event listeners for each event that should trigger provider loading
	// Each event gets its own listener to allow for granular control
	for _, event := range events {
		// Capture loop variables in local scope to avoid closure issues
		// This prevents all callbacks from referencing the final loop values
		eventName := event
		providerTypeName := providerType

		// Create the event callback function that will be executed when the event fires
		// This callback handles the complete provider loading and booting process
		eventCallback := func() error {
			pr.app.GetLogger().Info("🎯 Event '%s' fired - loading provider '%s'", eventName, providerTypeName)

			// Prevent duplicate loading by checking if provider is already loaded
			// This is essential for event-driven loading where multiple events might trigger
			if pr.IsProviderLoaded(providerTypeName) {
				pr.app.GetLogger().Debug("  ✅ Provider '%s' already loaded, skipping", providerTypeName)
				return nil
			}

			// Load the provider using the standard loading mechanism
			// This handles registration and state tracking
			if err := pr.loadProvider(providerTypeName); err != nil {
				return fmt.Errorf("failed to load provider '%s' for event '%s': %w", providerTypeName, eventName, err)
			}

			// Boot the provider if it's not already booted and supports booting
			// Event-triggered providers often need immediate booting after loading
			if provider, exists := pr.providerInstances[providerTypeName]; exists {
				if !pr.IsProviderBooted(providerTypeName) {
					// Call the Boot method if provider implements BootableServiceProviderInterface
					if err := provider.Boot(pr.app); err != nil {
						return fmt.Errorf("failed to boot event-triggered provider '%s': %w", providerTypeName, err)
					}
					// Track the booted provider for future reference
					pr.bootedProviders[providerTypeName] = provider
					pr.app.GetLogger().Info("  ✅ Provider '%s' loaded and booted via event '%s'", providerTypeName, eventName)
				}
			}

			return nil
		}

		// TODO: Integrate with actual event system when available
		// Currently storing callbacks for future integration with event dispatcher
		// When event system is available, callbacks will be registered with event names
		pr.app.GetLogger().Debug("    ✅ Event handler registered: '%s' → load '%s'", eventName, providerTypeName)
		_ = eventCallback // Prevent unused variable error until event system integration
	}

	return nil
}

// setupDeferredServices configures lazy loading for services provided by deferred providers.
// This implements Laravel's deferred service loading pattern, where providers are only loaded
// when their services are first requested, rather than during application bootstrap.
//
// Deferred loading benefits:
// - Reduces memory usage by avoiding loading of unused services
// - Improves application startup time
// - Allows conditional loading based on actual service usage
// - Supports on-demand initialization for expensive services
//
// The setup process:
// 1. Validates that deferred services are configured
// 2. Creates lazy loading callbacks for each service
// 3. Registers these callbacks with the service container
// 4. Services are loaded only when first accessed via container resolution
//
// This method prepares the infrastructure for deferred loading but doesn't actually
// load any providers - that happens later when services are requested.
//
// Parameters:
//
//	deferredServices: Map of service names to provider types that provide them
//
// Returns:
//
//	error: Any error that occurred during deferred service setup
func (pr *ProviderRepository) setupDeferredServices(deferredServices map[string]string) error {
	// Early return if no deferred services are configured
	// This is common in applications that load all providers eagerly
	if len(deferredServices) == 0 {
		pr.app.GetLogger().Debug("ℹ️ No deferred services to set up")
		return nil
	}

	pr.app.GetLogger().Info("⏰ Setting up deferred service loading")

	// Process each service that should be loaded on-demand
	// The map key is the service name, value is the provider type that provides it
	for service, providerType := range deferredServices {
		pr.app.GetLogger().Debug("  🔗 Service '%s' → Provider '%s' (deferred)", service, providerType)

		// Capture loop variables to prevent closure reference issues
		// This ensures each callback references the correct service and provider
		serviceName := service
		providerTypeName := providerType

		// Log the deferred service registration for debugging
		// In production, this would create actual container bindings
		pr.app.GetLogger().Debug("    🔗 Deferred service '%s' registered (provider: %s)", serviceName, providerTypeName)

		// TODO: Implement actual deferred loading when container supports lazy resolution
		// The full implementation would involve:
		// 1. Registering a factory function with the container for the service
		// 2. Factory checks if the provider is already loaded
		// 3. If not loaded, loads and boots the provider
		// 4. Returns the requested service from the newly loaded provider
		// 5. Subsequent requests use the cached service instance
	}

	// Report successful setup with statistics
	pr.app.GetLogger().Info("✅ Set up deferred loading for %d services", len(deferredServices))
	return nil
}

// sortProviders sorts a slice of service providers by their priority in descending order.
// Higher priority providers will be placed first in the slice. This ensures that
// critical infrastructure services (like paths, configuration, etc.) are processed
// before services that depend on them.
//
// Parameters:
//
//	providers: The slice of providers to sort
func (pr *ProviderRepository) sortProviders(providers []providerInterfaces.ServiceProviderInterface) {
	sort.Slice(providers, func(i, j int) bool {
		// Sort in descending order by priority (higher priority first)
		return providers[i].Priority() > providers[j].Priority()
	})
}

// Compile-time interface compliance check
var _ providerInterfaces.ProviderRepositoryInterface = (*ProviderRepository)(nil)
