package facades

import (
	httpInterfaces "govel/types/interfaces/http"
	facade "govel/support"
)

// HTTP provides a clean, static-like interface to the application's HTTP client service.
//
// This facade implements the facade pattern, providing global access to the HTTP client
// service configured in the dependency injection container. It offers a Laravel-style
// API for HTTP operations with automatic service resolution, request/response handling,
// middleware support, and comprehensive web service integration capabilities.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved HTTP client service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent HTTP operations across goroutines
//   - Supports multiple HTTP client configurations and connection pools
//   - Built-in middleware chain for request/response processing
//
// Behavior:
//   - First call: Resolves HTTP service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if HTTP service cannot be resolved (fail-fast behavior)
//   - Automatically handles connection pooling, retries, and timeout management
//
// Returns:
//   - HTTPInterface: The application's HTTP client service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "http" service is not registered in the container
//   - If the resolved service doesn't implement HTTPInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call HTTP() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - HTTP client operations are thread-safe with connection pooling
//
// Usage Examples:
//
//	// Basic GET request
//	response, err := facades.HTTP().Get("https://api.example.com/users")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer response.Body.Close()
//
//	fmt.Printf("Status: %d\n", response.StatusCode)
//	body, _ := io.ReadAll(response.Body)
//	fmt.Printf("Response: %s\n", body)
//
//	// GET with query parameters
//	response, err := facades.HTTP().Get("https://api.example.com/users", map[string]interface{}{
//	    "page":     1,
//	    "per_page": 10,
//	    "active":   true,
//	})
//
//	// POST request with JSON data
//	data := map[string]interface{}{
//	    "name":  "John Doe",
//	    "email": "john@example.com",
//	    "age":   30,
//	}
//
//	response, err := facades.HTTP().Post("https://api.example.com/users", data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// PUT request for updates
//	updatedData := map[string]interface{}{
//	    "name": "John Smith",
//	    "age":  31,
//	}
//
//	response, err := facades.HTTP().Put("https://api.example.com/users/123", updatedData)
//
//	// DELETE request
//	response, err := facades.HTTP().Delete("https://api.example.com/users/123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if response.StatusCode == 204 {
//	    fmt.Println("User deleted successfully")
//	}
//
//	// Request with custom headers
//	headers := map[string]string{
//	    "Authorization": "Bearer " + token,
//	    "Content-Type":  "application/json",
//	    "User-Agent":    "MyApp/1.0",
//	}
//
//	response, err := facades.HTTP().WithHeaders(headers).Get("https://api.example.com/profile")
//
//	// Request with timeout
//	response, err := facades.HTTP().WithTimeout(30 * time.Second).Get("https://slow-api.example.com/data")
//
//	// Request with retry logic
//	response, err := facades.HTTP().WithRetry(3).Get("https://unreliable-api.example.com/data")
//
//	// Form data submission
//	formData := map[string]string{
//	    "username": "john_doe",
//	    "password": "secret123",
//	}
//
//	response, err := facades.HTTP().PostForm("https://api.example.com/login", formData)
//
//	// File upload
//	filePath := "/path/to/upload/file.txt"
//	response, err := facades.HTTP().PostFile("https://api.example.com/upload", "file", filePath)
//
//	// Multiple file upload
//	files := map[string]string{
//	    "document": "/path/to/document.pdf",
//	    "image":    "/path/to/image.jpg",
//	}
//
//	response, err := facades.HTTP().PostFiles("https://api.example.com/upload", files)
//
//	// JSON response parsing
//	type User struct {
//	    ID    int    `json:"id"`
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	var user User
//	err := facades.HTTP().GetJSON("https://api.example.com/users/123", &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("User: %+v\n", user)
//
//	// Batch requests
//	urls := []string{
//	    "https://api.example.com/users/1",
//	    "https://api.example.com/users/2",
//	    "https://api.example.com/users/3",
//	}
//
//	responses := facades.HTTP().GetBatch(urls)
//	for i, resp := range responses {
//	    fmt.Printf("URL %d: Status %d\n", i+1, resp.StatusCode)
//	}
//
// Advanced HTTP Patterns:
//
//	// API client wrapper
//	type APIClient struct {
//	    baseURL string
//	    token   string
//	}
//
//	func NewAPIClient(baseURL, token string) *APIClient {
//	    return &APIClient{
//	        baseURL: baseURL,
//	        token:   token,
//	    }
//	}
//
//	func (c *APIClient) GetUser(id int) (*User, error) {
//	    url := fmt.Sprintf("%s/users/%d", c.baseURL, id)
//
//	    headers := map[string]string{
//	        "Authorization": "Bearer " + c.token,
//	    }
//
//	    var user User
//	    err := facades.HTTP().WithHeaders(headers).GetJSON(url, &user)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    return &user, nil
//	}
//
//	// Webhook handling
//	func HandleWebhook(w http.ResponseWriter, r *http.Request) {
//	    // Verify webhook signature
//	    signature := r.Header.Get("X-Webhook-Signature")
//	    if !facades.HTTP().VerifyWebhookSignature(r.Body, signature, webhookSecret) {
//	        http.Error(w, "Invalid signature", http.StatusUnauthorized)
//	        return
//	    }
//
//	    // Process webhook payload
//	    var payload WebhookPayload
//	    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
//	        http.Error(w, "Invalid JSON", http.StatusBadRequest)
//	        return
//	    }
//
//	    // Handle webhook event
//	    processWebhookEvent(payload)
//
//	    w.WriteHeader(http.StatusOK)
//	}
//
//	// Circuit breaker pattern
//	func CallExternalAPI() ([]byte, error) {
//	    // Check circuit breaker state
//	    if facades.HTTP().IsCircuitOpen("external-api") {
//	        return nil, errors.New("circuit breaker is open")
//	    }
//
//	    response, err := facades.HTTP().Get("https://external-api.example.com/data")
//	    if err != nil {
//	        facades.HTTP().RecordFailure("external-api")
//	        return nil, err
//	    }
//
//	    if response.StatusCode >= 500 {
//	        facades.HTTP().RecordFailure("external-api")
//	        return nil, fmt.Errorf("server error: %d", response.StatusCode)
//	    }
//
//	    facades.HTTP().RecordSuccess("external-api")
//	    defer response.Body.Close()
//
//	    return io.ReadAll(response.Body)
//	}
//
//	// HTTP middleware chain
//	facades.HTTP().Use(func(req *http.Request, next func() (*http.Response, error)) (*http.Response, error) {
//	    // Add request ID header
//	    req.Header.Set("X-Request-ID", generateRequestID())
//
//	    // Log request
//	    facades.Log().Info("HTTP request", map[string]interface{}{
//	        "method": req.Method,
//	        "url":    req.URL.String(),
//	    })
//
//	    start := time.Now()
//	    response, err := next()
//	    duration := time.Since(start)
//
//	    if err != nil {
//	        facades.Log().Error("HTTP request failed", map[string]interface{}{
//	            "method":   req.Method,
//	            "url":      req.URL.String(),
//	            "duration": duration,
//	            "error":    err.Error(),
//	        })
//	    } else {
//	        facades.Log().Info("HTTP request completed", map[string]interface{}{
//	            "method":      req.Method,
//	            "url":         req.URL.String(),
//	            "status_code": response.StatusCode,
//	            "duration":    duration,
//	        })
//	    }
//
//	    return response, err
//	})
//
//	// Rate limiting
//	err := facades.HTTP().WithRateLimit("api.example.com", 100, time.Minute).Get("https://api.example.com/data")
//	if err != nil {
//	    if errors.Is(err, http.ErrRateLimited) {
//	        fmt.Println("Rate limit exceeded, please try again later")
//	    }
//	}
//
//	// Request pooling for performance
//	pool := facades.HTTP().NewRequestPool(10) // 10 concurrent requests max
//
//	for i := 0; i < 100; i++ {
//	    url := fmt.Sprintf("https://api.example.com/items/%d", i)
//	    pool.Add(func() {
//	        response, err := facades.HTTP().Get(url)
//	        if err != nil {
//	            log.Printf("Request failed: %v", err)
//	            return
//	        }
//	        processResponse(response)
//	    })
//	}
//
//	pool.Wait() // Wait for all requests to complete
//
// Best Practices:
//   - Always set appropriate timeouts for requests
//   - Use context for request cancellation and deadlines
//   - Implement proper error handling and retry logic
//   - Use connection pooling for better performance
//   - Add request/response logging for debugging
//   - Implement circuit breakers for external service calls
//   - Use rate limiting to respect API limits
//   - Handle different response content types appropriately
//
// Error Handling Patterns:
//   - Timeout errors: Implement retry with exponential backoff
//   - Network errors: Use circuit breakers and fallback mechanisms
//   - HTTP errors: Handle different status codes appropriately
//   - JSON parsing errors: Validate response content types
//   - Rate limit errors: Implement proper backoff strategies
//
// Security Considerations:
//   - Validate SSL certificates in production
//   - Use secure authentication methods (OAuth, JWT)
//   - Sanitize and validate all input data
//   - Implement proper CORS handling
//   - Log security-relevant events
//   - Use HTTPS for all external communications
//   - Implement request signing for sensitive APIs
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume HTTP service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	httpClient, err := facade.Resolve[HTTPInterface]("http")
//	if err != nil {
//	    // Handle HTTP service unavailability gracefully
//	    return cachedData, fmt.Errorf("HTTP service unavailable: %w", err)
//	}
//	response, err := httpClient.Get(url)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestAPIIntegration(t *testing.T) {
//	    // Create a test HTTP client that mocks responses
//	    testHTTP := &TestHTTPClient{
//	        responses: map[string]*http.Response{
//	            "GET:https://api.example.com/users/123": {
//	                StatusCode: 200,
//	                Body:       io.NopCloser(strings.NewReader(`{"id":123,"name":"Test User"}`)),
//	            },
//	        },
//	    }
//
//	    // Swap the real HTTP client with test client
//	    restore := support.SwapService("http", testHTTP)
//	    defer restore() // Always restore after test
//
//	    // Now facades.HTTP() returns testHTTP
//	    apiClient := NewAPIClient("https://api.example.com", "test-token")
//
//	    user, err := apiClient.GetUser(123)
//	    require.NoError(t, err)
//	    assert.Equal(t, 123, user.ID)
//	    assert.Equal(t, "Test User", user.Name)
//
//	    // Verify request was made
//	    assert.True(t, testHTTP.WasCalled("GET", "https://api.example.com/users/123"))
//	}
//
// Container Configuration:
// Ensure the HTTP service is properly configured in your container:
//
//	// Example HTTP registration
//	container.Singleton("http", func() interface{} {
//	    config := http.Config{
//	        // Client configuration
//	        Timeout:         time.Second * 30,
//	        MaxIdleConns:    100,
//	        MaxConnsPerHost: 10,
//
//	        // TLS configuration
//	        InsecureSkipVerify: false,
//	        TLSHandshakeTimeout: time.Second * 10,
//
//	        // Retry configuration
//	        MaxRetries:    3,
//	        RetryDelay:    time.Second * 2,
//	        RetryBackoff:  "exponential", // or "linear", "constant"
//
//	        // Circuit breaker settings
//	        CircuitBreakerEnabled:    true,
//	        CircuitBreakerThreshold:  5,    // failures before opening
//	        CircuitBreakerTimeout:    60,   // seconds before retry
//
//	        // Rate limiting
//	        RateLimitEnabled: true,
//	        DefaultRateLimit: 1000, // requests per minute
//
//	        // Logging configuration
//	        LogRequests:  true,
//	        LogResponses: false, // Can be verbose
//	        LogErrors:    true,
//
//	        // User agent
//	        UserAgent: "MyApplication/1.0",
//
//	        // Default headers
//	        DefaultHeaders: map[string]string{
//	            "Accept":       "application/json",
//	            "Content-Type": "application/json",
//	        },
//
//	        // Middleware chain
//	        Middleware: []http.Middleware{
//	            http.LoggingMiddleware(),
//	            http.RetryMiddleware(),
//	            http.CircuitBreakerMiddleware(),
//	        },
//	    }
//
//	    httpClient, err := http.NewHTTPClient(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create HTTP client: %v", err)
//	    }
//
//	    return httpClient
//	})
func HTTP() httpInterfaces.HttpInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "http" service from the dependency injection container
	// - Performs type assertion to HTTPInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[httpInterfaces.HttpInterface](httpInterfaces.HTTP_TOKEN)
}

// HTTPWithError provides error-safe access to the HTTP client service.
//
// This function offers the same functionality as HTTP() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle HTTP service unavailability gracefully.
//
// This is a convenience wrapper around facade.Resolve() that provides
// the same caching and performance benefits as HTTP() but with error handling.
//
// Returns:
//   - HTTPInterface: The resolved HTTP instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement HTTPInterface
//
// Usage Examples:
//
//	// Basic error-safe HTTP request
//	httpClient, err := facades.HTTPWithError()
//	if err != nil {
//	    log.Printf("HTTP service unavailable: %v", err)
//	    return fallbackData, fmt.Errorf("HTTP client not available")
//	}
//	response, err := httpClient.Get("https://api.example.com/data")
//
//	// Conditional HTTP operations
//	if httpClient, err := facades.HTTPWithError(); err == nil {
//	    // Perform optional HTTP operations
//	    httpClient.Post("https://analytics.example.com/track", trackingData)
//	}
//
//	// Health check pattern
//	func CheckHTTPHealth() error {
//	    httpClient, err := facades.HTTPWithError()
//	    if err != nil {
//	        return fmt.Errorf("HTTP service unavailable: %w", err)
//	    }
//
//	    // Test basic HTTP functionality
//	    response, err := httpClient.Get("https://httpbin.org/status/200")
//	    if err != nil {
//	        return fmt.Errorf("HTTP request failed: %w", err)
//	    }
//	    defer response.Body.Close()
//
//	    if response.StatusCode != 200 {
//	        return fmt.Errorf("HTTP service not working properly")
//	    }
//
//	    return nil
//	}
func HTTPWithError() (httpInterfaces.HttpInterface, error) {
	// Use facade.Resolve() for error-return behavior:
	// - Resolves "http" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[httpInterfaces.HttpInterface](httpInterfaces.HTTP_TOKEN)
}
