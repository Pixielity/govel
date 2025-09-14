package config

// Database returns the database configuration map.
// This matches Laravel's database.php configuration structure exactly.
func Database() map[string]any {
	return map[string]any{

		// Default Database Connection Name
		//
		//
		// Here you may specify which of the database connections below you wish
		// to use as your default connection for database operations. This is
		// the connection which will be utilized unless another connection
		// is explicitly specified when you execute a query / statement.
		//
		"default": Env("DB_CONNECTION", "sqlite"),

		// Database Connections
		//
		//
		// Below are all of the database connections defined for your application.
		// An example configuration is provided for each database system which
		// is supported by Laravel. You're free to add / remove connections.
		//
		"connections": map[string]any{

			// SQLite Database Connection
			//
			// SQLite is a file-based database engine perfect for development and
			// small to medium-sized applications. It requires no server setup.
			"sqlite": map[string]any{
				// Database driver type
				"driver": "sqlite",

				// Database connection URL (optional)
				"url": Env("DB_URL", ""),

				// Database file path
				//
				// Path to the SQLite database file. Can be absolute or relative
				// to the application directory.
				"database": Env("DB_DATABASE", "database.sqlite"),

				// Table prefix for all database tables
				"prefix": "",

				// Foreign Key Constraints
				//
				// Enable or disable foreign key constraint enforcement.
				// SQLite requires this to be explicitly enabled.
				"foreign_key_constraints": Env("DB_FOREIGN_KEYS", true),

				// Busy timeout in milliseconds (nil for default)
				"busy_timeout": nil,

				// Journal mode (DELETE, TRUNCATE, PERSIST, MEMORY, WAL, OFF)
				"journal_mode": nil,

				// Synchronous mode (OFF, NORMAL, FULL, EXTRA)
				"synchronous": nil,

				// Transaction mode for this connection
				"transaction_mode": "DEFERRED",
			},

			// MySQL Database Connection
			//
			// MySQL is a popular open-source relational database management system.
			// It's widely used for web applications and supports ACID transactions.
			"mysql": map[string]any{
				// Database driver type
				"driver": "mysql",

				// Database connection URL (optional, overrides individual settings)
				"url": Env("DB_URL", ""),

				// Database server hostname or IP address
				"host": Env("DB_HOST", "127.0.0.1"),

				// Database server port number
				"port": Env("DB_PORT", "3306"),

				// Database name to connect to
				"database": Env("DB_DATABASE", "govel"),

				// Database username for authentication
				"username": Env("DB_USERNAME", "root"),

				// Database password for authentication
				"password": Env("DB_PASSWORD", ""),

				// Unix socket path (alternative to host/port)
				"unix_socket": Env("DB_SOCKET", ""),

				// Character set for database connection
				"charset": Env("DB_CHARSET", "utf8mb4"),

				// Collation for database connection
				"collation": Env("DB_COLLATION", "utf8mb4_unicode_ci"),

				// Table prefix for all database tables
				"prefix": "",

				// Enable prefix for indexes
				"prefix_indexes": true,

				// Enable strict SQL mode
				"strict": true,

				// MySQL storage engine (nil for default)
				"engine": nil,

				// Additional PDO options
				"options": map[string]any{},
			},

			// MariaDB Database Connection
			//
			// MariaDB is a MySQL-compatible database with additional features
			// and performance improvements. It's a drop-in MySQL replacement.
			"mariadb": map[string]any{
				// Database driver type
				"driver": "mariadb",

				// Database connection URL (optional, overrides individual settings)
				"url": Env("DB_URL", ""),

				// Database server hostname or IP address
				"host": Env("DB_HOST", "127.0.0.1"),

				// Database server port number
				"port": Env("DB_PORT", "3306"),

				// Database name to connect to
				"database": Env("DB_DATABASE", "govel"),

				// Database username for authentication
				"username": Env("DB_USERNAME", "root"),

				// Database password for authentication
				"password": Env("DB_PASSWORD", ""),

				// Unix socket path (alternative to host/port)
				"unix_socket": Env("DB_SOCKET", ""),

				// Character set for database connection
				"charset": Env("DB_CHARSET", "utf8mb4"),

				// Collation for database connection
				"collation": Env("DB_COLLATION", "utf8mb4_unicode_ci"),

				// Table prefix for all database tables
				"prefix": "",

				// Enable prefix for indexes
				"prefix_indexes": true,

				// Enable strict SQL mode
				"strict": true,

				// MariaDB storage engine (nil for default)
				"engine": nil,

				// Additional PDO options
				"options": map[string]any{},
			},

			// PostgreSQL Database Connection
			//
			// PostgreSQL is a powerful, open-source object-relational database
			// with excellent standards compliance and advanced features.
			"pgsql": map[string]any{
				// Database driver type
				"driver": "pgsql",

				// Database connection URL (optional, overrides individual settings)
				"url": Env("DB_URL", ""),

				// Database server hostname or IP address
				"host": Env("DB_HOST", "127.0.0.1"),

				// Database server port number (default PostgreSQL port is 5432)
				"port": Env("DB_PORT", "5432"),

				// Database name to connect to
				"database": Env("DB_DATABASE", "govel"),

				// Database username for authentication
				"username": Env("DB_USERNAME", "root"),

				// Database password for authentication
				"password": Env("DB_PASSWORD", ""),

				// Character encoding for database connection
				"charset": Env("DB_CHARSET", "utf8"),

				// Table prefix for all database tables
				"prefix": "",

				// Enable prefix for indexes
				"prefix_indexes": true,

				// Schema search path for PostgreSQL
				"search_path": "public",

				// SSL connection mode (disable, allow, prefer, require, verify-ca, verify-full)
				"sslmode": "prefer",
			},

			// SQL Server Database Connection
			//
			// Microsoft SQL Server is a relational database management system
			// developed by Microsoft for enterprise applications.
			"sqlsrv": map[string]any{
				// Database driver type
				"driver": "sqlsrv",

				// Database connection URL (optional, overrides individual settings)
				"url": Env("DB_URL", ""),

				// Database server hostname or IP address
				"host": Env("DB_HOST", "localhost"),

				// Database server port number (default SQL Server port is 1433)
				"port": Env("DB_PORT", "1433"),

				// Database name to connect to
				"database": Env("DB_DATABASE", "govel"),

				// Database username for authentication
				"username": Env("DB_USERNAME", "root"),

				// Database password for authentication
				"password": Env("DB_PASSWORD", ""),

				// Character encoding for database connection
				"charset": Env("DB_CHARSET", "utf8"),

				// Table prefix for all database tables
				"prefix": "",

				// Enable prefix for indexes
				"prefix_indexes": true,

				// Enable connection encryption (optional)
				// "encrypt": Env("DB_ENCRYPT", "yes"),

				// Trust server certificate (optional)
				// "trust_server_certificate": Env("DB_TRUST_SERVER_CERTIFICATE", "false"),
			},
		},

		// Migration Repository Table
		//
		//
		// This table keeps track of all the migrations that have already run for
		// your application. Using this information, we can determine which of
		// the migrations on disk haven't actually been run on the database.
		//
		"migrations": map[string]any{
			// Migration Table Name
			//
			// The name of the table that keeps track of executed migrations.
			// This table stores which migrations have been run.
			"table": "migrations",

			// Update Migration Date on Publish
			//
			// Whether to update the migration date when publishing migrations.
			// This helps track when migrations were actually deployed.
			"update_date_on_publish": true,
		},

		// Redis Databases
		//
		//
		// Redis is an open source, fast, and advanced key-value store that also
		// provides a richer body of commands than a typical key-value system
		// such as Memcached. You may define your connection settings here.
		//
		"redis": map[string]any{

			// Redis Client Type
			//
			// The Redis client library to use for connections. Different clients
			// may have different performance characteristics and feature support.
			"client": Env("REDIS_CLIENT", "phpredis"),

			// Global Redis Options
			//
			// These options apply to all Redis connections and control
			// cluster behavior, key prefixes, and connection persistence.
			"options": map[string]any{
				// Cluster Configuration
				//
				// Redis cluster mode setting. Use "redis" for standard mode
				// or "cluster" for Redis Cluster deployments.
				"cluster": Env("REDIS_CLUSTER", "redis"),

				// Key Prefix
				//
				// Prefix applied to all Redis keys to avoid conflicts with
				// other applications using the same Redis instance.
				"prefix": Env("REDIS_PREFIX", Slug(Env("APP_NAME", "govel").(string))+"-database-"),

				// Persistent Connections
				//
				// Whether to use persistent connections to Redis. This can
				// improve performance by reusing connections across requests.
				"persistent": Env("REDIS_PERSISTENT", false),
			},

			// Default Redis Connection
			//
			// The primary Redis connection used for general database operations,
			// queues, and other non-cache Redis usage.
			"default": map[string]any{
				// Redis Connection URL (optional, overrides individual settings)
				"url": Env("REDIS_URL", ""),

				// Redis Server Hostname
				//
				// The hostname or IP address of the Redis server.
				"host": Env("REDIS_HOST", "127.0.0.1"),

				// Redis Authentication Username (Redis 6.0+)
				//
				// Username for Redis ACL authentication. Leave empty if not using ACLs.
				"username": Env("REDIS_USERNAME", ""),

				// Redis Authentication Password
				//
				// Password for Redis authentication. Leave empty for no authentication.
				"password": Env("REDIS_PASSWORD", ""),

				// Redis Server Port
				//
				// The port number where Redis is listening (default is 6379).
				"port": Env("REDIS_PORT", "6379"),

				// Redis Database Number
				//
				// The Redis database number to use (0-15 by default).
				// Different databases provide logical separation.
				"database": Env("REDIS_DB", "0"),

				// Maximum Connection Retries
				//
				// Number of times to retry connection attempts before failing.
				"max_retries": Env("REDIS_MAX_RETRIES", 3),

				// Retry Backoff Algorithm
				//
				// Algorithm used for calculating retry delays. Options include
				// "fixed", "linear", "exponential", "decorrelated_jitter".
				"backoff_algorithm": Env("REDIS_BACKOFF_ALGORITHM", "decorrelated_jitter"),

				// Base Backoff Time (milliseconds)
				//
				// Base delay in milliseconds for retry backoff calculations.
				"backoff_base": Env("REDIS_BACKOFF_BASE", 100),

				// Maximum Backoff Time (milliseconds)
				//
				// Maximum delay in milliseconds for retry backoff calculations.
				"backoff_cap": Env("REDIS_BACKOFF_CAP", 1000),
			},

			// Cache Redis Connection
			//
			// A separate Redis connection optimized for caching operations.
			// Uses a different database to separate cache data from other Redis data.
			"cache": map[string]any{
				// Redis Connection URL (optional, overrides individual settings)
				"url": Env("REDIS_URL", ""),

				// Redis Server Hostname
				"host": Env("REDIS_HOST", "127.0.0.1"),

				// Redis Authentication Username (Redis 6.0+)
				"username": Env("REDIS_USERNAME", ""),

				// Redis Authentication Password
				"password": Env("REDIS_PASSWORD", ""),

				// Redis Server Port
				"port": Env("REDIS_PORT", "6379"),

				// Redis Database Number for Cache
				//
				// Uses database 1 by default to separate cache data
				// from other Redis operations.
				"database": Env("REDIS_CACHE_DB", "1"),

				// Maximum Connection Retries
				"max_retries": Env("REDIS_MAX_RETRIES", 3),

				// Retry Backoff Algorithm
				"backoff_algorithm": Env("REDIS_BACKOFF_ALGORITHM", "decorrelated_jitter"),

				// Base Backoff Time (milliseconds)
				"backoff_base": Env("REDIS_BACKOFF_BASE", 100),

				// Maximum Backoff Time (milliseconds)
				"backoff_cap": Env("REDIS_BACKOFF_CAP", 1000),
			},
		},
	}
}
