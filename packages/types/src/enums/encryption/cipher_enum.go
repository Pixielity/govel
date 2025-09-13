package enums

import (
	"errors"
	"strings"
)

// Cipher represents the encryption cipher algorithms.
type Cipher string

const (
	// CipherAES128CBC represents AES-128-CBC cipher
	CipherAES128CBC Cipher = "AES-128-CBC"

	// CipherAES256CBC represents AES-256-CBC cipher
	CipherAES256CBC Cipher = "AES-256-CBC"

	// CipherAES128GCM represents AES-128-GCM cipher
	CipherAES128GCM Cipher = "AES-128-GCM"

	// CipherAES256GCM represents AES-256-GCM cipher
	CipherAES256GCM Cipher = "AES-256-GCM"

	// CipherAES128CTR represents AES-128-CTR cipher
	CipherAES128CTR Cipher = "AES-128-CTR"

	// CipherAES256CTR represents AES-256-CTR cipher
	CipherAES256CTR Cipher = "AES-256-CTR"
)

// String returns the string representation of the Cipher.
func (c Cipher) String() string {
	return string(c)
}

// ValidateCipher validates if the given cipher is supported.
// Returns an error if the cipher is not recognized or supported.
//
// Parameters:
//   - cipher: The cipher algorithm to validate
//
// Returns:
//   - error: nil if valid, error describing the issue if invalid
func ValidateCipher(cipher Cipher) error {
	// Normalize cipher string to uppercase
	normalizedCipher := strings.ToUpper(cipher.String())

	switch normalizedCipher {
	case "AES-128-CBC", "AES-256-CBC", "AES-128-GCM", "AES-256-GCM", "AES-128-CTR", "AES-256-CTR":
		return nil
	default:
		return errors.New("unsupported cipher algorithm")
	}
}

// GetKeyLength returns the required key length in bytes for the given cipher.
// Returns 0 if the cipher is not recognized.
//
// Parameters:
//   - cipher: The cipher algorithm
//
// Returns:
//   - int: Key length in bytes (16 for AES-128, 32 for AES-256, 0 if unknown)
func GetKeyLength(cipher Cipher) int {
	// Normalize cipher string to uppercase
	normalizedCipher := strings.ToUpper(cipher.String())

	switch normalizedCipher {
	case "AES-128-CBC", "AES-128-GCM", "AES-128-CTR":
		return 16 // 128 bits
	case "AES-256-CBC", "AES-256-GCM", "AES-256-CTR":
		return 32 // 256 bits
	default:
		return 0 // Unknown cipher
	}
}

// RequiresMAC determines if the given cipher requires separate MAC validation.
// AEAD ciphers (like GCM) don't require separate MAC as they provide authentication.
//
// Parameters:
//   - cipher: The cipher algorithm
//
// Returns:
//   - bool: true if separate MAC is required, false for AEAD ciphers
func RequiresMAC(cipher Cipher) bool {
	// Normalize cipher string to uppercase
	normalizedCipher := strings.ToUpper(cipher.String())

	switch normalizedCipher {
	case "AES-128-GCM", "AES-256-GCM":
		return false // AEAD ciphers have built-in authentication
	case "AES-128-CBC", "AES-256-CBC", "AES-128-CTR", "AES-256-CTR":
		return true // Traditional ciphers require separate MAC
	default:
		return true // Conservative default - require MAC for unknown ciphers
	}
}

// IsAEAD determines if the given cipher is an Authenticated Encryption with Associated Data cipher.
// AEAD ciphers provide both encryption and authentication in a single operation.
//
// Parameters:
//   - cipher: The cipher algorithm identifier
//
// Returns:
//   - bool: true if the cipher is AEAD, false otherwise
func IsAEAD(cipher string) bool {
	// Normalize cipher string to uppercase
	normalizedCipher := strings.ToUpper(cipher)

	switch normalizedCipher {
	case "AES-128-GCM", "AES-256-GCM":
		return true // GCM modes are AEAD
	default:
		return false // Other modes are not AEAD
	}
}

// GetIVLength returns the required IV length in bytes for the given cipher.
// Different cipher modes require different IV lengths for security.
//
// Parameters:
//   - cipher: The cipher algorithm
//
// Returns:
//   - int: IV length in bytes (12 for GCM, 16 for CBC/CTR, 0 if unknown)
func GetIVLength(cipher Cipher) int {
	// Normalize cipher string to uppercase
	normalizedCipher := strings.ToUpper(cipher.String())

	switch normalizedCipher {
	case "AES-128-GCM", "AES-256-GCM":
		return 12 // 96 bits for GCM mode
	case "AES-128-CBC", "AES-256-CBC", "AES-128-CTR", "AES-256-CTR":
		return 16 // 128 bits for CBC and CTR modes
	default:
		return 0 // Unknown cipher
	}
}

// GetDefaultCipher returns the default cipher algorithm for the encryption system.
// This provides a centralized way to manage the default encryption algorithm.
//
// Returns:
//   - string: The default cipher algorithm identifier
func GetDefaultCipher() string {
	return string(CipherAES256GCM) // Default to AES-256-GCM for best security
}
