// Package providers contains service provider implementations for the pipeline package.
// Service providers are responsible for registering pipeline services in the
// dependency injection container and configuring them for use throughout the application.
package providers

import (
	"fmt"

	"govel/application/providers"
	applicationInterfaces "govel/types/interfaces/application/base"
	pipelineInterfaces "govel/types/interfaces/pipeline"

	pipeline "govel/pipeline"
)

// PipelineServiceProvider implements a Laravel-compatible service provider
// for the pipeline package. It registers pipeline services in the dependency
// injection container and provides factory methods for creating pipeline instances.
//
// This service provider is equivalent to Laravel's PipelineServiceProvider and
// follows the same patterns for service registration and binding.
//
// Services registered:
//   - PIPELINE_TOKEN: Main pipeline service for creating new Pipeline instances
//   - PIPELINE_FACTORY_TOKEN: Explicit pipeline factory for factory pattern usage
//   - PIPELINE_HUB_TOKEN: Singleton Hub instance for managing named pipelines
//   - PIPELINE_HUB_CONTRACT_TOKEN: Hub interface contract for dependency injection
//
// The provider implements deferred loading, meaning services are only
// instantiated when first requested from the container.
type PipelineServiceProvider struct {
	providers.ServiceProvider
}

// NewPipelineServiceProvider creates a new PipelineServiceProvider instance.
//
// Returns:
//   - *PipelineServiceProvider: New service provider instance ready for registration
//
// Example:
//
//	provider := NewPipelineServiceProvider()
//	err := provider.Register(application)
//	if err != nil {
//		log.Fatal("Failed to register pipeline services:", err)
//	}
func NewPipelineServiceProvider() *PipelineServiceProvider {
	return &PipelineServiceProvider{
		ServiceProvider: providers.ServiceProvider{},
	}
}

// Register registers all pipeline services in the dependency injection container.
// This method sets up the service bindings but does not instantiate the services
// immediately (deferred loading).
//
// Registered services:
//   - PIPELINE_TOKEN: Factory function that creates new Pipeline instances
//   - PIPELINE_FACTORY_TOKEN: Explicit pipeline factory for factory pattern usage
//   - PIPELINE_HUB_TOKEN: Singleton Hub instance for managing named pipelines
//   - PIPELINE_HUB_CONTRACT_TOKEN: Hub interface contract for dependency injection
//
// Returns:
//   - error: Any error that occurred during service registration
//
// Thread-safe: This method should only be called once during application bootstrap.
func (p *PipelineServiceProvider) Register(application applicationInterfaces.ApplicationInterface) error {
	// Call parent Register method to set the registered flag
	if err := p.ServiceProvider.Register(application); err != nil {
		return fmt.Errorf("failed to register base service provider: %w", err)
	}

	// Register pipeline factory as a non-singleton binding using token
	// Each time the pipeline token is resolved, a new Pipeline instance is created
	err := application.Bind(pipelineInterfaces.PIPELINE_TOKEN, p.createPipelineFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind pipeline factory: %w", err)
	}

	// Register pipeline factory with explicit factory token for factory pattern usage
	err = application.Bind(pipelineInterfaces.PIPELINE_FACTORY_TOKEN, p.createPipelineFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind pipeline factory token: %w", err)
	}

	// Register Hub as singleton using token
	err = application.Singleton(pipelineInterfaces.PIPELINE_HUB_TOKEN, p.createHubFactory(application))
	if err != nil {
		return fmt.Errorf("failed to bind pipeline hub: %w", err)
	}

	return nil
}

// Provides returns a list of service names that this provider offers.
// This is used by the container to determine which services are available
// and to implement deferred loading.
//
// Returns:
//   - []string: List of service names provided by this provider
func (p *PipelineServiceProvider) Provides() []string {
	return []string{
		pipelineInterfaces.PIPELINE_TOKEN.String(),         // Pipeline service
		pipelineInterfaces.PIPELINE_FACTORY_TOKEN.String(), // Pipeline factory
		pipelineInterfaces.PIPELINE_HUB_TOKEN.String(),     // Hub service
	}
}

/**
 * createPipelineFactory creates a factory function for Pipeline instances.
 * This factory returns a new Pipeline instance each time it's called,
 * enabling transient pipeline creation for different processing contexts.
 *
 * Each pipeline instance maintains its own state and can be configured
 * independently with different pipes, contexts, and destinations.
 *
 * @method createPipelineFactory
 * @param application applicationInterfaces.ApplicationInterface The application container instance
 * @return func() interface{} Factory function that creates new Pipeline instances
 */
func (p *PipelineServiceProvider) createPipelineFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		return pipeline.NewPipeline(application)
	}
}

/**
 * createHubFactory creates a factory function for the Hub singleton.
 * This factory returns the same Hub instance on subsequent calls,
 * ensuring a centralized registry for all named pipeline configurations.
 *
 * The Hub singleton maintains all registered pipeline configurations
 * and provides a consistent interface for pipeline execution across
 * the entire application lifecycle.
 *
 * @method createHubFactory
 * @param application applicationInterfaces.ApplicationInterface The application container instance
 * @return func() interface{} Factory function that returns the Hub singleton
 */
func (p *PipelineServiceProvider) createHubFactory(application applicationInterfaces.ApplicationInterface) func() interface{} {
	return func() interface{} {
		return pipeline.NewHub(application)
	}
}
