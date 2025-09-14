package exceptions

import "errors"

// ErrValueTooLong indicates that the input value exceeds the maximum allowed length.
// Prevents denial-of-service attacks and ensures reasonable memory usage during hashing.
//
// Common scenarios:
//   - Password input exceeds implementation-specific limits
//   - Malicious input attempting to consume excessive resources
//   - System resources would be exhausted by oversized input
//
// Algorithm limits: bcrypt (72 bytes), Argon2 (configurable, typically 1KB-4KB).
var ErrValueTooLong = errors.New("value is too long to hash")
