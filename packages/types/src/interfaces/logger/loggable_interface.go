package interfaces

// LoggableInterface defines the contract for objects that can manage logger instances.
// This interface provides logger management functionality for traits and other
// objects that need to maintain logger state.
//
// Key features:
//   - Logger instance management (get, set, check)
//   - Logger information retrieval
//   - Support for dependency injection and factory patterns
//   - Thread-safe logger state management
type LoggableInterface interface {
	// GetLogger returns the logger instance.
	// Provides access to the current logger for logging operations.
	//
	// Returns:
	//   - LoggerInterface: The current logger instance, may be nil if not set
	GetLogger() LoggerInterface

	// SetLogger sets the logger instance.
	// Updates the current logger, typically used for dependency injection.
	//
	// Parameters:
	//   - logger: The logger instance to set (interface{} for flexibility)
	SetLogger(logger interface{})

	// HasLogger returns whether a logger instance is set.
	// Useful for checking logger availability before attempting operations.
	//
	// Returns:
	//   - bool: true if a logger is configured, false otherwise
	HasLogger() bool

	// GetLoggerInfo returns information about the logger.
	// Provides metadata about logger configuration and state.
	//
	// Returns:
	//   - map[string]interface{}: Logger information including type, level, etc.
	GetLoggerInfo() map[string]interface{}
}
