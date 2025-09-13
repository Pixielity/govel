package events

import (
	"fmt"
	"time"

	"govel/bus/src/interfaces"
)

// BatchCancelledEvent is fired when a batch is cancelled
type BatchCancelledEvent struct {
	// Batch is the batch that was cancelled
	Batch interfaces.Batch

	// CancelledAt is when the batch was cancelled
	CancelledAt time.Time

	// Duration is how long the batch ran before being cancelled
	Duration time.Duration

	// CancelledBy is information about who/what cancelled the batch
	CancelledBy string

	// Reason is the reason for cancellation
	Reason string

	// ProcessedJobs is the number of jobs that were processed before cancellation
	ProcessedJobs int

	// FailedJobs is the number of jobs that failed before cancellation
	FailedJobs int

	// TotalJobs is the total number of jobs in the batch
	TotalJobs int

	// RemainingJobs is the number of jobs that were not processed
	RemainingJobs int

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewBatchCancelledEvent creates a new BatchCancelledEvent
func NewBatchCancelledEvent(batch interfaces.Batch, cancelledBy, reason string) *BatchCancelledEvent {
	var duration time.Duration
	if batch.CreatedAt() != nil {
		duration = time.Since(*batch.CreatedAt())
	}

	processedJobs := batch.ProcessedJobs()
	failedJobs := batch.FailedJobs()
	totalJobs := batch.TotalJobs()
	remainingJobs := totalJobs - processedJobs - failedJobs

	return &BatchCancelledEvent{
		Batch:         batch,
		CancelledAt:   time.Now(),
		Duration:      duration,
		CancelledBy:   cancelledBy,
		Reason:        reason,
		ProcessedJobs: processedJobs,
		FailedJobs:    failedJobs,
		TotalJobs:     totalJobs,
		RemainingJobs: remainingJobs,
		Options:       make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *BatchCancelledEvent) GetEventName() string {
	return "batch.cancelled"
}

// GetBatchID returns the batch ID
func (e *BatchCancelledEvent) GetBatchID() string {
	return e.Batch.ID()
}

// GetBatchName returns the batch name
func (e *BatchCancelledEvent) GetBatchName() string {
	return e.Batch.Name()
}

// HasCancelledBy returns true if cancelledBy is specified
func (e *BatchCancelledEvent) HasCancelledBy() bool {
	return e.CancelledBy != ""
}

// HasReason returns true if a reason is specified
func (e *BatchCancelledEvent) HasReason() bool {
	return e.Reason != ""
}

// GetCompletionRate returns the completion rate at time of cancellation (0.0 to 1.0)
func (e *BatchCancelledEvent) GetCompletionRate() float64 {
	if e.TotalJobs == 0 {
		return 0.0
	}
	return float64(e.ProcessedJobs+e.FailedJobs) / float64(e.TotalJobs)
}

// GetSuccessRate returns the success rate of completed jobs (0.0 to 1.0)
func (e *BatchCancelledEvent) GetSuccessRate() float64 {
	completedJobs := e.ProcessedJobs + e.FailedJobs
	if completedJobs == 0 {
		return 0.0
	}
	return float64(e.ProcessedJobs) / float64(completedJobs)
}

// WithOption adds an option to the event
func (e *BatchCancelledEvent) WithOption(key string, value interface{}) *BatchCancelledEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *BatchCancelledEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// ToMap converts the event to a map for serialization
func (e *BatchCancelledEvent) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"event":           e.GetEventName(),
		"batch_id":        e.GetBatchID(),
		"batch_name":      e.GetBatchName(),
		"cancelled_at":    e.CancelledAt,
		"duration_ms":     e.Duration.Milliseconds(),
		"processed_jobs":  e.ProcessedJobs,
		"failed_jobs":     e.FailedJobs,
		"total_jobs":      e.TotalJobs,
		"remaining_jobs":  e.RemainingJobs,
		"completion_rate": e.GetCompletionRate(),
		"success_rate":    e.GetSuccessRate(),
		"options":         e.Options,
	}

	if e.HasCancelledBy() {
		data["cancelled_by"] = e.CancelledBy
	}

	if e.HasReason() {
		data["reason"] = e.Reason
	}

	return data
}

// String returns a string representation of the event
func (e *BatchCancelledEvent) String() string {
	baseMsg := fmt.Sprintf("Batch '%s' (%s) was cancelled after %v. %d/%d jobs processed",
		e.GetBatchName(),
		e.GetBatchID(),
		e.Duration,
		e.ProcessedJobs,
		e.TotalJobs)

	if e.HasCancelledBy() {
		baseMsg += fmt.Sprintf(" by %s", e.CancelledBy)
	}

	if e.HasReason() {
		baseMsg += fmt.Sprintf(". Reason: %s", e.Reason)
	}

	return baseMsg
}
