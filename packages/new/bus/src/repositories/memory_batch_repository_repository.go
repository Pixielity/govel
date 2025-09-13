package repositories

import (
	"context"
	"fmt"
	"sync"
	"time"

	"govel/packages/bus/src/interfaces"
)

// MemoryBatchRepository implements BatchRepository using in-memory storage
type MemoryBatchRepository struct {
	mu      sync.RWMutex
	batches map[string]interfaces.BatchData
}

// NewMemoryBatchRepository creates a new memory-based batch repository
func NewMemoryBatchRepository() *MemoryBatchRepository {
	return &MemoryBatchRepository{
		batches: make(map[string]interfaces.BatchData),
	}
}

// Get retrieves a batch by its ID
func (r *MemoryBatchRepository) Get(ctx context.Context, batchID string) (interfaces.BatchData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return interfaces.BatchData{}, ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return interfaces.BatchData{}, fmt.Errorf("batch %s not found", batchID)
	}

	return batch, nil
}

// Store stores a new batch
func (r *MemoryBatchRepository) Store(ctx context.Context, pendingBatch interfaces.PendingBatchData) (interfaces.BatchData, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return interfaces.BatchData{}, ctx.Err()
	default:
	}

	// Check if batch already exists
	if _, exists := r.batches[pendingBatch.ID]; exists {
		return interfaces.BatchData{}, fmt.Errorf("batch %s already exists", pendingBatch.ID)
	}

	// Create batch data from pending batch
	batchData := interfaces.BatchData{
		ID:           pendingBatch.ID,
		Name:         pendingBatch.Name,
		TotalJobs:    len(pendingBatch.Jobs),
		PendingJobs:  len(pendingBatch.Jobs),
		FailedJobs:   0,
		FailedJobIDs: make([]string, 0),
		Options:      pendingBatch.Options,
		CreatedAt:    pendingBatch.CreatedAt,
		CancelledAt:  nil,
		FinishedAt:   nil,
	}

	r.batches[batchData.ID] = batchData
	return batchData, nil
}

// IncrementTotalJobs increments the total job count for a batch
func (r *MemoryBatchRepository) IncrementTotalJobs(ctx context.Context, batchID string, amount int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	batch.TotalJobs += amount
	batch.PendingJobs += amount // New jobs start as pending
	r.batches[batchID] = batch

	return nil
}

// DecrementPendingJobs decrements the pending job count for a batch
func (r *MemoryBatchRepository) DecrementPendingJobs(ctx context.Context, batchID string, amount int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	batch.PendingJobs -= amount
	if batch.PendingJobs < 0 {
		batch.PendingJobs = 0
	}

	r.batches[batchID] = batch
	return nil
}

// IncrementFailedJobs increments the failed job count and adds job ID
func (r *MemoryBatchRepository) IncrementFailedJobs(ctx context.Context, batchID string, jobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	batch.FailedJobs++
	batch.FailedJobIDs = append(batch.FailedJobIDs, jobID)
	r.batches[batchID] = batch

	return nil
}

// MarkAsFinished marks a batch as finished
func (r *MemoryBatchRepository) MarkAsFinished(ctx context.Context, batchID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	if batch.CancelledAt != nil {
		return fmt.Errorf("batch %s is already cancelled", batchID)
	}

	if batch.FinishedAt != nil {
		return fmt.Errorf("batch %s is already finished", batchID)
	}

	now := time.Now()
	batch.FinishedAt = &now
	r.batches[batchID] = batch

	return nil
}

// Cancel cancels a batch
func (r *MemoryBatchRepository) Cancel(ctx context.Context, batchID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	batch, exists := r.batches[batchID]
	if !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	if batch.FinishedAt != nil {
		return fmt.Errorf("batch %s is already finished", batchID)
	}

	if batch.CancelledAt != nil {
		return fmt.Errorf("batch %s is already cancelled", batchID)
	}

	now := time.Now()
	batch.CancelledAt = &now
	r.batches[batchID] = batch

	return nil
}

// Delete deletes a batch
func (r *MemoryBatchRepository) Delete(ctx context.Context, batchID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := r.batches[batchID]; !exists {
		return fmt.Errorf("batch %s not found", batchID)
	}

	delete(r.batches, batchID)
	return nil
}

// Transaction executes a function within a transaction
// For memory repository, we use a simple mutex lock
func (r *MemoryBatchRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return fn(ctx)
}

// GetAll returns all batches (for testing/debugging)
func (r *MemoryBatchRepository) GetAll(ctx context.Context) ([]interfaces.BatchData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	batches := make([]interfaces.BatchData, 0, len(r.batches))
	for _, batch := range r.batches {
		batches = append(batches, batch)
	}

	return batches, nil
}

// Count returns the total number of batches
func (r *MemoryBatchRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	return len(r.batches), nil
}

// Clear removes all batches (for testing)
func (r *MemoryBatchRepository) Clear(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.batches = make(map[string]interfaces.BatchData)
	return nil
}

// GetBatchesByStatus returns batches filtered by status
func (r *MemoryBatchRepository) GetBatchesByStatus(ctx context.Context, finished, cancelled bool) ([]interfaces.BatchData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	batches := make([]interfaces.BatchData, 0)
	for _, batch := range r.batches {
		isFinished := batch.FinishedAt != nil
		isCancelled := batch.CancelledAt != nil

		if finished && isFinished && !isCancelled {
			batches = append(batches, batch)
		} else if cancelled && isCancelled {
			batches = append(batches, batch)
		} else if !finished && !cancelled && !isFinished && !isCancelled {
			batches = append(batches, batch)
		}
	}

	return batches, nil
}

// Ensure MemoryBatchRepository implements BatchRepository interface
var _ interfaces.BatchRepository = (*MemoryBatchRepository)(nil)
