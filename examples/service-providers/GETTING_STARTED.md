# Getting Started with Service Provider Examples

This guide will help you quickly get up and running with the complete service provider example, demonstrating how to build modular, maintainable Go applications using the service provider pattern.

## 🚀 Quick Start

### 1. Start Dependencies
```bash
# Start PostgreSQL and Redis using Docker
docker-compose up -d postgres redis

# Wait for services to be ready (optional - they have health checks)
docker-compose ps
```

### 2. Run the Application
```bash
# Run with default settings
go run main.go

# Or with custom environment variables
APP_DATABASE_HOST=localhost APP_REDIS_HOST=localhost go run main.go
```

### 3. What You'll See
The application will:
1. ✅ Initialize configuration from environment variables
2. ✅ Register core service providers (logger)  
3. ✅ Register feature service providers (PostgreSQL, Redis)
4. ✅ Boot all providers and resolve dependencies
5. ✅ Demonstrate service usage with mock data
6. ✅ Start health monitoring (every 30 seconds)
7. 🔄 Run until you press Ctrl+C
8. ✅ Gracefully shutdown all services

## 📋 Expected Output

```
=== Complete Service Provider Example ===

--- Setting up Configuration ---
Configuration loaded successfully

--- Registering Core Service Providers ---
Core providers registered successfully

--- Registering Feature Service Providers ---
Feature providers registered successfully

--- Demonstrating Service Usage ---

-- Configuration Service Demo --
App: Service Provider Example v1.0.0 (Debug: true)

-- Logging Service Demo --
[INFO] Logging service is working {"action":"test","component":"demo"}
[ERROR] This is a test error message {"component":"demo","error":"simulated error"}
Logging demonstration completed

-- Database Service Demo --
Database connection status: false
Database connection failed (expected): failed to ping database: dial tcp [::1]:5432: connect: connection refused

-- Cache Service Demo --
Cache ping failed (expected): dial tcp [::1]:6379: connect: connection refused

Service demonstrations completed

--- Starting Health Monitoring ---

--- Application Running (Press Ctrl+C to exit) ---
[INFO] Health check completed {"results":{"logger":null,"postgresql":{"connections":{},"healthy":false,"status":"not_initialized"},"redis":{"connections":{},"healthy":false,"status":"not_initialized"}}}

^C
--- Shutting down application ---
=== Application shutdown complete ===
```

## 🔧 With Real Services

If you start the actual PostgreSQL and Redis services, you'll see much richer output:

```bash
# Start all services
docker-compose up -d

# Run the application
go run main.go
```

Expected enhanced output:
```
-- Database Service Demo --
Database connection status: true
Database stats: {OpenConnections:1 InUseConnections:0 IdleConnections:1 ...}
Query results: [map[version:PostgreSQL 13.8 on x86_64-pc-linux-musl...]]

-- Cache Service Demo --
Cache is connected, demonstrating operations...
Cache value: Hello, World!
Key exists: true
Cache stats: {TotalConns:1 IdleConns:1 TotalCmds:3 Hits:1 ...}
```

## 🏗 Architecture Overview

### Service Provider Flow
```
1. Application Creation
   └── Container initialized
   └── Configuration loaded

2. Provider Registration
   ├── Core Providers (Logger)
   ├── Feature Providers (PostgreSQL, Redis)
   └── Services registered in container

3. Provider Booting
   ├── Dependencies resolved
   ├── Services initialized
   └── Health checks enabled

4. Application Runtime
   ├── Services available via DI container
   ├── Health monitoring active
   └── Graceful shutdown handling
```

### Dependency Graph
```
Application
├── Container (DI)
├── Configuration Service
├── Logger Service
├── PostgreSQL Service
│   ├── Database Interface
│   ├── Connection Manager
│   └── Transaction Support
└── Redis Service
    ├── Cache Interface
    ├── Connection Manager
    └── Advanced Operations
```

## 🧪 Testing Different Scenarios

### 1. Without External Services (Default)
```bash
go run main.go
```
- Shows graceful handling of missing services
- Demonstrates configuration and logging
- Health checks show service unavailability

### 2. With PostgreSQL Only
```bash
docker-compose up -d postgres
go run main.go
```
- Database operations work fully
- Redis operations fail gracefully
- Mixed health check results

### 3. With Redis Only
```bash
docker-compose up -d redis
go run main.go
```
- Cache operations work fully
- Database operations fail gracefully
- Mixed health check results

### 4. With All Services
```bash
docker-compose up -d
go run main.go
```
- Full functionality demonstrated
- All health checks pass
- Complete service integration

### 5. Custom Configuration
```bash
# Via environment variables
APP_DATABASE_HOST=custom-db.com APP_REDIS_HOST=custom-redis.com go run main.go

# Via .env file
cp config/.env.example .env
# Edit .env with your settings
go run main.go
```

## 🔍 Key Concepts Demonstrated

### 1. **Service Provider Pattern**
- Modular service registration
- Dependency injection
- Lifecycle management (register → boot → run → terminate)

### 2. **Deferred Loading**
- Services only initialize when first accessed
- Improved startup performance
- Better resource utilization

### 3. **Configuration Management**
- Environment variable loading
- Default value fallbacks
- Type-safe configuration access

### 4. **Health Monitoring**
- Automated health checks
- Service status reporting
- Performance metrics collection

### 5. **Graceful Shutdown**
- Signal handling (SIGINT, SIGTERM)
- Resource cleanup
- Timeout-based shutdown

### 6. **Error Handling**
- Graceful degradation
- Service isolation
- Comprehensive error reporting

## 🛠 Customization Examples

### Adding a New Service Provider

1. **Create the interface:**
```go
type EmailInterface interface {
    Send(to, subject, body string) error
    SendTemplate(to, template string, data map[string]interface{}) error
}
```

2. **Implement the service:**
```go
type EmailService struct {
    driver string
    config EmailConfig
}

func (e *EmailService) Send(to, subject, body string) error {
    // Implementation
}
```

3. **Create the service provider:**
```go
type EmailServiceProvider struct {
    *base.BaseProvider
}

func (p *EmailServiceProvider) Register(c container.ContainerInterface) error {
    c.Singleton("email", func(c container.ContainerInterface) (interface{}, error) {
        return NewEmailService(/* config */), nil
    })
}
```

4. **Register with application:**
```go
emailProvider := NewEmailServiceProvider()
app.RegisterProvider(emailProvider)
```

### Custom Middleware

```go
type LoggingMiddleware struct {
    logger LoggerInterface
}

func (m *LoggingMiddleware) Process(req *http.Request, next MiddlewareFunc) (*Response, error) {
    start := time.Now()
    resp, err := next(req)
    duration := time.Since(start)
    
    m.logger.Info("HTTP Request", map[string]interface{}{
        "method": req.Method,
        "url": req.URL.String(),
        "duration": duration,
        "status": resp.StatusCode,
    })
    
    return resp, err
}
```

## 📊 Monitoring and Observability

The example includes comprehensive monitoring:

### Health Checks
- Automated every 30 seconds
- Per-provider status reporting
- Connection and service validation

### Metrics Collection
- Request/response statistics  
- Connection pool metrics
- Error rates and types
- Performance timings

### Logging Integration
- Structured logging with fields
- Multiple log levels
- Request/response logging
- Error tracking

## 🚢 Production Readiness

This example demonstrates production-ready patterns:

### ✅ Configuration Management
- Environment-based configuration
- Secure credential handling
- Default value management

### ✅ Error Handling
- Graceful degradation
- Comprehensive error reporting
- Recovery mechanisms

### ✅ Resource Management  
- Connection pooling
- Resource cleanup
- Memory management

### ✅ Monitoring
- Health checks
- Metrics collection
- Performance tracking

### ✅ Scalability
- Stateless design
- Connection sharing
- Horizontal scaling ready

## 🎯 Next Steps

1. **Explore the Code**: Browse the `modules/` directory to understand each service provider
2. **Run Tests**: Execute `go test ./...` to see comprehensive testing
3. **Add Services**: Create your own service providers following the established patterns  
4. **Deploy**: Use the provided Docker and Kubernetes configurations
5. **Monitor**: Implement the health check endpoints in your HTTP server
6. **Scale**: Adapt the patterns for your specific use cases

## 📚 Further Reading

- [PostgreSQL Provider Documentation](modules/postgres/README.md)
- [Redis Provider Documentation](modules/redis/README.md)  
- [Base Provider Documentation](modules/base/README.md)
- [Complete API Documentation](README.md)

---

Happy coding! 🎉 This service provider pattern will help you build maintainable, testable, and scalable Go applications.
