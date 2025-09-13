// Package application provides the core application functionality for the GoVel framework.
package application

import (
	"fmt"
	"path/filepath"

	"govel/packages/application/concerns"
	"govel/packages/application/core/maintenance"
	coreProviders "govel/packages/application/core/providers"
	"govel/packages/application/core/shutdown"
	"govel/packages/application/providers"
	"govel/packages/application/traits"
	containerTraits "govel/packages/container/traits"
	loggerTraits "govel/packages/logger/traits"
	configTraits "govel/packages/config/traits"
	applicationInterfaces "govel/packages/types/src/interfaces/application"
	"os"
	"sync"
)

// Application represents the GoVel application instance with Laravel-like features.
// It serves as the central coordinator for all application components and provides the
// primary interface for application lifecycle management using trait composition.
type Application struct {
	/**
	 * mutex provides thread-safe access to the app's internal state
	 */
	mutex sync.RWMutex

	// Application ISP trait embeddings - these provide specialized functionality

	/**
	 * Identity provides application identity management (name, version)
	 */
	*concerns.HasIdentity

	/**
	 * Runtime provides runtime state management (console mode, testing)
	 */
	*concerns.HasRuntime

	/**
	 * HasTiming provides timing management (start time, uptime)
	 */
	*concerns.HasTiming

	// Core trait embeddings - these provide all the specialized functionality

	/**
	 * DirectableTrait provides directory path management
	 */
	*traits.Directable

	/**
	 * Localizable provides internationalization functionality
	 */
	*traits.Localizable

	/**
	 * Environmentable provides environment management
	 */
	*traits.Environmentable

	/**
	 * LifecycleableTrait provides application lifecycle management
	 */
	*traits.Lifecycleable

	/**
	 * ShutdownableTrait provides shutdown management
	 */
	*shutdown.Shutdownable

	/**
	 * MaintainableTrait provides maintenance mode management
	 */
	*maintenance.Maintainable

	/**
	 * Serviceable provides service provider management
	 */
	*coreProviders.Serviceable

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

	// Create ISP trait instances (they handle environment variables internally)
	identityTrait := concerns.NewIdentity()
	runtimeTrait := concerns.NewRuntime()
	timingTrait := concerns.NewTiming()

	// Create core trait instances
	directoriesTrait := traits.NewDirectable(basePath)
	localeTrait := traits.NewLocalizable()
	environmentTrait := traits.NewEnvironmentable()
	lifecycleTrait := traits.NewLifecycleable()
	shutdownTrait := shutdown.NewShutdownable(nil)        // Will be initialized with provider repository
	containerTrait := containerTraits.NewContainable(nil) // Create default container

	// TODO: remove from here and move to the bootstrapping process
	// Register essential service providers before creating maintenance trait
	// Create a minimal application instance for service registration
	minimalApp := &Application{
		Directable:  directoriesTrait,
		Containable: containerTrait,
	}
	// Register the paths service provider to make paths available in the container
	pathsProvider := providers.NewPathsServiceProvider()
	if err := pathsProvider.Register(minimalApp); err != nil {
		// If paths registration fails, we can't continue
		panic(fmt.Sprintf("Failed to register essential paths service: %v", err))
	}

	// Create maintenance trait using container for path resolution (now paths are available)
	maintenanceTrait := maintenance.NewMaintainable(containerTrait.GetContainer())

	// Create logger and config traits
	hasLoggerTrait := loggerTraits.NewLoggableDefault()
	hasConfigTrait := configTraits.NewConfigurableWithEnvironment()

	// Initialize provider repository with manifest path
	manifestPath := filepath.Join(basePath, "bootstrap", "cache", "providers.json")
	// Create temporary application instance for provider repository
	tempApp := &Application{
		HasIdentity:     identityTrait,
		HasRuntime:      runtimeTrait,
		HasTiming:       timingTrait,
		Directable:      directoriesTrait,
		Localizable:     localeTrait,
		Environmentable: environmentTrait,
		Lifecycleable:   lifecycleTrait,
		Containable:     containerTrait,
		Loggable:        hasLoggerTrait,
		Configurable:    hasConfigTrait,
	}
	providerRepository := coreProviders.NewProviderRepository(tempApp, manifestPath)
	serviceableTrait := coreProviders.NewServiceable(tempApp, manifestPath)

	// Update shutdown trait with provider repository
	shutdownTrait = shutdown.NewShutdownable(providerRepository)

	application := &Application{
		// Embed ISP traits for specialized functionality
		HasIdentity: identityTrait,
		HasRuntime:  runtimeTrait,
		HasTiming:   timingTrait,

		// Embed core traits anonymously for method promotion
		Directable:      directoriesTrait,
		Localizable:     localeTrait,
		Environmentable: environmentTrait,
		Lifecycleable:   lifecycleTrait,
		Shutdownable:    shutdownTrait,
		Containable:     containerTrait,
		Maintainable:    maintenanceTrait,
		Serviceable:     serviceableTrait,
		Loggable:        hasLoggerTrait,
		Configurable:    hasConfigTrait,
	}

	return application
}

// GetApplicationInfo returns comprehensive application information.
// This method delegates to the provided info gathering function to avoid import cycles.
//
// Returns:
//
//	map[string]interface{}: A map containing detailed application information
//	including name, version, environment, runtime state, and trait data
//
// Example:
//
//	info := app.GetApplicationInfo()
//	fmt.Printf("App: %s v%s\n", info["name"], info["version"])
//	fmt.Printf("Environment: %s\n", info["environment"].(map[string]interface{})["environment"])
func (a *Application) GetApplicationInfo() map[string]interface{} {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	info := map[string]interface{}{
		"name":               a.GetName(),
		"version":            a.GetVersion(),
		"running_in_console": a.IsRunningInConsole(),
		"running_unit_tests": a.IsRunningUnitTests(),
		"start_time":         a.GetStartTime(),
		"uptime":             a.GetUptime(),

		// Information from traits
		"environment": a.Environmentable.GetEnvironmentInfo(),
		"directories": a.Directable.AllCustomPaths(),
		"locale": map[string]interface{}{
			"locale":          a.Localizable.GetLocale(),
			"fallback_locale": a.Localizable.GetFallbackLocale(),
			"timezone":        a.Localizable.GetTimezone(),
		},
		"lifecycle":   a.Lifecycleable.GetLifecycleInfo(),
		"shutdown":    a.Shutdownable.GetShutdownInfo(),
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
