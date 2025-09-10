# **GoVel Application Package - Complete Blueprint**

## **Package Overview**

The `application` package is the **central coordinator** for the GoVel framework, inspired by Laravel but designed for Go. It orchestrates all components using dependency injection, trait patterns, and the Interface Segregation Principle (ISP).

## **Core Architecture**

### **🎯 Main Purpose**

- **Application Lifecycle Management** (boot, run, shutdown)
- **Service Provider Coordination** (registration, priority, deferred loading)
- **Dependency Injection Container** (singleton, transient bindings)
- **Environment & Configuration Management**
- **Trait-based Feature Composition** (locale, directories, etc.)
- **Graceful Shutdown & Maintenance Mode**

---

## **📁 Directory Structure & File Responsibilities**

### **🏗️ Core Files**

```
application/
├── doc.go                     # Package documentation & examples
├── application.go            # Main Application struct & New() constructor
```

### **🏛️ Core Infrastructure (`core/`)**

```
core/
├── lifecycle.go              # Application lifecycle hooks (PreStart, PostStart, etc.)
├── shutdown.go              # Immediate shutdown coordination
├── container.go             # Service container methods (Bind, Singleton, Make)
├── bootable.go             # Bootable services management
├── hooks.go                # Event hooks system
├── shutdown/
│   └── manager.go          # Graceful shutdown orchestration
├── maintenance/
│   ├── manager.go          # Maintenance mode management
│   └── middleware.go       # Maintenance mode middleware
└── service_provider/
    ├── registry.go         # Service provider registration & priority
    ├── deferred_repository.go  # Deferred/lazy loading providers
    └── termination_manager.go   # Provider termination coordination
```

### **🧩 Interfaces (`interfaces/`)**

Following **Interface Segregation Principle** - small, focused interfaces:

```
interfaces/
├── interfaces.go           # Re-exports all interfaces for clean imports
├── application/           # Core application interfaces
├── providers/            # Service provider interfaces
└── traits/              # "Has*" capability interfaces
    ├── has_bootable_interface.go     # Boot capability
    ├── has_directories_interface.go   # Path management
    ├── has_environment_interface.go   # Environment handling
    ├── has_lifecycle_interface.go     # Lifecycle participation
    ├── has_locale_interface.go        # Localization
    ├── has_logger_interface.go        # Logging capability
    └── has_shutdown_interface.go      # Shutdown participation
```

### **🎨 Traits (`traits/`)**

**Trait Pattern Implementation** - business logic separated from application struct:

```
traits/
├── locale.go              # Locale management trait (dependency injection pattern)
├── directories.go         # Directory path management trait
├── environment.go         # Environment detection & configuration trait
└── application.go         # Application metadata methods
```

### **🏗️ Builders (`builders/`)**

**Builder Pattern** for fluent application creation:

```
builders/
└── app_builder.go         # Fluent API: NewApp().WithName().WithEnv().Build()
```

### **⚙️ Service Providers (`providers/`)**

**Template Method Pattern** service provider base classes:

```
providers/
├── base_service_provider.go         # Base implementation with defaults
├── deferred_service_provider.go     # Lazy-loading providers
└── terminatable_service_provider.go  # Graceful shutdown providers
```

### **📊 Enums (`enums/`)**

**Type-safe state management**:

```
enums/
├── app_state_enum.go      # Application states (uninitialized, booting, running, etc.)
├── environment_enum.go    # Environments (development, testing, production)
├── priority_enum.go       # Service provider priorities
├── status_enum.go         # Service statuses
└── timeout_enum.go        # Timeout configurations
```

### **🎯 Constants (`constants/`)**

**Default values & configuration**:

```
constants/
├── app_defaults.go        # Application defaults (name, version, locale)
├── directories.go         # Default directory paths
└── maintenance.go         # Maintenance mode constants
```

### **🔧 Internal Utilities (`internal/`)**

```
internal/
└── env_helper.go          # Environment variable helpers with fallbacks
```

### **📋 Types (`types/`)**

**Legacy/alternative type definitions**:

```
types/
├── application.go         # Alternative Application type definitions
├── app_builder.go        # Alternative builder types
└── [various provider types] # Alternative service provider types
```

---

## **🏗️ Architecture Patterns Used**

### **1. Dependency Injection Pattern**

- **Container**: Services registered via `Bind()` and `Singleton()`
- **Traits**: Accept interfaces, not concrete types
- **Providers**: Injectable into application lifecycle

### **2. Interface Segregation Principle (ISP)**

- Small, focused interfaces (`HasLocale`, `HasDirectories`)
- Components implement only needed interfaces
- Easy testing and mocking

### **3. Trait Pattern**

- Business logic separated from Application struct
- **HasLocale** example: Interface + Implementation + Integration
- Reusable across different application types

### **4. Builder Pattern**

- **AppBuilder**: Fluent API for application configuration
- Method chaining: `NewApp().WithName().WithEnv().Build()`

### **5. Template Method Pattern**

- **ServiceProvider**: Default implementations
- Concrete providers override only necessary methods

### **6. Observer Pattern**

- Lifecycle hooks (PreStart, PostStart, PreShutdown, PostShutdown)
- Services can participate in application events

---

## **🔄 Application Lifecycle Flow**

```
1. Creation        → application := application.New()
2. Configuration   → app.SetEnvironment(), app.UseConfigPath()
3. Registration    → app.Register(serviceProvider)
4. Booting         → app.Boot(ctx)
5. Running         → // Application serves requests
6. Shutdown        → app.ShutdownWithSignals(ctx)
```

---

## **📋 Implementation Plan (If Recreating)**

### **Phase 1: Core Structure**

1. **Create `application.go`** with main Application struct
   - Define all struct fields (basePath, environment, debug, etc.)
   - Implement `New()` constructor with sensible defaults
   - Add thread-safe field access methods with mutexes

2. **Define core interfaces** in `interfaces/traits/`
   - Create small, focused interfaces following ISP
   - `HasBootable`, `HasLifecycle`, `HasShutdown`, etc.
   - Each interface should have 1-3 methods maximum

3. **Implement container methods** in `core/container.go`
   - `Bind()`, `Singleton()`, `Make()`, `Has()`, `Forget()`, `Flush()`
   - Delegate to underlying container implementation
   - Add deferred provider loading support

4. **Basic lifecycle management** in `core/lifecycle.go`
   - `PreStart()`, `PostStart()`, `PreShutdown()`, `PostShutdown()`
   - `Shutdown()` and `GracefulShutdown()` methods
   - Execute hooks for all registered services

### **Phase 2: Service Provider System**

1. **Service provider registry** in `core/service_provider/`
   - `ServiceProviderRegistry` for managing providers
   - Priority-based sorting and registration
   - Boot state tracking

2. **Base provider classes** in `providers/`
   - `ServiceProvider` with default implementations
   - Template method pattern for `Register()` and `Boot()`
   - Priority system (0-99: core, 100-199: framework, 200+: app)

3. **Deferred loading system**
   - `DeferredProviderRepository` for lazy loading
   - Load providers only when their services are requested
   - Service name to provider mapping

4. **Priority-based ordering**
   - Sort providers by priority during registration
   - Core infrastructure providers load first
   - Application providers load last

### **Phase 3: Trait System**

1. **Choose a trait** (e.g., locale, directories)
   - Start with locale management as example
   - Define what functionality the trait should provide

2. **Create trait interface** (e.g., `LocaleAppInterface`)
   - Define minimal interface for Application to implement
   - Only include methods needed by the trait
   - Follow ISP - keep interface small and focused

3. **Implement trait struct** with business logic
   - `HasLocale` struct with all locale-related methods
   - Constructor that accepts the interface
   - All business logic contained within trait

4. **Create integration file** showing Application → Trait connection
   - Application methods that implement the interface
   - Convenience methods that delegate to trait
   - Example of how to use trait in Application

### **Phase 4: Builder & Configuration**

1. **AppBuilder pattern** for fluent creation
   - `NewApp()` to create builder
   - `WithName()`, `WithEnvironment()`, `WithDebug()` methods
   - `Build()` method that creates Application instance

2. **Environment detection** and configuration
   - Read from environment variables with fallbacks
   - Support for .env files
   - Environment-specific defaults (debug off in production)

3. **Default constants** and enums
   - Application defaults (name, version, locale)
   - Environment enums (development, testing, production)
   - Application state enums (uninitialized, booting, running)

### **Phase 5: Advanced Features**

1. **Maintenance mode** system
   - Put application into maintenance mode
   - IP and path bypasses
   - Custom maintenance messages
   - Middleware for checking maintenance status

2. **Graceful shutdown** coordination
   - Signal handling (SIGINT, SIGTERM)
   - Timeout-based shutdown
   - Service termination coordination
   - Drain mode for completing active operations

3. **Lifecycle hooks** system
   - Event registration and execution
   - Hook priorities and ordering
   - Error handling and aggregation

4. **Thread-safe operations** with mutexes
   - Protect shared state with RWMutex
   - Consistent locking patterns
   - Avoid deadlocks with proper lock ordering

---

## **🎯 Key Design Principles**

1. **Separation of Concerns**: Traits handle specific functionality
2. **Interface Segregation**: Small, focused interfaces
3. **Dependency Injection**: Services injected via container
4. **Laravel Inspiration**: Familiar API for Laravel developers
5. **Go Idioms**: Channels, context, error handling
6. **Thread Safety**: Mutexes protect shared state
7. **Testability**: Interfaces enable easy mocking

---

## **🔧 Development Guidelines**

### **File Organization**

- **One responsibility per file**: Don't mix concerns
- **Consistent naming**: Use kebab-case for files, PascalCase for types
- **Interface files**: Place in `interfaces/` with one interface per file
- **Constants**: Place in `constants/` with one constant group per file

### **Code Style**

- **Detailed docblocks**: Every public type, function, and method
- **Examples in docs**: Show real usage examples
- **Error handling**: Always handle and wrap errors appropriately
- **Context usage**: Pass context through all long-running operations

### **Testing Strategy**

- **Interface mocking**: Use interfaces to enable easy testing
- **Integration tests**: Test the full application lifecycle
- **Unit tests**: Test individual components in isolation
- **Benchmarks**: Measure performance of critical paths

---

## **📝 Usage Examples**

### **Basic Application Setup**

```go
package main

import (
    "context"
    "log"
    "time"
    
    "govel/packages/application"
    "govel/packages/application/builders"
)

func main() {
    // Create and configure application
    application := builders.NewApp().
        SetEnvironment("production").
        SetName("My GoVel Application").
        SetVersion("1.0.0").
        Build()
    
    // Boot the application
    ctx := context.Background()
    if err := application.Boot(ctx); err != nil {
        log.Fatalf("Failed to boot application: %v", err)
    }
    
    // Setup graceful shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Your application logic here...
    
    // Graceful shutdown
    if err := application.ShutdownWithSignals(shutdownCtx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
}
```

### **Service Provider Example**

```go
type DatabaseServiceProvider struct {
    service_providers.ServiceProvider
}

func (p *DatabaseServiceProvider) Register(application interfaces.ApplicationInterface) error {
    return application.Singleton("database", func() interface{} {
        return &DatabaseConnection{
            Host: application.Config().GetString("db.host", "localhost"),
            Port: application.Config().GetInt("db.port", 5432),
        }
    })
}

func (p *DatabaseServiceProvider) Priority() int {
    return 10 // High priority for infrastructure
}
```

### **Trait Usage Example**

```go
// Using locale trait
app.HasLocale().SetLocale("fr-CA")
app.HasLocale().SetFallbackLocale("fr")
app.HasLocale().SetTimezone("America/Montreal")

// Or using convenience methods
app.SetLocaleValue("fr-CA")
info := app.LocaleInfo()
fmt.Printf("Language: %s, Country: %s", info["language"], info["country"])
```

This architecture provides a robust, extensible foundation for Go web applications while maintaining familiar Laravel-like patterns and Go best practices.
