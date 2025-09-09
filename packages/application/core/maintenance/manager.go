package maintenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	"govel/packages/application/constants"
	"govel/packages/application/types"
)

// MaintenanceManager handles maintenance mode operations including
// enabling, disabling, and checking maintenance status.
type MaintenanceManager struct {
	// application holds a reference to the application instance
	application applicationInterfaces.ApplicationInterface

	// maintenanceFile is the path to the maintenance mode file
	maintenanceFile string

	// currentMode caches the current maintenance mode state
	currentMode *types.MaintenanceMode
}

// NewMaintenanceManager creates a new maintenance mode manager.
//
// Parameters:
//
//	application: The application instance
//
// Returns:
//
//	*MaintenanceManager: A new maintenance manager instance
func NewMaintenanceManager(application applicationInterfaces.ApplicationInterface) *MaintenanceManager {
	maintenanceFile := filepath.Join(application.StoragePath(), constants.DirectoryFramework, constants.MaintenanceFileName)
	return &MaintenanceManager{
		application:     application,
		maintenanceFile: maintenanceFile,
		currentMode:     nil,
	}
}

// IsDown returns whether the application is currently in maintenance mode.
// This method checks for the existence of the maintenance file and loads
// the maintenance configuration if it exists.
//
// Returns:
//
//	bool: true if the application is in maintenance mode, false otherwise
//
// Example:
//
//	if manager.IsDown() {
//	    // Show maintenance page
//	    return maintenanceResponse()
//	}
func (mm *MaintenanceManager) IsDown() bool {
	// Check if maintenance file exists
	if _, err := os.Stat(mm.maintenanceFile); os.IsNotExist(err) {
		mm.currentMode = nil
		return false
	}

	// Load maintenance configuration
	if mm.currentMode == nil {
		mm.loadMaintenanceMode()
	}

	return mm.currentMode != nil && mm.currentMode.Active
}

// IsUp returns whether the application is currently accessible (not in maintenance mode).
//
// Returns:
//
//	bool: true if the application is accessible, false if in maintenance mode
//
// Example:
//
//	if !manager.IsUp() {
//	    return errors.New("application is in maintenance mode")
//	}
func (mm *MaintenanceManager) IsUp() bool {
	return !mm.IsDown()
}

// Down puts the application into maintenance mode with the specified options.
// This method creates a maintenance file with the provided configuration.
//
// Parameters:
//
//	options: Maintenance mode configuration options
//
// Returns:
//
//	error: Any error that occurred while enabling maintenance mode
//
// Example:
//
//	err := manager.Down(&MaintenanceMode{
//	    Message:    "We're performing scheduled maintenance",
//	    RetryAfter: 3600, // 1 hour
//	    AllowedIPs: []string{"127.0.0.1", "::1"},
//	    Secret:     "secret-bypass-token",
//	})
func (mm *MaintenanceManager) Down(options *types.MaintenanceMode) error {
	if options == nil {
		options = &types.MaintenanceMode{}
	}

	// Set default values
	options.Active = constants.MaintenanceFileActiveValue
	options.StartTime = time.Now()

	if options.Message == "" {
		options.Message = constants.DefaultMaintenanceMessage
	}

	if options.RetryAfter == 0 {
		options.RetryAfter = constants.DefaultMaintenanceRetryAfter
	}

	// Ensure the storage/framework directory exists
	frameworkDir := filepath.Dir(mm.maintenanceFile)
	if err := os.MkdirAll(frameworkDir, 0755); err != nil {
		return fmt.Errorf("failed to create framework directory: %w", err)
	}

	// Write maintenance configuration to file
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance config: %w", err)
	}

	if err := os.WriteFile(mm.maintenanceFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write maintenance file: %w", err)
	}

	// Cache the current mode
	mm.currentMode = options

	return nil
}

// Up brings the application out of maintenance mode by removing the
// maintenance file and clearing the cached mode.
//
// Returns:
//
//	error: Any error that occurred while disabling maintenance mode
//
// Example:
//
//	err := manager.Up()
//	if err != nil {
//	    log.Printf("Failed to bring application up: %v", err)
//	}
func (mm *MaintenanceManager) Up() error {
	// Remove maintenance file
	if err := os.Remove(mm.maintenanceFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove maintenance file: %w", err)
	}

	// Clear cached mode
	mm.currentMode = nil

	return nil
}

// MaintenanceMode returns the current maintenance mode configuration.
// Returns nil if the application is not in maintenance mode.
//
// Returns:
//
//	*MaintenanceMode: Current maintenance configuration, or nil if not in maintenance
//
// Example:
//
//	mode := manager.MaintenanceMode()
//	if mode != nil {
//	    fmt.Printf("Maintenance message: %s\n", mode.Message)
//	    fmt.Printf("Retry after: %d seconds\n", mode.RetryAfter)
//	}
func (mm *MaintenanceManager) MaintenanceMode() *types.MaintenanceMode {
	if !mm.IsDown() {
		return nil
	}
	return mm.currentMode
}

// CanBypassMaintenance checks if the given request can bypass maintenance mode.
// This method checks IP addresses, paths, and secret tokens.
//
// Parameters:
//
//	clientIP: The client's IP address
//	path: The requested path
//	secret: The secret token provided (if any)
//
// Returns:
//
//	bool: true if the request can bypass maintenance mode, false otherwise
//
// Example:
//
//	if manager.CanBypassMaintenance(clientIP, requestPath, secretToken) {
//	    // Allow access
//	    return processRequest()
//	}
//	// Show maintenance page
//	return maintenanceResponse()
func (mm *MaintenanceManager) CanBypassMaintenance(clientIP, path, secret string) bool {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return true // Not in maintenance mode
	}

	// Check secret bypass
	if secret != "" && mode.Secret != "" && secret == mode.Secret {
		return true
	}

	// Check allowed IPs
	for _, allowedIP := range mode.AllowedIPs {
		if clientIP == allowedIP {
			return true
		}
	}

	// Check allowed paths
	for _, allowedPath := range mode.AllowedPaths {
		if path == allowedPath || (len(path) > len(allowedPath) &&
			path[:len(allowedPath)] == allowedPath && path[len(allowedPath)] == '/') {
			return true
		}
	}

	return false
}

// MaintenanceDuration returns how long the application has been in maintenance mode.
//
// Returns:
//
//	time.Duration: Duration since maintenance mode was activated, 0 if not in maintenance
//
// Example:
//
//	duration := manager.MaintenanceDuration()
//	if duration > 0 {
//	    fmt.Printf("Application has been in maintenance for %v\n", duration)
//	}
func (mm *MaintenanceManager) MaintenanceDuration() time.Duration {
	mode := mm.MaintenanceMode()
	if mode == nil || mode.StartTime.IsZero() {
		return 0
	}
	return time.Since(mode.StartTime)
}

// loadMaintenanceMode loads the maintenance configuration from the maintenance file.
func (mm *MaintenanceManager) loadMaintenanceMode() {
	data, err := os.ReadFile(mm.maintenanceFile)
	if err != nil {
		mm.currentMode = nil
		return
	}

	var mode types.MaintenanceMode
	if err := json.Unmarshal(data, &mode); err != nil {
		mm.currentMode = nil
		return
	}

	mm.currentMode = &mode
}
