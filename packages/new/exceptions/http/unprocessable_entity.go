package http

import (
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// UnprocessableEntityException represents a 422 Unprocessable Entity error.
// Used when the server understands the content but cannot process the instructions.
type UnprocessableEntityException struct {
	*core.Exception
}

// NewUnprocessableEntityException creates a new 422 Unprocessable Entity exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewUnprocessableEntityException("Validation failed")
func NewUnprocessableEntityException(message ...string) *UnprocessableEntityException {
	msg := "Unprocessable Entity"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &UnprocessableEntityException{
		Exception: core.NewException(msg, http.StatusUnprocessableEntity),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewValidationErrorSolution(nil))

	return exc
}

// Ensure UnprocessableEntityException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*UnprocessableEntityException)(nil)
