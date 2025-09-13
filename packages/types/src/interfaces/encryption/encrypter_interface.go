package interfaces

import (
	types "govel/types"
)

// EncrypterInterface defines the contract for encryption operations.
// This interface provides comprehensive encryption functionality including
// encryption, decryption, hash operations, payload validation, and configuration.
//
// Key features:
//   - Symmetric encryption with serialization support
//   - MAC-based authentication for data integrity
//   - Payload validation and parsing capabilities
//   - Key and cipher configuration management
//   - Support for multiple cipher algorithms (AES-GCM, AES-CBC, AES-CTR)
type EncrypterInterface interface {
	// Info gets comprehensive information about the given encrypted payload.
	// Extracts cipher metadata and parameters from encrypted payload format.
	//
	// Parameters:
	//   - payload: The encrypted payload string to analyze
	//
	// Returns:
	//   - CipherInfo: Structured information about cipher, key length, IV length, etc.
	Info(payload string) types.CipherInfo

	// Encrypt encrypts the given value with optional serialization.
	// Creates secure encrypted payload with MAC authentication.
	//
	// Parameters:
	//   - value: The value to encrypt (any serializable type)
	//   - serialize: Whether to JSON serialize the value before encryption
	//
	// Returns:
	//   - string: Base64-encoded encrypted payload with MAC
	//   - error: Encryption error if operation fails
	Encrypt(value interface{}, serialize bool) (string, error)

	// EncryptString encrypts a string value without serialization.
	// Optimized method for direct string encryption.
	//
	// Parameters:
	//   - value: The string to encrypt
	//
	// Returns:
	//   - string: Base64-encoded encrypted payload with MAC
	//   - error: Encryption error if operation fails
	EncryptString(value string) (string, error)

	// Decrypt decrypts the given payload with optional unserialization.
	// Validates MAC before decryption to prevent tampering.
	//
	// Parameters:
	//   - payload: The encrypted payload string to decrypt
	//   - unserialize: Whether to JSON deserialize after decryption
	//
	// Returns:
	//   - interface{}: Decrypted value (type depends on unserialize flag)
	//   - error: Decryption or validation error if operation fails
	Decrypt(payload string, unserialize bool) (interface{}, error)

	// DecryptString decrypts a payload back to a string without unserialization.
	// Optimized method for direct string decryption.
	//
	// Parameters:
	//   - payload: The encrypted payload string to decrypt
	//
	// Returns:
	//   - string: Decrypted string value
	//   - error: Decryption or validation error if operation fails
	DecryptString(payload string) (string, error)

	// Hash computes HMAC-SHA256 for the given parameters.
	// Used internally for MAC computation and validation.
	//
	// Parameters:
	//   - iv: Initialization vector bytes
	//   - value: The encrypted value bytes
	//   - key: The encryption key bytes
	//
	// Returns:
	//   - []byte: Computed HMAC hash
	//   - error: Hash computation error if operation fails
	Hash(iv, value, key []byte) ([]byte, error)

	// GetJsonPayload parses an encrypted payload string into its components.
	// Extracts cipher metadata, IV, MAC, and encrypted value from JSON structure.
	//
	// Parameters:
	//   - payload: The encrypted payload string to parse
	//
	// Returns:
	//   - map[string]interface{}: Component map with keys: cipher, iv, value, mac
	//   - error: Parsing error if payload format is invalid
	GetJsonPayload(payload string) (map[string]interface{}, error)

	// ValidPayload validates the structure and format of an encrypted payload.
	// Checks JSON structure, required fields, and Base64 encoding validity.
	//
	// Parameters:
	//   - payload: The encrypted payload string to validate
	//
	// Returns:
	//   - bool: true if payload is structurally valid, false otherwise
	ValidPayload(payload string) bool

	// ValidMac validates the MAC of a parsed payload using the default encryption key.
	// Performs constant-time comparison to prevent timing attacks.
	//
	// Parameters:
	//   - payload: Parsed payload components as map[string]interface{}
	//
	// Returns:
	//   - bool: true if MAC is valid and payload is authentic, false otherwise
	ValidMac(payload map[string]interface{}) bool

	// ValidMacForKey validates the MAC using a specific encryption key.
	// Allows MAC validation with custom key material for multi-key scenarios.
	//
	// Parameters:
	//   - payload: Parsed payload components as map[string]interface{}
	//   - key: The encryption key bytes to use for MAC validation
	//
	// Returns:
	//   - bool: true if MAC is valid with the specified key, false otherwise
	ValidMacForKey(payload map[string]interface{}, key []byte) bool

	// EnsureTagIsValid validates GCM authentication tag format and length.
	// Specific to AEAD ciphers that use authentication tags instead of separate MAC.
	//
	// Parameters:
	//   - tag: The authentication tag string to validate
	//
	// Returns:
	//   - error: nil if tag is valid, error describing validation failure otherwise
	EnsureTagIsValid(tag string) error

	// ShouldValidateMac determines if MAC validation should be performed based on cipher type.
	// AEAD ciphers (like GCM) don't require separate MAC as they have built-in authentication.
	//
	// Returns:
	//   - bool: true if MAC validation is required, false for AEAD ciphers
	ShouldValidateMac() bool

	// GetKey returns the current encryption key bytes.
	// Provides access to configured key material for advanced operations.
	//
	// Returns:
	//   - []byte: Current encryption key bytes, nil if not configured
	GetKey() []byte

	// GetCipher returns the current cipher algorithm identifier.
	// Returns cipher name like "aes-256-gcm", "aes-128-cbc", etc.
	//
	// Returns:
	//   - string: Current cipher identifier, empty string if not configured
	GetCipher() string

	// SetKey sets the encryption key with method chaining support.
	// Validates key length compatibility with current cipher algorithm.
	//
	// Parameters:
	//   - key: The encryption key bytes to configure
	//
	// Returns:
	//   - EncrypterInterface: Returns self for method chaining
	SetKey(key []byte) EncrypterInterface

	// SetCipher sets the cipher algorithm with method chaining support.
	// Updates cipher while preserving other configuration settings.
	//
	// Parameters:
	//   - cipher: The cipher algorithm identifier (e.g., "aes-256-gcm")
	//
	// Returns:
	//   - EncrypterInterface: Returns self for method chaining
	SetCipher(cipher string) EncrypterInterface

	// GenerateKey generates a cryptographically secure key for the specified cipher.
	// Uses system random number generator for high-entropy key material.
	//
	// Parameters:
	//   - cipher: The cipher algorithm to generate a key for
	//
	// Returns:
	//   - []byte: Generated key bytes with cipher-appropriate length
	//   - error: Key generation error if cipher is unsupported or RNG fails
	GenerateKey(cipher string) ([]byte, error)

	// VerifyConfiguration validates that the encrypter's current configuration
	// meets security requirements and is suitable for production use.
	// Checks key strength, cipher compatibility, and security parameters.
	//
	// Returns:
	//   - bool: true if configuration is valid and secure, false otherwise
	VerifyConfiguration() bool
}
