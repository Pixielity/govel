package providers

import (
	"strings"

	solutionInterface "govel/exceptions/interfaces/solution"
	httpSolutions "govel/exceptions/solutions/http"
)

// HTTPExceptionProvider provides solutions for common HTTP exceptions.
// This provider can provide solutions for standard HTTP error conditions.
type HTTPExceptionProvider struct{}

// NewHTTPExceptionProvider creates a new HTTP exception solution provider
func NewHTTPExceptionProvider() *HTTPExceptionProvider {
	return &HTTPExceptionProvider{}
}

// CanSolve determines if this provider can provide a solution for the given error
func (p *HTTPExceptionProvider) CanSolve(err error) bool {
	// Check if this is an HTTP-related exception by looking for HTTP status codes in the message
	message := strings.ToLower(err.Error())
	
	// Check for common HTTP errors
	return strings.Contains(message, "404") || strings.Contains(message, "not found") ||
		strings.Contains(message, "401") || strings.Contains(message, "unauthorized") ||
		strings.Contains(message, "403") || strings.Contains(message, "forbidden") ||
		strings.Contains(message, "400") || strings.Contains(message, "bad request") ||
		strings.Contains(message, "500") || strings.Contains(message, "internal server") ||
		strings.Contains(message, "422") || strings.Contains(message, "unprocessable") ||
		strings.Contains(message, "405") || strings.Contains(message, "method not allowed") ||
		strings.Contains(message, "429") || strings.Contains(message, "too many requests") ||
		strings.Contains(message, "503") || strings.Contains(message, "service unavailable") ||
		strings.Contains(message, "409") || strings.Contains(message, "conflict")
}

// GetSolutions returns solutions for the given HTTP error
func (p *HTTPExceptionProvider) GetSolutions(err error) []solutionInterface.Solution {
	message := strings.ToLower(err.Error())
	var solutions []solutionInterface.Solution
	
	if strings.Contains(message, "404") || strings.Contains(message, "not found") {
		solutions = append(solutions, httpSolutions.NewNotFoundSolution(""))
	}
	
	if strings.Contains(message, "401") || strings.Contains(message, "unauthorized") {
		solutions = append(solutions, httpSolutions.NewUnauthorizedSolution())
	}
	
	if strings.Contains(message, "403") || strings.Contains(message, "forbidden") {
		solutions = append(solutions, httpSolutions.NewForbiddenSolution(""))
	}
	
	if strings.Contains(message, "400") || strings.Contains(message, "bad request") {
		solutions = append(solutions, httpSolutions.NewBadRequestSolution())
	}
	
	if strings.Contains(message, "500") || strings.Contains(message, "internal server") {
		solutions = append(solutions, httpSolutions.NewInternalServerErrorSolution())
	}
	
	if strings.Contains(message, "422") || strings.Contains(message, "unprocessable") {
		solutions = append(solutions, httpSolutions.NewValidationErrorSolution(nil))
	}
	
	if strings.Contains(message, "405") || strings.Contains(message, "method not allowed") {
		solutions = append(solutions, httpSolutions.NewMethodNotAllowedSolution("", nil))
	}
	
	if strings.Contains(message, "429") || strings.Contains(message, "too many requests") {
		solutions = append(solutions, httpSolutions.NewTooManyRequestsSolution(0, 0))
	}
	
	if strings.Contains(message, "503") || strings.Contains(message, "service unavailable") {
		solutions = append(solutions, httpSolutions.NewServiceUnavailableSolution(0))
	}
	
	if strings.Contains(message, "409") || strings.Contains(message, "conflict") {
		solutions = append(solutions, httpSolutions.NewConflictSolution())
	}
	
	return solutions
}

// Ensure HTTPExceptionProvider implements the HasSolutionsForThrowable interface
var _ solutionInterface.HasSolutionsForThrowable = (*HTTPExceptionProvider)(nil)
