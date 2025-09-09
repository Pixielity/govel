# GoVel Framework - Basic Example

This example demonstrates all the core features of the GoVel framework in a comprehensive, real-world scenario.

## What This Example Demonstrates

### 🏗️ **AppBuilder Pattern**
- **Fluent Interface**: Method chaining for clean, readable configuration
- **Environment Configurations**: `ForProduction()`, `ForDevelopment()`, `ForTesting()`
- **Comprehensive Settings**: Name, version, environment, debug, locale, timezone, etc.
- **Runtime Modes**: Console mode and testing mode configuration

### ⚙️ **Configuration Management**
- **Multiple Data Types**: Strings, integers, booleans, slices
- **Dot Notation**: Hierarchical configuration keys (`database.host`, `server.port`)
- **Default Values**: Fallback values when configuration keys don't exist
- **Dynamic Configuration**: Setting and retrieving values at runtime

### 📦 **Dependency Injection Container**
- **Service Registration**: Both singleton and transient service bindings
- **Service Resolution**: Retrieving services with type assertions
- **Dependency Injection**: Services automatically receive configuration values
- **Singleton Testing**: Verifying that singletons return the same instance

### 📝 **Structured Logging**
- **Basic Logging**: Info, debug, warning, and error messages
- **Structured Fields**: Adding context with `WithField()` and `WithFields()`
- **Field Chaining**: Combining multiple context fields
- **Component Logging**: Service-specific logging with components

### ⏰ **Application Lifecycle**
- **Timing Tracking**: Start time and uptime calculation
- **Application Information**: Comprehensive application metadata
- **Runtime State**: Console mode, testing mode, debug mode
- **Graceful Shutdown**: Clean service cleanup on termination

## Services Implemented

### 🗄️ **Database Service** (`PostgreSQLService`)
- Connection management
- Query execution with mock results
- Clean connection closing

### 💾 **Cache Service** (`RedisCache`)
- Key-value storage with TTL support
- Get, set, and delete operations
- In-memory mock implementation

### 🌐 **HTTP Service** (`HTTPService`)  
- Server start/stop functionality
- Configuration from application settings
- Clean shutdown handling

## Running the Example

### Prerequisites
Make sure you're in the GoVel project root directory and all dependencies are available.

### Run the Example
```bash
cd examples/basic
go run main.go
```

### Expected Output
The example will display comprehensive output showing:

1. **Application Information** - Name, version, environment settings
2. **Configuration Testing** - Setting and retrieving various config values
3. **Container Testing** - Service registration and resolution
4. **Logger Testing** - Structured logging with fields
5. **Service Usage** - Actual usage of resolved services
6. **Timing Information** - Application start time and uptime
7. **Builder Variations** - Different AppBuilder configuration patterns
8. **Graceful Shutdown Setup** - Signal handling for clean termination

### Test Graceful Shutdown
Once running, press `Ctrl+C` to trigger the graceful shutdown sequence, which will:
- Log the shutdown process
- Close database connections
- Stop HTTP server
- Clear cache data
- Exit cleanly

## Code Structure

### Service Interfaces
The example defines clean interfaces for each service type, demonstrating Go's interface-based design patterns.

### Service Implementations
Mock implementations show how real services would integrate with the GoVel framework.

### Configuration-Driven Services
All services receive their configuration from the application's configuration system, demonstrating centralized config management.

### Error Handling
Proper error handling throughout the application lifecycle.

## Key GoVel Features Showcased

### ✨ **Fluent Configuration**
```go
app := builders.NewApp().
    WithName("GoVel Basic Example").
    WithVersion("1.0.0").
    WithEnvironment("development").
    WithDebug(true).
    ForDevelopment().
    Build()
```

### 🔧 **Service Registration**
```go
app.Singleton("database", func() interface{} {
    return &PostgreSQLService{
        Host:     app.GetString("database.host", "localhost"),
        Port:     app.GetInt("database.port", 5432),
        Database: app.GetString("database.name", "govel"),
    }
})
```

### 📊 **Structured Logging**  
```go
app.WithFields(map[string]interface{}{
    "service": "cache",
    "host":    "localhost", 
    "port":    6379,
}).Info("Cache service configured")
```

### 🎯 **Service Resolution**
```go
dbService, err := app.Make("database")
db := dbService.(DatabaseService)
db.Connect()
```

## Learning Outcomes

After running this example, you'll understand:

1. **How to build applications** using the GoVel AppBuilder pattern
2. **Configuration management** with multiple data types and hierarchical keys
3. **Dependency injection** for loose coupling and testability
4. **Structured logging** for better observability
5. **Service lifecycle management** from registration to cleanup
6. **Graceful shutdown patterns** for production applications
7. **Environment-specific configurations** for different deployment scenarios

This example serves as a comprehensive reference for building applications with the GoVel framework, demonstrating real-world patterns and best practices.
