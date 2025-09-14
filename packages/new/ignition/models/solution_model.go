package models

import (
	"strings"

	"govel/ignition/interfaces"
)

// Solution represents a suggested fix for an error matching Laravel Ignition format
type Solution struct {
	Class             string            `json:"class"`
	Title             string            `json:"title"`
	Links             map[string]string `json:"links"`
	Description       string            `json:"description"`
	IsRunnable        bool              `json:"is_runnable"`
	AiGenerated       bool              `json:"ai_generated"`
	ActionDescription string            `json:"action_description"`
	RunButtonText     string            `json:"run_button_text"`
	ExecuteEndpoint   string            `json:"execute_endpoint"`
	RunParameters     []interface{}     `json:"run_parameters"`
}

// NewSolution creates a new solution
func NewSolution(title, description string) *Solution {
	return &Solution{
		Class:             "GoVel\\Solutions\\Generic",
		Title:             title,
		Links:             make(map[string]string),
		Description:       description,
		IsRunnable:        false,
		AiGenerated:       false,
		ActionDescription: "",
		RunButtonText:     "",
		ExecuteEndpoint:   "",
		RunParameters:     []interface{}{},
	}
}

// NewGoDocumentationSolution creates a solution with Go documentation links
func NewGoDocumentationSolution(title, description string) *Solution {
	solution := NewSolution(title, description)
	solution.Class = "GoVel\\Solutions\\Documentation"
	solution.AddLink("Go Documentation", "https://golang.org/doc/")
	solution.AddLink("Go Standard Library", "https://golang.org/pkg/")
	return solution
}

// NewGoModuleSolution creates a solution for Go module issues
func NewGoModuleSolution(title, description string) *Solution {
	solution := NewSolution(title, description)
	solution.Class = "GoVel\\Solutions\\GoModule"
	solution.AddLink("Go Modules", "https://golang.org/doc/modules/gomod-ref")
	return solution
}

// NewRunnableSolution creates a runnable solution
func NewRunnableSolution(title, description, actionDesc, buttonText, endpoint string) *Solution {
	solution := NewSolution(title, description)
	solution.Class = "GoVel\\Solutions\\Runnable"
	solution.IsRunnable = true
	solution.ActionDescription = actionDesc
	solution.RunButtonText = buttonText
	solution.ExecuteEndpoint = endpoint
	return solution
}

// GetTitle returns the solution title
func (s *Solution) GetTitle() string {
	return s.Title
}

// SetTitle sets the solution title
func (s *Solution) SetTitle(title string) {
	s.Title = title
}

// GetDescription returns the solution description
func (s *Solution) GetDescription() string {
	return s.Description
}

// SetDescription sets the solution description
func (s *Solution) SetDescription(description string) {
	s.Description = description
}

// GetClass returns the solution class
func (s *Solution) GetClass() string {
	return s.Class
}

// SetClass sets the solution class
func (s *Solution) SetClass(class string) {
	s.Class = class
}

// GetLinks returns the documentation links
func (s *Solution) GetLinks() map[string]string {
	return s.Links
}

// SetLinks sets the documentation links
func (s *Solution) SetLinks(links map[string]string) {
	s.Links = links
}

// AddLink adds a documentation link
func (s *Solution) AddLink(label, url string) {
	if s.Links == nil {
		s.Links = make(map[string]string)
	}
	s.Links[label] = url
}

// GetIsRunnable returns whether the solution is runnable
func (s *Solution) GetIsRunnable() bool {
	return s.IsRunnable
}

// SetIsRunnable sets whether the solution is runnable
func (s *Solution) SetIsRunnable(runnable bool) {
	s.IsRunnable = runnable
}

// GetAiGenerated returns whether the solution is AI generated
func (s *Solution) GetAiGenerated() bool {
	return s.AiGenerated
}

// SetAiGenerated sets whether the solution is AI generated
func (s *Solution) SetAiGenerated(aiGenerated bool) {
	s.AiGenerated = aiGenerated
}

// GetActionDescription returns the action description
func (s *Solution) GetActionDescription() string {
	return s.ActionDescription
}

// SetActionDescription sets the action description
func (s *Solution) SetActionDescription(description string) {
	s.ActionDescription = description
}

// GetRunButtonText returns the run button text
func (s *Solution) GetRunButtonText() string {
	return s.RunButtonText
}

// SetRunButtonText sets the run button text
func (s *Solution) SetRunButtonText(text string) {
	s.RunButtonText = text
}

// GetExecuteEndpoint returns the execute endpoint
func (s *Solution) GetExecuteEndpoint() string {
	return s.ExecuteEndpoint
}

// SetExecuteEndpoint sets the execute endpoint
func (s *Solution) SetExecuteEndpoint(endpoint string) {
	s.ExecuteEndpoint = endpoint
}

// GetRunParameters returns the run parameters
func (s *Solution) GetRunParameters() []interface{} {
	return s.RunParameters
}

// SetRunParameters sets the run parameters
func (s *Solution) SetRunParameters(parameters []interface{}) {
	s.RunParameters = parameters
}

// HasLinks returns true if the solution has documentation links
func (s *Solution) HasLinks() bool {
	return len(s.Links) > 0
}

// GetLinkCount returns the number of documentation links
func (s *Solution) GetLinkCount() int {
	return len(s.Links)
}

// HasRunCode returns true if the solution is runnable
func (s *Solution) HasRunCode() bool {
	return s.IsRunnable
}

// IsEmpty returns true if the solution is empty
func (s *Solution) IsEmpty() bool {
	return s.Title == "" && s.Description == ""
}

// GetShortDescription returns a shortened version of the description
func (s *Solution) GetShortDescription(maxLength int) string {
	if len(s.Description) <= maxLength {
		return s.Description
	}

	// Find the last space before the max length
	shortened := s.Description[:maxLength]
	lastSpace := strings.LastIndex(shortened, " ")
	if lastSpace > 0 {
		shortened = shortened[:lastSpace]
	}

	return shortened + "..."
}

// ContainsKeyword returns true if the solution contains the given keyword
func (s *Solution) ContainsKeyword(keyword string) bool {
	keyword = strings.ToLower(keyword)
	return strings.Contains(strings.ToLower(s.Title), keyword) ||
		strings.Contains(strings.ToLower(s.Description), keyword)
}

// GetWordCount returns the word count of the description
func (s *Solution) GetWordCount() int {
	if s.Description == "" {
		return 0
	}
	return len(strings.Fields(s.Description))
}

// ToString returns a string representation of the solution
func (s *Solution) ToString() string {
	result := s.Title
	if s.Description != "" {
		result += ": " + s.Description
	}
	return result
}

// Compile-time interface compliance check
var _ interfaces.SolutionInterface = (*Solution)(nil)
