// Package interfaces defines the core contracts and interfaces for the middleware system.
// This package contains all interface definitions that form the foundation of the
// middleware architecture, following Go's interface segregation principle.
//
// The interfaces defined here enable:
// - Type-safe middleware execution with full generic support
// - Flexible middleware composition and chaining
// - Context-aware request/response processing
// - Extensible middleware patterns (AOP, decorators, filters)
// - Integration with external systems (pipelines, HTTP, gRPC)
//
// All interfaces are designed to be minimal, focused, and easily testable,
// following Go's philosophy of small, composable interfaces that can be
// implemented independently and combined as needed.
package interfaces

import (
	"context"
)

// Middleware defines the core contract that all middleware implementations must satisfy.
// This is the fundamental interface that enables the middleware pattern in the GoVel framework.
//
// The middleware pattern allows for cross-cutting concerns like authentication, logging,
// validation, caching, rate limiting, and request/response transformation to be applied
// in a composable, reusable manner across different parts of an application.
//
// Type Parameters:
//   - TRequest: The type of the incoming request data that flows through the middleware
//   - TResponse: The type of the outgoing response data returned by the middleware
//
// Design Philosophy:
//   The interface follows the "Russian Doll" or "Onion" pattern where each middleware
//   layer wraps the next, allowing for:
//   - Pre-processing logic before calling the next handler
//   - Post-processing logic after the next handler returns
//   - Error handling and transformation
//   - Short-circuiting the chain for early returns
//   - Context modification and value passing
//
// Implementation Patterns:
//
// 1. **Struct-based Middleware**: Define a struct that implements Handle method
//    - Best for complex middleware with configuration, dependencies, and state
//    - Enables dependency injection and comprehensive testing
//    - Provides clear separation of concerns and encapsulation
//
// 2. **Function-based Middleware**: Use function types that implement this interface
//    - Best for simple, stateless middleware operations
//    - Enables functional programming patterns and inline definitions
//    - Convenient for rapid prototyping and simple transformations
//
// Usage Examples:
//
//   // Struct-based middleware implementation
//   type LoggingMiddleware struct {
//       logger Logger
//       level  LogLevel
//   }
//   
//   func (m *LoggingMiddleware) Handle(ctx context.Context, req *http.Request, next Handler[*http.Request, *http.Response]) (*http.Response, error) {
//       start := time.Now()
//       m.logger.Info("Processing request", "method", req.Method, "path", req.URL.Path)
//       
//       resp, err := next(ctx, req)
//       
//       duration := time.Since(start)
//       if err != nil {
//           m.logger.Error("Request failed", "error", err, "duration", duration)
//       } else {
//           m.logger.Info("Request completed", "status", resp.StatusCode, "duration", duration)
//       }
//       
//       return resp, err
//   }
//
//   // Function-based middleware using MiddlewareFunc type
//   var authMiddleware MiddlewareFunc[*http.Request, *http.Response] = func(ctx context.Context, req *http.Request, next Handler[*http.Request, *http.Response]) (*http.Response, error) {
//       token := req.Header.Get("Authorization")
//       if !validateToken(token) {
//           return &http.Response{StatusCode: 401}, errors.New("unauthorized")
//       }
//       return next(ctx, req)
//   }
//
// Error Handling Strategies:
//   Middleware can handle errors in several ways:
//   - **Pass Through**: Return errors unchanged for upstream handling
//   - **Transform**: Wrap errors with additional context or convert error types
//   - **Recover**: Handle specific errors and return success responses
//   - **Short Circuit**: Return early without calling next handler
//
// Context Management:
//   Middleware should properly handle context:
//   - **Cancellation**: Check ctx.Done() for long-running operations
//   - **Timeouts**: Respect context deadlines and timeouts
//   - **Values**: Add context values for downstream middleware and handlers
//   - **Tracing**: Propagate tracing and correlation IDs
//
// Thread Safety:
//   Middleware implementations must be safe for concurrent use if they will be
//   called from multiple goroutines. This includes proper handling of shared
//   resources, configuration, and any mutable state.
type Middleware[TRequest, TResponse any] interface {
	// Handle processes the request through the middleware and calls the next handler in the chain.
	// This is the core method that implements the middleware logic and defines the middleware's behavior.
	//
	// The middleware implementation should follow this general pattern:
	// 1. **Pre-processing**: Perform any logic needed before the request continues (validation, logging, etc.)
	// 2. **Decision Point**: Decide whether to call the next handler or short-circuit
	// 3. **Delegation**: Call the next handler with the (possibly modified) request
	// 4. **Post-processing**: Perform any logic needed after the response is received (logging, transformation, etc.)
	// 5. **Return**: Return the final response and any errors
	//
	// Parameters:
	//   - ctx: The context for the request, containing cancellation signals, timeouts, and values
	//     * Should be checked for cancellation in long-running operations
	//     * Can be augmented with additional values for downstream handlers
	//     * Must be passed to the next handler (possibly modified)
	//   - request: The incoming request data of type TRequest to be processed
	//     * Can be inspected, validated, and modified before passing to next handler
	//     * Modifications should be made thoughtfully to maintain data integrity
	//   - next: The next handler in the middleware chain (from types/handlers package)
	//     * Represents "the rest of the chain" including all remaining middleware and final handler
	//     * Can be called zero, one, or multiple times depending on middleware logic
	//     * Should be called with proper context and request data
	//
	// Returns:
	//   - TResponse: The processed response data after middleware and handler execution
	//     * Can be the response from next handler, a transformed version, or completely new response
	//     * Should maintain expected response contract for the application
	//   - error: Any error that occurred during middleware processing or from downstream handlers
	//     * Can be errors from next handler, new errors, or transformed errors
	//     * Should provide sufficient context for debugging and error handling
	//
	// Middleware Responsibilities:
	//
	// **Request Processing**:
	//   - Validate request data and structure
	//   - Authenticate and authorize requests
	//   - Transform or enrich request data
	//   - Apply rate limiting and throttling
	//   - Log request details for auditing
	//
	// **Response Processing**:
	//   - Transform response data and format
	//   - Add security headers and metadata
	//   - Cache responses for future requests
	//   - Log response details and metrics
	//   - Handle error responses appropriately
	//
	// **Error Handling**:
	//   - Catch and handle specific error types
	//   - Transform errors for consistent API responses
	//   - Log errors with appropriate detail levels
	//   - Implement circuit breaker patterns
	//   - Provide fallback responses when appropriate
	//
	// **Context Management**:
	//   - Respect context cancellation and timeouts
	//   - Add tracing and correlation information
	//   - Pass authentication and authorization data
	//   - Provide request-scoped configuration
	//
	// **Performance Considerations**:
	//   - Minimize allocations in hot paths
	//   - Use efficient algorithms for request processing
	//   - Implement proper caching strategies
	//   - Profile and monitor middleware performance
	//   - Avoid blocking operations without timeout
	//
	// **Security Considerations**:
	//   - Validate all input data thoroughly
	//   - Sanitize data to prevent injection attacks
	//   - Implement proper authentication checks
	//   - Handle sensitive data securely
	//   - Log security events appropriately
	//
	// Thread Safety Requirements:
	//   Implementations must be safe for concurrent use across multiple goroutines.
	//   This includes proper synchronization of any shared state and thread-safe
	//   access to configuration and dependencies.
Handle(ctx context.Context, request TRequest, next Handler[TRequest, TResponse]) (TResponse, error)
}

// Handler represents the signature of functions that can process requests and return responses.
// This type is imported from the handlers package and represents the "next" function
// that middleware receives to continue the processing chain.
//
// The Handler type enables the middleware pattern by providing a consistent interface
// for calling the next step in the processing chain, whether that's another middleware
// or the final handler that produces the actual response.
//
// Type Parameters:
//   - TRequest: The type of request data that the handler can process
//   - TResponse: The type of response data that the handler returns
//
// Handler implementations should:
// - Process the request according to their specific logic
// - Respect context cancellation and timeouts
// - Return appropriate responses and meaningful errors
// - Be safe for concurrent use if needed
//
// This type alias provides a convenient reference to the handler signature
// used throughout the middleware system, ensuring consistency and type safety.
type Handler[TRequest, TResponse any] interface {
	// Execute processes the request and returns a response.
	// This method represents the actual processing logic that handles the request.
	//
	// Parameters:
	//   - ctx: Context for cancellation, timeouts, and values
	//   - request: The request data to process
	//
	// Returns:
	//   - TResponse: The processed response
	//   - error: Any error that occurred during processing
	Execute(ctx context.Context, request TRequest) (TResponse, error)
}
