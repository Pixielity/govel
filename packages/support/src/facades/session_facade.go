package facades

import (
	sessionInterfaces "govel/packages/types/src/interfaces/session"
	facade "govel/packages/support/src"
)

// Session provides a clean, static-like interface to the application's session management service.
//
// This facade implements the facade pattern, providing global access to the session
// service configured in the dependency injection container. It offers a Laravel-style
// API for session management with automatic service resolution, multiple storage drivers,
// flash data support, and comprehensive session security.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved session service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent session operations across goroutines
//   - Supports multiple storage drivers (file, database, Redis, memory, cookie)
//   - Built-in session security with CSRF protection and encryption
//
// Behavior:
//   - First call: Resolves session service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if session service cannot be resolved (fail-fast behavior)
//   - Automatically handles session lifecycle, data persistence, and cleanup
//
// Returns:
//   - SessionInterface: The application's session service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "session" service is not registered in the container
//   - If the resolved service doesn't implement SessionInterface
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
//   - Multiple goroutines can call Session() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Session operations are thread-safe with proper locking mechanisms
//
// Usage Examples:
//
//	// Basic session data management
//	facades.Session().Put("user_id", 123)
//	facades.Session().Put("username", "john_doe")
//	facades.Session().Put("preferences", map[string]interface{}{
//	    "theme":    "dark",
//	    "language": "en",
//	    "timezone": "UTC",
//	})
//
//	// Retrieve session data
//	userID := facades.Session().Get("user_id")
//	username := facades.Session().GetString("username")
//	preferences := facades.Session().GetMap("preferences")
//
//	if userID != nil {
//	    fmt.Printf("User %s (ID: %v) is logged in\n", username, userID)
//	}
//
//	// Type-safe retrieval with defaults
//	userID := facades.Session().GetInt("user_id", 0)
//	username := facades.Session().GetString("username", "guest")
//	isAdmin := facades.Session().GetBool("is_admin", false)
//	preferences := facades.Session().GetStringMap("preferences", map[string]string{})
//
//	// Check for existence
//	if facades.Session().Has("user_id") {
//	    fmt.Println("User is authenticated")
//	}
//
//	if facades.Session().Missing("shopping_cart") {
//	    facades.Session().Put("shopping_cart", []interface{}{})
//	}
//
//	// Flash data (available only for the next request)
//	facades.Session().Flash("success", "Profile updated successfully!")
//	facades.Session().Flash("errors", []string{
//	    "Email is required",
//	    "Password must be at least 8 characters",
//	})
//
//	// Retrieve flash data
//	successMessage := facades.Session().GetFlash("success")
//	errorMessages := facades.Session().GetFlashStringSlice("errors")
//
//	// Flash data is automatically removed after retrieval
//	if successMessage != nil {
//	    fmt.Printf("Success: %s\n", successMessage)
//	}
//
//	// Session management
//	sessionID := facades.Session().ID()
//	fmt.Printf("Current session ID: %s\n", sessionID)
//
//	// Regenerate session ID (security best practice)
//	facades.Session().Regenerate()
//
//	// Flush all session data
//	facades.Session().Flush()
//
//	// Invalidate session (regenerate ID and flush data)
//	facades.Session().Invalidate()
//
//	// Array/slice operations
//	// Add items to session array
//	facades.Session().Push("recent_pages", "/dashboard")
//	facades.Session().Push("recent_pages", "/profile")
//	facades.Session().Push("recent_pages", "/settings")
//
//	// Get array data
//	recentPages := facades.Session().GetSlice("recent_pages")
//	for _, page := range recentPages {
//	    fmt.Printf("Recent page: %v\n", page)
//	}
//
//	// Remove specific array item
//	facades.Session().Pull("recent_pages", "/profile")
//
//	// Increment/decrement operations
//	facades.Session().Increment("page_views")
//	facades.Session().IncrementBy("page_views", 5)
//	facades.Session().Decrement("attempts")
//	facades.Session().DecrementBy("attempts", 2)
//
//	pageViews := facades.Session().GetInt("page_views", 0)
//	attempts := facades.Session().GetInt("attempts", 0)
//
// Advanced Session Patterns:
//
//	// Session-based authentication
//	type UserAuth struct {
//	    ID       int    `json:"id"`
//	    Username string `json:"username"`
//	    Email    string `json:"email"`
//	    Role     string `json:"role"`
//	}
//
//	func LoginUser(user *UserAuth) {
//	    // Store user data in session
//	    facades.Session().Put("auth.user_id", user.ID)
//	    facades.Session().Put("auth.username", user.Username)
//	    facades.Session().Put("auth.email", user.Email)
//	    facades.Session().Put("auth.role", user.Role)
//	    facades.Session().Put("auth.logged_in", true)
//	    facades.Session().Put("auth.login_time", time.Now())
//
//	    // Regenerate session ID for security
//	    facades.Session().Regenerate()
//
//	    facades.Session().Flash("success", "Login successful!")
//	}
//
//	func LogoutUser() {
//	    facades.Session().Remove("auth.user_id")
//	    facades.Session().Remove("auth.username")
//	    facades.Session().Remove("auth.email")
//	    facades.Session().Remove("auth.role")
//	    facades.Session().Remove("auth.logged_in")
//	    facades.Session().Remove("auth.login_time")
//
//	    facades.Session().Flash("info", "You have been logged out.")
//
//	    // Invalidate entire session
//	    facades.Session().Invalidate()
//	}
//
//	func GetAuthenticatedUser() *UserAuth {
//	    if !facades.Session().GetBool("auth.logged_in", false) {
//	        return nil
//	    }
//
//	    return &UserAuth{
//	        ID:       facades.Session().GetInt("auth.user_id", 0),
//	        Username: facades.Session().GetString("auth.username", ""),
//	        Email:    facades.Session().GetString("auth.email", ""),
//	        Role:     facades.Session().GetString("auth.role", "user"),
//	    }
//	}
//
//	// Shopping cart session management
//	type CartItem struct {
//	    ID       int     `json:"id"`
//	    Name     string  `json:"name"`
//	    Price    float64 `json:"price"`
//	    Quantity int     `json:"quantity"`
//	}
//
//	func AddToCart(item CartItem) {
//	    cart := GetCart()
//
//	    // Check if item already exists
//	    for i, existingItem := range cart {
//	        if existingItem.ID == item.ID {
//	            cart[i].Quantity += item.Quantity
//	            facades.Session().Put("cart.items", cart)
//	            return
//	        }
//	    }
//
//	    // Add new item
//	    cart = append(cart, item)
//	    facades.Session().Put("cart.items", cart)
//
//	    // Update cart totals
//	    updateCartTotals(cart)
//	}
//
//	func GetCart() []CartItem {
//	    cartData := facades.Session().Get("cart.items")
//	    if cartData == nil {
//	        return []CartItem{}
//	    }
//
//	    cart, ok := cartData.([]CartItem)
//	    if !ok {
//	        return []CartItem{}
//	    }
//
//	    return cart
//	}
//
//	func updateCartTotals(cart []CartItem) {
//	    total := 0.0
//	    itemCount := 0
//
//	    for _, item := range cart {
//	        total += item.Price * float64(item.Quantity)
//	        itemCount += item.Quantity
//	    }
//
//	    facades.Session().Put("cart.total", total)
//	    facades.Session().Put("cart.item_count", itemCount)
//	}
//
//	func ClearCart() {
//	    facades.Session().Remove("cart.items")
//	    facades.Session().Remove("cart.total")
//	    facades.Session().Remove("cart.item_count")
//	}
//
//	// Multi-step form handling
//	type MultiStepForm struct {
//	    CurrentStep int                    `json:"current_step"`
//	    TotalSteps  int                    `json:"total_steps"`
//	    Data        map[string]interface{} `json:"data"`
//	    Completed   []int                  `json:"completed"`
//	}
//
//	func InitializeMultiStepForm(formName string, totalSteps int) {
//	    form := MultiStepForm{
//	        CurrentStep: 1,
//	        TotalSteps:  totalSteps,
//	        Data:        make(map[string]interface{}),
//	        Completed:   []int{},
//	    }
//
//	    facades.Session().Put(fmt.Sprintf("forms.%s", formName), form)
//	}
//
//	func SaveFormStep(formName string, stepData map[string]interface{}) {
//	    formData := facades.Session().Get(fmt.Sprintf("forms.%s", formName))
//	    if formData == nil {
//	        return
//	    }
//
//	    form, ok := formData.(MultiStepForm)
//	    if !ok {
//	        return
//	    }
//
//	    // Merge step data
//	    for key, value := range stepData {
//	        form.Data[key] = value
//	    }
//
//	    // Mark current step as completed
//	    form.Completed = append(form.Completed, form.CurrentStep)
//
//	    facades.Session().Put(fmt.Sprintf("forms.%s", formName), form)
//	}
//
//	func NextFormStep(formName string) bool {
//	    formData := facades.Session().Get(fmt.Sprintf("forms.%s", formName))
//	    if formData == nil {
//	        return false
//	    }
//
//	    form, ok := formData.(MultiStepForm)
//	    if !ok || form.CurrentStep >= form.TotalSteps {
//	        return false
//	    }
//
//	    form.CurrentStep++
//	    facades.Session().Put(fmt.Sprintf("forms.%s", formName), form)
//	    return true
//	}
//
// Session Security Patterns:
//
//	// CSRF token management
//	func GenerateCSRFToken() string {
//	    token := facades.Session().GetString("_csrf_token", "")
//	    if token == "" {
//	        // Generate new token
//	        token = generateRandomString(32)
//	        facades.Session().Put("_csrf_token", token)
//	    }
//	    return token
//	}
//
//	func ValidateCSRFToken(token string) bool {
//	    sessionToken := facades.Session().GetString("_csrf_token", "")
//	    return sessionToken != "" && sessionToken == token
//	}
//
//	// Session hijacking protection
//	func ValidateSessionFingerprint(r *http.Request) bool {
//	    currentFingerprint := generateSessionFingerprint(r)
//	    storedFingerprint := facades.Session().GetString("_fingerprint", "")
//
//	    if storedFingerprint == "" {
//	        facades.Session().Put("_fingerprint", currentFingerprint)
//	        return true
//	    }
//
//	    return storedFingerprint == currentFingerprint
//	}
//
//	func generateSessionFingerprint(r *http.Request) string {
//	    userAgent := r.Header.Get("User-Agent")
//	    acceptLang := r.Header.Get("Accept-Language")
//	    acceptEnc := r.Header.Get("Accept-Encoding")
//
//	    fingerprint := fmt.Sprintf("%s|%s|%s", userAgent, acceptLang, acceptEnc)
//	    return fmt.Sprintf("%x", sha256.Sum256([]byte(fingerprint)))
//	}
//
//	// Session timeout management
//	func CheckSessionTimeout() bool {
//	    lastActivity := facades.Session().GetTime("_last_activity", time.Time{})
//	    if lastActivity.IsZero() {
//	        facades.Session().Put("_last_activity", time.Now())
//	        return true
//	    }
//
//	    timeout := time.Duration(30) * time.Minute // 30 minutes timeout
//	    if time.Since(lastActivity) > timeout {
//	        facades.Session().Invalidate()
//	        return false
//	    }
//
//	    facades.Session().Put("_last_activity", time.Now())
//	    return true
//	}
//
// Session Data Serialization:
//
//	// Store complex objects
//	type UserProfile struct {
//	    ID          int                    `json:"id"`
//	    Name        string                 `json:"name"`
//	    Email       string                 `json:"email"`
//	    Preferences map[string]interface{} `json:"preferences"`
//	    CreatedAt   time.Time              `json:"created_at"`
//	}
//
//	func StoreUserProfile(profile UserProfile) error {
//	    return facades.Session().PutJSON("user_profile", profile)
//	}
//
//	func GetUserProfile() (*UserProfile, error) {
//	    var profile UserProfile
//	    err := facades.Session().GetJSON("user_profile", &profile)
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &profile, nil
//	}
//
// Session Driver Configuration:
//
//	// File-based session storage
//	func ConfigureFileSession() {
//	    config := session.FileConfig{
//	        StorePath:   "./storage/sessions",
//	        FilePrefix:  "sess_",
//	        FileMode:    0600,
//	        GCProbability: 0.01,
//	    }
//
//	    facades.Session().SetDriver("file", config)
//	}
//
//	// Database session storage
//	func ConfigureDatabaseSession() {
//	    config := session.DatabaseConfig{
//	        Connection: "default",
//	        Table:      "sessions",
//	        Lifetime:   120, // minutes
//	    }
//
//	    facades.Session().SetDriver("database", config)
//	}
//
//	// Redis session storage
//	func ConfigureRedisSession() {
//	    config := session.RedisConfig{
//	        Connection: "default",
//	        KeyPrefix:  "laravel_session:",
//	        Lifetime:   120, // minutes
//	    }
//
//	    facades.Session().SetDriver("redis", config)
//	}
//
// Session Middleware Integration:
//
//	// HTTP middleware for session management
//	func SessionMiddleware(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        // Start session
//	        facades.Session().Start(r, w)
//
//	        // Validate session security
//	        if !ValidateSessionFingerprint(r) {
//	            facades.Session().Invalidate()
//	            http.Redirect(w, r, "/login", http.StatusSeeOther)
//	            return
//	        }
//
//	        // Check session timeout
//	        if !CheckSessionTimeout() {
//	            facades.Session().Flash("warning", "Session expired. Please log in again.")
//	            http.Redirect(w, r, "/login", http.StatusSeeOther)
//	            return
//	        }
//
//	        // Continue with request
//	        next.ServeHTTP(w, r)
//
//	        // Save session data
//	        facades.Session().Save()
//	    })
//	}
//
// Testing Support:
//
//	// Test session management
//	func TestSessionOperations(t *testing.T) {
//	    // Create test session
//	    testSession := &TestSession{
//	        data: make(map[string]interface{}),
//	    }
//
//	    // Swap session service
//	    restore := support.SwapService("session", testSession)
//	    defer restore()
//
//	    // Test session operations
//	    facades.Session().Put("test_key", "test_value")
//	    value := facades.Session().GetString("test_key", "")
//	    assert.Equal(t, "test_value", value)
//
//	    // Test flash data
//	    facades.Session().Flash("message", "Flash message")
//	    flashMessage := facades.Session().GetFlash("message")
//	    assert.Equal(t, "Flash message", flashMessage)
//
//	    // Flash data should be gone after retrieval
//	    flashMessage = facades.Session().GetFlash("message")
//	    assert.Nil(t, flashMessage)
//	}
//
// Best Practices:
//   - Always regenerate session ID after login/privilege changes
//   - Use CSRF tokens for form submissions
//   - Implement session timeout and fingerprinting for security
//   - Store minimal data in sessions (use database for large objects)
//   - Use appropriate session drivers for your infrastructure
//   - Encrypt sensitive session data
//   - Implement proper session cleanup and garbage collection
//   - Use flash data for temporary messages and notifications
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume session service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	session, err := facade.TryResolve[SessionInterface]("session")
//	if err != nil {
//	    // Handle session service unavailability gracefully
//	    log.Printf("Session service unavailable: %v", err)
//	    return // Skip session operations
//	}
//	session.Put("key", "value")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestSessionBehavior(t *testing.T) {
//	    // Create a test session that records all operations
//	    testSession := &TestSession{
//	        data:       make(map[string]interface{}),
//	        operations: []string{},
//	    }
//
//	    // Swap the real session with test session
//	    restore := support.SwapService("session", testSession)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Session() returns testSession
//	    facades.Session().Put("test_key", "test_value")
//	    value := facades.Session().Get("test_key")
//
//	    // Verify session behavior
//	    assert.Equal(t, "test_value", value)
//	    assert.Contains(t, testSession.operations, "PUT:test_key")
//	    assert.Contains(t, testSession.operations, "GET:test_key")
//	}
//
// Container Configuration:
// Ensure the session service is properly configured in your container:
//
//	// Example session registration
//	container.Singleton("session", func() interface{} {
//	    config := session.Config{
//	        // Session driver (file, database, redis, cookie, memory)
//	        Driver: "file",
//
//	        // Session lifetime in minutes
//	        Lifetime: 120,
//
//	        // Session ID length
//	        IDLength: 40,
//
//	        // Cookie configuration
//	        Cookie: session.CookieConfig{
//	            Name:     "laravel_session",
//	            Path:     "/",
//	            Domain:   "",
//	            Secure:   true,
//	            HTTPOnly: true,
//	            SameSite: http.SameSiteStrictMode,
//	        },
//
//	        // Encryption
//	        Encrypt: true,
//
//	        // File driver specific
//	        Files: session.FileConfig{
//	            StorePath: "./storage/sessions",
//	            FileMode:  0600,
//	        },
//
//	        // Database driver specific
//	        Database: session.DatabaseConfig{
//	            Connection: "default",
//	            Table:      "sessions",
//	        },
//
//	        // Redis driver specific
//	        Redis: session.RedisConfig{
//	            Connection: "default",
//	            KeyPrefix:  "sess:",
//	        },
//
//	        // Garbage collection
//	        GarbageCollection: session.GCConfig{
//	            Probability: 0.02,
//	            Divisor:     100,
//	        },
//	    }
//
//	    sessionService, err := session.NewSessionService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create session service: %v", err)
//	    }
//
//	    return sessionService
//	})
func Session() sessionInterfaces.SessionInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "session" service from the dependency injection container
	// - Performs type assertion to SessionInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[sessionInterfaces.SessionInterface](sessionInterfaces.SESSION_TOKEN)
}

// SessionWithError provides error-safe access to the session service.
//
// This function offers the same functionality as Session() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle session service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Session() but with error handling.
//
// Returns:
//   - SessionInterface: The resolved session instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement SessionInterface
//
// Usage Examples:
//
//	// Basic error-safe session operations
//	session, err := facades.SessionWithError()
//	if err != nil {
//	    log.Printf("Session service unavailable: %v", err)
//	    return // Skip session operations
//	}
//	session.Put("user_id", 123)
//
//	// Conditional session usage
//	if session, err := facades.SessionWithError(); err == nil {
//	    // Store temporary flash message
//	    session.Flash("info", "Optional feature enabled")
//	}
func SessionWithError() (sessionInterfaces.SessionInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "session" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[sessionInterfaces.SessionInterface](sessionInterfaces.SESSION_TOKEN)
}
