package models

import (
	"fmt"
	"path/filepath"
	"strings"

	"govel/packages/ignition/constants"
)

// StackFrameTemplate represents a stack frame for template rendering
type StackFrameTemplate struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// NewStackFrameTemplate creates a new stack frame template
func NewStackFrameTemplate() *StackFrameTemplate {
	return &StackFrameTemplate{}
}

// NewStackFrameTemplateWithData creates a new stack frame template with data
func NewStackFrameTemplateWithData(function, file string, line int) *StackFrameTemplate {
	return &StackFrameTemplate{
		Function: function,
		File:     file,
		Line:     line,
	}
}

// GetFunction returns the function name
func (s *StackFrameTemplate) GetFunction() string {
	return s.Function
}

// SetFunction sets the function name
func (s *StackFrameTemplate) SetFunction(function string) {
	s.Function = function
}

// GetFile returns the file path
func (s *StackFrameTemplate) GetFile() string {
	return s.File
}

// SetFile sets the file path
func (s *StackFrameTemplate) SetFile(file string) {
	s.File = file
}

// GetLine returns the line number
func (s *StackFrameTemplate) GetLine() int {
	return s.Line
}

// SetLine sets the line number
func (s *StackFrameTemplate) SetLine(line int) {
	s.Line = line
}

// GetRelativeFile returns the file path relative to the application root
func (s *StackFrameTemplate) GetRelativeFile() string {
	// Remove leading slash for relative path display
	return strings.TrimPrefix(s.File, constants.UnixPathSeparator)
}

// GetFileName returns just the filename without the path
func (s *StackFrameTemplate) GetFileName() string {
	return filepath.Base(s.File)
}

// GetPackageName extracts the package name from the function
func (s *StackFrameTemplate) GetPackageName() string {
	parts := strings.Split(s.Function, ".")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

// GetMethodName extracts the method name from the function
func (s *StackFrameTemplate) GetMethodName() string {
	parts := strings.Split(s.Function, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return s.Function
}

// IsGoStdLib returns true if this frame is from Go standard library
func (s *StackFrameTemplate) IsGoStdLib() bool {
	return constants.IsGoStdLibPath(s.File)
}

// IsApplicationFrame returns true if this is likely an application frame
func (s *StackFrameTemplate) IsApplicationFrame() bool {
	return constants.IsApplicationPath(s.File)
}

// IsThirdParty returns true if this frame is from a third-party library
func (s *StackFrameTemplate) IsThirdParty() bool {
	return constants.IsThirdPartyPath(s.File)
}

// IsEmpty returns true if the stack frame template is empty
func (s *StackFrameTemplate) IsEmpty() bool {
	return s.Function == "" && s.File == "" && s.Line == 0
}

// HasValidLine returns true if the line number is valid (> 0)
func (s *StackFrameTemplate) HasValidLine() bool {
	return s.Line > 0
}

// GetShortFunction returns a shortened version of the function name
func (s *StackFrameTemplate) GetShortFunction(maxLength int) string {
	if len(s.Function) <= maxLength {
		return s.Function
	}

	// Try to keep the method name if possible
	parts := strings.Split(s.Function, ".")
	if len(parts) > 1 {
		methodName := parts[len(parts)-1]
		if len(methodName) < maxLength-3 {
			return "..." + methodName
		}
	}

	return s.Function[:maxLength-3] + "..."
}

// GetShortFile returns a shortened version of the file path
func (s *StackFrameTemplate) GetShortFile(maxLength int) string {
	if len(s.File) <= maxLength {
		return s.File
	}

	// Try to keep the filename if possible
	fileName := filepath.Base(s.File)
	if len(fileName) < maxLength-3 {
		return "..." + fileName
	}

	return s.File[:maxLength-3] + "..."
}

// ToString returns a string representation of the stack frame template
func (s *StackFrameTemplate) ToString() string {
	return fmt.Sprintf("%s:%d in %s", s.File, s.Line, s.Function)
}

// ToShortString returns a short string representation of the stack frame template
func (s *StackFrameTemplate) ToShortString() string {
	return fmt.Sprintf("%s:%d in %s", s.GetFileName(), s.Line, s.GetMethodName())
}

// Clone creates a copy of the stack frame template
func (s *StackFrameTemplate) Clone() *StackFrameTemplate {
	return &StackFrameTemplate{
		Function: s.Function,
		File:     s.File,
		Line:     s.Line,
	}
}
