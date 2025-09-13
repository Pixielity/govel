package events

import (
	"fmt"
	"time"

	"govel/packages/bus/src/interfaces"
)

// BatchFailedEvent is fired when a batch fails (reaches failure threshold)
type BatchFailedEvent struct {
	// Batch is the batch that failed
	Batch interfaces.Batch

	// FailedAt is when the batch was marked as failed
	FailedAt time.Time

	// Duration is how long the batch ran before failing
	Duration time.Duration

	// LastJobError is the error from the job that caused the batch to fail
	LastJobError error

	// FailedJobs is the number of jobs that failed
	FailedJobs int

	// ProcessedJobs is the number of jobs that were processed (successful + failed)
	ProcessedJobs int

	// TotalJobs is the total number of jobs in the batch
	TotalJobs int

	// RemainingJobs is the number of jobs that were not processed
	RemainingJobs int

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewBatchFailedEvent creates a new BatchFailedEvent
func NewBatchFailedEvent(batch interfaces.Batch, lastJobError error) *BatchFailedEvent {
	var duration time.Duration
	if batch.CreatedAt() != nil {
		duration = time.Since(*batch.CreatedAt())
	}

	processedJobs := batch.ProcessedJobs() + batch.FailedJobs()
	remainingJobs := batch.TotalJobs() - processedJobs

	return &BatchFailedEvent{
		Batch:         batch,
		FailedAt:      time.Now(),
		Duration:      duration,
		LastJobError:  lastJobError,
		FailedJobs:    batch.FailedJobs(),
		ProcessedJobs: processedJobs,
		TotalJobs:     batch.TotalJobs(),
		RemainingJobs: remainingJobs,
		Options:       make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *BatchFailedEvent) GetEventName() string {
	return "batch.failed"
}

// GetBatchID returns the batch ID
func (e *BatchFailedEvent) GetBatchID() string {
	return e.Batch.ID()
}

// GetBatchName returns the batch name
func (e *BatchFailedEvent) GetBatchName() string {
	return e.Batch.Name()
}

// HasLastJobError returns true if there is a last job error
func (e *BatchFailedEvent) HasLastJobError() bool {
	return e.LastJobError != nil
}

// GetLastJobErrorMessage returns the last job error message
func (e *BatchFailedEvent) GetLastJobErrorMessage() string {
	if e.LastJobError != nil {
		return e.LastJobError.Error()
	}
	return ""
}

// GetFailureRate returns the failure rate as a percentage (0.0 to 1.0)
func (e *BatchFailedEvent) GetFailureRate() float64 {
	if e.ProcessedJobs == 0 {
		return 0.0
	}
	return float64(e.FailedJobs) / float64(e.ProcessedJobs)
}

// GetCompletionRate returns the completion rate as a percentage (0.0 to 1.0)
func (e *BatchFailedEvent) GetCompletionRate() float64 {
	if e.TotalJobs == 0 {
		return 0.0
	}
	return float64(e.ProcessedJobs) / float64(e.TotalJobs)
}

// WithOption adds an option to the event
func (e *BatchFailedEvent) WithOption(key string, value interface{}) *BatchFailedEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *BatchFailedEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// ToMap converts the event to a map for serialization
func (e *BatchFailedEvent) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"event":           e.GetEventName(),
		"batch_id":        e.GetBatchID(),
		"batch_name":      e.GetBatchName(),
		"failed_at":       e.FailedAt,
		"duration_ms":     e.Duration.Milliseconds(),
		"failed_jobs":     e.FailedJobs,
		"processed_jobs":  e.ProcessedJobs,
		"total_jobs":      e.TotalJobs,
		"remaining_jobs":  e.RemainingJobs,
		"failure_rate":    e.GetFailureRate(),
		"completion_rate": e.GetCompletionRate(),
		"options":         e.Options,
	}

	if e.HasLastJobError() {
		data["last_job_error"] = e.GetLastJobErrorMessage()
	}

	return data
}

// String returns a string representation of the event
func (e *BatchFailedEvent) String() string {
	baseMsg := fmt.Sprintf("Batch '%s' (%s) failed after %v. %d jobs failed, %d/%d jobs processed",
		e.GetBatchName(),
		e.GetBatchID(),
		e.Duration,
		e.FailedJobs,
		e.ProcessedJobs,
		e.TotalJobs)

	if e.HasLastJobError() {
		baseMsg += fmt.Sprintf(". Last error: %s", e.GetLastJobErrorMessage())
	}

	return baseMsg
}
