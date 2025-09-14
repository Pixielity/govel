package http

import (
	"govel/exceptions/core/solution"
	solutionInterface "govel/exceptions/interfaces/solution"
)

// BadRequestSolution provides specific guidance for 400 Bad Request errors
type BadRequestSolution struct {
	*solution.BaseSolution
}

// NewBadRequestSolution creates a solution specifically for 400 errors
func NewBadRequestSolution() *BadRequestSolution {
	base := solution.NewBaseSolution("Bad Request").
		SetSolutionDescription("The server cannot process the request due to invalid syntax. Common issues:\n\n• Malformed JSON or XML in request body\n• Invalid query parameters\n• Missing required headers\n• Incorrect Content-Type header\n• Invalid URL encoding").
		AddDocumentationLink("HTTP 400 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400").
		AddDocumentationLink("GoVel Request Handling", "https://govel.dev/docs/requests").
		AddDocumentationLink("GoVel Validation", "https://govel.dev/docs/validation").
		AddDocumentationLink("JSON API Guidelines", "https://jsonapi.org/")

	return &BadRequestSolution{
		BaseSolution: base,
	}
}

// Ensure BadRequestSolution implements the Solution interface
var _ solutionInterface.Solution = (*BadRequestSolution)(nil)
