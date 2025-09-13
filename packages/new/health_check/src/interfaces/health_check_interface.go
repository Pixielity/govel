// Package interfaces provides the core interfaces for the GoVel health check system.
// These interfaces define contracts for health checks, results, and registry management.
package interfaces

import (
	"time"
)

// CheckInterface defines the contract that all health checks must implement.
// This interface provides the foundation for creating standardized health checks
// that can be registered and executed by the health check system.
//
// Example implementation:
//
//	type DatabaseCheck struct {
//		connectionString string
//		timeout          time.Duration
//	}
//
//	func (d *DatabaseCheck) Run(ctx context.Context) ResultInterface {
//		// Implementation here
//	}
//
//	func (d *DatabaseCheck) GetName() string {
//		return "database-check"
//	}
type CheckInterface interface {
	// Run executes the health check and returns a result.
	// The context can be used for timeouts and cancellation.
	//
	// Parameters:
	//   ctx: Context for timeout and cancellation control
	//
	// Returns:
	//   ResultInterface: The result of the health check execution
	Run() ResultInterface

	// GetName returns the unique name identifier for this health check.
	// This name is used for registration, logging, and result identification.
	//
	// Returns:
	//   string: Unique name for this health check
	GetName() string
}

// ResultInterface defines the contract for health check results.
// Results contain the status, messages, metadata, and timing information
// from health check executions.
type ResultInterface interface {
	// GetStatus returns the status of the health check.
	//
	// Returns:
	//   StatusInterface: The status enum value
	GetStatus() StatusInterface

	// SetStatus sets the status of the health check result.
	//
	// Parameters:
	//   status: The status to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetStatus(status StatusInterface) ResultInterface

	// GetNotificationMessage returns the message to be used in notifications.
	//
	// Returns:
	//   string: The notification message
	GetNotificationMessage() string

	// SetNotificationMessage sets the notification message.
	//
	// Parameters:
	//   message: The notification message to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetNotificationMessage(message string) ResultInterface

	// GetShortSummary returns a brief summary of the result.
	//
	// Returns:
	//   string: Brief summary text
	GetShortSummary() string

	// SetShortSummary sets a brief summary of the result.
	//
	// Parameters:
	//   summary: Brief summary text to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetShortSummary(summary string) ResultInterface

	// GetMeta returns metadata associated with the result.
	//
	// Returns:
	//   map[string]interface{}: Metadata key-value pairs
	GetMeta() map[string]interface{}

	// SetMeta sets metadata for the result.
	//
	// Parameters:
	//   meta: Metadata key-value pairs to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetMeta(meta map[string]interface{}) ResultInterface

	// GetCheck returns the health check that produced this result.
	//
	// Returns:
	//   CheckInterface: The originating health check
	GetCheck() CheckInterface

	// SetCheck sets the health check that produced this result.
	//
	// Parameters:
	//   check: The health check to associate with this result
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetCheck(check CheckInterface) ResultInterface

	// GetStartedAt returns when the health check started executing.
	//
	// Returns:
	//   *time.Time: Start time, nil if not set
	GetStartedAt() *time.Time

	// SetStartedAt sets when the health check started executing.
	//
	// Parameters:
	//   startedAt: The start time to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetStartedAt(startedAt time.Time) ResultInterface

	// GetEndedAt returns when the health check finished executing.
	//
	// Returns:
	//   *time.Time: End time, nil if not set
	GetEndedAt() *time.Time

	// SetEndedAt sets when the health check finished executing.
	//
	// Parameters:
	//   endedAt: The end time to set
	//
	// Returns:
	//   ResultInterface: Self for method chaining
	SetEndedAt(endedAt time.Time) ResultInterface

	// GetDuration returns the execution duration of the health check.
	//
	// Returns:
	//   time.Duration: Execution duration, 0 if times not set
	GetDuration() time.Duration
}

// StatusInterface defines the contract for health check status enums.
// This interface allows for different status implementations while
// maintaining consistent behavior.
type StatusInterface interface {
	// String returns the string representation of the status.
	//
	// Returns:
	//   string: Status as string (e.g., "ok", "warning", "failed")
	String() string

	// IsHealthy returns true if the status represents a healthy state.
	//
	// Returns:
	//   bool: true if healthy (ok), false otherwise
	IsHealthy() bool

	// IsWarning returns true if the status represents a warning state.
	//
	// Returns:
	//   bool: true if warning, false otherwise
	IsWarning() bool

	// IsFailed returns true if the status represents a failed state.
	//
	// Returns:
	//   bool: true if failed or crashed, false otherwise
	IsFailed() bool

	// GetSeverityLevel returns a numeric severity level for sorting/comparison.
	// Lower numbers indicate better health (0=ok, 1=warning, 2=failed, 3=crashed).
	//
	// Returns:
	//   int: Severity level for comparison
	GetSeverityLevel() int
}

// CheckResultsInterface defines the contract for collections of health check results.
// This interface provides methods for analyzing and manipulating groups of results.
type CheckResultsInterface interface {
	// GetResults returns all individual health check results.
	//
	// Returns:
	//   []ResultInterface: Slice of all results
	GetResults() []ResultInterface

	// AddResult adds a result to the collection.
	//
	// Parameters:
	//   result: The result to add
	//
	// Returns:
	//   CheckResultsInterface: Self for method chaining
	AddResult(result ResultInterface) CheckResultsInterface

	// GetResultByName returns a result by the name of its associated check.
	//
	// Parameters:
	//   name: The name of the check to find
	//
	// Returns:
	//   ResultInterface: The result if found, nil otherwise
	GetResultByName(name string) ResultInterface

	// ContainsFailingCheck returns true if any result has a failed status.
	//
	// Returns:
	//   bool: true if any check failed, false otherwise
	ContainsFailingCheck() bool

	// ContainsWarningCheck returns true if any result has a warning status.
	//
	// Returns:
	//   bool: true if any check has warnings, false otherwise
	ContainsWarningCheck() bool

	// GetOverallStatus returns the worst status among all results.
	//
	// Returns:
	//   StatusInterface: The worst status found
	GetOverallStatus() StatusInterface

	// FilterByStatus returns results that match the specified status.
	//
	// Parameters:
	//   status: The status to filter by
	//
	// Returns:
	//   []ResultInterface: Results matching the status
	FilterByStatus(status StatusInterface) []ResultInterface

	// ToJSON converts the results to JSON format.
	//
	// Returns:
	//   string: JSON representation of results
	//   error: Any error during JSON marshaling
	ToJSON() (string, error)

	// GetTotalDuration returns the sum of all check durations.
	//
	// Returns:
	//   time.Duration: Total execution time for all checks
	GetTotalDuration() time.Duration

	// GetExecutedAt returns when the checks were executed.
	//
	// Returns:
	//   *time.Time: Execution timestamp, nil if not set
	GetExecutedAt() *time.Time

	// SetExecutedAt sets when the checks were executed.
	//
	// Parameters:
	//   executedAt: The execution timestamp to set
	//
	// Returns:
	//   CheckResultsInterface: Self for method chaining
	SetExecutedAt(executedAt time.Time) CheckResultsInterface
}
