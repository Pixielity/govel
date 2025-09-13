# GoVel Health Check Package

Monitor the health of your GoVel application by registering configurable health checks. This package provides a comprehensive health monitoring system with support for multiple check types, HTTP endpoints, notifications, and flexible result storage.

## Features

- **Registry Pattern**: Central registration and management of health checks
- **Multiple Check Types**: Built-in support for database, Redis, filesystem, memory, and custom checks
- **HTTP Endpoints**: RESTful endpoints for health status monitoring
- **Rich Results**: Detailed health check results with metadata, timing, and error information
- **Notifications**: Built-in notification channels for failed health checks
- **Background Jobs**: Queue health checks for asynchronous execution
- **Flexible Storage**: Multiple result storage backends (cache, database, in-memory)
- **Configurable**: Extensive configuration options for timeouts, thresholds, and scheduling

## Installation

Add the health check package to your GoVel application:

```go
import "govel/packages/healthcheck/src"
```

## Quick Start

### Basic Usage

```go
package main

import (
    "govel/packages/healthcheck/src"
    "govel/healthcheck/checks/checks"
)

func main() {
    // Create health registry
    health := healthcheck.New()
    
    // Register health checks
    health.Checks([]healthcheck.CheckInterface{
        checks.NewPingCheck().
            Name("external-api").
            URL("https://api.example.com").
            Timeout(5 * time.Second),
            
        checks.NewUsedDiskSpaceCheck().
            Name("disk-space").
            Path("/").
            WarnWhenUsedSpaceIsAbovePercentage(70).
            FailWhenUsedSpaceIsAbovePercentage(90),
    })
}
```

### HTTP Endpoints

The package provides several HTTP endpoints:

- `GET /health` - HTML dashboard with all health check results
- `GET /health.json` - JSON API with all health check results  
- `GET /health/simple` - Simple text response (OK/FAILED)

### Custom Health Checks

Create your own health checks by implementing the `CheckInterface`:

```go
type CustomCheck struct {
    // ... your fields
}

func (c *CustomCheck) Run(ctx context.Context) *healthcheck.Result {
    // Your health check logic here
    return healthcheck.NewResult().OK("Everything is working fine")
}

func (c *CustomCheck) GetName() string {
    return "custom-check"
}
```

## Examples

See the `__examples__/` directory for comprehensive usage examples:

- `basic_health_check_example.go` - Basic usage and setup
- `custom_health_checker_example.go` - Creating custom health checks
- `database_health_check_example.go` - Database connectivity checks
- `redis_health_check_example.go` - Redis connectivity checks

## Configuration

Health checks can be configured through the GoVel config system:

```go
config.Set("healthcheck.default_timeout", "30s")
config.Set("healthcheck.result_store", "cache")
config.Set("healthcheck.notifications.enabled", true)
```

## Available Health Checks

### Built-in Checks

- **PingCheck**: HTTP/HTTPS endpoint availability
- **UsedDiskSpaceCheck**: Disk space usage monitoring
- **MemoryUsageCheck**: Memory usage monitoring  
- **EnvironmentCheck**: Environment variable validation

### Framework Integration Checks

- **DatabaseCheck**: Database connectivity (integrates with GoVel database)
- **RedisCheck**: Redis connectivity (integrates with GoVel cache)
- **CacheCheck**: Cache system health

## Result Storage

Choose from multiple result storage backends:

- **InMemoryStore**: Fast, temporary storage (default)
- **CacheStore**: Persistent cache-based storage
- **DatabaseStore**: Database-backed storage for long-term retention

## Notifications

Get notified when health checks fail:

- **Slack**: Send notifications to Slack channels
- **Email**: Send email notifications
- **Custom**: Implement your own notification channels

## Testing

The package includes testing utilities for easy mocking:

```go
import "govel/healthcheck/testing"

func TestMyApp(t *testing.T) {
    fakeCheck := testing.NewFakeCheck("test-check")
    fakeCheck.SetResult(healthcheck.NewResult().Failed("Simulated failure"))
    
    // Use in your tests
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This package is open-sourced software licensed under the [MIT license](LICENSE).

## Credits

Inspired by [Spatie's Laravel Health](https://github.com/spatie/laravel-health) package, adapted for the Go ecosystem and GoVel framework.
