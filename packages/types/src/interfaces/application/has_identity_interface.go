package interfaces

// HasIdentityInterface defines the contract for application identity management functionality.
type HasIdentityInterface interface {
	// GetName returns the application name
	GetName() string
	
	// SetName sets the application name
	SetName(name string)
	
	// GetVersion returns the application version
	GetVersion() string
	
	// SetVersion sets the application version
	SetVersion(version string)
}