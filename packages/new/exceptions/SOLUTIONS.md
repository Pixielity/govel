# GoVel Exception Solutions System

The GoVel exception solutions system is inspired by Laravel's Spatie Error Solutions package and provides helpful guidance and automated fixes for common errors in GoVel applications.

## Features

### 🎯 **Automatic Solutions for HTTP Exceptions**
- All HTTP exceptions (400, 401, 403, 404, 405, 422, 429, 500, 503) automatically include helpful solutions
- Solutions provide contextual guidance, documentation links, and troubleshooting steps
- Rich exception rendering includes solution information in JSON responses

### 🔧 **Runnable Solutions**
- Solutions that can automatically fix common development issues
- Generate application keys, create missing directories, install dependencies, fix permissions
- Safe execution with proper error handling and rollback capabilities

### 🏗️ **Solution Providers**
- Extensible provider system for adding custom solutions to any error type
- HTTP Exception Provider for common HTTP status code errors
- Common Runnable Solutions Provider for development environment issues
- Easy to create custom providers for domain-specific errors

### 📝 **Rich Documentation Integration**
- Solutions include relevant documentation links
- Contextual help with step-by-step instructions
- Links to GoVel documentation, MDN references, and best practices

## Quick Start

### Using Built-in HTTP Exception Solutions

```go
// HTTP exceptions automatically include solutions
exc := exceptions.NewNotFoundException("User not found")

if exc.HasSolution() {
    solution := exc.GetSolution()
    fmt.Println("Solution:", solution.GetSolutionTitle())
    fmt.Println("Description:", solution.GetSolutionDescription())
    
    // Get documentation links
    for name, url := range solution.GetDocumentationLinks() {
        fmt.Printf("📖 %s: %s\n", name, url)
    }
}
```

### Creating Custom Solutions

```go
// Create a custom solution
solution := solutions.NewBaseSolution("Fix Database Connection").
    SetSolutionDescription("Check your database credentials and ensure the database server is running.").
    AddDocumentationLink("Database Configuration", "https://govel.dev/docs/database").
    AddDocumentationLink("Connection Troubleshooting", "https://govel.dev/docs/db-troubleshooting")

// Attach to any exception
exc := exceptions.NewException("Database connection failed", 500)
exc.SetSolution(solution)
```

### Using Runnable Solutions

```go
// Create a runnable solution that can fix issues automatically
runnableSolution := solutions.NewGenerateAppKeySolution()

fmt.Println("Action:", runnableSolution.GetSolutionActionDescription())
fmt.Println("Button:", runnableSolution.GetRunButtonText())

// Execute the solution
err := runnableSolution.Run(map[string]interface{}{})
if err != nil {
    fmt.Printf("Failed to run solution: %v\n", err)
} else {
    fmt.Println("Solution executed successfully!")
}
```

### Using Solution Providers

```go
// Set up solution provider repository
repo := solutions.NewSolutionProviderRepository()

// Register built-in providers
repo.RegisterSolutionProvider(solutions.NewHTTPExceptionProvider())
repo.RegisterSolutionProvider(solutions.NewCommonRunnableSolutionsProvider())

// Find solutions for any error
err := fmt.Errorf("404 not found - resource does not exist")
solutionsFound := repo.GetSolutionsForError(err)

for _, solution := range solutionsFound {
    fmt.Printf("💡 %s\n", solution.GetSolutionTitle())
    fmt.Printf("   %s\n", solution.GetSolutionDescription())
}
```

### Creating Custom Solution Providers

```go
type CustomSolutionProvider struct{}

func (p *CustomSolutionProvider) CanSolve(err error) bool {
    return strings.Contains(strings.ToLower(err.Error()), "database")
}

func (p *CustomSolutionProvider) GetSolutions(err error) []solutions.Solution {
    return []solutions.Solution{
        solutions.NewBaseSolution("Database Issue Detected").
            SetSolutionDescription("This appears to be a database-related error. Check your connection settings and database status.").
            AddDocumentationLink("Database Docs", "https://govel.dev/docs/database"),
    }
}

// Register the custom provider
repo.RegisterSolutionProvider(&CustomSolutionProvider{})
```

## Built-in Solutions

### HTTP Exception Solutions

| Status Code | Solution Title | Key Features |
|-------------|----------------|--------------|
| 400 | Bad Request | Request syntax validation, JSON format checking |
| 401 | Authentication Required | Token validation, session management guidance |
| 403 | Access Forbidden | Permission troubleshooting, authorization guides |
| 404 | Resource Not Found | URL validation, routing configuration help |
| 405 | HTTP Method Not Allowed | Method verification, route definition guidance |
| 422 | Request Validation Failed | Validation error details, field requirement help |
| 429 | Rate Limit Exceeded | Rate limiting guidance, retry strategies |
| 500 | Internal Server Error | Server troubleshooting, logging guidance |
| 503 | Service Temporarily Unavailable | Maintenance mode, service health checks |

### Runnable Solutions

| Solution | Description | Use Case |
|----------|-------------|-----------|
| GenerateAppKeySolution | Creates new APP_KEY in .env file | Missing encryption key errors |
| CreateDirectorySolution | Creates missing directories with proper permissions | Storage/log directory not found |
| InstallDependencySolution | Runs commands to install missing dependencies | Missing Go modules, npm packages |
| FixPermissionsSolution | Sets correct file/directory permissions | Permission denied errors |

## Exception Rendering with Solutions

Exceptions now include solution information in their JSON rendering:

```json
{
  "error": true,
  "status_code": 404,
  "message": "User not found",
  "timestamp": "2025-09-11T01:27:11+04:00",
  "solution": {
    "title": "Resource Not Found",
    "description": "The requested resource could not be found. This typically happens when:\n\n• The URL path is incorrect\n• The resource has been moved or deleted\n• The route is not properly defined\n• The resource requires authentication",
    "runnable": false,
    "links": {
      "HTTP 404 Reference": "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404",
      "GoVel Routing Documentation": "https://govel.dev/docs/routing",
      "GoVel Resource Controllers": "https://govel.dev/docs/controllers"
    }
  }
}
```

## Architecture

### Core Interfaces

- **`Solution`** - Basic solution interface with title, description, and documentation links
- **`ProvidesSolution`** - Interface for exceptions that provide their own solutions
- **`HasSolutionsForThrowable`** - Interface for solution providers
- **`RunnableSolution`** - Extension for solutions that can be executed automatically
- **`SolutionProvider`** - Repository interface for managing solution providers

### Key Components

- **`BaseSolution`** - Base implementation with method chaining support
- **`SolutionProviderRepository`** - Central registry for solution providers
- **HTTP Exception Solutions** - Specific solutions for each HTTP status code
- **Runnable Solutions** - Automated fixes for common development issues
- **Solution Providers** - Pluggable system for domain-specific solutions

## Testing

The solution system includes comprehensive tests covering:

- ✅ Base solution functionality and method chaining
- ✅ Solution provider repository management
- ✅ HTTP exception solutions integration
- ✅ Runnable solution execution (mocked)
- ✅ Exception rendering with solution data
- ✅ Custom solution providers and exceptions
- ✅ Multi-provider solution aggregation

Run tests with:
```bash
go test ./packages/exceptions/tests/... -v
```

## Examples

See `/examples/solutions_example/main.go` for a comprehensive demonstration of all solution system features including:

1. **Basic HTTP exceptions with built-in solutions**
2. **Custom exceptions with custom solutions**  
3. **Solution providers in action**
4. **Runnable solutions demonstration**
5. **Exception rendering with solution data**

Run the example:
```bash
go run examples/solutions_example/main.go
```

## Benefits

### For Developers
- **Faster debugging** with contextual error guidance
- **Automated fixes** for common setup and configuration issues
- **Learning tool** with integrated documentation links
- **Consistent error handling** across the application

### For Users
- **Better error messages** with actionable advice
- **Self-service troubleshooting** with guided solutions
- **Reduced support requests** through clear guidance
- **Improved developer experience** with helpful error pages

### For Teams
- **Knowledge sharing** through documented solutions
- **Standardized error handling** patterns
- **Extensible architecture** for domain-specific errors
- **Reduced onboarding time** for new developers

## Inspiration

This solution system is heavily inspired by Laravel's [Spatie Error Solutions](https://github.com/spatie/laravel-error-solutions) package, adapting the concepts to Go's type system and GoVel's architecture while maintaining the core philosophy of helpful, actionable error guidance.

---

**The GoVel Exception Solutions System transforms errors from obstacles into opportunities for learning and improvement.** 🚀
