package http

import (
	"fmt"
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// ServiceUnavailableException represents a 503 Service Unavailable error.
// Used when the server is temporarily overloaded or under maintenance.
type ServiceUnavailableException struct {
	*core.Exception
}

// NewServiceUnavailableException creates a new 503 Service Unavailable exception.
//
// Parameters:
//   message: Optional custom error message
//   retryAfter: Optional retry-after value in seconds
//
// Example:
//   err := http.NewServiceUnavailableException("Server under maintenance", 3600)
func NewServiceUnavailableException(message string, retryAfter ...int) *ServiceUnavailableException {
	if message == "" {
		message = "Service Unavailable"
	}

	exception := core.NewException(message, http.StatusServiceUnavailable)

	// Add Retry-After header if specified
	if len(retryAfter) > 0 && retryAfter[0] > 0 {
		exception.WithHeader("Retry-After", fmt.Sprintf("%d", retryAfter[0]))
	}

	exc := &ServiceUnavailableException{
		Exception: exception,
	}

	// Set solution for this exception with retry after value
	retryAfterValue := 0
	if len(retryAfter) > 0 {
		retryAfterValue = retryAfter[0]
	}
	exc.Exception.SetSolution(httpSolutions.NewServiceUnavailableSolution(retryAfterValue))

	return exc
}

// Ensure ServiceUnavailableException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*ServiceUnavailableException)(nil)
