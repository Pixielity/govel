package interfaces

// Driver defines the interface that all configuration drivers must implement.
// This interface provides a unified way to interact with different configuration sources
// like files, environment variables, remote stores, etc.
type DriverInterface interface {
	// Get retrieves a configuration value by key.
	// Returns the value and an error if the operation fails.
	Get(key string) (interface{}, error)

	// Set stores a configuration value by key.
	// Returns an error if the operation fails.
	Set(key string, value interface{}) error

	// Load loads the configuration from the source.
	// This method should be called to initialize or refresh the configuration.
	Load() error

	// Watch starts watching for configuration changes.
	// The callback function will be called when changes are detected.
	// Returns an error if watching cannot be started.
	Watch(callback func()) error

	// Unwatch stops watching for configuration changes.
	// Returns an error if unwatching fails.
	Unwatch() error

	// Invalidate clears the configuration cache and forces a reload on next access.
	// Returns an error if invalidation fails.
	Invalidate() error

	// GetAll returns all configuration data as a map.
	// Returns the configuration map and an error if the operation fails.
	GetAll() (map[string]interface{}, error)

	// Has checks if a configuration key exists.
	// Returns true if the key exists, false otherwise.
	Has(key string) bool

	// Delete removes a configuration key.
	// Returns an error if the operation fails.
	Delete(key string) error
}
