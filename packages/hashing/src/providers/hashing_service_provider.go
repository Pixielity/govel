// Package providers contains service provider implementations for the hashing package.
// Service providers are responsible for registering hashing services in the
// dependency injection container and configuring them for use throughout the application.
package providers

import (
	"fmt"

	providers "govel/packages/application/providers"
	applicationInterfaces "govel/packages/types/src/interfaces/application/base"

	hashing "govel/packages/hashing/src"
	hashingInterfaces "govel/packages/types/src/interfaces/hashing"
)

// HashingServiceProvider implements a Laravel-compatible service provider
// for the hashing package. It registers hashing services in the dependency
// injection container and provides factory methods for creating hasher instances.
//
// This service provider is equivalent to Laravel's HashServiceProvider and
// follows the same patterns for service registration and binding.
//
// Services registered:
//   - "hash": Singleton HashManager instance for managing hash drivers
//   - "hash.driver": Factory for getting the default hasher driver
//   - hashingInterfaces.HasherInterface: Default hasher interface contract
//   - hashingInterfaces.FactoryInterface: Hash factory interface contract
//
// The provider implements deferred loading, meaning services are only
// instantiated when first requested from the container.
type HashingServiceProvider struct {
	providers.ServiceProvider
}

// NewHashingServiceProvider creates a new HashingServiceProvider instance.
//
// Returns:
//   - *HashingServiceProvider: New service provider instance ready for registration
//
// Example:
//
//	provider := NewHashingServiceProvider()
//	err := provider.Register(application)
//	if err != nil {
//		log.Fatal("Failed to register hashing services:", err)
//	}
func NewHashingServiceProvider() *HashingServiceProvider {
	return &HashingServiceProvider{
		ServiceProvider: providers.ServiceProvider{},
	}
}

// Register registers all hashing services in the dependency injection container.
// This method sets up the service bindings but does not instantiate the services
// immediately (deferred loading).
//
// Registered services:
//   - "hash": Singleton HashManager instance for managing hash drivers
//   - "hash.driver": Factory function that returns the default hasher driver
//   - "hash.contract": HashManager interface contract
//   - "hash.factory.contract": Hash factory interface contract
//
// Returns:
//   - error: Any error that occurred during service registration
//
// Thread-safe: This method should only be called once during application bootstrap.
func (h *HashingServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := h.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register HashManager as a singleton
	// This ensures only one HashManager instance exists throughout the application
	err := application.Singleton(hashingInterfaces.HASHING_TOKEN, h.createHashManagerFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind hash manager: %w", err)
	}

	// Register default hash driver factory
	// This provides quick access to the default hasher without manager overhead
	err = application.Bind(hashingInterfaces.HASH_DRIVER_TOKEN, h.createHashDriverFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind hash driver: %w", err)
	}

	// Register FactoryInterface contract using the HashManager
	err = application.Singleton(hashingInterfaces.HASH_FACTORY_TOKEN, h.createHashManagerFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind hash factory contract: %w", err)
	}

	return nil
}

// Provides returns a list of service tokens that this provider offers.
// This is used by the container to determine which services are available
// and to implement deferred loading.
//
// Returns:
//   - []interface{}: List of service tokens provided by this provider
func (h *HashingServiceProvider) Provides() []interface{} {
	return []interface{}{
		hashingInterfaces.HASHING_TOKEN,      // HashManager singleton
		hashingInterfaces.HASH_DRIVER_TOKEN,  // Default hash driver
		hashingInterfaces.HASH_FACTORY_TOKEN, // FactoryInterface contract
	}
}

// createHashManagerFactory creates a factory function for the HashManager singleton.
// This factory returns the same HashManager instance on subsequent calls.
func (h *HashingServiceProvider) createHashManagerFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		return hashing.NewHashManager(application)
	}
}

// createHashDriverFactory creates a factory function for the default hash driver.
// This factory returns the default hasher driver instance each time it's called.
func (h *HashingServiceProvider) createHashDriverFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		// Get the hash manager from the container
		hashManager, err := application.Make(hashingInterfaces.HASHING_TOKEN)
		if err != nil {
			// Return nil if we can't get the hash manager
			return nil
		}

		// Cast to HashManager and get the default driver
		if manager, ok := hashManager.(*hashing.HashManager); ok {
			driver, err := manager.Driver(manager.GetDefaultDriver())
			if err != nil {
				return nil
			}
			return driver
		}

		return nil
	}
}
