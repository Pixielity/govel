package interfaces

// SolutionProviderInterface interface for providing error solutions
type SolutionProviderInterface interface {
	// GetSolutions returns a list of possible solutions for the given error
	GetSolutions(error) []SolutionInterface
}
