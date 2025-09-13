package enums

// Algorithm represents the hashing algorithms.
type Algorithm string

const (
	// AlgorithmBcrypt represents bcrypt algorithm
	AlgorithmBcrypt Algorithm = "bcrypt"
	
	// AlgorithmArgon2i represents Argon2i algorithm
	AlgorithmArgon2i Algorithm = "argon2i"
	
	// AlgorithmArgon2id represents Argon2id algorithm
	AlgorithmArgon2id Algorithm = "argon2id"
	
	// AlgorithmSHA256 represents SHA-256 algorithm
	AlgorithmSHA256 Algorithm = "sha256"
	
	// AlgorithmSHA512 represents SHA-512 algorithm
	AlgorithmSHA512 Algorithm = "sha512"
)