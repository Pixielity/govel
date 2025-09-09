package interfaces

/**
 * DirectableInterface defines the contract for components that provide
 * directory path management functionality. This interface follows the Interface
 * Segregation Principle by focusing solely on directory-related operations.
 *
 * Components implementing this interface can provide standardized access to
 * application directories, custom path management, and directory utilities.
 */
type DirectableInterface interface {
	/**
	 * GetBasePath returns the base directory path of the application.
	 *
	 * @return string The absolute path to the application's base directory
	 */
	GetBasePath() string

	/**
	 * SetBasePath updates the base directory path for the application.
	 *
	 * @param path string The new base directory path
	 */
	SetBasePath(path string)

	/**
	 * PublicPath returns the path to the public directory.
	 *
	 * @return string The absolute path to the public directory
	 */
	PublicPath() string

	/**
	 * SetPublicPath sets a custom public directory path.
	 *
	 * @param path string The custom public directory path
	 */
	SetPublicPath(path string)

	/**
	 * StoragePath returns the path to the storage directory.
	 *
	 * @return string The absolute path to the storage directory
	 */
	StoragePath() string

	/**
	 * SetStoragePath sets a custom storage directory path.
	 *
	 * @param path string The custom storage directory path
	 */
	SetStoragePath(path string)

	/**
	 * ConfigPath returns the path to the configuration directory.
	 *
	 * @return string The absolute path to the config directory
	 */
	ConfigPath() string

	/**
	 * SetConfigPath sets a custom config directory path.
	 *
	 * @param path string The custom config directory path
	 */
	SetConfigPath(path string)

	/**
	 * LogPath returns the path to the logs directory.
	 *
	 * @return string The absolute path to the logs directory
	 */
	LogPath() string

	/**
	 * SetLogPath sets a custom logs directory path.
	 *
	 * @param path string The custom logs directory path
	 */
	SetLogPath(path string)

	/**
	 * CachePath returns the path to the cache directory.
	 *
	 * @return string The absolute path to the cache directory
	 */
	CachePath() string

	/**
	 * SetCachePath sets a custom cache directory path.
	 *
	 * @param path string The custom cache directory path
	 */
	SetCachePath(path string)

	/**
	 * ViewPath returns the path to the views directory.
	 *
	 * @return string The absolute path to the views directory
	 */
	ViewPath() string

	/**
	 * SetViewPath sets a custom views directory path.
	 *
	 * @param path string The custom views directory path
	 */
	SetViewPath(path string)

	/**
	 * GetCustomPath returns a custom directory path by key.
	 *
	 * @param key string The custom path key
	 * @return string The custom path or empty string if not found
	 */
	GetCustomPath(key string) string

	/**
	 * SetCustomPath sets a custom directory path.
	 *
	 * @param key string The custom path key
	 * @param path string The custom path value
	 */
	SetCustomPath(key, path string)

	/**
	 * GetAllCustomPaths returns all custom directory paths.
	 *
	 * @return map[string]string A copy of all custom paths
	 */
	GetAllCustomPaths() map[string]string

	/**
	 * ClearCustomPaths removes all custom directory paths.
	 */
	ClearCustomPaths()

	/**
	 * EnsureDirectoryExists creates the directory if it doesn't exist.
	 *
	 * @param path string The directory path to ensure exists
	 * @return error Any error that occurred during directory creation
	 */
	EnsureDirectoryExists(path string) error
}
