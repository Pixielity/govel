// Package interfaces - Adapter interface definition
// This file defines the AdapterInterface contract for framework adapters (GoFiber, Gin, Echo).
package interfaces

import "context"

// AdapterInterface defines how a framework-specific adapter integrates with the unified webserver API.
// Implementations translate the high-level WebserverInterface calls into framework-specific operations.
//
// Responsibilities:
//   - Register routes and handlers for HTTP methods
//   - Apply global and route-specific middleware
//   - Manage server lifecycle (start, stop, graceful shutdown)
//   - Provide access to request/response translations
//
// Note: Route grouping can be implemented internally by adapters or coordinated by the builder/server.
type AdapterInterface interface {
	// Initialization & Configuration
	
	// Init initializes the adapter with initial configuration and global middleware.
	//
	// Parameters:
	//   config: Arbitrary configuration map
	//   middleware: Global middleware to apply
	Init(config map[string]interface{}, middleware []MiddlewareInterface) error
	
	// SetConfig sets a configuration value for the adapter at runtime.
	SetConfig(key string, value interface{})
	
	// GetConfig retrieves a configuration value for the adapter.
	GetConfig(key string) interface{}
	
	// Middleware
	
	// Use registers global middleware for all routes.
	Use(middleware ...MiddlewareInterface)
	
	// Routing
	
	// Handle registers a route for a given HTTP method and path with the provided handler.
	// Middlewares provided here are route-specific and applied after global middleware.
	Handle(method, path string, handler HandlerInterface, middlewares ...MiddlewareInterface)
	
	// Group creates a route group with a prefix and returns a function to register routes within the group.
	// The returned function should be used to define routes that share the prefix and optional middleware.
	Group(prefix string, register func(), middlewares ...MiddlewareInterface)
	
	// Lifecycle
	
	// Listen starts the HTTP server on the specified address.
	Listen(addr string) error
	
	// ListenTLS starts the HTTPS server with TLS certificates.
	ListenTLS(addr, certFile, keyFile string) error
	
	// Shutdown gracefully shuts down the server.
	Shutdown(ctx context.Context) error
}
