package config

// Kafka returns the Kafka configuration map.
// This configuration handles Apache Kafka message streaming,
// producer/consumer settings, and topic management for event-driven architectures.
func Kafka() map[string]any {
	return map[string]any{

		// Default Kafka Connection Name
		//
		// The default Kafka connection to use when no specific connection
		// is specified. This connection will be used for producers and consumers.
		"default": Env("KAFKA_CONNECTION", "default"),

		// Kafka Connections
		//
		// Here you may configure multiple Kafka cluster connections.
		// Each connection can have different brokers, security settings,
		// and client configurations for different environments or use cases.
		"connections": map[string]any{

			// Default Kafka Connection
			//
			// The primary Kafka cluster connection used for most
			// message streaming operations.
			"default": map[string]any{

				// Kafka Broker Addresses
				//
				// Comma-separated list of Kafka broker addresses.
				// Format: "host1:port1,host2:port2,host3:port3"
				"brokers": Env("KAFKA_BROKERS", "localhost:9092"),

				// Client ID
				//
				// Identifier for this client instance. Used for logging
				// and metrics identification on the Kafka side.
				"client_id": Env("KAFKA_CLIENT_ID", Env("APP_NAME", "govel").(string)+"-client"),

				// Security Protocol
				//
				// Protocol used to communicate with brokers.
				// Options: "plaintext", "ssl", "sasl_plaintext", "sasl_ssl"
				"security_protocol": Env("KAFKA_SECURITY_PROTOCOL", "plaintext"),

				// SASL Configuration
				//
				// Authentication settings when using SASL security protocols.
				"sasl": map[string]any{
					// SASL Mechanism
					//
					// Authentication mechanism to use with SASL.
					// Options: "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512", "GSSAPI"
					"mechanism": Env("KAFKA_SASL_MECHANISM", "PLAIN"),

					// SASL Username
					//
					// Username for SASL authentication.
					"username": Env("KAFKA_SASL_USERNAME", ""),

					// SASL Password
					//
					// Password for SASL authentication.
					"password": Env("KAFKA_SASL_PASSWORD", ""),
				},

				// SSL Configuration
				//
				// SSL/TLS settings when using SSL-enabled security protocols.
				"ssl": map[string]any{
					// CA Certificate File Path
					//
					// Path to the Certificate Authority (CA) certificate file
					// for verifying broker certificates.
					"ca_cert_file": Env("KAFKA_SSL_CA_CERT_FILE", ""),

					// Client Certificate File Path
					//
					// Path to the client certificate file for mutual TLS.
					"cert_file": Env("KAFKA_SSL_CERT_FILE", ""),

					// Client Private Key File Path
					//
					// Path to the client private key file for mutual TLS.
					"key_file": Env("KAFKA_SSL_KEY_FILE", ""),

					// Skip Certificate Verification
					//
					// Whether to skip SSL certificate verification.
					// Set to true only for development/testing.
					"skip_verify": Env("KAFKA_SSL_SKIP_VERIFY", false),
				},

				// Producer Configuration
				//
				// Settings specific to Kafka message producers.
				"producer": map[string]any{
					// Acknowledgment Level
					//
					// Number of acknowledgments the producer requires before
					// considering a message successfully sent.
					// Options: "no_ack" (0), "leader_ack" (1), "all" (-1)
					"acks": Env("KAFKA_PRODUCER_ACKS", "all"),

					// Retry Attempts
					//
					// Maximum number of times to retry sending a message
					// before giving up.
					"retries": Env("KAFKA_PRODUCER_RETRIES", 3),

					// Batch Size
					//
					// Maximum number of bytes to include in a batch of messages.
					// Larger batches can improve throughput but increase latency.
					"batch_size": Env("KAFKA_PRODUCER_BATCH_SIZE", 16384),

					// Linger Time (milliseconds)
					//
					// Time to wait for additional messages to batch together.
					// Higher values improve throughput but increase latency.
					"linger_ms": Env("KAFKA_PRODUCER_LINGER_MS", 5),

					// Buffer Memory
					//
					// Total bytes of memory the producer can use to buffer
					// messages waiting to be sent.
					"buffer_memory": Env("KAFKA_PRODUCER_BUFFER_MEMORY", 33554432),

					// Compression Type
					//
					// Compression algorithm to use for message batches.
					// Options: "none", "gzip", "snappy", "lz4", "zstd"
					"compression_type": Env("KAFKA_PRODUCER_COMPRESSION", "snappy"),

					// Request Timeout (milliseconds)
					//
					// Maximum time to wait for a response from Kafka brokers.
					"request_timeout_ms": Env("KAFKA_PRODUCER_REQUEST_TIMEOUT_MS", 30000),

					// Delivery Timeout (milliseconds)
					//
					// Upper bound on time to report success or failure
					// after a message is sent.
					"delivery_timeout_ms": Env("KAFKA_PRODUCER_DELIVERY_TIMEOUT_MS", 120000),

					// Max In-Flight Requests
					//
					// Maximum number of unacknowledged requests per connection.
					// Setting to 1 ensures message ordering.
					"max_in_flight_requests": Env("KAFKA_PRODUCER_MAX_IN_FLIGHT", 5),

					// Enable Idempotence
					//
					// Whether to enable idempotent producer to avoid
					// duplicate messages on retry.
					"idempotence": Env("KAFKA_PRODUCER_IDEMPOTENCE", true),
				},

				// Consumer Configuration
				//
				// Settings specific to Kafka message consumers.
				"consumer": map[string]any{
					// Consumer Group ID
					//
					// Identifier for the consumer group. All consumers with
					// the same group ID will share partition consumption.
					"group_id": Env("KAFKA_CONSUMER_GROUP_ID", Env("APP_NAME", "govel").(string)+"-consumers"),

					// Auto Offset Reset
					//
					// Where to start consuming when no committed offset exists.
					// Options: "earliest", "latest", "none"
					"auto_offset_reset": Env("KAFKA_CONSUMER_AUTO_OFFSET_RESET", "latest"),

					// Enable Auto Commit
					//
					// Whether to automatically commit offsets periodically.
					"enable_auto_commit": Env("KAFKA_CONSUMER_AUTO_COMMIT", true),

					// Auto Commit Interval (milliseconds)
					//
					// Frequency of automatic offset commits when auto commit
					// is enabled.
					"auto_commit_interval_ms": Env("KAFKA_CONSUMER_AUTO_COMMIT_INTERVAL_MS", 5000),

					// Session Timeout (milliseconds)
					//
					// Time to wait for heartbeats before considering
					// a consumer dead and triggering rebalancing.
					"session_timeout_ms": Env("KAFKA_CONSUMER_SESSION_TIMEOUT_MS", 30000),

					// Heartbeat Interval (milliseconds)
					//
					// Frequency of heartbeats to the consumer coordinator.
					// Should be less than session timeout.
					"heartbeat_interval_ms": Env("KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS", 3000),

					// Max Poll Records
					//
					// Maximum number of messages returned in a single poll.
					"max_poll_records": Env("KAFKA_CONSUMER_MAX_POLL_RECORDS", 500),

					// Max Poll Interval (milliseconds)
					//
					// Maximum delay between polls before consumer is considered
					// failed and removed from group.
					"max_poll_interval_ms": Env("KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS", 300000),

					// Fetch Min Bytes
					//
					// Minimum amount of data the server should return for
					// a fetch request.
					"fetch_min_bytes": Env("KAFKA_CONSUMER_FETCH_MIN_BYTES", 1),

					// Fetch Max Wait (milliseconds)
					//
					// Maximum amount of time to wait for fetch min bytes
					// to be available.
					"fetch_max_wait_ms": Env("KAFKA_CONSUMER_FETCH_MAX_WAIT_MS", 500),

					// Fetch Max Bytes
					//
					// Maximum amount of data the server should return for
					// a fetch request.
					"fetch_max_bytes": Env("KAFKA_CONSUMER_FETCH_MAX_BYTES", 52428800),

					// Max Partition Fetch Bytes
					//
					// Maximum amount of data per partition the server
					// will return.
					"max_partition_fetch_bytes": Env("KAFKA_CONSUMER_MAX_PARTITION_FETCH_BYTES", 1048576),
				},
			},

			// Local Development Kafka Connection
			//
			// Simplified configuration for local development with minimal security.
			"local": map[string]any{
				// Local Kafka broker address
				"brokers": "localhost:9092",

				// Client identifier for local development
				"client_id": Env("APP_NAME", "govel").(string) + "-local",

				// No security for local development
				"security_protocol": "plaintext",

				// Local Producer Settings (simplified)
				"producer": map[string]any{
					// Leader acknowledgment only (faster for development)
					"acks": "1",

					// Minimal retries for faster feedback
					"retries": 1,

					// Small batch size for immediate sending
					"batch_size": 1024,

					// No batching delay for development
					"linger_ms": 0,

					// No compression for simplicity
					"compression_type": "none",

					// Disabled for simpler development setup
					"idempotence": false,
				},

				// Local Consumer Settings (simplified)
				"consumer": map[string]any{
					// Consumer group for local development
					"group_id": Env("APP_NAME", "govel").(string) + "-local-consumers",

					// Start from beginning for development
					"auto_offset_reset": "earliest",

					// Auto-commit for simplicity
					"enable_auto_commit": true,

					// Shorter timeout for faster development feedback
					"session_timeout_ms": 10000,

					// Frequent heartbeats for development
					"heartbeat_interval_ms": 3000,
				},
			},

			// Production Kafka Connection
			//
			// High-performance, secure configuration for production environments.
			"production": map[string]any{
				// Production Kafka broker addresses (comma-separated)
				"brokers": Env("KAFKA_PROD_BROKERS", ""),

				// Production client identifier
				"client_id": Env("APP_NAME", "govel").(string) + "-prod",

				// Secure protocol for production (SASL with SSL)
				"security_protocol": Env("KAFKA_PROD_SECURITY_PROTOCOL", "sasl_ssl"),

				// Production SASL Authentication
				"sasl": map[string]any{
					// Strong authentication mechanism
					"mechanism": Env("KAFKA_PROD_SASL_MECHANISM", "SCRAM-SHA-256"),

					// Production username
					"username": Env("KAFKA_PROD_SASL_USERNAME", ""),

					// Production password
					"password": Env("KAFKA_PROD_SASL_PASSWORD", ""),
				},

				// Production SSL Configuration
				"ssl": map[string]any{
					// CA certificate for broker verification
					"ca_cert_file": Env("KAFKA_PROD_SSL_CA_CERT", ""),

					// Client certificate for mutual TLS
					"cert_file": Env("KAFKA_PROD_SSL_CERT", ""),

					// Client private key for mutual TLS
					"key_file": Env("KAFKA_PROD_SSL_KEY", ""),

					// Never skip verification in production
					"skip_verify": false,
				},

				// Production Producer Settings (high reliability)
				"producer": map[string]any{
					// Wait for all replicas (strongest durability)
					"acks": "all",

					// Higher retries for production reliability
					"retries": 10,

					// Larger batch size for better throughput
					"batch_size": 32768,

					// Small linger for balance of latency/throughput
					"linger_ms": 10,

					// Efficient compression for production
					"compression_type": "lz4",

					// Enable idempotence for exactly-once semantics
					"idempotence": true,

					// Limit for message ordering guarantees
					"max_in_flight_requests": 1,

					// Extended timeout for production stability
					"delivery_timeout_ms": 300000,
				},

				// Production Consumer Settings (high reliability)
				"consumer": map[string]any{
					// Production consumer group identifier
					"group_id": Env("KAFKA_PROD_CONSUMER_GROUP", ""),

					// Start from earliest for data integrity
					"auto_offset_reset": "earliest",

					// Manual commit for better control
					"enable_auto_commit": false,

					// Longer timeout for production stability
					"session_timeout_ms": 45000,

					// Less frequent heartbeats in production
					"heartbeat_interval_ms": 15000,

					// Extended poll interval for production workloads
					"max_poll_interval_ms": 600000,

					// Larger fetch size for production efficiency
					"max_partition_fetch_bytes": 2097152,
				},
			},
		},

		// Topic Configuration
		//
		// Default settings for Kafka topics used by the application.
		"topics": map[string]any{
			// Default Replication Factor
			//
			// Number of replicas for each partition across brokers.
			// Should be at least 2 for production.
			"default_replication_factor": Env("KAFKA_DEFAULT_REPLICATION_FACTOR", 1),

			// Default Partition Count
			//
			// Number of partitions for new topics.
			// More partitions allow higher parallelism.
			"default_partitions": Env("KAFKA_DEFAULT_PARTITIONS", 3),

			// Topic Prefix
			//
			// Prefix to add to all topic names to avoid conflicts.
			"prefix": Env("KAFKA_TOPIC_PREFIX", Env("APP_NAME", "govel").(string)+"."),

			// Auto-create Topics
			//
			// Whether to automatically create topics when they don't exist.
			// Should be false in production.
			"auto_create": Env("KAFKA_AUTO_CREATE_TOPICS", false),
		},

		// Schema Registry Configuration
		//
		// Settings for Confluent Schema Registry integration for Avro/JSON schemas.
		"schema_registry": map[string]any{
			// Schema Registry URL
			//
			// Base URL of the schema registry service.
			"url": Env("KAFKA_SCHEMA_REGISTRY_URL", ""),

			// Authentication
			//
			// Credentials for schema registry access.
			"auth": map[string]any{
				"username": Env("KAFKA_SCHEMA_REGISTRY_USERNAME", ""),
				"password": Env("KAFKA_SCHEMA_REGISTRY_PASSWORD", ""),
			},

			// Subject Name Strategy
			//
			// Strategy for naming schema subjects.
			// Options: "topic", "record", "topic_record"
			"subject_name_strategy": Env("KAFKA_SCHEMA_SUBJECT_STRATEGY", "topic"),
		},

		// Monitoring and Metrics
		//
		// Configuration for Kafka client monitoring and metrics collection.
		"monitoring": map[string]any{
			// Enable Metrics
			//
			// Whether to collect and expose Kafka client metrics.
			"enabled": Env("KAFKA_MONITORING_ENABLED", true),

			// Metrics Reporting Interval (seconds)
			//
			// How often to report metrics to monitoring systems.
			"interval_seconds": Env("KAFKA_METRICS_INTERVAL", 30),

			// JMX Port
			//
			// Port for JMX metrics exposure (if using Java clients).
			"jmx_port": Env("KAFKA_JMX_PORT", 0),
		},

		// Logging Configuration
		//
		// Settings for Kafka client logging behavior.
		"logging": map[string]any{
			// Log Level
			//
			// Minimum log level for Kafka client logs.
			// Options: "debug", "info", "warn", "error"
			"level": Env("KAFKA_LOG_LEVEL", "info"),

			// Enable Debug
			//
			// Whether to enable detailed debug logging.
			// Should be false in production due to verbosity.
			"debug": Env("KAFKA_DEBUG", false),
		},
	}
}