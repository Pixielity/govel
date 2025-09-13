package interfaces

// Stackable defines the interface for exceptions that can provide stack trace information.
// This interface follows the Interface Segregation Principle (ISP) by focusing
// solely on stack trace-related functionality.
type Stackable interface {
	// GetStackTrace returns the stack trace where the exception occurred
	GetStackTrace() []string
}
