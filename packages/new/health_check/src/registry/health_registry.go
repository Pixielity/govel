// Package registry provides the health check registry implementation.
// The registry manages health check registration, execution, and result aggregation.
package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"govel/packages/healthcheck/src/checks"
	"govel/packages/healthcheck/src/enums"
	"govel/packages/healthcheck/src/interfaces"
	"govel/packages/healthcheck/src/types"
)

// HealthRegistry is the central registry for managing health checks.
// It handles registration, execution, and result aggregation with support
// for concurrent execution and configurable timeouts.
type HealthRegistry struct {
	// checks holds all registered health checks by name
	checks map[string]interfaces.CheckInterface

	// resultStore is the configured result storage backend
	resultStore interfaces.ResultStoreInterface

	// defaultTimeout is the default timeout for check execution
	defaultTimeout time.Duration

	// maxConcurrency limits the number of concurrent check executions
	maxConcurrency int

	// mutex provides thread-safe access to the registry
	mutex sync.RWMutex
}

// NewHealthRegistry creates a new health check registry with default settings.
//
// Returns:
//
//	*HealthRegistry: A new registry instance ready for use
func NewHealthRegistry() *HealthRegistry {
	return &HealthRegistry{
		checks:         make(map[string]interfaces.CheckInterface),
		defaultTimeout: 30 * time.Second,
		maxConcurrency: 0, // 0 means unlimited
	}
}

// NewHealthRegistryWithDefaults creates a new registry with custom defaults.
//
// Parameters:
//
//	timeout: Default timeout for check execution
//	maxConcurrency: Maximum number of concurrent executions (0 = unlimited)
//
// Returns:
//
//	*HealthRegistry: A new registry instance with custom defaults
func NewHealthRegistryWithDefaults(timeout time.Duration, maxConcurrency int) *HealthRegistry {
	return &HealthRegistry{
		checks:         make(map[string]interfaces.CheckInterface),
		defaultTimeout: timeout,
		maxConcurrency: maxConcurrency,
	}
}

// Register registers a single health check with the registry.
//
// Parameters:
//
//	name: Unique identifier for the health check
//	check: The health check implementation
//
// Returns:
//
//	error: Any error during registration (e.g., duplicate names)
func (hr *HealthRegistry) Register(name string, check interfaces.CheckInterface) error {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	if name == "" {
		return fmt.Errorf("health check name cannot be empty")
	}

	if check == nil {
		return fmt.Errorf("health check cannot be nil")
	}

	if _, exists := hr.checks[name]; exists {
		return fmt.Errorf("health check with name '%s' is already registered", name)
	}

	hr.checks[name] = check
	return nil
}

// RegisterMultiple registers multiple health checks at once.
//
// Parameters:
//
//	checks: Map of name to check interface pairs
//
// Returns:
//
//	error: Any error during registration
func (hr *HealthRegistry) RegisterMultiple(checks map[string]interfaces.CheckInterface) error {
	for name, check := range checks {
		if err := hr.Register(name, check); err != nil {
			return err
		}
	}
	return nil
}

// Checks registers multiple health checks from a slice.
// This method provides a fluent interface similar to Laravel Health.
//
// Parameters:
//
//	checks: Slice of health check implementations
//
// Returns:
//
//	interfaces.HealthRegistryInterface: Self for method chaining
//	error: Any error during registration
func (hr *HealthRegistry) Checks(checks []interfaces.CheckInterface) (interfaces.HealthRegistryInterface, error) {
	for _, check := range checks {
		if check == nil {
			continue
		}

		name := check.GetName()
		if err := hr.Register(name, check); err != nil {
			return hr, err
		}
	}
	return hr, nil
}

// Unregister removes a health check from the registry.
//
// Parameters:
//
//	name: Name of the health check to remove
//
// Returns:
//
//	bool: true if check was found and removed, false otherwise
func (hr *HealthRegistry) Unregister(name string) bool {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	if _, exists := hr.checks[name]; exists {
		delete(hr.checks, name)
		return true
	}
	return false
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
func (hr *HealthRegistry) Has(name string) bool {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	_, exists := hr.checks[name]
	return exists
}

// Get retrieves a registered health check by name.
//
// Parameters:
//
//	name: Name of the health check to retrieve
//
// Returns:
//
//	interfaces.CheckInterface: The health check if found, nil otherwise
func (hr *HealthRegistry) Get(name string) interfaces.CheckInterface {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	return hr.checks[name]
}

// GetAll returns all registered health checks.
//
// Returns:
//
//	map[string]interfaces.CheckInterface: Map of all registered checks
func (hr *HealthRegistry) GetAll() map[string]interfaces.CheckInterface {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	// Create a copy to prevent external modification
	result := make(map[string]interfaces.CheckInterface)
	for name, check := range hr.checks {
		result[name] = check
	}
	return result
}

// GetNames returns the names of all registered health checks.
//
// Returns:
//
//	[]string: Slice of all registered check names
func (hr *HealthRegistry) GetNames() []string {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	names := make([]string, 0, len(hr.checks))
	for name := range hr.checks {
		names = append(names, name)
	}
	return names
}

// Count returns the number of registered health checks.
//
// Returns:
//
//	int: Number of registered checks
func (hr *HealthRegistry) Count() int {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	return len(hr.checks)
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
func (hr *HealthRegistry) RunChecks(ctx context.Context) interfaces.CheckResultsInterface {
	hr.mutex.RLock()
	checks := make(map[string]interfaces.CheckInterface)
	for name, check := range hr.checks {
		checks[name] = check
	}
	hr.mutex.RUnlock()

	results := types.NewCheckResultsWithCapacity(len(checks))
	results.SetExecutedAt(time.Now())

	// If no checks registered, return empty results
	if len(checks) == 0 {
		return results
	}

	// Execute checks with concurrency control
	return hr.executeChecksWithConcurrency(ctx, checks, results)
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
func (hr *HealthRegistry) RunCheck(ctx context.Context, name string) (interfaces.ResultInterface, error) {
	check := hr.Get(name)
	if check == nil {
		return nil, fmt.Errorf("health check '%s' not found", name)
	}

	return hr.executeCheck(ctx, name, check), nil
}

// RunChecksWithNames executes only the specified health checks.
//
// Parameters:
//
//	ctx: Context for timeout and cancellation control
//	names: Slice of check names to execute
//
// Returns:
//
//	interfaces.CheckResultsInterface: Collection of specified check results
func (hr *HealthRegistry) RunChecksWithNames(ctx context.Context, names []string) interfaces.CheckResultsInterface {
	hr.mutex.RLock()
	checks := make(map[string]interfaces.CheckInterface)
	for _, name := range names {
		if check, exists := hr.checks[name]; exists {
			checks[name] = check
		}
	}
	hr.mutex.RUnlock()

	results := types.NewCheckResultsWithCapacity(len(checks))
	results.SetExecutedAt(time.Now())

	return hr.executeChecksWithConcurrency(ctx, checks, results)
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
func (hr *HealthRegistry) RunChecksWithTimeout(timeout time.Duration) interfaces.CheckResultsInterface {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return hr.RunChecks(ctx)
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
func (hr *HealthRegistry) RunChecksAsync(ctx context.Context) <-chan interfaces.CheckResultsInterface {
	resultChan := make(chan interfaces.CheckResultsInterface, 1)

	go func() {
		defer close(resultChan)
		results := hr.RunChecks(ctx)
		resultChan <- results
	}()

	return resultChan
}

// SetDefaultTimeout sets the default timeout for check execution.
//
// Parameters:
//
//	timeout: Default timeout duration
//
// Returns:
//
//	interfaces.HealthRegistryInterface: Self for method chaining
func (hr *HealthRegistry) SetDefaultTimeout(timeout time.Duration) interfaces.HealthRegistryInterface {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.defaultTimeout = timeout
	return hr
}

// GetDefaultTimeout returns the default timeout for check execution.
//
// Returns:
//
//	time.Duration: Default timeout duration
func (hr *HealthRegistry) GetDefaultTimeout() time.Duration {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	return hr.defaultTimeout
}

// SetMaxConcurrency sets the maximum number of checks to run concurrently.
//
// Parameters:
//
//	maxConcurrency: Maximum concurrent executions (0 = unlimited)
//
// Returns:
//
//	interfaces.HealthRegistryInterface: Self for method chaining
func (hr *HealthRegistry) SetMaxConcurrency(maxConcurrency int) interfaces.HealthRegistryInterface {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.maxConcurrency = maxConcurrency
	return hr
}

// GetMaxConcurrency returns the maximum number of concurrent executions.
//
// Returns:
//
//	int: Maximum concurrent executions
func (hr *HealthRegistry) GetMaxConcurrency() int {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	return hr.maxConcurrency
}

// Clear removes all registered health checks.
//
// Returns:
//
//	interfaces.HealthRegistryInterface: Self for method chaining
func (hr *HealthRegistry) Clear() interfaces.HealthRegistryInterface {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.checks = make(map[string]interfaces.CheckInterface)
	return hr
}

// Clone creates a copy of the registry with all its registered checks.
//
// Returns:
//
//	interfaces.HealthRegistryInterface: New registry instance with copied checks
func (hr *HealthRegistry) Clone() interfaces.HealthRegistryInterface {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	clone := &HealthRegistry{
		checks:         make(map[string]interfaces.CheckInterface),
		resultStore:    hr.resultStore,
		defaultTimeout: hr.defaultTimeout,
		maxConcurrency: hr.maxConcurrency,
	}

	for name, check := range hr.checks {
		clone.checks[name] = check
	}

	return clone
}

// WithResultStore sets the result store for persisting check results.
//
// Parameters:
//
//	store: Result store implementation
//
// Returns:
//
//	interfaces.HealthRegistryInterface: Self for method chaining
func (hr *HealthRegistry) WithResultStore(store interfaces.ResultStoreInterface) interfaces.HealthRegistryInterface {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.resultStore = store
	return hr
}

// GetResultStore returns the configured result store.
//
// Returns:
//
//	interfaces.ResultStoreInterface: The configured result store, nil if none set
func (hr *HealthRegistry) GetResultStore() interfaces.ResultStoreInterface {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	return hr.resultStore
}

// executeChecksWithConcurrency executes checks with concurrency control.
func (hr *HealthRegistry) executeChecksWithConcurrency(
	ctx context.Context,
	checks map[string]interfaces.CheckInterface,
	results *types.CheckResults,
) interfaces.CheckResultsInterface {
	if len(checks) == 0 {
		return results
	}

	// Channel to collect results
	resultChan := make(chan interfaces.ResultInterface, len(checks))

	// Semaphore for concurrency control
	var semaphore chan struct{}
	if hr.maxConcurrency > 0 {
		semaphore = make(chan struct{}, hr.maxConcurrency)
	}

	var wg sync.WaitGroup

	// Execute each check
	for name, check := range checks {
		wg.Add(1)
		go func(name string, check interfaces.CheckInterface) {
			defer wg.Done()

			// Acquire semaphore if concurrency is limited
			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}

			result := hr.executeCheck(ctx, name, check)
			resultChan <- result
		}(name, check)
	}

	// Close result channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect all results
	for result := range resultChan {
		results.AddResult(result)
	}

	// Store results if result store is configured
	if hr.resultStore != nil {
		if err := hr.resultStore.Store(results); err != nil {
			// Log error but don't fail the entire operation
			// In a real implementation, you'd use the GoVel logger here
		}
	}

	return results
}

// executeCheck executes a single health check with proper error handling.
func (hr *HealthRegistry) executeCheck(ctx context.Context, name string, check interfaces.CheckInterface) interfaces.ResultInterface {
	// Create a result with timing information
	result := checks.NewResult()
	result.SetCheck(check)
	result.SetStartedAt(time.Now())

	// Create timeout context
	checkCtx, cancel := context.WithTimeout(ctx, hr.defaultTimeout)
	defer cancel()

	// Execute the check with panic recovery
	func() {
		defer func() {
			result.SetEndedAt(time.Now())

			if r := recover(); r != nil {
				result.SetStatus(enums.StatusCrashed)
				result.SetNotificationMessage(fmt.Sprintf("Health check panicked: %v", r))
				result.SetShortSummary("Panicked")
			}
		}()

		// Check for context cancellation before execution
		select {
		case <-checkCtx.Done():
			result.SetStatus(enums.StatusFailed)
			result.SetNotificationMessage("Health check timed out")
			result.SetShortSummary("Timeout")
			return
		default:
		}

		// Execute the actual health check
		checkResult := check.Run()
		if checkResult != nil {
			// Copy the result from the check execution
			result.SetStatus(checkResult.GetStatus())
			result.SetNotificationMessage(checkResult.GetNotificationMessage())
			result.SetShortSummary(checkResult.GetShortSummary())
			result.SetMeta(checkResult.GetMeta())
		} else {
			// Handle nil result
			result.SetStatus(enums.StatusFailed)
			result.SetNotificationMessage("Health check returned nil result")
			result.SetShortSummary("Nil Result")
		}
	}()

	return result
}

// Compile-time interface compliance check
var _ interfaces.HealthRegistryInterface = (*HealthRegistry)(nil)
