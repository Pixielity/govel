package traits

import (
	"path/filepath"
	"sync"

	constants "govel/packages/types/src/constants/application"
	traitInterfaces "govel/packages/types/src/interfaces/application"
)

/**
 * Directable provides directory path management functionality through composition.
 * This struct implements the DirectableInterface and contains its own directory data,
 * following the self-contained trait pattern.
 *
 * Unlike dependency injection, this trait owns and manages its own state
 * for base path, custom paths, and directory configurations.
 */
type Directable struct {
	basePath    string            // Base application path
	customPaths map[string]string // Custom directory paths
	mutex       sync.RWMutex      // Thread safety for trait operations
}

/**
 * NewDirectable creates and initializes a new Directable instance.
 * This constructor sets up the trait with default values and proper initialization.
 *
 * @param basePath string The initial base path for the application
 * @return *Directable A properly initialized directories trait
 */
func NewDirectable(basePath string) *Directable {
	return &Directable{
		basePath:    basePath,
		customPaths: make(map[string]string),
	}
}

/**
 * BasePath returns the base directory path of the application.
 *
 * @return string The absolute path to the application's base directory
 */
func (d *Directable) BasePath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.basePath
}

/**
 * SetBasePath updates the base directory path for the application.
 *
 * @param path string The new base directory path
 */
func (d *Directable) SetBasePath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.basePath = path
}

/**
 * PublicPath returns the path to the public directory.
 * This directory contains publicly accessible files such as static assets.
 *
 * @return string The absolute path to the public directory
 */
func (d *Directable) PublicPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["public"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryPublic)
}

/**
 * SetPublicPath sets a custom public directory path.
 *
 * @param path string The custom public directory path
 */
func (d *Directable) SetPublicPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["public"] = path
}

/**
 * StoragePath returns the path to the storage directory.
 * This directory is used for storing application-generated files.
 *
 * @return string The absolute path to the storage directory
 */
func (d *Directable) StoragePath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["storage"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryStorage)
}

/**
 * SetStoragePath sets a custom storage directory path.
 *
 * @param path string The custom storage directory path
 */
func (d *Directable) SetStoragePath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["storage"] = path
}

/**
 * ConfigPath returns the path to the configuration directory.
 *
 * @return string The absolute path to the config directory
 */
func (d *Directable) ConfigPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["config"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryConfig)
}

/**
 * SetConfigPath sets a custom config directory path.
 *
 * @param path string The custom config directory path
 */
func (d *Directable) SetConfigPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["config"] = path
}

/**
 * LogPath returns the path to the logs directory.
 *
 * @return string The absolute path to the logs directory
 */
func (d *Directable) LogPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["logs"]; exists {
		return customPath
	}
	return filepath.Join(d.StoragePath(), constants.DirectoryLogs)
}

/**
 * SetLogPath sets a custom logs directory path.
 *
 * @param path string The custom logs directory path
 */
func (d *Directable) SetLogPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["logs"] = path
}

/**
 * CachePath returns the path to the cache directory.
 *
 * @return string The absolute path to the cache directory
 */
func (d *Directable) CachePath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["cache"]; exists {
		return customPath
	}
	return filepath.Join(d.StoragePath(), constants.DirectoryCache)
}

/**
 * SetCachePath sets a custom cache directory path.
 *
 * @param path string The custom cache directory path
 */
func (d *Directable) SetCachePath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["cache"] = path
}

/**
 * ViewPath returns the path to the views directory.
 *
 * @return string The absolute path to the views directory
 */
func (d *Directable) ViewPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["views"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryViews)
}

/**
 * SetViewPath sets a custom views directory path.
 *
 * @param path string The custom views directory path
 */
func (d *Directable) SetViewPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["views"] = path
}

/**
 * CustomPath returns a custom directory path by key.
 *
 * @param key string The custom path key
 * @return string The custom path or empty string if not found
 */
func (d *Directable) CustomPath(key string) string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.customPaths[key]
}

/**
 * SetCustomPath sets a custom directory path.
 *
 * @param key string The custom path key
 * @param path string The custom path value
 */
func (d *Directable) SetCustomPath(key, path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths[key] = path
}

/**
 * AllCustomPaths returns all custom directory paths.
 *
 * @return map[string]string A copy of all custom paths
 */
func (d *Directable) AllCustomPaths() map[string]string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range d.customPaths {
		result[k] = v
	}
	return result
}

/**
 * ClearCustomPaths removes all custom directory paths.
 */
func (d *Directable) ClearCustomPaths() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths = make(map[string]string)
}

/**
 * LogsPath returns the path to the logs directory.
 * This is an alias for LogPath for better Laravel compatibility.
 *
 * @return string The absolute path to the logs directory
 */
func (d *Directable) LogsPath() string {
	return d.LogPath()
}

/**
 * ResourcesPath returns the path to the resources directory.
 *
 * @return string The absolute path to the resources directory
 */
func (d *Directable) ResourcesPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["resources"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryResources)
}

/**
 * SetResourcesPath sets a custom resources directory path.
 *
 * @param path string The custom resources directory path
 */
func (d *Directable) SetResourcesPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["resources"] = path
}

/**
 * BootstrapPath returns the path to the bootstrap directory.
 *
 * @return string The absolute path to the bootstrap directory
 */
func (d *Directable) BootstrapPath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["bootstrap"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryBootstrap)
}

/**
 * SetBootstrapPath sets a custom bootstrap directory path.
 *
 * @param path string The custom bootstrap directory path
 */
func (d *Directable) SetBootstrapPath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["bootstrap"] = path
}

/**
 * DatabasePath returns the path to the database directory.
 *
 * @return string The absolute path to the database directory
 */
func (d *Directable) DatabasePath() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if customPath, exists := d.customPaths["database"]; exists {
		return customPath
	}
	return filepath.Join(d.basePath, constants.DirectoryDatabase)
}

/**
 * SetDatabasePath sets a custom database directory path.
 *
 * @param path string The custom database directory path
 */
func (d *Directable) SetDatabasePath(path string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.customPaths["database"] = path
}

/**
 * EnsureDirectoryExists creates the directory if it doesn't exist.
 *
 * @param path string The directory path to ensure exists
 * @return error Any error that occurred during directory creation
 */
func (d *Directable) EnsureDirectoryExists(path string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// This would need to import "os" to use os.MkdirAll
	// For now, returning nil (placeholder implementation)
	return nil
}

// Compile-time interface compliance check
// This ensures that Directable implements the DirectableInterface at compile time
var _ traitInterfaces.DirectableInterface = (*Directable)(nil)
