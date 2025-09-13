package interfaces

// HTTPable defines the interface for exceptions that have HTTP-related functionality.
// This interface follows the Interface Segregation Principle (ISP) by focusing
// solely on HTTP-related functionality.
type HTTPable interface {
	// GetStatusCode returns the HTTP status code for this exception
	GetStatusCode() int

	// GetMessage returns the exception message
	GetMessage() string

	// GetHeaders returns HTTP headers that should be sent with the response
	GetHeaders() map[string]string

	// IsHTTPException returns true if this is an HTTP-related exception
	IsHTTPException() bool

	// WithHeader adds an HTTP header to the exception
	WithHeader(key, value string) ExceptionInterface

	// WithHeaders adds multiple HTTP headers to the exception
	WithHeaders(headers map[string]string) ExceptionInterface

	// WithMessage sets a custom message for the exception
	WithMessage(message string) ExceptionInterface
}
