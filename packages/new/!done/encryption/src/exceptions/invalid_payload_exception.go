package exceptions

import "errors"

// ErrInvalidPayload indicates that the encrypted payload structure is invalid or corrupted.
// Occurs when payload format, encoding, or required fields are missing or malformed.
//
// Common causes:
//   - Payload is not valid base64-encoded JSON
//   - Required fields (iv, value, mac) are missing
//   - Payload structure doesn't match expected format
//   - Payload has been truncated or corrupted during transmission
//
// Ensure payload follows the expected encrypted format with all required fields.
var ErrInvalidPayload = errors.New("invalid encrypted payload")