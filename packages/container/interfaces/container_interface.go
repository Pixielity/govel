package interfaces

// ContainerInterface defines the service container contract.
// This interface provides dependency injection and service resolution capabilities.
//
// The container is responsible for:
// - Registering service bindings (abstract -> concrete mappings)
// - Managing singleton instances and their lifecycle
// - Resolving services and their dependencies
// - Providing introspection capabilities
//
// Example usage:
//
//	// Register a service
//	container.Bind("database", func() interface{} {
//		return &DatabaseConnection{
//			Host: "localhost",
//			Port: 5432,
//		}
//	})
//
//	// Register a singleton service
//	container.Singleton("logger", func() interface{} {
//		return &Logger{Level: "info"}
//	})
//
//	// Resolve a service
//	db, err := container.Make("database")
//	if err != nil {
//		return err
//	}
//
// The interface promotes:
// - Dependency injection patterns
// - Service location and resolution
// - Singleton lifecycle management
// - Testability through interface contracts
type ContainerInterface interface {
	// Bind registers a binding in the service container.
	// The binding maps an abstract service name to a concrete implementation.
	// Each call to Make() with this abstract will create a new instance.
	//
	// The concrete parameter can be:
	// - A function that returns an instance: func() interface{} { return &Service{} }
	// - A function with dependencies: func(dep1 Dependency1) interface{} { return &Service{dep1} }
	// - A struct type: &Service{}
	// - Any other value that should be returned as-is
	//
	// Parameters:
	//   abstract: The service name/key (e.g., "database", "logger", "mailer")
	//   concrete: The concrete implementation (function, struct, or instance)
	//
	// Returns:
	//   error: Any error that occurred during binding registration
	//
	// Example:
	//   err := container.Bind("mailer", func() interface{} {
	//       return &SMTPMailer{Host: "localhost", Port: 587}
	//   })
	Bind(abstract string, concrete interface{}) error

	// Singleton registers a singleton binding in the service container.
	// Singleton services are instantiated once and cached for subsequent requests.
	// This is useful for services that should be shared across the application.
	//
	// The first call to Make() will create the instance and cache it.
	// Subsequent calls will return the cached instance.
	//
	// Parameters:
	//   abstract: The service name/key
	//   concrete: The concrete implementation (function, struct, or instance)
	//
	// Returns:
	//   error: Any error that occurred during singleton registration
	//
	// Example:
	//   err := container.Singleton("database", func() interface{} {
	//       return &DatabaseConnection{
	//           Host: os.Getenv("DB_HOST"),
	//           Port: 5432,
	//       }
	//   })
	Singleton(abstract string, concrete interface{}) error

	// Make resolves a service from the container.
	// For regular bindings, creates a new instance each time.
	// For singletons, returns the cached instance or creates and caches a new one.
	//
	// The container will automatically resolve dependencies if the concrete
	// implementation is a function that requires parameters.
	//
	// Parameters:
	//   abstract: The service name/key to resolve
	//
	// Returns:
	//   interface{}: The resolved service instance
	//   error: Any error that occurred during resolution or instantiation
	//
	// Example:
	//   logger, err := container.Make("logger")
	//   if err != nil {
	//       return fmt.Errorf("failed to resolve logger: %w", err)
	//   }
	//   log := logger.(*Logger)
	Make(abstract string) (interface{}, error)

	// IsBound checks if a service is registered in the container.
	// This is useful for conditional service resolution or debugging.
	//
	// Parameters:
	//   abstract: The service name/key to check
	//
	// Returns:
	//   bool: true if the service is registered (either as binding or singleton), false otherwise
	//
	// Example:
	//   if container.IsBound("optional-service") {
	//       service, _ := container.Make("optional-service")
	//       // Use the service
	//   }
	IsBound(abstract string) bool

	// Forget removes a service binding from the container.
	// For singletons, this also removes the cached instance.
	// After calling Forget, the service will no longer be resolvable.
	//
	// Parameters:
	//   abstract: The service name/key to remove
	//
	// Example:
	//   container.Forget("temporary-service")
	//   // container.Has("temporary-service") will now return false
	Forget(abstract string)

	// Flush removes all service bindings and cached instances.
	// This effectively resets the container to its initial empty state.
	// This is primarily useful for testing scenarios where you need a clean container.
	//
	// After calling Flush:
	// - All bindings are removed
	// - All singleton instances are destroyed
	// - Has() will return false for all services
	//
	// Example:
	//   container.Flush()
	//   // Container is now empty and ready for new bindings
	FlushContainer()
}
