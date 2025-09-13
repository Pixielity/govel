package symbol

import (
	"sync"
)

// globalSymbolRegistry stores globally registered symbols by their string keys.
// This registry implements the TypeScript Symbol.for() functionality, ensuring
// that the same key always returns the same symbol instance across the entire application.
//
// The registry is thread-safe and uses a RWMutex for concurrent access:
// - Read operations (symbol lookup) can happen concurrently
// - Write operations (symbol creation/registration) are exclusive
//
// Key characteristics:
// - Singleton behavior: Same key = same symbol, always
// - Global scope: Accessible from any package in the application
// - Thread-safe: Safe for concurrent access from multiple goroutines
// - Persistent: Symbols remain in registry until explicitly cleared
//
// The registry uses string keys following a hierarchical naming convention:
// - "govel.interfaces.logger" for core framework interfaces
// - "govel.services.user" for application services
// - "app.custom.service" for application-specific services
// - "Symbol.iterator" for well-known symbols (similar to JavaScript)
var globalSymbolRegistry = make(map[string]SymbolType)

// registryMutex provides thread-safe access to the globalSymbolRegistry.
// It uses a RWMutex to allow concurrent reads while ensuring exclusive writes.
//
// Read operations (RLock):
// - symbol.For() when symbol exists
// - symbol.KeyFor() lookups
// - Registry introspection methods
//
// Write operations (Lock):
// - symbol.For() when creating new symbols
// - Registry modification methods
//
// This design optimizes for the common case where symbols are looked up
// more frequently than they are created.
var registryMutex sync.RWMutex

// IsRegistered checks if a key is already registered in the global registry.
// This method provides a fast way to check for key existence without
// creating a symbol if it doesn't exist.
//
// Parameters:
//
//	key: The key to check in the global registry
//
// Returns:
//
//	bool: true if the key exists in the registry, false otherwise
func IsRegistered(key string) bool {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	_, exists := globalSymbolRegistry[key]
	return exists
}

// GetRegisteredKeys returns all keys currently registered in the global registry.
// The returned slice is a copy, so modifications will not affect the registry.
//
// Returns:
//
//	[]string: Slice of all keys currently registered in the global registry
func GetRegisteredKeys() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	keys := make([]string, 0, len(globalSymbolRegistry))
	for key := range globalSymbolRegistry {
		keys = append(keys, key)
	}
	return keys
}

// GetRegistrySize returns the number of symbols in the global registry.
// This method provides a quick way to check registry size without
// copying the entire registry.
//
// Returns:
//
//	int: Number of symbols currently registered globally
func GetRegistrySize() int {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return len(globalSymbolRegistry)
}

// GetGlobalSymbols returns all symbols in the global registry.
// This is useful for debugging and introspection.
//
// Returns:
//
//	map[string]SymbolType: Copy of the global symbol registry
func GetGlobalSymbols() map[string]SymbolType {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	result := make(map[string]SymbolType, len(globalSymbolRegistry))
	for key, symbol := range globalSymbolRegistry {
		result[key] = symbol
	}
	return result
}

// ClearGlobalSymbols clears the global symbol registry.
// This is primarily useful for testing scenarios.
//
// Warning: This will affect any code relying on global symbols.
// Use with caution, typically only in test cleanup.
func ClearGlobalSymbols() {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	globalSymbolRegistry = make(map[string]SymbolType)
}
