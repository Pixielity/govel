package interfaces

// DeferredServiceProvider defines the interface for service providers that
// can defer their loading until one of their services is actually needed.
// This interface follows Laravel's deferred provider pattern.
//
// Deferring the loading of service providers improves application performance
// by reducing the number of providers loaded on every request. Only when one
// of the services provided by a deferred provider is requested does the
// framework load and register the provider.
//
// This pattern is particularly beneficial for:
// - Expensive services that are not used on every request
// - Third-party integrations that may be conditionally used
// - Services with heavy initialization costs
// - Optional features that can be lazily loaded
//
// Example usage:
//
//	type EmailServiceProvider struct {
//		service_providers.BaseServiceProvider
//	}
//
//	func (p *EmailServiceProvider) Register(application ContainableInterface) error {
//		return application.Singleton("mailer", func() interface{} {
//			return &MailerService{}
//		})
//	}
//
//	func (p *EmailServiceProvider) Provides() []string {
//		return []string{"mailer", "mail.manager"}
//	}
//
//	func (p *EmailServiceProvider) When() []string {
//		return []string{"mailer"} // Load only when mailer is requested
//	}
//
// The interface promotes:
// - Improved application performance through lazy loading
// - Better resource utilization
// - Modular service architecture
// - On-demand service resolution
// - Reduced memory footprint
type DeferredServiceProvider interface {
	ServiceProviderInterface

	// Provides returns the service container bindings provided by this provider.
	// The framework uses this information to determine when to load the provider.
	// Only when one of these services is requested will the provider be registered.
	//
	// The returned slice should contain all service identifiers that this
	// provider registers in the service container. These identifiers are used
	// as keys for service resolution and trigger the provider's registration
	// when requested.
	//
	// Best practices:
	// - Include all primary service identifiers
	// - Use consistent naming conventions
	// - Avoid overly generic service names
	// - Document service dependencies
	//
	// Returns:
	//   []string: Service identifiers provided by this provider
	//
	// Example:
	//   func (p *DatabaseServiceProvider) Provides() []string {
	//       return []string{
	//           "database",
	//           "database.connection",
	//           "database.manager",
	//           "db", // Common alias
	//       }
	//   }
	Provides() []string

	// When returns the service identifiers that should trigger the loading
	// of this deferred provider. This is typically the same as Provides(),
	// but can be customized for more granular control.
	//
	// This method allows providers to specify exactly which service requests
	// should cause them to be loaded, providing fine-grained control over
	// the deferred loading behavior. This is useful when:
	// - Some services should trigger loading while others shouldn't
	// - You want to provide aliases without triggering loading
	// - Conditional loading based on specific service requests
	//
	// Returns:
	//   []string: Service identifiers that trigger provider loading
	//
	// Example:
	//   func (p *CacheServiceProvider) When() []string {
	//       // Only trigger loading for primary cache service
	//       return []string{"cache", "cache.store"}
	//       // But not for "cache.config" which might be provided
	//       // without needing the full cache system
	//   }
	When() []string
}
