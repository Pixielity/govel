package traits

import (
	"fmt"
	"math"
	"time"

	carbon "govel/support/src/carbon"
)

// TimeInteractor interface defines methods for time-related operations
type TimeInteractor interface {
	SecondsUntil(delay interface{}) int64
	AvailableAt(delay ...interface{}) int64
	CurrentTime() int64
	RunTimeForHumans(startTime float64, endTime ...float64) string
}

// InteractsWithTime provides time-related functionality that can be embedded in structs
type InteractsWithTime struct{}

// SecondsUntil returns the number of seconds until the given DateTime or duration.
//
// Parameters:
//   - delay: can be *carbon.Carbon, time.Time, time.Duration, or int64 (seconds)
//
// Returns the number of seconds until the specified time, or 0 if the time is in the past.
func (t *InteractsWithTime) SecondsUntil(delay interface{}) int64 {
	parsed := t.parseDateInterval(delay)

	switch v := parsed.(type) {
	case *carbon.Carbon:
		return int64(math.Max(0, float64(v.Timestamp()-t.CurrentTime())))
	case time.Time:
		return int64(math.Max(0, float64(v.Unix()-t.CurrentTime())))
	case int64:
		return v
	case time.Duration:
		return int64(v.Seconds())
	default:
		return 0
	}
}

// AvailableAt returns the "available at" UNIX timestamp.
//
// Parameters:
//   - delay: optional delay parameter, defaults to 0 if not provided
//
// Returns the timestamp when something will be available.
func (t *InteractsWithTime) AvailableAt(delay ...interface{}) int64 {
	var delayValue interface{} = int64(0)
	if len(delay) > 0 {
		delayValue = delay[0]
	}

	parsed := t.parseDateInterval(delayValue)

	switch v := parsed.(type) {
	case *carbon.Carbon:
		return v.Timestamp()
	case time.Time:
		return v.Unix()
	case int64:
		return carbon.Now().AddSeconds(int(v)).Timestamp()
	case time.Duration:
		return carbon.Now().AddSeconds(int(v.Seconds())).Timestamp()
	default:
		return carbon.Now().Timestamp()
	}
}

// ParseDateInterval converts various time representations to a consistent format.
//
// Handles:
//   - *carbon.Carbon: returns as-is
//   - time.Time: returns as-is
//   - time.Duration: returns as-is
//   - int64: treats as seconds, returns as-is
//   - int: converts to int64 seconds
//   - string: attempts to parse as duration or timestamp
func (t *InteractsWithTime) parseDateInterval(delay interface{}) interface{} {
	switch v := delay.(type) {
	case *carbon.Carbon:
		return v
	case time.Time:
		return v
	case time.Duration:
		return v
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		// Try to parse as duration first
		if duration, err := time.ParseDuration(v); err == nil {
			return duration
		}
		// Try to parse as carbon timestamp
		parsed := carbon.Parse(v)
		if !parsed.IsZero() {
			return parsed
		}
		// Default to 0
		return int64(0)
	default:
		// For any other type, try to add it to current time if it looks like a duration
		return int64(0)
	}
}

// CurrentTime returns the current system time as a UNIX timestamp.
func (t *InteractsWithTime) CurrentTime() int64 {
	return carbon.Now().Timestamp()
}

// RunTimeForHumans formats the total run time between start and end times for human readability.
//
// Parameters:
//   - startTime: the start time as a float64 (typically from time.Now().UnixNano() / 1e9)
//   - endTime: optional end time, defaults to current time if not provided
//
// Returns a human-readable string representation of the runtime.
func (t *InteractsWithTime) RunTimeForHumans(startTime float64, endTime ...float64) string {
	var end float64
	if len(endTime) > 0 {
		end = endTime[0]
	} else {
		end = float64(time.Now().UnixNano()) / 1e9
	}

	runTime := (end - startTime) * 1000 // Convert to milliseconds

	if runTime > 1000 {
		// Convert to seconds and format appropriately
		seconds := runTime / 1000

		if seconds < 60 {
			return fmt.Sprintf("%.2fs", seconds)
		} else if seconds < 3600 {
			minutes := int(seconds / 60)
			remainingSeconds := seconds - float64(minutes*60)
			if remainingSeconds > 0 {
				return fmt.Sprintf("%dm %.1fs", minutes, remainingSeconds)
			}
			return fmt.Sprintf("%dm", minutes)
		} else {
			hours := int(seconds / 3600)
			remainingSeconds := seconds - float64(hours*3600)
			minutes := int(remainingSeconds / 60)
			remainingSeconds = remainingSeconds - float64(minutes*60)

			if remainingSeconds > 0 {
				return fmt.Sprintf("%dh %dm %.1fs", hours, minutes, remainingSeconds)
			} else if minutes > 0 {
				return fmt.Sprintf("%dh %dm", hours, minutes)
			}
			return fmt.Sprintf("%dh", hours)
		}
	}

	return fmt.Sprintf("%.2fms", runTime)
}

// Convenience methods for common time operations

// SecondsUntilNextMinute returns seconds until the next minute boundary
func (t *InteractsWithTime) SecondsUntilNextMinute() int64 {
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	return int64(nextMinute.Sub(now).Seconds())
}

// SecondsUntilNextHour returns seconds until the next hour boundary
func (t *InteractsWithTime) SecondsUntilNextHour() int64 {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(time.Hour)
	return int64(nextHour.Sub(now).Seconds())
}

// SecondsUntilNextDay returns seconds until the next day boundary
func (t *InteractsWithTime) SecondsUntilNextDay() int64 {
	now := carbon.Now()
	nextDay := now.Copy().AddDay().StartOfDay()
	return int64(nextDay.Timestamp() - now.Timestamp())
}

// SecondsUntilEndOfDay returns seconds until the end of the current day
func (t *InteractsWithTime) SecondsUntilEndOfDay() int64 {
	now := carbon.Now()
	endOfDay := now.Copy().EndOfDay()
	return int64(endOfDay.Timestamp() - now.Timestamp())
}

// AvailableAtNextMinute returns timestamp for the next minute boundary
func (t *InteractsWithTime) AvailableAtNextMinute() int64 {
	return t.AvailableAt(t.SecondsUntilNextMinute())
}

// AvailableAtNextHour returns timestamp for the next hour boundary
func (t *InteractsWithTime) AvailableAtNextHour() int64 {
	return t.AvailableAt(t.SecondsUntilNextHour())
}

// AvailableAtNextDay returns timestamp for the start of the next day
func (t *InteractsWithTime) AvailableAtNextDay() int64 {
	return t.AvailableAt(t.SecondsUntilNextDay())
}

// MicroTimeFloat returns the current time as a float64 (seconds.microseconds)
// Similar to PHP's microtime(true)
func (t *InteractsWithTime) MicroTimeFloat() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

// Timer provides a simple way to measure execution time
type Timer struct {
	InteractsWithTime
	startTime float64
}

// NewTimer creates a new timer and starts it
func NewTimer() *Timer {
	timer := &Timer{}
	timer.Start()
	return timer
}

// Start begins timing
func (timer *Timer) Start() {
	timer.startTime = timer.MicroTimeFloat()
}

// Stop returns the elapsed time as a human-readable string
func (timer *Timer) Stop() string {
	return timer.RunTimeForHumans(timer.startTime)
}

// Elapsed returns the elapsed time in seconds as float64
func (timer *Timer) Elapsed() float64 {
	return timer.MicroTimeFloat() - timer.startTime
}

// ElapsedMilliseconds returns the elapsed time in milliseconds
func (timer *Timer) ElapsedMilliseconds() float64 {
	return (timer.MicroTimeFloat() - timer.startTime) * 1000
}

// Helper functions for creating delays

// DelayUntil creates a delay until a specific carbon time
func DelayUntil(when *carbon.Carbon) interface{} {
	return when
}

// DelayFor creates a delay for a specific duration
func DelayFor(duration time.Duration) interface{} {
	return duration
}

// DelaySeconds creates a delay for a specific number of seconds
func DelaySeconds(seconds int64) interface{} {
	return seconds
}

// DelayMinutes creates a delay for a specific number of minutes
func DelayMinutes(minutes int64) interface{} {
	return minutes * 60
}

// DelayHours creates a delay for a specific number of hours
func DelayHours(hours int64) interface{} {
	return hours * 3600
}

// DelayDays creates a delay for a specific number of days
func DelayDays(days int64) interface{} {
	return days * 86400
}
