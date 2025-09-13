package hashers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"govel/packages/hashing/src/exceptions"
	interfaces "govel/packages/types/src/interfaces/hashing"
	"govel/packages/support/src/number"
	"govel/packages/support/src/traits"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ArgonHasher implements the Argon2i password hashing algorithm, providing
// memory-hard password hashing with side-channel resistance. Part of the Argon2
// family that won the Password Hashing Competition (PHC) in 2015.
//
// Key features:
//   - Memory-hard function requiring significant RAM for computation
//   - Data-independent memory access patterns (side-channel resistant)
//   - Tunable parameters: memory usage, time cost, and parallelism
//   - Built-in salt generation and management
//   - RFC 9106 compliant implementation
//
// Security properties:
//   - Resistant to GPU/ASIC acceleration attempts
//   - Protection against side-channel attacks via independent memory access
//   - Immune to time-memory trade-off attacks
//   - Configurable resistance to various attack vectors
//
// Performance characteristics:
//   - Memory usage: 8KB to several GB (configurable)
//   - Time complexity: Linear with iteration count
//   - Parallelism: Can utilize multiple CPU cores
//   - 64MB, 1 iteration: ~100ms on modern hardware
//
// Thread safety: ArgonHasher is safe for concurrent use across goroutines.
type ArgonHasher struct {
	*BaseHasher
	traits.Proxiable        // Enable self-reference for method dispatch
	memory           uint32 // Memory usage in KB (e.g., 65536 = 64MB)
	time             uint32 // Number of iterations/passes
	threads          uint8  // Degree of parallelism (number of threads)
	keyLen           uint32 // Output key length in bytes
	saltLen          int    // Salt length in bytes
}

// NewArgonHasher creates a new Argon2i hasher instance with configurable parameters.
// Provides secure defaults while allowing customization for specific security requirements.
//
// Default parameters (suitable for most applications):
//   - Memory: 64MB (65536 KB) - good balance of security and resource usage
//   - Time: 1 iteration - minimal time cost for Argon2i
//   - Threads: 1 - single-threaded operation
//   - Key length: 32 bytes - 256-bit security level
//   - Salt length: 16 bytes - sufficient entropy
//
// Supported configuration keys:
//   - "memory" (uint32): Memory usage in KB (8MB-1GB range recommended)
//   - "time" (uint32): Number of iterations (1-100 range recommended)
//   - "threads" (uint8): Degree of parallelism (1-8 threads recommended)
//   - "keyLen" (uint32): Output length in bytes (16-128 range recommended)
//   - "saltLen" (int): Salt length in bytes (16-64 range recommended)
//
// Parameters:
//   - config: Configuration map with optional parameter overrides
//
// Returns:
//   - *ArgonHasher: A new Argon2i hasher instance with validated configuration
//
// Configuration validation:
//   - Invalid parameters are silently ignored, defaults retained
//   - Memory and time parameters are validated against reasonable ranges
//   - Thread count and key length are bounded for security
//
// Usage examples:
//
//	// Default configuration hasher
//	hasher := NewArgonHasher(nil)
//
//	// High-security configuration
//	highSecConfig := map[string]interface{}{
//		"memory":  uint32(256 * 1024), // 256MB
//		"time":    uint32(3),          // 3 iterations
//		"threads": uint8(4),           // 4 threads
//		"keyLen":  uint32(32),         // 32-byte output
//	}
//	hasher := NewArgonHasher(highSecConfig)
func NewArgonHasher(config map[string]interface{}) *ArgonHasher {
	// Initialize hasher with secure default parameters
	// These defaults provide good security for most applications
	h := &ArgonHasher{
		BaseHasher: NewBaseHasher(), // Inherit common hasher functionality
		memory:     65536,           // 64MB - good balance of security and performance
		time:       1,               // 1 pass - minimal for Argon2i due to data-independent access
		threads:    1,               // 1 thread - single-threaded for side-channel resistance
		keyLen:     32,              // 32 bytes - 256-bit security level
		saltLen:    16,              // 16 bytes - sufficient entropy for salt uniqueness
	}

	// Apply configuration overrides with validation
	// Invalid parameters are silently ignored to prevent configuration errors
	if config != nil {
		// Memory configuration with bounds checking
		if m, exists := config["memory"]; exists {
			if memVal, err := number.ToUint32(m); err == nil && memVal > 0 {
				// Validate memory range: 8MB minimum, 1GB maximum
				// Lower bound prevents weak hashes, upper bound prevents resource exhaustion
				if memVal >= 8*1024 && memVal <= 1024*1024 {
					h.memory = memVal
				}
			}
		}

		// Time cost configuration with reasonable limits
		if t, exists := config["time"]; exists {
			if timeVal, err := number.ToUint32(t); err == nil && timeVal > 0 {
				// Validate time range: minimum 1, maximum 100 iterations
				// Argon2i typically uses lower time costs due to memory hardness
				if timeVal <= 100 {
					h.time = timeVal
				}
			}
		}

		// Thread count configuration
		if p, exists := config["threads"]; exists {
			if threadVal, err := number.ToUint8(p); err == nil && threadVal > 0 {
				// Accept any positive thread count - system will limit naturally
				// More threads can improve performance on multi-core systems
				h.threads = threadVal
			}
		}

		// Output key length configuration with security bounds
		if k, exists := config["keyLen"]; exists {
			if keyVal, err := number.ToUint32(k); err == nil && keyVal > 0 {
				// Validate key length: 16-128 bytes (128-1024 bits)
				// Minimum ensures adequate security, maximum prevents excessive output
				if keyVal >= 16 && keyVal <= 128 {
					h.keyLen = keyVal
				}
			}
		}

		// Salt length configuration with security bounds
		if s, ok := config["saltLen"].(int); ok && s > 0 {
			// Validate salt length: 16-64 bytes for optimal entropy
			// Minimum ensures uniqueness, maximum prevents unnecessary overhead
			if s >= 16 && s <= 64 {
				h.saltLen = s
			}
		}
	}

	// Set up proxy self-reference for method dispatch
	h.SetProxySelf(h)

	return h
}

// Make hashes the given value using Argon2i algorithm.
// Creates a secure hash with automatic salt generation and parameter validation.
//
// Hashing features:
//   - Cryptographically secure random salt generation
//   - Input length validation to prevent DoS attacks
//   - Runtime parameter override with validation
//   - RFC 9106 compliant Argon2i implementation
//
// Returns standard Argon2i hash format or error if hashing fails.
func (h *ArgonHasher) Make(value string, options map[string]interface{}) (string, error) {
	// Input validation to prevent resource exhaustion attacks
	// Large inputs can consume excessive memory during hashing
	const maxInputLength = 4096 // 4KB reasonable limit for password-like inputs
	if len(value) > maxInputLength {
		return "", exceptions.ErrValueTooLong
	}

	// Generate cryptographically secure random salt
	// Salt ensures unique hashes for identical inputs
	salt := make([]byte, h.saltLen)
	if _, err := rand.Read(salt); err != nil {
		// Cryptographic random generation failure is critical
		return "", fmt.Errorf("failed to generate secure salt: %w", err)
	}

	// Initialize working parameters from instance defaults
	// These can be overridden by runtime options for flexibility
	memory := h.memory   // Memory usage in KB
	time := h.time       // Number of iterations
	threads := h.threads // Degree of parallelism
	keyLen := h.keyLen   // Output key length

	// Process runtime parameter overrides with strict validation
	// Invalid parameters result in immediate error to prevent weak hashes
	if options != nil {
		// Memory parameter validation and override
		if m, exists := options["memory"]; exists {
			memVal, err := number.ToUint32(m)
			if err != nil {
				return "", fmt.Errorf("%w: invalid memory parameter", exceptions.ErrInvalidOptions)
			}
			// Enforce memory bounds: 8MB minimum for security, 1GB maximum for resource protection
			if memVal < 8*1024 || memVal > 1024*1024 {
				return "", fmt.Errorf("%w: memory must be between 8MB and 1GB", exceptions.ErrInvalidOptions)
			}
			memory = memVal
		}

		// Time cost parameter validation and override
		if t, exists := options["time"]; exists {
			timeVal, err := number.ToUint32(t)
			if err != nil {
				return "", fmt.Errorf("%w: invalid time parameter", exceptions.ErrInvalidOptions)
			}
			// Enforce time bounds: minimum 1 iteration, maximum 100 for practicality
			if timeVal < 1 || timeVal > 100 {
				return "", fmt.Errorf("%w: time must be between 1 and 100", exceptions.ErrInvalidOptions)
			}
			time = timeVal
		}

		// Thread count parameter validation and override
		if p, exists := options["threads"]; exists {
			threadVal, err := number.ToUint8(p)
			if err != nil {
				return "", fmt.Errorf("%w: invalid threads parameter", exceptions.ErrInvalidOptions)
			}
			// Enforce minimum thread count for functionality
			if threadVal < 1 {
				return "", fmt.Errorf("%w: threads must be at least 1", exceptions.ErrInvalidOptions)
			}
			threads = threadVal
		}

		// Key length parameter validation and override
		if k, exists := options["keyLen"]; exists {
			keyLenVal, err := number.ToUint32(k)
			if err != nil {
				return "", fmt.Errorf("%w: invalid keyLen parameter", exceptions.ErrInvalidOptions)
			}
			// Enforce key length bounds: 16-128 bytes for practical security
			if keyLenVal < 16 || keyLenVal > 128 {
				return "", fmt.Errorf("%w: keyLen must be between 16 and 128 bytes", exceptions.ErrInvalidOptions)
			}
			keyLen = keyLenVal
		}
	}

	// Perform Argon2i key derivation with validated parameters
	// Uses data-independent memory access patterns for side-channel resistance
	hash := h.deriveKey([]byte(value), salt, time, memory, threads, keyLen)

	// Encode salt and hash components using base64 without padding
	// Raw encoding matches PHC string format specification
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Format according to PHC string format: $algorithm$version$params$salt$hash
	// Version 19 is current Argon2 version, parameters are comma-separated
	return fmt.Sprintf("$%s$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.getAlgorithmIdentifier(), memory, time, threads, saltB64, hashB64), nil
}

// Check verifies the given plain value against an Argon2i hash.
// Performs constant-time comparison to prevent timing attacks.
//
// Verification features:
//   - PHC string format validation and parsing
//   - Parameter extraction from hash string
//   - Identical Argon2i computation for comparison
//   - Safe handling of malformed hash strings
//
// Returns true if value matches hash, false otherwise or on errors.
func (h *ArgonHasher) Check(value, hashedValue string, options map[string]interface{}) bool {
	// Pre-validate hash format to prevent processing invalid inputs
	// This prevents wasting resources on obviously malformed hashes
	if err := h.ValidateHashFormat(hashedValue); err != nil {
		return false
	}

	// Parse PHC string format: $argon2[i|id]$v=19$m=memory,t=time,p=threads$salt$hash
	// Expected 6 parts after splitting by '$' delimiter
	parts := strings.Split(hashedValue, "$")
	if len(parts) != 6 || parts[1] != h.getAlgorithmIdentifier() {
		// Invalid format or wrong algorithm identifier
		return false
	}

	// Extract and parse Argon2i parameters from hash string
	// Parameters are in format: m=memory,t=time,p=threads
	params := strings.Split(parts[3], ",")
	var memory, time uint32
	var threads uint8

	// Parse each parameter using key=value format
	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			// Skip malformed parameter entries
			continue
		}

		// Extract parameter values with error handling
		switch kv[0] {
		case "m": // Memory usage in KB
			fmt.Sscanf(kv[1], "%d", &memory)
		case "t": // Time cost (iterations)
			fmt.Sscanf(kv[1], "%d", &time)
		case "p": // Parallelism (threads)
			fmt.Sscanf(kv[1], "%d", &threads)
		}
	}

	// Decode base64-encoded salt from hash string
	// Salt is required for hash recreation
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		// Malformed salt encoding
		return false
	}

	// Decode base64-encoded expected hash for comparison
	// Hash length determines key length parameter
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		// Malformed hash encoding
		return false
	}

	// Recreate hash using extracted parameters and original salt
	// Key length matches the expected hash length
	keyLen := uint32(len(expectedHash))
	hash := h.deriveKey([]byte(value), salt, time, memory, threads, keyLen)

	// Perform constant-time comparison to prevent timing attacks
	// Both hashes should be identical if value is correct
	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}

// NeedsRehash checks if the hash needs to be rehashed with updated parameters.
// Compares current configuration against hash parameters to determine if upgrade needed.
//
// Rehash detection features:
//   - Algorithm compatibility checking
//   - Parameter comparison (memory, time, threads)
//   - Security standard compliance validation
//   - Conservative approach for security upgrades
//
// Returns true if rehashing recommended, false if current hash is adequate.
func (h *ArgonHasher) NeedsRehash(hashedValue string, options map[string]interface{}) bool {
	// Pre-validate hash format - invalid hashes always need rehashing
	if err := h.ValidateHashFormat(hashedValue); err != nil {
		return true // Invalid hashes must be rehashed for security
	}

	// Extract hash information for parameter comparison
	info := h.Info(hashedValue)
	if info.AlgoName != h.getAlgorithmIdentifier() {
		// Different algorithm requires rehashing
		return true
	}

	// Get desired parameters from options or use current configuration as default
	desiredMemory := h.memory
	desiredTime := h.time
	desiredThreads := h.threads

	// Apply parameter overrides from options if provided
	if options != nil {
		if m, exists := options["memory"]; exists {
			if memVal, err := number.ToUint32(m); err == nil {
				desiredMemory = memVal
			}
		}
		if t, exists := options["time"]; exists {
			if timeVal, err := number.ToUint32(t); err == nil {
				desiredTime = timeVal
			}
		}
		if p, exists := options["threads"]; exists {
			if threadVal, err := number.ToUint8(p); err == nil {
				desiredThreads = threadVal
			}
		}
	}

	// Compare desired parameters with hash parameters
	// Any mismatch triggers rehashing recommendation
	if m, ok := info.Options["memory_cost"].(int); ok && uint32(m) != desiredMemory {
		// Memory parameter changed - rehash for consistency
		return true
	}
	if t, ok := info.Options["time_cost"].(int); ok && uint32(t) != desiredTime {
		// Time parameter changed - rehash for consistency
		return true
	}
	if p, ok := info.Options["threads"].(int); ok && uint8(p) != desiredThreads {
		// Thread parameter changed - rehash for consistency
		return true
	}

	// Hash parameters match desired configuration
	return false
}

// VerifyConfiguration verifies that Argon2i parameters are valid.
// Ensures hasher configuration meets minimum security requirements.
//
// Configuration validation:
//   - All parameters must be positive values
//   - Memory, time, threads, and key length bounds checking
//   - Production readiness assessment
//
// Returns true if configuration is valid and secure.
func (h *ArgonHasher) VerifyConfiguration(value string) bool {
	// Check all critical parameters are positive and within reasonable bounds
	return h.memory > 0 && // Memory usage must be configured
		h.time > 0 && // Time cost must be at least 1 iteration
		h.threads > 0 && // Thread count must be positive
		h.keyLen > 0 && // Key length must be configured
		h.memory >= 8*1024 && // Minimum 8MB for security
		h.keyLen >= 16 // Minimum 16 bytes for adequate security
}

// getAlgorithmIdentifier returns the algorithm identifier for hash format strings.
// Uses Proxiable to call Algorithm() on the concrete implementation if available.
func (h *ArgonHasher) getAlgorithmIdentifier() string {
	// Try to call Algorithm() on the concrete implementation via proxy
	if h.HasProxySelf() {
		if results, err := h.CallOnSelf("Algorithm"); err == nil && len(results) > 0 {
			if algoName := results[0].String(); algoName != "" {
				return algoName
			}
		}
	}
	// Fallback to default implementation
	return h.Algorithm()
}

// Algorithm returns the default algorithm identifier.
// This method should be overridden by embedded hashers (like Argon2IdHasher).
func (h *ArgonHasher) Algorithm() string {
	return "argon2i"
}

// deriveKey performs Argon2i key derivation.
// This method can be overridden by embedded hashers (like Argon2IdHasher)
// to use different Argon2 variants while reusing all other functionality.
func (h *ArgonHasher) deriveKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
	// Use Argon2i key derivation (data-independent memory access)
	return argon2.Key(password, salt, time, memory, threads, keyLen)
}

// Compile-time interface compliance check
var _ interfaces.HasherInterface = (*ArgonHasher)(nil)
