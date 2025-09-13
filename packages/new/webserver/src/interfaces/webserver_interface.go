// Package interfaces provides all the contracts and interfaces used throughout the webserver package.
// These interfaces define the behavior and contracts that implementations must follow,
// ensuring consistency across different web framework adapters (GoFiber, Gin, Echo).
package interfaces

import (
	"context"
)

// WebserverInterface defines the main webserver contract that all webserver implementations must follow.
// This interface provides a Laravel-inspired API for web server functionality, abstracting away
// the underlying web framework (GoFiber, Gin, Echo) and providing a unified interface.
//
// The interface supports:
//   - HTTP route registration (GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD)
//   - Route grouping with prefixes
//   - Middleware registration and chaining
//   - Server lifecycle management (start, stop, graceful shutdown)
//   - Configuration management
//   - Static file serving
//   - Host and port configuration
//
// Example usage:
//   server := webserver.New()
//   server.Get("/users", userHandler).
//          Post("/users", createUserHandler).
//          Listen(":8080")
type WebserverInterface interface {
	// HTTP Route Registration Methods
	// These methods register handlers for specific HTTP methods and paths
	
	// Get registers a handler for HTTP GET requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern (e.g., "/users", "/users/:id")
	//   handler: The handler function to execute for matching requests
	//
	// Example:
	//   server.Get("/users", func(req *Request) *Response {
	//       return Json(map[string]string{"message": "Hello"})
	//   })
	Get(path string, handler HandlerInterface) WebserverInterface
	
	// Post registers a handler for HTTP POST requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Post(path string, handler HandlerInterface) WebserverInterface
	
	// Put registers a handler for HTTP PUT requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Put(path string, handler HandlerInterface) WebserverInterface
	
	// Patch registers a handler for HTTP PATCH requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Patch(path string, handler HandlerInterface) WebserverInterface
	
	// Delete registers a handler for HTTP DELETE requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Delete(path string, handler HandlerInterface) WebserverInterface
	
	// Options registers a handler for HTTP OPTIONS requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Options(path string, handler HandlerInterface) WebserverInterface
	
	// Head registers a handler for HTTP HEAD requests to the specified path.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   path: The URL path pattern
	//   handler: The handler function to execute for matching requests
	Head(path string, handler HandlerInterface) WebserverInterface
	
	// Route Grouping
	
	// Group creates a route group with the specified prefix.
	// All routes registered within the handler function will be prefixed with the given prefix.
	// This is useful for organizing related routes (e.g., API versioning, admin routes).
	//
	// Parameters:
	//   prefix: The URL prefix for all routes in this group (e.g., "/api/v1", "/admin")
	//   handler: A function that receives the group instance for route registration
	//
	// Example:
	//   server.Group("/api/v1", func(api WebserverInterface) {
	//       api.Get("/users", getUsersHandler)     // Maps to /api/v1/users
	//       api.Post("/users", createUserHandler)  // Maps to /api/v1/users
	//   })
	Group(prefix string, handler func(WebserverInterface)) WebserverInterface
	
	// Middleware Registration
	
	// Use registers middleware to be applied to all routes.
	// Middleware is executed in the order it was registered.
	// Returns the webserver instance for method chaining.
	//
	// Parameters:
	//   middleware: One or more middleware functions to register
	//
	// Example:
	//   server.Use(CorsMiddleware(), AuthMiddleware())
	Use(middleware ...MiddlewareInterface) WebserverInterface
	
	// Middleware is an alias for Use() for better readability.
	// Registers middleware to be applied to all routes.
	//
	// Parameters:
	//   middleware: One or more middleware functions to register
	Middleware(middleware ...MiddlewareInterface) WebserverInterface
	
	// Server Lifecycle Management
	
	// Listen starts the HTTP server on the specified address.
	// If no address is provided, defaults to ":8080".
	// This method blocks until the server is stopped or encounters an error.
	//
	// Parameters:
	//   addr: Optional server address in the format "host:port" or ":port"
	//
	// Returns:
	//   error: Any error that occurred while starting or running the server
	//
	// Example:
	//   err := server.Listen(":3000")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	Listen(addr ...string) error
	
	// ListenTLS starts the HTTPS server with the provided certificate files.
	// This method blocks until the server is stopped or encounters an error.
	//
	// Parameters:
	//   certFile: Path to the SSL certificate file
	//   keyFile: Path to the SSL private key file
	//   addr: Optional server address in the format "host:port" or ":port"
	//
	// Returns:
	//   error: Any error that occurred while starting or running the server
	//
	// Example:
	//   err := server.ListenTLS("/path/to/cert.pem", "/path/to/key.pem", ":443")
	ListenTLS(certFile, keyFile string, addr ...string) error
	
	// Shutdown gracefully shuts down the server with the provided context.
	// The server will stop accepting new requests and will attempt to finish
	// processing existing requests within the context timeout.
	//
	// Parameters:
	//   ctx: Context with timeout for graceful shutdown
	//
	// Returns:
	//   error: Any error that occurred during shutdown
	//
	// Example:
	//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//   defer cancel()
	//   err := server.Shutdown(ctx)
	Shutdown(ctx context.Context) error
	
	// Configuration Management
	
	// SetConfig sets a configuration value for the webserver.
	// This allows for runtime configuration of server behavior.
	//
	// Parameters:
	//   key: The configuration key (e.g., "timeout", "max_body_size")
	//   value: The configuration value
	//
	// Returns:
	//   WebserverInterface: The webserver instance for method chaining
	SetConfig(key string, value interface{}) WebserverInterface
	
	// GetConfig retrieves a configuration value by key.
	// Returns nil if the key doesn't exist.
	//
	// Parameters:
	//   key: The configuration key to retrieve
	//
	// Returns:
	//   interface{}: The configuration value, or nil if not found
	GetConfig(key string) interface{}
	
	// Static File Serving
	
	// Static registers a static file server for the given prefix and root directory.
	// This is useful for serving CSS, JavaScript, images, and other static assets.
	//
	// Parameters:
	//   prefix: The URL prefix for static files (e.g., "/static", "/assets")
	//   root: The file system path to the directory containing static files
	//
	// Returns:
	//   WebserverInterface: The webserver instance for method chaining
	//
	// Example:
	//   server.Static("/assets", "./public/assets")
	//   // Now files in ./public/assets/ are served at /assets/
	Static(prefix, root string) WebserverInterface
	
	// Host and Port Configuration
	
	// Port sets the default port for the server.
	// This port will be used if no port is specified in Listen() or ListenTLS().
	//
	// Parameters:
	//   port: The port number (e.g., 8080, 3000)
	//
	// Returns:
	//   WebserverInterface: The webserver instance for method chaining
	//
	// Example:
	//   server.Port(3000).Listen() // Will listen on :3000
	Port(port int) WebserverInterface
	
	// Host sets the default host for the server.
	// This host will be used if no host is specified in Listen() or ListenTLS().
	//
	// Parameters:
	//   host: The host address (e.g., "localhost", "0.0.0.0")
	//
	// Returns:
	//   WebserverInterface: The webserver instance for method chaining
	//
	// Example:
	//   server.Host("0.0.0.0").Port(8080).Listen() // Will listen on 0.0.0.0:8080
	Host(host string) WebserverInterface
}
