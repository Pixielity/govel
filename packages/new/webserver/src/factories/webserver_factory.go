// Package factories provides factory functions for creating webserver instances.
// This file contains the webserver factory that creates webserver instances from adapters.
package factories

import (
	webserver "govel/new/webserver"
	"govel/new/webserver/interfaces"
)

// WebserverFactory provides factory methods for creating webserver instances.
// This factory encapsulates the creation of webserver objects that wrap framework-specific
// adapters while exposing a unified interfaces.WebserverInterface.
//
// Responsibilities:
//   - Construct concrete webserver instances
//   - Initialize configuration, routing, and middleware containers
//   - Provide convenience constructors with pre-configured state
//
// Thread Safety:
//
//	The factory methods are safe for concurrent use. Each call returns a new
//	independent webserver instance without shared mutable state.
type WebserverFactory struct{}

// NewWebserverFactory creates a new webserver factory instance.
//
// Returns:
//
//	*WebserverFactory: A new webserver factory
func NewWebserverFactory() *WebserverFactory {
	return &WebserverFactory{}
}

// CreateWebserverInstance creates a new webserver instance from a framework adapter.
// This is the primary factory function that wraps adapters in concrete webserver
// implementations that conform to interfaces.WebserverInterface.
//
// The created webserver provides:
//   - HTTP route registration for all standard methods
//   - Route grouping with prefixes and group-level middleware
//   - Global middleware registration
//   - Server lifecycle management (Listen, ListenTLS, Shutdown)
//   - Configuration management (SetConfig, GetConfig)
//   - Static file serving (adapter-dependent)
//
// The function initializes:
//   - An empty configuration map for runtime settings
//   - A fresh routing collection for introspection and tooling
//   - An empty middleware slice for global middleware
//   - Running state and current address fields
//
// Parameters:
//
//	adapter: The initialized adapter to wrap
//
// Returns:
//
//	interfaces.WebserverInterface: A new webserver instance
func (f *WebserverFactory) CreateWebserverInstance(adapter interfaces.AdapterInterface) interfaces.WebserverInterface {
	// Use the constructor method NewWithAdapter instead of direct struct literal
	return webserver.NewWithAdapter(adapter)
}

// CreateWebserverWithConfig creates a webserver instance pre-populated with configuration.
// This convenience constructor applies configuration to both the webserver instance and
// underlying adapter, keeping configuration in sync.
//
// Configuration Keys:
//
//	Common keys include (but are not limited to):
//	- "host": string - Host address to bind (e.g., "0.0.0.0", "localhost")
//	- "port": int - Port to listen on (e.g., 8080)
//	- "timeout": time.Duration - Server timeout (adapter dependent)
//	- "max_body_size": int64 - Maximum request body size (adapter dependent)
//	- "debug": bool - Enable debug mode (adapter dependent)
//
// Notes:
//   - Unknown keys are passed through to the adapter as-is
//   - Adapters may support additional keys specific to their framework
//   - Invalid values may result in runtime errors from the adapter
//
// Parameters:
//
//	adapter: The initialized adapter to wrap
//	config: Initial configuration for the webserver
//
// Returns:
//
//	interfaces.WebserverInterface: A new webserver instance with configuration
func (f *WebserverFactory) CreateWebserverWithConfig(adapter interfaces.AdapterInterface, config map[string]interface{}) interfaces.WebserverInterface {
	// Create base webserver instance
	server := f.CreateWebserverInstance(adapter)

	// Apply configuration using the SetConfig method
	for key, value := range config {
		server.SetConfig(key, value)
	}

	return server
}

// CreateWebserverWithMiddleware creates a webserver instance with global middleware pre-registered.
// Middleware are registered in the order provided and will execute in that order for all routes.
//
// Middleware Execution Model:
//   - Before() methods execute from first to last
//   - Handle() calls wrap the next handler (outermost is first registered)
//   - After() methods execute from last to first
//
// Use Cases:
//   - Application-wide logging, recovery, CORS, auth, rate limiting, etc.
//   - Pre-configured server instances for testing environments
//
// Parameters:
//
//	adapter: The initialized adapter to wrap
//	middleware: Middleware to pre-register on the webserver
//
// Returns:
//
//	interfaces.WebserverInterface: A new webserver instance with middleware
func (f *WebserverFactory) CreateWebserverWithMiddleware(adapter interfaces.AdapterInterface, middleware ...interfaces.MiddlewareInterface) interfaces.WebserverInterface {
	// Create base webserver instance
	server := f.CreateWebserverInstance(adapter)

	// Apply middleware using the Use method
	server.Use(middleware...)

	return server
}

// Package-level factory function
// Convenience function to create a webserver instance without explicitly creating a factory.
var defaultFactory = NewWebserverFactory()

// CreateWebserverInstance creates a webserver instance using the default factory.
//
// Parameters:
//
//	adapter: The initialized adapter to wrap
//
// Returns:
//
//	interfaces.WebserverInterface: A new webserver instance
func CreateWebserverInstance(adapter interfaces.AdapterInterface) interfaces.WebserverInterface {
	return defaultFactory.CreateWebserverInstance(adapter)
}
