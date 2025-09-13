package enums

// Environment represents the application environment types.
// These constants define the standard environments that a GoVel application can run in.
type Environment string

const (
	// EnvironmentDevelopment represents the development environment.
	// Used for local development with debug features enabled.
	EnvironmentDevelopment Environment = "development"

	// EnvironmentTesting represents the testing environment.
	// Used for automated testing with test-specific configurations.
	EnvironmentTesting Environment = "testing"

	// EnvironmentStaging represents the staging environment.
	// Used for pre-production testing and integration testing.
	EnvironmentStaging Environment = "staging"

	// EnvironmentProduction represents the production environment.
	// Used for live production deployments with optimized settings.
	EnvironmentProduction Environment = "production"
)

// String returns the string representation of the Environment.
func (e Environment) String() string {
	return string(e)
}

// IsValid checks if the environment value is one of the defined constants.
func (e Environment) IsValid() bool {
	switch e {
	case EnvironmentDevelopment, EnvironmentTesting, EnvironmentStaging, EnvironmentProduction:
		return true
	default:
		return false
	}
}

// IsProduction returns true if this is the production environment.
func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

// IsDevelopment returns true if this is the development environment.
func (e Environment) IsDevelopment() bool {
	return e == EnvironmentDevelopment
}

// IsTesting returns true if this is the testing environment.
func (e Environment) IsTesting() bool {
	return e == EnvironmentTesting
}

// IsStaging returns true if this is the staging environment.
func (e Environment) IsStaging() bool {
	return e == EnvironmentStaging
}
