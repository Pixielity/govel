package exceptions

import "errors"

// ErrDecryptionFailed indicates that the decryption process failed.
// Occurs when payload cannot be decrypted due to various cryptographic issues.
//
// Common causes:
//   - Wrong decryption key used
//   - Payload has been tampered with or corrupted
//   - IV or encrypted data is invalid
//   - Padding is incorrect for block cipher modes
//
// Verify the correct key is being used and payload hasn't been modified.
var ErrDecryptionFailed = errors.New("decryption failed")