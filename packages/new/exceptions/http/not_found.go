package http

import (
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// NotFoundException represents a 404 Not Found error.
// Used when the requested resource could not be found on the server.
type NotFoundException struct {
	*core.Exception
}

// NewNotFoundException creates a new 404 Not Found exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewNotFoundException("User not found")
func NewNotFoundException(message ...string) *NotFoundException {
	msg := "Not Found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &NotFoundException{
		Exception: core.NewException(msg, http.StatusNotFound),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewNotFoundSolution(""))

	return exc
}

// Ensure NotFoundException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*NotFoundException)(nil)
