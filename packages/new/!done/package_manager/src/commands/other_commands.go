package commands

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"strings"
)

// ActivateCommand handles package activation
type ActivateCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewActivateCommand(pm interfaces.PackageManagerInterface) *ActivateCommand {
	return &ActivateCommand{pm: pm}
}

func (cmd *ActivateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]
	fmt.Printf("Activating package: %s\n", packageName)

	if err := cmd.pm.Activate(ctx, packageName); err != nil {
		return fmt.Errorf("failed to activate package: %w", err)
	}

	fmt.Printf("✅ Package '%s' activated successfully\n", packageName)
	return nil
}

// DeactivateCommand handles package deactivation
type DeactivateCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewDeactivateCommand(pm interfaces.PackageManagerInterface) *DeactivateCommand {
	return &DeactivateCommand{pm: pm}
}

func (cmd *DeactivateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]
	fmt.Printf("Deactivating package: %s\n", packageName)

	if err := cmd.pm.Deactivate(ctx, packageName); err != nil {
		return fmt.Errorf("failed to deactivate package: %w", err)
	}

	fmt.Printf("✅ Package '%s' deactivated successfully\n", packageName)
	return nil
}

// UninstallCommand handles package uninstallation
type UninstallCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewUninstallCommand(pm interfaces.PackageManagerInterface) *UninstallCommand {
	return &UninstallCommand{pm: pm}
}

func (cmd *UninstallCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]
	fmt.Printf("Uninstalling package: %s\n", packageName)

	if err := cmd.pm.Uninstall(ctx, packageName); err != nil {
		return fmt.Errorf("failed to uninstall package: %w", err)
	}

	fmt.Printf("✅ Package '%s' uninstalled successfully\n", packageName)
	return nil
}

// UpdateCommand handles package updates
type UpdateCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewUpdateCommand(pm interfaces.PackageManagerInterface) *UpdateCommand {
	return &UpdateCommand{pm: pm}
}

func (cmd *UpdateCommand) Execute(ctx context.Context, args []string) error {
	updateAll := false
	for _, arg := range args {
		if arg == "--all" {
			updateAll = true
			break
		}
	}

	if updateAll {
		fmt.Printf("Updating all packages...\n")
		results, err := cmd.pm.UpdateAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to update packages: %w", err)
		}

		for _, result := range results {
			if result.Success {
				fmt.Printf("✅ %s\n", result.Message)
			} else {
				fmt.Printf("❌ %s\n", result.Message)
			}
		}
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("package name is required (or use --all)")
	}

	packageName := args[0]
	fmt.Printf("Updating package: %s\n", packageName)

	result, err := cmd.pm.Update(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to update package: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
	} else {
		fmt.Printf("❌ %s\n", result.Message)
	}

	return nil
}

// SearchCommand handles package searching
type SearchCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewSearchCommand(pm interfaces.PackageManagerInterface) *SearchCommand {
	return &SearchCommand{pm: pm}
}

func (cmd *SearchCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("search query is required")
	}

	query := strings.Join(args, " ")
	fmt.Printf("Searching for: %s\n\n", query)

	packages, err := cmd.pm.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(packages) == 0 {
		fmt.Printf("No packages found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d package(s):\n\n", len(packages))
	for _, pkg := range packages {
		status := ""
		if pkg.IsActive {
			status = " ✅ (active)"
		} else if pkg.IsInstalled {
			status = " 📦 (installed)"
		}

		fmt.Printf("📦 %s@%s%s\n", pkg.Name, pkg.Version, status)
		if pkg.Description != "" {
			fmt.Printf("   %s\n", pkg.Description)
		}
		fmt.Printf("   Category: %s\n", pkg.GovelConfig.Category)
		fmt.Printf("\n")
	}

	return nil
}

// InfoCommand handles showing package information
type InfoCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewInfoCommand(pm interfaces.PackageManagerInterface) *InfoCommand {
	return &InfoCommand{pm: pm}
}

func (cmd *InfoCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]
	pkg, err := cmd.pm.Info(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to get package info: %w", err)
	}

	fmt.Printf("📦 Package Information\n")
	fmt.Printf("═══════════════════════\n\n")
	fmt.Printf("Name:         %s\n", pkg.Name)
	fmt.Printf("Version:      %s\n", pkg.Version)
	fmt.Printf("Description:  %s\n", pkg.Description)
	fmt.Printf("Author:       %s\n", pkg.Author)
	fmt.Printf("License:      %s\n", pkg.License)
	fmt.Printf("Category:     %s\n", pkg.GovelConfig.Category)
	fmt.Printf("Type:         %s\n", pkg.GovelConfig.Type)
	fmt.Printf("Homepage:     %s\n", pkg.Homepage)

	if pkg.IsInstalled {
		fmt.Printf("Status:       ✅ Installed")
		if pkg.IsActive {
			fmt.Printf(" & Active")
		}
		fmt.Printf("\n")
		fmt.Printf("Installed:    %s\n", pkg.InstalledAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:      %s\n", pkg.UpdatedAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("Status:       📦 Available\n")
	}

	if len(pkg.Keywords) > 0 {
		fmt.Printf("Keywords:     %s\n", strings.Join(pkg.Keywords, ", "))
	}

	if len(pkg.Dependencies) > 0 {
		fmt.Printf("\n🔗 Dependencies:\n")
		for dep, constraint := range pkg.Dependencies {
			fmt.Printf("  • %s %s\n", dep, constraint)
		}
	}

	if len(pkg.GovelConfig.Providers) > 0 {
		fmt.Printf("\n🔧 Service Providers:\n")
		for _, provider := range pkg.GovelConfig.Providers {
			fmt.Printf("  • %s\n", provider)
		}
	}

	if len(pkg.Scripts) > 0 {
		fmt.Printf("\n📝 Available Scripts:\n")
		for script, command := range pkg.Scripts {
			fmt.Printf("  • %s: %s\n", script, command)
		}
	}

	return nil
}

// RefreshCommand handles refreshing the package registry
type RefreshCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewRefreshCommand(pm interfaces.PackageManagerInterface) *RefreshCommand {
	return &RefreshCommand{pm: pm}
}

func (cmd *RefreshCommand) Execute(ctx context.Context, args []string) error {
	fmt.Printf("Refreshing package registry...\n")

	if err := cmd.pm.Refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh registry: %w", err)
	}

	fmt.Printf("✅ Package registry refreshed successfully\n")
	return nil
}

// ValidateCommand handles package validation
type ValidateCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewValidateCommand(pm interfaces.PackageManagerInterface) *ValidateCommand {
	return &ValidateCommand{pm: pm}
}

func (cmd *ValidateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package path is required")
	}

	packagePath := args[0]
	fmt.Printf("Validating package at: %s\n", packagePath)

	if err := cmd.pm.Validate(ctx, packagePath); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ Package validation successful\n")
	return nil
}

// RunCommand handles running package scripts
type RunCommand struct {
	pm interfaces.PackageManagerInterface
}

func NewRunCommand(pm interfaces.PackageManagerInterface) *RunCommand {
	return &RunCommand{pm: pm}
}

func (cmd *RunCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("package name and script name are required")
	}

	packageName := args[0]
	scriptName := args[1]

	fmt.Printf("Running script '%s' in package '%s'...\n", scriptName, packageName)

	result, err := cmd.pm.RunScript(ctx, packageName, scriptName)
	if err != nil {
		return fmt.Errorf("failed to run script: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ Script executed successfully\n")
		if result.Output != "" {
			fmt.Printf("Output:\n%s\n", result.Output)
		}
	} else {
		fmt.Printf("❌ Script failed with exit code %d\n", result.ExitCode)
		if result.Error != "" {
			fmt.Printf("Error: %s\n", result.Error)
		}
		if result.Output != "" {
			fmt.Printf("Output:\n%s\n", result.Output)
		}
	}

	fmt.Printf("Duration: %v\n", result.Duration)
	return nil
}
