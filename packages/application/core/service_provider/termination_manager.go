package service_provider

import (
	"context"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
)

// Type aliases for cleaner usage
type TerminatableServiceProvider = providerInterfaces.TerminatableServiceProvider

// TerminationManager manages the termination process for terminatable
// service providers. This component handles the ordered execution of
// terminate methods after response delivery.
type TerminationManager struct {
	// providers holds the registered terminatable providers
	providers []TerminatableServiceProvider

	// application holds a reference to the application
	application applicationInterfaces.ApplicationInterface
}

// NewTerminationManager creates a new termination manager.
//
// Parameters:
//
//	application: The application instance
//
// Returns:
//
//	*TerminationManager: A new termination manager instance
func NewTerminationManager(application applicationInterfaces.ApplicationInterface) *TerminationManager {
	return &TerminationManager{
		providers:   make([]TerminatableServiceProvider, 0),
		application: application,
	}
}

// Register adds a terminatable service provider to the termination manager.
// The provider will be included in the termination process when Terminate
// is called.
//
// Parameters:
//
//	provider: The terminatable service provider to register
//
// Example:
//
//	loggingProvider := &LoggingServiceProvider{}
//	terminationManager.Register(loggingProvider)
func (tm *TerminationManager) Register(provider TerminatableServiceProvider) {
	tm.providers = append(tm.providers, provider)
}

// Terminate executes the terminate methods of all registered providers
// in priority order. This method should be called after the response
// has been sent to the client.
//
// The termination process executes providers in priority order (lower
// priorities first) and handles errors gracefully, allowing other
// providers to complete their termination even if some fail.
//
// Parameters:
//
//	ctx: Context for controlling the termination process
//
// Returns:
//
//	[]error: Any errors that occurred during termination
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	errors := terminationManager.Terminate(ctx)
//	for _, err := range errors {
//	    log.Printf("Termination error: %v", err)
//	}
func (tm *TerminationManager) Terminate(ctx context.Context) []error {
	// Sort providers by termination priority
	sortedProviders := tm.getSortedProviders()

	var errors []error

	// Execute terminate for each provider
	for _, provider := range sortedProviders {
		select {
		case <-ctx.Done():
			// Context cancelled, stop termination process
			return errors
		default:
		}

		if err := provider.Terminate(ctx, tm.application); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// TerminateAsync executes the termination process asynchronously.
// This method returns immediately and performs termination in a
// separate goroutine. Errors are handled by the provided error handler.
//
// Parameters:
//
//	ctx: Context for controlling the termination process
//	errorHandler: Function to handle termination errors (can be nil)
//
// Example:
//
//	terminationManager.TerminateAsync(ctx, func(errors []error) {
//	    for _, err := range errors {
//	        log.Printf("Async termination error: %v", err)
//	    }
//	})
func (tm *TerminationManager) TerminateAsync(ctx context.Context, errorHandler func([]error)) {
	go func() {
		errors := tm.Terminate(ctx)
		if errorHandler != nil && len(errors) > 0 {
			errorHandler(errors)
		}
	}()
}

// Count returns the number of registered terminatable providers.
//
// Returns:
//
//	int: Number of registered providers
func (tm *TerminationManager) Count() int {
	return len(tm.providers)
}

// getSortedProviders returns providers sorted by termination priority.
//
// Returns:
//
//	[]TerminatableServiceProvider: Providers sorted by priority
func (tm *TerminationManager) getSortedProviders() []TerminatableServiceProvider {
	sorted := make([]TerminatableServiceProvider, len(tm.providers))
	copy(sorted, tm.providers)

	// Simple bubble sort by priority
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].TerminatePriority() > sorted[j+1].TerminatePriority() {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}
