package events

import (
	"fmt"
	"time"
)

// JobFailedEvent is fired when a single job fails to process
type JobFailedEvent struct {
	// JobID is the unique identifier for the job
	JobID string

	// JobType is the type/class name of the job
	JobType string

	// BatchID is the ID of the batch this job belongs to (if any)
	BatchID string

	// FailedAt is when the job failed
	FailedAt time.Time

	// Duration is how long the job ran before failing
	Duration time.Duration

	// Error is the error that caused the job to fail
	Error error

	// Queue is the queue the job was processed from
	Queue string

	// Attempts is the number of attempts made to process this job
	Attempts int

	// MaxAttempts is the maximum number of attempts allowed
	MaxAttempts int

	// WillRetry indicates if the job will be retried
	WillRetry bool

	// RetryDelay is the delay before the next retry (if applicable)
	RetryDelay time.Duration

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewJobFailedEvent creates a new JobFailedEvent
func NewJobFailedEvent(jobID, jobType string, err error) *JobFailedEvent {
	return &JobFailedEvent{
		JobID:       jobID,
		JobType:     jobType,
		Error:       err,
		FailedAt:    time.Now(),
		Attempts:    1,
		MaxAttempts: 1,
		WillRetry:   false,
		Options:     make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *JobFailedEvent) GetEventName() string {
	return "job.failed"
}

// WithBatch sets the batch information
func (e *JobFailedEvent) WithBatch(batchID string) *JobFailedEvent {
	e.BatchID = batchID
	return e
}

// WithDuration sets the processing duration
func (e *JobFailedEvent) WithDuration(duration time.Duration) *JobFailedEvent {
	e.Duration = duration
	return e
}

// WithQueue sets the queue information
func (e *JobFailedEvent) WithQueue(queue string) *JobFailedEvent {
	e.Queue = queue
	return e
}

// WithAttempts sets the attempt information
func (e *JobFailedEvent) WithAttempts(attempts, maxAttempts int) *JobFailedEvent {
	e.Attempts = attempts
	e.MaxAttempts = maxAttempts
	e.WillRetry = attempts < maxAttempts
	return e
}

// WithRetry sets the retry information
func (e *JobFailedEvent) WithRetry(willRetry bool, retryDelay time.Duration) *JobFailedEvent {
	e.WillRetry = willRetry
	e.RetryDelay = retryDelay
	return e
}

// WithOption adds an option to the event
func (e *JobFailedEvent) WithOption(key string, value interface{}) *JobFailedEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *JobFailedEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// HasBatch returns true if the job belongs to a batch
func (e *JobFailedEvent) HasBatch() bool {
	return e.BatchID != ""
}

// HasQueue returns true if queue information is available
func (e *JobFailedEvent) HasQueue() bool {
	return e.Queue != ""
}

// IsFinalFailure returns true if this is the final failure (no more retries)
func (e *JobFailedEvent) IsFinalFailure() bool {
	return !e.WillRetry
}

// GetErrorMessage returns the error message
func (e *JobFailedEvent) GetErrorMessage() string {
	if e.Error != nil {
		return e.Error.Error()
	}
	return ""
}

// GetRemainingAttempts returns the number of remaining attempts
func (e *JobFailedEvent) GetRemainingAttempts() int {
	return e.MaxAttempts - e.Attempts
}

// ToMap converts the event to a map for serialization
func (e *JobFailedEvent) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"event":             e.GetEventName(),
		"job_id":            e.JobID,
		"job_type":          e.JobType,
		"failed_at":         e.FailedAt,
		"duration_ms":       e.Duration.Milliseconds(),
		"attempts":          e.Attempts,
		"max_attempts":      e.MaxAttempts,
		"will_retry":        e.WillRetry,
		"remaining_attempts": e.GetRemainingAttempts(),
		"options":           e.Options,
	}

	if e.Error != nil {
		data["error"] = e.GetErrorMessage()
	}

	if e.HasBatch() {
		data["batch_id"] = e.BatchID
	}

	if e.HasQueue() {
		data["queue"] = e.Queue
	}

	if e.WillRetry && e.RetryDelay > 0 {
		data["retry_delay_ms"] = e.RetryDelay.Milliseconds()
	}

	return data
}

// String returns a string representation of the event
func (e *JobFailedEvent) String() string {
	baseMsg := fmt.Sprintf("Job %s (%s) failed", e.JobID, e.JobType)

	if e.Duration > 0 {
		baseMsg += fmt.Sprintf(" after %v", e.Duration)
	}

	if e.Error != nil {
		baseMsg += fmt.Sprintf(": %s", e.GetErrorMessage())
	}

	if e.Attempts > 1 {
		baseMsg += fmt.Sprintf(" (attempt %d/%d)", e.Attempts, e.MaxAttempts)
	} else if e.MaxAttempts > 1 {
		baseMsg += fmt.Sprintf(" (attempt %d/%d)", e.Attempts, e.MaxAttempts)
	}

	if e.WillRetry {
		if e.RetryDelay > 0 {
			baseMsg += fmt.Sprintf(", will retry in %v", e.RetryDelay)
		} else {
			baseMsg += ", will retry"
		}
	} else {
		baseMsg += ", no more retries"
	}

	if e.HasBatch() {
		baseMsg += fmt.Sprintf(" (batch: %s)", e.BatchID)
	}

	if e.HasQueue() {
		baseMsg += fmt.Sprintf(" from queue '%s'", e.Queue)
	}

	return baseMsg
}
