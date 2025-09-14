package commands

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"time"
)

// StatusCommand handles displaying package manager status
type StatusCommand struct {
	pm interfaces.PackageManagerInterface
}

// NewStatusCommand creates a new status command
func NewStatusCommand(pm interfaces.PackageManagerInterface) *StatusCommand {
	return &StatusCommand{pm: pm}
}

// Execute runs the status command
func (cmd *StatusCommand) Execute(ctx context.Context, args []string) error {
	fmt.Printf("📊 GoVel Package Manager Status\n")
	fmt.Printf("═══════════════════════════════\n\n")

	// Get overall status
	status, err := cmd.pm.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Display basic statistics
	fmt.Printf("📦 Package Statistics:\n")
	fmt.Printf("  Total packages:     %d\n", status.TotalCount)
	fmt.Printf("  Active packages:    %d\n", status.ActiveCount)
	fmt.Printf("  Inactive packages:  %d\n", status.TotalCount-status.ActiveCount)

	fmt.Printf("\n🔧 System Information:\n")
	fmt.Printf("  State version:      %s\n", status.Version)
	fmt.Printf("  Last updated:       %s (%s ago)\n",
		status.UpdatedAt.Format("2006-01-02 15:04:05"),
		formatDuration(time.Since(status.UpdatedAt)))

	// Show package breakdown by category
	if err := cmd.showCategoryBreakdown(ctx); err != nil {
		fmt.Printf("  Warning: Could not load category breakdown: %v\n", err)
	}

	// Show recent activities
	fmt.Printf("\n📋 Package Health:\n")
	if status.TotalCount == 0 {
		fmt.Printf("  ⚠️  No packages discovered\n")
	} else {
		activePercentage := float64(status.ActiveCount) / float64(status.TotalCount) * 100
		if activePercentage >= 80 {
			fmt.Printf("  ✅ System health: Excellent (%.1f%% active)\n", activePercentage)
		} else if activePercentage >= 60 {
			fmt.Printf("  🟡 System health: Good (%.1f%% active)\n", activePercentage)
		} else {
			fmt.Printf("  🔴 System health: Needs attention (%.1f%% active)\n", activePercentage)
		}
	}

	// Show recommendations
	cmd.showRecommendations(status)

	return nil
}

func (cmd *StatusCommand) showCategoryBreakdown(ctx context.Context) error {
	// Get all packages to analyze categories
	allPackages, err := cmd.pm.List(ctx, interfaces.ListOptions{Status: "all"})
	if err != nil {
		return err
	}

	categoryCount := make(map[string]int)
	categoryActive := make(map[string]int)

	for _, pkg := range allPackages {
		category := pkg.GovelConfig.Category
		if category == "" {
			category = "uncategorized"
		}
		categoryCount[category]++
		if pkg.IsActive {
			categoryActive[category]++
		}
	}

	if len(categoryCount) > 0 {
		fmt.Printf("\n📂 Packages by Category:\n")
		for category, total := range categoryCount {
			active := categoryActive[category]
			fmt.Printf("  %-15s %d total, %d active\n", category+":", total, active)
		}
	}

	return nil
}

func (cmd *StatusCommand) showRecommendations(status *models.PackageState) {
	fmt.Printf("\n💡 Recommendations:\n")

	if status.TotalCount == 0 {
		fmt.Printf("  • Run 'govel-pkg refresh' to scan for packages\n")
		return
	}

	inactiveCount := status.TotalCount - status.ActiveCount
	if inactiveCount > 0 {
		fmt.Printf("  • Consider activating %d inactive packages if needed\n", inactiveCount)
		fmt.Printf("  • Run 'govel-pkg list --status inactive' to see inactive packages\n")
	}

	if time.Since(status.UpdatedAt) > 24*time.Hour {
		fmt.Printf("  • State hasn't been updated in over 24 hours\n")
		fmt.Printf("  • Run 'govel-pkg refresh' to update package registry\n")
	}

	if status.ActiveCount > 0 {
		fmt.Printf("  • Run 'govel-pkg update --all' to update active packages\n")
	}

	fmt.Printf("  • Use 'govel-pkg search <query>' to find specific packages\n")
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	} else if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	} else {
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
}
