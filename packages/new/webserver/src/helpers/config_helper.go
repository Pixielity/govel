// Package helpers provides utility functions for webserver configuration and validation.
// This file contains configuration management helpers extracted from the builder pattern.
package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"govel/new/webserver/src/constants"
	"govel/new/webserver/src/enums"
)

// ConfigHelper provides utilities for webserver configuration management and validation.
type ConfigHelper struct{}

// NewConfigHelper creates a new configuration helper instance.
//
// Returns:
//
//	*ConfigHelper: A new configuration helper
func NewConfigHelper() *ConfigHelper {
	return &ConfigHelper{}
}

// SetDefaultConfig sets default configuration values if not already present.
//
// Parameters:
//
//	config: The configuration map to populate with defaults
//
// Returns:
//
//	map[string]interface{}: The configuration map with defaults applied
func (c *ConfigHelper) SetDefaultConfig(config map[string]interface{}) map[string]interface{} {
	if config == nil {
		config = make(map[string]interface{})
	}

	// Set default host if not configured
	if _, exists := config[constants.HOST]; !exists {
		config[constants.HOST] = constants.DEFAULT_HOST
	}

	// Set default port if not configured
	if _, exists := config[constants.PORT]; !exists {
		config[constants.PORT] = constants.DEFAULT_PORT
	}

	// Set default timeout if not configured (in seconds)
	if _, exists := config[constants.TIMEOUT]; !exists {
		config[constants.TIMEOUT] = constants.DEFAULT_TIMEOUT_SECONDS
	}

	// Set default max body size if not configured (in bytes)
	if _, exists := config[constants.MAX_BODY_SIZE]; !exists {
		config[constants.MAX_BODY_SIZE] = constants.DEFAULT_MAX_BODY_SIZE
	}

	// Set default debug mode if not configured
	if _, exists := config[constants.DEBUG]; !exists {
		config[constants.DEBUG] = constants.DEFAULT_DEBUG
	}

	return config
}

// MergeConfig merges source configuration into target configuration.
// Values in source will override values in target.
//
// Parameters:
//
//	target: The target configuration map
//	source: The source configuration map to merge
//
// Returns:
//
//	map[string]interface{}: The merged configuration
func (c *ConfigHelper) MergeConfig(target, source map[string]interface{}) map[string]interface{} {
	if target == nil {
		target = make(map[string]interface{})
	}

	for key, value := range source {
		target[key] = value
	}

	return target
}

// ParseAddress parses an address string and extracts host and port components.
//
// Parameters:
//
//	addr: The address string in "host:port" format
//
// Returns:
//
//	string: The host component
//	int: The port component
//	error: Any parsing error
func (c *ConfigHelper) ParseAddress(addr string) (string, int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("address must be in format 'host:port', got: %s", addr)
	}

	host := strings.TrimSpace(parts[0])
	portStr := strings.TrimSpace(parts[1])

	if portStr == "" {
		return host, 0, fmt.Errorf("port cannot be empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return host, 0, fmt.Errorf("invalid port: %s", portStr)
	}

	return host, port, nil
}

// ApplyAddressToConfig applies a parsed address to the configuration map.
//
// Parameters:
//
//	config: The configuration map to update
//	addr: The address string to parse and apply
//
// Returns:
//
//	error: Any parsing or validation error
func (c *ConfigHelper) ApplyAddressToConfig(config map[string]interface{}, addr string) error {
	host, port, err := c.ParseAddress(addr)
	if err != nil {
		return err
	}

	config["address"] = addr
	if host != "" {
		config["host"] = host
	}
	if port > 0 {
		config["port"] = port
	}

	return nil
}

// GetConfigString retrieves a configuration value as a string.
//
// Parameters:
//
//	config: The configuration map
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	string: The configuration value as string, or default value
func (c *ConfigHelper) GetConfigString(config map[string]interface{}, key string, defaultValue string) string {
	if value, exists := config[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetConfigInt retrieves a configuration value as an integer.
//
// Parameters:
//
//	config: The configuration map
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	int: The configuration value as integer, or default value
func (c *ConfigHelper) GetConfigInt(config map[string]interface{}, key string, defaultValue int) int {
	if value, exists := config[key]; exists {
		if intVal, ok := value.(int); ok {
			return intVal
		}
	}
	return defaultValue
}

// GetConfigBool retrieves a configuration value as a boolean.
//
// Parameters:
//
//	config: The configuration map
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	bool: The configuration value as boolean, or default value
//
// GetConfigBool retrieves a configuration value as a boolean.
// This method safely extracts a boolean value from the configuration map,
// providing type checking and default value fallback.
//
// Parameters:
//
//	config: The configuration map containing key-value pairs
//	key: The configuration key to retrieve
//	defaultValue: The default value returned if key is not found or invalid type
//
// Returns:
//
//	bool: The configuration value as boolean, or defaultValue if not found
//
// Example:
//
//	debug := helper.GetConfigBool(config, "debug", false)
//	compression := helper.GetConfigBool(config, "compression_enabled", true)
func (c *ConfigHelper) GetConfigBool(config map[string]interface{}, key string, defaultValue bool) bool {
	if value, exists := config[key]; exists {
		if boolVal, ok := value.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

// BuildListenAddress constructs the listening address from configuration.
// This method handles various address configuration patterns with intelligent fallbacks:
//   - "address" key takes precedence (e.g., "localhost:8080")
//   - Falls back to combining "host" and "port" keys
//   - Uses default values if neither is configured
//
// The method is particularly useful for webserver initialization where the listening
// address can be specified in multiple ways for flexibility.
//
// Parameters:
//
//	config: The configuration map containing network settings
//
// Returns:
//
//	string: The complete listening address in "host:port" format
//
// Example:
//
//	// With complete address
//	config["address"] = "0.0.0.0:8080"
//	addr := helper.BuildListenAddress(config) // Returns: "0.0.0.0:8080"
//
//	// With separate host and port
//	config["host"] = "127.0.0.1"
//	config["port"] = 3000
//	addr := helper.BuildListenAddress(config) // Returns: "127.0.0.1:3000"
//
//	// With defaults
//	addr := helper.BuildListenAddress(map[string]interface{}{}) // Returns: "localhost:8080"
func (c *ConfigHelper) BuildListenAddress(config map[string]interface{}) string {
	// Priority 1: Check if complete address is configured
	if addr := c.GetConfigString(config, constants.ADDRESS, ""); addr != "" {
		return addr
	}

	// Priority 2: Build address from separate host and port configuration
	host := c.GetConfigString(config, constants.HOST, constants.DEFAULT_HOST)
	port := c.GetConfigInt(config, constants.PORT, constants.DEFAULT_PORT)

	return fmt.Sprintf("%s:%d", host, port)
}

// Security & Authentication Configuration Methods
// These methods provide secure configuration management for various authentication
// and security mechanisms including TLS, JWT, CORS, and access control.

// SetTLS configures TLS/SSL certificate settings for HTTPS connections.
// This method enables secure communication by setting up the necessary certificate
// and private key file paths. Once configured, the webserver can serve HTTPS traffic.
//
// Security considerations:
//   - Certificate and key files should have proper file permissions (600 or 644)
//   - Files should be stored in secure locations outside the web root
//   - Consider using environment variables for file paths in production
//
// Parameters:
//
//	config: The configuration map to update with TLS settings
//	certFile: Absolute path to the SSL/TLS certificate file (usually .pem or .crt)
//	keyFile: Absolute path to the SSL/TLS private key file (usually .key or .pem)
//
// Example:
//
//	helper.SetTLS(config, "/etc/ssl/certs/server.pem", "/etc/ssl/private/server.key")
//	helper.SetTLS(config, "./certs/localhost.pem", "./certs/localhost-key.pem")
func (c *ConfigHelper) SetTLS(config map[string]interface{}, certFile, keyFile string) {
	config[constants.TLS_CERT_FILE] = certFile
	config[constants.TLS_KEY_FILE] = keyFile
	config[constants.TLS_ENABLED] = true
}

// GetTLS retrieves TLS/SSL certificate configuration settings.
// This method extracts the TLS configuration including certificate file path,
// private key file path, and whether TLS is enabled.
//
// Returns:
//
//	certFile: Path to the SSL certificate file, empty string if not configured
//	keyFile: Path to the SSL private key file, empty string if not configured
//	enabled: Boolean indicating if TLS is enabled (true if SetTLS was called)
//
// Example:
//
//	cert, key, enabled := helper.GetTLS(config)
//	if enabled {
//	    server.ListenTLS(cert, key, ":443")
//	}
func (c *ConfigHelper) GetTLS(config map[string]interface{}) (certFile, keyFile string, enabled bool) {
	certFile = c.GetConfigString(config, constants.TLS_CERT_FILE, "")
	keyFile = c.GetConfigString(config, constants.TLS_KEY_FILE, "")
	enabled = c.GetConfigBool(config, constants.TLS_ENABLED, false)
	return
}

// SetJWTSecret configures the JWT (JSON Web Token) signing secret.
// This secret is used to sign and verify JWT tokens for authentication.
// The secret should be cryptographically strong and kept confidential.
//
// Security best practices:
//   - Use a random, high-entropy secret (minimum 32 characters)
//   - Store the secret in environment variables, not in code
//   - Rotate the secret periodically
//   - Use different secrets for different environments
//
// Parameters:
//
//	config: The configuration map to update with JWT settings
//	secret: The secret key for JWT signing (should be strong and random)
//
// Example:
//
//	helper.SetJWTSecret(config, os.Getenv("JWT_SECRET"))
//	helper.SetJWTSecret(config, "your-256-bit-secret-here")
func (c *ConfigHelper) SetJWTSecret(config map[string]interface{}, secret string) {
	config[constants.JWT_SECRET] = secret
}

// GetJWTSecret retrieves the JWT signing secret from configuration.
// This method safely extracts the JWT secret used for token signing and verification.
//
// Returns:
//
//	string: The JWT signing secret, empty string if not configured
//
// Example:
//
//	secret := helper.GetJWTSecret(config)
//	if secret == "" {
//	    log.Fatal("JWT secret not configured")
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, _ := token.SignedString([]byte(secret))
func (c *ConfigHelper) GetJWTSecret(config map[string]interface{}) string {
	return c.GetConfigString(config, constants.JWT_SECRET, "")
}

// SetCORS configures Cross-Origin Resource Sharing (CORS) settings.
// CORS is a security feature implemented by web browsers to restrict cross-origin
// HTTP requests. This method enables and configures CORS policies for your API.
//
// CORS is essential for:
//   - Single Page Applications (SPAs) accessing APIs from different domains
//   - Mobile apps making API calls to your server
//   - Third-party integrations and webhooks
//   - Development environments with different ports
//
// Common CORS configuration keys:
//   - "allow_origins": []string or "*" - Allowed origin domains
//   - "allow_methods": []string - Permitted HTTP methods (GET, POST, etc.)
//   - "allow_headers": []string - Allowed request headers
//   - "allow_credentials": bool - Whether to allow cookies/auth headers
//   - "max_age": int - Preflight cache duration in seconds
//   - "expose_headers": []string - Headers exposed to the client
//
// Security considerations:
//   - Avoid using "*" for origins in production with credentials enabled
//   - Be specific about allowed origins to prevent unauthorized access
//   - Limit allowed methods to only what your API supports
//
// Parameters:
//
//	config: The configuration map to update with CORS settings
//	corsConfig: Map containing CORS policy configuration
//
// Example:
//
//	corsSettings := map[string]interface{}{
//	    "allow_origins": []string{"https://myapp.com", "https://admin.myapp.com"},
//	    "allow_methods": []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//	    "allow_headers": []string{"Content-Type", "Authorization", "X-API-Key"},
//	    "allow_credentials": true,
//	    "max_age": 86400, // 24 hours
//	}
//	helper.SetCORS(config, corsSettings)
func (c *ConfigHelper) SetCORS(config map[string]interface{}, corsConfig map[string]interface{}) {
	config[constants.CORS_CONFIG] = corsConfig
	config[constants.CORS_ENABLED] = true
}

// GetCORS retrieves the current CORS configuration settings.
// Returns both the CORS configuration map and whether CORS is enabled.
//
// Returns:
//
//	corsConfig: Map containing CORS policy settings, nil if not configured
//	enabled: Boolean indicating if CORS is enabled (true if SetCORS was called)
//
// Example:
//
//	corsConfig, enabled := helper.GetCORS(config)
//	if enabled {
//	    origins := corsConfig["allow_origins"]
//	    methods := corsConfig["allow_methods"]
//	    log.Printf("CORS enabled - Origins: %v, Methods: %v", origins, methods)
//	}
func (c *ConfigHelper) GetCORS(config map[string]interface{}) (corsConfig map[string]interface{}, enabled bool) {
	if cfg, exists := config[constants.CORS_CONFIG]; exists {
		if corsMap, ok := cfg.(map[string]interface{}); ok {
			corsConfig = corsMap
		}
	}
	enabled = c.GetConfigBool(config, constants.CORS_ENABLED, false)
	return
}

// SetSecurityHeaders configures HTTP security headers for enhanced protection.
// Security headers are crucial for protecting web applications against common
// attacks like XSS, clickjacking, MIME-type confusion, and protocol downgrade.
//
// Essential security headers:
//   - "X-Frame-Options": Prevents clickjacking ("DENY", "SAMEORIGIN")
//   - "X-XSS-Protection": Enables XSS filtering ("1; mode=block")
//   - "X-Content-Type-Options": Prevents MIME sniffing ("nosniff")
//   - "Strict-Transport-Security": Enforces HTTPS ("max-age=31536000")
//   - "Content-Security-Policy": Controls resource loading ("default-src 'self'")
//   - "Referrer-Policy": Controls referrer information ("strict-origin-when-cross-origin")
//   - "Permissions-Policy": Controls browser features access
//
// Security benefits:
//   - Protection against Cross-Site Scripting (XSS)
//   - Prevention of clickjacking attacks
//   - MIME-type confusion protection
//   - Enforced HTTPS connections
//   - Content injection prevention
//
// Parameters:
//
//	config: The configuration map to update with security header settings
//	headers: Map of HTTP header names to their values
//
// Example:
//
//	securityHeaders := map[string]string{
//	    "X-Frame-Options":           "DENY",
//	    "X-XSS-Protection":          "1; mode=block",
//	    "X-Content-Type-Options":    "nosniff",
//	    "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
//	    "Content-Security-Policy":   "default-src 'self'; script-src 'self' 'unsafe-inline'",
//	}
//	helper.SetSecurityHeaders(config, securityHeaders)
func (c *ConfigHelper) SetSecurityHeaders(config map[string]interface{}, headers map[string]string) {
	config[constants.SECURITY_HEADERS] = headers
}

// GetSecurityHeaders retrieves the configured HTTP security headers.
// Returns a map of all configured security headers that will be added to responses.
//
// Returns:
//
//	map[string]string: Map of header names to values, empty map if none configured
//
// Example:
//
//	headers := helper.GetSecurityHeaders(config)
//	for name, value := range headers {
//	    log.Printf("Security header: %s = %s", name, value)
//	}
//	if len(headers) == 0 {
//	    log.Println("No security headers configured")
//	}
func (c *ConfigHelper) GetSecurityHeaders(config map[string]interface{}) map[string]string {
	if headers, exists := config[constants.SECURITY_HEADERS]; exists {
		if headerMap, ok := headers.(map[string]string); ok {
			return headerMap
		}
	}
	return make(map[string]string)
}

// SetRateLimit configures request rate limiting to prevent abuse.
// Rate limiting is essential for protecting APIs from abuse, DDoS attacks,
// and ensuring fair usage across all clients. It controls how many requests
// a client can make within a specified time window.
//
// Rate limiting strategies:
//   - Fixed window: Reset count at fixed intervals
//   - Sliding window: More precise, considers request history
//   - Token bucket: Allows bursts within limits
//   - Per-IP, per-user, or per-API key limiting
//
// Implementation considerations:
//   - Choose appropriate limits based on your API capacity
//   - Consider different limits for different endpoints
//   - Implement proper error responses (HTTP 429)
//   - Use Redis or similar for distributed rate limiting
//
// Parameters:
//
//	config: The configuration map to update with rate limiting settings
//	requests: Maximum number of requests allowed within the time window
//	window: Time window duration for the request limit (e.g., time.Minute)
//
// Example:
//
//	helper.SetRateLimit(config, 100, time.Minute)     // 100 requests per minute
//	helper.SetRateLimit(config, 1000, time.Hour)      // 1000 requests per hour
//	helper.SetRateLimit(config, 10, time.Second)      // 10 requests per second
func (c *ConfigHelper) SetRateLimit(config map[string]interface{}, requests int, window time.Duration) {
	config[constants.RATE_LIMIT_REQUESTS] = requests
	config[constants.RATE_LIMIT_WINDOW] = window
	config[constants.RATE_LIMIT_ENABLED] = true
}

// GetRateLimit retrieves the current rate limiting configuration.
// Returns the request limit, time window, and whether rate limiting is enabled.
//
// Returns:
//
//	requests: Maximum requests allowed in the time window (0 if not configured)
//	window: Time window duration (zero duration if not configured)
//	enabled: Boolean indicating if rate limiting is enabled
//
// Example:
//
//	requests, window, enabled := helper.GetRateLimit(config)
//	if enabled {
//	    log.Printf("Rate limit: %d requests per %v", requests, window)
//	} else {
//	    log.Println("Rate limiting is disabled")
//	}
func (c *ConfigHelper) GetRateLimit(config map[string]interface{}) (requests int, window time.Duration, enabled bool) {
	requests = c.GetConfigInt(config, constants.RATE_LIMIT_REQUESTS, 0)
	if w, exists := config[constants.RATE_LIMIT_WINDOW]; exists {
		if duration, ok := w.(time.Duration); ok {
			window = duration
		}
	}
	enabled = c.GetConfigBool(config, constants.RATE_LIMIT_ENABLED, false)
	return
}

// SetBasicAuth configures HTTP Basic Authentication credentials.
// Basic Auth is a simple authentication scheme built into HTTP, where
// credentials are sent as base64-encoded username:password in the
// Authorization header. While simple, it should only be used over HTTPS.
//
// Security considerations:
//   - ALWAYS use HTTPS when Basic Auth is enabled
//   - Use strong, unique passwords
//   - Consider more secure alternatives (JWT, OAuth) for production
//   - Credentials are sent with every request
//   - No built-in logout mechanism
//
// Use cases:
//   - Simple admin interfaces
//   - Development/staging environments
//   - Legacy system integration
//   - Quick prototyping
//
// Parameters:
//
//	config: The configuration map to update with Basic Auth settings
//	username: Username for Basic Authentication
//	password: Password for Basic Authentication (store securely!)
//
// Example:
//
//	helper.SetBasicAuth(config, "admin", os.Getenv("ADMIN_PASSWORD"))
//	helper.SetBasicAuth(config, "api-user", generateSecurePassword())
//
// Security Note: Never hardcode passwords in source code!
func (c *ConfigHelper) SetBasicAuth(config map[string]interface{}, username, password string) {
	config[constants.BASIC_AUTH_USERNAME] = username
	config[constants.BASIC_AUTH_PASSWORD] = password
	config[constants.BASIC_AUTH_ENABLED] = true
}

// GetBasicAuth retrieves HTTP Basic Authentication configuration.
// Returns the configured username, password, and whether Basic Auth is enabled.
//
// Returns:
//
//	username: Configured username for Basic Auth, empty if not set
//	password: Configured password for Basic Auth, empty if not set
//	enabled: Boolean indicating if Basic Auth is enabled
//
// Example:
//
//	username, password, enabled := helper.GetBasicAuth(config)
//	if enabled {
//	    log.Printf("Basic Auth enabled for user: %s", username)
//	    // Note: Avoid logging passwords in production!
//	}
//
// Security Warning: Handle returned credentials carefully!
func (c *ConfigHelper) GetBasicAuth(config map[string]interface{}) (username, password string, enabled bool) {
	username = c.GetConfigString(config, constants.BASIC_AUTH_USERNAME, "")
	password = c.GetConfigString(config, constants.BASIC_AUTH_PASSWORD, "")
	enabled = c.GetConfigBool(config, constants.BASIC_AUTH_ENABLED, false)
	return
}

// SetAPIKeys configures API key authentication for secure API access.
// API key authentication provides a simple yet effective way to authenticate
// API clients using pre-shared secret keys. Keys can be passed via headers,
// query parameters, or request body depending on implementation.
//
// API key benefits:
//   - Simple implementation and usage
//   - Easy to revoke individual keys
//   - Can associate keys with specific clients/applications
//   - Suitable for server-to-server communication
//   - No user session management required
//
// Security best practices:
//   - Generate cryptographically strong keys (minimum 32 characters)
//   - Use HTTPS to protect keys in transit
//   - Implement key rotation policies
//   - Log and monitor key usage
//   - Store keys securely (hashed/encrypted)
//   - Implement rate limiting per key
//
// Parameters:
//
//	config: The configuration map to update with API key settings
//	keys: Slice of valid API keys that will be accepted for authentication
//
// Example:
//
//	apiKeys := []string{
//	    "ak_live_1234567890abcdef1234567890abcdef", // Production key
//	    "ak_test_abcdef1234567890abcdef1234567890", // Test key
//	    "ak_dev_9876543210fedcba9876543210fedcba",  // Development key
//	}
//	helper.SetAPIKeys(config, apiKeys)
//
// Note: In production, consider loading keys from environment variables or secure storage
func (c *ConfigHelper) SetAPIKeys(config map[string]interface{}, keys []string) {
	config[constants.API_KEYS] = keys
	config[constants.API_KEY_AUTH_ENABLED] = true
}

// GetAPIKeys retrieves the current API key authentication configuration.
// Returns the list of valid API keys and whether API key auth is enabled.
//
// Returns:
//
//	keys: Slice of configured API keys, nil if none configured
//	enabled: Boolean indicating if API key authentication is enabled
//
// Example:
//
//	keys, enabled := helper.GetAPIKeys(config)
//	if enabled {
//	    log.Printf("API key authentication enabled with %d keys", len(keys))
//	    // Security: Never log actual API keys in production!
//	} else {
//	    log.Println("API key authentication is disabled")
//	}
func (c *ConfigHelper) GetAPIKeys(config map[string]interface{}) (keys []string, enabled bool) {
	if apiKeys, exists := config[constants.API_KEYS]; exists {
		if keySlice, ok := apiKeys.([]string); ok {
			keys = keySlice
		}
	}
	enabled = c.GetConfigBool(config, constants.API_KEY_AUTH_ENABLED, false)
	return
}

// Performance & Optimization Configuration Methods
// These methods control timing, connection management, and resource limits
// that directly impact server performance and client experience.

// SetTimeout configures the general request timeout duration.
// This timeout applies to the overall request processing time from start to finish.
// It's a catch-all timeout that prevents requests from hanging indefinitely.
//
// Timeout considerations:
//   - Should be longer than expected response times for your slowest endpoint
//   - Consider database query times, external API calls, and processing time
//   - Balance user experience with server resource protection
//   - Typical values: 30-300 seconds depending on application type
//
// Use cases:
//   - API endpoints with database operations
//   - File upload/download processing
//   - Complex business logic operations
//   - Integration with slow external services
//
// Parameters:
//
//	config: The configuration map to update with timeout settings
//	timeout: Maximum duration for request processing
//
// Example:
//
//	helper.SetTimeout(config, 30*time.Second)  // Standard API timeout
//	helper.SetTimeout(config, 5*time.Minute)   // File processing timeout
//	helper.SetTimeout(config, 10*time.Second)  // Fast response requirement
func (c *ConfigHelper) SetTimeout(config map[string]interface{}, timeout time.Duration) {
	config[constants.REQUEST_TIMEOUT] = timeout
}

// GetTimeout retrieves the configured general request timeout.
// Returns the maximum duration allowed for request processing.
//
// Returns:
//
//	time.Duration: Request timeout duration (default: 30 seconds)
//
// Example:
//
//	timeout := helper.GetTimeout(config)
//	log.Printf("Request timeout set to: %v", timeout)
//	if timeout > 60*time.Second {
//	    log.Println("Warning: Long timeout configured")
//	}
func (c *ConfigHelper) GetTimeout(config map[string]interface{}) time.Duration {
	if t, exists := config[constants.REQUEST_TIMEOUT]; exists {
		if timeout, ok := t.(time.Duration); ok {
			return timeout
		}
	}
	return time.Duration(constants.DEFAULT_REQUEST_TIMEOUT) * time.Second
}

// SetReadTimeout configures the maximum duration for reading HTTP requests.
// This timeout covers reading request headers, body, and any trailers.
// It prevents slow or malicious clients from keeping connections open indefinitely.
//
// Read timeout scenarios:
//   - Slow network connections
//   - Large request bodies (file uploads)
//   - Malicious slow-read attacks
//   - Client connection issues
//
// Recommended values:
//   - Small requests: 5-15 seconds
//   - File uploads: 1-5 minutes
//   - Streaming data: 30+ seconds
//
// Parameters:
//
//	config: The configuration map to update with read timeout settings
//	timeout: Maximum duration for reading request data
//
// Example:
//
//	helper.SetReadTimeout(config, 10*time.Second) // Standard web requests
//	helper.SetReadTimeout(config, 2*time.Minute)  // File upload endpoints
//	helper.SetReadTimeout(config, 30*time.Second) // API with large payloads
func (c *ConfigHelper) SetReadTimeout(config map[string]interface{}, timeout time.Duration) {
	config[constants.READ_TIMEOUT] = timeout
}

// GetReadTimeout retrieves the configured HTTP request read timeout.
// Returns the maximum duration allowed for reading request data.
//
// Returns:
//
//	time.Duration: Read timeout duration (default: 10 seconds)
//
// Example:
//
//	readTimeout := helper.GetReadTimeout(config)
//	log.Printf("Read timeout: %v", readTimeout)
//	if readTimeout < 5*time.Second {
//	    log.Println("Warning: Very short read timeout may cause issues")
//	}
func (c *ConfigHelper) GetReadTimeout(config map[string]interface{}) time.Duration {
	if t, exists := config[constants.READ_TIMEOUT]; exists {
		if timeout, ok := t.(time.Duration); ok {
			return timeout
		}
	}
	return time.Duration(constants.DEFAULT_READ_TIMEOUT) * time.Second
}

// SetWriteTimeout configures the maximum duration for writing HTTP responses.
// This timeout covers writing response headers, body, and flushing data to the client.
// It prevents slow clients from causing resource exhaustion on the server.
//
// Write timeout scenarios:
//   - Slow client connections
//   - Large response bodies (file downloads)
//   - Network congestion
//   - Client connection drops
//
// Considerations:
//   - Should account for response size and client connection speed
//   - Larger files need longer timeouts
//   - Streaming responses may need extended timeouts
//   - Balance between user experience and resource protection
//
// Parameters:
//
//	config: The configuration map to update with write timeout settings
//	timeout: Maximum duration for writing response data
//
// Example:
//
//	helper.SetWriteTimeout(config, 10*time.Second) // API responses
//	helper.SetWriteTimeout(config, 5*time.Minute)  // File downloads
//	helper.SetWriteTimeout(config, 30*time.Second) // Large JSON responses
func (c *ConfigHelper) SetWriteTimeout(config map[string]interface{}, timeout time.Duration) {
	config[constants.WRITE_TIMEOUT] = timeout
}

// GetWriteTimeout retrieves the configured HTTP response write timeout.
// Returns the maximum duration allowed for writing response data.
//
// Returns:
//
//	time.Duration: Write timeout duration (default: 10 seconds)
//
// Example:
//
//	writeTimeout := helper.GetWriteTimeout(config)
//	log.Printf("Write timeout: %v", writeTimeout)
//	// Adjust based on expected response sizes
func (c *ConfigHelper) GetWriteTimeout(config map[string]interface{}) time.Duration {
	if t, exists := config[constants.WRITE_TIMEOUT]; exists {
		if timeout, ok := t.(time.Duration); ok {
			return timeout
		}
	}
	return time.Duration(constants.DEFAULT_WRITE_TIMEOUT) * time.Second
}

// SetIdleTimeout configures the maximum time to wait for next request on keep-alive connections.
// This timeout determines how long the server keeps a connection open waiting for
// the next request when HTTP keep-alive is enabled.
//
// Idle timeout benefits:
//   - Prevents resource exhaustion from idle connections
//   - Balances connection reuse with resource management
//   - Allows proper cleanup of abandoned connections
//   - Maintains server responsiveness under load
//
// Tuning considerations:
//   - Longer timeouts improve connection reuse (better performance)
//   - Shorter timeouts free up resources faster (better scalability)
//   - Consider your application's request patterns
//   - Balance with maximum concurrent connections limit
//
// Parameters:
//
//	config: The configuration map to update with idle timeout settings
//	timeout: Maximum duration to wait for next request on keep-alive connections
//
// Example:
//
//	helper.SetIdleTimeout(config, 60*time.Second)  // Standard keep-alive
//	helper.SetIdleTimeout(config, 30*time.Second)  // Conservative setting
//	helper.SetIdleTimeout(config, 2*time.Minute)   // High-performance apps
func (c *ConfigHelper) SetIdleTimeout(config map[string]interface{}, timeout time.Duration) {
	config[constants.IDLE_TIMEOUT] = timeout
}

// GetIdleTimeout retrieves the configured idle timeout for keep-alive connections.
// Returns the maximum duration the server waits for the next request.
//
// Returns:
//
//	time.Duration: Idle timeout duration (default: 60 seconds)
//
// Example:
//
//	idleTimeout := helper.GetIdleTimeout(config)
//	log.Printf("Keep-alive idle timeout: %v", idleTimeout)
func (c *ConfigHelper) GetIdleTimeout(config map[string]interface{}) time.Duration {
	if t, exists := config[constants.IDLE_TIMEOUT]; exists {
		if timeout, ok := t.(time.Duration); ok {
			return timeout
		}
	}
	return time.Duration(constants.DEFAULT_IDLE_TIMEOUT) * time.Second
}

// SetKeepAlive configures HTTP keep-alive connection reuse functionality.
// Keep-alive allows multiple HTTP requests to reuse the same TCP connection,
// reducing connection overhead and improving performance for clients making
// multiple requests.
//
// Keep-alive benefits:
//   - Reduced connection establishment overhead
//   - Lower latency for subsequent requests
//   - Reduced server resource usage
//   - Better client-side performance
//   - Improved throughput for APIs
//
// Considerations:
//   - Uses more server memory (connections stay open longer)
//   - Requires proper idle timeout configuration
//   - May not be beneficial for single-request clients
//   - Can help with HTTP/1.1 pipelining
//
// Parameters:
//
//	config: The configuration map to update with keep-alive settings
//	enabled: Boolean flag to enable (true) or disable (false) keep-alive
//
// Example:
//
//	helper.SetKeepAlive(config, true)  // Enable for better performance
//	helper.SetKeepAlive(config, false) // Disable for simpler connection management
func (c *ConfigHelper) SetKeepAlive(config map[string]interface{}, enabled bool) {
	config[constants.KEEP_ALIVE_ENABLED] = enabled
}

// GetKeepAlive retrieves the current HTTP keep-alive configuration.
// Returns whether connection reuse is enabled for HTTP requests.
//
// Returns:
//
//	bool: True if keep-alive is enabled, false otherwise (default: true)
//
// Example:
//
//	if helper.GetKeepAlive(config) {
//	    log.Println("HTTP keep-alive connections enabled")
//	} else {
//	    log.Println("Each request uses a new connection")
//	}
func (c *ConfigHelper) GetKeepAlive(config map[string]interface{}) bool {
	return c.GetConfigBool(config, constants.KEEP_ALIVE_ENABLED, constants.DEFAULT_KEEP_ALIVE)
}

// SetMaxBodySize configures the maximum allowed HTTP request body size.
// This setting protects the server from memory exhaustion attacks and
// accidental large uploads that could impact server performance.
//
// Body size considerations:
//   - Prevents memory exhaustion from large requests
//   - Must accommodate legitimate use cases (file uploads, data imports)
//   - Should consider available server memory
//   - May need different limits for different endpoints
//
// Common size limits:
//   - APIs: 1-10MB for JSON payloads
//   - File uploads: 10-100MB for documents
//   - Image uploads: 5-50MB for photos
//   - Data imports: 50-500MB for bulk operations
//
// Security implications:
//   - Prevents DoS attacks via large request bodies
//   - Protects against accidental resource consumption
//   - Should be enforced early in request processing
//
// Parameters:
//
//	config: The configuration map to update with body size settings
//	size: Maximum request body size in bytes
//
// Example:
//
//	helper.SetMaxBodySize(config, 10*1024*1024)   // 10MB limit
//	helper.SetMaxBodySize(config, 1024*1024)      // 1MB for APIs
//	helper.SetMaxBodySize(config, 100*1024*1024)  // 100MB for file uploads
func (c *ConfigHelper) SetMaxBodySize(config map[string]interface{}, size int64) {
	config[constants.MAX_BODY_SIZE] = size
}

// GetMaxBodySize retrieves the configured maximum request body size limit.
// Returns the maximum number of bytes allowed in HTTP request bodies.
//
// Returns:
//
//	int64: Maximum body size in bytes (default: 4MB)
//
// Example:
//
//	maxSize := helper.GetMaxBodySize(config)
//	log.Printf("Maximum request body size: %d bytes (%.2f MB)",
//	    maxSize, float64(maxSize)/(1024*1024))
//	if maxSize > 100*1024*1024 {
//	    log.Println("Warning: Large body size limit configured")
//	}
func (c *ConfigHelper) GetMaxBodySize(config map[string]interface{}) int64 {
	if size, exists := config[constants.MAX_BODY_SIZE]; exists {
		if sizeInt64, ok := size.(int64); ok {
			return sizeInt64
		}
		// Handle int to int64 conversion for compatibility
		if sizeInt, ok := size.(int); ok {
			return int64(sizeInt)
		}
	}
	return constants.DEFAULT_MAX_BODY_SIZE
}

// SetCompression configures HTTP response compression settings.
// Enables or disables automatic compression of HTTP responses using algorithms
// like gzip and deflate. Compression reduces bandwidth usage and improves
// transfer speeds for text-based content (HTML, CSS, JS, JSON, etc.).
//
// Performance considerations:
//   - Compression trades CPU usage for bandwidth savings
//   - Most effective for text-based content (can reduce size by 60-80%)
//   - Less effective for already compressed content (images, videos)
//   - Consider compression level settings for CPU vs compression ratio balance
//
// Parameters:
//
//	config: The configuration map to update with compression settings
//	enabled: Boolean flag to enable (true) or disable (false) response compression
//
// Example:
//
//	helper.SetCompression(config, true)  // Enable compression for bandwidth savings
//	helper.SetCompression(config, false) // Disable for CPU-sensitive applications
func (c *ConfigHelper) SetCompression(config map[string]interface{}, enabled bool) {
	config[constants.COMPRESSION_ENABLED] = enabled
}

// GetCompression retrieves the current HTTP response compression configuration.
// This method returns whether automatic response compression is enabled.
//
// Returns:
//
//	bool: True if compression is enabled, false otherwise (default: false)
//
// Example:
//
//	if helper.GetCompression(config) {
//	    log.Println("Response compression is enabled")
//	}
func (c *ConfigHelper) GetCompression(config map[string]interface{}) bool {
	return c.GetConfigBool(config, constants.COMPRESSION_ENABLED, constants.DEFAULT_COMPRESSION)
}

// SetConcurrency configures the maximum number of concurrent connections.
// This setting controls how many simultaneous client connections the server
// will accept, helping to manage resource usage and prevent overload.
//
// Performance implications:
//   - Higher limits allow more simultaneous users but consume more memory
//   - Lower limits reduce resource usage but may reject connections under load
//   - Consider your server's RAM and CPU capacity when setting this value
//   - Each connection typically uses 2-8KB of memory plus application overhead
//
// Parameters:
//
//	config: The configuration map to update with concurrency settings
//	limit: Maximum number of concurrent connections (recommended: 1000-10000)
//
// Example:
//
//	helper.SetConcurrency(config, 5000)  // Allow up to 5000 concurrent connections
//	helper.SetConcurrency(config, 1000)  // Conservative limit for resource-constrained servers
func (c *ConfigHelper) SetConcurrency(config map[string]interface{}, limit int) {
	config[constants.MAX_CONCURRENT_CONNECTIONS] = limit
}

// GetConcurrency retrieves the maximum concurrent connections limit.
// Returns the configured limit for simultaneous client connections.
//
// Returns:
//
//	int: Maximum concurrent connections limit (default: 10000)
//
// Example:
//
//	maxConns := helper.GetConcurrency(config)
//	log.Printf("Server configured for max %d concurrent connections", maxConns)
func (c *ConfigHelper) GetConcurrency(config map[string]interface{}) int {
	return c.GetConfigInt(config, constants.MAX_CONCURRENT_CONNECTIONS, constants.DEFAULT_MAX_CONCURRENT_CONNS)
}

// SetCaching configures response caching settings for improved performance.
// Enables caching with specified configuration including provider type, TTL,
// and storage limits. Caching can dramatically improve response times and
// reduce server load for frequently accessed content.
//
// Supported cache providers:
//   - "memory": In-memory caching (fastest, but limited by RAM)
//   - "redis": Redis-based caching (shared across instances)
//   - "file": File-system based caching (persistent, slower)
//
// Cache configuration keys:
//   - "provider": Cache storage backend ("memory", "redis", "file")
//   - "ttl": Time-to-live duration for cached items
//   - "max_size": Maximum cache size in bytes
//   - "redis_url": Redis connection URL (if using Redis provider)
//   - "cache_dir": Directory path for file-based cache
//
// Parameters:
//
//	config: The configuration map to update with caching settings
//	cacheConfig: Map containing cache provider and settings configuration
//
// Example:
//
//	cacheSettings := map[string]interface{}{
//	    "provider": "memory",
//	    "ttl": 5 * time.Minute,
//	    "max_size": 100 * 1024 * 1024, // 100MB
//	}
//	helper.SetCaching(config, cacheSettings)
func (c *ConfigHelper) SetCaching(config map[string]interface{}, cacheConfig map[string]interface{}) {
	config[constants.CACHE_CONFIG] = cacheConfig
	config[constants.CACHE_ENABLED] = true
}

// GetCaching retrieves the current caching configuration settings.
// Returns both the cache configuration map and whether caching is enabled.
//
// Returns:
//
//	cacheConfig: Map containing cache provider and settings, nil if not configured
//	enabled: Boolean indicating if caching is enabled (true if SetCaching was called)
//
// Example:
//
//	cacheConfig, enabled := helper.GetCaching(config)
//	if enabled {
//	    provider := cacheConfig["provider"]
//	    ttl := cacheConfig["ttl"]
//	    log.Printf("Caching enabled with provider: %v, TTL: %v", provider, ttl)
//	}
func (c *ConfigHelper) GetCaching(config map[string]interface{}) (cacheConfig map[string]interface{}, enabled bool) {
	if cfg, exists := config[constants.CACHE_CONFIG]; exists {
		if cacheMap, ok := cfg.(map[string]interface{}); ok {
			cacheConfig = cacheMap
		}
	}
	enabled = c.GetConfigBool(config, constants.CACHE_ENABLED, false)
	return
}

// Monitoring & Observability Configuration Methods
// These methods manage monitoring, health checks, and observability features
// that are essential for production deployments and operational visibility.

// SetHealthCheck configures a health check endpoint for monitoring.
// Health checks are essential for load balancers, container orchestrators,
// and monitoring systems to determine if the server is operational.
//
// Common use cases:
//   - Load balancer health probes
//   - Kubernetes liveness and readiness probes
//   - Monitoring system uptime checks
//   - Service discovery health verification
//
// Parameters:
//
//	config: The configuration map to update with health check settings
//	path: URL path for the health check endpoint (e.g., "/health", "/status")
//	handler: Handler function or interface{} that processes health check requests
//
// Example:
//
//	healthHandler := func() map[string]interface{} {
//	    return map[string]interface{}{
//	        "status": "healthy",
//	        "timestamp": time.Now(),
//	        "version": "1.0.0",
//	    }
//	}
//	helper.SetHealthCheck(config, "/health", healthHandler)
func (c *ConfigHelper) SetHealthCheck(config map[string]interface{}, path string, handler interface{}) {
	config[constants.HEALTH_CHECK_PATH] = path
	config[constants.HEALTH_CHECK_HANDLER] = handler
	config[constants.HEALTH_CHECK_ENABLED] = true
}

// GetHealthCheck retrieves health check endpoint configuration.
// Returns the configured path, handler, and whether health checks are enabled.
//
// Returns:
//
//	path: URL path for health check endpoint (default: "/health")
//	handler: The configured health check handler, nil if not set
//	enabled: Boolean indicating if health checks are configured
//
// Example:
//
//	path, handler, enabled := helper.GetHealthCheck(config)
//	if enabled {
//	    log.Printf("Health check enabled at %s", path)
//	}
func (c *ConfigHelper) GetHealthCheck(config map[string]interface{}) (path string, handler interface{}, enabled bool) {
	path = c.GetConfigString(config, constants.HEALTH_CHECK_PATH, constants.DEFAULT_HEALTH_CHECK_PATH)
	handler = config[constants.HEALTH_CHECK_HANDLER]
	enabled = c.GetConfigBool(config, constants.HEALTH_CHECK_ENABLED, false)
	return
}

// SetTracing configures distributed tracing capabilities.
// Distributed tracing helps track requests across multiple services,
// providing visibility into request flows, performance bottlenecks,
// and error propagation in microservice architectures.
//
// Tracing benefits:
//   - Request flow visualization across services
//   - Performance bottleneck identification
//   - Error tracking and root cause analysis
//   - Service dependency mapping
//   - Latency analysis and optimization
//
// Common tracing systems:
//   - Jaeger, Zipkin, AWS X-Ray, Google Cloud Trace
//
// Parameters:
//
//	config: The configuration map to update with tracing settings
//	enabled: Boolean flag to enable (true) or disable (false) distributed tracing
//
// Example:
//
//	helper.SetTracing(config, true)  // Enable for production observability
//	helper.SetTracing(config, false) // Disable for development/testing
func (c *ConfigHelper) SetTracing(config map[string]interface{}, enabled bool) {
	config[constants.TRACING_ENABLED] = enabled
}

// GetTracing retrieves the current distributed tracing configuration.
// Returns whether distributed tracing is enabled for the webserver.
//
// Returns:
//
//	bool: True if tracing is enabled, false otherwise (default: false)
//
// Example:
//
//	if helper.GetTracing(config) {
//	    // Initialize tracing middleware
//	    log.Println("Distributed tracing is enabled")
//	}
func (c *ConfigHelper) GetTracing(config map[string]interface{}) bool {
	return c.GetConfigBool(config, constants.TRACING_ENABLED, false)
}

// Static Content & Asset Management Configuration Methods
// These methods handle configuration for serving static assets, public files,
// and template rendering - essential for web applications with frontend assets.

// SetStaticFiles configures static file serving with URL prefix mapping.
// Maps URL prefixes to filesystem directories for efficient static asset serving.
// This method supports multiple static file routes with different prefixes.
//
// Use cases:
//   - CSS, JavaScript, and image assets
//   - Font files and icons
//   - PDF documents and downloads
//   - API documentation and static content
//
// Performance considerations:
//   - Static files are typically cached by browsers and CDNs
//   - Consider using a reverse proxy (nginx) for better static file performance
//   - Enable compression for text-based static assets
//
// Parameters:
//
//	config: The configuration map to update with static file settings
//	prefix: URL prefix for static files (e.g., "/assets", "/css", "/images")
//	directory: Filesystem directory containing the static files
//
// Example:
//
//	helper.SetStaticFiles(config, "/assets", "./public/assets")
//	helper.SetStaticFiles(config, "/css", "./public/stylesheets")
//	// Multiple calls add additional static file routes
func (c *ConfigHelper) SetStaticFiles(config map[string]interface{}, prefix, directory string) {
	// Initialize static files config map if it doesn't exist
	if config[constants.STATIC_FILES] == nil {
		config[constants.STATIC_FILES] = make(map[string]string)
	}
	staticFiles := config[constants.STATIC_FILES].(map[string]string)
	staticFiles[prefix] = directory
	config[constants.STATIC_FILES_ENABLED] = true
}

// GetStaticFiles retrieves all configured static file route mappings.
// Returns the complete map of URL prefixes to directory paths.
//
// Returns:
//
//	staticFiles: Map of URL prefixes to filesystem directories, nil if none configured
//	enabled: Boolean indicating if static file serving is enabled
//
// Example:
//
//	files, enabled := helper.GetStaticFiles(config)
//	if enabled {
//	    for prefix, directory := range files {
//	        log.Printf("Static files: %s -> %s", prefix, directory)
//	    }
//	}
func (c *ConfigHelper) GetStaticFiles(config map[string]interface{}) (staticFiles map[string]string, enabled bool) {
	if files, exists := config[constants.STATIC_FILES]; exists {
		if fileMap, ok := files.(map[string]string); ok {
			staticFiles = fileMap
		}
	}
	enabled = c.GetConfigBool(config, constants.STATIC_FILES_ENABLED, false)
	return
}

// SetPublicDirectory configures a public directory for direct asset serving.
// This is a convenience method for serving all files in a directory at the root URL.
// Ideal for simple applications with a single public assets directory.
//
// Security considerations:
//   - Ensure the directory only contains files safe for public access
//   - Avoid placing sensitive configuration files in the public directory
//   - Consider file permissions and directory traversal protection
//
// Parameters:
//
//	config: The configuration map to update with public directory settings
//	directory: Filesystem path to the public directory (e.g., "./public", "./www")
//
// Example:
//
//	helper.SetPublicDirectory(config, "./public")
//	// Files in ./public/ are served at the root URL
//	// ./public/index.html becomes available at /index.html
func (c *ConfigHelper) SetPublicDirectory(config map[string]interface{}, directory string) {
	config[constants.PUBLIC_DIRECTORY] = directory
	config[constants.PUBLIC_DIRECTORY_ENABLED] = true
}

// GetPublicDirectory retrieves the configured public directory path.
// Returns the directory path and whether public directory serving is enabled.
//
// Returns:
//
//	directory: Filesystem path to the public directory, empty if not configured
//	enabled: Boolean indicating if public directory serving is enabled
//
// Example:
//
//	dir, enabled := helper.GetPublicDirectory(config)
//	if enabled {
//	    log.Printf("Public directory configured: %s", dir)
//	}
func (c *ConfigHelper) GetPublicDirectory(config map[string]interface{}) (directory string, enabled bool) {
	directory = c.GetConfigString(config, constants.PUBLIC_DIRECTORY, "")
	enabled = c.GetConfigBool(config, constants.PUBLIC_DIRECTORY_ENABLED, false)
	return
}

// SetFileServer configures multiple file server routes simultaneously.
// This method allows bulk configuration of multiple static file routes,
// making it convenient to set up complex static asset serving schemes.
//
// Use cases:
//   - Organizing assets by type (CSS, JS, images, fonts)
//   - Serving different content types from different directories
//   - Multi-tenant applications with tenant-specific asset directories
//
// Parameters:
//
//	config: The configuration map to update with file server settings
//	routes: Map of URL prefixes to filesystem directory paths
//
// Example:
//
//	routes := map[string]string{
//	    "/css":     "./assets/stylesheets",
//	    "/js":      "./assets/javascript",
//	    "/images":  "./assets/images",
//	    "/uploads": "./storage/uploads",
//	}
//	helper.SetFileServer(config, routes)
func (c *ConfigHelper) SetFileServer(config map[string]interface{}, routes map[string]string) {
	config[constants.FILE_SERVER_ROUTES] = routes
	config[constants.FILE_SERVER_ENABLED] = true
}

// GetFileServer retrieves all configured file server route mappings.
// Returns the complete routing configuration for file serving.
//
// Returns:
//
//	routes: Map of URL prefixes to filesystem directories, nil if not configured
//	enabled: Boolean indicating if file server routes are enabled
//
// Example:
//
//	routes, enabled := helper.GetFileServer(config)
//	if enabled {
//	    for prefix, directory := range routes {
//	        log.Printf("File server route: %s -> %s", prefix, directory)
//	    }
//	}
func (c *ConfigHelper) GetFileServer(config map[string]interface{}) (routes map[string]string, enabled bool) {
	if fileRoutes, exists := config[constants.FILE_SERVER_ROUTES]; exists {
		if routeMap, ok := fileRoutes.(map[string]string); ok {
			routes = routeMap
		}
	}
	enabled = c.GetConfigBool(config, constants.FILE_SERVER_ENABLED, false)
	return
}

// SetTemplate configures server-side template rendering capabilities.
// Enables template processing with specified engine and template directory.
// Supports various template engines for dynamic content generation.
//
// Supported template engines:
//   - HTML: Plain HTML (no templating)
//   - JSX: React JSX for server-side rendering
//   - Vue: Vue.js templates
//   - Handlebars: Mustache-based templating
//   - Pug: Indentation-based templates (formerly Jade)
//   - EJS: Embedded JavaScript templates
//   - Go: Go's built-in html/template
//
// Parameters:
//
//	config: The configuration map to update with template settings
//	engine: Template engine enum (from enums.TemplateEngine)
//	engineStr: String representation of the engine name
//	directory: Filesystem directory containing template files
//
// Example:
//
//	helper.SetTemplate(config, enums.JSX, "jsx", "./templates")
//	helper.SetTemplate(config, enums.Handlebars, "handlebars", "./views")
func (c *ConfigHelper) SetTemplate(config map[string]interface{}, engine enums.TemplateEngine, engineStr, directory string) {
	config[constants.TEMPLATE_ENGINE] = engine
	config[constants.TEMPLATE_ENGINE_STRING] = engineStr
	config[constants.TEMPLATE_DIRECTORY] = directory
	config[constants.TEMPLATE_ENABLED] = true
}

// GetTemplate retrieves the complete template engine configuration.
// Returns the engine enum, string representation, directory, and enabled status.
//
// Returns:
//
//	engine: Template engine enum value (enums.TemplateEngine)
//	engineStr: String representation of the engine name
//	directory: Filesystem directory containing template files
//	enabled: Boolean indicating if template rendering is enabled
//
// Example:
//
//	engine, engineStr, directory, enabled := helper.GetTemplate(config)
//	if enabled {
//	    log.Printf("Template engine: %s (%s) in directory: %s",
//	        engine.DisplayName(), engineStr, directory)
//	}
func (c *ConfigHelper) GetTemplate(config map[string]interface{}) (engine enums.TemplateEngine, engineStr, directory string, enabled bool) {
	if eng, exists := config[constants.TEMPLATE_ENGINE]; exists {
		if templateEngine, ok := eng.(enums.TemplateEngine); ok {
			engine = templateEngine
		}
	}
	engineStr = c.GetConfigString(config, constants.TEMPLATE_ENGINE_STRING, "")
	directory = c.GetConfigString(config, constants.TEMPLATE_DIRECTORY, "")
	enabled = c.GetConfigBool(config, constants.TEMPLATE_ENABLED, false)
	return
}
