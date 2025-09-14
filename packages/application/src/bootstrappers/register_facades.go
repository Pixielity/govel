package bootstrappers

import (
	"fmt"
	"sync"

	applicationInterfaces "govel/types/interfaces/application/base"
)

// RegisterFacades bootstrapper is responsible for registering facade aliases
// and setting up the facade system within the application. This allows for
// clean, static-like access to services through facades.
//
// This is the Go equivalent of Laravel's RegisterFacades bootstrapper.
type RegisterFacades struct {
	// aliases holds the facade aliases registered with the system
	aliases map[string]interface{}

	// packageAliases holds package-level facade aliases
	packageAliases map[string]interface{}

	// facadeApplication holds the reference to the application for facades
	facadeApplication interface{}

	// mutex protects concurrent access to aliases
	mutex sync.RWMutex

	// resolvedInstances caches resolved facade instances for performance
	resolvedInstances map[string]interface{}
}

// NewRegisterFacades creates a new RegisterFacades bootstrapper.
//
// Returns:
//
//	*RegisterFacades: A new bootstrapper instance
//
// Example:
//
//	bootstrapper := NewRegisterFacades()
func NewRegisterFacades() *RegisterFacades {
	return &RegisterFacades{
		aliases:           make(map[string]interface{}),
		packageAliases:    make(map[string]interface{}),
		resolvedInstances: make(map[string]interface{}),
	}
}

// Bootstrap bootstraps the given application with facade registration.
// This method follows Laravel's facade registration pattern:
// 1. Clear previously resolved facade instances
// 2. Set the facade application reference
// 3. Load and register aliases from configuration and packages
//
// Parameters:
//
//	application: The application instance
//
// Returns:
//
//	error: Any error that occurred during bootstrapping
//
// Example:
//
//	err := bootstrapper.Bootstrap(application)
//	if err != nil {
//	    log.Fatalf("Facade registration failed: %v", err)
//	}
func (r *RegisterFacades) Bootstrap(application applicationInterfaces.ApplicationInterface) error {
	// Step 1: Clear previously resolved instances
	r.clearResolvedInstances()

	// Step 2: Set the facade application reference
	r.setFacadeApplication(application)

	// Step 3: Load aliases from configuration
	configAliases, err := r.getConfigAliases()
	if err != nil {
		return fmt.Errorf("failed to get config aliases: %w", err)
	}

	// Step 4: Load package aliases (equivalent to Laravel's PackageManifest)
	packageAliases, err := r.getPackageAliases(application)
	if err != nil {
		return fmt.Errorf("failed to get package aliases: %w", err)
	}

	// Step 5: Merge aliases
	allAliases := r.mergeAliases(configAliases, packageAliases)

	// Step 6: Register all aliases
	if err := r.registerAliases(allAliases); err != nil {
		return fmt.Errorf("failed to register aliases: %w", err)
	}

	return nil
}

// clearResolvedInstances clears all previously resolved facade instances.
// This ensures that facades will re-resolve their services on next access.
//
// Example:
//
//	bootstrapper.clearResolvedInstances()
func (r *RegisterFacades) clearResolvedInstances() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Clear the resolved instances cache
	r.resolvedInstances = make(map[string]interface{})
}

// setFacadeApplication sets the application instance for use by facades.
//
// Parameters:
//
//	application: The application instance
//
// Example:
//
//	bootstrapper.setFacadeApplication(application)
func (r *RegisterFacades) setFacadeApplication(application applicationInterfaces.ApplicationInterface) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.facadeApplication = application
}

// getConfigAliases retrieves facade aliases from the application configuration.
//
// Returns:
//
//	map[string]interface{}: Configuration-based aliases
//	error: Any error that occurred
//
// Example:
//
//	aliases, err := bootstrapper.getConfigAliases(application)
func (r *RegisterFacades) getConfigAliases() (map[string]interface{}, error) {
	// This would typically call application.Make("config").Get("application.aliases", [])
	// For now, return a default set of aliases

	defaultAliases := map[string]interface{}{
		"App":        "govel/application/facades.App",
		"Config":     "govel/config/facades.Config",
		"Container":  "govel/container/facades.Container",
		"Log":        "govel/logger/facades.Log",
		"Auth":       "govel/support/facades.Auth",
		"Cache":      "govel/support/facades.Cache",
		"DB":         "govel/support/facades.DB",
		"Event":      "govel/support/facades.Event",
		"Hash":       "govel/support/facades.Hash",
		"Mail":       "govel/support/facades.Mail",
		"Queue":      "govel/support/facades.Queue",
		"Route":      "govel/support/facades.Route",
		"Storage":    "govel/support/facades.Storage",
		"Validation": "govel/support/facades.Validation",
		"View":       "govel/support/facades.View",
	}

	return defaultAliases, nil
}

// getPackageAliases retrieves facade aliases from package manifests.
// This is equivalent to Laravel's PackageManifest::aliases() method.
//
// Parameters:
//
//	application: The application instance
//
// Returns:
//
//	map[string]interface{}: Package-based aliases
//	error: Any error that occurred
//
// Example:
//
//	aliases, err := bootstrapper.getPackageAliases(application)
func (r *RegisterFacades) getPackageAliases(application applicationInterfaces.ApplicationInterface) (map[string]interface{}, error) {
	// This would load aliases from package manifests or discovery
	// For now, return empty as packages would register their own aliases

	packageAliases := make(map[string]interface{})

	// Example of how packages might register aliases:
	// packageAliases["CustomFacade"] = "vendor/package/facades.CustomFacade"

	return packageAliases, nil
}

// mergeAliases merges configuration and package aliases.
//
// Parameters:
//
//	configAliases: Aliases from configuration
//	packageAliases: Aliases from packages
//
// Returns:
//
//	map[string]interface{}: Merged aliases (config takes precedence)
//
// Example:
//
//	merged := bootstrapper.mergeAliases(configAliases, packageAliases)
func (r *RegisterFacades) mergeAliases(configAliases, packageAliases map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Add package aliases first
	for alias, target := range packageAliases {
		merged[alias] = target
	}

	// Add config aliases (overrides package aliases)
	for alias, target := range configAliases {
		merged[alias] = target
	}

	return merged
}

// registerAliases registers facade aliases with the system.
//
// Parameters:
//
//	aliases: Map of alias name to facade target
//
// Returns:
//
//	error: Any registration error
//
// Example:
//
//	err := bootstrapper.registerAliases(aliases)
func (r *RegisterFacades) registerAliases(aliases map[string]interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Store aliases for later use
	for alias, target := range aliases {
		r.aliases[alias] = target
	}

	// In a real implementation, you might:
	// 1. Register with a global alias loader
	// 2. Set up import aliases
	// 3. Configure facade resolution

	return nil
}

// GetRegisteredAliases returns all registered facade aliases.
//
// Returns:
//
//	map[string]interface{}: All registered aliases
//
// Example:
//
//	aliases := bootstrapper.GetRegisteredAliases()
//	fmt.Printf("Registered %d facade aliases\n", len(aliases))
func (r *RegisterFacades) GetRegisteredAliases() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Return a copy to prevent external modification
	aliases := make(map[string]interface{}, len(r.aliases))
	for k, v := range r.aliases {
		aliases[k] = v
	}

	return aliases
}

// IsAliasRegistered checks if a facade alias is registered.
//
// Parameters:
//
//	alias: The alias name to check
//
// Returns:
//
//	bool: true if the alias is registered
//
// Example:
//
//	if bootstrapper.IsAliasRegistered("App") {
//	    fmt.Println("App facade alias is registered")
//	}
func (r *RegisterFacades) IsAliasRegistered(alias string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.aliases[alias]
	return exists
}

// GetAliasTarget returns the target for a given alias.
//
// Parameters:
//
//	alias: The alias name
//
// Returns:
//
//	interface{}: The alias target (could be string, function, etc.)
//	bool: true if the alias exists
//
// Example:
//
//	target, exists := bootstrapper.GetAliasTarget("App")
//	if exists {
//	    fmt.Printf("App alias points to: %v\n", target)
//	}
func (r *RegisterFacades) GetAliasTarget(alias string) (interface{}, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	target, exists := r.aliases[alias]
	return target, exists
}

// RegisterCustomAlias registers a custom facade alias.
//
// Parameters:
//
//	alias: The alias name
//	target: The target (typically a package path or function)
//
// Example:
//
//	bootstrapper.RegisterCustomAlias("MyFacade", "mypackage/facades.MyFacade")
func (r *RegisterFacades) RegisterCustomAlias(alias string, target interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.aliases[alias] = target
}

// UnregisterAlias removes a facade alias.
//
// Parameters:
//
//	alias: The alias name to remove
//
// Returns:
//
//	bool: true if the alias was found and removed
//
// Example:
//
//	removed := bootstrapper.UnregisterAlias("OldFacade")
//	if removed {
//	    fmt.Println("Old facade alias removed")
//	}
func (r *RegisterFacades) UnregisterAlias(alias string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.aliases[alias]; exists {
		delete(r.aliases, alias)
		return true
	}

	return false
}

// ResolveFacade resolves a facade instance using the registered alias.
// This method provides the core facade resolution functionality.
//
// Parameters:
//
//	alias: The facade alias name
//
// Returns:
//
//	interface{}: The resolved facade instance
//	error: Any resolution error
//
// Example:
//
//	facade, err := bootstrapper.ResolveFacade("App")
//	if err != nil {
//	    log.Printf("Failed to resolve App facade: %v", err)
//	}
func (r *RegisterFacades) ResolveFacade(alias string) (interface{}, error) {
	r.mutex.RLock()

	// Check if already resolved and cached
	if cached, exists := r.resolvedInstances[alias]; exists {
		r.mutex.RUnlock()
		return cached, nil
	}

	// Get the alias target
	target, exists := r.aliases[alias]
	if !exists {
		r.mutex.RUnlock()
		return nil, fmt.Errorf("facade alias '%s' is not registered", alias)
	}

	r.mutex.RUnlock()

	// Resolve the facade instance
	instance, err := r.resolveFacadeTarget(target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve facade '%s': %w", alias, err)
	}

	// Cache the resolved instance
	r.mutex.Lock()
	r.resolvedInstances[alias] = instance
	r.mutex.Unlock()

	return instance, nil
}

// resolveFacadeTarget resolves a facade target to an actual instance.
//
// Parameters:
//
//	target: The facade target (string path, function, etc.)
//
// Returns:
//
//	interface{}: The resolved facade instance
//	error: Any resolution error
func (r *RegisterFacades) resolveFacadeTarget(target interface{}) (interface{}, error) {
	switch t := target.(type) {
	case string:
		// String target - could be a package path like "govel/application/facades.App"
		return r.resolveStringTarget(t)

	case func() interface{}:
		// Function target - call the function to get the instance
		return t(), nil

	case interface{}:
		// Direct instance target
		return t, nil

	default:
		return nil, fmt.Errorf("unsupported facade target type: %T", target)
	}
}

// resolveStringTarget resolves a string-based facade target.
//
// Parameters:
//
//	target: The string target (typically a package path)
//
// Returns:
//
//	interface{}: The resolved facade instance
//	error: Any resolution error
func (r *RegisterFacades) resolveStringTarget(target string) (interface{}, error) {
	// This would typically use reflection or import resolution
	// For now, return a placeholder that indicates the target

	return fmt.Sprintf("Facade[%s]", target), nil
}

// GetFacadeApplication returns the application instance set for facades.
//
// Returns:
//
//	interface{}: The facade application instance
//
// Example:
//
//	application := bootstrapper.GetFacadeApplication()
func (r *RegisterFacades) GetFacadeApplication() interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.facadeApplication
}

// ClearAllAliases removes all registered facade aliases.
//
// Example:
//
//	bootstrapper.ClearAllAliases()
func (r *RegisterFacades) ClearAllAliases() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.aliases = make(map[string]interface{})
	r.resolvedInstances = make(map[string]interface{})
}

// GetAliasCount returns the number of registered aliases.
//
// Returns:
//
//	int: Number of registered aliases
//
// Example:
//
//	count := bootstrapper.GetAliasCount()
//	fmt.Printf("Registered %d facade aliases\n", count)
func (r *RegisterFacades) GetAliasCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.aliases)
}

// GetResolvedInstanceCount returns the number of cached facade instances.
//
// Returns:
//
//	int: Number of cached instances
//
// Example:
//
//	count := bootstrapper.GetResolvedInstanceCount()
//	fmt.Printf("Cached %d facade instances\n", count)
func (r *RegisterFacades) GetResolvedInstanceCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.resolvedInstances)
}

// RefreshFacade clears the cached instance for a specific facade alias.
//
// Parameters:
//
//	alias: The facade alias to refresh
//
// Returns:
//
//	bool: true if a cached instance was found and cleared
//
// Example:
//
//	refreshed := bootstrapper.RefreshFacade("App")
//	if refreshed {
//	    fmt.Println("App facade cache cleared")
//	}
func (r *RegisterFacades) RefreshFacade(alias string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.resolvedInstances[alias]; exists {
		delete(r.resolvedInstances, alias)
		return true
	}

	return false
}

// ValidateAliases validates that all registered aliases can be resolved.
//
// Returns:
//
//	map[string]error: Map of alias to validation error (empty if all valid)
//
// Example:
//
//	errors := bootstrapper.ValidateAliases()
//	if len(errors) > 0 {
//	    fmt.Printf("Found %d invalid aliases\n", len(errors))
//	}
func (r *RegisterFacades) ValidateAliases() map[string]error {
	r.mutex.RLock()
	aliases := make(map[string]interface{}, len(r.aliases))
	for k, v := range r.aliases {
		aliases[k] = v
	}
	r.mutex.RUnlock()

	errors := make(map[string]error)

	for alias, target := range aliases {
		// Try to resolve each alias
		if _, err := r.resolveFacadeTarget(target); err != nil {
			errors[alias] = err
		}
	}

	return errors
}

// GetStatistics returns usage statistics for the facade system.
//
// Returns:
//
//	map[string]interface{}: Statistics including alias count, cached instances, etc.
//
// Example:
//
//	stats := bootstrapper.GetStatistics()
//	fmt.Printf("Facade Statistics: %+v\n", stats)
func (r *RegisterFacades) GetStatistics() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return map[string]interface{}{
		"registered_aliases": len(r.aliases),
		"resolved_instances": len(r.resolvedInstances),
		"facade_application": r.facadeApplication != nil,
		"memory_usage_bytes": r.estimateMemoryUsage(),
	}
}

// estimateMemoryUsage estimates the memory usage of the facade system.
//
// Returns:
//
//	int: Estimated memory usage in bytes
func (r *RegisterFacades) estimateMemoryUsage() int {
	// Rough estimation
	aliasSize := len(r.aliases) * 64               // Approximate per-alias overhead
	instanceSize := len(r.resolvedInstances) * 128 // Approximate per-instance overhead

	return aliasSize + instanceSize
}
