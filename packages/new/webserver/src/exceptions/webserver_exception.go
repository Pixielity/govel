// Package exceptions - Webserver exceptions
package exceptions

import "fmt"

// WebserverException represents a generic webserver error.
type WebserverException struct {
	Message string
}

func (e *WebserverException) Error() string  { return e.Message }
func NewWebserverException(msg string) error { return &WebserverException{Message: msg} }
func NewWebserverExceptionf(format string, a ...interface{}) error {
	return &WebserverException{Message: fmt.Sprintf(format, a...)}
}
