package interfaces

// Contextable defines the interface for exceptions that can carry context information.
// This interface follows the Interface Segregation Principle (ISP) by focusing
// solely on context-related functionality.
type Contextable interface {
	// GetContext returns additional context information about the exception
	GetContext() map[string]interface{}

	// SetContext sets additional context information
	SetContext(context map[string]interface{}) ExceptionInterface

	// WithContext adds a single context key-value pair
	WithContext(key string, value interface{}) ExceptionInterface
}
