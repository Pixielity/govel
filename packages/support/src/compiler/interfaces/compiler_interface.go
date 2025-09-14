package interfaces

import (
	"context"

	"govel/support/compiler/types"
)

// CompilerInterface defines the main interface for Go code compilation and execution within the GoVel framework.
// Provides a clean API for compiling Go files and code strings with caching, metrics, and context support.
type CompilerInterface interface {
	// Compile compiles and executes a Go file, returning the execution result.
	// Reads, compiles, and executes the Go source file with automatic cleanup.
	//
	// Parameters:
	//
	//	filePath: The path to the Go source file to compile and execute
	//
	// Returns:
	//
	//	*types.Result: Comprehensive result with output, errors, and metrics
	Compile(filePath string) *types.Result

	// CompileWithContext compiles and executes a Go file with context support for cancellation and timeout control.
	// Extends Compile() with context-aware operation management throughout the compilation pipeline.
	//
	// Parameters:
	//
	//	ctx: The context for controlling compilation lifecycle, cancellation, and timeouts
	//	filePath: The path to the Go source file to compile and execute
	//
	// Returns:
	//
	//	*types.Result: Comprehensive result, with StatusTimeout if context cancelled/timed out
	CompileWithContext(ctx context.Context, filePath string) *types.Result

	// CompileCode compiles and executes Go code from a string, returning the execution result.
	// Creates temporary files, compiles the code, and executes with automatic cleanup.
	//
	// Parameters:
	//
	//	code: Complete Go source code string with main package and main function
	//
	// Returns:
	//
	//	*types.Result: Comprehensive result with execution details and metrics
	CompileCode(code string) *types.Result

	// CompileCodeWithContext compiles and executes Go code from a string with context support.
	// Combines CompileCode() functionality with context-aware cancellation and timeout handling.
	//
	// Parameters:
	//
	//	ctx: The context for controlling compilation lifecycle and request-scoped values
	//	code: Complete Go source code string with main package and main function
	//
	// Returns:
	//
	//	*types.Result: Comprehensive result with context-aware error handling
	CompileCodeWithContext(ctx context.Context, code string) *types.Result

	// GetConfig returns the current compiler configuration.
	// Provides read-only access to all active configuration settings.
	//
	// Returns:
	//
	//	*types.Config: A copy of the current configuration
	GetConfig() *types.Config

	// UpdateConfig updates the compiler configuration with new settings.
	// Merges new configuration with current settings, applying non-zero values.
	//
	// Parameters:
	//
	//	config: New configuration settings to apply
	//
	// Returns:
	//
	//	error: An error if the configuration is invalid or compiler is closed
	UpdateConfig(config *types.Config) error

	// Close cleans up compiler resources and shuts down the compiler gracefully.
	// Removes temporary directories, clears cache, and cancels in-progress operations.
	//
	// Returns:
	//
	//	error: An error if cleanup operations fail
	Close() error

	// GetMetrics returns current compilation metrics and statistics.
	// Provides comprehensive performance data including success rates, timing, and cache effectiveness.
	//
	// Returns:
	//
	//	*types.Metrics: Comprehensive metrics object with all collected statistics
	GetMetrics() *types.Metrics
}
