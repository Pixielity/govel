package interfaces

// ExceptionInterface defines the contract for all GoVel exceptions.
// This interface composes multiple smaller interfaces following the Interface Segregation Principle (ISP).
// Each component interface focuses on a specific aspect of exception functionality.
type ExceptionInterface interface {
	// Embed Go's built-in error interface
	error

	// HTTP-related functionality
	HTTPable

	// Context management functionality  
	Contextable

	// Rendering functionality
	Renderable

	// Stack trace functionality
	Stackable

	// Solution-related functionality
	Solutionable
}
