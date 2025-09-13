// Package interfaces - Response interface definition
// This file defines the ResponseInterface contract that provides a unified API
// for creating HTTP responses across different web frameworks (GoFiber, Gin, Echo).
package interfaces

import (
	"io"
	"net/http"
)

// ResponseInterface defines the contract for HTTP response handling in the webserver package.
// This interface enables building responses in a framework-agnostic manner with a Laravel-inspired API.
//
// The interface supports:
//   - Status codes and headers
//   - JSON, Text, HTML, and custom body responses
//   - File downloads and file responses
//   - Redirects and streaming
//   - Cookie management
//
// Example usage:
//   return res.Status(201).Json(map[string]string{"message": "created"})
//   return res.Header("X-Foo", "Bar").Text("ok")
//   return res.Download("/path/to/file.zip")
type ResponseInterface interface {
	// Status sets the HTTP status code for the response.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   code: The HTTP status code (e.g., 200, 404, 500)
	//
	// Returns:
	//   ResponseInterface: The response instance for chaining
	Status(code int) ResponseInterface
	
	// Header sets an HTTP header.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   key: The header name
	//   value: The header value
	Header(key, value string) ResponseInterface
	
	// Headers sets multiple headers at once from a map.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   headers: Map of header names to values
	Headers(headers map[string]string) ResponseInterface
	
	// RemoveHeader removes an HTTP header if present.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   key: The header name to remove
	RemoveHeader(key string) ResponseInterface
	
	// Json sends a JSON response.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   payload: Any value that can be serialized to JSON
	Json(payload interface{}) ResponseInterface
	
	// Text sends a plain text response.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   body: The plain text body
	Text(body string) ResponseInterface
	
	// HTML sends an HTML response.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   html: The HTML body
	HTML(html string) ResponseInterface
	
	// Send sends a custom body with a specified content type.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   body: The raw body bytes
	//   contentType: The Content-Type header value
	Send(body []byte, contentType string) ResponseInterface
	
	// Stream streams data to the client using an io.Reader.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   reader: The io.Reader providing the stream data
	//   contentType: The Content-Type header value
	Stream(reader io.Reader, contentType string) ResponseInterface
	
	// File serves a file from the filesystem.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   path: The filesystem path to the file
	File(path string) ResponseInterface
	
	// Download serves a file as an attachment, prompting browser download.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   path: The filesystem path to the file
	//   filename: Optional filename for the downloaded file
	Download(path string, filename ...string) ResponseInterface
	
	// Redirect sends an HTTP redirect to the specified URL with optional status code.
	// Defaults to 302 Found when no status is provided.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   url: The destination URL
	//   status: Optional status code (e.g., 301, 302, 307)
	Redirect(url string, status ...int) ResponseInterface
	
	// Cookie sets a cookie on the response.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   cookie: The http.Cookie to set
	Cookie(cookie *http.Cookie) ResponseInterface
	
	// ClearCookie removes a cookie by name.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   name: The cookie name to clear
	ClearCookie(name string) ResponseInterface
	
	// NoContent sets the response to have no body with the given status code (default 204).
	// Returns the response for method chaining.
	//
	// Parameters:
	//   status: Optional status code (default 204)
	NoContent(status ...int) ResponseInterface
	
	// RawWriter exposes a low-level writer callback to write directly to the underlying response.
	// Use with care as this bypasses higher-level abstractions.
	// Returns the response for method chaining.
	//
	// Parameters:
	//   writer: A function that receives the underlying http.ResponseWriter (or equivalent)
	RawWriter(writer func(w http.ResponseWriter)) ResponseInterface
}
