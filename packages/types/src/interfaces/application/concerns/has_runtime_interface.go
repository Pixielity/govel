package interfaces

// HasRuntimeInterface defines the contract for application runtime state management functionality.
type HasRuntimeInterface interface {
	// IsRunningInConsole returns whether the application is running in console mode
	IsRunningInConsole() bool
	
	// SetRunningInConsole sets whether the application is running in console mode
	SetRunningInConsole(console bool)
	
	// IsRunningUnitTests returns whether the application is running unit tests
	IsRunningUnitTests() bool
	
	// SetRunningUnitTests sets whether the application is running unit tests
	SetRunningUnitTests(testing bool)
}