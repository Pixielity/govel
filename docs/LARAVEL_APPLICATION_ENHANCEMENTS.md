# GoVel Laravel Application Class Enhancements

## Overview

This document summarizes the Laravel Application class compatibility enhancements added to the GoVel framework. These enhancements provide better migration path for developers familiar with Laravel and improve the overall developer experience.

## Enhanced Features

### 1. Booting and Booted Callback Support

**New Methods:**
- `Booting(callback func(ApplicationInterface))` - Register callbacks to execute before provider booting
- `Booted(callback func(ApplicationInterface))` - Register callbacks to execute after provider booting  
- `BootProvidersLaravel() error` - Laravel-compatible provider booting without context
- `IsBootedProviders() bool` - Check if providers have been booted

**Usage Example:**
```go
app.Booting(func(app ApplicationInterface) {
    fmt.Println("Application is about to boot...")
})

app.Booted(func(app ApplicationInterface) {
    fmt.Println("Application has finished booting!")
})

if err := app.BootProvidersLaravel(); err != nil {
    return fmt.Errorf("failed to boot application providers: %w", err)
}
```

### 2. Service Provider Lifecycle Management

**New Methods:**
- `GetLoadedProvidersMap() map[string]bool` - Get map of loaded providers with their status
- `ProviderIsLoaded(provider string) bool` - Check if a specific provider is loaded
- `GetDeferredServices() map[string]string` - Get map of deferred services to providers
- `SetDeferredServices(services map[string]string)` - Set deferred services mapping
- `IsDeferredService(service string) bool` - Check if a service is deferred
- `AddDeferredServices(services map[string]string)` - Add services to deferred mapping

**Usage Example:**
```go
// Check provider status
if app.ProviderIsLoaded("*modules.PostgreSQLServiceProvider") {
    fmt.Println("PostgreSQL provider is loaded")
}

// Manage deferred services
deferredServices := map[string]string{
    "redis": "*modules.RedisServiceProvider",
    "queue": "*modules.QueueServiceProvider",
}
app.SetDeferredServices(deferredServices)

if app.IsDeferredService("redis") {
    fmt.Println("Redis service is deferred")
}
```

### 3. Terminating Callback Support

**New Methods:**
- `Terminating(callback func(ApplicationInterface)) interface{}` - Register termination callbacks
- Enhanced `Terminate(ctx context.Context) []error` - Execute callbacks before provider termination

**Usage Example:**
```go
app.Terminating(func(app ApplicationInterface) {
    fmt.Println("Application is terminating...")
    // Cleanup code here
})

shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

errors := app.Terminate(shutdownCtx)
if len(errors) > 0 {
    for _, err := range errors {
        log.Printf("Termination error: %v", err)
    }
}
```

### 4. Cache Path Management

**New Methods:**
- `GetCachedServicesPath() string` - Get path to cached services file
- `GetCachedPackagesPath() string` - Get path to cached packages file
- `GetCachedConfigPath() string` - Get path to cached config file
- `GetCachedRoutesPath() string` - Get path to cached routes file
- `GetCachedEventsPath() string` - Get path to cached events file

**Usage Example:**
```go
configPath := app.GetCachedConfigPath()
fmt.Printf("Config cache path: %s\n", configPath)

servicesPath := app.GetCachedServicesPath()
fmt.Printf("Services cache path: %s\n", servicesPath)
```

### 5. Cache Status Checks

**New Methods:**
- `ConfigurationIsCached() bool` - Check if configuration is cached
- `RoutesAreCached() bool` - Check if routes are cached
- `EventsAreCached() bool` - Check if events are cached
- `ServicesAreCached() bool` - Check if services are cached
- `PackagesAreCached() bool` - Check if packages are cached

**Usage Example:**
```go
if app.ConfigurationIsCached() {
    fmt.Println("Configuration is cached, loading from cache")
}

if app.RoutesAreCached() {
    fmt.Println("Routes are cached, loading from cache")
}
```

### 6. HTTP Exception Support

**New Types:**
- `HTTPException` - HTTP exception with status code and message
- `NotFoundHTTPException` - Specific 404 Not Found exception

**New Methods:**
- `Abort(code int, message string, headers map[string]string) error` - Throw HTTP exception
- `AbortIf(condition bool, code int, message string, headers map[string]string) error` - Conditional abort
- `AbortUnless(condition bool, code int, message string, headers map[string]string) error` - Conditional abort unless

**Usage Example:**
```go
if !authorized {
    return app.Abort(403, "Access denied", nil)
}

// Conditional abort
if err := app.AbortIf(!validated, 400, "Invalid input", nil); err != nil {
    return err
}

// Conditional abort unless
if err := app.AbortUnless(authenticated, 401, "Authentication required", nil); err != nil {
    return err
}
```

### 7. Request Handling (Placeholder)

**New Methods:**
- `HandleRequest(request interface{}) (interface{}, error)` - Handle HTTP requests
- `HandleCommand(input interface{}) (int, error)` - Handle console commands
- `RunningConsoleCommand(commands ...string) bool` - Check if running console commands

**Note:** These are placeholder implementations that require HTTP/Console kernel integration.

### 8. Maintenance Mode Support

**New Methods:**
- `IsDownForMaintenance() bool` - Check if application is in maintenance mode

**Usage Example:**
```go
if app.IsDownForMaintenance() {
    fmt.Println("Application is currently down for maintenance")
    return
}
```

## Interface Segregation Principle (ISP) Architecture

To maintain clean architecture, new functionality is organized into focused interfaces:

### New ISP Interfaces

1. **ApplicationCallbackInterface** - Booting/booted/terminating callbacks
2. **ApplicationLaravelBootInterface** - Laravel-compatible provider booting
3. **ApplicationCacheInterface** - Cache path management and status checks
4. **ApplicationHTTPExceptionInterface** - HTTP exception handling
5. **ApplicationRequestHandlingInterface** - HTTP and console request handling
6. **ApplicationLaravelProviderInterface** - Laravel-compatible provider lifecycle

### Updated Main Interface

The `ApplicationInterface` now includes all new ISP interfaces while maintaining backward compatibility:

```go
type ApplicationInterface interface {
    // Existing trait interfaces
    traitInterfaces.DirectableInterface
    traitInterfaces.EnvironmentableInterface
    // ... other existing interfaces
    
    // New Laravel-compatible enhancement interfaces
    ApplicationCallbackInterface
    ApplicationLaravelBootInterface
    ApplicationCacheInterface
    ApplicationHTTPExceptionInterface
    ApplicationRequestHandlingInterface
    ApplicationLaravelProviderInterface
    
    // Additional maintenance mode support
    IsDownForMaintenance() bool
}
```

## Implementation Notes

### Threading Safety
- All new methods use proper mutex locking for thread safety
- Callback arrays are properly synchronized
- Provider tracking maps are protected with read/write locks

### Laravel Compatibility
- Method names and signatures match Laravel's Application class where possible
- Return types and error handling follow Laravel patterns
- Callback execution order matches Laravel's behavior

### Error Handling
- HTTP exceptions provide proper error interfaces
- Panic recovery in callback execution prevents crashes  
- Comprehensive error reporting for termination issues

### Environment Variable Support
- Cache paths support environment variable overrides
- Follows Laravel's caching configuration patterns
- Proper path normalization for cross-platform compatibility

## Future Enhancements

### Pending Features
1. **HTTP Kernel Integration** - Complete implementation of HandleRequest method
2. **Console Kernel Integration** - Complete implementation of HandleCommand method
3. **Command Detection** - Proper console command name detection for RunningConsoleCommand
4. **Bootstrap Process** - Laravel-style bootstrap with phase callbacks
5. **Environment Pattern Matching** - `environment(...string) bool` method

### Recommended Next Steps
1. Implement HTTP kernel for complete request handling
2. Add console kernel for command processing
3. Enhance environment detection with pattern matching
4. Add version information methods
5. Implement mock updates for testing compatibility

## Migration Benefits

### For Laravel Developers
- Familiar method names and patterns
- Compatible callback registration and execution
- Similar provider lifecycle management
- Consistent cache path handling

### For GoVel Projects  
- Improved service provider management
- Better application lifecycle control
- Enhanced debugging and introspection capabilities
- Flexible cache management

## Conclusion

These enhancements significantly improve GoVel's compatibility with Laravel concepts while maintaining its Go-specific advantages. The Interface Segregation Principle architecture ensures maintainability and allows selective feature adoption. The new functionality provides a solid foundation for Laravel developers transitioning to GoVel.
