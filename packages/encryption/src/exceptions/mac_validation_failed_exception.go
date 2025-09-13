package exceptions

import "errors"

// ErrMacValidationFailed indicates that MAC (Message Authentication Code) validation failed.
// This is a security-critical error suggesting payload tampering or key mismatch.
//
// Common causes:
//   - Payload has been tampered with or modified
//   - Wrong MAC key used for validation
//   - MAC computation algorithm mismatch
//   - Payload corruption during transmission or storage
//
// This error should be treated as a security issue - the payload should not be trusted.
var ErrMacValidationFailed = errors.New("MAC validation failed")