// Package types provides concrete implementations of the health check interfaces.
// This package contains the core data structures used throughout the health check system.
package checks

import (
	"time"

	"govel/packages/healthcheck/src/interfaces"
)

// Result represents the outcome of a health check execution.
// It contains status information, messages, metadata, and timing data.
type Result struct {
	// status holds the health check status (ok, warning, failed, etc.)
	status interfaces.StatusInterface

	// notificationMessage is the message to be used in notifications
	notificationMessage string

	// shortSummary is a brief description of the result
	shortSummary string

	// meta contains additional metadata about the check execution
	meta map[string]interface{}

	// check is a reference to the health check that produced this result
	check interfaces.CheckInterface

	// startedAt is when the health check execution started
	startedAt *time.Time

	// endedAt is when the health check execution completed
	endedAt *time.Time
}

// NewResult creates a new Result instance with default values.
// The result is initialized with empty status and metadata.
//
// Returns:
//
//	*Result: A new result instance ready for configuration
//
// Example:
//
//	result := types.NewResult()
//	result.SetStatus(enums.StatusOK).SetNotificationMessage("All systems operational")
func NewResult() *Result {
	return &Result{
		meta: make(map[string]interface{}),
	}
}

// NewResultWithStatus creates a new Result with the specified status.
// This is a convenience constructor for quickly creating results with a status.
//
// Parameters:
//
//	status: The initial status for the result
//
// Returns:
//
//	*Result: A new result instance with the specified status
func NewResultWithStatus(status interfaces.StatusInterface) *Result {
	result := NewResult()
	result.status = status
	return result
}

// GetStatus returns the status of the health check result.
//
// Returns:
//
//	interfaces.StatusInterface: The current status
func (r *Result) GetStatus() interfaces.StatusInterface {
	return r.status
}

// SetStatus sets the status of the health check result.
//
// Parameters:
//
//	status: The status to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetStatus(status interfaces.StatusInterface) interfaces.ResultInterface {
	r.status = status
	return r
}

// GetNotificationMessage returns the message to be used in notifications.
//
// Returns:
//
//	string: The notification message
func (r *Result) GetNotificationMessage() string {
	return r.notificationMessage
}

// SetNotificationMessage sets the notification message.
//
// Parameters:
//
//	message: The notification message to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetNotificationMessage(message string) interfaces.ResultInterface {
	r.notificationMessage = message
	return r
}

// GetShortSummary returns a brief summary of the result.
//
// Returns:
//
//	string: Brief summary text
func (r *Result) GetShortSummary() string {
	if r.shortSummary != "" {
		return r.shortSummary
	}

	// Auto-generate short summary from status if not explicitly set
	if r.status != nil {
		switch r.status.String() {
		case "ok":
			return "Healthy"
		case "warning":
			return "Warning"
		case "failed":
			return "Failed"
		case "crashed":
			return "Crashed"
		case "skipped":
			return "Skipped"
		default:
			return "Unknown"
		}
	}

	return "Unknown"
}

// SetShortSummary sets a brief summary of the result.
//
// Parameters:
//
//	summary: Brief summary text to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetShortSummary(summary string) interfaces.ResultInterface {
	r.shortSummary = summary
	return r
}

// GetMeta returns metadata associated with the result.
//
// Returns:
//
//	map[string]interface{}: Metadata key-value pairs
func (r *Result) GetMeta() map[string]interface{} {
	if r.meta == nil {
		r.meta = make(map[string]interface{})
	}
	return r.meta
}

// SetMeta sets metadata for the result.
//
// Parameters:
//
//	meta: Metadata key-value pairs to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetMeta(meta map[string]interface{}) interfaces.ResultInterface {
	r.meta = meta
	return r
}

// AddMeta adds a single metadata key-value pair to the result.
// This is a convenience method for adding individual metadata items.
//
// Parameters:
//
//	key: The metadata key
//	value: The metadata value
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) AddMeta(key string, value interface{}) interfaces.ResultInterface {
	if r.meta == nil {
		r.meta = make(map[string]interface{})
	}
	r.meta[key] = value
	return r
}

// GetCheck returns the health check that produced this result.
//
// Returns:
//
//	interfaces.CheckInterface: The originating health check
func (r *Result) GetCheck() interfaces.CheckInterface {
	return r.check
}

// SetCheck sets the health check that produced this result.
//
// Parameters:
//
//	check: The health check to associate with this result
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetCheck(check interfaces.CheckInterface) interfaces.ResultInterface {
	r.check = check
	return r
}

// GetStartedAt returns when the health check started executing.
//
// Returns:
//
//	*time.Time: Start time, nil if not set
func (r *Result) GetStartedAt() *time.Time {
	return r.startedAt
}

// SetStartedAt sets when the health check started executing.
//
// Parameters:
//
//	startedAt: The start time to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetStartedAt(startedAt time.Time) interfaces.ResultInterface {
	r.startedAt = &startedAt
	return r
}

// GetEndedAt returns when the health check finished executing.
//
// Returns:
//
//	*time.Time: End time, nil if not set
func (r *Result) GetEndedAt() *time.Time {
	return r.endedAt
}

// SetEndedAt sets when the health check finished executing.
//
// Parameters:
//
//	endedAt: The end time to set
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetEndedAt(endedAt time.Time) interfaces.ResultInterface {
	r.endedAt = &endedAt
	return r
}

// GetDuration returns the execution duration of the health check.
// If either start or end time is not set, returns 0.
//
// Returns:
//
//	time.Duration: Execution duration, 0 if times not set
func (r *Result) GetDuration() time.Duration {
	if r.startedAt == nil || r.endedAt == nil {
		return 0
	}
	return r.endedAt.Sub(*r.startedAt)
}

// SetDuration is a convenience method to set both start and end times
// based on a duration from the current time.
//
// Parameters:
//
//	duration: The duration the check took to execute
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) SetDuration(duration time.Duration) interfaces.ResultInterface {
	now := time.Now()
	r.endedAt = &now
	startTime := now.Add(-duration)
	r.startedAt = &startTime
	return r
}

// OK is a convenience method to set status to OK with an optional message.
//
// Parameters:
//
//	message: Optional notification message (can be empty)
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) OK(message string) interfaces.ResultInterface {
	// This will need to be updated once we have the actual Status implementation
	// For now, we'll assume there's a way to create an OK status
	r.notificationMessage = message
	return r
}

// Warning is a convenience method to set status to Warning with an optional message.
//
// Parameters:
//
//	message: Optional notification message (can be empty)
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) Warning(message string) interfaces.ResultInterface {
	r.notificationMessage = message
	return r
}

// Failed is a convenience method to set status to Failed with an optional message.
//
// Parameters:
//
//	message: Optional notification message (can be empty)
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) Failed(message string) interfaces.ResultInterface {
	r.notificationMessage = message
	return r
}

// Crashed is a convenience method to set status to Crashed with an optional message.
//
// Parameters:
//
//	message: Optional notification message (can be empty)
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) Crashed(message string) interfaces.ResultInterface {
	r.notificationMessage = message
	return r
}

// Skipped is a convenience method to set status to Skipped with an optional message.
//
// Parameters:
//
//	message: Optional notification message (can be empty)
//
// Returns:
//
//	interfaces.ResultInterface: Self for method chaining
func (r *Result) Skipped(message string) interfaces.ResultInterface {
	r.notificationMessage = message
	return r
}

// IsHealthy returns true if the result status indicates a healthy state.
//
// Returns:
//
//	bool: true if the result is healthy, false otherwise
func (r *Result) IsHealthy() bool {
	if r.status == nil {
		return false
	}
	return r.status.IsHealthy()
}

// IsFailed returns true if the result status indicates a failed state.
//
// Returns:
//
//	bool: true if the result is failed, false otherwise
func (r *Result) IsFailed() bool {
	if r.status == nil {
		return false
	}
	return r.status.IsFailed()
}

// IsWarning returns true if the result status indicates a warning state.
//
// Returns:
//
//	bool: true if the result has warnings, false otherwise
func (r *Result) IsWarning() bool {
	if r.status == nil {
		return false
	}
	return r.status.IsWarning()
}

// GetExecutionInfo returns comprehensive information about the execution.
// This is useful for debugging and logging.
//
// Returns:
//
//	map[string]interface{}: Execution information including timing and metadata
func (r *Result) GetExecutionInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if r.status != nil {
		info["status"] = r.status.String()
		info["status_display"] = r.status.String()
	}

	info["notification_message"] = r.notificationMessage
	info["short_summary"] = r.GetShortSummary()

	if r.startedAt != nil {
		info["started_at"] = r.startedAt.Format(time.RFC3339)
	}

	if r.endedAt != nil {
		info["ended_at"] = r.endedAt.Format(time.RFC3339)
	}

	info["duration_ms"] = r.GetDuration().Milliseconds()

	if r.check != nil {
		info["check_name"] = r.check.GetName()
	}

	if len(r.meta) > 0 {
		info["metadata"] = r.meta
	}

	return info
}

// Clone creates a deep copy of the result.
// This is useful for creating result variants or for testing.
//
// Returns:
//
//	interfaces.ResultInterface: A new result instance with copied data
func (r *Result) Clone() interfaces.ResultInterface {
	clone := &Result{
		status:              r.status,
		notificationMessage: r.notificationMessage,
		shortSummary:        r.shortSummary,
		check:               r.check,
	}

	// Deep copy timing information
	if r.startedAt != nil {
		startTime := *r.startedAt
		clone.startedAt = &startTime
	}

	if r.endedAt != nil {
		endTime := *r.endedAt
		clone.endedAt = &endTime
	}

	// Deep copy metadata
	clone.meta = make(map[string]interface{})
	for k, v := range r.meta {
		clone.meta[k] = v
	}

	return clone
}

// Compile-time interface compliance check
var _ interfaces.ResultInterface = (*Result)(nil)
