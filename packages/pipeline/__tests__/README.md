# Pipeline Package Tests

This directory contains comprehensive test cases for the GoVel Pipeline package.

## Test Structure

Each major function has its own test file:

### Pipeline Tests

- `new_pipeline_test.go` - Tests for `NewPipeline()` constructor
- `pipeline_send_test.go` - Tests for `Send()` method
- `pipeline_through_test.go` - Tests for `Through()` method
- `pipeline_then_test.go` - Tests for `Then()` method

### Hub Tests  

- `hub_new_hub_test.go` - Tests for `NewHub()` constructor
- `hub_pipeline_test.go` - Tests for `Pipeline()` method

### Context Tests

- `context_new_test.go` - Tests for `NewPipelineContext()` constructor

### Service Provider Tests

- `service_provider_test.go` - Tests for service provider functionality

## Running Tests

To run all tests:

```bash
go test ./...
```

To run tests with verbose output:

```bash
go test -v ./...
```

To run tests for a specific file:

```bash
go test -v -run TestNewPipeline ./...
```

To run tests with race condition detection:

```bash
go test -race ./...
```

## Test Coverage

To generate test coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Categories

### Unit Tests

- Test individual functions in isolation
- Mock dependencies using MockContainer
- Verify expected behavior and error conditions

### Integration Tests  

- Test interaction between components
- Verify thread safety
- Test complex scenarios

### Performance Tests

- Concurrent access tests
- Load testing with multiple goroutines
- Memory leak detection

## Mock Objects

The tests use several mock implementations:

- `MockContainer` - Simple container implementation for testing
- `MockMiddleware` - Middleware that can simulate success/failure
- `MockApplicationInterface` - Application interface for service provider tests

## Test Patterns

All tests follow Go testing best practices:

- Use table-driven tests where appropriate
- Test both success and failure cases
- Include edge cases and boundary conditions  
- Verify thread safety with concurrent tests
- Use descriptive test names
- Include setup/teardown as needed

## Contributing

When adding new tests:

1. Create a separate test file for each major function
2. Use the existing mock objects or create new ones as needed
3. Follow the naming convention: `TestFunctionName`
4. Include both positive and negative test cases
5. Add concurrent tests for thread-safety verification
6. Document complex test scenarios
