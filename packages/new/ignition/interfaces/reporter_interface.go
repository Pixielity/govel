package interfaces

// ReporterInterface interface for different reporting backends
type ReporterInterface interface {
	// Report sends the error report to the reporting backend
	// Returns an error if the reporting fails
	Report(ErrorReportInterface) error
	
	// IsEnabled returns true if the reporter is enabled
	IsEnabled() bool
	
	// SetEnabled enables or disables the reporter
	SetEnabled(bool)
	
	// GetName returns the name of the reporter
	GetName() string
}
