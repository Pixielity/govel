package interfaces

import "time"

// HasTimingInterface defines the contract for application timing management functionality.
type HasTimingInterface interface {
	// GetStartTime returns when the application was started
	GetStartTime() time.Time
	
	// SetStartTime sets when the application was started
	SetStartTime(startTime time.Time)
	
	// GetUptime returns how long the application has been running
	GetUptime() time.Duration
}