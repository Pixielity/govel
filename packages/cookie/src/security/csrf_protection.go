package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	cookie "govel/cookie"
	cookieInterfaces "govel/cookie/interfaces"
)

// CSRFProtection provides Laravel-compatible CSRF protection for cookies and forms.
// This implementation follows Laravel's CSRF token generation and validation patterns,
// providing protection against Cross-Site Request Forgery attacks.
//
// Features:
//   - Automatic CSRF token generation
//   - Laravel-compatible token format
//   - Multiple token validation methods (header, form, query)
//   - Configurable token lifetime
//   - Double-submit cookie pattern support
//   - Integration with cookie jar for token storage
//
// The CSRF protection works by:
//  1. Generating a unique token for each session
//  2. Storing the token in a secure cookie
//  3. Requiring the token in forms/AJAX requests
//  4. Validating the token on submission
type CSRFProtection struct {
	// jar manages cookie operations
	jar cookieInterfaces.JarInterface

	// tokenName is the name of the CSRF token field/cookie
	tokenName string

	// cookieName is the name of the cookie storing the CSRF token
	cookieName string

	// headerName is the name of the HTTP header containing the CSRF token
	headerName string

	// tokenLength is the length of generated tokens in bytes
	tokenLength int

	// tokenLifetime is how long tokens remain valid
	tokenLifetime time.Duration

	// except contains routes/paths that should be excluded from CSRF protection
	except []string

	// methods contains HTTP methods that require CSRF protection
	methods []string
}

// NewCSRFProtection creates a new CSRF protection instance with Laravel defaults.
//
// Default configuration:
//   - Token name: "_token" (Laravel standard)
//   - Cookie name: "XSRF-TOKEN" (Laravel standard for AJAX)
//   - Header name: "X-CSRF-TOKEN" (Laravel standard)
//   - Token length: 32 bytes (256 bits)
//   - Token lifetime: 2 hours (Laravel session lifetime)
//   - Protected methods: POST, PUT, PATCH, DELETE
//
// Parameters:
//   - jar: Cookie jar for token storage
//   - options: Configuration options
func NewCSRFProtection(jar cookieInterfaces.JarInterface, options ...CSRFOption) *CSRFProtection {
	csrf := &CSRFProtection{
		jar:           jar,
		tokenName:     "_token",       // Laravel default
		cookieName:    "XSRF-TOKEN",   // Laravel AJAX token
		headerName:    "X-CSRF-TOKEN", // Laravel header
		tokenLength:   32,             // 256 bits
		tokenLifetime: 2 * time.Hour,  // Laravel session lifetime
		except:        []string{},
		methods:       []string{"POST", "PUT", "PATCH", "DELETE"},
	}

	// Apply configuration options
	for _, option := range options {
		option(csrf)
	}

	return csrf
}

// GenerateToken creates a new CSRF token.
// The token is cryptographically secure and follows Laravel's format.
//
// Returns:
//   - string: Base64-encoded token
//   - error: Any error during token generation
func (c *CSRFProtection) GenerateToken() (string, error) {
	// Generate random bytes
	tokenBytes := make([]byte, c.tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	// Encode to base64 for safe transmission
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

// GetTokenFromRequest extracts the CSRF token from the HTTP request.
// This method checks multiple sources in order of priority:
//  1. X-CSRF-TOKEN header
//  2. XSRF-TOKEN header (for AJAX compatibility)
//  3. _token form field
//  4. token query parameter
//
// Returns:
//   - string: The extracted token (empty if not found)
func (c *CSRFProtection) GetTokenFromRequest(r *http.Request) string {
	// Check headers first (most secure)
	if token := r.Header.Get(c.headerName); token != "" {
		return token
	}

	// Check XSRF-TOKEN header for AJAX requests
	if token := r.Header.Get("X-XSRF-TOKEN"); token != "" {
		// Decode from URL encoding (JavaScript's encodeURIComponent)
		if decoded, err := base64.URLEncoding.DecodeString(token); err == nil {
			return string(decoded)
		}
		return token
	}

	// Check form fields
	if err := r.ParseForm(); err == nil {
		if token := r.Form.Get(c.tokenName); token != "" {
			return token
		}
	}

	// Check query parameters (least secure)
	if token := r.URL.Query().Get(c.tokenName); token != "" {
		return token
	}

	return ""
}

// GetTokenFromCookie extracts the CSRF token from cookies.
// This is used for the double-submit cookie pattern validation.
//
// Returns:
//   - string: The token from cookies (empty if not found)
func (c *CSRFProtection) GetTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(c.cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ValidateToken validates a CSRF token against the stored token.
// This method implements constant-time comparison to prevent timing attacks.
//
// Parameters:
//   - requestToken: Token from the request
//   - storedToken: Token stored in session/cookie
//
// Returns:
//   - bool: True if tokens match, false otherwise
func (c *CSRFProtection) ValidateToken(requestToken, storedToken string) bool {
	if requestToken == "" || storedToken == "" {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(requestToken), []byte(storedToken)) == 1
}

// SetTokenCookie sets the CSRF token in a cookie.
// This is used for the double-submit cookie pattern and AJAX requests.
//
// Parameters:
//   - w: HTTP response writer
//   - token: CSRF token to set
func (c *CSRFProtection) SetTokenCookie(w http.ResponseWriter, token string) {
	// Create cookie with Laravel-compatible settings
	cookie := c.jar.Make(c.cookieName, token,
		cookie.WithPath("/"),                      // Available site-wide
		cookie.WithSecure(false),                  // Allow HTTP in development
		cookie.WithHttpOnly(false),                // Allow JavaScript access
		cookie.WithSameSite(http.SameSiteLaxMode), // Laravel default
		cookie.WithMaxAge(int(c.tokenLifetime.Seconds())),
	)

	// Set the cookie immediately
	http.SetCookie(w, cookie)
}

// ShouldValidate determines if a request should be validated for CSRF.
// This method checks the HTTP method and exclusion list.
//
// Parameters:
//   - r: HTTP request to check
//
// Returns:
//   - bool: True if CSRF validation is required
func (c *CSRFProtection) ShouldValidate(r *http.Request) bool {
	// Check if method requires CSRF protection
	methodRequiresCSRF := false
	for _, method := range c.methods {
		if r.Method == method {
			methodRequiresCSRF = true
			break
		}
	}

	if !methodRequiresCSRF {
		return false
	}

	// Check if path is excluded
	requestPath := r.URL.Path
	for _, exceptPath := range c.except {
		if c.matchesPath(requestPath, exceptPath) {
			return false
		}
	}

	return true
}

// matchesPath checks if a request path matches an exception pattern.
// Supports wildcards and exact matches.
func (c *CSRFProtection) matchesPath(requestPath, pattern string) bool {
	// Exact match
	if requestPath == pattern {
		return true
	}

	// Wildcard match (simple implementation)
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(requestPath, prefix)
	}

	return false
}

// Middleware creates an HTTP middleware for CSRF protection.
// This middleware automatically validates CSRF tokens and generates new ones as needed.
//
// The middleware:
//  1. Checks if the request requires CSRF validation
//  2. Extracts and validates the CSRF token
//  3. Returns 403 Forbidden for invalid tokens
//  4. Generates and sets new tokens for valid requests
//
// Usage:
//
//	app.Use(csrfProtection.Middleware())
func (c *CSRFProtection) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this request needs CSRF validation
			if !c.ShouldValidate(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Get tokens from request and cookie
			requestToken := c.GetTokenFromRequest(r)
			storeToken := c.GetTokenFromCookie(r)

			// Validate the token
			if !c.ValidateToken(requestToken, storeToken) {
				// CSRF validation failed
				http.Error(w, "CSRF token mismatch", http.StatusForbidden)
				return
			}

			// Generate new token for the response
			newToken, err := c.GenerateToken()
			if err != nil {
				// Token generation failed - continue but log the error
				// In production, you might want to return an error
				next.ServeHTTP(w, r)
				return
			}

			// Set the new token in a cookie
			c.SetTokenCookie(w, newToken)

			// Continue processing the request
			next.ServeHTTP(w, r)
		})
	}
}

// Configuration options for CSRF protection

// CSRFOption defines a configuration function for CSRF protection.
type CSRFOption func(*CSRFProtection)

// WithTokenName sets the name of the CSRF token field.
func WithTokenName(name string) CSRFOption {
	return func(c *CSRFProtection) {
		c.tokenName = name
	}
}

// WithCookieName sets the name of the CSRF token cookie.
func WithCookieName(name string) CSRFOption {
	return func(c *CSRFProtection) {
		c.cookieName = name
	}
}

// WithHeaderName sets the name of the CSRF token header.
func WithHeaderName(name string) CSRFOption {
	return func(c *CSRFProtection) {
		c.headerName = name
	}
}

// WithTokenLength sets the length of generated tokens.
func WithTokenLength(length int) CSRFOption {
	return func(c *CSRFProtection) {
		c.tokenLength = length
	}
}

// WithTokenLifetime sets how long tokens remain valid.
func WithTokenLifetime(lifetime time.Duration) CSRFOption {
	return func(c *CSRFProtection) {
		c.tokenLifetime = lifetime
	}
}

// WithExcept sets paths that should be excluded from CSRF protection.
func WithExcept(paths []string) CSRFOption {
	return func(c *CSRFProtection) {
		c.except = paths
	}
}

// WithMethods sets HTTP methods that require CSRF protection.
func WithMethods(methods []string) CSRFOption {
	return func(c *CSRFProtection) {
		c.methods = methods
	}
}
