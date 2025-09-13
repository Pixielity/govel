package solutions

import (
	coreSolution "govel/exceptions/core/solution"
	solutionInterfaces "govel/exceptions/interfaces/solution"
	httpSolutions "govel/exceptions/solutions/http"
	"govel/exceptions/solutions/providers"
	"govel/exceptions/solutions/runnable"
)

// Re-export core types
type Solution = solutionInterfaces.Solution
type RunnableSolution = solutionInterfaces.RunnableSolution
type HasSolutionsForThrowable = solutionInterfaces.HasSolutionsForThrowable
type SolutionProvider = solutionInterfaces.SolutionProvider
type BaseSolution = coreSolution.BaseSolution
type SolutionProviderRepository = coreSolution.SolutionProviderRepository

// Re-export core constructors
var NewBaseSolution = coreSolution.NewBaseSolution
var NewSolutionProviderRepository = coreSolution.NewSolutionProviderRepository

// Re-export HTTP solutions
var NewNotFoundSolution = httpSolutions.NewNotFoundSolution
var NewUnauthorizedSolution = httpSolutions.NewUnauthorizedSolution
var NewForbiddenSolution = httpSolutions.NewForbiddenSolution
var NewBadRequestSolution = httpSolutions.NewBadRequestSolution
var NewInternalServerErrorSolution = httpSolutions.NewInternalServerErrorSolution
var NewServiceUnavailableSolution = httpSolutions.NewServiceUnavailableSolution
var NewTooManyRequestsSolution = httpSolutions.NewTooManyRequestsSolution
var NewConflictSolution = httpSolutions.NewConflictSolution

// Re-export providers
var NewHTTPExceptionProvider = providers.NewHTTPExceptionProvider
var NewCommonRunnableSolutionsProvider = providers.NewCommonRunnableSolutionsProvider

// Re-export runnable solutions
var NewGenerateAppKeySolution = runnable.NewGenerateAppKeySolution
var NewCreateDirectorySolution = runnable.NewCreateDirectorySolution
var NewInstallDependencySolution = runnable.NewInstallDependencySolution
var NewFixPermissionsSolution = runnable.NewFixPermissionsSolution
