package helpers

import (
	"govel/exceptions/core"
	"govel/exceptions/interfaces"
)

// Abort creates and returns a new exception with the given status code and message.
// This provides Laravel-style exception throwing functionality.
//
// Parameters:
//   statusCode: The HTTP status code for the exception
//   message: The error message for the exception
//
// Returns:
//   interfaces.ExceptionInterface: A new exception instance
//
// Example:
//   err := helpers.Abort(404, "Resource not found")
func Abort(statusCode int, message string) interfaces.ExceptionInterface {
	return core.NewException(message, statusCode)
}

// AbortIf creates and returns a new exception if the given condition is true.
// If the condition is false, it returns nil.
//
// Parameters:
//   condition: Boolean condition to check
//   statusCode: The HTTP status code for the exception
//   message: The error message for the exception
//
// Returns:
//   interfaces.ExceptionInterface: A new exception if condition is true, nil otherwise
//
// Example:
//   err := helpers.AbortIf(user == nil, 404, "User not found")
//   if err != nil {
//       return err
//   }
func AbortIf(condition bool, statusCode int, message string) interfaces.ExceptionInterface {
	if condition {
		return Abort(statusCode, message)
	}
	return nil
}

// AbortUnless creates and returns a new exception if the given condition is false.
// If the condition is true, it returns nil.
//
// Parameters:
//   condition: Boolean condition to check
//   statusCode: The HTTP status code for the exception
//   message: The error message for the exception
//
// Returns:
//   interfaces.ExceptionInterface: A new exception if condition is false, nil otherwise
//
// Example:
//   err := helpers.AbortUnless(user.IsActive(), 403, "Account is inactive")
//   if err != nil {
//       return err
//   }
func AbortUnless(condition bool, statusCode int, message string) interfaces.ExceptionInterface {
	if !condition {
		return Abort(statusCode, message)
	}
	return nil
}
