package bus

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"govel/new/bus/interfaces"
)

// PendingBatch represents a batch of jobs waiting to be dispatched
type PendingBatch struct {
	id                string
	name              string
	jobs              []interface{}
	options           map[string]interface{}
	thenCallbacks     []func(ctx context.Context, batch interfaces.Batch) error
	catchCallbacks    []func(ctx context.Context, batch interfaces.Batch, err error) error
	finallyCallbacks  []func(ctx context.Context, batch interfaces.Batch) error
	batchRepository   interfaces.BatchRepository
	queueDispatcher   interface{} // Would be a queue dispatcher in full implementation
}

// NewPendingBatch creates a new pending batch
func NewPendingBatch(jobs []interface{}, repository interfaces.BatchRepository) *PendingBatch {
	return &PendingBatch{
		id:              uuid.New().String(),
		jobs:            jobs,
		options:         make(map[string]interface{}),
		thenCallbacks:   make([]func(ctx context.Context, batch interfaces.Batch) error, 0),
		catchCallbacks:  make([]func(ctx context.Context, batch interfaces.Batch, err error) error, 0),
		finallyCallbacks: make([]func(ctx context.Context, batch interfaces.Batch) error, 0),
		batchRepository: repository,
	}
}

// Name sets the name of the batch
func (pb *PendingBatch) Name(name string) interfaces.PendingBatch {
	pb.name = name
	return pb
}

// OnConnection sets the connection for the batch
func (pb *PendingBatch) OnConnection(connection string) interfaces.PendingBatch {
	pb.options["connection"] = connection
	return pb
}

// OnQueue sets the queue for the batch
func (pb *PendingBatch) OnQueue(queue string) interfaces.PendingBatch {
	pb.options["queue"] = queue
	return pb
}

// Then sets a callback to run after the batch completes successfully
func (pb *PendingBatch) Then(callback func(ctx context.Context, batch interfaces.Batch) error) interfaces.PendingBatch {
	pb.thenCallbacks = append(pb.thenCallbacks, callback)
	return pb
}

// Catch sets a callback to run if the batch fails
func (pb *PendingBatch) Catch(callback func(ctx context.Context, batch interfaces.Batch, err error) error) interfaces.PendingBatch {
	pb.catchCallbacks = append(pb.catchCallbacks, callback)
	return pb
}

// Finally sets a callback to run regardless of success or failure
func (pb *PendingBatch) Finally(callback func(ctx context.Context, batch interfaces.Batch) error) interfaces.PendingBatch {
	pb.finallyCallbacks = append(pb.finallyCallbacks, callback)
	return pb
}

// AllowFailures allows the batch to continue even if some jobs fail
func (pb *PendingBatch) AllowFailures(allowed bool) interfaces.PendingBatch {
	pb.options["allow_failures"] = allowed
	return pb
}

// Delay sets a delay before the batch starts processing
func (pb *PendingBatch) Delay(delay time.Duration) interfaces.PendingBatch {
	pb.options["delay"] = delay
	return pb
}

// SetOption sets a custom option
func (pb *PendingBatch) SetOption(key string, value interface{}) interfaces.PendingBatch {
	pb.options[key] = value
	return pb
}

// Dispatch dispatches the batch
func (pb *PendingBatch) Dispatch(ctx context.Context) (interfaces.Batch, error) {
	if len(pb.jobs) == 0 {
		return nil, fmt.Errorf("cannot dispatch empty batch")
	}

	// Create pending batch data
	pendingData := interfaces.PendingBatchData{
		ID:        pb.id,
		Name:      pb.name,
		Jobs:      pb.jobs,
		Options:   pb.options,
		CreatedAt: time.Now(),
	}

	// Store the batch
	batchData, err := pb.batchRepository.Store(ctx, pendingData)
	if err != nil {
		return nil, fmt.Errorf("failed to store batch: %w", err)
	}

	// Create batch instance
	batch := NewBatch(batchData, pb.batchRepository)

	// Store callbacks in batch options for later execution
	if len(pb.thenCallbacks) > 0 {
		batch.data.Options["then_callbacks"] = pb.thenCallbacks
	}
	if len(pb.catchCallbacks) > 0 {
		batch.data.Options["catch_callbacks"] = pb.catchCallbacks
	}
	if len(pb.finallyCallbacks) > 0 {
		batch.data.Options["finally_callbacks"] = pb.finallyCallbacks
	}

	// In a full implementation, this would dispatch jobs to the queue
	// For now, we'll just return the batch
	return batch, nil
}

// Add adds more jobs to the pending batch
func (pb *PendingBatch) Add(jobs ...interface{}) interfaces.PendingBatch {
	pb.jobs = append(pb.jobs, jobs...)
	return pb
}

// Jobs returns the jobs in the pending batch
func (pb *PendingBatch) Jobs() []interface{} {
	return pb.jobs
}

// JobCount returns the number of jobs in the pending batch
func (pb *PendingBatch) JobCount() int {
	return len(pb.jobs)
}

// GetID returns the batch ID
func (pb *PendingBatch) GetID() string {
	return pb.id
}

// GetName returns the batch name
func (pb *PendingBatch) GetName() string {
	return pb.name
}

// GetOptions returns the batch options
func (pb *PendingBatch) GetOptions() map[string]interface{} {
	return pb.options
}

// GetOption returns a specific option value
func (pb *PendingBatch) GetOption(key string) interface{} {
	return pb.options[key]
}

// HasOption checks if an option exists
func (pb *PendingBatch) HasOption(key string) bool {
	_, exists := pb.options[key]
	return exists
}

// ToMap returns a map representation of the pending batch
func (pb *PendingBatch) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         pb.id,
		"name":       pb.name,
		"job_count":  pb.JobCount(),
		"options":    pb.options,
		"then_callbacks":    len(pb.thenCallbacks),
		"catch_callbacks":   len(pb.catchCallbacks),
		"finally_callbacks": len(pb.finallyCallbacks),
	}
}

// Clone creates a copy of the pending batch
func (pb *PendingBatch) Clone() *PendingBatch {
	// Deep copy jobs
	jobsCopy := make([]interface{}, len(pb.jobs))
	copy(jobsCopy, pb.jobs)

	// Deep copy options
	optionsCopy := make(map[string]interface{})
	for k, v := range pb.options {
		optionsCopy[k] = v
	}

	// Create new pending batch
	newPB := &PendingBatch{
		id:              uuid.New().String(),
		name:            pb.name,
		jobs:            jobsCopy,
		options:         optionsCopy,
		thenCallbacks:   make([]func(ctx context.Context, batch interfaces.Batch) error, len(pb.thenCallbacks)),
		catchCallbacks:  make([]func(ctx context.Context, batch interfaces.Batch, err error) error, len(pb.catchCallbacks)),
		finallyCallbacks: make([]func(ctx context.Context, batch interfaces.Batch) error, len(pb.finallyCallbacks)),
		batchRepository: pb.batchRepository,
		queueDispatcher: pb.queueDispatcher,
	}

	// Copy callbacks
	copy(newPB.thenCallbacks, pb.thenCallbacks)
	copy(newPB.catchCallbacks, pb.catchCallbacks)
	copy(newPB.finallyCallbacks, pb.finallyCallbacks)

	return newPB
}

// Validate validates the pending batch before dispatch
func (pb *PendingBatch) Validate() error {
	if len(pb.jobs) == 0 {
		return fmt.Errorf("batch cannot be empty")
	}

	if pb.batchRepository == nil {
		return fmt.Errorf("batch repository is required")
	}

	return nil
}

// Ensure PendingBatch implements the PendingBatch interface
var _ interfaces.PendingBatch = (*PendingBatch)(nil)
