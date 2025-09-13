// Package providers contains service provider implementations for the encryption package.
// Service providers are responsible for registering encryption services in the
// dependency injection container and configuring them for use throughout the application.
package providers

import (
	"fmt"

	providers "govel/application/providers"
	applicationInterfaces "govel/types/src/interfaces/application"

	encryption "govel/encryption/src"
	encryptionInterfaces "govel/types/src/interfaces/encryption"
)

// EncryptionServiceProvider implements a Laravel-compatible service provider
// for the encryption package. It registers encryption services in the dependency
// injection container and provides factory methods for creating encryptioner instances.
//
// This service provider is equivalent to Laravel's EncryptionServiceProvider and
// follows the same patterns for service registration and binding.
//
// Services registered:
//   - "encryption": Singleton EncryptionManager instance for managing encryption drivers
//   - "encryption.driver": Factory for getting the default encryptioner driver
//   - encryptionInterfaces.EncryptionerInterface: Default encryptioner interface contract
//   - encryptionInterfaces.FactoryInterface: Encryption factory interface contract
//
// The provider implements deferred loading, meaning services are only
// instantiated when first requested from the container.
type EncryptionServiceProvider struct {
	providers.ServiceProvider
}

// NewEncryptionServiceProvider creates a new EncryptionServiceProvider instance.
//
// Returns:
//   - *EncryptionServiceProvider: New service provider instance ready for registration
//
// Example:
//
//	provider := NewEncryptionServiceProvider()
//	err := provider.Register(application)
//	if err != nil {
//		log.Fatal("Failed to register encryption services:", err)
//	}
func NewEncryptionServiceProvider() *EncryptionServiceProvider {
	return &EncryptionServiceProvider{
		ServiceProvider: providers.ServiceProvider{},
	}
}

// Register registers all encryption services in the dependency injection container.
// This method sets up the service bindings but does not instantiate the services
// immediately (deferred loading).
//
// Registered services:
//   - "encryption": Singleton EncryptionManager instance for managing encryption drivers
//   - "encryption.driver": Factory function that returns the default encryptioner driver
//   - "encryption.contract": EncryptionManager interface contract
//   - "encryption.factory.contract": Encryption factory interface contract
//
// Returns:
//   - error: Any error that occurred during service registration
//
// Thread-safe: This method should only be called once during application bootstrap.
func (h *EncryptionServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := h.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register EncryptionManager as a singleton
	// This ensures only one EncryptionManager instance exists throughout the application
	err := application.Singleton(encryptionInterfaces.ENCRYPTION_TOKEN, h.createEncryptionManagerFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind encryption manager: %w", err)
	}

	// Register default encryption driver factory
	// This provides quick access to the default encryptioner without manager overhead
	err = application.Bind(encryptionInterfaces.ENCRYPTION_DRIVER_TOKEN, h.createEncryptionDriverFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind encryption driver: %w", err)
	}

	// Register FactoryInterface contract using the EncryptionManager
	err = application.Singleton(encryptionInterfaces.ENCRYPTION_FACTORY_TOKEN, h.createEncryptionManagerFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind encryption factory contract: %w", err)
	}

	return nil
}

// Provides returns a list of service tokens that this provider offers.
// This is used by the container to determine which services are available
// and to implement deferred loading.
//
// Returns:
//   - []interface{}: List of service tokens provided by this provider
func (h *EncryptionServiceProvider) Provides() []interface{} {
	return []interface{}{
		encryptionInterfaces.ENCRYPTION_TOKEN,         // EncryptionManager singleton
		encryptionInterfaces.ENCRYPTION_DRIVER_TOKEN,  // Default encryption driver
		encryptionInterfaces.ENCRYPTION_FACTORY_TOKEN, // FactoryInterface contract
	}
}

// createEncryptionManagerFactory creates a factory function for the EncryptionManager singleton.
// This factory returns the same EncryptionManager instance on subsequent calls.
func (h *EncryptionServiceProvider) createEncryptionManagerFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		return encryption.NewEncryptionManager(application)
	}
}

// createEncryptionDriverFactory creates a factory function for the default encryption driver.
// This factory returns the default encryptioner driver instance each time it's called.
func (h *EncryptionServiceProvider) createEncryptionDriverFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		// Get the encryption manager from the container
		encryptionManager, err := application.Make(encryptionInterfaces.ENCRYPTION_TOKEN)
		if err != nil {
			// Return nil if we can't get the encryption manager
			return nil
		}

		// Cast to EncryptionManager and get the default driver
		if manager, ok := encryptionManager.(*encryption.EncryptionManager); ok {
			driver, err := manager.Driver(manager.GetDefaultDriver())
			if err != nil {
				return nil
			}
			return driver
		}

		return nil
	}
}
