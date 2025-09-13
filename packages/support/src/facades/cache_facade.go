package facades

import (
	cacheInterfaces "govel/types/interfaces/cache"
	facade "govel/support"
)

// Cache provides a clean, static-like interface to the application's caching service.
//
// This facade implements the facade pattern, providing global access to the cache
// service configured in the dependency injection container. It offers a Laravel-style
// API for caching operations with automatic service resolution, type safety, and
// high-performance cache access patterns.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved cache service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent access across goroutines
//   - Supports multiple cache drivers (Redis, Memory, File, etc.)
//
// Behavior:
//   - First call: Resolves cache service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if cache service cannot be resolved (fail-fast behavior)
//   - Automatically handles service lifecycle and connection management
//
// Returns:
//   - CacheInterface: The application's cache service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "cache" service is not registered in the container
//   - If the resolved service doesn't implement CacheInterface
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
//   - Multiple goroutines can call Cache() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Cache updates are atomic and consistent
//
// Usage Examples:
//
//	// Basic cache operations
//	facades.Cache().Set("user:123", userData, 3600) // Cache for 1 hour
//	user := facades.Cache().Get("user:123")
//	facades.Cache().Delete("user:123")
//
//	// Cache with default values
//	value := facades.Cache().GetWithDefault("settings:theme", "dark")
//
//	// Batch operations for efficiency
//	values := facades.Cache().GetMultiple([]string{"key1", "key2", "key3"})
//	facades.Cache().SetMultiple(map[string]interface{}{
//	    "key1": "value1",
//	    "key2": "value2",
//	    "key3": "value3",
//	}, 1800) // 30 minutes TTL
//
//	// Conditional operations
//	if facades.Cache().Has("expensive_computation") {
//	    result = facades.Cache().Get("expensive_computation")
//	} else {
//	    result = performExpensiveComputation()
//	    facades.Cache().Set("expensive_computation", result, 7200) // 2 hours
//	}
//
//	// Cache-aside pattern
//	func GetUser(id int) (*User, error) {
//	    cacheKey := fmt.Sprintf("user:%d", id)
//
//	    // Try cache first
//	    if cached := facades.Cache().Get(cacheKey); cached != nil {
//	        if user, ok := cached.(*User); ok {
//	            return user, nil
//	        }
//	    }
//
//	    // Cache miss - fetch from database
//	    user, err := db.GetUser(id)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Cache for future requests
//	    facades.Cache().Set(cacheKey, user, 1800) // 30 minutes
//	    return user, nil
//	}
//
//	// Write-through caching
//	func UpdateUser(user *User) error {
//	    // Update database first
//	    if err := db.UpdateUser(user); err != nil {
//	        return err
//	    }
//
//	    // Update cache
//	    cacheKey := fmt.Sprintf("user:%d", user.ID)
//	    facades.Cache().Set(cacheKey, user, 1800)
//	    return nil
//	}
//
//	// Cache invalidation patterns
//	func InvalidateUserCache(userID int) {
//	    facades.Cache().Delete(fmt.Sprintf("user:%d", userID))
//	    facades.Cache().DeletePattern(fmt.Sprintf("user:%d:*", userID))
//	    facades.Cache().ClearTag("user_data")
//	}
//
//	// Atomic operations for counters
//	views := facades.Cache().Increment("page:views", 1)
//	likes := facades.Cache().Decrement("post:123:likes", 1)
//
//	// Time-based expiration
//	facades.Cache().SetWithExpiration("session:abc123", sessionData, time.Now().Add(2*time.Hour))
//	facades.Cache().SetUntil("flash_message", "Success!", time.Now().Add(5*time.Minute))
//
//	// Remember pattern (get or set)
//	result := facades.Cache().Remember("expensive_query", 3600, func() interface{} {
//	    return database.RunExpensiveQuery()
//	})
//
// Best Practices:
//   - Use meaningful, hierarchical cache keys ("user:123", "post:456:comments")
//   - Set appropriate TTL values based on data volatility
//   - Consider cache warming for critical data
//   - Implement proper cache invalidation strategies
//   - Use batch operations for multiple keys to reduce network overhead
//   - Monitor cache hit ratios and adjust TTL accordingly
//   - Avoid caching large objects that exceed memory limits
//
// Cache Patterns:
//
//	// 1. Cache-Aside (Lazy Loading)
//	data := facades.Cache().Get(key)
//	if data == nil {
//	    data = loadFromDatabase(key)
//	    facades.Cache().Set(key, data, ttl)
//	}
//
//	// 2. Write-Through
//	facades.Cache().Set(key, data, ttl)
//	database.Save(data)
//
//	// 3. Write-Behind (Write-Back)
//	facades.Cache().Set(key, data, ttl)
//	// Asynchronously write to database later
//
//	// 4. Refresh-Ahead
//	if facades.Cache().TTL(key) < refreshThreshold {
//	    go refreshCache(key) // Async refresh
//	}
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume caching always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	cache, err := facade.TryResolve[CacheInterface]("cache")
//	if err != nil {
//	    // Handle cache unavailability gracefully
//	    return fallbackValue, nil
//	}
//	cache.Set("key", "value", 3600)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserService(t *testing.T) {
//	    // Create a test cache that tracks operations
//	    testCache := &TestCache{}
//
//	    // Swap the real cache with test cache
//	    restore := support.SwapService("cache", testCache)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Cache() returns testCache
//	    userService := NewUserService()
//	    userService.GetUser(123)
//
//	    // Verify caching behavior
//	    assert.True(t, testCache.WasCalled("Get", "user:123"))
//	    assert.True(t, testCache.WasCalled("Set", "user:123"))
//	}
//
// Container Configuration:
// Ensure the cache service is properly configured in your container:
//
//	// Example cache registration
//	container.Singleton("cache", func() interface{} {
//	    config := cache.Config{
//	        Driver:     "redis",           // redis, memory, file, etc.
//	        Connection: "default",        // connection name
//	        Prefix:     "myapp:",          // key prefix
//	        DefaultTTL: 3600,              // 1 hour default
//	        Serializer: "json",           // json, gob, msgpack
//	        Options: map[string]interface{}{
//	            "redis_host": "localhost:6379",
//	            "redis_db":   0,
//	            "pool_size": 10,
//	        },
//	    }
//	    return cache.NewCache(config)
//	})
func Cache() cacheInterfaces.CacheInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves cache service using type-safe token from the dependency injection container
	// - Performs type assertion to CacheInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[cacheInterfaces.CacheInterface](cacheInterfaces.CACHE_TOKEN)
}

// CacheWithError provides error-safe access to the cache service.
//
// This function offers the same functionality as Cache() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle cache unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Cache() but with error handling.
//
// Returns:
//   - CacheInterface: The resolved cache instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement CacheInterface
//
// Usage Examples:
//
//	// Basic error-safe caching
//	cache, err := facades.CacheWithError()
//	if err != nil {
//	    log.Printf("Cache unavailable: %v", err)
//	    return fallbackValue // Use non-cached fallback
//	}
//	cache.Set("key", "value", 3600)
//
//	// Conditional caching
//	if cache, err := facades.CacheWithError(); err == nil {
//	    cache.Set("optional_cache_key", expensiveData, 1800)
//	}
func CacheWithError() (cacheInterfaces.CacheInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves cache service using type-safe token from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[cacheInterfaces.CacheInterface](cacheInterfaces.CACHE_TOKEN)
}
