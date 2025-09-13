// Package checks provides built-in health check implementations.
package checks

import (
	"fmt"
	"os"

	"govel/healthcheck/checks"
	"govel/healthcheck/enums"
	"govel/healthcheck/interfaces"
)

// EnvironmentCheck validates the application environment.
// It closely mirrors the Laravel health EnvironmentCheck pattern.
type EnvironmentCheck struct {
	*checks.BaseCheck

	// expectedEnvironment is the expected environment name (e.g., "production", "development")
	expectedEnvironment string
}

// NewEnvironmentCheck creates a new environment check with default settings.
//
// Returns:
//
//	*EnvironmentCheck: A new environment check instance
func NewEnvironmentCheck() *EnvironmentCheck {
	return &EnvironmentCheck{
		BaseCheck:           checks.NewBaseCheck(),
		expectedEnvironment: "production",
	}
}

// ExpectEnvironment sets the expected environment name.
//
// Parameters:
//
//	expectedEnvironment: The expected environment name
//
// Returns:
//
//	*EnvironmentCheck: Self for method chaining
func (ec *EnvironmentCheck) ExpectEnvironment(expectedEnvironment string) *EnvironmentCheck {
	ec.expectedEnvironment = expectedEnvironment
	return ec
}

// Run performs the environment health check.
//
// Returns:
//
//	interfaces.ResultInterface: The health check result
func (ec *EnvironmentCheck) Run() interfaces.ResultInterface {
	actualEnvironment := ec.environment()

	result := checks.NewResult().
		SetMeta(map[string]interface{}{
			"actual":   actualEnvironment,
			"expected": ec.expectedEnvironment,
		}).
		SetShortSummary(actualEnvironment)

	if ec.expectedEnvironment == actualEnvironment {
		return result.SetStatus(enums.StatusOK)
	}

	return result.
		SetStatus(enums.StatusFailed).
		SetNotificationMessage(fmt.Sprintf("The environment was expected to be `%s`, but actually was `%s`", ec.expectedEnvironment, actualEnvironment))
}

// environment returns the current application environment.
// This mirrors Laravel's app()->environment() method.
func (ec *EnvironmentCheck) environment() string {
	// Check APP_ENV first (Laravel standard)
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}

	// Fall back to other common environment variables
	envVars := []string{"ENVIRONMENT", "ENV", "GO_ENV"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			return value
		}
	}

	// Default to "production" if no environment variable is set
	return "production"
}

// Compile-time interface compliance check
var _ interfaces.CheckInterface = (*EnvironmentCheck)(nil)
