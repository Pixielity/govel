// Package container provides the service container implementation for the GoVel framework.
// This package implements dependency injection and service resolution following
// Laravel's container patterns.
//
// The container supports:
// - Service binding with closures and concrete instances
// - Singleton services with automatic caching
// - Automatic dependency resolution
// - Thread-safe operations
package container

import (
	"fmt"
	"sync"
)

// Container represents a service container for dependency injection.
// It manages service bindings, singleton instances, and provides
// thread-safe service resolution.
//
// The container follows Laravel's container patterns:
// - Bind: Register a service binding
// - Singleton: Register a singleton service
// - Make: Resolve a service from the container
type Container struct {
	// bindings holds service bindings and singletons
	bindings map[string]interface{}

	// singletonInstances caches singleton instances
	singletonInstances map[string]interface{}

	// mutex provides thread-safe access to container state
	mutex sync.RWMutex
}

// New creates a new service container instance.
//
// Returns:
//
//	*Container: A new container ready for service registration
//
// Example:
//
//	container := container.New()
//	container.Bind("logger", func() interface{} {
//	    return &Logger{}
//	})
func New() *Container {
	return &Container{
		bindings:           make(map[string]interface{}),
		singletonInstances: make(map[string]interface{}),
	}
}

// Bind registers a binding in the service container.
// The binding maps an abstract service name to a concrete implementation.
// Each time the service is resolved, a new instance will be created.
//
// Parameters:
//
//	abstract: The service name/key
//	concrete: The concrete implementation (function, struct, or instance)
//
// Returns:
//
//	error: Any error that occurred during binding
//
// Example:
//
//	container.Bind("database", func() interface{} {
//	    return &DatabaseConnection{}
//	})
//
//	container.Bind("config", &ConfigStruct{})
func (c *Container) Bind(abstract string, concrete interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if abstract == "" {
		return fmt.Errorf("abstract service name cannot be empty")
	}

	c.bindings[abstract] = concrete
	return nil
}

// Singleton registers a singleton binding in the service container.
// Singleton services are instantiated once and cached for subsequent requests.
// This is useful for services that should maintain state across requests.
//
// Parameters:
//
//	abstract: The service name/key
//	concrete: The concrete implementation (function, struct, or instance)
//
// Returns:
//
//	error: Any error that occurred during binding
//
// Example:
//
//	container.Singleton("logger", func() interface{} {
//	    return &Logger{level: "info"}
//	})
func (c *Container) Singleton(abstract string, concrete interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if abstract == "" {
		return fmt.Errorf("abstract service name cannot be empty")
	}

	// Mark as singleton by prefixing the key
	c.bindings["singleton:"+abstract] = concrete
	return nil
}

// Make resolves a service from the container.
// For regular bindings, creates a new instance each time.
// For singletons, returns the cached instance or creates and caches a new one.
//
// Parameters:
//
//	abstract: The service name/key to resolve
//
// Returns:
//
//	interface{}: The resolved service instance
//	error: Any error that occurred during resolution
//
// Example:
//
//	logger, err := container.Make("logger")
//	if err != nil {
//	    return err
//	}
//	log := logger.(*Logger)
func (c *Container) Make(abstract string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if abstract == "" {
		return nil, fmt.Errorf("abstract service name cannot be empty")
	}

	// Check for singleton first
	singletonKey := "singleton:" + abstract
	if concrete, exists := c.bindings[singletonKey]; exists {
		// Check if instance is cached
		if instance, cached := c.singletonInstances[abstract]; cached {
			return instance, nil
		}

		// Create and cache instance
		instance, err := c.resolveService(concrete)
		if err != nil {
			return nil, err
		}

		c.singletonInstances[abstract] = instance
		return instance, nil
	}

	// Check for regular binding
	if concrete, exists := c.bindings[abstract]; exists {
		return c.resolveService(concrete)
	}

	return nil, fmt.Errorf("service '%s' not found in container", abstract)
}

// IsBound checks if a service is registered in the container.
//
// Parameters:
//
//	abstract: The service name/key to check
//
// Returns:
//
//	bool: true if the service is registered, false otherwise
//
// Example:
//
//	if container.IsBound("logger") {
//	    logger, _ := container.Make("logger")
//	}
func (c *Container) IsBound(abstract string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if abstract == "" {
		return false
	}

	// Check for singleton
	if _, exists := c.bindings["singleton:"+abstract]; exists {
		return true
	}

	// Check for regular binding
	_, exists := c.bindings[abstract]
	return exists
}

// IsSingleton checks if a service is registered as a singleton.
//
// Parameters:
//
//	abstract: The service name/key to check
//
// Returns:
//
//	bool: true if the service is a singleton, false otherwise
//
// Example:
//
//	if container.IsSingleton("logger") {
//	    // Logger will be reused across requests
//	}
func (c *Container) IsSingleton(abstract string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if abstract == "" {
		return false
	}

	_, exists := c.bindings["singleton:"+abstract]
	return exists
}

// Forget removes a service binding from the container.
// For singletons, also removes any cached instance.
//
// Parameters:
//
//	abstract: The service name/key to remove
//
// Example:
//
//	container.Forget("logger")
//	// Logger service is no longer available
func (c *Container) Forget(abstract string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if abstract == "" {
		return
	}

	// Remove regular binding
	delete(c.bindings, abstract)

	// Remove singleton binding and cached instance
	delete(c.bindings, "singleton:"+abstract)
	delete(c.singletonInstances, abstract)
}

// Flush removes all service bindings and cached instances.
// This is primarily useful for testing scenarios.
//
// Example:
//
//	// In test cleanup
//	defer container.Flush()
func (c *Container) FlushContainer() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.bindings = make(map[string]interface{})
	c.singletonInstances = make(map[string]interface{})
}

// Count returns the total number of registered services.
// This includes both regular bindings and singletons.
//
// Returns:
//
//	int: Total number of registered services
//
// Example:
//
//	fmt.Printf("Container has %d registered services\n", container.Count())
func (c *Container) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Count unique service names (excluding singleton: prefix)
	services := make(map[string]bool)
	for key := range c.bindings {
		if key[:10] == "singleton:" {
			services[key[10:]] = true
		} else {
			services[key] = true
		}
	}

	return len(services)
}

// RegisteredServices returns a list of all registered service names.
//
// Returns:
//
//	[]string: List of registered service names
//
// Example:
//
//	services := container.RegisteredServices()
//	for _, service := range services {
//	    fmt.Println("Registered service:", service)
//	}
func (c *Container) RegisteredServices() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	services := make(map[string]bool)
	for key := range c.bindings {
		if len(key) > 10 && key[:10] == "singleton:" {
			services[key[10:]] = true
		} else {
			services[key] = true
		}
	}

	var result []string
	for service := range services {
		result = append(result, service)
	}

	return result
}

// resolveService resolves a concrete service implementation.
// Handles both function-based and direct instance bindings.
func (c *Container) resolveService(concrete interface{}) (interface{}, error) {
	// If it's a function, call it to get the instance
	if fn, ok := concrete.(func() interface{}); ok {
		instance := fn()
		if instance == nil {
			return nil, fmt.Errorf("service factory returned nil")
		}
		return instance, nil
	}

	// Return the concrete instance directly
	return concrete, nil
}
