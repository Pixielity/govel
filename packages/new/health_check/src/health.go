// Package healthcheck provides a comprehensive health monitoring system for GoVel applications.
// This package offers a registry-based approach to health checks with support for multiple
// check types, HTTP endpoints, notifications, and flexible result storage.
//
// Key Features:
// - Registry pattern for health check management
// - Multiple built-in health check types
// - HTTP endpoints for monitoring integration
// - Configurable timeouts and concurrency
// - Rich result metadata and JSON/HTML output
// - Background job support and caching
// - Integration with GoVel config and logger systems
//
// Example usage:
//
//	health := healthcheck.New()
//	health.Checks([]healthcheck.CheckInterface{
//	    checks.NewPingCheck().Name("api").URL("https://api.example.com"),
//	    checks.NewUsedDiskSpaceCheck().Path("/").WarnAt(70).FailAt(90),
//	})
//
//	results := health.RunChecks(context.Background())
//	if results.ContainsFailingCheck() {
//	    log.Println("Health checks failed!")
//	}
package healthcheck

import (
	"context"
	"time"

	"govel/healthcheck/http/controllers"
	"govel/healthcheck/interfaces"
	"govel/healthcheck/registry"
	"govel/healthcheck/types"
)

// Health is the main facade for the health check system.
// It provides a convenient interface for registering and executing health checks.
type Health struct {
	registry interfaces.HealthRegistryInterface
}

// New creates a new Health instance with default configuration.
// The instance comes with a configured registry and default settings.
//
// Returns:
//
//	*Health: A new health check instance ready for use
//
// Example:
//
//	health := healthcheck.New()
//	health.Register("database", &DatabaseCheck{})
func New() *Health {
	return &Health{
		registry: registry.NewHealthRegistry(),
	}
}

// NewWithRegistry creates a new Health instance with a custom registry.
// This allows for advanced configuration of the underlying registry.
//
// Parameters:
//
//	registry: A custom health registry implementation
//
// Returns:
//
//	*Health: A new health check instance with the custom registry
func NewWithRegistry(registry interfaces.HealthRegistryInterface) *Health {
	return &Health{
		registry: registry,
	}
}

// Register registers a single health check with a custom name.
//
// Parameters:
//
//	name: Unique identifier for the health check
//	check: The health check implementation
//
// Returns:
//
//	error: Any error during registration
//
// Example:
//
//	err := health.Register("database", &DatabaseCheck{
//	    ConnectionString: "postgres://...",
//	})
func (h *Health) Register(name string, check interfaces.CheckInterface) error {
	return h.registry.Register(name, check)
}

// Checks registers multiple health checks from a slice.
// This is the fluent interface method similar to Laravel Health.
//
// Parameters:
//
//	checks: Slice of health check implementations
//
// Returns:
//
//	*Health: Self for method chaining
//	error: Any error during registration
//
// Example:
//
//	_, err := health.Checks([]healthcheck.CheckInterface{
//	    checks.NewPingCheck().Name("api").URL("https://api.example.com"),
//	    checks.NewDiskCheck().Path("/").WarnAt(70).FailAt(90),
//	})
func (h *Health) Checks(checks []interfaces.CheckInterface) (*Health, error) {
	_, err := h.registry.Checks(checks)
	return h, err
}

// RunChecks executes all registered health checks.
//
// Parameters:
//
//	ctx: Context for timeout and cancellation control
//
// Returns:
//
//	interfaces.CheckResultsInterface: Collection of all check results
//
// Example:
//
//	results := health.RunChecks(context.Background())
//	if results.ContainsFailingCheck() {
//	    log.Println("Health checks failed")
//	}
func (h *Health) RunChecks(ctx context.Context) interfaces.CheckResultsInterface {
	return h.registry.RunChecks(ctx)
}

// RunCheck executes a specific health check by name.
//
// Parameters:
//
//	ctx: Context for timeout and cancellation control
//	name: Name of the health check to run
//
// Returns:
//
//	interfaces.ResultInterface: The result of the check execution
//	error: Any error during execution or if check not found
func (h *Health) RunCheck(ctx context.Context, name string) (interfaces.ResultInterface, error) {
	return h.registry.RunCheck(ctx, name)
}

// RunChecksWithTimeout executes all checks with a global timeout.
//
// Parameters:
//
//	timeout: Maximum time to wait for all checks to complete
//
// Returns:
//
//	interfaces.CheckResultsInterface: Collection of all check results
func (h *Health) RunChecksWithTimeout(timeout time.Duration) interfaces.CheckResultsInterface {
	return h.registry.RunChecksWithTimeout(timeout)
}

// RunChecksAsync executes all checks concurrently and returns a channel.
//
// Parameters:
//
//	ctx: Context for timeout and cancellation control
//
// Returns:
//
//	<-chan interfaces.CheckResultsInterface: Channel that will receive the results
func (h *Health) RunChecksAsync(ctx context.Context) <-chan interfaces.CheckResultsInterface {
	return h.registry.RunChecksAsync(ctx)
}

// GetRegistry returns the underlying health registry.
// This allows access to advanced registry features.
//
// Returns:
//
//	interfaces.HealthRegistryInterface: The health registry
func (h *Health) GetRegistry() interfaces.HealthRegistryInterface {
	return h.registry
}

// SetDefaultTimeout sets the default timeout for check execution.
//
// Parameters:
//
//	timeout: Default timeout duration
//
// Returns:
//
//	*Health: Self for method chaining
func (h *Health) SetDefaultTimeout(timeout time.Duration) *Health {
	h.registry.SetDefaultTimeout(timeout)
	return h
}

// SetMaxConcurrency sets the maximum number of checks to run concurrently.
//
// Parameters:
//
//	maxConcurrency: Maximum concurrent executions (0 = unlimited)
//
// Returns:
//
//	*Health: Self for method chaining
func (h *Health) SetMaxConcurrency(maxConcurrency int) *Health {
	h.registry.SetMaxConcurrency(maxConcurrency)
	return h
}

// WithResultStore sets the result store for persisting check results.
//
// Parameters:
//
//	store: Result store implementation
//
// Returns:
//
//	*Health: Self for method chaining
func (h *Health) WithResultStore(store interfaces.ResultStoreInterface) *Health {
	h.registry.WithResultStore(store)
	return h
}

// NewController creates a new HTTP controller for this health instance.
// The controller provides HTTP endpoints for health monitoring.
//
// Returns:
//
//	*controllers.HealthController: A new HTTP controller
//
// Example:
//
//	controller := health.NewController()
//	mux := http.NewServeMux()
//	controller.RegisterRoutes(mux)
func (h *Health) NewController() *controllers.HealthController {
	return controllers.NewHealthController(h.registry)
}

// GetNames returns the names of all registered health checks.
//
// Returns:
//
//	[]string: Slice of all registered check names
func (h *Health) GetNames() []string {
	return h.registry.GetNames()
}

// Count returns the number of registered health checks.
//
// Returns:
//
//	int: Number of registered checks
func (h *Health) Count() int {
	return h.registry.Count()
}

// Has checks if a health check is registered.
//
// Parameters:
//
//	name: Name of the health check to check
//
// Returns:
//
//	bool: true if check is registered, false otherwise
func (h *Health) Has(name string) bool {
	return h.registry.Has(name)
}

// Clear removes all registered health checks.
//
// Returns:
//
//	*Health: Self for method chaining
func (h *Health) Clear() *Health {
	h.registry.Clear()
	return h
}

// Clone creates a copy of the health instance with all its registered checks.
//
// Returns:
//
//	*Health: New health instance with copied checks
func (h *Health) Clone() *Health {
	return &Health{
		registry: h.registry.Clone(),
	}
}

// Package-level convenience functions
// ================================================================================

// NewResult creates a new result instance.
// This is a convenience function for creating results in custom health checks.
//
// Returns:
//
//	interfaces.ResultInterface: A new result instance
//
// Example:
//
//	result := healthcheck.NewResult()
//	return result.SetStatus(enums.StatusOK).SetMessage("All good!")
func NewResult() interfaces.ResultInterface {
	return types.NewResult()
}

// NewCheckResults creates a new check results collection.
// This is useful for creating custom result collections.
//
// Returns:
//
//	interfaces.CheckResultsInterface: A new check results collection
func NewCheckResults() interfaces.CheckResultsInterface {
	return types.NewCheckResults()
}

// NewRegistry creates a new health check registry.
// This is useful for creating custom registries with specific configurations.
//
// Returns:
//
//	interfaces.HealthRegistryInterface: A new health registry
func NewRegistry() interfaces.HealthRegistryInterface {
	return registry.NewHealthRegistry()
}

// NewRegistryWithDefaults creates a new registry with custom defaults.
//
// Parameters:
//
//	timeout: Default timeout for check execution
//	maxConcurrency: Maximum number of concurrent executions (0 = unlimited)
//
// Returns:
//
//	interfaces.HealthRegistryInterface: A new registry with custom defaults
func NewRegistryWithDefaults(timeout time.Duration, maxConcurrency int) interfaces.HealthRegistryInterface {
	return registry.NewHealthRegistryWithDefaults(timeout, maxConcurrency)
}

// Global health instance for convenience
// ================================================================================

var globalHealth *Health

// init initializes the global health instance
func init() {
	globalHealth = New()
}

// GetGlobalHealth returns the global health instance.
// This provides a singleton pattern for simple use cases.
//
// Returns:
//
//	*Health: The global health instance
//
// Example:
//
//	health := healthcheck.GetGlobalHealth()
//	health.Register("ping", &PingCheck{URL: "https://example.com"})
func GetGlobalHealth() *Health {
	return globalHealth
}

// SetGlobalHealth sets a custom global health instance.
// This is useful for dependency injection scenarios.
//
// Parameters:
//
//	health: The health instance to use as global
func SetGlobalHealth(health *Health) {
	globalHealth = health
}

// Global convenience functions that operate on the global instance
// ================================================================================

// RegisterGlobal registers a health check with the global instance.
//
// Parameters:
//
//	name: Unique identifier for the health check
//	check: The health check implementation
//
// Returns:
//
//	error: Any error during registration
func RegisterGlobal(name string, check interfaces.CheckInterface) error {
	return globalHealth.Register(name, check)
}

// ChecksGlobal registers multiple health checks with the global instance.
//
// Parameters:
//
//	checks: Slice of health check implementations
//
// Returns:
//
//	error: Any error during registration
func ChecksGlobal(checks []interfaces.CheckInterface) error {
	_, err := globalHealth.Checks(checks)
	return err
}

// RunChecksGlobal executes all checks registered with the global instance.
//
// Parameters:
//
//	ctx: Context for timeout and cancellation control
//
// Returns:
//
//	interfaces.CheckResultsInterface: Collection of all check results
func RunChecksGlobal(ctx context.Context) interfaces.CheckResultsInterface {
	return globalHealth.RunChecks(ctx)
}

// NewControllerGlobal creates a controller for the global health instance.
//
// Returns:
//
//	*controllers.HealthController: A new HTTP controller
func NewControllerGlobal() *controllers.HealthController {
	return globalHealth.NewController()
}
