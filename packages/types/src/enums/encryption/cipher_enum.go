package enums

// Cipher represents the encryption cipher algorithms.
type Cipher string

const (
	// CipherAES128CBC represents AES-128-CBC cipher
	CipherAES128CBC Cipher = "AES-128-CBC"
	
	// CipherAES256CBC represents AES-256-CBC cipher
	CipherAES256CBC Cipher = "AES-256-CBC"
	
	// CipherAES128GCM represents AES-128-GCM cipher
	CipherAES128GCM Cipher = "AES-128-GCM"
	
	// CipherAES256GCM represents AES-256-GCM cipher
	CipherAES256GCM Cipher = "AES-256-GCM"
	
	// CipherAES128CTR represents AES-128-CTR cipher
	CipherAES128CTR Cipher = "AES-128-CTR"
	
	// CipherAES256CTR represents AES-256-CTR cipher
	CipherAES256CTR Cipher = "AES-256-CTR"
)