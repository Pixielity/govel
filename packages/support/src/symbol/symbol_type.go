package symbol

// SymbolType represents a unique symbol identifier similar to TypeScript's Symbol.
// Each symbol is guaranteed to be globally unique, even if created with the same description.
//
// Symbols are immutable and thread-safe, making them safe for concurrent access.
// They can be used as unique keys for maps, service identifiers in dependency injection,
// event types, or any scenario requiring guaranteed uniqueness.
//
// Key characteristics:
// - Immutable: Once created, a symbol's properties cannot be changed
// - Unique: Each symbol has a unique ID, even with identical descriptions
// - Thread-safe: Safe for concurrent access across goroutines
// - Comparable: Symbols can be compared for equality using == or the Equals method
// - Serializable: Symbols have a consistent string representation
//
// The symbol contains two main properties:
// - id: A unique numeric identifier (auto-incrementing)
// - description: An optional human-readable description for debugging
//
// Example:
//
//	// Create symbols
//	logger1 := symbol.New("Logger")
//	logger2 := symbol.New("Logger")
//
//	// They are different despite same description
//	fmt.Println(logger1 == logger2) // false
//	fmt.Println(logger1.Equals(logger2)) // false
//
//	// But have readable descriptions
//	fmt.Println(logger1.String()) // "Symbol(1: Logger)"
//	fmt.Println(logger2.String()) // "Symbol(2: Logger)"
type SymbolType struct {
	// id is the unique numeric identifier for this symbol.
	// This ID is automatically assigned and incremented atomically
	// to ensure global uniqueness across all symbols.
	id uint64

	// description is an optional human-readable description of the symbol.
	// This is used for debugging and logging purposes and does not affect
	// the symbol's uniqueness. Multiple symbols can have the same description
	// but will still be unique based on their ID.
	description string
}
