package core

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"path/filepath"
	"strings"
	"time"
)

// PackageManager implements the main package manager functionality
type PackageManager struct {
	rootPath           string
	registry           interfaces.RegistryInterface
	parser             interfaces.ParserInterface
	executor           interfaces.ExecutorInterface
	dependencyResolver interfaces.DependencyResolverInterface
	stateManager       interfaces.StateManagerInterface
}

// NewPackageManager creates a new package manager instance
func NewPackageManager(
	rootPath string,
	registry interfaces.RegistryInterface,
	parser interfaces.ParserInterface,
	executor interfaces.ExecutorInterface,
	dependencyResolver interfaces.DependencyResolverInterface,
	stateManager interfaces.StateManagerInterface,
) interfaces.PackageManagerInterface {
	return &PackageManager{
		rootPath:           rootPath,
		registry:           registry,
		parser:             parser,
		executor:           executor,
		dependencyResolver: dependencyResolver,
		stateManager:       stateManager,
	}
}

// Install installs a package with the given options
func (pm *PackageManager) Install(ctx context.Context, packageName string, options interfaces.InstallOptions) (*models.InstallResult, error) {
	startTime := time.Now()

	result := &models.InstallResult{
		Success:  false,
		HooksRun: []string{},
		Duration: 0,
	}

	// Find the package
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		result.Message = fmt.Sprintf("Package '%s' not found: %v", packageName, err)
		result.Duration = time.Since(startTime)
		return result, err
	}
	result.Package = pkg

	// Check if already installed (unless force is specified)
	if pm.IsInstalled(packageName) && !options.Force {
		result.Success = true
		result.Message = fmt.Sprintf("Package '%s' is already installed", packageName)
		result.Duration = time.Since(startTime)
		return result, nil
	}

	// Install dependencies first (unless skipped)
	if !options.SkipDeps {
		if err := pm.installDependencies(ctx, pkg, options); err != nil {
			result.Message = fmt.Sprintf("Failed to install dependencies: %v", err)
			result.Duration = time.Since(startTime)
			return result, err
		}
	}

	// Execute pre-install hooks
	if !options.SkipHooks && len(pkg.Hooks.PreInstall) > 0 {
		if err := pm.executeHooks(ctx, pkg, pkg.Hooks.PreInstall, "pre-install"); err != nil {
			result.Message = fmt.Sprintf("Pre-install hooks failed: %v", err)
			result.Duration = time.Since(startTime)
			return result, err
		}
		result.HooksRun = append(result.HooksRun, "pre-install")
	}

	// Mark as installed and active
	pkg.IsInstalled = true
	pkg.IsActive = true
	pkg.InstalledAt = time.Now()
	pkg.UpdatedAt = time.Now()

	// Save package state
	if err := pm.registry.Save(ctx, pkg); err != nil {
		result.Message = fmt.Sprintf("Failed to save package state: %v", err)
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Update global state
	if err := pm.stateManager.UpdatePackageState(packageName, true); err != nil {
		result.Message = fmt.Sprintf("Failed to update package state: %v", err)
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Execute post-install hooks
	if !options.SkipHooks && len(pkg.Hooks.PostInstall) > 0 {
		if err := pm.executeHooks(ctx, pkg, pkg.Hooks.PostInstall, "post-install"); err != nil {
			result.Message = fmt.Sprintf("Post-install hooks failed: %v", err)
			result.Duration = time.Since(startTime)
			return result, err
		}
		result.HooksRun = append(result.HooksRun, "post-install")
	}

	result.Success = true
	result.Message = fmt.Sprintf("Package '%s' installed successfully", packageName)
	result.Duration = time.Since(startTime)

	return result, nil
}

// Uninstall removes a package from the system
func (pm *PackageManager) Uninstall(ctx context.Context, packageName string) error {
	// Find the package
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	// Check if package is installed
	if !pkg.IsInstalled {
		return fmt.Errorf("package '%s' is not installed", packageName)
	}

	// Check for dependents
	dependents, err := pm.findDependents(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to check dependents: %w", err)
	}

	if len(dependents) > 0 {
		return fmt.Errorf("cannot uninstall '%s': still required by %s", packageName, strings.Join(dependents, ", "))
	}

	// Mark as uninstalled and inactive
	pkg.IsInstalled = false
	pkg.IsActive = false
	pkg.UpdatedAt = time.Now()

	// Save package state
	if err := pm.registry.Save(ctx, pkg); err != nil {
		return fmt.Errorf("failed to save package state: %w", err)
	}

	// Update global state
	if err := pm.stateManager.UpdatePackageState(packageName, false); err != nil {
		return fmt.Errorf("failed to update package state: %w", err)
	}

	return nil
}

// Update updates a package to the latest version
func (pm *PackageManager) Update(ctx context.Context, packageName string) (*models.InstallResult, error) {
	// For now, treat update the same as install with force=true
	options := interfaces.InstallOptions{
		Force:   true,
		Verbose: true,
	}

	result, err := pm.Install(ctx, packageName, options)
	if err != nil {
		return result, err
	}

	if result.Success {
		result.Message = fmt.Sprintf("Package '%s' updated successfully", packageName)
	}

	return result, nil
}

// UpdateAll updates all installed packages
func (pm *PackageManager) UpdateAll(ctx context.Context) ([]*models.InstallResult, error) {
	packages, err := pm.registry.GetInstalledPackages()
	if err != nil {
		return nil, fmt.Errorf("failed to get installed packages: %w", err)
	}

	var results []*models.InstallResult
	for _, pkg := range packages {
		result, err := pm.Update(ctx, pkg.Name)
		if err != nil {
			result = &models.InstallResult{
				Package: pkg,
				Success: false,
				Message: fmt.Sprintf("Update failed: %v", err),
			}
		}
		results = append(results, result)
	}

	return results, nil
}

// Activate activates an installed package
func (pm *PackageManager) Activate(ctx context.Context, packageName string) error {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	if !pkg.IsInstalled {
		return fmt.Errorf("package '%s' is not installed", packageName)
	}

	if pkg.IsActive {
		return fmt.Errorf("package '%s' is already active", packageName)
	}

	pkg.IsActive = true
	pkg.UpdatedAt = time.Now()

	if err := pm.registry.Save(ctx, pkg); err != nil {
		return fmt.Errorf("failed to save package state: %w", err)
	}

	if err := pm.stateManager.UpdatePackageState(packageName, true); err != nil {
		return fmt.Errorf("failed to update package state: %w", err)
	}

	return nil
}

// Deactivate deactivates an active package
func (pm *PackageManager) Deactivate(ctx context.Context, packageName string) error {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	if !pkg.IsActive {
		return fmt.Errorf("package '%s' is not active", packageName)
	}

	// Check for active dependents
	dependents, err := pm.findActiveDependents(ctx, packageName)
	if err != nil {
		return fmt.Errorf("failed to check dependents: %w", err)
	}

	if len(dependents) > 0 {
		return fmt.Errorf("cannot deactivate '%s': still required by active packages %s", packageName, strings.Join(dependents, ", "))
	}

	pkg.IsActive = false
	pkg.UpdatedAt = time.Now()

	if err := pm.registry.Save(ctx, pkg); err != nil {
		return fmt.Errorf("failed to save package state: %w", err)
	}

	if err := pm.stateManager.UpdatePackageState(packageName, false); err != nil {
		return fmt.Errorf("failed to update package state: %w", err)
	}

	return nil
}

// IsActive checks if a package is active
func (pm *PackageManager) IsActive(packageName string) bool {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return false
	}
	return pkg.IsActive
}

// IsInstalled checks if a package is installed
func (pm *PackageManager) IsInstalled(packageName string) bool {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return false
	}
	return pkg.IsInstalled
}

// List returns a list of packages based on the given options
func (pm *PackageManager) List(ctx context.Context, options interfaces.ListOptions) ([]*models.Package, error) {
	var packages []*models.Package
	var err error

	switch options.Status {
	case "active":
		packages, err = pm.registry.GetActivePackages()
	case "inactive":
		allPackages, allErr := pm.registry.GetInstalledPackages()
		if allErr != nil {
			return nil, allErr
		}
		for _, pkg := range allPackages {
			if !pkg.IsActive {
				packages = append(packages, pkg)
			}
		}
	default:
		// "all" should list all discovered packages, not only installed ones
		packages, err = pm.registry.Scan(ctx, pm.rootPath)
	}

	if err != nil {
		return nil, err
	}

	// Apply filters
	if options.Category != "" {
		packages = pm.filterByCategory(packages, options.Category)
	}

	if options.Search != "" {
		packages = pm.filterBySearch(packages, options.Search)
	}

	return packages, nil
}

// Search searches for packages by name or description
func (pm *PackageManager) Search(ctx context.Context, query string) ([]*models.Package, error) {
	allPackages, err := pm.registry.Scan(ctx, pm.rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan packages: %w", err)
	}

	return pm.filterBySearch(allPackages, query), nil
}

// Info returns detailed information about a package
func (pm *PackageManager) Info(ctx context.Context, packageName string) (*models.Package, error) {
	return pm.registry.FindPackage(packageName)
}

// Status returns the current status of all packages
func (pm *PackageManager) Status(ctx context.Context) (*models.PackageState, error) {
	return pm.stateManager.LoadState()
}

// ResolveDependencies resolves dependencies for a package
func (pm *PackageManager) ResolveDependencies(ctx context.Context, packageName string) (*models.DependencyGraph, error) {
	return pm.dependencyResolver.Resolve(ctx, packageName)
}

// CheckDependencies checks for dependency issues
func (pm *PackageManager) CheckDependencies(ctx context.Context) ([]interfaces.DependencyIssue, error) {
	// Implementation would check for missing dependencies, circular dependencies, version conflicts, etc.
	return []interfaces.DependencyIssue{}, nil
}

// RunScript runs a script defined in a package
func (pm *PackageManager) RunScript(ctx context.Context, packageName string, scriptName string) (*models.CommandResult, error) {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	script, exists := pkg.Scripts[scriptName]
	if !exists {
		return nil, fmt.Errorf("script '%s' not found in package '%s'", scriptName, packageName)
	}

	return pm.executor.ExecuteScript(ctx, script, pkg.Path)
}

// RunHook runs a specific hook for a package
func (pm *PackageManager) RunHook(ctx context.Context, packageName string, hookName string) (*models.CommandResult, error) {
	pkg, err := pm.registry.FindPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	var hooks []string
	switch hookName {
	case "pre-install":
		hooks = pkg.Hooks.PreInstall
	case "post-install":
		hooks = pkg.Hooks.PostInstall
	case "pre-update":
		hooks = pkg.Hooks.PreUpdate
	case "post-update":
		hooks = pkg.Hooks.PostUpdate
	case "pre-build":
		hooks = pkg.Hooks.PreBuild
	case "post-build":
		hooks = pkg.Hooks.PostBuild
	case "pre-test":
		hooks = pkg.Hooks.PreTest
	case "post-test":
		hooks = pkg.Hooks.PostTest
	case "pre-publish":
		hooks = pkg.Hooks.PrePublish
	case "post-publish":
		hooks = pkg.Hooks.PostPublish
	default:
		return nil, fmt.Errorf("unknown hook '%s'", hookName)
	}

	if len(hooks) == 0 {
		return &models.CommandResult{
			Command:  hookName,
			Success:  true,
			Output:   fmt.Sprintf("No hooks defined for %s", hookName),
			ExitCode: 0,
		}, nil
	}

	results, err := pm.executor.ExecuteHooks(ctx, hooks, pkg.Path)
	if err != nil {
		return nil, err
	}

	// Return the last result (or combine them)
	if len(results) > 0 {
		return results[len(results)-1], nil
	}

	return &models.CommandResult{
		Command: hookName,
		Success: true,
		Output:  "Hooks executed successfully",
	}, nil
}

// Refresh refreshes the package registry
func (pm *PackageManager) Refresh(ctx context.Context) error {
	_, err := pm.registry.Scan(ctx, pm.rootPath)
	return err
}

// Validate validates a package at the given path
func (pm *PackageManager) Validate(ctx context.Context, packagePath string) error {
	moduleFile := filepath.Join(packagePath, "module.json")
	pkg, err := pm.parser.ParseFile(moduleFile)
	if err != nil {
		return fmt.Errorf("failed to parse module.json: %w", err)
	}

	return pm.parser.ValidatePackage(pkg)
}

// Private helper methods

func (pm *PackageManager) installDependencies(ctx context.Context, pkg *models.Package, options interfaces.InstallOptions) error {
	if options.Verbose {
		fmt.Printf("🔍 Installing dependencies for package: %s\n", pkg.Name)
		fmt.Printf("📋 Found %d dependencies\n", len(pkg.Dependencies))
	}
	
	for depName, depVersion := range pkg.Dependencies {
		if options.Verbose {
			fmt.Printf("  📦 Processing dependency: %s@%s\n", depName, depVersion)
		}
		
		// Check if it's a GoVel package dependency (starts with @govel/)
		if strings.HasPrefix(depName, "@govel/") {
			if options.Verbose {
				fmt.Printf("    🔍 GoVel package dependency\n")
			}
			// Handle GoVel package dependencies
			if !pm.IsInstalled(depName) {
				if _, err := pm.Install(ctx, depName, options); err != nil {
					return fmt.Errorf("failed to install GoVel dependency '%s': %w", depName, err)
				}
			}
		} else if strings.Contains(depName, "/") {
			if options.Verbose {
				fmt.Printf("    🌐 External Go module dependency\n")
			}
			// Handle external Go module dependencies (e.g., github.com/google/uuid)
			if err := pm.installExternalDependency(ctx, pkg, depName, depVersion, options); err != nil {
				return fmt.Errorf("failed to install external dependency '%s': %w", depName, err)
			}
		} else {
			if options.Verbose {
				fmt.Printf("    ⚠️  Unknown dependency type: %s\n", depName)
			}
		}
	}
	return nil
}

// installExternalDependency adds an external Go module dependency to the package's go.mod
func (pm *PackageManager) installExternalDependency(ctx context.Context, pkg *models.Package, depName, depVersion string, options interfaces.InstallOptions) error {
	// Convert version constraint (^1.6.0) to Go module version (v1.6.0)
	goVersion := pm.convertVersionConstraint(depVersion)
	
	if options.Verbose {
			fmt.Printf("    📥 Installing %s@%s (converted from %s)\n", depName, goVersion, depVersion)
	}
	
	// Add the dependency using go get
	goModPath := filepath.Join(pkg.Path, "src")
	command := fmt.Sprintf("go get %s@%s", depName, goVersion)
	
	if options.Verbose {
			fmt.Printf("    🔧 Running: %s in %s\n", command, goModPath)
	}
	
	result, err := pm.executor.ExecuteCommand(ctx, command, goModPath)
	if err != nil {
		return fmt.Errorf("failed to execute go get: %w", err)
	}
	
	if !result.Success {
		return fmt.Errorf("go get failed: %s\nOutput: %s", result.Error, result.Output)
	}
	
	if options.Verbose {
			fmt.Printf("    ✅ Successfully installed %s@%s\n", depName, goVersion)
	}
	
	return nil
}

// convertVersionConstraint converts npm-style version constraints to Go module versions
func (pm *PackageManager) convertVersionConstraint(constraint string) string {
	// Remove constraint prefixes (^, ~, >=, etc.)
	version := strings.TrimPrefix(constraint, "^")
	version = strings.TrimPrefix(version, "~")
	version = strings.TrimPrefix(version, ">=")
	version = strings.TrimPrefix(version, "<=")
	version = strings.TrimPrefix(version, ">")
	version = strings.TrimPrefix(version, "<")
	version = strings.TrimSpace(version)
	
	// Add 'v' prefix if not present
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	
	return version
}

func (pm *PackageManager) executeHooks(ctx context.Context, pkg *models.Package, hooks []string, hookType string) error {
	if len(hooks) == 0 {
		return nil
	}

	results, err := pm.executor.ExecuteHooks(ctx, hooks, pkg.Path)
	if err != nil {
		// Provide detailed error information
		if len(results) > 0 {
			lastResult := results[len(results)-1]
			return fmt.Errorf("%s hook failed: %s\nCommand: %s\nOutput: %s\nError: %s\nExit Code: %d",
				hookType, err.Error(), lastResult.Command, lastResult.Output, lastResult.Error, lastResult.ExitCode)
		}
		return fmt.Errorf("%s hook failed: %s", hookType, err.Error())
	}
	return nil
}

func (pm *PackageManager) findDependents(ctx context.Context, packageName string) ([]string, error) {
	var dependents []string
	packages, err := pm.registry.GetInstalledPackages()
	if err != nil {
		return nil, err
	}

	for _, pkg := range packages {
		for depName := range pkg.Dependencies {
			if depName == packageName {
				dependents = append(dependents, pkg.Name)
				break
			}
		}
	}

	return dependents, nil
}

func (pm *PackageManager) findActiveDependents(ctx context.Context, packageName string) ([]string, error) {
	var dependents []string
	packages, err := pm.registry.GetActivePackages()
	if err != nil {
		return nil, err
	}

	for _, pkg := range packages {
		for depName := range pkg.Dependencies {
			if depName == packageName {
				dependents = append(dependents, pkg.Name)
				break
			}
		}
	}

	return dependents, nil
}

func (pm *PackageManager) filterByCategory(packages []*models.Package, category string) []*models.Package {
	var filtered []*models.Package
	for _, pkg := range packages {
		if pkg.GovelConfig.Category == category {
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}

func (pm *PackageManager) filterBySearch(packages []*models.Package, query string) []*models.Package {
	var filtered []*models.Package
	query = strings.ToLower(query)

	for _, pkg := range packages {
		if strings.Contains(strings.ToLower(pkg.Name), query) ||
			strings.Contains(strings.ToLower(pkg.Description), query) {
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}
