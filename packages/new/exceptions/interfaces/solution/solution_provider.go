package solution

// HasSolutionsForThrowable interface is used for solution providers that can provide solutions
// for specific types of errors/exceptions.
// This is inspired by Laravel's Spatie\ErrorSolutions\Interfaces\HasSolutionsForThrowable interface.
type HasSolutionsForThrowable interface {
	// CanSolve determines if this provider can provide a solution for the given error
	CanSolve(err error) bool

	// GetSolutions returns an array of solutions for the given error
	GetSolutions(err error) []Solution
}

// SolutionProvider interface defines methods for managing solution providers
type SolutionProvider interface {
	// RegisterSolutionProviders registers multiple solution providers
	RegisterSolutionProviders(providers []HasSolutionsForThrowable)

	// RegisterSolutionProvider registers a single solution provider
	RegisterSolutionProvider(provider HasSolutionsForThrowable)

	// GetSolutionsForError returns all available solutions for the given error
	GetSolutionsForError(err error) []Solution

	// GetProviders returns all registered solution providers
	GetProviders() []HasSolutionsForThrowable
}
