package enums

import (
	"errors"
	"strings"
)

// Algorithm represents the hashing algorithms.
type Algorithm string

const (
	// AlgorithmBcrypt represents bcrypt algorithm
	AlgorithmBcrypt Algorithm = "bcrypt"

	// AlgorithmArgon2i represents Argon2i algorithm
	AlgorithmArgon2i Algorithm = "argon2i"

	// AlgorithmArgon2id represents Argon2id algorithm
	AlgorithmArgon2id Algorithm = "argon2id"

	// AlgorithmSHA256 represents SHA-256 algorithm
	AlgorithmSHA256 Algorithm = "sha256"

	// AlgorithmSHA512 represents SHA-512 algorithm
	AlgorithmSHA512 Algorithm = "sha512"
)

// String returns the string representation of the Algorithm.
func (a Algorithm) String() string {
	return string(a)
}

// ValidateAlgorithm validates if the given algorithm is supported.
// Returns an error if the algorithm is not recognized or supported.
//
// Parameters:
//   - algorithm: The hashing algorithm to validate
//
// Returns:
//   - error: nil if valid, error describing the issue if invalid
func ValidateAlgorithm(algorithm Algorithm) error {
	// Convert algorithm to string for validation
	normalizedAlgorithm := strings.ToLower(algorithm.String())

	switch normalizedAlgorithm {
	case "bcrypt", "argon2i", "argon2id", "sha256", "sha512":
		return nil
	default:
		return errors.New("unsupported hashing algorithm")
	}
}

// GetDefaultAlgorithm returns the default hashing algorithm for the system.
// This provides a centralized way to manage the default hashing algorithm.
//
// Returns:
//   - string: The default hashing algorithm identifier
func GetDefaultAlgorithm() string {
	return string(AlgorithmBcrypt) // Default to bcrypt for compatibility
}
