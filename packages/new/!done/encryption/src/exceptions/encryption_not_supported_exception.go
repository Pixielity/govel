package exceptions

import "errors"

// ErrEncryptionNotSupported indicates that the requested encryption algorithm is not available.
// Occurs when algorithms are missing due to build configuration, dependencies, or system limitations.
//
// Common causes:
//   - Required cryptographic libraries not installed or linked
//   - Algorithm not compiled into binary (build tags, CGO disabled)
//   - System lacks necessary hardware features or OS support
//   - Library versions incompatible or outdated
//
// Use fallback algorithms or install missing dependencies to resolve.
var ErrEncryptionNotSupported = errors.New("encryption algorithm not supported")