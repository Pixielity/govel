package exceptions

import "errors"

// ErrInvalidKey indicates that the provided encryption key is invalid or inappropriate.
// Occurs when keys have incorrect length, format, or cryptographic properties.
//
// Common causes:
//   - Key length doesn't match cipher requirements (e.g., 16 bytes for AES-128, 32 for AES-256)
//   - Key contains invalid characters or encoding
//   - Key is all zeros or has insufficient entropy
//   - Key format doesn't match expected structure
//
// Ensure key meets cipher-specific length and entropy requirements.
var ErrInvalidKey = errors.New("invalid encryption key")