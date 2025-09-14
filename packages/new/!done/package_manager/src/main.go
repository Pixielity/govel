package main

import (
	"context"
	"fmt"
	"govel/package_manager/commands"
	"govel/package_manager/core"
	"govel/package_manager/parsers"
	"govel/package_manager/registry"
	"govel/package_manager/services"
	"govel/package_manager/utils"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppName    = "govel-pkg"
	AppVersion = "1.0.0"
)

func main() {
	// Initialize the application
	app := &CLI{
		Name:    AppName,
		Version: AppVersion,
		Usage:   "GoVel Package Manager - Manage GoVel framework packages",
	}

	// Set up dependencies
	if err := setupDependencies(app); err != nil {
		fmt.Printf("Error initializing package manager: %v\n", err)
		os.Exit(1)
	}

	// Register commands
	registerCommands(app)

	// Run the application
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// CLI represents the command-line interface
type CLI struct {
	Name     string
	Version  string
	Usage    string
	Commands []Command
	pm       *core.PackageManager
}

// Command represents a CLI command
type Command struct {
	Name        string
	Usage       string
	Description string
	Action      func(ctx context.Context, args []string) error
	Flags       []Flag
}

// Flag represents a command-line flag
type Flag struct {
	Name     string
	Usage    string
	Value    interface{}
	Required bool
}

// Run runs the CLI application with the given arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 2 {
		c.showHelp()
		return nil
	}

	commandName := args[1]

	// Handle global commands
	switch commandName {
	case "help", "-h", "--help":
		c.showHelp()
		return nil
	case "version", "-v", "--version":
		fmt.Printf("%s version %s\n", c.Name, c.Version)
		return nil
	}

	// Find and execute the command
	for _, cmd := range c.Commands {
		if cmd.Name == commandName {
			ctx := context.Background()
			return cmd.Action(ctx, args[2:])
		}
	}

	fmt.Printf("Unknown command: %s\n\n", commandName)
	c.showHelp()
	return fmt.Errorf("unknown command: %s", commandName)
}

// showHelp displays the help information
func (c *CLI) showHelp() {
	fmt.Printf("%s - %s\n\n", c.Name, c.Usage)
	fmt.Printf("VERSION:\n  %s\n\n", c.Version)
	fmt.Printf("USAGE:\n  %s [COMMAND] [OPTIONS]\n\n", c.Name)
	fmt.Printf("COMMANDS:\n")

	for _, cmd := range c.Commands {
		fmt.Printf("  %-12s %s\n", cmd.Name, cmd.Description)
	}

	fmt.Printf("\nGLOBAL OPTIONS:\n")
	fmt.Printf("  %-12s %s\n", "--help, -h", "Show help")
	fmt.Printf("  %-12s %s\n", "--version, -v", "Show version")
	fmt.Printf("\nUse '%s [COMMAND] --help' for more information about a command.\n", c.Name)
}

// setupDependencies initializes all the required dependencies
func setupDependencies(app *CLI) error {
	// Get the current working directory or GoVel root
	rootPath, err := findGovelRoot()
	if err != nil {
		return fmt.Errorf("failed to find GoVel root: %w", err)
	}

	// Initialize components
	parser := parsers.NewModuleParser()
	executor := services.NewCommandExecutor()
	stateManager := services.NewStateManager(filepath.Join(rootPath, ".govel"))
	packageRegistry := registry.NewPackageRegistry(parser, stateManager)
	dependencyResolver := services.NewDependencyResolver(packageRegistry)

	// Create the package manager
	packageManager := core.NewPackageManager(
		rootPath,
		packageRegistry,
		parser,
		executor,
		dependencyResolver,
		stateManager,
	)

	app.pm = packageManager.(*core.PackageManager)

	return nil
}

// registerCommands registers all available CLI commands
func registerCommands(app *CLI) {
	app.Commands = []Command{
		{
			Name:        "install",
			Usage:       "install [PACKAGE_NAME] [OPTIONS]",
			Description: "Install a package",
			Action:      commands.NewInstallCommand(app.pm).Execute,
			Flags: []Flag{
				{Name: "--force", Usage: "Force installation even if already installed"},
				{Name: "--skip-hooks", Usage: "Skip pre/post install hooks"},
				{Name: "--skip-deps", Usage: "Skip dependency installation"},
				{Name: "--verbose", Usage: "Verbose output"},
			},
		},
		{
			Name:        "uninstall",
			Usage:       "uninstall [PACKAGE_NAME]",
			Description: "Uninstall a package",
			Action:      commands.NewUninstallCommand(app.pm).Execute,
		},
		{
			Name:        "list",
			Usage:       "list [OPTIONS]",
			Description: "List installed packages",
			Action:      commands.NewListCommand(app.pm).Execute,
			Flags: []Flag{
				{Name: "--status", Usage: "Filter by status (active, inactive, all)"},
				{Name: "--category", Usage: "Filter by category"},
				{Name: "--format", Usage: "Output format (table, json, simple)"},
				{Name: "--search", Usage: "Search packages"},
			},
		},
		{
			Name:        "activate",
			Usage:       "activate [PACKAGE_NAME]",
			Description: "Activate an installed package",
			Action:      commands.NewActivateCommand(app.pm).Execute,
		},
		{
			Name:        "deactivate",
			Usage:       "deactivate [PACKAGE_NAME]",
			Description: "Deactivate an active package",
			Action:      commands.NewDeactivateCommand(app.pm).Execute,
		},
		{
			Name:        "status",
			Usage:       "status",
			Description: "Show package manager status",
			Action:      commands.NewStatusCommand(app.pm).Execute,
		},
		{
			Name:        "update",
			Usage:       "update [PACKAGE_NAME]",
			Description: "Update a package or all packages",
			Action:      commands.NewUpdateCommand(app.pm).Execute,
			Flags: []Flag{
				{Name: "--all", Usage: "Update all packages"},
			},
		},
		{
			Name:        "search",
			Usage:       "search [QUERY]",
			Description: "Search for packages",
			Action:      commands.NewSearchCommand(app.pm).Execute,
		},
		{
			Name:        "info",
			Usage:       "info [PACKAGE_NAME]",
			Description: "Show detailed package information",
			Action:      commands.NewInfoCommand(app.pm).Execute,
		},
		{
			Name:        "run",
			Usage:       "run [PACKAGE_NAME] [SCRIPT_NAME]",
			Description: "Run a package script",
			Action:      commands.NewRunCommand(app.pm).Execute,
		},
		{
			Name:        "validate",
			Usage:       "validate [PACKAGE_PATH]",
			Description: "Validate a package module.json",
			Action:      commands.NewValidateCommand(app.pm).Execute,
		},
		{
			Name:        "refresh",
			Usage:       "refresh",
			Description: "Refresh the package registry",
			Action:      commands.NewRefreshCommand(app.pm).Execute,
		},
	}
}

// findGovelRoot finds the root directory of the GoVel project
func findGovelRoot() (string, error) {
	// Start from current directory and walk up
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Look for directory that contains a packages subdirectory
	dir := currentDir
	for {
		// Check if this directory contains a packages directory
		packagesDir := filepath.Join(dir, "packages")
		if utils.IsDirectory(packagesDir) {
			// Verify this packages directory contains GoVel packages by checking for module.json files
			moduleFiles, err := utils.FindFiles(packagesDir, "module.json")
			if err == nil && len(moduleFiles) > 0 {
				return dir, nil
			}
		}

		// Check for go.mod with govel module references
		goMod := filepath.Join(dir, "go.mod")
		if utils.IsFile(goMod) {
			content, err := os.ReadFile(goMod)
			if err == nil {
				if utils.ContainsGovelModule(string(content)) {
					// Also check if packages directory exists
					packagesDir := filepath.Join(dir, "packages")
					if utils.IsDirectory(packagesDir) {
						return dir, nil
					}
				}
			}
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root directory
			break
		}
		dir = parent
	}

	// Fallback: if we're inside a packages directory structure, try to find the root
	if strings.Contains(currentDir, "/packages/") {
		// Extract the part before /packages/
		parts := strings.Split(currentDir, "/packages/")
		if len(parts) > 0 && parts[0] != "" {
			potentialRoot := parts[0]
			packagesDir := filepath.Join(potentialRoot, "packages")
			if utils.IsDirectory(packagesDir) {
				return potentialRoot, nil
			}
		}
	}

	// Final fallback to current directory
	return currentDir, nil
}
