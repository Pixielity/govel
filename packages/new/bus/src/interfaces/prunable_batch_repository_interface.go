package interfaces

import (
	"context"
	"time"
)

// PrunableBatchRepository extends BatchRepository with pruning capabilities
type PrunableBatchRepository interface {
	BatchRepository

	// Prune removes batches that are older than the specified time
	Prune(ctx context.Context, before time.Time) (int, error)

	// PruneUnfinished removes unfinished batches older than the specified time
	PruneUnfinished(ctx context.Context, before time.Time) (int, error)

	// PruneCancelled removes cancelled batches older than the specified time
	PruneCancelled(ctx context.Context, before time.Time) (int, error)

	// CountBatches returns the total number of batches
	CountBatches(ctx context.Context) (int, error)

	// CountFinished returns the number of finished batches
	CountFinished(ctx context.Context) (int, error)

	// CountCancelled returns the number of cancelled batches
	CountCancelled(ctx context.Context) (int, error)

	// CountPending returns the number of pending batches
	CountPending(ctx context.Context) (int, error)

	// GetOldest gets the oldest batches up to the specified limit
	GetOldest(ctx context.Context, limit int) ([]BatchData, error)

	// GetBatchesBefore gets batches created before the specified time
	GetBatchesBefore(ctx context.Context, before time.Time, limit int) ([]BatchData, error)
}

// PruneOptions contains options for batch pruning
type PruneOptions struct {
	// Hours specifies how old batches should be before pruning
	Hours int

	// KeepFinished whether to keep finished batches
	KeepFinished bool

	// KeepCancelled whether to keep cancelled batches
	KeepCancelled bool

	// KeepUnfinished whether to keep unfinished batches
	KeepUnfinished bool

	// MaxBatches maximum number of batches to keep
	MaxBatches int

	// DryRun if true, only count what would be pruned without actually pruning
	DryRun bool
}

// PruneResult contains the results of a pruning operation
type PruneResult struct {
	// TotalPruned is the total number of batches that were pruned
	TotalPruned int

	// FinishedPruned is the number of finished batches that were pruned
	FinishedPruned int

	// CancelledPruned is the number of cancelled batches that were pruned
	CancelledPruned int

	// UnfinishedPruned is the number of unfinished batches that were pruned
	UnfinishedPruned int

	// ErrorsEncountered is the number of errors encountered during pruning
	ErrorsEncountered int

	// Duration is how long the pruning operation took
	Duration time.Duration

	// LastError is the last error encountered, if any
	LastError error
}

// PruneStrategy defines different pruning strategies
type PruneStrategy string

const (
	// PruneAll removes all batches older than the specified time
	PruneAll PruneStrategy = "all"

	// PruneFinished removes only finished batches
	PruneFinished PruneStrategy = "finished"

	// PruneCancelled removes only cancelled batches
	PruneCancelled PruneStrategy = "cancelled"

	// PruneUnfinished removes only unfinished batches
	PruneUnfinished PruneStrategy = "unfinished"

	// PruneByCount removes oldest batches to maintain a maximum count
	PruneByCount PruneStrategy = "count"
)

// Advanced pruning interface for more complex pruning scenarios
type AdvancedPrunableBatchRepository interface {
	PrunableBatchRepository

	// PruneWithStrategy prunes batches using the specified strategy
	PruneWithStrategy(ctx context.Context, strategy PruneStrategy, options PruneOptions) (PruneResult, error)

	// PruneWithCallback prunes batches and calls the callback for each batch
	PruneWithCallback(ctx context.Context, before time.Time, callback func(batch BatchData) bool) (PruneResult, error)

	// GetPruningCandidates returns batches that would be pruned without actually pruning them
	GetPruningCandidates(ctx context.Context, options PruneOptions) ([]BatchData, error)
}
