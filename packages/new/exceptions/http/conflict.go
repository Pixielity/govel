package http

import (
	"net/http"

	"govel/exceptions/core"
	"govel/exceptions/interfaces"
	httpSolutions "govel/exceptions/solutions/http"
)

// ConflictException represents a 409 Conflict error.
// Used when the request conflicts with the current state of the server.
type ConflictException struct {
	*core.Exception
}

// NewConflictException creates a new 409 Conflict exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewConflictException("Resource already exists")
func NewConflictException(message ...string) *ConflictException {
	msg := "Conflict"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &ConflictException{
		Exception: core.NewException(msg, http.StatusConflict),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewConflictSolution())

	return exc
}

// Ensure ConflictException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*ConflictException)(nil)
