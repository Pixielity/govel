package exceptions

import "errors"

// ErrInvalidHash indicates that a malformed or invalid hash string was provided.
// Occurs when hash fails to match expected format patterns or contains invalid characters.
//
// Common causes:
//   - Truncated or incomplete hash strings
//   - Incorrect delimiter usage (missing $ separators)
//   - Invalid base64 encoding in salt/hash segments
//   - Corrupted hash variants or mixed algorithm components
//
// Prevents malformed input exploitation and distinguishes format errors from verification failures.
var ErrInvalidHash = errors.New("invalid hash format")