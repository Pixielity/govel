package services

import (
	"encoding/json"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"govel/package_manager/utils"
	"os"
	"path/filepath"
	"time"
)

const (
	StateFileName = "package_state.json"
	LockFileName  = "package_lock.json"
)

// StateManager implements StateManagerInterface for managing package state
type StateManager struct {
	stateDir string
}

// NewStateManager creates a new state manager instance
func NewStateManager(stateDir string) interfaces.StateManagerInterface {
	return &StateManager{
		stateDir: stateDir,
	}
}

// LoadState loads the current package state from disk
func (sm *StateManager) LoadState() (*models.PackageState, error) {
	// Ensure state directory exists
	if err := utils.EnsureDirectory(sm.stateDir); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	stateFile := filepath.Join(sm.stateDir, StateFileName)

	// If state file doesn't exist, return default state
	if !utils.FileExists(stateFile) {
		return sm.createDefaultState(), nil
	}

	// Read state file
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state models.PackageState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	// Validate and update state
	sm.updateStateCounts(&state)

	return &state, nil
}

// SaveState saves the current package state to disk
func (sm *StateManager) SaveState(state *models.PackageState) error {
	if state == nil {
		return fmt.Errorf("state cannot be nil")
	}

	// Ensure state directory exists
	if err := utils.EnsureDirectory(sm.stateDir); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Update counts and timestamp
	sm.updateStateCounts(state)
	state.UpdatedAt = time.Now()

	// Marshal state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write to state file
	stateFile := filepath.Join(sm.stateDir, StateFileName)
	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// UpdatePackageState updates the state of a specific package
func (sm *StateManager) UpdatePackageState(packageName string, isActive bool) error {
	state, err := sm.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	pkg, exists := state.Packages[packageName]
	if !exists {
		return fmt.Errorf("package '%s' not found in state", packageName)
	}

	pkg.IsActive = isActive
	pkg.UpdatedAt = time.Now()
	state.Packages[packageName] = pkg

	return sm.SaveState(state)
}

// GetPackageState returns the state of a specific package
func (sm *StateManager) GetPackageState(packageName string) (*models.Package, error) {
	state, err := sm.LoadState()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	pkg, exists := state.Packages[packageName]
	if !exists {
		return nil, fmt.Errorf("package '%s' not found in state", packageName)
	}

	return &pkg, nil
}

// CreateLockFile creates a lock file with the specified packages
func (sm *StateManager) CreateLockFile(packages []*models.Package) error {
	if packages == nil {
		packages = []*models.Package{}
	}

	// Create lock file structure
	lockData := struct {
		Version   string                     `json:"version"`
		Generated time.Time                  `json:"generated"`
		Packages  map[string]*models.Package `json:"packages"`
		Count     int                        `json:"count"`
	}{
		Version:   "1.0.0",
		Generated: time.Now(),
		Packages:  make(map[string]*models.Package),
		Count:     len(packages),
	}

	// Add packages to lock file
	for _, pkg := range packages {
		lockData.Packages[pkg.Name] = pkg
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(lockData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	// Ensure state directory exists
	if err := utils.EnsureDirectory(sm.stateDir); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Write lock file
	lockFile := filepath.Join(sm.stateDir, LockFileName)
	if err := os.WriteFile(lockFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

// ReadLockFile reads and returns packages from the lock file
func (sm *StateManager) ReadLockFile() ([]*models.Package, error) {
	lockFile := filepath.Join(sm.stateDir, LockFileName)

	// If lock file doesn't exist, return empty slice
	if !utils.FileExists(lockFile) {
		return []*models.Package{}, nil
	}

	// Read lock file
	data, err := os.ReadFile(lockFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lockData struct {
		Version   string                     `json:"version"`
		Generated time.Time                  `json:"generated"`
		Packages  map[string]*models.Package `json:"packages"`
		Count     int                        `json:"count"`
	}

	if err := json.Unmarshal(data, &lockData); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	// Convert to slice
	var packages []*models.Package
	for _, pkg := range lockData.Packages {
		packages = append(packages, pkg)
	}

	return packages, nil
}

// AddPackageToState adds a new package to the state
func (sm *StateManager) AddPackageToState(pkg *models.Package) error {
	if pkg == nil {
		return fmt.Errorf("package cannot be nil")
	}

	state, err := sm.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Add package to state
	state.Packages[pkg.Name] = *pkg

	return sm.SaveState(state)
}

// RemovePackageFromState removes a package from the state
func (sm *StateManager) RemovePackageFromState(packageName string) error {
	state, err := sm.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	delete(state.Packages, packageName)

	return sm.SaveState(state)
}

// GetStateFilePath returns the path to the state file
func (sm *StateManager) GetStateFilePath() string {
	return filepath.Join(sm.stateDir, StateFileName)
}

// GetLockFilePath returns the path to the lock file
func (sm *StateManager) GetLockFilePath() string {
	return filepath.Join(sm.stateDir, LockFileName)
}

// BackupState creates a backup of the current state
func (sm *StateManager) BackupState() error {
	state, err := sm.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(sm.stateDir, fmt.Sprintf("package_state_backup_%s.json", timestamp))

	// Marshal state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write backup file
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// RestoreState restores state from a backup file
func (sm *StateManager) RestoreState(backupFile string) error {
	// Read backup file
	data, err := os.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var state models.PackageState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to parse backup file: %w", err)
	}

	// Save restored state
	return sm.SaveState(&state)
}

// CleanupOldBackups removes old backup files (older than 30 days)
func (sm *StateManager) CleanupOldBackups() error {
	files, err := os.ReadDir(sm.stateDir)
	if err != nil {
		return fmt.Errorf("failed to read state directory: %w", err)
	}

	cutoff := time.Now().AddDate(0, 0, -30) // 30 days ago

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" &&
			(filepath.Base(file.Name()) != StateFileName && filepath.Base(file.Name()) != LockFileName) {

			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				filePath := filepath.Join(sm.stateDir, file.Name())
				if err := os.Remove(filePath); err != nil {
					fmt.Printf("Warning: failed to remove old backup file %s: %v\n", filePath, err)
				}
			}
		}
	}

	return nil
}

// ValidateState validates the integrity of the state
func (sm *StateManager) ValidateState() error {
	state, err := sm.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if state.Packages == nil {
		return fmt.Errorf("packages map is nil")
	}

	// Validate package consistency
	for name, pkg := range state.Packages {
		if pkg.Name != name {
			return fmt.Errorf("package name mismatch: map key '%s' != package name '%s'", name, pkg.Name)
		}

		if pkg.Name == "" {
			return fmt.Errorf("package has empty name")
		}

		if pkg.Version == "" {
			return fmt.Errorf("package '%s' has empty version", pkg.Name)
		}
	}

	return nil
}

// Private helper methods

func (sm *StateManager) createDefaultState() *models.PackageState {
	return &models.PackageState{
		Version:     "1.0.0",
		Packages:    make(map[string]models.Package),
		ActiveCount: 0,
		TotalCount:  0,
		UpdatedAt:   time.Now(),
	}
}

func (sm *StateManager) updateStateCounts(state *models.PackageState) {
	if state.Packages == nil {
		state.Packages = make(map[string]models.Package)
	}

	activeCount := 0
	totalCount := len(state.Packages)

	for _, pkg := range state.Packages {
		if pkg.IsActive {
			activeCount++
		}
	}

	state.ActiveCount = activeCount
	state.TotalCount = totalCount
}
