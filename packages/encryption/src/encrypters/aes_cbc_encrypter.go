package encrypters

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"govel/packages/encryption/src/exceptions"
	enums "govel/packages/types/src/enums/encryption"
	encryptionInterfaces "govel/packages/types/src/interfaces/encryption"
)

// AESCBCEncrypter implements AES encryption in CBC (Cipher Block Chaining) mode.
// This encrypter is compatible with Laravel's default encryption and provides
// authenticated encryption using HMAC-SHA256 for integrity verification.
//
// Key features:
//   - AES-128-CBC and AES-256-CBC support
//   - PKCS#7 padding for block alignment
//   - HMAC-SHA256 for authentication
//   - Laravel-compatible encrypted payload format
//   - Cryptographically secure IV generation
//
// Security considerations:
//   - Uses random IV for each encryption operation
//   - Validates MAC before decryption to prevent tampering
//   - Constant-time MAC comparison to prevent timing attacks
//   - Proper error handling to prevent information leakage
type AESCBCEncrypter struct {
	*BaseEncrypter
}

// NewAESCBCEncrypter creates a new AES-CBC encrypter with the given key and cipher.
// Validates that the cipher is a supported CBC variant and the key length is appropriate.
//
// Supported ciphers:
//   - AES-128-CBC: 16-byte key
//   - AES-256-CBC: 32-byte key
//
// Parameters:
//   - key: The encryption key bytes
//   - cipher: The cipher algorithm identifier (must be CBC variant)
//   - config: Optional configuration parameters
//
// Returns a configured AESCBCEncrypter or an error if validation fails.
func NewAESCBCEncrypter(key []byte, cipher string, config map[string]interface{}) (*AESCBCEncrypter, error) {
	// Validate that this is a CBC cipher
	if cipher != string(enums.CipherAES128CBC) && cipher != string(enums.CipherAES256CBC) {
		return nil, fmt.Errorf("cipher %s is not supported by AES-CBC encrypter", cipher)
	}

	// Create base encrypter
	baseEncrypter, err := NewBaseEncrypter(key, cipher, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create base encrypter: %w", err)
	}

	return &AESCBCEncrypter{
		BaseEncrypter: baseEncrypter,
	}, nil
}

// Encrypt encrypts the given value using AES-CBC mode with comprehensive security features.
// Implements authenticated encryption using PKCS#7 padding and HMAC-SHA256 for integrity protection.
//
// AES-CBC encryption features:
//   - PKCS#7 padding for proper block alignment
//   - Cryptographically secure random IV generation for each operation
//   - HMAC-SHA256 authentication to prevent tampering
//   - Laravel-compatible encrypted payload format
//   - Support for both serialized and direct string encryption
//
// Returns base64-encoded JSON payload or error if any encryption step fails.
func (e *AESCBCEncrypter) Encrypt(value interface{}, serialize bool) (string, error) {
	// Serialize the value based on configuration
	// Handles both direct string encryption and JSON serialization
	plaintext, err := e.serializeValue(value, serialize)
	if err != nil {
		return "", fmt.Errorf("serialization failed: %w", err)
	}

	// Generate cryptographically secure random IV
	// Each encryption operation must use a unique IV for semantic security
	iv, err := e.generateIV()
	if err != nil {
		return "", fmt.Errorf("IV generation failed: %w", err)
	}

	// Encrypt the plaintext using AES-CBC with PKCS#7 padding
	encrypted, err := e.encryptCBC(plaintext, iv)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Compute HMAC-SHA256 MAC for authentication and integrity protection
	// MAC covers both IV and encrypted data to prevent tampering
	mac, err := e.Hash(iv, encrypted, e.key)
	if err != nil {
		return "", fmt.Errorf("MAC computation failed: %w", err)
	}

	// Create Laravel-compatible base64-encoded JSON payload
	return e.createPayloadString(iv, encrypted, mac, nil)
}

// EncryptString encrypts a string value directly without JSON serialization overhead.
// Optimized for string data encryption with AES-CBC mode and HMAC authentication.
//
// String encryption features:
//   - Direct string-to-bytes conversion without serialization
//   - Full AES-CBC security with PKCS#7 padding and MAC
//   - Optimal performance for string data types
//   - Laravel-compatible encrypted payload format
//
// Returns base64-encoded JSON payload or error if encryption fails.
func (e *AESCBCEncrypter) EncryptString(value string) (string, error) {
	// Delegate to main encrypt method with serialization disabled
	// This avoids JSON overhead for simple string encryption
	return e.Encrypt(value, false)
}

// Decrypt decrypts an AES-CBC encrypted payload with comprehensive security validation.
// Performs MAC authentication before decryption to prevent tampering and chosen-ciphertext attacks.
//
// AES-CBC decryption features:
//   - MAC validation before decryption to prevent tampering
//   - Payload structure validation and component extraction
//   - AES-CBC decryption with PKCS#7 padding removal
//   - Optional JSON deserialization support
//   - Protection against padding oracle attacks through MAC verification
//
// Returns decrypted value or error if validation or decryption fails.
func (e *AESCBCEncrypter) Decrypt(payload string, unserialize bool) (interface{}, error) {
	// Parse the Laravel-compatible base64-encoded JSON payload
	parsedPayload, err := e.parsePayloadString(payload)
	if err != nil {
		return nil, fmt.Errorf("payload parsing failed: %w", err)
	}

	// Validate the payload structure for required components
	if !parsedPayload.IsValid() {
		return nil, exceptions.ErrInvalidPayload
	}

	// CBC mode requires MAC validation for authenticated encryption
	// Without MAC, CBC is vulnerable to padding oracle attacks
	if !parsedPayload.HasMac() {
		return nil, fmt.Errorf("MAC is required for CBC mode")
	}

	// Extract and validate IV component
	iv, err := parsedPayload.GetIVBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid IV: %w", err)
	}

	// Extract and validate encrypted value component
	encrypted, err := parsedPayload.GetValueBytes()
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted value: %w", err)
	}

	// Validate HMAC-SHA256 MAC before attempting decryption
	// This prevents padding oracle attacks and ensures payload integrity
	payloadMap := map[string]interface{}{
		"iv":    parsedPayload.IV,
		"value": parsedPayload.Value,
		"mac":   parsedPayload.MAC,
	}

	if !e.ValidMac(payloadMap) {
		return nil, exceptions.ErrMacValidationFailed
	}

	// Decrypt the data using AES-CBC and remove PKCS#7 padding
	plaintext, err := e.decryptCBC(encrypted, iv)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Deserialize the result based on configuration
	return e.deserializeValue(plaintext, unserialize)
}

// DecryptString decrypts an AES-CBC encrypted payload directly to string format.
// Optimized for string data decryption without JSON deserialization overhead.
//
// String decryption features:
//   - Full AES-CBC security with MAC validation and padding removal
//   - Direct string result without deserialization
//   - Type-safe string conversion with fallback formatting
//   - Optimal performance for string data types
//
// Returns decrypted string or error if decryption fails.
func (e *AESCBCEncrypter) DecryptString(payload string) (string, error) {
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

// encryptCBC performs AES-CBC encryption with PKCS#7 padding for block alignment.
// Implements the core CBC encryption logic with proper padding and block processing.
//
// CBC encryption features:
//   - AES cipher initialization with configured key
//   - PKCS#7 padding application for block size alignment
//   - CBC mode encryption with initialization vector
//   - Secure block-by-block processing for large plaintexts
//
// Returns encrypted ciphertext bytes or error if encryption fails.
func (e *AESCBCEncrypter) encryptCBC(plaintext, iv []byte) ([]byte, error) {
	// Create AES cipher block with the encryption key
	// This validates key length and initializes the cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Apply PKCS#7 padding to ensure plaintext fits block boundaries
	// CBC requires input length to be multiple of block size (16 bytes)
	paddedPlaintext := pkcs7Pad(plaintext, aes.BlockSize)

	// Create CBC encrypter mode with the IV
	// IV must be exactly one block size (16 bytes) for AES
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the padded plaintext using CBC mode
	// CBC processes data in blocks, chaining each block with the previous
	encrypted := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(encrypted, paddedPlaintext)

	return encrypted, nil
}

// decryptCBC performs AES-CBC decryption and removes PKCS#7 padding to recover plaintext.
// Implements the core CBC decryption logic with validation and padding removal.
//
// CBC decryption features:
//   - Input length validation for block size alignment
//   - IV length validation for CBC requirements
//   - AES cipher initialization and CBC mode setup
//   - Block-by-block decryption processing
//   - PKCS#7 padding validation and removal
//
// Returns decrypted plaintext bytes or error if decryption or padding removal fails.
func (e *AESCBCEncrypter) decryptCBC(encrypted, iv []byte) ([]byte, error) {
	// Validate encrypted data length must be multiple of block size
	// CBC requires aligned block boundaries for proper decryption
	if len(encrypted)%aes.BlockSize != 0 {
		return nil, exceptions.ErrDecryptionFailed
	}

	// Validate IV length must exactly match AES block size
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV length: expected %d, got %d", aes.BlockSize, len(iv))
	}

	// Create AES cipher block with the decryption key
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create CBC decrypter mode with the IV
	// CBC decryption reverses the chaining process used in encryption
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the ciphertext using CBC mode
	// CBC decrypts each block and applies XOR with previous ciphertext block
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove PKCS#7 padding and validate padding correctness
	// Invalid padding may indicate tampering or corruption
	plaintext, err := pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return nil, exceptions.ErrDecryptionFailed
	}

	return plaintext, nil
}

// pkcs7Pad applies PKCS#7 padding to align data with block boundaries.
// Implements RFC 5652 PKCS#7 padding standard for block cipher compatibility.
//
// PKCS#7 padding features:
//   - Deterministic padding length calculation based on data length
//   - Padding bytes contain the padding length value
//   - Always adds at least 1 byte of padding (1-16 bytes for AES)
//   - Enables unambiguous padding removal during decryption
//
// Returns padded data bytes aligned to block size boundaries.
func pkcs7Pad(data []byte, blockSize int) []byte {
	// Calculate padding length needed to reach block boundary
	// Padding is always 1-blockSize bytes, never 0
	padding := blockSize - (len(data) % blockSize)

	// Create padding bytes, each containing the padding length
	// This allows unambiguous padding detection during removal
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}

	// Append padding to original data
	return append(data, padText...)
}

// pkcs7Unpad removes and validates PKCS#7 padding from decrypted data.
// Implements secure padding validation to detect tampering and ensure data integrity.
//
// PKCS#7 padding validation features:
//   - Data length validation for proper block alignment
//   - Padding length extraction from final byte
//   - Range validation for padding length (1-blockSize)
//   - Byte-by-byte padding content verification
//   - Protection against padding oracle attacks through consistent validation
//
// Returns unpadded plaintext bytes or error if padding is invalid or corrupted.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	// Validate input data length and block alignment
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data length")
	}

	// Get the padding length from the last byte (PKCS#7 standard)
	// All padding bytes contain the padding length value
	paddingLen := int(data[len(data)-1])

	// Validate padding length is within acceptable range
	// Padding must be 1-blockSize bytes for valid PKCS#7
	if paddingLen == 0 || paddingLen > blockSize || paddingLen > len(data) {
		return nil, fmt.Errorf("invalid padding length: %d", paddingLen)
	}

	// Verify all padding bytes contain the correct padding length value
	// This protects against padding corruption and tampering
	for i := len(data) - paddingLen; i < len(data); i++ {
		if data[i] != byte(paddingLen) {
			return nil, fmt.Errorf("invalid padding at position %d", i)
		}
	}

	// Remove validated padding to recover original plaintext
	return data[:len(data)-paddingLen], nil
}

// Compile-time interface compliance check
var _ encryptionInterfaces.EncrypterInterface = (*AESCBCEncrypter)(nil)
