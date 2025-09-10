package interfaces

// interfaces.go re-exports all interfaces from subdirectories to provide
// a clean, unified interface for importing. This allows users to import
// interfaces directly from the interfaces package without navigating
// subdirectories.
//
// Usage:
//   import "govel/packages/application/interfaces"
//
//   var application applicationInterfaces.ApplicationInterface

import (
	// Import all interface packages
	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	traitInterfaces "govel/packages/application/interfaces/traits"

	// External trait interfaces
	configInterfaces "govel/packages/config/interfaces"
	containerInterfaces "govel/packages/container/interfaces"
	loggerInterfaces "govel/packages/logger/interfaces"
)

// Main Application Interfaces
// These are the primary application interfaces

// ApplicationInterface represents the main application interface
type ApplicationInterface = applicationInterfaces.ApplicationInterface

// Application ISP (Interface Segregation Principle) Interfaces
// These interfaces provide focused, single-responsibility contracts

// ApplicationIdentityInterface manages application name and version
type ApplicationIdentityInterface = applicationInterfaces.ApplicationIdentityInterface

// ApplicationRuntimeInterface manages runtime state (console mode, testing mode)
type ApplicationRuntimeInterface = applicationInterfaces.ApplicationRuntimeInterface

// ApplicationTimingInterface manages application timing (start time, uptime)
type ApplicationTimingInterface = applicationInterfaces.ApplicationTimingInterface

// ApplicationInfoInterface provides comprehensive application information
type ApplicationInfoInterface = applicationInterfaces.ApplicationInfoInterface

// ApplicationProviderInterface manages service provider lifecycle operations
type ApplicationProviderInterface = applicationInterfaces.ApplicationProviderInterface

// Core Trait Interfaces
// These interfaces are for application-specific traits

// DirectableInterface represents components that manage directory paths
type DirectableInterface = traitInterfaces.DirectableInterface

// EnvironmentableInterface represents components that manage environment settings
type EnvironmentableInterface = traitInterfaces.EnvironmentableInterface

// HookableInterface represents components that can manage hooks
type HookableInterface = traitInterfaces.HookableInterface

// LifecycleableInterface represents components that participate in application lifecycle
type LifecycleableInterface = traitInterfaces.LifecycleableInterface

// LocalizableInterface represents components that manage localization settings
type LocalizableInterface = traitInterfaces.LocalizableInterface

// ShutdownableInterface represents components that can be gracefully shut down
type ShutdownableInterface = traitInterfaces.ShutdownableInterface

// MaintainableInterface represents components that manage maintenance mode
type MaintainableInterface = traitInterfaces.MaintainableInterface

// External Trait Interfaces
// These interfaces are from external packages

// ConfigurableInterface represents configuration functionality
type ConfigurableInterface = configInterfaces.ConfigurableInterface

// ContainableInterface represents dependency injection container functionality
type ContainableInterface = containerInterfaces.ContainableInterface

// LoggableInterface represents logging functionality
type LoggableInterface = loggerInterfaces.LoggableInterface

// Deprecated aliases for backward compatibility
// These follow the "Has*" pattern but are deprecated in favor of the specific interfaces above

// HasHookable represents components that can be booted (deprecated: use HookableInterface)
type HasHookable = traitInterfaces.HookableInterface

// Directable represents components that manage directory paths (deprecated: use DirectableInterface)
type Directable = traitInterfaces.DirectableInterface

// Environmentable represents components that manage environment settings (deprecated: use EnvironmentableInterface)
type Environmentable = traitInterfaces.EnvironmentableInterface

// Lifecycleable represents components that participate in application lifecycle (deprecated: use LifecycleableInterface)
type Lifecycleable = traitInterfaces.LifecycleableInterface

// Localizable represents components that manage localization settings (deprecated: use LocalizableInterface)
type Localizable = traitInterfaces.LocalizableInterface

// Shutdownable represents components that can be gracefully shut down (deprecated: use ShutdownableInterface)
type Shutdownable = traitInterfaces.ShutdownableInterface

// Maintainable represents components that manage maintenance mode (deprecated: use MaintainableInterface)
type Maintainable = traitInterfaces.MaintainableInterface

// Service Provider Interfaces
// These interfaces define contracts for service provider management

// ServiceProviderInterface defines the contract for all service providers
type ServiceProviderInterface = providerInterfaces.ServiceProviderInterface

// ProviderRepositoryInterface manages provider registration and lifecycle
type ProviderRepositoryInterface = providerInterfaces.ProviderRepositoryInterface

// TerminatableProvider defines providers that require graceful termination
type TerminatableProvider = providerInterfaces.TerminatableProvider

// DeferrableProvider defines providers that can be deferred
type DeferrableProvider = providerInterfaces.DeferrableProvider

// EventTriggeredProvider defines providers that can be triggered by events
type EventTriggeredProvider = providerInterfaces.EventTriggeredProvider
