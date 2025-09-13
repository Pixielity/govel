package facades

import (
	interfaces "govel/packages/types/src/interfaces/encryption"
	facade "govel/packages/support/src"
)

// Crypt provides a clean, static-like interface to the application's cryptography service.
//
// This facade implements the facade pattern, providing global access to the cryptographic
// service configured in the dependency injection container. It offers a Laravel-style
// API for encryption, decryption, hashing, and other cryptographic operations with
// automatic service resolution, key management, and security best practices.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved cryptographic service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent cryptographic operations across goroutines
//   - Supports multiple encryption ciphers (AES-256-GCM, ChaCha20-Poly1305, etc.)
//   - Built-in key rotation and secure key management
//
// Behavior:
//   - First call: Resolves crypt service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if crypt service cannot be resolved (fail-fast behavior)
//   - Automatically handles encryption keys, initialization vectors, and authentication
//
// Returns:
//   - CryptInterface: The application's cryptographic service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "crypt" service is not registered in the container
//   - If the resolved service doesn't implement CryptInterface
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
//   - Multiple goroutines can call Crypt() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Cryptographic operations are thread-safe and stateless
//
// Usage Examples:
//
//	// Basic encryption and decryption
//	plaintext := "sensitive user data"
//	encrypted, err := facades.Crypt().Encrypt(plaintext)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Encrypted: %s\n", encrypted)
//
//	decrypted, err := facades.Crypt().Decrypt(encrypted)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Decrypted: %s\n", decrypted)
//
//	// Encrypt with custom payload
//	payload := map[string]interfaces.EncrypterInterface{
//	    "user_id": 123,
//	    "email":   "user@example.com",
//	    "expires": time.Now().Add(24 * time.Hour).Unix(),
//	}
//
//	encryptedPayload, err := facades.Crypt().EncryptString(payload)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Decrypt payload back to map
//	var decryptedPayload map[string]interfaces.EncrypterInterface
//	err = facades.Crypt().DecryptString(encryptedPayload, &decryptedPayload)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Password hashing (one-way)
//	password := "user_password_123"
//	hashed, err := facades.Crypt().Hash(password)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Hashed password: %s\n", hashed)
//
//	// Password verification
//	if facades.Crypt().CheckHash(password, hashed) {
//	    fmt.Println("Password is correct")
//	} else {
//	    fmt.Println("Invalid password")
//	}
//
//	// Generate secure random strings
//	apiKey := facades.Crypt().GenerateRandomString(32)
//	sessionToken := facades.Crypt().GenerateRandomString(64)
//	csrfToken := facades.Crypt().GenerateRandomString(40)
//
//	fmt.Printf("API Key: %s\n", apiKey)
//	fmt.Printf("Session Token: %s\n", sessionToken)
//	fmt.Printf("CSRF Token: %s\n", csrfToken)
//
//	// Generate cryptographically secure random bytes
//	randomBytes, err := facades.Crypt().GenerateRandomBytes(32)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Random bytes: %x\n", randomBytes)
//
//	// HMAC signing and verification
//	message := "important message to sign"
//	signature := facades.Crypt().Sign(message)
//	fmt.Printf("Signature: %s\n", signature)
//
//	if facades.Crypt().VerifySignature(message, signature) {
//	    fmt.Println("Signature is valid")
//	} else {
//	    fmt.Println("Invalid signature")
//	}
//
//	// JWT-style tokens with expiration
//	claims := map[string]interfaces.EncrypterInterface{
//	    "user_id": 123,
//	    "role":    "admin",
//	    "exp":     time.Now().Add(2 * time.Hour).Unix(),
//	}
//
//	token, err := facades.Crypt().CreateToken(claims)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Token: %s\n", token)
//
//	// Verify and parse token
//	verifiedClaims, err := facades.Crypt().VerifyToken(token)
//	if err != nil {
//	    log.Printf("Token verification failed: %v", err)
//	} else {
//	    fmt.Printf("User ID: %v\n", verifiedClaims["user_id"])
//	    fmt.Printf("Role: %v\n", verifiedClaims["role"])
//	}
//
//	// File encryption
//	plaintextFile := "/path/to/sensitive/file.txt"
//	encryptedFile := "/path/to/encrypted/file.enc"
//
//	err = facades.Crypt().EncryptFile(plaintextFile, encryptedFile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	decryptedFile := "/path/to/decrypted/file.txt"
//	err = facades.Crypt().DecryptFile(encryptedFile, decryptedFile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Key derivation for passwords
//	password := "user_password"
//	salt := facades.Crypt().GenerateRandomBytes(32)
//	derivedKey := facades.Crypt().DeriveKey(password, salt, 32) // 32 bytes = 256 bits
//	fmt.Printf("Derived key: %x\n", derivedKey)
//
// Advanced Cryptographic Patterns:
//
//	// Envelope encryption (encrypt data key with master key)
//	func EncryptLargeData(data []byte) ([]byte, error) {
//	    // Generate a random data encryption key
//	    dataKey, err := facades.Crypt().GenerateRandomBytes(32)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Encrypt the actual data with the data key
//	    encryptedData, err := facades.Crypt().EncryptWithKey(data, dataKey)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Encrypt the data key with the master key
//	    encryptedKey, err := facades.Crypt().Encrypt(dataKey)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Combine encrypted key and encrypted data
//	    envelope := map[string]interfaces.EncrypterInterface{
//	        "key":  encryptedKey,
//	        "data": encryptedData,
//	    }
//
//	    return json.Marshal(envelope)
//	}
//
//	// Database field encryption
//	type User struct {
//	    ID    int    `json:"id"`
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	    SSN   string `json:"ssn,omitempty"`
//	}
//
//	func (u *User) EncryptSensitiveFields() error {
//	    if u.SSN != "" {
//	        encrypted, err := facades.Crypt().Encrypt(u.SSN)
//	        if err != nil {
//	            return err
//	        }
//	        u.SSN = encrypted
//	    }
//	    return nil
//	}
//
//	func (u *User) DecryptSensitiveFields() error {
//	    if u.SSN != "" {
//	        decrypted, err := facades.Crypt().Decrypt(u.SSN)
//	        if err != nil {
//	            return err
//	        }
//	        u.SSN = decrypted
//	    }
//	    return nil
//	}
//
//	// Secure session data encryption
//	func EncryptSessionData(sessionData map[string]interfaces.EncrypterInterface) (string, error) {
//	    // Add timestamp and nonce for security
//	    sessionData["timestamp"] = time.Now().Unix()
//	    sessionData["nonce"] = facades.Crypt().GenerateRandomString(16)
//
//	    // Convert to JSON and encrypt
//	    jsonData, err := json.Marshal(sessionData)
//	    if err != nil {
//	        return "", err
//	    }
//
//	    return facades.Crypt().Encrypt(string(jsonData))
//	}
//
//	// API request signing
//	func SignAPIRequest(method, url string, body []byte, timestamp int64) string {
//	    message := fmt.Sprintf("%s\n%s\n%s\n%d", method, url, string(body), timestamp)
//	    return facades.Crypt().Sign(message)
//	}
//
// Best Practices:
//   - Always use authenticated encryption (GCM mode) for confidentiality and integrity
//   - Never reuse initialization vectors (IVs) or nonces
//   - Use secure random number generation for all cryptographic material
//   - Implement proper key rotation and key management
//   - Use constant-time operations to prevent timing attacks
//   - Validate all inputs before cryptographic operations
//   - Use appropriate key derivation functions (PBKDF2, Argon2, scrypt)
//   - Store keys securely (HSM, key management service, or encrypted storage)
//
// Security Considerations:
//   - Protect encryption keys with proper access controls
//   - Use hardware security modules (HSMs) for production key storage
//   - Implement proper key rotation schedules
//   - Monitor for cryptographic failures and attacks
//   - Use timing-safe comparison functions
//   - Implement proper error handling without information leakage
//   - Use well-vetted cryptographic libraries and algorithms
//   - Regularly update cryptographic dependencies
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume cryptographic service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	crypt, err := facade.Resolve[CryptInterface]("crypt")
//	if err != nil {
//	    // Handle crypt service unavailability gracefully
//	    return plaintext, fmt.Errorf("encryption unavailable: %w", err)
//	}
//	encrypted, err := crypt.Encrypt(plaintext)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestEncryption(t *testing.T) {
//	    // Create a test crypto service with known keys
//	    testCrypt := &TestCrypt{
//	        encryptionKey: []byte("test-key-32-bytes-long-for-aes"),
//	        signingKey:    []byte("test-signing-key-for-hmac-ops"),
//	    }
//
//	    // Swap the real crypt with test crypt
//	    restore := support.SwapService("crypt", testCrypt)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Crypt() returns testCrypt
//	    plaintext := "test data"
//	    encrypted, err := facades.Crypt().Encrypt(plaintext)
//	    require.NoError(t, err)
//
//	    decrypted, err := facades.Crypt().Decrypt(encrypted)
//	    require.NoError(t, err)
//	    assert.Equal(t, plaintext, decrypted)
//	}
//
// Container Configuration:
// Ensure the cryptographic service is properly configured in your container:
//
//	// Example crypt registration
//	container.Singleton("crypt", func() interfaces.EncrypterInterface {
//	    config := crypt.Config{
//	        // Encryption configuration
//	        Cipher:     "aes-256-gcm",        // or "chacha20-poly1305"
//	        Key:        os.Getenv("APP_KEY"), // 32-byte base64 encoded key
//
//	        // Key derivation configuration
//	        KDF: crypt.KDFConfig{
//	            Algorithm: "argon2id",   // or "pbkdf2", "scrypt"
//	            Memory:    64 * 1024,    // 64MB for Argon2
//	            Time:      3,            // 3 iterations
//	            Threads:   4,            // 4 parallel threads
//	            KeyLength: 32,           // 32 bytes output
//	        },
//
//	        // Hashing configuration
//	        Hash: crypt.HashConfig{
//	            Algorithm: "argon2id",   // for password hashing
//	            Memory:    64 * 1024,    // 64MB
//	            Time:      3,            // 3 iterations
//	            Threads:   4,            // 4 threads
//	            SaltLength: 16,          // 16 bytes salt
//	        },
//
//	        // Signing configuration
//	        Signing: crypt.SigningConfig{
//	            Algorithm: "hmac-sha256",
//	            Key:       os.Getenv("APP_SIGNING_KEY"),
//	        },
//
//	        // Token configuration
//	        Token: crypt.TokenConfig{
//	            Algorithm:    "HS256",              // or "RS256" for RSA
//	            Secret:       os.Getenv("JWT_SECRET"),
//	            DefaultTTL:   time.Hour * 24,       // 24 hours
//	            Issuer:       "myapp",
//	            ValidateExp:  true,
//	            ValidateNbf:  true,
//	        },
//	    }
//
//	    cryptService, err := crypt.NewCryptService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create crypt service: %v", err)
//	    }
//
//	    return cryptService
//	})
func Crypt() interfaces.EncrypterInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "crypt" service from the dependency injection container
	// - Performs type assertion to CryptInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[interfaces.EncrypterInterface](interfaces.ENCRYPTION_TOKEN)
}

// CryptWithError provides error-safe access to the cryptographic service.
//
// This function offers the same functionality as Crypt() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle cryptographic service unavailability gracefully.
//
// This is a convenience wrapper around facade.Resolve() that provides
// the same caching and performance benefits as Crypt() but with error handling.
//
// Returns:
//   - CryptInterface: The resolved crypt instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement CryptInterface
//
// Usage Examples:
//
//	// Basic error-safe encryption
//	crypt, err := facades.CryptWithError()
//	if err != nil {
//	    log.Printf("Crypto service unavailable: %v", err)
//	    return plaintext, fmt.Errorf("encryption disabled")
//	}
//	encrypted, err := crypt.Encrypt(plaintext)
//
//	// Conditional cryptographic operations
//	if crypt, err := facades.CryptWithError(); err == nil {
//	    // Perform optional encryption
//	    if sensitiveData != "" {
//	        encrypted, _ := crypt.Encrypt(sensitiveData)
//	        saveEncryptedData(encrypted)
//	    }
//	}
//
//	// Health check pattern
//	func CheckCryptHealth() error {
//	    crypt, err := facades.CryptWithError()
//	    if err != nil {
//	        return fmt.Errorf("crypto service unavailable: %w", err)
//	    }
//
//	    // Test basic crypto functionality
//	    testData := "health-check"
//	    encrypted, err := crypt.Encrypt(testData)
//	    if err != nil {
//	        return fmt.Errorf("encryption failed: %w", err)
//	    }
//
//	    decrypted, err := crypt.Decrypt(encrypted)
//	    if err != nil || decrypted != testData {
//	        return fmt.Errorf("crypto service not working properly")
//	    }
//
//	    return nil
//	}
func CryptWithError() (interfaces.EncrypterInterface, error) {
	// Use facade.Resolve() for error-return behavior:
	// - Resolves "crypt" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[interfaces.EncrypterInterface]("crypt")
}
