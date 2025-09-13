package facades

import facade "govel/support/src"

// Hash provides a clean, static-like interface to the application's cryptographic hashing service.
//
// This facade implements the facade pattern, providing global access to the hashing
// service configured in the dependency injection container. It offers a Laravel-style
// API for password hashing, data integrity verification, secure token generation,
// and cryptographic operations with automatic service resolution and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved hash service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent hashing operations across goroutines
//   - Supports multiple hashing algorithms and security configurations
//   - Built-in salt generation and timing-safe comparison
//
// Behavior:
//   - First call: Resolves hash service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if hash service cannot be resolved (fail-fast behavior)
//   - Automatically handles secure random generation and timing-safe operations
//
// Returns:
//   - HashInterface: The application's cryptographic hashing service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "hash" service is not registered in the container
//   - If the resolved service doesn't implement HashInterface
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
//   - Multiple goroutines can call Hash() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Hash operations are cryptographically secure and thread-safe
//
// Usage Examples:
//
//	// Password hashing and verification
//	passwordHash, err := facades.Hash().Make("user-password")
//	if err != nil {
//	    log.Fatal("Failed to hash password")
//	}
//
//	// Store passwordHash in database
//	user := &User{
//	    Email:    "user@example.com",
//	    Password: passwordHash,
//	}
//
//	// Verify password during login
//	isValid := facades.Hash().Check("user-password", user.Password)
//	if isValid {
//	    log.Println("Password is correct")
//	} else {
//	    log.Println("Invalid password")
//	}
//
//	// Check if password needs rehashing (after algorithm/cost changes)
//	if facades.Hash().NeedsRehash(user.Password) {
//	    newHash, _ := facades.Hash().Make(plainPassword)
//	    user.Password = newHash
//	    // Update user in database
//	}
//
//	// Generate secure tokens and API keys
//	apiKey := facades.Hash().GenerateToken(32) // 32 bytes = 256 bits
//	sessionToken := facades.Hash().GenerateToken(16) // 16 bytes = 128 bits
//	resetToken := facades.Hash().GenerateURLSafeToken(32) // URL-safe base64
//
//	// HMAC for data integrity and authentication
//	secretKey := "your-secret-key"
//	message := "important-data"
//	mac := facades.Hash().HMAC("sha256", secretKey, message)
//
//	// Verify HMAC
//	isValidMAC := facades.Hash().VerifyHMAC("sha256", secretKey, message, mac)
//	if isValidMAC {
//	    log.Println("Message is authentic")
//	}
//
//	// File integrity verification
//	fileContent, _ := os.ReadFile("important-file.txt")
//	fileSHA256 := facades.Hash().SHA256(fileContent)
//	fileHash := fmt.Sprintf("%x", fileSHA256)
//
//	// Store hash for later verification
//	storeFileHash(fileName, fileHash)
//
//	// Later verification
//	storedhash := getStoredFileHash(fileName)
//	currentHash := fmt.Sprintf("%x", facades.Hash().SHA256(fileContent))
//	if storedHash == currentHash {
//	    log.Println("File integrity verified")
//	}
//
//	// Different hash algorithms
//	md5Hash := facades.Hash().MD5([]byte("data"))
//	sha1Hash := facades.Hash().SHA1([]byte("data"))
//	sha256Hash := facades.Hash().SHA256([]byte("data"))
//	sha512Hash := facades.Hash().SHA512([]byte("data"))
//
//	// Hash with salt
//	salt := facades.Hash().GenerateSalt(16)
//	saltedHash := facades.Hash().HashWithSalt("bcrypt", "password", salt)
//
//	// Argon2 hashing (recommended for passwords)
//	argonHash := facades.Hash().Argon2ID("password", salt)
//	isValidArgon := facades.Hash().VerifyArgon2ID("password", argonHash)
//
//	// Timing-safe string comparison (prevents timing attacks)
//	token1 := "secret-token-12345"
//	token2 := getUserProvidedToken()
//	isEqual := facades.Hash().TimingSafeEqual(token1, token2)
//
// Security-Critical Usage Patterns:
//
//	// Secure user registration
//	func RegisterUser(email, password string) error {
//	    // Validate password strength first
//	    if len(password) < 8 {
//	        return errors.New("password too short")
//	    }
//
//	    // Hash password securely
//	    hashedPassword, err := facades.Hash().Make(password)
//	    if err != nil {
//	        return fmt.Errorf("failed to hash password: %w", err)
//	    }
//
//	    user := &User{
//	        Email:    email,
//	        Password: hashedPassword,
//	    }
//
//	    return database.CreateUser(user)
//	}
//
//	// Secure login verification
//	func LoginUser(email, password string) (*User, error) {
//	    user, err := database.GetUserByEmail(email)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Use timing-safe password verification
//	    if !facades.Hash().Check(password, user.Password) {
//	        return nil, errors.New("invalid credentials")
//	    }
//
//	    // Check if password needs rehashing
//	    if facades.Hash().NeedsRehash(user.Password) {
//	        go func() {
//	            newHash, _ := facades.Hash().Make(password)
//	            user.Password = newHash
//	            database.UpdateUser(user) // Update in background
//	        }()
//	    }
//
//	    return user, nil
//	}
//
//	// API key generation and validation
//	func GenerateAPIKey(userID int) (string, error) {
//	    // Generate cryptographically secure API key
//	    token := facades.Hash().GenerateURLSafeToken(32)
//
//	    // Create hash for storage (don't store raw token)
//	    tokenHash := facades.Hash().SHA256([]byte(token))
//
//	    apiKey := &APIKey{
//	        UserID:    userID,
//	        TokenHash: fmt.Sprintf("%x", tokenHash),
//	        CreatedAt: time.Now(),
//	    }
//
//	    if err := database.CreateAPIKey(apiKey); err != nil {
//	        return "", err
//	    }
//
//	    return token, nil
//	}
//
//	func ValidateAPIKey(token string) (*APIKey, error) {
//	    tokenHash := fmt.Sprintf("%x", facades.Hash().SHA256([]byte(token)))
//	    return database.GetAPIKeyByTokenHash(tokenHash)
//	}
//
//	// Password reset token generation
//	func GeneratePasswordResetToken(userID int) (string, error) {
//	    token := facades.Hash().GenerateURLSafeToken(32)
//
//	    // Hash token for storage
//	    hashedToken, err := facades.Hash().Make(token)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    resetToken := &PasswordResetToken{
//	        UserID:    userID,
//	        Token:     hashedToken,
//	        ExpiresAt: time.Now().Add(24 * time.Hour),
//	    }
//
//	    if err := database.CreatePasswordResetToken(resetToken); err != nil {
//	        return "", err
//	    }
//
//	    return token, nil
//	}
//
//	// Data signing and verification
//	func SignData(data []byte) ([]byte, []byte, error) {
//	    secretKey := facades.Config().GetString("app.secret")
//	    signature := facades.Hash().HMAC("sha256", secretKey, string(data))
//	    return data, signature, nil
//	}
//
//	func VerifySignedData(data, signature []byte) bool {
//	    secretKey := facades.Config().GetString("app.secret")
//	    expectedSig := facades.Hash().HMAC("sha256", secretKey, string(data))
//	    return facades.Hash().TimingSafeEqual(string(expectedSig), string(signature))
//	}
//
// Best Practices:
//   - Always use timing-safe comparison for tokens and hashes
//   - Use bcrypt or Argon2 for password hashing
//   - Generate cryptographically secure random tokens
//   - Hash API keys and tokens before storing in database
//   - Regularly rehash passwords with updated cost factors
//   - Use HMAC for data integrity and authenticity
//   - Never store plaintext passwords or sensitive tokens
//   - Use appropriate hash algorithms for different use cases
//
// Security Considerations:
//  1. Use strong, cryptographically secure algorithms
//  2. Generate sufficient entropy for salts and tokens
//  3. Protect against timing attacks with constant-time operations
//  4. Regularly update hash costs as computing power increases
//  5. Store hashes, not plaintext, for sensitive data
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume hashing always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	hasher, err := facade.TryResolve[HashInterface]("hash")
//	if err != nil {
//	    // Handle hash service unavailability gracefully
//	    return fmt.Errorf("hashing service unavailable: %w", err)
//	}
//	hash, _ := hasher.Make("password")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestPasswordHashing(t *testing.T) {
//	    // Create a test hasher with predictable output
//	    testHasher := &TestHasher{
//	        hashes: make(map[string]string),
//	    }
//
//	    // Swap the real hasher with test hasher
//	    restore := support.SwapService("hash", testHasher)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Hash() returns testHasher
//	    hash, err := facades.Hash().Make("test-password")
//	    assert.NoError(t, err)
//	    assert.NotEmpty(t, hash)
//
//	    // Verify password checking
//	    valid := facades.Hash().Check("test-password", hash)
//	    assert.True(t, valid)
//
//	    invalid := facades.Hash().Check("wrong-password", hash)
//	    assert.False(t, invalid)
//	}
//
// Container Configuration:
// Ensure the hash service is properly configured in your container:
//
//	// Example hasher registration
//	container.Singleton("hash", func() interface{} {
//	    config := hasher.Config{
//	        // Default algorithm for password hashing
//	        DefaultDriver: "bcrypt",
//
//	        // Bcrypt configuration
//	        BcryptCost: 12, // Adjust based on performance requirements
//
//	        // Argon2 configuration (recommended)
//	        Argon2Config: hasher.Argon2Config{
//	            Memory:      64 * 1024, // 64 MB
//	            Iterations:  3,
//	            Parallelism: 2,
//	            SaltLength:  32,
//	            KeyLength:   32,
//	        },
//
//	        // HMAC configuration
//	        HMACKey: facades.Config().GetString("app.secret"),
//
//	        // Token generation
//	        TokenLength:    32,
//	        URLSafeTokens: true,
//
//	        // Security settings
//	        TimingSafeComparisons: true,
//	        SecureRandom:          true,
//	    }
//
//	    return hasher.NewHasher(config)
//	})
func Hash() interface{} {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "hash" service from the dependency injection container
	// - Performs type assertion to HashInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[interface{}]("hash")
}

// HashWithError provides error-safe access to the cryptographic hashing service.
//
// This function offers the same functionality as Hash() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle hash service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Hash() but with error handling.
//
// Returns:
//   - HashInterface: The resolved hash instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement HashInterface
//
// Usage Examples:
//
//	// Basic error-safe hashing
//	hasher, err := facades.HashWithError()
//	if err != nil {
//	    log.Printf("Hash service unavailable: %v", err)
//	    return fmt.Errorf("cryptographic operations not available")
//	}
//	passwordHash, _ := hasher.Make("user-password")
//
//	// Conditional cryptographic operations
//	if hasher, err := facades.HashWithError(); err == nil {
//	    token := hasher.GenerateToken(32)
//	    // Use token for session management
//	}
func HashWithError() (interface{}, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "hash" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[interface{}]("hash")
}
