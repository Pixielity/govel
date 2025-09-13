// Package core provides the base exception implementation for GoVel applications.
// This package follows Laravel's exception pattern, providing HTTP-aware exceptions with
// centralized handling, flexible rendering, and conditional throwing capabilities.
package core

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"govel/packages/exceptions/interfaces"
	solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// Exception is the base exception struct that implements ExceptionInterface.
// It provides Laravel-like exception functionality adapted for Go.
type Exception struct {
	// StatusCode is the HTTP status code for this exception
	StatusCode int

	// Message is the human-readable error message
	Message string

	// Headers contains HTTP headers to be sent with the response
	Headers map[string]string

	// Context contains additional context information about the exception
	Context map[string]interface{}

	// StackTrace contains the call stack where the exception occurred
	StackTrace []string

	// Timestamp is when the exception was created
	Timestamp time.Time

	// IsHTTP indicates if this is an HTTP-related exception
	IsHTTP bool

	// Solution is the solution for this exception
	Solution solutionInterface.Solution
}

// NewException creates a new base exception with the given message and status code.
//
// Parameters:
//   message: The error message
//   statusCode: The HTTP status code (defaults to 500 if not provided)
//
// Returns:
//   *Exception: A new exception instance
//
// Example:
//   err := core.NewException("Something went wrong", 500)
func NewException(message string, statusCode ...int) *Exception {
	code := http.StatusInternalServerError
	if len(statusCode) > 0 && statusCode[0] > 0 {
		code = statusCode[0]
	}

	return &Exception{
		StatusCode: code,
		Message:    message,
		Headers:    make(map[string]string),
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(),
		Timestamp:  time.Now(),
		IsHTTP:     isHTTPStatusCode(code),
	}
}

// =============================================================================
// Error Interface Implementation
// =============================================================================

// Error implements the error interface, returning the exception message.
func (e *Exception) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// =============================================================================
// HTTPable Interface Implementation
// =============================================================================

// GetStatusCode returns the HTTP status code for this exception.
func (e *Exception) GetStatusCode() int {
	return e.StatusCode
}

// GetMessage returns the exception message.
func (e *Exception) GetMessage() string {
	return e.Message
}

// GetHeaders returns HTTP headers that should be sent with the response.
func (e *Exception) GetHeaders() map[string]string {
	if e.Headers == nil {
		e.Headers = make(map[string]string)
	}
	return e.Headers
}

// IsHTTPException returns true if this is an HTTP-related exception.
func (e *Exception) IsHTTPException() bool {
	return e.IsHTTP
}

// WithHeader adds an HTTP header to the exception and returns the exception.
func (e *Exception) WithHeader(key, value string) interfaces.ExceptionInterface {
	if e.Headers == nil {
		e.Headers = make(map[string]string)
	}
	e.Headers[key] = value
	return e
}

// WithHeaders adds multiple HTTP headers to the exception and returns the exception.
func (e *Exception) WithHeaders(headers map[string]string) interfaces.ExceptionInterface {
	if e.Headers == nil {
		e.Headers = make(map[string]string)
	}
	for key, value := range headers {
		e.Headers[key] = value
	}
	return e
}

// WithMessage sets a custom message for the exception and returns the exception.
func (e *Exception) WithMessage(message string) interfaces.ExceptionInterface {
	e.Message = message
	return e
}

// =============================================================================
// Contextable Interface Implementation
// =============================================================================

// GetContext returns additional context information about the exception.
func (e *Exception) GetContext() map[string]interface{} {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	return e.Context
}

// SetContext sets additional context information.
func (e *Exception) SetContext(context map[string]interface{}) interfaces.ExceptionInterface {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	for key, value := range context {
		e.Context[key] = value
	}
	return e
}

// WithContext adds context information to the exception and returns the exception.
func (e *Exception) WithContext(key string, value interface{}) interfaces.ExceptionInterface {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// =============================================================================
// Stackable Interface Implementation
// =============================================================================

// GetStackTrace returns the stack trace where the exception occurred.
func (e *Exception) GetStackTrace() []string {
	return e.StackTrace
}

// =============================================================================
// Solutionable Interface Implementation
// =============================================================================

// GetSolution returns the solution for this exception.
func (e *Exception) GetSolution() solutionInterface.Solution {
	return e.Solution
}

// HasSolution returns true if this exception has an associated solution.
func (e *Exception) HasSolution() bool {
	return e.Solution != nil
}

// SetSolution sets a solution for this exception.
func (e *Exception) SetSolution(solution solutionInterface.Solution) interfaces.ExceptionInterface {
	e.Solution = solution
	return e
}

// WithSolution sets a solution for this exception and returns the exception.
func (e *Exception) WithSolution(solution solutionInterface.Solution) interfaces.ExceptionInterface {
	e.Solution = solution
	return e
}

// =============================================================================
// Renderable Interface Implementation
// =============================================================================

// Render returns a response representation of the exception.
// This method provides a Laravel-like rendering system for exceptions.
func (e *Exception) Render() map[string]interface{} {
	response := map[string]interface{}{
		"error":       true,
		"status_code": e.StatusCode,
		"message":     e.GetMessage(),
		"timestamp":   e.Timestamp.Format(time.RFC3339),
	}

	// Add status text if message is empty
	if e.Message == "" {
		response["status_text"] = http.StatusText(e.StatusCode)
	}

	// Add context if available
	if len(e.Context) > 0 {
		response["context"] = e.Context
	}

	// Add headers if available
	if len(e.Headers) > 0 {
		response["headers"] = e.Headers
	}

	// Add solution if available
	if e.Solution != nil {
		response["solution"] = map[string]interface{}{
			"title":       e.Solution.GetSolutionTitle(),
			"description": e.Solution.GetSolutionDescription(),
			"links":       e.Solution.GetDocumentationLinks(),
		}

		// Add runnable solution information if applicable
		if runnableSolution, ok := e.Solution.(solutionInterface.RunnableSolution); ok {
			response["solution"].(map[string]interface{})["runnable"] = true
			response["solution"].(map[string]interface{})["action_description"] = runnableSolution.GetSolutionActionDescription()
			response["solution"].(map[string]interface{})["run_button_text"] = runnableSolution.GetRunButtonText()
			response["solution"].(map[string]interface{})["run_parameters"] = runnableSolution.GetRunParameters()
		} else {
			response["solution"].(map[string]interface{})["runnable"] = false
		}
	}

	return response
}

// =============================================================================
// Helper Functions
// =============================================================================

// captureStackTrace captures the current stack trace, excluding exception package frames.
func captureStackTrace() []string {
	var stack []string

	// Get up to 32 stack frames
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs) // Skip runtime.Callers, captureStackTrace, and NewException

	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()

		// Skip frames from the exceptions package itself
		if !strings.Contains(frame.File, "exceptions/") {
			stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		}

		if !more {
			break
		}
	}

	return stack
}

// isHTTPStatusCode checks if the given code is a valid HTTP status code.
func isHTTPStatusCode(code int) bool {
	return code >= 100 && code < 600
}

// Ensure Exception implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*Exception)(nil)
