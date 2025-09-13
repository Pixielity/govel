package types

import (
	"govel/packages/support/src/symbol"
)

// ServiceIdentifier represents a service identifier in the container.
type ServiceIdentifier interface{}

// ToKey converts a ServiceIdentifier to a string key.
// It handles both string and Symbol types.
func ToKey(identifier ServiceIdentifier) string {
	switch v := identifier.(type) {
	case string:
		return v
	case *symbol.SymbolType:
		return v.String()
	case symbol.SymbolType:
		return v.String()
	default:
		// For any other type, try to convert to string
		if stringer, ok := v.(interface{ String() string }); ok {
			return stringer.String()
		}
		// Fallback: return empty string or could panic/error
		return ""
	}
}
