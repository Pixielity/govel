package http

import (
	"net/http"

	"govel/exceptions/core"
	"govel/exceptions/interfaces"
	httpSolutions "govel/exceptions/solutions/http"
)

// UnauthorizedException represents a 401 Unauthorized error.
// Used when authentication is required and has failed or not been provided.
type UnauthorizedException struct {
	*core.Exception
}

// NewUnauthorizedException creates a new 401 Unauthorized exception.
//
// Parameters:
//
//	message: Optional custom error message
//
// Example:
//
//	err := http.NewUnauthorizedException("Authentication required")
func NewUnauthorizedException(message ...string) *UnauthorizedException {
	msg := "Unauthorized"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	exception := core.NewException(msg, http.StatusUnauthorized)
	// Add WWW-Authenticate header for 401 responses
	exception.WithHeader("WWW-Authenticate", "Bearer")

	exc := &UnauthorizedException{
		Exception: exception,
	}

	// Set solution for this exception
	exc.Exception.SetSolution(httpSolutions.NewUnauthorizedSolution())

	return exc
}

// Ensure UnauthorizedException implements the ExceptionInterface
var _ interfaces.ExceptionInterface = (*UnauthorizedException)(nil)
