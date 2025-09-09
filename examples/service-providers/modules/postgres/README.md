# PostgreSQL Service Provider Module

This module provides a comprehensive PostgreSQL database service provider for the Govel framework. It demonstrates advanced service provider patterns including deferred loading, configuration management, connection pooling, and proper resource cleanup.

## Features

- **Database Interface**: Clean, testable database abstraction
- **Connection Management**: Multi-connection support with health monitoring
- **Connection Pooling**: Configurable connection pool settings
- **Transaction Support**: Safe transaction handling with automatic rollback
- **Configuration Integration**: Flexible configuration loading from config services
- **Health Monitoring**: Built-in health checks and statistics
- **Resource Cleanup**: Proper cleanup during application termination
- **Deferred Loading**: Efficient loading only when services are needed

## Components

### 1. Interfaces (`interfaces.go`)

- **DatabaseInterface**: Core database operations (connect, query, execute, transactions)
- **TransactionInterface**: Transaction-specific operations
- **ConnectionManagerInterface**: Multi-connection management
- **DatabaseConfig**: Configuration structure
- **DatabaseStats**: Connection statistics

### 2. Database Implementation (`database.go`)

The `Database` struct implements `DatabaseInterface` and provides:

- Thread-safe connection management
- Context-aware query execution
- Connection pooling configuration
- Timeout handling
- Transaction support with proper error handling

```go
// Create database instance
db := NewDatabase(DatabaseConfig{
    Host:            "localhost",
    Port:            5432,
    Database:        "myapp",
    Username:        "dbuser",
    Password:        "dbpass",
    SSLMode:         "prefer",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnectTimeout:  30 * time.Second,
    QueryTimeout:    30 * time.Second,
})

// Connect to database
ctx := context.Background()
err := db.Connect(ctx)

// Execute queries
results, err := db.Query(ctx, "SELECT * FROM users WHERE active = $1", true)

// Use transactions
err = db.Transaction(ctx, func(tx TransactionInterface) error {
    err := tx.Execute("INSERT INTO users (name) VALUES ($1)", "John")
    if err != nil {
        return err
    }
    // Automatically commits if no error, rolls back if error
    return nil
})
```

### 3. Connection Manager (`connection-manager.go`)

The `ConnectionManager` handles multiple database connections:

- Connection registration and retrieval
- Health monitoring across all connections
- Statistics collection
- Bulk operations (connect all, close all)

```go
// Create connection manager
manager := NewConnectionManager()

// Create connections
db1, err := manager.CreateConnection(config1)
db2, err := manager.CreateConnection(config2)

// Health check all connections
healthResults := manager.HealthCheck(ctx)

// Get statistics
stats := manager.GetAllConnectionStats()
```

### 4. Service Provider (`service-provider.go`)

The `PostgreSQLServiceProvider` implements the service provider pattern:

- **Deferred Loading**: Services are only initialized when first accessed
- **Configuration Integration**: Loads settings from config services
- **Dependency Injection**: Registers services with the container
- **Resource Management**: Proper cleanup during termination

#### Registered Services

- `postgres.database`: Database factory (creates new connections)
- `postgres.connection_manager`: Singleton connection manager
- `postgres.default_connection`: Default database connection
- `database`: Alias for default connection

#### Configuration Keys

The provider reads these configuration keys:

- `database.host`: Database host (default: localhost)
- `database.port`: Database port (default: 5432)
- `database.database`: Database name
- `database.username`: Database username
- `database.password`: Database password
- `database.ssl_mode`: SSL mode (default: prefer)
- `database.max_open_connections`: Max open connections (default: 25)
- `database.max_idle_connections`: Max idle connections (default: 5)
- `database.connection_lifetime`: Connection lifetime (default: 5m)
- `database.connect_timeout`: Connection timeout (default: 30s)
- `database.query_timeout`: Query timeout (default: 30s)

## Usage Example

```go
package main

import (
    "context"
    "log"
    "github.com/yourorg/govel/internal/container"
    "github.com/yourorg/govel/examples/service-providers/modules/postgres"
)

func main() {
    // Create container
    c := container.NewContainer()
    
    // Register config and logger services
    // ... (see example.go for complete setup)
    
    // Create and register PostgreSQL provider
    provider := postgres.NewPostgreSQLServiceProvider()
    
    if err := provider.Register(c); err != nil {
        log.Fatal(err)
    }
    
    if err := provider.Boot(c); err != nil {
        log.Fatal(err)
    }
    
    // Use database service
    dbService, err := c.Make("database")
    if err != nil {
        log.Fatal(err)
    }
    
    db := dbService.(postgres.DatabaseInterface)
    
    // Execute query
    results, err := db.Query(context.Background(), "SELECT version()")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Query results: %+v", results)
    
    // Clean up
    provider.Terminate(c)
}
```

## Health Monitoring

The provider includes comprehensive health monitoring:

```go
// Get health status
healthData := provider.HealthCheck(context.Background())

// Returns:
// {
//   "status": "initialized",
//   "healthy": true,
//   "connections": {
//     "localhost_myapp_5432": true
//   },
//   "stats": {
//     "localhost_myapp_5432": {
//       "OpenConnections": 2,
//       "InUseConnections": 0,
//       "IdleConnections": 2,
//       ...
//     }
//   }
// }
```

## Testing

The module includes mock implementations for testing:

- `MockConfig`: Mock configuration service
- `MockLogger`: Mock logger service

See `example.go` for complete testing examples.

## Dependencies

- `github.com/lib/pq`: PostgreSQL driver
- Internal container package for dependency injection
- Base provider package for common functionality

## Security Considerations

- **Connection Security**: Supports SSL/TLS connections
- **Password Handling**: Passwords are not logged or exposed in health checks
- **Query Safety**: Uses parameterized queries to prevent SQL injection
- **Connection Limits**: Configurable connection limits to prevent resource exhaustion

## Performance Features

- **Connection Pooling**: Efficient connection reuse
- **Context Support**: Proper timeout and cancellation handling
- **Statistics**: Detailed connection and query statistics
- **Concurrent Safety**: Thread-safe operations throughout

## Error Handling

The module provides comprehensive error handling:

- Connection failures with proper cleanup
- Transaction rollback on errors
- Timeout handling for all operations
- Detailed error messages with context

## Best Practices

1. **Configuration**: Use configuration files or environment variables for database settings
2. **Connections**: Use the connection manager for multiple database scenarios
3. **Transactions**: Always use transactions for multi-statement operations
4. **Context**: Pass context for proper timeout and cancellation handling
5. **Cleanup**: Ensure proper cleanup by calling `Terminate()` during shutdown
6. **Health Checks**: Implement regular health checking in production
7. **Monitoring**: Monitor connection statistics for performance optimization

## Extensions

The module is designed for extensibility:

- Add custom query builders
- Implement caching layers
- Add database migration support
- Integrate with ORM libraries
- Add query logging and metrics
