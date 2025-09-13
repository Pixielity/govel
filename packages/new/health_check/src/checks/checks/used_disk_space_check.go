// Package checks provides built-in health check implementations.
package checks

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"govel/healthcheck/src/checks"
	"govel/healthcheck/src/enums"
	"govel/healthcheck/src/interfaces"
)

// UsedDiskSpaceCheck monitors disk space usage and alerts on configurable thresholds.
// It closely mirrors the Laravel health UsedDiskSpaceCheck pattern.
type UsedDiskSpaceCheck struct {
	*checks.BaseCheck

	// warningThreshold is the percentage at which to issue warnings (0-100)
	warningThreshold int

	// errorThreshold is the percentage at which to fail the check (0-100)
	errorThreshold int

	// filesystemName is the filesystem/path to check (nil means current directory)
	filesystemName *string
}

// NewUsedDiskSpaceCheck creates a new disk space check instance with default settings.
//
// Returns:
//
//	*UsedDiskSpaceCheck: A new disk space check instance
func NewUsedDiskSpaceCheck() *UsedDiskSpaceCheck {
	return &UsedDiskSpaceCheck{
		BaseCheck:        checks.NewBaseCheck(),
		warningThreshold: 70,
		errorThreshold:   90,
		filesystemName:   nil,
	}
}

// FilesystemName sets the filesystem or path to check.
//
// Parameters:
//
//	filesystemName: The filesystem or path to check
//
// Returns:
//
//	*UsedDiskSpaceCheck: Self for method chaining
func (dsc *UsedDiskSpaceCheck) FilesystemName(filesystemName string) *UsedDiskSpaceCheck {
	dsc.filesystemName = &filesystemName
	return dsc
}

// WarnWhenUsedSpaceIsAbovePercentage sets the warning threshold percentage.
//
// Parameters:
//
//	percentage: The percentage (0-100) at which to issue warnings
//
// Returns:
//
//	*UsedDiskSpaceCheck: Self for method chaining
func (dsc *UsedDiskSpaceCheck) WarnWhenUsedSpaceIsAbovePercentage(percentage int) *UsedDiskSpaceCheck {
	dsc.warningThreshold = percentage
	return dsc
}

// FailWhenUsedSpaceIsAbovePercentage sets the failure threshold percentage.
//
// Parameters:
//
//	percentage: The percentage (0-100) at which to fail the check
//
// Returns:
//
//	*UsedDiskSpaceCheck: Self for method chaining
func (dsc *UsedDiskSpaceCheck) FailWhenUsedSpaceIsAbovePercentage(percentage int) *UsedDiskSpaceCheck {
	dsc.errorThreshold = percentage
	return dsc
}

// Run performs the disk space health check.
//
// Returns:
//
//	interfaces.ResultInterface: The health check result
func (dsc *UsedDiskSpaceCheck) Run() interfaces.ResultInterface {
	diskSpaceUsedPercentage := dsc.getDiskUsagePercentage()

	result := checks.NewResult().
		SetMeta(map[string]interface{}{"disk_space_used_percentage": diskSpaceUsedPercentage}).
		SetShortSummary(fmt.Sprintf("%d%%", diskSpaceUsedPercentage))

	if diskSpaceUsedPercentage > dsc.errorThreshold {
		return result.
			SetStatus(enums.StatusFailed).
			SetNotificationMessage(fmt.Sprintf("The disk is almost full (%d%% used).", diskSpaceUsedPercentage))
	}

	if diskSpaceUsedPercentage > dsc.warningThreshold {
		return result.
			SetStatus(enums.StatusWarning).
			SetNotificationMessage(fmt.Sprintf("The disk is almost full (%d%% used).", diskSpaceUsedPercentage))
	}

	return result.SetStatus(enums.StatusOK)
}

// getDiskUsagePercentage gets the disk usage percentage using df command.
// This mirrors the Laravel implementation exactly.
func (dsc *UsedDiskSpaceCheck) getDiskUsagePercentage() int {
	// Build the df command - use current directory if no filesystem specified
	path := "."
	if dsc.filesystemName != nil {
		path = *dsc.filesystemName
	}

	// Execute df -P command (POSIX format for consistent parsing)
	cmd := exec.Command("df", "-P", path)
	output, err := cmd.Output()
	if err != nil {
		// If df command fails, return 0 to avoid false alarms
		return 0
	}

	// Parse the percentage from df output using regex
	// df output format: "Filesystem 1024-blocks Used Available Capacity Mounted"
	// We need to extract the percentage (e.g., "85%")
	re := regexp.MustCompile(`(\d+)%`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		// If we can't parse the percentage, return 0
		return 0
	}

	// Convert percentage string to int
	percentage, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return percentage
}

// Compile-time interface compliance check
var _ interfaces.CheckInterface = (*UsedDiskSpaceCheck)(nil)
