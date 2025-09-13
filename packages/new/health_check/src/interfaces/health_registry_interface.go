// Package interfaces provides registry interface definitions for the GoVel health check system.
package interfaces

import (
	"context"
	"time"
)

// HealthRegistryInterface defines the contract for health check registry implementations.
// The registry is responsible for managing health check registration, execution,
// and result aggregation.
//
// Example usage:
//
//	registry := healthcheck.NewRegistry()
//	registry.Register("database", &DatabaseCheck{})
//	registry.Register("redis", &RedisCheck{})
//	results := registry.RunChecks(context.Background())
type HealthRegistryInterface interface {
	// Register registers a single health check with the registry.
	//
	// Parameters:
	//   name: Unique identifier for the health check
	//   check: The health check implementation
	//
	// Returns:
	//   error: Any error during registration (e.g., duplicate names)
	Register(name string, check CheckInterface) error

	// RegisterMultiple registers multiple health checks at once.
	//
	// Parameters:
	//   checks: Map of name to check interface pairs
	//
	// Returns:
	//   error: Any error during registration
	RegisterMultiple(checks map[string]CheckInterface) error

	// Checks registers multiple health checks from a slice.
	// This method provides a fluent interface similar to Laravel Health.
	//
	// Parameters:
	//   checks: Slice of health check implementations
	//
	// Returns:
	//   HealthRegistryInterface: Self for method chaining
	//   error: Any error during registration
	Checks(checks []CheckInterface) (HealthRegistryInterface, error)

	// Unregister removes a health check from the registry.
	//
	// Parameters:
	//   name: Name of the health check to remove
	//
	// Returns:
	//   bool: true if check was found and removed, false otherwise
	Unregister(name string) bool

	// Has checks if a health check is registered.
	//
	// Parameters:
	//   name: Name of the health check to check
	//
	// Returns:
	//   bool: true if check is registered, false otherwise
	Has(name string) bool

	// Get retrieves a registered health check by name.
	//
	// Parameters:
	//   name: Name of the health check to retrieve
	//
	// Returns:
	//   CheckInterface: The health check if found, nil otherwise
	Get(name string) CheckInterface

	// GetAll returns all registered health checks.
	//
	// Returns:
	//   map[string]CheckInterface: Map of all registered checks
	GetAll() map[string]CheckInterface

	// GetNames returns the names of all registered health checks.
	//
	// Returns:
	//   []string: Slice of all registered check names
	GetNames() []string

	// Count returns the number of registered health checks.
	//
	// Returns:
	//   int: Number of registered checks
	Count() int

	// RunChecks executes all registered health checks.
	//
	// Parameters:
	//   ctx: Context for timeout and cancellation control
	//
	// Returns:
	//   CheckResultsInterface: Collection of all check results
	RunChecks(ctx context.Context) CheckResultsInterface

	// RunCheck executes a specific health check by name.
	//
	// Parameters:
	//   ctx: Context for timeout and cancellation control
	//   name: Name of the health check to run
	//
	// Returns:
	//   ResultInterface: The result of the check execution
	//   error: Any error during execution or if check not found
	RunCheck(ctx context.Context, name string) (ResultInterface, error)

	// RunChecksWithNames executes only the specified health checks.
	//
	// Parameters:
	//   ctx: Context for timeout and cancellation control
	//   names: Slice of check names to execute
	//
	// Returns:
	//   CheckResultsInterface: Collection of specified check results
	RunChecksWithNames(ctx context.Context, names []string) CheckResultsInterface

	// RunChecksWithTimeout executes all checks with a global timeout.
	//
	// Parameters:
	//   timeout: Maximum time to wait for all checks to complete
	//
	// Returns:
	//   CheckResultsInterface: Collection of all check results
	RunChecksWithTimeout(timeout time.Duration) CheckResultsInterface

	// RunChecksAsync executes all checks concurrently.
	//
	// Parameters:
	//   ctx: Context for timeout and cancellation control
	//
	// Returns:
	//   <-chan CheckResultsInterface: Channel that will receive the results
	RunChecksAsync(ctx context.Context) <-chan CheckResultsInterface

	// SetDefaultTimeout sets the default timeout for check execution.
	//
	// Parameters:
	//   timeout: Default timeout duration
	//
	// Returns:
	//   HealthRegistryInterface: Self for method chaining
	SetDefaultTimeout(timeout time.Duration) HealthRegistryInterface

	// GetDefaultTimeout returns the default timeout for check execution.
	//
	// Returns:
	//   time.Duration: Default timeout duration
	GetDefaultTimeout() time.Duration

	// SetMaxConcurrency sets the maximum number of checks to run concurrently.
	//
	// Parameters:
	//   maxConcurrency: Maximum concurrent executions (0 = unlimited)
	//
	// Returns:
	//   HealthRegistryInterface: Self for method chaining
	SetMaxConcurrency(maxConcurrency int) HealthRegistryInterface

	// GetMaxConcurrency returns the maximum number of concurrent executions.
	//
	// Returns:
	//   int: Maximum concurrent executions
	GetMaxConcurrency() int

	// Clear removes all registered health checks.
	//
	// Returns:
	//   HealthRegistryInterface: Self for method chaining
	Clear() HealthRegistryInterface

	// Clone creates a copy of the registry with all its registered checks.
	//
	// Returns:
	//   HealthRegistryInterface: New registry instance with copied checks
	Clone() HealthRegistryInterface

	// WithResultStore sets the result store for persisting check results.
	//
	// Parameters:
	//   store: Result store implementation
	//
	// Returns:
	//   HealthRegistryInterface: Self for method chaining
	WithResultStore(store ResultStoreInterface) HealthRegistryInterface

	// GetResultStore returns the configured result store.
	//
	// Returns:
	//   ResultStoreInterface: The configured result store, nil if none set
	GetResultStore() ResultStoreInterface
}

// ResultStoreInterface defines the contract for result storage backends.
// Result stores handle persistence and retrieval of health check results.
type ResultStoreInterface interface {
	// Store persists health check results.
	//
	// Parameters:
	//   results: The check results to store
	//
	// Returns:
	//   error: Any error during storage
	Store(results CheckResultsInterface) error

	// Get retrieves the latest stored results.
	//
	// Returns:
	//   CheckResultsInterface: Latest stored results, nil if none found
	//   error: Any error during retrieval
	Get() (CheckResultsInterface, error)

	// GetByTimestamp retrieves results from a specific timestamp.
	//
	// Parameters:
	//   timestamp: The timestamp to retrieve results for
	//
	// Returns:
	//   CheckResultsInterface: Results from the specified timestamp, nil if none found
	//   error: Any error during retrieval
	GetByTimestamp(timestamp time.Time) (CheckResultsInterface, error)

	// GetHistory retrieves historical results within a time range.
	//
	// Parameters:
	//   from: Start time for the range
	//   to: End time for the range
	//
	// Returns:
	//   []CheckResultsInterface: Historical results within the range
	//   error: Any error during retrieval
	GetHistory(from, to time.Time) ([]CheckResultsInterface, error)

	// Clear removes all stored results.
	//
	// Returns:
	//   error: Any error during clearing
	Clear() error

	// GetStorageInfo returns information about the storage backend.
	//
	// Returns:
	//   map[string]interface{}: Storage information and statistics
	GetStorageInfo() map[string]interface{}
}

// HealthControllerInterface defines the contract for HTTP health check controllers.
// Controllers handle HTTP requests and provide health status endpoints.
type HealthControllerInterface interface {
	// HandleHealthCheck handles the main health check endpoint.
	// Typically responds with HTML dashboard or JSON based on Accept header.
	//
	// Parameters:
	//   w: HTTP response writer
	//   r: HTTP request
	HandleHealthCheck(w ResponseWriter, r RequestInterface)

	// HandleHealthCheckJSON handles JSON health check endpoint.
	// Always responds with JSON format.
	//
	// Parameters:
	//   w: HTTP response writer
	//   r: HTTP request
	HandleHealthCheckJSON(w ResponseWriter, r RequestInterface)

	// HandleSimpleHealthCheck handles simple text health check endpoint.
	// Responds with simple "OK" or "FAILED" text.
	//
	// Parameters:
	//   w: HTTP response writer
	//   r: HTTP request
	HandleSimpleHealthCheck(w ResponseWriter, r RequestInterface)

	// HandleReadinessCheck handles Kubernetes-style readiness probe.
	//
	// Parameters:
	//   w: HTTP response writer
	//   r: HTTP request
	HandleReadinessCheck(w ResponseWriter, r RequestInterface)

	// HandleLivenessCheck handles Kubernetes-style liveness probe.
	//
	// Parameters:
	//   w: HTTP response writer
	//   r: HTTP request
	HandleLivenessCheck(w ResponseWriter, r RequestInterface)

	// SetRegistry sets the health registry for the controller.
	//
	// Parameters:
	//   registry: The health registry to use
	//
	// Returns:
	//   HealthControllerInterface: Self for method chaining
	SetRegistry(registry HealthRegistryInterface) HealthControllerInterface

	// GetRegistry returns the configured health registry.
	//
	// Returns:
	//   HealthRegistryInterface: The configured registry
	GetRegistry() HealthRegistryInterface

	// SetResultStore sets the result store for the controller.
	//
	// Parameters:
	//   store: The result store to use
	//
	// Returns:
	//   HealthControllerInterface: Self for method chaining
	SetResultStore(store ResultStoreInterface) HealthControllerInterface
}

// ResponseWriter defines the contract for HTTP response writing.
// This interface allows for testing and different HTTP frameworks.
type ResponseWriter interface {
	// Header returns the response header map.
	Header() map[string][]string

	// Write writes data to the response body.
	Write([]byte) (int, error)

	// WriteHeader sends an HTTP response header with the provided status code.
	WriteHeader(statusCode int)
}

// RequestInterface defines the contract for HTTP requests.
// This interface allows for testing and different HTTP frameworks.
type RequestInterface interface {
	// GetMethod returns the HTTP method (GET, POST, etc.).
	GetMethod() string

	// GetURL returns the request URL.
	GetURL() string

	// GetHeader returns the value of a request header.
	GetHeader(key string) string

	// GetQuery returns the value of a query parameter.
	GetQuery(key string) string

	// GetContext returns the request context.
	GetContext() context.Context
}
