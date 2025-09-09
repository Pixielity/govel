package interfaces

import (
	traitInterfaces "govel/packages/application/interfaces/traits"
	configInterfaces "govel/packages/config/interfaces"
	containerInterfaces "govel/packages/container/interfaces"
	loggerInterfaces "govel/packages/logger/interfaces"
)

// ApplicationInterface defines the comprehensive interface that service providers
// and other components need from the application instance. This allows for better
// testing, decoupling from the concrete Application type, and provides access to all
// essential application services.
//
// This interface combines multiple focused interfaces following the Interface
// Segregation Principle, providing a complete application contract while
// maintaining clear separation of concerns.
type ApplicationInterface interface {
	// Core trait interfaces for application functionality
	traitInterfaces.DirectableInterface
	traitInterfaces.EnvironmentableInterface
	traitInterfaces.HookableInterface
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
	ApplicationIdentityInterface
	ApplicationRuntimeInterface
	ApplicationTimingInterface
	ApplicationInfoInterface
}
