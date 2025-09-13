package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	encryptionInterfaces "govel/types/interfaces/encryption"
)

// EncryptCookies middleware provides Laravel-compatible cookie encryption and decryption.
// This middleware automatically encrypts cookies when they are sent to the client
// and decrypts them when they are received from the client.
//
// Features:
//   - Automatic cookie encryption/decryption
//   - Selective encryption (only specified cookies)
//   - Laravel-compatible payload format
//   - Base64 encoding for safe HTTP transmission
//   - Error handling for invalid or tampered cookies
//   - Support for cookie whitelisting (unencrypted cookies)
//
// The middleware works by:
//  1. Decrypting incoming cookies that should be encrypted
//  2. Allowing normal request processing
//  3. Encrypting outgoing cookies before sending to client
//
// Usage:
//
//	app.Use(middlewares.NewEncryptCookies(encrypter, []string{"user_session", "preferences"}))
type EncryptCookies struct {
	// encrypter handles the actual encryption/decryption operations
	encrypter encryptionInterfaces.EncrypterInterface

	// encryptedCookies lists cookie names that should be encrypted
	// If empty, all cookies are encrypted except those in except list
	encryptedCookies []string

	// except lists cookie names that should never be encrypted
	// This is useful for cookies that need to be readable by JavaScript
	// or external services (like CSRF tokens, language preferences, etc.)
	except []string

	// encryptAll determines if all cookies should be encrypted by default
	// If true, all cookies are encrypted except those in the except list
	// If false, only cookies in encryptedCookies list are encrypted
	encryptAll bool
}

// NewEncryptCookies creates a new cookie encryption middleware.
//
// Parameters:
//   - encrypter: The encryption service to use for cookie encryption
//   - options: Configuration options for the middleware
//
// Example:
//
//	// Encrypt specific cookies
//	middleware := NewEncryptCookies(encrypter,
//	    WithEncryptedCookies([]string{"user_session", "preferences"}),
//	    WithExceptCookies([]string{"csrf_token", "language"}),
//	)
func NewEncryptCookies(encrypter encryptionInterfaces.EncrypterInterface, options ...EncryptOption) *EncryptCookies {
	middleware := &EncryptCookies{
		encrypter:        encrypter,
		encryptedCookies: []string{},
		except:           []string{},
		encryptAll:       true, // Laravel default: encrypt all cookies
	}

	// Apply configuration options
	for _, option := range options {
		option(middleware)
	}

	return middleware
}

// Handle processes the HTTP request and response for cookie encryption.
// This method implements the middleware pattern by wrapping the next handler.
func (m *EncryptCookies) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Decrypt incoming cookies
		m.decryptRequestCookies(r)

		// Step 2: Wrap response writer to encrypt outgoing cookies
		respWrapper := &encryptingResponseWriter{
			ResponseWriter: w,
			middleware:     m,
			encrypted:      make(map[string]bool),
		}

		// Step 3: Continue with request processing
		next.ServeHTTP(respWrapper, r)
	})
}

// decryptRequestCookies decrypts cookies in the incoming request.
// Only cookies that should be encrypted are decrypted.
func (m *EncryptCookies) decryptRequestCookies(r *http.Request) {
	for _, cookie := range r.Cookies() {
		if m.shouldEncrypt(cookie.Name) {
			// Attempt to decrypt the cookie value
			decrypted, err := m.decryptCookieValue(cookie.Value)
			if err != nil {
				// If decryption fails, remove the cookie (it's invalid)
				// Laravel does this to handle tampered or corrupted cookies
				cookie.Value = ""
				cookie.MaxAge = -1
				continue
			}

			// Replace the encrypted value with decrypted value
			cookie.Value = decrypted
		}
	}
}

// decryptCookieValue decrypts a single cookie value.
// The value is expected to be base64 encoded encrypted data.
func (m *EncryptCookies) decryptCookieValue(encryptedValue string) (string, error) {
	if encryptedValue == "" {
		return "", nil
	}

	// Decode from base64
	decodedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 cookie value: %w", err)
	}

	// Decrypt the payload
	decrypted, err := m.encrypter.DecryptString(string(decodedBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decrypt cookie value: %w", err)
	}

	return decrypted, nil
}

// encryptCookieValue encrypts a single cookie value.
// The result is base64 encoded for safe HTTP transmission.
func (m *EncryptCookies) encryptCookieValue(plainValue string) (string, error) {
	if plainValue == "" {
		return "", nil
	}

	// Encrypt the value
	encrypted, err := m.encrypter.EncryptString(plainValue)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt cookie value: %w", err)
	}

	// Encode to base64 for safe HTTP transmission
	encoded := base64.StdEncoding.EncodeToString([]byte(encrypted))
	return encoded, nil
}

// shouldEncrypt determines if a cookie should be encrypted.
// This follows Laravel's logic for determining which cookies to encrypt.
func (m *EncryptCookies) shouldEncrypt(cookieName string) bool {
	// First check if it's explicitly excluded
	for _, exceptName := range m.except {
		if exceptName == cookieName {
			return false
		}
	}

	// If encryptAll is true, encrypt everything not in except list
	if m.encryptAll {
		return true
	}

	// Otherwise, only encrypt cookies explicitly listed
	for _, encryptName := range m.encryptedCookies {
		if encryptName == cookieName {
			return true
		}
	}

	return false
}

// encryptingResponseWriter wraps http.ResponseWriter to encrypt cookies.
type encryptingResponseWriter struct {
	http.ResponseWriter
	middleware *EncryptCookies
	encrypted  map[string]bool // Track which cookies have been encrypted
}

// Header returns the header map for the response.
func (w *encryptingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write writes data to the response.
func (w *encryptingResponseWriter) Write(data []byte) (int, error) {
	// Before writing, process any cookies that were set
	w.encryptResponseCookies()
	return w.ResponseWriter.Write(data)
}

// WriteHeader writes the HTTP response header with the given status code.
func (w *encryptingResponseWriter) WriteHeader(statusCode int) {
	// Before writing headers, process any cookies that were set
	w.encryptResponseCookies()
	w.ResponseWriter.WriteHeader(statusCode)
}

// encryptResponseCookies encrypts cookies in the response headers.
func (w *encryptingResponseWriter) encryptResponseCookies() {
	headers := w.Header()
	cookieHeaders := headers["Set-Cookie"]

	if len(cookieHeaders) == 0 {
		return
	}

	// Process each Set-Cookie header
	updatedHeaders := make([]string, 0, len(cookieHeaders))
	for _, cookieHeader := range cookieHeaders {
		updatedHeader := w.encryptCookieHeader(cookieHeader)
		updatedHeaders = append(updatedHeaders, updatedHeader)
	}

	// Replace the headers with encrypted versions
	headers["Set-Cookie"] = updatedHeaders
}

// encryptCookieHeader encrypts the value in a Set-Cookie header.
func (w *encryptingResponseWriter) encryptCookieHeader(cookieHeader string) string {
	// Parse the cookie header to extract name and value
	parts := strings.Split(cookieHeader, ";")
	if len(parts) == 0 {
		return cookieHeader
	}

	// Extract name=value part
	nameValue := strings.TrimSpace(parts[0])
	equalIndex := strings.Index(nameValue, "=")
	if equalIndex == -1 {
		return cookieHeader
	}

	cookieName := nameValue[:equalIndex]
	cookieValue := nameValue[equalIndex+1:]

	// Check if this cookie should be encrypted and hasn't been already
	if !w.middleware.shouldEncrypt(cookieName) {
		return cookieHeader
	}

	// Avoid double-encryption
	key := cookieName + "=" + cookieValue
	if w.encrypted[key] {
		return cookieHeader
	}
	w.encrypted[key] = true

	// Encrypt the cookie value
	encryptedValue, err := w.middleware.encryptCookieValue(cookieValue)
	if err != nil {
		// If encryption fails, return original header
		// In production, you might want to log this error
		return cookieHeader
	}

	// Reconstruct the cookie header with encrypted value
	newNameValue := cookieName + "=" + encryptedValue
	if len(parts) > 1 {
		return newNameValue + ";" + strings.Join(parts[1:], ";")
	}
	return newNameValue
}

// Configuration options for EncryptCookies middleware

// EncryptOption defines a configuration function for the middleware.
type EncryptOption func(*EncryptCookies)

// WithEncryptedCookies specifies which cookies should be encrypted.
// If provided, only these cookies will be encrypted (plus any not in except list if encryptAll is true).
func WithEncryptedCookies(cookies []string) EncryptOption {
	return func(m *EncryptCookies) {
		m.encryptedCookies = cookies
	}
}

// WithExceptCookies specifies cookies that should never be encrypted.
// These cookies will always be sent in plain text.
func WithExceptCookies(cookies []string) EncryptOption {
	return func(m *EncryptCookies) {
		m.except = cookies
	}
}

// WithEncryptAll sets whether all cookies should be encrypted by default.
// If true (default), all cookies are encrypted except those in the except list.
// If false, only cookies in the encryptedCookies list are encrypted.
func WithEncryptAll(encryptAll bool) EncryptOption {
	return func(m *EncryptCookies) {
		m.encryptAll = encryptAll
	}
}

// Common cookie names that are typically excluded from encryption
// These are provided as a convenience for common use cases
var (
	// DefaultExceptCookies contains cookie names commonly excluded from encryption
	DefaultExceptCookies = []string{
		"csrf_token",     // CSRF tokens need to be readable
		"XSRF-TOKEN",     // Laravel's CSRF token for AJAX
		"language",       // Language preferences
		"locale",         // Locale preferences
		"timezone",       // Timezone preferences
		"theme",          // UI theme preferences
		"cookie_consent", // Cookie consent status
		"debug_mode",     // Development debugging flags
	}

	// SessionCookies contains cookie names typically used for sessions
	SessionCookies = []string{
		"laravel_session",
		"session_id",
		"user_session",
		"remember_token",
	}

	// UserDataCookies contains cookie names typically containing user data
	UserDataCookies = []string{
		"user_preferences",
		"shopping_cart",
		"recent_items",
		"favorites",
		"user_settings",
	}
)
