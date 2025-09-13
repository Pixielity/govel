package http

import (
	"govel/packages/exceptions/core/solution"
	solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// UnauthorizedSolution provides specific guidance for 401 Unauthorized errors
type UnauthorizedSolution struct {
	*solution.BaseSolution
}

// NewUnauthorizedSolution creates a solution specifically for 401 errors
func NewUnauthorizedSolution() *UnauthorizedSolution {
	base := solution.NewBaseSolution("Authentication Required").
		SetSolutionDescription("This request requires authentication. Common causes and solutions:\n\n• Missing authentication token - Include your API key or bearer token\n• Expired session - Login again or refresh your token\n• Invalid credentials - Verify your username/password\n• Missing authentication headers - Check Authorization header format").
		AddDocumentationLink("HTTP 401 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401").
		AddDocumentationLink("GoVel Authentication", "https://govel.dev/docs/authentication").
		AddDocumentationLink("GoVel API Authentication", "https://govel.dev/docs/api-authentication").
		AddDocumentationLink("JWT Token Guide", "https://govel.dev/docs/jwt-tokens")

	return &UnauthorizedSolution{
		BaseSolution: base,
	}
}

// Ensure UnauthorizedSolution implements the Solution interface
var _ solutionInterface.Solution = (*UnauthorizedSolution)(nil)
