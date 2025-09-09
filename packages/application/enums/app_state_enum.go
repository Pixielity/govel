package enums

// AppState represents the current state of the GoVel application
type AppState string

const (
	// AppStateUninitialized indicates the application has not been initialized
	AppStateUninitialized AppState = "uninitialized"

	// AppStateRegistering indicates the application is registering service providers
	AppStateRegistering AppState = "registering"

	// AppStateBooting indicates the application is in the boot process
	AppStateBooting AppState = "booting"

	// AppStateRunning indicates the application is running normally
	AppStateRunning AppState = "running"

	// AppStateShuttingDown indicates the application is shutting down
	AppStateShuttingDown AppState = "shutting_down"

	// AppStateShutdown indicates the application has completed shutdown
	AppStateShutdown AppState = "shutdown"
)

// String returns the string representation of the application state
func (s AppState) String() string {
	return string(s)
}

// IsValid checks if the application state is valid
func (s AppState) IsValid() bool {
	switch s {
	case AppStateUninitialized, AppStateRegistering, AppStateBooting,
		AppStateRunning, AppStateShuttingDown, AppStateShutdown:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the current state can transition to the target state
func (s AppState) CanTransitionTo(target AppState) bool {
	validTransitions := map[AppState][]AppState{
		AppStateUninitialized: {AppStateRegistering},
		AppStateRegistering:   {AppStateBooting, AppStateShuttingDown},
		AppStateBooting:       {AppStateRunning, AppStateShuttingDown},
		AppStateRunning:       {AppStateShuttingDown},
		AppStateShuttingDown:  {AppStateShutdown},
		AppStateShutdown:      {}, // Terminal state
	}

	allowedStates, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == target {
			return true
		}
	}
	return false
}

// AllAppStates returns all valid application states
func AllAppStates() []AppState {
	return []AppState{
		AppStateUninitialized,
		AppStateRegistering,
		AppStateBooting,
		AppStateRunning,
		AppStateShuttingDown,
		AppStateShutdown,
	}
}
