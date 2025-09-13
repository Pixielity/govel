package cookie

import (
	"net/http"
	"sync"
	"time"

	"govel/packages/cookie/src/interfaces"
)

// CookieJar implements Laravel-style cookie management with queuing support.
// This struct provides a complete cookie handling system that mirrors Laravel's
// Cookie facade and CookieJar functionality.
//
// Features:
//   - Laravel-compatible cookie creation and management
//   - Cookie queuing for batch processing
//   - Configurable default settings (path, domain, secure, SameSite)
//   - Thread-safe operations with internal synchronization
//   - Support for session cookies, persistent cookies, and "forever" cookies
//   - Cookie forgetting (deletion) functionality
//
// The CookieJar maintains default configuration that can be set once
// and applied to all new cookies unless explicitly overridden.
type CookieJar struct {
	// mu provides thread-safety for all cookie operations
	mu sync.RWMutex

	// queuedCookies stores cookies that are queued for later processing
	// Key format: "name@path" to allow multiple cookies with same name but different paths
	queuedCookies map[string]*QueuedCookie

	// Default configuration settings applied to new cookies
	defaultPath     string
	defaultDomain   string
	defaultSecure   bool
	defaultSameSite http.SameSite
}

// NewCookieJar creates a new cookie jar instance with default configuration.
// The jar is initialized with sensible defaults that match Laravel's behavior:
//   - Path: "/" (root path)
//   - Domain: "" (current domain)
//   - Secure: false (allow HTTP)
//   - SameSite: http.SameSiteLaxMode (Laravel's default)
//
// Returns a fully initialized CookieJar ready for use.
func NewCookieJar() *CookieJar {
	return &CookieJar{
		queuedCookies:   make(map[string]*QueuedCookie),
		defaultPath:     "/",
		defaultDomain:   "",
		defaultSecure:   false,
		defaultSameSite: http.SameSiteLaxMode,
	}
}

// Make creates a new cookie with the specified name and value.
// This method applies default configuration settings and then
// processes any provided options to customize the cookie.
//
// The created cookie uses Laravel-compatible defaults:
//   - Path: Uses jar's default path (typically "/")
//   - Domain: Uses jar's default domain
//   - Secure: Uses jar's default secure setting
//   - HttpOnly: true (secure by default)
//   - SameSite: Uses jar's default SameSite policy
//   - Expires: Not set (session cookie by default)
//
// Examples:
//
//	// Basic session cookie
//	cookie := jar.Make("user_id", "123")
//
//	// Persistent cookie with custom domain
//	cookie := jar.Make("theme", "dark",
//	    interfaces.WithExpiry(time.Now().Add(30*24*time.Hour)),
//	    interfaces.WithDomain(".example.com"),
//	)
func (jar *CookieJar) Make(name, value string, options ...types.CookieOption) *http.Cookie {
	jar.mu.RLock()
	defaultPath := jar.defaultPath
	defaultDomain := jar.defaultDomain
	defaultSecure := jar.defaultSecure
	defaultSameSite := jar.defaultSameSite
	jar.mu.RUnlock()

	// Create cookie with Laravel-compatible defaults
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     defaultPath,
		Domain:   defaultDomain,
		Secure:   defaultSecure,
		HttpOnly: true, // Secure by default
		SameSite: defaultSameSite,
		// Expires is not set - creates a session cookie
	}

	// Apply any provided options
	for _, option := range options {
		option(cookie)
	}

	return cookie
}

// Forever creates a cookie that lasts "forever" (actually 5 years).
// This matches Laravel's cookie()->forever() method behavior.
// The cookie will expire 5 years from the current time.
//
// This is useful for:
//   - User preferences that should persist long-term
//   - "Remember me" functionality
//   - Application settings that rarely change
//
// Examples:
//
//	// Remember user login
//	cookie := jar.Forever("remember_token", token)
//
//	// Persistent user preferences
//	cookie := jar.Forever("language", "en",
//	    interfaces.WithDomain(".example.com"),
//	)
func (jar *CookieJar) Forever(name, value string, options ...types.CookieOption) *http.Cookie {
	// Create base cookie
	cookie := jar.Make(name, value, options...)

	// Set expiration to 5 years from now (Laravel's "forever" duration)
	cookie.Expires = time.Now().Add(5 * 365 * 24 * time.Hour)

	return cookie
}

// Forget creates a cookie that immediately expires, effectively deleting it.
// This is Laravel's standard way of removing cookies from the client.
//
// The method creates a cookie with the same name and path but sets:
//   - Value: "" (empty string)
//   - Expires: January 1, 1970 (Unix epoch)
//   - MaxAge: -1 (immediate deletion)
//
// This ensures the cookie is deleted regardless of the client's implementation.
//
// Examples:
//
//	// Delete a session cookie
//	cookie := jar.Forget("user_session")
//
//	// Delete a cookie with specific path
//	cookie := jar.Forget("admin_token",
//	    interfaces.WithPath("/admin"),
//	)
func (jar *CookieJar) Forget(name string, options ...types.CookieOption) *http.Cookie {
	// Create a cookie with empty value and past expiration
	cookie := jar.Make(name, "", options...)

	// Set expiration to Unix epoch (January 1, 1970)
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1 // Immediate deletion

	return cookie
}

// Queue adds a cookie to the queue for later processing.
// Queued cookies are typically sent with the HTTP response by middleware.
// If a cookie with the same name and path is already queued, it will be replaced.
//
// The cookie is stored with additional metadata including queue timestamp
// and priority for advanced processing scenarios.
//
// Examples:
//
//	// Queue a simple cookie
//	cookie := jar.Make("flash_message", "Welcome!")
//	jar.Queue(cookie)
//
//	// Queue multiple cookies
//	jar.Queue(jar.Make("user_id", "123"))
//	jar.Queue(jar.Make("theme", "dark"))
//	jar.Queue(jar.Make("language", "en"))
func (jar *CookieJar) Queue(cookie *http.Cookie) {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	// Create queued cookie with metadata
	queuedCookie := &QueuedCookie{
		Cookie:   cookie,
		QueuedAt: time.Now().Unix(),
		Priority: 0,
		Metadata: make(map[string]interface{}),
	}

	// Store using name@path key for uniqueness
	key := queuedCookie.GetKey()
	jar.queuedCookies[key] = queuedCookie
}

// Unqueue removes a cookie from the queue.
// If path is not provided, uses the default path ("/").
// This prevents a previously queued cookie from being sent with the response.
//
// Examples:
//
//	// Remove a cookie from default path
//	jar.Unqueue("temporary_token")
//
//	// Remove a cookie from specific path
//	jar.Unqueue("admin_session", "/admin")
func (jar *CookieJar) Unqueue(name string, path ...string) {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	// Determine the path to use
	cookiePath := jar.defaultPath
	if len(path) > 0 && path[0] != "" {
		cookiePath = path[0]
	}

	// Create the key and remove from queue
	key := name + "@" + cookiePath
	delete(jar.queuedCookies, key)
}

// HasQueued checks if a cookie with the given name and path is queued.
// If path is not provided, uses the default path ("/").
//
// This is useful for conditional logic and avoiding duplicate queuing.
//
// Examples:
//
//	// Check if cookie is queued
//	if jar.HasQueued("user_preferences") {
//	    // Cookie is already queued
//	}
//
//	// Check specific path
//	if jar.HasQueued("admin_token", "/admin") {
//	    // Admin cookie is queued
//	}
func (jar *CookieJar) HasQueued(name string, path ...string) bool {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	// Determine the path to use
	cookiePath := jar.defaultPath
	if len(path) > 0 && path[0] != "" {
		cookiePath = path[0]
	}

	// Check if the cookie exists in the queue
	key := name + "@" + cookiePath
	_, exists := jar.queuedCookies[key]
	return exists
}

// Queued retrieves a queued cookie by name and optional path.
// If path is not provided, uses the default path ("/").
// Returns nil if the cookie is not found in the queue.
//
// This allows you to inspect or modify queued cookies before they're sent.
//
// Examples:
//
//	// Get a queued cookie
//	cookie := jar.Queued("user_session")
//	if cookie != nil {
//	    fmt.Println("Session:", cookie.Value)
//	}
//
//	// Get from specific path
//	adminCookie := jar.Queued("admin_token", "/admin")
func (jar *CookieJar) Queued(name string, path ...string) *http.Cookie {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	// Determine the path to use
	cookiePath := jar.defaultPath
	if len(path) > 0 && path[0] != "" {
		cookiePath = path[0]
	}

	// Look up the cookie in the queue
	key := name + "@" + cookiePath
	if queuedCookie, exists := jar.queuedCookies[key]; exists {
		return queuedCookie.Cookie
	}

	return nil
}

// GetQueuedCookies returns all cookies currently in the queue.
// This method is typically used by middleware to process all queued cookies
// and add them to the HTTP response.
//
// The returned cookies maintain their original configuration and can be
// directly added to an HTTP response.
//
// Example:
//
//	// Process all queued cookies in middleware
//	cookies := jar.GetQueuedCookies()
//	for _, cookie := range cookies {
//	    http.SetCookie(w, cookie)
//	}
func (jar *CookieJar) GetQueuedCookies() []*http.Cookie {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	cookies := make([]*http.Cookie, 0, len(jar.queuedCookies))
	for _, queuedCookie := range jar.queuedCookies {
		cookies = append(cookies, queuedCookie.Cookie)
	}

	return cookies
}

// FlushQueuedCookies removes all cookies from the queue.
// This is useful for clearing the queue after processing or in error conditions.
//
// After calling this method, GetQueuedCookies() will return an empty slice
// and HasQueued() will return false for all cookies.
//
// Example:
//
//	// Clear the queue after processing
//	cookies := jar.GetQueuedCookies()
//	// ... process cookies ...
//	jar.FlushQueuedCookies()
func (jar *CookieJar) FlushQueuedCookies() {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	// Clear the queue by creating a new map
	jar.queuedCookies = make(map[string]*QueuedCookie)
}

// Default configuration getters and setters
// These methods provide Laravel-style configuration management

// SetDefaultPath sets the default path for new cookies.
// Returns the jar instance for method chaining.
func (jar *CookieJar) SetDefaultPath(path string) interfaces.JarInterface {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	jar.defaultPath = path
	return jar
}

// GetDefaultPath returns the current default path for cookies.
func (jar *CookieJar) GetDefaultPath() string {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	return jar.defaultPath
}

// SetDefaultDomain sets the default domain for new cookies.
// Returns the jar instance for method chaining.
func (jar *CookieJar) SetDefaultDomain(domain string) interfaces.JarInterface {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	jar.defaultDomain = domain
	return jar
}

// GetDefaultDomain returns the current default domain for cookies.
func (jar *CookieJar) GetDefaultDomain() string {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	return jar.defaultDomain
}

// SetDefaultSecure sets the default secure flag for new cookies.
// Returns the jar instance for method chaining.
func (jar *CookieJar) SetDefaultSecure(secure bool) interfaces.JarInterface {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	jar.defaultSecure = secure
	return jar
}

// GetDefaultSecure returns the current default secure flag for cookies.
func (jar *CookieJar) GetDefaultSecure() bool {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	return jar.defaultSecure
}

// SetDefaultSameSite sets the default SameSite attribute for new cookies.
// Returns the jar instance for method chaining.
func (jar *CookieJar) SetDefaultSameSite(sameSite http.SameSite) interfaces.JarInterface {
	jar.mu.Lock()
	defer jar.mu.Unlock()

	jar.defaultSameSite = sameSite
	return jar
}

// GetDefaultSameSite returns the current default SameSite attribute for cookies.
func (jar *CookieJar) GetDefaultSameSite() http.SameSite {
	jar.mu.RLock()
	defer jar.mu.RUnlock()

	return jar.defaultSameSite
}

// Compile-time interface compliance checks
var _ interfaces.JarInterface = (*CookieJar)(nil)
var _ interfaces.QueueingInterface = (*CookieJar)(nil)
