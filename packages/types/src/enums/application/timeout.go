package enums

import "time"

// TimeoutType represents different types of timeouts in the GoVel application
type TimeoutType string

const (
	// TimeoutShutdown represents shutdown timeout
	TimeoutShutdown TimeoutType = "shutdown"

	// TimeoutDrain represents drain timeout
	TimeoutDrain TimeoutType = "drain"

	// TimeoutBoot represents application boot timeout
	TimeoutBoot TimeoutType = "boot"

	// TimeoutLifecycle represents lifecycle operation timeout
	TimeoutLifecycle TimeoutType = "lifecycle"

	// TimeoutTermination represents service provider termination timeout
	TimeoutTermination TimeoutType = "termination"

	// TimeoutHTTP represents HTTP client timeout
	TimeoutHTTP TimeoutType = "http"

	// TimeoutDatabase represents database operation timeout
	TimeoutDatabase TimeoutType = "database"

	// TimeoutCache represents cache operation timeout
	TimeoutCache TimeoutType = "cache"

	// TimeoutContext represents context operation timeout
	TimeoutContext TimeoutType = "context"
)

// TimeoutDuration represents different timeout durations
type TimeoutDuration time.Duration

const (
	// TimeoutVeryShort represents very short operations (1 second)
	TimeoutVeryShort TimeoutDuration = TimeoutDuration(1 * time.Second)

	// TimeoutShort represents short operations (5 seconds)
	TimeoutShort TimeoutDuration = TimeoutDuration(5 * time.Second)

	// TimeoutMedium represents medium operations (15 seconds)
	TimeoutMedium TimeoutDuration = TimeoutDuration(15 * time.Second)

	// TimeoutLong represents long operations (30 seconds)
	TimeoutLong TimeoutDuration = TimeoutDuration(30 * time.Second)

	// TimeoutVeryLong represents very long operations (60 seconds)
	TimeoutVeryLong TimeoutDuration = TimeoutDuration(60 * time.Second)

	// TimeoutExtended represents extended operations (5 minutes)
	TimeoutExtended TimeoutDuration = TimeoutDuration(300 * time.Second)
)

// DefaultTimeouts maps timeout types to their default durations
var DefaultTimeouts = map[TimeoutType]TimeoutDuration{
	TimeoutShutdown:    TimeoutLong,   // 30 seconds
	TimeoutDrain:       TimeoutShort,  // 5 seconds (was 10, but 5 is more appropriate)
	TimeoutBoot:        TimeoutLong,   // 30 seconds
	TimeoutLifecycle:   TimeoutLong,   // 30 seconds
	TimeoutTermination: TimeoutMedium, // 15 seconds (was 20)
	TimeoutHTTP:        TimeoutLong,   // 30 seconds
	TimeoutDatabase:    TimeoutMedium, // 15 seconds
	TimeoutCache:       TimeoutShort,  // 5 seconds
	TimeoutContext:     TimeoutMedium, // 15 seconds
}

// String returns the string representation of the timeout type
func (t TimeoutType) String() string {
	return string(t)
}

// Duration converts TimeoutDuration to time.Duration
func (td TimeoutDuration) Duration() time.Duration {
	return time.Duration(td)
}

// String returns the string representation of the timeout duration
func (td TimeoutDuration) String() string {
	return time.Duration(td).String()
}

// Seconds returns the duration in seconds
func (td TimeoutDuration) Seconds() float64 {
	return time.Duration(td).Seconds()
}

// GetDefaultTimeout returns the default timeout for a given timeout type
func GetDefaultTimeout(timeoutType TimeoutType) TimeoutDuration {
	if duration, exists := DefaultTimeouts[timeoutType]; exists {
		return duration
	}
	return TimeoutMedium // Default fallback
}

// IsValid checks if the timeout type is valid
func (t TimeoutType) IsValid() bool {
	_, exists := DefaultTimeouts[t]
	return exists
}

// AllTimeoutTypes returns all valid timeout types
func AllTimeoutTypes() []TimeoutType {
	types := make([]TimeoutType, 0, len(DefaultTimeouts))
	for timeoutType := range DefaultTimeouts {
		types = append(types, timeoutType)
	}
	return types
}

// TimeoutCategory represents categories of timeout durations
type TimeoutCategory string

const (
	// CategoryQuick represents quick operations
	CategoryQuick TimeoutCategory = "quick"

	// CategoryStandard represents standard operations
	CategoryStandard TimeoutCategory = "standard"

	// CategoryLengthy represents lengthy operations
	CategoryLengthy TimeoutCategory = "lengthy"
)

// String returns the string representation of the TimeoutCategory.
func (tc TimeoutCategory) String() string {
	return string(tc)
}

// GetTimeoutCategory returns the category for a given timeout duration
func GetTimeoutCategory(duration TimeoutDuration) TimeoutCategory {
	d := duration.Duration()
	switch {
	case d <= 10*time.Second:
		return CategoryQuick
	case d <= 30*time.Second:
		return CategoryStandard
	default:
		return CategoryLengthy
	}
}