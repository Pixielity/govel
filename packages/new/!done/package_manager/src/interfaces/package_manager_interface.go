package interfaces

import (
	"context"
	"govel/package_manager/models"
)

// PackageManagerInterface defines the main package manager operations
type PackageManagerInterface interface {
	// Package operations
	Install(ctx context.Context, packageName string, options InstallOptions) (*models.InstallResult, error)
	Uninstall(ctx context.Context, packageName string) error
	Update(ctx context.Context, packageName string) (*models.InstallResult, error)
	UpdateAll(ctx context.Context) ([]*models.InstallResult, error)

	// Package state management
	Activate(ctx context.Context, packageName string) error
	Deactivate(ctx context.Context, packageName string) error
	IsActive(packageName string) bool
	IsInstalled(packageName string) bool

	// Package information
	List(ctx context.Context, options ListOptions) ([]*models.Package, error)
	Search(ctx context.Context, query string) ([]*models.Package, error)
	Info(ctx context.Context, packageName string) (*models.Package, error)
	Status(ctx context.Context) (*models.PackageState, error)

	// Dependency management
	ResolveDependencies(ctx context.Context, packageName string) (*models.DependencyGraph, error)
	CheckDependencies(ctx context.Context) ([]DependencyIssue, error)

	// Script and hook execution
	RunScript(ctx context.Context, packageName string, scriptName string) (*models.CommandResult, error)
	RunHook(ctx context.Context, packageName string, hookName string) (*models.CommandResult, error)

	// Registry operations
	Refresh(ctx context.Context) error
	Validate(ctx context.Context, packagePath string) error
}

// RegistryInterface defines package registry operations
type RegistryInterface interface {
	// Registry management
	Scan(ctx context.Context, rootPath string) ([]*models.Package, error)
	Load(ctx context.Context, packagePath string) (*models.Package, error)
	Save(ctx context.Context, pkg *models.Package) error
	Remove(ctx context.Context, packageName string) error

	// Package discovery
	FindPackage(packageName string) (*models.Package, error)
	FindPackages(pattern string) ([]*models.Package, error)
	GetInstalledPackages() ([]*models.Package, error)
	GetActivePackages() ([]*models.Package, error)
}

// ParserInterface defines module.json parsing operations
type ParserInterface interface {
	ParseFile(filePath string) (*models.Package, error)
	ParseBytes(data []byte) (*models.Package, error)
	ValidatePackage(pkg *models.Package) error
	WriteFile(pkg *models.Package, filePath string) error
}

// ExecutorInterface defines command and script execution
type ExecutorInterface interface {
	ExecuteCommand(ctx context.Context, command string, workDir string) (*models.CommandResult, error)
	ExecuteScript(ctx context.Context, script string, workDir string) (*models.CommandResult, error)
	ExecuteHooks(ctx context.Context, hooks []string, workDir string) ([]*models.CommandResult, error)
}

// DependencyResolverInterface defines dependency resolution
type DependencyResolverInterface interface {
	Resolve(ctx context.Context, packageName string) (*models.DependencyGraph, error)
	GetInstallOrder(ctx context.Context, packages []string) ([]string, error)
	CheckCircularDependencies(ctx context.Context, packages []string) error
	ValidateConstraints(ctx context.Context, dependencies map[string]string) error
}

// StateManagerInterface defines package state management
type StateManagerInterface interface {
	LoadState() (*models.PackageState, error)
	SaveState(state *models.PackageState) error
	UpdatePackageState(packageName string, isActive bool) error
	GetPackageState(packageName string) (*models.Package, error)
	AddPackageToState(pkg *models.Package) error
	RemovePackageFromState(packageName string) error
	CreateLockFile(packages []*models.Package) error
	ReadLockFile() ([]*models.Package, error)
}

// InstallOptions defines options for package installation
type InstallOptions struct {
	Force       bool
	SkipHooks   bool
	Development bool
	SkipDeps    bool
	Verbose     bool
}

// ListOptions defines options for listing packages
type ListOptions struct {
	Category string
	Status   string // "active", "inactive", "all"
	Format   string // "table", "json", "simple"
	Search   string
}

// DependencyIssue represents a dependency problem
type DependencyIssue struct {
	Package     string `json:"package"`
	Issue       string `json:"issue"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}
