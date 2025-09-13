// Package support provides Laravel-style Manager pattern implementation for Go.
//
// The Manager pattern is a core architectural pattern in Laravel used for services
// like Cache, Queue, Mail, Database, and others where multiple "drivers" or
// implementations can be dynamically created and managed.
//
// Key Features:
//   - Laravel-compatible driver discovery via reflection
//   - Thread-safe driver caching and management
//   - Runtime extensibility via custom creator functions
//   - Dependency injection support through containers
//   - Method naming conventions: Create{Driver}Driver
//
// Usage Example:
//
//	// Create a concrete manager that embeds Manager
//	type HashManager struct {
//	    *Manager
//	    defaultDriver string
//	}
//
//	// Implement GetDefaultDriver (required by ManagerInterface)
//	func (h *HashManager) GetDefaultDriver() string {
//	    return h.defaultDriver
//	}
//
//	// Add driver creation methods using Laravel naming convention
//	func (h *HashManager) CreateBcryptDriver() (interface{}, error) {
//	    return &BcryptHashDriver{cost: 10}, nil
//	}
//
//	// Initialize and use
//	hashManager := &HashManager{Manager: NewManager(container), defaultDriver: "bcrypt"}
//	hashManager.SetSelf(hashManager) // Required for embedded structs
//	driver, _ := hashManager.Driver() // Gets bcrypt driver
package support

import (
	"fmt"
	"reflect"
	"sync"

	configInterfaces "govel/types/src/interfaces/config"
	containerInterfaces "govel/types/src/interfaces/container"

	"govel/support/src/str"
	"govel/support/src/traits"
	"govel/types/src/types/support"
)

// Manager provides the base implementation for Laravel-style driver management.
// It implements the Manager pattern commonly used in Laravel for services like
// Cache, Queue, Mail, etc. where multiple "drivers" (implementations) can be
// dynamically created and cached.
//
// The Manager uses reflection to discover and call Create{Driver}Driver methods
// on concrete manager implementations, mirroring Laravel's approach of using
// PHP's method_exists() and call_user_func() for dynamic method invocation.
//
// Key features:
//   - Thread-safe driver caching with sync.RWMutex
//   - Laravel-style method discovery via reflection
//   - Support for custom driver creators via Extend()
//   - Embedded struct pattern support via self-reference
type Manager struct {
	// container provides dependency injection capabilities
	container containerInterfaces.ContainerInterface

	// config provides access to application configuration
	config configInterfaces.ConfigInterface

	// customCreators stores user-defined driver creation functions
	// mapped by driver name for runtime extension
	customCreators map[string]types.DriverCreator

	// drivers caches created driver instances to avoid repeated instantiation
	// Key: driver name, Value: driver instance
	drivers map[string]interface{}

	// mutex provides thread-safe access to the drivers cache and customCreators
	mutex sync.RWMutex

	// Proxiable provides self-reference capabilities for proper method resolution
	// when Manager is embedded in concrete types (e.g., HashManager)
	traits.Proxiable
}

// NewManager creates a new Manager instance with the provided container.
//
// The container must provide a "config" binding that implements ConfigInterface.
// This follows Laravel's pattern where managers receive the application container
// and resolve dependencies from it.
//
// Parameters:
//   - container: The dependency injection container
//
// Returns:
//   - *Manager: A new Manager instance ready for driver management
//
// Panics:
//   - If the container cannot resolve the "config" binding
//
// Example:
//
//	container := containerMocks.NewMockContainer()
//	config := configMocks.NewMockConfig()
//	container.Bind("config", config)
//	manager := NewManager(container)
func NewManager(container containerInterfaces.ContainerInterface) *Manager {
	// Resolve configuration from container - required for all managers
	config, err := container.Make("config")
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve config from container: %v", err))
	}

	return &Manager{
		container:      container,
		config:         config.(configInterfaces.ConfigInterface),
		customCreators: make(map[string]types.DriverCreator),
		drivers:        make(map[string]interface{}),
		// Proxiable will be initialized by concrete managers via SetProxySelf()
	}
}

// GetDefaultDriver provides a base implementation that returns an empty string.
//
// Concrete manager implementations should override this method to return their
// actual default driver name. This method implements the ManagerInterface
// requirement and follows Laravel's abstract method pattern.
//
// Returns:
//   - string: Empty string (concrete managers should override this)
//
// Example Override:
//
//	func (h *HashManager) GetDefaultDriver() string {
//	    return "bcrypt"
//	}
func (m *Manager) GetDefaultDriver() string {
	panic("GetDefaultDriver method must be implemented by concrete manager")
}

// Driver retrieves a driver instance by name, creating and caching it if necessary.
//
// This is the main entry point for driver access, following Laravel's pattern where
// calling driver() without arguments uses the default driver, while driver('name')
// uses a specific driver.
//
// The method implements a thread-safe caching mechanism to avoid recreating drivers
// on subsequent calls. Driver resolution follows this order:
//  1. Check cache for existing instance
//  2. Resolve driver name (provided or default)
//  3. Create new instance via createDriver()
//  4. Cache and return the instance
//
// Parameters:
//   - driver: Optional driver name. If empty, uses GetDefaultDriver()
//
// Returns:
//   - interface{}: The driver instance (should be cast to expected type)
//   - error: Any error that occurred during driver resolution or creation
//
// Thread Safety:
//
//	This method is fully thread-safe using RWMutex for cache access.
//
// Example:
//
//	// Get default driver
//	driver, err := manager.Driver()
//
//	// Get specific driver
//	driver, err := manager.Driver("redis")
func (m *Manager) Driver(driver ...string) (interface{}, error) {
	var driverName string

	// Resolve the driver name - either provided or default
	if len(driver) > 0 && driver[0] != "" {
		// Use explicitly provided driver name
		driverName = driver[0]
	} else {
		// Use default driver from concrete manager implementation via proxy
		if m.HasProxySelf() {
			// Use proxy to call GetDefaultDriver on concrete implementation
			results, err := m.CallOnSelf("GetDefaultDriver")
			if err == nil && len(results) > 0 {
				driverName = results[0].String()
			} else {
				// Fallback if proxy call fails
				driverName = m.GetDefaultDriver()
			}
		} else {
			// Fallback to base implementation (usually returns empty string)
			driverName = m.GetDefaultDriver()
		}
	}

	// Validate that we have a driver name to work with
	if driverName == "" {
		return nil, fmt.Errorf("unable to resolve empty driver for [%s]", reflect.TypeOf(m).String())
	}

	// Check cache first for existing driver instance
	m.mutex.RLock()
	if instance, exists := m.drivers[driverName]; exists {
		// Cache hit - return existing instance
		m.mutex.RUnlock()
		return instance, nil
	}
	m.mutex.RUnlock()

	// Cache miss - create new driver instance
	instance, err := m.createDriver(driverName)
	if err != nil {
		return nil, err
	}

	// Cache the newly created instance for future use
	m.mutex.Lock()
	m.drivers[driverName] = instance
	m.mutex.Unlock()

	return instance, nil
}

// createDriver creates a new driver instance using Laravel-style reflection or custom creators.
//
// This method implements the core Laravel Manager pattern logic for driver instantiation.
// It follows Laravel's approach of using reflection to discover and call Create{Driver}Driver
// methods on the manager class.
//
// Driver creation priority:
//  1. Custom creators registered via Extend() (highest priority)
//  2. Reflection-based method discovery for Create{Driver}Driver methods
//  3. Return "not supported" error if neither found
//
// Method Discovery Process:
//   - Convert driver name to StudlyCase (e.g., "my_driver" -> "MyDriver")
//   - Build method name: "Create" + StudlyCase + "Driver"
//   - Use reflection to find and call the method on the target manager
//   - Validate method signature: func() (interface{}, error)
//
// Parameters:
//   - driver: The name of the driver to create (e.g., "redis", "database")
//
// Returns:
//   - interface{}: The created driver instance
//   - error: Any error that occurred during creation
//
// Example Method Discovery:
//
//	driver="bcrypt" -> method="CreateBcryptDriver"
//	driver="my_custom" -> method="CreateMyCustomDriver"
func (m *Manager) createDriver(driver string) (interface{}, error) {
	// Priority 1: Check for custom driver creators (registered via Extend)
	if creator, exists := m.customCreators[driver]; exists {
		return m.callCustomCreator(creator)
	}

	// Priority 2: Use Laravel-style reflection to find Create{Driver}Driver methods
	// Convert driver name to StudlyCase and build method name
	studlyDriver := str.Studly(driver)
	methodName := "Create" + studlyDriver + "Driver"

	// Determine the target for reflection - prefer proxy self-reference for embedded structs
	var managerValue reflect.Value
	if m.HasProxySelf() {
		// Use proxy self-reference to ensure we see methods on the concrete manager type
		// (e.g., HashManager) rather than just the base Manager type
		managerValue = reflect.ValueOf(m.GetProxySelf())
	} else {
		// Fallback to the current manager instance
		managerValue = reflect.ValueOf(m)
	}

	// Attempt to find the method using reflection
	method := managerValue.MethodByName(methodName)

	if method.IsValid() {
		// Method found - call it and validate the response
		results := method.Call([]reflect.Value{})

		// Validate method signature: must return (interface{}, error)
		if len(results) != 2 {
			return nil, fmt.Errorf("create method %s must return (interface{}, error)", methodName)
		}

		// Check if the method returned an error
		errorResult := results[1]
		if !errorResult.IsNil() {
			return nil, errorResult.Interface().(error)
		}

		// Return the created driver instance
		return results[0].Interface(), nil
	}

	// No custom creator or Create{Driver}Driver method found
	return nil, fmt.Errorf("driver [%s] not supported", driver)
}

// callCustomCreator executes a custom driver creation function.
//
// Custom creators are user-defined functions that can create driver instances
// without following the reflection-based Create{Driver}Driver method pattern.
// This provides flexibility for runtime driver registration and dependency injection.
//
// Parameters:
//   - creator: The function that creates the driver instance
//
// Returns:
//   - interface{}: The created driver instance
//   - error: Any error returned by the creator function
func (m *Manager) callCustomCreator(creator types.DriverCreator) (interface{}, error) {
	return creator(m.container)
}

// Extend registers a custom driver creator function for runtime driver extension.
//
// This method allows users to add new drivers without modifying the manager class,
// following Laravel's pattern of runtime extensibility. Custom creators take
// precedence over reflection-based method discovery.
//
// The creator function receives the container and should return a driver instance
// and any error that occurred during creation.
//
// Parameters:
//   - driver: The name of the driver to register
//   - creator: The function that will create instances of this driver
//
// Returns:
//   - *Manager: The manager instance for method chaining
//
// Thread Safety:
//
//	This method is thread-safe using mutex protection.
//
// Example:
//
//	manager.Extend("custom", func(container containerInterfaces.ContainerInterface) (interface{}, error) {
//	    return &CustomDriver{}, nil
//	})
func (m *Manager) Extend(driver string, creator types.DriverCreator) *Manager {
	m.mutex.Lock()
	m.customCreators[driver] = creator
	m.mutex.Unlock()
	return m
}

// GetDrivers returns a copy of all currently cached driver instances.
//
// This method provides access to the internal driver cache for inspection
// purposes. The returned map is a copy to prevent external modification
// of the internal cache state.
//
// Returns:
//   - map[string]interface{}: A copy of the driver cache (driver name -> instance)
//
// Thread Safety:
//
//	This method is thread-safe using RWMutex for read access.
func (m *Manager) GetDrivers() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification of the internal cache
	drivers := make(map[string]interface{})
	for k, v := range m.drivers {
		drivers[k] = v
	}
	return drivers
}

// GetContainer returns the dependency injection container used by this manager.
//
// The container is used for resolving dependencies during driver creation
// and accessing application services.
//
// Returns:
//   - containerInterfaces.ContainerInterface: The current container instance
func (m *Manager) GetContainer() containerInterfaces.ContainerInterface {
	return m.container
}

// SetContainer updates the dependency injection container used by this manager.
//
// When a new container is set, the method automatically attempts to resolve
// a new "config" instance from the container to keep the configuration
// reference in sync.
//
// Parameters:
//   - container: The new container instance to use
//
// Returns:
//   - *Manager: The manager instance for method chaining
//
// Note:
//
//	If the new container cannot provide a "config" binding, the config
//	reference will remain unchanged (no error is thrown).
func (m *Manager) SetContainer(container containerInterfaces.ContainerInterface) *Manager {
	m.container = container

	// Attempt to update config reference from the new container
	if config, err := container.Make("config"); err == nil {
		m.config = config.(configInterfaces.ConfigInterface)
	}

	return m
}

// ForgetDrivers clears all cached driver instances, forcing recreation on next access.
//
// This method is useful for testing or when you need to ensure fresh driver
// instances are created (e.g., after configuration changes). All cached drivers
// are removed from the internal cache.
//
// Returns:
//   - *Manager: The manager instance for method chaining
//
// Thread Safety:
//
//	This method is thread-safe using mutex protection.
//
// Note:
//
//	This only clears the cache; it does not affect custom creators registered
//	via Extend().
func (m *Manager) ForgetDrivers() *Manager {
	m.mutex.Lock()
	m.drivers = make(map[string]interface{})
	m.mutex.Unlock()
	return m
}

// GetConfig returns the configuration interface used by this manager.
//
// The configuration interface provides access to application configuration
// values and is used throughout the manager for driver creation and management.
//
// Returns:
//   - configInterfaces.ConfigInterface: The current configuration instance
func (m *Manager) GetConfig() configInterfaces.ConfigInterface {
	return m.config
}
