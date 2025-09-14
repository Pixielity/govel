package concerns

import (
	"sync"

	"govel/application/helpers"
	concernsInterfaces "govel/types/interfaces/application/concerns"
)

/**
 * Identity provides application identity management functionality including
 * name and version handling. This trait implements the HasIdentityInterface
 * and manages application identification information.
 *
 * Features:
 * - Application name management with thread safety
 * - Version management and tracking
 * - Laravel-compatible API with aliases
 * - Thread-safe access to identity properties
 */
type HasIdentity struct {
	/**
	 * name holds the application name
	 */
	name string

	/**
	 * version holds the current application version
	 */
	version string

	/**
	 * mutex provides thread-safe access to identity fields
	 */
	mutex sync.RWMutex
}

/**
 * NewHasIdentity creates a new identity trait with optional parameters.
 * If values are not provided, they will be read from environment variables.
 *
 * Parameters:
 *   options[0]: Optional application name (string)
 *   options[1]: Optional application version (string)
 *
 * Returns:
 *   *HasIdentity: A new identity trait instance
 *
 * Example:
 *   // Using environment variables
 *   identity := NewIdentity()
 *   // Providing explicit values (name="MyApp", version="1.0.0")
 *   identity := NewIdentity("MyApp", "1.0.0")
 */
func NewIdentity(options ...string) *HasIdentity {
	envHelper := helpers.NewEnvHelper()

	// Use provided options or fallback to environment variables
	appName := envHelper.GetAppName()       // Default from environment
	appVersion := envHelper.GetAppVersion() // Default from environment

	// If options are provided:
	// First is name
	// Second is version
	if len(options) > 0 && options[0] != "" {
		appName = options[0]
	}
	if len(options) > 1 && options[1] != "" {
		appVersion = options[1]
	}

	return &HasIdentity{
		name:    appName,
		version: appVersion,
	}
}

// GetName returns the application name.
//
// Returns:
//
//	string: The current application name
func (i *HasIdentity) GetName() string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.name
}

// SetName sets the application name.
//
// Parameters:
//
//	name: The new application name
func (i *HasIdentity) SetName(name string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.name = name
}

// GetVersion returns the application version.
//
// Returns:
//
//	string: The current application version
func (i *HasIdentity) GetVersion() string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.version
}

// SetVersion sets the application version.
//
// Parameters:
//
//	version: The new application version
func (i *HasIdentity) SetVersion(version string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.version = version
}

// Compile-time interface compliance check
var _ concernsInterfaces.HasIdentityInterface = (*HasIdentity)(nil)
