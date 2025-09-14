package config

// Queue returns the queue configuration map.
// This matches Laravel's queue.php configuration structure exactly.
// This configuration handles background job processing, queue drivers,
// and job batching for asynchronous task execution.
func Queue() map[string]any {
	return map[string]any{

		// Default Queue Connection Name
		//
		//
		// The framework's queue supports a variety of backends via a single, unified
		// API, giving you convenient access to each backend using identical
		// syntax for each. The default queue connection is defined below.
		//
		"default": Env("QUEUE_CONNECTION", "database"),

		// Queue Connections
		//
		//
		// Here you may configure the connection options for every queue backend
		// used by your application. An example configuration is provided for
		// each backend supported by the framework. You're also free to add more.
		//
		// Drivers: "sync", "database", "beanstalkd", "sqs", "redis", "null"
		//
		"connections": map[string]any{

			// Synchronous Queue Connection
			//
			// Processes jobs immediately in the same request. Useful for
			// testing and development where you want jobs to run synchronously.
			"sync": map[string]any{
				// Queue driver type - processes jobs immediately
				"driver": "sync",
			},

			// Database Queue Connection
			//
			// Uses your application's database to store queued jobs.
			// This is reliable and works well for most applications.
			"database": map[string]any{
				// Queue driver type
				"driver": "database",

				// Database connection to use (empty for default)
				"connection": Env("DB_QUEUE_CONNECTION", ""),

				// Database table name for storing jobs
				"table": Env("DB_QUEUE_TABLE", "jobs"),

				// Default queue name for jobs
				"queue": Env("DB_QUEUE", "default"),

				// Seconds after which jobs are considered failed and retried
				"retry_after": Env("DB_QUEUE_RETRY_AFTER", 90),

				// Whether to wait for database transaction commit before processing
				"after_commit": false,
			},

			// Beanstalkd Queue Connection
			//
			// A simple, fast work queue service. Good for high-throughput
			// job processing with built-in job priorities and delays.
			"beanstalkd": map[string]any{
				// Queue driver type
				"driver": "beanstalkd",

				// Beanstalkd server hostname or IP address
				"host": Env("BEANSTALKD_QUEUE_HOST", "localhost"),

				// Default tube (queue) name for jobs
				"queue": Env("BEANSTALKD_QUEUE", "default"),

				// Seconds after which jobs are considered failed and retried
				"retry_after": Env("BEANSTALKD_QUEUE_RETRY_AFTER", 90),

				// Seconds to block waiting for jobs (0 = don't block)
				"block_for": 0,

				// Whether to wait for database transaction commit
				"after_commit": false,
			},

			// Amazon SQS Queue Connection
			//
			// AWS Simple Queue Service provides a managed queue service
			// with high availability and scalability.
			"sqs": map[string]any{
				// Queue driver type
				"driver": "sqs",

				// AWS Access Key ID for authentication
				"key": Env("AWS_ACCESS_KEY_ID", ""),

				// AWS Secret Access Key for authentication
				"secret": Env("AWS_SECRET_ACCESS_KEY", ""),

				// SQS Queue URL prefix (base URL for your queues)
				"prefix": Env("SQS_PREFIX", "https://sqs.us-east-1.amazonaws.com/your-account-id"),

				// Default SQS queue name
				"queue": Env("SQS_QUEUE", "default"),

				// Optional suffix to append to queue URLs
				"suffix": Env("SQS_SUFFIX", ""),

				// AWS region where your SQS queues are located
				"region": Env("AWS_DEFAULT_REGION", "us-east-1"),

				// Whether to wait for database transaction commit
				"after_commit": false,
			},

			// Redis Queue Connection
			//
			// Uses Redis as a fast, in-memory queue backend.
			// Excellent for high-performance job processing.
			"redis": map[string]any{
				// Queue driver type
				"driver": "redis",

				// Redis connection name from database config
				"connection": Env("REDIS_QUEUE_CONNECTION", "default"),

				// Default Redis queue/list name for jobs
				"queue": Env("REDIS_QUEUE", "default"),

				// Seconds after which jobs are considered failed and retried
				"retry_after": Env("REDIS_QUEUE_RETRY_AFTER", 90),

				// Seconds to block waiting for jobs (nil = default blocking)
				"block_for": nil,

				// Whether to wait for database transaction commit
				"after_commit": false,
			},
		},

		// Job Batching
		//
		//
		// The following options configure the database and table that store job
		// batching information. These options can be updated to any database
		// connection and table which has been defined by your application.
		//
		"batching": map[string]any{
			// Database connection to use for job batching
			"database": Env("DB_CONNECTION", "sqlite"),

			// Table name for storing job batch information
			"table": "job_batches",
		},

		// Failed Queue Jobs
		//
		//
		// These options configure the behavior of failed queue job logging so you
		// can control how and where failed jobs are stored. The framework ships with
		// support for storing failed jobs in a simple file or in a database.
		//
		// Supported drivers: "database-uuids", "dynamodb", "file", "null"
		//
		"failed": map[string]any{
			// Driver for storing failed jobs
			//
			// Options: "database-uuids", "dynamodb", "file", "null"
			"driver": Env("QUEUE_FAILED_DRIVER", "database-uuids"),

			// Database connection for failed job storage
			"database": Env("DB_CONNECTION", "sqlite"),

			// Table name for storing failed job information
			"table": "failed_jobs",
		},
	}
}
