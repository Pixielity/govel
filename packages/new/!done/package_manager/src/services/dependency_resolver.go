package services

import (
	"context"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"sort"
	"strings"
)

// DependencyResolver implements DependencyResolverInterface for resolving package dependencies
type DependencyResolver struct {
	registry interfaces.RegistryInterface
}

// NewDependencyResolver creates a new dependency resolver instance
func NewDependencyResolver(registry interfaces.RegistryInterface) interfaces.DependencyResolverInterface {
	return &DependencyResolver{
		registry: registry,
	}
}

// Resolve resolves dependencies for a package and returns a dependency graph
func (dr *DependencyResolver) Resolve(ctx context.Context, packageName string) (*models.DependencyGraph, error) {
	graph := &models.DependencyGraph{
		Nodes: make(map[string]*models.DependencyNode),
		Edges: []models.DependencyEdge{},
	}

	// Find the root package
	rootPkg, err := dr.registry.FindPackage(packageName)
	if err != nil {
		return nil, fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	// Build the dependency graph recursively
	if err := dr.buildDependencyGraph(ctx, rootPkg, graph, 0, make(map[string]bool)); err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Calculate dependency levels
	dr.calculateLevels(graph)

	return graph, nil
}

// GetInstallOrder returns the order in which packages should be installed
func (dr *DependencyResolver) GetInstallOrder(ctx context.Context, packages []string) ([]string, error) {
	// Build a combined dependency graph for all packages
	combinedGraph := &models.DependencyGraph{
		Nodes: make(map[string]*models.DependencyNode),
		Edges: []models.DependencyEdge{},
	}

	visited := make(map[string]bool)
	for _, packageName := range packages {
		if visited[packageName] {
			continue
		}

		pkg, err := dr.registry.FindPackage(packageName)
		if err != nil {
			return nil, fmt.Errorf("package '%s' not found: %w", packageName, err)
		}

		if err := dr.buildDependencyGraph(ctx, pkg, combinedGraph, 0, visited); err != nil {
			return nil, fmt.Errorf("failed to build dependency graph for '%s': %w", packageName, err)
		}
	}

	// Check for circular dependencies
	if err := dr.checkCircularDependenciesInGraph(combinedGraph); err != nil {
		return nil, err
	}

	// Calculate levels and sort by dependency order
	dr.calculateLevels(combinedGraph)
	return dr.topologicalSort(combinedGraph), nil
}

// CheckCircularDependencies checks for circular dependencies among packages
func (dr *DependencyResolver) CheckCircularDependencies(ctx context.Context, packages []string) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, packageName := range packages {
		if !visited[packageName] {
			if err := dr.hasCircularDependency(ctx, packageName, visited, recStack); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateConstraints validates version constraints for dependencies
func (dr *DependencyResolver) ValidateConstraints(ctx context.Context, dependencies map[string]string) error {
	for packageName, constraint := range dependencies {
		pkg, err := dr.registry.FindPackage(packageName)
		if err != nil {
			return fmt.Errorf("dependency '%s' not found: %w", packageName, err)
		}

		if !dr.satisfiesConstraint(pkg.Version, constraint) {
			return fmt.Errorf("package '%s' version '%s' does not satisfy constraint '%s'",
				packageName, pkg.Version, constraint)
		}
	}

	return nil
}

// Private helper methods

func (dr *DependencyResolver) buildDependencyGraph(ctx context.Context, pkg *models.Package, graph *models.DependencyGraph, level int, visited map[string]bool) error {
	// Check if context is canceled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Check if already processed
	if visited[pkg.Name] {
		return nil
	}
	visited[pkg.Name] = true

	// Create node for this package
	node := &models.DependencyNode{
		Package:      pkg,
		Dependencies: []string{},
		Dependents:   []string{},
		Level:        level,
	}

	// Process dependencies
	for depName, constraint := range pkg.Dependencies {
		// Find dependency package
		depPkg, err := dr.registry.FindPackage(depName)
		if err != nil {
			// Skip missing dependencies with a warning
			fmt.Printf("Warning: dependency '%s' not found for package '%s'\n", depName, pkg.Name)
			continue
		}

		// Add dependency relationship
		node.Dependencies = append(node.Dependencies, depName)

		// Create edge
		edge := models.DependencyEdge{
			From:       pkg.Name,
			To:         depName,
			Constraint: constraint,
			Required:   true,
		}
		graph.Edges = append(graph.Edges, edge)

		// Add to dependents of the dependency
		if depNode, exists := graph.Nodes[depName]; exists {
			depNode.Dependents = append(depNode.Dependents, pkg.Name)
		}

		// Recursively process dependencies
		if err := dr.buildDependencyGraph(ctx, depPkg, graph, level+1, visited); err != nil {
			return err
		}
	}

	// Add node to graph
	graph.Nodes[pkg.Name] = node

	return nil
}

func (dr *DependencyResolver) calculateLevels(graph *models.DependencyGraph) {
	// Reset all levels
	for _, node := range graph.Nodes {
		node.Level = 0
	}

	// Calculate levels using topological sorting approach
	for {
		changed := false
		for name, node := range graph.Nodes {
			maxDepLevel := -1
			for _, depName := range node.Dependencies {
				if depNode, exists := graph.Nodes[depName]; exists {
					if depNode.Level > maxDepLevel {
						maxDepLevel = depNode.Level
					}
				}
			}

			newLevel := maxDepLevel + 1
			if newLevel != node.Level {
				node.Level = newLevel
				graph.Nodes[name] = node
				changed = true
			}
		}

		if !changed {
			break
		}
	}
}

func (dr *DependencyResolver) topologicalSort(graph *models.DependencyGraph) []string {
	// Create a slice of packages with their levels
	type packageLevel struct {
		name  string
		level int
	}

	var packages []packageLevel
	for name, node := range graph.Nodes {
		packages = append(packages, packageLevel{name: name, level: node.Level})
	}

	// Sort by level (dependencies first)
	sort.Slice(packages, func(i, j int) bool {
		if packages[i].level == packages[j].level {
			return packages[i].name < packages[j].name // Secondary sort by name for consistency
		}
		return packages[i].level < packages[j].level
	})

	// Extract package names
	var result []string
	for _, pkg := range packages {
		result = append(result, pkg.name)
	}

	return result
}

func (dr *DependencyResolver) hasCircularDependency(ctx context.Context, packageName string, visited, recStack map[string]bool) error {
	// Check if context is canceled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visited[packageName] = true
	recStack[packageName] = true

	// Find the package
	pkg, err := dr.registry.FindPackage(packageName)
	if err != nil {
		return fmt.Errorf("package '%s' not found: %w", packageName, err)
	}

	// Check all dependencies
	for depName := range pkg.Dependencies {
		if !visited[depName] {
			if err := dr.hasCircularDependency(ctx, depName, visited, recStack); err != nil {
				return err
			}
		} else if recStack[depName] {
			return fmt.Errorf("circular dependency detected: %s -> %s", packageName, depName)
		}
	}

	recStack[packageName] = false
	return nil
}

func (dr *DependencyResolver) checkCircularDependenciesInGraph(graph *models.DependencyGraph) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range graph.Nodes {
		if !visited[name] {
			if err := dr.detectCycleInGraph(name, graph, visited, recStack); err != nil {
				return err
			}
		}
	}

	return nil
}

func (dr *DependencyResolver) detectCycleInGraph(packageName string, graph *models.DependencyGraph, visited, recStack map[string]bool) error {
	visited[packageName] = true
	recStack[packageName] = true

	node, exists := graph.Nodes[packageName]
	if !exists {
		return nil
	}

	for _, depName := range node.Dependencies {
		if !visited[depName] {
			if err := dr.detectCycleInGraph(depName, graph, visited, recStack); err != nil {
				return err
			}
		} else if recStack[depName] {
			return fmt.Errorf("circular dependency detected: %s -> %s", packageName, depName)
		}
	}

	recStack[packageName] = false
	return nil
}

func (dr *DependencyResolver) satisfiesConstraint(version, constraint string) bool {
	// Simple constraint checking - can be enhanced with proper semver logic
	if constraint == "" || constraint == "*" {
		return true
	}

	// Handle simple version constraints
	if strings.HasPrefix(constraint, "^") {
		// Caret range (compatible version)
		requiredVersion := strings.TrimPrefix(constraint, "^")
		return dr.isCompatibleVersion(version, requiredVersion)
	}

	if strings.HasPrefix(constraint, "~") {
		// Tilde range (reasonably close version)
		requiredVersion := strings.TrimPrefix(constraint, "~")
		return dr.isReasonablyCloseVersion(version, requiredVersion)
	}

	if strings.HasPrefix(constraint, ">=") {
		requiredVersion := strings.TrimPrefix(constraint, ">=")
		return dr.compareVersions(version, requiredVersion) >= 0
	}

	if strings.HasPrefix(constraint, "<=") {
		requiredVersion := strings.TrimPrefix(constraint, "<=")
		return dr.compareVersions(version, requiredVersion) <= 0
	}

	if strings.HasPrefix(constraint, ">") {
		requiredVersion := strings.TrimPrefix(constraint, ">")
		return dr.compareVersions(version, requiredVersion) > 0
	}

	if strings.HasPrefix(constraint, "<") {
		requiredVersion := strings.TrimPrefix(constraint, "<")
		return dr.compareVersions(version, requiredVersion) < 0
	}

	// Exact version match
	return version == constraint
}

func (dr *DependencyResolver) isCompatibleVersion(version, required string) bool {
	// Simple caret range implementation
	versionParts := strings.Split(version, ".")
	requiredParts := strings.Split(required, ".")

	if len(versionParts) < 3 || len(requiredParts) < 3 {
		return version == required
	}

	// Major version must match
	if versionParts[0] != requiredParts[0] {
		return false
	}

	// Version should be >= required
	return dr.compareVersions(version, required) >= 0
}

func (dr *DependencyResolver) isReasonablyCloseVersion(version, required string) bool {
	// Simple tilde range implementation
	versionParts := strings.Split(version, ".")
	requiredParts := strings.Split(required, ".")

	if len(versionParts) < 3 || len(requiredParts) < 3 {
		return version == required
	}

	// Major and minor versions must match
	if versionParts[0] != requiredParts[0] || versionParts[1] != requiredParts[1] {
		return false
	}

	// Patch version should be >= required
	return dr.compareVersions(version, required) >= 0
}

func (dr *DependencyResolver) compareVersions(v1, v2 string) int {
	// Simple version comparison - can be enhanced with proper semver logic
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	// Pad with zeros if necessary
	for len(parts1) < maxLen {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < maxLen {
		parts2 = append(parts2, "0")
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		fmt.Sscanf(parts1[i], "%d", &n1)
		fmt.Sscanf(parts2[i], "%d", &n2)

		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}

	return 0
}

// GetDependencyTree returns a formatted dependency tree for a package
func (dr *DependencyResolver) GetDependencyTree(ctx context.Context, packageName string) (string, error) {
	graph, err := dr.Resolve(ctx, packageName)
	if err != nil {
		return "", err
	}

	rootNode, exists := graph.Nodes[packageName]
	if !exists {
		return "", fmt.Errorf("root package not found in dependency graph")
	}

	var tree strings.Builder
	dr.buildTreeString(&tree, rootNode, graph, "", make(map[string]bool))

	return tree.String(), nil
}

func (dr *DependencyResolver) buildTreeString(tree *strings.Builder, node *models.DependencyNode, graph *models.DependencyGraph, prefix string, visited map[string]bool) {
	tree.WriteString(fmt.Sprintf("%s%s (%s)\n", prefix, node.Package.Name, node.Package.Version))

	if visited[node.Package.Name] {
		return // Avoid infinite recursion
	}
	visited[node.Package.Name] = true

	for i, depName := range node.Dependencies {
		isLast := i == len(node.Dependencies)-1

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if depNode, exists := graph.Nodes[depName]; exists {
			dr.buildTreeString(tree, depNode, graph, prefix+connector, visited)
		} else {
			tree.WriteString(fmt.Sprintf("%s%s%s (not found)\n", prefix, connector, depName))
		}
	}
}
