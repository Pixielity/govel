package http

import (
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// MethodNotAllowedException represents a 405 Method Not Allowed error.
// Used when the HTTP method is not supported for the requested resource.
type MethodNotAllowedException struct {
	*core.Exception
}

// NewMethodNotAllowedException creates a new 405 Method Not Allowed exception.
//
// Parameters:
//   message: Optional custom error message
//   allowedMethods: List of allowed HTTP methods
//
// Example:
//   err := http.NewMethodNotAllowedException("Method not allowed", "GET", "POST")
func NewMethodNotAllowedException(message string, allowedMethods ...string) *MethodNotAllowedException {
	if message == "" {
		message = "Method Not Allowed"
	}

	exception := core.NewException(message, http.StatusMethodNotAllowed)

	// Add Allow header with permitted methods
	if len(allowedMethods) > 0 {
		allow := ""
		for i, method := range allowedMethods {
			if i > 0 {
				allow += ", "
			}
			allow += method
		}
		exception.WithHeader("Allow", allow)
	}

	exc := &MethodNotAllowedException{
		Exception: exception,
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewMethodNotAllowedSolution("", allowedMethods))

	return exc
}

// Ensure MethodNotAllowedException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*MethodNotAllowedException)(nil)
