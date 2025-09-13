// Package factories provides factory functions and classes for creating framework-specific implementations.
// This package implements the Factory pattern to abstract the creation of complex objects
// and provide a centralized way to instantiate webserver components.
//
// The adapter factory specifically handles the creation of HTTP framework adapters
// that bridge different Go web frameworks (Gin, Echo, Fiber, net/http) to a unified interface.
//
// Key Features:
//   - Centralized adapter creation and management
//   - Plugin-style architecture for adding new framework support
//   - Registry-based adapter discovery and instantiation
//   - Error handling for unsupported or missing adapters
//   - Type-safe adapter creation with interface contracts
//
// Design Patterns Used:
//   - Factory Pattern: For creating adapter instances
//   - Registry Pattern: For managing available adapters
//   - Strategy Pattern: Different adapters for different frameworks
//   - Abstract Factory: Consistent creation interface across adapters
package factories

import (
	"fmt"

	"govel/packages/new/webserver/src/adapters"
	"govel/packages/new/webserver/src/enums"
	"govel/packages/new/webserver/src/interfaces"
)

// AdapterFactory provides methods for creating HTTP framework adapter instances.
// This factory encapsulates the logic for instantiating different types of adapters
// based on the selected web framework engine.
//
// The factory uses a registry-based approach where adapters register themselves
// with the global adapter registry, and the factory looks up and instantiates
// the appropriate adapter based on the engine type.
//
// Thread Safety:
//
//	The AdapterFactory is thread-safe for read operations (creating adapters).
//	The underlying adapter registry should be populated during application
//	initialization before concurrent access begins.
//
// Usage Examples:
//
//	// Create factory instance
//	factory := NewAdapterFactory()
//
//	// Create Gin adapter
//	adapter, err := factory.CreateAdapter("gin")
//	if err != nil {
//	    log.Fatal("Failed to create adapter:", err)
//	}
//
//	// Use with different engines
//	fiberAdapter, _ := factory.CreateAdapter("fiber")
//	echoAdapter, _ := factory.CreateAdapter("echo")
type AdapterFactory struct {
	// registry holds references to adapter factory functions
	// This could be extended in the future to hold additional metadata
	registry map[enums.Engine]func() interfaces.AdapterInterface
}

// NewAdapterFactory creates a new adapter factory instance.
// The factory is initialized with access to the global adapter registry
// and provides methods for creating adapter instances.
//
// Returns:
//   - *AdapterFactory: A new adapter factory instance ready for use
//
// Example:
//
//	factory := NewAdapterFactory()
//	adapter, err := factory.CreateAdapter("gin")
func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{
		registry: adapters.AdapterRegistry,
	}
}

// CreateAdapter creates a new adapter instance for the specified engine name.
// This method performs engine name validation, registry lookup, and adapter instantiation.
//
// The creation process follows these steps:
//  1. Parse and validate the engine name string
//  2. Look up the adapter factory function in the registry
//  3. Instantiate the adapter using the factory function
//  4. Return the initialized adapter or error
//
// Parameters:
//   - engineName: String name of the web framework ("gin", "echo", "fiber", "net/http")
//     The name is case-insensitive and will be normalized internally.
//
// Returns:
//   - interfaces.AdapterInterface: A new adapter instance configured for the specified engine
//   - error: Error if the engine is unsupported, not registered, or instantiation fails
//
// Possible Errors:
//   - "unsupported engine: <name>" - The engine name is not recognized
//   - "no adapter registered for engine: <name>" - The engine is valid but no adapter is available
//   - Adapter-specific initialization errors from the factory function
//
// Thread Safety:
//
//	This method is safe for concurrent use as it only performs read operations
//	on the registry and creates new instances without shared state.
//
// Usage Examples:
//
//	// Basic usage
//	adapter, err := factory.CreateAdapter("gin")
//	if err != nil {
//	    return fmt.Errorf("adapter creation failed: %w", err)
//	}
//
//	// Error handling for unsupported engines
//	adapter, err := factory.CreateAdapter("unsupported")
//	if err != nil {
//	    log.Printf("Engine not available: %v", err)
//	    // Fallback to default engine
//	    adapter, _ = factory.CreateAdapter("net/http")
//	}
//
//	// Case-insensitive engine names
//	ginAdapter, _ := factory.CreateAdapter("GIN")     // Works
//	fiberAdapter, _ := factory.CreateAdapter("Fiber") // Works
func (f *AdapterFactory) CreateAdapter(engine enums.Engine) (interfaces.AdapterInterface, error) {
	// Look up the adapter factory function in the registry
	// The registry maps engine enums to factory functions
	factory, exists := f.registry[engine]
	if !exists {
		return nil, fmt.Errorf("no adapter registered for engine: %s", engine.Name())
	}

	// Instantiate and return the adapter
	// The factory function handles all engine-specific initialization
	return factory(), nil
}

// GetSupportedEngines returns a list of all engines that have registered adapters.
// This method is useful for runtime discovery of available web framework support
// and for providing user feedback about supported options.
//
// Returns:
//   - []enums.Engine: Slice of engine enums that have registered adapters
//   - []string: Slice of engine names as strings for display purposes
//
// Usage Examples:
//
//	// Get supported engines for validation
//	engines, names := factory.GetSupportedEngines()
//	fmt.Printf("Supported engines: %v\n", names)
//
//	// Validate user input
//	if !contains(names, userChoice) {
//	    return fmt.Errorf("unsupported engine: %s, available: %v", userChoice, names)
//	}
func (f *AdapterFactory) GetSupportedEngines() ([]enums.Engine, []string) {
	engines := make([]enums.Engine, 0, len(f.registry))
	names := make([]string, 0, len(f.registry))

	// Extract engine enums and their string representations
	for engine := range f.registry {
		engines = append(engines, engine)
		names = append(names, engine.Name())
	}

	return engines, names
}

// HasAdapter checks if an adapter is registered for the specified engine.
// This is useful for conditional logic and validation without creating instances.
//
// Parameters:
//   - engineName: String name of the engine to check
//
// Returns:
//   - bool: True if an adapter is available for the engine, false otherwise
//
// Example:
//
//	if factory.HasAdapter("gin") {
//	    adapter, _ := factory.CreateAdapter("gin")
//	} else {
//	    log.Warn("Gin adapter not available, using fallback")
//	}
func (f *AdapterFactory) HasAdapter(engineName string) bool {
	engine, valid := enums.ParseEngine(engineName)
	if !valid {
		return false
	}

	_, exists := f.registry[engine]
	return exists
}

// Package-level convenience functions for backward compatibility
// and simplified usage patterns.

// CreateAdapter creates a new adapter instance using the default factory.
// This is a convenience function that creates a factory internally and
// delegates to the factory's CreateAdapter method.
//
// Parameters:
//   - engineName: String name of the web framework engine
//
// Returns:
//   - interfaces.AdapterInterface: A new adapter instance
//   - error: Error if creation fails
//
// Note: This function creates a new factory instance on each call.
//
//	For multiple adapter creations, consider using a factory instance directly.
//
// Example:
//
//	adapter, err := factories.CreateAdapter("gin")
//	if err != nil {
//	    log.Fatal("Failed to create adapter:", err)
//	}
func CreateAdapter(engineName enums.Engine) (interfaces.AdapterInterface, error) {
	factory := NewAdapterFactory()
	return factory.CreateAdapter(engineName)
}

// CreateAdapterFromEnum creates an adapter using an engine enum directly.
// This function provides type-safe adapter creation for internal use
// where the engine type is already known and validated.
//
// Parameters:
//   - engine: The engine enum value
//
// Returns:
//   - interfaces.AdapterInterface: A new adapter instance
//   - error: Error if the engine is not registered or creation fails
//
// Example:
//
//	adapter, err := CreateAdapterFromEnum(enums.GoFiber)
//	if err != nil {
//	    return fmt.Errorf("fiber adapter creation failed: %w", err)
//	}
func CreateAdapterFromEnum(engine enums.Engine) (interfaces.AdapterInterface, error) {
	factory := NewAdapterFactory()
	return factory.CreateAdapter(engine)
}
