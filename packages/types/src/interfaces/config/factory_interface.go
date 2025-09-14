package interfaces

import enums "govel/types/enums/config"

// FactoryInterface defines the contract for configuration driver factory operations.
// This interface provides methods for creating and managing different configuration drivers.
type FactoryInterface interface {
	// Config gets a configuration driver instance by driver name (optional).
	// If no name is provided, uses the default driver.
	// Returns nil if driver is invalid or creation fails.
	Config(name ...enums.Driver) ConfigInterface
}
