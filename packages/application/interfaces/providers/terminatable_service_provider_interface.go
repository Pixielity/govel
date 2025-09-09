package interfaces

import (
	"context"
	applicationInterfaces "govel/packages/application/interfaces/application"
)

// TerminatableServiceProvider defines the interface for service providers that
// need to perform cleanup operations after the application response has been
// sent to the client. This interface follows Laravel's terminatable provider pattern.
//
// Terminatable providers are useful for performing tasks that should happen
// after the response has been delivered, such as:
// - Logging request/response data
// - Cleaning up temporary resources
// - Sending metrics or analytics
// - Performing background cleanup tasks
// - Closing connections that aren't needed for the response
// - Flushing caches and buffers
// - Updating statistics and counters
//
// This pattern allows for:
// - Faster response times by deferring non-critical operations
// - Better user experience through reduced response latency
// - Proper resource management and cleanup
// - Background task processing without blocking responses
//
// Example usage:
//
//	type LoggingServiceProvider struct {
//		service_providers.BaseServiceProvider
//	}
//
//	func (p *LoggingServiceProvider) Register(application ContainableInterface) error {
//		return application.Singleton("logger", func() interface{} {
//			return &Logger{}
//		})
//	}
//
//	func (p *LoggingServiceProvider) Terminate(ctx context.Context, application ContainableInterface) error {
//		logger, _ := application.Make("logger")
//		return logger.(*Logger).Flush()
//	}
//
//	func (p *LoggingServiceProvider) TerminatePriority() int {
//		return 10 // High priority for logging cleanup
//	}
//
// The interface promotes:
// - Clean separation between response generation and cleanup
// - Better user experience through faster response times
// - Proper resource management and cleanup
// - Background task processing
// - Ordered cleanup operations
type TerminatableServiceProvider interface {
	ServiceProviderInterface

	// Terminate is called after the response has been sent to the client.
	// This method should perform any cleanup operations, background tasks,
	// or post-response processing that the provider needs to handle.
	//
	// The method is called with a context that may have a timeout, so
	// implementations should respect context cancellation and perform
	// cleanup operations efficiently.
	//
	// Implementations should:
	// - Respect the provided context for cancellation
	// - Perform cleanup operations quickly and efficiently
	// - Handle errors gracefully without causing application issues
	// - Avoid long-running operations that could block shutdown
	// - Use timeouts for external service calls
	//
	// Common use cases for terminate:
	// - Flushing logs and metrics to external systems
	// - Cleaning up temporary files and resources
	// - Closing database connections and network resources
	// - Sending analytics data to tracking services
	// - Performing garbage collection operations
	// - Updating caches asynchronously
	// - Sending notifications about request completion
	//
	// Parameters:
	//   ctx: Context for controlling the termination process
	//   application: The application instance for service access
	//
	// Returns:
	//   error: Any error that occurred during termination, nil if successful
	//
	// Example:
	//   func (p *DatabaseServiceProvider) Terminate(ctx context.Context, application ContainableInterface) error {
	//       db, err := application.Make("database")
	//       if err != nil {
	//           return err
	//       }
	//
	//       select {
	//       case <-ctx.Done():
	//           return ctx.Err()
	//       default:
	//           return db.(*Database).CloseIdleConnections()
	//       }
	//   }
	Terminate(ctx context.Context, application applicationInterfaces.ApplicationInterface) error

	// TerminatePriority returns the priority for termination operations.
	// Lower values terminate first. This allows for proper cleanup ordering
	// when some providers depend on others for termination.
	//
	// Priority ordering ensures that:
	// - Critical infrastructure cleanup happens first
	// - Dependent services are cleaned up before their dependencies
	// - Non-critical tasks are handled last
	// - System resources are released in the correct order
	//
	// Priority levels (suggested convention):
	// - 0-99: Core infrastructure cleanup (logging, monitoring)
	// - 100-199: Application service cleanup (caches, sessions)
	// - 200-299: Database and storage cleanup (connections, files)
	// - 300-399: External service cleanup (APIs, third-party services)
	// - 400+: Non-critical background tasks (analytics, metrics)
	//
	// Returns:
	//   int: Priority level for termination ordering (lower = higher priority)
	//
	// Example:
	//   func (p *LoggingServiceProvider) TerminatePriority() int {
	//       return 10 // High priority for logging cleanup
	//   }
	//
	//   func (p *AnalyticsServiceProvider) TerminatePriority() int {
	//       return 450 // Low priority for analytics data
	//   }
	TerminatePriority() int
}
