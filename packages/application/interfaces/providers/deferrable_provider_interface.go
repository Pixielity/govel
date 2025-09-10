package interfaces

// DeferrableProvider defines the contract for service providers that can be deferred.
// This interface follows Laravel's DeferrableProvider pattern, allowing providers
// to specify which services they provide and when they should be loaded.
//
// Deferred providers are only loaded when one of their services is actually needed,
// which can significantly improve application boot time by delaying the loading
// of providers until they are required.
type DeferrableProvider interface {
	// Provides returns the services provided by the provider.
	// This method must return a slice of service identifiers that this provider
	// can resolve. The provider will only be loaded when one of these services
	// is requested from the container.
	//
	// Returns:
	//   []string: A slice of service identifiers that this provider can resolve
	//
	// Example:
	//   func (p *MyProvider) Provides() []string {
	//       return []string{"cache", "cache.store", "memcached"}
	//   }
	Provides() []string

	// IsDeferred determines if the provider is deferred.
	// Deferred providers are only loaded when one of their services is actually needed.
	// This can improve application boot time by delaying the loading of providers
	// until they are required.
	//
	// Returns:
	//   bool: true if the provider should be deferred, false otherwise
	IsDeferred() bool
}
