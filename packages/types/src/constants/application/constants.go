package constants

import "time"

// Application default constants
const (
	// DefaultName is the default application name
	DefaultName = "Govel Application"

	// DefaultVersion is the default application version
	DefaultVersion = "1.0.0"

	// DefaultDebug is the default debug mode
	DefaultDebug = false

	// DefaultLocale is the default application locale
	DefaultLocale = "en"

	// DefaultFallbackLocale is the default fallback locale
	DefaultFallbackLocale = "en"

	// DefaultTimezone is the default application timezone
	DefaultTimezone = "UTC"

	// DefaultRunningInConsole indicates if running in console mode by default
	DefaultRunningInConsole = false

	// DefaultRunningUnitTests indicates if running unit tests by default
	DefaultRunningUnitTests = false

	// DefaultShutdownTimeout is the default graceful shutdown timeout
	DefaultShutdownTimeout = 30 * time.Second
)

// Directory constants
const (
	// DirectoryPublic is the public directory name
	DirectoryPublic = "public"

	// DirectoryStorage is the storage directory name
	DirectoryStorage = "storage"

	// DirectoryConfig is the config directory name
	DirectoryConfig = "config"

	// DirectoryLogs is the logs directory name
	DirectoryLogs = "logs"

	// DirectoryCache is the cache directory name
	DirectoryCache = "cache"

	// DirectoryViews is the views directory name
	DirectoryViews = "resources/views"

	// DirectoryResources is the resources directory name
	DirectoryResources = "resources"

	// DirectoryBootstrap is the bootstrap directory name
	DirectoryBootstrap = "bootstrap"

	// DirectoryDatabase is the database directory name
	DirectoryDatabase = "database"

	// DirectoryFramework is the framework directory name
	DirectoryFramework = "storage/framework"
)

// Maintenance constants
const (
	// MaintenanceFileName is the name of the maintenance file
	MaintenanceFileName = "down"

	// MaintenanceFileActiveValue is the value indicating maintenance is active
	MaintenanceFileActiveValue = "1"
)
