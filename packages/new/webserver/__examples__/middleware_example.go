package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"govel/new/webserver"
	"govel/new/webserver/builders"
	"govel/new/webserver/enums"
	"govel/new/webserver/interfaces"
	"govel/new/webserver/logging"
	"govel/new/webserver/types"
)

// Comprehensive middleware usage example
// This example demonstrates various middleware patterns and chaining
func main() {
	fmt.Println("Starting Middleware Example Server...")

	// Create custom middleware implementations
	loggingMiddleware := NewLoggingMiddleware()
	corsMiddleware := NewCORSMiddleware()
	authMiddleware := NewAuthMiddleware()
	metricsMiddleware := NewMetricsMiddleware()
	rateLimit := NewRateLimitMiddleware(10) // 10 requests per minute
	errorHandler := NewErrorHandlerMiddleware()

	// Create webserver with global middleware
	server := builders.NewWebserverBuilder().
		WithEngine(enums.GoFiber).
		WithPort(8084).
		WithHost("localhost").
		WithMiddleware(
			loggingMiddleware,    // Log all requests
			errorHandler,        // Handle errors gracefully
			metricsMiddleware,    // Collect metrics
			corsMiddleware,       // Enable CORS
		).
		Build()

	// Public endpoints (no authentication required)
	server.Get("/", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"message": "Welcome to Middleware Example Server!",
			"middleware_applied": []string{
				"Logging", "Error Handler", "Metrics", "CORS",
			},
			"public": true,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}))

	// Health check endpoint
	server.Get("/health", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"status": "healthy",
			"middleware_status": "operational",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}))

	// Middleware information endpoint
	server.Get("/middleware/info", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"global_middleware": []map[string]interface{}{
				{"name": "Logging", "description": "Logs all incoming requests", "priority": 1},
				{"name": "Error Handler", "description": "Catches and formats errors", "priority": 2},
				{"name": "Metrics", "description": "Collects performance metrics", "priority": 3},
				{"name": "CORS", "description": "Handles cross-origin requests", "priority": 4},
			},
			"route_specific": []string{
				"Authentication", "Rate Limiting", "Caching", "Validation",
			},
			"middleware_patterns": []string{
				"Chain of Responsibility", "Decorator", "Pipeline", "Filter",
			},
		})
	}))

	// Protected routes with authentication middleware
	server.Group("/protected", func(protected interfaces.WebserverInterface) {
		// Apply authentication middleware to this group
		// In a real implementation, this would be done via route-specific middleware
		
		protected.Get("/profile", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			// Simulate authentication check
			token := req.Header("Authorization", "")
			if !isValidToken(token) {
				return webserver.NewResponse().
					Status(401).
					Json(map[string]interface{}{
						"error": "Unauthorized",
						"message": "Valid authentication token required",
						"middleware": "Authentication",
					})
			}

			return webserver.NewResponse().Json(map[string]interface{}{
				"user": map[string]interface{}{
					"id": "user123",
					"name": "John Doe",
					"email": "john@example.com",
					"role": "user",
				},
				"middleware_applied": []string{
					"Logging", "Error Handler", "Metrics", "CORS", "Authentication",
				},
				"access_time": time.Now().Format(time.RFC3339),
			})
		}))

		protected.Post("/data", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			// Check authentication and rate limiting
			token := req.Header("Authorization", "")
			if !isValidToken(token) {
				return webserver.NewResponse().
					Status(401).
					Json(map[string]string{"error": "Unauthorized"})
			}

			// Simulate rate limiting check
			clientIP := req.IP()
			if isRateLimited(clientIP) {
				return webserver.NewResponse().
					Status(429).
					Json(map[string]interface{}{
						"error": "Rate limit exceeded",
						"message": "Too many requests",
						"middleware": "Rate Limiter",
						"retry_after": "60 seconds",
					})
			}

			// Process the request
			data := req.Input("data", "")
			if data == "" {
				return webserver.NewResponse().
					Status(400).
					Json(map[string]string{"error": "Data field is required"})
			}

			return webserver.NewResponse().
				Status(201).
				Json(map[string]interface{}{
					"message": "Data processed successfully",
					"processed_data": data,
					"middleware_chain": []string{
						"Logging", "Error Handler", "Metrics", "CORS",
						"Authentication", "Rate Limiter", "Validation",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				})
		}))
	})

	// Admin routes with additional middleware
	server.Group("/admin", func(admin interfaces.WebserverInterface) {
		admin.Get("/users", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			// Check for admin token
			token := req.Header("Authorization", "")
			if !isAdminToken(token) {
				return webserver.NewResponse().
					Status(403).
					Json(map[string]interface{}{
						"error": "Forbidden",
						"message": "Admin access required",
						"middleware": "Admin Authorization",
					})
			}

			return webserver.NewResponse().Json(map[string]interface{}{
				"users": []map[string]interface{}{
					{"id": 1, "name": "Alice", "role": "admin"},
					{"id": 2, "name": "Bob", "role": "user"},
					{"id": 3, "name": "Charlie", "role": "user"},
				},
				"admin_access": true,
				"middleware_stack": []string{
					"Logging", "Error Handler", "Metrics", "CORS", "Admin Auth",
				},
			})
		}))

		admin.Post("/system/restart", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			token := req.Header("Authorization", "")
			if !isAdminToken(token) {
				return webserver.NewResponse().
					Status(403).
					Json(map[string]string{"error": "Admin access required"})
			}

			return webserver.NewResponse().Json(map[string]interface{}{
				"message": "System restart initiated",
				"admin_action": true,
				"initiated_by": "admin",
				"scheduled_for": time.Now().Add(5 * time.Minute).Format(time.RFC3339),
				"security_middleware": "Admin authorization passed",
			})
		}))
	})

	// Error demonstration endpoints
	server.Get("/errors/panic", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		// This would normally cause a panic, but error middleware should catch it
		return webserver.NewResponse().
			Status(500).
			Json(map[string]interface{}{
				"error": "Simulated panic",
				"message": "This demonstrates error handling middleware",
				"middleware": "Error Handler caught this",
				"recovered": true,
			})
	}))

	server.Get("/errors/timeout", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		// Simulate a long-running request
		time.Sleep(100 * time.Millisecond) // Simulate work
		return webserver.NewResponse().
			Status(408).
			Json(map[string]interface{}{
				"error": "Request timeout simulation",
				"middleware": "Timeout handler would catch this in real scenario",
				"processing_time": "100ms",
			})
	}))

	// Metrics endpoint
	server.Get("/metrics", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().Json(map[string]interface{}{
			"middleware_metrics": map[string]interface{}{
				"total_requests": 42,
				"successful_requests": 38,
				"failed_requests": 4,
				"average_response_time": "150ms",
				"middleware_overhead": "5ms",
			},
			"middleware_performance": map[string]string{
				"logging": "2ms",
				"cors": "1ms",
				"auth": "10ms",
				"metrics": "1ms",
				"rate_limit": "1ms",
			},
			"collected_at": time.Now().Format(time.RFC3339),
		})
	}))

	// CORS preflight handling example
	server.Options("/*", types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
		return webserver.NewResponse().
			Status(200).
			Header("Access-Control-Allow-Origin", "*").
			Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS").
			Header("Access-Control-Allow-Headers", "Content-Type, Authorization").
			Json(map[string]interface{}{
				"cors_preflight": "handled",
				"middleware": "CORS",
			})
	}))

	// Start server with middleware information
	fmt.Println("Starting Middleware Example Server...")
	fmt.Println("")
	fmt.Println("Middleware Stack (Global):")
	fmt.Println("  1. Logging Middleware      - Logs all requests")
	fmt.Println("  2. Error Handler          - Catches and handles errors")
	fmt.Println("  3. Metrics Collector      - Collects performance metrics")
	fmt.Println("  4. CORS Handler           - Handles cross-origin requests")
	fmt.Println("")
	fmt.Println("Route-Specific Middleware:")
	fmt.Println("  • Authentication         - /protected/* routes")
	fmt.Println("  • Rate Limiting          - /protected/data")
	fmt.Println("  • Admin Authorization    - /admin/* routes")
	fmt.Println("")
	
	// Display all registered routes with clickable URLs
	logging.DisplayRoutesClickable(server, "localhost", 8084)
	
	fmt.Println("Test with:")
	fmt.Println("  curl http://localhost:8084/")
	fmt.Println("  curl -H \"Authorization: Bearer valid-token\" http://localhost:8084/protected/profile")
	fmt.Println("  curl -H \"Authorization: Bearer admin-token\" http://localhost:8084/admin/users")
	fmt.Println("")
	fmt.Println("Press Ctrl+C to stop the server")

	// Graceful shutdown
	go func() {
		time.Sleep(60 * time.Second) // Run longer for middleware testing
		fmt.Println("\nShutting down middleware server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Middleware server shutdown error: %v", err)
		}
	}()

	// Start listening
	if err := server.Listen(":8084"); err != nil {
		log.Printf("Middleware server failed to start: %v", err)
	}

	fmt.Println("Middleware server stopped.")
}

// Custom middleware implementations

// LoggingMiddleware logs all incoming requests
type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (m *LoggingMiddleware) Before(req interfaces.RequestInterface) error {
	fmt.Printf("[LOG] %s %s from %s\n", req.Method(), req.Path(), req.IP())
	return nil
}

func (m *LoggingMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	start := time.Now()
	resp := next.Handle(req)
	duration := time.Since(start)
	fmt.Printf("[LOG] %s %s -> %d (%v)\n", req.Method(), req.Path(), resp.StatusCode(), duration)
	return resp
}

func (m *LoggingMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *LoggingMiddleware) Priority() int {
	return 1 // High priority (execute first)
}

// CORSMiddleware handles cross-origin requests
type CORSMiddleware struct{}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (m *CORSMiddleware) Before(req interfaces.RequestInterface) error {
	return nil
}

func (m *CORSMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	resp := next.Handle(req)
	
	// Add CORS headers
	resp.Header("Access-Control-Allow-Origin", "*")
	resp.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	resp.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	return resp
}

func (m *CORSMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *CORSMiddleware) Priority() int {
	return 4
}

// AuthMiddleware handles authentication
type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Before(req interfaces.RequestInterface) error {
	return nil
}

func (m *AuthMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	// Check if this is a protected route
	if strings.HasPrefix(req.Path(), "/protected") {
		token := req.Header("Authorization", "")
		if !isValidToken(token) {
			return webserver.NewResponse().
				Status(401).
				Json(map[string]string{"error": "Unauthorized"})
		}
	}
	
	return next.Handle(req)
}

func (m *AuthMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *AuthMiddleware) Priority() int {
	return 5
}

// MetricsMiddleware collects performance metrics
type MetricsMiddleware struct {
	requestCount int
	errorCount   int
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

func (m *MetricsMiddleware) Before(req interfaces.RequestInterface) error {
	m.requestCount++
	return nil
}

func (m *MetricsMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	resp := next.Handle(req)
	
	if resp.StatusCode() >= 400 {
		m.errorCount++
	}
	
	// Add metrics headers
	resp.Header("X-Request-Count", fmt.Sprintf("%d", m.requestCount))
	resp.Header("X-Error-Count", fmt.Sprintf("%d", m.errorCount))
	
	return resp
}

func (m *MetricsMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *MetricsMiddleware) Priority() int {
	return 3
}

// RateLimitMiddleware implements rate limiting
type RateLimitMiddleware struct {
	limit int
	requests map[string]int
}

func NewRateLimitMiddleware(limit int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limit: limit,
		requests: make(map[string]int),
	}
}

func (m *RateLimitMiddleware) Before(req interfaces.RequestInterface) error {
	return nil
}

func (m *RateLimitMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	clientIP := req.IP()
	
	if m.requests[clientIP] >= m.limit {
		return webserver.NewResponse().
			Status(429).
			Json(map[string]string{"error": "Rate limit exceeded"})
	}
	
	m.requests[clientIP]++
	return next.Handle(req)
}

func (m *RateLimitMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *RateLimitMiddleware) Priority() int {
	return 6
}

// ErrorHandlerMiddleware catches and handles errors
type ErrorHandlerMiddleware struct{}

func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}

func (m *ErrorHandlerMiddleware) Before(req interfaces.RequestInterface) error {
	return nil
}

func (m *ErrorHandlerMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	// In a real implementation, this would use recover() to catch panics
	resp := next.Handle(req)
	
	// Add error handling headers
	if resp.StatusCode() >= 400 {
		resp.Header("X-Error-Handled", "true")
	}
	
	return resp
}

func (m *ErrorHandlerMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

func (m *ErrorHandlerMiddleware) Priority() int {
	return 2 // High priority for error handling
}

// Helper functions

func isValidToken(token string) bool {
	// Simple token validation for demo
	return strings.HasPrefix(token, "Bearer ") && len(token) > 7
}

func isAdminToken(token string) bool {
	// Simple admin token validation for demo
	return token == "Bearer admin-token"
}

func isRateLimited(clientIP string) bool {
	// Simple rate limiting simulation
	return false // Allow all for demo
}
