package enums

// Priority represents priority levels for various operations in the GoVel framework.
// These constants provide standardized priority values for consistent ordering.
type Priority int

const (
	// PrioritySystem represents the highest priority for system-level operations.
	// Used for critical system components that must be initialized first.
	PrioritySystem Priority = 0

	// PriorityCore represents high priority for core framework components.
	// Used for essential framework services like configuration, logging, etc.
	PriorityCore Priority = 10

	// PriorityInfrastructure represents priority for infrastructure services.
	// Used for database connections, cache systems, message queues, etc.
	PriorityInfrastructure Priority = 50

	// PriorityFramework represents priority for framework services.
	// Used for routing, middleware, authentication, etc.
	PriorityFramework Priority = 100

	// PriorityApplication represents priority for application-specific services.
	// Used for business logic services, controllers, etc.
	PriorityApplication Priority = 200

	// PriorityExtensions represents priority for extensions and plugins.
	// Used for optional features and third-party integrations.
	PriorityExtensions Priority = 300

	// PriorityLow represents low priority for non-critical operations.
	// Used for background tasks, cleanup operations, etc.
	PriorityLow Priority = 500

	// PriorityDefault represents the default priority when none is specified.
	// This is typically used as a fallback value.
	PriorityDefault Priority = PriorityApplication
)

// String returns the string representation of the Priority.
func (p Priority) String() string {
	switch p {
	case PrioritySystem:
		return "system"
	case PriorityCore:
		return "core"
	case PriorityInfrastructure:
		return "infrastructure"
	case PriorityFramework:
		return "framework"
	case PriorityApplication:
		return "application"
	case PriorityExtensions:
		return "extensions"
	case PriorityLow:
		return "low"
	default:
		return "custom"
	}
}

// IsValid checks if the priority value is within reasonable bounds.
// Valid priorities are typically between 0 and 1000.
func (p Priority) IsValid() bool {
	return p >= 0 && p <= 1000
}

// IsHigherThan returns true if this priority is higher than the other.
// Lower numeric values represent higher priority.
func (p Priority) IsHigherThan(other Priority) bool {
	return p < other
}

// IsLowerThan returns true if this priority is lower than the other.
// Higher numeric values represent lower priority.
func (p Priority) IsLowerThan(other Priority) bool {
	return p > other
}

// Equal returns true if this priority equals the other priority.
func (p Priority) Equal(other Priority) bool {
	return p == other
}