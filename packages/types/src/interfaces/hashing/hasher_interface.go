package interfaces

import (
	types "govel/types"
)

// HasherInterface defines the contract for hash operations.
// This interface provides comprehensive password hashing functionality including
// hash creation, verification, rehash detection, and configuration validation.
//
// Key features:
//   - Secure password hashing with multiple algorithms (Argon2, bcrypt)
//   - Constant-time verification to prevent timing attacks
//   - Rehash detection for security parameter updates
//   - Hash format detection and validation
//   - Production-ready configuration verification
type HasherInterface interface {
	// Info gets comprehensive information about the given hashed value.
	// Analyzes hash format to extract algorithm metadata and parameters.
	//
	// Parameters:
	//   - hashedValue: The hash string to analyze
	//
	// Returns:
	//   - HashInfo: Structured information about algorithm, cost parameters, etc.
	Info(hashedValue string) types.HashInfo

	// Make hashes the given value with optional parameters.
	// Creates secure hash using the configured default algorithm.
	//
	// Parameters:
	//   - value: The plaintext value to hash (typically password)
	//   - options: Algorithm-specific options (cost, memory, time, threads)
	//
	// Returns:
	//   - string: Secure hash string with embedded algorithm and parameters
	//   - error: Hashing error if operation fails or options are invalid
	Make(value string, options map[string]interface{}) (string, error)

	// Check verifies the given plain value against a hash.
	// Performs constant-time comparison to prevent timing attacks.
	//
	// Parameters:
	//   - value: The plaintext value to verify
	//   - hashedValue: The hash string to verify against
	//   - options: Optional verification parameters
	//
	// Returns:
	//   - bool: true if value matches hash, false otherwise or on errors
	Check(value, hashedValue string, options map[string]interface{}) bool

	// NeedsRehash checks if the given hash needs rehashing with updated options.
	// Compares current hash parameters with recommended security standards.
	//
	// Parameters:
	//   - hashedValue: The hash string to analyze
	//   - options: Desired hash parameters for comparison
	//
	// Returns:
	//   - bool: true if rehashing is recommended, false if current hash is adequate
	NeedsRehash(hashedValue string, options map[string]interface{}) bool

	// IsHashed determines if a given string is already a hash.
	// Analyzes format patterns to detect hash vs plaintext strings.
	//
	// Parameters:
	//   - value: The string to analyze
	//
	// Returns:
	//   - bool: true if string appears to be a hash, false otherwise
	IsHashed(value string) bool

	// VerifyConfiguration verifies that the current configuration is valid.
	// Checks parameter ranges, system resources, and production readiness.
	//
	// Parameters:
	//   - value: Optional test value for configuration validation
	//
	// Returns:
	//   - bool: true if configuration is valid and production-ready, false otherwise
	VerifyConfiguration(value string) bool
}
