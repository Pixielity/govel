package encryption

import (
	"govel/packages/encryption/src/encrypters"
	enums "govel/packages/types/src/enums/encryption"
	containerInterfaces "govel/packages/types/src/interfaces/container"
	encrypterInterfaces "govel/packages/types/src/interfaces/encryption"
	support "govel/packages/support/src"
)

// EncryptionManager provides centralized encryption management with multiple algorithm support.
// Extends support.Manager implementing Laravel-style driver pattern for encryption operations.
//
// Key features:
//   - Multi-algorithm support (AES-CBC, AES-GCM, AES-CTR)
//   - Laravel-style driver creation and management
//   - Container-based dependency injection
//   - Automatic algorithm validation and fallback
//
// Implements both EncrypterInterface and FactoryInterface for complete encryption functionality.
type EncryptionManager struct {
	*support.Manager
}

// NewEncryptionManager creates a new encryption manager with Laravel-style functionality.
// Initializes the manager with container dependency injection and proxy self-reference.
//
// Constructor features:
//   - Container-based dependency injection for configuration
//   - Proxy self-reference for proper method resolution
//   - Laravel-style driver pattern implementation
//   - Automatic driver discovery via reflection
//
// Returns configured EncryptionManager ready for encryption operations.
func NewEncryptionManager(container containerInterfaces.ContainerInterface) *EncryptionManager {
	// Create base manager with dependency injection container
	// This provides core driver management functionality
	baseManager := support.NewManager(container)

	// Initialize encryption manager wrapping the base manager
	// Inherits all driver management capabilities
	encryptionManager := &EncryptionManager{
		Manager: baseManager,
	}

	// Set up proxy self-reference for proper method resolution
	// This enables the Manager's reflection to find EncryptionManager's CreateXXXDriver methods
	// Critical for Laravel-style automatic driver discovery
	baseManager.SetProxySelf(encryptionManager)

	return encryptionManager
}

// Encrypter gets an encrypter instance by cipher name (optional).
// If no name is provided, uses the default driver.
// Provides type-safe access to specific encryption algorithms with validation.
//
// Method features:
//   - Optional cipher name (uses default if not provided)
//   - Cipher validation before driver creation (only if name provided)
//   - Automatic driver instantiation and caching
//   - Type-safe EncrypterInterface casting
//   - Graceful error handling with nil returns
//
// Returns nil if cipher is invalid or driver creation fails.
func (e *EncryptionManager) Encrypter(name ...string) encrypterInterfaces.EncrypterInterface {
	var driverName string

	// Determine which driver to use
	if len(name) > 0 && name[0] != "" {
		// Name provided - validate it
		driverName = name[0]
		if err := enums.ValidateCipher(driverName); err != nil {
			return nil // Invalid cipher returns nil for graceful handling
		}
	}

	// Get driver instance using base manager's driver resolution
	// This triggers CreateXXXDriver methods via reflection
	driver, err := e.Driver(driverName)
	if err != nil {
		// Driver creation failed - return nil for consistent error handling
		return nil
	}

	// Type-safe casting to EncrypterInterface
	// Ensures returned driver implements required encryption operations
	encrypter, ok := driver.(encrypterInterfaces.EncrypterInterface)
	if !ok {
		// Driver doesn't implement EncrypterInterface - should never happen
		return nil
	}

	return encrypter
}

// Driver creation methods - Laravel reflection-style method naming
// These methods are automatically discovered by the base Manager via reflection.
// Method naming convention: Create{CipherName}Driver() for automatic discovery.
// Each method handles configuration loading and encrypter instantiation for its cipher.

// CreateAES256CBCDriver creates AES-256-CBC driver instance.
// Automatically discovered by base Manager for "AES-256-CBC" cipher requests.
//
// Extracts key and cipher from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES256CBCDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-256-cbc"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key and cipher from config or use defaults
	var key []byte
	cipher := enums.CipherAES256CBC

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESCBCEncrypter(key, cipher, config)
}

// CreateAES256GCMDriver creates AES-256-GCM driver instance.
// Automatically discovered by base Manager for "AES-256-GCM" cipher requests.
//
// Extracts key and cipher from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES256GCMDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-256-gcm"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key and cipher from config or use defaults
	var key []byte
	cipher := enums.CipherAES256GCM

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESGCMEncrypter(key, cipher, config)
}

// CreateAES256CTRDriver creates AES-256-CTR driver instance.
// Automatically discovered by base Manager for "AES-256-CTR" cipher requests.
//
// Extracts key from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES256CTRDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-256-ctr"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key from config or use defaults
	var key []byte
	cipher := enums.CipherAES256CTR

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESCTREncrypter(key, cipher)
}

// CreateAES128CBCDriver creates AES-128-CBC driver instance.
// Automatically discovered by base Manager for "AES-128-CBC" cipher requests.
//
// Extracts key and cipher from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES128CBCDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-128-cbc"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key and cipher from config or use defaults
	var key []byte
	cipher := enums.CipherAES128CBC

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESCBCEncrypter(key, cipher, config)
}

// CreateAES128GCMDriver creates AES-128-GCM driver instance.
// Automatically discovered by base Manager for "AES-128-GCM" cipher requests.
//
// Extracts key and cipher from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES128GCMDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-128-gcm"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key and cipher from config or use defaults
	var key []byte
	cipher := enums.CipherAES128GCM

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESGCMEncrypter(key, cipher, config)
}

// CreateAES128CTRDriver creates AES-128-CTR driver instance.
// Automatically discovered by base Manager for "AES-128-CTR" cipher requests.
//
// Extracts key from configuration and creates encrypter instance.
func (e *EncryptionManager) CreateAES128CTRDriver() (interface{}, error) {
	// Get configuration from container
	var config map[string]interface{}
	if configValue, exists := e.GetConfig().Get("encryption.aes-128-ctr"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}

	// Extract key from config or use defaults
	var key []byte
	cipher := enums.CipherAES128CTR

	// Get key from config or global config
	if keyValue, exists := config["key"]; exists {
		if keyBytes, ok := keyValue.([]byte); ok {
			key = keyBytes
		} else if keyStr, ok := keyValue.(string); ok {
			key = []byte(keyStr)
		}
	}

	// Fallback to global app key if not found
	if len(key) == 0 {
		if appKeyValue, exists := e.GetConfig().Get("app.key"); exists {
			if keyStr, ok := appKeyValue.(string); ok {
				key = []byte(keyStr)
			}
		}
	}

	return encrypters.NewAESCTREncrypter(key, cipher)
}

// These methods provide a unified interface for encryption operations by delegating
// to the current default driver. This allows the manager to act as both a
// factory (via Encrypter method) and a direct encrypter (via these delegate methods).

// GetDefaultDriver implements the ManagerInterface requirement.
// Returns the default cipher identifier for driver resolution and delegation.
//
// Default driver features:
//   - Centralized default cipher configuration
//   - Consistent behavior across all manager operations
//   - Easy configuration override via enum constants
//   - Laravel-style manager pattern compliance
//
// Returns the default cipher name from enums package.
func (e *EncryptionManager) GetDefaultDriver() string {
	// Delegate to enums package for centralized default management
	// This ensures consistency across the entire encryption system
	return enums.GetDefaultCipher()
}

// Compile-time interface compliance checks
// These ensure EncryptionManager properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ encrypterInterfaces.FactoryInterface = (*EncryptionManager)(nil) // Encrypter factory operations
