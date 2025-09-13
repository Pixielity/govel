package solution

import "govel/packages/exceptions/interfaces/solution"

// BaseSolution provides a basic implementation of the Solution interface.
// This is inspired by Laravel's Spatie\ErrorSolutions\Interfaces\BaseSolution class.
type BaseSolution struct {
	title       string
	description string
	links       map[string]string
}

// NewBaseSolution creates a new BaseSolution with the given title
func NewBaseSolution(title string) *BaseSolution {
	return &BaseSolution{
		title:       title,
		description: "",
		links:       make(map[string]string),
	}
}

// GetSolutionTitle returns the title of the solution
func (s *BaseSolution) GetSolutionTitle() string {
	return s.title
}

// SetSolutionTitle sets the title of the solution (for method chaining)
func (s *BaseSolution) SetSolutionTitle(title string) *BaseSolution {
	s.title = title
	return s
}

// GetSolutionDescription returns the description of the solution
func (s *BaseSolution) GetSolutionDescription() string {
	return s.description
}

// SetSolutionDescription sets the description of the solution (for method chaining)
func (s *BaseSolution) SetSolutionDescription(description string) *BaseSolution {
	s.description = description
	return s
}

// GetDocumentationLinks returns the documentation links
func (s *BaseSolution) GetDocumentationLinks() map[string]string {
	return s.links
}

// SetDocumentationLinks sets the documentation links (for method chaining)
func (s *BaseSolution) SetDocumentationLinks(links map[string]string) *BaseSolution {
	s.links = links
	return s
}

// AddDocumentationLink adds a single documentation link (for method chaining)
func (s *BaseSolution) AddDocumentationLink(name, url string) *BaseSolution {
	if s.links == nil {
		s.links = make(map[string]string)
	}
	s.links[name] = url
	return s
}

// Ensure BaseSolution implements the Solution interface
var _ solution.Solution = (*BaseSolution)(nil)
