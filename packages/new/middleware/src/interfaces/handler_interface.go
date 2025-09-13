// Package interfaces - Handler interface definitions
// This file contains interfaces related to request handlers that work with the middleware system.
// Handlers represent the final processing step in a middleware chain and define how requests
// are ultimately processed to produce responses.
package interfaces

import (
	"context"
	"time"
)

// RequestHandler defines the interface for components that can handle requests.
// This interface provides a more comprehensive contract than the basic Handler type,
// including support for metadata, configuration, and lifecycle management.
//
// RequestHandler is typically used for final handlers in middleware chains,
// providing additional capabilities beyond simple request processing.
//
// Type Parameters:
//   - TRequest: The type of request data that this handler can process
//   - TResponse: The type of response data that this handler produces
//
// Design Philosophy:
//   RequestHandler represents the "destination" of a middleware chain - the component
//   that performs the actual business logic after all middleware has been applied.
//   It provides a rich interface for handlers that need more than simple processing.
type RequestHandler[TRequest, TResponse any] interface {
	// Handle processes the request and produces a response.
	// This is the primary method for request processing.
	//
	// Parameters:
	//   - ctx: Context for cancellation, timeouts, and values
	//   - request: The request data to process
	//
	// Returns:
	//   - TResponse: The processed response
	//   - error: Any error that occurred during processing
	Handle(ctx context.Context, request TRequest) (TResponse, error)

	// CanHandle determines if this handler can process the given request type.
	// This method enables dynamic handler selection and routing.
	//
	// Parameters:
	//   - request: The request to evaluate
	//
	// Returns:
	//   - bool: true if this handler can process the request
	CanHandle(request TRequest) bool

	// GetMetadata returns metadata about this handler.
	// This information can be used for routing, monitoring, and documentation.
	//
	// Returns:
	//   - HandlerMetadata: Metadata describing this handler
	GetMetadata() HandlerMetadata
}

// HandlerChain defines the interface for managing chains of handlers.
// This interface enables building complex handler hierarchies and routing logic.
//
// Type Parameters:
//   - TRequest: The type of request data flowing through the chain
//   - TResponse: The type of response data returned by the chain
//
// HandlerChain provides capabilities for:
// - Managing multiple handlers in sequence or parallel
// - Conditional handler execution based on request properties
// - Load balancing and failover between handlers
// - Aggregating responses from multiple handlers
type HandlerChain[TRequest, TResponse any] interface {
	// AddHandler adds a handler to the chain.
	// The handler will be considered for future requests.
	//
	// Parameters:
	//   - handler: The handler to add to the chain
	//   - priority: Priority for handler selection (higher numbers = higher priority)
	//
	// Returns:
	//   - error: Any error that occurred adding the handler
	AddHandler(handler RequestHandler[TRequest, TResponse], priority int) error

	// RemoveHandler removes a handler from the chain.
	//
	// Parameters:
	//   - handlerID: Unique identifier for the handler to remove
	//
	// Returns:
	//   - bool: true if the handler was found and removed
	RemoveHandler(handlerID string) bool

	// Execute processes a request through the handler chain.
	// The chain will select appropriate handlers and coordinate execution.
	//
	// Parameters:
	//   - ctx: Context for cancellation, timeouts, and values
	//   - request: The request to process
	//
	// Returns:
	//   - TResponse: The response from processing
	//   - error: Any error that occurred during processing
	Execute(ctx context.Context, request TRequest) (TResponse, error)

	// GetHandlers returns all handlers in the chain.
	//
	// Returns:
	//   - []RequestHandler[TRequest, TResponse]: All handlers in the chain
	GetHandlers() []RequestHandler[TRequest, TResponse]
}

// AsyncHandler defines the interface for handlers that support asynchronous processing.
// This interface enables non-blocking request processing with callback-based responses.
//
// Type Parameters:
//   - TRequest: The type of request data to process
//   - TResponse: The type of response data to return
//
// AsyncHandler is useful for:
// - Long-running operations that shouldn't block the caller
// - Integration with event-driven architectures
// - Batch processing and background tasks
// - Streaming responses and real-time updates
type AsyncHandler[TRequest, TResponse any] interface {
	// HandleAsync processes the request asynchronously.
	// The result will be delivered via the provided callback function.
	//
	// Parameters:
	//   - ctx: Context for cancellation, timeouts, and values
	//   - request: The request to process
	//   - callback: Function to call when processing completes
	//
	// Returns:
	//   - string: Unique identifier for tracking this async operation
	//   - error: Any immediate error (not processing errors, which go to callback)
	HandleAsync(ctx context.Context, request TRequest, callback AsyncCallback[TResponse]) (string, error)

	// GetStatus returns the status of an async operation.
	//
	// Parameters:
	//   - operationID: Unique identifier returned by HandleAsync
	//
	// Returns:
	//   - AsyncStatus: Current status of the operation
	//   - error: Any error getting the status
	GetStatus(operationID string) (AsyncStatus, error)

	// Cancel attempts to cancel an async operation.
	//
	// Parameters:
	//   - operationID: Unique identifier for the operation to cancel
	//
	// Returns:
	//   - bool: true if the operation was successfully cancelled
	Cancel(operationID string) bool
}

// AsyncCallback defines the function signature for async handler callbacks.
// This function will be called when async processing completes.
//
// Type Parameters:
//   - TResponse: The type of response data
//
// Parameters:
//   - response: The response from processing (may be nil if error occurred)
//   - err: Any error that occurred during processing
type AsyncCallback[TResponse any] func(response TResponse, err error)

// AsyncStatus represents the status of an asynchronous operation.
type AsyncStatus int

const (
	// AsyncStatusPending indicates the operation is still in progress
	AsyncStatusPending AsyncStatus = iota
	// AsyncStatusCompleted indicates the operation completed successfully
	AsyncStatusCompleted
	// AsyncStatusFailed indicates the operation failed with an error
	AsyncStatusFailed
	// AsyncStatusCancelled indicates the operation was cancelled
	AsyncStatusCancelled
	// AsyncStatusTimeout indicates the operation timed out
	AsyncStatusTimeout
)

// String returns a string representation of the AsyncStatus.
func (s AsyncStatus) String() string {
	switch s {
	case AsyncStatusPending:
		return "pending"
	case AsyncStatusCompleted:
		return "completed"
	case AsyncStatusFailed:
		return "failed"
	case AsyncStatusCancelled:
		return "cancelled"
	case AsyncStatusTimeout:
		return "timeout"
	default:
		return "unknown"
	}
}

// HandlerMetadata contains descriptive information about a handler.
// This metadata is used for documentation, monitoring, and routing decisions.
type HandlerMetadata struct {
	// ID is a unique identifier for this handler
	ID string

	// Name is a human-readable name for this handler
	Name string

	// Description provides detailed information about what this handler does
	Description string

	// Version indicates the version of this handler implementation
	Version string

	// Tags are arbitrary key-value pairs for categorizing and filtering handlers
	Tags map[string]string

	// SupportedMethods lists the types of requests this handler can process
	SupportedMethods []string

	// MaxConcurrency indicates the maximum number of concurrent requests this handler can process
	// A value of 0 indicates no limit
	MaxConcurrency int

	// AverageResponseTime provides performance information for monitoring
	AverageResponseTime time.Duration

	// SuccessRate indicates the percentage of requests that complete successfully
	SuccessRate float64

	// CreatedAt indicates when this handler was created or registered
	CreatedAt time.Time

	// LastUsed indicates when this handler was last used to process a request
	LastUsed time.Time
}

// HandlerFactory defines the interface for creating handlers dynamically.
// This interface enables dependency injection and configuration-based handler creation.
//
// Type Parameters:
//   - TRequest: The type of request data the created handlers will process
//   - TResponse: The type of response data the created handlers will return
//
// HandlerFactory is useful for:
// - Creating handlers with dependency injection
// - Dynamic handler configuration from external sources
// - Handler pooling and lifecycle management
// - Testing with mock handler creation
type HandlerFactory[TRequest, TResponse any] interface {
	// CreateHandler creates a new handler instance.
	//
	// Parameters:
	//   - config: Configuration data for the handler
	//
	// Returns:
	//   - RequestHandler[TRequest, TResponse]: The created handler
	//   - error: Any error that occurred during creation
	CreateHandler(config map[string]interface{}) (RequestHandler[TRequest, TResponse], error)

	// GetSupportedTypes returns the types of handlers this factory can create.
	//
	// Returns:
	//   - []string: List of handler type names
	GetSupportedTypes() []string

	// ValidateConfig validates configuration before creating a handler.
	//
	// Parameters:
	//   - handlerType: The type of handler to validate config for
	//   - config: The configuration to validate
	//
	// Returns:
	//   - error: Any validation error
	ValidateConfig(handlerType string, config map[string]interface{}) error
}
