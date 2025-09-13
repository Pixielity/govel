package maintenance

import (
	"time"

	"govel/types/src/types/application"
	applicationInterfaces "govel/types/src/interfaces/application"
	containerInterfaces "govel/types/src/interfaces/container"
)

/**
 * Maintainable provides application maintenance mode functionality by wrapping
 * the core maintenance manager. This trait follows the self-contained pattern
 * and delegates all operations to the underlying manager.
 *
 * The trait implements the MaintainableInterface and serves as a lightweight
 * wrapper around the MaintenanceManager, ensuring all maintenance functionality
 * is handled by the manager while providing a clean trait interface.
 */
type Maintainable struct {
	/**
	 * manager is the underlying maintenance manager instance that handles
	 * all maintenance operations including file I/O, state management,
	 * and bypass logic.
	 */
	manager *MaintenanceManager
}

// NewMaintainable creates a new maintainable trait with the provided container.
//
// Parameters:
//
//	container: The dependency injection container
//
// Returns:
//
//	*Maintainable: A new maintainable trait instance
//
// Example:
//
//	maintainable := NewMaintainable(container)
func NewMaintainable(container containerInterfaces.ContainerInterface) *Maintainable {
	return &Maintainable{
		manager: NewMaintenanceManager(container),
	}
}

// IsDownForMaintenance returns whether the application is currently down for maintenance.
// This method is provided for Laravel compatibility.
// Delegates to the manager.
func (m *Maintainable) IsDownForMaintenance() bool {
	return m.manager.IsDown()
}

// IsInMaintenanceMode returns whether the application is currently in maintenance mode.
// Delegates to the manager.
func (m *Maintainable) IsInMaintenanceMode() bool {
	return m.manager.IsDown()
}

// IsMaintenanceModeOff returns whether maintenance mode is currently disabled.
// Delegates to the manager.
func (m *Maintainable) IsMaintenanceModeOff() bool {
	return m.manager.IsUp()
}

// GetMaintenanceMode returns the current maintenance mode configuration.
// Delegates to the manager.
func (m *Maintainable) GetMaintenanceMode() *types.MaintenanceMode {
	return m.manager.MaintenanceMode()
}

// CanBypass checks if a request can bypass maintenance mode.
// Delegates to the manager.
func (m *Maintainable) CanBypass(clientIP, path, secret string) bool {
	return m.manager.CanBypassMaintenance(clientIP, path, secret)
}

// GetMaintenanceDuration returns how long the application has been in maintenance mode.
// Delegates to the manager.
func (m *Maintainable) GetMaintenanceDuration() time.Duration {
	return m.manager.MaintenanceDuration()
}

// SetMaintenanceMessage updates the maintenance message.
// Delegates to the manager.
func (m *Maintainable) SetMaintenanceMessage(message string) error {
	return m.manager.SetMaintenanceMessage(message)
}

// SetRetryAfter updates the retry-after value.
// Delegates to the manager.
func (m *Maintainable) SetRetryAfter(seconds int) error {
	return m.manager.SetRetryAfter(seconds)
}

// AddAllowedIP adds an IP address to the bypass list.
// Delegates to the manager.
func (m *Maintainable) AddAllowedIP(ip string) error {
	return m.manager.AddAllowedIP(ip)
}

// RemoveAllowedIP removes an IP address from the bypass list.
// Delegates to the manager.
func (m *Maintainable) RemoveAllowedIP(ip string) error {
	return m.manager.RemoveAllowedIP(ip)
}

// AddAllowedPath adds a URL path to the bypass list.
// Delegates to the manager.
func (m *Maintainable) AddAllowedPath(path string) error {
	return m.manager.AddAllowedPath(path)
}

// RemoveAllowedPath removes a URL path from the bypass list.
// Delegates to the manager.
func (m *Maintainable) RemoveAllowedPath(path string) error {
	return m.manager.RemoveAllowedPath(path)
}

// SetMaintenanceData sets custom maintenance data.
// Delegates to the manager.
func (m *Maintainable) SetMaintenanceData(key string, value interface{}) error {
	return m.manager.SetMaintenanceData(key, value)
}

// GetMaintenanceData gets custom maintenance data.
// Delegates to the manager.
func (m *Maintainable) GetMaintenanceData(key string) interface{} {
	return m.manager.GetMaintenanceData(key)
}

// GetMaintenanceInfo returns comprehensive maintenance information.
// Delegates to the manager.
func (m *Maintainable) GetMaintenanceInfo() map[string]interface{} {
	return m.manager.GetMaintenanceInfo()
}

// Compile-time interface compliance check
var _ applicationInterfaces.MaintainableInterface = (*Maintainable)(nil)
