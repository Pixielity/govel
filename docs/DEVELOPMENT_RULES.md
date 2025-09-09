# GoVel Development Rules

## **📁 Project Structure Rules**

### **Rule 1: ALL Code in src/ Folder**

- **NEVER** place main.go, worker.go, or any Go code in project root
- **ALL** application code must live in `src/` directory
- Only configuration files (README.md, .env, docker-compose.yml) in root

```
✅ CORRECT:
project/
├── README.md
├── go.mod
└── src/
    ├── main.go              # Entry point in src/
    ├── jobs/
    │   └── worker.go        # Background jobs in src/jobs/
    ├── entities/
    └── services/

❌ WRONG:
project/
├── main.go                 # Should be in src/
├── worker.go               # Should be in src/jobs/
└── src/
    └── services/
```

### **Rule 2: ALL Folder Names are PLURAL**

- **ALWAYS** use plural folder names for consistency
- This includes: `entities/`, `services/`, `controllers/`, `middlewares/`, `repositories/`, `validators/`, `factories/`, `configs/`, `databases/`, `interfaces/`, `constants/`, `exceptions/`, `enums/`

```
✅ CORRECT: entities/, services/, controllers/, middlewares/
❌ WRONG: entity/, service/, controller/, middleware/
```

### **Rule 3: types/ Folder - STRUCT DEFINITIONS ONLY**

- **types/** folder contains **ONLY** struct definitions
- **NO** methods, functions, or business logic allowed in types/
- Use for DTOs, request/response types, configuration structs

```go
// ✅ CORRECT: types/user_types.go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ❌ WRONG: types/user_types.go - NO methods allowed
func (u *User) Validate() error {  // ❌ NO METHODS IN types/
    return nil
}

func CreateUser() *User {  // ❌ NO FUNCTIONS IN types/
    return &User{}
}
```

### **Rule 4: Domain-Based Microservice Structure**

- Use domain-based structure for microservices: `user/`, `auth/`, `product/`
- Each domain is 100% self-contained with its own `entities/`, `services/`, `controllers/`, etc.
- Shared code goes in dedicated packages: `shared/`, `common/`

### **Rule 5: Background Jobs in src/jobs/**

- **ALL** background jobs live in `src/jobs/` folder
- Main worker entry point: `jobs/worker.go`
- Individual jobs: `*_job.go` pattern (e.g., `email_job.go`)

### **Rule 6: Laravel-style configs/ Folder**

- Configuration split by concern: `app.go`, `database.go`, `cache.go`, `auth.go`
- **NOT** single monolithic config file

---

## **📄 File Naming Rules**

### **Rule 7: Use snake_case for File Names**

- **ALWAYS** use snake_case: `user_service.go`, `cache_interface.go`
- **NEVER** use kebab-case: `user-service.go` ❌
- **NEVER** use camelCase: `userService.go` ❌

### **Rule 8: Consistent File Suffixes**

- `_interface.go` - Interface definitions
- `_impl.go` - Concrete implementations
- `_test.go` - Test files
- `_mock.go` - Mock implementations
- `_stub.go` - Stub implementations
- `_enum.go` - Enum-like constants
- `_constants.go` - Constants
- `_types.go` - Type definitions
- `_entity.go` - Domain entities
- `_service.go` - Business logic services
- `_repository.go` - Data access layer
- `_controller.go` - HTTP controllers
- `_middleware.go` - HTTP middlewares
- `_validator.go` - Input validation
- `_helper.go` - Utility functions
- `_mapper.go` - Entity/DTO conversion
- `_factory.go` - Object creation
- `_builder.go` - Builder pattern
- `_exception.go` - Custom errors

### **Rule 9: Package Names Match Directory Names**

- Package name should be same as directory name
- Lowercase, no underscores preferred (single word if possible)
- Avoid stutter: if package is `user`, type should be `Service`, not `UserService`

---

## **🎨 Trait Pattern Rules**

### **Rule 10: Self-Contained Traits**

- **EACH** trait owns its data and behavior
- **NO** dependencies on external structs
- Trait manages its own state and lifecycle

```go
// ✅ CORRECT: Self-contained trait
type HasLocale struct {
    locale         string  // Owns its data
    fallbackLocale string
    timezone       string
    mutex          sync.RWMutex // Owns synchronization
}

// ❌ WRONG: Dependent trait
type HasLocale struct {
    app *Application  // Creates tight coupling
}
```

### **Rule 11: Dependency Injection Pattern**

- Traits accept **interfaces**, not concrete types
- Use minimal interfaces following ISP
- Constructor functions accept interface

```go
// ✅ CORRECT: Trait accepts interface
type LocaleAppInterface interface {
    GetLocale() string
    SetLocale(locale string)
    GetMutex() *sync.RWMutex
}

func NewHasLocale(app LocaleAppInterface) *HasLocale {
    return &HasLocale{
        app: app,
    }
}

// ❌ WRONG: Trait depends on concrete type
func NewHasLocale(app *Application) *HasLocale {
    return &HasLocale{
        app: app,
    }
}
```

### **Rule 12: Integration Pattern**

- Application implements trait interface
- Convenience methods delegate to trait
- Lazy initialization of trait instances

```go
// Application implements interface
func (a *Application) GetLocale() string {
    return a.locale
}

// Convenience method delegates to trait
func (a *Application) Locale() string {
    return a.HasLocale().GetLocale()
}
```

---

## **🧩 Interface Design Rules**

### **Rule 13: Interface Segregation Principle (ISP)**

- **SMALL** interfaces (1-3 methods maximum)
- **FOCUSED** on single responsibility
- **NO** god interfaces with mixed concerns

```go
// ✅ CORRECT: Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// ❌ WRONG: God interface
type DatabaseManager interface {
    // Connection methods
    Connect() error
    Disconnect() error
    
    // CRUD methods
    Create(entity interface{}) error
    Read(id string) (interface{}, error)
    Update(id string, entity interface{}) error
    Delete(id string) error
    
    // Transaction methods
    BeginTransaction() Transaction
    Commit() error
    Rollback() error
    
    // ... 20+ more methods
}
```

### **Rule 14: Interface Composition**

- Compose larger interfaces from smaller ones
- Use embedding for logical grouping
- Progressive complexity

```go
// Build complex interfaces from simple ones
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### **Rule 15: Context-Specific Interfaces**

- Define interfaces where they're **used**, not where they're **implemented**
- Different contexts may need different interface views

```go
// Define in consuming package
package userservice

type UserRepository interface {
    GetUser(id string) (*User, error)  // Only what service needs
}

// Implementation can be elsewhere
package database

func (r *DatabaseUserRepository) GetUser(id string) (*User, error) {
    // Implementation
}
```

---

## **✅ Interface Compliance Rules**

### **Rule 16: Compile-Time Interface Checks**

- **ALWAYS** add interface compliance checks
- Place at bottom of files
- Use descriptive comments

```go
// ✅ CORRECT: Interface compliance checks
var (
    // Ensure UserService implements all required business logic methods
    _ UserService = (*userServiceImpl)(nil)
    
    // Verify HTTP handler compliance for REST API endpoints  
    _ http.Handler = (*UserHandler)(nil)
    
    // Check database repository implements full CRUD interface
    _ UserRepository = (*sqlUserRepository)(nil)
)
```

### **Rule 17: Pointer vs Value Receiver Checks**

- Check correct receiver type (pointer vs value)
- Understand method set implications

```go
// Method has pointer receiver
func (u *User) Validate() error {
    return nil
}

// ✅ CORRECT: Check pointer type
var _ Validator = (*User)(nil)

// ❌ WRONG: This would fail
// var _ Validator = User{}
```

---

## **🧪 Testing Rules**

### **Rule 18: Interface-Based Testing**

- Test via interfaces, not concrete types
- Easy mocking with small interfaces
- Each interface tested independently

```go
// ✅ CORRECT: Test via interface
func TestUserService(t *testing.T) {
    mock := &MockUserRepository{
        users: map[string]*User{
            "1": {ID: "1", Name: "John"},
        },
    }
    
    service := NewUserService(mock)  // Accepts interface
    user, err := service.GetUser("1")
    
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### **Rule 19: Test Organization**

- Unit tests in domain folders: `user/tests/`
- Integration tests in root `tests/integrations/`
- Mocks in `tests/mocks/` or domain `mocks/`
- Fixtures in `tests/fixtures/`

---

## **📖 Documentation Rules**

### **Rule 20: Detailed Docblocks**

- **EVERY** public type, function, and method has docblocks
- Include parameter descriptions
- Provide usage examples
- Document return values and errors

```go
// CreateUser creates a new user account with the provided information.
// It validates the input data, checks for duplicate email addresses,
// and stores the user in the database.
//
// Parameters:
//   request: The user creation request containing name, email, and other details
//
// Returns:
//   *User: The created user with generated ID and timestamps
//   error: Validation errors, duplicate email errors, or database errors
//
// Example:
//   request := &CreateUserRequest{
//       Name:  "John Doe",
//       Email: "john@example.com",
//   }
//   user, err := service.CreateUser(request)
//   if err != nil {
//       log.Printf("Failed to create user: %v", err)
//       return err
//   }
//   fmt.Printf("Created user: %s\n", user.Name)
func (s *UserService) CreateUser(request *CreateUserRequest) (*User, error) {
    // Implementation
}
```

### **Rule 21: Package Documentation**

- Each package has `doc.go` with overview
- Explain package purpose and main concepts
- Provide usage examples
- Document architectural patterns used

---

## **🔧 Code Quality Rules**

### **Rule 22: DRY Principle**

- Extract common behavior into reusable components
- Use middleware for cross-cutting concerns
- Centralize constants and types
- Avoid copy-paste programming

### **Rule 23: Error Handling**

- Custom error types in `exceptions/` folder
- Wrap errors with context
- Use error interfaces for type checking

```go
// Custom error types
type ValidationException struct {
    Field   string
    Message string
}

func (e *ValidationException) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// Wrap errors with context
func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
    if err := s.validator.Validate(req); err != nil {
        return nil, fmt.Errorf("user creation validation failed: %w", err)
    }
    // ...
}
```

### **Rule 24: Generic Functions (Go 1.18+)**

- Use generics for type-safe reusable code
- Prefer type constraints over `interface{}`
- Use standard library constraints when possible

```go
import "golang.org/x/exp/constraints"

// Generic map function
func Map[T any, R any](slice []T, fn func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Constrained generic function
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

---

## **🏗️ Architecture Rules**

### **Rule 25: Separation of Concerns**

- Each layer has single responsibility
- No business logic in controllers
- No HTTP concerns in services
- No database details in entities

### **Rule 26: Dependency Direction**

- Dependencies point inward (toward business logic)
- Services depend on repository interfaces
- Controllers depend on service interfaces
- No circular dependencies

### **Rule 27: Builder Pattern**

- Use builders for complex object construction
- Fluent API for configuration
- Validation in Build() method

```go
type AppBuilder struct {
    name        string
    environment string
    debug       bool
}

func NewApp() *AppBuilder {
    return &AppBuilder{
        name:        "GoVel App",
        environment: "development",
        debug:       true,
    }
}

func (b *AppBuilder) WithName(name string) *AppBuilder {
    b.name = name
    return b
}

func (b *AppBuilder) WithEnvironment(env string) *AppBuilder {
    b.environment = env
    return b
}

func (b *AppBuilder) Build() (*Application, error) {
    // Validation and construction
    if b.name == "" {
        return nil, errors.New("application name is required")
    }
    
    return &Application{
        name:        b.name,
        environment: b.environment,
        debug:       b.debug,
    }, nil
}
```

---

## **🚀 Performance Rules**

### **Rule 28: Memory Management**

- Use object pools for frequently allocated objects
- Prefer value receivers for small structs
- Use pointer receivers for large structs or when mutation needed

### **Rule 29: Concurrency Safety**

- Protect shared state with mutexes
- Use channels for communication
- Prefer sync.RWMutex for read-heavy scenarios

```go
type SafeCounter struct {
    mu    sync.RWMutex
    count int64
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}
```

---

## **🔧 Tool Integration Rules**

### **Rule 30: Linting and Formatting**

- Use golangci-lint with comprehensive configuration
- Enable gofmt, goimports, gosec, staticcheck
- Pre-commit hooks for code quality

### **Rule 31: Git Commit Standards**

- Use Conventional Commits format
- feat, fix, docs, style, refactor, perf, test, chore
- Descriptive commit messages with examples

```
feat: add user authentication service with JWT tokens
fix: resolve memory leak in cache implementation  
docs: update API documentation for user endpoints
refactor: extract common validation logic into helper
test: add integration tests for payment processing
```

### **Rule 32: CI/CD Integration**

- Automated testing on all commits
- Security scanning (gosec, govulncheck)
- Code coverage reporting
- Automated releases with semantic versioning

---

## **📋 Summary Checklist**

When implementing any feature in GoVel, ensure:

- [ ] Code lives in `src/` folder
- [ ] Folder names are plural
- [ ] `types/` contains only struct definitions
- [ ] File names use snake_case with appropriate suffixes
- [ ] Interfaces are small and focused (ISP)
- [ ] Traits are self-contained with dependency injection
- [ ] Interface compliance checks added
- [ ] Comprehensive docblocks with examples
- [ ] Error handling with custom exceptions
- [ ] Tests use interface mocking
- [ ] Generic functions where appropriate
- [ ] Builder patterns for complex construction
- [ ] Thread-safe concurrent access
- [ ] Linting passes with no warnings
- [ ] Conventional commit format used

---

## **🎯 Implementation Priority**

1. **High Priority**: Project structure, file naming, interface design
2. **Medium Priority**: Trait patterns, testing strategy, documentation
3. **Low Priority**: Performance optimizations, advanced generics

Following these rules ensures consistent, maintainable, and scalable Go applications that leverage modern language features and proven architectural patterns.
