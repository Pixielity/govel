// Package factories provides factory functions and utilities for creating middleware and handler components.
// This package implements factory patterns for middleware creation, composition, and chain management,
// providing a functional programming approach to HTTP middleware construction.
//
// The middleware factory offers several key capabilities:
//   - Function-to-interface conversion for handlers and middleware
//   - Middleware composition and chaining utilities
//   - Type-safe middleware creation with proper interface compliance
//   - Functional programming patterns for middleware construction
//   - Builder patterns for complex middleware chains
//
// Design Patterns Implemented:
//   - Factory Pattern: Creating middleware and handler instances
//   - Builder Pattern: Constructing middleware chains
//   - Decorator Pattern: Composing middleware with additional behavior
//   - Functional Programming: Function-based middleware creation
//
// Usage Patterns:
//
//	The factory supports both functional and object-oriented middleware styles:
//
//	Functional Style:
//	  middleware := factories.Middleware(func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
//	      // middleware logic
//	      return next.Handle(req)
//	  })
//
//	Chain Building:
//	  chain := factories.Chain(
//	      middleware.Logger(),
//	      middleware.Auth(),
//	      middleware.CORS(),
//	  )
//
//	Composition:
//	  composed := factories.Compose(authMiddleware, loggingMiddleware)
package factories

import (
	webserver "govel/packages/new/webserver/src"
	"govel/packages/new/webserver/src/interfaces"
	"govel/packages/new/webserver/src/types"
)

// MiddlewareFactory provides methods for creating and managing middleware instances.
// This factory encapsulates middleware creation patterns and provides utilities
// for building complex middleware chains and compositions.
//
// The factory supports multiple middleware creation patterns:
//   - Function-based middleware creation
//   - Middleware chain construction
//   - Middleware composition and decoration
//   - Handler function wrapping
//
// Thread Safety:
//
//	All factory methods are thread-safe and can be called concurrently.
//	The created middleware instances should be designed to handle concurrent requests.
type MiddlewareFactory struct{}

// NewMiddlewareFactory creates a new middleware factory instance.
//
// Returns:
//   - *MiddlewareFactory: A new factory instance for creating middleware
//
// Example:
//
//	factory := NewMiddlewareFactory()
//	middleware := factory.CreateFromFunction(myMiddlewareFunction)
func NewMiddlewareFactory() *MiddlewareFactory {
	return &MiddlewareFactory{}
}

// CreateFromFunction wraps a middleware function into a MiddlewareInterface.
// This method converts function-based middleware into interface-compliant objects
// that can be used throughout the webserver middleware system.
//
// The function signature must match:
//
//	func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface
//
// Parameters:
//   - fn: The middleware function to wrap
//
// Returns:
//   - interfaces.MiddlewareInterface: A middleware instance that implements the full interface
//
// Usage Examples:
//
//	// Basic middleware function
//	loggerFunc := func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
//	    start := time.Now()
//	    response := next.Handle(req)
//	    log.Printf("Request processed in %v", time.Since(start))
//	    return response
//	}
//
//	// Convert to middleware interface
//	factory := NewMiddlewareFactory()
//	middleware := factory.CreateFromFunction(loggerFunc)
//
//	// Use in webserver
//	server.Use(middleware)
func (f *MiddlewareFactory) CreateFromFunction(fn func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface) interfaces.MiddlewareInterface {
	return types.MiddlewareFunc(fn)
}

// CreateChain builds a MiddlewareChain from provided middleware instances.
// This method creates an execution chain where middleware are executed in order,
// with each middleware having the opportunity to process the request before and after
// the next middleware in the chain.
//
// Parameters:
//   - middleware: Variable number of middleware instances to chain together
//
// Returns:
//   - *types.MiddlewareChain: A middleware chain that can execute all middleware in sequence
//
// Chain Execution Order:
//
//	Middleware are executed in the order provided:
//	1. First middleware's Before() method
//	2. First middleware's Handle() method (which calls next)
//	3. Second middleware's Handle() method, etc.
//	4. Final handler execution
//	5. Middleware After() methods in reverse order
//
// Usage Examples:
//
//	// Create chain with multiple middleware
//	factory := NewMiddlewareFactory()
//	chain := factory.CreateChain(
//	    middleware.Recovery(),  // Executes first (outermost)
//	    middleware.Logger(),    // Executes second
//	    middleware.Auth(),      // Executes third
//	    middleware.CORS(),      // Executes fourth (innermost)
//	)
//
//	// Execute chain with final handler
//	response := chain.Execute(request, finalHandler)
func (f *MiddlewareFactory) CreateChain(middleware ...interfaces.MiddlewareInterface) *webserver.MiddlewareChain {
	return webserver.NewMiddlewareChain(middleware...)
}

// ComposeMiddleware composes two middleware into a single middleware.
// This method creates a new middleware that executes the first middleware,
// which then executes the second middleware, creating a nested execution pattern.
//
// The composition follows the mathematical function composition pattern:
//
//	compose(f, g)(x) = f(g(x))
//	In middleware terms: a.Handle(req, b.Handle)
//
// Parameters:
//   - outer: The outer middleware that executes first
//   - inner: The inner middleware that executes second
//
// Returns:
//   - interfaces.MiddlewareInterface: A composed middleware that combines both behaviors
//
// Execution Flow:
//  1. outer.Handle() is called with the request
//  2. outer decides whether to call next (which is inner.Handle)
//  3. inner.Handle() is called if outer proceeds
//  4. inner decides whether to call the actual next handler
//  5. Response flows back through inner, then outer
//
// Usage Examples:
//
//	// Compose authentication and logging
//	factory := NewMiddlewareFactory()
//	authThenLog := factory.ComposeMiddleware(
//	    middleware.Auth(),    // Executes first (outer)
//	    middleware.Logger(),  // Executes second (inner)
//	)
//
//	// Use composed middleware
//	server.Use(authThenLog)
//
//	// Chain multiple compositions
//	final := factory.ComposeMiddleware(
//	    middleware.Recovery(),
//	    factory.ComposeMiddleware(
//	        middleware.Auth(),
//	        middleware.Logger(),
//	    ),
//	)
func (f *MiddlewareFactory) ComposeMiddleware(outer, inner interfaces.MiddlewareInterface) interfaces.MiddlewareInterface {
	return types.MiddlewareFunc(func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
		// Create a handler that wraps the inner middleware
		innerHandler := types.HandlerFunc(func(r interfaces.RequestInterface) interfaces.ResponseInterface {
			return inner.Handle(r, next)
		})

		// Execute outer middleware with inner as the "next" handler
		return outer.Handle(req, innerHandler)
	})
}

// Package-level convenience functions that provide a simplified API
// for common middleware creation patterns.

// Middleware wraps a function into a MiddlewareInterface using the default factory.
// This is a convenience function for quick middleware creation without explicitly
// creating a factory instance.
//
// Function Signature:
//
//	The provided function must match:
//	func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface
//
// Parameters:
//   - fn: The middleware function to wrap into an interface
//
// Returns:
//   - interfaces.MiddlewareInterface: A middleware interface implementation
//
// Usage Examples:
//
//	// Simple request timing middleware
//	timingMiddleware := factories.Middleware(func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
//	    start := time.Now()
//	    resp := next.Handle(req)
//	    duration := time.Since(start)
//	    resp.Header("X-Response-Time", duration.String())
//	    return resp
//	})
//
//	// Authentication middleware
//	authMiddleware := factories.Middleware(func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
//	    token := req.Header("Authorization")
//	    if !isValidToken(token) {
//	        return factories.NewResponse().Status(401).Json(map[string]string{"error": "unauthorized"})
//	    }
//	    return next.Handle(req)
//	})
//
//	// Request modification middleware
//	enrichMiddleware := factories.Middleware(func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
//	    req.SetHeader("X-Request-ID", generateRequestID())
//	    req.SetHeader("X-Timestamp", time.Now().Format(time.RFC3339))
//	    return next.Handle(req)
//	})
func Middleware(fn func(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface) interfaces.MiddlewareInterface {
	factory := NewMiddlewareFactory()
	return factory.CreateFromFunction(fn)
}

// Handler wraps a function into a HandlerInterface.
// This convenience function converts handler functions into interface-compliant objects
// that can be used as route handlers or final handlers in middleware chains.
//
// Function Signature:
//
//	The provided function must match:
//	func(req interfaces.RequestInterface) interfaces.ResponseInterface
//
// Parameters:
//   - fn: The handler function to wrap into an interface
//
// Returns:
//   - interfaces.HandlerInterface: A handler interface implementation
//
// Usage Examples:
//
//	// Simple JSON API handler
//	userHandler := factories.Handler(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
//	    userID := req.Param("id")
//	    user, err := getUserByID(userID)
//	    if err != nil {
//	        return factories.NewResponse().Status(404).Json(map[string]string{"error": "user not found"})
//	    }
//	    return factories.NewResponse().Json(user)
//	})
//
//	// File upload handler
//	uploadHandler := factories.Handler(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
//	    file, err := req.File("upload")
//	    if err != nil {
//	        return factories.NewResponse().Status(400).Json(map[string]string{"error": "no file provided"})
//	    }
//
//	    if err := saveFile(file); err != nil {
//	        return factories.NewResponse().Status(500).Json(map[string]string{"error": "failed to save file"})
//	    }
//
//	    return factories.NewResponse().Status(201).Json(map[string]string{"message": "file uploaded successfully"})
//	})
//
//	// Health check handler
//	healthHandler := factories.Handler(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
//	    return factories.NewResponse().Json(map[string]interface{}{
//	        "status": "healthy",
//	        "timestamp": time.Now(),
//	        "version": "1.0.0",
//	    })
//	})
func Handler(fn func(req interfaces.RequestInterface) interfaces.ResponseInterface) interfaces.HandlerInterface {
	return types.HandlerFunc(fn)
}

// Chain builds a MiddlewareChain from provided middleware using the default factory.
// This convenience function creates a middleware chain without explicitly creating
// a factory instance.
//
// Parameters:
//   - m: Variable number of middleware instances to chain together
//
// Returns:
//   - *types.MiddlewareChain: An executable middleware chain
//
// Chain Characteristics:
//   - Middleware execute in the order provided
//   - Each middleware can short-circuit the chain
//   - Response flows back through middleware in reverse order
//   - Thread-safe for concurrent execution
//
// Usage Examples:
//
//	// Standard web application middleware chain
//	standardChain := factories.Chain(
//	    middleware.Recovery(),     // Catch panics (outermost)
//	    middleware.Logger(),       // Log requests and responses
//	    middleware.CORS(),         // Handle cross-origin requests
//	    middleware.RateLimit(),    // Rate limiting
//	)
//
//	// API-specific middleware chain
//	apiChain := factories.Chain(
//	    middleware.Recovery(),
//	    middleware.Logger(),
//	    middleware.Auth(),         // Authentication required
//	    middleware.Validate(),     // Request validation
//	)
//
//	// Execute chain with final handler
//	response := standardChain.Execute(request, finalHandler)
func Chain(m ...interfaces.MiddlewareInterface) *webserver.MiddlewareChain {
	factory := NewMiddlewareFactory()
	return factory.CreateChain(m...)
}

// Compose composes two middleware into one using the default factory.
// This convenience function creates composed middleware without explicitly
// creating a factory instance.
//
// Parameters:
//   - a: The outer middleware (executes first)
//   - b: The inner middleware (executes second)
//
// Returns:
//   - interfaces.MiddlewareInterface: A composed middleware combining both behaviors
//
// Composition Benefits:
//   - Creates reusable middleware combinations
//   - Reduces boilerplate in middleware registration
//   - Enables functional middleware construction patterns
//   - Maintains clean separation of concerns
//
// Usage Examples:
//
//	// Create a composed authentication + authorization middleware
//	authzMiddleware := factories.Compose(
//	    middleware.Auth(),     // Outer: authenticate first
//	    middleware.Authz(),    // Inner: then authorize
//	)
//
//	// Compose logging with request ID generation
//	trackedMiddleware := factories.Compose(
//	    middleware.RequestID(), // Outer: generate request ID first
//	    middleware.Logger(),    // Inner: then log with ID
//	)
//
//	// Chain composed middleware
//	server.Use(
//	    middleware.Recovery(),
//	    authzMiddleware,        // Uses composed auth + authz
//	    trackedMiddleware,      // Uses composed request ID + logging
//	)
func Compose(a, b interfaces.MiddlewareInterface) interfaces.MiddlewareInterface {
	factory := NewMiddlewareFactory()
	return factory.ComposeMiddleware(a, b)
}
