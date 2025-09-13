package enums

// Status represents the application status states.
// These constants define the possible states of a GoVel application during its lifecycle.
type Status string

const (
	// StatusUnbooted represents an application that hasn't been booted yet.
	// This is the initial state after application creation but before Boot() is called.
	StatusUnbooted Status = "unbooted"

	// StatusBooting represents an application that is currently in the boot process.
	// This state is active during the execution of Boot() method.
	StatusBooting Status = "booting"

	// StatusRunning represents a fully booted and running application.
	// This is the normal operational state after successful boot.
	StatusRunning Status = "running"

	// StatusShuttingDown represents an application that is in the shutdown process.
	// This state is active during graceful shutdown procedures.
	StatusShuttingDown Status = "shutting_down"

	// StatusShutdown represents an application that has completed shutdown.
	// This is the final state after all cleanup operations are complete.
	StatusShutdown Status = "shutdown"

	// StatusMaintenance represents an application in maintenance mode.
	// The application is running but not accepting normal requests.
	StatusMaintenance Status = "maintenance"

	// StatusError represents an application that encountered a critical error.
	// This state indicates the application cannot continue normal operation.
	StatusError Status = "error"
)

// String returns the string representation of the Status.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status value is one of the defined constants.
func (s Status) IsValid() bool {
	switch s {
	case StatusUnbooted, StatusBooting, StatusRunning, StatusShuttingDown,
		StatusShutdown, StatusMaintenance, StatusError:
		return true
	default:
		return false
	}
}

// IsOperational returns true if the application is in an operational state.
// Operational states are those where the application can handle requests.
func (s Status) IsOperational() bool {
	return s == StatusRunning
}

// IsTransitional returns true if the application is transitioning between states.
// Transitional states are temporary states during lifecycle changes.
func (s Status) IsTransitional() bool {
	return s == StatusBooting || s == StatusShuttingDown
}

// IsFinal returns true if the application is in a final state.
// Final states are those where no further transitions are expected.
func (s Status) IsFinal() bool {
	return s == StatusShutdown || s == StatusError
}