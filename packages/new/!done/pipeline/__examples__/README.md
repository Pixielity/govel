# Pipeline Package Examples

This directory contains comprehensive examples demonstrating real-world usage of the GoVel Pipeline package.

## Example Structure

### Basic Examples
- `example_new_pipeline_example.go` - Basic pipeline creation and usage
- `example_pipeline_through_example.go` - Using the Through method with middleware

### Real-World Examples
- `http_middleware_example.go` - HTTP middleware processing pipeline
- `hub_named_pipelines_example.go` - Using Hub with named pipelines for different order types
- `context_timeout_example.go` - Context usage with timeouts, cancellation, and metadata
- `integration_service_provider_example.go` - Full integration with service provider

## Running Examples

Each example can be run directly:

```bash
# Basic pipeline example
go run __examples__/example_new_pipeline_example.go

# HTTP middleware example
go run __examples__/http_middleware_example.go

# Hub named pipelines example
go run __examples__/hub_named_pipelines_example.go

# Context timeout example  
go run __examples__/context_timeout_example.go

# Service provider integration example
go run __examples__/integration_service_provider_example.go
```

## Example Categories

### 1. HTTP Middleware Pipeline
**File:** `http_middleware_example.go`

Demonstrates:
- Authentication middleware
- Logging middleware  
- Rate limiting middleware
- Error handling in pipelines
- Request/response processing

**Use Case:** Web API request processing

### 2. Hub Named Pipelines
**File:** `hub_named_pipelines_example.go`

Demonstrates:
- Multiple named pipeline configurations
- Order processing workflows (standard, express, digital)
- Different middleware chains per pipeline type
- Error handling per pipeline

**Use Case:** E-commerce order processing system

### 3. Context Management
**File:** `context_timeout_example.go`

Demonstrates:
- Context with timeouts
- Context cancellation
- Metadata storage and retrieval
- Context propagation through pipelines

**Use Case:** Task processing with timeouts and cancellation

### 4. Service Provider Integration
**File:** `integration_service_provider_example.go`

Demonstrates:
- Full dependency injection setup
- Service provider registration
- Container-based service resolution
- Integration with application lifecycle

**Use Case:** Complete application integration

## Key Concepts Demonstrated

### Pipeline Patterns
- **Middleware Chain**: Sequential processing with early termination on errors
- **Russian Doll**: Nested execution where each middleware wraps the next
- **Context Propagation**: Passing metadata and cancellation signals through the pipeline

### Hub Patterns
- **Named Pipelines**: Different processing workflows for different scenarios
- **Default Pipeline**: Fallback processing when no specific pipeline is named
- **Dynamic Pipeline**: Runtime pipeline selection based on input data

### Service Integration
- **Dependency Injection**: Using containers to resolve pipeline dependencies
- **Service Provider**: Laravel-style service registration and bootstrapping
- **Factory Pattern**: Creating pipeline instances through factories

## Common Use Cases

### 1. HTTP Request Processing
```go
// Authentication -> Validation -> Rate Limiting -> Handler
result, err := pipeline.
    Send(request).
    Through([]interface{}{auth, validate, rateLimit}).
    Then(handler)
```

### 2. Data Processing Pipeline
```go
// Parse -> Transform -> Validate -> Store
result, err := pipeline.
    Send(data).
    Through([]interface{}{parser, transformer, validator}).
    Then(storer)
```

### 3. Order Processing
```go
// Inventory -> Payment -> Shipping -> Notification
result, err := hub.Pipe(order, "standard")
```

### 4. Background Task Processing
```go
// With context timeout and cancellation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
result, err := pipeline.
    Send(task).
    WithContext(ctx).
    Through(processors).
    Then(finalizer)
```

## Error Handling Patterns

### 1. Graceful Degradation
```go
result, err := pipeline.Then(func(passable interface{}) interface{} {
    // Return partial result on error
    if someCondition {
        return PartialResult{Data: passable, Error: "warning"}
    }
    return FullResult{Data: processedData}
})
```

### 2. Error Wrapping
```go
middleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
    result, err := next(passable)
    if err != nil {
        return nil, fmt.Errorf("middleware failed: %w", err)
    }
    return result, nil
}
```

### 3. Circuit Breaker
```go
middleware := func(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
    if circuitBreaker.IsOpen() {
        return nil, errors.New("circuit breaker is open")
    }
    return next(passable)
}
```

## Performance Considerations

### Thread Safety
All examples demonstrate thread-safe usage:
- Pipeline instances are safe for concurrent reads
- Hub instances handle concurrent pipeline registration
- Context instances are immutable after creation

### Memory Management
- Pipelines don't retain references to processed data
- Middleware should avoid memory leaks
- Context metadata is cleaned up automatically

### Scalability
- Pipeline creation is lightweight
- Hub can manage thousands of named pipelines
- Context propagation has minimal overhead

## Contributing

When adding new examples:

1. Focus on real-world use cases
2. Include comprehensive error handling
3. Demonstrate best practices
4. Add detailed comments explaining the concepts
5. Show both success and failure scenarios
6. Include performance considerations
7. Document the use case and target audience
