package interfaces

// ApplicationRuntimeInterface defines the contract for components that manage
// application runtime state information such as console mode and testing mode.
// This interface follows the Interface Segregation Principle by focusing
// solely on runtime state operations.
type ApplicationRuntimeInterface interface {
	// IsRunningInConsole returns whether the application is running in console mode.
	//
	// Returns:
	//   bool: true if running in console mode, false otherwise
	IsRunningInConsole() bool

	// SetRunningInConsole sets whether the application is running in console mode.
	//
	// Parameters:
	//   console: true if running in console mode, false otherwise
	SetRunningInConsole(console bool)

	// IsRunningUnitTests returns whether the application is running unit tests.
	//
	// Returns:
	//   bool: true if running unit tests, false otherwise
	IsRunningUnitTests() bool

	// SetRunningUnitTests sets whether the application is running unit tests.
	//
	// Parameters:
	//   testing: true if running unit tests, false otherwise
	SetRunningUnitTests(testing bool)
}
