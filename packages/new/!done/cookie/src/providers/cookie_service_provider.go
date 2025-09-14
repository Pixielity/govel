// Package providers contains service provider implementations for the cookie package.
// Service providers are responsible for registering cookie services in the
// dependency injection container and configuring them for use throughout the application.
package providers

import (
	"fmt"

	"govel/application/providers"
	applicationInterfaces "govel/types/interfaces/application/base"

	cookie "govel/cookie"
	cookieInterfaces "govel/cookie/interfaces"
)

// CookieServiceProvider implements a Laravel-compatible service provider
// for the cookie package. It registers cookie services in the dependency
// injection container and provides factory methods for creating cookie jar instances.
//
// This service provider is equivalent to Laravel's CookieServiceProvider and
// follows the same patterns for service registration and binding.
//
// Services registered:
//   - "cookie": Singleton CookieJar instance for managing cookie creation and queuing
//   - "cookie.jar": Factory for getting the default cookie jar driver
//   - cookieInterfaces.JarInterface: Default cookie jar interface contract
//   - cookieInterfaces.QueueingInterface: Cookie queuing interface contract
//
// The provider implements deferred loading, meaning services are only
// instantiated when first requested from the container.
type CookieServiceProvider struct {
	providers.ServiceProvider
}

// NewCookieServiceProvider creates a new CookieServiceProvider instance.
//
// Returns:
//   - *CookieServiceProvider: New service provider instance ready for registration
//
// Example:
//
//	provider := NewCookieServiceProvider()
//	err := provider.Register(application)
//	if err != nil {
//		log.Fatal("Failed to register cookie services:", err)
//	}
func NewCookieServiceProvider() *CookieServiceProvider {
	return &CookieServiceProvider{
		ServiceProvider: providers.ServiceProvider{},
	}
}

// Register registers all cookie services in the dependency injection container.
// This method sets up the service bindings but does not instantiate the services
// immediately (deferred loading).
//
// Registered services:
//   - "cookie": Singleton CookieJar instance for managing cookie creation and queuing
//   - "cookie.jar": Factory function that returns the default cookie jar driver
//   - "cookie.contract": CookieJar interface contract
//   - "cookie.queuing.contract": Cookie queuing interface contract
//
// Returns:
//   - error: Any error that occurred during service registration
//
// Thread-safe: This method should only be called once during application bootstrap.
func (h *CookieServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := h.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register CookieJar as a singleton
	// This ensures only one CookieJar instance exists throughout the application
	err := application.Singleton(cookieInterfaces.COOKIE_JAR_TOKEN, h.createCookieJarFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind cookie jar: %w", err)
	}

	// Register default cookie jar factory
	// This provides quick access to the default cookie jar without manager overhead
	err = application.Bind(cookieInterfaces.COOKIE_JAR_TOKEN, h.createCookieJarFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind cookie jar driver: %w", err)
	}

	// Register JarInterface contract using the CookieJar
	err = application.Singleton(cookieInterfaces.COOKIE_JAR_TOKEN, h.createCookieJarFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind cookie jar contract: %w", err)
	}

	return nil
}

// Provides returns a list of service tokens that this provider offers.
// This is used by the container to determine which services are available
// and to implement deferred loading.
//
// Returns:
//   - []interface{}: List of service tokens provided by this provider
func (h *CookieServiceProvider) Provides() []interface{} {
	return []interface{}{
		cookieInterfaces.COOKIE_JAR_TOKEN, // CookieJar singleton
	}
}

// createCookieJarFactory creates a factory function for the CookieJar singleton.
// This factory returns the same CookieJar instance on subsequent calls.
func (h *CookieServiceProvider) createCookieJarFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		return cookie.NewCookieJar()
	}
}
