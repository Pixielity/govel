// Package application provides the core application functionality for the GoVel framework.
package application

import (
	"context"
	"fmt"
	"govel/packages/application/core/maintenance"
	"govel/packages/application/core/service_provider"
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

	// Service management properties

	/**
	 * serviceProviderRegistry manages registered service providers
	 */
	serviceProviderRegistry *service_provider.ServiceProviderRegistry

	/**
	 * deferredProviderRepository manages deferred service providers
	 */
	deferredProviderRepository *service_provider.DeferredProviderRepository

	/**
	 * terminationManager manages terminatable service providers
	 */
	terminationManager *service_provider.TerminationManager

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

	// Initialize service provider management
	application.serviceProviderRegistry = service_provider.NewServiceProviderRegistry()
	application.deferredProviderRepository = service_provider.NewDeferredProviderRepository(application)
	application.terminationManager = service_provider.NewTerminationManager(application)

	// Set up maintenance manager with application reference (once it's created)
	maintenanceManager := maintenance.NewMaintenanceManager(application)
	application.Maintainable.SetManager(maintenanceManager)

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

// Service Provider Management

// RegisterProvider registers a service provider with the application.
// TODO: Fix interface compliance to re-enable provider registration
func (a *Application) RegisterProvider(provider providerInterfaces.ServiceProviderInterface) error {
	a.serviceProviderRegistry.Register(provider)

	if err := provider.Register(a); err != nil {
		return fmt.Errorf("service provider registration failed: %w", err)
	}

	return nil
}

// BootProviders boots all registered service providers.
// TODO: Fix interface compliance to re-enable provider booting
func (a *Application) BootProviders(ctx context.Context) error {
	providers := a.serviceProviderRegistry.AllServiceProviders()

	for _, provider := range providers {
		if err := provider.Boot(a); err != nil {
			return fmt.Errorf("failed to boot service provider %T: %w", provider, err)
		}
	}

	return nil
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

// Compile-time interface compliance check
// This ensures that Application implements the ApplicationInterface at compile time
var _ applicationInterfaces.ApplicationInterface = (*Application)(nil)
