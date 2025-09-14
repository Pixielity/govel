package commands

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"strings"
)

// InstallCommand handles package installation
type InstallCommand struct {
	pm interfaces.PackageManagerInterface
}

// NewInstallCommand creates a new install command
func NewInstallCommand(pm interfaces.PackageManagerInterface) *InstallCommand {
	return &InstallCommand{pm: pm}
}

// Execute runs the install command
func (cmd *InstallCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("package name is required")
	}

	packageName := args[0]
	options := cmd.parseOptions(args[1:])

	fmt.Printf("Installing package: %s\n", packageName)

	if options.Verbose {
		fmt.Printf("Options: force=%v, skip-hooks=%v, skip-deps=%v\n",
			options.Force, options.SkipHooks, options.SkipDeps)
	}

	result, err := cmd.pm.Install(ctx, packageName, options)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
		if len(result.HooksRun) > 0 {
			fmt.Printf("🔧 Hooks executed: %s\n", strings.Join(result.HooksRun, ", "))
		}
		fmt.Printf("⏱️  Completed in %v\n", result.Duration)
	} else {
		fmt.Printf("❌ %s\n", result.Message)
		if options.Verbose && result.Package != nil {
			fmt.Printf("📦 Package: %s v%s\n", result.Package.Name, result.Package.Version)
		}
	}

	return nil
}

func (cmd *InstallCommand) parseOptions(args []string) interfaces.InstallOptions {
	options := interfaces.InstallOptions{}

	for _, arg := range args {
		switch arg {
		case "--force", "-f":
			options.Force = true
		case "--skip-hooks":
			options.SkipHooks = true
		case "--skip-deps":
			options.SkipDeps = true
		case "--verbose", "-v":
			options.Verbose = true
		case "--dev":
			options.Development = true
		}
	}

	return options
}
