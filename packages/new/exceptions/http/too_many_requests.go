package http

import (
	"fmt"
	"net/http"

	"govel/exceptions/core"
	"govel/exceptions/interfaces"
	httpSolutions "govel/exceptions/solutions/http"
)

// TooManyRequestsException represents a 429 Too Many Requests error.
// Used when the user has sent too many requests in a given time frame.
type TooManyRequestsException struct {
	*core.Exception
}

// NewTooManyRequestsException creates a new 429 Too Many Requests exception.
//
// Parameters:
//   message: Optional custom error message
//   retryAfter: Optional retry-after value in seconds
//
// Example:
//   err := http.NewTooManyRequestsException("Rate limit exceeded", 60)
func NewTooManyRequestsException(message string, retryAfter ...int) *TooManyRequestsException {
	if message == "" {
		message = "Too Many Requests"
	}

	exception := core.NewException(message, http.StatusTooManyRequests)

	// Add Retry-After header if specified
	if len(retryAfter) > 0 && retryAfter[0] > 0 {
		exception.WithHeader("Retry-After", fmt.Sprintf("%d", retryAfter[0]))
	}

	exc := &TooManyRequestsException{
		Exception: exception,
	}

	// Set solution for this exception with retry after value
	retryAfterValue := 0
	if len(retryAfter) > 0 {
		retryAfterValue = retryAfter[0]
	}
	exc.Exception.SetSolution(httpSolutions.NewTooManyRequestsSolution(retryAfterValue, 0))

	return exc
}

// Ensure TooManyRequestsException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*TooManyRequestsException)(nil)
