package providers

import (
	"fmt"
	"os"

	serviceProviders "govel/packages/application/providers"
	"govel/packages/logger"
	applicationInterfaces "govel/packages/types/src/interfaces/application"
	interfaces "govel/packages/types/src/interfaces/logger"
	loggerInterfaces "govel/packages/types/src/interfaces/logger"
)

/**
 * LoggerServiceProvider provides comprehensive logging services with Laravel-style configuration.
 *
 * This service provider implements real-world logging functionality, binding the LoggerInterface
 * to the concrete Logger implementation. It provides configurable logging levels, output destinations,
 * and context-aware logging for development and production environments.
 *
 * Features:
 * - Multi-level logging (Debug, Info, Warn, Error, Fatal)
 * - Configurable output destinations (stdout, files, custom writers)
 * - Environment-specific logging configuration
 * - Context-aware logging with structured fields
 * - Thread-safe operations for concurrent usage
 * - Configuration-driven logger setup
 * - Development vs production logging modes
 * - Custom formatter support
 * - Log level filtering and management
 *
 * Binding Strategy:
 * - Binds "logger" abstract to LoggerInterface implementation
 * - Registers as singleton for application-wide logging consistency
 * - Provides factory for creating context-specific loggers
 * - Supports configuration-based logger initialization
 *
 * Similar to Laravel's LogServiceProvider, this implementation:
 * - Configures logging based on environment and configuration
 * - Provides single logger instance across the application
 * - Supports multiple logging channels and drivers
 * - Integrates with application configuration system
 * - Provides structured logging capabilities
 */
type LoggerServiceProvider struct {
	serviceProviders.ServiceProvider
}

// NewLoggerServiceProvider creates a new logger service provider with default settings.
// This constructor initializes the provider ready for registration with the application.
//
// Returns:
//
//	*LoggerServiceProvider: A new logging service provider ready for registration
//
// Example:
//
//	loggerProvider := providers.NewLoggerServiceProvider()
//	app.RegisterProvider(loggerProvider)
func NewLoggerServiceProvider() *LoggerServiceProvider {
	return &LoggerServiceProvider{
		ServiceProvider: serviceProviders.ServiceProvider{},
	}
}

// Register binds the logging service into the application container.
// This method implements Laravel-style service registration, binding the LoggerInterface
// abstract to the concrete Logger implementation as a singleton service.
//
// Registration Process:
// 1. Calls parent registration to set provider state
// 2. Binds "logger" abstract to LoggerInterface implementation
// 3. Configures logger based on environment variables and settings
// 4. Sets up logging level and output destination
// 5. Registers logger factory for context-specific instances
//
// The logging service is registered as a singleton to ensure:
// - Consistent logging behavior across the application
// - Memory efficiency by sharing the same logger instance
// - Centralized log level and configuration management
// - Thread-safe access to logging functionality
//
// Parameters:
//
//	application: The application instance with service container access
//
// Returns:
//
//	error: Any error that occurred during registration, nil if successful
//
// Post-Registration Usage:
//
//	// Resolve logger service from container
//	loggerService, err := application.Make("logger")
//	if err != nil {
//	    return fmt.Errorf("failed to resolve logger service: %w", err)
//	}
//
//	// Cast to interface and use
//	logger := loggerService.(loggerInterfaces.LoggerInterface)
//	logger.Info("Application started successfully")
//	logger.Error("Database connection failed: %v", err)
//	logger.WithField("user_id", 123).Info("User authenticated")
func (p *LoggerServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register the logging service as a singleton
	// This binds the logger token to the LoggerInterface implementation
	if err := application.Singleton(interfaces.LOGGER_TOKEN, p.createLoggerFactory()); err != nil {
		return fmt.Errorf("failed to register logger singleton: %w", err)
	}

	// Register logger factory for creating context-specific logger instances
	if err := application.Bind(interfaces.LOGGER_FACTORY_TOKEN, p.createLoggerFactoryMethod()); err != nil {
		return fmt.Errorf("failed to register logger factory: %w", err)
	}

	// Register current log level resolver for debugging and monitoring
	if err := application.Bind(interfaces.LOGGER_LEVEL_TOKEN, func() interface{} {
		return p.getConfiguredLogLevel()
	}); err != nil {
		return fmt.Errorf("failed to register logger level resolver: %w", err)
	}

	return nil
}

// createLoggerFactory creates the main logger service factory function.
// This factory is responsible for creating and configuring the Logger instance
// with environment-specific settings, output destinations, and log levels.
//
// Returns:
//
//	func() interface{}: Factory function that creates LoggerInterface instance
func (p *LoggerServiceProvider) createLoggerFactory() func() interface{} {
	return func() interface{} {
		// Create a new logger instance
		loggerInstance := logger.New()

		// Configure the logger based on environment and configuration
		p.configureLogger(loggerInstance)

		// Return as LoggerInterface to maintain interface segregation
		return loggerInterface(loggerInstance)
	}
}

// createLoggerFactoryMethod creates a factory method for creating additional logger instances.
// This is useful for creating context-specific or scoped logger instances.
//
// Returns:
//
//	func() interface{}: Factory method that returns a logger creation function
func (p *LoggerServiceProvider) createLoggerFactoryMethod() func() interface{} {
	return func() interface{} {
		return func(level string) loggerInterfaces.LoggerInterface {
			loggerInstance := logger.New()
			loggerInstance.SetLevel(level)
			return loggerInterface(loggerInstance)
		}
	}
}

// configureLogger configures the logger instance based on environment variables and settings.
// This method sets up logging level, output format, and other logger-specific configuration.
//
// Parameters:
//
//	loggerInstance: The logger instance to configure
func (p *LoggerServiceProvider) configureLogger(loggerInstance *logger.Logger) {
	// Set logging level based on environment
	logLevel := p.getConfiguredLogLevel()
	loggerInstance.SetLevel(logLevel)

	// Configure output destination based on environment
	// In production, you might want to log to files or external services
	// In development, logging to stdout is usually preferred
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("GOVEL_ENV")
	}
	if env == "" {
		env = "development"
	}

	// Set appropriate log level for environment
	switch env {
	case "production":
		// Production: Only log important messages
		if logLevel == "" {
			loggerInstance.SetLevel("warn")
		}
	case "testing":
		// Testing: Minimal logging to avoid cluttering test output
		if logLevel == "" {
			loggerInstance.SetLevel("error")
		}
	default:
		// Development: Verbose logging for debugging
		if logLevel == "" {
			loggerInstance.SetLevel("debug")
		}
	}
}

// getConfiguredLogLevel determines the appropriate logging level from environment variables.
// This method checks various environment variables to determine the desired log level.
//
// Returns:
//
//	string: The configured log level
func (p *LoggerServiceProvider) getConfiguredLogLevel() string {
	// Check various environment variable patterns
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = os.Getenv("APP_LOG_LEVEL")
	}
	if logLevel == "" {
		logLevel = os.Getenv("GOVEL_LOG_LEVEL")
	}

	// Default to info level if not specified
	if logLevel == "" {
		logLevel = "info"
	}

	return logLevel
}

// loggerInterface is a type assertion helper that ensures the logger instance
// implements the LoggerInterface. This provides compile-time safety for the binding.
//
// Parameters:
//
//	loggerInstance: The concrete Logger instance
//
// Returns:
//
//	loggerInterfaces.LoggerInterface: The logger instance as an interface
func loggerInterface(loggerInstance *logger.Logger) loggerInterfaces.LoggerInterface {
	return loggerInstance
}

// Priority returns the registration priority for this service provider.
// Logger services have high priority since many other services depend on logging.
//
// Returns:
//
//	int: Priority level 50 for high-priority infrastructure services
func (p *LoggerServiceProvider) Priority() int {
	return 50 // High priority - logging is fundamental infrastructure
}
