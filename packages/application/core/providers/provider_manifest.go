// Package providers provides core functionality for managing service provider manifests in the GoVel framework.
// This package handles the provider manifest structure and operations, following Laravel's provider manifest pattern.
package providers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	providerInterfaces "govel/packages/application/interfaces/providers"
	baseProviders "govel/packages/application/providers"
)

// ProviderManifest represents the provider manifest structure.
// This follows Laravel's provider manifest pattern, storing information about
// which providers are eager, deferred, and their provided services.
//
// The manifest is used to optimize application boot time by:
// - Loading only eager providers at startup
// - Deferring other providers until their services are requested
// - Mapping services to their providers for lazy loading
// - Storing event triggers for provider loading
type ProviderManifest struct {
	// Providers contains the complete list of registered provider types
	Providers []string `json:"providers"`

	// Eager contains providers that should be loaded immediately
	Eager []string `json:"eager"`

	// Deferred maps service names to their provider types
	Deferred map[string]string `json:"deferred"`

	// When maps provider types to events that trigger their loading
	When map[string][]string `json:"when"`
}

// ProviderManifestManager handles loading, saving, and managing provider manifests.
type ProviderManifestManager struct {
	// manifestPath is the path to the provider manifest file
	manifestPath string
}

// NewProviderManifestManager creates a new provider manifest manager.
//
// Parameters:
//   manifestPath: Path to the provider manifest file
//
// Returns:
//   *ProviderManifestManager: A new manifest manager instance
func NewProviderManifestManager(manifestPath string) *ProviderManifestManager {
	return &ProviderManifestManager{
		manifestPath: manifestPath,
	}
}

// LoadManifest loads the provider manifest from the file system.
// If the manifest file doesn't exist, returns a fresh empty manifest.
//
// Returns:
//   *ProviderManifest: The loaded manifest or a fresh one if file doesn't exist
//   error: Any error that occurred during manifest loading
func (pmm *ProviderManifestManager) LoadManifest() (*ProviderManifest, error) {
	if _, err := os.Stat(pmm.manifestPath); os.IsNotExist(err) {
		// Return fresh manifest if file doesn't exist
		return pmm.createFreshManifest(), nil
	}

	data, err := os.ReadFile(pmm.manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest ProviderManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	// Ensure When map exists
	if manifest.When == nil {
		manifest.When = make(map[string][]string)
	}

	// Ensure Deferred map exists
	if manifest.Deferred == nil {
		manifest.Deferred = make(map[string]string)
	}

	return &manifest, nil
}

// SaveManifest writes the provider manifest to disk.
//
// Parameters:
//   manifest: The manifest to write
//
// Returns:
//   error: Any error that occurred during writing
func (pmm *ProviderManifestManager) SaveManifest(manifest *ProviderManifest) error {
	// Ensure directory exists
	dir := filepath.Dir(pmm.manifestPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create manifest directory: %w", err)
	}

	data, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(pmm.manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}

// CompileManifest creates a fresh provider manifest by analyzing all providers.
// This method inspects each provider to determine if it's deferred and what services it provides.
//
// Parameters:
//   providers: List of provider instances to analyze
//
// Returns:
//   *ProviderManifest: The compiled manifest
//   error: Any error that occurred during compilation
func (pmm *ProviderManifestManager) CompileManifest(providers []providerInterfaces.ServiceProviderInterface) (*ProviderManifest, error) {
	manifest := &ProviderManifest{
		Providers: make([]string, 0, len(providers)),
		Eager:     make([]string, 0),
		Deferred:  make(map[string]string),
		When:      make(map[string][]string),
	}

	for _, provider := range providers {
		providerType := fmt.Sprintf("%T", provider)
		manifest.Providers = append(manifest.Providers, providerType)

		if baseProviders.IsProviderDeferred(provider) {
			// Handle deferred provider
			services := provider.GetProvides()
			for _, service := range services {
				manifest.Deferred[service] = providerType
			}

			// Check if provider implements EventTriggeredProvider for event triggers
			if eventTriggered, ok := provider.(providerInterfaces.EventTriggeredProvider); ok {
				events := eventTriggered.When()
				if len(events) > 0 {
					manifest.When[providerType] = events
				}
			}
		} else {
			// Handle eager provider
			manifest.Eager = append(manifest.Eager, providerType)
		}
	}

	// Save manifest to disk
	if err := pmm.SaveManifest(manifest); err != nil {
		return nil, fmt.Errorf("failed to save manifest: %w", err)
	}

	return manifest, nil
}

// ShouldRecompile determines if the manifest should be recompiled.
// This compares the current provider list with the manifest's provider list.
//
// Parameters:
//   manifest: The current manifest
//   providerTypes: The list of provider types to check against
//
// Returns:
//   bool: true if recompilation is needed
func (pmm *ProviderManifestManager) ShouldRecompile(manifest *ProviderManifest, providerTypes []string) bool {
	if manifest == nil {
		return true
	}

	// Check if provider lists match
	if len(manifest.Providers) != len(providerTypes) {
		return true
	}

	providerMap := make(map[string]bool)
	for _, provider := range providerTypes {
		providerMap[provider] = true
	}

	for _, provider := range manifest.Providers {
		if !providerMap[provider] {
			return true
		}
	}

	return false
}

// createFreshManifest creates a new empty provider manifest.
//
// Returns:
//   *ProviderManifest: A fresh empty manifest
func (pmm *ProviderManifestManager) createFreshManifest() *ProviderManifest {
	return &ProviderManifest{
		Providers: make([]string, 0),
		Eager:     make([]string, 0),
		Deferred:  make(map[string]string),
		When:      make(map[string][]string),
	}
}

// GetManifestPath returns the path to the manifest file.
//
// Returns:
//   string: The manifest file path
func (pmm *ProviderManifestManager) GetManifestPath() string {
	return pmm.manifestPath
}

// SetManifestPath sets the path to the manifest file.
//
// Parameters:
//   path: The new manifest file path
func (pmm *ProviderManifestManager) SetManifestPath(path string) {
	pmm.manifestPath = path
}

// ValidateManifest validates the integrity of a provider manifest.
//
// Parameters:
//   manifest: The manifest to validate
//
// Returns:
//   error: Any validation error found
func (pmm *ProviderManifestManager) ValidateManifest(manifest *ProviderManifest) error {
	if manifest == nil {
		return fmt.Errorf("manifest is nil")
	}

	// Validate that all deferred services have providers
	for service, providerType := range manifest.Deferred {
		if providerType == "" {
			return fmt.Errorf("service '%s' has empty provider type", service)
		}
	}

	// Validate that event-triggered providers exist in the provider list
	providerExists := make(map[string]bool)
	for _, providerType := range manifest.Providers {
		providerExists[providerType] = true
	}

	for providerType := range manifest.When {
		if !providerExists[providerType] {
			return fmt.Errorf("event-triggered provider '%s' not found in provider list", providerType)
		}
	}

	return nil
}

// GetEagerProviders returns the list of eager provider types from the manifest.
//
// Parameters:
//   manifest: The manifest to read from
//
// Returns:
//   []string: List of eager provider types
func (pmm *ProviderManifestManager) GetEagerProviders(manifest *ProviderManifest) []string {
	if manifest == nil {
		return []string{}
	}
	return manifest.Eager
}

// GetDeferredServices returns the deferred services map from the manifest.
//
// Parameters:
//   manifest: The manifest to read from
//
// Returns:
//   map[string]string: Map of service names to provider types
func (pmm *ProviderManifestManager) GetDeferredServices(manifest *ProviderManifest) map[string]string {
	if manifest == nil {
		return make(map[string]string)
	}
	return manifest.Deferred
}

// GetEventTriggers returns the event triggers map from the manifest.
//
// Parameters:
//   manifest: The manifest to read from
//
// Returns:
//   map[string][]string: Map of provider types to their triggering events
func (pmm *ProviderManifestManager) GetEventTriggers(manifest *ProviderManifest) map[string][]string {
	if manifest == nil {
		return make(map[string][]string)
	}
	return manifest.When
}

// GetProviderForService returns the provider type that provides a specific service.
//
// Parameters:
//   manifest: The manifest to search in
//   service: The service name to look up
//
// Returns:
//   string: The provider type that provides the service
//   bool: true if the service was found in the manifest
func (pmm *ProviderManifestManager) GetProviderForService(manifest *ProviderManifest, service string) (string, bool) {
	if manifest == nil {
		return "", false
	}
	
	providerType, exists := manifest.Deferred[service]
	return providerType, exists
}
