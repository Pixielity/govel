package traits

import (
	"sync"

	"govel/packages/container"
	"govel/packages/container/interfaces"
)

/**
 * Containable provides container management functionality in a thread-safe manner.
 * This trait follows the self-contained pattern and manages its own container instance.
 *
 * The embedded container field automatically promotes all ContainerInterface methods,
 * so we don't need to implement them manually - Go's embedding handles this.
 */
type Containable struct {
	/**
	 * mutex provides thread-safe access to container properties
	 */
	mutex sync.RWMutex

	/**
	 * container holds the container instance and is embedded to promote its methods
	 */
	*container.Container
}

/**
 * NewContainable creates a new Containable instance with a container.
 *
 * @param containerInstance *container.Container The container instance to use
 * @return *Containable The newly created trait instance
 */
func NewContainable(containerInstance *container.Container) *Containable {
	if containerInstance == nil {
		containerInstance = container.New()
	}

	return &Containable{
		Container: containerInstance,
	}
}

/**
 * GetContainer returns the container instance.
 *
 * @return interfaces.ContainerInterface The container instance
 */
func (t *Containable) GetContainer() interfaces.ContainerInterface {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Container
}

/**
 * SetContainer sets the container instance.
 *
 * @param containerInterface interface{} The container instance to set
 */
func (t *Containable) SetContainer(containerInterface interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Type assertion to ensure it's a container
	if containerInstance, ok := containerInterface.(*container.Container); ok {
		t.Container = containerInstance
	} else if containerInterface == nil {
		t.Container = container.New()
	}
}

/**
 * HasContainer returns whether a container instance is set.
 *
 * @return bool true if a container is set
 */
func (t *Containable) HasContainer() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Container != nil
}

/**
 * GetContainerInfo returns information about the container.
 *
 * @return map[string]interface{} Container information
 */
func (t *Containable) GetContainerInfo() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	info := map[string]interface{}{
		"has_container": t.Container != nil,
	}

	if t.Container != nil {
		info["container_type"] = "default"
		info["service_count"] = t.Container.Count()
		info["registered_services"] = t.Container.RegisteredServices()
	}

	return info
}

// Compile-time interface compliance check
var _ interfaces.ContainableInterface = (*Containable)(nil)
