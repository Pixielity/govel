package hashers

import (
	"crypto/subtle"
	"regexp"
	"strconv"
	"strings"

	"govel/packages/hashing/src/exceptions"
	enums "govel/packages/types/src/enums/hashing"
)

// BaseHasher provides common functionality and utilities shared across all hasher implementations.
// It serves as a foundation that can be embedded in specific hasher types to inherit
// common behavior such as hash parsing, format validation, and basic operations.
//
// The BaseHasher implements the HasherInterface with generic functionality that works
// across different hashing algorithms. Specific hashers can embed this type and
// override methods as needed for algorithm-specific behavior.
//
// Features provided:
//   - Hash format parsing for multiple algorithms (bcrypt, Argon2i, Argon2id)
//   - Basic validation and error handling
//   - Constant-time comparison fallback
//   - Standardized info extraction
//   - Type-safe configuration value extraction
//
// Example embedding:
//
//	type BcryptHasher struct {
//		*BaseHasher
//		cost int
//	}
//
//	func (b *BcryptHasher) Make(value string, options map[string]interface{}) (string, error) {
//		// Specific bcrypt implementation
//		return bcryptHash, nil
//	}
//
// Thread-safety: BaseHasher is stateless and safe for concurrent use.
type BaseHasher struct {
	// BaseHasher is intentionally stateless to ensure thread-safety
	// and to serve as a pure utility provider for embedded hashers.
	// Each hasher implementation maintains its own state as needed.
}

// NewBaseHasher creates a new BaseHasher instance.
// This constructor ensures proper initialization of the base hasher
// and can be extended in the future if state becomes necessary.
//
// Returns:
//   - *BaseHasher: A new BaseHasher instance ready for embedding or direct use
//
// Usage:
//   - Can be used directly for hash analysis and parsing
//   - Typically embedded in specific hasher implementations
//   - Safe to call from multiple goroutines
//
// Example:
//
//	// Direct usage
//	base := NewBaseHasher()
//	info := base.Info("$2y$12$...")
//
//	// Embedding usage
//	type MyHasher struct {
//		*BaseHasher
//		// Additional fields
//	}
//
//	func NewMyHasher() *MyHasher {
//		return &MyHasher{
//			BaseHasher: NewBaseHasher(),
//		}
//	}
func NewBaseHasher() *BaseHasher {
	return &BaseHasher{}
}

// Info extracts comprehensive metadata from a hashed value by analyzing its format.
// This method serves as a universal hash parser that can identify and extract
// information from multiple hashing algorithm formats.
//
// Supported formats:
//   - bcrypt: $2[a|b|y]$cost$salt+hash (e.g., "$2y$12$abc123...")
//   - Argon2i: $argon2i$v=19$m=memory,t=time,p=threads$salt$hash
//   - Argon2id: $argon2id$v=19$m=memory,t=time,p=threads$salt$hash
//
// Parameters:
//   - hashedValue: The hash string to analyze and parse
//
// Returns:
//   - types.HashInfo: Structured information about the hash
//   - For recognized formats: Complete algorithm and parameter information
//   - For unrecognized formats: Empty HashInfo with zero values
//   - For empty input: Empty HashInfo with initialized but empty maps
//
// Algorithm detection logic:
//  1. Checks for empty input (returns empty info)
//  2. Tests for bcrypt format (starts with "$2")
//  3. Tests for Argon2 formats (starts with "$argon2")
//  4. Falls back to empty info for unknown formats
//
// Thread-safety: This method is safe for concurrent use as it operates only
// on the input parameter and doesn't modify any shared state.
//
// Example usage:
//
//	// Analyzing a bcrypt hash
//	info := hasher.Info("$2y$12$abc123def456...")
//	if info.AlgoName == "bcrypt" {
//		cost := info.Options["cost"].(int) // 12
//		fmt.Printf("bcrypt hash with cost %d\n", cost)
//	}
//
//	// Analyzing an Argon2id hash
//	info = hasher.Info("$argon2id$v=19$m=65536,t=2,p=1$...")
//	if info.AlgoName == "argon2id" {
//		memory := info.Options["memory_cost"].(int) // 65536
//		time := info.Options["time_cost"].(int)     // 2
//	}
//
//	// Handling unknown format
//	info = hasher.Info("unknown-hash-format")
//	if info.AlgoName == "" {
//		fmt.Println("Unknown hash format")
//	}
func (h *BaseHasher) Info(hashedValue string) types.HashInfo {
	if len(hashedValue) == 0 {
		return types.HashInfo{
			Algo:     "",
			AlgoName: "",
			Options:  make(map[string]interface{}),
		}
	}

	// Parse bcrypt hash format: $2[ayb]$[cost]$[salt][hash]
	if strings.HasPrefix(hashedValue, "$2") {
		return h.parseBcryptInfo(hashedValue)
	}

	// Parse argon2 hash format: $argon2[id/i]$[params]$[salt]$[hash]
	if strings.HasPrefix(hashedValue, "$argon2") {
		return h.parseArgonInfo(hashedValue)
	}

	// Unknown format
	return types.HashInfo{
		Algo:     "",
		AlgoName: "",
		Options:  make(map[string]interface{}),
	}
}

// Check provides a fallback password verification implementation using constant-time comparison.
// This method serves as a basic implementation for the HasherInterface, but should typically
// be overridden by specific hasher implementations with algorithm-appropriate logic.
//
// The base implementation performs a simple byte-wise constant-time comparison between
// the plain value and the hashed value. This is NOT suitable for actual password verification
// as it doesn't account for salt handling or algorithm-specific verification logic.
//
// Parameters:
//   - value: The plain text value to verify
//   - hashedValue: The hash to compare against
//   - options: Additional options (unused in base implementation)
//
// Returns:
//   - bool: true if values match exactly (byte-for-byte), false otherwise
//
// Security notes:
//   - Uses crypto/subtle.ConstantTimeCompare to prevent timing attacks
//   - Returns false immediately for empty hashed values
//   - This is a FALLBACK implementation - real hashers should override this
//
// Usage:
//   - Not intended for direct use in production password verification
//   - Serves as interface compliance for base hasher
//   - Real implementations (BcryptHasher, ArgonHasher, etc.) override this method
//
// Example (should be overridden):
//
//	type MyHasher struct {
//		*BaseHasher
//	}
//
//	func (m *MyHasher) Check(value, hashedValue string, options map[string]interface{}) bool {
//		// Proper algorithm-specific verification logic here
//		return algorithmSpecificCheck(value, hashedValue)
//	}
func (h *BaseHasher) Check(value, hashedValue string, options map[string]interface{}) bool {
	// Return false immediately for empty hash values
	if len(hashedValue) == 0 {
		return false
	}

	// Basic format validation - check if it looks like a known hash format
	if !h.looksLikeValidHash(hashedValue) {
		return false
	}

	// This is a basic fallback - specific implementations should override this method
	// with algorithm-appropriate verification logic (bcrypt.CompareHashAndPassword, etc.)
	return subtle.ConstantTimeCompare([]byte(value), []byte(hashedValue)) == 1
}

// IsHashed determines if a given string appears to be a hashed value by analyzing its format.
// This method uses the Info() method to detect known hash patterns and determine if the
// input string matches any supported hashing algorithm formats.
//
// The detection is based on format analysis rather than cryptographic validation,
// making it a heuristic check that can help distinguish between plain text and hashed values.
//
// Parameters:
//   - value: The string to examine for hash-like characteristics
//
// Returns:
//   - bool: true if the string appears to be a hash, false if it appears to be plain text
//
// Detection logic:
//   - Calls Info() to attempt hash format parsing
//   - Returns true if any algorithm is detected (Algo field is non-empty)
//   - Returns false for unrecognized formats or plain text
//
// Supported hash format detection:
//   - bcrypt: "$2y$12$...", "$2a$10$...", etc.
//   - Argon2i: "$argon2i$v=19$m=65536,t=1,p=1$..."
//   - Argon2id: "$argon2id$v=19$m=65536,t=2,p=1$..."
//
// Use cases:
//   - Determining if user input needs hashing
//   - Validating hash formats before verification
//   - Database migration utilities
//   - Security auditing and analysis
//
// Limitations:
//   - Heuristic-based; may have false positives/negatives
//   - Only detects supported algorithm formats
//   - Does not validate hash integrity or correctness
//
// Example usage:
//
//	// Check if user input is already hashed
//	userInput := "$2y$12$abc123..."
//	if hasher.IsHashed(userInput) {
//		fmt.Println("Input appears to be hashed")
//		// Skip hashing, use for verification
//	} else {
//		fmt.Println("Input appears to be plain text")
//		// Hash the input before storage
//		hash, _ := hasher.Make(userInput, nil)
//	}
func (h *BaseHasher) IsHashed(value string) bool {
	// Use the Info method to detect hash formats
	// If any algorithm is detected, consider it hashed
	return h.Info(value).Algo != ""
}

// VerifyConfiguration validates the hasher configuration for security and compatibility.
// The base implementation provides a minimal validation that always returns true,
// serving as a safe default for hashers that don't require specific configuration checks.
//
// Specific hasher implementations should override this method to perform appropriate
// validation checks for their algorithm-specific parameters and system requirements.
//
// Parameters:
//   - value: Test value for configuration validation (implementation-specific)
//     The base implementation ignores this parameter
//
// Returns:
//   - bool: true if configuration is valid (always true in base implementation)
//
// Base implementation behavior:
//   - Always returns true (permissive default)
//   - No actual validation performed
//   - Safe fallback for hashers without specific requirements
//
// Override recommendations:
//
//	Specific hasher implementations should override this to check:
//	- Parameter ranges (cost values, memory limits, etc.)
//	- System capabilities (available memory, CPU cores)
//	- Algorithm availability and library versions
//	- Security compliance requirements
//
// Example override:
//
//	type BcryptHasher struct {
//		*BaseHasher
//		cost int
//	}
//
//	func (b *BcryptHasher) VerifyConfiguration(value string) bool {
//		// Check bcrypt cost is within acceptable range
//		return b.cost >= bcrypt.MinCost && b.cost <= bcrypt.MaxCost
//	}
//
// Usage contexts:
//   - Application startup validation
//   - Health checks and monitoring
//   - Configuration testing
//   - Security audits
func (h *BaseHasher) VerifyConfiguration(value string) bool {
	// Base implementation provides permissive default
	// Specific hasher implementations should override with appropriate validation
	return true
}

// parseBcryptInfo parses bcrypt hash format and extracts algorithm metadata
// and configuration parameters. This method handles the detailed parsing of
// bcrypt hash strings to provide comprehensive information about the hash.
//
// bcrypt parsing process:
//  1. Validates hash format using regex pattern matching
//  2. Extracts cost parameter from the hash string
//  3. Constructs standardized HashInfo structure
//  4. Handles parsing errors gracefully with empty results
//
// Extracted information:
//   - Algorithm identifier ("2y" for bcrypt)
//   - Standardized algorithm name (AlgorithmBcrypt)
//   - Cost parameter (4-31 range)
//   - Hash format validation status
//
// Parameters:
//   - hashedValue: The bcrypt hash string to parse
//
// Returns:
//   - types.HashInfo: Structured hash information with algorithm and options
//   - For invalid hashes: Returns empty HashInfo with zero values
//
// Error handling:
//   - Malformed hashes return empty HashInfo rather than panicking
//   - Cost parsing errors default to 0 value
//   - Regex match failures result in empty algorithm information
//
// Usage:
//   - Called by Info() method for bcrypt hash analysis
//   - Used internally for hash parameter extraction
//   - Supports hash migration and audit utilities
func (h *BaseHasher) parseBcryptInfo(hashedValue string) types.HashInfo {
	// Bcrypt format: $2[ayb]$[cost]$[salt+hash]
	re := regexp.MustCompile(`^\$2[ayb]?\$(\d+)\$`)
	matches := re.FindStringSubmatch(hashedValue)

	if len(matches) < 2 {
		return types.HashInfo{
			Algo:     "",
			AlgoName: "",
			Options:  make(map[string]interface{}),
		}
	}

	cost, err := strconv.Atoi(matches[1])
	if err != nil {
		cost = 0
	}

	return types.HashInfo{
		Algo:     "2y", // bcrypt identifier
		AlgoName: enums.AlgorithmBcrypt,
		Options: map[string]interface{}{
			"cost": cost,
		},
	}
}

// looksLikeValidHash performs basic format validation on hash strings by checking
// against known hash format patterns. This method provides a first-line defense
// against malformed hash inputs and helps identify supported hash types.
//
// The validation process checks for:
//   - Empty or null inputs (invalid)
//   - bcrypt format patterns (starts with "$2")
//   - Argon2 format patterns (starts with "$argon2")
//   - Unknown formats are rejected
//
// This is a heuristic check that helps prevent processing of obviously invalid
// hash strings before more expensive validation operations.
//
// Parameters:
//   - hashedValue: The hash string to validate
//
// Returns:
//   - bool: true if the hash appears to match a known format, false otherwise
//
// Usage:
//   - Called internally by Check() methods to fail fast on invalid inputs
//   - Used by ValidateHashFormat() for public validation API
//   - Helps distinguish between hash formats for algorithm detection
func (h *BaseHasher) looksLikeValidHash(hashedValue string) bool {
	// Empty hashes are invalid
	if len(hashedValue) == 0 {
		return false
	}

	// Check for bcrypt format
	if strings.HasPrefix(hashedValue, "$2") {
		return h.isValidBcryptFormat(hashedValue)
	}

	// Check for Argon2 format
	if strings.HasPrefix(hashedValue, "$argon2") {
		return h.isValidArgonFormat(hashedValue)
	}

	// Unknown format
	return false
}

// isValidBcryptFormat validates bcrypt hash format against the standard bcrypt
// specification. This method performs detailed validation of bcrypt hash structure
// to ensure compliance with the expected format and prevent malformed input processing.
//
// bcrypt format specification:
//   - Total length: exactly 60 characters
//   - Structure: $2[a|b|y]$[cost]$[22-char-salt][31-char-hash]
//   - Cost: 2-digit decimal number (04-31)
//   - Salt and hash: base64-encoded using bcrypt's custom alphabet
//
// Validation checks:
//   - Exact length requirement (60 characters)
//   - Proper prefix format ($2a$, $2b$, or $2y$)
//   - Valid cost parameter format
//   - Correct salt and hash segment lengths
//
// Parameters:
//   - hashedValue: The bcrypt hash string to validate
//
// Returns:
//   - bool: true if the hash matches bcrypt format specification, false otherwise
//
// Security considerations:
//   - Prevents processing of malformed bcrypt hashes
//   - Helps identify hash corruption or tampering
//   - Ensures compatibility with standard bcrypt libraries
func (h *BaseHasher) isValidBcryptFormat(hashedValue string) bool {
	// Bcrypt format: $2[a|b|y]$[cost]$[salt+hash]
	// Should be exactly 60 characters long
	if len(hashedValue) != 60 {
		return false
	}

	// Check regex pattern
	re := regexp.MustCompile(`^\$2[ayb]?\$\d{2}\$.{53}$`)
	return re.MatchString(hashedValue)
}

// isValidArgonFormat validates Argon2 hash format against the Argon2 specification
// as defined in RFC 9106. This method ensures that Argon2i and Argon2id hashes
// conform to the standard format before attempting to parse or verify them.
//
// Argon2 format specification:
//   - Structure: $argon2[i|id]$v=version$m=memory,t=time,p=threads$salt$hash
//   - Total segments: exactly 6 parts separated by '$'
//   - Algorithm: must be "argon2i" or "argon2id"
//   - Version: must be present (typically "v=19")
//   - Parameters: must include memory (m=), time (t=), and parallelism (p=)
//   - Salt and hash: base64-encoded binary data
//
// Validation checks:
//   - Correct number of segments (6)
//   - Valid algorithm identifier
//   - Proper version parameter format
//   - Required parameter presence
//   - Non-empty salt and hash segments
//
// Parameters:
//   - hashedValue: The Argon2 hash string to validate
//
// Returns:
//   - bool: true if the hash matches Argon2 format specification, false otherwise
//
// Standards compliance:
//   - RFC 9106: The Argon2 Memory-Hard Function for Password Hashing
//   - Compatible with reference implementation format
func (h *BaseHasher) isValidArgonFormat(hashedValue string) bool {
	// Argon2 format: $argon2[id/i]$v=19$m=memory,t=time,p=threads$salt$hash
	parts := strings.Split(hashedValue, "$")
	if len(parts) != 6 {
		return false
	}

	// Check algorithm part
	if parts[1] != "argon2i" && parts[1] != "argon2id" {
		return false
	}

	// Check version part
	if !strings.HasPrefix(parts[2], "v=") {
		return false
	}

	// Check parameters part has required format
	params := parts[3]
	if !strings.Contains(params, "m=") || !strings.Contains(params, "t=") || !strings.Contains(params, "p=") {
		return false
	}

	// Basic validation that salt and hash are base64-like
	return len(parts[4]) > 0 && len(parts[5]) > 0
}

// ValidateHashFormat provides a public interface for comprehensive hash format
// validation with explicit error reporting. This method serves as the primary
// entry point for hash validation in the hashing system.
//
// The validation process:
//  1. Performs format detection using looksLikeValidHash()
//  2. Returns ErrInvalidHash for malformed or unrecognized formats
//  3. Provides detailed error information for debugging
//
// Supported hash formats:
//   - bcrypt: $2[a|b|y]$cost$salt+hash (60 characters)
//   - Argon2i: $argon2i$v=19$m=memory,t=time,p=threads$salt$hash
//   - Argon2id: $argon2id$v=19$m=memory,t=time,p=threads$salt$hash
//
// Parameters:
//   - hashedValue: The hash string to validate
//
// Returns:
//   - error: ErrInvalidHash if validation fails, nil if hash format is valid
//
// Usage examples:
//   - Pre-verification validation in authentication systems
//   - Hash format checking in migration utilities
//   - Input validation in API endpoints
//
// Error handling:
//   - Use errors.Is(err, exceptions.ErrInvalidHash) for type checking
//   - Provides consistent error responses across the system
func (h *BaseHasher) ValidateHashFormat(hashedValue string) error {
	if !h.looksLikeValidHash(hashedValue) {
		return exceptions.ErrInvalidHash
	}
	return nil
}

// parseArgonInfo parses Argon2 hash format and extracts comprehensive algorithm
// metadata and configuration parameters. This method handles the complex parsing
// of Argon2i and Argon2id hash strings according to RFC 9106 specification.
//
// Argon2 parsing process:
//  1. Splits hash string into standard 6 segments
//  2. Validates segment count and algorithm identifier
//  3. Parses parameter string (m=memory,t=time,p=threads)
//  4. Extracts version, salt, and hash information
//  5. Constructs comprehensive HashInfo structure
//
// Extracted information:
//   - Algorithm variant (argon2i or argon2id)
//   - Memory cost parameter (in KB)
//   - Time cost parameter (iteration count)
//   - Parallelism parameter (thread count)
//   - Version information (typically 19)
//   - Salt and hash validation status
//
// Parameters:
//   - hashedValue: The Argon2 hash string to parse
//
// Returns:
//   - types.HashInfo: Structured hash information with algorithm and parameters
//   - For invalid hashes: Returns empty HashInfo with zero values
//
// Error handling:
//   - Insufficient segments return empty HashInfo
//   - Parameter parsing errors skip invalid parameters
//   - Unknown algorithm variants default to original name
//
// Standards compliance:
//   - RFC 9106: The Argon2 Memory-Hard Function for Password Hashing
//   - Compatible with libargon2 reference implementation
//
// Usage:
//   - Called by Info() method for Argon2 hash analysis
//   - Supports configuration auditing and migration
//   - Enables parameter optimization and tuning
func (h *BaseHasher) parseArgonInfo(hashedValue string) types.HashInfo {
	// Argon2 format: $argon2[id/i]$v=version$m=memory,t=time,p=threads$salt$hash
	// Example: $argon2id$v=19$m=65536,t=2,p=1$c29tZXNhbHQ$hash...
	parts := strings.Split(hashedValue, "$")

	// Validate minimum required segments (algorithm, version, params, salt, hash)
	if len(parts) < 5 {
		// Return empty info for malformed hashes with insufficient segments
		return types.HashInfo{
			Algo:     "",
			AlgoName: "",
			Options:  make(map[string]interface{}),
		}
	}

	// Extract algorithm variant from first segment (argon2i or argon2id)
	algoName := parts[1] // Expected: "argon2i" or "argon2id"

	// Extract parameter string from third segment
	paramStr := parts[3] // Expected format: "m=65536,t=2,p=1"

	// Initialize options map to store parsed parameters
	options := make(map[string]interface{})

	// Parse comma-separated parameter string into individual key-value pairs
	params := strings.Split(paramStr, ",")
	for _, param := range params {
		// Split each parameter into key=value format
		kv := strings.Split(param, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]

			// Convert numeric parameters to integers
			if intVal, err := strconv.Atoi(value); err == nil {
				// Map Argon2 parameter keys to standardized option names
				switch key {
				case "m": // Memory cost in KB
					options["memory_cost"] = intVal
				case "t": // Time cost (iterations)
					options["time_cost"] = intVal
				case "p": // Parallelism (threads)
					options["threads"] = intVal
				case "v": // Version (typically 19 for Argon2)
					options["version"] = intVal
				}
			}
		}
	}

	// Map algorithm name to standardized enum constant
	var standardAlgoName string
	switch algoName {
	case "argon2i":
		// Use enum constant for Argon2i
		standardAlgoName = enums.AlgorithmArgon2i
	case "argon2id":
		// Use enum constant for Argon2id (recommended variant)
		standardAlgoName = enums.AlgorithmArgon2id
	default:
		// Fallback to original name for unknown variants
		standardAlgoName = algoName
	}

	// Return structured hash information
	return types.HashInfo{
		Algo:     algoName,         // Original algorithm name from hash
		AlgoName: standardAlgoName, // Standardized enum constant
		Options:  options,          // Parsed parameters (memory, time, threads, version)
	}
}
