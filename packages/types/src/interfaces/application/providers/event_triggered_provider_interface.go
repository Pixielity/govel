package interfaces

// EventTriggeredProvider defines the contract for providers that can be triggered by events.
// This is an optional interface that can be implemented by providers that need to be loaded
// in response to specific events rather than just service resolution.
//
// This interface follows the Interface Segregation Principle (ISP) by separating event-based
// triggering concerns from the core DeferrableProvider interface, allowing providers to
// implement only the functionality they need.
type EventTriggeredProvider interface {
	// When returns the events that should trigger this provider to be loaded.
	// This allows the provider to be loaded in response to specific events
	// rather than just service resolution.
	//
	// The returned slice should contain event names that are meaningful to the
	// application's event system. When any of these events are fired, the
	// provider will be loaded if it hasn't been loaded already.
	//
	// Returns:
	//   []string: A slice of event names that should trigger provider loading
	//
	// Example:
	//   func (p *MyProvider) When() []string {
	//       return []string{"booting", "request.started", "user.authenticated"}
	//   }
	//
	// Note:
	//   - Returning an empty slice or nil means the provider won't be triggered by events
	//   - Event names should be consistent with the application's event naming conventions
	//   - The provider must also implement DeferrableProvider to be eligible for deferred loading
	When() []string
}
