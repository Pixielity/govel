package traits

import (
	"sync"

	"govel/packages/application/constants"
	traitInterfaces "govel/packages/application/interfaces/traits"
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
	environment string       // Current environment (e.g., "development", "production")
	debug       bool         // Debug mode state
	mutex       sync.RWMutex // Thread safety for trait operations
}

/**
 * NewEnvironmentable creates and initializes a new Environmentable instance.
 *
 * @param environment string The initial environment
 * @param debug bool The initial debug state
 * @return *Environmentable A properly initialized environment trait
 */
func NewEnvironmentable(environment string, debug bool) *Environmentable {
	return &Environmentable{
		environment: environment,
		debug:       debug,
	}
}

/**
 * GetEnvironment returns the current application environment.
 *
 * @return string The current environment name
 */
func (e *Environmentable) GetEnvironment() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.environment
}

/**
 * SetEnvironment sets the application environment.
 *
 * @param env string The environment name to set
 */
func (e *Environmentable) SetEnvironment(env string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.environment = env

	// Automatically disable debug in production
	if env == constants.ProductionEnvironment {
		e.debug = false
	}
}

/**
 * IsProduction returns whether the application is running in production environment.
 *
 * @return bool true if environment is "production"
 */
func (e *Environmentable) IsProduction() bool {
	return e.GetEnvironment() == constants.ProductionEnvironment
}

/**
 * IsDevelopment returns whether the application is running in development environment.
 *
 * @return bool true if environment is "development"
 */
func (e *Environmentable) IsDevelopment() bool {
	return e.GetEnvironment() == constants.DevelopmentEnvironment
}

/**
 * IsTesting returns whether the application is running in testing environment.
 *
 * @return bool true if environment is "testing"
 */
func (e *Environmentable) IsTesting() bool {
	return e.GetEnvironment() == constants.TestingEnvironment
}

/**
 * IsStaging returns whether the application is running in staging environment.
 *
 * @return bool true if environment is "staging"
 */
func (e *Environmentable) IsStaging() bool {
	return e.GetEnvironment() == constants.StagingEnvironment
}

/**
 * IsDebug returns whether the application is running in debug mode.
 *
 * @return bool true if debug mode is enabled
 */
func (e *Environmentable) IsDebug() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.debug
}

/**
 * SetDebug sets the debug mode state.
 *
 * @param debug bool Whether to enable debug mode
 */
func (e *Environmentable) SetDebug(debug bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.debug = debug
}

/**
 * EnableDebug enables debug mode.
 */
func (e *Environmentable) EnableDebug() {
	e.SetDebug(true)
}

/**
 * DisableDebug disables debug mode.
 */
func (e *Environmentable) DisableDebug() {
	e.SetDebug(false)
}

/**
 * IsEnvironment checks if the current environment matches the given environment.
 *
 * @param env string The environment to check against
 * @return bool true if the current environment matches
 */
func (e *Environmentable) IsEnvironment(env string) bool {
	return e.GetEnvironment() == env
}

/**
 * GetEnvironmentInfo returns comprehensive environment information.
 *
 * @return map[string]interface{} Environment details
 */
func (e *Environmentable) GetEnvironmentInfo() map[string]interface{} {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return map[string]interface{}{
		"environment":    e.environment,
		"debug":          e.debug,
		"is_production":  e.environment == constants.ProductionEnvironment,
		"is_development": e.environment == constants.DevelopmentEnvironment,
		"is_testing":     e.environment == constants.TestingEnvironment,
		"is_staging":     e.environment == constants.StagingEnvironment,
	}
}

/**
 * IsValidEnvironment checks if an environment name is valid.
 *
 * @param env string The environment name to validate
 * @return bool true if the environment is valid
 */
func (e *Environmentable) IsValidEnvironment(env string) bool {
	validEnvironments := []string{
		constants.DevelopmentEnvironment,
		constants.TestingEnvironment,
		constants.StagingEnvironment,
		constants.ProductionEnvironment,
		constants.LocalEnvironment,
	}

	for _, validEnv := range validEnvironments {
		if env == validEnv {
			return true
		}
	}
	return false
}

// Compile-time interface compliance check
// This ensures that Environmentable implements the EnvironmentableInterface at compile time
var _ traitInterfaces.EnvironmentableInterface = (*Environmentable)(nil)
