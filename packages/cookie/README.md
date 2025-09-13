# GoVel Cookie Package

A Laravel-compatible cookie management package for the GoVel framework, providing comprehensive cookie handling with encryption, CSRF protection, and advanced security features.

## Features

- **Laravel-Compatible API**: Full compatibility with Laravel's Cookie facade and functionality
- **Cookie Encryption**: Automatic cookie encryption/decryption with multiple cipher support
- **Cookie Queuing**: Batch cookie processing with middleware integration
- **CSRF Protection**: Laravel-style CSRF token generation and validation
- **SameSite Policies**: Advanced SameSite attribute management with browser compatibility
- **Thread-Safe Operations**: Concurrent cookie operations with internal synchronization
- **Flexible Configuration**: Extensive configuration options with sensible defaults
- **Service Provider Integration**: Full dependency injection support

## Installation

```bash
go get -u govel/packages/cookie
```

## Quick Start

### Basic Usage

```go
package main

import (
    "net/http"
    "govel/packages/cookie/src/facades"
    cookieInterfaces "govel/packages/cookie/src/interfaces"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Create a simple session cookie
    cookie := facades.Make("user_id", "123")
    http.SetCookie(w, cookie)
    
    // Create a persistent cookie
    cookie = facades.Forever("theme", "dark", 
        cookieInterfaces.WithDomain(".example.com"),
        cookieInterfaces.WithSecure(true),
    )
    http.SetCookie(w, cookie)
    
    // Queue cookies for batch processing
    facades.Queue(facades.Make("flash_message", "Welcome!"))
    
    // Delete a cookie
    cookie = facades.Forget("old_session")
    http.SetCookie(w, cookie)
}
```

### Service Registration

```go
package main

import (
    cookieProviders "govel/packages/cookie/src/providers"
    "govel/packages/container"
    "govel/packages/config"
)

func main() {
    container := container.New()
    config := config.New()
    
    // Register cookie services
    err := cookieProviders.RegisterCookieServices(container, config)
    if err != nil {
        panic(err)
    }
    
    // Cookie services are now available via facades or container resolution
}
```

## Core Components

### Cookie Jar

The main cookie management service providing cookie creation, queuing, and configuration:

```go
import cookie "govel/packages/cookie/src"

// Create a new cookie jar
jar := cookie.NewCookieJar()

// Configure default settings
jar.SetDefaultPath("/api")
jar.SetDefaultDomain(".example.com")
jar.SetDefaultSecure(true)

// Create cookies with defaults applied
sessionCookie := jar.Make("session", sessionID)
persistentCookie := jar.Forever("preferences", userPrefs)

// Queue cookies for batch processing
jar.Queue(sessionCookie)
jar.Queue(persistentCookie)

// Process queued cookies
cookies := jar.GetQueuedCookies()
for _, cookie := range cookies {
    http.SetCookie(w, cookie)
}
```

### Facades

Laravel-style static access to cookie functionality:

```go
import "govel/packages/cookie/src/facades"

// Direct cookie operations
cookie := facades.Make("name", "value")
cookie = facades.Forever("remember", token)
cookie = facades.Forget("session")

// Queue operations
facades.Queue(cookie)
facades.QueueForever("preferences", data)
facades.QueueSession("flash", message)

// Configuration
facades.WithDefaultSecure(true).Make("secure_cookie", value)
```

### Middleware

Automatic cookie encryption and queue processing:

```go
import (
    cookieMiddlewares "govel/packages/cookie/src/middlewares"
    encrypter "govel/packages/encryption/src"
)

// Cookie encryption middleware
encryptMiddleware := cookieMiddlewares.NewEncryptCookies(
    encrypter,
    cookieMiddlewares.WithEncryptAll(true),
    cookieMiddlewares.WithExceptCookies([]string{"csrf_token"}),
)

// Queue processing middleware
queueMiddleware := cookieMiddlewares.NewAddQueuedCookiesToResponse(cookieJar)

// Apply middleware to your router
router.Use(encryptMiddleware.Handle)
router.Use(queueMiddleware.Handle)
```

## Advanced Features

### Cookie Encryption

Automatic encryption/decryption of cookie values:

```go
// Configure encryption
middleware := cookieMiddlewares.NewEncryptCookies(
    encrypter,
    cookieMiddlewares.WithEncryptedCookies([]string{
        "user_session", 
        "shopping_cart",
        "user_preferences",
    }),
    cookieMiddlewares.WithExceptCookies([]string{
        "csrf_token",
        "language",
        "theme",
    }),
)

// Cookies are automatically encrypted on response
// and decrypted on incoming requests
```

### CSRF Protection

Laravel-compatible CSRF protection:

```go
import cookieSecurity "govel/packages/cookie/src/security"

// Create CSRF protection
csrf := cookieSecurity.NewCSRFProtection(cookieJar,
    cookieSecurity.WithTokenName("_token"),
    cookieSecurity.WithCookieName("XSRF-TOKEN"),
    cookieSecurity.WithHeaderName("X-CSRF-TOKEN"),
    cookieSecurity.WithExcept([]string{"/api/public/*"}),
)

// Generate tokens
token, err := csrf.GenerateToken()

// Validate tokens
isValid := csrf.ValidateToken(requestToken, storedToken)

// Use as middleware
router.Use(csrf.Middleware())
```

### SameSite Policies

Advanced SameSite attribute management:

```go
import cookieSecurity "govel/packages/cookie/src/security"

// Create SameSite manager
manager := cookieSecurity.NewSameSiteManager(
    cookieSecurity.WithDefaultPolicy(cookieSecurity.SameSiteLax),
    cookieSecurity.WithCookiePolicy("admin_session", cookieSecurity.SameSiteStrict),
    cookieSecurity.WithCookiePolicy("api_token", cookieSecurity.SameSiteStrict),
    cookieSecurity.WithEnforceSecure(true),
    cookieSecurity.WithCheckUserAgent(true),
)

// Apply policies to cookies
manager.ApplySameSitePolicy(cookie, request)
```

### Cookie Options

Flexible cookie configuration using functional options:

```go
import cookieInterfaces "govel/packages/cookie/src/interfaces"

cookie := jar.Make("advanced_cookie", "value",
    // Expiration options
    cookieInterfaces.WithExpiry(time.Now().Add(24*time.Hour)),
    cookieInterfaces.WithMaxAge(86400), // 24 hours in seconds
    
    // Domain and path
    cookieInterfaces.WithDomain(".example.com"),
    cookieInterfaces.WithPath("/app"),
    
    // Security options
    cookieInterfaces.WithSecure(true),
    cookieInterfaces.WithHttpOnly(true),
    cookieInterfaces.WithSameSite(http.SameSiteStrictMode),
)
```

## Laravel Compatibility

This package maintains full compatibility with Laravel's cookie system:

### Method Compatibility

| Laravel Method | GoVel Equivalent | Description |
|---|---|---|
| `Cookie::make()` | `facades.Make()` | Create a new cookie |
| `Cookie::forever()` | `facades.Forever()` | Create a "forever" cookie |
| `Cookie::forget()` | `facades.Forget()` | Delete a cookie |
| `Cookie::queue()` | `facades.Queue()` | Queue a cookie |
| `Cookie::unqueue()` | `facades.Unqueue()` | Remove from queue |
| `Cookie::queued()` | `facades.Queued()` | Get queued cookie |

### Middleware Compatibility

| Laravel Middleware | GoVel Equivalent |
|---|---|
| `EncryptCookies` | `cookieMiddlewares.NewEncryptCookies()` |
| `AddQueuedCookiesToResponse` | `cookieMiddlewares.NewAddQueuedCookiesToResponse()` |

### Payload Format

Cookies use Laravel's JSON payload format for cross-platform compatibility:

```json
{
  "iv": "base64_encoded_iv",
  "value": "base64_encoded_encrypted_data", 
  "mac": "hmac_sha256_signature"
}
```

## Configuration

### Environment Variables

```env
# Cookie encryption
COOKIE_ENCRYPTION_KEY=base64:your_encryption_key_here
COOKIE_CIPHER=AES-128-CBC

# Default cookie settings
COOKIE_PATH=/
COOKIE_DOMAIN=.example.com
COOKIE_SECURE=true
COOKIE_HTTP_ONLY=true
COOKIE_SAME_SITE=lax

# Session cookie settings
SESSION_LIFETIME=120
SESSION_EXPIRE_ON_CLOSE=false

# CSRF protection
CSRF_TOKEN_NAME=_token
CSRF_COOKIE_NAME=XSRF-TOKEN
CSRF_HEADER_NAME=X-CSRF-TOKEN
```

### Cookie Jar Configuration

```go
// Configure default cookie settings
jar := cookie.NewCookieJar()

// Path and domain settings
jar.SetDefaultPath("/api/v1")
jar.SetDefaultDomain(".example.com")

// Security settings
jar.SetDefaultSecure(true)
jar.SetDefaultHttpOnly(true)
jar.SetDefaultSameSite(http.SameSiteLaxMode)

// Expiration settings
jar.SetDefaultMaxAge(3600) // 1 hour
```

## Testing

### Unit Tests

Run the complete test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Example Tests

```go
package cookie_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    
    cookie "govel/packages/cookie/src"
    "github.com/stretchr/testify/assert"
)

func TestCookieJar_Make(t *testing.T) {
    jar := cookie.NewCookieJar()
    
    cookie := jar.Make("test_cookie", "test_value")
    
    assert.Equal(t, "test_cookie", cookie.Name)
    assert.Equal(t, "test_value", cookie.Value)
    assert.Equal(t, "/", cookie.Path)
}

func TestCookieJar_Forever(t *testing.T) {
    jar := cookie.NewCookieJar()
    
    cookie := jar.Forever("persistent_cookie", "value")
    
    assert.True(t, cookie.Expires.After(time.Now().Add(5*365*24*time.Hour)))
}

func TestCookieJar_Queue(t *testing.T) {
    jar := cookie.NewCookieJar()
    
    cookie1 := jar.Make("cookie1", "value1")
    cookie2 := jar.Make("cookie2", "value2")
    
    jar.Queue(cookie1)
    jar.Queue(cookie2)
    
    queued := jar.GetQueuedCookies()
    assert.Len(t, queued, 2)
    assert.Equal(t, "cookie1", queued[0].Name)
    assert.Equal(t, "cookie2", queued[1].Name)
}

func TestEncryptCookies_Middleware(t *testing.T) {
    // Test cookie encryption middleware
    // Implementation depends on your encryption package
}

func TestCSRFProtection(t *testing.T) {
    jar := cookie.NewCookieJar()
    csrf := security.NewCSRFProtection(jar)
    
    token, err := csrf.GenerateToken()
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    isValid := csrf.ValidateToken(token, token)
    assert.True(t, isValid)
}
```

## Performance Considerations

### Thread Safety

The cookie jar is thread-safe and can handle concurrent operations:

```go
// Safe for concurrent use
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        cookie := jar.Make(fmt.Sprintf("cookie_%d", id), "value")
        jar.Queue(cookie)
    }(i)
}
wg.Wait()
```

### Memory Management

Queued cookies are automatically cleared after processing:

```go
// Queue cookies
jar.Queue(cookie1)
jar.Queue(cookie2)

// Process and clear queue
cookies := jar.GetQueuedCookies()
jar.ClearQueue()
```

### Encryption Performance

For high-throughput applications, consider:

- Limiting encrypted cookies to sensitive data only
- Using `WithExceptCookies()` to exclude non-sensitive cookies
- Implementing cookie value caching where appropriate

```go
middleware := cookieMiddlewares.NewEncryptCookies(
    encrypter,
    // Only encrypt specific cookies
    cookieMiddlewares.WithEncryptedCookies([]string{
        "user_session",
        "auth_token", 
        "user_data",
    }),
    // Exclude performance-critical cookies
    cookieMiddlewares.WithExceptCookies([]string{
        "analytics_id",
        "theme_preference",
        "language",
    }),
)
```

## Security Best Practices

### Cookie Security

1. **Always use HTTPS in production**:
   ```go
   jar.SetDefaultSecure(true)
   ```

2. **Enable HttpOnly for session cookies**:
   ```go
   jar.SetDefaultHttpOnly(true)
   ```

3. **Use appropriate SameSite policies**:
   ```go
   jar.SetDefaultSameSite(http.SameSiteStrictMode) // For auth cookies
   jar.SetDefaultSameSite(http.SameSiteLaxMode)    // For general use
   ```

4. **Set secure cookie domains**:
   ```go
   jar.SetDefaultDomain(".yourdomain.com") // Subdomain access
   // jar.SetDefaultDomain("yourdomain.com") // Single domain only
   ```

### Encryption Security

1. **Use strong encryption keys**:
   ```bash
   # Generate a secure key
   openssl rand -base64 32
   ```

2. **Rotate encryption keys periodically**
3. **Encrypt sensitive cookie data**:
   ```go
   encryptedCookies := []string{
       "user_session",
       "auth_token",
       "payment_info",
       "personal_data",
   }
   ```

### CSRF Protection

1. **Enable CSRF protection for state-changing operations**
2. **Use secure token generation**
3. **Validate tokens on all POST/PUT/DELETE requests**

```go
csrf := security.NewCSRFProtection(jar,
    security.WithTokenLength(32),
    security.WithTokenLifetime(time.Hour),
    security.WithExcept([]string{"/api/public/*"}),
)
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Write tests** for your changes
4. **Follow Go conventions** and run `gofmt`
5. **Update documentation** as needed
6. **Submit a pull request**

### Development Setup

```bash
# Clone the repository
git clone https://github.com/govel/packages.git
cd packages/cookie

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run linting
golangci-lint run
```

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Write comprehensive tests
- Follow Laravel patterns where applicable

## License

This package is open-sourced software licensed under the [MIT License](LICENSE).

## Support

For support, please:

1. **Check the documentation** above
2. **Search existing issues** on GitHub
3. **Create a new issue** with:
   - Go version
   - Package version
   - Minimal reproduction case
   - Expected vs actual behavior

## Changelog

### v1.0.0 (Latest)

**Added:**
- Laravel-compatible Cookie facade and functionality
- Comprehensive cookie encryption with multiple cipher support
- Cookie queuing system with middleware integration
- CSRF protection with Laravel-style token generation
- Advanced SameSite policy management
- Thread-safe cookie operations
- Service provider with dependency injection
- Extensive configuration options
- Full test coverage
- Comprehensive documentation

**Features:**
- Cookie Jar with Laravel-compatible methods
- EncryptCookies middleware for automatic encryption/decryption
- AddQueuedCookiesToResponse middleware for batch processing
- CSRF protection middleware
- SameSite policy manager with browser compatibility
- Functional options pattern for flexible configuration
- Laravel payload format compatibility

---

**GoVel Cookie Package** - Building Laravel-compatible cookie management for Go applications.
