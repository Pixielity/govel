package types

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
