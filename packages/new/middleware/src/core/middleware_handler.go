// Package core provides concrete implementations of the middleware system components.
// This file contains the core middleware handler implementation that manages
// middleware execution, request processing, and response handling with full
// support for performance monitoring and error management.
package core

import (
	"context"
	"fmt"
	"govel/middleware/interfaces"
	"govel/middleware/types"
	"sync"
	"time"
)

// MiddlewareHandler is the primary implementation for handling middleware execution.
// It provides a complete middleware management system with support for dynamic
// middleware registration, execution tracking, performance monitoring, and
// sophisticated error handling patterns.
//
// Key Features:
// - Thread-safe middleware registration and execution
// - Performance metrics collection and reporting
// - Configurable timeout and concurrency management
// - Circuit breaker pattern for fault tolerance
// - Distributed tracing and logging integration
// - Pipeline integration for advanced workflows
//
// Type Parameters:
//   - TRequest: The type of request data processed by this handler
//   - TResponse: The type of response data returned by this handler
type MiddlewareHandler[TRequest, TResponse any] struct {
	// middlewares holds the registered middleware in execution order
	middlewares []interfaces.Middleware[TRequest, TResponse]

	// config contains the handler configuration
	config *MiddlewareHandlerConfig

	// metrics tracks execution performance and health
	metrics *HandlerMetrics

	// contextFactory creates execution contexts
	contextFactory interfaces.ContextFactory

	// executionChain manages actual middleware execution
	executionChain *types.ExecutionChain[TRequest, TResponse]

	// mu protects concurrent access to handler state
	mu sync.RWMutex
}

// MiddlewareHandlerConfig contains configuration options for the middleware handler.
type MiddlewareHandlerConfig struct {
	// Name is a human-readable identifier for this handler
	Name string `json:"name" yaml:"name"`

	// EnableMetrics indicates whether to collect performance metrics
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`

	// EnableTracing indicates whether to enable distributed tracing
	EnableTracing bool `json:"enable_tracing" yaml:"enable_tracing"`

	// EnableProfiling indicates whether to collect detailed performance profiles
	EnableProfiling bool `json:"enable_profiling" yaml:"enable_profiling"`

	// DefaultTimeout specifies the default timeout for middleware execution
	DefaultTimeout time.Duration `json:"default_timeout" yaml:"default_timeout"`

	// MaxConcurrency limits the number of concurrent middleware executions
	MaxConcurrency int `json:"max_concurrency" yaml:"max_concurrency"`

	// CircuitBreakerConfig configures circuit breaker behavior
	CircuitBreakerConfig *types.CircuitBreakerConfig `json:"circuit_breaker" yaml:"circuit_breaker"`

	// RetryPolicy configures automatic retry behavior
	RetryPolicy *types.RetryPolicy `json:"retry_policy" yaml:"retry_policy"`

	// Properties contains additional configuration properties
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
}

// HandlerMetrics tracks detailed performance and health metrics for the handler.
type HandlerMetrics struct {
	// totalRequests is the total number of requests processed
	totalRequests int64

	// successfulRequests is the number of requests that completed successfully
	successfulRequests int64

	// failedRequests is the number of requests that failed
	failedRequests int64

	// totalProcessingTime is the cumulative time spent processing requests
	totalProcessingTime time.Duration

	// averageProcessingTime is the average time to process a request
	averageProcessingTime time.Duration

	// minProcessingTime is the minimum time to process a request
	minProcessingTime time.Duration

	// maxProcessingTime is the maximum time to process a request
	maxProcessingTime time.Duration

	// currentConcurrency is the current number of concurrent executions
	currentConcurrency int64

	// maxConcurrency is the peak number of concurrent executions
	maxConcurrency int64

	// lastRequestTime is the timestamp of the last processed request
	lastRequestTime time.Time

	// mu protects concurrent access to metrics
	mu sync.RWMutex
}

// NewMiddlewareHandler creates a new middleware handler with the specified configuration.
// This function initializes a complete middleware management system with all the
// necessary components for robust middleware execution.
//
// Type Parameters:
//   - TRequest: The type of request data this handler will process
//   - TResponse: The type of response data this handler will return
//
// Parameters:
//   - config: Configuration for the handler behavior and features
//
// Returns:
//   - *MiddlewareHandler[TRequest, TResponse]: Configured middleware handler
//
// Usage:
//
//	handler := NewMiddlewareHandler[*http.Request, *http.Response](&MiddlewareHandlerConfig{
//	    Name: "HTTP Handler",
//	    EnableMetrics: true,
//	    EnableTracing: true,
//	    DefaultTimeout: 30 * time.Second,
//	    MaxConcurrency: 100,
//	})
func NewMiddlewareHandler[TRequest, TResponse any](
	config *MiddlewareHandlerConfig,
) *MiddlewareHandler[TRequest, TResponse] {
	// Apply default configuration if needed
	if config == nil {
		config = &MiddlewareHandlerConfig{
			Name:           "DefaultMiddlewareHandler",
			EnableMetrics:  true,
			DefaultTimeout: 30 * time.Second,
			MaxConcurrency: 0, // unlimited
		}
	}

	// Create execution chain configuration
	executionConfig := types.ExecutionConfig{
		Name:            config.Name,
		EnableProfiling: config.EnableProfiling,
		EnableTracing:   config.EnableTracing,
		EnableMetrics:   config.EnableMetrics,
		Timeout:         config.DefaultTimeout,
		MaxConcurrency:  config.MaxConcurrency,
		RetryPolicy:     config.RetryPolicy,
		CircuitBreaker:  config.CircuitBreakerConfig,
		Properties:      config.Properties,
	}

	return &MiddlewareHandler[TRequest, TResponse]{
		middlewares:    make([]interfaces.Middleware[TRequest, TResponse], 0),
		config:         config,
		metrics:        &HandlerMetrics{},
		contextFactory: types.NewContextFactory(),
		executionChain: types.NewExecutionChain[TRequest, TResponse](executionConfig),
	}
}

// AddMiddleware registers one or more middleware with the handler.
// Middleware will be executed in the order they are added, with each middleware
// wrapping the next in the execution chain.
//
// Parameters:
//   - middleware: One or more middleware instances to register
//
// Thread Safety:
//
//	This method is safe for concurrent use and will not affect ongoing executions.
func (mh *MiddlewareHandler[TRequest, TResponse]) AddMiddleware(middleware ...interfaces.Middleware[TRequest, TResponse]) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// Add to local middleware list
	mh.middlewares = append(mh.middlewares, middleware...)

	// Update execution chain
	mh.executionChain.AddMiddleware(middleware...)
}

// PrependMiddleware adds middleware to the beginning of the execution chain.
// This is useful for adding "outer" middleware like error handlers or loggers
// that should wrap all other middleware.
//
// Parameters:
//   - middleware: One or more middleware instances to prepend
//
// Thread Safety:
//
//	This method is safe for concurrent use and will not affect ongoing executions.
func (mh *MiddlewareHandler[TRequest, TResponse]) PrependMiddleware(middleware ...interfaces.Middleware[TRequest, TResponse]) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// Create new slice with prepended middleware
	newMiddlewares := make([]interfaces.Middleware[TRequest, TResponse], len(middleware)+len(mh.middlewares))
	copy(newMiddlewares, middleware)
	copy(newMiddlewares[len(middleware):], mh.middlewares)
	mh.middlewares = newMiddlewares

	// Update execution chain
	mh.executionChain.PrependMiddleware(middleware...)
}

// Handle processes a request through the registered middleware chain.
// This is the main entry point for request processing, providing comprehensive
// error handling, performance monitoring, and execution tracking.
//
// Parameters:
//   - ctx: Context for the request, including cancellation and timeout
//   - request: The request data to process
//   - finalHandler: The handler to execute after all middleware
//
// Returns:
//   - TResponse: The processed response data
//   - error: Any error that occurred during processing
//
// Error Handling:
//   - Context cancellation and timeout are respected
//   - Middleware errors are propagated with additional context
//   - Panic recovery is handled gracefully
//   - Circuit breaker patterns are applied if configured
//
// Performance Monitoring:
//   - Execution time is tracked and reported
//   - Concurrency levels are monitored
//   - Success/failure rates are calculated
//   - Detailed metrics are available via GetMetrics()
func (mh *MiddlewareHandler[TRequest, TResponse]) Handle(
	ctx context.Context,
	request TRequest,
	finalHandler interfaces.Handler[TRequest, TResponse],
) (TResponse, error) {
	start := time.Now()
	var zeroResponse TResponse

	// Update concurrency metrics
	mh.updateConcurrency(1)
	defer mh.updateConcurrency(-1)

	// Create execution context if tracing is enabled
	var executionCtx interfaces.ExecutionContext
	if mh.config.EnableTracing {
		executionCtx = mh.contextFactory.CreateExecutionContext(ctx)
		ctx = executionCtx
	}

	// Execute the middleware chain with error recovery
	response, err := mh.executeWithRecovery(ctx, request, finalHandler)

	// Update performance metrics
	duration := time.Since(start)
	mh.updateMetrics(err == nil, duration)

	// Record execution in tracing context
	if executionCtx != nil {
		if err != nil {
			executionCtx.RecordError(err, "MiddlewareHandler")
		}
		executionCtx.RecordMetric("execution_time", duration.Milliseconds(), map[string]string{"unit": "ms"})
	}

	return response, err
}

// executeWithRecovery executes the middleware chain with panic recovery.
// This private method ensures that panics in middleware don't crash the entire application.
func (mh *MiddlewareHandler[TRequest, TResponse]) executeWithRecovery(
	ctx context.Context,
	request TRequest,
	finalHandler interfaces.Handler[TRequest, TResponse],
) (response TResponse, err error) {
	// Recover from panics in middleware
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("middleware panic: %v", r)
			// Could add stack trace here if needed
		}
	}()

	// Execute through the execution chain
	return mh.executionChain.Execute(ctx, request, finalHandler)
}

// updateConcurrency updates the current concurrency count in a thread-safe manner.
func (mh *MiddlewareHandler[TRequest, TResponse]) updateConcurrency(delta int64) {
	mh.metrics.mu.Lock()
	defer mh.metrics.mu.Unlock()

	mh.metrics.currentConcurrency += delta
	if mh.metrics.currentConcurrency > mh.metrics.maxConcurrency {
		mh.metrics.maxConcurrency = mh.metrics.currentConcurrency
	}
}

// updateMetrics updates the handler performance metrics in a thread-safe manner.
func (mh *MiddlewareHandler[TRequest, TResponse]) updateMetrics(success bool, duration time.Duration) {
	mh.metrics.mu.Lock()
	defer mh.metrics.mu.Unlock()

	// Update request counts
	mh.metrics.totalRequests++
	if success {
		mh.metrics.successfulRequests++
	} else {
		mh.metrics.failedRequests++
	}

	// Update timing metrics
	mh.metrics.totalProcessingTime += duration
	mh.metrics.averageProcessingTime = mh.metrics.totalProcessingTime / time.Duration(mh.metrics.totalRequests)
	mh.metrics.lastRequestTime = time.Now()

	// Update min/max timing
	if mh.metrics.totalRequests == 1 || duration < mh.metrics.minProcessingTime {
		mh.metrics.minProcessingTime = duration
	}
	if duration > mh.metrics.maxProcessingTime {
		mh.metrics.maxProcessingTime = duration
	}
}

// GetMetrics returns a snapshot of current handler performance metrics.
// The returned metrics provide detailed information about handler performance,
// including request counts, timing statistics, and concurrency levels.
//
// Returns:
//   - HandlerMetricsSnapshot: Current performance metrics
//
// Thread Safety:
//
//	This method is safe for concurrent use and returns a consistent snapshot
//	of metrics at the time of the call.
func (mh *MiddlewareHandler[TRequest, TResponse]) GetMetrics() HandlerMetricsSnapshot {
	mh.metrics.mu.RLock()
	defer mh.metrics.mu.RUnlock()

	var successRate float64
	if mh.metrics.totalRequests > 0 {
		successRate = float64(mh.metrics.successfulRequests) / float64(mh.metrics.totalRequests) * 100
	}

	return HandlerMetricsSnapshot{
		HandlerName:           mh.config.Name,
		TotalRequests:         mh.metrics.totalRequests,
		SuccessfulRequests:    mh.metrics.successfulRequests,
		FailedRequests:        mh.metrics.failedRequests,
		SuccessRate:           successRate,
		TotalProcessingTime:   mh.metrics.totalProcessingTime,
		AverageProcessingTime: mh.metrics.averageProcessingTime,
		MinProcessingTime:     mh.metrics.minProcessingTime,
		MaxProcessingTime:     mh.metrics.maxProcessingTime,
		CurrentConcurrency:    mh.metrics.currentConcurrency,
		MaxConcurrency:        mh.metrics.maxConcurrency,
		LastRequestTime:       mh.metrics.lastRequestTime,
	}
}

// HandlerMetricsSnapshot provides a point-in-time view of handler performance metrics.
type HandlerMetricsSnapshot struct {
	HandlerName           string        `json:"handler_name"`
	TotalRequests         int64         `json:"total_requests"`
	SuccessfulRequests    int64         `json:"successful_requests"`
	FailedRequests        int64         `json:"failed_requests"`
	SuccessRate           float64       `json:"success_rate"`
	TotalProcessingTime   time.Duration `json:"total_processing_time"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	MinProcessingTime     time.Duration `json:"min_processing_time"`
	MaxProcessingTime     time.Duration `json:"max_processing_time"`
	CurrentConcurrency    int64         `json:"current_concurrency"`
	MaxConcurrency        int64         `json:"max_concurrency"`
	LastRequestTime       time.Time     `json:"last_request_time"`
}

// GetMiddlewareCount returns the number of registered middleware.
func (mh *MiddlewareHandler[TRequest, TResponse]) GetMiddlewareCount() int {
	mh.mu.RLock()
	defer mh.mu.RUnlock()
	return len(mh.middlewares)
}

// GetConfig returns the current handler configuration.
func (mh *MiddlewareHandler[TRequest, TResponse]) GetConfig() *MiddlewareHandlerConfig {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *mh.config
	return &configCopy
}

// UpdateConfig updates the handler configuration with new settings.
// Some configuration changes may require restarting active operations to take effect.
//
// Parameters:
//   - config: New configuration to apply
//
// Thread Safety:
//
//	This method is safe for concurrent use, though configuration changes
//	may not affect already-running middleware executions.
func (mh *MiddlewareHandler[TRequest, TResponse]) UpdateConfig(config *MiddlewareHandlerConfig) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// Update configuration
	mh.config = config

	// Update execution chain configuration if needed
	executionConfig := types.ExecutionConfig{
		Name:            config.Name,
		EnableProfiling: config.EnableProfiling,
		EnableTracing:   config.EnableTracing,
		EnableMetrics:   config.EnableMetrics,
		Timeout:         config.DefaultTimeout,
		MaxConcurrency:  config.MaxConcurrency,
		RetryPolicy:     config.RetryPolicy,
		CircuitBreaker:  config.CircuitBreakerConfig,
		Properties:      config.Properties,
	}

	// Create new execution chain with updated config
	newExecutionChain := types.NewExecutionChain[TRequest, TResponse](executionConfig)
	newExecutionChain.AddMiddleware(mh.middlewares...)
	mh.executionChain = newExecutionChain
}

// Reset clears all registered middleware and resets performance metrics.
// This method is primarily useful for testing and development scenarios.
//
// Thread Safety:
//
//	This method is safe for concurrent use, though it may affect ongoing executions.
func (mh *MiddlewareHandler[TRequest, TResponse]) Reset() {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	// Clear middleware
	mh.middlewares = mh.middlewares[:0]

	// Reset metrics
	mh.metrics.mu.Lock()
	mh.metrics.totalRequests = 0
	mh.metrics.successfulRequests = 0
	mh.metrics.failedRequests = 0
	mh.metrics.totalProcessingTime = 0
	mh.metrics.averageProcessingTime = 0
	mh.metrics.minProcessingTime = 0
	mh.metrics.maxProcessingTime = 0
	mh.metrics.currentConcurrency = 0
	mh.metrics.maxConcurrency = 0
	mh.metrics.mu.Unlock()

	// Create new execution chain
	executionConfig := types.ExecutionConfig{
		Name:            mh.config.Name,
		EnableProfiling: mh.config.EnableProfiling,
		EnableTracing:   mh.config.EnableTracing,
		EnableMetrics:   mh.config.EnableMetrics,
		Timeout:         mh.config.DefaultTimeout,
		MaxConcurrency:  mh.config.MaxConcurrency,
		RetryPolicy:     mh.config.RetryPolicy,
		CircuitBreaker:  mh.config.CircuitBreakerConfig,
		Properties:      mh.config.Properties,
	}
	mh.executionChain = types.NewExecutionChain[TRequest, TResponse](executionConfig)
}
