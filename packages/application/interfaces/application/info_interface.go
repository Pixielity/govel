package interfaces

// ApplicationInfoInterface defines the contract for components that provide
// comprehensive application information and introspection capabilities.
// This interface follows the Interface Segregation Principle by focusing
// solely on information retrieval operations.
type ApplicationInfoInterface interface {
	// GetApplicationInfo returns comprehensive application information.
	// This includes data from all traits and core application properties.
	//
	// Returns:
	//   map[string]interface{}: A map containing detailed application information
	//     including name, version, environment, runtime state, and trait data
	GetApplicationInfo() map[string]interface{}
}
