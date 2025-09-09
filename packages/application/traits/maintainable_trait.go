package traits

import (
	"fmt"
	"time"

	"govel/packages/application/core/maintenance"
	traitInterfaces "govel/packages/application/interfaces/traits"
	"govel/packages/application/types"
)

/**
 * Maintainable provides application maintenance mode functionality by wrapping
 * the core maintenance manager. This trait follows the self-contained pattern
 * and delegates all operations to the underlying manager.
 */
type Maintainable struct {
	/**
	 * manager is the underlying maintenance manager instance
	 */
	manager *maintenance.MaintenanceManager

	/**
	 * storagePath is the path where maintenance files are stored
	 */
	storagePath string
}

/**
 * NewMaintainable creates a new Maintainable instance.
 *
 * @param storagePath string The storage path for maintenance files
 * @return *Maintainable The newly created trait instance
 */
func NewMaintainable(storagePath string) *Maintainable {
	return &Maintainable{
		storagePath: storagePath,
		// manager will be created with application reference when needed
	}
}

/**
 * NewMaintainableWithManager creates a new Maintainable with an existing manager.
 *
 * @param manager *maintenance.MaintenanceManager The maintenance manager to wrap
 * @return *Maintainable The newly created trait instance
 */
func NewMaintainableWithManager(manager *maintenance.MaintenanceManager) *Maintainable {
	return &Maintainable{
		manager: manager,
	}
}

/**
 * SetManager sets the maintenance manager for this trait.
 * This is typically called once the application instance is available.
 *
 * @param manager *maintenance.MaintenanceManager The manager to use
 */
func (t *Maintainable) SetManager(manager *maintenance.MaintenanceManager) {
	t.manager = manager
}

/**
 * GetManager returns the underlying maintenance manager instance.
 *
 * @return *maintenance.MaintenanceManager The underlying manager
 */
func (t *Maintainable) GetManager() *maintenance.MaintenanceManager {
	return t.manager
}

/**
 * IsDown returns whether the application is currently in maintenance mode.
 *
 * @return bool true if the application is in maintenance mode
 */
func (t *Maintainable) IsDown() bool {
	if t.manager == nil {
		return false
	}
	return t.manager.IsDown()
}

/**
 * IsUp returns whether the application is currently accessible (not in maintenance mode).
 *
 * @return bool true if the application is accessible
 */
func (t *Maintainable) IsUp() bool {
	if t.manager == nil {
		return true
	}
	return t.manager.IsUp()
}

/**
 * Down puts the application into maintenance mode.
 *
 * @param mode *MaintenanceMode The maintenance mode configuration
 * @return error Any error that occurred while enabling maintenance mode
 */
func (t *Maintainable) Down(mode *types.MaintenanceMode) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	// Note: This would need type conversion between types.MaintenanceMode
	// and the actual type expected by the manager
	// For now, we'll assume they're compatible or handle the conversion
	return t.manager.Down(mode)
}

/**
 * Up brings the application out of maintenance mode.
 *
 * @return error Any error that occurred while disabling maintenance mode
 */
func (t *Maintainable) Up() error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}
	return t.manager.Up()
}

/**
 * GetMaintenanceMode returns the current maintenance mode configuration.
 *
 * @return *MaintenanceMode Current maintenance configuration, or nil if not in maintenance
 */
func (t *Maintainable) GetMaintenanceMode() *types.MaintenanceMode {
	if t.manager == nil {
		return nil
	}

	// Note: This would need type conversion between the manager's type
	// and types.MaintenanceMode
	mode := t.manager.MaintenanceMode()
	if mode == nil {
		return nil
	}

	// Convert to interface type (simplified - actual conversion would depend on types)
	return &types.MaintenanceMode{
		Active:            mode.Active,
		Message:           mode.Message,
		RetryAfter:        mode.RetryAfter,
		AllowedIPs:        mode.AllowedIPs,
		AllowedPaths:      mode.AllowedPaths,
		Secret:            mode.Secret,
		StartTime:         mode.StartTime,
		EstimatedDuration: mode.EstimatedDuration,
		MaintenanceType:   mode.MaintenanceType,
		Data:              mode.Data,
	}
}

/**
 * CanBypass checks if a request can bypass maintenance mode.
 *
 * @param clientIP string The client's IP address
 * @param path string The requested path
 * @param secret string The secret token provided (if any)
 * @return bool true if the request can bypass maintenance mode
 */
func (t *Maintainable) CanBypass(clientIP, path, secret string) bool {
	if t.manager == nil {
		return true
	}
	return t.manager.CanBypassMaintenance(clientIP, path, secret)
}

/**
 * GetMaintenanceDuration returns how long the application has been in maintenance mode.
 *
 * @return time.Duration Duration since maintenance mode was activated, 0 if not in maintenance
 */
func (t *Maintainable) GetMaintenanceDuration() time.Duration {
	if t.manager == nil {
		return 0
	}
	return t.manager.MaintenanceDuration()
}

/**
 * SetMaintenanceMessage updates the maintenance message.
 * Note: This is a simplified implementation - the actual manager may not support this directly.
 *
 * @param message string The new maintenance message
 * @return error Any error that occurred while updating the message
 */
func (t *Maintainable) SetMaintenanceMessage(message string) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	// This would require getting current mode, updating message, and saving
	// For now, return not implemented
	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * SetRetryAfter updates the retry-after value.
 * Note: This is a simplified implementation - the actual manager may not support this directly.
 *
 * @param seconds int The new retry-after value in seconds
 * @return error Any error that occurred while updating the retry-after value
 */
func (t *Maintainable) SetRetryAfter(seconds int) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	// This would require getting current mode, updating retry-after, and saving
	// For now, return not implemented
	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * AddAllowedIP adds an IP address to the bypass list.
 * Note: This is a simplified implementation - the actual manager may not support this directly.
 *
 * @param ip string The IP address to allow
 * @return error Any error that occurred while adding the IP
 */
func (t *Maintainable) AddAllowedIP(ip string) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	// This would require getting current mode, updating IPs, and saving
	// For now, return not implemented
	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * RemoveAllowedIP removes an IP address from the bypass list.
 *
 * @param ip string The IP address to remove
 * @return error Any error that occurred while removing the IP
 */
func (t *Maintainable) RemoveAllowedIP(ip string) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * AddAllowedPath adds a URL path to the bypass list.
 *
 * @param path string The path to allow
 * @return error Any error that occurred while adding the path
 */
func (t *Maintainable) AddAllowedPath(path string) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * RemoveAllowedPath removes a URL path from the bypass list.
 *
 * @param path string The path to remove
 * @return error Any error that occurred while removing the path
 */
func (t *Maintainable) RemoveAllowedPath(path string) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * SetMaintenanceData sets custom maintenance data.
 *
 * @param key string The data key
 * @param value interface{} The data value
 * @return error Any error that occurred while setting the data
 */
func (t *Maintainable) SetMaintenanceData(key string, value interface{}) error {
	if t.manager == nil {
		return fmt.Errorf("maintenance manager not initialized")
	}

	return fmt.Errorf("not implemented - use Down() with updated configuration")
}

/**
 * GetMaintenanceData gets custom maintenance data.
 *
 * @param key string The data key
 * @return interface{} The data value, or nil if not found
 */
func (t *Maintainable) GetMaintenanceData(key string) interface{} {
	if t.manager == nil {
		return nil
	}

	mode := t.GetMaintenanceMode()
	if mode == nil || mode.Data == nil {
		return nil
	}

	return mode.Data[key]
}

/**
 * GetMaintenanceInfo returns comprehensive maintenance information.
 *
 * @return map[string]interface{} Maintenance details
 */
func (t *Maintainable) GetMaintenanceInfo() map[string]interface{} {
	if t.manager == nil {
		return map[string]interface{}{
			"is_down":        false,
			"is_up":          true,
			"has_config":     false,
			"manager_loaded": false,
		}
	}

	info := map[string]interface{}{
		"is_down":        t.manager.IsDown(),
		"is_up":          t.manager.IsUp(),
		"manager_loaded": true,
	}

	mode := t.GetMaintenanceMode()
	if mode != nil {
		info["has_config"] = true
		info["message"] = mode.Message
		info["retry_after"] = mode.RetryAfter
		info["allowed_ips_count"] = len(mode.AllowedIPs)
		info["allowed_paths_count"] = len(mode.AllowedPaths)
		info["has_secret"] = mode.Secret != ""
		info["start_time"] = mode.StartTime
		info["duration"] = t.GetMaintenanceDuration()
		info["maintenance_type"] = mode.MaintenanceType
		info["estimated_duration"] = mode.EstimatedDuration

		if mode.Data != nil {
			info["custom_data_keys"] = len(mode.Data)
		}
	} else {
		info["has_config"] = false
	}

	return info
}

// Compile-time interface compliance check
var _ traitInterfaces.MaintainableInterface = (*Maintainable)(nil)
