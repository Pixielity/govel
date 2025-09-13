package providers

import (
	"fmt"
	serviceProviders "govel/application/providers"
	"govel/container"
	applicationInterfaces "govel/types/src/interfaces/application"
	containerInterfaces "govel/types/src/interfaces/container"
)

/**
 * ContainerServiceProvider provides comprehensive dependency injection container services.
 *
 * This service provider implements Laravel-style container functionality, binding the ContainerInterface
 * to the concrete Container implementation. It provides the core dependency injection infrastructure
 * that other service providers and application components depend on.
 *
 * Features:
 * - Service binding and resolution (Bind/Make patterns)
 * - Singleton service management with automatic caching
 * - Thread-safe container operations for concurrent usage
 * - Automatic dependency resolution and injection
 * - Service introspection capabilities (IsBound, etc.)
 * - Container lifecycle management (Flush, Forget)
 * - Memory-efficient singleton instance management
 * - Support for factory functions and concrete instances
 * - Error handling and service resolution validation
 *
 * Binding Strategy:
 * - Binds "container" abstract to ContainerInterface implementation
 * - Registers as singleton since container should be application-wide
 * - Provides the foundational IoC container for other services
 * - Supports self-registration for container accessibility
 *
 * This service provider is critical infrastructure:
 * - Must be registered first before other service providers
 * - Provides the IoC container that other providers use for registration
 * - Enables dependency injection throughout the application
 * - Supports Laravel-style service provider patterns
 *
 * Similar to Laravel's container functionality, this implementation:
 * - Provides the core IoC container for dependency injection
 * - Supports service binding with closures and concrete instances
 * - Manages singleton services with automatic caching
 * - Enables automatic dependency resolution
 * - Provides introspection methods for service discovery
 */
type ContainerServiceProvider struct {
	serviceProviders.ServiceProvider
}

// NewContainerServiceProvider creates a new container service provider with default settings.
// This constructor initializes the provider ready for registration with the application.
//
// Returns:
//
//	*ContainerServiceProvider: A new container service provider ready for registration
//
// Example:
//
//	containerProvider := providers.NewContainerServiceProvider()
//	app.RegisterProvider(containerProvider)
func NewContainerServiceProvider() *ContainerServiceProvider {
	return &ContainerServiceProvider{
		ServiceProvider: serviceProviders.ServiceProvider{},
	}
}

// Register binds the container service into the application container.
// This method implements the core container registration, making the container itself
// available as a service that can be resolved by other services.
//
// Registration Process:
// 1. Calls parent registration to set provider state
// 2. Binds "container" abstract to ContainerInterface implementation
// 3. Registers container as singleton for application-wide access
// 4. Sets up container introspection utilities
// 5. Provides container factory for creating scoped containers
//
// The container service is registered as a singleton to ensure:
// - Single source of truth for service resolution
// - Consistent service bindings across the application
// - Memory efficiency by sharing the same container instance
// - Thread-safe access to dependency injection functionality
//
// Parameters:
//
//	application: The application instance with service container access
//
// Returns:
//
//	error: Any error that occurred during registration, nil if successful
//
// Post-Registration Usage:
//
//	// Resolve container service from container (self-reference)
//	containerService, err := application.Make("container")
//	if err != nil {
//	    return fmt.Errorf("failed to resolve container service: %w", err)
//	}
//
//	// Cast to interface and use for additional service registration
//	container := containerService.(containerInterfaces.ContainerInterface)
//	container.Bind("my-service", func() interface{} {
//	    return &MyService{}
//	})
//
//	// Resolve services from the container
//	myService, err := container.Make("my-service")
//	if err != nil {
//	    return fmt.Errorf("failed to resolve my-service: %w", err)
//	}
func (p *ContainerServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register the container service as a singleton
	// This creates a self-referential binding where the container can resolve itself
	if err := application.Singleton(containerInterfaces.CONTAINER_TOKEN, p.createContainerFactory()); err != nil {
		return fmt.Errorf("failed to register container singleton: %w", err)
	}

	// Register container factory for creating scoped or temporary containers
	if err := application.Bind(containerInterfaces.CONTAINER_FACTORY_TOKEN, p.createContainerFactoryMethod()); err != nil {
		return fmt.Errorf("failed to register container factory: %w", err)
	}

	// Register container utilities for introspection and debugging
	if err := application.Bind(containerInterfaces.CONTAINER_BINDINGS_TOKEN, func() interface{} {
		// Resolve the container and call its GetBindings method
		containerService, err := application.Make(containerInterfaces.CONTAINER_TOKEN)
		if err != nil {
			return map[string]interface{}{"error": "failed to resolve container"}
		}
		container := containerService.(containerInterfaces.ContainerInterface)
		return container.GetBindings()
	}); err != nil {
		return fmt.Errorf("failed to register container bindings introspector: %w", err)
	}

	// Register container statistics for monitoring
	if err := application.Bind(containerInterfaces.CONTAINER_STATS_TOKEN, func() interface{} {
		// Resolve the container and call its GetStatistics method
		containerService, err := application.Make(containerInterfaces.CONTAINER_TOKEN)
		if err != nil {
			return map[string]interface{}{"error": "failed to resolve container"}
		}
		container := containerService.(containerInterfaces.ContainerInterface)
		return container.GetStatistics()
	}); err != nil {
		return fmt.Errorf("failed to register container statistics: %w", err)
	}

	return nil
}

// createContainerFactory creates the main container service factory function.
// This factory creates a new container instance or returns the existing application container.
//
// Returns:
//
//	func() interface{}: Factory function that creates ContainerInterface instance
func (p *ContainerServiceProvider) createContainerFactory() func() interface{} {
	return func() interface{} {
		// Create a new container instance
		containerInstance := container.New()

		// Return as ContainerInterface to maintain interface segregation
		return containerInterface(containerInstance)
	}
}

// createContainerFactoryMethod creates a factory method for creating additional container instances.
// This is useful for creating scoped containers or isolated dependency injection contexts.
//
// Returns:
//
//	func() interface{}: Factory method that returns a container creation function
func (p *ContainerServiceProvider) createContainerFactoryMethod() func() interface{} {
	return func() interface{} {
		return func() containerInterfaces.ContainerInterface {
			containerInstance := container.New()
			return containerInterface(containerInstance)
		}
	}
}

// containerInterface is a type assertion helper that ensures the container instance
// implements the ContainerInterface. This provides compile-time safety for the binding.
//
// Parameters:
//
//	containerInstance: The concrete ServiceContainer instance
//
// Returns:
//
//	containerInterfaces.ContainerInterface: The container instance as an interface
func containerInterface(containerInstance *container.ServiceContainer) containerInterfaces.ContainerInterface {
	return containerInstance
}

// Priority returns the registration priority for this service provider.
// Container services have the highest priority since all other services depend on the container.
//
// Returns:
//
//	int: Priority level 10 for critical infrastructure services
func (p *ContainerServiceProvider) Priority() int {
	return 10 // Highest priority - container is fundamental infrastructure
}
