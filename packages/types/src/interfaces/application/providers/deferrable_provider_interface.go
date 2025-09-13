package interfaces

// DeferrableProvider defines the contract for service providers that can be deferred.
// Deferred providers are only loaded when one of their services is actually needed,
// which can improve application boot time by delaying the loading of providers
// until they are required.
//
// This interface follows the Interface Segregation Principle (ISP) by separating
// deferred loading concerns from the core ServiceProviderInterface, allowing
// providers to implement only the functionality they need.
type DeferrableProvider interface {
	ServiceProviderInterface

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
