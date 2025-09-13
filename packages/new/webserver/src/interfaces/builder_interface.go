// Package interfaces - Builder interface definition
// This file defines the BuilderInterface contract for configuring and constructing the webserver.
package interfaces

import (
	"time"
)

// BuilderInterface provides a fluent API to configure and build a webserver instance
// with a selected engine (adapter), global middleware, and configuration.
//
// Typical usage:
//   server := NewBuilder().
//      WithEngine(enums.GoFiber).
//      WithPort(8080).
//      WithHost("0.0.0.0").
//      WithMiddleware(Logger(), Recovery()).
//      Build()
type BuilderInterface interface {
	// Engine Selection
	
	// WithEngine selects the engine/adapter to use (e.g., GoFiber, Gin, Echo).
	// The engine is typically provided as an enum in the enums package.
	WithEngine(engine interface{}) BuilderInterface
	
	// WithAdapter explicitly sets a custom adapter implementation.
	// This overrides WithEngine and is useful for testing or custom engines.
	WithAdapter(adapter AdapterInterface) BuilderInterface
	
	// Configuration
	
	// WithConfig sets configuration values for the webserver and underlying adapter.
	// These may include timeouts, read/write limits, etc.
	WithConfig(config map[string]interface{}) BuilderInterface
	
	// Set sets a single configuration key-value pair.
	Set(key string, value interface{}) BuilderInterface
	
	// Network Configuration
	
	// WithPort sets the preferred listening port.
	WithPort(port int) BuilderInterface
	
	// WithHost sets the preferred listening host (e.g., 0.0.0.0).
	WithHost(host string) BuilderInterface
	
	// WithAddress sets the full address (host:port). Overrides separate host/port if provided.
	WithAddress(addr string) BuilderInterface
	
	// Middleware
	
	// WithMiddleware registers one or more global middleware to be applied to all routes.
	WithMiddleware(middleware ...MiddlewareInterface) BuilderInterface
	
	// Security & Authentication Methods
	
	// WithTLS configures TLS/SSL certificates for HTTPS.
	// Parameters:
	//   certFile: Path to the SSL certificate file
	//   keyFile: Path to the SSL private key file
	WithTLS(certFile, keyFile string) BuilderInterface
	
	// WithJWTSecret sets the JWT signing secret for authentication.
	// Parameters:
	//   secret: The secret key used for signing JWT tokens
	WithJWTSecret(secret string) BuilderInterface
	
	// WithCORS configures Cross-Origin Resource Sharing settings.
	// Parameters:
	//   config: CORS configuration map with keys like "allow_origins", "allow_methods", etc.
	WithCORS(config map[string]interface{}) BuilderInterface
	
	// WithSecurityHeaders adds security headers to all responses.
	// Parameters:
	//   headers: Map of security headers to add (e.g., "X-Frame-Options", "X-XSS-Protection")
	WithSecurityHeaders(headers map[string]string) BuilderInterface
	
	// WithRateLimit configures rate limiting for requests.
	// Parameters:
	//   requests: Maximum number of requests allowed
	//   window: Time window for the rate limit (e.g., time.Minute, time.Hour)
	WithRateLimit(requests int, window time.Duration) BuilderInterface
	
	// WithBasicAuth enables HTTP Basic Authentication.
	// Parameters:
	//   username: The username for basic auth
	//   password: The password for basic auth
	WithBasicAuth(username, password string) BuilderInterface
	
	// WithAPIKeys configures API key authentication.
	// Parameters:
	//   keys: Slice of valid API keys
	WithAPIKeys(keys []string) BuilderInterface
	
	// Performance & Optimization Methods
	
	// WithTimeout sets the general request timeout.
	// Parameters:
	//   timeout: Request timeout duration
	WithTimeout(timeout time.Duration) BuilderInterface
	
	// WithReadTimeout sets the maximum duration for reading the entire request.
	// Parameters:
	//   timeout: Read timeout duration
	WithReadTimeout(timeout time.Duration) BuilderInterface
	
	// WithWriteTimeout sets the maximum duration before timing out writes of the response.
	// Parameters:
	//   timeout: Write timeout duration
	WithWriteTimeout(timeout time.Duration) BuilderInterface
	
	// WithIdleTimeout sets the maximum amount of time to wait for the next request.
	// Parameters:
	//   timeout: Idle timeout duration
	WithIdleTimeout(timeout time.Duration) BuilderInterface
	
	// WithKeepAlive configures HTTP keep-alive connections.
	// Parameters:
	//   enabled: Whether to enable keep-alive connections
	WithKeepAlive(enabled bool) BuilderInterface
	
	// WithMaxBodySize sets the maximum allowed request body size.
	// Parameters:
	//   size: Maximum body size in bytes
	WithMaxBodySize(size int64) BuilderInterface
	
	// WithCompression enables or disables response compression.
	// Parameters:
	//   enabled: Whether to enable response compression (gzip, deflate)
	WithCompression(enabled bool) BuilderInterface
	
	// WithConcurrency sets the maximum number of concurrent connections.
	// Parameters:
	//   limit: Maximum concurrent connections
	WithConcurrency(limit int) BuilderInterface
	
	// WithCaching configures caching settings.
	// Parameters:
	//   config: Caching configuration map with keys like "provider", "ttl", "max_size", etc.
	WithCaching(config map[string]interface{}) BuilderInterface
	
	// Monitoring Methods
	
	// WithHealthCheck adds a health check endpoint.
	// Parameters:
	//   path: The path for the health check endpoint (e.g., "/health")
	//   handler: The health check handler function
	WithHealthCheck(path string, handler interface{}) BuilderInterface
	
	// WithTracing enables or disables distributed tracing.
	// Parameters:
	//   enabled: Whether to enable distributed tracing
	WithTracing(enabled bool) BuilderInterface
	
	// Static Content & Assets Methods
	
	// WithStaticFiles configures static file serving.
	// Parameters:
	//   prefix: URL prefix for static files (e.g., "/static")
	//   directory: File system directory containing static files
	WithStaticFiles(prefix, directory string) BuilderInterface
	
	// WithPublicDirectory sets the public assets directory.
	// Parameters:
	//   directory: Directory path for public assets
	WithPublicDirectory(directory string) BuilderInterface
	
	// WithFileServer configures multiple file servers with different routes.
	// Parameters:
	//   routes: Map of URL prefixes to directory paths
	WithFileServer(routes map[string]string) BuilderInterface
	
	// WithTemplate configures the template engine and directory.
	// Parameters:
	//   engine: Template engine identifier (use enums.TemplateEngine)
	//   directory: Directory containing template files
	WithTemplate(engine interface{}, directory string) BuilderInterface
	
	// Build constructs and returns a fully configured webserver instance.
	// Building does not start the server; use Listen() on the returned instance.
	Build() WebserverInterface
}
