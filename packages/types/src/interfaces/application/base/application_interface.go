package interfaces

import (
	concernsInterfaces "govel/types/interfaces/application/concerns"
	traitInterfaces "govel/types/interfaces/application/traits"
	configInterfaces "govel/types/interfaces/config"
	containerInterfaces "govel/types/interfaces/container"
	loggerInterfaces "govel/types/interfaces/logger"
)

// ApplicationInterface defines the contract for the GoVel application.
// This interface combines multiple focused interfaces following the Interface
// Segregation Principle, providing a complete application contract while
// maintaining clear separation of concerns.
type ApplicationInterface interface {
	// Core trait interfaces for application functionality
	traitInterfaces.DirectableInterface
	traitInterfaces.EnvironmentableInterface
	traitInterfaces.LifecycleableInterface
	traitInterfaces.LocalizableInterface
	traitInterfaces.ShutdownableInterface
	traitInterfaces.MaintainableInterface

	// External trait interfaces
	loggerInterfaces.LoggableInterface
	configInterfaces.ConfigurableInterface
	containerInterfaces.ContainableInterface

	// Core application ISP interfaces
	// These interfaces follow the Interface Segregation Principle
	concernsInterfaces.HasInfoInterface
	concernsInterfaces.HasTimingInterface
	concernsInterfaces.HasRuntimeInterface
	concernsInterfaces.HasIdentityInterface
	concernsInterfaces.HasBootstrapInterface
	concernsInterfaces.ApplicationProviderInterface
}
