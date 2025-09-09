# Complete Service Provider Examples

This comprehensive example demonstrates the full service provider pattern implementation in Go, showcasing how to build modular, maintainable applications with proper dependency injection, configuration management, and lifecycle handling.

## 🚀 Features

### Core Patterns
- **Service Provider Architecture**: Modular service registration and initialization
- **Dependency Injection**: Container-based dependency management
- **Configuration Management**: Flexible configuration with environment variables
- **Lifecycle Management**: Proper bootstrapping, running, and shutdown phases
- **Health Monitoring**: Built-in health checks for all services
- **Deferred Loading**: Efficient service loading only when needed

### Service Providers Included

#### 1. PostgreSQL Provider (`modules/postgres/`)
- **Database Interface**: Clean abstraction over database operations
- **Connection Management**: Multi-connection support with health monitoring
- **Transaction Support**: Safe transaction handling with rollback
- **Connection Pooling**: Configurable connection pool settings
- **Statistics & Monitoring**: Detailed connection and query statistics

#### 2. Redis Provider (`modules/redis/`)
- **Cache Interface**: Comprehensive Redis operations (strings, lists, hashes, sets)
- **Multiple Redis Modes**: Support for standalone, sentinel, and cluster
- **Connection Management**: Multiple Redis instance support
- **Advanced Operations**: Pipelining, transactions, pub/sub ready
- **Statistics & Health**: Real-time Redis statistics and health monitoring

#### 3. Base Provider (`modules/base/`)
- **Common Functionality**: Shared base implementation for all providers
- **Lifecycle Management**: Standard registration, boot, and termination
- **Configuration Support**: Built-in configuration loading patterns

## 📁 Project Structure

```
examples/service-providers/
├── main.go                          # Main application demonstrating all providers
├── README.md                        # This documentation
├── config/                          # Configuration files
│   ├── app.yaml                     # Application configuration
│   ├── database.yaml                # Database settings
│   ├── redis.yaml                   # Redis settings
│   └── .env.example                 # Environment variables example
├── internal/                        # Core framework components
│   ├── application/                 # Application lifecycle management
│   ├── container/                   # Dependency injection container
│   ├── config/                      # Configuration management
│   └── logger/                      # Logging infrastructure
├── modules/                         # Service provider modules
│   ├── base/                        # Base provider functionality
│   │   ├── provider.go              # Base provider implementation
│   │   └── README.md                # Base provider documentation
│   ├── postgres/                    # PostgreSQL service provider
│   │   ├── interfaces.go            # Database interfaces
│   │   ├── database.go              # Database implementation
│   │   ├── connection-manager.go    # Connection management
│   │   ├── service-provider.go      # Service provider registration
│   │   ├── example.go               # Usage examples
│   │   └── README.md                # PostgreSQL documentation
│   ├── redis/                       # Redis service provider
│   │   ├── interfaces.go            # Cache interfaces
│   │   ├── cache.go                 # Redis implementation
│   │   ├── connection-manager.go    # Redis connection management
│   │   ├── service-provider.go      # Service provider registration
│   │   └── README.md                # Redis documentation
│   └── http/                        # HTTP client provider (interfaces)
│       ├── interfaces.go            # HTTP client interfaces
│       └── README.md                # HTTP client documentation
└── tests/                           # Integration tests
    ├── integration/                 # Full integration tests
    ├── providers/                   # Provider-specific tests
    └── examples/                    # Example usage tests
```

## 🛠 Installation & Setup

### Prerequisites
- Go 1.19 or later
- PostgreSQL 12+ (optional, for database examples)
- Redis 6+ (optional, for cache examples)

### Dependencies
```bash
go mod tidy
```

Key dependencies:
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/go-redis/redis/v8` - Redis client
- Standard library for HTTP client

### Environment Setup
Copy the example environment file and configure:
```bash
cp config/.env.example .env
```

Edit `.env` with your settings:
```env
# Application
APP_NAME="Service Provider Example"
APP_VERSION="1.0.0"
APP_ENVIRONMENT="development"
APP_DEBUG=true

# Database
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_DATABASE=example_db
APP_DATABASE_USERNAME=postgres
APP_DATABASE_PASSWORD=password
APP_DATABASE_SSL_MODE=disable

# Redis
APP_REDIS_HOST=localhost
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=
APP_REDIS_DATABASE=0

# Logging
APP_LOGGING_LEVEL=info
APP_LOGGING_FORMAT=json
```

## 🚀 Running the Example

### Basic Usage
```bash
go run main.go
```

This will:
1. Initialize all service providers
2. Demonstrate service usage
3. Start health monitoring
4. Wait for shutdown signal (Ctrl+C)

### With Real Services
To see full functionality, start PostgreSQL and Redis:

```bash
# Start PostgreSQL (using Docker)
docker run --name postgres-example -e POSTGRES_PASSWORD=password -e POSTGRES_DB=example_db -p 5432:5432 -d postgres:13

# Start Redis (using Docker)
docker run --name redis-example -p 6379:6379 -d redis:6-alpine

# Run the application
go run main.go
```

### Configuration Options

#### Via Environment Variables
```bash
APP_DATABASE_HOST=192.168.1.100 APP_REDIS_HOST=192.168.1.101 go run main.go
```

#### Via Configuration Files
Edit `config/app.yaml`:
```yaml
app:
  name: "Custom Service Provider Example"
  version: "1.0.0"
  environment: "production"
  debug: false

database:
  host: "production-db.example.com"
  port: 5432
  database: "prod_db"
  username: "prod_user"
  ssl_mode: "require"

redis:
  host: "redis-cluster.example.com"
  port: 6379
  pool_size: 20
```

## 💡 Usage Examples

### Working with Database Service

```go
package main

import (
    "context"
    "log"
    "./modules/postgres"
)

func databaseExample(app ApplicationInterface) {
    container := app.GetContainer()
    
    // Get database service
    dbService, err := container.Make("database")
    if err != nil {
        log.Fatal(err)
    }
    
    db := dbService.(postgres.DatabaseInterface)
    
    // Connect and use
    ctx := context.Background()
    if err := db.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Execute query
    results, err := db.Query(ctx, "SELECT * FROM users WHERE active = $1", true)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use transaction
    err = db.Transaction(ctx, func(tx postgres.TransactionInterface) error {
        if err := tx.Execute("UPDATE users SET last_login = NOW() WHERE id = $1", userID); err != nil {
            return err
        }
        return tx.Execute("INSERT INTO login_logs (user_id, timestamp) VALUES ($1, NOW())", userID)
    })
}
```

### Working with Cache Service

```go
package main

import (
    "context"
    "time"
    "./modules/redis"
)

func cacheExample(app ApplicationInterface) {
    container := app.GetContainer()
    
    // Get cache service
    cacheService, err := container.Make("cache")
    if err != nil {
        log.Fatal(err)
    }
    
    cache := cacheService.(redis.CacheInterface)
    ctx := context.Background()
    
    // Basic operations
    err = cache.Set(ctx, "user:123", "John Doe", 1*time.Hour)
    value, err := cache.Get(ctx, "user:123")
    
    // Hash operations
    err = cache.HashSet(ctx, "user:123:profile", "name", "John Doe")
    err = cache.HashSet(ctx, "user:123:profile", "email", "john@example.com")
    profile, err := cache.HashGetAll(ctx, "user:123:profile")
    
    // List operations
    cache.ListPush(ctx, "notifications:123", "New message", "Friend request")
    notifications, err := cache.ListRange(ctx, "notifications:123", 0, -1)
    
    // Advanced operations
    exists, err := cache.Exists(ctx, "user:123")
    ttl, err := cache.TTL(ctx, "user:123")
    counter, err := cache.Increment(ctx, "page:views", 1)
}
```

### Creating Custom Service Providers

```go
package custom

import (
    "context"
    "../base"
    "../internal/container"
)

type EmailServiceProvider struct {
    *base.BaseProvider
}

func NewEmailServiceProvider() *EmailServiceProvider {
    return &EmailServiceProvider{
        BaseProvider: base.NewBaseProvider("email", true), // Deferred
    }
}

func (p *EmailServiceProvider) Provides() []string {
    return []string{"email", "mailer"}
}

func (p *EmailServiceProvider) Register(c container.ContainerInterface) error {
    p.SetContainer(c)
    
    err := c.Singleton("email", func(c container.ContainerInterface) (interface{}, error) {
        // Get config
        configService, _ := c.Make("config")
        config := configService.(ConfigInterface)
        
        // Create email service based on config
        return NewEmailService(EmailConfig{
            Driver:   config.GetString("mail.driver", "smtp"),
            Host:     config.GetString("mail.host", "localhost"),
            Port:     config.GetInt("mail.port", 587),
            Username: config.GetString("mail.username", ""),
            Password: config.GetString("mail.password", ""),
        }), nil
    })
    
    if err != nil {
        return err
    }
    
    p.SetRegistered(true)
    return nil
}

func (p *EmailServiceProvider) Boot(c container.ContainerInterface) error {
    p.SetBooted(true)
    return nil
}

func (p *EmailServiceProvider) HealthCheck(ctx context.Context) map[string]interface{} {
    // Implement health check logic
    return map[string]interface{}{
        "status": "healthy",
        "smtp_connection": true,
    }
}
```

## 🔧 Configuration Reference

### Database Configuration
```yaml
database:
  host: "localhost"                    # Database host
  port: 5432                          # Database port
  database: "myapp"                   # Database name
  username: "dbuser"                  # Database username
  password: "dbpass"                  # Database password
  ssl_mode: "prefer"                  # SSL mode: disable, require, prefer
  max_open_connections: 25            # Max open connections
  max_idle_connections: 5             # Max idle connections
  connection_lifetime: "5m"           # Connection max lifetime
  connect_timeout: "30s"              # Connection timeout
  query_timeout: "30s"                # Query timeout
```

### Redis Configuration
```yaml
redis:
  host: "localhost"                   # Redis host
  port: 6379                         # Redis port
  password: ""                       # Redis password (optional)
  database: 0                        # Redis database number
  pool_size: 10                      # Connection pool size
  min_idle_connections: 2            # Minimum idle connections
  max_retries: 3                     # Maximum retry attempts
  dial_timeout: "5s"                 # Connection timeout
  read_timeout: "3s"                 # Read timeout
  write_timeout: "3s"                # Write timeout
  pool_timeout: "4s"                 # Pool timeout
  
  # TLS Configuration (optional)
  tls:
    enabled: false
    insecure_skip_verify: false
    cert_file: ""
    key_file: ""
    ca_file: ""
  
  # Sentinel Configuration (optional)
  sentinel:
    enabled: false
    master_name: "mymaster"
    addresses:
      - "sentinel1:26379"
      - "sentinel2:26379"
    password: ""
  
  # Cluster Configuration (optional)
  cluster:
    enabled: false
    addresses:
      - "cluster1:6379"
      - "cluster2:6379"
    max_redirects: 3
    read_only: false
```

## 🧪 Testing

### Run All Tests
```bash
go test ./...
```

### Integration Tests
```bash
go test ./tests/integration/
```

### Provider-Specific Tests
```bash
go test ./modules/postgres/
go test ./modules/redis/
```

### With Test Dependencies
```bash
# Start test databases
docker-compose -f docker-compose.test.yml up -d

# Run tests
go test ./... -tags=integration

# Cleanup
docker-compose -f docker-compose.test.yml down
```

## 🔍 Monitoring & Health Checks

The application includes comprehensive health monitoring:

### Health Check Endpoints
```go
func healthCheckHandler(app ApplicationInterface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
        defer cancel()
        
        providers := app.GetProviders()
        health := make(map[string]interface{})
        
        for _, provider := range providers {
            if hp, ok := provider.(HealthCheckProvider); ok {
                health[provider.Name()] = hp.HealthCheck(ctx)
            }
        }
        
        json.NewEncoder(w).Encode(health)
    }
}
```

### Statistics
Each provider exposes detailed statistics:
- Connection pool metrics
- Query/operation counts
- Error rates and types
- Performance metrics
- Resource utilization

## 🚀 Production Deployment

### Environment Configuration
```bash
# Production environment variables
export APP_ENVIRONMENT=production
export APP_DEBUG=false
export DATABASE_SSL_MODE=require
export REDIS_POOL_SIZE=50
```

### Docker Deployment
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o service-provider-example main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/service-provider-example .
COPY --from=builder /app/config ./config
CMD ["./service-provider-example"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: service-provider-example
spec:
  replicas: 3
  selector:
    matchLabels:
      app: service-provider-example
  template:
    metadata:
      labels:
        app: service-provider-example
    spec:
      containers:
      - name: app
        image: service-provider-example:latest
        env:
        - name: DATABASE_HOST
          value: "postgres.default.svc.cluster.local"
        - name: REDIS_HOST
          value: "redis.default.svc.cluster.local"
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow the established patterns in existing providers
- Add comprehensive tests for new providers
- Update documentation
- Ensure health checks are implemented
- Maintain backward compatibility

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Laravel's service provider pattern
- Built with Go's excellent standard library
- PostgreSQL and Redis communities for excellent documentation
- The Go community for best practices and patterns
