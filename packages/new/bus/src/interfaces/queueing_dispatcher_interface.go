package interfaces

import (
	"context"
	"time"
)

// QueueingDispatcher extends Dispatcher with queueing capabilities
type QueueingDispatcher interface {
	Dispatcher

	// DispatchToQueue dispatches a command to a queue
	DispatchToQueue(ctx context.Context, command interface{}) error

	// DispatchAfterResponse dispatches a command after the HTTP response is sent
	DispatchAfterResponse(ctx context.Context, command interface{}) error

	// Batch creates a new batch of queueable jobs
	Batch(jobs []interface{}) PendingBatch

	// Chain creates a new chain of queueable jobs
	Chain(jobs []interface{}) PendingChain

	// FindBatch finds a batch by its ID
	FindBatch(ctx context.Context, batchID string) (Batch, error)
}

// PendingBatch represents a batch of jobs waiting to be dispatched
type PendingBatch interface {
	// Name sets the name of the batch
	Name(name string) PendingBatch

	// OnConnection sets the connection for the batch
	OnConnection(connection string) PendingBatch

	// OnQueue sets the queue for the batch
	OnQueue(queue string) PendingBatch

	// Then sets a callback to run after the batch completes
	Then(callback func(ctx context.Context, batch Batch) error) PendingBatch

	// Catch sets a callback to run if the batch fails
	Catch(callback func(ctx context.Context, batch Batch, err error) error) PendingBatch

	// Finally sets a callback to run regardless of success or failure
	Finally(callback func(ctx context.Context, batch Batch) error) PendingBatch

	// AllowFailures allows the batch to continue even if some jobs fail
	AllowFailures(allowed bool) PendingBatch

	// Dispatch dispatches the batch
	Dispatch(ctx context.Context) (Batch, error)
}

// PendingChain represents a chain of jobs waiting to be dispatched
type PendingChain interface {
	// OnConnection sets the connection for the chain
	OnConnection(connection string) PendingChain

	// OnQueue sets the queue for the chain
	OnQueue(queue string) PendingChain

	// Delay sets a delay before the chain starts
	Delay(delay time.Duration) PendingChain

	// Catch sets a callback to run if the chain fails
	Catch(callback func(ctx context.Context, err error) error) PendingChain

	// Dispatch dispatches the chain
	Dispatch(ctx context.Context) error
}

// Batch represents a collection of jobs that can be processed together
type Batch interface {
	// ID returns the batch ID
	ID() string

	// Name returns the batch name
	Name() string

	// TotalJobs returns the total number of jobs in the batch
	TotalJobs() int

	// PendingJobs returns the number of pending jobs
	PendingJobs() int

	// ProcessedJobs returns the number of processed jobs
	ProcessedJobs() int

	// Progress returns the progress percentage (0-100)
	Progress() float64

	// FailedJobs returns the number of failed jobs
	FailedJobs() int

	// HasFailures returns true if the batch has any failed jobs
	HasFailures() bool

	// AllowsFailures returns true if the batch allows failures
	AllowsFailures() bool

	// HasPendingJobs returns true if the batch has pending jobs
	HasPendingJobs() bool

	// Finished returns true if the batch is finished
	Finished() bool

	// Cancelled returns true if the batch was cancelled
	Cancelled() bool

	// CreatedAt returns when the batch was created
	CreatedAt() time.Time

	// FinishedAt returns when the batch finished (nil if not finished)
	FinishedAt() *time.Time

	// CancelledAt returns when the batch was cancelled (nil if not cancelled)
	CancelledAt() *time.Time

	// Add adds more jobs to the batch
	Add(ctx context.Context, jobs []interface{}) error

	// Cancel cancels the batch
	Cancel(ctx context.Context) error

	// Delete deletes the batch
	Delete(ctx context.Context) error

	// Fresh returns a fresh instance of the batch
	Fresh(ctx context.Context) (Batch, error)
}
