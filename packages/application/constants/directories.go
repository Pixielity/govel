package constants

// Directory and Path Constants
// These constants define standard directory names and path segments used
// throughout the GoVel application for organizing files and resources.

const (
	// DirectoryStorage is the default storage directory name
	DirectoryStorage = "storage"

	// DirectoryFramework is the framework directory within storage
	DirectoryFramework = "framework"

	// DirectoryLogs is the logs directory name
	DirectoryLogs = "logs"

	// DirectoryCache is the cache directory name
	DirectoryCache = "cache"

	// DirectoryViews is the views/templates directory name
	DirectoryViews = "views"

	// DirectoryPublic is the public assets directory name
	DirectoryPublic = "public"

	// DirectoryConfig is the configuration directory name
	DirectoryConfig = "config"

	// DirectoryResources is the resources directory name
	DirectoryResources = "resources"

	// DirectoryBootstrap is the bootstrap directory name
	DirectoryBootstrap = "bootstrap"

	// DirectoryDatabase is the database directory name (migrations, seeds, etc.)
	DirectoryDatabase = "database"
)

// Path Separators and Extensions
const (
	// PathSeparator is the standard path separator (use filepath.Separator in actual code)
	PathSeparatorStr = "/"

	// ConfigFileExtension is the default configuration file extension
	ConfigFileExtension = ".env"

	// LogFileExtension is the default log file extension
	LogFileExtension = ".log"

	// CacheFileExtension is the default cache file extension
	CacheFileExtension = ".cache"
)

// Special File Names
const (
	// MaintenanceFileName is the maintenance mode file name
	MaintenanceFileName = "down"

	// EnvironmentFileName is the environment configuration file name
	EnvironmentFileName = ".env"

	// AppConfigFileName is the main application configuration file name
	AppConfigFileName = "application.json"

	// DatabaseConfigFileName is the database configuration file name
	DatabaseConfigFileName = "database.json"

	// CacheConfigFileName is the cache configuration file name
	CacheConfigFileName = "cache.json"

	// LogConfigFileName is the logging configuration file name
	LogConfigFileName = "logging.json"
)

// URL and Route Constants
const (
	// StaticPathPrefix is the default prefix for static file routes
	StaticPathPrefix = "/static/"

	// APIPathPrefix is the default prefix for API routes
	APIPathPrefix = "/api/"

	// AdminPathPrefix is the default prefix for admin routes
	AdminPathPrefix = "/admin/"

	// HealthCheckPath is the default health check endpoint path
	HealthCheckPath = "/health"

	// MetricsPath is the default metrics endpoint path
	MetricsPath = "/metrics"
)

// Custom Path Keys
// These are used as keys in the customPaths map for overriding default directories
const (
	// CustomPathKeyStorage is the key for custom storage path
	CustomPathKeyStorage = "storage"

	// CustomPathKeyLogs is the key for custom logs path
	CustomPathKeyLogs = "logs"

	// CustomPathKeyCache is the key for custom cache path
	CustomPathKeyCache = "cache"

	// CustomPathKeyViews is the key for custom views path
	CustomPathKeyViews = "views"

	// CustomPathKeyPublic is the key for custom public path
	CustomPathKeyPublic = "public"

	// CustomPathKeyConfig is the key for custom config path
	CustomPathKeyConfig = "config"
)
