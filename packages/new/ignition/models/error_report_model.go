package models

import (
	"time"
)

// ErrorReport represents a structured error report
type ErrorReport struct {
	Message   string       `json:"message"`
	ErrorType string       `json:"type"`
	File      string       `json:"file"`
	Line      int          `json:"line"`
	Stack     []StackFrame `json:"stack"`
	Context   ErrorContext `json:"context"`
	Solutions []Solution   `json:"solutions"`
	Timestamp time.Time    `json:"timestamp"`
}

// NewErrorReport creates a new error report
func NewErrorReport() *ErrorReport {
	return &ErrorReport{
		Timestamp: time.Now(),
		Solutions: []Solution{},
		Stack:     []StackFrame{},
	}
}

// GetMessage returns the error message
func (e *ErrorReport) GetMessage() string {
	return e.Message
}

// SetMessage sets the error message
func (e *ErrorReport) SetMessage(message string) {
	e.Message = message
}

// GetType returns the error type
func (e *ErrorReport) GetType() string {
	return e.ErrorType
}

// SetType sets the error type
func (e *ErrorReport) SetType(errorType string) {
	e.ErrorType = errorType
}

// GetFile returns the file where the error occurred
func (e *ErrorReport) GetFile() string {
	return e.File
}

// SetFile sets the file where the error occurred
func (e *ErrorReport) SetFile(file string) {
	e.File = file
}

// GetLine returns the line number where the error occurred
func (e *ErrorReport) GetLine() int {
	return e.Line
}

// SetLine sets the line number where the error occurred
func (e *ErrorReport) SetLine(line int) {
	e.Line = line
}

// GetStack returns the stack trace
func (e *ErrorReport) GetStack() []StackFrame {
	return e.Stack
}

// SetStack sets the stack trace
func (e *ErrorReport) SetStack(stack []StackFrame) {
	e.Stack = stack
}

// AddStackFrame adds a frame to the stack trace
func (e *ErrorReport) AddStackFrame(frame StackFrame) {
	e.Stack = append(e.Stack, frame)
}

// GetContext returns the error context
func (e *ErrorReport) GetContext() ErrorContext {
	return e.Context
}

// SetContext sets the error context
func (e *ErrorReport) SetContext(context ErrorContext) {
	e.Context = context
}

// GetSolutions returns the list of solutions
func (e *ErrorReport) GetSolutions() []Solution {
	return e.Solutions
}

// SetSolutions sets the list of solutions
func (e *ErrorReport) SetSolutions(solutions []Solution) {
	e.Solutions = solutions
}

// AddSolution adds a solution to the report
func (e *ErrorReport) AddSolution(solution Solution) {
	e.Solutions = append(e.Solutions, solution)
}

// GetTimestamp returns the timestamp when the error occurred
func (e *ErrorReport) GetTimestamp() time.Time {
	return e.Timestamp
}

// SetTimestamp sets the timestamp when the error occurred
func (e *ErrorReport) SetTimestamp(timestamp time.Time) {
	e.Timestamp = timestamp
}

// IsEmpty returns true if the error report is empty
func (e *ErrorReport) IsEmpty() bool {
	return e.Message == "" && e.ErrorType == ""
}

// HasSolutions returns true if the report has solutions
func (e *ErrorReport) HasSolutions() bool {
	return len(e.Solutions) > 0
}

// GetStackFrameCount returns the number of stack frames
func (e *ErrorReport) GetStackFrameCount() int {
	return len(e.Stack)
}

// GetApplicationFrames returns only application frames from the stack
func (e *ErrorReport) GetApplicationFrames() []StackFrame {
	var appFrames []StackFrame
	for _, frame := range e.Stack {
		if frame.IsApplicationFrame() {
			appFrames = append(appFrames, frame)
		}
	}
	return appFrames
}

// Note: Compile-time interface compliance check removed due to circular dependencies
// The ErrorReport interface references concrete types instead of interfaces
