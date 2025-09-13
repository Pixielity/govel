package exceptions

import "errors"

// ErrInvalidCipher indicates that an unknown or malformed cipher identifier was provided.
// Distinguishes between unavailable ciphers and completely unrecognized identifiers.
//
// Common causes:
//   - Typos in cipher names ("AES-265-CBC" instead of "AES-256-CBC")
//   - Deprecated or obsolete cipher identifiers
//   - Malformed cipher specifications or empty cipher names
//   - Invalid configuration file specifications
//
// Valid ciphers: "AES-256-CBC", "AES-256-GCM", "AES-256-CTR", "AES-128-CBC", "AES-128-GCM".
var ErrInvalidCipher = errors.New("invalid encryption cipher")