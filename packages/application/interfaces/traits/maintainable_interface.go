package interfaces

import (
	"govel/packages/application/types"
	"time"
)

/**
 * MaintainableInterface defines the contract for components that provide
 * application maintenance mode functionality. This interface follows the
 * Interface Segregation Principle by focusing solely on maintenance operations.
 */
type MaintainableInterface interface {
	/**
	 * IsDown returns whether the application is currently in maintenance mode.
	 *
	 * @return bool true if the application is in maintenance mode
	 */
	IsDown() bool

	/**
	 * IsUp returns whether the application is currently accessible (not in maintenance mode).
	 *
	 * @return bool true if the application is accessible
	 */
	IsUp() bool

	/**
	 * Down puts the application into maintenance mode.
	 *
	 * @param mode *MaintenanceMode The maintenance mode configuration
	 * @return error Any error that occurred while enabling maintenance mode
	 */
	Down(mode *types.MaintenanceMode) error

	/**
	 * Up brings the application out of maintenance mode.
	 *
	 * @return error Any error that occurred while disabling maintenance mode
	 */
	Up() error

	/**
	 * GetMaintenanceMode returns the current maintenance mode configuration.
	 *
	 * @return *MaintenanceMode Current maintenance configuration, or nil if not in maintenance
	 */
	GetMaintenanceMode() *types.MaintenanceMode

	/**
	 * CanBypass checks if a request can bypass maintenance mode.
	 *
	 * @param clientIP string The client's IP address
	 * @param path string The requested path
	 * @param secret string The secret token provided (if any)
	 * @return bool true if the request can bypass maintenance mode
	 */
	CanBypass(clientIP, path, secret string) bool

	/**
	 * GetMaintenanceDuration returns how long the application has been in maintenance mode.
	 *
	 * @return time.Duration Duration since maintenance mode was activated, 0 if not in maintenance
	 */
	GetMaintenanceDuration() time.Duration

	/**
	 * SetMaintenanceMessage updates the maintenance message.
	 *
	 * @param message string The new maintenance message
	 * @return error Any error that occurred while updating the message
	 */
	SetMaintenanceMessage(message string) error

	/**
	 * SetRetryAfter updates the retry-after value.
	 *
	 * @param seconds int The new retry-after value in seconds
	 * @return error Any error that occurred while updating the retry-after value
	 */
	SetRetryAfter(seconds int) error

	/**
	 * AddAllowedIP adds an IP address to the bypass list.
	 *
	 * @param ip string The IP address to allow
	 * @return error Any error that occurred while adding the IP
	 */
	AddAllowedIP(ip string) error

	/**
	 * RemoveAllowedIP removes an IP address from the bypass list.
	 *
	 * @param ip string The IP address to remove
	 * @return error Any error that occurred while removing the IP
	 */
	RemoveAllowedIP(ip string) error

	/**
	 * AddAllowedPath adds a URL path to the bypass list.
	 *
	 * @param path string The path to allow
	 * @return error Any error that occurred while adding the path
	 */
	AddAllowedPath(path string) error

	/**
	 * RemoveAllowedPath removes a URL path from the bypass list.
	 *
	 * @param path string The path to remove
	 * @return error Any error that occurred while removing the path
	 */
	RemoveAllowedPath(path string) error

	/**
	 * SetMaintenanceData sets custom maintenance data.
	 *
	 * @param key string The data key
	 * @param value interface{} The data value
	 * @return error Any error that occurred while setting the data
	 */
	SetMaintenanceData(key string, value interface{}) error

	/**
	 * GetMaintenanceData gets custom maintenance data.
	 *
	 * @param key string The data key
	 * @return interface{} The data value, or nil if not found
	 */
	GetMaintenanceData(key string) interface{}

	/**
	 * GetMaintenanceInfo returns comprehensive maintenance information.
	 *
	 * @return map[string]interface{} Maintenance details
	 */
	GetMaintenanceInfo() map[string]interface{}
}
