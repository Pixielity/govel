package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"govel/middleware/interfaces"
	"govel/middleware/types"
)

// contextManager manages middleware contexts and their lifecycle
type contextManager struct {
	contexts   map[string]*types.MiddlewareContext[any, any]
	factory    interfaces.ContextFactoryInterface[any, any]
	mu         sync.RWMutex
	logger     interfaces.LoggerInterface
	cleanupTTL time.Duration
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewContextManager creates a new context manager
func NewContextManager(
	factory interfaces.ContextFactoryInterface[any, any],
	logger interfaces.LoggerInterface,
	cleanupTTL time.Duration,
) interfaces.ContextManagerInterface[any, any] {
	cm := &contextManager{
		contexts:   make(map[string]*types.MiddlewareContext[any, any]),
		factory:    factory,
		logger:     logger,
		cleanupTTL: cleanupTTL,
		stopChan:   make(chan struct{}),
	}

	// Start cleanup goroutine if TTL is configured
	if cleanupTTL > 0 {
		cm.startCleanupRoutine()
	}

	return cm
}

// CreateContext creates a new middleware context
func (cm *contextManager) CreateContext(
	ctx context.Context,
	request any,
	metadata map[string]interface{},
) *types.MiddlewareContext[any, any] {
	// Use factory to create context
	mwCtx := cm.factory.CreateMiddlewareContext(ctx, request, metadata)

	// Store context for management if it has an ID
	if mwCtx.ID != "" {
		cm.mu.Lock()
		cm.contexts[mwCtx.ID] = mwCtx
		cm.mu.Unlock()

		cm.logDebug("Created and stored middleware context", map[string]interface{}{
			"context_id":     mwCtx.ID,
			"total_contexts": len(cm.contexts),
		})
	}

	return mwCtx
}

// GetContext retrieves a context by ID
func (cm *contextManager) GetContext(id string) (*types.MiddlewareContext[any, any], bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ctx, exists := cm.contexts[id]
	return ctx, exists
}

// UpdateContext updates context metadata
func (cm *contextManager) UpdateContext(id string, metadata map[string]interface{}) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ctx, exists := cm.contexts[id]
	if !exists {
		return fmt.Errorf("context with ID %s not found", id)
	}

	// Update metadata
	for key, value := range metadata {
		ctx.Metadata[key] = value
	}

	cm.logDebug("Updated context metadata", map[string]interface{}{
		"context_id":    id,
		"metadata_keys": len(metadata),
	})

	return nil
}

// RemoveContext removes a context from management
func (cm *contextManager) RemoveContext(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.contexts[id]; exists {
		delete(cm.contexts, id)
		cm.logDebug("Removed context from management", map[string]interface{}{
			"context_id":         id,
			"remaining_contexts": len(cm.contexts),
		})
	}
}

// GetActiveContexts returns the number of active contexts
func (cm *contextManager) GetActiveContexts() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.contexts)
}

// CleanupExpiredContexts removes expired contexts
func (cm *contextManager) CleanupExpiredContexts() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	expired := make([]string, 0)

	for id, ctx := range cm.contexts {
		// Check if context is expired based on creation time + TTL
		if now.Sub(ctx.CreatedAt) > cm.cleanupTTL {
			expired = append(expired, id)
		}

		// Also check if the underlying Go context is done
		select {
		case <-ctx.Context.Done():
			expired = append(expired, id)
		default:
		}
	}

	// Remove expired contexts
	for _, id := range expired {
		delete(cm.contexts, id)
	}

	if len(expired) > 0 {
		cm.logDebug("Cleaned up expired contexts", map[string]interface{}{
			"expired_count":      len(expired),
			"remaining_contexts": len(cm.contexts),
		})
	}

	return len(expired)
}

// CreateRequestContext creates a request-specific context
func (cm *contextManager) CreateRequestContext(
	ctx context.Context,
	request any,
	metadata map[string]interface{},
) *types.RequestContext[any] {
	return cm.factory.CreateRequestContext(ctx, request, metadata)
}

// CreateExecutionContext creates an execution-specific context
func (cm *contextManager) CreateExecutionContext(
	ctx context.Context,
	config *types.ExecutionConfig,
) *types.ExecutionContext {
	return cm.factory.CreateExecutionContext(ctx, config)
}

// WithTimeout creates a context with timeout
func (cm *contextManager) WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// WithCancel creates a context with cancellation
func (cm *contextManager) WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

// WithDeadline creates a context with deadline
func (cm *contextManager) WithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, deadline)
}

// WithValue creates a context with a key-value pair
func (cm *contextManager) WithValue(parent context.Context, key, value interface{}) context.Context {
	return context.WithValue(parent, key, value)
}

// Shutdown gracefully shuts down the context manager
func (cm *contextManager) Shutdown() error {
	// Stop cleanup routine
	if cm.cleanupTTL > 0 {
		close(cm.stopChan)
		cm.wg.Wait()
	}

	// Clear all contexts
	cm.mu.Lock()
	clearedCount := len(cm.contexts)
	cm.contexts = make(map[string]*types.MiddlewareContext[any, any])
	cm.mu.Unlock()

	cm.logDebug("Context manager shutdown", map[string]interface{}{
		"cleared_contexts": clearedCount,
	})

	return nil
}

// GetStats returns context manager statistics
func (cm *contextManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return map[string]interface{}{
		"active_contexts": len(cm.contexts),
		"cleanup_ttl":     cm.cleanupTTL,
		"cleanup_enabled": cm.cleanupTTL > 0,
	}
}

// startCleanupRoutine starts the background cleanup routine
func (cm *contextManager) startCleanupRoutine() {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()

		ticker := time.NewTicker(cm.cleanupTTL / 2) // Cleanup twice as often as TTL
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cm.CleanupExpiredContexts()
			case <-cm.stopChan:
				return
			}
		}
	}()

	cm.logDebug("Started context cleanup routine", map[string]interface{}{
		"cleanup_ttl": cm.cleanupTTL,
	})
}

// Helper logging methods
func (cm *contextManager) logDebug(message string, context map[string]interface{}) {
	if cm.logger != nil {
		cm.logger.Debug(message, context)
	}
}

// pooledContextManager extends contextManager with object pooling for performance
type pooledContextManager struct {
	*contextManager
	contextPool   *sync.Pool
	requestPool   *sync.Pool
	executionPool *sync.Pool
}

// NewPooledContextManager creates a context manager with object pooling
func NewPooledContextManager(
	factory interfaces.ContextFactoryInterface[any, any],
	logger interfaces.LoggerInterface,
	cleanupTTL time.Duration,
) interfaces.ContextManagerInterface[any, any] {
	cm := &pooledContextManager{
		contextManager: NewContextManager(factory, logger, cleanupTTL).(*contextManager),
	}

	// Initialize pools
	cm.contextPool = &sync.Pool{
		New: func() interface{} {
			return &types.MiddlewareContext[any, any]{
				Metadata: make(map[string]interface{}),
			}
		},
	}

	cm.requestPool = &sync.Pool{
		New: func() interface{} {
			return &types.RequestContext[any]{
				Metadata: make(map[string]interface{}),
			}
		},
	}

	cm.executionPool = &sync.Pool{
		New: func() interface{} {
			return &types.ExecutionContext{}
		},
	}

	return cm
}

// CreateContext creates a context using object pooling
func (pcm *pooledContextManager) CreateContext(
	ctx context.Context,
	request any,
	metadata map[string]interface{},
) *types.MiddlewareContext[any, any] {
	// Get from pool
	mwCtx := pcm.contextPool.Get().(*types.MiddlewareContext[any, any])

	// Reset and initialize
	mwCtx.Context = ctx
	mwCtx.Request = request
	mwCtx.CreatedAt = time.Now()
	mwCtx.ID = fmt.Sprintf("ctx_%d", time.Now().UnixNano())

	// Clear and set metadata
	for k := range mwCtx.Metadata {
		delete(mwCtx.Metadata, k)
	}
	for k, v := range metadata {
		mwCtx.Metadata[k] = v
	}

	// Store for management
	pcm.mu.Lock()
	pcm.contexts[mwCtx.ID] = mwCtx
	pcm.mu.Unlock()

	return mwCtx
}

// RemoveContext removes and returns context to pool
func (pcm *pooledContextManager) RemoveContext(id string) {
	pcm.mu.Lock()
	mwCtx, exists := pcm.contexts[id]
	if exists {
		delete(pcm.contexts, id)
	}
	pcm.mu.Unlock()

	if exists {
		// Return to pool
		pcm.contextPool.Put(mwCtx)

		pcm.logDebug("Returned context to pool", map[string]interface{}{
			"context_id": id,
		})
	}
}

// CreateRequestContext creates a request context using pooling
func (pcm *pooledContextManager) CreateRequestContext(
	ctx context.Context,
	request any,
	metadata map[string]interface{},
) *types.RequestContext[any] {
	// Get from pool
	reqCtx := pcm.requestPool.Get().(*types.RequestContext[any])

	// Reset and initialize
	reqCtx.Context = ctx
	reqCtx.Request = request
	reqCtx.CreatedAt = time.Now()

	// Clear and set metadata
	for k := range reqCtx.Metadata {
		delete(reqCtx.Metadata, k)
	}
	for k, v := range metadata {
		reqCtx.Metadata[k] = v
	}

	return reqCtx
}

// CreateExecutionContext creates an execution context using pooling
func (pcm *pooledContextManager) CreateExecutionContext(
	ctx context.Context,
	config *types.ExecutionConfig,
) *types.ExecutionContext {
	// Get from pool
	execCtx := pcm.executionPool.Get().(*types.ExecutionContext)

	// Reset and initialize
	execCtx = types.NewExecutionContext(ctx, config)

	return execCtx
}

// ReturnToPool returns contexts back to their respective pools
func (pcm *pooledContextManager) ReturnToPool(ctx interface{}) {
	switch c := ctx.(type) {
	case *types.MiddlewareContext[any, any]:
		pcm.contextPool.Put(c)
	case *types.RequestContext[any]:
		pcm.requestPool.Put(c)
	case *types.ExecutionContext:
		pcm.executionPool.Put(c)
	}
}
