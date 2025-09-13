package http

import (
	"net/http"

	"govel/exceptions/core"
	"govel/exceptions/interfaces"
	httpSolutions "govel/exceptions/solutions/http"
)

// BadRequestException represents a 400 Bad Request error.
// Used when the server cannot understand the request due to invalid syntax.
type BadRequestException struct {
	*core.Exception
}

// NewBadRequestException creates a new 400 Bad Request exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewBadRequestException("Invalid JSON format")
func NewBadRequestException(message ...string) *BadRequestException {
	msg := "Bad Request"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &BadRequestException{
		Exception: core.NewException(msg, http.StatusBadRequest),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewBadRequestSolution())

	return exc
}

// Ensure BadRequestException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*BadRequestException)(nil)
