package types

// CipherInfo contains comprehensive information about a cipher algorithm.
// This structure provides metadata about cipher capabilities, parameters,
// and security characteristics used throughout the encryption system.
//
// Key features:
//   - Cipher algorithm identification and parameters
//   - Key length, IV length, and mode information
//   - AEAD capability detection and MAC requirements
//   - Security level and cryptographic properties
//   - Options for cipher-specific parameters
type CipherInfo struct {
	// Cipher is the cipher algorithm identifier (e.g., "aes-256-gcm")
	Cipher string `json:"cipher"`
	
	// KeyLength is the required key length in bytes
	KeyLength int `json:"key_length"`
	
	// IVLength is the required initialization vector length in bytes
	IVLength int `json:"iv_length"`
	
	// Mode is the cipher mode of operation (e.g., "gcm", "cbc", "ctr")
	Mode string `json:"mode"`
	
	// IsAEAD indicates if the cipher provides authenticated encryption
	IsAEAD bool `json:"is_aead"`
	
	// RequiresMAC indicates if the cipher requires separate MAC validation
	RequiresMAC bool `json:"requires_mac"`
	
	// Options contains cipher-specific parameters and settings
	Options map[string]interface{} `json:"options"`
}

// HashInfo contains comprehensive information about a hash algorithm.
// This structure provides metadata about hash algorithm capabilities,
// parameters, and security characteristics used throughout the hashing system.
//
// Key features:
//   - Algorithm identification and parameters
//   - Cost, memory, time, and parallelism settings
//   - Algorithm-specific option support
//   - Security parameter tracking
//   - Hash format and version information
type HashInfo struct {
	// Algo is the hash algorithm identifier (e.g., "argon2id", "bcrypt")
	Algo string `json:"algo"`
	
	// AlgoName is the human-readable algorithm name
	AlgoName string `json:"algo_name"`
	
	// Options contains algorithm-specific parameters (cost, memory, time, etc.)
	Options map[string]interface{} `json:"options"`
}

// NewCipherInfo creates a new CipherInfo instance for the specified cipher.
// Automatically populates cipher metadata based on algorithm characteristics.
//
// This function analyzes the cipher string and sets appropriate values for:
//   - Key length requirements
//   - IV length requirements  
//   - Cipher mode identification
//   - AEAD capability detection
//   - MAC requirement determination
//
// Parameters:
//   - cipher: The cipher algorithm identifier
//
// Returns:
//   - CipherInfo: Populated cipher information structure
func NewCipherInfo(cipher string) CipherInfo {
	info := CipherInfo{
		Cipher:  cipher,
		Options: make(map[string]interface{}),
	}
	
	// Parse cipher characteristics based on algorithm identifier
	switch cipher {
	case "aes-128-gcm", "AES-128-GCM":
		info.KeyLength = 16   // 128 bits
		info.IVLength = 12    // 96 bits for GCM
		info.Mode = "gcm"
		info.IsAEAD = true
		info.RequiresMAC = false
		
	case "aes-256-gcm", "AES-256-GCM":
		info.KeyLength = 32   // 256 bits
		info.IVLength = 12    // 96 bits for GCM
		info.Mode = "gcm"
		info.IsAEAD = true
		info.RequiresMAC = false
		
	case "aes-128-cbc", "AES-128-CBC":
		info.KeyLength = 16   // 128 bits
		info.IVLength = 16    // 128 bits for CBC
		info.Mode = "cbc"
		info.IsAEAD = false
		info.RequiresMAC = true
		
	case "aes-256-cbc", "AES-256-CBC":
		info.KeyLength = 32   // 256 bits
		info.IVLength = 16    // 128 bits for CBC
		info.Mode = "cbc"
		info.IsAEAD = false
		info.RequiresMAC = true
		
	case "aes-128-ctr", "AES-128-CTR":
		info.KeyLength = 16   // 128 bits
		info.IVLength = 16    // 128 bits for CTR
		info.Mode = "ctr"
		info.IsAEAD = false
		info.RequiresMAC = true
		
	case "aes-256-ctr", "AES-256-CTR":
		info.KeyLength = 32   // 256 bits
		info.IVLength = 16    // 128 bits for CTR
		info.Mode = "ctr"
		info.IsAEAD = false
		info.RequiresMAC = true
		
	default:
		// Unknown cipher - return empty info
		info.KeyLength = 0
		info.IVLength = 0
		info.Mode = ""
		info.IsAEAD = false
		info.RequiresMAC = true // Conservative default
	}
	
	return info
}

// NewHashInfo creates a new HashInfo instance for the specified algorithm.
// Automatically populates algorithm metadata based on hash characteristics.
//
// Parameters:
//   - algo: The hash algorithm identifier
//   - options: Algorithm-specific options
//
// Returns:
//   - HashInfo: Populated hash information structure
func NewHashInfo(algo string, options map[string]interface{}) HashInfo {
	if options == nil {
		options = make(map[string]interface{})
	}
	
	info := HashInfo{
		Algo:    algo,
		Options: options,
	}
	
	// Set human-readable algorithm names
	switch algo {
	case "argon2i":
		info.AlgoName = "Argon2i"
	case "argon2id":
		info.AlgoName = "Argon2id"
	case "bcrypt":
		info.AlgoName = "bcrypt"
	default:
		info.AlgoName = algo
	}
	
	return info
}