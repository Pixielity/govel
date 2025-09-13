package concerns

import (
	"sync"
	"time"

	applicationInterfaces "govel/packages/types/src/interfaces/application"
)

/**
 * HasTiming provides application timing management functionality including
 * start time tracking and uptime calculation. This trait implements the HasTimingInterface
 * and manages application timing information for monitoring and lifecycle management.
 *
 * Features:
 * - Application start time tracking
 * - Uptime calculation and monitoring
 * - Thread-safe access to timing information
 * - Performance and lifecycle monitoring support
 */
type HasTiming struct {
	/**
	 * startTime records when the application was started
	 */
	startTime time.Time

	/**
	 * mutex provides thread-safe access to timing fields
	 */
	mutex sync.RWMutex
}

// NewTiming creates a new timing trait with an optional start time parameter.
// If start time is not provided, it will be initialized to zero time.
//
// Parameters:
//
//	startTime: Optional start time (variadic, first non-zero value used if provided)
//
// Returns:
//
//	*HasTiming: A new timing trait instance
//
// Example:
//
//	// Using default zero time
//	timing := NewTiming()
//	// Providing explicit start time
//	timing := NewTiming(time.Now())
func NewTiming(startTime ...time.Time) *HasTiming {
	// Use provided start time or fallback to zero time
	initTime := time.Time{} // Default to zero time

	if len(startTime) > 0 && !startTime[0].IsZero() {
		initTime = startTime[0]
	}

	return &HasTiming{
		startTime: initTime,
	}
}

// GetStartTime returns when the application was started.
//
// Returns:
//
//	time.Time: The time when the application was started
//
// Example:
//
//	startTime := app.GetStartTime()
//	if startTime.IsZero() {
//	    fmt.Println("Application not yet started")
//	} else {
//	    fmt.Printf("Started at: %s\n", startTime.Format(time.RFC3339))
//	}
func (t *HasTiming) GetStartTime() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.startTime
}

// SetStartTime sets when the application was started.
//
// Parameters:
//
//	startTime: The time when the application was started
//
// Example:
//
//	app.SetStartTime(time.Now()) // Mark application as started now
func (t *HasTiming) SetStartTime(startTime time.Time) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.startTime = startTime
}

// GetUptime returns how long the application has been running.
//
// Returns:
//
//	time.Duration: The duration since the application was started
//
// Example:
//
//	uptime := app.GetUptime()
//	if uptime > 0 {
//	    fmt.Printf("Application has been running for: %s\n", uptime)
//	} else {
//	    fmt.Println("Application not yet started or start time not set")
//	}
func (t *HasTiming) GetUptime() time.Duration {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if t.startTime.IsZero() {
		return 0
	}
	return time.Since(t.startTime)
}

// Compile-time interface compliance check
var _ applicationInterfaces.HasTimingInterface = (*HasTiming)(nil)
