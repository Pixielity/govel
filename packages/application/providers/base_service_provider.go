package providers

import (
	applicationInterfaces "govel/packages/application/interfaces/application"
)

/**
 * BaseServiceProvider provides the fundamental implementation for all service providers.
 *
 * This service provider binds the ServiceProviderInterface to the BaseProvider implementation,
 * enabling dependency injection and proper lifecycle management for the base provider functionality.
 *
 * Features:
 * - Common service provider lifecycle methods
 * - Container management
 * - Registration and boot state tracking
 * - Base implementation for all other providers
 */
type BaseServiceProvider struct {
	registered bool
	booted     bool
}

// Register provides a default implementation of the Register method.
// Base implementation does nothing - concrete providers should override this method
// to perform their specific service registration logic.
//
// The registration phase should only bind services into the container.
// Do not perform any operations that depend on other services being available,
// as the boot order is not guaranteed during registration.
//
// Parameters:
//
//	application: The application instance for service container access
//
// Returns:
//
//	error: Always returns nil in the base implementation
//
// Example:
//
//	func (p *MyServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
//	    return application.Singleton("my-service", func() interface{} {
//	        return &MyService{}
//	    })
//	}
func (p *BaseServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	p.registered = true
	return nil
}

// Boot provides a default implementation of the Boot method.
// Base implementation does nothing - concrete providers should override this method
// to perform their specific bootstrap logic.
//
// The boot phase is called after all providers have been registered,
// so it's safe to resolve services from the container and perform
// initialization that depends on other services.
//
// Parameters:
//
//	application: The application instance
//
// Returns:
//
//	error: Always returns nil in the base implementation
//
// Example:
//
//	func (p *MyServiceProvider) Boot(application applicationInterfaces.ApplicationInterface) error {
//	    myService, _ := application.Make("my-service")
//	    return myService.(*MyService).Initialize()
//	}
func (p *BaseServiceProvider) Boot(application applicationInterfaces.ApplicationInterface) error {
	p.booted = true
	return nil
}
