package providers

import (
	"os"
	"strings"

	solutionInterface "govel/packages/exceptions/interfaces/solution"
	"govel/packages/exceptions/solutions/runnable"
)

// CommonRunnableSolutionsProvider provides runnable solutions for common development issues
type CommonRunnableSolutionsProvider struct{}

// NewCommonRunnableSolutionsProvider creates a new common runnable solutions provider
func NewCommonRunnableSolutionsProvider() *CommonRunnableSolutionsProvider {
	return &CommonRunnableSolutionsProvider{}
}

// CanSolve determines if this provider can provide a solution for the given error
func (p *CommonRunnableSolutionsProvider) CanSolve(err error) bool {
	message := err.Error()
	
	// Check for common development issues
	return contains(message, "app key", "encryption key", "missing key") ||
		contains(message, "permission denied", "access denied") ||
		contains(message, "directory", "not found", "no such file") ||
		contains(message, "missing dependency", "module not found", "package not installed")
}

// GetSolutions returns runnable solutions for common development issues
func (p *CommonRunnableSolutionsProvider) GetSolutions(err error) []solutionInterface.Solution {
	message := err.Error()
	var solutions []solutionInterface.Solution
	
	if contains(message, "app key", "encryption key", "missing key") {
		solutions = append(solutions, runnable.NewGenerateAppKeySolution())
	}
	
	if contains(message, "storage", "logs", "directory") && contains(message, "not found", "no such file") {
		// Common directories that might be missing
		commonDirs := []string{"storage/logs", "storage/app", "storage/framework", "bootstrap/cache"}
		for _, dir := range commonDirs {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				solutions = append(solutions, runnable.NewCreateDirectorySolution(dir))
			}
		}
	}
	
	if contains(message, "go.mod", "module") && contains(message, "not found") {
		solutions = append(solutions, runnable.NewInstallDependencySolution("go modules", "go mod tidy"))
	}
	
	if contains(message, "permission", "chmod") && contains(message, "denied", "failed") {
		// Common permission fixes
		commonPaths := []string{"storage", "bootstrap/cache", "logs"}
		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				solutions = append(solutions, runnable.NewFixPermissionsSolution(path, 0755))
			}
		}
	}
	
	return solutions
}

// Helper function to check if message contains any of the given substrings
func contains(message string, substrings ...string) bool {
	message = strings.ToLower(message)
	for _, substr := range substrings {
		if len(substr) > 0 && strings.Contains(message, strings.ToLower(substr)) {
			return true
		}
	}
	return false
}

// Ensure CommonRunnableSolutionsProvider implements the HasSolutionsForThrowable interface
var _ solutionInterface.HasSolutionsForThrowable = (*CommonRunnableSolutionsProvider)(nil)
