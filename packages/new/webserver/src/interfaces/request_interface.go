// Package interfaces - Request interface definition
// This file defines the RequestInterface contract that provides a unified API
// for handling HTTP requests across different web frameworks (GoFiber, Gin, Echo).
package interfaces

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// RequestInterface defines the contract for HTTP request handling in the webserver package.
// This interface provides a Laravel-inspired API for accessing request data, headers, parameters,
// files, and other HTTP request information in a framework-agnostic way.
//
// The interface supports:
//   - URL parameters and route parameters
//   - Query string parameters  
//   - Request body parsing (JSON, form data, raw)
//   - File uploads (single and multiple)
//   - HTTP headers access
//   - Request metadata (method, URL, IP, user agent)
//   - Cookies handling
//   - Request validation helpers
//
// Example usage:
//   func userHandler(req RequestInterface) ResponseInterface {
//       id := req.Param("id")
//       name := req.Input("name")
//       file := req.File("avatar")
//       return Json(map[string]string{"user_id": id})
//   }
type RequestInterface interface {
	// URL Parameters and Route Parameters
	
	// Param retrieves a route parameter by name.
	// Route parameters are defined in the route path (e.g., "/users/:id").
	// Returns an empty string if the parameter doesn't exist.
	//
	// Parameters:
	//   key: The parameter name (without the ":" prefix)
	//
	// Returns:
	//   string: The parameter value, or empty string if not found
	//
	// Example:
	//   // For route "/users/:id" and request "/users/123"
	//   id := req.Param("id") // Returns "123"
	Param(key string) string
	
	// ParamInt retrieves a route parameter as an integer.
	// Returns 0 if the parameter doesn't exist or cannot be converted to int.
	//
	// Parameters:
	//   key: The parameter name
	//
	// Returns:
	//   int: The parameter value as integer, or 0 if not found/invalid
	ParamInt(key string) int
	
	// Params retrieves all route parameters as a map.
	//
	// Returns:
	//   map[string]string: All route parameters
	Params() map[string]string
	
	// Query String Parameters
	
	// Query retrieves a query string parameter by name.
	// Returns an empty string if the parameter doesn't exist.
	//
	// Parameters:
	//   key: The query parameter name
	//   defaultValue: Optional default value if parameter is not found
	//
	// Returns:
	//   string: The query parameter value, or default value if provided, or empty string
	//
	// Example:
	//   // For URL "/users?page=2&limit=10"
	//   page := req.Query("page")     // Returns "2"
	//   sort := req.Query("sort", "id") // Returns "id" (default)
	Query(key string, defaultValue ...string) string
	
	// QueryInt retrieves a query parameter as an integer.
	// Returns the default value if parameter doesn't exist or cannot be converted.
	//
	// Parameters:
	//   key: The query parameter name
	//   defaultValue: Optional default value
	//
	// Returns:
	//   int: The query parameter value as integer, or default value
	QueryInt(key string, defaultValue ...int) int
	
	// QueryBool retrieves a query parameter as a boolean.
	// Accepts "true", "1", "yes", "on" as true values (case-insensitive).
	//
	// Parameters:
	//   key: The query parameter name
	//   defaultValue: Optional default value
	//
	// Returns:
	//   bool: The query parameter value as boolean, or default value
	QueryBool(key string, defaultValue ...bool) bool
	
	// Queries retrieves all query parameters as a map.
	//
	// Returns:
	//   url.Values: All query parameters
	Queries() url.Values
	
	// Request Body and Input Data
	
	// Input retrieves input data from the request body (JSON, form data).
	// This method automatically detects the content type and parses accordingly.
	// Returns an empty string if the key doesn't exist.
	//
	// Parameters:
	//   key: The input field name
	//   defaultValue: Optional default value if field is not found
	//
	// Returns:
	//   string: The input value, or default value, or empty string
	//
	// Example:
	//   name := req.Input("name")
	//   email := req.Input("email", "unknown@example.com")
	Input(key string, defaultValue ...string) string
	
	// InputInt retrieves input data as an integer.
	//
	// Parameters:
	//   key: The input field name
	//   defaultValue: Optional default value
	//
	// Returns:
	//   int: The input value as integer, or default value
	InputInt(key string, defaultValue ...int) int
	
	// InputBool retrieves input data as a boolean.
	//
	// Parameters:
	//   key: The input field name
	//   defaultValue: Optional default value
	//
	// Returns:
	//   bool: The input value as boolean, or default value
	InputBool(key string, defaultValue ...bool) bool
	
	// All retrieves all input data as a map.
	// This includes both query parameters and request body data.
	//
	// Returns:
	//   map[string]interface{}: All input data
	All() map[string]interface{}
	
	// Only retrieves only the specified input fields.
	//
	// Parameters:
	//   keys: The field names to retrieve
	//
	// Returns:
	//   map[string]interface{}: Only the specified input fields
	//
	// Example:
	//   data := req.Only("name", "email", "age")
	Only(keys ...string) map[string]interface{}
	
	// Except retrieves all input data except the specified fields.
	//
	// Parameters:
	//   keys: The field names to exclude
	//
	// Returns:
	//   map[string]interface{}: All input data except specified fields
	//
	// Example:
	//   data := req.Except("password", "confirm_password")
	Except(keys ...string) map[string]interface{}
	
	// Has checks if the request contains the specified input field.
	//
	// Parameters:
	//   key: The field name to check
	//
	// Returns:
	//   bool: True if the field exists, false otherwise
	Has(key string) bool
	
	// Filled checks if the request contains the specified field and it's not empty.
	//
	// Parameters:
	//   key: The field name to check
	//
	// Returns:
	//   bool: True if the field exists and is not empty, false otherwise
	Filled(key string) bool
	
	// JSON Body Parsing
	
	// Json retrieves the JSON body and unmarshals it into the provided struct.
	// The target parameter should be a pointer to the struct.
	//
	// Parameters:
	//   target: Pointer to the struct to unmarshal JSON into
	//
	// Returns:
	//   error: Any error that occurred during JSON parsing
	//
	// Example:
	//   var user User
	//   err := req.Json(&user)
	Json(target interface{}) error
	
	// Body retrieves the raw request body as bytes.
	//
	// Returns:
	//   []byte: The raw request body
	//   error: Any error that occurred while reading the body
	Body() ([]byte, error)
	
	// BodyString retrieves the request body as a string.
	//
	// Returns:
	//   string: The request body as string
	//   error: Any error that occurred while reading the body
	BodyString() (string, error)
	
	// BodyReader returns an io.Reader for the request body.
	// This is useful for streaming large request bodies.
	//
	// Returns:
	//   io.Reader: A reader for the request body
	BodyReader() io.Reader
	
	// File Uploads
	
	// File retrieves an uploaded file by field name.
	// Returns nil if no file was uploaded with the specified name.
	//
	// Parameters:
	//   key: The form field name for the file upload
	//
	// Returns:
	//   *multipart.FileHeader: The uploaded file, or nil if not found
	//   error: Any error that occurred while accessing the file
	//
	// Example:
	//   file, err := req.File("avatar")
	//   if err == nil && file != nil {
	//       // Process the uploaded file
	//   }
	File(key string) (*multipart.FileHeader, error)
	
	// Files retrieves all uploaded files for a field name.
	// This is useful for multiple file uploads with the same field name.
	//
	// Parameters:
	//   key: The form field name for the file uploads
	//
	// Returns:
	//   []*multipart.FileHeader: Array of uploaded files
	//   error: Any error that occurred while accessing the files
	Files(key string) ([]*multipart.FileHeader, error)
	
	// AllFiles retrieves all uploaded files from the request.
	//
	// Returns:
	//   map[string][]*multipart.FileHeader: Map of field names to uploaded files
	//   error: Any error that occurred while accessing the files
	AllFiles() (map[string][]*multipart.FileHeader, error)
	
	// HasFile checks if a file was uploaded with the specified field name.
	//
	// Parameters:
	//   key: The form field name to check
	//
	// Returns:
	//   bool: True if a file exists for the field, false otherwise
	HasFile(key string) bool
	
	// HTTP Headers
	
	// Header retrieves an HTTP header value by name.
	// Header names are case-insensitive.
	//
	// Parameters:
	//   key: The header name
	//   defaultValue: Optional default value if header is not found
	//
	// Returns:
	//   string: The header value, or default value, or empty string
	//
	// Example:
	//   auth := req.Header("Authorization")
	//   userAgent := req.Header("User-Agent")
	Header(key string, defaultValue ...string) string
	
	// Headers retrieves all HTTP headers.
	//
	// Returns:
	//   http.Header: All HTTP headers
	Headers() http.Header
	
	// HasHeader checks if an HTTP header exists.
	//
	// Parameters:
	//   key: The header name to check
	//
	// Returns:
	//   bool: True if the header exists, false otherwise
	HasHeader(key string) bool
	
	// Bearer retrieves the Bearer token from the Authorization header.
	// Returns an empty string if no Bearer token is found.
	//
	// Returns:
	//   string: The Bearer token without the "Bearer " prefix
	//
	// Example:
	//   token := req.Bearer() // For "Authorization: Bearer abc123" returns "abc123"
	Bearer() string
	
	// Request Metadata
	
	// Method retrieves the HTTP method (GET, POST, PUT, etc.).
	//
	// Returns:
	//   string: The HTTP method in uppercase
	Method() string
	
	// URL retrieves the full request URL.
	//
	// Returns:
	//   *url.URL: The parsed request URL
	URL() *url.URL
	
	// Path retrieves the request path without query string.
	//
	// Returns:
	//   string: The request path
	//
	// Example:
	//   // For "/users/123?page=1" returns "/users/123"
	//   path := req.Path()
	Path() string
	
	// FullURL retrieves the complete URL including scheme, host, path, and query.
	//
	// Returns:
	//   string: The complete URL
	FullURL() string
	
	// Scheme retrieves the URL scheme (http, https).
	//
	// Returns:
	//   string: The URL scheme
	Scheme() string
	
	// Host retrieves the request host (hostname:port).
	//
	// Returns:
	//   string: The request host
	Host() string
	
	// Hostname retrieves just the hostname without port.
	//
	// Returns:
	//   string: The hostname
	Hostname() string
	
	// Port retrieves the port number.
	// Returns 80 for HTTP or 443 for HTTPS if no explicit port is specified.
	//
	// Returns:
	//   int: The port number
	Port() int
	
	// IP retrieves the client IP address.
	// This method considers X-Forwarded-For and X-Real-IP headers for proxy scenarios.
	//
	// Returns:
	//   string: The client IP address
	IP() string
	
	// UserAgent retrieves the User-Agent header.
	//
	// Returns:
	//   string: The User-Agent header value
	UserAgent() string
	
	// Referer retrieves the Referer header.
	//
	// Returns:
	//   string: The Referer header value
	Referer() string
	
	// ContentType retrieves the Content-Type header.
	//
	// Returns:
	//   string: The Content-Type header value
	ContentType() string
	
	// ContentLength retrieves the Content-Length as an integer.
	//
	// Returns:
	//   int64: The content length in bytes, or 0 if not specified
	ContentLength() int64
	
	// Cookies Handling
	
	// Cookie retrieves a cookie value by name.
	//
	// Parameters:
	//   name: The cookie name
	//   defaultValue: Optional default value if cookie is not found
	//
	// Returns:
	//   string: The cookie value, or default value, or empty string
	Cookie(name string, defaultValue ...string) string
	
	// Cookies retrieves all cookies.
	//
	// Returns:
	//   []*http.Cookie: Array of all cookies
	Cookies() []*http.Cookie
	
	// HasCookie checks if a cookie exists.
	//
	// Parameters:
	//   name: The cookie name to check
	//
	// Returns:
	//   bool: True if the cookie exists, false otherwise
	HasCookie(name string) bool
	
	// Request Validation Helpers
	
	// IsJson checks if the request has JSON content type.
	//
	// Returns:
	//   bool: True if content type is application/json
	IsJson() bool
	
	// IsXml checks if the request has XML content type.
	//
	// Returns:
	//   bool: True if content type is application/xml or text/xml
	IsXml() bool
	
	// IsForm checks if the request has form data content type.
	//
	// Returns:
	//   bool: True if content type is application/x-www-form-urlencoded
	IsForm() bool
	
	// IsMultipart checks if the request has multipart form data content type.
	//
	// Returns:
	//   bool: True if content type is multipart/form-data
	IsMultipart() bool
	
	// IsSecure checks if the request was made over HTTPS.
	//
	// Returns:
	//   bool: True if the request is HTTPS
	IsSecure() bool
	
	// IsAjax checks if the request was made via AJAX.
	// This checks for X-Requested-With header with value "XMLHttpRequest".
	//
	// Returns:
	//   bool: True if the request is an AJAX request
	IsAjax() bool
	
	// Accepts checks if the request accepts the specified content type.
	//
	// Parameters:
	//   contentType: The content type to check (e.g., "application/json")
	//
	// Returns:
	//   bool: True if the content type is accepted
	Accepts(contentType string) bool
	
	// WantsJson checks if the request wants a JSON response.
	// This checks the Accept header for JSON content types.
	//
	// Returns:
	//   bool: True if JSON response is preferred
	WantsJson() bool
	
	// Context and Advanced Features
	
	// Context retrieves the request context.
	// This can be used for request-scoped values and cancellation.
	//
	// Returns:
	//   context.Context: The request context
	Context() interface{}
	
	// SetContext sets a value in the request context.
	//
	// Parameters:
	//   key: The context key
	//   value: The context value
	SetContext(key string, value interface{})
	
	// GetContext retrieves a value from the request context.
	//
	// Parameters:
	//   key: The context key
	//
	// Returns:
	//   interface{}: The context value, or nil if not found
	GetContext(key string) interface{}
	
	// Fresh checks if the request is fresh based on If-None-Match and If-Modified-Since headers.
	// This is useful for HTTP caching.
	//
	// Returns:
	//   bool: True if the request is fresh (not modified)
	Fresh() bool
	
	// Stale is the opposite of Fresh().
	//
	// Returns:
	//   bool: True if the request is stale (modified)
	Stale() bool
	
	// IfModifiedSince retrieves the If-Modified-Since header as a time.
	//
	// Returns:
	//   time.Time: The If-Modified-Since time, or zero time if not present
	IfModifiedSince() time.Time
	
	// IfNoneMatch retrieves the If-None-Match header.
	//
	// Returns:
	//   string: The If-None-Match header value
	IfNoneMatch() string
}
