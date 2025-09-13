// Package types - Context-related type definitions.
// This file defines types for request context and context management.
package webserver

import (
	"context"
	"sync"
	"time"
)

// RequestContext represents the context associated with an HTTP request.
// This provides a way to store and retrieve request-scoped values, timeouts, and cancellation.
type RequestContext struct {
	// ctx is the underlying Go context
	ctx context.Context

	// cancel is the cancellation function for the context
	cancel context.CancelFunc

	// values stores request-scoped key-value pairs
	values map[string]interface{}

	// mutex protects access to values map
	mutex sync.RWMutex

	// startTime is when the request context was created
	startTime time.Time
}

// NewRequestContext creates a new request context.
//
// Returns:
//
//	*RequestContext: A new request context instance
func NewRequestContext() *RequestContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &RequestContext{
		ctx:       ctx,
		cancel:    cancel,
		values:    make(map[string]interface{}),
		startTime: time.Now(),
	}
}

// NewRequestContextWithTimeout creates a new request context with a timeout.
//
// Parameters:
//
//	timeout: The timeout duration for the context
//
// Returns:
//
//	*RequestContext: A new request context with timeout
func NewRequestContextWithTimeout(timeout time.Duration) *RequestContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &RequestContext{
		ctx:       ctx,
		cancel:    cancel,
		values:    make(map[string]interface{}),
		startTime: time.Now(),
	}
}

// NewRequestContextWithDeadline creates a new request context with a deadline.
//
// Parameters:
//
//	deadline: The deadline for the context
//
// Returns:
//
//	*RequestContext: A new request context with deadline
func NewRequestContextWithDeadline(deadline time.Time) *RequestContext {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	return &RequestContext{
		ctx:       ctx,
		cancel:    cancel,
		values:    make(map[string]interface{}),
		startTime: time.Now(),
	}
}

// Context returns the underlying Go context.
//
// Returns:
//
//	context.Context: The underlying context
func (rc *RequestContext) Context() context.Context {
	return rc.ctx
}

// Set stores a value in the request context.
//
// Parameters:
//
//	key: The key to store the value under
//	value: The value to store
func (rc *RequestContext) Set(key string, value interface{}) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.values[key] = value
}

// Get retrieves a value from the request context.
//
// Parameters:
//
//	key: The key to retrieve
//
// Returns:
//
//	interface{}: The stored value, or nil if not found
//	bool: True if the key exists, false otherwise
func (rc *RequestContext) Get(key string) (interface{}, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	value, exists := rc.values[key]
	return value, exists
}

// GetString retrieves a string value from the request context.
//
// Parameters:
//
//	key: The key to retrieve
//
// Returns:
//
//	string: The stored value as string, or empty string if not found/invalid type
//	bool: True if the key exists and is a string, false otherwise
func (rc *RequestContext) GetString(key string) (string, bool) {
	value, exists := rc.Get(key)
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt retrieves an integer value from the request context.
//
// Parameters:
//
//	key: The key to retrieve
//
// Returns:
//
//	int: The stored value as int, or 0 if not found/invalid type
//	bool: True if the key exists and is an int, false otherwise
func (rc *RequestContext) GetInt(key string) (int, bool) {
	value, exists := rc.Get(key)
	if !exists {
		return 0, false
	}
	if i, ok := value.(int); ok {
		return i, true
	}
	return 0, false
}

// GetBool retrieves a boolean value from the request context.
//
// Parameters:
//
//	key: The key to retrieve
//
// Returns:
//
//	bool: The stored value as bool, or false if not found/invalid type
//	bool: True if the key exists and is a bool, false otherwise
func (rc *RequestContext) GetBool(key string) (bool, bool) {
	value, exists := rc.Get(key)
	if !exists {
		return false, false
	}
	if b, ok := value.(bool); ok {
		return b, true
	}
	return false, false
}

// Has checks if a key exists in the request context.
//
// Parameters:
//
//	key: The key to check
//
// Returns:
//
//	bool: True if the key exists, false otherwise
func (rc *RequestContext) Has(key string) bool {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	_, exists := rc.values[key]
	return exists
}

// Delete removes a value from the request context.
//
// Parameters:
//
//	key: The key to remove
func (rc *RequestContext) Delete(key string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	delete(rc.values, key)
}

// Keys returns all keys in the request context.
//
// Returns:
//
//	[]string: All keys in the context
func (rc *RequestContext) Keys() []string {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	keys := make([]string, 0, len(rc.values))
	for key := range rc.values {
		keys = append(keys, key)
	}
	return keys
}

// Clear removes all values from the request context.
func (rc *RequestContext) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.values = make(map[string]interface{})
}

// Cancel cancels the request context.
// This will trigger cancellation for any operations using this context.
func (rc *RequestContext) Cancel() {
	if rc.cancel != nil {
		rc.cancel()
	}
}

// Done returns a channel that is closed when the context is canceled.
//
// Returns:
//
//	<-chan struct{}: A channel that closes when the context is done
func (rc *RequestContext) Done() <-chan struct{} {
	return rc.ctx.Done()
}

// Err returns the error that caused the context to be canceled.
//
// Returns:
//
//	error: The cancellation error, or nil if not canceled
func (rc *RequestContext) Err() error {
	return rc.ctx.Err()
}

// Deadline returns the deadline for the context, if any.
//
// Returns:
//
//	time.Time: The deadline time
//	bool: True if a deadline is set, false otherwise
func (rc *RequestContext) Deadline() (time.Time, bool) {
	return rc.ctx.Deadline()
}

// StartTime returns when the request context was created.
//
// Returns:
//
//	time.Time: The creation time of the context
func (rc *RequestContext) StartTime() time.Time {
	return rc.startTime
}

// Duration returns how long the request context has been active.
//
// Returns:
//
//	time.Duration: The duration since context creation
func (rc *RequestContext) Duration() time.Duration {
	return time.Since(rc.startTime)
}

// Clone creates a copy of the request context with the same values.
// The underlying Go context is not cloned.
//
// Returns:
//
//	*RequestContext: A new context with copied values
func (rc *RequestContext) Clone() *RequestContext {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	newCtx := NewRequestContext()
	for key, value := range rc.values {
		newCtx.values[key] = value
	}

	return newCtx
}
