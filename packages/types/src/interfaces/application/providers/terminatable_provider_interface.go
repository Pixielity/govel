package interfaces

import (
	"context"
	applicationInterfaces "govel/packages/types/src/interfaces/application"
)

// TerminatableProvider defines the contract for service providers that require
// graceful termination. This interface allows providers to clean up resources,
// close connections, and perform other shutdown procedures when the application
// is terminating.
//
// This is particularly useful for providers that manage:
// - Database connections
// - Cache connections
// - Message queues
// - Background workers
// - File handles
// - Network connections
type TerminatableProvider interface {
	ServiceProviderInterface

	// Terminate performs graceful shutdown operations for the service provider.
	// This method is called when the application is shutting down and should be used
	// to clean up resources, close connections, and perform any necessary cleanup.
	//
	// The context parameter allows for timeout control and cancellation during
	// the termination process. Providers should respect the context deadline
	// and attempt to complete their cleanup within the given time.
	//
	// Parameters:
	//   ctx: Context for controlling termination timeout and cancellation
	//   app: The application instance for accessing services during termination
	//
	// Returns:
	//   error: Any error that occurred during termination
	//
	// Example:
	//   func (p *DatabaseProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	//       db, _ := app.GetContainer().Make("database")
	//       return db.Close()
	//   }
	Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error
}
