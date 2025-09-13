package encrypters

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"govel/packages/encryption/src/exceptions"
	enums "govel/packages/types/src/enums/encryption"
	encryptionInterfaces "govel/packages/types/src/interfaces/encryption"
)

// AESGCMEncrypter implements AES encryption in GCM (Galois/Counter Mode) mode.
// This encrypter provides authenticated encryption with associated data (AEAD),
// combining encryption and authentication in a single operation.
//
// Key features:
//   - AES-128-GCM and AES-256-GCM support
//   - Built-in authentication (no separate MAC required)
//   - Resistant to padding oracle attacks
//   - High performance with parallel processing capabilities
//   - 96-bit (12-byte) IV for optimal security and performance
//
// Security considerations:
//   - Uses random IV for each encryption operation
//   - Authentication tag provides integrity and authenticity
//   - No padding required (stream cipher properties)
//   - Resistant to chosen-ciphertext attacks
type AESGCMEncrypter struct {
	*BaseEncrypter
}

// NewAESGCMEncrypter creates a new AES-GCM encrypter with the given key and cipher.
// Validates that the cipher is a supported GCM variant and the key length is appropriate.
//
// Supported ciphers:
//   - AES-128-GCM: 16-byte key
//   - AES-256-GCM: 32-byte key
//
// Parameters:
//   - key: The encryption key bytes
//   - cipher: The cipher algorithm identifier (must be GCM variant)
//   - config: Optional configuration parameters
//
// Returns a configured AESGCMEncrypter or an error if validation fails.
func NewAESGCMEncrypter(key []byte, cipher string, config map[string]interface{}) (*AESGCMEncrypter, error) {
	// Validate that this is a GCM cipher
	if cipher != string(enums.CipherAES128GCM) && cipher != string(enums.CipherAES256GCM) {
		return nil, fmt.Errorf("cipher %s is not supported by AES-GCM encrypter", cipher)
	}

	// Create base encrypter
	baseEncrypter, err := NewBaseEncrypter(key, cipher, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create base encrypter: %w", err)
	}

	return &AESGCMEncrypter{
		BaseEncrypter: baseEncrypter,
	}, nil
}

// Encrypt encrypts the given value using AES-GCM authenticated encryption with comprehensive security.
// Implements AEAD (Authenticated Encryption with Associated Data) providing both confidentiality and integrity.
//
// AES-GCM encryption features:
//   - AEAD properties with built-in authentication (no separate MAC required)
//   - Optimal 12-byte IV for maximum security and performance
//   - 16-byte authentication tag for integrity verification
//   - No padding required (stream cipher properties)
//   - Laravel-compatible encrypted payload format
//   - Support for both serialized and direct data encryption
//
// Returns base64-encoded JSON payload or error if any encryption step fails.
func (e *AESGCMEncrypter) Encrypt(value interface{}, serialize bool) (string, error) {
	// Serialize the value based on configuration
	// Handles both direct conversion and JSON serialization
	plaintext, err := e.serializeValue(value, serialize)
	if err != nil {
		return "", fmt.Errorf("serialization failed: %w", err)
	}

	// Generate cryptographically secure random IV (12 bytes optimal for GCM)
	// 96-bit IV provides optimal security-performance balance for GCM mode
	iv, err := e.generateIV()
	if err != nil {
		return "", fmt.Errorf("IV generation failed: %w", err)
	}

	// Encrypt the plaintext using AES-GCM authenticated encryption
	// GCM provides both confidentiality and authenticity in a single operation
	encrypted, tag, err := e.encryptGCM(plaintext, iv)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Create Laravel-compatible payload with authentication tag (no separate MAC needed)
	return e.createPayloadString(iv, encrypted, nil, tag)
}

// EncryptString encrypts a string value directly using AES-GCM without serialization overhead.
// Optimized for string data encryption with AEAD properties and built-in authentication.
//
// String encryption features:
//   - Direct string-to-bytes conversion without JSON serialization
//   - Full AES-GCM AEAD security with built-in authentication
//   - Optimal performance for string data types
//   - 16-byte authentication tag for integrity verification
//   - Laravel-compatible encrypted payload format
//
// Returns base64-encoded JSON payload or error if encryption fails.
func (e *AESGCMEncrypter) EncryptString(value string) (string, error) {
	// Delegate to main encrypt method with serialization disabled
	// Provides direct string encryption without JSON overhead
	return e.Encrypt(value, false)
}

// Decrypt decrypts an AES-GCM encrypted payload with comprehensive AEAD security validation.
// Performs authentication tag verification during decryption for built-in integrity protection.
//
// AES-GCM decryption features:
//   - AEAD authentication tag validation during decryption process
//   - Payload structure validation and component extraction
//   - GCM-specific IV length validation (12 bytes)
//   - Authentication tag format and length validation
//   - Optional JSON deserialization support
//   - Automatic failure on authentication mismatch (tampered data)
//
// Returns decrypted value or error if validation or decryption fails.
func (e *AESGCMEncrypter) Decrypt(payload string, unserialize bool) (interface{}, error) {
	// Parse the Laravel-compatible base64-encoded JSON payload
	parsedPayload, err := e.parsePayloadString(payload)
	if err != nil {
		return nil, fmt.Errorf("payload parsing failed: %w", err)
	}

	// Validate the payload structure for required GCM components
	if !parsedPayload.IsValid() {
		return nil, exceptions.ErrInvalidPayload
	}

	// GCM mode requires authentication tag for AEAD operation
	// Tag provides both authenticity and integrity verification
	if !parsedPayload.HasTag() {
		return nil, fmt.Errorf("authentication tag is required for GCM mode")
	}

	// Validate authentication tag format and length (16 bytes)
	if err := e.EnsureTagIsValid(parsedPayload.Tag); err != nil {
		return nil, fmt.Errorf("invalid authentication tag: %w", err)
	}

	// Extract and validate IV component (12 bytes for GCM)
	iv, err := parsedPayload.GetIVBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid IV: %w", err)
	}

	// Extract and validate encrypted value component
	encrypted, err := parsedPayload.GetValueBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted value: %w", err)
	}

	// Extract and validate authentication tag component
	tag, err := parsedPayload.GetTagBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid authentication tag: %w", err)
	}

	// Decrypt and authenticate the data in single GCM operation
	// GCM automatically validates authentication tag during decryption
	plaintext, err := e.decryptGCM(encrypted, iv, tag)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Deserialize the result based on configuration
	return e.deserializeValue(plaintext, unserialize)
}

// DecryptString decrypts an AES-GCM encrypted payload directly to string format.
// Optimized for string data decryption without JSON deserialization overhead.
//
// String decryption features:
//   - Full AES-GCM AEAD security with built-in authentication validation
//   - Direct string result without deserialization
//   - Type-safe string conversion with fallback formatting
//   - Optimal performance for string data types
//   - Automatic authentication tag verification during decryption
//
// Returns decrypted string or error if decryption fails.
func (e *AESGCMEncrypter) DecryptString(payload string) (string, error) {
	// Decrypt using main method with deserialization disabled
	result, err := e.Decrypt(payload, false)
	if err != nil {
		return "", err
	}

	// Perform type-safe string conversion
	if str, ok := result.(string); ok {
		return str, nil
	}

	// Fallback to formatted string representation
	return fmt.Sprintf("%v", result), nil
}

// encryptGCM performs AES-GCM authenticated encryption with integrated authentication.
// Implements the core GCM AEAD operation providing both confidentiality and authenticity.
//
// GCM encryption features:
//   - AES cipher initialization with configured key
//   - GCM mode creation for authenticated encryption
//   - Single-operation encryption and authentication (Seal)
//   - Automatic authentication tag generation (16 bytes)
//   - Ciphertext and tag separation for Laravel compatibility
//
// Returns encrypted ciphertext, authentication tag, or error if encryption fails.
func (e *AESGCMEncrypter) encryptGCM(plaintext, iv []byte) ([]byte, []byte, error) {
	// Create AES cipher block with the encryption key
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM authenticated encryption mode
	// GCM provides both encryption and authentication in single mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Encrypt and authenticate in single operation using Seal
	// GCM automatically appends authentication tag to ciphertext
	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	// Split ciphertext and authentication tag (tag is appended by Seal)
	// GCM tag size is typically 16 bytes (128 bits)
	tagSize := gcm.Overhead() // Usually 16 bytes for GCM
	if len(ciphertext) < tagSize {
		return nil, nil, fmt.Errorf("ciphertext too short")
	}

	// Separate encrypted data and authentication tag
	encrypted := ciphertext[:len(ciphertext)-tagSize]
	tag := ciphertext[len(ciphertext)-tagSize:]

	return encrypted, tag, nil
}

// decryptGCM performs AES-GCM authenticated decryption with integrated authentication verification.
// Implements the core GCM AEAD operation providing both decryption and authentication validation.
//
// GCM decryption features:
//   - IV length validation (12 bytes optimal for GCM)
//   - AES cipher initialization with decryption key
//   - GCM mode creation for authenticated decryption
//   - Ciphertext reconstruction (encrypted data + authentication tag)
//   - Single-operation decryption and authentication verification (Open)
//   - Automatic failure on authentication mismatch
//
// Returns decrypted plaintext or error if decryption or authentication fails.
func (e *AESGCMEncrypter) decryptGCM(encrypted, iv, tag []byte) ([]byte, error) {
	// Validate IV length for GCM (12 bytes provides optimal security)
	// 96-bit IV is recommended for GCM mode
	if len(iv) != 12 {
		return nil, fmt.Errorf("invalid IV length for GCM: expected 12 bytes, got %d", len(iv))
	}

	// Create AES cipher block with the decryption key
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM authenticated decryption mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Reconstruct the sealed ciphertext (encrypted data + authentication tag)
	// GCM Open expects ciphertext with appended authentication tag
	ciphertext := append(encrypted, tag...)

	// Decrypt and verify authentication in single operation using Open
	// GCM automatically validates authentication tag during decryption
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		// GCM Open() will fail if authentication tag verification fails
		// This indicates data tampering or corruption
		return nil, exceptions.ErrDecryptionFailed
	}

	return plaintext, nil
}

// ShouldValidateMac overrides the base implementation for GCM AEAD characteristics.
// GCM doesn't require separate MAC validation as it provides built-in authentication.
//
// GCM authentication features:
//   - AEAD properties with integrated authentication
//   - Authentication tag verification during decryption
//   - No separate HMAC validation required
//   - Built-in protection against tampering and forgery
//
// Returns false as GCM has built-in authentication, eliminating separate MAC validation.
func (e *AESGCMEncrypter) ShouldValidateMac() bool {
	// GCM provides built-in authentication through AEAD properties
	// No separate MAC validation is needed or desired
	return false
}

// ValidMac overrides the base implementation for GCM AEAD authentication model.
// GCM doesn't validate MAC separately as authentication is integrated into the cipher operation.
//
// GCM authentication handling:
//   - Authentication is built into GCM decryption process
//   - No separate MAC computation or validation required
//   - Authentication tag verification occurs during GCM Open operation
//   - Always returns true as MAC validation is not applicable to GCM
//
// Returns true as GCM handles authentication internally, not through separate MAC validation.
func (e *AESGCMEncrypter) ValidMac(payload map[string]interface{}) bool {
	// GCM doesn't use separate MAC validation
	// Authentication is handled automatically during decryption process
	return true
}

// ValidMacForKey overrides the base implementation for GCM AEAD authentication model.
// GCM doesn't validate MAC with specific keys as authentication is integrated into cipher operation.
//
// GCM key-specific authentication handling:
//   - Authentication key is the same as encryption key in GCM
//   - No separate MAC computation with different keys
//   - Authentication tag verification is built into GCM decryption
//   - Always returns true as separate MAC validation is not applicable
//
// Returns true as GCM handles authentication internally without separate key-based MAC validation.
func (e *AESGCMEncrypter) ValidMacForKey(payload map[string]interface{}, key []byte) bool {
	// GCM doesn't use separate MAC validation with specific keys
	// Authentication is integrated into the cipher operation
	return true
}

// EnsureTagIsValid validates the GCM authentication tag format and characteristics.
// Performs comprehensive validation specific to GCM AEAD authentication requirements.
//
// GCM tag validation features:
//   - Non-empty tag requirement for AEAD security
//   - Base64 encoding format verification
//   - Tag length validation (16 bytes for GCM standard)
//   - GCM-specific authentication tag characteristics
//
// Returns nil if tag is valid for GCM mode, error if tag format or length is invalid.
func (e *AESGCMEncrypter) EnsureTagIsValid(tag string) error {
	// Ensure authentication tag is present (required for GCM AEAD)
	if tag == "" {
		return fmt.Errorf("authentication tag is required for GCM mode")
	}

	// Delegate to base encrypter for standard tag validation
	// This validates base64 encoding and 16-byte length
	if err := e.BaseEncrypter.EnsureTagIsValid(tag); err != nil {
		return err
	}

	return nil
}

// Compile-time interface compliance check
var _ encryptionInterfaces.EncrypterInterface = (*AESGCMEncrypter)(nil)
