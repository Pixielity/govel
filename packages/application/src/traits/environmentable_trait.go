package traits

import (
	"sync"

	"govel/packages/application/helpers"
	enums "govel/types/enums/application"
	traitInterfaces "govel/types/interfaces/application/traits"
)

/**
 * Environmentable provides environment management functionality through composition.
 * This struct implements the EnvironmentableInterface and contains its own environment data,
 * following the self-contained trait pattern.
 *
 * Unlike dependency injection, this trait owns and manages its own state
 * for environment, debug mode, and environment-specific configurations.
 */
type Environmentable struct {
	environment enums.Environment // Current environment using enum type
	debug       bool              // Debug mode state
	mutex       sync.RWMutex      // Thread safety for trait operations
}

/**
 * NewEnvironmentable creates and initializes a new Environmentable instance with optional parameters.
 * If values are not provided, they will be read from environment variables.
 *
 * Parameters:
 *   options[0]: Optional environment setting (string)
 *   options[1]: Optional debug state (bool)
 *
 * Returns:
 *   *Environmentable: A properly initialized environment trait
 *
 * Example:
 *   // Using environment variables
 *   env := NewEnvironmentable()
 *   // Providing explicit values (environment="development", debug=true)
 *   env := NewEnvironmentable("development", true)
 */
func NewEnvironmentable(options ...interface{}) *Environmentable {
	envHelper := helpers.NewEnvHelper()

	// Use provided options or fallback to environment variables
	appEnvStr := envHelper.GetAppEnvironment() // Default from environment
	appDebug := envHelper.GetAppDebug()        // Default from environment

	// Convert string to Environment enum, defaulting to development
	appEnv := enums.EnvironmentDevelopment
	if envEnum := enums.Environment(appEnvStr); envEnum.IsValid() {
		appEnv = envEnum
	}

	// If options are provided:
	// First can be environment (string or Environment enum)
	// Second is debug (bool)
	if len(options) > 0 {
		switch v := options[0].(type) {
		case string:
			if v != "" {
				if envEnum := enums.Environment(v); envEnum.IsValid() {
					appEnv = envEnum
				}
			}
		case enums.Environment:
			if v.IsValid() {
				appEnv = v
			}
		}
	}
	if len(options) > 1 {
		if debugBool, ok := options[1].(bool); ok {
			appDebug = debugBool
		}
	}

	return &Environmentable{
		environment: appEnv,
		debug:       appDebug,
	}
}

// GetEnvironment returns the current application environment as a string.
//
// Returns:
//
//	string: The current environment name
func (e *Environmentable) GetEnvironment() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.environment.String()
}

// GetEnvironmentEnum returns the current application environment as an enum.
//
// Returns:
//
//	enums.Environment: The current environment enum
func (e *Environmentable) GetEnvironmentEnum() enums.Environment {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.environment
}

// SetEnvironment sets the application environment.
//
// Parameters:
//
//	env: The environment name to set (string)
func (e *Environmentable) SetEnvironment(env string) {
	e.SetEnvironmentEnum(enums.Environment(env))
}

// SetEnvironmentEnum sets the application environment using an enum.
//
// Parameters:
//
//	env: The environment enum to set
func (e *Environmentable) SetEnvironmentEnum(env enums.Environment) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !env.IsValid() {
		// Default to development if invalid
		env = enums.EnvironmentDevelopment
	}

	e.environment = env

	// Automatically disable debug in production
	if env.IsProduction() {
		e.debug = false
	}
}

// IsProduction returns whether the application is running in production environment.
//
// Returns:
//
//	bool: true if environment is "production"
func (e *Environmentable) IsProduction() bool {
	return e.GetEnvironmentEnum().IsProduction()
}

// IsDevelopment returns whether the application is running in development environment.
//
// Returns:
//
//	bool: true if environment is "development"
func (e *Environmentable) IsDevelopment() bool {
	return e.GetEnvironmentEnum().IsDevelopment()
}

// IsTesting returns whether the application is running in testing environment.
//
// Returns:
//
//	bool: true if environment is "testing"
func (e *Environmentable) IsTesting() bool {
	return e.GetEnvironmentEnum().IsTesting()
}

// IsStaging returns whether the application is running in staging environment.
//
// Returns:
//
//	bool: true if environment is "staging"
func (e *Environmentable) IsStaging() bool {
	return e.GetEnvironmentEnum().IsStaging()
}

// IsDebug returns whether the application is running in debug mode.
//
// Returns:
//
//	bool: true if debug mode is enabled
func (e *Environmentable) IsDebug() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.debug
}

// SetDebug sets the debug mode state.
//
// Parameters:
//
//	debug: Whether to enable debug mode
func (e *Environmentable) SetDebug(debug bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.debug = debug
}

// EnableDebug enables debug mode.
func (e *Environmentable) EnableDebug() {
	e.SetDebug(true)
}

// DisableDebug disables debug mode.
func (e *Environmentable) DisableDebug() {
	e.SetDebug(false)
}

// IsEnvironment checks if the current environment matches the given environment.
//
// Parameters:
//
//	env: The environment to check against
//
// Returns:
//
//	bool: true if the current environment matches
func (e *Environmentable) IsEnvironment(env string) bool {
	return e.GetEnvironment() == env
}

// GetEnvironmentInfo returns comprehensive environment information.
//
// Returns:
//
//	map[string]interface{}: Environment details
func (e *Environmentable) GetEnvironmentInfo() map[string]interface{} {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return map[string]interface{}{
		"environment":      e.environment.String(),
		"environment_enum": e.environment,
		"debug":            e.debug,
		"is_production":    e.environment.IsProduction(),
		"is_development":   e.environment.IsDevelopment(),
		"is_testing":       e.environment.IsTesting(),
		"is_staging":       e.environment.IsStaging(),
	}
}

// IsValidEnvironment checks if an environment name is valid.
//
// Parameters:
//
//	env: The environment name to validate
//
// Returns:
//
//	bool: true if the environment is valid
func (e *Environmentable) IsValidEnvironment(env string) bool {
	return enums.Environment(env).IsValid()
}

// Compile-time interface compliance check
var _ traitInterfaces.EnvironmentableInterface = (*Environmentable)(nil)
