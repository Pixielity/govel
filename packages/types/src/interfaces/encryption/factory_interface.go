package interfaces

import enums "govel/types/enums/encryption"

// FactoryInterface defines the contract for encryption factory functionality.
// This interface provides encrypter instance creation and management capabilities.
type FactoryInterface interface {
	// Encrypter gets an encrypter instance by cipher name (optional).
	// If no name is provided, uses the default driver.
	Encrypter(name ...enums.Cipher) EncrypterInterface
}
