// Package types provides supporting types for the GoVel compiler system.
//
// This file defines cache entries, metrics, and validation result types
// that support the main compiler functionality with comprehensive documentation
// and detailed field descriptions following Go best practices.
package types

import (
	"time"
)

// CacheEntry represents an entry in the compilation cache with metadata.
// Contains cached compilation result and tracking information for cache management.
type CacheEntry struct {
	// Result contains the cached compilation result.
	Result *Result `json:"result"`

	// Hash is the unique identifier for this cache entry.
	Hash string `json:"hash"`

	// Timestamp indicates when this entry was created.
	Timestamp time.Time `json:"timestamp"`

	// AccessCount tracks how many times this entry has been accessed.
	AccessCount int64 `json:"access_count"`

	// LastAccessed indicates when this entry was last accessed.
	LastAccessed time.Time `json:"last_accessed"`

	// Size represents the approximate size of this cache entry in bytes.
	Size int64 `json:"size"`
}

// IsExpired checks if the cache entry has exceeded the given TTL.
// Compares entry creation time against the provided TTL duration.
//
// Parameters:
//
//	ttl: The time-to-live duration for cache entries
//
// Returns:
//
//	bool: true if the entry has expired, false otherwise
func (e *CacheEntry) IsExpired(ttl time.Duration) bool {
	return time.Since(e.Timestamp) > ttl
}

// Touch updates the access count and last accessed time.
// Should be called when the cache entry is accessed for accurate statistics.
func (e *CacheEntry) Touch() {
	e.AccessCount++
	e.LastAccessed = time.Now()
}
