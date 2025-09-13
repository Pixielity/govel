package interfaces

import (
	types "govel/types/types/container"
)

// ContainerInterface defines the contract for dependency injection containers.
type ContainerInterface interface {
	// Bind registers a binding with the container.
	Bind(abstract types.ServiceIdentifier, concrete interface{}) error

	// Singleton registers a shared binding in the container.
	Singleton(abstract types.ServiceIdentifier, concrete interface{}) error

	// Make resolves and returns an instance from the container.
	Make(abstract types.ServiceIdentifier) (interface{}, error)

	// IsBound checks if a service is registered in the container.
	IsBound(abstract types.ServiceIdentifier) bool

	// IsSingleton checks if a service is registered as a singleton.
	IsSingleton(abstract types.ServiceIdentifier) bool

	// Forget removes a service binding from the container.
	Forget(abstract types.ServiceIdentifier)

	// FlushContainer removes all service bindings and cached instances.
	FlushContainer()

	// Count returns the total number of registered services.
	Count() int

	// RegisteredServices returns a list of all registered service names.
	RegisteredServices() []string

	// GetBindings returns detailed information about all service bindings.
	GetBindings() map[string]interface{}

	// GetStatistics returns container usage statistics and performance metrics.
	GetStatistics() map[string]interface{}
}
