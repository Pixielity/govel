package http

import (
	"net/http"

	"govel/packages/exceptions/core"
	"govel/packages/exceptions/interfaces"
	httpSolutions "govel/packages/exceptions/solutions/http"
)

// ForbiddenException represents a 403 Forbidden error.
// Used when the server understands the request but refuses to authorize it.
type ForbiddenException struct {
	*core.Exception
}

// NewForbiddenException creates a new 403 Forbidden exception.
//
// Parameters:
//   message: Optional custom error message
//
// Example:
//   err := http.NewForbiddenException("Access denied to this resource")
func NewForbiddenException(message ...string) *ForbiddenException {
	msg := "Forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exc := &ForbiddenException{
		Exception: core.NewException(msg, http.StatusForbidden),
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewForbiddenSolution(""))

	return exc
}

// Ensure ForbiddenException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*ForbiddenException)(nil)
