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

// ClientDeferredServiceProvider provides deferred client services to the application.
// This provider implements the DeferrableProvider interface, which means it will only
// be loaded and registered when one of its services is actually requested.
//
// This can improve application boot time by delaying the loading of services
// until they are actually needed.
type ClientDeferredServiceProvider struct {
	*baseProviders.ServiceProvider
	isRegistered bool
	isBooted     bool
	logger       interface{}
}

// NewClientDeferredServiceProvider creates a new deferred client service provider.
//
// Returns:
//
//	*ClientDeferredServiceProvider: A new deferred client service provider instance
func NewClientDeferredServiceProvider() *ClientDeferredServiceProvider {
	fmt.Println("[DEBUG] Creating new ClientDeferredServiceProvider (Deferred)")
	return &ClientDeferredServiceProvider{
		ServiceProvider: baseProviders.NewServiceProvider(nil), // Will be set later
		isRegistered:    false,
		isBooted:        false,
	}
}

// Register registers client services with the application container.
// This method is called only when one of the provided services is requested.
func (p *ClientDeferredServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] ClientDeferredServiceProvider.Register() called (DEFERRED LOADING)")
	if p.isRegistered {
		fmt.Println("[DEBUG] ClientDeferredServiceProvider already registered, skipping")
		return nil
	}

	// Get container and logger
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	// Register the deferred client service
	if err := container.Bind("client.deferred.service", func() interface{} {
		logger.Info("Creating Deferred ClientService instance")
		return clientServices.NewClientService(logger)
	}); err != nil {
		return fmt.Errorf("failed to register deferred client service: %w", err)
	}

	logger.Info("ClientDeferredServiceProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization for the deferred client services.
// This method is called after the provider has been registered due to service request.
func (p *ClientDeferredServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting ClientDeferredServiceProvider...")

	// Perform any boot-time initialization here
	// This happens only when the service is actually needed

	p.isBooted = true
	logger.Info("ClientDeferredServiceProvider booted successfully")
	return nil
}

// GetProvides returns the services provided by this provider.
// This is used by the base ServiceProvider.IsDeferred() method.
func (p *ClientDeferredServiceProvider) GetProvides() []string {
	return []string{"client.deferred.service"}
}

// Provides returns the services provided by the deferred provider.
// This method is required by the DeferrableProvider interface.
func (p *ClientDeferredServiceProvider) Provides() []string {
	return []string{"client.deferred.service"}
}

// IsDeferred determines if the provider is deferred.
// This method is required by the DeferrableProvider interface.
func (p *ClientDeferredServiceProvider) IsDeferred() bool {
	return true // This provider should always be deferred
}

// Terminate gracefully shuts down the deferred client service provider.
func (p *ClientDeferredServiceProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating ClientDeferredServiceProvider...")

	select {
	case <-time.After(100 * time.Millisecond):
		logger.Info("ClientDeferredServiceProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("ClientDeferredServiceProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("ClientDeferredServiceProvider terminated successfully")
	return nil
}

// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*ClientDeferredServiceProvider)(nil)
var _ providerInterfaces.DeferrableProvider = (*ClientDeferredServiceProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*ClientDeferredServiceProvider)(nil)
