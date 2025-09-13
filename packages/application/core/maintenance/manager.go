package maintenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	containerInterfaces "govel/types/src/interfaces/container"
	"govel/types/src/constants/application"
	"govel/types/src/types/application"
)

// MaintenanceManager handles maintenance mode operations including
// enabling, disabling, and checking maintenance status.
type MaintenanceManager struct {
	// container provides access to dependency injection services
	container containerInterfaces.ContainerInterface

	// maintenanceFile is the path to the maintenance mode file
	maintenanceFile string

	// currentMode caches the current maintenance mode state
	currentMode *types.MaintenanceMode
}

// NewMaintenanceManager creates a new maintenance mode manager.
//
// Parameters:
//
//	container: The dependency injection container
//
// Returns:
//
//	*MaintenanceManager: A new maintenance manager instance
func NewMaintenanceManager(container containerInterfaces.ContainerInterface) *MaintenanceManager {
	// Resolve storage path from container
	storagePathService, err := container.Make("paths.storage")
	if err != nil {
		// If container resolution fails, use a default storage path
		// This should not happen if paths service provider is properly registered
		panic(fmt.Sprintf("failed to resolve storage path from container: %v", err))
	}

	storagePath, ok := storagePathService.(string)
	if !ok {
		// Type assertion should not fail if paths service provider returns string
		panic(fmt.Sprintf("storage path service returned unexpected type: %T", storagePathService))
	}

	maintenanceFile := filepath.Join(storagePath, constants.DirectoryFramework, constants.MaintenanceFileName)
	return &MaintenanceManager{
		container:       container,
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

// SetMaintenanceMessage updates the maintenance message.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	message: The new maintenance message
//
// Returns:
//
//	error: Any error that occurred while updating the message
//
// Example:
//
//	err := manager.SetMaintenanceMessage("System upgrade in progress")
func (mm *MaintenanceManager) SetMaintenanceMessage(message string) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	mode.Message = message
	return mm.saveMaintenanceMode(mode)
}

// SetRetryAfter updates the retry-after value.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	seconds: The new retry-after value in seconds
//
// Returns:
//
//	error: Any error that occurred while updating the retry-after value
//
// Example:
//
//	err := manager.SetRetryAfter(7200) // 2 hours
func (mm *MaintenanceManager) SetRetryAfter(seconds int) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	mode.RetryAfter = seconds
	return mm.saveMaintenanceMode(mode)
}

// AddAllowedIP adds an IP address to the bypass list.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	ip: The IP address to allow
//
// Returns:
//
//	error: Any error that occurred while adding the IP
//
// Example:
//
//	err := manager.AddAllowedIP("192.168.1.100")
func (mm *MaintenanceManager) AddAllowedIP(ip string) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	// Check if IP is already in the list
	for _, allowedIP := range mode.AllowedIPs {
		if allowedIP == ip {
			return nil // Already exists
		}
	}

	mode.AllowedIPs = append(mode.AllowedIPs, ip)
	return mm.saveMaintenanceMode(mode)
}

// RemoveAllowedIP removes an IP address from the bypass list.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	ip: The IP address to remove
//
// Returns:
//
//	error: Any error that occurred while removing the IP
//
// Example:
//
//	err := manager.RemoveAllowedIP("192.168.1.100")
func (mm *MaintenanceManager) RemoveAllowedIP(ip string) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	// Find and remove the IP
	for i, allowedIP := range mode.AllowedIPs {
		if allowedIP == ip {
			// Remove IP by slicing
			mode.AllowedIPs = append(mode.AllowedIPs[:i], mode.AllowedIPs[i+1:]...)
			return mm.saveMaintenanceMode(mode)
		}
	}

	return nil // IP was not in the list
}

// AddAllowedPath adds a URL path to the bypass list.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	path: The path to allow
//
// Returns:
//
//	error: Any error that occurred while adding the path
//
// Example:
//
//	err := manager.AddAllowedPath("/api/health")
func (mm *MaintenanceManager) AddAllowedPath(path string) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	// Check if path is already in the list
	for _, allowedPath := range mode.AllowedPaths {
		if allowedPath == path {
			return nil // Already exists
		}
	}

	mode.AllowedPaths = append(mode.AllowedPaths, path)
	return mm.saveMaintenanceMode(mode)
}

// RemoveAllowedPath removes a URL path from the bypass list.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	path: The path to remove
//
// Returns:
//
//	error: Any error that occurred while removing the path
//
// Example:
//
//	err := manager.RemoveAllowedPath("/api/health")
func (mm *MaintenanceManager) RemoveAllowedPath(path string) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	// Find and remove the path
	for i, allowedPath := range mode.AllowedPaths {
		if allowedPath == path {
			// Remove path by slicing
			mode.AllowedPaths = append(mode.AllowedPaths[:i], mode.AllowedPaths[i+1:]...)
			return mm.saveMaintenanceMode(mode)
		}
	}

	return nil // Path was not in the list
}

// SetMaintenanceData sets custom maintenance data.
// If not in maintenance mode, this method returns an error.
//
// Parameters:
//
//	key: The data key
//	value: The data value
//
// Returns:
//
//	error: Any error that occurred while setting the data
//
// Example:
//
//	err := manager.SetMaintenanceData("progress", 75)
func (mm *MaintenanceManager) SetMaintenanceData(key string, value interface{}) error {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return fmt.Errorf("application is not in maintenance mode")
	}

	if mode.Data == nil {
		mode.Data = make(map[string]interface{})
	}

	mode.Data[key] = value
	return mm.saveMaintenanceMode(mode)
}

// GetMaintenanceData gets custom maintenance data.
// Returns nil if the key is not found or not in maintenance mode.
//
// Parameters:
//
//	key: The data key
//
// Returns:
//
//	interface{}: The data value, or nil if not found
//
// Example:
//
//	progress := manager.GetMaintenanceData("progress")
func (mm *MaintenanceManager) GetMaintenanceData(key string) interface{} {
	mode := mm.MaintenanceMode()
	if mode == nil || mode.Data == nil {
		return nil
	}

	return mode.Data[key]
}

// GetMaintenanceInfo returns comprehensive maintenance information.
//
// Returns:
//
//	map[string]interface{}: Maintenance details
//
// Example:
//
//	info := manager.GetMaintenanceInfo()
//	fmt.Printf("Maintenance active: %v\n", info["active"])
func (mm *MaintenanceManager) GetMaintenanceInfo() map[string]interface{} {
	mode := mm.MaintenanceMode()
	if mode == nil {
		return map[string]interface{}{
			"active":              false,
			"maintenance_file":    mm.maintenanceFile,
			"file_exists":         false,
		}
	}

	duration := mm.MaintenanceDuration()
	info := map[string]interface{}{
		"active":              mode.Active,
		"message":             mode.Message,
		"retry_after":         mode.RetryAfter,
		"start_time":          mode.StartTime,
		"duration":            duration.String(),
		"duration_seconds":    int(duration.Seconds()),
		"estimated_duration":  mode.EstimatedDuration.String(),
		"maintenance_type":    mode.MaintenanceType,
		"allowed_ips_count":   len(mode.AllowedIPs),
		"allowed_paths_count": len(mode.AllowedPaths),
		"has_secret":          mode.Secret != "",
		"maintenance_file":    mm.maintenanceFile,
		"file_exists":         true,
	}

	if mode.AllowedIPs != nil {
		info["allowed_ips"] = mode.AllowedIPs
	}

	if mode.AllowedPaths != nil {
		info["allowed_paths"] = mode.AllowedPaths
	}

	if mode.Data != nil {
		info["data"] = mode.Data
	}

	return info
}

// saveMaintenanceMode saves the current maintenance mode to the maintenance file.
func (mm *MaintenanceManager) saveMaintenanceMode(mode *types.MaintenanceMode) error {
	data, err := json.MarshalIndent(mode, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance config: %w", err)
	}

	if err := os.WriteFile(mm.maintenanceFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write maintenance file: %w", err)
	}

	// Update cached mode
	mm.currentMode = mode
	return nil
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
