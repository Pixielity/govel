package bus

import (
	"context"
	"fmt"
	"time"

	"govel/packages/bus/src/interfaces"
)

// Batch represents a collection of jobs that can be processed together
type Batch struct {
	data       interfaces.BatchData
	repository interfaces.BatchRepository
}

// NewBatch creates a new batch instance
func NewBatch(data interfaces.BatchData, repository interfaces.BatchRepository) *Batch {
	return &Batch{
		data:       data,
		repository: repository,
	}
}

// ID returns the batch ID
func (b *Batch) ID() string {
	return b.data.ID
}

// Name returns the batch name
func (b *Batch) Name() string {
	return b.data.Name
}

// TotalJobs returns the total number of jobs in the batch
func (b *Batch) TotalJobs() int {
	return b.data.TotalJobs
}

// PendingJobs returns the number of pending jobs
func (b *Batch) PendingJobs() int {
	return b.data.PendingJobs
}

// ProcessedJobs returns the number of processed jobs
func (b *Batch) ProcessedJobs() int {
	return b.data.ProcessedJobs()
}

// Progress returns the progress percentage (0-100)
func (b *Batch) Progress() float64 {
	return b.data.Progress()
}

// FailedJobs returns the number of failed jobs
func (b *Batch) FailedJobs() int {
	return b.data.FailedJobs
}

// HasFailures returns true if the batch has any failed jobs
func (b *Batch) HasFailures() bool {
	return b.data.HasFailures()
}

// AllowsFailures returns true if the batch allows failures
func (b *Batch) AllowsFailures() bool {
	return b.data.AllowsFailures()
}

// HasPendingJobs returns true if the batch has pending jobs
func (b *Batch) HasPendingJobs() bool {
	return b.data.HasPendingJobs()
}

// Finished returns true if the batch is finished
func (b *Batch) Finished() bool {
	return b.data.Finished()
}

// Cancelled returns true if the batch was cancelled
func (b *Batch) Cancelled() bool {
	return b.data.Cancelled()
}

// CreatedAt returns when the batch was created
func (b *Batch) CreatedAt() time.Time {
	return b.data.CreatedAt
}

// FinishedAt returns when the batch finished (nil if not finished)
func (b *Batch) FinishedAt() *time.Time {
	return b.data.FinishedAt
}

// CancelledAt returns when the batch was cancelled (nil if not cancelled)
func (b *Batch) CancelledAt() *time.Time {
	return b.data.CancelledAt
}

// Add adds more jobs to the batch
func (b *Batch) Add(ctx context.Context, jobs []interface{}) error {
	if b.Finished() || b.Cancelled() {
		return ErrBatchNotModifiable{BatchID: b.ID(), State: b.getState()}
	}

	// This would integrate with the queue system to add jobs
	// For now, we'll just increment the job count via repository
	return b.repository.IncrementTotalJobs(ctx, b.ID(), len(jobs))
}

// Cancel cancels the batch
func (b *Batch) Cancel(ctx context.Context) error {
	if b.Finished() || b.Cancelled() {
		return ErrBatchNotModifiable{BatchID: b.ID(), State: b.getState()}
	}

	return b.repository.Cancel(ctx, b.ID())
}

// Delete deletes the batch
func (b *Batch) Delete(ctx context.Context) error {
	return b.repository.Delete(ctx, b.ID())
}

// Fresh returns a fresh instance of the batch
func (b *Batch) Fresh(ctx context.Context) (interfaces.Batch, error) {
	data, err := b.repository.Get(ctx, b.ID())
	if err != nil {
		return nil, err
	}

	return NewBatch(data, b.repository), nil
}

// MarkJobAsProcessed marks a job as processed (decrements pending count)
func (b *Batch) MarkJobAsProcessed(ctx context.Context) error {
	return b.repository.DecrementPendingJobs(ctx, b.ID(), 1)
}

// MarkJobAsFailed marks a job as failed
func (b *Batch) MarkJobAsFailed(ctx context.Context, jobID string) error {
	// First decrement pending jobs
	if err := b.repository.DecrementPendingJobs(ctx, b.ID(), 1); err != nil {
		return err
	}

	// Then increment failed jobs
	return b.repository.IncrementFailedJobs(ctx, b.ID(), jobID)
}

// ShouldMarkAsFinished determines if the batch should be marked as finished
func (b *Batch) ShouldMarkAsFinished() bool {
	return !b.Finished() && !b.Cancelled() && b.data.PendingJobs <= 0
}

// MarkAsFinished marks the batch as finished if all jobs are complete
func (b *Batch) MarkAsFinished(ctx context.Context) error {
	if !b.ShouldMarkAsFinished() {
		return nil
	}

	return b.repository.MarkAsFinished(ctx, b.ID())
}

// GetOptions returns the batch options
func (b *Batch) GetOptions() map[string]interface{} {
	return b.data.Options
}

// GetOption returns a specific option value
func (b *Batch) GetOption(key string) interface{} {
	return b.data.Options[key]
}

// GetFailedJobIDs returns the IDs of failed jobs
func (b *Batch) GetFailedJobIDs() []string {
	return b.data.FailedJobIDs
}

// ToMap returns a map representation of the batch
func (b *Batch) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"id":            b.ID(),
		"name":          b.Name(),
		"total_jobs":    b.TotalJobs(),
		"pending_jobs":  b.PendingJobs(),
		"processed_jobs": b.ProcessedJobs(),
		"failed_jobs":   b.FailedJobs(),
		"progress":      b.Progress(),
		"finished":      b.Finished(),
		"cancelled":     b.Cancelled(),
		"created_at":    b.CreatedAt(),
		"options":       b.GetOptions(),
		"failed_job_ids": b.GetFailedJobIDs(),
	}

	if finishedAt := b.FinishedAt(); finishedAt != nil {
		result["finished_at"] = *finishedAt
	}

	if cancelledAt := b.CancelledAt(); cancelledAt != nil {
		result["cancelled_at"] = *cancelledAt
	}

	return result
}

// getState returns the current state of the batch
func (b *Batch) getState() string {
	if b.Cancelled() {
		return "cancelled"
	}
	if b.Finished() {
		return "finished"
	}
	if b.HasPendingJobs() {
		return "processing"
	}
	return "unknown"
}

// Error types

// ErrBatchNotModifiable is returned when trying to modify a finished or cancelled batch
type ErrBatchNotModifiable struct {
	BatchID string
	State   string
}

func (e ErrBatchNotModifiable) Error() string {
	return fmt.Sprintf("batch %s is %s and cannot be modified", e.BatchID, e.State)
}

// ErrBatchNotFound is returned when a batch cannot be found
type ErrBatchNotFound struct {
	BatchID string
}

func (e ErrBatchNotFound) Error() string {
	return fmt.Sprintf("batch %s not found", e.BatchID)
}

// Ensure Batch implements the Batch interface
var _ interfaces.Batch = (*Batch)(nil)
