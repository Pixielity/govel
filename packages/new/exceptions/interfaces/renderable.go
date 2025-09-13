package interfaces

// Renderable defines the interface for exceptions that can be rendered to various formats.
// This interface follows the Interface Segregation Principle (ISP) by focusing
// solely on rendering-related functionality.
type Renderable interface {
	// Render returns a response representation of the exception
	// This method provides a Laravel-like rendering system for exceptions
	Render() map[string]interface{}
}
