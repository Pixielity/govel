// Package interfaces - Context interface definitions
// This file contains interfaces related to context management within the middleware system.
// Context management is crucial for passing data, handling cancellation, timeouts,
// and maintaining request-scoped information throughout the middleware chain.
package interfaces

import (
	"context"
	"time"
)

// MiddlewareContext defines the interface for enhanced context management in middleware chains.
// This interface extends the standard Go context with middleware-specific functionality
// for managing request metadata, tracing information, and middleware state.
//
// MiddlewareContext provides:
// - Request-scoped data storage and retrieval
// - Middleware execution tracing and debugging
// - Performance monitoring and metrics collection
// - Error context and stack trace management
// - Security context and authentication state
//
// Design Philosophy:
//   MiddlewareContext follows the principle of explicit context management,
//   where each middleware can inspect, modify, and enhance the context
//   while maintaining immutability and thread safety.
type MiddlewareContext interface {
	context.Context // Embed standard Go context

	// WithValue creates a new context with the given key-value pair.
	// This method follows the same semantics as context.WithValue but
	// provides type safety and middleware-specific optimizations.
	//
	// Parameters:
	//   - key: The key for the value (should be comparable)
	//   - value: The value to associate with the key
	//
	// Returns:
	//   - MiddlewareContext: New context with the added value
	WithValue(key interface{}, value interface{}) MiddlewareContext

	// GetValue retrieves a value from the context with type safety.
	// This method provides a safer alternative to direct context.Value() calls.
	//
	// Parameters:
	//   - key: The key to look up
	//
	// Returns:
	//   - interface{}: The value associated with the key
	//   - bool: true if the key was found
	GetValue(key interface{}) (interface{}, bool)

	// WithTimeout creates a new context with the specified timeout.
	// This method wraps context.WithTimeout but maintains the MiddlewareContext interface.
	//
	// Parameters:
	//   - timeout: Duration after which the context should be cancelled
	//
	// Returns:
	//   - MiddlewareContext: New context with timeout
	//   - context.CancelFunc: Function to cancel the context early
	WithTimeout(timeout time.Duration) (MiddlewareContext, context.CancelFunc)

	// WithCancel creates a new context with cancellation capability.
	// This method wraps context.WithCancel but maintains the MiddlewareContext interface.
	//
	// Returns:
	//   - MiddlewareContext: New context with cancellation
	//   - context.CancelFunc: Function to cancel the context
	WithCancel() (MiddlewareContext, context.CancelFunc)

	// Clone creates a copy of the current context without parent relationships.
	// This is useful for creating independent contexts for parallel processing.
	//
	// Returns:
	//   - MiddlewareContext: Independent copy of the context
	Clone() MiddlewareContext
}

// RequestContext defines the interface for managing request-specific context data.
// This interface provides structured access to common request metadata and
// enables middleware to share information efficiently.
//
// RequestContext is designed for:
// - HTTP request metadata (headers, method, path, etc.)
// - Authentication and authorization information
// - Request tracing and correlation IDs
// - Performance metrics and timing data
// - Custom application-specific data
type RequestContext interface {
	MiddlewareContext // Extend middleware context

	// GetRequestID returns a unique identifier for this request.
	// This ID is used for logging, tracing, and correlation across services.
	//
	// Returns:
	//   - string: Unique request identifier
	GetRequestID() string

	// SetRequestID sets the unique identifier for this request.
	//
	// Parameters:
	//   - requestID: Unique identifier to set
	SetRequestID(requestID string)

	// GetUserID returns the ID of the authenticated user, if available.
	//
	// Returns:
	//   - string: User identifier
	//   - bool: true if user ID is available
	GetUserID() (string, bool)

	// SetUserID sets the authenticated user ID.
	//
	// Parameters:
	//   - userID: User identifier to set
	SetUserID(userID string)

	// GetMetadata returns all request metadata as a map.
	// This provides access to custom metadata set by middleware.
	//
	// Returns:
	//   - map[string]interface{}: Request metadata
	GetMetadata() map[string]interface{}

	// SetMetadata sets a metadata key-value pair.
	//
	// Parameters:
	//   - key: Metadata key
	//   - value: Metadata value
	SetMetadata(key string, value interface{})

	// GetStartTime returns when request processing started.
	//
	// Returns:
	//   - time.Time: Request start time
	GetStartTime() time.Time

	// SetStartTime sets when request processing started.
	//
	// Parameters:
	//   - startTime: Request start time
	SetStartTime(startTime time.Time)

	// GetTraceID returns the distributed tracing ID.
	//
	// Returns:
	//   - string: Trace ID
	//   - bool: true if trace ID is available
	GetTraceID() (string, bool)

	// SetTraceID sets the distributed tracing ID.
	//
	// Parameters:
	//   - traceID: Trace ID to set
	SetTraceID(traceID string)
}

// ExecutionContext defines the interface for managing middleware execution state.
// This interface tracks the execution flow and provides debugging information
// about the middleware chain processing.
//
// ExecutionContext enables:
// - Middleware execution tracing and profiling
// - Error context and debugging information
// - Performance monitoring and bottleneck detection
// - Audit trails and compliance logging
// - Circuit breaker and resilience patterns
type ExecutionContext interface {
	MiddlewareContext // Extend middleware context

	// GetExecutionStack returns the current middleware execution stack.
	// This is useful for debugging and understanding execution flow.
	//
	// Returns:
	//   - []string: Names of middleware in execution order
	GetExecutionStack() []string

	// PushMiddleware adds a middleware to the execution stack.
	//
	// Parameters:
	//   - middlewareName: Name of the middleware being executed
	PushMiddleware(middlewareName string)

	// PopMiddleware removes the last middleware from the execution stack.
	//
	// Returns:
	//   - string: Name of the popped middleware
	//   - bool: true if a middleware was popped
	PopMiddleware() (string, bool)

	// GetCurrentMiddleware returns the name of the currently executing middleware.
	//
	// Returns:
	//   - string: Current middleware name
	//   - bool: true if there is a current middleware
	GetCurrentMiddleware() (string, bool)

	// RecordError records an error in the execution context.
	// This maintains a history of errors for debugging.
	//
	// Parameters:
	//   - err: The error that occurred
	//   - middlewareName: Name of the middleware where error occurred
	RecordError(err error, middlewareName string)

	// GetErrors returns all errors recorded during execution.
	//
	// Returns:
	//   - []ExecutionError: List of execution errors
	GetErrors() []ExecutionError

	// HasErrors returns true if any errors have been recorded.
	//
	// Returns:
	//   - bool: true if errors exist
	HasErrors() bool

	// GetExecutionDuration returns how long the request has been processing.
	//
	// Returns:
	//   - time.Duration: Time since execution started
	GetExecutionDuration() time.Duration

	// RecordMetric records a performance metric.
	//
	// Parameters:
	//   - name: Metric name
	//   - value: Metric value
	//   - tags: Additional tags for the metric
	RecordMetric(name string, value interface{}, tags map[string]string)

	// GetMetrics returns all recorded metrics.
	//
	// Returns:
	//   - map[string]ExecutionMetric: Recorded metrics
	GetMetrics() map[string]ExecutionMetric
}

// ExecutionError represents an error that occurred during middleware execution.
type ExecutionError struct {
	// Error is the actual error that occurred
	Error error

	// MiddlewareName is the name of the middleware where the error occurred
	MiddlewareName string

	// Timestamp is when the error occurred
	Timestamp time.Time

	// StackTrace contains the call stack when the error occurred
	StackTrace string

	// Context contains additional context about the error
	Context map[string]interface{}
}

// ExecutionMetric represents a performance metric recorded during execution.
type ExecutionMetric struct {
	// Name is the metric name
	Name string

	// Value is the metric value
	Value interface{}

	// Timestamp is when the metric was recorded
	Timestamp time.Time

	// Tags provide additional metadata for the metric
	Tags map[string]string

	// MiddlewareName is the middleware that recorded the metric
	MiddlewareName string
}

// ContextFactory defines the interface for creating context instances.
// This interface enables dependency injection and configuration-based
// context creation with consistent initialization.
//
// ContextFactory is useful for:
// - Creating contexts with pre-configured defaults
// - Injecting logging, tracing, and monitoring capabilities
// - Setting up security and authentication contexts
// - Providing test contexts with mock implementations
type ContextFactory interface {
	// CreateMiddlewareContext creates a new middleware context.
	//
	// Parameters:
	//   - parent: Parent context to derive from
	//
	// Returns:
	//   - MiddlewareContext: New middleware context
	CreateMiddlewareContext(parent context.Context) MiddlewareContext

	// CreateRequestContext creates a new request context.
	//
	// Parameters:
	//   - parent: Parent context to derive from
	//   - requestID: Unique identifier for the request
	//
	// Returns:
	//   - RequestContext: New request context
	CreateRequestContext(parent context.Context, requestID string) RequestContext

	// CreateExecutionContext creates a new execution context.
	//
	// Parameters:
	//   - parent: Parent context to derive from
	//
	// Returns:
	//   - ExecutionContext: New execution context
	CreateExecutionContext(parent context.Context) ExecutionContext

	// WithDefaults configures the factory with default settings.
	//
	// Parameters:
	//   - config: Configuration map for defaults
	//
	// Returns:
	//   - ContextFactory: Configured factory instance
	WithDefaults(config map[string]interface{}) ContextFactory
}
