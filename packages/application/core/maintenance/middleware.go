package maintenance

import "govel/packages/application/types"

// MaintenanceMiddleware provides middleware functionality for handling
// maintenance mode in HTTP applications.
type MaintenanceMiddleware struct {
	manager *MaintenanceManager
}

// NewMaintenanceMiddleware creates a new maintenance middleware.
//
// Parameters:
//
//	manager: The maintenance manager instance
//
// Returns:
//
//	*MaintenanceMiddleware: A new maintenance middleware instance
func NewMaintenanceMiddleware(manager *MaintenanceManager) *MaintenanceMiddleware {
	return &MaintenanceMiddleware{
		manager: manager,
	}
}

// CheckMaintenance checks if a request should be blocked due to maintenance mode.
// This method is designed to be used in HTTP middleware chains.
//
// Parameters:
//
//	clientIP: The client's IP address
//	path: The requested path
//	secret: The secret token provided (if any)
//
// Returns:
//
//	bool: true if request should be blocked, false if it can proceed
//	*MaintenanceMode: Maintenance configuration if blocked, nil otherwise
//
// Example:
//
//	blocked, mode := middleware.CheckMaintenance(clientIP, path, secret)
//	if blocked {
//	    return renderMaintenancePage(mode)
//	}
//	// Continue with normal request processing
func (mm *MaintenanceMiddleware) CheckMaintenance(clientIP, path, secret string) (bool, *types.MaintenanceMode) {
	if !mm.manager.IsDown() {
		return false, nil
	}

	if mm.manager.CanBypassMaintenance(clientIP, path, secret) {
		return false, nil
	}

	return true, mm.manager.MaintenanceMode()
}
