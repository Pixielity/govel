package interfaces

import "govel/packages/exceptions/interfaces/solution"

// Solutionable defines the interface for exceptions that can provide or be associated with solutions.
// This interface follows the Interface Segregation Principle (ISP) by focusing
// solely on solution-related functionality.
type Solutionable interface {
	// GetSolution returns the solution for this exception (if it implements ProvidesSolution)
	GetSolution() solution.Solution

	// HasSolution returns true if this exception has an associated solution
	HasSolution() bool

	// SetSolution sets a solution for this exception
	SetSolution(sol solution.Solution) ExceptionInterface

	// WithSolution sets a solution for this exception and returns the exception (for method chaining)
	WithSolution(sol solution.Solution) ExceptionInterface
}
