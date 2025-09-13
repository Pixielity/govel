package interfaces

// LoggerInterface defines the contract for logger operations.
// This interface provides comprehensive logging functionality including
// multiple log levels, contextual fields, and configuration management.
//
// Key features:
//   - Multiple log levels (Debug, Info, Warning, Error, Fatal)
//   - Contextual field support for structured logging
//   - Thread-safe operations
//   - Configurable log levels and output destinations
//   - Method chaining for fluent API
type LoggerInterface interface {
	// Debug logs a message at debug level.
	// Debug messages are typically only enabled during development.
	//
	// Parameters:
	//   - format: Printf-style format string
	//   - args: Arguments for the format string
	Debug(format string, args ...interface{})

	// Info logs a message at info level.
	// Info messages represent general application flow.
	//
	// Parameters:
	//   - format: Printf-style format string
	//   - args: Arguments for the format string
	Info(format string, args ...interface{})

	// Warn logs a message at warning level.
	// Warning messages indicate something unexpected but not necessarily problematic.
	//
	// Parameters:
	//   - format: Printf-style format string
	//   - args: Arguments for the format string
	Warn(format string, args ...interface{})

	// Error logs a message at error level.
	// Error messages indicate something has gone wrong but the application can continue.
	//
	// Parameters:
	//   - format: Printf-style format string
	//   - args: Arguments for the format string
	Error(format string, args ...interface{})

	// Fatal logs a message at fatal level and exits the program.
	// Fatal messages indicate the application cannot continue and must terminate.
	//
	// Parameters:
	//   - format: Printf-style format string
	//   - args: Arguments for the format string
	Fatal(format string, args ...interface{})

	// WithFields returns a new logger instance with additional contextual fields.
	// The original logger is not modified.
	//
	// Parameters:
	//   - fields: Map of field names to values to include in log entries
	//
	// Returns:
	//   - LoggerInterface: A new logger instance with additional fields
	WithFields(fields map[string]interface{}) LoggerInterface

	// WithField returns a new logger instance with an additional contextual field.
	// This is a convenience method for adding a single field.
	//
	// Parameters:
	//   - key: The field name
	//   - value: The field value
	//
	// Returns:
	//   - LoggerInterface: A new logger instance with the additional field
	WithField(key string, value interface{}) LoggerInterface

	// SetLevel sets the logging level using string representation.
	// Controls which log messages are output based on severity.
	//
	// Parameters:
	//   - level: The logging level as a string ("debug", "info", "warn", "error", "fatal")
	SetLevel(level string)

	// GetLevel returns the current logging level as string.
	// Returns the minimum level that will be output.
	//
	// Returns:
	//   - string: The current logging level in lowercase
	GetLevel() string
}
