package interfaces

import (
	"context"
	"time"
)

// BatchRepository defines the contract for batch storage and retrieval
type BatchRepository interface {
	// Get retrieves a batch by its ID
	Get(ctx context.Context, batchID string) (BatchData, error)

	// Store stores a new batch
	Store(ctx context.Context, batch PendingBatchData) (BatchData, error)

	// IncrementTotalJobs increments the total job count for a batch
	IncrementTotalJobs(ctx context.Context, batchID string, amount int) error

	// DecrementPendingJobs decrements the pending job count for a batch
	DecrementPendingJobs(ctx context.Context, batchID string, amount int) error

	// IncrementFailedJobs increments the failed job count and adds job ID
	IncrementFailedJobs(ctx context.Context, batchID string, jobID string) error

	// MarkAsFinished marks a batch as finished
	MarkAsFinished(ctx context.Context, batchID string) error

	// Cancel cancels a batch
	Cancel(ctx context.Context, batchID string) error

	// Delete deletes a batch
	Delete(ctx context.Context, batchID string) error

	// Transaction executes a function within a transaction
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BatchData represents the stored batch data
type BatchData struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	TotalJobs    int               `json:"total_jobs"`
	PendingJobs  int               `json:"pending_jobs"`
	FailedJobs   int               `json:"failed_jobs"`
	FailedJobIDs []string          `json:"failed_job_ids"`
	Options      map[string]interface{} `json:"options"`
	CreatedAt    time.Time         `json:"created_at"`
	CancelledAt  *time.Time        `json:"cancelled_at,omitempty"`
	FinishedAt   *time.Time        `json:"finished_at,omitempty"`
}

// PendingBatchData represents a batch that's about to be stored
type PendingBatchData struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Jobs      []interface{}          `json:"jobs"`
	Options   map[string]interface{} `json:"options"`
	CreatedAt time.Time              `json:"created_at"`
}

// ProcessedJobs returns the number of processed jobs
func (b BatchData) ProcessedJobs() int {
	return b.TotalJobs - b.PendingJobs
}

// Progress returns the batch progress percentage
func (b BatchData) Progress() float64 {
	if b.TotalJobs == 0 {
		return 0.0
	}
	return float64(b.ProcessedJobs()) / float64(b.TotalJobs) * 100.0
}

// HasFailures returns true if the batch has failed jobs
func (b BatchData) HasFailures() bool {
	return b.FailedJobs > 0
}

// HasPendingJobs returns true if the batch has pending jobs
func (b BatchData) HasPendingJobs() bool {
	return b.PendingJobs > 0
}

// Finished returns true if the batch is finished
func (b BatchData) Finished() bool {
	return b.FinishedAt != nil
}

// Cancelled returns true if the batch was cancelled
func (b BatchData) Cancelled() bool {
	return b.CancelledAt != nil
}

// AllowsFailures returns true if the batch allows failures
func (b BatchData) AllowsFailures() bool {
	if allowFailures, exists := b.Options["allow_failures"]; exists {
		if allowed, ok := allowFailures.(bool); ok {
			return allowed
		}
	}
	return false
}
