package enums

// Environment represents the application environment
type Environment string

// Application environment constants
const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentTesting     Environment = "testing"
	EnvironmentStaging     Environment = "staging" 
	EnvironmentProduction  Environment = "production"
	EnvironmentLocal       Environment = "local"
)

// IsValid checks if the environment value is valid
func (e Environment) IsValid() bool {
	switch e {
	case EnvironmentDevelopment, EnvironmentTesting, EnvironmentStaging, EnvironmentProduction, EnvironmentLocal:
		return true
	default:
		return false
	}
}

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// IsDevelopment checks if the environment is development
func (e Environment) IsDevelopment() bool {
	return e == EnvironmentDevelopment
}

// IsTesting checks if the environment is testing
func (e Environment) IsTesting() bool {
	return e == EnvironmentTesting
}

// IsStaging checks if the environment is staging
func (e Environment) IsStaging() bool {
	return e == EnvironmentStaging
}

// IsProduction checks if the environment is production
func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

// IsLocal checks if the environment is local
func (e Environment) IsLocal() bool {
	return e == EnvironmentLocal
}