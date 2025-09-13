package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"govel/packages/new/webserver/src/builders"
	"govel/packages/new/webserver/src/enums"
)

// Enhanced builder example demonstrating new security, performance, and static file methods
func main() {
	fmt.Println("Starting Enhanced Webserver Builder Example...")

	// Create webserver with comprehensive configuration using the new builder methods
	server := builders.Configure().
		// Engine selection
		WithEngine(enums.GoFiber).

		// Network configuration
		WithHost("0.0.0.0").
		WithPort(8085).

		// Security & Authentication
		WithTLS("/path/to/cert.pem", "/path/to/key.pem").
		WithJWTSecret("your-super-secret-jwt-signing-key").
		WithCORS(map[string]interface{}{
			"allow_origins":     []string{"https://example.com", "https://app.com"},
			"allow_methods":     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			"allow_headers":     []string{"Content-Type", "Authorization", "X-API-Key"},
			"allow_credentials": true,
			"max_age":          86400, // 24 hours
		}).
		WithSecurityHeaders(map[string]string{
			"X-Frame-Options":           "DENY",
			"X-XSS-Protection":          "1; mode=block",
			"X-Content-Type-Options":    "nosniff",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
			"Content-Security-Policy":   "default-src 'self'",
		}).
		WithRateLimit(100, time.Minute). // 100 requests per minute
		WithBasicAuth("admin", "secure-password").
		WithAPIKeys([]string{
			"api-key-1234567890",
			"api-key-abcdefghij",
			"api-key-0987654321",
		}).

		// Performance & Optimization
		WithTimeout(30 * time.Second).
		WithReadTimeout(10 * time.Second).
		WithWriteTimeout(10 * time.Second).
		WithIdleTimeout(60 * time.Second).
		WithKeepAlive(true).
		WithMaxBodySize(10 * 1024 * 1024). // 10MB
		WithCompression(true).
		WithConcurrency(1000).
		WithCaching(map[string]interface{}{
			"provider": "memory",
			"ttl":      5 * time.Minute,
			"max_size": 100 * 1024 * 1024, // 100MB
		}).

		// Monitoring
		WithHealthCheck("/health", nil). // Handler would be provided in real implementation
		WithTracing(true).

		// Static Content & Templates
		WithStaticFiles("/css", "./assets/css").
		WithStaticFiles("/js", "./assets/js").
		WithStaticFiles("/images", "./assets/images").
		WithPublicDirectory("./public").
		WithFileServer(map[string]string{
			"/uploads":   "./storage/uploads",
			"/downloads": "./storage/downloads",
			"/docs":      "./documentation",
		}).
		WithTemplate(enums.JSX, "./templates").

		// Additional configuration
		Set("custom_setting", "custom_value").
		Set("api_version", "v1.0.0").
		Set("environment", "development").

		// Build the server
		Build()

	// Demonstrate that the configuration was applied
	fmt.Println("Enhanced webserver built successfully with configuration:")
	fmt.Printf("  Engine: %v\n", server.GetConfig("engine"))
	fmt.Printf("  Host: %v\n", server.GetConfig("host"))
	fmt.Printf("  Port: %v\n", server.GetConfig("port"))
	fmt.Printf("  TLS Enabled: %v\n", server.GetConfig("tls_enabled"))
	fmt.Printf("  JWT Secret Set: %v\n", server.GetConfig("jwt_secret") != nil)
	fmt.Printf("  CORS Enabled: %v\n", server.GetConfig("cors_enabled"))
	fmt.Printf("  Security Headers Set: %v\n", server.GetConfig("security_headers") != nil)
	fmt.Printf("  Rate Limiting: %v requests per %v\n", 
		server.GetConfig("rate_limit_requests"), 
		server.GetConfig("rate_limit_window"))
	fmt.Printf("  Basic Auth Enabled: %v\n", server.GetConfig("basic_auth_enabled"))
	fmt.Printf("  API Keys Count: %d\n", len(server.GetConfig("api_keys").([]string)))
	fmt.Printf("  Request Timeout: %v\n", server.GetConfig("request_timeout"))
	fmt.Printf("  Max Body Size: %v bytes\n", server.GetConfig("max_body_size"))
	fmt.Printf("  Compression Enabled: %v\n", server.GetConfig("compression_enabled"))
	fmt.Printf("  Max Concurrent Connections: %v\n", server.GetConfig("max_concurrent_connections"))
	fmt.Printf("  Caching Enabled: %v\n", server.GetConfig("cache_enabled"))
	fmt.Printf("  Health Check Path: %v\n", server.GetConfig("health_check_path"))
	fmt.Printf("  Tracing Enabled: %v\n", server.GetConfig("tracing_enabled"))
	fmt.Printf("  Static Files Enabled: %v\n", server.GetConfig("static_files_enabled"))
	fmt.Printf("  Public Directory: %v\n", server.GetConfig("public_directory"))
	fmt.Printf("  File Server Enabled: %v\n", server.GetConfig("file_server_enabled"))
	fmt.Printf("  Template Engine: %v\n", server.GetConfig("template_engine_string"))
	fmt.Printf("  Template Directory: %v\n", server.GetConfig("template_directory"))
	fmt.Printf("  Custom Setting: %v\n", server.GetConfig("custom_setting"))

	fmt.Println("\nConfiguration features demonstrated:")
	fmt.Println("✓ TLS/SSL certificate configuration")
	fmt.Println("✓ JWT secret for authentication")
	fmt.Println("✓ Comprehensive CORS settings")
	fmt.Println("✓ Security headers protection")
	fmt.Println("✓ Rate limiting configuration")
	fmt.Println("✓ HTTP Basic Authentication")
	fmt.Println("✓ API key authentication")
	fmt.Println("✓ Request/Read/Write/Idle timeouts")
	fmt.Println("✓ HTTP Keep-alive settings")
	fmt.Println("✓ Request body size limits")
	fmt.Println("✓ Response compression")
	fmt.Println("✓ Concurrency limits")
	fmt.Println("✓ Caching configuration")
	fmt.Println("✓ Health check endpoints")
	fmt.Println("✓ Distributed tracing")
	fmt.Println("✓ Static file serving")
	fmt.Println("✓ Public directory serving")
	fmt.Println("✓ Multi-route file servers")
	fmt.Printf("✓ Template engine (%s) configuration\n", enums.JSX.DisplayName())

	fmt.Println("\nTemplate Engine Examples:")
	fmt.Printf("  HTML: %s (%s)\n", enums.HTML.DisplayName(), enums.HTML.FileExtension())
	fmt.Printf("  JSX: %s (%s)\n", enums.JSX.DisplayName(), enums.JSX.FileExtension())
	fmt.Printf("  Vue: %s (%s)\n", enums.Vue.DisplayName(), enums.Vue.FileExtension())
	fmt.Printf("  Handlebars: %s (%s)\n", enums.Handlebars.DisplayName(), enums.Handlebars.FileExtension())
	fmt.Printf("  Pug: %s (%s)\n", enums.Pug.DisplayName(), enums.Pug.FileExtension())
	fmt.Printf("  Go Template: %s (%s)\n", enums.GoTemplate.DisplayName(), enums.GoTemplate.FileExtension())

	fmt.Println("\nThis example demonstrates the fluent builder pattern with:")
	fmt.Println("• Method chaining for clean, readable configuration")
	fmt.Println("• Type-safe enum usage for template engines")
	fmt.Println("• Comprehensive validation of all settings")
	fmt.Println("• Enterprise-ready security and performance options")
	fmt.Println("• Static asset and template rendering capabilities")

	// Note: In a real application, you would add routes and call server.Listen()
	fmt.Println("\nBuilder pattern completed successfully!")
	fmt.Println("Server is configured and ready for route registration and startup.")

	// Graceful shutdown example
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("\nInitiating graceful shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	fmt.Println("Enhanced builder example completed.")
}
