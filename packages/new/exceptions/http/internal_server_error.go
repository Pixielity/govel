package http

import (
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// InternalServerErrorException represents a 500 Internal Server Error.
// Used when the server encountered an unexpected condition.
type InternalServerErrorException struct {
	*core.Exception
}

// NewInternalServerErrorException creates a new 500 Internal Server Error exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewInternalServerErrorException("Database connection failed")
func NewInternalServerErrorException(message ...string) *InternalServerErrorException {
	msg := "Internal Server Error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &InternalServerErrorException{
		Exception: core.NewException(msg, http.StatusInternalServerError),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewInternalServerErrorSolution())

	return exc
}

// Ensure InternalServerErrorException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*InternalServerErrorException)(nil)
