package encryption

import (
	"fmt"
	"govel/types/src/types/encryption"
	encrypterInterfaces "govel/types/src/interfaces/encryption"
)

// Info gets comprehensive information about the given encrypted payload.
// Delegates to the default driver to extract cipher metadata and parameters.
//
// Method features:
//   - Automatic cipher detection from payload format
//   - Parameter extraction (cipher, key length, IV length)
//   - Graceful error handling with empty CipherInfo fallback
//   - Support for all configured encryption algorithms
//
// Returns CipherInfo with cipher details or empty info if parsing fails.
func (e *EncryptionManager) Info(payload string) types.CipherInfo {
	// Get default driver for encryption analysis
	// Uses configured default cipher for consistent behavior
	driver, err := e.Driver()
	if err != nil {
		// Return empty CipherInfo on driver resolution failure
		// Ensures consistent behavior when no drivers are available
		return types.CipherInfo{
			Cipher:      "",
			KeyLength:   0,
			IVLength:    0,
			Mode:        "",
			IsAEAD:      false,
			RequiresMAC: false,
			Options:     make(map[string]interface{}),
		}
	}

	// Type-safe casting to EncrypterInterface
	// Ensures driver supports encryption information extraction
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return empty CipherInfo if driver doesn't implement EncrypterInterface
		// Should rarely happen but provides safety net
		return types.CipherInfo{
			Cipher:      "",
			KeyLength:   0,
			IVLength:    0,
			Mode:        "",
			IsAEAD:      false,
			RequiresMAC: false,
			Options:     make(map[string]interface{}),
		}
	}

	// Delegate to actual encrypter implementation for encryption analysis
	return encrypter.Info(payload)
}

// Encrypt encrypts the given value with optional serialization.
// Delegates to the default driver to create secure encrypted payload with MAC.
//
// Method features:
//   - Secure encryption using default cipher
//   - Optional JSON serialization support
//   - Automatic IV generation and MAC computation
//   - Error propagation for invalid inputs or configurations
//
// Returns encrypted payload string or error if encryption fails.
func (e *EncryptionManager) Encrypt(value interface{}, serialize bool) (string, error) {
	// Get default driver for encryption
	// Uses configured default cipher for consistent encryption
	driver, err := e.Driver()
	if err != nil {
		// Propagate driver resolution errors
		return "", err
	}

	// Type-safe casting to EncrypterInterface
	// Ensures driver supports encryption operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return descriptive error for interface compliance issues
		return "", fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to actual encrypter implementation with parameters
	// Handles cipher-specific options and secure payload generation
	return encrypter.Encrypt(value, serialize)
}

// EncryptString encrypts a string value without serialization.
// Delegates to the default driver for direct string encryption.
//
// Encryption features:
//   - Direct string encryption without serialization overhead
//   - Secure IV generation and MAC computation
//   - Optimized for string data types
//   - Safe handling of empty strings
//
// Returns encrypted payload string or error if encryption fails.
func (e *EncryptionManager) EncryptString(value string) (string, error) {
	// Get default driver for encryption
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return "", err
	}

	// Ensure driver implements EncrypterInterface
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support encryption
		return "", fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for secure string encryption
	return encrypter.EncryptString(value)
}

// Decrypt decrypts the given payload with optional unserialization.
// Delegates to the default driver for secure payload decryption and validation.
//
// Decryption features:
//   - MAC validation before decryption to prevent tampering
//   - Optional JSON deserialization support
//   - Support for all configured cipher algorithms
//   - Safe handling of malformed or invalid payloads
//
// Returns decrypted value or error if decryption/validation fails.
func (e *EncryptionManager) Decrypt(payload string, unserialize bool) (interface{}, error) {
	// Get default driver for decryption
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return nil, err
	}

	// Ensure driver implements EncrypterInterface
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support decryption
		return nil, fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for secure decryption
	return encrypter.Decrypt(payload, unserialize)
}

// DecryptString decrypts a payload back to a string without unserialization.
// Delegates to the default driver for direct string decryption.
//
// Decryption features:
//   - Direct string decryption without deserialization
//   - MAC validation to prevent payload tampering
//   - Optimized for string data types
//   - Consistent error handling
//
// Returns decrypted string or error if decryption fails.
func (e *EncryptionManager) DecryptString(payload string) (string, error) {
	// Get default driver for decryption
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return "", err
	}

	// Ensure driver implements EncrypterInterface
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support decryption
		return "", fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for string decryption
	return encrypter.DecryptString(payload)
}

// Hash computes HMAC-SHA256 for the given parameters.
// Delegates to the default driver for MAC computation.
//
// Hash features:
//   - HMAC-SHA256 computation for integrity protection
//   - Proper key and data handling
//   - Cryptographically secure hash generation
//   - Support for all cipher configurations
//
// Returns computed HMAC hash or error if computation fails.
func (e *EncryptionManager) Hash(iv, value, key []byte) ([]byte, error) {
	// Get default driver for hash computation
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return nil, err
	}

	// Ensure driver implements EncrypterInterface
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support hashing
		return nil, fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for hash computation
	return encrypter.Hash(iv, value, key)
}

// GetJsonPayload parses an encrypted payload string into its components.
// Delegates to the default driver for secure payload parsing and structure extraction.
//
// Parsing features:
//   - JSON payload deserialization into component map
//   - Extraction of cipher metadata, IV, MAC, and encrypted value
//   - Graceful handling of malformed or invalid payload formats
//   - Support for all cipher-specific payload structures
//
// Returns component map with keys: cipher, iv, value, mac or error if parsing fails.
func (e *EncryptionManager) GetJsonPayload(payload string) (map[string]interface{}, error) {
	// Get default driver for payload parsing
	// Uses configured cipher for format-specific parsing logic
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return nil, err
	}

	// Ensure driver implements EncrypterInterface for payload operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support payload parsing
		return nil, fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for secure payload component extraction
	return encrypter.GetJsonPayload(payload)
}

// ValidPayload validates the structure and format of an encrypted payload.
// Delegates to the default driver for comprehensive payload validation.
//
// Validation features:
//   - JSON structure verification for proper payload format
//   - Required field presence checking (cipher, iv, value, mac)
//   - Base64 encoding validation for binary components
//   - Cipher-specific format requirements verification
//
// Returns true if payload is structurally valid, false otherwise.
func (e *EncryptionManager) ValidPayload(payload string) bool {
	// Get default driver for payload validation
	// Uses cipher-specific validation rules and format requirements
	driver, err := e.Driver()
	if err != nil {
		// Return false on driver errors - invalid configuration
		return false
	}

	// Ensure driver implements EncrypterInterface for validation operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return false if driver doesn't support payload validation
		return false
	}

	// Delegate to encrypter for comprehensive payload structure validation
	return encrypter.ValidPayload(payload)
}

// ValidMac validates the MAC of a parsed payload using the default encryption key.
// Delegates to the default driver for cryptographic MAC verification.
//
// MAC validation features:
//   - HMAC-SHA256 computation and comparison
//   - Constant-time comparison to prevent timing attacks
//   - Automatic key derivation from configured encryption key
//   - Protection against payload tampering and forgery
//
// Returns true if MAC is valid and payload is authentic, false otherwise.
func (e *EncryptionManager) ValidMac(payload map[string]interface{}) bool {
	// Get default driver for MAC validation
	// Uses configured key and cipher for secure MAC verification
	driver, err := e.Driver()
	if err != nil {
		// Return false on driver errors - cannot validate securely
		return false
	}

	// Ensure driver implements EncrypterInterface for MAC operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return false if driver doesn't support MAC validation
		return false
	}

	// Delegate to encrypter for cryptographic MAC verification
	return encrypter.ValidMac(payload)
}

// ValidMacForKey validates the MAC using a specific encryption key.
// Delegates to the default driver for key-specific cryptographic MAC verification.
//
// Key-specific MAC validation features:
//   - HMAC-SHA256 computation with custom key material
//   - Constant-time comparison to prevent timing attacks
//   - Independent validation without affecting current configuration
//   - Support for multi-key scenarios and key rotation
//
// Returns true if MAC is valid with the specified key, false otherwise.
func (e *EncryptionManager) ValidMacForKey(payload map[string]interface{}, key []byte) bool {
	// Get default driver for key-specific MAC validation
	// Uses provided key instead of configured default key
	driver, err := e.Driver()
	if err != nil {
		// Return false on driver errors - cannot validate securely
		return false
	}

	// Ensure driver implements EncrypterInterface for MAC operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return false if driver doesn't support MAC validation
		return false
	}

	// Delegate to encrypter for key-specific MAC verification
	return encrypter.ValidMacForKey(payload, key)
}

// EnsureTagIsValid validates GCM authentication tag format and length.
// Delegates to the default driver for AEAD-specific tag validation.
//
// GCM tag validation features:
//   - Authentication tag format verification (Base64 encoding)
//   - Tag length validation for GCM cipher requirements
//   - AEAD-specific security parameter checking
//   - Protection against malformed or truncated tags
//
// Returns nil if tag is valid, error if tag format or length is invalid.
func (e *EncryptionManager) EnsureTagIsValid(tag string) error {
	// Get default driver for AEAD tag validation
	// Uses cipher-specific requirements for tag format and length
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return err
	}

	// Ensure driver implements EncrypterInterface for tag validation
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support tag validation
		return fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for AEAD tag format and length validation
	return encrypter.EnsureTagIsValid(tag)
}

// ShouldValidateMac determines if MAC validation should be performed based on cipher type.
// Delegates to the default driver for cipher-specific MAC requirements and AEAD detection.
//
// MAC requirement determination features:
//   - AEAD cipher detection (GCM modes don't require separate MAC)
//   - Traditional cipher identification (CBC, CTR modes require MAC)
//   - Conservative security approach with default validation
//   - Cipher-specific authentication mechanism selection
//
// Returns true if MAC validation is required, false for AEAD ciphers with built-in authentication.
func (e *EncryptionManager) ShouldValidateMac() bool {
	// Get default driver for cipher-specific MAC requirements
	// Different cipher modes have different authentication requirements
	driver, err := e.Driver()
	if err != nil {
		// Conservative approach - validate by default on errors
		// Ensures security when driver configuration is uncertain
		return true
	}

	// Ensure driver implements EncrypterInterface for MAC requirements
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Conservative approach - validate by default if interface unavailable
		return true
	}

	// Delegate to encrypter for cipher-specific MAC requirement determination
	return encrypter.ShouldValidateMac()
}

// GetKey returns the current encryption key from the default driver.
// Delegates to the default driver for secure key material access.
//
// Key retrieval features:
//   - Access to configured encryption key material
//   - Safe handling of key material without exposure
//   - Support for key rotation and multi-key scenarios
//   - Consistent key format across all operations
//
// Returns current encryption key bytes or nil if key is not configured.
func (e *EncryptionManager) GetKey() []byte {
	// Get default driver for key access
	// Retrieves currently configured encryption key material
	driver, err := e.Driver()
	if err != nil {
		// Return nil on driver errors - no valid key available
		return nil
	}

	// Ensure driver implements EncrypterInterface for key operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return nil if driver doesn't support key access
		return nil
	}

	// Delegate to encrypter for secure key material retrieval
	return encrypter.GetKey()
}

// GetCipher returns the current cipher algorithm identifier from the default driver.
// Delegates to the default driver for cipher configuration access.
//
// Cipher identification features:
//   - Current cipher algorithm name retrieval (e.g., "aes-256-gcm")
//   - Support for all configured cipher types
//   - Consistent cipher naming across encryption operations
//   - Integration with cipher-specific parameter selection
//
// Returns cipher identifier string or empty string if not configured.
func (e *EncryptionManager) GetCipher() string {
	// Get default driver for cipher identification
	// Retrieves currently configured cipher algorithm
	driver, err := e.Driver()
	if err != nil {
		// Return empty string on driver errors - no valid cipher
		return ""
	}

	// Ensure driver implements EncrypterInterface for cipher operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return empty string if driver doesn't support cipher access
		return ""
	}

	// Delegate to encrypter for cipher algorithm identification
	return encrypter.GetCipher()
}

// SetKey sets the encryption key for the default encrypter with method chaining support.
// Delegates to the default driver for secure key configuration and returns manager for fluent interface.
//
// Key configuration features:
//   - Secure key material storage and configuration
//   - Key length validation for cipher compatibility
//   - Support for key rotation and dynamic key updates
//   - Method chaining for fluent configuration API
//
// Returns EncryptionManager instance for method chaining, preserves existing config on errors.
func (e *EncryptionManager) SetKey(key []byte) encrypterInterfaces.EncrypterInterface {
	// Get default driver for key configuration
	// Updates encryption key while preserving other settings
	driver, err := e.Driver()
	if err != nil {
		// Return manager unchanged on driver errors
		// Preserves existing configuration for graceful degradation
		return e
	}

	// Ensure driver implements EncrypterInterface for key operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return manager unchanged if driver doesn't support key setting
		return e
	}

	// Delegate to encrypter for secure key configuration
	// Key validation and storage handled by specific driver implementation
	encrypter.SetKey(key)
	// Return manager for method chaining support
	return e
}

// SetCipher sets the cipher algorithm for the default encrypter with method chaining support.
// Delegates to the default driver for cipher configuration and returns manager for fluent interface.
//
// Cipher configuration features:
//   - Dynamic cipher algorithm selection (AES-256-GCM, AES-128-CBC, etc.)
//   - Cipher compatibility validation with current key length
//   - Automatic parameter adjustment for cipher-specific requirements
//   - Method chaining for fluent configuration API
//
// Returns EncryptionManager instance for method chaining, preserves existing config on errors.
func (e *EncryptionManager) SetCipher(cipher string) encrypterInterfaces.EncrypterInterface {
	// Get default driver for cipher configuration
	// Updates cipher algorithm while preserving other settings
	driver, err := e.Driver()
	if err != nil {
		// Return manager unchanged on driver errors
		// Preserves existing configuration for graceful degradation
		return e
	}

	// Ensure driver implements EncrypterInterface for cipher operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return manager unchanged if driver doesn't support cipher setting
		return e
	}

	// Delegate to encrypter for cipher algorithm configuration
	// Cipher validation and parameter setup handled by specific driver
	encrypter.SetCipher(cipher)
	// Return manager for method chaining support
	return e
}

// GenerateKey generates a cryptographically secure key appropriate for the specified cipher.
// Delegates to the default driver for cipher-specific key generation with proper entropy.
//
// Key generation features:
//   - Cryptographically secure random key generation
//   - Cipher-specific key length determination (128, 192, or 256 bits for AES)
//   - High-entropy key material from system random number generator
//   - Compliance with cryptographic standards and best practices
//
// Returns generated key bytes or error if key generation fails or cipher is unsupported.
func (e *EncryptionManager) GenerateKey(cipher string) ([]byte, error) {
	// Get default driver for cipher-specific key generation
	// Uses secure random number generation for cryptographic strength
	driver, err := e.Driver()
	if err != nil {
		// Return error on driver resolution failure
		return nil, err
	}

	// Ensure driver implements EncrypterInterface for key generation
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return error if driver doesn't support key generation
		return nil, fmt.Errorf("driver does not implement EncrypterInterface")
	}

	// Delegate to encrypter for secure key generation with cipher-appropriate length
	return encrypter.GenerateKey(cipher)
}

// VerifyConfiguration validates that the encrypter's current configuration
// meets security requirements and is suitable for production use.
// Delegates to the default driver for configuration validation.
func (e *EncryptionManager) VerifyConfiguration() bool {
	// Get default driver for configuration verification
	driver, err := e.Driver()
	if err != nil {
		// Return false on driver errors indicating configuration issues
		return false
	}

	// Ensure driver implements EncrypterInterface
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Return false if driver doesn't support configuration verification
		return false
	}

	// Delegate to encrypter for configuration validation
	return encrypter.VerifyConfiguration()
}

// Compile-time interface compliance checks
// These ensure EncryptionManager properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ encrypterInterfaces.EncrypterInterface = (*EncryptionManager)(nil) // Direct encryption operations
