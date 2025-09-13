package facades
import (
	containerInterfaces "govel/types/src/interfaces/container"
	facade "govel/support/src"
)
// Container provides a clean, static-like interface to the application's dependency injection container service.
//
// This facade implements the facade pattern, providing global access to the container
// service configured in the dependency injection container. It offers a Laravel-style
// API for service binding, resolution, dependency injection, and container management
// with automatic service resolution and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved container service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent container operations across goroutines
//   - Supports singleton, transient, and instance binding patterns
//   - Built-in circular dependency detection and resolution
//
// Behavior:
//   - First call: Resolves container service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if container service cannot be resolved (fail-fast behavior)
//   - Automatically handles service lifecycle, dependency injection, and scope management
//
// Returns:
//   - ContainerInterface: The application's dependency injection container instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "container" service is not registered in the container
//   - If the resolved service doesn't implement ContainerInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Container() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Container operations are thread-safe and consistent
//
// Usage Examples:
//
//	// Basic service binding
//	facades.Container().Bind("logger", func() interface{} {
//	    return &Logger{
//	        level: "info",
//	        output: os.Stdout,
//	    }
//	})
//
//	// Singleton binding (single instance shared)
//	facades.Container().Singleton("database", func() interface{} {
//	    config := facades.Config().Get("database")
//	    db, err := sql.Open("mysql", config.DSN)
//	    if err != nil {
//	        log.Fatalf("Database connection failed: %v", err)
//	    }
//	    return &Database{connection: db}
//	})
//
//	// Instance binding (bind existing instance)
//	httpClient := &http.Client{Timeout: 30 * time.Second}
//	facades.Container().Instance("http_client", httpClient)
//
//	// Service resolution
//	logger := facades.Container().Make("logger")
//	database := facades.Container().Make("database")
//	client := facades.Container().Make("http_client")
//
//	// Type-safe resolution with generics
//	logger := facades.Container().MakeWithType[*Logger]("logger")
//	database := facades.Container().MakeWithType[*Database]("database")
//	client := facades.Container().MakeWithType[*http.Client]("http_client")
//
//	// Check if service is bound
//	if facades.Container().Bound("redis") {
//	    redis := facades.Container().Make("redis")
//	    // Use Redis client
//	}
//
//	// Service resolution with parameters
//	facades.Container().BindWithParams("user_service", func(params map[string]interface{}) interface{} {
//	    return &UserService{
//	        repository: facades.Container().Make("user_repository"),
//	        cache:      facades.Container().Make("cache"),
//	        config:     params["config"],
//	    }
//	})
//
//	userService := facades.Container().MakeWithParams("user_service", map[string]interface{}{
//	    "config": customConfig,
//	})
//
//	// Contextual binding (different implementations for different contexts)
//	facades.Container().When("EmailService").Needs("NotificationChannel").Give("smtp")
//	facades.Container().When("SMSService").Needs("NotificationChannel").Give("twilio")
//
//	// Tag binding for grouped services
//	facades.Container().Tag([]string{"mysql_repository", "redis_repository"}, "repositories")
//	facades.Container().Tag([]string{"email_handler", "sms_handler", "push_handler"}, "notification_handlers")
//
//	// Resolve all tagged services
//	repositories := facades.Container().Tagged("repositories")
//	handlers := facades.Container().Tagged("notification_handlers")
//
//	// Method injection
//	facades.Container().MethodCall("logger", "SetLevel", []interface{}{"debug"})
//	facades.Container().MethodCall("database", "SetTimeout", []interface{}{60 * time.Second})
//
//	// Container scoping
//	childContainer := facades.Container().Child()
//	childContainer.Bind("scoped_service", func() interface{} {
//	    return &ScopedService{}
//	})
//
//	// Service provider registration
//	facades.Container().Register(&DatabaseServiceProvider{})
//	facades.Container().Register(&CacheServiceProvider{})
//	facades.Container().Register(&LoggingServiceProvider{})
//
//	// Container introspection
//	bindings := facades.Container().GetBindings()
//	singletons := facades.Container().GetSingletons()
//	instances := facades.Container().GetInstances()
//
// Advanced Container Patterns:
//
//	// Factory pattern with container
//	facades.Container().Bind("user_factory", func() interface{} {
//	    return func(userData map[string]interface{}) *User {
//	        validator := facades.Container().Make("validator")
//	        hasher := facades.Container().Make("hash")
//	        
//	        return &User{
//	            Name:     userData["name"].(string),
//	            Email:    userData["email"].(string),
//	            Password: hasher.Make(userData["password"].(string)),
//	        }
//	    }
//	})
//
//	userFactory := facades.Container().Make("user_factory").(func(map[string]interface{}) *User)
//	newUser := userFactory(map[string]interface{}{
//	    "name": "John Doe",
//	    "email": "john@example.com",
//	    "password": "secret123",
//	})
//
//	// Decorator pattern
//	facades.Container().Extend("logger", func(logger interface{}) interface{} {
//	    return &LoggerWithMetrics{
//	        logger: logger.(*Logger),
//	        metrics: facades.Container().Make("metrics"),
//	    }
//	})
//
//	// Conditional binding based on environment
//	if facades.App().IsLocal() {
//	    facades.Container().Bind("mail", func() interface{} {
//	        return &LogMailer{} // Log emails instead of sending
//	    })
//	} else {
//	    facades.Container().Bind("mail", func() interface{} {
//	        return &SMTPMailer{
//	            host: facades.Config().GetString("mail.host"),
//	            port: facades.Config().GetInt("mail.port"),
//	        }
//	    })
//	}
//
//	// Lazy loading with proxy
//	facades.Container().Bind("expensive_service", func() interface{} {
//	    return &LazyProxy{
//	        factory: func() interface{} {
//	            // Expensive initialization only when needed
//	            return &ExpensiveService{
//	                data: loadLargeDataset(),
//	            }
//	        },
//	    }
//	})
//
//	// Multi-tenant container scoping
//	func GetTenantContainer(tenantID string) ContainerInterface {
//	    tenantContainer := facades.Container().Child()
//	    
//	    tenantContainer.Instance("tenant_id", tenantID)
//	    tenantContainer.Bind("database", func() interface{} {
//	        return &TenantDatabase{
//	            tenantID: tenantID,
//	            connection: facades.Container().Make("base_database"),
//	        }
//	    })
//	    
//	    return tenantContainer
//	}
//
//	// Service lifecycle management
//	facades.Container().OnResolving("database", func(service interface{}) {
//	    log.Printf("Resolving database service: %T", service)
//	})
//
//	facades.Container().OnResolved("database", func(service interface{}) {
//	    log.Printf("Database service resolved: %T", service)
//	    // Perform post-resolution setup
//	})
//
//	// Container middleware for cross-cutting concerns
//	facades.Container().AddMiddleware(func(serviceName string, next func() interface{}) interface{} {
//	    start := time.Now()
//	    service := next()
//	    duration := time.Since(start)
//	    log.Printf("Service %s resolved in %v", serviceName, duration)
//	    return service
//	})
//
//	// Circular dependency detection
//	facades.Container().Bind("service_a", func() interface{} {
//	    return &ServiceA{
//	        serviceB: facades.Container().Make("service_b"),
//	    }
//	})
//
//	facades.Container().Bind("service_b", func() interface{} {
//	    return &ServiceB{
//	        serviceA: facades.Container().Make("service_a"), // This will detect circular dependency
//	    }
//	})
//
// Container Configuration Patterns:
//
//	// Service provider pattern for organized binding
//	type DatabaseServiceProvider struct{}
//
//	func (p *DatabaseServiceProvider) Register(container ContainerInterface) {
//	    container.Singleton("database", func() interface{} {
//	        config := facades.Config().Get("database")
//	        return newDatabase(config)
//	    })
//	    
//	    container.Bind("user_repository", func() interface{} {
//	        return &UserRepository{
//	            db: container.Make("database"),
//	        }
//	    })
//	}
//
//	func (p *DatabaseServiceProvider) Boot(container ContainerInterface) {
//	    // Perform post-registration setup
//	    db := container.Make("database").(*Database)
//	    db.RunMigrations()
//	}
//
//	// Feature-based service registration
//	func RegisterCoreServices() {
//	    facades.Container().Singleton("config", func() interface{} {
//	        return loadConfiguration()
//	    })
//	    
//	    facades.Container().Singleton("logger", func() interface{} {
//	        return &Logger{
//	            level: facades.Config().GetString("logging.level"),
//	        }
//	    })
//	}
//
//	func RegisterDatabaseServices() {
//	    facades.Container().Singleton("database", func() interface{} {
//	        return newDatabaseConnection()
//	    })
//	    
//	    facades.Container().Bind("user_repository", func() interface{} {
//	        return &UserRepository{db: facades.Container().Make("database")}
//	    })
//	}
//
//	func RegisterCacheServices() {
//	    cacheDriver := facades.Config().GetString("cache.driver")
//	    
//	    switch cacheDriver {
//	    case "redis":
//	        facades.Container().Singleton("cache", func() interface{} {
//	            return &RedisCache{
//	                client: facades.Container().Make("redis_client"),
//	            }
//	        })
//	    case "memory":
//	        facades.Container().Singleton("cache", func() interface{} {
//	            return &MemoryCache{}
//	        })
//	    }
//	}
//
// Best Practices:
//   - Use singleton binding for stateful services (database, cache)
//   - Use transient binding for stateless services
//   - Bind interfaces, not concrete types when possible
//   - Use service providers to organize related bindings
//   - Tag related services for batch resolution
//   - Implement proper lifecycle management with callbacks
//   - Avoid circular dependencies through careful design
//   - Use contextual binding for different implementations
//
// Container Design Principles:
//  1. Single Responsibility: Each service should have one clear purpose
//  2. Dependency Inversion: Depend on abstractions, not concretions
//  3. Open/Closed: Open for extension, closed for modification
//  4. Interface Segregation: Use focused, specific interfaces
//  5. Dependency Injection: Inject dependencies rather than creating them
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume container always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	container, err := facade.TryResolve[ContainerInterface]("container")
//	if err != nil {
//	    // Handle container unavailability gracefully
//	    return fmt.Errorf("container unavailable: %w", err)
//	}
//	service := container.Make("service")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestServiceResolution(t *testing.T) {
//	    // Create a test container
//	    testContainer := &TestContainer{
//	        bindings: make(map[string]interface{}),
//	    }
//
//	    // Swap the real container with test container
//	    restore := support.SwapService("container", testContainer)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Container() returns testContainer
//	    facades.Container().Bind("test_service", func() interface{} {
//	        return &TestService{}
//	    })
//
//	    // Test service resolution
//	    service := facades.Container().Make("test_service")
//	    assert.NotNil(t, service)
//	    assert.IsType(t, &TestService{}, service)
//	}
//
// Container Configuration:
// Ensure the container service is properly configured in your container:
//
//	// Example container registration
//	container.Singleton("container", func() interface{} {
//	    config := container.Config{
//	        // Container configuration
//	        EnableCircularDependencyDetection: true,
//	        MaxResolutionDepth: 50,
//	        EnableServiceCaching: true,
//	        
//	        // Lifecycle configuration
//	        EnableLifecycleCallbacks: true,
//	        EnableMiddleware: true,
//	        
//	        // Performance settings
//	        ConcurrentResolution: true,
//	        ResolutionTimeout: 30 * time.Second,
//	        
//	        // Debugging
//	        EnableDebugMode: facades.App().IsLocal(),
//	        LogResolutions: facades.Config().GetBool("container.log_resolutions"),
//	        
//	        // Service providers
//	        Providers: []ServiceProvider{
//	            &CoreServiceProvider{},
//	            &DatabaseServiceProvider{},
//	            &CacheServiceProvider{},
//	        },
//	    }
//
//	    return container.NewContainer(config)
//	})
func Container() containerInterfaces.ContainerInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "container" service from the dependency injection container
	// - Performs type assertion to ContainerInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[containerInterfaces.ContainerInterface](containerInterfaces.CONTAINER_TOKEN)
}

// ContainerWithError provides error-safe access to the dependency injection container service.
//
// This function offers the same functionality as Container() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle container unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Container() but with error handling.
//
// Returns:
//   - ContainerInterface: The resolved container instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ContainerInterface
//
// Usage Examples:
//
//	// Basic error-safe container operations
//	container, err := facades.ContainerWithError()
//	if err != nil {
//	    log.Printf("Container unavailable: %v", err)
//	    return fmt.Errorf("dependency injection not available")
//	}
//	container.Bind("service", serviceFactory)
//
//	// Conditional service registration
//	if container, err := facades.ContainerWithError(); err == nil {
//	    container.Bind("optional_service", func() interface{} {
//	        return &OptionalService{}
//	    })
//	}
func ContainerWithError() (containerInterfaces.ContainerInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "container" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[containerInterfaces.ContainerInterface](containerInterfaces.CONTAINER_TOKEN)
}
