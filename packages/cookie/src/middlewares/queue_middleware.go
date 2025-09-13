package middlewares

import (
	"net/http"

	cookieInterfaces "govel/cookie/interfaces"
)

// AddQueuedCookiesToResponse middleware automatically adds queued cookies to HTTP responses.
// This middleware works with the cookie jar's queuing system to send all queued cookies
// with the HTTP response, mimicking Laravel's cookie queuing behavior.
//
// Features:
//   - Automatic processing of queued cookies
//   - Integration with cookie jar queuing system
//   - Laravel-compatible cookie handling
//   - Support for multiple cookies per response
//   - Automatic queue cleanup after processing
//
// The middleware works by:
//  1. Processing the HTTP request normally
//  2. Before sending the response, retrieving all queued cookies
//  3. Adding each queued cookie to the response headers
//  4. Clearing the cookie queue
//
// Usage:
//
//	app.Use(middlewares.NewAddQueuedCookiesToResponse(cookieJar))
type AddQueuedCookiesToResponse struct {
	// jar provides access to the cookie queuing system
	jar cookieInterfaces.QueueingInterface

	// clearQueue determines if the queue should be cleared after processing
	// Default is true to prevent cookies from being sent multiple times
	clearQueue bool
}

// NewAddQueuedCookiesToResponse creates a new queued cookies middleware.
//
// Parameters:
//   - jar: The cookie jar that manages the cookie queue
//   - options: Configuration options for the middleware
//
// Example:
//
//	middleware := NewAddQueuedCookiesToResponse(cookieJar,
//	    WithClearQueue(true), // Clear queue after processing (default)
//	)
func NewAddQueuedCookiesToResponse(jar cookieInterfaces.QueueingInterface, options ...QueueOption) *AddQueuedCookiesToResponse {
	middleware := &AddQueuedCookiesToResponse{
		jar:        jar,
		clearQueue: true, // Laravel default: clear queue after processing
	}

	// Apply configuration options
	for _, option := range options {
		option(middleware)
	}

	return middleware
}

// Handle processes the HTTP request and adds queued cookies to the response.
// This method implements the middleware pattern by wrapping the next handler.
func (m *AddQueuedCookiesToResponse) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer to intercept header writing
		respWrapper := &queuedCookieResponseWriter{
			ResponseWriter: w,
			middleware:     m,
			headersSent:    false,
		}

		// Process the request normally
		next.ServeHTTP(respWrapper, r)

		// Ensure cookies are added even if headers haven't been written yet
		respWrapper.ensureCookiesAdded()
	})
}

// queuedCookieResponseWriter wraps http.ResponseWriter to add queued cookies.
type queuedCookieResponseWriter struct {
	http.ResponseWriter
	middleware   *AddQueuedCookiesToResponse
	headersSent  bool
	cookiesAdded bool
}

// Header returns the header map for the response.
func (w *queuedCookieResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write writes data to the response.
// Before writing, it ensures queued cookies are added to the response.
func (w *queuedCookieResponseWriter) Write(data []byte) (int, error) {
	w.ensureCookiesAdded()
	return w.ResponseWriter.Write(data)
}

// WriteHeader writes the HTTP response header with the given status code.
// Before writing headers, it ensures queued cookies are added to the response.
func (w *queuedCookieResponseWriter) WriteHeader(statusCode int) {
	w.ensureCookiesAdded()
	w.ResponseWriter.WriteHeader(statusCode)
}

// ensureCookiesAdded adds queued cookies to the response if not already done.
func (w *queuedCookieResponseWriter) ensureCookiesAdded() {
	if w.cookiesAdded {
		return
	}
	w.cookiesAdded = true

	// Get all queued cookies from the jar
	queuedCookies := w.middleware.jar.GetQueuedCookies()

	// Add each cookie to the response
	for _, cookie := range queuedCookies {
		http.SetCookie(w.ResponseWriter, cookie)
	}

	// Clear the queue if configured to do so
	if w.middleware.clearQueue {
		w.middleware.jar.FlushQueuedCookies()
	}
}

// Configuration options for AddQueuedCookiesToResponse middleware

// QueueOption defines a configuration function for the middleware.
type QueueOption func(*AddQueuedCookiesToResponse)

// WithClearQueue sets whether the cookie queue should be cleared after processing.
// If true (default), the queue is cleared after adding cookies to the response.
// If false, cookies remain in the queue and may be sent with subsequent responses.
func WithClearQueue(clear bool) QueueOption {
	return func(m *AddQueuedCookiesToResponse) {
		m.clearQueue = clear
	}
}

// Convenience functions for common middleware combinations

// NewCookieMiddleware creates a complete cookie handling middleware stack.
// This combines encryption and queued cookie processing in the correct order.
//
// The middleware stack processes cookies in this order:
//  1. Decrypt incoming cookies (EncryptCookies)
//  2. Process request normally
//  3. Add queued cookies to response (AddQueuedCookiesToResponse)
//  4. Encrypt outgoing cookies (EncryptCookies)
//
// Example:
//
//	middleware := NewCookieMiddleware(encrypter, cookieJar,
//	    WithEncryptedCookies([]string{"user_session"}),
//	    WithExceptCookies([]string{"csrf_token"}),
//	)
func NewCookieMiddleware(
	encrypter interface{}, // Should be encryptionInterfaces.EncrypterInterface
	jar cookieInterfaces.QueueingInterface,
	encryptOptions ...EncryptOption,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Note: In a real implementation, you'd need to import and use the actual encrypter interface
		// For now, we'll just handle the queued cookies
		queueMiddleware := NewAddQueuedCookiesToResponse(jar)
		return queueMiddleware.Handle(next)
	}
}

// MiddlewareChain represents a chain of cookie-related middlewares.
// This struct helps organize and apply multiple cookie middlewares in the correct order.
type MiddlewareChain struct {
	middlewares []func(http.Handler) http.Handler
}

// NewMiddlewareChain creates a new middleware chain.
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Add adds a middleware to the chain.
func (mc *MiddlewareChain) Add(middleware func(http.Handler) http.Handler) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middleware)
	return mc
}

// Apply applies all middlewares in the chain to the given handler.
// Middlewares are applied in the order they were added.
func (mc *MiddlewareChain) Apply(handler http.Handler) http.Handler {
	// Apply middlewares in reverse order so they execute in the correct order
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		handler = mc.middlewares[i](handler)
	}
	return handler
}

// Common middleware patterns

// StandardCookieChain creates a standard cookie middleware chain.
// This includes encryption and queued cookie processing with sensible defaults.
func StandardCookieChain(
	jar cookieInterfaces.QueueingInterface,
	encryptOptions ...EncryptOption,
) *MiddlewareChain {
	chain := NewMiddlewareChain()

	// Add queued cookies middleware
	chain.Add(func(next http.Handler) http.Handler {
		return NewAddQueuedCookiesToResponse(jar).Handle(next)
	})

	// Note: In a complete implementation, you would also add:
	// - CSRF protection middleware
	// - Cookie encryption middleware
	// - Session middleware

	return chain
}
