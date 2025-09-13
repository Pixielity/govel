package facades

import (
	facade "govel/support"
	loggerInterfaces "govel/types/interfaces/logger"
)

// Log provides a clean, static-like interface to the application's logging service.
//
// This facade implements the facade pattern, providing global access to the logger
// service configured in the dependency injection container. It offers a Laravel-style
// API for logging operations with automatic service resolution, caching, and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved logger for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent access across goroutines
//
// Behavior:
//   - First call: Resolves logger from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached logger instance (extremely fast)
//   - Panics if logger cannot be resolved (fail-fast behavior)
//   - Automatically handles service lifecycle and caching
//
// Returns:
//   - LoggerInterface: The application's logger service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or facade.SetContainer()
//   - If "log" service is not registered in the container
//   - If the resolved service doesn't implement LoggerInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Log() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Cache updates are atomic and consistent
//
// Usage Examples:
//
//	// Basic logging at different levels
//	facades.Log().Debug("Detailed debugging info: %+v", debugData)
//	facades.Log().Info("User %d logged in successfully", userID)
//	facades.Log().Warn("API rate limit approaching: %d/%d", current, limit)
//	facades.Log().Error("Database query failed: %v", err)
//	facades.Log().Fatal("Critical system failure: %v", criticalErr)
//
//	// Structured logging with contextual fields
//	facades.Log().WithField("user_id", 123).Info("User updated profile")
//	facades.Log().WithField("request_id", reqID).Error("Request processing failed")
//
//	// Multiple contextual fields for rich logging
//	facades.Log().WithFields(map[string]interface{}{
//	    "user_id":    123,
//	    "session_id": "abc-123",
//	    "ip_address": "192.168.1.1",
//	    "user_agent": "Mozilla/5.0...",
//	}).Info("User authentication successful")
//
//	// HTTP request logging
//	facades.Log().WithFields(map[string]interface{}{
//	    "method":      r.Method,
//	    "url":         r.URL.String(),
//	    "status_code": statusCode,
//	    "duration_ms": duration.Milliseconds(),
//	    "bytes":       responseBytes,
//	}).Info("HTTP request processed")
//
//	// Error logging with stack traces and context
//	if err != nil {
//	    facades.Log().WithFields(map[string]interface{}{
//	        "error":      err.Error(),
//	        "stack_trace": fmt.Sprintf("%+v", err), // If using pkg/errors
//	        "function":    "ProcessUser",
//	        "user_id":     userID,
//	    }).Error("User processing failed")
//	}
//
//	// Conditional logging based on environment
//	if config.Environment == "development" {
//	    facades.Log().Debug("Development-only debug info: %+v", sensitiveData)
//	}
//
// Best Practices:
//   - Use appropriate log levels (Debug < Info < Warn < Error < Fatal)
//   - Include relevant context with WithField/WithFields for better searchability
//   - Use structured logging for easier parsing and analysis
//   - Avoid logging sensitive information (passwords, tokens, etc.)
//   - Use consistent field names across the application
//   - Consider performance impact of logging in hot code paths
//
// Integration with Monitoring:
// The logger can be integrated with various monitoring and logging systems:
//
//	// Example: Structured logging for ELK stack
//	facades.Log().WithFields(map[string]interface{}{
//	    "service":     "user-service",
//	    "version":     "1.2.3",
//	    "environment": "production",
//	    "trace_id":    traceID,
//	}).Info("Service operation completed")
//
// Container Configuration:
// Ensure the logger service is properly configured in your container:
//
//	// Example logger registration
//	container.Singleton("log", func() interface{} {
//	    config := logger.Config{
//	        Level:      "info",
//	        Format:     "json",         // json, text, or custom
//	        Output:     os.Stdout,      // or file, syslog, etc.
//	        TimeFormat: time.RFC3339,   // consistent timestamp format
//	        Fields: map[string]interface{}{
//	            "service": "my-application",
//	            "version": version.Get(),
//	        },
//	    }
//	    return logger.NewLogger(config)
//	})
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserService(t *testing.T) {
//	    // Create a test logger that captures logs
//	    testLogger := &TestLogger{}
//
//	    // Swap the real logger with test logger
//	    restore := facade.SwapService("log", testLogger)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Log() returns testLogger
//	    userService := NewUserService()
//	    userService.CreateUser(userData)
//
//	    // Verify logging behavior
//	    logs := testLogger.GetLogs()
//	    assert.Contains(t, logs, "User created successfully")
//	    assert.Equal(t, "info", logs[0].Level)
//	}
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume logging always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	logger, err := facade.TryResolve[LoggerInterface]("log")
//	if err != nil {
//	    // Handle logger unavailability gracefully
//	    fmt.Printf("Logger not available: %v\n", err)
//	    return
//	}
//	logger.Info("Using error-safe logger access")
func Log() loggerInterfaces.LoggerInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "log" service from the dependency injection container
	// - Performs type assertion to LoggerInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[loggerInterfaces.LoggerInterface](loggerInterfaces.LOGGER_TOKEN)
}

// LogWithError provides error-safe access to the logger service.
//
// This function offers the same functionality as Log() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle logger unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Log() but with error handling.
//
// Returns:
//   - LoggerInterface: The resolved logger instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - facade.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement LoggerInterface
//
// Usage Examples:
//
//	// Basic error-safe logging
//	logger, err := facades.LogWithError()
//	if err != nil {
//	    fmt.Printf("Logger unavailable: %v\n", err)
//	    return // or use alternative logging
//	}
//	logger.Info("Application started")
//
//	// Conditional logging
//	if logger, err := facades.LogWithError(); err == nil {
//	    logger.Debug("Optional debug information")
//	}
func LogWithError() (loggerInterfaces.LoggerInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "log" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[loggerInterfaces.LoggerInterface](loggerInterfaces.LOGGER_TOKEN)
}
