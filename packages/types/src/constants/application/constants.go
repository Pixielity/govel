package constants

import "time"

// Application default constants
const (
	// DefaultName is the default application name
	DefaultName = "Govel Application"

	// DefaultVersion is the default application version
	DefaultVersion = "1.0.0"

	// DefaultEnvironment is the default application environment
	DefaultEnvironment = "development"

	// DefaultDebug is the default debug mode
	DefaultDebug = true

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

// Path Separators and Extensions
const (
	// PathSeparatorStr is the standard path separator
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
