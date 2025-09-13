package constants

import "time"

// Application default constants
const (
	// DefaultName is the default application name
	DefaultName = "Govel Application"

	// DefaultVersion is the default application version
	DefaultVersion = "1.0.0"

	// DefaultDebug is the default debug mode
	DefaultDebug = false

	// DefaultLocale is the default application locale
	DefaultLocale = "en"

	// DefaultFallbackLocale is the default fallback locale
	DefaultFallbackLocale = "en"

	// DefaultTimezone is the default application timezone
	DefaultTimezone = "UTC"

	// DefaultRunningInConsole indicates if running in console mode by default
	DefaultRunningInConsole = false

	// DefaultRunningUnitTests indicates if running unit tests by default
	DefaultRunningUnitTests = false

	// DefaultShutdownTimeout is the default graceful shutdown timeout
	DefaultShutdownTimeout = 30 * time.Second
)