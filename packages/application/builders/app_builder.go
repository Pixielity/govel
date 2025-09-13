package builders

import (
	"time"

	"govel/packages/application"
	"govel/packages/application/helpers"
	"govel/packages/container"
	enums "govel/packages/types/src/enums/application"
	providerInterfaces "govel/packages/types/src/interfaces/application/providers"
)

// app_builder.go implements the Builder pattern for GoVel application creation.
// This provides a fluent API for configuring applications similar to modern frameworks.
//
// The builder allows for method chaining to create and configure applications:
//   application := NewApp().
//       WithName("My API").
//       WithVersion("1.2.0").
//       WithEnvironment("production").
//       WithDebug(false).
//       Build()

// AppBuilder provides a fluent interface for building GoVel applications.
// It follows the Builder pattern to allow method chaining for clean,
// readable application configuration.
type AppBuilder struct {
	// basePath holds the base directory path for the application
	basePath string

	// environment holds the application environment
	environment string

	// debug indicates whether debug mode should be enabled
	debug bool

	// name holds the application name
	name string

	// version holds the application version
	version string

	// locale holds the application locale
	locale string

	// fallbackLocale holds the fallback locale
	fallbackLocale string

	// timezone holds the application timezone
	timezone string

	// runningInConsole indicates console mode
	runningInConsole bool

	// runningUnitTests indicates test mode
	runningUnitTests bool

	// shutdownTimeout specifies graceful shutdown timeout
	shutdownTimeout time.Duration

	// customContainer allows injection of a custom container
	customContainer *container.ServiceContainer

	// serviceProviders holds service providers to register with the application
	serviceProviders []providerInterfaces.ServiceProviderInterface
}

// NewApp creates a new AppBuilder with sensible defaults.
// This is the starting point for fluent application configuration.
//
// Returns:
//
//	*AppBuilder: A new builder instance with default values
//
// Example:
//
//	application := NewApp().
//	    WithName("My API").
//	    WithEnvironment("production").
//	    Build()
func NewApp() *AppBuilder {
	// Initialize environment helper for reading env vars with fallbacks
	envHelper := helpers.NewEnvHelper()

	return &AppBuilder{
		basePath:         ".",                              // Will be resolved to actual path in Build()
		environment:      envHelper.GetAppEnvironment(),    // From APP_ENV or default
		debug:            envHelper.GetAppDebug(),          // From APP_DEBUG or default
		name:             envHelper.GetAppName(),           // From APP_NAME or default
		version:          envHelper.GetAppVersion(),        // From APP_VERSION or default
		locale:           envHelper.GetAppLocale(),         // From APP_LOCALE or default
		fallbackLocale:   envHelper.GetAppFallbackLocale(), // From APP_FALLBACK_LOCALE or default
		timezone:         envHelper.GetAppTimezone(),       // From APP_TIMEZONE or default
		runningInConsole: envHelper.GetRunningInConsole(),  // From APP_CONSOLE or default
		runningUnitTests: envHelper.GetRunningUnitTests(),  // From APP_TESTING or default
		shutdownTimeout:  envHelper.GetShutdownTimeout(),   // From APP_SHUTDOWN_TIMEOUT or default
	}
}

// WithBasePath sets the base directory path for the application.
//
// Parameters:
//
//	path: The base directory path
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithBasePath("/opt/myapp").Build()
func (b *AppBuilder) WithBasePath(path string) *AppBuilder {
	b.basePath = path
	return b
}

// WithEnvironment sets the application environment.
//
// Parameters:
//
//	env: The environment name (development, testing, staging, production)
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithEnvironment("production").Build()
func (b *AppBuilder) WithEnvironment(env string) *AppBuilder {
	b.environment = env
	// Automatically disable debug in production
	if env == string(enums.EnvironmentProduction) {
		b.debug = false
	}
	return b
}

// WithDebug sets the debug mode state.
//
// Parameters:
//
//	debug: Whether to enable debug mode
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithDebug(false).Build()
func (b *AppBuilder) WithDebug(debug bool) *AppBuilder {
	b.debug = debug
	return b
}

// WithName sets the application name.
//
// Parameters:
//
//	name: The application name
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithName("My API Server").Build()
func (b *AppBuilder) WithName(name string) *AppBuilder {
	b.name = name
	return b
}

// WithVersion sets the application version.
//
// Parameters:
//
//	version: The application version
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithVersion("2.1.0").Build()
func (b *AppBuilder) WithVersion(version string) *AppBuilder {
	b.version = version
	return b
}

// WithLocale sets the application locale.
//
// Parameters:
//
//	locale: The application locale
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithLocale("fr").Build()
func (b *AppBuilder) WithLocale(locale string) *AppBuilder {
	b.locale = locale
	return b
}

// WithFallbackLocale sets the fallback locale.
//
// Parameters:
//
//	fallbackLocale: The fallback locale
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithFallbackLocale("en").Build()
func (b *AppBuilder) WithFallbackLocale(fallbackLocale string) *AppBuilder {
	b.fallbackLocale = fallbackLocale
	return b
}

// WithTimezone sets the application timezone.
//
// Parameters:
//
//	timezone: The application timezone
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithTimezone("America/New_York").Build()
func (b *AppBuilder) WithTimezone(timezone string) *AppBuilder {
	b.timezone = timezone
	return b
}

// InConsole marks the application as running in console mode.
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().InConsole().Build()
func (b *AppBuilder) InConsole() *AppBuilder {
	b.runningInConsole = true
	return b
}

// InTesting marks the application as running in test mode.
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().InTesting().Build()
func (b *AppBuilder) InTesting() *AppBuilder {
	b.runningUnitTests = true
	return b
}

// WithShutdownTimeout sets the graceful shutdown timeout.
//
// Parameters:
//
//	timeout: Maximum time to wait for graceful shutdown
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().WithShutdownTimeout(60 * time.Second).Build()
func (b *AppBuilder) WithShutdownTimeout(timeout time.Duration) *AppBuilder {
	b.shutdownTimeout = timeout
	return b
}

// WithContainer injects a custom service container.
//
// Parameters:
//
//	container: Custom container instance
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	customContainer := container.New()
//	application := NewApp().WithContainer(customContainer).Build()
func (b *AppBuilder) WithContainer(container *container.ServiceContainer) *AppBuilder {
	b.customContainer = container
	return b
}

// ForProduction configures the application for production environment.
// This is a convenience method that sets multiple production-appropriate values.
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().ForProduction().Build()
func (b *AppBuilder) ForProduction() *AppBuilder {
	return b.WithEnvironment("production").
		WithDebug(false).
		WithShutdownTimeout(60 * time.Second)
}

// ForDevelopment configures the application for development environment.
// This is a convenience method that sets multiple development-appropriate values.
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().ForDevelopment().Build()
func (b *AppBuilder) ForDevelopment() *AppBuilder {
	return b.WithEnvironment("development").
		WithDebug(true).
		WithShutdownTimeout(10 * time.Second)
}

// ForTesting configures the application for testing environment.
// This is a convenience method that sets multiple test-appropriate values.
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	application := NewApp().ForTesting().Build()
func (b *AppBuilder) ForTesting() *AppBuilder {
	return b.WithEnvironment("testing").
		WithDebug(true).
		InTesting().
		WithShutdownTimeout(5 * time.Second)
}

// WithServiceProviders sets the service providers to register with the application.
// These providers will be registered after the application is built.
//
// Parameters:
//
//	providers: A slice of service provider instances to register
//
// Returns:
//
//	*AppBuilder: The builder instance for method chaining
//
// Example:
//
//	providers := []providerInterfaces.ServiceProviderInterface{
//	    userProviders.NewUserServiceProvider(),
//	    productProviders.NewProductServiceProvider(),
//	}
//	application := NewApp().WithServiceProviders(providers).Build()
func (b *AppBuilder) WithServiceProviders(providers []providerInterfaces.ServiceProviderInterface) *AppBuilder {
	b.serviceProviders = providers
	return b
}

// Build creates and returns the configured Application instance.
// This is the final method in the builder chain and creates
// the actual application with all configured values.
//
// Returns:
//
//	*Application: A fully configured application instance
//
// Example:
//
//	application := NewApp().
//	    WithName("My API").
//	    WithEnvironment("production").
//	    Build()
func (b *AppBuilder) Build() *application.Application {
	// Create the base application (use existing constructor but override values)
	app := application.New()

	// Set base path if different from default
	if b.basePath != "." {
		app.SetBasePath(b.basePath)
	}

	// Apply all builder configurations using trait methods
	app.SetEnvironment(b.environment)
	app.SetDebug(b.debug)
	app.SetName(b.name)
	app.SetVersion(b.version)
	app.SetLocale(b.locale)
	app.SetFallbackLocale(b.fallbackLocale)
	app.SetTimezone(b.timezone)
	app.SetRunningInConsole(b.runningInConsole)
	app.SetRunningUnitTests(b.runningUnitTests)
	app.SetShutdownTimeout(b.shutdownTimeout)

	// Register service providers if any were provided
	if len(b.serviceProviders) > 0 {
		// Convert interface{} slice to individual provider registrations
		for _, provider := range b.serviceProviders {
			if err := app.RegisterProvider(provider); err != nil {
				// Log the error but don't panic - let the application start
				app.GetLogger().Error("Failed to register service provider %T during build: %v", provider, err)
			}
		}
	}

	return app
}
