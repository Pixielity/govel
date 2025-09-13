package interfaces

import "time"

// MaintainableInterface defines the contract for maintenance mode functionality.
type MaintainableInterface interface {
	// IsDownForMaintenance returns whether the application is currently down for maintenance
	IsDownForMaintenance() bool
	
	// IsInMaintenanceMode returns whether the application is currently in maintenance mode
	IsInMaintenanceMode() bool
	
	// IsMaintenanceModeOff returns whether maintenance mode is currently disabled
	IsMaintenanceModeOff() bool
	
	// GetMaintenanceMode returns the current maintenance mode configuration
	GetMaintenanceMode() interface{} // Using interface{} to avoid circular imports
	
	// CanBypass checks if a request can bypass maintenance mode
	CanBypass(clientIP, path, secret string) bool
	
	// GetMaintenanceDuration returns how long the application has been in maintenance mode
	GetMaintenanceDuration() time.Duration
	
	// SetMaintenanceMessage updates the maintenance message
	SetMaintenanceMessage(message string) error
	
	// SetRetryAfter updates the retry-after value
	SetRetryAfter(seconds int) error
	
	// AddAllowedIP adds an IP address to the bypass list
	AddAllowedIP(ip string) error
	
	// RemoveAllowedIP removes an IP address from the bypass list
	RemoveAllowedIP(ip string) error
	
	// AddAllowedPath adds a URL path to the bypass list
	AddAllowedPath(path string) error
	
	// RemoveAllowedPath removes a URL path from the bypass list
	RemoveAllowedPath(path string) error
	
	// SetMaintenanceData sets custom maintenance data
	SetMaintenanceData(key string, value interface{}) error
	
	// GetMaintenanceData gets custom maintenance data
	GetMaintenanceData(key string) interface{}
	
	// GetMaintenanceInfo returns comprehensive maintenance information
	GetMaintenanceInfo() map[string]interface{}
}
