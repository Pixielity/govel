package events

import (
	"fmt"
	"time"

	"govel/packages/bus/src/interfaces"
)

// BatchFinishedEvent is fired when a batch is finished (all jobs completed)
type BatchFinishedEvent struct {
	// Batch is the batch that was finished
	Batch interfaces.Batch

	// FinishedAt is when the batch finished
	FinishedAt time.Time

	// Duration is how long the batch took to complete
	Duration time.Duration

	// SuccessfulJobs is the number of jobs that succeeded
	SuccessfulJobs int

	// FailedJobs is the number of jobs that failed
	FailedJobs int

	// TotalJobs is the total number of jobs in the batch
	TotalJobs int

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewBatchFinishedEvent creates a new BatchFinishedEvent
func NewBatchFinishedEvent(batch interfaces.Batch) *BatchFinishedEvent {
	var duration time.Duration
	if batch.CreatedAt() != nil && batch.FinishedAt() != nil {
		duration = batch.FinishedAt().Sub(*batch.CreatedAt())
	}

	return &BatchFinishedEvent{
		Batch:          batch,
		FinishedAt:     time.Now(),
		Duration:       duration,
		SuccessfulJobs: batch.ProcessedJobs(),
		FailedJobs:     batch.FailedJobs(),
		TotalJobs:      batch.TotalJobs(),
		Options:        make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *BatchFinishedEvent) GetEventName() string {
	return "batch.finished"
}

// GetBatchID returns the batch ID
func (e *BatchFinishedEvent) GetBatchID() string {
	return e.Batch.ID()
}

// GetBatchName returns the batch name
func (e *BatchFinishedEvent) GetBatchName() string {
	return e.Batch.Name()
}

// IsSuccessful returns true if all jobs completed successfully
func (e *BatchFinishedEvent) IsSuccessful() bool {
	return e.FailedJobs == 0
}

// HasFailures returns true if any jobs failed
func (e *BatchFinishedEvent) HasFailures() bool {
	return e.FailedJobs > 0
}

// SuccessRate returns the success rate as a percentage (0.0 to 1.0)
func (e *BatchFinishedEvent) SuccessRate() float64 {
	if e.TotalJobs == 0 {
		return 0.0
	}
	return float64(e.SuccessfulJobs) / float64(e.TotalJobs)
}

// WithOption adds an option to the event
func (e *BatchFinishedEvent) WithOption(key string, value interface{}) *BatchFinishedEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *BatchFinishedEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// ToMap converts the event to a map for serialization
func (e *BatchFinishedEvent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"event":          e.GetEventName(),
		"batch_id":       e.GetBatchID(),
		"batch_name":     e.GetBatchName(),
		"finished_at":    e.FinishedAt,
		"duration_ms":    e.Duration.Milliseconds(),
		"successful_jobs": e.SuccessfulJobs,
		"failed_jobs":    e.FailedJobs,
		"total_jobs":     e.TotalJobs,
		"success_rate":   e.SuccessRate(),
		"options":        e.Options,
	}
}

// String returns a string representation of the event
func (e *BatchFinishedEvent) String() string {
	status := "successfully"
	if e.HasFailures() {
		status = fmt.Sprintf("with %d failures", e.FailedJobs)
	}

	return fmt.Sprintf("Batch '%s' (%s) finished %s. %d/%d jobs completed in %v",
		e.GetBatchName(),
		e.GetBatchID(),
		status,
		e.SuccessfulJobs,
		e.TotalJobs,
		e.Duration)
}
