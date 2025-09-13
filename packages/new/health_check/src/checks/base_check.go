// Package checks provides base functionality and common utilities for health checks.
// This package contains the base check implementation that mirrors Laravel's Check class.
package checks

import (
	"reflect"
	"strings"

	"govel/healthcheck/src/enums"
	"govel/healthcheck/src/interfaces"
)

// ConditionFunc represents a condition function that returns whether a check should run
type ConditionFunc func() bool

// BaseCheck provides common functionality for health check implementations.
// It closely mirrors the Laravel Health Check class pattern.
type BaseCheck struct {
	// expression is the cron-like schedule expression (simplified)
	expression string

	// name is the unique identifier for this health check
	name *string

	// label is the display label for this health check
	label *string

	// timezone for scheduling (simplified to just store the name)
	timezone string

	// shouldRun contains conditions that determine if the check should execute
	shouldRun []interface{} // can be bool or ConditionFunc
}

// NewBaseCheck creates a new BaseCheck instance with default values.
//
// Returns:
//
//	*BaseCheck: A new base check instance
func NewBaseCheck() *BaseCheck {
	return &BaseCheck{
		expression: "* * * * *", // Every minute (simplified)
		name:       nil,
		label:      nil,
		timezone:   "UTC",
		shouldRun:  make([]interface{}, 0),
	}
}

// Name sets the name of the health check.
//
// Parameters:
//
//	name: The name to set
//
// Returns:
//
//	*BaseCheck: Self for method chaining
func (bc *BaseCheck) Name(name string) *BaseCheck {
	bc.name = &name
	return bc
}

// Label sets the display label for the health check.
//
// Parameters:
//
//	label: The label to set
//
// Returns:
//
//	*BaseCheck: Self for method chaining
func (bc *BaseCheck) Label(label string) *BaseCheck {
	bc.label = &label
	return bc
}

// GetLabel returns the display label for the health check.
// If no label is set, it generates one from the name.
//
// Returns:
//
//	string: The display label
func (bc *BaseCheck) GetLabel() string {
	if bc.label != nil {
		return *bc.label
	}

	name := bc.GetName()
	// Convert snake_case to Title Case (simplified version)
	words := strings.Split(strings.ReplaceAll(name, "_", " "), " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// GetName returns the name of the health check.
// If no name is set, it generates one from the struct type name.
//
// Returns:
//
//	string: The check name
func (bc *BaseCheck) GetName() string {
	if bc.name != nil {
		return *bc.name
	}

	// Get the struct type name and remove "Check" suffix
	// This is a simplified version of Laravel's class_basename logic
	typeName := reflect.TypeOf(bc).Elem().Name()
	typeName = strings.TrimSuffix(typeName, "Check")
	return typeName
}

// GetRunConditions returns the slice of run conditions.
//
// Returns:
//
//	[]interface{}: Slice of conditions (bool or ConditionFunc)
func (bc *BaseCheck) GetRunConditions() []interface{} {
	return bc.shouldRun
}

// ShouldRun determines if this check should be executed based on its conditions.
// This mirrors Laravel's shouldRun method logic.
//
// Returns:
//
//	bool: true if the check should run, false otherwise
func (bc *BaseCheck) ShouldRun() bool {
	for _, condition := range bc.shouldRun {
		var shouldRun bool

		switch v := condition.(type) {
		case bool:
			shouldRun = v
		case ConditionFunc:
			shouldRun = v()
		case func() bool:
			shouldRun = v()
		default:
			shouldRun = false
		}

		if !shouldRun {
			return false
		}
	}

	// Always return true for now (simplified scheduling)
	// In a full implementation, you would check the cron expression against current time
	return true
}

// If adds a condition that must be true for the check to run.
//
// Parameters:
//
//	condition: A boolean value or function that returns boolean
//
// Returns:
//
//	*BaseCheck: Self for method chaining
func (bc *BaseCheck) If(condition interface{}) *BaseCheck {
	bc.shouldRun = append(bc.shouldRun, condition)
	return bc
}

// Unless adds a condition that must be false for the check to run.
//
// Parameters:
//
//	condition: A boolean value or function that returns boolean
//
// Returns:
//
//	*BaseCheck: Self for method chaining
func (bc *BaseCheck) Unless(condition interface{}) *BaseCheck {
	// Convert "unless" to "if not"
	switch v := condition.(type) {
	case bool:
		bc.shouldRun = append(bc.shouldRun, !v)
	case ConditionFunc:
		bc.shouldRun = append(bc.shouldRun, ConditionFunc(func() bool {
			return !v()
		}))
	case func() bool:
		bc.shouldRun = append(bc.shouldRun, ConditionFunc(func() bool {
			return !v()
		}))
	}
	return bc
}

// MarkAsCrashed creates a result marked as crashed.
// This mirrors Laravel's markAsCrashed method.
//
// Returns:
//
//	interfaces.ResultInterface: A result with crashed status
func (bc *BaseCheck) MarkAsCrashed() interfaces.ResultInterface {
	return types.NewResult().SetStatus(enums.StatusCrashed)
}

// OnTerminate is called when the check terminates (placeholder for Laravel compatibility).
//
// Parameters:
//
//	request: The request object (unused in Go implementation)
//	response: The response object (unused in Go implementation)
func (bc *BaseCheck) OnTerminate(request, response interface{}) {
	// Placeholder method for Laravel compatibility
	// In Laravel, this is called during application termination
}
