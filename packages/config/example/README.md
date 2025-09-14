# GoVel Configuration Package - Examples

This directory contains comprehensive examples demonstrating the GoVel configuration package capabilities with different file formats and drivers.

## Files Overview

### Configuration Files

#### 1. `config/config.yaml`
YAML format configuration file demonstrating:
- Nested configuration structure
- Application settings (name, debug, version)
- Server configuration with SSL and CORS
- Database configurations
- Cache settings
- Logging configuration
- External services configuration

#### 2. `config/config.json`
JSON format configuration file with the same structure as YAML, showing:
- Complete application configuration
- Database connection pools
- Queue configurations
- Session management
- Broadcasting settings
- Mail service configurations

#### 3. `config/config.toml`
TOML format configuration file demonstrating:
- Human-readable configuration format
- Hierarchical structure with sections
- Array and nested object support
- Comments for documentation

### Go Example Code

#### `main.go`
Comprehensive demonstration of all config package features:

1. **File Driver Examples** - YAML, JSON, TOML file loading
2. **Environment Driver** - Environment variable configuration
3. **Memory Driver** - In-memory configuration with runtime changes
4. **Remote Driver** - Remote configuration (stub implementation)
5. **Configuration Watching** - File change detection
6. **Type-safe Access** - All supported data types with `spf13/cast`
7. **Driver Comparison** - Side-by-side comparison of all drivers

## Running the Examples

### Prerequisites

1. Ensure Go 1.21+ is installed
2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Environment Setup

The example automatically sets up environment variables for demonstration:
- `APP_NAME=GoVel from Environment`
- `APP_DEBUG=false`
- `APP_PORT=9090`
- `APP_DATABASE_HOST=env-database-host`
- `APP_DATABASE_PORT=5432`
- `APP_CACHE_ENABLED=true`
- `APP_CACHE_TTL=7200`

### Running the Demo

```bash
cd /path/to/govel/packages/config/example
go run main.go
```

## Example Output

The demo will produce output like:

```
=== GoVel Configuration Package Demo ===

1. File Driver Example (YAML)
  App Name: GoVel Config Example
  Debug Mode: true
  Server Port: 8080
  Timeout: 30s
  Database: localhost:5432 (max connections: 100)
  Social Logins: [google github facebook]

2. File Driver Example (JSON)
  App Name: GoVel Config Example
  CORS Enabled: true
  CORS Origins: [http://localhost:3000 https://example.com]
  Default DB: mysql
  MySQL: localhost:3306

3. File Driver Example (TOML)
  App Name: GoVel Config Example (TOML)
  Default Cache: redis
  Redis: localhost:6379
  Log Level: debug
  Log Format: json

4. Environment Driver Example
  App Name: GoVel from Environment
  Debug Mode: false
  Port: 9090
  Database: env-database-host:5432
  Cache Enabled: true (TTL: 7200s)
  Runtime Value: environment value

5. Memory Driver Example
  App Name: Memory Config App
  Port: 9000
  DB Host: memory-db-host
  Users: [admin user guest]
  Runtime Timestamp: 2023-12-07T10:30:45Z
  Runtime Counter: 42
  Analytics: true

6. Remote Driver Example
  Remote config load failed (expected): remote driver load not implemented
  Demo Value: remote value
  Timestamp: 1701944445
  Note: This is a stub implementation. Real remote drivers would connect to etcd, Consul, etc.

7. Configuration Watching Example
  Setting up file watching...
  ✅ Watching enabled. Try modifying the config/config.yaml file.
  ⏰ Waiting 3 seconds to demonstrate...
  ⏹️  Watching stopped.

8. Type-safe Configuration Access
  String Value: hello world
  Int Value: 42
  Int64 Value: 9223372036854775807
  Float Value: 3.14159
  Bool Value: true
  Duration Value: 30s
  Time Value: 2023-01-01 00:00:00 +0000 UTC
  String Slice: [apple banana cherry]
  Int Slice: [1 2 3 4 5]
  String Map: map[key1:value1 key2:value2]
  String Map String: map[email:john@example.com name:John]

  With defaults:
  Missing String: default value
  Missing Int: 100
  Missing Bool: false

  Type conversion safety:
  String as Int: 0
  Int as String: 42
  Float as Int: 3

9. Driver Comparison
  Comparing different drivers with the same interface:

  📁 File Driver - App Name: GoVel Config Example
  🌍 Env Driver - App Name: GoVel from Environment
  🧠 Memory Driver - App Name: Memory Config App
  🌐 Remote Driver - App Name: Remote Config App

  All drivers implement the same interface and provide consistent access patterns!

=== Demo Complete ===
```

## Key Features Demonstrated

### 1. Multiple File Formats
- **YAML**: Human-readable, great for complex configurations
- **JSON**: Machine-readable, widely supported
- **TOML**: Simple and readable, good for application configs

### 2. Driver Architecture
- **Consistent Interface**: All drivers implement the same interface
- **Swappable**: Easy to switch between different configuration sources
- **Type Safety**: All values are type-converted using `spf13/cast`

### 3. Configuration Access Patterns
```go
// Basic access with defaults
appName := config.GetString("app.name", "Default App")
port := config.GetInt("server.port", 8080)
debug := config.GetBool("app.debug", false)

// Nested configuration access
dbHost := config.GetString("database.host", "localhost")
maxConnections := config.GetInt("database.max_connections", 100)

// Array/slice access
socialLogins := config.GetStringSlice("features.social_login", []string{})
corsOrigins := config.GetStringSlice("server.cors.origins", []string{})

// Complex type conversion
timeout := config.GetDuration("server.timeout", 30*time.Second)
timestamp := config.GetTime("created_at", time.Now())
```

### 4. Environment Variable Mapping
Environment variables are automatically mapped using dot notation:
- `APP_NAME` → `name`
- `APP_DATABASE_HOST` → `database.host`
- `APP_CACHE_ENABLED` → `cache.enabled`

### 5. Runtime Configuration Changes
```go
// Memory driver supports runtime changes
config.Set("feature.enabled", true)
config.Set("limits.requests", 1000)

// Changes trigger watch callbacks
config.Watch(func() {
    fmt.Println("Configuration changed!")
})
```

### 6. File Watching
Real-time configuration reloading when files change:
```go
config.Watch(func() {
    // Reload application settings
    updateApplicationSettings()
})
```

## Integration Examples

### Web Server Configuration
```go
// Load server configuration
serverConfig := config.GetStringMap("server")
port := config.GetInt("server.port", 8080)
host := config.GetString("server.host", "localhost")
timeout := config.GetDuration("server.timeout", 30*time.Second)

// Configure CORS
corsEnabled := config.GetBool("server.cors.enabled", false)
corsOrigins := config.GetStringSlice("server.cors.origins", []string{})
```

### Database Configuration
```go
// Database connection settings
dbConfig := config.GetStringMap("database.connections.mysql")
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
    config.GetString("database.connections.mysql.username"),
    config.GetString("database.connections.mysql.password"),
    config.GetString("database.connections.mysql.host"),
    config.GetInt("database.connections.mysql.port"),
    config.GetString("database.connections.mysql.database"),
)
```

### Cache Configuration
```go
// Cache settings
cacheDriver := config.GetString("cache.default", "memory")
redisHost := config.GetString("cache.stores.redis.host", "localhost")
redisPort := config.GetInt("cache.stores.redis.port", 6379)
ttl := config.GetInt("cache.stores.redis.ttl", 3600)
```

## Advanced Usage

### Environment-specific Configurations
Create different config files for different environments:
- `config.yaml` - Base configuration
- `config.development.yaml` - Development overrides
- `config.production.yaml` - Production overrides

### Configuration Validation
```go
// Validate required configuration
requiredKeys := []string{
    "app.name", "database.host", "cache.default",
}

for _, key := range requiredKeys {
    if !config.Has(key) {
        log.Fatalf("Required configuration key missing: %s", key)
    }
}
```

### Dynamic Configuration Updates
```go
// Watch for configuration changes
config.Watch(func() {
    // Reload specific components
    reloadDatabaseConnections()
    updateCacheSettings()
    refreshExternalServiceConfigs()
})
```

This example demonstrates the full power and flexibility of the GoVel configuration package, showing how it can handle complex, real-world configuration scenarios with multiple data sources and formats.