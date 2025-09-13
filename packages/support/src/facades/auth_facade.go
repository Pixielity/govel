package facades

import (
	authInterfaces "govel/packages/types/src/interfaces/auth"
	facade "govel/packages/support/src"
)

// Auth provides a clean, static-like interface to the application's authentication service.
//
// This facade implements the facade pattern, providing global access to the authentication
// service configured in the dependency injection container. It offers a Laravel-style
// API for authentication operations with automatic service resolution, session management,
// guard selection, and security features.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved authentication service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent authentication operations across goroutines
//   - Supports multiple authentication guards (session, token, JWT, etc.)
//   - Built-in session management and CSRF protection
//
// Behavior:
//   - First call: Resolves auth service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if auth service cannot be resolved (fail-fast behavior)
//   - Automatically handles session state, token validation, and guard switching
//
// Returns:
//   - AuthInterface: The application's authentication service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "auth" service is not registered in the container
//   - If the resolved service doesn't implement AuthInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Auth() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Authentication state is properly isolated per session/context
//
// Usage Examples:
//
//	// Check if user is authenticated
//	if facades.Auth().Check() {
//	    fmt.Println("User is logged in")
//	} else {
//	    fmt.Println("User is not authenticated")
//	}
//
//	// Get the currently authenticated user
//	user := facades.Auth().User()
//	if user != nil {
//	    fmt.Printf("Logged in as: %s\n", user.Name)
//	}
//
//	// Get user ID of currently authenticated user
//	userID := facades.Auth().ID()
//	if userID != nil {
//	    fmt.Printf("Current user ID: %d\n", *userID)
//	}
//
//	// Login with credentials
//	credentials := map[string]interface{}{
//	    "email":    "user@example.com",
//	    "password": "secret123",
//	}
//
//	if facades.Auth().Attempt(credentials, false) {
//	    fmt.Println("Login successful")
//	} else {
//	    fmt.Println("Invalid credentials")
//	}
//
//	// Login with "Remember Me" functionality
//	if facades.Auth().Attempt(credentials, true) {
//	    fmt.Println("Login successful with remember token")
//	}
//
//	// Login a user directly (bypass password check)
//	user := GetUserFromDatabase(123)
//	facades.Auth().Login(user, false) // false = don't remember
//
//	// Login and remember user
//	facades.Auth().Login(user, true) // true = set remember token
//
//	// Login user for a single request (stateless)
//	facades.Auth().Once(credentials)
//
//	// Logout current user
//	facades.Auth().Logout()
//	fmt.Println("User logged out")
//
//	// Logout all sessions for current user
//	facades.Auth().LogoutOtherDevices("current_password")
//
//	// Working with guards
//	// Switch to API guard for token-based auth
//	apiAuth := facades.Auth().Guard("api")
//	if apiAuth.Check() {
//	    apiUser := apiAuth.User()
//	    fmt.Printf("API user: %s\n", apiUser.Email)
//	}
//
//	// Use session guard explicitly
//	sessionAuth := facades.Auth().Guard("session")
//	if sessionAuth.Attempt(credentials, false) {
//	    fmt.Println("Session login successful")
//	}
//
//	// JWT guard example
//	jwtAuth := facades.Auth().Guard("jwt")
//	token := jwtAuth.Attempt(credentials)
//	if token != "" {
//	    fmt.Printf("JWT token: %s\n", token)
//	}
//
//	// Middleware integration
//	func AuthMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        if !facades.Auth().Check() {
//	            http.Redirect(w, r, "/login", http.StatusFound)
//	            return
//	        }
//	        next.ServeHTTP(w, r)
//	    })
//	}
//
//	// Role and permission checking
//	if facades.Auth().User().HasRole("admin") {
//	    fmt.Println("User is an admin")
//	}
//
//	if facades.Auth().User().Can("edit-posts") {
//	    fmt.Println("User can edit posts")
//	}
//
//	// Two-factor authentication
//	if facades.Auth().User().HasTwoFactorEnabled() {
//	    if !facades.Auth().ConfirmTwoFactor(totpCode) {
//	        fmt.Println("Invalid 2FA code")
//	        return
//	    }
//	}
//
//	// Password verification
//	if facades.Auth().ValidateCredentials(user, "password123") {
//	    fmt.Println("Password is correct")
//	}
//
//	// Generate password reset tokens
//	token := facades.Auth().CreatePasswordResetToken(user)
//	fmt.Printf("Password reset token: %s\n", token)
//
//	// Verify password reset token
//	if facades.Auth().ValidatePasswordResetToken(user, token) {
//	    // Allow password reset
//	    facades.Auth().ResetPassword(user, "new_password")
//	}
//
// Advanced Authentication Patterns:
//
//	// Multi-guard authentication
//	func GetAuthenticatedUser() (User, string) {
//	    // Try session authentication first
//	    if facades.Auth().Guard("session").Check() {
//	        return facades.Auth().Guard("session").User(), "session"
//	    }
//
//	    // Fall back to API token authentication
//	    if facades.Auth().Guard("api").Check() {
//	        return facades.Auth().Guard("api").User(), "api"
//	    }
//
//	    // Try JWT authentication
//	    if facades.Auth().Guard("jwt").Check() {
//	        return facades.Auth().Guard("jwt").User(), "jwt"
//	    }
//
//	    return nil, ""
//	}
//
//	// Rate limiting failed login attempts
//	func AttemptLogin(credentials map[string]interface{}) error {
//	    email := credentials["email"].(string)
//
//	    // Check if too many failed attempts
//	    if facades.Auth().HasTooManyLoginAttempts(email) {
//	        return ErrTooManyAttempts
//	    }
//
//	    if facades.Auth().Attempt(credentials, false) {
//	        facades.Auth().ClearLoginAttempts(email)
//	        return nil
//	    }
//
//	    facades.Auth().IncrementLoginAttempts(email)
//	    return ErrInvalidCredentials
//	}
//
//	// Session regeneration for security
//	func SecureLogin(credentials map[string]interface{}) error {
//	    if facades.Auth().Attempt(credentials, false) {
//	        // Regenerate session ID to prevent session fixation
//	        facades.Auth().RegenerateSession()
//	        return nil
//	    }
//	    return ErrInvalidCredentials
//	}
//
//	// Impersonation functionality
//	func ImpersonateUser(adminUser, targetUser User) {
//	    if !adminUser.HasRole("admin") {
//	        panic("Insufficient permissions")
//	    }
//
//	    facades.Auth().Impersonate(targetUser)
//	}
//
//	func StopImpersonating() {
//	    facades.Auth().StopImpersonation()
//	}
//
// Best Practices:
//   - Always check authentication status before accessing protected resources
//   - Use appropriate guards for different authentication methods
//   - Implement proper session management and CSRF protection
//   - Use "Remember Me" functionality judiciously
//   - Implement rate limiting for login attempts
//   - Regenerate session IDs after authentication
//   - Use secure password hashing and validation
//   - Implement proper logout functionality that clears all session data
//
// Security Considerations:
//   - Protect against session fixation attacks
//   - Implement CSRF protection for state-changing operations
//   - Use secure, httpOnly cookies for session management
//   - Implement proper password policies and validation
//   - Consider implementing two-factor authentication
//   - Log authentication events for security monitoring
//   - Implement account lockout after failed login attempts
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume authentication service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	auth, err := facade.TryResolve[AuthInterface]("auth")
//	if err != nil {
//	    // Handle auth service unavailability gracefully
//	    return nil, false // Return unauthenticated state
//	}
//	isAuthenticated := auth.Check()
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserController(t *testing.T) {
//	    // Create a test auth service
//	    testAuth := &TestAuth{
//	        authenticated: true,
//	        user: &User{
//	            ID:    123,
//	            Name:  "Test User",
//	            Email: "test@example.com",
//	        },
//	    }
//
//	    // Swap the real auth with test auth
//	    restore := support.SwapService("auth", testAuth)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Auth() returns testAuth
//	    controller := NewUserController()
//
//	    // Test authenticated behavior
//	    assert.True(t, facades.Auth().Check())
//	    user := facades.Auth().User()
//	    assert.Equal(t, "Test User", user.Name)
//
//	    // Test controller behavior
//	    response := controller.Profile()
//	    assert.Equal(t, "Test User", response.Name)
//	}
//
//	// Test unauthenticated state
//	func TestUnauthenticatedAccess(t *testing.T) {
//	    testAuth := &TestAuth{authenticated: false}
//	    restore := support.SwapService("auth", testAuth)
//	    defer restore()
//
//	    assert.False(t, facades.Auth().Check())
//	    assert.Nil(t, facades.Auth().User())
//	}
//
// Container Configuration:
// Ensure the authentication service is properly configured in your container:
//
//	// Example auth registration
//	container.Singleton("auth", func() interface{} {
//	    config := auth.Config{
//	        // Default guard
//	        DefaultGuard: "session",
//
//	        // Guard configurations
//	        Guards: map[string]auth.GuardConfig{
//	            "session": {
//	                Driver: "session",
//	                Provider: "users",
//	            },
//	            "api": {
//	                Driver: "token",
//	                Provider: "users",
//	                StorageKey: "api_token",
//	                Hash: false,
//	            },
//	            "jwt": {
//	                Driver: "jwt",
//	                Provider: "users",
//	                Secret: "your-secret-key",
//	                TTL: time.Hour * 24, // 24 hours
//	            },
//	        },
//
//	        // User providers
//	        Providers: map[string]auth.ProviderConfig{
//	            "users": {
//	                Driver: "database",
//	                Model: "User",
//	                Table: "users",
//	            },
//	        },
//
//	        // Password configuration
//	        Passwords: auth.PasswordConfig{
//	            ResetTable: "password_resets",
//	            ExpireMinutes: 60,
//	            Throttle: 60, // seconds between reset attempts
//	        },
//
//	        // Session configuration
//	        Session: auth.SessionConfig{
//	            Key: "auth_session",
//	            Lifetime: time.Hour * 24 * 30, // 30 days
//	            ExpireOnClose: false,
//	            Encrypt: true,
//	        },
//
//	        // Remember me configuration
//	        RememberMe: auth.RememberConfig{
//	            Key: "remember_token",
//	            Lifetime: time.Hour * 24 * 365, // 1 year
//	        },
//	    }
//
//	    authService, err := auth.NewAuthManager(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create auth service: %v", err)
//	    }
//
//	    return authService
//	})
func Auth() authInterfaces.AuthInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves authentication service using type-safe token from the dependency injection container
	// - Performs type assertion to AuthInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[authInterfaces.AuthInterface](authInterfaces.AUTH_TOKEN)
}

// AuthWithError provides error-safe access to the authentication service.
//
// This function offers the same functionality as Auth() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle authentication service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Auth() but with error handling.
//
// Returns:
//   - AuthInterface: The resolved auth instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement AuthInterface
//
// Usage Examples:
//
//	// Basic error-safe authentication access
//	auth, err := facades.AuthWithError()
//	if err != nil {
//	    log.Printf("Auth service unavailable: %v", err)
//	    return false // Return unauthenticated state
//	}
//	isAuthenticated := auth.Check()
//
//	// Conditional authentication operations
//	if auth, err := facades.AuthWithError(); err == nil {
//	    if auth.Check() {
//	        // Perform authenticated user operations
//	        auth.User().UpdateLastSeen()
//	    }
//	}
//
//	// Health check pattern
//	func CheckAuthHealth() error {
//	    auth, err := facades.AuthWithError()
//	    if err != nil {
//	        return fmt.Errorf("auth service unavailable: %w", err)
//	    }
//
//	    // Test basic auth functionality
//	    if !auth.CanPerformOperations() {
//	        return fmt.Errorf("auth service not operational")
//	    }
//
//	    return nil
//	}
func AuthWithError() (authInterfaces.AuthInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves authentication service using type-safe token from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[authInterfaces.AuthInterface](authInterfaces.AUTH_TOKEN)
}
