package interfaces

import "time"

// ApplicationTimingInterface defines the contract for components that manage
// application timing information such as start time and uptime.
// This interface follows the Interface Segregation Principle by focusing
// solely on timing-related operations.
type ApplicationTimingInterface interface {
	// GetStartTime returns when the application was started.
	//
	// Returns:
	//   time.Time: The time when the application was started
	GetStartTime() time.Time

	// SetStartTime sets when the application was started.
	//
	// Parameters:
	//   startTime: The time when the application was started
	SetStartTime(startTime time.Time)

	// GetUptime returns how long the application has been running.
	//
	// Returns:
	//   time.Duration: The duration since the application was started
	GetUptime() time.Duration
}
