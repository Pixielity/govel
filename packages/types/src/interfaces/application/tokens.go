package interfaces

// Standard tokens for application package
const (
	// APPLICATION_TOKEN is the main service token for application
	APPLICATION_TOKEN = "govel.application"

	// APPLICATION_FACTORY_TOKEN is the factory token for application
	APPLICATION_FACTORY_TOKEN = "govel.application.factory"

	// APPLICATION_MANAGER_TOKEN is the manager token for application
	APPLICATION_MANAGER_TOKEN = "govel.application.manager"

	// APPLICATION_INTERFACE_TOKEN is the interface token for application
	APPLICATION_INTERFACE_TOKEN = "govel.application.interface"

	// APPLICATION_CONFIG_TOKEN is the config token for application
	APPLICATION_CONFIG_TOKEN = "govel.application.config"
)

// Path-related tokens used by PathsServiceProvider
const (
	// PATHS_BASE_TOKEN is the token for base path
	PATHS_BASE_TOKEN = "paths.base"

	// PATHS_STORAGE_TOKEN is the token for storage path
	PATHS_STORAGE_TOKEN = "paths.storage"

	// PATHS_CONFIG_TOKEN is the token for config path
	PATHS_CONFIG_TOKEN = "paths.config"

	// PATHS_CACHE_TOKEN is the token for cache path
	PATHS_CACHE_TOKEN = "paths.cache"

	// PATHS_LOGS_TOKEN is the token for logs path
	PATHS_LOGS_TOKEN = "paths.logs"

	// PATHS_RESOURCES_TOKEN is the token for resources path
	PATHS_RESOURCES_TOKEN = "paths.resources"

	// PATHS_PUBLIC_TOKEN is the token for public path
	PATHS_PUBLIC_TOKEN = "paths.public"

	// PATHS_BOOTSTRAP_TOKEN is the token for bootstrap path
	PATHS_BOOTSTRAP_TOKEN = "paths.bootstrap"

	// PATHS_DATABASE_TOKEN is the token for database path
	PATHS_DATABASE_TOKEN = "paths.database"

	// PATHS_ALL_TOKEN is the token for all paths as a map
	PATHS_ALL_TOKEN = "paths.all"
)
