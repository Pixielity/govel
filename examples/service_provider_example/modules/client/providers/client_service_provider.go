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

// ClientServiceProvider provides standard (eager) client services to the application.
// This provider registers immediately during application startup and provides
// all client-related services without any deferral or event triggers.
//
// This is the simplest type of service provider - it gets loaded, registered,
// and booted immediately when the application starts up.
type ClientServiceProvider struct {
	*baseProviders.ServiceProvider
	isRegistered bool
	isBooted     bool
	logger       interface{}
}

// NewClientServiceProvider creates a new standard client service provider.
//
// Returns:
//   *ClientServiceProvider: A new client service provider instance
func NewClientServiceProvider() *ClientServiceProvider {
	fmt.Println("[DEBUG] Creating new ClientServiceProvider (Standard/Eager)")
	return &ClientServiceProvider{
		ServiceProvider: baseProviders.NewServiceProvider(nil), // Will be set later
		isRegistered:    false,
		isBooted:        false,
	}
}

// Register registers client services with the application container.
// This method is called immediately during application startup for eager providers.
func (p *ClientServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] ClientServiceProvider.Register() called")
	if p.isRegistered {
		fmt.Println("[DEBUG] ClientServiceProvider already registered, skipping")
		return nil
	}

	// Get container and logger
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	// Register the client service
	if err := container.Bind("client.service", func() interface{} {
		logger.Info("Creating ClientService instance")
		return clientServices.NewClientService(logger)
	}); err != nil {
		return fmt.Errorf("failed to register client service: %w", err)
	}

	logger.Info("ClientServiceProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization for the client services.
// This method is called after all providers have been registered.
func (p *ClientServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting ClientServiceProvider...")

	// Perform any boot-time initialization here
	// For example, you might warm up caches, establish connections, etc.

	p.isBooted = true
	logger.Info("ClientServiceProvider booted successfully")
	return nil
}

// GetProvides returns the services provided by this provider.
func (p *ClientServiceProvider) GetProvides() []string {
	return []string{"client.service"}
}

// Terminate gracefully shuts down the client service provider.
func (p *ClientServiceProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating ClientServiceProvider...")

	select {
	case <-time.After(100 * time.Millisecond):
		logger.Info("ClientServiceProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("ClientServiceProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("ClientServiceProvider terminated successfully")
	return nil
}

// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*ClientServiceProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*ClientServiceProvider)(nil)
