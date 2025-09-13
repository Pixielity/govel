package interfaces

// ContextBuilderInterface interface for building error context
type ContextBuilderInterface interface {
	// BuildContext builds error context from the provided data
	BuildContext(interface{}) ErrorContextInterface
	
	// BuildRequestContext builds request-specific context
	BuildRequestContext(interface{}) RequestContextInterface
	
	// BuildEnvironmentContext builds environment-specific context  
	BuildEnvironmentContext() EnvContextInterface
}
