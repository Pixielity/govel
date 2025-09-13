package hashing

import (
	"fmt"
	"govel/types/src/types/hashing"
	hashingInterfaces "govel/types/src/interfaces/hashing"
)

// HasherInterface implementation - delegates to default driver
// These methods provide a unified interface for hash operations by delegating
// to the current default driver. This allows the manager to act as both a
// factory (via Hasher method) and a direct hasher (via these delegate methods).

// Info gets comprehensive information about the given hashed value.
// Delegates to the default driver to extract algorithm metadata and parameters.
//
// Method features:
//   - Automatic algorithm detection from hash format
//   - Parameter extraction (cost, memory, time, threads)
//   - Graceful error handling with empty HashInfo fallback
//   - Support for all configured hash algorithms
//
// Returns HashInfo with algorithm details or empty info if parsing fails.
func (h *HashManager) Info(hashedValue string) types.HashInfo {
	// Get default driver for hash analysis
	// Uses configured default algorithm for consistent behavior
	driver, err := h.Driver()
	if err != nil {
		// Return empty HashInfo on driver resolution failure
		// Ensures consistent behavior when no drivers are available
		return types.HashInfo{
			Algo:     "",
			AlgoName: "",
			Options:  make(map[string]interface{}),
		}
	}

	// Type-safe casting to HasherInterface
	// Ensures driver supports hash information extraction
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return empty HashInfo if driver doesn't implement HasherInterface
		// Should rarely happen but provides safety net
		return types.HashInfo{
			Algo:     "",
			AlgoName: "",
			Options:  make(map[string]interface{}),
		}
	}

	// Delegate to actual hasher implementation for hash analysis
	return hasher.Info(hashedValue)
}

// Make hashes the given value with optional parameters.
// Delegates to the default driver to create secure hash with specified options.
//
// Method features:
//   - Secure hash generation using default algorithm
//   - Optional parameter support for algorithm-specific settings
//   - Automatic salt generation and secure random handling
//   - Error propagation for invalid inputs or configurations
//
// Returns secure hash string or error if hashing fails.
func (h *HashManager) Make(value string, options map[string]interface{}) (string, error) {
	// Get default driver for hash creation
	// Uses configured default algorithm for consistent hashing
	driver, err := h.Driver()
	if err != nil {
		// Propagate driver resolution errors
		return "", err
	}

	// Type-safe casting to HasherInterface
	// Ensures driver supports hash generation
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return descriptive error for interface compliance issues
		return "", fmt.Errorf("driver does not implement HasherInterface")
	}

	// Delegate to actual hasher implementation with parameters
	// Handles algorithm-specific options and secure hash generation
	return hasher.Make(value, options)
}

// Check verifies the given plain value against a hash.
// Delegates to the default driver for secure password verification.
//
// Verification features:
//   - Constant-time comparison to prevent timing attacks
//   - Support for all configured hash algorithms
//   - Optional parameter support for verification settings
//   - Safe handling of malformed or invalid hashes
//
// Returns true if value matches hash, false otherwise or on errors.
func (h *HashManager) Check(value, hashedValue string, options map[string]interface{}) bool {
	// Get default driver for hash verification
	driver, err := h.Driver()
	if err != nil {
		// Return false on driver resolution failure for security
		return false
	}

	// Ensure driver implements HasherInterface
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return false if driver doesn't support verification
		return false
	}

	// Delegate to hasher for secure verification
	return hasher.Check(value, hashedValue, options)
}

// NeedsRehash checks if the given hash needs rehashing with updated options.
// Delegates to the default driver to determine if hash parameters are outdated.
//
// Rehash detection features:
//   - Algorithm parameter comparison (cost, memory, time)
//   - Security standard compliance checking
//   - Migration support between algorithms
//   - Conservative approach returning true on errors
//
// Returns true if rehashing recommended, false if current hash is adequate.
func (h *HashManager) NeedsRehash(hashedValue string, options map[string]interface{}) bool {
	// Get default driver for rehash analysis
	driver, err := h.Driver()
	if err != nil {
		// Return true on errors to trigger rehashing for security
		return true
	}

	// Ensure driver implements HasherInterface
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return true if driver doesn't support rehash detection
		return true
	}

	// Delegate to hasher for rehash analysis
	return hasher.NeedsRehash(hashedValue, options)
}

// IsHashed determines if a given string is already a hash.
// Delegates to the default driver for hash format detection.
//
// Detection features:
//   - Hash format pattern recognition
//   - Algorithm-specific prefix detection
//   - Length and character validation
//   - Safe handling of arbitrary input strings
//
// Returns true if string appears to be a hash, false otherwise.
func (h *HashManager) IsHashed(value string) bool {
	// Get default driver for hash format detection
	driver, err := h.Driver()
	if err != nil {
		// Return false on driver errors for safety
		return false
	}

	// Ensure driver implements HasherInterface
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return false if driver doesn't support detection
		return false
	}

	// Delegate to hasher for format detection
	return hasher.IsHashed(value)
}

// VerifyConfiguration verifies that the current configuration is valid.
// Delegates to the default driver for configuration validation.
//
// Validation features:
//   - Parameter range checking (cost, memory, threads)
//   - System resource availability verification
//   - Algorithm-specific constraint validation
//   - Production readiness assessment
//
// Returns true if configuration is valid and production-ready.
func (h *HashManager) VerifyConfiguration(value string) bool {
	// Get default driver for configuration verification
	driver, err := h.Driver()
	if err != nil {
		// Return false on driver errors indicating configuration issues
		return false
	}

	// Ensure driver implements HasherInterface
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Return false if driver doesn't support configuration verification
		return false
	}

	// Delegate to hasher for configuration validation
	return hasher.VerifyConfiguration(value)
}

// Compile-time interface compliance checks
// These ensure HashManager properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ hashingInterfaces.HasherInterface = (*HashManager)(nil) // Direct hash operations
