package concerns

import (
	"fmt"
	"reflect"
	"sync"

	applicationInterfaces "govel/types/interfaces/application/base"
	concernsInterfaces "govel/types/interfaces/application/concerns"
)

/**
 * HasBootstrap provides application bootstrap management functionality.
 * This trait implements the HasBootstrapInterface and manages the bootstrap
 * classes and bootstrapping process for the application.
 *
 * Features:
 * - Bootstrap class registration and management
 * - Complete application bootstrapping process
 * - Selective bootstrapping with custom bootstrappers
 * - Bootstrap without service providers
 * - Bootstrap state tracking
 * - Thread-safe bootstrap operations
 */
type HasBootstrap struct {
	/**
	 * bootstrappers holds the bootstrap classes for the application
	 */
	bootstrappers []interface{}

	/**
	 * hasBeenBootstrapped tracks whether the application has been bootstrapped
	 */
	hasBeenBootstrapped bool

	/**
	 * mutex provides thread-safe access to bootstrap fields
	 */
	mutex sync.RWMutex

	/**
	 * applicationRef holds a reference to the application instance for self-referential operations
	 * This allows the concern to call methods on the application that embeds it
	 */
	applicationRef interface{}
}

// NewBootstrap creates a new bootstrap concern with optional bootstrap classes.
//
// Parameters:
//
//	bootstrappers: Optional slice of bootstrap class instances
//
// Returns:
//
//	*HasBootstrap: A new bootstrap concern instance
//
// Example:
//
//	// Basic bootstrap concern
//	bootstrap := NewBootstrap()
//	// With predefined bootstrappers
//	bootstrap := NewBootstrap([]interface{}{
//	    &ConfigBootstrapper{},
//	    &DatabaseBootstrapper{},
//	})
func NewBootstrap(bootstrappers ...[]interface{}) *HasBootstrap {
	var bootstrapperList []interface{}

	if len(bootstrappers) > 0 && bootstrappers[0] != nil {
		bootstrapperList = bootstrappers[0]
	} else {
		bootstrapperList = make([]interface{}, 0)
	}

	return &HasBootstrap{
		bootstrappers:       bootstrapperList,
		hasBeenBootstrapped: false,
	}
}

// SetApplicationRef sets a reference to the application instance.
// This allows the concern to perform operations that require access to the full application.
//
// Parameters:
//
//	application: The application instance reference
//
// Example:
//
//	bootstrap.SetApplicationRef(application)
func (b *HasBootstrap) SetApplicationRef(application applicationInterfaces.ApplicationInterface) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.applicationRef = application
}

// Bootstrap bootstraps the application with all registered bootstrappers.
// This method checks if the application has been bootstrapped and performs
// the complete bootstrap process including deferred provider loading.
//
// Returns:
//
//	error: Any error that occurred during bootstrapping
//
// Example:
//
//	err := application.Bootstrap()
//	if err != nil {
//	    log.Fatalf("Application bootstrap failed: %v", err)
//	}
func (b *HasBootstrap) Bootstrap() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if already bootstrapped
	if b.hasBeenBootstrapped {
		return nil
	}

	// Bootstrap with registered bootstrappers
	if err := b.bootstrapWithInternal(b.bootstrappers); err != nil {
		return fmt.Errorf("bootstrap failed: %w", err)
	}

	// Load deferred providers if application reference is available
	if err := b.loadDeferredProviders(); err != nil {
		return fmt.Errorf("failed to load deferred providers: %w", err)
	}

	// Mark as bootstrapped
	b.hasBeenBootstrapped = true
	return nil
}

// BootstrapWith bootstraps the application with specific bootstrappers.
// This method allows selective bootstrapping with a custom set of bootstrappers.
//
// Parameters:
//
//	bootstrappers: A slice of bootstrapper instances to use for bootstrapping
//
// Returns:
//
//	error: Any error that occurred during bootstrapping
//
// Example:
//
//	customBootstrappers := []interface{}{
//	    &ConfigBootstrapper{},
//	    &DatabaseBootstrapper{},
//	}
//	err := application.BootstrapWith(customBootstrappers)
func (b *HasBootstrap) BootstrapWith(bootstrappers []interface{}) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.bootstrapWithInternal(bootstrappers); err != nil {
		return fmt.Errorf("bootstrap with custom bootstrappers failed: %w", err)
	}

	// Mark as bootstrapped
	b.hasBeenBootstrapped = true
	return nil
}

// BootstrapWithoutProviders bootstraps the application without booting service providers.
// This method performs bootstrapping but excludes provider-related bootstrappers.
//
// Returns:
//
//	error: Any error that occurred during bootstrapping
//
// Example:
//
//	err := application.BootstrapWithoutProviders()
//	if err != nil {
//	    log.Printf("Bootstrap without providers failed: %v", err)
//	}
func (b *HasBootstrap) BootstrapWithoutProviders() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Filter out provider bootstrappers
	filteredBootstrappers := make([]interface{}, 0, len(b.bootstrappers))
	for _, bootstrapper := range b.bootstrappers {
		if !b.isProviderBootstrapper(bootstrapper) {
			filteredBootstrappers = append(filteredBootstrappers, bootstrapper)
		}
	}

	if err := b.bootstrapWithInternal(filteredBootstrappers); err != nil {
		return fmt.Errorf("bootstrap without providers failed: %w", err)
	}

	// Mark as bootstrapped
	b.hasBeenBootstrapped = true
	return nil
}

// HasBeenBootstrapped returns whether the application has been bootstrapped.
//
// Returns:
//
//	bool: true if the application has been bootstrapped, false otherwise
//
// Example:
//
//	if !application.HasBeenBootstrapped() {
//	    application.Bootstrap()
//	}
func (b *HasBootstrap) HasBeenBootstrapped() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.hasBeenBootstrapped
}

// SetBootstrapped marks the application as bootstrapped or not bootstrapped.
//
// Parameters:
//
//	bootstrapped: Whether the application should be marked as bootstrapped
//
// Example:
//
//	application.SetBootstrapped(true)  // Mark as bootstrapped
//	application.SetBootstrapped(false) // Reset bootstrap state
func (b *HasBootstrap) SetBootstrapped(bootstrapped bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.hasBeenBootstrapped = bootstrapped
}

// GetBootstrappers returns the list of registered bootstrap classes.
//
// Returns:
//
//	[]interface{}: List of bootstrap class instances
//
// Example:
//
//	bootstrappers := application.GetBootstrappers()
//	fmt.Printf("Registered %d bootstrappers\n", len(bootstrappers))
func (b *HasBootstrap) GetBootstrappers() []interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Return a copy to prevent external modification
	bootstrappers := make([]interface{}, len(b.bootstrappers))
	copy(bootstrappers, b.bootstrappers)
	return bootstrappers
}

// SetBootstrappers sets the bootstrap classes for the application.
//
// Parameters:
//
//	bootstrappers: A slice of bootstrap class instances
//
// Example:
//
//	bootstrappers := []interface{}{
//	    &ConfigBootstrapper{},
//	    &DatabaseBootstrapper{},
//	    &CacheBootstrapper{},
//	}
//	application.SetBootstrappers(bootstrappers)
func (b *HasBootstrap) SetBootstrappers(bootstrappers []interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Create a copy to prevent external modification
	b.bootstrappers = make([]interface{}, len(bootstrappers))
	copy(b.bootstrappers, bootstrappers)
}

// AddBootstrapper adds a single bootstrap class to the application.
//
// Parameters:
//
//	bootstrapper: The bootstrap class instance to add
//
// Example:
//
//	application.AddBootstrapper(&LoggingBootstrapper{})
func (b *HasBootstrap) AddBootstrapper(bootstrapper interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if bootstrapper != nil {
		b.bootstrappers = append(b.bootstrappers, bootstrapper)
	}
}

// RemoveBootstrapper removes a bootstrap class from the application.
//
// Parameters:
//
//	bootstrapper: The bootstrap class instance to remove
//
// Returns:
//
//	bool: true if the bootstrapper was found and removed, false otherwise
//
// Example:
//
//	removed := application.RemoveBootstrapper(&LoggingBootstrapper{})
//	if removed {
//	    fmt.Println("Bootstrapper removed successfully")
//	}
func (b *HasBootstrap) RemoveBootstrapper(bootstrapper interface{}) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, existing := range b.bootstrappers {
		if b.bootstrappersEqual(existing, bootstrapper) {
			// Remove by slicing
			b.bootstrappers = append(b.bootstrappers[:i], b.bootstrappers[i+1:]...)
			return true
		}
	}
	return false
}

// ClearBootstrappers removes all bootstrap classes.
//
// Example:
//
//	application.ClearBootstrappers() // Reset all bootstrappers
func (b *HasBootstrap) ClearBootstrappers() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bootstrappers = make([]interface{}, 0)
}

// Helper methods

// bootstrapWithInternal performs the actual bootstrapping process
func (b *HasBootstrap) bootstrapWithInternal(bootstrappers []interface{}) error {
	for _, bootstrapper := range bootstrappers {
		if err := b.executeBootstrapper(bootstrapper); err != nil {
			return fmt.Errorf("bootstrapper %T failed: %w", bootstrapper, err)
		}
	}
	return nil
}

// executeBootstrapper executes a single bootstrapper
func (b *HasBootstrap) executeBootstrapper(bootstrapper interface{}) error {
	// This is where you would implement the actual bootstrapper execution logic
	// The specific implementation depends on your bootstrapper interface design

	// Example pattern (adapt based on your actual bootstrapper interface):
	// if executableBootstrapper, ok := bootstrapper.(BootstrapperInterface); ok {
	//     return executableBootstrapper.Bootstrap(b.applicationRef)
	// }

	// For now, this is a placeholder that always succeeds
	return nil
}

// loadDeferredProviders loads deferred providers if the application reference supports it
func (b *HasBootstrap) loadDeferredProviders() error {
	if b.applicationRef == nil {
		return nil // No application reference available
	}

	// Try to call LoadDeferredProviders method on the application if it exists
	// This uses reflection to avoid circular dependencies
	appValue := reflect.ValueOf(b.applicationRef)
	if appValue.Kind() == reflect.Ptr {
		appValue = appValue.Elem()
	}

	// Look for LoadDeferredProviders method
	method := appValue.MethodByName("LoadDeferredProviders")
	if method.IsValid() && method.Type().NumIn() == 0 {
		results := method.Call(nil)
		if len(results) > 0 && !results[0].IsNil() {
			if err, ok := results[0].Interface().(error); ok {
				return err
			}
		}
	}

	return nil
}

// isProviderBootstrapper checks if a bootstrapper is provider-related
func (b *HasBootstrap) isProviderBootstrapper(bootstrapper interface{}) bool {
	// This would need to be adapted based on your actual provider bootstrapper types
	// For example, you might check the type name or interface implementation

	bootstrapperType := reflect.TypeOf(bootstrapper).String()

	// Example patterns - adapt based on your actual provider bootstrapper names
	providerPatterns := []string{
		"BootProviders",
		"ProviderBootstrapper",
		"ServiceProviderBootstrapper",
	}

	for _, pattern := range providerPatterns {
		if bootstrapperType == pattern ||
			(len(bootstrapperType) > len(pattern) &&
				bootstrapperType[len(bootstrapperType)-len(pattern):] == pattern) {
			return true
		}
	}

	return false
}

// bootstrappersEqual compares two bootstrappers for equality
func (b *HasBootstrap) bootstrappersEqual(a, bootstrapper interface{}) bool {
	// Simple type-based comparison - you might want more sophisticated comparison
	return reflect.TypeOf(a) == reflect.TypeOf(bootstrapper)
}

// RegisterDefaultBootstrappers returns the default bootstrap classes for the application.
// This method should be overridden in the actual application implementation to return
// the specific bootstrappers needed for that application.
//
// Returns:
//
//	[]interface{}: List of default bootstrap class instances (empty by default)
//
// Example:
//
//	bootstrappers := application.RegisterDefaultBootstrappers()
//	fmt.Printf("Default bootstrappers: %d\n", len(bootstrappers))
//
// Note: This base implementation returns an empty slice. Applications should
// override this method to return their specific bootstrap requirements.
func (b *HasBootstrap) RegisterDefaultBootstrappers() []interface{} {
	// Return empty slice by default - applications will override this
	return []interface{}{}
}

// Compile-time interface compliance check
var _ concernsInterfaces.HasBootstrapInterface = (*HasBootstrap)(nil)
