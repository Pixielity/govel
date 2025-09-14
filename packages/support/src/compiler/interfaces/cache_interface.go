package interfaces

import "govel/support/compiler/types"

// CacheInterface defines the interface for caching compiled results within the GoVel compiler system.
// It provides thread-safe key-value storage with automatic expiration and LRU eviction policies.
type CacheInterface interface {
	// Get retrieves a cached result by its hash key.
	// Returns the cache entry and a boolean indicating if the key exists.
	//
	// Parameters:
	//
	//	key: The unique identifier for the cached entry
	//
	// Returns:
	//
	//	*types.CacheEntry: The cached entry, or nil if not found
	//	bool: true if the key exists in the cache, false otherwise
	Get(key string) (*types.CacheEntry, bool)

	// Set stores a result in the cache with the specified key.
	// Evicts older entries if the cache reaches its size limit.
	//
	// Parameters:
	//
	//	key: The unique identifier for the cache entry
	//	entry: The cache entry to store
	//
	// Returns:
	//
	//	error: An error if the entry cannot be stored
	Set(key string, entry *types.CacheEntry) error

	// Delete removes a cache entry by its key.
	// Returns nil if the key doesn't exist or deletion succeeds.
	//
	// Parameters:
	//
	//	key: The unique identifier of the cache entry to remove
	//
	// Returns:
	//
	//	error: An error if the deletion operation fails
	Delete(key string) error

	// Clear removes all cache entries and resets the cache to an empty state.
	// Useful for cache maintenance, testing, and memory pressure recovery.
	//
	// Returns:
	//
	//	error: An error if the clear operation fails
	Clear() error

	// Size returns the current number of entries in the cache.
	// Includes both active and expired entries that haven't been cleaned up.
	//
	// Returns:
	//
	//	int: The total number of cache entries currently stored
	Size() int
}
