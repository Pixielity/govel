package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"strings"
)

// ListCommand handles package listing
type ListCommand struct {
	pm interfaces.PackageManagerInterface
}

// NewListCommand creates a new list command
func NewListCommand(pm interfaces.PackageManagerInterface) *ListCommand {
	return &ListCommand{pm: pm}
}

// Execute runs the list command
func (cmd *ListCommand) Execute(ctx context.Context, args []string) error {
	options := cmd.parseOptions(args)

	packages, err := cmd.pm.List(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	if len(packages) == 0 {
		fmt.Printf("No packages found matching the criteria.\n")
		return nil
	}

	switch options.Format {
	case "json":
		return cmd.displayJSON(packages)
	case "simple":
		return cmd.displaySimple(packages)
	default:
		return cmd.displayTable(packages)
	}
}

func (cmd *ListCommand) parseOptions(args []string) interfaces.ListOptions {
	options := interfaces.ListOptions{
		Status:   "all",
		Format:   "table",
		Category: "",
		Search:   "",
	}

	for i, arg := range args {
		switch arg {
		case "--status":
			if i+1 < len(args) {
				options.Status = args[i+1]
			}
		case "--category":
			if i+1 < len(args) {
				options.Category = args[i+1]
			}
		case "--format":
			if i+1 < len(args) {
				options.Format = args[i+1]
			}
		case "--search":
			if i+1 < len(args) {
				options.Search = args[i+1]
			}
		}
	}

	return options
}

func (cmd *ListCommand) displayTable(packages []*models.Package) error {
	// Print header
	fmt.Printf("┌─────────────────────────────────────────┬─────────┬───────────────┬────────┬───────────┐\n")
	fmt.Printf("│ %-39s │ %-7s │ %-13s │ %-6s │ %-9s │\n", "Package Name", "Version", "Category", "Status", "Installed")
	fmt.Printf("├─────────────────────────────────────────┼─────────┼───────────────┼────────┼───────────┤\n")

	for _, pkg := range packages {
		name := pkg.Name
		if len(name) > 39 {
			name = name[:36] + "..."
		}

		version := pkg.Version
		if len(version) > 7 {
			version = version[:7]
		}

		category := pkg.GovelConfig.Category
		if len(category) > 13 {
			category = category[:10] + "..."
		}

		status := "inactive"
		if pkg.IsActive {
			status = "active"
		}

		installed := "No"
		if pkg.IsInstalled {
			installed = "Yes"
		}

		// Color coding
		statusColor := ""
		if pkg.IsActive {
			statusColor = "🟢 "
		} else if pkg.IsInstalled {
			statusColor = "🟡 "
		} else {
			statusColor = "🔴 "
		}

		fmt.Printf("│ %-39s │ %-7s │ %-13s │ %s%-6s │ %-9s │\n",
			name, version, category, statusColor, status, installed)
	}

	fmt.Printf("└─────────────────────────────────────────┴─────────┴───────────────┴────────┴───────────┘\n")
	fmt.Printf("\nTotal: %d packages\n", len(packages))

	// Show summary
	active := 0
	installed := 0
	for _, pkg := range packages {
		if pkg.IsActive {
			active++
		}
		if pkg.IsInstalled {
			installed++
		}
	}

	fmt.Printf("Active: %d, Installed: %d, Available: %d\n", active, installed, len(packages)-installed)

	return nil
}

func (cmd *ListCommand) displaySimple(packages []*models.Package) error {
	for _, pkg := range packages {
		status := ""
		if pkg.IsActive {
			status = " (active)"
		} else if pkg.IsInstalled {
			status = " (installed)"
		}
		fmt.Printf("%s@%s%s\n", pkg.Name, pkg.Version, status)
	}
	return nil
}

func (cmd *ListCommand) displayJSON(packages []*models.Package) error {
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// Helper methods for different command variations

func (cmd *ListCommand) ShowHelp() {
	help := `
Usage: govel-pkg list [OPTIONS]

List installed packages with various filtering and formatting options.

OPTIONS:
  --status STATUS    Filter by status: active, inactive, all (default: all)
  --category CAT     Filter by category (core, web, security, etc.)
  --format FORMAT    Output format: table, json, simple (default: table)
  --search QUERY     Search packages by name or description

EXAMPLES:
  govel-pkg list                           # List all packages in table format
  govel-pkg list --status active          # List only active packages
  govel-pkg list --category web            # List web packages
  govel-pkg list --format json            # Output as JSON
  govel-pkg list --search cookie           # Search for packages containing "cookie"
  govel-pkg list --status active --format simple  # Active packages in simple format

STATUS COLORS:
  🟢 active    - Package is installed and active
  🟡 installed - Package is installed but inactive
  🔴 available - Package is available but not installed
`
	fmt.Print(strings.TrimSpace(help))
}
