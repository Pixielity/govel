package events

import (
	"fmt"
	"time"
)

// JobProcessedEvent is fired when a single job is processed successfully
type JobProcessedEvent struct {
	// JobID is the unique identifier for the job
	JobID string

	// JobType is the type/class name of the job
	JobType string

	// BatchID is the ID of the batch this job belongs to (if any)
	BatchID string

	// ProcessedAt is when the job was processed
	ProcessedAt time.Time

	// Duration is how long the job took to process
	Duration time.Duration

	// Result is the result data from processing (if any)
	Result interface{}

	// Queue is the queue the job was processed from
	Queue string

	// Attempts is the number of attempts made to process this job
	Attempts int

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewJobProcessedEvent creates a new JobProcessedEvent
func NewJobProcessedEvent(jobID, jobType string) *JobProcessedEvent {
	return &JobProcessedEvent{
		JobID:       jobID,
		JobType:     jobType,
		ProcessedAt: time.Now(),
		Attempts:    1,
		Options:     make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *JobProcessedEvent) GetEventName() string {
	return "job.processed"
}

// WithBatch sets the batch information
func (e *JobProcessedEvent) WithBatch(batchID string) *JobProcessedEvent {
	e.BatchID = batchID
	return e
}

// WithDuration sets the processing duration
func (e *JobProcessedEvent) WithDuration(duration time.Duration) *JobProcessedEvent {
	e.Duration = duration
	return e
}

// WithResult sets the processing result
func (e *JobProcessedEvent) WithResult(result interface{}) *JobProcessedEvent {
	e.Result = result
	return e
}

// WithQueue sets the queue information
func (e *JobProcessedEvent) WithQueue(queue string) *JobProcessedEvent {
	e.Queue = queue
	return e
}

// WithAttempts sets the number of attempts
func (e *JobProcessedEvent) WithAttempts(attempts int) *JobProcessedEvent {
	e.Attempts = attempts
	return e
}

// WithOption adds an option to the event
func (e *JobProcessedEvent) WithOption(key string, value interface{}) *JobProcessedEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *JobProcessedEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// HasBatch returns true if the job belongs to a batch
func (e *JobProcessedEvent) HasBatch() bool {
	return e.BatchID != ""
}

// HasResult returns true if there is a processing result
func (e *JobProcessedEvent) HasResult() bool {
	return e.Result != nil
}

// HasQueue returns true if queue information is available
func (e *JobProcessedEvent) HasQueue() bool {
	return e.Queue != ""
}

// ToMap converts the event to a map for serialization
func (e *JobProcessedEvent) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"event":        e.GetEventName(),
		"job_id":       e.JobID,
		"job_type":     e.JobType,
		"processed_at": e.ProcessedAt,
		"duration_ms":  e.Duration.Milliseconds(),
		"attempts":     e.Attempts,
		"options":      e.Options,
	}

	if e.HasBatch() {
		data["batch_id"] = e.BatchID
	}

	if e.HasResult() {
		data["result"] = e.Result
	}

	if e.HasQueue() {
		data["queue"] = e.Queue
	}

	return data
}

// String returns a string representation of the event
func (e *JobProcessedEvent) String() string {
	baseMsg := fmt.Sprintf("Job %s (%s) processed successfully", e.JobID, e.JobType)

	if e.Duration > 0 {
		baseMsg += fmt.Sprintf(" in %v", e.Duration)
	}

	if e.Attempts > 1 {
		baseMsg += fmt.Sprintf(" after %d attempts", e.Attempts)
	}

	if e.HasBatch() {
		baseMsg += fmt.Sprintf(" (batch: %s)", e.BatchID)
	}

	if e.HasQueue() {
		baseMsg += fmt.Sprintf(" from queue '%s'", e.Queue)
	}

	return baseMsg
}
