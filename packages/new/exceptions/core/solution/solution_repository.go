package solution

import (
	"sync"

	solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// SolutionProviderRepository manages solution providers and provides solutions for errors.
// This is inspired by Laravel's solution provider system.
type SolutionProviderRepository struct {
	providers []solutionInterface.HasSolutionsForThrowable
	mutex     sync.RWMutex
}

// NewSolutionProviderRepository creates a new solution provider repository
func NewSolutionProviderRepository() *SolutionProviderRepository {
	return &SolutionProviderRepository{
		providers: make([]solutionInterface.HasSolutionsForThrowable, 0),
	}
}

// RegisterSolutionProviders registers multiple solution providers
func (repo *SolutionProviderRepository) RegisterSolutionProviders(providers []solutionInterface.HasSolutionsForThrowable) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	
	repo.providers = append(repo.providers, providers...)
}

// RegisterSolutionProvider registers a single solution provider
func (repo *SolutionProviderRepository) RegisterSolutionProvider(provider solutionInterface.HasSolutionsForThrowable) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	
	repo.providers = append(repo.providers, provider)
}

// GetSolutionsForError returns all available solutions for the given error
func (repo *SolutionProviderRepository) GetSolutionsForError(err error) []solutionInterface.Solution {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	
	var solutions []solutionInterface.Solution
	
	// First check if the error itself provides a solution
	if providesSolution, ok := err.(solutionInterface.ProvidesSolution); ok {
		solution := providesSolution.GetSolution()
		if solution != nil {
			solutions = append(solutions, solution)
		}
	}
	
	// Then check all registered providers
	for _, provider := range repo.providers {
		if provider.CanSolve(err) {
			providerSolutions := provider.GetSolutions(err)
			solutions = append(solutions, providerSolutions...)
		}
	}
	
	return solutions
}

// GetProviders returns all registered solution providers (mainly for testing)
func (repo *SolutionProviderRepository) GetProviders() []solutionInterface.HasSolutionsForThrowable {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	providers := make([]solutionInterface.HasSolutionsForThrowable, len(repo.providers))
	copy(providers, repo.providers)
	return providers
}

// Ensure SolutionProviderRepository implements the SolutionProvider interface
var _ solutionInterface.SolutionProvider = (*SolutionProviderRepository)(nil)
