package types

import (
	containerInterfaces "govel/packages/types/src/interfaces/container"
)

// DriverCreator represents a function that creates drivers for the Manager.
// It receives the container for dependency injection, returning the created driver instance and any error.
type DriverCreator func(container containerInterfaces.ContainerInterface) (interface{}, error)
