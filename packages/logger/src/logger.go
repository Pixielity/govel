// Package logger provides logging functionality for the GoVel framework.
// This package implements a flexible logging system with multiple levels,
// output destinations, and formatting options.
//
// The logger supports:
// - Multiple log levels (Debug, Info, Warning, Error, Fatal)
// - Custom formatters and outputs
// - Thread-safe operations
// - Context-based logging
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	interfaces "govel/types/interfaces/logger"
)

// LogLevel represents the logging level.
type LogLevel int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel LogLevel = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance with configurable level and output.
// The logger is thread-safe and can be used concurrently from multiple goroutines.
//
// Example usage:
//
//	logger := logger.New()
//	logger.Info("Application started")
//	logger.Error("Database connection failed: %v", err)
//	logger.WithFields(map[string]interface{}{
//		"user_id": 123,
//		"action":  "login",
//	}).Info("User logged in")
type Logger struct {
	// level is the minimum log level that will be output
	level LogLevel

	// output is the destination for log output
	output io.Writer

	// logger is the underlying Go logger instance
	logger *log.Logger

	// mutex provides thread-safe access to logger state
	mutex sync.RWMutex

	// fields contains contextual fields to be included in log entries
	fields map[string]interface{}
}

// New creates a new logger instance with default configuration.
// Default configuration:
// - Level: InfoLevel
// - Output: os.Stdout
// - Format: Standard Go log format with timestamp
//
// Returns:
//
//	*Logger: A new logger instance ready for use
//
// Example:
//
//	logger := logger.New()
//	logger.Info("Logger initialized")
func New() *Logger {
	return &Logger{
		level:  InfoLevel,
		output: os.Stdout,
		logger: log.New(os.Stdout, "", log.LstdFlags),
		fields: make(map[string]interface{}),
	}
}

// NewWithOutput creates a new logger instance with custom output destination.
//
// Parameters:
//
//	output: The io.Writer to write log entries to
//
// Returns:
//
//	*Logger: A new logger instance with custom output
//
// Example:
//
//	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		panic(err)
//	}
//	logger := logger.NewWithOutput(file)
func NewWithOutput(output io.Writer) *Logger {
	return &Logger{
		level:  InfoLevel,
		output: output,
		logger: log.New(output, "", log.LstdFlags),
		fields: make(map[string]interface{}),
	}
}

// SetLevelEnum sets the minimum log level for this logger using LogLevel enum.
// Log entries below this level will be discarded.
//
// Parameters:
//
//	level: The minimum log level to output
//
// Example:
//
//	logger.SetLevelEnum(logger.DebugLevel) // Enable debug logging
//	logger.SetLevelEnum(logger.ErrorLevel) // Only show errors and fatal
func (l *Logger) SetLevelEnum(level LogLevel) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// SetLevel sets the logging level using string representation (interface-compatible method).
// This method is required by the LoggerInterface.
//
// Parameters:
//
//	levelStr: The logging level as a string ("debug", "info", "warn", "error", "fatal")
//
// Example:
//
//	logger.SetLevel("debug") // Enable debug logging
//	logger.SetLevel("error") // Only show errors and fatal
func (l *Logger) SetLevel(levelStr string) {
	var level LogLevel
	switch levelStr {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn", "warning":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	default:
		level = InfoLevel // default to info level
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// GetLevel returns the current logging level as string (interface-compatible method).
//
// Returns:
//
//	string: The current logging level in lowercase
func (l *Logger) GetLevel() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	switch l.level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "info"
	}
}

// SetOutput sets the output destination for this logger.
//
// Parameters:
//
//	output: The io.Writer to write log entries to
//
// Example:
//
//	logger.SetOutput(os.Stderr) // Log to stderr instead of stdout
func (l *Logger) SetOutput(output io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.output = output
	l.logger = log.New(output, "", log.LstdFlags)
}

// WithFields returns a new logger instance with additional contextual fields.
// The original logger is not modified.
//
// Parameters:
//
//	fields: Map of field names to values to include in log entries
//
// Returns:
//
//	interfaces.LoggerInterface: A new logger instance with additional fields
//
// Example:
//
//	requestLogger := logger.WithFields(map[string]interface{}{
//		"request_id": "req-123",
//		"user_id":    456,
//	})
//	requestLogger.Info("Processing request")
func (l *Logger) WithFields(fields map[string]interface{}) interfaces.LoggerInterface {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	// Create a new logger instance with merged fields
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		level:  l.level,
		output: l.output,
		logger: l.logger,
		fields: newFields,
	}
}

// WithField returns a new logger instance with an additional contextual field.
// This is a convenience method for adding a single field.
//
// Parameters:
//
//	key: The field name
//	value: The field value
//
// Returns:
//
//	interfaces.LoggerInterface: A new logger instance with the additional field
//
// Example:
//
//	userLogger := logger.WithField("user_id", 123)
//	userLogger.Info("User action performed")
func (l *Logger) WithField(key string, value interface{}) interfaces.LoggerInterface {
	return l.WithFields(map[string]interface{}{key: value})
}

// Debug logs a message at debug level.
// Debug messages are typically only enabled during development.
//
// Parameters:
//
//	format: Printf-style format string
//	args: Arguments for the format string
//
// Example:
//
//	logger.Debug("Processing item %d of %d", current, total)
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs a message at info level.
// Info messages represent general application flow.
//
// Parameters:
//
//	format: Printf-style format string
//	args: Arguments for the format string
//
// Example:
//
//	logger.Info("Server started on port %d", port)
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a message at warning level.
// Warning messages indicate something unexpected but not necessarily problematic.
//
// Parameters:
//
//	format: Printf-style format string
//	args: Arguments for the format string
//
// Example:
//
//	logger.Warn("Deprecated configuration option used: %s", option)
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs a message at error level.
// Error messages indicate something has gone wrong but the application can continue.
//
// Parameters:
//
//	format: Printf-style format string
//	args: Arguments for the format string
//
// Example:
//
//	logger.Error("Failed to connect to database: %v", err)
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal logs a message at fatal level and exits the program.
// Fatal messages indicate the application cannot continue and must terminate.
//
// Parameters:
//
//	format: Printf-style format string
//	args: Arguments for the format string
//
// Example:
//
//	logger.Fatal("Critical system failure: %v", err)
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
	os.Exit(1)
}

// log is the internal logging method that handles level checking and formatting.
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mutex.RLock()
	currentLevel := l.level
	fields := l.fields
	logger := l.logger
	l.mutex.RUnlock()

	// Check if we should log at this level
	if level < currentLevel {
		return
	}

	// Format the message
	message := fmt.Sprintf(format, args...)

	// Add fields to the message if any
	if len(fields) > 0 {
		fieldsStr := ""
		for k, v := range fields {
			if fieldsStr != "" {
				fieldsStr += " "
			}
			fieldsStr += fmt.Sprintf("%s=%v", k, v)
		}
		message = fmt.Sprintf("[%s] %s [%s]", level.String(), message, fieldsStr)
	} else {
		message = fmt.Sprintf("[%s] %s", level.String(), message)
	}

	// Output the log entry
	logger.Println(message)
}

// IsDebugEnabled returns true if debug logging is enabled.
//
// Returns:
//
//	bool: true if debug logging is enabled, false otherwise
//
// Example:
//
//	if logger.IsDebugEnabled() {
//		// Expensive debug operation
//		logger.Debug("Debug info: %s", expensiveDebugData())
//	}
func (l *Logger) IsDebugEnabled() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.level <= DebugLevel
}

// IsInfoEnabled returns true if info logging is enabled.
//
// Returns:
//
//	bool: true if info logging is enabled, false otherwise
func (l *Logger) IsInfoEnabled() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.level <= InfoLevel
}

// IsWarnEnabled returns true if warning logging is enabled.
//
// Returns:
//
//	bool: true if warning logging is enabled, false otherwise
func (l *Logger) IsWarnEnabled() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.level <= WarnLevel
}

// IsErrorEnabled returns true if error logging is enabled.
//
// Returns:
//
//	bool: true if error logging is enabled, false otherwise
func (l *Logger) IsErrorEnabled() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.level <= ErrorLevel
}

// Flush ensures all pending log entries are written to the output destination.
// This is useful when using buffered writers or when gracefully shutting down.
//
// Returns:
//
//	error: Any error that occurred during flushing, nil if successful
//
// Example:
//
//	defer logger.Flush()
func (l *Logger) FlushLogger() error {
	// If the output implements a Flush method, call it
	if flusher, ok := l.output.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// Compile-time interface compliance checks
// These ensure Logger properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ interfaces.LoggerInterface = (*Logger)(nil) // Direct logging operations
