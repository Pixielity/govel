package providers

import (
	"fmt"
	applicationInterfaces "govel/packages/types/src/interfaces/application"
)

/**
 * PathsServiceProvider provides application path services through the container.
 *
 * This service provider binds all application directory paths to the container,
 * making them available for dependency injection. This follows the Laravel pattern
 * of making common application resources available via the service container.
 *
 * Features:
 * - Binds all application paths as container services
 * - Provides centralized path management
 * - Enables dependency injection of paths
 * - Supports path resolution via container.Make()
 *
 * Bound Services:
 * - "paths.base": Application base directory path
 * - "paths.storage": Storage directory path
 * - "paths.config": Configuration directory path
 * - "paths.cache": Cache directory path
 * - "paths.logs": Logs directory path
 * - "paths.resources": Resources directory path
 * - "paths.public": Public directory path
 * - "paths.bootstrap": Bootstrap directory path
 * - "paths.database": Database directory path
 *
 * Usage:
 *   storagePath, err := container.Make("paths.storage")
 *   configPath, err := container.Make("paths.config")
 */
type PathsServiceProvider struct {
	ServiceProvider
}

// NewPathsServiceProvider creates a new paths service provider.
//
// Returns:
//
//	*PathsServiceProvider: A new paths service provider ready for registration
func NewPathsServiceProvider() *PathsServiceProvider {
	return &PathsServiceProvider{
		ServiceProvider: ServiceProvider{},
	}
}

// Register binds all application paths into the service container.
// This method makes all application directory paths available for dependency injection.
//
// Parameters:
//
//	application: The application instance with directory path access
//
// Returns:
//
//	error: Any error that occurred during registration, nil if successful
func (p *PathsServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Bind base path
	if err := application.Bind("paths.base", func() interface{} {
		return application.BasePath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.base: %w", err)
	}

	// Bind storage path
	if err := application.Bind("paths.storage", func() interface{} {
		return application.StoragePath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.storage: %w", err)
	}

	// Bind config path
	if err := application.Bind("paths.config", func() interface{} {
		return application.ConfigPath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.config: %w", err)
	}

	// Bind cache path
	if err := application.Bind("paths.cache", func() interface{} {
		return application.CachePath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.cache: %w", err)
	}

	// Bind logs path
	if err := application.Bind("paths.logs", func() interface{} {
		return application.LogsPath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.logs: %w", err)
	}

	// Bind resources path
	if err := application.Bind("paths.resources", func() interface{} {
		return application.ResourcesPath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.resources: %w", err)
	}

	// Bind public path
	if err := application.Bind("paths.public", func() interface{} {
		return application.PublicPath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.public: %w", err)
	}

	// Bind bootstrap path
	if err := application.Bind("paths.bootstrap", func() interface{} {
		return application.BootstrapPath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.bootstrap: %w", err)
	}

	// Bind database path
	if err := application.Bind("paths.database", func() interface{} {
		return application.DatabasePath()
	}); err != nil {
		return fmt.Errorf("failed to register paths.database: %w", err)
	}

	// Bind all paths as a map for convenience
	if err := application.Bind("paths.all", func() interface{} {
		return map[string]string{
			"base":      application.BasePath(),
			"storage":   application.StoragePath(),
			"config":    application.ConfigPath(),
			"cache":     application.CachePath(),
			"logs":      application.LogsPath(),
			"resources": application.ResourcesPath(),
			"public":    application.PublicPath(),
			"bootstrap": application.BootstrapPath(),
			"database":  application.DatabasePath(),
		}
	}); err != nil {
		return fmt.Errorf("failed to register paths.all: %w", err)
	}

	return nil
}
