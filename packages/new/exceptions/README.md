# GoVel Exception System

A comprehensive, ISP-compliant exception handling system for GoVel applications, inspired by Laravel's exception handling with Go-specific optimizations and comprehensive solution support.

## 🏗️ Architecture Overview

The GoVel exception system follows the **Interface Segregation Principle (ISP)** and **Separation of Concerns**, providing a clean, modular architecture that's easy to extend and maintain.

### File Structure

```
packages/exceptions/
├── README.md                        # This file
├── SOLUTIONS.md                     # Solution system documentation
├── exceptions.go                    # Main package with backward compatibility exports
├── interfaces/                     # ISP-compliant interface definitions
│   ├── core_interfaces.go          # Main ExceptionInterface composition
│   ├── contextable.go              # Context management interface
│   ├── renderable.go               # Rendering interface
│   ├── httpable.go                 # HTTP functionality interface
│   ├── solutionable.go             # Solution-related interface
│   ├── stackable.go                # Stack trace interface
│   └── solution/                   # Solution-specific interfaces
│       ├── solution.go             # Basic Solution interface
│       ├── provides_solution.go    # ProvidesSolution interface
│       ├── runnable_solution.go    # RunnableSolution interface
│       └── solution_provider.go    # Solution provider interfaces
├── core/                           # Core implementations
│   ├── exception.go                # Base Exception struct (ISP-compliant)
│   └── solution/                   # Solution core implementations
│       ├── base_solution.go        # BaseSolution implementation
│       └── solution_repository.go  # Solution provider repository
├── http/                           # HTTP exception types
│   ├── not_found.go               # 404 Not Found exception
│   └── [other HTTP exceptions...]  # Additional HTTP exceptions
├── solutions/                      # Solution implementations
│   ├── http/                       # HTTP-specific solutions
│   │   ├── not_found_solution.go   # 404 solution implementation
│   │   └── [other solutions...]    # Additional HTTP solutions
│   ├── runnable/                   # Runnable solution implementations
│   │   └── [runnable solutions...]
│   └── providers/                  # Solution provider implementations
├── helpers/                        # Utility functions
│   ├── abort.go                   # Abort, AbortIf, AbortUnless functions
│   └── shortcuts.go               # Shortcut functions (Abort400, etc.)
├── tests/                         # Comprehensive tests
│   ├── core_test.go               # Core functionality tests
│   ├── compatibility_test.go      # Backward compatibility tests
│   └── [other test files...]
└── examples/                      # Usage examples
    ├── basic_example.go
    └── isp_exceptions_example/
        └── main.go                # ISP usage demonstration
```

## 🎯 Key Features

### ✅ **Interface Segregation Principle (ISP) Compliance**
- **Small, focused interfaces** for specific functionality
- **Compose larger interfaces** from smaller ones
- **Implement only what you need** - no forced dependencies

### ✅ **Backward Compatibility**
- **Existing code works unchanged** through package-level exports
- **Same API** as before, with new capabilities added transparently

### ✅ **Comprehensive Solution System**
- **Automatic solutions** for HTTP exceptions
- **Runnable solutions** for common development issues
- **Custom solution providers** for domain-specific problems

### ✅ **Clean Architecture**
- **Separation of concerns** with dedicated packages
- **Clear dependencies** between packages
- **Easy to extend** and modify

## 🚀 Quick Start

### Basic Usage (Backward Compatible)

```go
package main

import "govel/packages/exceptions"

func main() {
    // All existing code works exactly the same
    exc := exceptions.NewException("Something went wrong", 500)
    notFound := exceptions.NewNotFoundException("User not found")
    aborted := exceptions.Abort404("Page not found")
    
    // New solution functionality works automatically
    if notFound.HasSolution() {
        solution := notFound.GetSolution()
        fmt.Println("Solution:", solution.GetSolutionTitle())
    }
}
```

### ISP Interface Usage

```go
package main

import (
    "govel/packages/exceptions/core"
    "govel/packages/exceptions/interfaces"
)

func handleHTTPError(exc interfaces.HTTPable) {
    fmt.Printf("Status: %d, Message: %s\n", 
        exc.GetStatusCode(), exc.GetMessage())
    
    exc.WithHeader("X-Error-ID", "12345")
}

func addContext(exc interfaces.Contextable) {
    exc.WithContext("user_id", 123)
    exc.WithContext("action", "update_profile")
}

func renderError(exc interfaces.Renderable) map[string]interface{} {
    return exc.Render()
}

func main() {
    exc := core.NewException("Test error", 400)
    
    // Use specific interfaces for specific functionality
    handleHTTPError(exc)   // Only HTTP functionality
    addContext(exc)        // Only context functionality
    response := renderError(exc) // Only rendering functionality
}
```

### Direct Package Usage

```go
package main

import (
    "govel/packages/exceptions/core"
    "govel/packages/exceptions/core/solution"
    httpExceptions "govel/packages/exceptions/http"
    "govel/packages/exceptions/helpers"
)

func main() {
    // Use packages directly for specific functionality
    coreExc := core.NewException("Core error", 500)
    httpExc := httpExceptions.NewNotFoundException("HTTP error")
    helperExc := helpers.Abort(422, "Helper error")
    
    // Create custom solutions
    sol := solution.NewBaseSolution("Custom Solution").
        SetSolutionDescription("Fix this specific issue").
        AddDocumentationLink("Docs", "https://govel.dev/docs")
    
    coreExc.SetSolution(sol)
}
```

### Custom Implementation

```go
package main

import (
    "govel/packages/exceptions/interfaces"
    solutionInterface "govel/packages/exceptions/interfaces/solution"
)

// Custom exception implementing only needed interfaces
type CustomException struct {
    message string
    code    int
}

func (e *CustomException) Error() string {
    return e.message
}

func (e *CustomException) GetStatusCode() int {
    return e.code
}

func (e *CustomException) GetMessage() string {
    return e.message
}

// Implement only the methods you need from HTTPable interface
// ... (other required HTTPable methods)

// Custom solution
type CustomSolution struct {
    title string
}

func (s *CustomSolution) GetSolutionTitle() string {
    return s.title
}

func (s *CustomSolution) GetSolutionDescription() string {
    return "Custom solution description"
}

func (s *CustomSolution) GetDocumentationLinks() map[string]string {
    return map[string]string{
        "Custom Docs": "https://example.com/docs",
    }
}

// Ensure interface compliance
var _ interfaces.HTTPable = (*CustomException)(nil)
var _ solutionInterface.Solution = (*CustomSolution)(nil)
```

## 📋 Interface Reference

### Core Interfaces

| Interface | Purpose | Key Methods |
|-----------|---------|-------------|
| `ExceptionInterface` | Main composition interface | Combines all other interfaces |
| `HTTPable` | HTTP functionality | `GetStatusCode()`, `GetMessage()`, `WithHeader()` |
| `Contextable` | Context management | `GetContext()`, `WithContext()`, `SetContext()` |
| `Renderable` | Response rendering | `Render()` |
| `Stackable` | Stack trace handling | `GetStackTrace()` |
| `Solutionable` | Solution support | `GetSolution()`, `HasSolution()`, `SetSolution()` |

### Solution Interfaces

| Interface | Purpose | Key Methods |
|-----------|---------|-------------|
| `Solution` | Basic solution info | `GetSolutionTitle()`, `GetSolutionDescription()` |
| `ProvidesSolution` | Exception provides solution | `GetSolution()` |
| `RunnableSolution` | Executable solutions | `Run()`, `GetRunButtonText()` |
| `HasSolutionsForThrowable` | Solution provider | `CanSolve()`, `GetSolutions()` |

## 🎨 Benefits of the ISP Design

### 1. **Flexible Implementation**
```go
// Only implement what you need
type SimpleHTTPError struct {
    message string
    code    int
}

// Implement only HTTPable interface
func (e *SimpleHTTPError) GetStatusCode() int { return e.code }
func (e *SimpleHTTPError) GetMessage() string { return e.message }
// ... other HTTPable methods
```

### 2. **Clear Dependencies**
```go
// Functions depend only on what they use
func logHTTPError(err interfaces.HTTPable) {
    log.Printf("HTTP %d: %s", err.GetStatusCode(), err.GetMessage())
}

func addMetadata(err interfaces.Contextable) {
    err.WithContext("logged_at", time.Now())
}
```

### 3. **Easy Testing**
```go
// Mock only the interface you need
type mockHTTPable struct {
    code    int
    message string
}

func (m *mockHTTPable) GetStatusCode() int { return m.code }
func (m *mockHTTPable) GetMessage() string { return m.message }
// ... implement only needed methods

func TestHTTPHandler(t *testing.T) {
    mock := &mockHTTPable{code: 404, message: "Not found"}
    handleHTTPError(mock) // Works with any HTTPable implementation
}
```

### 4. **Extensible Architecture**
```go
// Add new interfaces without breaking existing code
type Loggable interface {
    GetLogLevel() string
    GetLogMessage() string
}

// Extend existing exceptions
type LoggableException struct {
    *core.Exception
    logLevel string
}

func (e *LoggableException) GetLogLevel() string {
    return e.logLevel
}
```

## 🔄 Migration from Old Structure

The new structure is **100% backward compatible**. Existing code continues to work without changes:

### Before (still works)
```go
import "govel/packages/exceptions"

exc := exceptions.NewException("Error", 500)
notFound := exceptions.NewNotFoundException("Not found")
aborted := exceptions.Abort404("Page not found")
```

### After (optional, for new features)
```go
// Use specific packages for specific needs
import (
    "govel/packages/exceptions/core"
    "govel/packages/exceptions/interfaces"
    httpExceptions "govel/packages/exceptions/http"
)

exc := core.NewException("Error", 500)
notFound := httpExceptions.NewNotFoundException("Not found")

// Use interfaces for parameters
func handleError(err interfaces.HTTPable) {
    // Handle any HTTPable implementation
}
```

## 📚 Further Reading

- [SOLUTIONS.md](SOLUTIONS.md) - Comprehensive solution system documentation
- [examples/isp_exceptions_example/main.go](examples/isp_exceptions_example/main.go) - Complete ISP usage example
- [Interface Segregation Principle](https://en.wikipedia.org/wiki/Interface_segregation_principle) - ISP explanation

## 🧪 Testing

The system includes comprehensive tests for all components:

```bash
# Test all components
go test ./packages/exceptions/tests/... -v

# Test specific components
go test ./packages/exceptions/tests/core_test.go -v
go test ./packages/exceptions/tests/compatibility_test.go -v
```

---

**The GoVel Exception System: Clean architecture meets practical functionality.** 🚀
