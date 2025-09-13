package interfaces

import (
	"context"
	"time"
)

// PipelineContextInterface defines the contract for pipeline execution context.
// This interface extends the standard Go context.Context with pipeline-specific
// functionality for tracking execution state, timing, and current pipe information.
//
// Key features:
//   - Standard Go context functionality (Done, Err, Value)
//   - Execution timing and performance tracking
//   - Current pipe tracking for debugging
//   - Context value management for pipeline state
//   - Cancellation and timeout support
type PipelineContextInterface interface {
	// Embed standard Go context for cancellation, deadlines, and values
	context.Context

	// SetExecutionStartTime sets the pipeline execution start time.
	// Used for performance tracking and timeout calculations.
	//
	// Parameters:
	//   - startTime: The time when pipeline execution began
	//
	// Returns:
	//   - PipelineContextInterface: Returns self for method chaining
	SetExecutionStartTime(startTime time.Time) PipelineContextInterface

	// GetExecutionStartTime returns the pipeline execution start time.
	// Returns zero time if not set.
	//
	// Returns:
	//   - time.Time: The pipeline execution start time
	GetExecutionStartTime() time.Time

	// SetCurrentPipe sets the name of the currently executing pipe.
	// Used for debugging and error reporting.
	//
	// Parameters:
	//   - pipeName: The name or identifier of the current pipe
	//
	// Returns:
	//   - PipelineContextInterface: Returns self for method chaining
	SetCurrentPipe(pipeName string) PipelineContextInterface

	// GetCurrentPipe returns the name of the currently executing pipe.
	// Returns empty string if not set.
	//
	// Returns:
	//   - string: The name of the current pipe
	GetCurrentPipe() string

	// GetExecutionDuration returns how long the pipeline has been executing.
	// Returns zero duration if execution hasn't started.
	//
	// Returns:
	//   - time.Duration: The elapsed execution time
	GetExecutionDuration() time.Duration
}
