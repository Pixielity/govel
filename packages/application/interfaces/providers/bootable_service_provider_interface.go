package interfaces

import (
	"context"
	applicationInterfaces "govel/packages/application/interfaces/application"
)

// Use ContainableInterface from the service_provider_interface.go file to avoid redeclaration

// HasBootableServiceProvider defines the interface for service providers that
// need to participate in the application boot process with context support.
// This extends the basic ServiceProvider with context-aware boot methods.
//
// Service providers implementing this interface can perform context-aware
// initialization during the application boot process. This is particularly
// useful for providers that need to:
// - Respect cancellation and timeouts during boot
// - Perform asynchronous initialization tasks
// - Handle boot-time dependency resolution
// - Participate in graceful application startup
//
// Example usage:
//
//	type DatabaseServiceProvider struct {
//		service_providers.BaseServiceProvider
//	}
//
//	func (p *DatabaseServiceProvider) Register(application ContainableInterface) error {
//		return application.Singleton("database", func() interface{} {
//			return &DatabaseConnection{}
//		})
//	}
//
//	func (p *DatabaseServiceProvider) BootWithContext(ctx context.Context, application ContainableInterface) error {
//		db, _ := application.Make("database")
//		return db.(*DatabaseConnection).ConnectWithContext(ctx)
//	}
//
// The interface promotes:
// - Context-aware initialization
// - Graceful handling of boot timeouts
// - Proper resource management during startup
// - Better error handling in boot processes
type HasBootableServiceProvider interface {
	ServiceProviderInterface

	// BootWithContext is called during the boot process with a context.
	// This allows providers to respect cancellation and timeouts during boot.
	//
	// The method is called after all service providers have been registered
	// but before the application starts serving requests. This is the ideal
	// place to perform initialization tasks that require other services to
	// be available in the container.
	//
	// Implementations should:
	// - Respect the provided context for cancellation
	// - Perform initialization tasks efficiently
	// - Return errors for any boot failures
	// - Avoid blocking operations without context checking
	//
	// Parameters:
	//   ctx: Context for controlling the boot process, may include timeout
	//   application: The application instance for service access
	//
	// Returns:
	//   error: Any error that occurred during boot, nil if successful
	//
	// Example:
	//   func (p *CacheServiceProvider) BootWithContext(ctx context.Context, application ContainableInterface) error {
	//       cache, err := application.Make("cache")
	//       if err != nil {
	//           return err
	//       }
	//
	//       select {
	//       case <-ctx.Done():
	//           return ctx.Err()
	//       default:
	//           return cache.(*Cache).Initialize()
	//       }
	//   }
	BootWithContext(ctx context.Context, application applicationInterfaces.ApplicationInterface) error
}
