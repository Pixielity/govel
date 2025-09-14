package concerns

import (
	"sync"

	"govel/application/helpers"
	concernsInterfaces "govel/types/interfaces/application/concerns"
)

/**
 * Runtime provides application runtime state management functionality including
 * console mode and testing mode detection. This trait implements the HasRuntimeInterface
 * and manages runtime execution context information.
 *
 * Features:
 * - Console mode detection and management
 * - Unit testing mode detection and management
 * - Thread-safe access to runtime state
 * - Runtime context awareness for conditional behavior
 */
type HasRuntime struct {
	/**
	 * runningInConsole indicates whether the application is running in console mode
	 */
	runningInConsole bool

	/**
	 * runningUnitTests indicates whether the application is running unit tests
	 */
	runningUnitTests bool

	/**
	 * mutex provides thread-safe access to runtime fields
	 */
	mutex sync.RWMutex
}

// NewHasRuntime creates a new runtime trait with optional mode parameters.
// If values are not provided, they will be read from environment variables.
//
// Parameters:
//
//	options: Optional boolean flags [console, testing]. If provided:
//	        - First value sets console mode
//	        - Second value sets testing mode
//
// Returns:
//
//	*HasRuntime: A new runtime trait instance
//
// Example:
//
//	// Using environment variables
//	runtime := NewRuntime()
//	// Providing explicit values (console=true, testing=false)
//	runtime := NewRuntime(true, false)
func NewRuntime(options ...bool) *HasRuntime {
	envHelper := helpers.NewEnvHelper()

	// Use provided options or fallback to environment
	isConsole := envHelper.GetRunningInConsole() // Default from environment
	isTesting := envHelper.GetRunningUnitTests() // Default from environment

	// If options are provided, first is console mode, second is testing mode
	if len(options) > 0 {
		isConsole = options[0]
	}
	if len(options) > 1 {
		isTesting = options[1]
	}

	return &HasRuntime{
		runningInConsole: isConsole,
		runningUnitTests: isTesting,
	}
}

// IsRunningInConsole returns whether the application is running in console mode.
//
// Returns:
//
//	bool: true if running in console mode, false otherwise
//
// Example:
//
//	if app.IsRunningInConsole() {
//	    fmt.Println("Running CLI command")
//	} else {
//	    fmt.Println("Running web request")
//	}
func (r *HasRuntime) IsRunningInConsole() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.runningInConsole
}

// SetRunningInConsole sets whether the application is running in console mode.
//
// Parameters:
//
//	console: true if running in console mode, false otherwise
//
// Example:
//
//	app.SetRunningInConsole(true) // Mark as console mode
func (r *HasRuntime) SetRunningInConsole(console bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.runningInConsole = console
}

// IsRunningUnitTests returns whether the application is running unit tests.
//
// Returns:
//
//	bool: true if running unit tests, false otherwise
//
// Example:
//
//	if app.IsRunningUnitTests() {
//	    // Use test database
//	    config.SetDatabase("test_db")
//	} else {
//	    // Use production/development database
//	    config.SetDatabase("main_db")
//	}
func (r *HasRuntime) IsRunningUnitTests() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.runningUnitTests
}

// SetRunningUnitTests sets whether the application is running unit tests.
//
// Parameters:
//
//	testing: true if running unit tests, false otherwise
//
// Example:
//
//	app.SetRunningUnitTests(true) // Mark as test mode
func (r *HasRuntime) SetRunningUnitTests(testing bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.runningUnitTests = testing
}

// Compile-time interface compliance check
var _ concernsInterfaces.HasRuntimeInterface = (*HasRuntime)(nil)
