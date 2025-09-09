package helpers

import (
	"os"
	"strconv"
	"time"

	"govel/packages/application/constants"
	"govel/packages/application/enums"
)

// EnvHelper provides utilities for reading environment variables
// with intelligent fallbacks to constants and enums.
type EnvHelper struct{}

// NewEnvHelper creates a new environment helper instance
func NewEnvHelper() *EnvHelper {
	return &EnvHelper{}
}

// GetAppName returns the application name from APP_NAME env var,
// falling back to the default constant
func (e *EnvHelper) GetAppName() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return constants.DefaultName
}

// GetAppVersion returns the application version from APP_VERSION env var,
// falling back to the default constant
func (e *EnvHelper) GetAppVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return constants.DefaultVersion
}

// GetAppEnvironment returns the application environment from APP_ENV env var,
// falling back to the default enum value
func (e *EnvHelper) GetAppEnvironment() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		// Validate that it's a known environment
		envEnum := enums.Environment(env)
		if envEnum.IsValid() {
			return string(envEnum)
		}
	}
	return string(enums.EnvironmentDevelopment)
}

// GetAppDebug returns the debug mode from APP_DEBUG env var,
// falling back to the default constant
func (e *EnvHelper) GetAppDebug() bool {
	if debug := os.Getenv("APP_DEBUG"); debug != "" {
		return e.parseBool(debug, constants.DefaultDebug)
	}
	return constants.DefaultDebug
}

// GetAppLocale returns the application locale from APP_LOCALE env var,
// falling back to the default constant
func (e *EnvHelper) GetAppLocale() string {
	if locale := os.Getenv("APP_LOCALE"); locale != "" {
		return locale
	}
	return constants.DefaultLocale
}

// GetAppFallbackLocale returns the fallback locale from APP_FALLBACK_LOCALE env var,
// falling back to the default constant
func (e *EnvHelper) GetAppFallbackLocale() string {
	if locale := os.Getenv("APP_FALLBACK_LOCALE"); locale != "" {
		return locale
	}
	return constants.DefaultFallbackLocale
}

// GetAppTimezone returns the application timezone from APP_TIMEZONE env var,
// falling back to the default constant
func (e *EnvHelper) GetAppTimezone() string {
	if timezone := os.Getenv("APP_TIMEZONE"); timezone != "" {
		return timezone
	}
	return constants.DefaultTimezone
}

// GetRunningInConsole returns whether running in console mode from APP_CONSOLE env var,
// falling back to the default constant
func (e *EnvHelper) GetRunningInConsole() bool {
	if console := os.Getenv("APP_CONSOLE"); console != "" {
		return e.parseBool(console, constants.DefaultRunningInConsole)
	}
	return constants.DefaultRunningInConsole
}

// GetRunningUnitTests returns whether running unit tests from APP_TESTING env var,
// falling back to the default constant
func (e *EnvHelper) GetRunningUnitTests() bool {
	if testing := os.Getenv("APP_TESTING"); testing != "" {
		return e.parseBool(testing, constants.DefaultRunningUnitTests)
	}
	return constants.DefaultRunningUnitTests
}

// GetShutdownTimeout returns the shutdown timeout as time.Duration from APP_SHUTDOWN_TIMEOUT env var,
// falling back to the default enum timeout
func (e *EnvHelper) GetShutdownTimeout() time.Duration {
	if timeout := os.Getenv("APP_SHUTDOWN_TIMEOUT"); timeout != "" {
		if parsed, err := strconv.Atoi(timeout); err == nil && parsed > 0 {
			return time.Duration(parsed) * time.Second
		}
	}
	return enums.GetDefaultTimeout(enums.TimeoutShutdown).Duration()
}

// parseBool parses common boolean representations with a fallback
func (e *EnvHelper) parseBool(value string, fallback bool) bool {
	switch value {
	case "true", "1", "yes", "on", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		return fallback
	}
}

// Environment variable names used by the helper
const (
	EnvAppName            = "APP_NAME"
	EnvAppVersion         = "APP_VERSION"
	EnvAppEnvironment     = "APP_ENV"
	EnvAppDebug           = "APP_DEBUG"
	EnvAppLocale          = "APP_LOCALE"
	EnvAppFallbackLocale  = "APP_FALLBACK_LOCALE"
	EnvAppTimezone        = "APP_TIMEZONE"
	EnvAppConsole         = "APP_CONSOLE"
	EnvAppTesting         = "APP_TESTING"
	EnvAppShutdownTimeout = "APP_SHUTDOWN_TIMEOUT"
)

// GetAllAppDefaults returns a complete set of application defaults
// from environment variables with fallbacks to constants/enums
func GetAllAppDefaults() map[string]interface{} {
	helper := NewEnvHelper()

	return map[string]interface{}{
		"name":               helper.GetAppName(),
		"version":            helper.GetAppVersion(),
		"environment":        helper.GetAppEnvironment(),
		"debug":              helper.GetAppDebug(),
		"locale":             helper.GetAppLocale(),
		"fallback_locale":    helper.GetAppFallbackLocale(),
		"timezone":           helper.GetAppTimezone(),
		"running_in_console": helper.GetRunningInConsole(),
		"running_unit_tests": helper.GetRunningUnitTests(),
		"shutdown_timeout":   helper.GetShutdownTimeout(),
	}
}
