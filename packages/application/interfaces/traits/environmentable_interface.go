package interfaces

/**
 * EnvironmentableInterface defines the contract for components that provide
 * environment management functionality. This interface follows the Interface
 * Segregation Principle by focusing solely on environment-related operations.
 */
type EnvironmentableInterface interface {
	/**
	 * GetEnvironment returns the current application environment.
	 *
	 * @return string The current environment name
	 */
	GetEnvironment() string

	/**
	 * SetEnvironment sets the application environment.
	 *
	 * @param env string The environment name to set
	 */
	SetEnvironment(env string)

	/**
	 * IsProduction returns whether the application is running in production environment.
	 *
	 * @return bool true if environment is "production"
	 */
	IsProduction() bool

	/**
	 * IsDevelopment returns whether the application is running in development environment.
	 *
	 * @return bool true if environment is "development"
	 */
	IsDevelopment() bool

	/**
	 * IsTesting returns whether the application is running in testing environment.
	 *
	 * @return bool true if environment is "testing"
	 */
	IsTesting() bool

	/**
	 * IsStaging returns whether the application is running in staging environment.
	 *
	 * @return bool true if environment is "staging"
	 */
	IsStaging() bool

	/**
	 * IsDebug returns whether the application is running in debug mode.
	 *
	 * @return bool true if debug mode is enabled
	 */
	IsDebug() bool

	/**
	 * SetDebug sets the debug mode state.
	 *
	 * @param debug bool Whether to enable debug mode
	 */
	SetDebug(debug bool)

	/**
	 * EnableDebug enables debug mode.
	 */
	EnableDebug()

	/**
	 * DisableDebug disables debug mode.
	 */
	DisableDebug()

	/**
	 * IsEnvironment checks if the current environment matches the given environment.
	 *
	 * @param env string The environment to check against
	 * @return bool true if the current environment matches
	 */
	IsEnvironment(env string) bool

	/**
	 * GetEnvironmentInfo returns comprehensive environment information.
	 *
	 * @return map[string]interface{} Environment details
	 */
	GetEnvironmentInfo() map[string]interface{}

	/**
	 * IsValidEnvironment checks if an environment name is valid.
	 *
	 * @param env string The environment name to validate
	 * @return bool true if the environment is valid
	 */
	IsValidEnvironment(env string) bool
}
