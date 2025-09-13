package encrypters

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"

	types "govel/types"
	enums "govel/types/enums/encryption"
	encryptionInterfaces "govel/types/interfaces/encryption"
)

// AESCTREncrypter provides AES encryption in CTR mode with HMAC-SHA256 MAC authentication.
//
// This implementation supports:
// - AES-128-CTR and AES-256-CTR encryption
// - HMAC-SHA256 MAC for authentication
// - Laravel-compatible payload format
// - Constant-time MAC validation
// - IV generation and management
//
// CTR mode provides:
// - Stream cipher characteristics (no padding needed)
// - Parallel encryption/decryption
// - Random access to encrypted data
// - Deterministic encryption with proper IV handling
type AESCTREncrypter struct {
	*BaseEncrypter
}

// NewAESCTREncrypter creates a new AES-CTR encrypter with comprehensive validation.
// Initializes the encrypter with proper configuration and cipher validation.
//
// Constructor features:
//   - Key length validation for AES-128 (16 bytes) and AES-256 (32 bytes)
//   - Cipher string validation against supported CTR variants
//   - Base encrypter initialization with shared functionality
//   - Comprehensive error handling for invalid configurations
//
// Supported cipher configurations:
//   - AES-128-CTR: 16-byte key with 16-byte IV
//   - AES-256-CTR: 32-byte key with 16-byte IV
//
// Returns configured AESCTREncrypter or error if validation fails.
func NewAESCTREncrypter(key []byte, cipher string) (*AESCTREncrypter, error) {
	// Initialize base encrypter with shared functionality
	base, err := NewBaseEncrypter(key, cipher, nil)
	if err != nil {
		return nil, err
	}

	// Validate that the cipher is a supported CTR variant
	// Only AES-128-CTR and AES-256-CTR are currently supported
	if cipher != string(enums.CipherAES256CTR) && cipher != string(enums.CipherAES128CTR) {
		return nil, fmt.Errorf(
			"Unsupported CTR cipher: %s. Only AES-128-CTR and AES-256-CTR are supported", cipher)
	}

	// Return configured CTR encrypter instance
	return &AESCTREncrypter{
		BaseEncrypter: base,
	}, nil
}

// Encrypt encrypts a value using AES-CTR mode with comprehensive security features.
// Implements stream cipher encryption with HMAC-SHA256 authentication for integrity protection.
//
// AES-CTR encryption features:
//   - Stream cipher properties (no padding required)
//   - Cryptographically secure random IV generation for each operation
//   - XOR-based encryption allowing parallel processing
//   - HMAC-SHA256 authentication to prevent tampering
//   - Laravel-compatible encrypted payload format
//   - Support for both serialized and direct data encryption
//
// CTR mode encryption process:
//  1. Serialize value if requested (JSON or direct conversion)
//  2. Generate cryptographically secure random IV
//  3. Initialize AES cipher with encryption key
//  4. Create CTR mode stream cipher with IV
//  5. XOR plaintext with cipher stream to produce ciphertext
//  6. Compute HMAC-SHA256 MAC over IV and ciphertext
//  7. Create base64-encoded JSON payload
//
// Returns base64-encoded JSON payload or error if any encryption step fails.
func (e *AESCTREncrypter) Encrypt(value interface{}, serialize bool) (string, error) {
	// Serialize the value based on configuration
	// Handles both direct conversion and JSON serialization
	plaintext, err := e.serializeValue(value, serialize)
	if err != nil {
		return "", err
	}

	// Generate cryptographically secure random IV
	// Each encryption operation must use a unique IV for semantic security
	iv, err := e.generateIV()
	if err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Create AES cipher block with the encryption key
	block, err := aes.NewCipher(e.GetKey())
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create CTR mode stream cipher with the IV
	// CTR mode converts block cipher into stream cipher
	stream := cipher.NewCTR(block, iv)

	// Encrypt the plaintext using XOR with cipher stream
	// CTR mode allows parallel encryption and random access
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// Compute HMAC-SHA256 MAC for authentication and integrity
	// MAC covers both IV and ciphertext to prevent tampering
	macHash, err := e.Hash(iv, ciphertext, e.GetKey())
	if err != nil {
		return "", fmt.Errorf("failed to compute MAC: %w", err)
	}

	// Create Laravel-compatible base64-encoded JSON payload
	return e.createPayloadString(iv, ciphertext, macHash, nil)
}

// EncryptString encrypts a string value directly using AES-CTR mode without serialization.
// Optimized for string data encryption with stream cipher properties and MAC authentication.
//
// String encryption features:
//   - Direct string-to-bytes conversion without JSON serialization overhead
//   - Full AES-CTR security with stream cipher properties
//   - HMAC-SHA256 authentication for integrity protection
//   - Optimal performance for string data types
//   - Laravel-compatible encrypted payload format
//
// Returns base64-encoded JSON payload or error if encryption fails.
func (e *AESCTREncrypter) EncryptString(value string) (string, error) {
	// Delegate to main encrypt method with serialization disabled
	// This provides direct string encryption without JSON overhead
	return e.Encrypt(value, false)
}

// Decrypt decrypts an AES-CTR encrypted payload with comprehensive security validation.
// Performs MAC authentication before decryption to prevent tampering and ensure integrity.
//
// AES-CTR decryption features:
//   - MAC validation before decryption to prevent tampering
//   - Payload structure validation and component extraction
//   - Stream cipher decryption with XOR operations
//   - IV length validation for cipher compatibility
//   - Optional JSON deserialization support
//   - Protection against chosen-ciphertext attacks through MAC verification
//
// CTR mode decryption process:
//  1. Parse and validate base64-encoded JSON payload structure
//  2. Verify HMAC-SHA256 MAC for authentication and integrity
//  3. Decode base64-encoded IV and ciphertext components
//  4. Validate IV length matches cipher requirements
//  5. Initialize AES cipher and CTR mode stream
//  6. XOR ciphertext with cipher stream to recover plaintext
//  7. Deserialize result if requested
//
// Returns decrypted value or error if validation or decryption fails.
func (e *AESCTREncrypter) Decrypt(payload string, unserialize bool) (interface{}, error) {
	// Parse the Laravel-compatible base64-encoded JSON payload
	parsedPayload, err := e.GetJsonPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate payload structure for required components
	if !e.ValidPayload(payload) {
		return nil, fmt.Errorf("Invalid payload format")
	}

	// Validate HMAC-SHA256 MAC before attempting decryption
	// This prevents chosen-ciphertext attacks and ensures payload integrity
	if e.ShouldValidateMac() && !e.ValidMac(parsedPayload) {
		return nil, fmt.Errorf("MAC validation failed")
	}

	// Decode base64-encoded IV component
	iv, err := base64.StdEncoding.DecodeString(parsedPayload["iv"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	// Decode base64-encoded ciphertext component
	ciphertext, err := base64.StdEncoding.DecodeString(parsedPayload["value"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Validate IV length matches cipher requirements
	// CTR mode typically uses 16-byte IVs for AES
	expectedIVLength := enums.GetIVLength(enums.Cipher(e.GetCipher()))
	if len(iv) != expectedIVLength {
		return nil, fmt.Errorf(
			"Invalid IV length: expected %d, got %d", expectedIVLength, len(iv))
	}

	// Create AES cipher block with the decryption key
	block, err := aes.NewCipher(e.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create CTR mode stream cipher with the IV
	// CTR decryption uses the same operation as encryption (XOR)
	stream := cipher.NewCTR(block, iv)

	// Decrypt the ciphertext using XOR with cipher stream
	// CTR mode decryption is identical to encryption operation
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	// Deserialize the result based on configuration
	return e.deserializeValue(plaintext, unserialize)
}

// DecryptString decrypts an AES-CTR encrypted payload directly to string format.
// Optimized for string data decryption without JSON deserialization overhead.
//
// String decryption features:
//   - Full AES-CTR security with MAC validation and stream cipher decryption
//   - Direct string result without deserialization
//   - Type-safe string conversion with strict validation
//   - Optimal performance for string data types
//   - Error handling for type mismatches
//
// Returns decrypted string or error if decryption fails or result is not a string.
func (e *AESCTREncrypter) DecryptString(payload string) (string, error) {
	// Decrypt using main method with deserialization disabled
	result, err := e.Decrypt(payload, false)
	if err != nil {
		return "", err
	}

	// Perform strict type checking for string result
	if str, ok := result.(string); ok {
		return str, nil
	}

	// Return error for non-string results to maintain type safety
	return "", fmt.Errorf("Decrypted value is not a string")
}

// Hash computes HMAC-SHA256 hash for MAC validation in CTR mode.
//
// The MAC is computed over base64(IV || ciphertext) using the encryption key.
// This provides authentication for CTR mode which doesn't have built-in authentication.
//
// Parameters:
//   - iv: Initialization vector bytes
//   - value: Ciphertext bytes
//   - key: HMAC key (encryption key)
//
// Returns:
//   - []byte: HMAC-SHA256 hash bytes
//   - error: Any hashing error
func (e *AESCTREncrypter) Hash(iv, value, key []byte) ([]byte, error) {
	// Create HMAC-SHA256 hasher
	var hasher hash.Hash = hmac.New(sha256.New, key)

	// Compute MAC over base64(IV || ciphertext)
	macData := base64.StdEncoding.EncodeToString(append(iv, value...))
	hasher.Write([]byte(macData))

	return hasher.Sum(nil), nil
}

// ShouldValidateMac returns true since CTR mode requires MAC validation.
//
// Unlike GCM mode which has built-in authentication, CTR mode requires
// separate MAC validation for message integrity and authenticity.
//
// Returns:
//   - bool: Always true for CTR mode
func (e *AESCTREncrypter) ShouldValidateMac() bool {
	return true
}

// VerifyConfiguration validates the AES-CTR encrypter configuration for security and correctness.
// Performs comprehensive validation beyond base encrypter checks for CTR-specific requirements.
//
// CTR-specific validation features:
//   - Base encrypter validation (key length, cipher support, weak key detection)
//   - CTR cipher variant validation (AES-128-CTR, AES-256-CTR)
//   - Key-cipher compatibility verification
//   - Security parameter consistency checking
//
// Returns true if configuration is secure and valid for CTR mode, false if any issues detected.
func (e *AESCTREncrypter) VerifyConfiguration() bool {
	// Perform base encrypter validation first
	// This checks key length, cipher validity, and weak key detection
	if !e.BaseEncrypter.VerifyConfiguration() {
		return false
	}

	// Additional CTR-specific validation
	// Ensure cipher is a supported CTR variant
	cipher := e.GetCipher()
	if cipher != string(enums.CipherAES256CTR) && cipher != string(enums.CipherAES128CTR) {
		return false
	}

	return true
}

// Info returns comprehensive information about the AES-CTR encryption configuration.
// Provides detailed metadata about cipher capabilities, security parameters, and operational characteristics.
//
// CTR configuration information features:
//   - Current cipher algorithm identification (AES-128-CTR, AES-256-CTR)
//   - Key length, IV length, and block size information
//   - Stream cipher properties and MAC requirements
//   - Security level and cryptographic characteristics
//   - Operational mode details (counter mode, parallel processing)
//
// Returns CipherInfo struct with comprehensive CTR mode metadata.
func (e *AESCTREncrypter) Info(payload string) types.CipherInfo {
	// Return comprehensive cipher information for CTR mode
	// Includes all relevant parameters and capabilities
	return types.NewCipherInfo(enums.Cipher(e.GetCipher()))
}

// SetKey sets a new encryption key with validation and method chaining support.
// Delegates to base encrypter for key validation and secure storage.
//
// Key configuration features:
//   - Key length validation for CTR cipher compatibility
//   - Secure key storage with defensive copying
//   - Method chaining support for fluent configuration API
//   - Graceful handling of invalid keys
//
// Returns AESCTREncrypter instance for method chaining.
func (e *AESCTREncrypter) SetKey(key []byte) encryptionInterfaces.EncrypterInterface {
	// Delegate to base encrypter for key validation and storage
	e.BaseEncrypter.SetKey(key)
	return e
}

// SetCipher sets a new cipher algorithm with validation and method chaining support.
// Delegates to base encrypter for cipher validation and configuration updates.
//
// Cipher configuration features:
//   - CTR cipher variant validation (AES-128-CTR, AES-256-CTR)
//   - Key compatibility checking with new cipher requirements
//   - Method chaining support for fluent configuration API
//   - Graceful handling of invalid ciphers
//
// Returns AESCTREncrypter instance for method chaining.
func (e *AESCTREncrypter) SetCipher(cipher string) encryptionInterfaces.EncrypterInterface {
	// Delegate to base encrypter for cipher validation and configuration
	e.BaseEncrypter.SetCipher(cipher)
	return e
}

// Parameters:
//   - cipher: The cipher name
//
// Returns:
//   - []byte: Generated key bytes
//   - error: Any generation error
func (e *AESCTREncrypter) GenerateKey(cipher string) ([]byte, error) {
	return e.BaseEncrypter.GenerateKey(cipher)
}

// Compile-time interface compliance check
// Ensures AESCTREncrypter properly implements the EncrypterInterface
// Prevents runtime errors from missing method implementations
var _ encryptionInterfaces.EncrypterInterface = (*AESCTREncrypter)(nil)
