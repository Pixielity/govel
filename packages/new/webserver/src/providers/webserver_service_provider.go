// Package providers contains service providers for integrating the webserver package with the Govel framework.
// This file implements the WebserverServiceProvider that registers webserver services into the Govel container,
// aligned in structure and style with the PathsServiceProvider.
package providers

import (
	"fmt"

	"govel/application/providers"
	"govel/new/webserver/enums"
	"govel/new/webserver/factories"
	applicationInterfaces "govel/types/interfaces/application/base"
)

// WebserverServiceProvider provides webserver services through the container.
//
// This provider binds webserver-related factories and helpers to the container,
// enabling dependency injection and consistent access across the application.
// It mirrors the style of PathsServiceProvider: concise registration via application.Bind
// with clear service keys and closures.
//
// Bound Services:
//   - "webserver.factory": WebserverFactory for creating webserver instances from adapters
//   - "webserver.adapter.factory": AdapterFactory for creating adapters by engine
//   - "webserver.middleware.factory": MiddlewareFactory for creating middleware and chains
//   - "webserver.create": Function(engine string, config map[string]interface{}) (interfaces.WebserverInterface, error)
//   - "webserver.default": Function() interfaces.WebserverInterface — creates a default-engine server
//
// Usage:
//
//	webFactoryAny, _ := container.Make("webserver.factory")
//	webFactory := webFactoryAny.(*factories.WebserverFactory)
//	adapterFactoryAny, _ := container.Make("webserver.adapter.factory")
//	adapterFactory := adapterFactoryAny.(*factories.AdapterFactory)
//
//	adapter, _ := adapterFactory.CreateAdapter("gin")
//	server := webFactory.CreateWebserverInstance(adapter)
//	server.Listen(":8080")
type WebserverServiceProvider struct {
	providers.ServiceProvider
}

// NewWebserverServiceProvider creates a new webserver service provider.
//
// Returns:
//   - *WebserverServiceProvider: A new provider ready for registration
func NewWebserverServiceProvider() *WebserverServiceProvider {
	return &WebserverServiceProvider{
		ServiceProvider: providers.ServiceProvider{},
	}
}

// Register binds webserver services into the service container.
// This method mirrors PathsServiceProvider by using application.Bind with closures.
//
// Parameters:
//   - application: The application instance with binding capabilities
//
// Returns:
//   - error: Any error that occurred during registration
func (p *WebserverServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register to set registered flag consistently
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Bind webserver factory
	if err := application.Bind("webserver.factory", func() interface{} {
		return factories.NewWebserverFactory()
	}); err != nil {
		return fmt.Errorf("failed to register webserver.factory: %w", err)
	}

	// Bind adapter factory
	if err := application.Bind("webserver.adapter.factory", func() interface{} {
		return factories.NewAdapterFactory()
	}); err != nil {
		return fmt.Errorf("failed to register webserver.adapter.factory: %w", err)
	}

	// Bind middleware factory
	if err := application.Bind("webserver.middleware.factory", func() interface{} {
		return factories.NewMiddlewareFactory()
	}); err != nil {
		return fmt.Errorf("failed to register webserver.middleware.factory: %w", err)
	}

	// Bind webserver.create helper
	if err := application.Bind("webserver.create", func() interface{} {
		// Return a function for on-demand creation with engine and config
		return func(engine enums.Engine, config map[string]interface{}) (interface{}, error) {
			adapterFactory := factories.NewAdapterFactory()
			adapter, err := adapterFactory.CreateAdapter(engine)
			if err != nil {
				return nil, err
			}
			webFactory := factories.NewWebserverFactory()
			server := webFactory.CreateWebserverInstance(adapter)
			// Apply optional config
			if config != nil {
				for k, v := range config {
					server.SetConfig(k, v)
				}
			}
			return server, nil
		}
	}); err != nil {
		return fmt.Errorf("failed to register webserver.create: %w", err)
	}

	// Bind webserver.default helper (uses default engine)
	if err := application.Bind("webserver.default", func() interface{} {
		return func() interface{} {
			adapterFactory := factories.NewAdapterFactory()
			adapter, err := adapterFactory.CreateAdapter(enums.DefaultEngine())
			if err != nil {
				// On error, return nil to avoid panics; caller should handle nil
				return nil
			}
			webFactory := factories.NewWebserverFactory()
			return webFactory.CreateWebserverInstance(adapter)
		}
	}); err != nil {
		return fmt.Errorf("failed to register webserver.default: %w", err)
	}

	return nil
}

// Provides returns the list of services provided by this provider.
// This helps the application understand what can be resolved from the container.
//
// Returns:
//   - []string: The service identifiers bound by this provider
func (p *WebserverServiceProvider) Provides() []string {
	return []string{
		"webserver.factory",
		"webserver.adapter.factory",
		"webserver.middleware.factory",
		"webserver.create",
		"webserver.default",
	}
}
