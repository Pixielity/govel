package traits

import (
	"context"
)

// Batchable provides batch processing capabilities to jobs
// This is implemented as a struct that can be embedded in other structs
type Batchable struct {
	// BatchID is the ID of the batch this job belongs to
	BatchID string `json:"batch_id,omitempty"`

	// JobID is the unique identifier for this job within the batch
	JobID string `json:"job_id,omitempty"`

	// BatchRepository is injected at runtime to interact with batches
	// This would be set by the job processor
	batchRepository interface{} `json:"-"`

	// FailureCallbacks are called when the job fails
	FailureCallbacks []func(ctx context.Context, err error) error `json:"-"`

	// SuccessCallbacks are called when the job succeeds
	SuccessCallbacks []func(ctx context.Context, result interface{}) error `json:"-"`
}

// NewBatchable creates a new Batchable instance
func NewBatchable() *Batchable {
	return &Batchable{
		FailureCallbacks: make([]func(ctx context.Context, err error) error, 0),
		SuccessCallbacks: make([]func(ctx context.Context, result interface{}) error, 0),
	}
}

// WithBatchID sets the batch ID
func (b *Batchable) WithBatchID(batchID string) *Batchable {
	b.BatchID = batchID
	return b
}

// WithJobID sets the job ID
func (b *Batchable) WithJobID(jobID string) *Batchable {
	b.JobID = jobID
	return b
}

// GetBatchID returns the batch ID
func (b *Batchable) GetBatchID() string {
	return b.BatchID
}

// GetJobID returns the job ID
func (b *Batchable) GetJobID() string {
	return b.JobID
}

// HasBatch returns true if this job belongs to a batch
func (b *Batchable) HasBatch() bool {
	return b.BatchID != ""
}

// OnBatchSuccess adds a callback to run when the entire batch succeeds
func (b *Batchable) OnBatchSuccess(callback func(ctx context.Context, result interface{}) error) *Batchable {
	b.SuccessCallbacks = append(b.SuccessCallbacks, callback)
	return b
}

// OnBatchFailure adds a callback to run when the batch fails
func (b *Batchable) OnBatchFailure(callback func(ctx context.Context, err error) error) *Batchable {
	b.FailureCallbacks = append(b.FailureCallbacks, callback)
	return b
}

// NotifyBatchSuccess notifies all success callbacks
func (b *Batchable) NotifyBatchSuccess(ctx context.Context, result interface{}) error {
	for _, callback := range b.SuccessCallbacks {
		if err := callback(ctx, result); err != nil {
			return err
		}
	}
	return nil
}

// NotifyBatchFailure notifies all failure callbacks
func (b *Batchable) NotifyBatchFailure(ctx context.Context, err error) error {
	for _, callback := range b.FailureCallbacks {
		if callbackErr := callback(ctx, err); callbackErr != nil {
			return callbackErr
		}
	}
	return nil
}

// Batch returns the current batch information (if available)
func (b *Batchable) Batch(ctx context.Context) (BatchInfo, error) {
	if !b.HasBatch() {
		return BatchInfo{}, nil
	}

	// In a real implementation, this would use the batch repository
	// to fetch current batch information
	return BatchInfo{
		ID:     b.BatchID,
		JobID:  b.JobID,
		Exists: true,
	}, nil
}

// MarkJobAsFailed marks this specific job as failed within the batch
func (b *Batchable) MarkJobAsFailed(ctx context.Context, err error) error {
	if !b.HasBatch() {
		return nil
	}

	// In a real implementation, this would interact with the batch repository
	// to mark the job as failed and update batch counters
	return nil
}

// MarkJobAsFinished marks this specific job as finished within the batch
func (b *Batchable) MarkJobAsFinished(ctx context.Context, result interface{}) error {
	if !b.HasBatch() {
		return nil
	}

	// In a real implementation, this would interact with the batch repository
	// to mark the job as finished and update batch counters
	return nil
}

// BatchInfo contains information about a batch
type BatchInfo struct {
	ID           string `json:"id"`
	JobID        string `json:"job_id"`
	Exists       bool   `json:"exists"`
	TotalJobs    int    `json:"total_jobs"`
	PendingJobs  int    `json:"pending_jobs"`
	FailedJobs   int    `json:"failed_jobs"`
	Cancelled    bool   `json:"cancelled"`
	Finished     bool   `json:"finished"`
	AllowFailures bool  `json:"allow_failures"`
}

// Progress returns the batch progress percentage
func (bi BatchInfo) Progress() float64 {
	if bi.TotalJobs == 0 {
		return 0.0
	}
	processed := bi.TotalJobs - bi.PendingJobs
	return float64(processed) / float64(bi.TotalJobs) * 100.0
}

// BatchableInterface defines the interface for batchable jobs
type BatchableInterface interface {
	WithBatchID(batchID string) *Batchable
	WithJobID(jobID string) *Batchable
	GetBatchID() string
	GetJobID() string
	HasBatch() bool
	OnBatchSuccess(callback func(ctx context.Context, result interface{}) error) *Batchable
	OnBatchFailure(callback func(ctx context.Context, err error) error) *Batchable
	Batch(ctx context.Context) (BatchInfo, error)
	MarkJobAsFailed(ctx context.Context, err error) error
	MarkJobAsFinished(ctx context.Context, result interface{}) error
}
