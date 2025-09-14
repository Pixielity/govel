package config

// Cache returns the cache configuration map.
// This matches Laravel's cache.php configuration structure exactly.
func Cache() map[string]any {
	return map[string]any{

		// Default Cache Store
		//
		// This option controls the default cache store that will be used by the
		// framework. This connection is utilized if another isn't explicitly
		// specified when running a cache operation inside the application.
		"default": Env("CACHE_STORE", "database"),

		// Cache Stores
		//
		// Here you may define all of the cache "stores" for your application as
		// well as their drivers. You may even define multiple stores for the
		// same cache driver to group types of items stored in your caches.
		//
		// Supported drivers: "array", "database", "file", "memcached",
		//                    "redis", "dynamodb", "octane", "null"

		"stores": map[string]any{

			// Array Cache Store
			//
			// The array cache store keeps data in memory for the duration of the
			// request only. This is useful for testing or when you need a simple
			// cache that doesn't persist between requests.
			"array": map[string]any{
				// Cache Driver Type
				"driver": "array",

				// Serialize Cache Values
				//
				// Whether to serialize cached values before storing them.
				// Set to false for better performance with simple data types.
				"serialize": false,
			},

			// Database Cache Store
			//
			// The database cache store uses your application's database to store
			// cached data. This is a good default choice that works across
			// multiple servers and provides persistence.
			"database": map[string]any{
				// Cache Driver Type
				"driver": "database",

				// Database Connection
				//
				// The database connection to use for caching. Leave empty to use
				// the default connection from your database configuration.
				"connection": Env("DB_CACHE_CONNECTION", ""),

				// Cache Table Name
				//
				// The database table where cache entries will be stored.
				// This table will be created automatically during migration.
				"table": Env("DB_CACHE_TABLE", "cache"),

				// Lock Connection
				//
				// The database connection to use for cache locks. Leave empty
				// to use the same connection as the cache table.
				"lock_connection": Env("DB_CACHE_LOCK_CONNECTION", ""),

				// Lock Table Name
				//
				// The database table where cache locks will be stored.
				// Leave empty to use the same table as cache entries.
				"lock_table": Env("DB_CACHE_LOCK_TABLE", ""),
			},

			// File Cache Store
			//
			// The file cache store saves cached data to the local filesystem.
			// This is simple and fast for single-server applications but
			// doesn't work well with multiple servers.
			"file": map[string]any{
				// Cache Driver Type
				"driver": "file",

				// Cache Storage Path
				//
				// The directory where cache files will be stored.
				// Make sure this directory is writable by the application.
				"path": StoragePath("framework/cache/data"),

				// Lock File Path
				//
				// The directory where cache lock files will be stored.
				// Usually the same as the cache path.
				"lock_path": StoragePath("framework/cache/data"),
			},

			// Memcached Cache Store
			//
			// Memcached is a high-performance, distributed memory caching system.
			// It's excellent for applications that need fast cache access across
			// multiple servers.
			"memcached": map[string]any{
				// Cache Driver Type
				"driver": "memcached",

				// Persistent Connection ID
				//
				// An optional identifier for persistent connections. This allows
				// connection reuse across requests for better performance.
				"persistent_id": Env("MEMCACHED_PERSISTENT_ID", ""),

				// SASL Authentication
				//
				// Username and password for SASL authentication if your
				// Memcached server requires authentication.
				"sasl": []string{
					Env("MEMCACHED_USERNAME", "").(string),
					Env("MEMCACHED_PASSWORD", "").(string),
				},

				// Memcached Options
				//
				// Additional options for fine-tuning Memcached behavior.
				// Add timeout, compression, and other settings as needed.
				"options": map[string]any{
					// Connection timeout in milliseconds
					// "connect_timeout": 2000,
				},

				// Memcached Servers
				//
				// List of Memcached servers to connect to. You can specify
				// multiple servers for redundancy and load distribution.
				"servers": []map[string]any{
					{
						// Server hostname or IP address
						"host": Env("MEMCACHED_HOST", "127.0.0.1"),

						// Server port number
						"port": Env("MEMCACHED_PORT", 11211),

						// Server weight for load balancing (higher = more traffic)
						"weight": 100,
					},
				},
			},

			// Redis Cache Store
			//
			// Redis is an advanced key-value store that can be used as a cache.
			// It offers features like data persistence, pub/sub, and atomic operations.
			"redis": map[string]any{
				// Cache Driver Type
				"driver": "redis",

				// Redis Connection
				//
				// The Redis connection to use for caching. This should match
				// a connection defined in your database configuration.
				"connection": Env("REDIS_CACHE_CONNECTION", "cache"),

				// Lock Connection
				//
				// The Redis connection to use for cache locks. Usually different
				// from the main cache connection to avoid conflicts.
				"lock_connection": Env("REDIS_CACHE_LOCK_CONNECTION", "default"),
			},

			// DynamoDB Cache Store
			//
			// Amazon DynamoDB is a NoSQL database service that can be used
			// for caching. It's serverless and scales automatically.
			"dynamodb": map[string]any{
				// Cache Driver Type
				"driver": "dynamodb",

				// AWS Access Key ID
				//
				// Your AWS access key for authenticating with DynamoDB.
				// It's recommended to use IAM roles instead of hardcoded keys.
				"key": Env("AWS_ACCESS_KEY_ID", ""),

				// AWS Secret Access Key
				//
				// Your AWS secret key for authenticating with DynamoDB.
				// Keep this secure and consider using environment variables.
				"secret": Env("AWS_SECRET_ACCESS_KEY", ""),

				// AWS Region
				//
				// The AWS region where your DynamoDB table is located.
				// Choose a region close to your application for better performance.
				"region": Env("AWS_DEFAULT_REGION", "us-east-1"),

				// DynamoDB Table Name
				//
				// The name of the DynamoDB table to use for caching.
				// This table will be created automatically if it doesn't exist.
				"table": Env("DYNAMODB_CACHE_TABLE", "cache"),

				// DynamoDB Endpoint
				//
				// Optional custom endpoint for DynamoDB. Useful for local development
				// with DynamoDB Local or when using VPC endpoints.
				"endpoint": Env("DYNAMODB_ENDPOINT", ""),
			},
		},

		//stores, there might be other applications using the same cache. For
		// that reason, you may prefix every cache key to avoid collisions.
		"prefix": Env("CACHE_PREFIX", Slug(Env("APP_NAME", "govel").(string))+"-cache-"),
	}
}
