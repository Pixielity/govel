package config

// Grpc returns the gRPC configuration map.
// This configuration handles gRPC server settings, TLS configuration,
// interceptors, and service discovery for high-performance RPC communication.
func Grpc() map[string]any {
	return map[string]any{

		// Default gRPC Server Configuration
		//
		//
		// This configuration defines the default gRPC server settings including
		// host, port, TLS settings, and other server-specific options.
		//
		"host": Env("GRPC_HOST", "127.0.0.1"),
		"port": Env("GRPC_PORT", 9090),

		// TLS Configuration
		//
		//
		// These settings control TLS encryption for gRPC connections. You can
		// enable TLS and specify certificate and key file paths for secure
		// communication between services.
		//
		"tls": map[string]any{
			"enabled":   Env("GRPC_TLS_ENABLED", false),
			"cert_file": Env("GRPC_TLS_CERT_FILE", ""),
			"key_file":  Env("GRPC_TLS_KEY_FILE", ""),
			"ca_file":   Env("GRPC_TLS_CA_FILE", ""),
		},

		// Connection Settings
		//
		//
		// These settings control various connection parameters for the gRPC
		// server including timeouts, keepalive settings, and connection limits.
		//
		"connection": map[string]any{
			"max_receive_message_size": Env("GRPC_MAX_RECEIVE_SIZE", 4194304), // 4MB
			"max_send_message_size":    Env("GRPC_MAX_SEND_SIZE", 4194304),    // 4MB
			"connection_timeout":       Env("GRPC_CONNECTION_TIMEOUT", 30),    // seconds
			"keepalive_time":           Env("GRPC_KEEPALIVE_TIME", 30),        // seconds
			"keepalive_timeout":        Env("GRPC_KEEPALIVE_TIMEOUT", 5),      // seconds
			"max_connection_idle":      Env("GRPC_MAX_CONNECTION_IDLE", 300),  // seconds
			"max_connection_age":       Env("GRPC_MAX_CONNECTION_AGE", 300),   // seconds
		},

		// Interceptors
		//
		//
		// Here you may specify the interceptors that should be applied to all
		// gRPC requests. Interceptors can be used for authentication, logging,
		// metrics collection, and other cross-cutting concerns.
		//
		"interceptors": map[string]any{
			"unary": []string{
				"logging",
				"metrics",
				"auth",
				"recovery",
			},
			"stream": []string{
				"logging",
				"metrics",
				"auth",
				"recovery",
			},
		},

		// Service Discovery
		//
		//
		// Configuration for service discovery and registration with external
		// service registries such as Consul, etcd, or Kubernetes.
		//
		"service_discovery": map[string]any{
			"enabled": Env("GRPC_SERVICE_DISCOVERY_ENABLED", false),
			"registry": map[string]any{
				"type":     Env("GRPC_REGISTRY_TYPE", "consul"),
				"address":  Env("GRPC_REGISTRY_ADDRESS", "127.0.0.1:8500"),
				"username": Env("GRPC_REGISTRY_USERNAME", ""),
				"password": Env("GRPC_REGISTRY_PASSWORD", ""),
			},
			"service": map[string]any{
				"name":    Env("GRPC_SERVICE_NAME", "govel-service"),
				"version": Env("GRPC_SERVICE_VERSION", "1.0.0"),
				"tags":    []string{"grpc", "api"},
			},
		},

		// Reflection
		//
		//
		// Enable gRPC reflection for development and debugging. This allows
		// tools like grpcurl to introspect your gRPC services.
		//
		"reflection": map[string]any{
			"enabled": Env("GRPC_REFLECTION_ENABLED", false),
		},

		// Health Check
		//
		//
		// Configuration for gRPC health checking service. This is useful for
		// load balancers and service mesh configurations.
		//
		"health_check": map[string]any{
			"enabled": Env("GRPC_HEALTH_CHECK_ENABLED", true),
			"path":    Env("GRPC_HEALTH_CHECK_PATH", "/health"),
		},
	}
}
