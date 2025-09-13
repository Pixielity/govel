package events

import (
	"fmt"
	"time"

	"govel/packages/bus/src/interfaces"
)

// BatchDispatchedEvent is fired when a batch is dispatched
type BatchDispatchedEvent struct {
	// Batch is the batch that was dispatched
	Batch interfaces.Batch

	// DispatchedAt is when the batch was dispatched
	DispatchedAt time.Time

	// JobCount is the number of jobs in the batch
	JobCount int

	// Options contains any additional options or metadata
	Options map[string]interface{}
}

// NewBatchDispatchedEvent creates a new BatchDispatchedEvent
func NewBatchDispatchedEvent(batch interfaces.Batch) *BatchDispatchedEvent {
	return &BatchDispatchedEvent{
		Batch:        batch,
		DispatchedAt: time.Now(),
		JobCount:     batch.TotalJobs(),
		Options:      make(map[string]interface{}),
	}
}

// GetEventName returns the event name
func (e *BatchDispatchedEvent) GetEventName() string {
	return "batch.dispatched"
}

// GetBatchID returns the batch ID
func (e *BatchDispatchedEvent) GetBatchID() string {
	return e.Batch.ID()
}

// GetBatchName returns the batch name
func (e *BatchDispatchedEvent) GetBatchName() string {
	return e.Batch.Name()
}

// WithOption adds an option to the event
func (e *BatchDispatchedEvent) WithOption(key string, value interface{}) *BatchDispatchedEvent {
	e.Options[key] = value
	return e
}

// GetOption retrieves an option from the event
func (e *BatchDispatchedEvent) GetOption(key string) interface{} {
	return e.Options[key]
}

// ToMap converts the event to a map for serialization
func (e *BatchDispatchedEvent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"event":        e.GetEventName(),
		"batch_id":     e.GetBatchID(),
		"batch_name":   e.GetBatchName(),
		"job_count":    e.JobCount,
		"dispatched_at": e.DispatchedAt,
		"options":      e.Options,
	}
}

// String returns a string representation of the event
func (e *BatchDispatchedEvent) String() string {
	return fmt.Sprintf("Batch '%s' (%s) with %d jobs dispatched at %s",
		e.GetBatchName(),
		e.GetBatchID(),
		e.JobCount,
		e.DispatchedAt.Format(time.RFC3339))
}
