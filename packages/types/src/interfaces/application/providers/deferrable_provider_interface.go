package interfaces

// DeferrableProviderInterface defines the contract for providers that can be deferred.
type DeferrableProviderInterface interface {
	// IsDeferred returns whether the provider is deferred.
	IsDeferred() bool

	// Provides returns the services that this provider provides.
	Provides() []string

	// When returns a callback that determines when the provider should be loaded.
	When() func() bool
}