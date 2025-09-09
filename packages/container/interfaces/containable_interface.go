package interfaces

/**
 * ContainableInterface defines the contract for components that provide
 * container functionality. This interface follows the Interface Segregation
 * Principle by focusing solely on container operations.
 *
 * By embedding ContainerInterface, this interface provides all container methods
 * plus trait-specific management methods.
 */
type ContainableInterface interface {
	// Embed ContainerInterface to get all container methods
	ContainerInterface

	/**
	 * GetContainer returns the container instance.
	 *
	 * @return ContainerInterface The container instance
	 */
	GetContainer() ContainerInterface

	/**
	 * SetContainer sets the container instance.
	 *
	 * @param container interface{} The container instance to set (using interface{} to avoid circular import)
	 */
	SetContainer(container interface{})

	/**
	 * HasContainer returns whether a container instance is set.
	 *
	 * @return bool true if a container is set
	 */
	HasContainer() bool

	/**
	 * GetContainerInfo returns information about the container.
	 *
	 * @return map[string]interface{} Container information
	 */
	GetContainerInfo() map[string]interface{}
}
