package http

import (
	"govel/exceptions/core/solution"
	solutionInterface "govel/exceptions/interfaces/solution"
)

// ForbiddenSolution provides specific guidance for 403 Forbidden errors
type ForbiddenSolution struct {
	*solution.BaseSolution
}

// NewForbiddenSolution creates a solution specifically for 403 errors
func NewForbiddenSolution(action string) *ForbiddenSolution {
	title := "Access Forbidden"
	description := "You don't have permission to access this resource. Common causes:\n\n• Insufficient user permissions\n• Resource requires admin access\n• Missing authorization scope\n• IP address restrictions"

	if action != "" {
		title = "Access Forbidden: " + action
		description = "You don't have permission to " + action + ". " + description
	}

	base := solution.NewBaseSolution(title).
		SetSolutionDescription(description).
		AddDocumentationLink("HTTP 403 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403").
		AddDocumentationLink("GoVel Authorization", "https://govel.dev/docs/authorization").
		AddDocumentationLink("GoVel Policies", "https://govel.dev/docs/policies").
		AddDocumentationLink("GoVel Gates", "https://govel.dev/docs/gates")

	return &ForbiddenSolution{
		BaseSolution: base,
	}
}

// Ensure ForbiddenSolution implements the Solution interface
var _ solutionInterface.Solution = (*ForbiddenSolution)(nil)
