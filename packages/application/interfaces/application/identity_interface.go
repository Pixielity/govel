package interfaces

// ApplicationIdentityInterface defines the contract for components that manage
// application identity information such as name and version.
// This interface follows the Interface Segregation Principle by focusing
// solely on identity-related operations.
type ApplicationIdentityInterface interface {
	// GetName returns the application name.
	//
	// Returns:
	//   string: The current application name
	GetName() string

	// SetName sets the application name.
	//
	// Parameters:
	//   name: The new application name
	SetName(name string)

	// Name returns the application name (Laravel-like API).
	// This is an alias for GetName() for API consistency.
	//
	// Returns:
	//   string: The current application name
	Name() string

	// GetVersion returns the application version.
	//
	// Returns:
	//   string: The current application version
	GetVersion() string

	// SetVersion sets the application version.
	//
	// Parameters:
	//   version: The new application version
	SetVersion(version string)

	// Version returns the application version (Laravel-like API).
	// This is an alias for GetVersion() for API consistency.
	//
	// Returns:
	//   string: The current application version
	Version() string
}
