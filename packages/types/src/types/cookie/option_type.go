package cookie

import "net/http"

/**
 * CookieOption defines a function type for configuring HTTP cookies.
 *
 * This type represents a functional option pattern for cookie configuration,
 * allowing flexible and extensible cookie customization. Each CookieOption
 * function receives a pointer to an http.Cookie and modifies its properties.
 *
 * Usage Example:
 *   cookie := &http.Cookie{Name: "session", Value: "abc123"}
 *   
 *   // Apply options to configure the cookie
 *   options := []CookieOption{
 *       WithDomain("example.com"),
 *       WithPath("/api"),
 *       WithSecure(true),
 *       WithHttpOnly(true),
 *   }
 *   
 *   for _, option := range options {
 *       option(cookie)
 *   }
 *
 * The functional option pattern provides several benefits:
 *   - Type-safe configuration with compile-time checking
 *   - Flexible combination of options in any order
 *   - Extensible design allowing new options without breaking changes
 *   - Clear, readable configuration code
 *   - Optional parameters with sensible defaults
 *
 * Common cookie options that implement this type include:
 *   - WithExpiry(time.Time) - Sets cookie expiration time
 *   - WithMaxAge(int) - Sets cookie max age in seconds
 *   - WithDomain(string) - Sets cookie domain restriction
 *   - WithPath(string) - Sets cookie path scope
 *   - WithSecure(bool) - Enables/disables HTTPS-only transmission
 *   - WithHttpOnly(bool) - Enables/disables JavaScript access prevention
 *   - WithSameSite(http.SameSite) - Sets SameSite attribute for CSRF protection
 *
 * Thread Safety:
 *   Individual CookieOption functions should be thread-safe when operating
 *   on the provided cookie pointer. The cookie being configured should not
 *   be accessed concurrently during option application.
 *
 * Performance Considerations:
 *   CookieOption functions are typically lightweight and execute quickly.
 *   The functional call overhead is minimal compared to cookie processing
 *   and network transmission costs.
 */
type CookieOption func(*http.Cookie)

/**
 * ApplyOptions is a utility function that applies multiple CookieOption functions
 * to a given http.Cookie in sequence.
 *
 * This function provides a convenient way to apply a slice of options to a cookie,
 * ensuring all configurations are applied in the order specified.
 *
 * Parameters:
 *   cookie *http.Cookie - The cookie to configure (must not be nil)
 *   options []CookieOption - Slice of option functions to apply
 *
 * Example Usage:
 *   cookie := &http.Cookie{Name: "token", Value: "xyz789"}
 *   options := []CookieOption{
 *       WithExpiry(time.Now().Add(24 * time.Hour)),
 *       WithSecure(true),
 *       WithHttpOnly(true),
 *   }
 *   ApplyOptions(cookie, options)
 *
 * Error Handling:
 *   This function does not return errors. Individual option functions
 *   should handle invalid inputs gracefully or panic if appropriate.
 *   Nil cookie pointers will cause a panic.
 *
 * Thread Safety:
 *   This function is not thread-safe. The provided cookie should not be
 *   accessed concurrently during option application.
 */
func ApplyOptions(cookie *http.Cookie, options []CookieOption) {
	for _, option := range options {
		if option != nil {
			option(cookie)
		}
	}
}

/**
 * ChainOptions combines multiple CookieOption functions into a single option.
 *
 * This utility function allows composition of multiple options into one,
 * which can be useful for creating reusable option combinations or
 * building complex configuration patterns.
 *
 * Parameters:
 *   options ...CookieOption - Variable number of options to chain together
 *
 * Returns:
 *   CookieOption - A single option function that applies all provided options
 *
 * Example Usage:
 *   // Create a secure session cookie option
 *   secureSession := ChainOptions(
 *       WithSecure(true),
 *       WithHttpOnly(true),
 *       WithSameSite(http.SameSiteStrictMode),
 *   )
 *   
 *   // Use the chained option
 *   cookie := &http.Cookie{Name: "session", Value: "token"}
 *   secureSession(cookie)
 *
 * Performance:
 *   The returned function has minimal overhead and applies options efficiently.
 *   Chaining is resolved at creation time, not at application time.
 *
 * Thread Safety:
 *   The returned CookieOption function is thread-safe for read operations
 *   but should not be applied to the same cookie concurrently.
 */
func ChainOptions(options ...CookieOption) CookieOption {
	return func(cookie *http.Cookie) {
		ApplyOptions(cookie, options)
	}
}

/**
 * ConditionalOption applies an option only if a condition is met.
 *
 * This utility function provides conditional option application,
 * allowing dynamic cookie configuration based on runtime conditions.
 *
 * Parameters:
 *   condition bool - Whether to apply the option
 *   option CookieOption - The option to apply if condition is true
 *
 * Returns:
 *   CookieOption - A conditional option function
 *
 * Example Usage:
 *   isProduction := os.Getenv("ENV") == "production"
 *   conditionalSecure := ConditionalOption(isProduction, WithSecure(true))
 *   
 *   cookie := &http.Cookie{Name: "session", Value: "token"}
 *   conditionalSecure(cookie) // Only sets Secure=true in production
 *
 * Use Cases:
 *   - Environment-specific cookie configuration
 *   - Feature flag based cookie settings  
 *   - Dynamic security policy application
 *   - Conditional debugging or development options
 */
func ConditionalOption(condition bool, option CookieOption) CookieOption {
	return func(cookie *http.Cookie) {
		if condition && option != nil {
			option(cookie)
		}
	}
}
