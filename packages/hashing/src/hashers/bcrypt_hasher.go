package hashers

import (
	"govel/packages/hashing/src/exceptions"
	interfaces "govel/packages/types/src/interfaces/hashing"
	"govel/packages/support/src/number"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements the bcrypt password hashing algorithm, providing
// a secure and widely-supported solution for password storage and verification.
// Built on the Blowfish cipher with adaptive cost scaling for future-proofing.
//
// Key features:
//   - Adaptive cost parameter (4-31 range, exponential scaling)
//   - Built-in salt generation and management
//   - Constant-time verification to prevent timing attacks
//   - Cross-platform compatibility and wide library support
//   - Proven security track record since 1999
//
// Security properties:
//   - Resistant to rainbow table attacks via salting
//   - Configurable work factor to counter hardware improvements
//   - No vulnerability to length-extension attacks
//   - Suitable for general-purpose password hashing
//
// Performance characteristics:
//   - Cost 10: ~70ms on modern hardware
//   - Cost 12: ~280ms on modern hardware
//   - Each cost increment doubles computation time
//   - Lower memory usage compared to Argon2 variants
//
// Thread safety: BcryptHasher is safe for concurrent use across goroutines.
type BcryptHasher struct {
	*BaseHasher
	cost int // Cost parameter (4-31), higher values increase security and computation time
}

// NewBcryptHasher creates a new bcrypt hasher instance with configuration from a map.
// The configuration controls the computational complexity and security level of the hashing.
//
// Configuration parameters:
//   - "cost" (int): The bcrypt cost parameter (4-31 range)
//   - Minimum: 4 (very fast, low security - not recommended for production)
//   - Default: 10 (good balance for most applications)
//   - Recommended: 12+ (higher security for sensitive applications)
//   - Maximum: 31 (extremely slow, maximum security)
//
// Parameters:
//   - config: Configuration map with bcrypt parameters. Invalid values use defaults.
//
// Returns:
//   - *BcryptHasher: A new bcrypt hasher instance ready for password operations
//
// Usage examples:
//
//	// Standard security hasher
//	config := map[string]interface{}{"cost": 12}
//	hasher := NewBcryptHasher(config)
//
//	// High security hasher (slower)
//	config := map[string]interface{}{"cost": 15}
//	secureHasher := NewBcryptHasher(config)
//
//	// Default cost hasher (nil or empty config uses defaults)
//	defaultHasher := NewBcryptHasher(nil)
func NewBcryptHasher(config map[string]interface{}) *BcryptHasher {
	// Start with default cost
	cost := bcrypt.DefaultCost

	// Override with config if provided and valid
	if config != nil {
		if c, exists := config["cost"]; exists {
			if costVal, err := number.ToInt(c); err == nil {
				// Validate cost is within acceptable range
				if costVal >= bcrypt.MinCost && costVal <= bcrypt.MaxCost {
					cost = costVal
				}
			}
		}
	}

	return &BcryptHasher{
		BaseHasher: NewBaseHasher(),
		cost:       cost,
	}
}

// Make hashes the given password value using the bcrypt algorithm with comprehensive
// input validation and configurable parameters. This method provides secure password
// hashing with protection against common attack vectors.
//
// Input validation:
//   - Rejects passwords longer than 72 bytes (bcrypt limitation)
//   - Validates cost parameters from options
//   - Returns ErrValueTooLong for oversized inputs
//   - Returns ErrInvalidOptions for invalid cost values
//
// Supported options:
//   - "cost" (int): Override the instance cost parameter (4-31 range)
//
// Parameters:
//   - value: The password string to hash (max 72 bytes)
//   - options: Optional parameters map for cost override
//
// Returns:
//   - string: The bcrypt hash in standard format ($2y$cost$salt+hash)
//   - error: ErrValueTooLong, ErrInvalidOptions, or bcrypt generation errors
//
// Security considerations:
//   - Automatically generates cryptographically secure salt
//   - Uses constant-time operations internally
//   - Cost parameter balances security vs. performance
//
// Example usage:
//
//	// Standard hashing
//	hash, err := hasher.Make("password123", nil)
//
//	// High security hashing
//	hash, err := hasher.Make("password123", map[string]interface{}{
//		"cost": 15,
//	})
func (h *BcryptHasher) Make(value string, options map[string]interface{}) (string, error) {
	// Validate input length - bcrypt effectively limits password length to 72 bytes
	if len(value) > 72 {
		return "", exceptions.ErrValueTooLong
	}

	cost := h.cost

	// Validate and override cost from options if provided
	if c, exists := options["cost"]; exists {
		costInt, err := number.ToInt(c)
		if err != nil {
			return "", exceptions.ErrInvalidOptions
		}
		if costInt < bcrypt.MinCost || costInt > bcrypt.MaxCost {
			return "", exceptions.ErrInvalidOptions
		}
		cost = costInt
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(value), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// Check verifies a plain text password against a bcrypt hash using constant-time
// comparison to prevent timing attacks. Includes comprehensive input validation
// to ensure hash integrity before verification.
//
// Verification process:
//  1. Validates hash format using ValidateHashFormat()
//  2. Performs constant-time password comparison
//  3. Returns false for any validation or verification failures
//
// Security features:
//   - Constant-time comparison prevents timing attacks
//   - Hash format validation prevents malformed input processing
//   - Safe handling of invalid hashes without exceptions
//
// Parameters:
//   - value: The plain text password to verify
//   - hashedValue: The bcrypt hash to verify against
//   - options: Reserved for future use (currently unused)
//
// Returns:
//   - bool: true if password matches hash, false otherwise
//
// Error handling:
//   - Invalid hash formats return false (no exceptions)
//   - Verification failures return false
//   - All errors are handled gracefully
//
// Usage examples:
//
//	// Basic password verification
//	if hasher.Check("password123", storedHash, nil) {
//		fmt.Println("Password verified successfully")
//	} else {
//		fmt.Println("Invalid password")
//	}
func (h *BcryptHasher) Check(value, hashedValue string, options map[string]interface{}) bool {
	// Validate hash format first
	if err := h.ValidateHashFormat(hashedValue); err != nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}

// NeedsRehash determines if a bcrypt hash should be regenerated with updated
// cost parameters. This method enables password hash migration and security
// upgrades without forcing user re-authentication.
//
// Rehashing scenarios:
//   - Hash format validation fails (corrupted or invalid hashes)
//   - Current hash cost differs from required cost
//   - Hash was created with outdated security parameters
//   - Algorithm migration or security policy updates
//
// Evaluation process:
//  1. Validates hash format (returns true if invalid)
//  2. Extracts current cost from hash
//  3. Compares with required cost from options or instance default
//  4. Returns true if costs don't match
//
// Supported options:
//   - "cost" (int): Target cost parameter for comparison
//
// Parameters:
//   - hashedValue: The bcrypt hash to evaluate
//   - options: Optional parameters with target cost
//
// Returns:
//   - bool: true if hash should be regenerated, false if current hash is acceptable
//
// Usage examples:
//
//	// Check if hash needs security upgrade
//	if hasher.NeedsRehash(userHash, map[string]interface{}{"cost": 15}) {
//		// Regenerate hash with higher cost
//		newHash, _ := hasher.Make(password, map[string]interface{}{"cost": 15})
//		// Update stored hash
//	}
func (h *BcryptHasher) NeedsRehash(hashedValue string, options map[string]interface{}) bool {
	// Validate hash format first
	if err := h.ValidateHashFormat(hashedValue); err != nil {
		return true // Invalid hashes need rehashing
	}

	info := h.Info(hashedValue)
	currentCost, ok := info.Options["cost"].(int)
	if !ok {
		return true
	}

	requiredCost := h.cost
	if c, ok := options["cost"].(int); ok {
		requiredCost = c
	}

	return currentCost != requiredCost
}

// VerifyConfiguration validates the bcrypt hasher's configuration parameters
// to ensure they are within acceptable ranges and suitable for secure operation.
// This method provides configuration validation for security auditing and setup.
//
// Configuration checks:
//   - Cost parameter within valid range (bcrypt.MinCost to bcrypt.MaxCost)
//   - Ensures configuration meets minimum security requirements
//   - Validates against bcrypt library constraints
//
// Security validation:
//   - Minimum cost: 4 (very fast, not recommended for production)
//   - Maximum cost: 31 (extremely slow, maximum theoretical security)
//   - Recommended minimum: 10+ for production systems
//
// Parameters:
//   - value: Test value for configuration validation (unused in bcrypt implementation)
//
// Returns:
//   - bool: true if configuration is valid and secure, false otherwise
//
// Usage examples:
//
//	// Validate hasher configuration during setup
//	if !hasher.VerifyConfiguration("") {
//		log.Error("Invalid bcrypt configuration detected")
//		return errors.New("hasher configuration failed validation")
//	}
//
//	// System health check
//	func (s *SecurityService) ValidateHashers() error {
//		for name, hasher := range s.hashers {
//			if !hasher.VerifyConfiguration("") {
//				return fmt.Errorf("hasher %s failed configuration check", name)
//			}
//		}
//		return nil
//	}
func (h *BcryptHasher) VerifyConfiguration(value string) bool {
	// For bcrypt, we just need to ensure the cost is within bounds
	return h.cost >= bcrypt.MinCost && h.cost <= bcrypt.MaxCost
}

// Compile-time interface compliance check
var _ interfaces.HasherInterface = (*BcryptHasher)(nil)
