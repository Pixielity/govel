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
	"govel/types/src/types"
	"sync"
)

// ServiceContainer represents a service container for dependency injection.
// It manages service bindings, singleton instances, and provides
// thread-safe service resolution with statistics tracking.
//
// The container follows Laravel's container patterns:
// - Bind: Register a service binding
// - Singleton: Register a singleton service
// - Make: Resolve a service from the container
// - GetBindings: Introspect container bindings
// - GetStatistics: Monitor container usage
type ServiceContainer struct {
	// bindings holds service bindings and singletons
	bindings map[string]interface{}

	// singletonInstances caches singleton instances
	singletonInstances map[string]interface{}

	// resolutionCount tracks how many times each service has been resolved
	resolutionCount map[string]int

	// totalResolutions tracks the total number of service resolutions
	totalResolutions int

	// mutex provides thread-safe access to container state
	mutex sync.RWMutex
}

// New creates a new service container instance.
//
// Returns:
//
//	*ServiceContainer: A new container ready for service registration
//
// Example:
//
//	container := container.New()
//	container.Bind("logger", func() interface{} {
//	    return &Logger{}
//	})
func New() *ServiceContainer {
	return &ServiceContainer{
		bindings:           make(map[string]interface{}),
		singletonInstances: make(map[string]interface{}),
		resolutionCount:    make(map[string]int),
		totalResolutions:   0,
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
func (c *ServiceContainer) Bind(abstract types.ServiceIdentifier, concrete interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := types.ToKey(abstract)
	if key == "" {
		return fmt.Errorf("abstract service name cannot be empty")
	}

	c.bindings[key] = concrete
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
func (c *ServiceContainer) Singleton(abstract types.ServiceIdentifier, concrete interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := types.ToKey(abstract)
	if key == "" {
		return fmt.Errorf("abstract service name cannot be empty")
	}

	// Mark as singleton by prefixing the key
	c.bindings["singleton:"+key] = concrete
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
func (c *ServiceContainer) Make(abstract types.ServiceIdentifier) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := types.ToKey(abstract)
	if key == "" {
		return nil, fmt.Errorf("abstract service name cannot be empty")
	}

	// Check for singleton first
	singletonKey := "singleton:" + key
	if concrete, exists := c.bindings[singletonKey]; exists {
		// Check if instance is cached
		if instance, cached := c.singletonInstances[key]; cached {
			// Track resolution
			c.resolutionCount[key]++
			c.totalResolutions++
			return instance, nil
		}

		// Create and cache instance
		instance, err := c.resolveService(concrete)
		if err != nil {
			return nil, err
		}

		c.singletonInstances[key] = instance
		// Track resolution
		c.resolutionCount[key]++
		c.totalResolutions++
		return instance, nil
	}

	// Check for regular binding
	if concrete, exists := c.bindings[key]; exists {
		instance, err := c.resolveService(concrete)
		if err != nil {
			return nil, err
		}
		// Track resolution
		c.resolutionCount[key]++
		c.totalResolutions++
		return instance, nil
	}

	return nil, fmt.Errorf("service '%s' not found in container", key)
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
func (c *ServiceContainer) IsBound(abstract types.ServiceIdentifier) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key := types.ToKey(abstract)
	if key == "" {
		return false
	}

	// Check for singleton
	if _, exists := c.bindings["singleton:"+key]; exists {
		return true
	}

	// Check for regular binding
	_, exists := c.bindings[key]
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
func (c *ServiceContainer) IsSingleton(abstract types.ServiceIdentifier) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key := types.ToKey(abstract)
	if key == "" {
		return false
	}

	_, exists := c.bindings["singleton:"+key]
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
func (c *ServiceContainer) Forget(abstract types.ServiceIdentifier) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := types.ToKey(abstract)
	if key == "" {
		return
	}

	// Remove regular binding
	delete(c.bindings, key)

	// Remove singleton binding and cached instance
	delete(c.bindings, "singleton:"+key)
	delete(c.singletonInstances, key)
}

// Flush removes all service bindings and cached instances.
// This is primarily useful for testing scenarios.
//
// Example:
//
//	// In test cleanup
//	defer container.Flush()
func (c *ServiceContainer) FlushContainer() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.bindings = make(map[string]interface{})
	c.singletonInstances = make(map[string]interface{})
	c.resolutionCount = make(map[string]int)
	c.totalResolutions = 0
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
func (c *ServiceContainer) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Count unique service names (excluding singleton: prefix)
	services := make(map[string]bool)
	for key := range c.bindings {
		if len(key) > 10 && key[:10] == "singleton:" {
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
func (c *ServiceContainer) RegisteredServices() []string {
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

// GetBindings returns detailed information about all service bindings in the container.
// This method provides introspection capabilities for debugging and monitoring purposes.
//
// Returns:
//
//	map[string]interface{}: Map of service names to their binding information
//
// Example:
//
//	bindings := container.GetBindings()
//	for serviceName, info := range bindings {
//	    fmt.Printf("Service '%s': %+v\n", serviceName, info)
//	}
func (c *ServiceContainer) GetBindings() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	bindings := make(map[string]interface{})
	services := make(map[string]bool) // Track unique service names

	// Process all bindings
	for key, concrete := range c.bindings {
		var serviceName string
		var bindingType string

		if len(key) > 10 && key[:10] == "singleton:" {
			serviceName = key[10:]
			bindingType = "singleton"
		} else {
			serviceName = key
			bindingType = "regular"
		}

		// Skip if we've already processed this service
		if services[serviceName] {
			continue
		}
		services[serviceName] = true

		// Determine concrete type
		concreteType := "unknown"
		if concrete != nil {
			switch concrete.(type) {
			case func() interface{}:
				concreteType = "function"
			default:
				concreteType = fmt.Sprintf("%T", concrete)
			}
		}

		// Check if singleton is cached
		cached := false
		if bindingType == "singleton" {
			_, cached = c.singletonInstances[serviceName]
		}

		// Get resolution count
		resolvedCount := c.resolutionCount[serviceName]

		bindings[serviceName] = map[string]interface{}{
			"type":           bindingType,
			"concrete":       concreteType,
			"cached":         cached,
			"resolved_count": resolvedCount,
		}
	}

	return bindings
}

// GetStatistics returns container usage statistics and performance metrics.
// This method provides monitoring data about container performance and usage patterns.
//
// Returns:
//
//	map[string]interface{}: Container statistics and metrics
//
// Example:
//
//	stats := container.GetStatistics()
//	fmt.Printf("Total bindings: %d\n", stats["total_bindings"])
//	fmt.Printf("Cached singletons: %d\n", stats["cached_singletons"])
func (c *ServiceContainer) GetStatistics() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Count different types of bindings
	var singletonBindings, regularBindings int
	services := make(map[string]bool)

	for key := range c.bindings {
		if len(key) > 10 && key[:10] == "singleton:" {
			serviceName := key[10:]
			if !services[serviceName] {
				singletonBindings++
				services[serviceName] = true
			}
		} else {
			if !services[key] {
				regularBindings++
				services[key] = true
			}
		}
	}

	// Count cached singletons
	cachedSingletons := len(c.singletonInstances)

	// Find most resolved services
	mostResolved := c.getMostResolvedServices(5) // Top 5

	return map[string]interface{}{
		"total_bindings":     singletonBindings + regularBindings,
		"singleton_bindings": singletonBindings,
		"regular_bindings":   regularBindings,
		"cached_singletons":  cachedSingletons,
		"total_resolutions":  c.totalResolutions,
		"most_resolved":      mostResolved,
		"memory_usage":       "tracking not implemented", // Could be implemented with runtime.ReadMemStats
	}
}

// getMostResolvedServices returns the most frequently resolved services.
// This is a helper method for GetStatistics.
//
// Parameters:
//
//	limit: Maximum number of services to return
//
// Returns:
//
//	[]map[string]interface{}: List of services with their resolution counts
func (c *ServiceContainer) getMostResolvedServices(limit int) []map[string]interface{} {
	type serviceCount struct {
		name  string
		count int
	}

	// Convert resolution count map to slice for sorting
	var services []serviceCount
	for name, count := range c.resolutionCount {
		services = append(services, serviceCount{name: name, count: count})
	}

	// Simple bubble sort (could use sort.Slice for better performance)
	for i := 0; i < len(services)-1; i++ {
		for j := 0; j < len(services)-i-1; j++ {
			if services[j].count < services[j+1].count {
				services[j], services[j+1] = services[j+1], services[j]
			}
		}
	}

	// Limit results
	if len(services) > limit {
		services = services[:limit]
	}

	// Convert to map format
	result := make([]map[string]interface{}, len(services))
	for i, service := range services {
		result[i] = map[string]interface{}{
			"name":  service.name,
			"count": service.count,
		}
	}

	return result
}

// resolveService resolves a concrete service implementation.
// Handles both function-based and direct instance bindings.
func (c *ServiceContainer) resolveService(concrete interface{}) (interface{}, error) {
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

