// Package pipeline provides pipeline context implementation for Go.
package pipeline

import (
	"context"
	"sync"
	"time"

	interfaces "govel/packages/types/src/interfaces/pipeline"
)

// PipelineContext implements the PipelineContextInterface interface.
// It extends Go's standard context.Context with pipeline-specific functionality
// for managing pipeline state, cancellation, and metadata.
//
// This implementation is thread-safe and can be used across multiple goroutines.
type PipelineContext struct {
	// ctx is the underlying Go context
	ctx context.Context

	// metadata stores pipeline-specific metadata
	metadata map[string]interface{}

	// pipelineID is a unique identifier for the pipeline execution
	pipelineID string

	// currentPipe is the name of the currently executing pipe
	currentPipe string

	// startTime is when the pipeline execution started
	startTime time.Time

	// mutex protects concurrent access to mutable fields
	mutex sync.RWMutex
}

// NewPipelineContext creates a new PipelineContext wrapping the given context.
// The provided context can be any Go context (Background, WithTimeout, etc.).
//
// Parameters:
//   - ctx: The underlying Go context to wrap
//
// Returns:
//   - interfaces.PipelineContextInterface: New pipeline context instance
//
// Example:
//
//	ctx := NewPipelineContext(context.Background())
//	ctx.SetMetadata("request_id", "abc123")
//
//	// Use with timeout
//	ctxWithTimeout := NewPipelineContext(context.WithTimeout(context.Background(), 30*time.Second))
func NewPipelineContext(ctx context.Context) interfaces.PipelineContextInterface {
	if ctx == nil {
		ctx = context.Background()
	}

	return &PipelineContext{
		ctx:      ctx,
		metadata: make(map[string]interface{}),
		mutex:    sync.RWMutex{},
	}
}

// Deadline returns the deadline for the context, if any.
// This is part of the context.Context interface.
func (pc *PipelineContext) Deadline() (deadline time.Time, ok bool) {
	return pc.ctx.Deadline()
}

// Done returns a channel that's closed when the context is cancelled.
// This is part of the context.Context interface.
func (pc *PipelineContext) Done() <-chan struct{} {
	return pc.ctx.Done()
}

// Err returns the error that caused the context to be cancelled.
// This is part of the context.Context interface.
func (pc *PipelineContext) Err() error {
	return pc.ctx.Err()
}

// Value returns the value associated with the given key in the context.
// This is part of the context.Context interface.
// It first checks pipeline metadata, then delegates to the underlying context.
func (pc *PipelineContext) Value(key interface{}) interface{} {
	// Check if the key is a string and exists in our metadata
	if strKey, ok := key.(string); ok {
		if value, exists := pc.GetMetadata(strKey); exists {
			return value
		}
	}

	// Delegate to the underlying context
	return pc.ctx.Value(key)
}

// WithTimeout creates a new context with the specified timeout.
// If the pipeline execution exceeds this timeout, it will be cancelled.
//
// Parameters:
//   - timeout: Maximum duration for pipeline execution
//
// Returns:
//   - interfaces.PipelineContextInterface: New context with timeout applied
func (pc *PipelineContext) WithTimeout(timeout time.Duration) interfaces.PipelineContextInterface {
	ctx, _ := context.WithTimeout(pc.ctx, timeout)
	newPipelineCtx := NewPipelineContext(ctx).(*PipelineContext)

	// Copy existing metadata and state
	pc.mutex.RLock()
	newPipelineCtx.metadata = pc.copyMetadata()
	newPipelineCtx.pipelineID = pc.pipelineID
	newPipelineCtx.currentPipe = pc.currentPipe
	newPipelineCtx.startTime = pc.startTime
	pc.mutex.RUnlock()

	return newPipelineCtx
}

// WithDeadline creates a new context with the specified deadline.
// If the pipeline execution hasn't completed by the deadline, it will be cancelled.
//
// Parameters:
//   - deadline: Absolute time when the context should be cancelled
//
// Returns:
//   - interfaces.PipelineContextInterface: New context with deadline applied
func (pc *PipelineContext) WithDeadline(deadline time.Time) interfaces.PipelineContextInterface {
	ctx, _ := context.WithDeadline(pc.ctx, deadline)
	newPipelineCtx := NewPipelineContext(ctx).(*PipelineContext)

	// Copy existing metadata and state
	pc.mutex.RLock()
	newPipelineCtx.metadata = pc.copyMetadata()
	newPipelineCtx.pipelineID = pc.pipelineID
	newPipelineCtx.currentPipe = pc.currentPipe
	newPipelineCtx.startTime = pc.startTime
	pc.mutex.RUnlock()

	return newPipelineCtx
}

// WithCancel creates a new context that can be cancelled manually.
// This returns both the context and a cancel function.
//
// Returns:
//   - interfaces.PipelineContextInterface: New cancellable context
//   - context.CancelFunc: Function to cancel the context
func (pc *PipelineContext) WithCancel() (interfaces.PipelineContextInterface, context.CancelFunc) {
	ctx, cancel := context.WithCancel(pc.ctx)
	newPipelineCtx := NewPipelineContext(ctx).(*PipelineContext)

	// Copy existing metadata and state
	pc.mutex.RLock()
	newPipelineCtx.metadata = pc.copyMetadata()
	newPipelineCtx.pipelineID = pc.pipelineID
	newPipelineCtx.currentPipe = pc.currentPipe
	newPipelineCtx.startTime = pc.startTime
	pc.mutex.RUnlock()

	return newPipelineCtx, cancel
}

// SetMetadata stores a key-value pair in the context.
// This metadata can be accessed by pipes during execution.
//
// Parameters:
//   - key: The metadata key
//   - value: The metadata value
//
// Returns:
//   - interfaces.PipelineContextInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) SetMetadata(key string, value interface{}) interfaces.PipelineContextInterface {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.metadata[key] = value
	return pc
}

// GetMetadata retrieves a metadata value by key.
// If the key doesn't exist, returns nil and false.
//
// Parameters:
//   - key: The metadata key to retrieve
//
// Returns:
//   - interface{}: The metadata value (nil if not found)
//   - bool: True if the key exists, false otherwise
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) GetMetadata(key string) (interface{}, bool) {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	value, exists := pc.metadata[key]
	return value, exists
}

// GetAllMetadata returns all metadata as a map.
// This is useful for debugging or passing context to external systems.
//
// Returns:
//   - map[string]interface{}: All metadata key-value pairs (copy)
//
// Thread-safe: This method returns a copy of the metadata to prevent race conditions.
func (pc *PipelineContext) GetAllMetadata() map[string]interface{} {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return pc.copyMetadata()
}

// Clone creates a copy of the current context with the same metadata.
// This is useful when you need to branch context for parallel processing.
//
// Returns:
//   - interfaces.PipelineContextInterface: New context with copied metadata
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) Clone() interfaces.PipelineContextInterface {
	newPipelineCtx := NewPipelineContext(pc.ctx).(*PipelineContext)

	// Copy all state
	pc.mutex.RLock()
	newPipelineCtx.metadata = pc.copyMetadata()
	newPipelineCtx.pipelineID = pc.pipelineID
	newPipelineCtx.currentPipe = pc.currentPipe
	newPipelineCtx.startTime = pc.startTime
	pc.mutex.RUnlock()

	return newPipelineCtx
}

// SetPipelineID sets a unique identifier for the current pipeline execution.
// This is useful for tracing and logging pipeline execution.
//
// Parameters:
//   - id: Unique identifier for the pipeline execution
//
// Returns:
//   - interfaces.PipelineContextInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) SetPipelineID(id string) interfaces.PipelineContextInterface {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.pipelineID = id
	return pc
}

// GetPipelineID returns the pipeline execution identifier.
// If no ID has been set, returns an empty string.
//
// Returns:
//   - string: The pipeline execution identifier
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) GetPipelineID() string {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return pc.pipelineID
}

// SetCurrentPipe sets the name of the currently executing pipe.
// This is useful for debugging and error reporting.
//
// Parameters:
//   - pipeName: Name of the currently executing pipe
//
// Returns:
//   - interfaces.PipelineContextInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) SetCurrentPipe(pipeName string) interfaces.PipelineContextInterface {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.currentPipe = pipeName
	return pc
}

// GetCurrentPipe returns the name of the currently executing pipe.
// If no pipe name has been set, returns an empty string.
//
// Returns:
//   - string: Name of the currently executing pipe
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) GetCurrentPipe() string {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return pc.currentPipe
}

// SetExecutionStartTime records when the pipeline execution started.
// This is useful for performance monitoring and timeout calculations.
//
// Parameters:
//   - startTime: When the pipeline execution started
//
// Returns:
//   - interfaces.PipelineContextInterface: Returns self for method chaining
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) SetExecutionStartTime(startTime time.Time) interfaces.PipelineContextInterface {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.startTime = startTime
	return pc
}

// GetExecutionStartTime returns when the pipeline execution started.
// If no start time has been set, returns zero time.
//
// Returns:
//   - time.Time: When the pipeline execution started
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) GetExecutionStartTime() time.Time {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return pc.startTime
}

// GetExecutionDuration returns how long the pipeline has been executing.
// This calculates the duration from the start time to now.
//
// Returns:
//   - time.Duration: How long the pipeline has been executing
//
// Thread-safe: This method is safe for concurrent use.
func (pc *PipelineContext) GetExecutionDuration() time.Duration {
	pc.mutex.RLock()
	startTime := pc.startTime
	pc.mutex.RUnlock()

	if startTime.IsZero() {
		return 0
	}

	return time.Since(startTime)
}

// IsTimedOut checks if the context has timed out.
// This is a convenience method that checks if Err() returns context.DeadlineExceeded.
//
// Returns:
//   - bool: True if the context has timed out
func (pc *PipelineContext) IsTimedOut() bool {
	return pc.ctx.Err() == context.DeadlineExceeded
}

// IsCancelled checks if the context has been cancelled.
// This is a convenience method that checks if Err() returns context.Canceled.
//
// Returns:
//   - bool: True if the context has been cancelled
func (pc *PipelineContext) IsCancelled() bool {
	return pc.ctx.Err() == context.Canceled
}

// copyMetadata creates a deep copy of the metadata map.
// This is used internally to ensure thread safety when sharing metadata.
func (pc *PipelineContext) copyMetadata() map[string]interface{} {
	copy := make(map[string]interface{}, len(pc.metadata))
	for k, v := range pc.metadata {
		copy[k] = v
	}
	return copy
}
