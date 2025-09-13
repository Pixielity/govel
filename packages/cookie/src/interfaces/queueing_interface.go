package interfaces

import "net/http"

// QueueingInterface defines the contract for cookie queuing operations.
// This interface provides methods for managing cookies that should be
// sent with HTTP responses, supporting Laravel's cookie queuing system.
//
// In Laravel, cookies can be queued during request processing and then
// automatically attached to the response by middleware. This interface
// provides the same functionality for Go applications.
//
// Features:
//   - Queue cookies for batch processing
//   - Retrieve queued cookies by name and path
//   - Check if specific cookies are queued
//   - Clear individual or all queued cookies
//   - Support for path-specific cookie queuing
type QueueingInterface interface {
	/**
	 * Queue a cookie to be sent with the next response.
	 *
	 * Adds a cookie to the internal queue. The cookie will be
	 * processed and sent with the HTTP response by middleware.
	 *
	 * This method allows you to queue cookies during request processing
	 * without immediately sending them. This is useful for:
	 *   - Deferring cookie setting until response preparation
	 *   - Batch processing multiple cookies
	 *   - Conditional cookie setting based on request flow
	 *
	 * @param cookie *http.Cookie The cookie to add to the queue
	 */
	Queue(cookie *http.Cookie)

	/**
	 * Remove a cookie from the queue.
	 *
	 * Removes a previously queued cookie by name and optional path.
	 * If path is not provided, removes the cookie from the default path ("/").
	 *
	 * This is useful for:
	 *   - Cancelling cookies that were queued earlier in the request
	 *   - Removing cookies based on conditional logic
	 *   - Cleaning up queued cookies in error conditions
	 *
	 * @param name string The name of the cookie to remove
	 * @param path string The path of the cookie (optional, defaults to "/")
	 */
	Unqueue(name string, path ...string)

	/**
	 * Get all queued cookies.
	 *
	 * Returns a slice containing all cookies currently in the queue.
	 * This method is typically used by middleware to retrieve all
	 * queued cookies and add them to the HTTP response.
	 *
	 * The returned cookies maintain their original configuration
	 * including expiration, domain, path, secure, and SameSite settings.
	 *
	 * @return []*http.Cookie A slice of all queued cookies
	 */
	GetQueuedCookies() []*http.Cookie

	/**
	 * Flush all queued cookies.
	 *
	 * Removes all cookies from the queue, effectively clearing
	 * the entire queue. This is useful for:
	 *   - Resetting the queue state
	 *   - Cleaning up after processing all queued cookies
	 *   - Handling error conditions where queued cookies should be discarded
	 *
	 * After calling this method, GetQueuedCookies() will return an empty slice.
	 */
	FlushQueuedCookies()

	/**
	 * Determine if a cookie has been queued.
	 *
	 * Checks whether a cookie with the specified name and path
	 * is currently in the queue. If path is not provided,
	 * checks for the cookie in the default path ("/").
	 *
	 * This method is useful for:
	 *   - Avoiding duplicate cookie queuing
	 *   - Conditional cookie processing
	 *   - Debugging cookie queue state
	 *
	 * @param name string The name of the cookie to check
	 * @param path string The path of the cookie (optional, defaults to "/")
	 * @return bool True if the cookie is queued, false otherwise
	 */
	HasQueued(name string, path ...string) bool

	/**
	 * Get a queued cookie instance.
	 *
	 * Retrieves a specific cookie from the queue by name and optional path.
	 * If path is not provided, looks for the cookie in the default path ("/").
	 *
	 * This method returns the actual cookie instance, allowing you to:
	 *   - Inspect queued cookie properties
	 *   - Modify queued cookies before they're sent
	 *   - Retrieve cookie values for processing
	 *
	 * @param name string The name of the cookie to retrieve
	 * @param path string The path of the cookie (optional, defaults to "/")
	 * @return *http.Cookie The queued cookie instance, or nil if not found
	 */
	Queued(name string, path ...string) *http.Cookie
}
