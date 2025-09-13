package exceptions

import "errors"

// ErrHashingNotSupported indicates that the requested hashing algorithm is not available.
// Occurs when algorithms are missing due to build configuration, dependencies, or system limitations.
//
// Common causes:
//   - Required cryptographic libraries not installed or linked
//   - Algorithm not compiled into binary (build tags, CGO disabled)
//   - System lacks necessary hardware features or OS support
//   - Library versions incompatible or outdated
//
// Use fallback algorithms or install missing dependencies to resolve.
var ErrHashingNotSupported = errors.New("hashing algorithm not supported")