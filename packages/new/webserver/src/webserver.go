// Package webserver provides the main webserver implementation that unifies multiple web frameworks
// under a single Laravel-inspired API. This is the primary package users interact with.
package webserver

import (
	"context"
	"fmt"
	"strings"

	routing "govel/new/routing/src"
	"govel/new/webserver/src/enums"
	"govel/new/webserver/src/interfaces"

	"govel/new/webserver/src/adapters"
)

// Webserver is the main webserver struct that provides a unified interface across multiple web frameworks.
// It wraps framework-specific adapters (GoFiber, Gin, Echo) and provides a Laravel-inspired API
// for route registration, middleware handling, and server lifecycle management.
//
// The webserver acts as a facade over the underlying adapter, translating high-level API calls
// into framework-specific operations while maintaining a consistent interface.
//
// Example usage:
//
//	server := webserver.New()
//	server.Get("/users", getUsersHandler)
//	server.Post("/users", createUserHandler)
//	server.Listen(":8080")
type Webserver struct {
	// adapter is the underlying framework-specific adapter (GoFiber, Gin, Echo)
	adapter interfaces.AdapterInterface

	// config stores webserver configuration values
	config map[string]interface{}

	// routes stores registered routes for introspection and URL generation
	routes *routing.RouteCollection

	// middleware stores global middleware that applies to all routes
	middleware []interfaces.MiddlewareInterface

	// running indicates whether the server is currently running
	running bool

	// address stores the current listening address
	address string
}

// NewWithEngine creates a new webserver instance with the specified engine.
// This is the main factory function that creates adapters and initializes the webserver.
//
// Parameters:
//
//	engine: The engine name ("gin", "echo", "fiber", "net/http")
//	config: Optional configuration parameters
//
// Returns:
//
//	interfaces.WebserverInterface: A new webserver instance
//	error: Error if engine is unsupported or initialization fails
//
// Example:
//
//	server, err := webserver.NewWithEngine("gin")
//	if err != nil {
//	    log.Fatal("Failed to create webserver:", err)
//	}
func NewWithEngine(engine string, config ...map[string]interface{}) (interfaces.WebserverInterface, error) {
	cfg := map[string]interface{}{}
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	adapter, err := adapters.Create(engine)
	if err != nil {
		return nil, err
	}

	ws := &Webserver{
		adapter:    adapter,
		config:     cfg,
		middleware: []interfaces.MiddlewareInterface{},
		routes:     routing.NewRouteCollection(),
		running:    false,
	}

	// Initialize the adapter with config and global middleware
	if err := adapter.Init(cfg, ws.middleware); err != nil {
		return nil, fmt.Errorf("failed to initialize adapter: %v", err)
	}

	return ws, nil
}

// New creates a new webserver instance with default configuration and no adapter.
// This is primarily used for testing or when the adapter will be set later.
//
// Returns:
//
//	*Webserver: A new webserver instance
//
// Example:
//
//	server := webserver.New()
//	server.SetAdapter(myAdapter)
func New() *Webserver {
	return &Webserver{
		config:     make(map[string]interface{}),
		routes:     routing.NewRouteCollection(),
		middleware: make([]interfaces.MiddlewareInterface, 0),
		running:    false,
	}
}

// NewWithAdapter creates a new webserver instance with a specific adapter.
// This allows injecting custom adapters or using non-default engines.
//
// Parameters:
//
//	adapter: The adapter to use for the webserver
//
// Returns:
//
//	*Webserver: A new webserver instance using the specified adapter
//
// Example:
//
//	adapter := gofiber.NewGoFiberAdapter()
//	server := webserver.NewWithAdapter(adapter)
func NewWithAdapter(adapter interfaces.AdapterInterface) *Webserver {
	server := New()
	server.adapter = adapter
	return server
}

// HTTP Route Registration Methods
// These methods provide a Laravel-inspired API for registering HTTP routes

// Get registers a GET route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern (e.g., "/users", "/users/:id")
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
//
// Example:
//
//	server.Get("/users/:id", func(req interfaces.RequestInterface) interfaces.ResponseInterface {
//	    id := req.Param("id")
//	    return NewResponse().Json(map[string]string{"user_id": id})
//	})
func (w *Webserver) Get(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.GET, path, handler)
}

// Post registers a POST route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Post(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.POST, path, handler)
}

// Put registers a PUT route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Put(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.PUT, path, handler)
}

// Patch registers a PATCH route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Patch(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.PATCH, path, handler)
}

// Delete registers a DELETE route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Delete(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.DELETE, path, handler)
}

// Options registers an OPTIONS route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Options(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.OPTIONS, path, handler)
}

// Head registers a HEAD route handler.
// Returns the webserver instance for method chaining.
//
// Parameters:
//
//	path: The URL path pattern
//	handler: The handler function to execute for matching requests
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Head(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return w.registerRoute(enums.HEAD, path, handler)
}

// Route Grouping

// Group creates a route group with the specified prefix.
// All routes registered within the handler function will be prefixed with the given prefix.
//
// Parameters:
//
//	prefix: The URL prefix for all routes in this group
//	handler: A function that receives the group instance for route registration
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
//
// Example:
//
//	server.Group("/api/v1", func(api interfaces.WebserverInterface) {
//	    api.Get("/users", getUsersHandler)     // Maps to /api/v1/users
//	    api.Post("/users", createUserHandler)  // Maps to /api/v1/users
//	})
func (w *Webserver) Group(prefix string, handler func(interfaces.WebserverInterface)) interfaces.WebserverInterface {
	// Create a group wrapper that prefixes all routes
	groupWrapper := &routeGroup{
		webserver: w,
		prefix:    strings.TrimSuffix(prefix, "/"),
	}

	// Execute the handler with the group wrapper
	handler(groupWrapper)

	return w
}

// Middleware Registration

// Use registers middleware to be applied to all routes.
// Middleware is executed in the order it was registered.
//
// Parameters:
//
//	middleware: One or more middleware functions to register
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
//
// Example:
//
//	server.Use(corsMiddleware, authMiddleware, loggerMiddleware)
func (w *Webserver) Use(middleware ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	w.middleware = append(w.middleware, middleware...)

	// Also register with the adapter if it exists
	if w.adapter != nil {
		w.adapter.Use(middleware...)
	}

	return w
}

// Middleware is an alias for Use() for better readability.
// Registers middleware to be applied to all routes.
//
// Parameters:
//
//	middleware: One or more middleware functions to register
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Middleware(middleware ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	return w.Use(middleware...)
}

// Server Lifecycle Management

// Listen starts the HTTP server on the specified address.
// If no address is provided, defaults to ":8080".
//
// Parameters:
//
//	addr: Optional server address in the format "host:port" or ":port"
//
// Returns:
//
//	error: Any error that occurred while starting or running the server
//
// Example:
//
//	err := server.Listen(":3000")
//	if err != nil {
//	    log.Fatal("Server failed to start:", err)
//	}
func (w *Webserver) Listen(addr ...string) error {
	if w.adapter == nil {
		return fmt.Errorf("no adapter configured - use webserver builder to create webserver")
	}

	// Determine listening address
	listenAddr := ":8080"
	if len(addr) > 0 && addr[0] != "" {
		listenAddr = addr[0]
	}

	w.address = listenAddr
	w.running = true

	// Start the underlying adapter
	return w.adapter.Listen(listenAddr)
}

// ListenTLS starts the HTTPS server with the provided certificate files.
//
// Parameters:
//
//	certFile: Path to the SSL certificate file
//	keyFile: Path to the SSL private key file
//	addr: Optional server address in the format "host:port" or ":port"
//
// Returns:
//
//	error: Any error that occurred while starting or running the server
//
// Example:
//
//	err := server.ListenTLS("/path/to/cert.pem", "/path/to/key.pem", ":443")
func (w *Webserver) ListenTLS(certFile, keyFile string, addr ...string) error {
	if w.adapter == nil {
		return fmt.Errorf("no adapter configured - use webserver builder to create webserver")
	}

	// Determine listening address
	listenAddr := ":443"
	if len(addr) > 0 && addr[0] != "" {
		listenAddr = addr[0]
	}

	w.address = listenAddr
	w.running = true

	// Start the underlying adapter with TLS
	return w.adapter.ListenTLS(listenAddr, certFile, keyFile)
}

// Shutdown gracefully shuts down the server with the provided context.
//
// Parameters:
//
//	ctx: Context with timeout for graceful shutdown
//
// Returns:
//
//	error: Any error that occurred during shutdown
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	err := server.Shutdown(ctx)
func (w *Webserver) Shutdown(ctx context.Context) error {
	if w.adapter == nil {
		return fmt.Errorf("no adapter configured")
	}

	w.running = false
	return w.adapter.Shutdown(ctx)
}

// Configuration Management

// SetConfig sets a configuration value for the webserver.
//
// Parameters:
//
//	key: The configuration key
//	value: The configuration value
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) SetConfig(key string, value interface{}) interfaces.WebserverInterface {
	w.config[key] = value

	// Also set on adapter if it exists
	if w.adapter != nil {
		w.adapter.SetConfig(key, value)
	}

	return w
}

// GetConfig retrieves a configuration value by key.
//
// Parameters:
//
//	key: The configuration key to retrieve
//
// Returns:
//
//	interface{}: The configuration value, or nil if not found
func (w *Webserver) GetConfig(key string) interface{} {
	return w.config[key]
}

// Static File Serving

// Static registers a static file server for the given prefix and root directory.
//
// Parameters:
//
//	prefix: The URL prefix for static files
//	root: The file system path to the directory containing static files
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
//
// Example:
//
//	server.Static("/assets", "./public/assets")
func (w *Webserver) Static(prefix, root string) interfaces.WebserverInterface {
	// TODO: Implement static file serving
	// This would typically involve registering a special handler with the adapter
	// that serves files from the specified directory
	return w
}

// Host and Port Configuration

// Port sets the default port for the server.
//
// Parameters:
//
//	port: The port number
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Port(port int) interfaces.WebserverInterface {
	return w.SetConfig("port", port)
}

// Host sets the default host for the server.
//
// Parameters:
//
//	host: The host address
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) Host(host string) interfaces.WebserverInterface {
	return w.SetConfig("host", host)
}

// Utility Methods

// GetRoutes returns all registered routes.
//
// Returns:
//
//	*routing.RouteCollection: The collection of registered routes
func (w *Webserver) GetRoutes() *routing.RouteCollection {
	return w.routes
}

// IsRunning returns whether the server is currently running.
//
// Returns:
//
//	bool: True if the server is running, false otherwise
func (w *Webserver) IsRunning() bool {
	return w.running
}

// GetAddress returns the current listening address.
//
// Returns:
//
//	string: The listening address, or empty string if not listening
func (w *Webserver) GetAddress() string {
	return w.address
}

// GetAdapter returns the underlying adapter.
//
// Returns:
//
//	interfaces.AdapterInterface: The underlying adapter
func (w *Webserver) GetAdapter() interfaces.AdapterInterface {
	return w.adapter
}

// SetAdapter sets the underlying adapter.
// This should typically only be called by the builder.
//
// Parameters:
//
//	adapter: The adapter to use
func (w *Webserver) SetAdapter(adapter interfaces.AdapterInterface) {
	w.adapter = adapter

	// Apply existing middleware to the new adapter
	if len(w.middleware) > 0 {
		adapter.Use(w.middleware...)
	}

	// Apply existing configuration to the new adapter
	for key, value := range w.config {
		adapter.SetConfig(key, value)
	}
}

// Private helper methods

// registerRoute is a helper method that registers a route with the specified method and path.
// This method handles the common logic for all HTTP method registration functions.
//
// Parameters:
//
//	method: The HTTP method enum
//	path: The URL path pattern
//	handler: The handler function
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance for chaining
func (w *Webserver) registerRoute(method enums.HTTPMethod, path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	// Create route object for tracking
	route := routing.NewRoute(method, path, handler)
	w.routes.AddRoute(route)

	// Register with adapter if available
	if w.adapter != nil {
		w.adapter.Handle(method.String(), path, handler)
	}

	return w
}

// Route group wrapper
// This struct implements WebserverInterface but prefixes all route registrations

type routeGroup struct {
	webserver *Webserver
	prefix    string
}

// Route registration methods for route group (these prefix the paths)

func (rg *routeGroup) Get(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Get(rg.prefix+path, handler)
}

func (rg *routeGroup) Post(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Post(rg.prefix+path, handler)
}

func (rg *routeGroup) Put(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Put(rg.prefix+path, handler)
}

func (rg *routeGroup) Patch(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Patch(rg.prefix+path, handler)
}

func (rg *routeGroup) Delete(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Delete(rg.prefix+path, handler)
}

func (rg *routeGroup) Options(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Options(rg.prefix+path, handler)
}

func (rg *routeGroup) Head(path string, handler interfaces.HandlerInterface) interfaces.WebserverInterface {
	return rg.webserver.Head(rg.prefix+path, handler)
}

// Other methods delegate to the main webserver

func (rg *routeGroup) Group(prefix string, handler func(interfaces.WebserverInterface)) interfaces.WebserverInterface {
	return rg.webserver.Group(rg.prefix+prefix, handler)
}

func (rg *routeGroup) Use(middleware ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	return rg.webserver.Use(middleware...)
}

func (rg *routeGroup) Middleware(middleware ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	return rg.webserver.Middleware(middleware...)
}

func (rg *routeGroup) Listen(addr ...string) error {
	return rg.webserver.Listen(addr...)
}

func (rg *routeGroup) ListenTLS(certFile, keyFile string, addr ...string) error {
	return rg.webserver.ListenTLS(certFile, keyFile, addr...)
}

func (rg *routeGroup) Shutdown(ctx context.Context) error {
	return rg.webserver.Shutdown(ctx)
}

func (rg *routeGroup) SetConfig(key string, value interface{}) interfaces.WebserverInterface {
	return rg.webserver.SetConfig(key, value)
}

func (rg *routeGroup) GetConfig(key string) interface{} {
	return rg.webserver.GetConfig(key)
}

func (rg *routeGroup) Static(prefix, root string) interfaces.WebserverInterface {
	return rg.webserver.Static(prefix, root)
}

func (rg *routeGroup) Port(port int) interfaces.WebserverInterface {
	return rg.webserver.Port(port)
}

func (rg *routeGroup) Host(host string) interfaces.WebserverInterface {
	return rg.webserver.Host(host)
}

// Ensure Webserver implements WebserverInterface at compile time
var _ interfaces.WebserverInterface = (*Webserver)(nil)
var _ interfaces.WebserverInterface = (*routeGroup)(nil)
