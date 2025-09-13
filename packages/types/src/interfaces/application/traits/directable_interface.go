package interfaces

// DirectableInterface defines the contract for directory path management functionality.
type DirectableInterface interface {
	// BasePath returns the base directory path of the application
	BasePath() string
	// SetBasePath updates the base directory path for the application
	SetBasePath(path string)
	
	// PublicPath returns the path to the public directory
	PublicPath() string
	// SetPublicPath sets a custom public directory path
	SetPublicPath(path string)
	
	// StoragePath returns the path to the storage directory
	StoragePath() string
	// SetStoragePath sets a custom storage directory path
	SetStoragePath(path string)
	
	// ConfigPath returns the path to the configuration directory
	ConfigPath() string
	// SetConfigPath sets a custom config directory path
	SetConfigPath(path string)
	
	// LogPath returns the path to the logs directory
	LogPath() string
	// SetLogPath sets a custom logs directory path
	SetLogPath(path string)
	
	// CachePath returns the path to the cache directory
	CachePath() string
	// SetCachePath sets a custom cache directory path
	SetCachePath(path string)
	
	// ViewPath returns the path to the views directory
	ViewPath() string
	// SetViewPath sets a custom views directory path
	SetViewPath(path string)
	
	// ResourcesPath returns the path to the resources directory
	ResourcesPath() string
	// SetResourcesPath sets a custom resources directory path
	SetResourcesPath(path string)
	
	// BootstrapPath returns the path to the bootstrap directory
	BootstrapPath() string
	// SetBootstrapPath sets a custom bootstrap directory path
	SetBootstrapPath(path string)
	
	// DatabasePath returns the path to the database directory
	DatabasePath() string
	// SetDatabasePath sets a custom database directory path
	SetDatabasePath(path string)
	
	// LogsPath returns the path to the logs directory (alias for LogPath)
	LogsPath() string
	
	// CustomPath returns a custom directory path by key
	CustomPath(key string) string
	// SetCustomPath sets a custom directory path
	SetCustomPath(key, path string)
	// AllCustomPaths returns all custom directory paths
	AllCustomPaths() map[string]string
	// ClearCustomPaths removes all custom directory paths
	ClearCustomPaths()
	
	// EnsureDirectoryExists creates the directory if it doesn't exist
	EnsureDirectoryExists(path string) error
}
