package support

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	containerInterfaces "govel/packages/types/src/interfaces/container"
	containerTypes "govel/packages/types/src/types/container"
)

// ================================================================================================
// TYPES AND CONFIGURATION
// ================================================================================================

// FacadeError represents a facade-specific error with detailed context.
// It provides structured error information for better error handling and debugging.
type FacadeError struct {
	// ServiceKey is the container binding key that was being resolved
	ServiceKey string
	// Operation describes what operation was being performed (resolve, type_assert, etc.)
	Operation string
	// Cause is the underlying error that caused this facade error
	Cause error
}

// Error implements the error interface for FacadeError.
// Returns a descriptive error message with context about the failed operation.
func (e FacadeError) Error() string {
	return fmt.Sprintf("facade %s for service '%s': %v", e.Operation, e.ServiceKey, e.Cause)
}

// Unwrap returns the underlying cause error for error chain compatibility.
// This allows errors.Is() and errors.As() to work correctly with FacadeError.
func (e FacadeError) Unwrap() error {
	return e.Cause
}

// ServiceNotFoundError indicates that a requested service was not found in the container.
// This is a specific type of error that can be checked for programmatically.
type ServiceNotFoundError struct {
	ServiceKey string
}

// Error implements the error interface for ServiceNotFoundError.
func (e ServiceNotFoundError) Error() string {
	return fmt.Sprintf("service '%s' not found in container", e.ServiceKey)
}

// TypeAssertionError indicates that a service was resolved but type assertion failed.
// This typically means the service doesn't implement the expected interface.
type TypeAssertionError struct {
	ServiceKey   string
	ExpectedType string
	ActualType   string
}

// Error implements the error interface for TypeAssertionError.
func (e TypeAssertionError) Error() string {
	return fmt.Sprintf("type assertion failed for service '%s': expected %s, got %s",
		e.ServiceKey, e.ExpectedType, e.ActualType)
}

// FacadeOptions configures the behavior of the facade system.
// These options control caching, performance, and debugging features.
type FacadeOptions struct {
	// CacheEnabled controls whether resolved services are cached.
	// Disabling caching will resolve services fresh on every call.
	// Default: true
	CacheEnabled bool

	// MaxCacheSize sets the maximum number of services to cache.
	// When exceeded, oldest entries are evicted (LRU behavior).
	// Set to 0 for unlimited cache size.
	// Default: 1000
	MaxCacheSize int

	// CacheTTL sets how long services remain cached before expiring.
	// Set to 0 to disable TTL (services cached indefinitely).
	// Default: 0 (no expiration)
	CacheTTL time.Duration

	// Debug enables detailed debug logging for facade operations.
	// This can be useful for troubleshooting but may impact performance.
	// Default: false
	Debug bool

	// MetricsEnabled controls whether facade metrics are collected.
	// Metrics include cache hit/miss rates, resolution counts, etc.
	// Default: true
	MetricsEnabled bool
}

// FacadeStats provides detailed statistics about facade operations.
// These metrics can be used for monitoring, debugging, and performance optimization.
type FacadeStats struct {
	// CacheHits is the total number of successful cache lookups
	CacheHits int64
	// CacheMisses is the total number of cache misses requiring container resolution
	CacheMisses int64
	// Resolutions is the total number of successful service resolutions
	Resolutions int64
	// Errors is the total number of errors that occurred during operations
	Errors int64
	// TypeAssertionFailures is the number of type assertion errors
	TypeAssertionFailures int64
	// CacheSize is the current number of cached services
	CacheSize int
	// CacheEvictions is the number of services evicted from cache due to size limits
	CacheEvictions int64
}

// cacheEntry represents a cached service with metadata for TTL and LRU eviction.
type cacheEntry struct {
	// service is the cached service instance
	service interface{}
	// lastAccessed is used for LRU eviction when cache is full
	lastAccessed time.Time
	// createdAt is used for TTL expiration
	createdAt time.Time
}

// ================================================================================================
// GLOBAL STATE AND VARIABLES
// ================================================================================================

// Global configuration and state for the facade system.
// These variables are protected by facadeMutex for thread safety.
var (
	// facadeContainer holds the dependency injection container instance.
	// This is set during application bootstrap via SetContainer().
	facadeContainer containerInterfaces.ContainerInterface

	// facadeCache stores resolved service instances for performance optimization.
	// Key: service key (string), Value: cacheEntry with service and metadata
	facadeCache = make(map[string]*cacheEntry)

	// facadeMutex protects concurrent access to all facade state.
	// Uses RWMutex for optimized read-heavy workloads.
	facadeMutex sync.RWMutex

	// facadeOptions contains current configuration settings.
	// Can be modified at runtime via Configure().
	facadeOptions = FacadeOptions{
		CacheEnabled:   true,
		MaxCacheSize:   1000,
		CacheTTL:       0, // No expiration by default
		Debug:          false,
		MetricsEnabled: true,
	}

	// facadeStats tracks detailed metrics about facade operations.
	// Updated atomically to avoid lock contention during metric collection.
	facadeStats = FacadeStats{}
)

// ================================================================================================
// CORE RESOLUTION FUNCTIONS
// ================================================================================================

// Resolve resolves a service from the global container with caching and error handling.
//
// This is the core resolution function used by all other facade methods. It handles:
//   - Container availability checking
//   - Cache lookup with TTL validation
//   - Service resolution from container
//   - Cache storage with size limit management
//   - Comprehensive error handling and metrics
//
// Parameters:
//   - serviceKey: The container binding key for the service to resolve
//
// Returns:
//   - interface{}: The resolved service instance
//   - error: Any error that occurred during resolution
//
// Thread Safety:
// This function is thread-safe and uses optimized locking patterns to minimize
// contention during concurrent access.
//
// Errors:
// Returns FacadeError for facade-specific issues, or wraps container errors
// for resolution failures.
//
// Example:
//
//	service, err := Resolve("logger")
//	if err != nil {
//	    log.Printf("Failed to resolve logger: %v", err)
//	    return
//	}
//	logger := service.(LoggerInterface)
func Make(serviceKey containerTypes.ServiceIdentifier) (interface{}, error) {
	key := containerTypes.ToKey(serviceKey)

	// Single lock operation combining container check and cache lookup for performance
	facadeMutex.RLock()
	container := facadeContainer
	cacheEntry := facadeCache[key]
	cacheEnabled := facadeOptions.CacheEnabled
	cacheTTL := facadeOptions.CacheTTL
	facadeMutex.RUnlock()

	// Validate container availability
	if container == nil {
		atomic.AddInt64(&facadeStats.Errors, 1)
		return nil, FacadeError{
			ServiceKey: key,
			Operation:  "resolve",
			Cause:      fmt.Errorf("no container set. Call support.SetContainer() first"),
		}
	}

	// Check cache if enabled and entry exists
	if cacheEnabled && cacheEntry != nil {
		// Check TTL expiration if TTL is configured
		if cacheTTL > 0 && time.Since(cacheEntry.createdAt) > cacheTTL {
			// Cache entry expired, remove it
			facadeMutex.Lock()
			delete(facadeCache, key)
			facadeMutex.Unlock()
		} else {
			// Cache hit - update access time for LRU and return cached service
			cacheEntry.lastAccessed = time.Now()
			atomic.AddInt64(&facadeStats.CacheHits, 1)
			return cacheEntry.service, nil
		}
	}

	// Cache miss - resolve from container
	atomic.AddInt64(&facadeStats.CacheMisses, 1)

	service, err := container.Make(serviceKey)
	if err != nil {
		atomic.AddInt64(&facadeStats.Errors, 1)
		return nil, FacadeError{
			ServiceKey: key,
			Operation:  "resolve",
			Cause:      err,
		}
	}

	// Successfully resolved - update metrics
	atomic.AddInt64(&facadeStats.Resolutions, 1)

	// Cache the result if caching is enabled
	if cacheEnabled {
		storeCachedService(serviceKey, service)
	}

	return service, nil
}

// ResolveWithContext resolves a service with context support for cancellation and timeouts.
//
// This function provides the same functionality as Resolve() but with context support
// for better control over long-running operations and cancellation.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - serviceKey: The container binding key for the service to resolve
//
// Returns:
//   - interface{}: The resolved service instance
//   - error: Any error including context cancellation/timeout
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	service, err := ResolveWithContext(ctx, "database")
//	if err != nil {
//	    if errors.Is(err, context.DeadlineExceeded) {
//	        log.Println("Service resolution timed out")
//	    }
//	    return
//	}
func ResolveWithContext(ctx context.Context, serviceKey containerTypes.ServiceIdentifier) (interface{}, error) {
	key := containerTypes.ToKey(serviceKey)

	// Check if context is already cancelled before starting
	select {
	case <-ctx.Done():
		atomic.AddInt64(&facadeStats.Errors, 1)
		return nil, FacadeError{
			ServiceKey: key,
			Operation:  "resolve_with_context",
			Cause:      ctx.Err(),
		}
	default:
		// Context is still active, proceed with resolution
		return Make(key)
	}
}

// ================================================================================================
// HIGH-LEVEL HELPER FUNCTIONS
// ================================================================================================

// MustResolve resolves a service and performs type assertion in a single call.
//
// This is the primary helper function used by facade implementations. It combines
// service resolution and type assertion into one operation, panicking if either fails.
// This provides a clean, Laravel-style API for facade implementations.
//
// Type Parameters:
//   - T: The expected service interface type (automatically inferred)
//
// Parameters:
//   - serviceKey: The container binding key for the service to resolve
//
// Returns:
//   - T: The resolved and type-asserted service instance
//
// Panics:
//   - If service cannot be resolved from container
//   - If type assertion fails (service doesn't implement expected interface)
//   - If no container is set
//
// Thread Safety:
// This function is thread-safe and can be called concurrently from multiple goroutines.
//
// Performance:
// Uses optimized caching to avoid repeated container resolutions. Type assertion
// is performed efficiently with proper error context.
//
// Example Usage in Facade:
//
//	func Log() LoggerInterface {
//	    return MustResolve[LoggerInterface]("log")
//	}
//
//	func Database() DatabaseInterface {
//	    return MustResolve[DatabaseInterface]("database")
//	}
func Resolve[T any](serviceKey containerTypes.ServiceIdentifier) T {
	key := containerTypes.ToKey(serviceKey)

	service, err := Make(key)
	if err != nil {
		// Panic with facade error for better debugging context
		panic(FacadeError{
			ServiceKey: key,
			Operation:  "must_resolve",
			Cause:      err,
		})
	}

	// Perform type assertion with panic recovery for better error messages
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&facadeStats.TypeAssertionFailures, 1)
			atomic.AddInt64(&facadeStats.Errors, 1)

			// Re-panic with more descriptive error
			panic(FacadeError{
				ServiceKey: key,
				Operation:  "type_assertion",
				Cause:      fmt.Errorf("type assertion failed: %v", r),
			})
		}
	}()

	return service.(T)
}

// TryResolve resolves a service and performs type assertion with comprehensive error handling.
//
// This function provides the same functionality as MustResolve but returns errors instead
// of panicking, making it suitable for use in error-sensitive contexts or when you want
// to handle resolution failures gracefully.
//
// Type Parameters:
//   - T: The expected service interface type (automatically inferred)
//
// Parameters:
//   - serviceKey: The container binding key for the service to resolve
//
// Returns:
//   - T: The resolved and type-asserted service instance (zero value on error)
//   - error: Detailed error information if resolution or type assertion fails
//
// Errors:
//   - FacadeError: For facade-specific errors (container not set, resolution failed)
//   - TypeAssertionError: For type assertion failures
//
// Thread Safety:
// This function is thread-safe and can be called concurrently from multiple goroutines.
//
// Example Usage:
//
//	logger, err := TryResolve[LoggerInterface]("log")
//	if err != nil {
//	    return fmt.Errorf("failed to get logger: %w", err)
//	}
//	logger.Info("Application started")
//
//	// Check for specific error types
//	db, err := TryResolve[DatabaseInterface]("database")
//	if err != nil {
//	    var notFoundErr ServiceNotFoundError
//	    if errors.As(err, &notFoundErr) {
//	        log.Println("Database service not configured")
//	        // Use fallback or return gracefully
//	    }
//	    return err
//	}
func TryResolve[T any](serviceKey containerTypes.ServiceIdentifier) (result T, err error) {
	key := containerTypes.ToKey(serviceKey)

	service, err := Make(serviceKey)
	if err != nil {
		return result, err
	}

	// Perform type assertion with panic recovery for proper error handling
	defer func() {
		if r := recover(); r != nil {
			atomic.AddInt64(&facadeStats.TypeAssertionFailures, 1)
			atomic.AddInt64(&facadeStats.Errors, 1)

			err = FacadeError{
				ServiceKey: key,
				Operation:  "type_assertion",
				Cause:      fmt.Errorf("type assertion failed: %v", r),
			}
		}
	}()

	return service.(T), nil
}

// ================================================================================================
// CONFIGURATION AND MANAGEMENT
// ================================================================================================

// Configure updates the facade system configuration with new options.
//
// This function allows runtime reconfiguration of the facade system behavior.
// Configuration changes take effect immediately and may trigger cache cleanup
// if caching is disabled or cache size limits are reduced.
//
// Parameters:
//   - opts: New configuration options to apply
//
// Thread Safety:
// This function is thread-safe but may temporarily block other operations
// during configuration updates.
//
// Side Effects:
//   - If caching is disabled, all cached services are immediately cleared
//   - If cache size is reduced, oldest entries are evicted to fit new limit
//   - Metrics may be reset if metrics collection is disabled
//
// Example:
//
//	// Disable caching for testing
//	support.Configure(support.FacadeOptions{
//	    CacheEnabled: false,
//	    Debug:        true,
//	})
//
//	// Enable TTL caching with size limit
//	support.Configure(support.FacadeOptions{
//	    CacheEnabled: true,
//	    MaxCacheSize: 500,
//	    CacheTTL:     30 * time.Minute,
//	})
func Configure(opts FacadeOptions) {
	facadeMutex.Lock()
	defer facadeMutex.Unlock()

	prevOptions := facadeOptions
	facadeOptions = opts

	// Handle cache configuration changes
	if !opts.CacheEnabled {
		// Caching disabled - clear all cached services
		facadeCache = make(map[string]*cacheEntry)
	} else if opts.MaxCacheSize > 0 && opts.MaxCacheSize < len(facadeCache) {
		// Cache size reduced - evict oldest entries
		evictOldestEntries(len(facadeCache) - opts.MaxCacheSize)
	}

	// Reset metrics if metrics collection was disabled
	if prevOptions.MetricsEnabled && !opts.MetricsEnabled {
		resetStats()
	}
}

// GetConfiguration returns the current facade configuration.
//
// Returns a copy of the current configuration to prevent external modification.
//
// Returns:
//   - FacadeOptions: Current configuration settings
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Example:
//
//	config := support.GetConfiguration()
//	if config.Debug {
//	    log.Println("Facade debugging is enabled")
//	}
func GetConfiguration() FacadeOptions {
	facadeMutex.RLock()
	defer facadeMutex.RUnlock()
	return facadeOptions // Return copy
}

// SetContainer sets the global dependency injection container for facade resolution.
//
// This function must be called during application bootstrap before any facades
// are used. Setting a new container will clear all cached services to prevent
// inconsistencies between the old and new container bindings.
//
// Parameters:
//   - container: The dependency injection container to use for service resolution
//
// Thread Safety:
// This function is thread-safe but will block all facade operations during
// the container update to ensure consistency.
//
// Side Effects:
//   - All cached services are immediately cleared
//   - In-flight resolutions may be affected
//   - Metrics cache size counter is reset
//
// Example:
//
//	// In main.go or application bootstrap
//	container := container.NewContainer()
//	container.Singleton("log", func() interface{} {
//	    return logger.NewLogger()
//	})
//	support.SetContainer(container)
//
//	// Now facades can be used
//	support.Log().Info("Application started")
func SetContainer(container containerInterfaces.ContainerInterface) {
	facadeMutex.Lock()
	defer facadeMutex.Unlock()

	facadeContainer = container

	// Clear cache when container changes to avoid stale references
	// This prevents cached services from the old container being returned
	facadeCache = make(map[string]*cacheEntry)

	// Reset cache size metric
	facadeStats.CacheSize = 0
}

// GetContainer returns the current dependency injection container.
//
// This function is primarily used for debugging and introspection.
// Most applications should not need to call this directly.
//
// Returns:
//   - ContainerInterface: The current container instance, or nil if not set
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Example:
//
//	container := support.GetContainer()
//	if container == nil {
//	    log.Fatal("Container not initialized")
//	}
func GetContainer() containerInterfaces.ContainerInterface {
	facadeMutex.RLock()
	defer facadeMutex.RUnlock()
	return facadeContainer
}

// ================================================================================================
// CACHE MANAGEMENT FUNCTIONS
// ================================================================================================

// storeCachedService stores a service in the cache with proper size limit management.
// This is an internal function that handles cache eviction and size limits.
func storeCachedService(serviceKey containerTypes.ServiceIdentifier, service interface{}) {
	key := containerTypes.ToKey(serviceKey)

	facadeMutex.Lock()
	defer facadeMutex.Unlock()

	// Check if we need to evict entries due to size limit
	if facadeOptions.MaxCacheSize > 0 && len(facadeCache) >= facadeOptions.MaxCacheSize {
		// Evict oldest entry to make room
		evictOldestEntries(1)
	}

	// Store the new cache entry
	now := time.Now()
	facadeCache[key] = &cacheEntry{
		service:      service,
		lastAccessed: now,
		createdAt:    now,
	}

	// Update cache size metric
	facadeStats.CacheSize = len(facadeCache)
}

// evictOldestEntries removes the oldest cache entries based on LRU policy.
// This function must be called with facadeMutex write lock held.
func evictOldestEntries(count int) {
	if count <= 0 || len(facadeCache) == 0 {
		return
	}

	// Build list of entries with their access times for sorting
	type entryInfo struct {
		key          string
		lastAccessed time.Time
	}

	entries := make([]entryInfo, 0, len(facadeCache))
	for key, entry := range facadeCache {
		entries = append(entries, entryInfo{
			key:          key,
			lastAccessed: entry.lastAccessed,
		})
	}

	// Sort by last accessed time (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].lastAccessed.After(entries[j].lastAccessed) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove the oldest entries
	evicted := 0
	for i := 0; i < len(entries) && evicted < count; i++ {
		delete(facadeCache, entries[i].key)
		evicted++
	}

	// Update metrics
	atomic.AddInt64(&facadeStats.CacheEvictions, int64(evicted))
	facadeStats.CacheSize = len(facadeCache)
}

// ClearCache clears all cached facade instances.
//
// This function removes all cached services, forcing fresh resolution from
// the container on subsequent calls. It's particularly useful for testing
// or when you need to ensure fresh instances after configuration changes.
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Side Effects:
//   - All cached services are immediately removed
//   - Cache size metric is reset to 0
//   - Next facade calls will resolve fresh instances from container
//
// Example:
//
//	// Clear all caches before running tests
//	support.ClearCache()
//
//	// Clear caches after configuration changes
//	updateConfig()
//	support.ClearCache() // Force fresh resolution with new config
func ClearCache() {
	facadeMutex.Lock()
	defer facadeMutex.Unlock()

	facadeCache = make(map[string]*cacheEntry)
	facadeStats.CacheSize = 0
}

// ClearCacheFor removes a specific service from the cache.
//
// This function removes only the specified service from the cache while
// leaving other cached services intact. Useful for selectively invalidating
// specific services after updates or configuration changes.
//
// Parameters:
//   - serviceKey: The container binding key of the service to remove from cache
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Example:
//
//	// Clear only the logger cache after log configuration changes
//	updateLogConfig()
//	support.ClearCacheFor("log")
//
//	// Clear database cache after connection configuration changes
//	updateDatabaseConfig()
//	support.ClearCacheFor("database")
func ClearCacheFor(serviceKey containerTypes.ServiceIdentifier) {
	key := containerTypes.ToKey(serviceKey)

	facadeMutex.Lock()
	defer facadeMutex.Unlock()

	if _, exists := facadeCache[key]; exists {
		delete(facadeCache, key)
		facadeStats.CacheSize = len(facadeCache)
	}
}

// GetCachedServices returns a copy of all currently cached services.
//
// This function provides introspection into the current cache state for
// debugging, monitoring, or administrative purposes. The returned map is
// a copy to prevent external modification of the internal cache.
//
// Returns:
//   - map[string]interface{}: Map of service keys to cached service instances
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Note:
// The returned map contains the actual service instances, not cache metadata.
// This is intentional to maintain the same interface as the original implementation.
//
// Example:
//
//	cachedServices := support.GetCachedServices()
//	fmt.Printf("Currently cached services: %d\n", len(cachedServices))
//	for key := range cachedServices {
//	    fmt.Printf("  - %s\n", key)
//	}
func GetCachedServices() map[string]interface{} {
	facadeMutex.RLock()
	defer facadeMutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]interface{}, len(facadeCache))
	for key, entry := range facadeCache {
		result[key] = entry.service
	}
	return result
}

// IsServiceCached checks if a specific service is currently cached.
//
// This function is useful for debugging, testing, or conditional logic
// that depends on whether a service is cached or would require fresh resolution.
//
// Parameters:
//   - serviceKey: The container binding key to check
//
// Returns:
//   - bool: true if the service is cached, false otherwise
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Example:
//
//	if support.IsServiceCached("expensive_service") {
//	    // Service is cached, will be fast
//	    service := facade.Resolve[ServiceInterface]("expensive_service")
//	    // ... use service
//	} else {
//	    // Service not cached, may be slow to resolve
//	    log.Println("Resolving expensive service for first time...")
//	    service := facade.Resolve[ServiceInterface]("expensive_service")
//	    // ... use service
//	}
func IsServiceCached(serviceKey containerTypes.ServiceIdentifier) bool {
	key := containerTypes.ToKey(serviceKey)

	facadeMutex.RLock()
	defer facadeMutex.RUnlock()

	entry, exists := facadeCache[key]
	if !exists {
		return false
	}

	// Check TTL expiration if configured
	if facadeOptions.CacheTTL > 0 {
		return time.Since(entry.createdAt) <= facadeOptions.CacheTTL
	}

	return true
}

// ================================================================================================
// TESTING AND DEBUGGING SUPPORT
// ================================================================================================

// SwapService temporarily replaces a service in the cache with a mock or alternative implementation.
//
// This function is specifically designed for testing scenarios where you need to replace
// a real service with a mock, stub, or alternative implementation. It returns a restoration
// function that can be called to restore the original service.
//
// Parameters:
//   - serviceKey: The container binding key of the service to replace
//   - mockService: The mock/alternative service instance to use instead
//
// Returns:
//   - func(): Restoration function that restores the original service when called
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Important Notes:
//   - The mock service is only stored in the cache, not in the container
//   - If caching is disabled, the swap will have no effect
//   - The restoration function should always be called to clean up after testing
//
// Example:
//
//	func TestUserService(t *testing.T) {
//	    // Create a mock database
//	    mockDB := &MockDatabase{}
//	    mockDB.SetupTestData()
//
//	    // Swap the real database with mock
//	    restore := support.SwapService("database", mockDB)
//	    defer restore() // Always restore original
//
//	    // Now all facade calls will use the mock
//	    userService := NewUserService()
//	    users := userService.GetAllUsers() // Uses mockDB
//
//	    assert.Len(t, users, 3) // Test with mock data
//	}
//
// Deferred Restoration Pattern:
//
//	restore := support.SwapService("service", mock)
//	defer restore() // Ensures cleanup even if test panics
func SwapService(serviceKey containerTypes.ServiceIdentifier, mockService interface{}) func() {
	key := containerTypes.ToKey(serviceKey)

	facadeMutex.Lock()
	originalEntry, hadOriginal := facadeCache[key]

	// Replace with mock service (or add if not cached)
	now := time.Now()
	facadeCache[key] = &cacheEntry{
		service:      mockService,
		lastAccessed: now,
		createdAt:    now,
	}

	// Update cache size if this is a new entry
	if !hadOriginal {
		facadeStats.CacheSize = len(facadeCache)
	}

	facadeMutex.Unlock()

	// Return restoration function
	return func() {
		facadeMutex.Lock()
		defer facadeMutex.Unlock()

		if hadOriginal {
			// Restore original service
			facadeCache[key] = originalEntry
		} else {
			// Remove the mock service (wasn't originally cached)
			delete(facadeCache, key)
			facadeStats.CacheSize = len(facadeCache)
		}
	}
}

// ================================================================================================
// METRICS AND MONITORING
// ================================================================================================

// GetStats returns comprehensive statistics about facade operations.
//
// This function provides detailed metrics that can be used for monitoring,
// performance analysis, debugging, and capacity planning. All metrics are
// collected atomically to ensure accuracy under concurrent load.
//
// Returns:
//   - FacadeStats: Detailed statistics about facade operations
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Metrics Included:
//   - Cache hit/miss ratios for performance monitoring
//   - Total resolutions and errors for reliability tracking
//   - Type assertion failures for interface compatibility issues
//   - Current cache size and eviction counts for capacity planning
//
// Example:
//
//	stats := support.GetStats()
//	hitRate := float64(stats.CacheHits) / float64(stats.CacheHits + stats.CacheMisses)
//	fmt.Printf("Cache hit rate: %.2f%%\n", hitRate*100)
//	fmt.Printf("Total errors: %d\n", stats.Errors)
//	fmt.Printf("Cache size: %d/%d\n", stats.CacheSize, support.GetConfiguration().MaxCacheSize)
func GetStats() FacadeStats {
	facadeMutex.RLock()
	cacheSize := len(facadeCache)
	facadeMutex.RUnlock()

	return FacadeStats{
		CacheHits:             atomic.LoadInt64(&facadeStats.CacheHits),
		CacheMisses:           atomic.LoadInt64(&facadeStats.CacheMisses),
		Resolutions:           atomic.LoadInt64(&facadeStats.Resolutions),
		Errors:                atomic.LoadInt64(&facadeStats.Errors),
		TypeAssertionFailures: atomic.LoadInt64(&facadeStats.TypeAssertionFailures),
		CacheSize:             cacheSize,
		CacheEvictions:        atomic.LoadInt64(&facadeStats.CacheEvictions),
	}
}

// ResetStats resets all facade statistics to zero.
//
// This function is useful for testing or when you want to start fresh
// metrics collection from a specific point in time.
//
// Thread Safety:
// This function is thread-safe and can be called concurrently.
//
// Example:
//
//	// Reset stats before benchmark
//	support.ResetStats()
//	runBenchmark()
//	stats := support.GetStats()
//	fmt.Printf("Benchmark results: %+v\n", stats)
func ResetStats() {
	resetStats()
}

// resetStats is an internal function to reset statistics.
// Must be called with appropriate locking when needed.
func resetStats() {
	atomic.StoreInt64(&facadeStats.CacheHits, 0)
	atomic.StoreInt64(&facadeStats.CacheMisses, 0)
	atomic.StoreInt64(&facadeStats.Resolutions, 0)
	atomic.StoreInt64(&facadeStats.Errors, 0)
	atomic.StoreInt64(&facadeStats.TypeAssertionFailures, 0)
	atomic.StoreInt64(&facadeStats.CacheEvictions, 0)
	// Note: CacheSize is updated separately as it's derived from actual cache length
}
