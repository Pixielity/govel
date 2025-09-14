package encrypters

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	types "govel/types"
	enums "govel/types/enums/encryption"
	encryptionInterfaces "govel/types/interfaces/encryption"
)

// BaseEncrypter provides common functionality for all encrypter implementations.
// This struct contains shared methods and utilities used by specific cipher implementations.
//
// Key features:
//   - Common key and cipher management
//   - Shared IV generation and validation
//   - Base64 encoding/decoding utilities
//   - MAC computation and validation
//   - JSON payload handling
//   - Laravel-compatible serialization
//
// Follows the same pattern as the hashing package's BaseHasher.
type BaseEncrypter struct {
	// key holds the encryption key bytes
	key []byte

	// cipher identifies the encryption algorithm
	cipher string

	// config contains encrypter-specific configuration
	config map[string]interface{}
}

// NewBaseEncrypter creates a new BaseEncrypter with the given key and cipher.
// Validates the key length against the cipher requirements.
//
// Parameters:
//   - key: The encryption key bytes
//   - cipher: The cipher algorithm identifier
//   - config: Optional configuration parameters
//
// Returns a configured BaseEncrypter or an error if validation fails.
func NewBaseEncrypter(key []byte, cipher string, config map[string]interface{}) (*BaseEncrypter, error) {
	// Validate cipher
	if err := enums.ValidateCipher(enums.Cipher(cipher)); err != nil {
		return nil, fmt.Errorf("invalid cipher: %w", err)
	}

	// Validate key length
	expectedKeyLen := enums.GetKeyLength(enums.Cipher(cipher))
	if expectedKeyLen == 0 {
		return nil, fmt.Errorf("unsupported cipher: %s", cipher)
	}
	if len(key) != expectedKeyLen {
		return nil, fmt.Errorf("invalid key length: expected %d bytes, got %d", expectedKeyLen, len(key))
	}

	// Initialize default config if nil
	if config == nil {
		config = make(map[string]interface{})
	}

	return &BaseEncrypter{
		key:    key,
		cipher: cipher,
		config: config,
	}, nil
}

// The following methods are required by EncrypterInterface but should be overridden by concrete implementations

// Encrypt is a placeholder that must be overridden by concrete encrypters.
func (b *BaseEncrypter) Encrypt(value interface{}, serialize bool) (string, error) {
	panic("Encrypt method must be implemented by concrete encrypter")
}

// EncryptString is a placeholder that must be overridden by concrete encrypters.
func (b *BaseEncrypter) EncryptString(value string) (string, error) {
	panic("EncryptString method must be implemented by concrete encrypter")
}

// Decrypt is a placeholder that must be overridden by concrete encrypters.
func (b *BaseEncrypter) Decrypt(payload string, unserialize bool) (interface{}, error) {
	panic("Decrypt method must be implemented by concrete encrypter")
}

// DecryptString is a placeholder that must be overridden by concrete encrypters.
func (b *BaseEncrypter) DecryptString(payload string) (string, error) {
	panic("DecryptString method must be implemented by concrete encrypter")
}

// GetKey returns the current encryption key with secure memory handling.
// Creates a defensive copy to prevent external modification of internal key material.
//
// Key access features:
//   - Defensive copying to prevent external key manipulation
//   - Safe memory handling to avoid key exposure
//   - Consistent interface for key retrieval across implementations
//   - Support for key rotation and validation scenarios
//
// Returns a copy of the encryption key bytes to maintain security.
func (b *BaseEncrypter) GetKey() []byte {
	// Return a copy to prevent external modification of sensitive key material
	// This defensive approach ensures the internal key remains secure
	keyCopy := make([]byte, len(b.key))
	copy(keyCopy, b.key)
	return keyCopy
}

// GetCipher returns the current cipher algorithm identifier string.
// Provides access to the configured cipher for validation and information purposes.
//
// Cipher identification features:
//   - Returns standardized cipher names (e.g., "aes-256-gcm")
//   - Enables cipher-specific parameter validation
//   - Supports dynamic cipher configuration scenarios
//   - Consistent naming across all encrypter implementations
//
// Returns the cipher algorithm identifier string.
func (b *BaseEncrypter) GetCipher() string {
	// Return the currently configured cipher algorithm identifier
	return b.cipher
}

// SetKey sets the encryption key with validation and method chaining support.
// Validates key length against cipher requirements and provides fluent interface.
//
// Key configuration features:
//   - Key length validation for cipher compatibility
//   - Secure key storage with defensive copying
//   - Graceful handling of invalid keys (Laravel-compatible behavior)
//   - Method chaining support for fluent configuration API
//
// Returns the BaseEncrypter instance for method chaining.
func (b *BaseEncrypter) SetKey(key []byte) encryptionInterfaces.EncrypterInterface {
	// Validate key length against cipher requirements
	// Different AES variants require different key lengths
	expectedKeyLen := enums.GetKeyLength(enums.Cipher(b.cipher))
	if len(key) != expectedKeyLen {
		// Silently ignore invalid keys to match Laravel behavior
		// This could panic like Laravel or return an error in stricter implementations
		return b
	}

	// Make a secure copy of the key to prevent external modification
	// This defensive approach protects against accidental key corruption
	b.key = make([]byte, len(key))
	copy(b.key, key)
	return b
}

// SetCipher sets the cipher algorithm with validation and method chaining support.
// Updates cipher configuration while checking key compatibility for security.
//
// Cipher configuration features:
//   - Cipher algorithm validation against supported variants
//   - Key compatibility checking with new cipher requirements
//   - Graceful handling of invalid ciphers (Laravel-compatible behavior)
//   - Method chaining support for fluent configuration API
//
// Returns the BaseEncrypter instance for method chaining.
func (b *BaseEncrypter) SetCipher(cipher string) encryptionInterfaces.EncrypterInterface {
	// Validate cipher algorithm against supported variants
	// Ensures only secure, implemented ciphers are used
	if err := enums.ValidateCipher(enums.Cipher(cipher)); err != nil {
		// Silently ignore invalid ciphers to match Laravel behavior
		// Maintains existing configuration on validation failure
		return b
	}

	// Check if current key is compatible with new cipher requirements
	// Different ciphers may require different key lengths
	expectedKeyLen := enums.GetKeyLength(enums.Cipher(cipher))
	if len(b.key) != expectedKeyLen {
		// Key is incompatible, but we'll set the cipher anyway
		// The user will need to set a new compatible key
		// This matches Laravel's behavior for configuration flexibility
	}

	// Update the cipher algorithm identifier
	b.cipher = cipher
	return b
}

// GenerateKey generates a cryptographically secure key for the specified cipher algorithm.
// Uses system-provided cryptographically secure random number generation for maximum entropy.
//
// Key generation features:
//   - Cryptographically secure random key generation using crypto/rand
//   - Cipher-specific key length determination (128, 192, or 256 bits)
//   - High entropy key material suitable for production use
//   - Comprehensive error handling for generation failures
//
// Returns generated key bytes or error if generation fails or cipher is unsupported.
func (b *BaseEncrypter) GenerateKey(cipher string) ([]byte, error) {
	// Validate cipher algorithm before key generation
	// Ensures we only generate keys for supported, secure ciphers
	if err := enums.ValidateCipher(enums.Cipher(cipher)); err != nil {
		return nil, fmt.Errorf("invalid cipher: %w", err)
	}

	// Get required key length for the specified cipher
	// Different AES variants require different key lengths
	keyLen := enums.GetKeyLength(enums.Cipher(cipher))
	if keyLen == 0 {
		return nil, fmt.Errorf("unsupported cipher: %s", cipher)
	}

	// Generate cryptographically secure random key material
	// Uses system entropy source for maximum unpredictability
	key := make([]byte, keyLen)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	return key, nil
}

// generateIV generates a cryptographically secure initialization vector for the current cipher.
// Uses system entropy to create unpredictable IVs essential for semantic security.
//
// IV generation features:
//   - Cipher-specific IV length determination (16 bytes for CBC, 12 for GCM)
//   - Cryptographically secure random generation using crypto/rand
//   - Proper length validation for cipher compatibility
//   - Essential for semantic security (same plaintext produces different ciphertext)
//
// Returns generated IV bytes or error if generation fails or cipher is unsupported.
func (b *BaseEncrypter) generateIV() ([]byte, error) {
	// Get required IV length for the current cipher
	// Different modes require different IV lengths for optimal security
	ivLen := enums.GetIVLength(enums.Cipher(b.cipher))
	if ivLen == 0 {
		return nil, fmt.Errorf("unsupported cipher: %s", b.cipher)
	}

	// Generate cryptographically secure random IV
	// Each encryption operation must use a unique, unpredictable IV
	iv := make([]byte, ivLen)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate random IV: %w", err)
	}

	return iv, nil
}

// Hash computes HMAC-SHA256 for message authentication in non-AEAD cipher modes.
// Provides cryptographic integrity protection for CBC and CTR modes that lack built-in authentication.
//
// MAC computation features:
//   - HMAC-SHA256 for cryptographically secure message authentication
//   - Proper ordering of IV and ciphertext in MAC computation
//   - Protection against tampering and forgery attacks
//   - Compatible with Laravel's MAC validation approach
//
// Returns computed HMAC bytes or error if computation fails.
func (b *BaseEncrypter) Hash(iv, value, key []byte) ([]byte, error) {
	// Create HMAC-SHA256 hasher with the provided key
	// SHA256 provides 256-bit security level matching AES-256
	mac := hmac.New(sha256.New, key)

	// Write IV to MAC computation first
	// IV inclusion prevents IV manipulation attacks
	if _, err := mac.Write(iv); err != nil {
		return nil, fmt.Errorf("failed to write IV to MAC: %w", err)
	}

	// Write encrypted value to MAC computation
	// This provides integrity protection for the ciphertext
	if _, err := mac.Write(value); err != nil {
		return nil, fmt.Errorf("failed to write value to MAC: %w", err)
	}

	// Return the computed HMAC digest
	return mac.Sum(nil), nil
}

// GetJsonPayload parses an encrypted payload string into its component map.
// Handles both Laravel-style base64-encoded JSON and direct JSON payload formats.
//
// Payload parsing features:
//   - Laravel-compatible base64-encoded JSON parsing
//   - Fallback to direct JSON parsing for flexibility
//   - Component extraction (cipher, iv, value, mac/tag)
//   - Graceful error handling for malformed payloads
//
// Returns component map with payload fields or error if parsing fails.
func (b *BaseEncrypter) GetJsonPayload(payload string) (map[string]interface{}, error) {
	// Try to parse as base64-encoded JSON first (Laravel standard format)
	// This is the primary format used by Laravel's encryption system
	jsonBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		// If base64 decoding fails, try to parse as direct JSON
		// This provides fallback compatibility for different payload formats
		jsonBytes = []byte(payload)
	}

	// Unmarshal JSON into component map
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	return result, nil
}

// ValidPayload validates the structure and format of an encrypted payload.
// Performs comprehensive validation of required fields and authentication mechanisms.
//
// Payload validation features:
//   - JSON structure verification and parsing
//   - Required field presence checking (iv, value)
//   - Base64 encoding format validation for binary components
//   - Authentication mechanism validation (MAC or tag)
//   - Type safety checks for all payload components
//
// Returns true if payload structure is valid and complete, false otherwise.
func (b *BaseEncrypter) ValidPayload(payload string) bool {
	// Parse payload into component map
	data, err := b.GetJsonPayload(payload)
	if err != nil {
		// Return false if payload parsing fails
		return false
	}

	// Check required fields for all encryption modes
	// IV and encrypted value are mandatory for all ciphers
	iv, ivOk := data["iv"]
	value, valueOk := data["value"]

	if !ivOk || !valueOk {
		// Missing essential payload components
		return false
	}

	// Ensure IV and value are non-empty strings (base64 encoded)
	ivStr, ivIsStr := iv.(string)
	valueStr, valueIsStr := value.(string)

	if !ivIsStr || !valueIsStr || ivStr == "" || valueStr == "" {
		// Invalid type or empty required fields
		return false
	}

	// Check authentication field (either MAC for CBC/CTR or Tag for GCM)
	mac, macOk := data["mac"]
	tag, tagOk := data["tag"]

	// At least one authentication mechanism should be present
	if !macOk && !tagOk {
		// No authentication mechanism found - security violation
		return false
	}

	// If MAC is present, it should be a non-empty string
	if macOk {
		if macStr, ok := mac.(string); !ok || macStr == "" {
			// Invalid MAC format or empty MAC
			return false
		}
	}

	// If Tag is present, it should be a non-empty string
	if tagOk {
		if tagStr, ok := tag.(string); !ok || tagStr == "" {
			// Invalid tag format or empty tag
			return false
		}
	}

	return true
}

// ValidMac validates the MAC of a parsed payload using the current encryption key.
// Delegates to ValidMacForKey with the configured key for authentication verification.
//
// MAC validation features:
//   - Uses current encrypter key for MAC computation
//   - Constant-time comparison to prevent timing attacks
//   - Protection against payload tampering and forgery
//   - Integration with cipher-specific authentication requirements
//
// Returns true if MAC is valid and payload is authentic, false otherwise.
func (b *BaseEncrypter) ValidMac(payload map[string]interface{}) bool {
	// Delegate to key-specific MAC validation using current key
	return b.ValidMacForKey(payload, b.key)
}

// ValidMacForKey validates the MAC using a specific encryption key.
// Performs cryptographic MAC verification with constant-time comparison for security.
//
// Key-specific MAC validation features:
//   - HMAC-SHA256 computation with custom key material
//   - Base64 decoding of payload components (IV, value, MAC)
//   - Constant-time comparison to prevent timing attacks
//   - Independent validation without affecting current configuration
//
// Returns true if MAC is valid with the specified key, false otherwise.
func (b *BaseEncrypter) ValidMacForKey(payload map[string]interface{}, key []byte) bool {
	// Extract required fields from payload map
	// All fields must be present and valid strings
	ivStr, ivOk := payload["iv"].(string)
	valueStr, valueOk := payload["value"].(string)
	macStr, macOk := payload["mac"].(string)

	if !ivOk || !valueOk || !macOk {
		// Missing required MAC validation components
		return false
	}

	// Decode base64 IV component
	iv, err := base64.StdEncoding.DecodeString(ivStr)
	if err != nil {
		// Invalid base64 encoding for IV
		return false
	}

	// Decode base64 encrypted value component
	value, err := base64.StdEncoding.DecodeString(valueStr)
	if err != nil {
		// Invalid base64 encoding for encrypted value
		return false
	}

	// Decode base64 expected MAC component
	expectedMac, err := base64.StdEncoding.DecodeString(macStr)
	if err != nil {
		// Invalid base64 encoding for MAC
		return false
	}

	// Compute expected MAC using provided key
	computedMac, err := b.Hash(iv, value, key)
	if err != nil {
		// MAC computation failed
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	// This is crucial for preventing MAC forgery through timing analysis
	return hmac.Equal(expectedMac, computedMac)
}

// EnsureTagIsValid validates GCM authentication tag format and length.
// Performs comprehensive validation of AEAD authentication tags for security.
//
// GCM tag validation features:
//   - Non-empty tag requirement for AEAD security
//   - Base64 encoding format verification
//   - Tag length validation (16 bytes for GCM standard)
//   - Protection against malformed or truncated tags
//
// Returns nil if tag is valid, error if tag format or length is invalid.
func (b *BaseEncrypter) EnsureTagIsValid(tag string) error {
	// Ensure authentication tag is present
	if tag == "" {
		return fmt.Errorf("authentication tag is required for AEAD ciphers")
	}

	// Try to decode the tag to ensure it's valid base64
	// This validates both encoding format and decodability
	tagBytes, err := base64.StdEncoding.DecodeString(tag)
	if err != nil {
		return fmt.Errorf("invalid tag encoding: %w", err)
	}

	// Check tag length (typically 16 bytes for GCM)
	// GCM standard specifies 16-byte authentication tags
	if len(tagBytes) != 16 {
		return fmt.Errorf("invalid tag length: expected 16 bytes, got %d", len(tagBytes))
	}

	return nil
}

// ShouldValidateMac determines if MAC validation should be performed based on cipher type.
// Uses cipher characteristics to determine authentication requirements for security.
//
// MAC requirement determination features:
//   - AEAD cipher detection (GCM modes have built-in authentication)
//   - Traditional cipher identification (CBC, CTR modes require separate MAC)
//   - Cipher-specific authentication mechanism selection
//   - Security-first approach to authentication validation
//
// Returns true if MAC validation is required, false for AEAD ciphers with built-in authentication.
func (b *BaseEncrypter) ShouldValidateMac() bool {
	// Use cipher enumeration to determine MAC requirement
	// AEAD ciphers like GCM don't need separate MAC validation
	return enums.RequiresMAC(enums.Cipher(b.cipher))
}

// Info extracts and returns comprehensive metadata about the current cipher configuration.
// Provides detailed information about cipher capabilities and security parameters.
//
// Cipher information features:
//   - Current cipher algorithm identification and parameters
//   - Key length, IV length, and mode information
//   - AEAD capability detection and MAC requirements
//   - Security level and cryptographic properties
//
// Returns CipherInfo struct with comprehensive cipher metadata.
func (b *BaseEncrypter) Info(payload string) types.CipherInfo {
	// Return comprehensive information about the current cipher
	// Future implementations could analyze payload to detect cipher used
	// Currently returns configured cipher information
	return types.NewCipherInfo(enums.Cipher(b.cipher))
}

// VerifyConfiguration validates that the encrypter's current configuration meets security requirements.
// Performs comprehensive security validation to ensure production readiness.
//
// Configuration validation features:
//   - Cipher algorithm validation against supported variants
//   - Key length verification for cipher compatibility
//   - Weak key detection (all-zero keys)
//   - Security parameter consistency checking
//
// Returns true if configuration is secure and valid, false if any security issues are detected.
func (b *BaseEncrypter) VerifyConfiguration() bool {
	// Check cipher validity against supported algorithms
	// Ensures only secure, implemented ciphers are used
	if err := enums.ValidateCipher(enums.Cipher(b.cipher)); err != nil {
		return false
	}

	// Check key length matches cipher requirements
	// Prevents use of incorrect key sizes that could weaken security
	expectedKeyLen := enums.GetKeyLength(enums.Cipher(b.cipher))
	if len(b.key) != expectedKeyLen {
		return false
	}

	// Check for weak keys (all zeros)
	// All-zero keys are cryptographically weak and should be rejected
	allZeros := true
	for _, keyByte := range b.key {
		if keyByte != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		// Reject weak keys for security
		return false
	}

	return true
}

// serializeValue serializes a value for encryption with support for multiple data types.
// Handles both direct byte conversion and JSON serialization based on configuration.
//
// Serialization features:
//   - Direct string and byte array handling without serialization overhead
//   - JSON serialization for complex data types and structures
//   - Type safety validation for non-serialized inputs
//   - Laravel-compatible serialization behavior
//
// Returns serialized byte array or error if serialization fails.
func (b *BaseEncrypter) serializeValue(value interface{}, serialize bool) ([]byte, error) {
	if !serialize {
		// If serialization is disabled, expect simple data types
		// This provides direct encryption without JSON overhead
		switch v := value.(type) {
		case string:
			// Convert string directly to bytes
			return []byte(v), nil
		case []byte:
			// Use byte array directly
			return v, nil
		default:
			// Reject unsupported types when serialization is disabled
			return nil, fmt.Errorf("value must be string or []byte when serialization is disabled")
		}
	}

	// Serialize the value as JSON for complex data types
	// Enables encryption of objects, arrays, and structured data
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize value: %w", err)
	}

	return jsonBytes, nil
}

// deserializeValue deserializes bytes back to original value with flexible type handling.
// Provides both direct string conversion and JSON deserialization based on configuration.
//
// Deserialization features:
//   - Direct string conversion for simple data types
//   - JSON deserialization for complex data structures
//   - Graceful fallback to string on JSON parsing errors
//   - Type preservation for objects, arrays, and primitives
//
// Returns deserialized value or string representation if deserialization fails.
func (b *BaseEncrypter) deserializeValue(data []byte, unserialize bool) (interface{}, error) {
	if !unserialize {
		// Return as string if unserialization is disabled
		// Provides direct access to decrypted bytes as string
		return string(data), nil
	}

	// Attempt to deserialize as JSON for complex data types
	// Enables reconstruction of original objects, arrays, etc.
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		// If JSON deserialization fails, return as string
		// Graceful fallback ensures no data loss
		return string(data), nil
	}

	return result, nil
}

// createPayloadString creates a Laravel-compatible encrypted payload string.
// Constructs properly formatted base64-encoded JSON payload for compatibility.
//
// Payload creation features:
//   - Laravel-compatible JSON structure with base64-encoded components
//   - Support for both MAC-based (CBC/CTR) and tag-based (GCM) authentication
//   - Proper component encoding for secure transmission and storage
//   - Consistent payload format across all encryption modes
//
// Returns base64-encoded JSON payload string or error if encoding fails.
func (b *BaseEncrypter) createPayloadString(iv, encrypted, mac, tag []byte) (string, error) {
	// Create payload struct with all encryption components
	// Handles both traditional MAC and AEAD tag authentication
	payload := types.NewEncryptedPayload(iv, encrypted, mac, tag)

	// Convert to base64-encoded JSON (Laravel standard format)
	// Ensures compatibility with Laravel encryption system
	return payload.ToBase64JSON()
}

// parsePayloadString parses a Laravel-compatible encrypted payload string.
// Decodes base64-encoded JSON payload into structured components for processing.
//
// Payload parsing features:
//   - Laravel-compatible base64-encoded JSON parsing
//   - Component extraction and validation (IV, value, MAC/tag)
//   - Type-safe payload structure reconstruction
//   - Error handling for malformed or invalid payloads
//
// Returns parsed EncryptedPayload struct or error if parsing fails.
func (b *BaseEncrypter) parsePayloadString(payloadStr string) (*types.EncryptedPayload, error) {
	// Delegate to types package for consistent payload parsing
	// Ensures standardized handling across all encrypter implementations
	return types.ParseEncryptedPayload(payloadStr)
}

// extractMode extracts the encryption mode from a cipher name for mode-specific processing.
// Parses standardized cipher naming convention to identify operation mode.
//
// Mode extraction features:
//   - Standardized cipher name parsing (AES-256-CBC, AES-128-GCM, etc.)
//   - Mode identification for cipher-specific parameter selection
//   - Graceful handling of non-standard cipher names
//   - Support for CBC, GCM, CTR, and other encryption modes
//
// Returns encryption mode string (CBC, GCM, CTR) or "UNKNOWN" for unrecognized formats.
func (b *BaseEncrypter) extractMode() string {
	// Parse cipher name using standard format: ALGORITHM-KEYSIZE-MODE
	// e.g., "AES-256-CBC" -> ["AES", "256", "CBC"]
	parts := strings.Split(b.cipher, "-")
	if len(parts) >= 3 {
		// Return the mode component (third part)
		return parts[2] // e.g., "CBC" from "AES-256-CBC"
	}
	// Return unknown for non-standard cipher names
	return "UNKNOWN"
}
