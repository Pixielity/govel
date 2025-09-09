package interfaces

import "context"

// ServiceProviderManagementInterface defines the contract for components that manage
// service providers within the application.
// This interface follows the Interface Segregation Principle by focusing
// solely on service provider management operations.
type ServiceProviderManagementInterface interface {
	// RegisterProvider registers a service provider with the application.
	//
	// Parameters:
	//   provider: The service provider to register
	//
	// Returns:
	//   error: Any error that occurred during registration, nil if successful
	RegisterProvider(provider ServiceProviderInterface) error

	// BootProviders boots all registered service providers.
	//
	// Parameters:
	//   ctx: Context for the boot operation
	//
	// Returns:
	//   error: Any error that occurred during booting, nil if successful
	BootProviders(ctx context.Context) error
}
