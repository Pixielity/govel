package providers

import (
	"context"
	"fmt"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	baseProviders "govel/packages/application/providers"

	clientServices "service_provider_example/modules/client/services"
)

// ClientEventServiceProvider provides event-triggered deferred client services.
// This provider implements both DeferrableProvider and EventTriggeredProvider interfaces,
// which means it will only be loaded when:
// 1. One of its services is requested, OR
// 2. One of its trigger events is fired
//
// This provides the most flexibility for when services should be loaded,
// allowing both on-demand service requests and event-based triggers.
type ClientEventServiceProvider struct {
	*baseProviders.ServiceProvider
	isRegistered bool
	isBooted     bool
	logger       interface{}
}

// NewClientEventServiceProvider creates a new event-triggered deferred client service provider.
//
// Returns:
//
//	*ClientEventServiceProvider: A new event-triggered client service provider instance
func NewClientEventServiceProvider() *ClientEventServiceProvider {
	fmt.Println("[DEBUG] Creating new ClientEventServiceProvider (Deferred + Event Triggered)")
	return &ClientEventServiceProvider{
		ServiceProvider: baseProviders.NewServiceProvider(nil), // Will be set later
		isRegistered:    false,
		isBooted:        false,
	}
}

// Register registers client services with the application container.
// This method is called only when:
// - One of the provided services is requested, OR
// - One of the trigger events is fired
func (p *ClientEventServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] ClientEventServiceProvider.Register() called (DEFERRED + EVENT TRIGGERED)")
	if p.isRegistered {
		fmt.Println("[DEBUG] ClientEventServiceProvider already registered, skipping")
		return nil
	}

	// Get container and logger
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	// Register the event-triggered client service
	if err := container.Bind("client.event.service", func() interface{} {
		logger.Info("Creating Event-Triggered ClientService instance")
		return clientServices.NewClientService(logger)
	}); err != nil {
		return fmt.Errorf("failed to register event-triggered client service: %w", err)
	}

	// Also register a specialized analytics service for event handling
	if err := container.Bind("client.analytics.service", func() interface{} {
		logger.Info("Creating Client Analytics Service instance")
		// In a real application, this might be a different service implementation
		return clientServices.NewClientService(logger)
	}); err != nil {
		return fmt.Errorf("failed to register client analytics service: %w", err)
	}

	logger.Info("ClientEventServiceProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization for the event-triggered client services.
// This method is called after the provider has been registered due to service request or event.
func (p *ClientEventServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting ClientEventServiceProvider...")

	// Perform any boot-time initialization here
	// This might include setting up event listeners, warming caches, etc.

	p.isBooted = true
	logger.Info("ClientEventServiceProvider booted successfully")
	return nil
}

// GetProvides returns the services provided by this provider.
// This is used by the base ServiceProvider.IsDeferred() method.
func (p *ClientEventServiceProvider) GetProvides() []string {
	return []string{"client.event.service", "client.analytics.service"}
}

// Provides returns the services provided by the deferred provider.
// This method is required by the DeferrableProvider interface.
func (p *ClientEventServiceProvider) Provides() []string {
	return []string{"client.event.service", "client.analytics.service"}
}

// IsDeferred determines if the provider is deferred.
// This method is required by the DeferrableProvider interface.
func (p *ClientEventServiceProvider) IsDeferred() bool {
	return true // This provider should always be deferred
}

// When returns the events that should trigger this provider to be loaded.
// This method is required by the EventTriggeredProvider interface.
func (p *ClientEventServiceProvider) When() []string {
	return []string{
		"client.requested",
		"analytics.requested",
		"user.created",
		"application.boot.complete",
	}
}

// GetLoadEvents returns events that should trigger loading this provider.
// This is an alias for When() to maintain consistency with other patterns.
func (p *ClientEventServiceProvider) GetLoadEvents() []string {
	return p.When()
}

// Terminate gracefully shuts down the event-triggered client service provider.
func (p *ClientEventServiceProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating ClientEventServiceProvider...")

	select {
	case <-time.After(100 * time.Millisecond):
		logger.Info("ClientEventServiceProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("ClientEventServiceProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("ClientEventServiceProvider terminated successfully")
	return nil
}

// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*ClientEventServiceProvider)(nil)
var _ providerInterfaces.DeferrableProvider = (*ClientEventServiceProvider)(nil)
var _ providerInterfaces.EventTriggeredProvider = (*ClientEventServiceProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*ClientEventServiceProvider)(nil)
