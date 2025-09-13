// Package interfaces - Pipeline-aware interface definitions
// This file contains interfaces for middleware components that can integrate
// with the GoVel Pipeline system, enabling seamless interoperability between
// the middleware and pipeline execution models.
package interfaces

import (
	"context"
)

// PipelineAware defines the interface for middleware that can integrate with pipelines.
// This interface enables middleware to be executed within pipeline chains while
// maintaining their middleware semantics and providing pipeline-specific capabilities.
//
// PipelineAware middleware can:
// - Be converted to pipeline pipes for use in pipeline chains
// - Access pipeline context and hub configurations
// - Participate in pipeline transactions and error handling
// - Utilize pipeline tracing and debugging features
// - Benefit from pipeline performance optimizations
//
// Type Parameters:
//   - TRequest: The type of request data flowing through the pipeline
//   - TResponse: The type of response data returned by the pipeline
//
// Design Philosophy:
//   PipelineAware middleware acts as a bridge between the middleware pattern
//   (focused on request/response processing) and the pipeline pattern
//   (focused on data transformation workflows). This enables the best of
//   both worlds while maintaining clean separation of concerns.
type PipelineAware[TRequest, TResponse any] interface {
	Middleware[TRequest, TResponse] // Extend core middleware interface

	// ToPipe converts this middleware into a pipeline-compatible pipe function.
	// The returned pipe can be used directly in pipeline chains while maintaining
	// the middleware's original behavior and error handling semantics.
	//
	// The conversion process:
	// 1. Wraps the middleware's Handle method in a pipe-compatible function
	// 2. Handles type conversion between pipeline and middleware signatures
	// 3. Maintains error propagation and context management
	// 4. Preserves middleware-specific behavior and state
	//
	// Returns:
	//   - interface{}: A pipeline-compatible pipe function
	//   - error: Any error that occurred during conversion
	//
	// Usage:
	//   pipe, err := middleware.ToPipe()
	//   if err == nil {
	//       result, err := pipeline.Send(data).Through(pipe).Then(handler)
	//   }
	//
	// Error Conditions:
	//   - Type incompatibility between middleware and pipeline signatures
	//   - Missing dependencies or configuration for pipeline integration
	//   - Internal state that cannot be represented in pipeline format
	ToPipe() (interface{}, error)

	// GetPipelineConfig returns configuration specific to pipeline integration.
	// This configuration controls how the middleware behaves when used in pipelines.
	//
	// Configuration options may include:
	// - Transaction participation settings
	// - Error handling preferences
	// - Performance optimization flags
	// - Tracing and debugging options
	// - Pipeline hub preferences
	//
	// Returns:
	//   - PipelineConfig: Configuration for pipeline integration
	GetPipelineConfig() PipelineConfig

	// SetPipelineConfig updates the pipeline integration configuration.
	// This allows dynamic reconfiguration of pipeline behavior.
	//
	// Parameters:
	//   - config: New configuration for pipeline integration
	//
	// Returns:
	//   - error: Any error that occurred updating the configuration
	SetPipelineConfig(config PipelineConfig) error

	// SupportsPipelineFeature checks if this middleware supports a specific pipeline feature.
	// This enables dynamic capability detection for advanced pipeline workflows.
	//
	// Parameters:
	//   - feature: The pipeline feature to check for
	//
	// Returns:
	//   - bool: true if the feature is supported
	//
	// Common features:
	//   - "transactions": Database transaction support
	//   - "tracing": Distributed tracing integration
	//   - "caching": Response caching capabilities
	//   - "async": Asynchronous processing support
	//   - "batching": Batch processing capabilities
	SupportsPipelineFeature(feature string) bool
}

// PipelineIntegrator defines the interface for components that can facilitate
// integration between middleware systems and pipeline systems.
//
// PipelineIntegrator provides:
// - Bidirectional conversion between middleware and pipes
// - Pipeline execution of middleware chains
// - Configuration management for pipeline integration
// - Performance optimization for pipeline-middleware hybrid workflows
//
// Type Parameters:
//   - TRequest: The type of request data for integration
//   - TResponse: The type of response data for integration
type PipelineIntegrator[TRequest, TResponse any] interface {
	// ExecuteMiddlewareAsPipeline runs middleware through the pipeline system.
	// This method provides access to pipeline features like transactions,
	// advanced error handling, and hub-based configuration.
	//
	// Parameters:
	//   - ctx: Context for execution
	//   - middleware: The middleware to execute
	//   - request: The request to process
	//   - handler: The final handler to execute
	//
	// Returns:
	//   - TResponse: The processed response
	//   - error: Any error that occurred during execution
	ExecuteMiddlewareAsPipeline(ctx context.Context, middleware Middleware[TRequest, TResponse], request TRequest, handler Handler[TRequest, TResponse]) (TResponse, error)

	// ConvertMiddlewareToPipe converts middleware to a pipeline pipe.
	//
	// Parameters:
	//   - middleware: The middleware to convert
	//
	// Returns:
	//   - interface{}: Pipeline-compatible pipe
	//   - error: Any error during conversion
	ConvertMiddlewareToPipe(middleware Middleware[TRequest, TResponse]) (interface{}, error)

	// ConvertPipeToMiddleware converts a pipeline pipe to middleware.
	//
	// Parameters:
	//   - pipe: The pipeline pipe to convert
	//
	// Returns:
	//   - Middleware[TRequest, TResponse]: Middleware-compatible wrapper
	//   - error: Any error during conversion
	ConvertPipeToMiddleware(pipe interface{}) (Middleware[TRequest, TResponse], error)

	// GetPipelineHub returns the pipeline hub for advanced pipeline operations.
	//
	// Returns:
	//   - interface{}: Pipeline hub instance
	//   - error: Any error accessing the hub
	GetPipelineHub() (interface{}, error)

	// WithPipelineHub sets the pipeline hub for this integrator.
	//
	// Parameters:
	//   - hub: The pipeline hub to use
	//
	// Returns:
	//   - PipelineIntegrator[TRequest, TResponse]: Updated integrator
	WithPipelineHub(hub interface{}) PipelineIntegrator[TRequest, TResponse]
}

// PipelineConfig contains configuration options for pipeline integration.
type PipelineConfig struct {
	// EnableTransactions indicates whether to wrap execution in database transactions
	EnableTransactions bool

	// TransactionIsolation specifies the isolation level for transactions
	TransactionIsolation string

	// EnableTracing indicates whether to enable distributed tracing
	EnableTracing bool

	// TracingServiceName specifies the service name for tracing
	TracingServiceName string

	// EnableErrorRecovery indicates whether to enable automatic error recovery
	EnableErrorRecovery bool

	// MaxRetries specifies the maximum number of retry attempts
	MaxRetries int

	// RetryDelay specifies the delay between retry attempts
	RetryDelay string // Duration string like "100ms", "1s"

	// EnableCaching indicates whether to enable response caching
	EnableCaching bool

	// CacheTimeout specifies how long to cache responses
	CacheTimeout string // Duration string

	// EnableMetrics indicates whether to collect performance metrics
	EnableMetrics bool

	// MetricPrefix specifies the prefix for metric names
	MetricPrefix string

	// PipelineName specifies which named pipeline to use (empty for default)
	PipelineName string

	// HubConfig contains hub-specific configuration options
	HubConfig map[string]interface{}

	// CustomProperties contains custom configuration properties
	CustomProperties map[string]interface{}
}

// PipelineFeature represents a capability that pipeline-aware middleware can support.
type PipelineFeature string

const (
	// PipelineFeatureTransactions indicates support for database transactions
	PipelineFeatureTransactions PipelineFeature = "transactions"

	// PipelineFeatureTracing indicates support for distributed tracing
	PipelineFeatureTracing PipelineFeature = "tracing"

	// PipelineFeatureCaching indicates support for response caching
	PipelineFeatureCaching PipelineFeature = "caching"

	// PipelineFeatureAsync indicates support for asynchronous processing
	PipelineFeatureAsync PipelineFeature = "async"

	// PipelineFeatureBatching indicates support for batch processing
	PipelineFeatureBatching PipelineFeature = "batching"

	// PipelineFeatureRetry indicates support for automatic retry logic
	PipelineFeatureRetry PipelineFeature = "retry"

	// PipelineFeatureCircuitBreaker indicates support for circuit breaker pattern
	PipelineFeatureCircuitBreaker PipelineFeature = "circuit_breaker"

	// PipelineFeatureMetrics indicates support for performance metrics
	PipelineFeatureMetrics PipelineFeature = "metrics"

	// PipelineFeatureValidation indicates support for input/output validation
	PipelineFeatureValidation PipelineFeature = "validation"

	// PipelineFeatureCompression indicates support for data compression
	PipelineFeatureCompression PipelineFeature = "compression"
)

// String returns the string representation of a PipelineFeature.
func (f PipelineFeature) String() string {
	return string(f)
}

// PipelineCapability represents the level of support for a pipeline feature.
type PipelineCapability int

const (
	// PipelineCapabilityNone indicates no support for the feature
	PipelineCapabilityNone PipelineCapability = iota

	// PipelineCapabilityBasic indicates basic support for the feature
	PipelineCapabilityBasic

	// PipelineCapabilityAdvanced indicates advanced support for the feature
	PipelineCapabilityAdvanced

	// PipelineCapabilityFull indicates full support for the feature
	PipelineCapabilityFull
)

// String returns the string representation of a PipelineCapability.
func (c PipelineCapability) String() string {
	switch c {
	case PipelineCapabilityNone:
		return "none"
	case PipelineCapabilityBasic:
		return "basic"
	case PipelineCapabilityAdvanced:
		return "advanced"
	case PipelineCapabilityFull:
		return "full"
	default:
		return "unknown"
	}
}

// PipelineAwareFactory defines the interface for creating pipeline-aware middleware.
// This factory enables dynamic creation of middleware that can integrate with pipelines.
//
// Type Parameters:
//   - TRequest: The type of request data for created middleware
//   - TResponse: The type of response data for created middleware
type PipelineAwareFactory[TRequest, TResponse any] interface {
	// CreatePipelineAware creates pipeline-aware middleware.
	//
	// Parameters:
	//   - config: Configuration for the middleware
	//   - pipelineConfig: Pipeline-specific configuration
	//
	// Returns:
	//   - PipelineAware[TRequest, TResponse]: Created middleware
	//   - error: Any error during creation
	CreatePipelineAware(config map[string]interface{}, pipelineConfig PipelineConfig) (PipelineAware[TRequest, TResponse], error)

	// GetSupportedFeatures returns the features supported by created middleware.
	//
	// Returns:
	//   - map[PipelineFeature]PipelineCapability: Supported features and their capability levels
	GetSupportedFeatures() map[PipelineFeature]PipelineCapability

	// ValidatePipelineConfig validates pipeline configuration.
	//
	// Parameters:
	//   - config: Pipeline configuration to validate
	//
	// Returns:
	//   - error: Any validation error
	ValidatePipelineConfig(config PipelineConfig) error
}
