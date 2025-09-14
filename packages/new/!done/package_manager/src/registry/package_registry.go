package registry

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"govel/package_manager/utils"
	"path/filepath"
	"strings"
	"time"
)

// PackageRegistry implements RegistryInterface for managing local packages
type PackageRegistry struct {
	parser       interfaces.ParserInterface
	stateManager interfaces.StateManagerInterface
	cache        map[string]*models.Package
	lastScan     time.Time
}

// NewPackageRegistry creates a new package registry instance
func NewPackageRegistry(parser interfaces.ParserInterface, stateManager interfaces.StateManagerInterface) interfaces.RegistryInterface {
	return &PackageRegistry{
		parser:       parser,
		stateManager: stateManager,
		cache:        make(map[string]*models.Package),
		lastScan:     time.Time{},
	}
}

// Scan scans the root path for packages and returns them
func (pr *PackageRegistry) Scan(ctx context.Context, rootPath string) ([]*models.Package, error) {
	var packages []*models.Package

	// Find all module.json files
	moduleFiles, err := utils.FindFiles(rootPath, "module.json")
	if err != nil {
		return nil, fmt.Errorf("failed to find module.json files: %w", err)
	}

	// Parse each module.json file
	for _, moduleFile := range moduleFiles {
		// Skip if context is canceled
		if ctx.Err() != nil {
			return packages, ctx.Err()
		}

		pkg, err := pr.parser.ParseFile(moduleFile)
		if err != nil {
			// Log error but continue with other packages
			fmt.Printf("Warning: failed to parse %s: %v\n", moduleFile, err)
			continue
		}

		// Validate package
		if err := pr.parser.ValidatePackage(pkg); err != nil {
			fmt.Printf("Warning: invalid package %s: %v\n", moduleFile, err)
			continue
		}

		// Set package path
		pkg.Path = filepath.Dir(moduleFile)

		// Load state information if available
		if statePackage, err := pr.stateManager.GetPackageState(pkg.Name); err == nil {
			pkg.IsActive = statePackage.IsActive
			pkg.IsInstalled = statePackage.IsInstalled
			pkg.InstalledAt = statePackage.InstalledAt
			pkg.UpdatedAt = statePackage.UpdatedAt
		}

		packages = append(packages, pkg)
		pr.cache[pkg.Name] = pkg
	}

	pr.lastScan = time.Now()

	// Update state with discovered packages
	if err := pr.updateStateWithPackages(packages); err != nil {
		fmt.Printf("Warning: failed to update state with discovered packages: %v\n", err)
	}

	return packages, nil
}

// Load loads a single package from the specified path
func (pr *PackageRegistry) Load(ctx context.Context, packagePath string) (*models.Package, error) {
	moduleFile := filepath.Join(packagePath, "module.json")

	if !utils.FileExists(moduleFile) {
		return nil, fmt.Errorf("module.json not found in %s", packagePath)
	}

	pkg, err := pr.parser.ParseFile(moduleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse module.json: %w", err)
	}

	// Validate package
	if err := pr.parser.ValidatePackage(pkg); err != nil {
		return nil, fmt.Errorf("invalid package: %w", err)
	}

	// Set package path
	pkg.Path = packagePath

	// Load state information if available
	if statePackage, err := pr.stateManager.GetPackageState(pkg.Name); err == nil {
		pkg.IsActive = statePackage.IsActive
		pkg.IsInstalled = statePackage.IsInstalled
		pkg.InstalledAt = statePackage.InstalledAt
		pkg.UpdatedAt = statePackage.UpdatedAt
	}

	// Update cache
	pr.cache[pkg.Name] = pkg

	return pkg, nil
}

// Save saves a package's state
func (pr *PackageRegistry) Save(ctx context.Context, pkg *models.Package) error {
	if pkg == nil {
		return fmt.Errorf("package cannot be nil")
	}

	// Update cache
	pr.cache[pkg.Name] = pkg

	// Add to state manager
	if err := pr.stateManager.AddPackageToState(pkg); err != nil {
		return fmt.Errorf("failed to save package state: %w", err)
	}

	return nil
}

// Remove removes a package from the registry
func (pr *PackageRegistry) Remove(ctx context.Context, packageName string) error {
	// Remove from cache
	delete(pr.cache, packageName)

	// Remove from state
	if err := pr.stateManager.RemovePackageFromState(packageName); err != nil {
		return fmt.Errorf("failed to remove package from state: %w", err)
	}

	return nil
}

// FindPackage finds a package by name
func (pr *PackageRegistry) FindPackage(packageName string) (*models.Package, error) {
	// Check cache first
	if pkg, exists := pr.cache[packageName]; exists {
		return pkg, nil
	}

	// Try to find in state
	if statePackage, err := pr.stateManager.GetPackageState(packageName); err == nil {
		pr.cache[packageName] = statePackage
		return statePackage, nil
	}

	// Package not found
	return nil, fmt.Errorf("package '%s' not found", packageName)
}

// FindPackages finds packages matching a pattern
func (pr *PackageRegistry) FindPackages(pattern string) ([]*models.Package, error) {
	var matchingPackages []*models.Package

	// Search in cache
	for name, pkg := range pr.cache {
		if pr.matchesPattern(name, pattern) || pr.matchesPattern(pkg.Description, pattern) {
			matchingPackages = append(matchingPackages, pkg)
		}
	}

	// If no matches in cache, try to load from state
	if len(matchingPackages) == 0 {
		state, err := pr.stateManager.LoadState()
		if err == nil {
			for _, pkg := range state.Packages {
				if pr.matchesPattern(pkg.Name, pattern) || pr.matchesPattern(pkg.Description, pattern) {
					pkgCopy := pkg
					matchingPackages = append(matchingPackages, &pkgCopy)
					pr.cache[pkg.Name] = &pkgCopy
				}
			}
		}
	}

	return matchingPackages, nil
}

// GetInstalledPackages returns all installed packages
func (pr *PackageRegistry) GetInstalledPackages() ([]*models.Package, error) {
	var installedPackages []*models.Package

	// Check cache first
	for _, pkg := range pr.cache {
		if pkg.IsInstalled {
			installedPackages = append(installedPackages, pkg)
		}
	}

	// If cache is empty, load from state
	if len(installedPackages) == 0 {
		state, err := pr.stateManager.LoadState()
		if err != nil {
			return nil, fmt.Errorf("failed to load state: %w", err)
		}

		for _, pkg := range state.Packages {
			if pkg.IsInstalled {
				pkgCopy := pkg
				installedPackages = append(installedPackages, &pkgCopy)
				pr.cache[pkg.Name] = &pkgCopy
			}
		}
	}

	return installedPackages, nil
}

// GetActivePackages returns all active packages
func (pr *PackageRegistry) GetActivePackages() ([]*models.Package, error) {
	var activePackages []*models.Package

	// Check cache first
	for _, pkg := range pr.cache {
		if pkg.IsActive {
			activePackages = append(activePackages, pkg)
		}
	}

	// If cache is empty, load from state
	if len(activePackages) == 0 {
		state, err := pr.stateManager.LoadState()
		if err != nil {
			return nil, fmt.Errorf("failed to load state: %w", err)
		}

		for _, pkg := range state.Packages {
			if pkg.IsActive {
				pkgCopy := pkg
				activePackages = append(activePackages, &pkgCopy)
				pr.cache[pkg.Name] = &pkgCopy
			}
		}
	}

	return activePackages, nil
}

// RefreshCache refreshes the package cache by re-scanning
func (pr *PackageRegistry) RefreshCache(ctx context.Context, rootPath string) error {
	// Clear cache
	pr.cache = make(map[string]*models.Package)

	// Re-scan packages
	_, err := pr.Scan(ctx, rootPath)
	return err
}

// GetCacheStats returns cache statistics
func (pr *PackageRegistry) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cached_packages": len(pr.cache),
		"last_scan":       pr.lastScan,
		"cache_age":       time.Since(pr.lastScan),
	}
}

// ValidatePackages validates all packages in the registry
func (pr *PackageRegistry) ValidatePackages() []error {
	var errors []error

	for name, pkg := range pr.cache {
		if err := pr.parser.ValidatePackage(pkg); err != nil {
			errors = append(errors, fmt.Errorf("package '%s': %w", name, err))
		}
	}

	return errors
}

// GetPackagesByCategory returns packages filtered by category
func (pr *PackageRegistry) GetPackagesByCategory(category string) ([]*models.Package, error) {
	var packages []*models.Package

	for _, pkg := range pr.cache {
		if pkg.GovelConfig.Category == category {
			packages = append(packages, pkg)
		}
	}

	// If cache is empty, load from state
	if len(packages) == 0 {
		state, err := pr.stateManager.LoadState()
		if err != nil {
			return nil, fmt.Errorf("failed to load state: %w", err)
		}

		for _, pkg := range state.Packages {
			if pkg.GovelConfig.Category == category {
				pkgCopy := pkg
				packages = append(packages, &pkgCopy)
				pr.cache[pkg.Name] = &pkgCopy
			}
		}
	}

	return packages, nil
}

// GetPackagesByType returns packages filtered by type
func (pr *PackageRegistry) GetPackagesByType(packageType string) ([]*models.Package, error) {
	var packages []*models.Package

	for _, pkg := range pr.cache {
		if pkg.GovelConfig.Type == packageType {
			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// Private helper methods

func (pr *PackageRegistry) matchesPattern(text, pattern string) bool {
	// Simple pattern matching - can be enhanced with regex or glob patterns
	lowerText := strings.ToLower(text)
	lowerPattern := strings.ToLower(pattern)

	return strings.Contains(lowerText, lowerPattern)
}

func (pr *PackageRegistry) updateStateWithPackages(packages []*models.Package) error {
	state, err := pr.stateManager.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Update state with discovered packages
	for _, pkg := range packages {
		if existingPkg, exists := state.Packages[pkg.Name]; exists {
			// Preserve state information
			pkg.IsActive = existingPkg.IsActive
			pkg.IsInstalled = existingPkg.IsInstalled
			pkg.InstalledAt = existingPkg.InstalledAt
			pkg.UpdatedAt = existingPkg.UpdatedAt
		}

		state.Packages[pkg.Name] = *pkg
	}

	return pr.stateManager.SaveState(state)
}

// GetPackagesWithProviders returns packages that provide specific services
func (pr *PackageRegistry) GetPackagesWithProviders(providerName string) ([]*models.Package, error) {
	var packages []*models.Package

	for _, pkg := range pr.cache {
		for _, provider := range pkg.GovelConfig.Providers {
			if provider == providerName {
				packages = append(packages, pkg)
				break
			}
		}
	}

	return packages, nil
}

// GetPackagesWithFacades returns packages that provide specific facades
func (pr *PackageRegistry) GetPackagesWithFacades(facadeName string) ([]*models.Package, error) {
	var packages []*models.Package

	for _, pkg := range pr.cache {
		for _, facade := range pkg.GovelConfig.Facades {
			if facade == facadeName {
				packages = append(packages, pkg)
				break
			}
		}
	}

	return packages, nil
}
