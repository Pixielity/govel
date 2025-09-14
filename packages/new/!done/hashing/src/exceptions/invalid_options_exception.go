package exceptions

import "errors"

// ErrInvalidOptions indicates that invalid or out-of-range configuration parameters were provided.
// Helps catch configuration errors early and provides clear feedback about parameter validation failures.
//
// Common parameter issues:
//   - bcrypt cost outside range (4-31) or memory/time parameters too extreme
//   - Argon2 memory exceeding system limits or thread count beyond CPU cores
//   - Incompatible option combinations or wrong parameter types
//   - System resource constraints violated
//
// Validate parameters before hashing operations to ensure proper configuration.
var ErrInvalidOptions = errors.New("invalid options provided")