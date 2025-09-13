package solution

// ProvidesSolution interface should be implemented by exceptions that can provide their own solution.
// This is inspired by Laravel's Spatie\ErrorSolutions\Interfaces\ProvidesSolution interface.
type ProvidesSolution interface {
	// GetSolution returns the solution provided by this exception
	GetSolution() Solution
}
