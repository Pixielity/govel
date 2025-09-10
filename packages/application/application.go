// Package application provides the core application functionality for the GoVel framework.
package application

import (
	"context"
	"fmt"
	"path/filepath"

	"govel/packages/application/core/maintenance"
	"govel/packages/application/core/providers"
	"govel/packages/application/helpers"
	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	"govel/packages/application/traits"
	configTraits "govel/packages/config/traits"
	containerTraits "govel/packages/container/traits"
	loggerTraits "govel/packages/logger/traits"
	"os"
	"sync"
	"time"
)

// Application represents the GoVel application instance with Laravel-like features.
// It serves as the central coordinator for all application components and provides the
// primary interface for application lifecycle management using trait composition.
type Application struct {
	// Core application properties that don't belong to traits

	/**
	 * name holds the application name
	 */
	name string

	/**
	 * version holds the current application version
	 */
	version string

	/**
	 * runningInConsole indicates whether the application is running in console mode
	 */
	runningInConsole bool

	/**
	 * runningUnitTests indicates whether the application is running unit tests
	 */
	runningUnitTests bool

	/**
	 * startTime records when the application was started
	 */
	startTime time.Time

	/**
	 * mutex provides thread-safe access to the app's internal state
	 */
	mutex sync.RWMutex

	// Service Provider Management

	/**
	 * providerRepository manages service provider registration and lifecycle
	 */
	providerRepository *providers.ProviderRepository

	/**
	 * terminatableProviders holds providers that need graceful termination
	 */
	terminatableProviders []providerInterfaces.TerminatableProvider

	/**
	 * bootingCallbacks holds callbacks to execute before booting providers
	 */
	bootingCallbacks []func(ApplicationInterface)

	/**
	 * bootedCallbacks holds callbacks to execute after booting providers
	 */
	bootedCallbacks []func(ApplicationInterface)

	/**
	 * terminatingCallbacks holds callbacks to execute during application termination
	 */
	terminatingCallbacks []func(ApplicationInterface)

	/**
	 * booted indicates whether the application has been booted
	 */
	booted bool

	/**
	 * loadedProviders tracks which providers have been loaded by type
	 */
	loadedProviders map[string]bool

	/**
	 * deferredServices holds the mapping of services to their provider types
	 */
	deferredServices map[string]string

	// Trait embeddings - these provide all the specialized functionality

	/**
	 * DirectableTrait provides directory path management
	 */
	*traits.Directable

	/**
	 * LocalizableTrait provides internationalization functionality
	 */
	*traits.Localizable

	/**
	 * EnvironmentableTrait provides environment management
	 */
	*traits.Environmentable

	/**
	 * LifecycleableTrait provides application lifecycle management
	 */
	*traits.Lifecycleable

	/**
	 * ShutdownableTrait provides shutdown management
	 */
	*traits.Shutdownable

	/**
	 * HookableTrait provides hook management
	 */
	*traits.Hookable

	/**
	 * MaintainableTrait provides maintenance mode management
	 */
	*traits.Maintainable

	/**
	 * LoggableTrait provides logger management
	 */
	*loggerTraits.Loggable

	/**
	 * ConfigurableTrait provides configuration management
	 */
	*configTraits.Configurable

	/**
	 * ContainableTrait provides dependency injection container
	 */
	*containerTraits.Containable
}

// New creates a new GoVel application instance with sensible defaults.
// The application is initialized with the current working directory as
// the base path and development environment settings.
//
// Returns:
//
//	*Application: A new application instance ready for configuration and booting
func New() *Application {
	// Get the current working directory as the default base path
	basePath, err := os.Getwd()
	if err != nil {
		// Fall back to current directory if unable to determine working directory
		basePath = "."
	}

	// Initialize environment helper for reading env vars with fallbacks
	envHelper := helpers.NewEnvHelper()

	// Create trait instances
	directoriesTrait := traits.NewDirectable(basePath)
	localeTrait := traits.NewLocalizable(
		envHelper.GetAppLocale(),         // From APP_LOCALE or default
		envHelper.GetAppFallbackLocale(), // From APP_FALLBACK_LOCALE or default
		envHelper.GetAppTimezone(),       // From APP_TIMEZONE or default
	)
	environmentTrait := traits.NewEnvironmentable(
		envHelper.GetAppEnvironment(), // From APP_ENV or default
		envHelper.GetAppDebug(),       // From APP_DEBUG or default
	)
	lifecycleTrait := traits.NewLifecycleable()
	shutdownTrait := traits.NewShutdownable()
	hooksTrait := traits.NewHookable()
	containerTrait := containerTraits.NewContainable(nil) // Create default container

	// Create storage path for maintenance
	storagePath := directoriesTrait.StoragePath()
	maintenanceTrait := traits.NewMaintainable(storagePath)

	// Create logger and config traits
	hasLoggerTrait := loggerTraits.NewLoggableDefault()
	hasConfigTrait := configTraits.NewConfigurableWithEnvironment(envHelper.GetAppEnvironment())

	application := &Application{
		name:             envHelper.GetAppName(),          // From APP_NAME or default
		version:          envHelper.GetAppVersion(),       // From APP_VERSION or default
		runningInConsole: envHelper.GetRunningInConsole(), // From APP_CONSOLE or default
		runningUnitTests: envHelper.GetRunningUnitTests(), // From APP_TESTING or default
		startTime:        time.Time{},                     // Will be set when the application starts

		// Initialize service provider management
		terminatableProviders: make([]providerInterfaces.TerminatableProvider, 0),

		// Initialize callback tracking
		bootingCallbacks:      make([]func(ApplicationInterface), 0),
		bootedCallbacks:       make([]func(ApplicationInterface), 0),
		terminatingCallbacks:  make([]func(ApplicationInterface), 0),
		booted:               false,
		loadedProviders:      make(map[string]bool),
		deferredServices:     make(map[string]string),

		// Embed traits anonymously for method promotion
		Directable:      directoriesTrait,
		Localizable:     localeTrait,
		Environmentable: environmentTrait,
		Lifecycleable:   lifecycleTrait,
		Shutdownable:    shutdownTrait,
		Hookable:        hooksTrait,
		Containable:     containerTrait,
		Maintainable:    maintenanceTrait,
		Loggable:        hasLoggerTrait,
		Configurable:    hasConfigTrait,
	}

	// Set up maintenance manager with application reference (once it's created)
	maintenanceManager := maintenance.NewMaintenanceManager(application)
	application.Maintainable.SetManager(maintenanceManager)

	// Initialize provider repository with manifest path
	manifestPath := filepath.Join(basePath, "bootstrap", "cache", "providers.json")
	application.providerRepository = providers.NewProviderRepository(application, manifestPath)

	return application
}

// Application name and version methods

// GetName returns the application name.
func (a *Application) GetName() string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.name
}

// SetName sets the application name.
func (a *Application) SetName(name string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.name = name
}

// Name returns the application name (Laravel-like API).
func (a *Application) Name() string {
	return a.GetName()
}

// GetVersion returns the application version.
func (a *Application) GetVersion() string {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.version
}

// SetVersion sets the application version.
func (a *Application) SetVersion(version string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.version = version
}

// Version returns the application version (Laravel-like API).
func (a *Application) Version() string {
	return a.GetVersion()
}

// Console and test mode methods

// IsRunningInConsole returns whether the application is running in console mode.
func (a *Application) IsRunningInConsole() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.runningInConsole
}

// SetRunningInConsole sets whether the application is running in console mode.
func (a *Application) SetRunningInConsole(console bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.runningInConsole = console
}

// IsRunningUnitTests returns whether the application is running unit tests.
func (a *Application) IsRunningUnitTests() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.runningUnitTests
}

// SetRunningUnitTests sets whether the application is running unit tests.
func (a *Application) SetRunningUnitTests(testing bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.runningUnitTests = testing
}

// Application timing methods

// GetStartTime returns when the application was started.
func (a *Application) GetStartTime() time.Time {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.startTime
}

// SetStartTime sets when the application was started.
func (a *Application) SetStartTime(startTime time.Time) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.startTime = startTime
}

// GetUptime returns how long the application has been running.
func (a *Application) GetUptime() time.Duration {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if a.startTime.IsZero() {
		return 0
	}
	return time.Since(a.startTime)
}

// Application Information

// GetApplicationInfo returns comprehensive application information.
func (a *Application) GetApplicationInfo() map[string]interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	info := map[string]interface{}{
		"name":               a.name,
		"version":            a.version,
		"running_in_console": a.runningInConsole,
		"running_unit_tests": a.runningUnitTests,
		"start_time":         a.startTime,
		"uptime":             a.GetUptime(),

		// Information from traits
		"environment": a.Environmentable.GetEnvironmentInfo(),
		"directories": a.Directable.GetAllCustomPaths(),
		"locale": map[string]interface{}{
			"locale":          a.Localizable.GetLocale(),
			"fallback_locale": a.Localizable.GetFallbackLocale(),
			"timezone":        a.Localizable.GetTimezone(),
		},
		"lifecycle":   a.Lifecycleable.GetLifecycleInfo(),
		"shutdown":    a.Shutdownable.GetShutdownInfo(),
		"hooks":       a.Hookable.GetHooksInfo(),
		"container":   a.Containable.GetContainerInfo(),
		"maintenance": a.Maintainable.GetMaintenanceInfo(),
		"logger":      a.Loggable.GetLoggerInfo(),
		"config":      a.Configurable.GetConfigInfo(),
	}

	return info
}

// Service Provider Management

// RegisterProvider registers a service provider with the application.
// This method stores the provider instance for later loading and adds it
// to the terminatable providers list if it implements TerminatableProvider.
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
func (a *Application) RegisterProvider(provider interface{}) error {
	// Cast to ServiceProviderInterface
	serviceProvider, ok := provider.(providerInterfaces.ServiceProviderInterface)
	if !ok {
		return fmt.Errorf("provider must implement ServiceProviderInterface, got %T", provider)
	}

	if err := a.providerRepository.RegisterProvider(serviceProvider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}

	// Check if provider is terminatable and store it
	if terminatable, ok := provider.(providerInterfaces.TerminatableProvider); ok {
		a.terminatableProviders = append(a.terminatableProviders, terminatable)
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
func (a *Application) RegisterProviders(providers []interface{}) error {
	for _, provider := range providers {
		if err := a.RegisterProvider(provider); err != nil {
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
func (a *Application) BootProviders(ctx context.Context) error {
	// Load providers using the repository (this handles eager/deferred loading)
	registeredProviders := a.providerRepository.GetRegisteredProviders()
	if err := a.providerRepository.LoadProviders(registeredProviders); err != nil {
		return fmt.Errorf("failed to load providers: %w", err)
	}

	// Boot all loaded providers
	if err := a.providerRepository.BootProviders(); err != nil {
		return fmt.Errorf("failed to boot providers: %w", err)
	}

	return nil
}

// TerminateProviders gracefully terminates all terminatable providers.
// This method is called during application shutdown to allow providers
// to clean up resources, close connections, and perform graceful shutdown.
//
// Parameters:
//
//	ctx: Context for the termination process (with timeout)
//
// Returns:
//
//	[]error: A slice of errors that occurred during termination
//
// Example:
//
//	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	errors := app.TerminateProviders(shutdownCtx)
//	if len(errors) > 0 {
//	    for _, err := range errors {
//	        log.Printf("Termination error: %v", err)
//	    }
//	}
func (a *Application) TerminateProviders(ctx context.Context) []error {
	var errors []error

	for _, provider := range a.terminatableProviders {
		if err := provider.Terminate(ctx, a); err != nil {
			errors = append(errors, fmt.Errorf("provider %T termination failed: %w", provider, err))
		}
	}

	return errors
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
func (a *Application) LoadDeferredProvider(service string) error {
	return a.providerRepository.LoadDeferredProvider(service)
}

// GetProviderRepository returns the provider repository for advanced operations.
// This allows access to provider introspection and management functionality.
//
// Returns:
//
//	interface{}: The provider repository instance (cast to *providers.ProviderRepository when needed)
//
// Example:
//
//	repo := app.GetProviderRepository().(*providers.ProviderRepository)
//	if repo.IsProviderLoaded("*modules.PostgreSQLServiceProvider") {
//	    fmt.Println("PostgreSQL provider is loaded")
//	}
func (a *Application) GetProviderRepository() interface{} {
	return a.providerRepository
}

// GetRegisteredProviders returns all registered service provider instances.
// This is useful for introspection and debugging.
//
// Returns:
//
//	[]interface{}: List of registered providers (cast to []ServiceProviderInterface when needed)
//
// Example:
//
//	providers := app.GetRegisteredProviders()
//	for _, provider := range providers {
//	    fmt.Printf("Registered provider: %T\n", provider)
//	}
func (a *Application) GetRegisteredProviders() []interface{} {
	registeredProviders := a.providerRepository.GetRegisteredProviders()
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
func (a *Application) GetLoadedProviders() []string {
	return a.providerRepository.GetLoadedProviders()
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
func (a *Application) GetBootedProviders() []string {
	return a.providerRepository.GetBootedProviders()
}

// Compile-time interface compliance check
// This ensures that Application implements the ApplicationInterface at compile time
var _ applicationInterfaces.ApplicationInterface = (*Application)(nil)
