package main

import (
	"fmt"
	"log"

	"govel/new/pipeline/src"
)

// HTTPRequest represents an HTTP request for the pipeline
type HTTPRequest struct {
	Method   string
	Path     string
	Headers  map[string]string
	Body     string
	UserID   int
	Metadata map[string]interface{}
}

// AuthMiddleware checks authentication
type AuthMiddleware struct{}

func (a *AuthMiddleware) Handle(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	req := passable.(*HTTPRequest)
	
	fmt.Printf("[Auth] Checking authentication for %s %s\n", req.Method, req.Path)
	
	// Simulate authentication check
	if authHeader, exists := req.Headers["Authorization"]; exists && authHeader == "Bearer valid-token" {
		req.UserID = 12345 // Set authenticated user ID
		fmt.Println("[Auth] Authentication successful")
	} else {
		return nil, fmt.Errorf("authentication failed: missing or invalid token")
	}
	
	return next(passable)
}

// LoggingMiddleware logs the request
type LoggingMiddleware struct{}

func (l *LoggingMiddleware) Handle(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	req := passable.(*HTTPRequest)
	
	fmt.Printf("[Logging] Processing %s %s for user %d\n", req.Method, req.Path, req.UserID)
	
	// Call next middleware
	result, err := next(passable)
	
	if err != nil {
		fmt.Printf("[Logging] Request failed: %v\n", err)
	} else {
		fmt.Printf("[Logging] Request completed successfully\n")
	}
	
	return result, err
}

// RateLimitMiddleware checks rate limits
func RateLimitMiddleware(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	req := passable.(*HTTPRequest)
	
	fmt.Printf("[RateLimit] Checking rate limits for user %d\n", req.UserID)
	
	// Simulate rate limit check (always pass in this example)
	fmt.Println("[RateLimit] Rate limit check passed")
	
	return next(passable)
}

// ExampleHTTPMiddlewarePipeline demonstrates using pipelines for HTTP middleware processing
func ExampleHTTPMiddlewarePipeline() {
	fmt.Println("=== HTTP Middleware Pipeline Example ===")
	
	// Create pipeline
	p := pipeline.NewPipeline(nil)
	
	// Create middleware instances
	auth := &AuthMiddleware{}
	logger := &LoggingMiddleware{}
	
	// Create test request
	request := &HTTPRequest{
		Method: "GET",
		Path:   "/api/users/profile",
		Headers: map[string]string{
			"Authorization": "Bearer valid-token",
			"Content-Type":  "application/json",
		},
		Body:     "",
		Metadata: make(map[string]interface{}),
	}
	
	// Process through middleware pipeline
	result, err := p.
		Send(request).
		Through([]interface{}{auth, logger, RateLimitMiddleware}).
		Then(func(passable interface{}) interface{} {
			req := passable.(*HTTPRequest)
			fmt.Printf("[Handler] Processing final request for user %d\n", req.UserID)
			
			// Simulate API response
			return map[string]interface{}{
				"user_id": req.UserID,
				"profile": map[string]interface{}{
					"name":  "John Doe",
					"email": "john@example.com",
				},
				"status": "success",
			}
		})
	
	if err != nil {
		log.Printf("Pipeline failed: %v", err)
		return
	}
	
	fmt.Printf("\nFinal result: %+v\n", result)
	
	fmt.Println("\n=== Trying with invalid token ===")
	
	// Test with invalid token
	invalidRequest := &HTTPRequest{
		Method: "GET",
		Path:   "/api/users/profile",
		Headers: map[string]string{
			"Authorization": "Bearer invalid-token",
		},
	}
	
	_, err = p.
		Send(invalidRequest).
		Through([]interface{}{auth, logger, RateLimitMiddleware}).
		ThenReturn()
	
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
	
	// Output:
	// === HTTP Middleware Pipeline Example ===
	// [Auth] Checking authentication for GET /api/users/profile
	// [Auth] Authentication successful
	// [Logging] Processing GET /api/users/profile for user 12345
	// [RateLimit] Checking rate limits for user 12345
	// [RateLimit] Rate limit check passed
	// [Handler] Processing final request for user 12345
	// [Logging] Request completed successfully
	// 
	// Final result: map[profile:map[email:john@example.com name:John Doe] status:success user_id:12345]
	// 
	// === Trying with invalid token ===
	// [Auth] Checking authentication for GET /api/users/profile
	// Expected error: authentication failed: missing or invalid token
}

func main() {
	ExampleHTTPMiddlewarePipeline()
}
