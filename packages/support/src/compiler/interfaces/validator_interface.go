package interfaces

import "govel/support/compiler/types"

// ValidatorInterface defines the interface for validating Go code within the GoVel compiler system.
// Provides comprehensive analysis including syntax validation, security checks, and complexity analysis.
type ValidatorInterface interface {
	// Validate checks if the provided Go code is safe and valid for compilation.
	// Performs syntax validation, security analysis, and complexity assessment.
	//
	// Parameters:
	//
	//	code: Complete Go source code as a string
	//
	// Returns:
	//
	//	*types.ValidationResult: Comprehensive result with errors, warnings, and metrics
	Validate(code string) *types.ValidationResult

	// ValidateFile validates a Go source file by reading and analyzing its contents.
	// Extends Validate() with file-specific checks including path security and encoding validation.
	//
	// Parameters:
	//
	//	filePath: Path to a Go source file that exists and is readable
	//
	// Returns:
	//
	//	*types.ValidationResult: Comprehensive result with file-specific information
	ValidateFile(filePath string) *types.ValidationResult
}
