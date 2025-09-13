package helpers

import (
	"net/http"

	"govel/exceptions/interfaces"
)

// Abort400 creates a new 400 Bad Request exception.
//
// Parameters:
//   message: Optional custom error message
//
// Returns:
//   interfaces.ExceptionInterface: A new 400 exception
//
// Example:
//   err := helpers.Abort400("Invalid request format")
func Abort400(message ...string) interfaces.ExceptionInterface {
	msg := "Bad Request"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Abort(http.StatusBadRequest, msg)
}

// Abort401 creates a new 401 Unauthorized exception.
//
// Parameters:
//   message: Optional custom error message
//
// Returns:
//   interfaces.ExceptionInterface: A new 401 exception
//
// Example:
//   err := helpers.Abort401("Authentication required")
func Abort401(message ...string) interfaces.ExceptionInterface {
	msg := "Unauthorized"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Abort(http.StatusUnauthorized, msg)
}

// Abort403 creates a new 403 Forbidden exception.
//
// Parameters:
//   message: Optional custom error message
//
// Returns:
//   interfaces.ExceptionInterface: A new 403 exception
//
// Example:
//   err := helpers.Abort403("Access denied")
func Abort403(message ...string) interfaces.ExceptionInterface {
	msg := "Forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Abort(http.StatusForbidden, msg)
}

// Abort404 creates a new 404 Not Found exception.
//
// Parameters:
//   message: Optional custom error message
//
// Returns:
//   interfaces.ExceptionInterface: A new 404 exception
//
// Example:
//   err := helpers.Abort404("Resource not found")
func Abort404(message ...string) interfaces.ExceptionInterface {
	msg := "Not Found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Abort(http.StatusNotFound, msg)
}

// Abort500 creates a new 500 Internal Server Error exception.
//
// Parameters:
//   message: Optional custom error message
//
// Returns:
//   interfaces.ExceptionInterface: A new 500 exception
//
// Example:
//   err := helpers.Abort500("Server error occurred")
func Abort500(message ...string) interfaces.ExceptionInterface {
	msg := "Internal Server Error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return Abort(http.StatusInternalServerError, msg)
}
