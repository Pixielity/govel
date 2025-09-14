package interfaces

// LoggerInterface defines the interface for logging within the GoVel compiler system.
// Provides thread-safe logging with multiple levels and printf-style formatting for debugging and monitoring.
type LoggerInterface interface {
	// Debug logs detailed debug information for troubleshooting and development.
	// Typically disabled in production for performance reasons.
	//
	// Parameters:
	//
	//	msg: The format string for the debug message
	//	args: Optional variadic arguments for printf-style formatting
	Debug(msg string, args ...interface{})

	// Info logs general informational messages about normal system operation.
	// Enabled in both development and production for operational visibility.
	//
	// Parameters:
	//
	//	msg: The format string for the informational message
	//	args: Optional variadic arguments for printf-style formatting
	Info(msg string, args ...interface{})

	// Warn logs warning messages about potentially problematic situations.
	// Indicates non-optimal conditions that should be investigated but don't prevent system function.
	//
	// Parameters:
	//
	//	msg: The format string for the warning message
	//	args: Optional variadic arguments for printf-style formatting
	Warn(msg string, args ...interface{})

	// Error logs error messages about conditions that require attention.
	// Indicates serious problems that should be monitored and addressed promptly.
	//
	// Parameters:
	//
	//	msg: The format string for the error message
	//	args: Optional variadic arguments for printf-style formatting
	Error(msg string, args ...interface{})
}
