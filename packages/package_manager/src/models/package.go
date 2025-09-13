package models

import "time"

// Package represents a GoVel package with its metadata and state
type Package struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	License      string            `json:"license"`
	Keywords     []string          `json:"keywords"`
	Repository   Repository        `json:"repository"`
	Bugs         Bugs              `json:"bugs"`
	Homepage     string            `json:"homepage"`
	Dependencies map[string]string `json:"dependencies"`
	Scripts      map[string]string `json:"scripts"`
	Hooks        Hooks             `json:"hooks"`
	Engines      map[string]string `json:"engines"`
	Files        []string          `json:"files"`
	GovelConfig  GovelConfig       `json:"govel"`
	
	// Package manager state
	Path         string            `json:"path"`
	IsActive     bool              `json:"is_active"`
	IsInstalled  bool              `json:"is_installed"`
	InstalledAt  time.Time         `json:"installed_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// Repository represents the repository information
type Repository struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Directory string `json:"directory"`
}

// Bugs represents bug tracking information
type Bugs struct {
	URL string `json:"url"`
}

// Hooks represents lifecycle hooks for the package
type Hooks struct {
	PreInstall   []string `json:"pre-install"`
	PostInstall  []string `json:"post-install"`
	PreUpdate    []string `json:"pre-update"`
	PostUpdate   []string `json:"post-update"`
	PreBuild     []string `json:"pre-build"`
	PostBuild    []string `json:"post-build"`
	PreTest      []string `json:"pre-test"`
	PostTest     []string `json:"post-test"`
	PrePublish   []string `json:"pre-publish"`
	PostPublish  []string `json:"post-publish"`
}

// GovelConfig represents GoVel-specific configuration
type GovelConfig struct {
	Type         string   `json:"type"`
	Category     string   `json:"category"`
	Providers    []string `json:"providers"`
	Facades      []string `json:"facades,omitempty"`
	Middleware   []string `json:"middleware,omitempty"`
	Commands     []string `json:"commands,omitempty"`
	Provides     []string `json:"provides,omitempty"`
}

// PackageState represents the current state of packages
type PackageState struct {
	Version      string             `json:"version"`
	Packages     map[string]Package `json:"packages"`
	ActiveCount  int                `json:"active_count"`
	TotalCount   int                `json:"total_count"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// DependencyGraph represents the dependency relationships
type DependencyGraph struct {
	Nodes map[string]*DependencyNode `json:"nodes"`
	Edges []DependencyEdge           `json:"edges"`
}

// DependencyNode represents a package in the dependency graph
type DependencyNode struct {
	Package      *Package  `json:"package"`
	Dependencies []string  `json:"dependencies"`
	Dependents   []string  `json:"dependents"`
	Level        int       `json:"level"`
}

// DependencyEdge represents a dependency relationship
type DependencyEdge struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Constraint  string `json:"constraint"`
	Required    bool   `json:"required"`
}

// InstallResult represents the result of a package installation
type InstallResult struct {
	Package     *Package `json:"package"`
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Duration    time.Duration `json:"duration"`
	HooksRun    []string `json:"hooks_run"`
}

// CommandResult represents the result of running a command or script
type CommandResult struct {
	Command     string        `json:"command"`
	Success     bool          `json:"success"`
	Output      string        `json:"output"`
	Error       string        `json:"error"`
	ExitCode    int           `json:"exit_code"`
	Duration    time.Duration `json:"duration"`
}