package interfaces

import (
	"net/http"

	"govel/packages/support/src/symbol"
)

// Service registration tokens for dependency injection
var (
	// COOKIE_JAR_TOKEN identifies the main cookie jar service
	COOKIE_JAR_TOKEN = symbol.For("govel.cookie.jar")

	// COOKIE_CSRF_TOKEN identifies the CSRF protection service
	COOKIE_CSRF_TOKEN = symbol.For("govel.cookie.csrf")

	// COOKIE_SAMESITE_TOKEN identifies the SameSite manager service
	COOKIE_SAMESITE_TOKEN = symbol.For("govel.cookie.samesite")

	// COOKIE_CONFIG_TOKEN identifies the cookie configuration
	COOKIE_CONFIG_TOKEN = symbol.For("govel.cookie.config")
)

// JarInterface defines the contract for cookie jar implementations.
// This interface mirrors Laravel's Cookie\CookieJar contract, providing
// methods for creating, retrieving, and managing HTTP cookies.
//
// The interface supports Laravel-compatible cookie operations including:
//   - Cookie creation with various options (domain, path, secure, httpOnly, sameSite)
//   - Cookie encryption and value serialization
//   - Cookie queuing for batch operations
//   - Laravel-style helper methods for common cookie operations
//
// This interface is designed to work seamlessly with Go's http.Cookie
// while providing Laravel's familiar API and advanced features.
type JarInterface interface {
	/**
	 * Create a new cookie instance.
	 *
	 * Creates a cookie with the specified name and value, using default
	 * configuration for expiration, path, domain, and security settings.
	 *
	 * @param name string The cookie name
	 * @param value string The cookie value
	 * @param options ...types.CookieOption Optional configuration functions
	 * @return *http.Cookie The created cookie instance
	 */
	Make(name, value string, options ...types.CookieOption) *http.Cookie

	/**
	 * Create a cookie that lasts "forever" (5 years).
	 *
	 * This is equivalent to Laravel's cookie()->forever() method.
	 * The cookie will expire 5 years from creation time.
	 *
	 * @param name string The cookie name
	 * @param value string The cookie value
	 * @param options ...types.CookieOption Optional configuration functions
	 * @return *http.Cookie The created cookie instance
	 */
	Forever(name, value string, options ...types.CookieOption) *http.Cookie

	/**
	 * Create a cookie that expires when the browser closes.
	 *
	 * This creates a session cookie with no explicit expiration time.
	 * The browser will delete the cookie when the session ends.
	 *
	 * @param name string The cookie name
	 * @param value string The cookie value
	 * @param options ...types.CookieOption Optional configuration functions
	 * @return *http.Cookie The created cookie instance
	 */
	Forget(name string, options ...types.CookieOption) *http.Cookie

	/**
	 * Determine if a cookie has been queued.
	 *
	 * Checks if a cookie with the given name and path combination
	 * has been added to the queue for later processing.
	 *
	 * @param name string The cookie name to check
	 * @param path string The cookie path (optional, defaults to "/")
	 * @return bool True if the cookie is queued, false otherwise
	 */
	HasQueued(name string, path ...string) bool

	/**
	 * Get a queued cookie instance.
	 *
	 * Retrieves a previously queued cookie by name and optional path.
	 * Returns nil if the cookie is not found in the queue.
	 *
	 * @param name string The cookie name to retrieve
	 * @param path string The cookie path (optional, defaults to "/")
	 * @return *http.Cookie The queued cookie or nil if not found
	 */
	Queued(name string, path ...string) *http.Cookie

	/**
	 * Queue a cookie to be sent with the next response.
	 *
	 * Adds a cookie to the queue for batch processing. Queued cookies
	 * are typically sent all at once by middleware processing.
	 *
	 * @param cookie *http.Cookie The cookie to queue
	 */
	Queue(cookie *http.Cookie)

	/**
	 * Remove a cookie from the queue.
	 *
	 * Removes a previously queued cookie by name and optional path.
	 * This prevents the cookie from being sent with the response.
	 *
	 * @param name string The cookie name to remove
	 * @param path string The cookie path (optional, defaults to "/")
	 */
	Unqueue(name string, path ...string)

	/**
	 * Get all queued cookies.
	 *
	 * Returns a slice of all cookies currently in the queue.
	 * This is typically used by middleware to process all queued cookies.
	 *
	 * @return []*http.Cookie Slice of all queued cookies
	 */
	GetQueuedCookies() []*http.Cookie

	/**
	 * Flush all queued cookies.
	 *
	 * Removes all cookies from the queue. This is useful for
	 * clearing the queue after processing or in error conditions.
	 */
	FlushQueuedCookies()

	/**
	 * Set the default path for cookies.
	 *
	 * Changes the default path that will be used for new cookies
	 * when no explicit path is provided.
	 *
	 * @param path string The default cookie path
	 * @return JarInterface Returns self for method chaining
	 */
	SetDefaultPath(path string) JarInterface

	/**
	 * Get the default path for cookies.
	 *
	 * Returns the currently configured default path for cookies.
	 *
	 * @return string The default cookie path
	 */
	GetDefaultPath() string

	/**
	 * Set the default domain for cookies.
	 *
	 * Changes the default domain that will be used for new cookies
	 * when no explicit domain is provided.
	 *
	 * @param domain string The default cookie domain
	 * @return JarInterface Returns self for method chaining
	 */
	SetDefaultDomain(domain string) JarInterface

	/**
	 * Get the default domain for cookies.
	 *
	 * Returns the currently configured default domain for cookies.
	 *
	 * @return string The default cookie domain
	 */
	GetDefaultDomain() string

	/**
	 * Set the default secure flag for cookies.
	 *
	 * Changes the default secure setting that will be used for new cookies
	 * when no explicit secure flag is provided.
	 *
	 * @param secure bool Whether cookies should be secure by default
	 * @return JarInterface Returns self for method chaining
	 */
	SetDefaultSecure(secure bool) JarInterface

	/**
	 * Get the default secure flag for cookies.
	 *
	 * Returns the currently configured default secure setting for cookies.
	 *
	 * @return bool The default secure flag
	 */
	GetDefaultSecure() bool

	/**
	 * Set the default SameSite attribute for cookies.
	 *
	 * Changes the default SameSite policy that will be used for new cookies
	 * when no explicit SameSite attribute is provided.
	 *
	 * @param sameSite http.SameSite The default SameSite policy
	 * @return JarInterface Returns self for method chaining
	 */
	SetDefaultSameSite(sameSite http.SameSite) JarInterface

	/**
	 * Get the default SameSite attribute for cookies.
	 *
	 * Returns the currently configured default SameSite policy for cookies.
	 *
	 * @return http.SameSite The default SameSite policy
	 */
	GetDefaultSameSite() http.SameSite
}
