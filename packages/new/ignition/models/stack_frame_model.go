package models

import (
	"fmt"
	"path/filepath"
	"strings"

	"govel/ignition/constants"
	"govel/ignition/interfaces"
)

// StackFrame represents a single frame in the stack trace
type StackFrame struct {
	Function         string            `json:"function"`
	File             string            `json:"file"`
	Line             int               `json:"line"`
	Code             map[string]string `json:"code"`
	ApplicationFrame bool              `json:"application_frame"`
}

// NewStackFrame creates a new stack frame
func NewStackFrame() *StackFrame {
	return &StackFrame{
		Code: make(map[string]string),
	}
}

// GetFunction returns the function name
func (s *StackFrame) GetFunction() string {
	return s.Function
}

// SetFunction sets the function name
func (s *StackFrame) SetFunction(function string) {
	s.Function = function
}

// GetFile returns the file path
func (s *StackFrame) GetFile() string {
	return s.File
}

// SetFile sets the file path
func (s *StackFrame) SetFile(file string) {
	s.File = file
}

// GetLine returns the line number
func (s *StackFrame) GetLine() int {
	return s.Line
}

// SetLine sets the line number
func (s *StackFrame) SetLine(line int) {
	s.Line = line
}

// GetCode returns the code snippet
func (s *StackFrame) GetCode() map[string]string {
	return s.Code
}

// SetCode sets the code snippet
func (s *StackFrame) SetCode(code map[string]string) {
	s.Code = code
}

// AddCodeLine adds a single line of code to the snippet
func (s *StackFrame) AddCodeLine(lineNumber string, lineCode string) {
	if s.Code == nil {
		s.Code = make(map[string]string)
	}
	s.Code[lineNumber] = lineCode
}

// IsApplicationFrame returns true if this is an application frame
func (s *StackFrame) IsApplicationFrame() bool {
	return s.ApplicationFrame
}

// SetApplicationFrame sets whether this is an application frame
func (s *StackFrame) SetApplicationFrame(isAppFrame bool) {
	s.ApplicationFrame = isAppFrame
}

// GetRelativeFile returns the file path relative to the application root
func (s *StackFrame) GetRelativeFile() string {
	// Remove leading slash for relative path display
	return strings.TrimPrefix(s.File, constants.UnixPathSeparator)
}

// GetFileName returns just the filename without the path
func (s *StackFrame) GetFileName() string {
	return filepath.Base(s.File)
}

// GetPackageName extracts the package name from the function
func (s *StackFrame) GetPackageName() string {
	parts := strings.Split(s.Function, ".")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

// GetMethodName extracts the method name from the function
func (s *StackFrame) GetMethodName() string {
	parts := strings.Split(s.Function, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return s.Function
}

// IsGoStdLib returns true if this frame is from Go standard library
func (s *StackFrame) IsGoStdLib() bool {
	return constants.IsGoStdLibPath(s.File)
}

// IsThirdParty returns true if this frame is from a third-party library
func (s *StackFrame) IsThirdParty() bool {
	return constants.IsThirdPartyPath(s.File)
}

// HasCode returns true if this frame has code snippet
func (s *StackFrame) HasCode() bool {
	return len(s.Code) > 0
}

// GetCodeLineCount returns the number of code lines
func (s *StackFrame) GetCodeLineCount() int {
	return len(s.Code)
}

// ToString returns a string representation of the stack frame
func (s *StackFrame) ToString() string {
	return fmt.Sprintf("%s:%d in %s", s.File, s.Line, s.Function)
}

// Compile-time interface compliance check
var _ interfaces.StackFrameInterface = (*StackFrame)(nil)
