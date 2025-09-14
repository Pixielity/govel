package hashers

import (
	interfaces "govel/types/interfaces/hashing"

	"golang.org/x/crypto/argon2"
)

// Argon2IdHasher implements the Argon2id password hashing algorithm.
// Extends ArgonHasher and only overrides the key derivation method to use Argon2id
// instead of Argon2i, following the DRY principle.
//
// Key differences from Argon2i:
//   - Uses hybrid memory access pattern (data-independent + data-dependent)
//   - Better time-memory trade-off resistance than Argon2i
//   - Recommended variant for new applications (RFC 9106)
//   - Typically uses higher time costs (2+ iterations vs 1 for Argon2i)
//
// Security advantages:
//   - Winner of Password Hashing Competition (PHC) 2015
//   - Best overall resistance to multiple attack vectors
//   - OWASP recommended for new password hashing implementations
//   - Future-proof design adapting to hardware improvements
//
// Thread safety: Argon2IdHasher is safe for concurrent use across goroutines.
type Argon2IdHasher struct {
	*ArgonHasher // Embed ArgonHasher to inherit all functionality
}

// NewArgon2IdHasher creates a new Argon2id hasher instance with Argon2id-specific defaults.
// Uses ArgonHasher as the base with optimized parameters for the hybrid Argon2id algorithm.
//
// Argon2id-specific defaults (more secure than Argon2i):
//   - Time: 2 iterations (higher than Argon2i's default of 1)
//   - All other parameters inherit from ArgonHasher with same validation
//
// The key difference is only in the key derivation function - all other functionality
// including parameter validation, salt generation, and configuration handling is
// inherited from ArgonHasher.
//
// Parameters:
//   - config: Configuration map with optional parameter overrides
//
// Returns:
//   - *Argon2IdHasher: A new Argon2id hasher instance with validated configuration
func NewArgon2IdHasher(config map[string]interface{}) *Argon2IdHasher {
	// Create base Argon hasher with Argon2id-optimized defaults
	// Override time parameter to 2 for better security with hybrid approach
	if config == nil {
		config = make(map[string]interface{})
	}

	// Set Argon2id-specific default (higher time cost than Argon2i)
	if _, exists := config["time"]; !exists {
		config["time"] = 2 // 2 iterations vs 1 for Argon2i
	}

	// Create base ArgonHasher with modified config
	argonHasher := NewArgonHasher(config)

	// Wrap in Argon2IdHasher to override key derivation method
	argon2IdHasher := &Argon2IdHasher{
		ArgonHasher: argonHasher,
	}

	// Update proxy to point to the concrete implementation
	argonHasher.SetProxySelf(argon2IdHasher)

	return argon2IdHasher
}

// Algorithm returns "argon2id" for proper hash format generation.
// Overrides ArgonHasher's "argon2i" identifier.
func (h *Argon2IdHasher) Algorithm() string {
	return "argon2id"
}

// deriveKey overrides ArgonHasher's key derivation to use Argon2id instead of Argon2i.
// This is the only method that needs to be different between the two implementations.
// All other functionality (validation, salt generation, formatting, etc.) is inherited.
func (h *Argon2IdHasher) deriveKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
	// Use Argon2id key derivation instead of Argon2i
	// This provides hybrid security benefits combining Argon2i and Argon2d
	return argon2.IDKey(password, salt, time, memory, threads, keyLen)
}

// Compile-time interface compliance check
var _ interfaces.HasherInterface = (*Argon2IdHasher)(nil)
