package types

import (
	applicationInterfaces "govel/packages/types/src/interfaces/application/base"
)

// InstanceCreator represents a function that creates instances for the MultipleInstanceManager.
// It receives the application container and configuration map, returning the created instance and any error.
type InstanceCreator func(app applicationInterfaces.ApplicationInterface, config map[string]interface{}) (interface{}, error)
