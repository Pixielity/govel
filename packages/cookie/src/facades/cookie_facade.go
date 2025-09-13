package facades

import (
	cookieInterfaces "govel/cookie/src/interfaces"
	facade "govel/support/src"
)

// Cookie returns the cookie jar instance from the service container.
// This is the main entry point for cookie operations via the facade.
//
// The method resolves the cookie jar service from the dependency injection
// container and provides type-safe access to cookie functionality.
//
// Returns:
//   - cookieInterfaces.JarInterface: The cookie jar service instance
//
// Panics:
//   - If the cookie service is not registered in the container
//   - If the container is not properly configured
//   - If the resolved service doesn't implement JarInterface
//
// Example:
//
//	jar := facades.Cookie()
//	cookie := jar.Make("session_id", sessionID)
func Cookie() cookieInterfaces.JarInterface {
	return facade.Resolve[cookieInterfaces.JarInterface](cookieInterfaces.COOKIE_JAR_TOKEN)
}

// CookieWithError provides error-safe access to the cookie jar service.
// This method returns errors instead of panicking, making it suitable for
// error-sensitive contexts where cookie service unavailability should be
// handled gracefully.
//
// Returns:
//   - cookieInterfaces.JarInterface: The cookie jar service instance (nil on error)
//   - error: Detailed error information if service resolution fails
//
// Example:
//
//	jar, err := facades.CookieWithError()
//	if err != nil {
//	    log.Printf("Cookie service unavailable: %v", err)
//	    return handleCookieError()
//	}
//	return jar.Make("fallback_cookie", "value")
func CookieWithError() (cookieInterfaces.JarInterface, error) {
	return facade.TryResolve[cookieInterfaces.JarInterface](cookieInterfaces.COOKIE_JAR_TOKEN)
}
