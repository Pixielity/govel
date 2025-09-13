package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/govel/middleware/src/interfaces"
	"github.com/govel/middleware/src/types"
)

// executionChain implements the middleware execution chain using the Russian Doll pattern
type executionChain struct {
	middlewares []interfaces.MiddlewareInterface[any, any]
	config      *types.ExecutionConfig
	metrics     *types.ExecutionMetrics
	mu          sync.RWMutex
	logger      interfaces.LoggerInterface
}

// NewExecutionChain creates a new middleware execution chain
func NewExecutionChain(config *types.ExecutionConfig) interfaces.ExecutionChainInterface[any, any] {
	return &executionChain{
		middlewares: make([]interfaces.MiddlewareInterface[any, any], 0),
		config:      config,
		metrics:     types.NewExecutionMetrics(),
		logger:      config.Logger,
	}
}

// AddMiddleware adds a middleware to the chain
func (ec *executionChain) AddMiddleware(middleware interfaces.MiddlewareInterface[any, any]) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	ec.middlewares = append(ec.middlewares, middleware)
	ec.logDebug("Added middleware to chain", map[string]interface{}{
		"middleware_count": len(ec.middlewares),
	})
}

// PrependMiddleware adds a middleware to the beginning of the chain
func (ec *executionChain) PrependMiddleware(middleware interfaces.MiddlewareInterface[any, any]) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	ec.middlewares = append([]interfaces.MiddlewareInterface[any, any]{middleware}, ec.middlewares...)
	ec.logDebug("Prepended middleware to chain", map[string]interface{}{
		"middleware_count": len(ec.middlewares),
	})
}

// Execute runs the middleware chain with the given context and request
func (ec *executionChain) Execute(ctx context.Context, request any) (any, error) {
	startTime := time.Now()
	
	// Update metrics
	ec.metrics.IncrementExecutions()
	
	defer func() {
		duration := time.Since(startTime)
		ec.metrics.AddExecutionTime(duration)
		
		if r := recover(); r != nil {
			ec.metrics.IncrementPanics()
			ec.logError("Panic during execution chain", map[string]interface{}{
				"panic": r,
				"duration": duration,
			})
			panic(r)
		}
	}()
	
	// Check if chain is empty
	if len(ec.middlewares) == 0 {
		ec.logDebug("Empty middleware chain, returning request as response", nil)
		return request, nil
	}
	
	// Apply timeout if configured
	if ec.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ec.config.Timeout)
		defer cancel()
	}
	
	// Create execution context
	execCtx := types.NewExecutionContext(ctx, ec.config)
	
	// Execute the chain using the Russian Doll pattern
	response, err := ec.executeChain(execCtx, request, 0)
	
	if err != nil {
		ec.metrics.IncrementErrors()
		ec.logError("Error in execution chain", map[string]interface{}{
			"error": err.Error(),
			"duration": time.Since(startTime),
		})
		return nil, err
	}
	
	ec.logDebug("Successfully executed middleware chain", map[string]interface{}{
		"middleware_count": len(ec.middlewares),
		"duration": time.Since(startTime),
	})
	
	return response, nil
}

// executeChain recursively executes the middleware chain
func (ec *executionChain) executeChain(ctx *types.ExecutionContext, request any, index int) (any, error) {
	// Check for context cancellation
	if err := ctx.GetContext().Err(); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}
	
	// If we've reached the end of the chain, return the request as response
	if index >= len(ec.middlewares) {
		return request, nil
	}
	
	// Get the current middleware
	middleware := ec.middlewares[index]
	
	// Create a next function for the current middleware
	next := func(nextCtx context.Context, nextRequest any) (any, error) {
		// Update the context if it changed
		if nextCtx != ctx.GetContext() {
			ctx = ctx.WithContext(nextCtx)
		}
		return ec.executeChain(ctx, nextRequest, index+1)
	}
	
	// Execute the middleware
	return ec.executeMiddleware(ctx, middleware, request, next)
}

// executeMiddleware executes a single middleware with error handling and retry logic
func (ec *executionChain) executeMiddleware(
	ctx *types.ExecutionContext,
	middleware interfaces.MiddlewareInterface[any, any],
	request any,
	next types.NextFunc[any, any],
) (any, error) {
	var lastErr error
	maxRetries := ec.config.RetryPolicy.MaxRetries
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check circuit breaker
		if ec.config.CircuitBreaker.IsOpen() {
			return nil, errors.New("circuit breaker is open")
		}
		
		// Add delay for retry attempts
		if attempt > 0 {
			delay := ec.config.RetryPolicy.BackoffStrategy(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.GetContext().Done():
				return nil, ctx.GetContext().Err()
			}
		}
		
		// Execute middleware with panic recovery
		response, err := ec.executeWithRecovery(ctx, middleware, request, next)
		
		if err == nil {
			// Success - record in circuit breaker
			ec.config.CircuitBreaker.RecordSuccess()
			return response, nil
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !ec.isRetryableError(err) {
			break
		}
		
		// Record failure in circuit breaker
		ec.config.CircuitBreaker.RecordFailure()
		
		ec.logWarning("Middleware execution failed, retrying", map[string]interface{}{
			"attempt": attempt + 1,
			"max_retries": maxRetries,
			"error": err.Error(),
		})
	}
	
	return nil, fmt.Errorf("middleware execution failed after %d attempts: %w", maxRetries+1, lastErr)
}

// executeWithRecovery executes middleware with panic recovery
func (ec *executionChain) executeWithRecovery(
	ctx *types.ExecutionContext,
	middleware interfaces.MiddlewareInterface[any, any],
	request any,
	next types.NextFunc[any, any],
) (response any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in middleware: %v", r)
			ec.logError("Panic recovered in middleware", map[string]interface{}{
				"panic": r,
				"middleware": fmt.Sprintf("%T", middleware),
			})
		}
	}()
	
	return middleware.Handle(ctx.GetContext(), request, next)
}

// isRetryableError determines if an error is retryable based on configuration
func (ec *executionChain) isRetryableError(err error) bool {
	if ec.config.RetryPolicy.RetryableErrors == nil {
		return true // Retry all errors by default
	}
	
	for _, retryableErr := range ec.config.RetryPolicy.RetryableErrors {
		if errors.Is(err, retryableErr) {
			return true
		}
	}
	
	return false
}

// GetMetrics returns the current execution metrics
func (ec *executionChain) GetMetrics() *types.ExecutionMetrics {
	return ec.metrics
}

// GetConfig returns the current execution configuration
func (ec *executionChain) GetConfig() *types.ExecutionConfig {
	return ec.config
}

// UpdateConfig updates the execution configuration
func (ec *executionChain) UpdateConfig(config *types.ExecutionConfig) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	ec.config = config
	ec.logger = config.Logger
	
	ec.logDebug("Updated execution chain configuration", nil)
}

// Clear removes all middlewares from the chain
func (ec *executionChain) Clear() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	
	ec.middlewares = ec.middlewares[:0]
	ec.metrics.Reset()
	
	ec.logDebug("Cleared middleware chain", nil)
}

// Count returns the number of middlewares in the chain
func (ec *executionChain) Count() int {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	
	return len(ec.middlewares)
}

// GetMiddlewares returns a copy of the current middleware slice
func (ec *executionChain) GetMiddlewares() []interfaces.MiddlewareInterface[any, any] {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	
	// Return a copy to prevent external modification
	middlewares := make([]interfaces.MiddlewareInterface[any, any], len(ec.middlewares))
	copy(middlewares, ec.middlewares)
	
	return middlewares
}

// Helper logging methods
func (ec *executionChain) logDebug(message string, context map[string]interface{}) {
	if ec.logger != nil {
		ec.logger.Debug(message, context)
	}
}

func (ec *executionChain) logWarning(message string, context map[string]interface{}) {
	if ec.logger != nil {
		ec.logger.Warning(message, context)
	}
}

func (ec *executionChain) logError(message string, context map[string]interface{}) {
	if ec.logger != nil {
		ec.logger.Error(message, context)
	}
}

// ConditionalExecutionChain extends execution chain with conditional logic
type conditionalExecutionChain struct {
	*executionChain
	condition types.ConditionFunc[any]
}

// NewConditionalExecutionChain creates a new conditional execution chain
func NewConditionalExecutionChain(
	config *types.ExecutionConfig,
	condition types.ConditionFunc[any],
) interfaces.ExecutionChainInterface[any, any] {
	return &conditionalExecutionChain{
		executionChain: NewExecutionChain(config).(*executionChain),
		condition:      condition,
	}
}

// Execute runs the chain only if the condition is met
func (cec *conditionalExecutionChain) Execute(ctx context.Context, request any) (any, error) {
	if !cec.condition(ctx, request) {
		cec.logDebug("Condition not met, skipping middleware chain", nil)
		return request, nil
	}
	
	return cec.executionChain.Execute(ctx, request)
}

// ParallelExecutionChain executes multiple chains in parallel
type parallelExecutionChain struct {
	chains []interfaces.ExecutionChainInterface[any, any]
	config *types.ExecutionConfig
	logger interfaces.LoggerInterface
}

// NewParallelExecutionChain creates a new parallel execution chain
func NewParallelExecutionChain(
	config *types.ExecutionConfig,
	chains ...interfaces.ExecutionChainInterface[any, any],
) interfaces.ExecutionChainInterface[any, any] {
	return &parallelExecutionChain{
		chains: chains,
		config: config,
		logger: config.Logger,
	}
}

// Execute runs all chains in parallel and aggregates results
func (pec *parallelExecutionChain) Execute(ctx context.Context, request any) (any, error) {
	if len(pec.chains) == 0 {
		return request, nil
	}
	
	if len(pec.chains) == 1 {
		return pec.chains[0].Execute(ctx, request)
	}
	
	type result struct {
		response any
		err      error
		index    int
	}
	
	results := make(chan result, len(pec.chains))
	
	// Execute all chains in parallel
	for i, chain := range pec.chains {
		go func(idx int, c interfaces.ExecutionChainInterface[any, any]) {
			defer func() {
				if r := recover(); r != nil {
					results <- result{
						response: nil,
						err:      fmt.Errorf("panic in parallel chain %d: %v", idx, r),
						index:    idx,
					}
				}
			}()
			
			response, err := c.Execute(ctx, request)
			results <- result{
				response: response,
				err:      err,
				index:    idx,
			}
		}(i, chain)
	}
	
	// Collect results
	responses := make([]any, len(pec.chains))
	var firstError error
	
	for i := 0; i < len(pec.chains); i++ {
		select {
		case res := <-results:
			responses[res.index] = res.response
			if res.err != nil && firstError == nil {
				firstError = res.err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	if firstError != nil {
		return nil, firstError
	}
	
	// Return the first non-nil response or the original request
	for _, response := range responses {
		if response != nil {
			return response, nil
		}
	}
	
	return request, nil
}

// Implement other required methods for parallel execution chain
func (pec *parallelExecutionChain) AddMiddleware(middleware interfaces.MiddlewareInterface[any, any]) {
	// Add to all chains
	for _, chain := range pec.chains {
		chain.AddMiddleware(middleware)
	}
}

func (pec *parallelExecutionChain) PrependMiddleware(middleware interfaces.MiddlewareInterface[any, any]) {
	// Prepend to all chains
	for _, chain := range pec.chains {
		chain.PrependMiddleware(middleware)
	}
}

func (pec *parallelExecutionChain) GetMetrics() *types.ExecutionMetrics {
	// Aggregate metrics from all chains
	aggregated := types.NewExecutionMetrics()
	
	for _, chain := range pec.chains {
		metrics := chain.GetMetrics()
		aggregated.Executions += metrics.Executions
		aggregated.Errors += metrics.Errors
		aggregated.Panics += metrics.Panics
		aggregated.TotalDuration += metrics.TotalDuration
	}
	
	return aggregated
}

func (pec *parallelExecutionChain) GetConfig() *types.ExecutionConfig {
	return pec.config
}

func (pec *parallelExecutionChain) UpdateConfig(config *types.ExecutionConfig) {
	pec.config = config
	pec.logger = config.Logger
	
	// Update all chains
	for _, chain := range pec.chains {
		chain.UpdateConfig(config)
	}
}

func (pec *parallelExecutionChain) Clear() {
	for _, chain := range pec.chains {
		chain.Clear()
	}
}

func (pec *parallelExecutionChain) Count() int {
	total := 0
	for _, chain := range pec.chains {
		total += chain.Count()
	}
	return total
}

func (pec *parallelExecutionChain) GetMiddlewares() []interfaces.MiddlewareInterface[any, any] {
	var all []interfaces.MiddlewareInterface[any, any]
	for _, chain := range pec.chains {
		all = append(all, chain.GetMiddlewares()...)
	}
	return all
}
