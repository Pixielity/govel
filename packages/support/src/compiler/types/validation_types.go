// Package types provides supporting types for the GoVel compiler system.
//
// This file defines cache entries, metrics, and validation result types
// that support the main compiler functionality with comprehensive documentation
// and detailed field descriptions following Go best practices.
package types

import (
	"time"
)

// ValidationResult represents the comprehensive result of code validation.
// Contains validation errors, warnings, security issues, and complexity analysis.
type ValidationResult struct {
	// Valid indicates if the code passed all validation checks.
	Valid bool `json:"valid"`

	// Errors contains any validation errors found.
	Errors []string `json:"errors,omitempty"`

	// Warnings contains any validation warnings.
	Warnings []string `json:"warnings,omitempty"`

	// UnsafeImports lists any potentially dangerous imports found.
	UnsafeImports []string `json:"unsafe_imports,omitempty"`

	// ComplexityScore is a rough measure of code complexity (0-100).
	ComplexityScore int `json:"complexity_score"`

	// ValidationTime is the time taken to perform validation.
	ValidationTime time.Duration `json:"validation_time"`
}

// NewValidationResult creates a new ValidationResult with default values.
// Initializes with valid status and empty slices for errors, warnings, and unsafe imports.
//
// Returns:
//
//	*ValidationResult: A new validation result ready for use
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:         true,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
		UnsafeImports: make([]string, 0),
	}
}

// AddError adds an error to the validation result and marks it as invalid.
// Automatically sets valid flag to false since any error invalidates code.
//
// Parameters:
//
//	err: The error message to add
func (v *ValidationResult) AddError(err string) {
	v.Valid = false
	v.Errors = append(v.Errors, err)
}

// AddWarning adds a warning to the validation result.
// Does not affect valid status since warnings are non-blocking.
//
// Parameters:
//
//	warning: The warning message to add
func (v *ValidationResult) AddWarning(warning string) {
	v.Warnings = append(v.Warnings, warning)
}

// AddUnsafeImport adds an unsafe import to the list.
// Records a potentially dangerous import detected during validation.
//
// Parameters:
//
//	importPath: The unsafe import path that was detected
func (v *ValidationResult) AddUnsafeImport(importPath string) {
	v.UnsafeImports = append(v.UnsafeImports, importPath)
}

// HasErrors returns true if there are any validation errors.
// Provides convenient check for presence of validation errors.
//
// Returns:
//
//	bool: true if any errors exist, false otherwise
func (v *ValidationResult) HasErrors() bool {
	return len(v.Errors) > 0
}

// HasWarnings returns true if there are any validation warnings.
// Provides convenient check for presence of validation warnings.
//
// Returns:
//
//	bool: true if any warnings exist, false otherwise
func (v *ValidationResult) HasWarnings() bool {
	return len(v.Warnings) > 0
}

// HasUnsafeImports returns true if any unsafe imports were detected.
// Provides convenient check for presence of potentially dangerous imports.
//
// Returns:
//
//	bool: true if any unsafe imports exist, false otherwise
func (v *ValidationResult) HasUnsafeImports() bool {
	return len(v.UnsafeImports) > 0
}
