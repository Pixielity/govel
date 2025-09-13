package interfaces

// EventTriggeredProviderInterface defines the contract for providers that are triggered by events.
type EventTriggeredProviderInterface interface {
	// ListenForEvents returns the events that this provider should listen for.
	ListenForEvents() []string

	// HandleEvent handles the given event.
	HandleEvent(event string, payload interface{}) error

	// ShouldRegister determines if the provider should be registered based on the event.
	ShouldRegister(event string) bool
}