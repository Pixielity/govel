package solution

// Solution defines the basic interface for providing solutions to exceptions.
// This is inspired by Laravel's Spatie\ErrorSolutions\Interfaces\Solution interface.
type Solution interface {
	// GetSolutionTitle returns the title of the solution
	GetSolutionTitle() string

	// GetSolutionDescription returns a detailed description of the solution
	GetSolutionDescription() string

	// GetDocumentationLinks returns a map of documentation links related to the solution
	GetDocumentationLinks() map[string]string
}
