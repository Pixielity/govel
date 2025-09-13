package interfaces

import (
	traitInterfaces "govel/packages/types/src/interfaces/application/traits"
	loggerInterfaces "govel/packages/types/src/interfaces/logger"
	configInterfaces "govel/packages/types/src/interfaces/config"
	containerInterfaces "govel/packages/types/src/interfaces/container"
)

// ApplicationInterface defines the contract for the GoVel application.
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
	HasIdentityInterface
	HasRuntimeInterface
	HasTimingInterface
	ApplicationInfoInterface
	ApplicationProviderInterface
}

// Note: HasIdentityInterface, HasRuntimeInterface, and HasTimingInterface
// are defined in separate files for better organization

// ApplicationInfoInterface provides comprehensive application information
type ApplicationInfoInterface interface {
	// GetApplicationInfo returns comprehensive application information
	GetApplicationInfo() map[string]interface{}
}

// ApplicationProviderInterface provides service provider management
type ApplicationProviderInterface interface {
	// RegisterProvider registers a service provider
	RegisterProvider(provider interface{}) error
	// GetProviders returns all registered providers
	GetProviders() []interface{}
}
