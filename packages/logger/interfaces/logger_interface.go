package interfaces

// LoggerInterface defines the contract for logging functionality.
// This interface provides a standardized logging API with level-based logging,
// contextual fields, and configurable outputs.
//
// The logger supports multiple log levels and allows for structured logging
// with contextual fields, making it suitable for both development and production
// environments.
//
// Example usage:
//
//	logger := &Logger{}
//	logger.Info("Application started")
//	logger.Error("Database connection failed: %v", err)
//	logger.WithField("user_id", 123).Info("User logged in")
//	logger.WithFields(map[string]interface{}{
//		"request_id": "abc123",
//		"endpoint": "/api/users",
//	}).Warn("Rate limit exceeded")
//
// The interface promotes:
// - Standardized logging across the application
// - Level-based log filtering
// - Structured logging with context
// - Easy testing with mock loggers
// - Flexible logger implementations
type LoggerInterface interface {
	// Debug logs a message at debug level.
	// Debug messages are typically used for detailed diagnostic information
	// that is most useful when diagnosing problems.
	//
	// Parameters:
	//   format: Printf-style format string
	//   args: Arguments for the format string
	//
	// Example:
	//   logger.Debug("Processing user %d with email %s", userID, email)
	Debug(format string, args ...interface{})

	// Info logs a message at info level.
	// Info messages are used for general informational messages that
	// highlight the progress of the application.
	//
	// Parameters:
	//   format: Printf-style format string
	//   args: Arguments for the format string
	//
	// Example:
	//   logger.Info("Application started successfully")
	Info(format string, args ...interface{})

	// Warn logs a message at warning level.
	// Warning messages are used for potentially harmful situations or
	// unexpected conditions that don't prevent the application from working.
	//
	// Parameters:
	//   format: Printf-style format string
	//   args: Arguments for the format string
	//
	// Example:
	//   logger.Warn("Failed to connect to cache, falling back to database")
	Warn(format string, args ...interface{})

	// Error logs a message at error level.
	// Error messages are used for error events that might still allow
	// the application to continue running.
	//
	// Parameters:
	//   format: Printf-style format string
	//   args: Arguments for the format string
	//
	// Example:
	//   logger.Error("Failed to process order %d: %v", orderID, err)
	Error(format string, args ...interface{})

	// Fatal logs a message at fatal level and then calls os.Exit(1).
	// Fatal messages are used for severe error events that will lead
	// to application termination.
	//
	// Parameters:
	//   format: Printf-style format string
	//   args: Arguments for the format string
	//
	// Example:
	//   logger.Fatal("Unable to connect to database: %v", err)
	Fatal(format string, args ...interface{})

	// WithField returns a new logger instance with a single contextual field.
	// This allows for adding structured context to log messages.
	//
	// Parameters:
	//   key: The field key
	//   value: The field value
	//
	// Returns:
	//   LoggerInterface: A new logger instance with the added field
	//
	// Example:
	//   contextLogger := logger.WithField("user_id", 123)
	//   contextLogger.Info("User action performed")
	WithField(key string, value interface{}) LoggerInterface

	// WithFields returns a new logger instance with multiple contextual fields.
	// This allows for adding multiple structured context fields to log messages.
	//
	// Parameters:
	//   fields: A map of field keys to values
	//
	// Returns:
	//   LoggerInterface: A new logger instance with the added fields
	//
	// Example:
	//   contextLogger := logger.WithFields(map[string]interface{}{
	//       "user_id": 123,
	//       "request_id": "abc123",
	//       "endpoint": "/api/users",
	//   })
	//   contextLogger.Info("API request processed")
	WithFields(fields map[string]interface{}) LoggerInterface

	// SetLevel sets the logging level for the logger.
	// Messages below this level will not be logged.
	//
	// Common levels: "debug", "info", "warn", "error", "fatal"
	//
	// Parameters:
	//   level: The minimum logging level
	//
	// Example:
	//   logger.SetLevel("warn") // Only warn, error, and fatal messages will be logged
	SetLevel(level string)

	// GetLevel returns the current logging level.
	//
	// Returns:
	//   string: The current logging level
	//
	// Example:
	//   currentLevel := logger.GetLevel()
	//   if currentLevel == "debug" {
	//       // Debug logging is enabled
	//   }
	GetLevel() string
}
