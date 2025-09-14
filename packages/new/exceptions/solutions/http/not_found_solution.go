package http

import (
	"govel/exceptions/core/solution"
	solutionInterface "govel/exceptions/interfaces/solution"
)

// NotFoundSolution provides specific guidance for 404 Not Found errors
type NotFoundSolution struct {
	*solution.BaseSolution
}

// NewNotFoundSolution creates a solution specifically for 404 errors
func NewNotFoundSolution(resource string) *NotFoundSolution {
	title := "Resource Not Found"
	if resource != "" {
		title = "Resource '" + resource + "' Not Found"
	}

	base := solution.NewBaseSolution(title).
		SetSolutionDescription("The requested resource could not be found. This typically happens when:\n\n• The URL path is incorrect\n• The resource has been moved or deleted\n• The route is not properly defined\n• The resource requires authentication").
		AddDocumentationLink("HTTP 404 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404").
		AddDocumentationLink("GoVel Routing Documentation", "https://govel.dev/docs/routing").
		AddDocumentationLink("GoVel Resource Controllers", "https://govel.dev/docs/controllers")

	return &NotFoundSolution{
		BaseSolution: base,
	}
}

// Ensure NotFoundSolution implements the Solution interface
var _ solutionInterface.Solution = (*NotFoundSolution)(nil)
