package enums

// LifecycleState represents the lifecycle state of the application
type LifecycleState string

const (
	// StateInitializing indicates the application is being initialized
	StateInitializing LifecycleState = "initializing"

	// StateBooting indicates the application is in the boot process
	StateBooting LifecycleState = "booting"

	// StateBooted indicates the application has completed booting
	StateBooted LifecycleState = "booted"

	// StateStarting indicates the application is starting
	StateStarting LifecycleState = "starting"

	// StateRunning indicates the application is running normally
	StateRunning LifecycleState = "running"

	// StateStopping indicates the application is stopping
	StateStopping LifecycleState = "stopping"

	// StateStopped indicates the application has stopped
	StateStopped LifecycleState = "stopped"

	// StateTerminating indicates the application is terminating
	StateTerminating LifecycleState = "terminating"

	// StateTerminated indicates the application has terminated
	StateTerminated LifecycleState = "terminated"

	// StateMaintenance indicates the application is in maintenance mode
	StateMaintenance LifecycleState = "maintenance"

	// StateError indicates the application has encountered an error
	StateError LifecycleState = "error"

	// StateShuttingDown indicates the application is shutting down gracefully
	StateShuttingDown LifecycleState = "shutting_down"
)

// String returns the string representation of the lifecycle state
func (s LifecycleState) String() string {
	return string(s)
}

// IsValid checks if the lifecycle state is valid
func (s LifecycleState) IsValid() bool {
	switch s {
	case StateInitializing, StateBooting, StateBooted, StateStarting,
		StateRunning, StateStopping, StateStopped, StateTerminating,
		StateTerminated, StateMaintenance, StateError, StateShuttingDown:
		return true
	default:
		return false
	}
}

// IsOperational returns true if the application is in an operational state
func (s LifecycleState) IsOperational() bool {
	return s == StateRunning
}

// IsTransitional returns true if the application is in a transitional state
func (s LifecycleState) IsTransitional() bool {
	return s == StateBooting || s == StateStarting || s == StateStopping ||
		s == StateTerminating || s == StateShuttingDown
}

// IsFinal returns true if the application is in a final state
func (s LifecycleState) IsFinal() bool {
	return s == StateTerminated || s == StateError
}

// CanTransitionTo checks if the current state can transition to the target state
func (s LifecycleState) CanTransitionTo(target LifecycleState) bool {
	validTransitions := map[LifecycleState][]LifecycleState{
		StateInitializing: {StateBooting, StateError},
		StateBooting:      {StateBooted, StateError, StateShuttingDown},
		StateBooted:       {StateStarting, StateError, StateShuttingDown},
		StateStarting:     {StateRunning, StateError, StateShuttingDown},
		StateRunning:      {StateMaintenance, StateStopping, StateShuttingDown, StateError},
		StateMaintenance:  {StateRunning, StateStopping, StateShuttingDown, StateError},
		StateStopping:     {StateStopped, StateError},
		StateStopped:      {StateStarting, StateTerminating, StateShuttingDown},
		StateTerminating:  {StateTerminated, StateError},
		StateShuttingDown: {StateStopped, StateTerminated, StateError},
		StateTerminated:   {}, // Terminal state
		StateError:        {}, // Terminal state
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

// FromString converts a string to a LifecycleState with validation
func FromString(s string) (LifecycleState, bool) {
	state := LifecycleState(s)
	return state, state.IsValid()
}

// AllLifecycleStates returns all valid lifecycle states
func AllLifecycleStates() []LifecycleState {
	return []LifecycleState{
		StateInitializing,
		StateBooting,
		StateBooted,
		StateStarting,
		StateRunning,
		StateStopping,
		StateStopped,
		StateTerminating,
		StateTerminated,
		StateMaintenance,
		StateError,
		StateShuttingDown,
	}
}