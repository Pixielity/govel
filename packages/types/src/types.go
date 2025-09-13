package types

import (
	"encoding/base64"
	"fmt"

	enums "govel/packages/types/src/enums/encryption"
	pipelineInterfaces "govel/packages/types/src/interfaces/pipeline"
)

// PipelineCallback represents a function that configures and executes a pipeline.
// This type is used by the pipeline hub to store named pipeline configurations.
//
// Parameters:
//   - pipeline: The pipeline instance to configure
//   - passable: The object to pass through the pipeline
//
// Returns:
//   - interface{}: The result of pipeline execution
type PipelineCallback func(pipeline pipelineInterfaces.PipelineInterface, passable interface{}) interface{}

// EncryptedPayload represents a structured encrypted payload.
// This type is used to represent the parsed components of an encrypted payload.
type EncryptedPayload struct {
	// Cipher is the encryption algorithm used
	Cipher string `json:"cipher"`

	// IV is the Base64-encoded initialization vector
	IV string `json:"iv"`

	// Value is the Base64-encoded encrypted value
	Value string `json:"value"`

	// MAC is the Base64-encoded message authentication code
	MAC string `json:"mac,omitempty"`

	// Tag is the Base64-encoded authentication tag (for AEAD ciphers like GCM)
	Tag string `json:"tag,omitempty"`
}

// NewEncryptedPayload creates a new EncryptedPayload from the given components.
// Automatically encodes binary data to Base64 format for JSON serialization.
//
// Parameters:
//   - iv: Initialization vector bytes
//   - encrypted: Encrypted data bytes
//   - mac: MAC bytes (optional for AEAD ciphers)
//   - tag: Authentication tag bytes (optional for traditional ciphers)
//
// Returns:
//   - *EncryptedPayload: Properly formatted payload with Base64-encoded components
func NewEncryptedPayload(iv, encrypted, mac, tag []byte) *EncryptedPayload {
	payload := &EncryptedPayload{
		IV:    encodeBase64(iv),
		Value: encodeBase64(encrypted),
	}

	if len(mac) > 0 {
		payload.MAC = encodeBase64(mac)
	}

	if len(tag) > 0 {
		payload.Tag = encodeBase64(tag)
	}

	return payload
}

// ParseEncryptedPayload parses a Base64-encoded JSON payload string.
// Decodes Laravel-compatible encrypted payload format.
//
// Parameters:
//   - payloadStr: Base64-encoded JSON payload string
//
// Returns:
//   - *EncryptedPayload: Parsed payload structure
//   - error: Parsing error if payload format is invalid
func ParseEncryptedPayload(payloadStr string) (*EncryptedPayload, error) {
	// Implementation would decode Base64 and parse JSON
	// This is a placeholder for the actual implementation
	return nil, fmt.Errorf("ParseEncryptedPayload not yet implemented")
}

// ToBase64JSON converts the payload to a Base64-encoded JSON string.
// Creates Laravel-compatible payload format.
//
// Returns:
//   - string: Base64-encoded JSON representation
//   - error: Encoding error if serialization fails
func (p *EncryptedPayload) ToBase64JSON() (string, error) {
	// Implementation would serialize to JSON and encode as Base64
	// This is a placeholder for the actual implementation
	return "", fmt.Errorf("ToBase64JSON not yet implemented")
}

// IsValid checks if the payload has required components.
// Validates payload structure for encryption/decryption operations.
//
// Returns:
//   - bool: true if payload has required IV and Value, false otherwise
func (p *EncryptedPayload) IsValid() bool {
	return p.IV != "" && p.Value != ""
}

// HasMac checks if the payload has a MAC component.
// Used to determine if MAC validation is required.
//
// Returns:
//   - bool: true if MAC is present, false otherwise
func (p *EncryptedPayload) HasMac() bool {
	return p.MAC != ""
}

// HasTag checks if the payload has an authentication tag.
// Used to determine if AEAD tag validation is required.
//
// Returns:
//   - bool: true if tag is present, false otherwise
func (p *EncryptedPayload) HasTag() bool {
	return p.Tag != ""
}

// GetIVBytes decodes and returns the IV as bytes.
// Converts Base64-encoded IV to binary format.
//
// Returns:
//   - []byte: Decoded IV bytes
//   - error: Decoding error if IV is malformed
func (p *EncryptedPayload) GetIVBytes() ([]byte, error) {
	return decodeBase64(p.IV)
}

// GetValueBytes decodes and returns the encrypted value as bytes.
// Converts Base64-encoded value to binary format.
//
// Returns:
//   - []byte: Decoded encrypted value bytes
//   - error: Decoding error if value is malformed
func (p *EncryptedPayload) GetValueBytes() ([]byte, error) {
	return decodeBase64(p.Value)
}

// GetMacBytes decodes and returns the MAC as bytes.
// Converts Base64-encoded MAC to binary format.
//
// Returns:
//   - []byte: Decoded MAC bytes
//   - error: Decoding error if MAC is malformed
func (p *EncryptedPayload) GetMacBytes() ([]byte, error) {
	return decodeBase64(p.MAC)
}

// GetTagBytes decodes and returns the authentication tag as bytes.
// Converts Base64-encoded tag to binary format.
//
// Returns:
//   - []byte: Decoded tag bytes
//   - error: Decoding error if tag is malformed
func (p *EncryptedPayload) GetTagBytes() ([]byte, error) {
	return decodeBase64(p.Tag)
}

// CipherInfo contains comprehensive information about a cipher algorithm.
// This structure provides metadata about cipher capabilities, parameters,
// and security characteristics used throughout the encryption system.
//
// Key features:
//   - Cipher algorithm identification and parameters
//   - Key length, IV length, and mode information
//   - AEAD capability detection and MAC requirements
//   - Security level and cryptographic properties
//   - Options for cipher-specific parameters
type CipherInfo struct {
	// Cipher is the cipher algorithm identifier (e.g., "aes-256-gcm")
	Cipher string `json:"cipher"`

	// KeyLength is the required key length in bytes
	KeyLength int `json:"key_length"`

	// IVLength is the required initialization vector length in bytes
	IVLength int `json:"iv_length"`

	// Mode is the cipher mode of operation (e.g., "gcm", "cbc", "ctr")
	Mode string `json:"mode"`

	// IsAEAD indicates if the cipher provides authenticated encryption
	IsAEAD bool `json:"is_aead"`

	// RequiresMAC indicates if the cipher requires separate MAC validation
	RequiresMAC bool `json:"requires_mac"`

	// Options contains cipher-specific parameters and settings
	Options map[string]interface{} `json:"options"`
}

// HashInfo contains comprehensive information about a hash algorithm.
// This structure provides metadata about hash algorithm capabilities,
// parameters, and security characteristics used throughout the hashing system.
//
// Key features:
//   - Algorithm identification and parameters
//   - Cost, memory, time, and parallelism settings
//   - Algorithm-specific option support
//   - Security parameter tracking
//   - Hash format and version information
type HashInfo struct {
	// Algo is the hash algorithm identifier (e.g., "argon2id", "bcrypt")
	Algo string `json:"algo"`

	// AlgoName is the human-readable algorithm name
	AlgoName string `json:"algo_name"`

	// Options contains algorithm-specific parameters (cost, memory, time, etc.)
	Options map[string]interface{} `json:"options"`
}

// NewCipherInfo creates a new CipherInfo instance for the specified cipher.
// Automatically populates cipher metadata based on algorithm characteristics.
//
// This function analyzes the cipher string and sets appropriate values for:
//   - Key length requirements
//   - IV length requirements
//   - Cipher mode identification
//   - AEAD capability detection
//   - MAC requirement determination
//
// Parameters:
//   - cipher: The cipher algorithm identifier
//
// Returns:
//   - CipherInfo: Populated cipher information structure
func NewCipherInfo(cipher enums.Cipher) CipherInfo {
	info := CipherInfo{
		Cipher:  cipher.String(),
		Options: make(map[string]interface{}),
	}

	// Parse cipher characteristics based on algorithm identifier
	switch cipher {
	case "aes-128-gcm", "AES-128-GCM":
		info.KeyLength = 16 // 128 bits
		info.IVLength = 12  // 96 bits for GCM
		info.Mode = "gcm"
		info.IsAEAD = true
		info.RequiresMAC = false

	case "aes-256-gcm", "AES-256-GCM":
		info.KeyLength = 32 // 256 bits
		info.IVLength = 12  // 96 bits for GCM
		info.Mode = "gcm"
		info.IsAEAD = true
		info.RequiresMAC = false

	case "aes-128-cbc", "AES-128-CBC":
		info.KeyLength = 16 // 128 bits
		info.IVLength = 16  // 128 bits for CBC
		info.Mode = "cbc"
		info.IsAEAD = false
		info.RequiresMAC = true

	case "aes-256-cbc", "AES-256-CBC":
		info.KeyLength = 32 // 256 bits
		info.IVLength = 16  // 128 bits for CBC
		info.Mode = "cbc"
		info.IsAEAD = false
		info.RequiresMAC = true

	case "aes-128-ctr", "AES-128-CTR":
		info.KeyLength = 16 // 128 bits
		info.IVLength = 16  // 128 bits for CTR
		info.Mode = "ctr"
		info.IsAEAD = false
		info.RequiresMAC = true

	case "aes-256-ctr", "AES-256-CTR":
		info.KeyLength = 32 // 256 bits
		info.IVLength = 16  // 128 bits for CTR
		info.Mode = "ctr"
		info.IsAEAD = false
		info.RequiresMAC = true

	default:
		// Unknown cipher - return empty info
		info.KeyLength = 0
		info.IVLength = 0
		info.Mode = ""
		info.IsAEAD = false
		info.RequiresMAC = true // Conservative default
	}

	return info
}

// encodeBase64 encodes bytes to Base64 string.
// Helper function for payload encoding.
//
// Parameters:
//   - data: Bytes to encode
//
// Returns:
//   - string: Base64-encoded string
func encodeBase64(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

// decodeBase64 decodes Base64 string to bytes.
// Helper function for payload decoding.
//
// Parameters:
//   - encoded: Base64-encoded string
//
// Returns:
//   - []byte: Decoded bytes
//   - error: Decoding error if string is malformed
func decodeBase64(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(encoded)
}

// NewHashInfo creates a new HashInfo instance for the specified algorithm.
// Automatically populates algorithm metadata based on hash characteristics.
//
// Parameters:
//   - algo: The hash algorithm identifier
//   - options: Algorithm-specific options
//
// Returns:
//   - HashInfo: Populated hash information structure
func NewHashInfo(algo string, options map[string]interface{}) HashInfo {
	if options == nil {
		options = make(map[string]interface{})
	}

	info := HashInfo{
		Algo:    algo,
		Options: options,
	}

	// Set human-readable algorithm names
	switch algo {
	case "argon2i":
		info.AlgoName = "Argon2i"
	case "argon2id":
		info.AlgoName = "Argon2id"
	case "bcrypt":
		info.AlgoName = "bcrypt"
	default:
		info.AlgoName = algo
	}

	return info
}
