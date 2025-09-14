package config

// Cors returns the CORS configuration map.
// This configuration handles Cross-Origin Resource Sharing settings.
func Cors() map[string]any {
	return map[string]any{

		// CORS Paths
		//
		// Here you may specify which paths should be subject to CORS handling.
		// You can use wildcards (*) to match multiple paths or specify exact
		// paths for fine-grained control over your CORS configuration.
		"paths": []string{"api/*", "sanctum/csrf-cookie"},

		// Allowed HTTP Methods
		//
		// This option controls which HTTP methods are allowed for cross-origin
		// requests. You may specify individual methods or use '*' to allow all
		// common methods. OPTIONS requests are handled automatically.
		"allowed_methods": []string{"*"},

		// Allowed Origins
		//
		// Here you may specify which origins are allowed to access your
		// application. Use '*' for all origins or specify exact domains.
		// For security, avoid wildcards in production environments.
		"allowed_origins": []string{"*"},

		// Allowed Origin Patterns
		//
		// If you need more flexible origin matching, specify regex patterns here.
		// This is useful for multiple subdomains or dynamic origins that follow
		// a predictable pattern. Leave empty if not using pattern matching.
		"allowed_origins_patterns": []string{},

		// Allowed Headers
		//
		// Specify which headers are allowed in cross-origin requests. Use '*'
		// to allow all headers or list specific headers like 'Content-Type',
		// 'Authorization'. Custom headers must be explicitly listed.
		"allowed_headers": []string{"*"},

		// Exposed Headers
		//
		// List response headers that should be accessible to client-side scripts.
		// By default, only basic headers are exposed. Add custom headers here
		// if your JavaScript needs to access them from API responses.
		"exposed_headers": []string{},

		// Preflight Max Age
		//
		// How long (in seconds) browsers should cache preflight request results.
		// Longer cache times reduce preflight requests and improve performance.
		// Set to 0 to disable caching or use 86400 for 24-hour caching.
		"max_age": 0,

		// Support Credentials
		//
		// Set to true to allow credentials (cookies, auth headers, certificates)
		// in cross-origin requests. When enabled, you cannot use wildcards for
		// origins and must specify exact domains for security.
		"supports_credentials": false,
	}
}
