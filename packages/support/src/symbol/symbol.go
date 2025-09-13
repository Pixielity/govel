package symbol

import (
	"fmt"
	"sync/atomic"
)

// symbolCounter provides atomic incrementation for unique symbol IDs
var symbolCounter uint64

// New creates a new unique symbol with an optional description.
// This is the Go equivalent of TypeScript's Symbol() constructor.
//
// IMPORTANT: Unlike symbol.For(), symbols created with New() are NOT added
// to the global registry. Each call to New() returns a completely unique symbol,
// even if the description is identical. These symbols are "private" and cannot
// be retrieved from other packages.
//
// Key characteristics:
// - Always unique: Each call creates a new symbol, even with same description
// - Not globally registered: Cannot be accessed via symbol.For() or symbol.KeyFor()
// - Thread-safe: Uses atomic counter for ID generation
// - Local scope: Perfect for package-internal or module-internal identifiers
// - Memory efficient: No registry overhead
//
// When to use New() vs For():
// - Use New() for: Internal package symbols, temporary identifiers, private keys
// - Use For() for: Cross-package symbols, dependency injection tokens, shared identifiers
//
// The symbol will have:
// - A unique numeric ID (auto-incremented atomically)
// - An optional description for debugging
// - All standard SymbolType methods (String, ID, Description, Equals)
// - IsGlobal() will return false
//
// Performance characteristics:
// - Very fast: Just atomic increment + struct creation
// - No locking: No registry access means no mutex overhead
// - No memory overhead: Symbol only exists in your variable
//
// Parameters:
//
//	description: Optional human-readable description for debugging and logging.
//	             This does not affect uniqueness - multiple symbols can have
//	             the same description but will still be unique.
//
// Returns:
//
//	SymbolType: A new unique symbol that is NOT registered globally
//
// Example:
//
//	// Each symbol is unique, even with same description
//	logger1 := symbol.New("Logger")
//	logger2 := symbol.New("Logger")
//	database := symbol.New("Database")
//
//	fmt.Println(logger1 == logger2)      // false - different symbols
//	fmt.Println(logger1.Equals(logger2)) // false - different IDs
//
//	// They have readable string representations
//	fmt.Println(logger1.String()) // "Symbol(1: Logger)"
//	fmt.Println(logger2.String()) // "Symbol(2: Logger)"
//	fmt.Println(database.String()) // "Symbol(3: Database)"
//
//	// But they're not globally accessible
//	fmt.Println(logger1.IsGlobal()) // false
//	fmt.Println(symbol.KeyFor(logger1)) // "" (empty string)
//
//	// Perfect for internal package use
//	var (
//	    internalCacheKey = symbol.New("internal.cache")
//	    tempProcessKey = symbol.New("temp.process")
//	)
//
//	// Use as map keys (guaranteed unique)
//	cache := make(map[string]interface{})
//	cache[internalCacheKey.String()] = "some data"
func New(description ...string) SymbolType {
	id := atomic.AddUint64(&symbolCounter, 1)

	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}

	return SymbolType{
		id:          id,
		description: desc,
	}
}

// For retrieves or creates a symbol in the global symbol registry.
// If a symbol with the given key already exists, it returns that symbol.
// Otherwise, it creates a new symbol and registers it with the key.
//
// This is equivalent to TypeScript's Symbol.for() method.
// Unlike New(), For() can return the same symbol for the same key.
//
// Key behavior:
// - With key: Returns same symbol for same key (shared across packages)
// - Without key: Generates unique key and registers symbol globally
//
// Parameters:
//
//	key: Optional global key to register/retrieve the symbol.
//	     If not provided, generates a unique key automatically.
//
// Returns:
//
//	SymbolType: The symbol associated with the key
//
// Example:
//
//	// With explicit key (recommended for shared symbols)
//	config1 := symbol.For("app.config")
//	config2 := symbol.For("app.config")
//	fmt.Println(config1 == config2) // true - same symbol returned
//
//	// With auto-generated key (unique but globally registered)
//	unique1 := symbol.For() // Auto-generates key like "symbol.auto.1"
//	unique2 := symbol.For() // Auto-generates key like "symbol.auto.2"
//	fmt.Println(unique1 == unique2) // false - different auto keys
//	fmt.Println(unique1.IsGlobal()) // true - both are globally registered
func For(key ...string) SymbolType {
	// Determine the key to use
	var actualKey string
	if len(key) > 0 && key[0] != "" {
		// Use provided key
		actualKey = key[0]
	} else {
		// Auto-generate unique key
		id := atomic.AddUint64(&symbolCounter, 1)
		actualKey = fmt.Sprintf("symbol.%d", id)
	}

	// Check if symbol already exists
	registryMutex.RLock()
	if symbol, exists := globalSymbolRegistry[actualKey]; exists {
		registryMutex.RUnlock()
		return symbol
	}
	registryMutex.RUnlock()

	// Need to create new symbol
	registryMutex.Lock()
	defer registryMutex.Unlock()

	// Double-check in case another goroutine created it
	if symbol, exists := globalSymbolRegistry[actualKey]; exists {
		return symbol
	}

	// Create new symbol with the key as description
	symbol := SymbolType{
		id:          atomic.AddUint64(&symbolCounter, 1),
		description: actualKey,
	}

	globalSymbolRegistry[actualKey] = symbol
	return symbol
}

// KeyFor returns the key associated with a global symbol.
// If the symbol is not in the global registry, returns an empty string.
//
// This is equivalent to TypeScript's Symbol.keyFor() method.
//
// Parameters:
//
//	sym: The symbol to look up
//
// Returns:
//
//	string: The key if found in global registry, empty string otherwise
//
// Example:
//
//	config := symbol.For("app.config")
//	key := symbol.KeyFor(config)
//	fmt.Println(key) // "app.config"
//
//	local := symbol.New("local")
//	key2 := symbol.KeyFor(local)
//	fmt.Println(key2) // "" - not in global registry
func KeyFor(sym SymbolType) string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	for key, registeredSymbol := range globalSymbolRegistry {
		if registeredSymbol.id == sym.id {
			return key
		}
	}
	return ""
}

// String returns a string representation of the symbol.
// This provides a consistent way to use symbols as string keys.
//
// The format is "Symbol(<id>)" or "Symbol(<id>: <description>)" if description exists.
//
// Returns:
//
//	string: String representation of the symbol
//
// Example:
//
//	logger := symbol.New("Logger")
//	fmt.Println(logger.String()) // "Symbol(1: Logger)"
//
//	anonymous := symbol.New()
//	fmt.Println(anonymous.String()) // "Symbol(2)"
func (s SymbolType) String() string {
	if s.description != "" {
		return fmt.Sprintf("Symbol(%d: %s)", s.id, s.description)
	}
	return fmt.Sprintf("Symbol(%d)", s.id)
}

// ID returns the unique numeric identifier of the symbol.
// This can be useful for debugging or internal operations.
//
// Returns:
//
//	uint64: The unique ID of this symbol
func (s SymbolType) ID() uint64 {
	return s.id
}

// Description returns the description of the symbol.
// Returns an empty string if no description was provided.
//
// Returns:
//
//	string: The symbol's description
func (s SymbolType) Description() string {
	return s.description
}

// Equals checks if two symbols are the same.
// Symbols are equal only if they have the same ID.
//
// Parameters:
//
//	other: The other symbol to compare with
//
// Returns:
//
//	bool: true if symbols are equal, false otherwise
func (s SymbolType) Equals(other SymbolType) bool {
	return s.id == other.id
}

// IsGlobal checks if this symbol is registered in the global registry.
//
// Returns:
//
//	bool: true if the symbol is in the global registry, false otherwise
func (s SymbolType) IsGlobal() bool {
	return KeyFor(s) != ""
}
