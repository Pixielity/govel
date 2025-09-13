// Package support provides Laravel-style MultipleInstanceManager for managing
// multiple named instances of the same driver type with different configurations.
//
// The MultipleInstanceManager pattern is used in Laravel for services that need
// multiple configured instances, such as multiple database connections, multiple
// filesystem disks, or multiple queue connections.
package support

import (
	"fmt"
	"reflect"
	"sync"

	applicationInterfaces "govel/packages/types/src/interfaces/application/base"
	configInterfaces "govel/packages/types/src/interfaces/config"
	types "govel/packages/types/src/types/support"

	"govel/packages/support/src/str"
	"govel/packages/support/src/traits"
)

// MultipleInstanceManager provides Laravel-style management of multiple named instances.
//
// Unlike the regular Manager which maintains one instance per driver type, the
// MultipleInstanceManager allows multiple instances of the same driver type with
// different configurations. This is essential for services like databases where
// you might need connections to different databases using the same driver.
//
// Key features:
//   - Multiple named instances with individual configurations
//   - Thread-safe instance caching and management
//   - Laravel-style reflection-based method discovery
//   - Runtime extensibility via custom creators
//   - Configurable driver key for flexibility
//
// Example use cases:
//   - Multiple database connections (primary, analytics, logging)
//   - Multiple filesystem disks (local, s3, cdn)
//   - Multiple queue connections (default, high-priority, delayed)
type MultipleInstanceManager struct {
	// app provides access to the application container for dependency injection
	app applicationInterfaces.ApplicationInterface

	// config provides access to application configuration
	config configInterfaces.ConfigInterface

	// instances caches created instances by name to avoid repeated creation
	// Key: instance name, Value: instance object
	instances map[string]interface{}

	// customCreators stores user-defined instance creation functions
	// mapped by driver name for runtime extension
	customCreators map[string]types.InstanceCreator

	// driverKey specifies the configuration key that contains the driver name
	// Default: "driver", but can be customized (e.g., "type", "engine")
	driverKey string

	// mutex provides thread-safe access to instances and customCreators
	mutex sync.RWMutex

	// Proxiable provides self-reference capabilities for proper method resolution
	// when MultipleInstanceManager is embedded in concrete types (e.g., DatabaseManager)
	traits.Proxiable
}

// NewMultipleInstanceManager creates a new MultipleInstanceManager instance.
//
// The application container must provide a "config" binding that implements
// ConfigInterface. This constructor initializes the manager with default settings
// and an empty instance cache.
//
// Note: After creating a manager, concrete types that embed MultipleInstanceManager
// should call SetProxySelf(self) to enable proper method resolution for abstract methods.
//
// Parameters:
//   - app: The application container for dependency injection
//
// Returns:
//   - *MultipleInstanceManager: A new manager ready for instance management
//
// Panics:
//   - If the application cannot resolve the "config" binding
//
// Example:
//
//	app := containerMocks.NewMockContainer()
//	config := configMocks.NewMockConfig()
//	app.Bind("config", config)
//	manager := NewMultipleInstanceManager(app)
func NewMultipleInstanceManager(app applicationInterfaces.ApplicationInterface) *MultipleInstanceManager {
	// Resolve configuration from application container - required for all managers
	config, err := app.Make("config")
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve config from application: %v", err))
	}

	// Safely type-assert the config with nil check
	configInterface, ok := config.(configInterfaces.ConfigInterface)
	if !ok {
		panic(fmt.Sprintf("Config service does not implement ConfigInterface, got %T", config))
	}

	return &MultipleInstanceManager{
		app:            app,
		config:         configInterface,
		instances:      make(map[string]interface{}),
		customCreators: make(map[string]types.InstanceCreator),
		driverKey:      "driver", // Default driver configuration key
	}
}

// GetDefaultInstance returns the name of the default instance.
//
// This method must be implemented by concrete managers to specify which
// instance should be used when no specific instance name is provided.
// Following Laravel's abstract method pattern.
//
// Returns:
//   - string: The name of the default instance
//
// Panics:
//   - Always, as concrete managers must override this method
//
// Example Implementation:
//
//	func (d *DatabaseManager) GetDefaultInstance() string {
//	    return "default"
//	}
func (m *MultipleInstanceManager) GetDefaultInstance() string {
	panic("GetDefaultInstance method must be implemented by concrete manager")
}

// SetDefaultInstance sets the name of the default instance.
//
// This method must be implemented by concrete managers to allow changing
// the default instance at runtime. Following Laravel's abstract method pattern.
//
// Parameters:
//   - name: The name of the instance to set as default
//
// Panics:
//   - Always, as concrete managers must override this method
//
// Example Implementation:
//
//	func (d *DatabaseManager) SetDefaultInstance(name string) {
//	    d.defaultInstance = name
//	}
func (m *MultipleInstanceManager) SetDefaultInstance(name string) {
	panic("SetDefaultInstance method must be implemented by concrete manager")
}

// GetInstanceConfig retrieves the configuration for a specific instance.
//
// This method must be implemented by concrete managers to provide the
// configuration map for instance creation. The configuration should include
// the driver key and any driver-specific settings.
//
// Parameters:
//   - name: The name of the instance to get configuration for
//
// Returns:
//   - map[string]interface{}: The configuration map for the instance
//
// Panics:
//   - Always, as concrete managers must override this method
//
// Example Implementation:
//
//	func (d *DatabaseManager) GetInstanceConfig(name string) map[string]interface{} {
//	    return d.configurations[name]
//	}
func (m *MultipleInstanceManager) GetInstanceConfig(name string) map[string]interface{} {
	panic("GetInstanceConfig method must be implemented by concrete manager")
}

// Instance retrieves a named instance, creating and caching it if necessary.
//
// This is the main entry point for instance access. It follows Laravel's pattern
// where calling Instance() without arguments uses the default instance, while
// Instance('name') uses a specific named instance.
//
// The method implements thread-safe caching to avoid recreating instances on
// subsequent calls. Instance resolution follows this order:
//  1. Check cache for existing instance
//  2. Resolve instance name (provided or default)
//  3. Create new instance via resolve()
//  4. Cache and return the instance
//
// Parameters:
//   - name: Optional instance name. If empty, uses GetDefaultInstance()
//
// Returns:
//   - interface{}: The instance (should be cast to expected type)
//   - error: Any error that occurred during resolution or creation
//
// Thread Safety:
//
//	This method is fully thread-safe using RWMutex for cache access.
//
// Example:
//
//	// Get default instance
//	instance, err := manager.Instance()
//
//	// Get specific instance
//	instance, err := manager.Instance("analytics")
func (m *MultipleInstanceManager) Instance(name ...string) (interface{}, error) {
	var instanceName string

	// Resolve the instance name - either provided or default
	if len(name) > 0 && name[0] != "" {
		// Use explicitly provided instance name
		instanceName = name[0]
	} else {
		// Use default instance from concrete manager implementation via proxy
		if m.HasProxySelf() {
			results, err := m.CallOnSelf("GetDefaultInstance")
			if err == nil && len(results) > 0 {
				instanceName = results[0].String()
			} else {
				return nil, fmt.Errorf("GetDefaultInstance method must be implemented by concrete manager: %v", err)
			}
		} else {
			return nil, fmt.Errorf("GetDefaultInstance method must be implemented by concrete manager")
		}
	}

	// Get the instance (from cache or create new) - get() handles its own caching
	return m.get(instanceName)
}

// get attempts to retrieve an instance from cache or creates it via resolve.
//
// This internal method implements the caching layer for instance management.
// It first checks the instance cache, and if not found, delegates to resolve()
// to create a new instance.
//
// Parameters:
//   - name: The name of the instance to retrieve
//
// Returns:
//   - interface{}: The instance from cache or newly created
//   - error: Any error from instance resolution
//
// Thread Safety:
//
//	Uses RWMutex for safe concurrent access to the instance cache.
func (m *MultipleInstanceManager) get(name string) (interface{}, error) {
	// Check cache first for existing instance
	m.mutex.RLock()
	if instance, exists := m.instances[name]; exists {
		// Cache hit - return existing instance
		m.mutex.RUnlock()
		return instance, nil
	}
	m.mutex.RUnlock()

	// Cache miss - create new instance
	instance, err := m.resolve(name)
	if err != nil {
		return nil, err
	}

	// Cache the newly created instance (with double-checked locking pattern)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double-check if instance was created by another goroutine while we were waiting
	if existingInstance, exists := m.instances[name]; exists {
		return existingInstance, nil
	}

	// Cache the instance
	m.instances[name] = instance
	return instance, nil
}

// resolve creates a new instance using Laravel-style reflection or custom creators.
//
// This method implements the core instance creation logic, following Laravel's
// approach of using reflection to discover and call Create{Driver}{DriverKey}
// methods on the manager class.
//
// Instance creation priority:
//  1. Custom creators registered via Extend() (highest priority)
//  2. Reflection-based method discovery for Create{Driver}{DriverKey} methods
//  3. Return "not supported" error if neither found
//
// Method Discovery Process:
//   - Extract driver name from instance configuration
//   - Try multiple naming conventions (Title and Studly case)
//   - Use reflection to find and call the method with configuration
//   - Validate method signature: func(map[string]interface{}) (interface{}, error)
//
// Parameters:
//   - name: The name of the instance to create
//
// Returns:
//   - interface{}: The created instance
//   - error: Any error that occurred during creation
//
// Example Method Discovery:
//
//	instance="analytics", driver="mysql" -> method="CreateMysqlDriver"
//	instance="cache", driver="redis_cluster" -> method="CreateRedisClusterDriver"
func (m *MultipleInstanceManager) resolve(name string) (interface{}, error) {
	// Get configuration for the named instance via proxy
	var config map[string]interface{}
	if m.HasProxySelf() {
		results, err := m.CallOnSelf("GetInstanceConfig", name)
		if err == nil && len(results) > 0 {
			if configVal := results[0].Interface(); configVal != nil {
				config = configVal.(map[string]interface{})
			}
		} else {
			return nil, fmt.Errorf("GetInstanceConfig method must be implemented by concrete manager")
		}
	} else {
		return nil, fmt.Errorf("GetInstanceConfig method must be implemented by concrete manager")
	}
	if config == nil {
		return nil, fmt.Errorf("instance [%s] is not defined", name)
	}

	// Extract driver name from configuration (thread-safe access to driverKey)
	m.mutex.RLock()
	currentDriverKey := m.driverKey
	m.mutex.RUnlock()

	driverName, exists := config[currentDriverKey]
	if !exists {
		return nil, fmt.Errorf("instance [%s] does not specify a %s", name, currentDriverKey)
	}

	// Validate driver name is a string
	driverNameStr, ok := driverName.(string)
	if !ok {
		return nil, fmt.Errorf("instance [%s] %s must be a string", name, currentDriverKey)
	}

	// Priority 1: Check for custom creator (registered via Extend)
	m.mutex.RLock()
	creator, exists := m.customCreators[driverNameStr]
	m.mutex.RUnlock()

	if exists {
		return m.callCustomCreator(config, creator)
	}

	// Priority 2: Use Laravel-style reflection to find Create{Driver}{DriverKey} methods
	// Try different method naming conventions for compatibility
	createMethods := []string{
		fmt.Sprintf("Create%s%s", str.Title(driverNameStr), str.Title(currentDriverKey)),
		fmt.Sprintf("Create%s%s", str.Studly(driverNameStr), str.Title(currentDriverKey)),
	}

	// Use reflection to find and call the appropriate method on the concrete manager
	var managerValue reflect.Value
	if m.HasProxySelf() {
		managerValue = reflect.ValueOf(m.GetProxySelf())
	} else {
		managerValue = reflect.ValueOf(m)
	}
	for _, methodName := range createMethods {
		method := managerValue.MethodByName(methodName)
		if method.IsValid() {
			// Method found - call it with the configuration
			configValue := reflect.ValueOf(config)
			results := method.Call([]reflect.Value{configValue})

			// Validate method signature: must return (interface{}, error)
			if len(results) != 2 {
				return nil, fmt.Errorf("create method %s must return (interface{}, error)", methodName)
			}

			// Validate that the second return value is an error type
			if results[1].Type() != reflect.TypeOf((*error)(nil)).Elem() {
				return nil, fmt.Errorf("create method %s second return value must be error type, got %s", methodName, results[1].Type())
			}

			// Check if the method returned an error
			if !results[1].IsNil() {
				if err, ok := results[1].Interface().(error); ok {
					return nil, err
				}
				return nil, fmt.Errorf("create method %s returned non-error as error: %v", methodName, results[1].Interface())
			}

			// Return the created instance
			return results[0].Interface(), nil
		}
	}

	// No custom creator or Create{Driver}{DriverKey} method found
	return nil, fmt.Errorf("instance %s [%s] is not supported", currentDriverKey, driverNameStr)
}

// callCustomCreator executes a custom instance creation function.
//
// Custom creators are user-defined functions that can create instance objects
// without following the reflection-based Create{Driver}{DriverKey} method pattern.
// This provides flexibility for runtime instance registration and dependency injection.
//
// Parameters:
//   - config: The configuration map for the instance
//   - creator: The function that creates the instance
//
// Returns:
//   - interface{}: The created instance
//   - error: Any error returned by the creator function
func (m *MultipleInstanceManager) callCustomCreator(config map[string]interface{}, creator types.InstanceCreator) (interface{}, error) {
	return creator(m.app, config)
}

// ForgetInstance removes specified instances from the cache, forcing recreation on next access.
//
// This method allows selective removal of cached instances. If no names are provided,
// it removes the default instance. Multiple instance names can be provided to remove
// several instances at once.
//
// Parameters:
//   - names: Optional instance names to remove. If empty, removes default instance
//
// Returns:
//   - *MultipleInstanceManager: The manager instance for method chaining
//
// Thread Safety:
//
//	This method is thread-safe using mutex protection.
//
// Example:
//
//	// Remove default instance
//	manager.ForgetInstance()
//
//	// Remove specific instances
//	manager.ForgetInstance("analytics", "logging")
func (m *MultipleInstanceManager) ForgetInstance(names ...string) *MultipleInstanceManager {
	// Determine which instances to remove
	var instanceNames []string
	if len(names) == 0 {
		// No names provided - remove default instance
		if m.HasProxySelf() {
			results, err := m.CallOnSelf("GetDefaultInstance")
			if err == nil && len(results) > 0 {
				instanceNames = []string{results[0].String()}
			} else {
				return m // Cannot get default instance, return unchanged
			}
		} else {
			return m // No proxy set, return unchanged
		}
	} else {
		// Use provided instance names
		instanceNames = names
	}

	// Remove instances from cache
	m.mutex.Lock()
	for _, instanceName := range instanceNames {
		delete(m.instances, instanceName)
	}
	m.mutex.Unlock()

	return m
}

// Purge disconnects and removes a single instance from the local cache.
//
// This method removes a specific instance from the cache, forcing recreation
// on the next access. If no name is provided, it purges the default instance.
//
// Parameters:
//   - name: Optional instance name to purge. If empty, purges default instance
//
// Thread Safety:
//
//	This method is thread-safe using mutex protection.
//
// Example:
//
//	// Purge default instance
//	manager.Purge()
//
//	// Purge specific instance
//	manager.Purge("analytics")
func (m *MultipleInstanceManager) Purge(name ...string) {
	// Determine which instance to purge
	var instanceName string
	if len(name) > 0 && name[0] != "" {
		// Use provided instance name
		instanceName = name[0]
	} else {
		// Use default instance via proxy
		if m.HasProxySelf() {
			results, err := m.CallOnSelf("GetDefaultInstance")
			if err == nil && len(results) > 0 {
				instanceName = results[0].String()
			} else {
				return // Cannot get default instance, do nothing
			}
		} else {
			return // No proxy set, do nothing
		}
	}

	// Remove instance from cache
	m.mutex.Lock()
	delete(m.instances, instanceName)
	m.mutex.Unlock()
}

// Extend registers a custom instance creator function for runtime extension.
//
// This method allows users to add new drivers without modifying the manager class,
// following Laravel's pattern of runtime extensibility. Custom creators take
// precedence over reflection-based method discovery.
//
// The creator function receives the application container and configuration,
// allowing for flexible instance creation with access to dependencies.
//
// Parameters:
//   - name: The name of the driver to register
//   - creator: The function that will create instances of this driver
//
// Returns:
//   - *MultipleInstanceManager: The manager instance for method chaining
//
// Thread Safety:
//
//	This method is thread-safe using mutex protection.
//
// Example:
//
//	manager.Extend("custom", func(app Application, config map[string]interface{}) (interface{}, error) {
//	    return &CustomInstance{Config: config}, nil
//	})
func (m *MultipleInstanceManager) Extend(name string, creator types.InstanceCreator) *MultipleInstanceManager {
	m.mutex.Lock()
	m.customCreators[name] = creator
	m.mutex.Unlock()
	return m
}

// SetApplication updates the application container used by this manager.
//
// When a new application is set, the method automatically attempts to resolve
// a new "config" instance from the application to keep the configuration
// reference in sync.
//
// Parameters:
//   - app: The new application instance to use
//
// Returns:
//   - *MultipleInstanceManager: The manager instance for method chaining
//
// Note:
//
//	If the new application cannot provide a "config" binding, the config
//	reference will remain unchanged (no error is thrown).
func (m *MultipleInstanceManager) SetApplication(app applicationInterfaces.ApplicationInterface) *MultipleInstanceManager {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.app = app

	// Attempt to update config reference from the new application
	if config, err := app.Make("config"); err == nil {
		// Safely type-assert the config
		if configInterface, ok := config.(configInterfaces.ConfigInterface); ok {
			m.config = configInterface
		}
	}

	return m
}

// GetApplication returns the application container used by this manager.
//
// The application container is used for resolving dependencies during instance
// creation and accessing application services.
//
// Returns:
//   - Application: The current application instance
func (m *MultipleInstanceManager) GetApplication() applicationInterfaces.ApplicationInterface {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.app
}

// GetInstances returns a copy of all currently cached instances.
//
// This method provides access to the internal instance cache for inspection
// purposes. The returned map is a copy to prevent external modification
// of the internal cache state.
//
// Returns:
//   - map[string]interface{}: A copy of the instance cache (name -> instance)
//
// Thread Safety:
//
//	This method is thread-safe using RWMutex for read access.
func (m *MultipleInstanceManager) GetInstances() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification of the internal cache
	instances := make(map[string]interface{})
	for k, v := range m.instances {
		instances[k] = v
	}
	return instances
}

// SetDriverKey sets the configuration key name used to identify the driver type.
//
// By default, the manager looks for a "driver" key in instance configurations.
// This method allows customization of that key name for different naming conventions.
//
// Parameters:
//   - key: The configuration key name (e.g., "driver", "type", "engine")
//
// Returns:
//   - *MultipleInstanceManager: The manager instance for method chaining
//
// Example:
//
//	// Use "type" instead of "driver"
//	manager.SetDriverKey("type")
//	// Now configurations should have: {"type": "mysql", ...}
func (m *MultipleInstanceManager) SetDriverKey(key string) *MultipleInstanceManager {
	m.mutex.Lock()
	m.driverKey = key
	m.mutex.Unlock()
	return m
}

// GetDriverKey returns the current configuration key name used for driver identification.
//
// Returns:
//   - string: The current driver key (default: "driver")
func (m *MultipleInstanceManager) GetDriverKey() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.driverKey
}

// GetConfig returns the configuration interface used by this manager.
//
// The configuration interface provides access to application configuration
// values and is used throughout the manager for instance creation and management.
//
// Returns:
//   - configInterfaces.ConfigInterface: The current configuration instance
func (m *MultipleInstanceManager) GetConfig() configInterfaces.ConfigInterface {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.config
}
