package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResultStatus represents the status of a compilation result.
// Provides high-level categorization of compilation outcome for error handling.
type ResultStatus string

const (
	// StatusSuccess indicates successful compilation and execution.
	StatusSuccess ResultStatus = "success"

	// StatusCompileError indicates a compilation error occurred.
	StatusCompileError ResultStatus = "compile_error"

	// StatusRuntimeError indicates a runtime error during execution.
	StatusRuntimeError ResultStatus = "runtime_error"

	// StatusTimeout indicates the operation exceeded the timeout limit.
	StatusTimeout ResultStatus = "timeout"

	// StatusValidationError indicates the code failed validation checks.
	StatusValidationError ResultStatus = "validation_error"

	// StatusCacheHit indicates the result was retrieved from cache.
	StatusCacheHit ResultStatus = "cache_hit"
)

// Result represents the comprehensive outcome of a Go compilation and execution operation.
// Contains status, output content, error details, performance metrics, and metadata.
type Result struct {
	// Status indicates the overall result status using predefined constants.
	Status ResultStatus `json:"status"`

	// Success is a boolean flag for quick success/failure checking.
	Success bool `json:"success"`

	// Content contains the output content from the executed program.
	Content []byte `json:"content,omitempty"`

	// ErrorOutput contains any error messages from stderr.
	ErrorOutput []byte `json:"error_output,omitempty"`

	// CompilationError contains detailed compilation error information.
	CompilationError string `json:"compilation_error,omitempty"`

	// RuntimeError contains runtime error details if execution failed.
	RuntimeError string `json:"runtime_error,omitempty"`

	// ExitCode is the exit code returned by the executed program.
	ExitCode int `json:"exit_code"`

	// Duration represents the total time taken for compilation and execution.
	Duration time.Duration `json:"duration"`

	// CompileTime represents the time spent on compilation only.
	CompileTime time.Duration `json:"compile_time"`

	// ExecutionTime represents the time spent on program execution.
	ExecutionTime time.Duration `json:"execution_time"`

	// MemoryUsed represents the peak memory usage during compilation/execution.
	MemoryUsed int64 `json:"memory_used"`

	// CacheHit indicates if this result was retrieved from cache.
	CacheHit bool `json:"cache_hit"`

	// FilePath is the path of the compiled file (if applicable).
	FilePath string `json:"file_path,omitempty"`

	// Hash is a unique identifier for this compilation (used for caching).
	Hash string `json:"hash,omitempty"`

	// Timestamp indicates when this result was generated.
	Timestamp time.Time `json:"timestamp"`

	// Metadata contains additional contextual information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewResult creates a new Result with default values and current timestamp.
// Initializes Status as success, timestamp as now, and empty metadata map.
//
// Returns:
//
//	*Result: A new Result instance ready for population
func NewResult() *Result {
	return &Result{
		Status:    StatusSuccess,
		Success:   true,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NewErrorResult creates a Result representing an error condition.
//
// This function creates a pre-configured Result for error scenarios,
// setting appropriate status and error information.
//
// Parameters:
//
//	status: The specific error status to set
//	err: The error that occurred
//
// Returns:
//
//	*Result: A Result configured for the error condition
func NewErrorResult(status ResultStatus, err error) *Result {
	result := NewResult()
	result.Status = status
	result.Success = false
	result.RuntimeError = err.Error()
	return result
}

// NewCompileErrorResult creates a Result for compilation errors.
//
// This function creates a specialized Result for compilation failures,
// setting the appropriate status and error fields.
//
// Parameters:
//
//	err: The compilation error that occurred
//
// Returns:
//
//	*Result: A Result configured for compilation failure
func NewCompileErrorResult(err error) *Result {
	result := NewResult()
	result.Status = StatusCompileError
	result.Success = false
	result.CompilationError = err.Error()
	return result
}

// GetContent returns the content as a string, providing a convenient accessor.
//
// This method converts the binary Content field to a string for easy
// access when the content is known to be text-based.
//
// Returns:
//
//	string: The content as a UTF-8 string
func (r *Result) GetContent() string {
	return string(r.Content)
}

// GetError returns the appropriate error message based on the result status.
//
// This method provides a unified way to access error information regardless
// of the specific error type, simplifying error handling logic.
//
// Returns:
//
//	string: The most relevant error message, or empty string if no error
func (r *Result) GetError() string {
	switch r.Status {
	case StatusCompileError:
		return r.CompilationError
	case StatusRuntimeError:
		return r.RuntimeError
	case StatusTimeout:
		return "Operation timed out"
	case StatusValidationError:
		return "Code validation failed"
	default:
		if len(r.ErrorOutput) > 0 {
			return string(r.ErrorOutput)
		}
		return ""
	}
}

// HasError returns true if the result represents any kind of error condition.
//
// This method provides a simple boolean check for error conditions,
// useful in conditional logic and error handling.
//
// Returns:
//
//	bool: true if any error occurred, false for successful operations
func (r *Result) HasError() bool {
	return !r.Success || r.Status != StatusSuccess
}

// SetMetadata sets a metadata value with the given key.
//
// This method provides a safe way to add metadata to the result,
// ensuring the metadata map is initialized if needed.
//
// Parameters:
//
//	key: The metadata key
//	value: The metadata value (any JSON-serializable type)
func (r *Result) SetMetadata(key string, value interface{}) {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
}

// GetMetadata retrieves a metadata value by key, returning nil if not found.
//
// This method provides safe access to metadata values without panicking
// if the key doesn't exist or the metadata map is nil.
//
// Parameters:
//
//	key: The metadata key to retrieve
//
// Returns:
//
//	interface{}: The metadata value, or nil if not found
func (r *Result) GetMetadata(key string) interface{} {
	if r.Metadata == nil {
		return nil
	}
	return r.Metadata[key]
}

// String provides a human-readable string representation of the Result.
//
// This method implements the Stringer interface and provides a concise
// summary of the result for logging and debugging purposes.
//
// Returns:
//
//	string: A formatted summary of the result
func (r *Result) String() string {
	if r.Success {
		return fmt.Sprintf("Result{Status: %s, Duration: %v, Content: %d bytes}",
			r.Status, r.Duration, len(r.Content))
	}
	return fmt.Sprintf("Result{Status: %s, Error: %s, Duration: %v}",
		r.Status, r.GetError(), r.Duration)
}

// ToJSON serializes the Result to JSON format.
//
// This method provides formatted JSON output for the result, useful for
// logging, API responses, and persistence.
//
// Returns:
//
//	[]byte: Pretty-printed JSON representation
//	error: Serialization error, if any
func (r *Result) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// FromJSON deserializes a Result from JSON format.
//
// This function reconstructs a Result instance from JSON data,
// useful for loading persisted results or API communication.
//
// Parameters:
//
//	data: JSON-encoded result data
//
// Returns:
//
//	*Result: The deserialized result
//	error: Deserialization error, if any
func FromJSON(data []byte) (*Result, error) {
	var result Result
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	return &result, nil
}
