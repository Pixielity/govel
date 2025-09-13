// Package checks provides built-in health check implementations.
package checks

import (
	"fmt"
	"runtime"

	"govel/packages/healthcheck/src/checks"
	"govel/packages/healthcheck/src/enums"
	"govel/packages/healthcheck/src/interfaces"
)

// MemoryUsageCheck monitors memory usage and alerts on configurable thresholds.
// It closely mirrors the Laravel health pattern for resource monitoring.
type MemoryUsageCheck struct {
	*checks.BaseCheck

	// warningThreshold is the memory usage percentage at which to issue warnings (0-100)
	warningThreshold int

	// errorThreshold is the memory usage percentage at which to fail the check (0-100)
	errorThreshold int
}

// NewMemoryUsageCheck creates a new memory usage check with default settings.
//
// Returns:
//
//	*MemoryUsageCheck: A new memory usage check instance
func NewMemoryUsageCheck() *MemoryUsageCheck {
	return &MemoryUsageCheck{
		BaseCheck:        checks.NewBaseCheck(),
		warningThreshold: 80, // 80% warning threshold
		errorThreshold:   90, // 90% error threshold
	}
}

// WarnWhenMemoryUsageIsAbove sets the warning threshold percentage.
//
// Parameters:
//
//	percentage: The percentage (0-100) at which to issue warnings
//
// Returns:
//
//	*MemoryUsageCheck: Self for method chaining
func (mc *MemoryUsageCheck) WarnWhenMemoryUsageIsAbove(percentage int) *MemoryUsageCheck {
	mc.warningThreshold = percentage
	return mc
}

// FailWhenMemoryUsageIsAbove sets the failure threshold percentage.
//
// Parameters:
//
//	percentage: The percentage (0-100) at which to fail the check
//
// Returns:
//
//	*MemoryUsageCheck: Self for method chaining
func (mc *MemoryUsageCheck) FailWhenMemoryUsageIsAbove(percentage int) *MemoryUsageCheck {
	mc.errorThreshold = percentage
	return mc
}

// Run performs the memory usage health check.
//
// Returns:
//
//	interfaces.ResultInterface: The health check result
func (mc *MemoryUsageCheck) Run() interfaces.ResultInterface {
	memoryUsagePercentage := mc.getMemoryUsagePercentage()

	result := checks.NewResult().
		SetMeta(map[string]interface{}{"memory_usage_percentage": memoryUsagePercentage}).
		SetShortSummary(fmt.Sprintf("%d%%", memoryUsagePercentage))

	if memoryUsagePercentage > mc.errorThreshold {
		return result.
			SetStatus(enums.StatusFailed).
			SetNotificationMessage(fmt.Sprintf("Memory usage is critical (%d%% used).", memoryUsagePercentage))
	}

	if memoryUsagePercentage > mc.warningThreshold {
		return result.
			SetStatus(enums.StatusWarning).
			SetNotificationMessage(fmt.Sprintf("Memory usage is high (%d%% used).", memoryUsagePercentage))
	}

	return result.SetStatus(enums.StatusOK)
}

// getMemoryUsagePercentage calculates the current memory usage percentage.
// This uses Go's runtime.MemStats to get heap usage information.
func (mc *MemoryUsageCheck) getMemoryUsagePercentage() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate usage percentage based on heap allocation vs system memory
	// Note: This is a simplified calculation. In production, you might want
	// to use more sophisticated methods to get system memory limits.

	// For now, we'll calculate based on heap in use vs heap system allocation
	if m.HeapSys == 0 {
		return 0
	}

	// Calculate percentage: (HeapInuse / HeapSys) * 100
	usagePercentage := int((float64(m.HeapInuse) / float64(m.HeapSys)) * 100)

	// Ensure we don't exceed 100%
	if usagePercentage > 100 {
		usagePercentage = 100
	}

	return usagePercentage
}

// Compile-time interface compliance check
var _ interfaces.CheckInterface = (*MemoryUsageCheck)(nil)
