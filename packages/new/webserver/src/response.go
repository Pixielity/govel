// Package webserver - Unified Response implementation
// This file implements a concrete Response that satisfies interfaces.ResponseInterface.
// Adapters are responsible for writing this response to their native response writers.
package webserver

import (
	"encoding/json"
	"govel/packages/new/webserver/src/interfaces"
	"io"
	"net/http"
)

// Response is a framework-agnostic HTTP response implementing interfaces.ResponseInterface.
// It supports fluent method chaining and stores headers, cookies, body, and status code.
type Response struct {
	status      int
	headers     map[string]string
	body        []byte
	contentType string
	cookies     []*http.Cookie
	noContent   bool
}

// NewResponse creates an empty response with defaults (status 200, no headers).
func NewResponse() *Response {
	return &Response{
		status:  200,
		headers: map[string]string{},
		cookies: []*http.Cookie{},
	}
}

// Status sets the HTTP status code.
func (r *Response) Status(code int) interfaces.ResponseInterface { r.status = code; return r }

// Header sets a single header value.
func (r *Response) Header(key, value string) interfaces.ResponseInterface {
	r.headers[key] = value
	return r
}

// Headers sets multiple headers.
func (r *Response) Headers(h map[string]string) interfaces.ResponseInterface {
	for k, v := range h {
		r.headers[k] = v
	}
	return r
}

// RemoveHeader removes a header if present.
func (r *Response) RemoveHeader(key string) interfaces.ResponseInterface {
	delete(r.headers, key)
	return r
}

// Json sets a JSON body and content type.
func (r *Response) Json(payload interface{}) interfaces.ResponseInterface {
	b, _ := json.Marshal(payload)
	r.body = b
	r.contentType = "application/json"
	return r
}

// Text sets a plain text body and content type.
func (r *Response) Text(body string) interfaces.ResponseInterface {
	r.body = []byte(body)
	r.contentType = "text/plain; charset=utf-8"
	return r
}

// HTML sets an HTML body and content type.
func (r *Response) HTML(html string) interfaces.ResponseInterface {
	r.body = []byte(html)
	r.contentType = "text/html; charset=utf-8"
	return r
}

// Send sets a raw body and explicit content type.
func (r *Response) Send(body []byte, contentType string) interfaces.ResponseInterface {
	r.body = body
	r.contentType = contentType
	return r
}

// Stream is a placeholder (adapters may implement streaming separately).
func (r *Response) Stream(reader io.Reader, contentType string) interfaces.ResponseInterface {
	// Not storing reader here; adapters can provide streaming-specific responses if needed.
	r.contentType = contentType
	return r
}

// File serves a file path (adapters should implement actual file serving).
func (r *Response) File(_ string) interfaces.ResponseInterface { return r }

// Download sets disposition header for an attachment; adapters should send the file.
func (r *Response) Download(_ string, _ ...string) interfaces.ResponseInterface { return r }

// Redirect sets the Location header and status code (default 302).
func (r *Response) Redirect(url string, status ...int) interfaces.ResponseInterface {
	code := 302
	if len(status) > 0 {
		code = status[0]
	}
	r.Status(code).Header("Location", url)
	return r
}

// Cookie adds a Set-Cookie entry.
func (r *Response) Cookie(cookie *http.Cookie) interfaces.ResponseInterface {
	if cookie != nil {
		r.cookies = append(r.cookies, cookie)
	}
	return r
}

// ClearCookie adds a clearing cookie with MaxAge < 0.
func (r *Response) ClearCookie(name string) interfaces.ResponseInterface {
	r.cookies = append(r.cookies, &http.Cookie{Name: name, MaxAge: -1, Path: "/"})
	return r
}

// NoContent marks response as having no body with a status (default 204).
func (r *Response) NoContent(status ...int) interfaces.ResponseInterface {
	code := 204
	if len(status) > 0 {
		code = status[0]
	}
	r.Status(code)
	r.noContent = true
	r.body = nil
	return r
}

// RawWriter is a placeholder for adapters to hook into low-level writers.
func (r *Response) RawWriter(_ func(w http.ResponseWriter)) interfaces.ResponseInterface { return r }

// Accessors for adapters
func (r *Response) StatusCode() int               { return r.status }
func (r *Response) HeadersMap() map[string]string { return r.headers }
func (r *Response) Body() []byte                  { return r.body }
func (r *Response) ContentType() string           { return r.contentType }
func (r *Response) Cookies() []*http.Cookie       { return r.cookies }
func (r *Response) IsNoContent() bool             { return r.noContent }

// Compile-time interface compliance check
var _ interfaces.ResponseInterface = (*Response)(nil)
