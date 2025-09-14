package bootstrap

import interfaces "govel/types/application"

// Providers returns all service provider instances that should be registered with the application.
//
// Providers defines the complete list of service providers in their registration
// order. Dependencies between providers should be considered when ordering to
// ensure proper initialization sequencing. The framework automatically handles
// the two-phase registration and booting process.
//
// Registration process:
//  1. Each provider's Register method is called in order
//  2. After all registrations complete, each provider's Boot method is called
//  3. Any registration or boot errors halt the application startup
//
// Provider categories:
//   - Core providers: Essential framework services (routing, database, cache)
//   - Application providers: Custom business logic and domain services
//   - Third-party providers: External package integrations
//
// Returns a slice of service provider instances in registration order.
//
// Example usage:
//
//	providers := bootstrap.Providers()
//	for _, provider := range providers {
//	    if err := app.Register(provider); err != nil {
//	        log.Fatal(err)
//	    }
//	}
func Providers() []interfaces.ServiceProvider {
	return []interfaces.ServiceProvider{
		// Example feature providers:
		// mail.NewProvider(),
		// queue.NewProvider(),
		// broadcast.NewProvider(),
		// storage.NewProvider(),
	}
}
