package interfaces

// ContainableInterface defines the contract for components that can work with a container.
type ContainableInterface interface {
	ContainerInterface

	// SetContainer sets the container instance.
	SetContainer(container ContainerInterface)

	// GetContainer returns the container instance.
	GetContainer() ContainerInterface

	// HasContainer returns whether a container is set.
	HasContainer() bool
}
