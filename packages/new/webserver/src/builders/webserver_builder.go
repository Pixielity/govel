// Package builders provides construction utilities for building webserver instances.
// This file implements the WebserverBuilder which uses the builder pattern to provide
// a fluent API for configuring and constructing webserver instances.
package builders

import (
	"fmt"
	"time"

	"govel/packages/new/webserver/src/enums"
	"govel/packages/new/webserver/src/factories"
	"govel/packages/new/webserver/src/helpers"
	"govel/packages/new/webserver/src/interfaces"
)

// WebserverBuilder implements the builder pattern for constructing webserver instances.
// It provides a Laravel-inspired fluent API that allows method chaining for configuration.
//
// The builder supports:
//   - Engine/adapter selection (GoFiber, Gin, Echo)
//   - Host and port configuration
//   - Custom configuration values
//   - Global middleware registration
//   - Custom adapter injection for testing
//
// Example usage:
//
//	server := Configure().
//	    WithEngine(enums.GoFiber).
//	    WithPort(8080).
//	    WithHost("0.0.0.0").
//	    WithMiddleware(cors.New(), logger.New()).
//	    Set("timeout", 30).
//	    Build()
//
// Thread safety: WebserverBuilder is not thread-safe. Each goroutine should use its own builder instance.
type WebserverBuilder struct {
	// engine stores the selected web framework engine (GoFiber, Gin, Echo)
	engine enums.Engine

	// adapter stores a custom adapter implementation if provided via WithAdapter()
	// This takes precedence over the engine selection
	adapter interfaces.AdapterInterface

	// config stores configuration key-value pairs for the webserver and underlying adapter
	// Common keys include: "host", "port", "address", "timeout", "max_body_size", etc.
	config map[string]interface{}

	// middleware stores global middleware that will be applied to all routes
	// Middleware is executed in the order it was added to this slice
	middleware []interfaces.MiddlewareInterface

	// built indicates whether Build() has been called on this builder
	// This prevents multiple calls to Build() which could cause issues
	built bool

	// configHelper provides configuration utilities
	configHelper *helpers.ConfigHelper

	// validationHelper provides validation utilities
	validationHelper *helpers.ValidationHelper
}

// Configure creates a new WebserverBuilder instance with default configuration.
// The builder is initialized with sensible defaults:
//   - Engine: Default engine from enums.DefaultEngine() (typically GoFiber)
//   - Host: "localhost"
//   - Port: 8080
//   - Empty middleware stack
//   - Empty custom configuration
//
// Returns:
//
//	*WebserverBuilder: A new builder instance ready for configuration
//
// Example:
//
//	builder := Configure()
//	server := builder.WithPort(3000).Build()
func Configure() *WebserverBuilder {
	return &WebserverBuilder{
		engine:           enums.DefaultEngine(),
		adapter:          nil,
		config:           make(map[string]interface{}),
		middleware:       make([]interfaces.MiddlewareInterface, 0),
		built:            false,
		configHelper:     helpers.NewConfigHelper(),
		validationHelper: helpers.NewValidationHelper(),
	}
}

// Engine Selection Methods

// WithEngine sets the web framework engine/adapter to use for the webserver.
// This method accepts an enums.Engine value and will create the appropriate adapter during Build().
//
// Parameters:
//
//	engine: The engine to use (enums.GoFiber, enums.Gin, enums.Echo)
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithEngine(enums.Gin)
//
// Note: If WithAdapter() is called after this, the custom adapter will take precedence.
func (b *WebserverBuilder) WithEngine(engine interface{}) interfaces.BuilderInterface {
	if eng, ok := engine.(enums.Engine); ok {
		b.engine = eng
		// Clear any custom adapter since we're using engine-based selection
		b.adapter = nil
	}
	return b
}

// WithAdapter sets a custom adapter implementation for the webserver.
// This method allows injection of custom or test adapters, overriding engine selection.
//
// Parameters:
//
//	adapter: The custom adapter implementation
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	customAdapter := &MyCustomAdapter{}
//	builder.WithAdapter(customAdapter)
//
// Note: This overrides any engine set via WithEngine().
func (b *WebserverBuilder) WithAdapter(adapter interfaces.AdapterInterface) interfaces.BuilderInterface {
	b.adapter = adapter
	return b
}

// Configuration Methods

// WithConfig sets multiple configuration values at once from a map.
// Existing configuration values are preserved unless overridden.
//
// Parameters:
//
//	config: Map of configuration key-value pairs
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	config := map[string]interface{}{
//	    "timeout":       30,
//	    "max_body_size": 1024 * 1024, // 1MB
//	    "enable_cors":   true,
//	}
//	builder.WithConfig(config)
func (b *WebserverBuilder) WithConfig(config map[string]interface{}) interfaces.BuilderInterface {
	for key, value := range config {
		b.config[key] = value
	}
	return b
}

// Set sets a single configuration key-value pair.
// This is a convenience method for setting individual configuration options.
//
// Parameters:
//
//	key: The configuration key
//	value: The configuration value
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.Set("timeout", 30).
//	        Set("debug", true).
//	        Set("api_key", "secret123")
func (b *WebserverBuilder) Set(key string, value interface{}) interfaces.BuilderInterface {
	b.config[key] = value
	return b
}

// Network Configuration Methods

// WithPort sets the listening port for the webserver.
// This is a convenience method that sets the "port" configuration key.
//
// Parameters:
//
//	port: The port number (1-65535)
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithPort(8080)  // Listen on port 8080
//	builder.WithPort(443)   // Listen on port 443 (HTTPS)
//
// Note: Port validation is performed during Build(), not here.
func (b *WebserverBuilder) WithPort(port int) interfaces.BuilderInterface {
	b.config["port"] = port
	return b
}

// WithHost sets the listening host/address for the webserver.
// This is a convenience method that sets the "host" configuration key.
//
// Parameters:
//
//	host: The host address (e.g., "localhost", "0.0.0.0", "127.0.0.1")
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithHost("0.0.0.0")      // Listen on all interfaces
//	builder.WithHost("localhost")    // Listen only on localhost
//	builder.WithHost("192.168.1.10") // Listen on specific IP
func (b *WebserverBuilder) WithHost(host string) interfaces.BuilderInterface {
	b.config["host"] = host
	return b
}

// WithAddress sets the complete listening address (host:port) for the webserver.
// This method parses the address and sets both host and port configuration.
// If provided, this overrides separate host/port settings.
//
// Parameters:
//
//	addr: The complete address in "host:port" format
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithAddress("localhost:8080")
//	builder.WithAddress("0.0.0.0:443")
//	builder.WithAddress(":3000")         // Empty host defaults to "0.0.0.0"
//
// Note: Invalid address formats will cause Build() to fail.
func (b *WebserverBuilder) WithAddress(addr string) interfaces.BuilderInterface {
	// Use config helper to parse and apply address
	if err := b.configHelper.ApplyAddressToConfig(b.config, addr); err != nil {
		// Log error but don't fail - validation will catch it later
		// In a real implementation, might want to handle this differently
		b.config["address"] = addr
	}
	return b
}

// Middleware Configuration

// WithMiddleware registers one or more global middleware functions.
// Global middleware is applied to all routes and executed in registration order.
//
// Parameters:
//
//	middleware: One or more middleware implementations
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	corsMiddleware := &CorsMiddleware{}
//	authMiddleware := &AuthMiddleware{}
//	loggerMiddleware := &LoggerMiddleware{}
//
//	builder.WithMiddleware(corsMiddleware, authMiddleware, loggerMiddleware)
//
// Note: Middleware order matters - they execute in the order provided.
func (b *WebserverBuilder) WithMiddleware(middleware ...interfaces.MiddlewareInterface) interfaces.BuilderInterface {
	b.middleware = append(b.middleware, middleware...)
	return b
}

// Security & Authentication Methods

// WithTLS configures TLS/SSL certificates for HTTPS.
// This method sets the certificate and key file paths for secure connections.
//
// Parameters:
//
//	certFile: Path to the SSL certificate file
//	keyFile: Path to the SSL private key file
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithTLS("/path/to/cert.pem", "/path/to/key.pem")
func (b *WebserverBuilder) WithTLS(certFile, keyFile string) interfaces.BuilderInterface {
	b.config["tls_cert_file"] = certFile
	b.config["tls_key_file"] = keyFile
	b.config["tls_enabled"] = true
	return b
}

// WithJWTSecret sets the JWT signing secret for authentication.
// This secret is used for signing and verifying JWT tokens.
//
// Parameters:
//
//	secret: The secret key used for signing JWT tokens
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithJWTSecret("your-super-secret-key")
func (b *WebserverBuilder) WithJWTSecret(secret string) interfaces.BuilderInterface {
	b.config["jwt_secret"] = secret
	return b
}

// WithCORS configures Cross-Origin Resource Sharing settings.
// This method sets up CORS policies for handling cross-origin requests.
//
// Parameters:
//
//	config: CORS configuration map with keys like:
//	        - "allow_origins": []string or "*" for allowed origins
//	        - "allow_methods": []string for allowed HTTP methods
//	        - "allow_headers": []string for allowed headers
//	        - "allow_credentials": bool for credential support
//	        - "max_age": int for preflight cache duration
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	corsConfig := map[string]interface{}{
//	    "allow_origins": []string{"https://example.com", "https://app.com"},
//	    "allow_methods": []string{"GET", "POST", "PUT", "DELETE"},
//	    "allow_headers": []string{"Content-Type", "Authorization"},
//	    "allow_credentials": true,
//	    "max_age": 86400,
//	}
//	builder.WithCORS(corsConfig)
func (b *WebserverBuilder) WithCORS(config map[string]interface{}) interfaces.BuilderInterface {
	b.config["cors_config"] = config
	b.config["cors_enabled"] = true
	return b
}

// WithSecurityHeaders adds security headers to all responses.
// These headers help protect against common web vulnerabilities.
//
// Parameters:
//
//	headers: Map of security headers to add, such as:
//	         - "X-Frame-Options": "DENY" or "SAMEORIGIN"
//	         - "X-XSS-Protection": "1; mode=block"
//	         - "X-Content-Type-Options": "nosniff"
//	         - "Strict-Transport-Security": "max-age=31536000; includeSubDomains"
//	         - "Content-Security-Policy": "default-src 'self'"
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	securityHeaders := map[string]string{
//	    "X-Frame-Options": "DENY",
//	    "X-XSS-Protection": "1; mode=block",
//	    "X-Content-Type-Options": "nosniff",
//	}
//	builder.WithSecurityHeaders(securityHeaders)
func (b *WebserverBuilder) WithSecurityHeaders(headers map[string]string) interfaces.BuilderInterface {
	b.config["security_headers"] = headers
	return b
}

// WithRateLimit configures rate limiting for requests.
// This helps prevent abuse and ensures fair resource usage.
//
// Parameters:
//
//	requests: Maximum number of requests allowed within the time window
//	window: Time window for the rate limit (e.g., time.Minute, time.Hour)
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithRateLimit(100, time.Minute) // 100 requests per minute
//	builder.WithRateLimit(1000, time.Hour)  // 1000 requests per hour
func (b *WebserverBuilder) WithRateLimit(requests int, window time.Duration) interfaces.BuilderInterface {
	b.config["rate_limit_requests"] = requests
	b.config["rate_limit_window"] = window
	b.config["rate_limit_enabled"] = true
	return b
}

// WithBasicAuth enables HTTP Basic Authentication.
// This provides a simple username/password authentication mechanism.
//
// Parameters:
//
//	username: The username for basic authentication
//	password: The password for basic authentication
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithBasicAuth("admin", "secret123")
//
// Note: Consider using more secure authentication methods in production.
func (b *WebserverBuilder) WithBasicAuth(username, password string) interfaces.BuilderInterface {
	b.config["basic_auth_username"] = username
	b.config["basic_auth_password"] = password
	b.config["basic_auth_enabled"] = true
	return b
}

// WithAPIKeys configures API key authentication.
// This enables authentication using API keys in headers or query parameters.
//
// Parameters:
//
//	keys: Slice of valid API keys that will be accepted
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	apiKeys := []string{"key1", "key2", "key3"}
//	builder.WithAPIKeys(apiKeys)
func (b *WebserverBuilder) WithAPIKeys(keys []string) interfaces.BuilderInterface {
	b.config["api_keys"] = keys
	b.config["api_key_auth_enabled"] = true
	return b
}

// Performance & Optimization Methods

// WithTimeout sets the general request timeout.
// This applies to the overall request processing time.
//
// Parameters:
//
//	timeout: Request timeout duration
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithTimeout(30 * time.Second)
func (b *WebserverBuilder) WithTimeout(timeout time.Duration) interfaces.BuilderInterface {
	b.config["request_timeout"] = timeout
	return b
}

// WithReadTimeout sets the maximum duration for reading the entire request.
// This includes the request headers and body.
//
// Parameters:
//
//	timeout: Read timeout duration
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithReadTimeout(10 * time.Second)
func (b *WebserverBuilder) WithReadTimeout(timeout time.Duration) interfaces.BuilderInterface {
	b.config["read_timeout"] = timeout
	return b
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
// This applies to response writing operations.
//
// Parameters:
//
//	timeout: Write timeout duration
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithWriteTimeout(10 * time.Second)
func (b *WebserverBuilder) WithWriteTimeout(timeout time.Duration) interfaces.BuilderInterface {
	b.config["write_timeout"] = timeout
	return b
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request.
// This applies to keep-alive connections.
//
// Parameters:
//
//	timeout: Idle timeout duration
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithIdleTimeout(60 * time.Second)
func (b *WebserverBuilder) WithIdleTimeout(timeout time.Duration) interfaces.BuilderInterface {
	b.config["idle_timeout"] = timeout
	return b
}

// WithKeepAlive configures HTTP keep-alive connections.
// Keep-alive allows connection reuse for multiple requests.
//
// Parameters:
//
//	enabled: Whether to enable keep-alive connections
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithKeepAlive(true)  // Enable keep-alive
//	builder.WithKeepAlive(false) // Disable keep-alive
func (b *WebserverBuilder) WithKeepAlive(enabled bool) interfaces.BuilderInterface {
	b.config["keep_alive_enabled"] = enabled
	return b
}

// WithMaxBodySize sets the maximum allowed request body size.
// This helps prevent memory exhaustion from large requests.
//
// Parameters:
//
//	size: Maximum body size in bytes
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithMaxBodySize(10 * 1024 * 1024) // 10MB limit
//	builder.WithMaxBodySize(100 * 1024)       // 100KB limit
func (b *WebserverBuilder) WithMaxBodySize(size int64) interfaces.BuilderInterface {
	b.config["max_body_size"] = size
	return b
}

// WithCompression enables or disables response compression.
// Compression reduces bandwidth usage for text-based responses.
//
// Parameters:
//
//	enabled: Whether to enable response compression (gzip, deflate)
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithCompression(true)  // Enable compression
//	builder.WithCompression(false) // Disable compression
func (b *WebserverBuilder) WithCompression(enabled bool) interfaces.BuilderInterface {
	b.config["compression_enabled"] = enabled
	return b
}

// WithConcurrency sets the maximum number of concurrent connections.
// This helps control server resource usage under high load.
//
// Parameters:
//
//	limit: Maximum concurrent connections
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithConcurrency(1000) // Allow up to 1000 concurrent connections
func (b *WebserverBuilder) WithConcurrency(limit int) interfaces.BuilderInterface {
	b.config["max_concurrent_connections"] = limit
	return b
}

// WithCaching configures caching settings.
// This sets up response caching to improve performance.
//
// Parameters:
//
//	config: Caching configuration map with keys like:
//	        - "provider": string ("memory", "redis", "file")
//	        - "ttl": time.Duration for default cache time-to-live
//	        - "max_size": int64 for maximum cache size in bytes
//	        - "redis_url": string for Redis connection (if using Redis)
//	        - "cache_dir": string for file cache directory (if using file cache)
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	cacheConfig := map[string]interface{}{
//	    "provider": "memory",
//	    "ttl": 5 * time.Minute,
//	    "max_size": 100 * 1024 * 1024, // 100MB
//	}
//	builder.WithCaching(cacheConfig)
func (b *WebserverBuilder) WithCaching(config map[string]interface{}) interfaces.BuilderInterface {
	b.config["cache_config"] = config
	b.config["cache_enabled"] = true
	return b
}

// Monitoring Methods

// WithHealthCheck adds a health check endpoint.
// This provides a standard endpoint for monitoring system health.
//
// Parameters:
//
//	path: The path for the health check endpoint (e.g., "/health", "/status")
//	handler: The health check handler function
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithHealthCheck("/health", myHealthHandler)
func (b *WebserverBuilder) WithHealthCheck(path string, handler interface{}) interfaces.BuilderInterface {
	b.config["health_check_path"] = path
	b.config["health_check_handler"] = handler
	b.config["health_check_enabled"] = true
	return b
}

// WithTracing enables or disables distributed tracing.
// Tracing helps with debugging and monitoring request flows across services.
//
// Parameters:
//
//	enabled: Whether to enable distributed tracing
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithTracing(true)  // Enable tracing
//	builder.WithTracing(false) // Disable tracing
func (b *WebserverBuilder) WithTracing(enabled bool) interfaces.BuilderInterface {
	b.config["tracing_enabled"] = enabled
	return b
}

// Static Content & Assets Methods

// WithStaticFiles configures static file serving.
// This sets up a file server for static assets like CSS, JS, images.
//
// Parameters:
//
//	prefix: URL prefix for static files (e.g., "/static", "/assets")
//	directory: File system directory containing static files
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithStaticFiles("/assets", "./public/assets")
//	// Files in ./public/assets/ will be served at /assets/*
func (b *WebserverBuilder) WithStaticFiles(prefix, directory string) interfaces.BuilderInterface {
	// Initialize static files config if it doesn't exist
	if b.config["static_files"] == nil {
		b.config["static_files"] = make(map[string]string)
	}
	staticFiles := b.config["static_files"].(map[string]string)
	staticFiles[prefix] = directory
	b.config["static_files_enabled"] = true
	return b
}

// WithPublicDirectory sets the public assets directory.
// This is a convenience method for serving files from a public directory at the root.
//
// Parameters:
//
//	directory: Directory path for public assets
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithPublicDirectory("./public")
//	// Files in ./public/ will be served at /*
func (b *WebserverBuilder) WithPublicDirectory(directory string) interfaces.BuilderInterface {
	b.config["public_directory"] = directory
	b.config["public_directory_enabled"] = true
	return b
}

// WithFileServer configures multiple file servers with different routes.
// This allows serving different directories at different URL paths.
//
// Parameters:
//
//	routes: Map of URL prefixes to directory paths
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	fileRoutes := map[string]string{
//	    "/css": "./assets/css",
//	    "/js": "./assets/js",
//	    "/images": "./assets/images",
//	    "/uploads": "./storage/uploads",
//	}
//	builder.WithFileServer(fileRoutes)
func (b *WebserverBuilder) WithFileServer(routes map[string]string) interfaces.BuilderInterface {
	b.config["file_server_routes"] = routes
	b.config["file_server_enabled"] = true
	return b
}

// WithTemplate configures the template engine and directory.
// This sets up server-side template rendering capabilities.
//
// Parameters:
//
//	engine: Template engine identifier (use enums.TemplateEngine or string)
//	directory: Directory containing template files
//
// Returns:
//
//	interfaces.BuilderInterface: The builder instance for method chaining
//
// Example:
//
//	builder.WithTemplate(enums.JSX, "./templates")
//	builder.WithTemplate("handlebars", "./views")
func (b *WebserverBuilder) WithTemplate(engine interface{}, directory string) interfaces.BuilderInterface {
	// Handle both enum and string types
	var templateEngine enums.TemplateEngine
	var engineStr string
	
	switch e := engine.(type) {
	case enums.TemplateEngine:
		templateEngine = e
		engineStr = e.String()
	case string:
		var valid bool
		templateEngine, valid = enums.ParseTemplateEngine(e)
		if !valid {
			// Default to HTML if invalid
			templateEngine = enums.DefaultTemplateEngine()
		}
		engineStr = e
	default:
		// Default to HTML template engine
		templateEngine = enums.DefaultTemplateEngine()
		engineStr = templateEngine.String()
	}
	
	b.config["template_engine"] = templateEngine
	b.config["template_engine_string"] = engineStr
	b.config["template_directory"] = directory
	b.config["template_enabled"] = true
	return b
}

// Build constructs and returns a fully configured webserver instance.
// This method:
//  1. Validates the current configuration
//  2. Creates or uses the specified adapter
//  3. Initializes the adapter with configuration and middleware
//  4. Returns a WebserverInterface implementation
//
// Returns:
//
//	interfaces.WebserverInterface: The configured webserver instance
//
// Panics:
//   - If Build() has already been called on this builder
//   - If the selected engine is invalid or unsupported
//   - If adapter creation/initialization fails
//   - If configuration validation fails
//
// Example:
//
//	server := builder.
//	    WithEngine(enums.GoFiber).
//	    WithPort(8080).
//	    WithMiddleware(corsMiddleware).
//	    Build()
//
//	// Now use the server
//	server.Get("/", homeHandler)
//	server.Listen()
//
// Note: After calling Build(), the builder instance should not be reused.
func (b *WebserverBuilder) Build() interfaces.WebserverInterface {
	// Prevent multiple calls to Build()
	if b.built {
		panic("WebserverBuilder.Build() called multiple times. Create a new builder instance for each webserver.")
	}
	b.built = true

	// Validate configuration before building
	if err := b.validateConfiguration(); err != nil {
		panic(fmt.Sprintf("Invalid webserver configuration: %v", err))
	}

	// Determine which adapter to use
	var adapter interfaces.AdapterInterface
	if b.adapter != nil {
		// Use custom adapter if provided
		adapter = b.adapter
	} else {
		// Create adapter based on engine selection
		var err error
		adapter, err = b.createAdapterFromEngine()
		if err != nil {
			panic(fmt.Sprintf("Failed to create adapter: %v", err))
		}
	}

	// Initialize the adapter with configuration and middleware
	if err := adapter.Init(b.config, b.middleware); err != nil {
		panic(fmt.Sprintf("Failed to initialize adapter: %v", err))
	}

	// Create and return the webserver instance
	// The actual webserver implementation will wrap the adapter
	return b.createWebserverInstance(adapter)
}

// Private helper methods

// validateConfiguration performs validation on the current configuration.
// This method checks for common configuration errors and ensures the webserver
// can be built successfully.
//
// Returns:
//
//	error: Validation error, or nil if configuration is valid
//
// Validates:
//   - Port numbers are in valid range (1-65535)
//   - Host addresses are properly formatted
//   - Engine selection is valid
//   - Required configuration keys are present
func (b *WebserverBuilder) validateConfiguration() error {
	// Use validation helper for comprehensive validation
	var engine *enums.Engine
	if b.adapter == nil {
		engine = &b.engine
	}
	return b.validationHelper.ValidateConfiguration(b.config, engine)
}

// validateAddress is deprecated - use validation helper instead
// This method is kept for backwards compatibility but delegates to the helper
func (b *WebserverBuilder) validateAddress(addr string) error {
	return b.validationHelper.ValidateAddressFormat(addr)
}

// createAdapterFromEngine creates an adapter instance based on the selected engine.
// This method uses the adapter factory to create the appropriate adapter implementation.
//
// Returns:
//
//	interfaces.AdapterInterface: The created adapter
//	error: Creation error, or nil on success
func (b *WebserverBuilder) createAdapterFromEngine() (interfaces.AdapterInterface, error) {
	// Use the adapter factory to create the adapter
	// The factory handles the engine-specific creation logic
	return factories.CreateAdapter(b.engine)
}

// createWebserverInstance creates the final webserver instance that wraps the adapter.
// This method creates the high-level WebserverInterface implementation that provides
// the unified API while delegating to the underlying adapter.
//
// Parameters:
//
//	adapter: The initialized adapter to wrap
//
// Returns:
//
//	interfaces.WebserverInterface: The webserver instance
//
// Note: The actual implementation will be in the main webserver package
func (b *WebserverBuilder) createWebserverInstance(adapter interfaces.AdapterInterface) interfaces.WebserverInterface {
	// Use the webserver factory to create the instance
	return factories.CreateWebserverInstance(adapter)
}

// Utility methods for configuration access

// GetConfig retrieves a configuration value by key.
// This method is primarily for internal use and testing.
//
// Parameters:
//
//	key: The configuration key
//
// Returns:
//
//	interface{}: The configuration value, or nil if not found
func (b *WebserverBuilder) GetConfig(key string) interface{} {
	return b.config[key]
}

// GetEngine returns the currently selected engine.
// This method is primarily for internal use and testing.
//
// Returns:
//
//	enums.Engine: The selected engine
func (b *WebserverBuilder) GetEngine() enums.Engine {
	return b.engine
}

// GetMiddleware returns a copy of the registered middleware stack.
// This method is primarily for internal use and testing.
//
// Returns:
//
//	[]interfaces.MiddlewareInterface: Copy of the middleware stack
func (b *WebserverBuilder) GetMiddleware() []interfaces.MiddlewareInterface {
	// Return a copy to prevent external modification
	middlewareCopy := make([]interfaces.MiddlewareInterface, len(b.middleware))
	copy(middlewareCopy, b.middleware)
	return middlewareCopy
}

// IsBuilt returns whether Build() has been called on this builder.
// This method is primarily for internal use and testing.
//
// Returns:
//
//	bool: True if Build() has been called, false otherwise
func (b *WebserverBuilder) IsBuilt() bool {
	return b.built
}

// Reset resets the builder to its initial state, allowing it to be reused.
// This clears all configuration, middleware, and resets the build state.
// Use with caution as this will lose all current configuration.
//
// Example:
//
//	builder := Configure().
//	    WithPort(8080).
//	    WithMiddleware(middleware1)
//
//	server1 := builder.Build() // This will panic on second call
//
//	// To reuse:
//	builder.Reset()
//	server2 := builder.WithPort(9090).Build() // Now works
func (b *WebserverBuilder) Reset() {
	b.engine = enums.DefaultEngine()
	b.adapter = nil
	b.config = make(map[string]interface{})
	b.middleware = make([]interfaces.MiddlewareInterface, 0)
	b.built = false
}

// Compile-time interface compliance check
var _ interfaces.BuilderInterface = (*WebserverBuilder)(nil)
