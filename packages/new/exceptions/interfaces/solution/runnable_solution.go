package solution

// RunnableSolution extends Solution to provide solutions that can be executed automatically.
// This is inspired by Laravel's Spatie\ErrorSolutions\Interfaces\RunnableSolution interface.
type RunnableSolution interface {
	Solution

	// GetSolutionActionDescription returns a description of what the runnable action will do
	GetSolutionActionDescription() string

	// GetRunButtonText returns the text to display on the run button
	GetRunButtonText() string

	// Run executes the solution with the given parameters
	Run(parameters map[string]interface{}) error

	// GetRunParameters returns the parameters that the solution expects
	GetRunParameters() map[string]interface{}
}
