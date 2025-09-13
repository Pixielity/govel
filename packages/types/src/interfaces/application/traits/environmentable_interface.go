package interfaces

// EnvironmentableInterface defines the contract for environment management functionality.
type EnvironmentableInterface interface {
	// GetEnvironment returns the current application environment as a string
	GetEnvironment() string
	
	// SetEnvironment sets the application environment
	SetEnvironment(env string)
	
	// IsProduction returns whether the application is running in production environment
	IsProduction() bool
	
	// IsDevelopment returns whether the application is running in development environment
	IsDevelopment() bool
	
	// IsTesting returns whether the application is running in testing environment
	IsTesting() bool
	
	// IsStaging returns whether the application is running in staging environment
	IsStaging() bool
	
	// IsDebug returns whether the application is running in debug mode
	IsDebug() bool
	
	// SetDebug sets the debug mode state
	SetDebug(debug bool)
	
	// EnableDebug enables debug mode
	EnableDebug()
	
	// DisableDebug disables debug mode
	DisableDebug()
	
	// IsEnvironment checks if the current environment matches the given environment
	IsEnvironment(env string) bool
	
	// GetEnvironmentInfo returns comprehensive environment information
	GetEnvironmentInfo() map[string]interface{}
	
	// IsValidEnvironment checks if an environment name is valid
	IsValidEnvironment(env string) bool
}
