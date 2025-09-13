package interfaces

import enums "govel/types/enums/hashing"

// FactoryInterface defines the contract for hashing factory functionality.
// This interface provides hasher instance creation and management capabilities.
//
// Key features:
//   - Multiple hashing algorithm support
//   - Laravel-style driver pattern implementation
//   - Default algorithm management
//   - Type-safe hasher instance creation
type FactoryInterface interface {
	// Hasher gets a hasher instance by algorithm name (optional).
	// If no name is provided, uses the default driver.
	//
	// Parameters:
	//   - name: Optional algorithm name ("bcrypt", "argon2i", "argon2id")
	//
	// Returns:
	//   - HasherInterface: A hasher instance for the specified algorithm
	Hasher(name ...enums.Algorithm) HasherInterface
}
