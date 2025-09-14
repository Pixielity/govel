package exceptions

import "errors"

// ErrInvalidAlgorithm indicates that an unknown or malformed algorithm identifier was provided.
// Distinguishes between unavailable algorithms and completely unrecognized identifiers.
//
// Common causes:
//   - Typos in algorithm names ("bcyrpt" instead of "bcrypt")
//   - Deprecated or obsolete algorithm identifiers
//   - Malformed hash prefixes or empty algorithm names
//   - Invalid configuration file specifications
//
// Valid algorithms: "bcrypt", "argon2i", "argon2id".
var ErrInvalidAlgorithm = errors.New("invalid hashing algorithm")