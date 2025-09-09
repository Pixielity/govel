package providers

import (
	"context"
	"govel/packages/application/interfaces"
)

// TerminatableServiceProvider is an alias for the interface in the interfaces package
type TerminatableServiceProvider = interfaces.TerminatableServiceProvider

// BaseTerminatableServiceProvider provides a base implementation for
// terminatable service providers. Most terminatable providers should
// embed this struct and override the Terminate method.
//
// Example usage:
//
//	type MetricsServiceProvider struct {
//		service_providers.BaseTerminatableServiceProvider
//	}
//
//	func (p *MetricsServiceProvider) Register(application ApplicationInterface) error {
//		return app.Singleton("metrics", func() interface{} {
//			return &MetricsCollector{}
//		})
//	}
//
//	func (p *MetricsServiceProvider) Terminate(ctx context.Context, application ApplicationInterface) error {
//		metrics, _ := app.Make("metrics")
//		return metrics.(*MetricsCollector).Flush()
//	}
type BaseTerminatableServiceProvider struct {
	BaseServiceProvider
}

// Terminate provides a default implementation that does nothing.
// Concrete terminatable providers should override this method to
// perform their specific cleanup operations.
//
// Parameters:
//
//	ctx: Context for controlling the termination process
//	app: The application instance
//
// Returns:
//
//	error: Always returns nil in the base implementation
func (p *BaseTerminatableServiceProvider) Terminate(ctx context.Context, application interfaces.ApplicationInterface) error {
	return nil
}

// TerminatePriority returns the default termination priority.
// Base implementation returns 200 (application service priority).
// Concrete providers can override this to specify their termination priority.
//
// Returns:
//
//	int: Default priority of 200
func (p *BaseTerminatableServiceProvider) TerminatePriority() int {
	return 200 // Default application service provider priority
}
