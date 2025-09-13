package hashing

import (
	"govel/hashing/hashers"
	support "govel/support"
	enums "govel/types/enums/hashing"
	containerInterfaces "govel/types/interfaces/container"
	hashingInterfaces "govel/types/interfaces/hashing"
)

// HashManager provides centralized hash management with multiple algorithm support.
// Extends support.Manager implementing Laravel-style driver pattern for hash operations.
//
// Key features:
//   - Multi-algorithm support (bcrypt, Argon2i, Argon2id)
//   - Laravel-style driver creation and management
//   - Container-based dependency injection
//   - Automatic algorithm validation and fallback
//
// Implements both HasherInterface and FactoryInterface for complete hash functionality.
type HashManager struct {
	*support.Manager
}

// NewHashManager creates a new hash manager with Laravel-style functionality.
// Initializes the manager with container dependency injection and proxy self-reference.
//
// Constructor features:
//   - Container-based dependency injection for configuration
//   - Proxy self-reference for proper method resolution
//   - Laravel-style driver pattern implementation
//   - Automatic driver discovery via reflection
//
// Returns configured HashManager ready for hash operations.
func NewHashManager(container containerInterfaces.ContainerInterface) *HashManager {
	// Create base manager with dependency injection container
	// This provides core driver management functionality
	baseManager := support.NewManager(container)

	// Initialize hash manager wrapping the base manager
	// Inherits all driver management capabilities
	hashManager := &HashManager{
		Manager: baseManager,
	}

	// Set up proxy self-reference for proper method resolution
	// This enables the Manager's reflection to find HashManager's CreateXXXDriver methods
	// Critical for Laravel-style automatic driver discovery
	baseManager.SetProxySelf(hashManager)

	return hashManager
}

// Hasher gets a hasher instance by algorithm name (optional).
// If no name is provided, uses the default driver.
// Provides type-safe access to specific hashing algorithms with validation.
//
// Method features:
//   - Optional algorithm name (uses default if not provided)
//   - Algorithm validation before driver creation (only if name provided)
//   - Automatic driver instantiation and caching
//   - Type-safe HasherInterface casting
//   - Graceful error handling with nil returns
//
// Returns nil if algorithm is invalid or driver creation fails.
func (h *HashManager) Hasher(name ...enums.Algorithm) hashingInterfaces.HasherInterface {
	var driverName string

	// Determine which driver to use
	if len(name) > 0 && name[0] != "" {
		// Name provided - validate it
		driverName = name[0].String()
		if err := enums.ValidateAlgorithm(name[0]); err != nil {
			return nil // Invalid algorithm returns nil for graceful handling
		}
	}

	// Get driver instance using base manager's driver resolution
	// This triggers CreateXXXDriver methods via reflection
	driver, err := h.Driver(driverName)
	if err != nil {
		// Driver creation failed - return nil for consistent error handling
		return nil
	}

	// Type-safe casting to HasherInterface
	// Ensures returned driver implements required hash operations
	hasher, ok := driver.(hashingInterfaces.HasherInterface)
	if !ok {
		// Driver doesn't implement HasherInterface - should never happen
		return nil
	}

	return hasher
}

// Driver creation methods - Laravel reflection-style method naming
// These methods are automatically discovered by the base Manager via reflection.
// Method naming convention: Create{AlgorithmName}Driver() for automatic discovery.
// Each method handles configuration loading and hasher instantiation for its algorithm.

// CreateBcryptDriver creates bcrypt driver instance.
// Automatically discovered by base Manager for "bcrypt" algorithm requests.
//
// Simply passes configuration from container to hasher constructor.
// All defaults, validation, and parameter handling is done by the hasher itself.
func (h *HashManager) CreateBcryptDriver() (interface{}, error) {
	// Get configuration from container and pass directly to hasher
	var config map[string]interface{}
	if configValue, exists := h.GetConfig().Get("hashing.bcrypt"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}
	return hashers.NewBcryptHasher(config), nil
}

// CreateArgon2iDriver creates Argon2i driver instance.
// Automatically discovered by base Manager for "argon2i" algorithm requests.
//
// Simply passes configuration from container to hasher constructor.
// All defaults, validation, and parameter handling is done by the hasher itself.
func (h *HashManager) CreateArgon2iDriver() (interface{}, error) {
	// Get configuration from container and pass directly to hasher
	var config map[string]interface{}
	if configValue, exists := h.GetConfig().Get("hashing.argon2i"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}
	return hashers.NewArgonHasher(config), nil
}

// CreateArgon2idDriver creates Argon2id driver instance.
// Automatically discovered by base Manager for "argon2id" algorithm requests.
//
// Simply passes configuration from container to hasher constructor.
// All defaults, validation, and parameter handling is done by the hasher itself.
func (h *HashManager) CreateArgon2idDriver() (interface{}, error) {
	// Get configuration from container and pass directly to hasher
	var config map[string]interface{}
	if configValue, exists := h.GetConfig().Get("hashing.argon2id"); exists {
		if configMap, ok := configValue.(map[string]interface{}); ok {
			config = configMap
		}
	}
	return hashers.NewArgon2IdHasher(config), nil
}

// GetDefaultDriver implements the ManagerInterface requirement.
// Returns the default algorithm identifier for driver resolution and delegation.
//
// Default driver features:
//   - Centralized default algorithm configuration
//   - Consistent behavior across all manager operations
//   - Easy configuration override via enum constants
//   - Laravel-style manager pattern compliance
//
// Returns the default algorithm name from enums package.
func (h *HashManager) GetDefaultDriver() string {
	// Delegate to enums package for centralized default management
	// This ensures consistency across the entire hashing system
	return enums.GetDefaultAlgorithm()
}

// Compile-time interface compliance checks
// These ensure HashManager properly implements required interfaces
// Prevents runtime errors from missing method implementations
var _ hashingInterfaces.FactoryInterface = (*HashManager)(nil) // Hasher factory operations
