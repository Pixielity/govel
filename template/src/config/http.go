package config

// Http returns the HTTP server configuration map.
// This configuration handles HTTP server settings, middleware,
// timeout configurations, and request/response limits.
func Http() map[string]any {
	return map[string]any{

		// HTTP Server Configuration
		//
		//
		// These settings control the basic HTTP server configuration including
		// host, port, and protocol settings for your web application.
		//
		"host": Env("HTTP_HOST", "127.0.0.1"),
		"port": Env("HTTP_PORT", 8080),

		// TLS/HTTPS Configuration
		//
		//
		// Configuration for HTTPS including certificate files, TLS versions,
		// and SSL-specific settings for secure HTTP communication.
		//
		"tls": map[string]any{
			"enabled":     Env("HTTP_TLS_ENABLED", false),
			"cert_file":   Env("HTTP_TLS_CERT_FILE", ""),
			"key_file":    Env("HTTP_TLS_KEY_FILE", ""),
			"min_version": Env("HTTP_TLS_MIN_VERSION", "1.2"),
			"max_version": Env("HTTP_TLS_MAX_VERSION", "1.3"),
		},

		// Timeouts
		//
		//
		// HTTP server timeout settings to prevent hanging connections and
		// ensure responsive behavior under load.
		//
		"timeouts": map[string]any{
			"read":       Env("HTTP_READ_TIMEOUT", 30),      // seconds
			"write":      Env("HTTP_WRITE_TIMEOUT", 30),     // seconds
			"idle":       Env("HTTP_IDLE_TIMEOUT", 60),      // seconds
			"shutdown":   Env("HTTP_SHUTDOWN_TIMEOUT", 30),  // seconds
			"keep_alive": Env("HTTP_KEEP_ALIVE_TIMEOUT", 3), // minutes
		},

		// Request Limits
		//
		//
		// Settings to control request size limits, header limits, and other
		// constraints on incoming HTTP requests.
		//
		"limits": map[string]any{
			"max_request_size":   Env("HTTP_MAX_REQUEST_SIZE", 32<<20),    // 32MB
			"max_header_size":    Env("HTTP_MAX_HEADER_SIZE", 1<<20),      // 1MB
			"max_form_size":      Env("HTTP_MAX_FORM_SIZE", 32<<20),       // 32MB
			"max_multipart_size": Env("HTTP_MAX_MULTIPART_SIZE", 128<<20), // 128MB
		},

		// Middleware Configuration
		//
		//
		// Global middleware that will be applied to all HTTP requests.
		// These middleware components handle cross-cutting concerns.
		//
		"middleware": map[string]any{
			"global": []string{
				"recovery",
				"logger",
				"cors",
				"secure_headers",
				"rate_limit",
			},
			"api": []string{
				"throttle",
				"auth:api",
				"bindings",
			},
			"web": []string{
				"web",
				"csrf",
				"auth",
			},
		},

		// Compression
		//
		//
		// HTTP response compression settings to reduce bandwidth usage
		// and improve response times.
		//
		"compression": map[string]any{
			"enabled": Env("HTTP_COMPRESSION_ENABLED", true),
			"level":   Env("HTTP_COMPRESSION_LEVEL", 6), // 1-9, 6 is balanced
			"types": []string{
				"text/html",
				"text/css",
				"text/javascript",
				"application/javascript",
				"application/json",
				"application/xml",
				"image/svg+xml",
			},
		},

		// Static File Serving
		//
		//
		// Configuration for serving static files including cache headers,
		// compression, and security settings.
		//
		"static": map[string]any{
			"enabled":       Env("HTTP_STATIC_ENABLED", true),
			"prefix":        Env("HTTP_STATIC_PREFIX", "/static"),
			"root":          Env("HTTP_STATIC_ROOT", "./public"),
			"index":         Env("HTTP_STATIC_INDEX", "index.html"),
			"cache_control": Env("HTTP_STATIC_CACHE_CONTROL", "public, max-age=31536000"), // 1 year
		},

		// Rate Limiting
		//
		//
		// HTTP rate limiting configuration to prevent abuse and ensure
		// fair usage of your application resources.
		//
		"rate_limit": map[string]any{
			"enabled":             Env("HTTP_RATE_LIMIT_ENABLED", true),
			"requests_per_minute": Env("HTTP_RATE_LIMIT_RPM", 60),
			"burst":               Env("HTTP_RATE_LIMIT_BURST", 10),
			"cleanup_interval":    Env("HTTP_RATE_LIMIT_CLEANUP", 300), // seconds
		},

		// Health Check
		//
		//
		// HTTP health check endpoint configuration for load balancers
		// and monitoring systems.
		//
		"health_check": map[string]any{
			"enabled": Env("HTTP_HEALTH_CHECK_ENABLED", true),
			"path":    Env("HTTP_HEALTH_CHECK_PATH", "/health"),
		},
	}
}
