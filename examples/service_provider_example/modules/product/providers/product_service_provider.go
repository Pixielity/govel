package providers

import (
	"context"
	"fmt"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	"service_provider_example/modules/product/repositories"
	"service_provider_example/modules/product/services"
)

// ProductServiceProvider registers product-related services
type ProductServiceProvider struct {
	isRegistered bool
	isBooted     bool
	logger       interface{} // Will hold logger instance
}

// NewProductServiceProvider creates a new product service provider
func NewProductServiceProvider() *ProductServiceProvider {
	fmt.Println("[DEBUG] Creating new ProductServiceProvider")
	return &ProductServiceProvider{
		isRegistered: false,
		isBooted:     false,
	}
}

// Register registers product services with the application container
func (p *ProductServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] ProductServiceProvider.Register() called")
	if p.isRegistered {
		fmt.Println("[DEBUG] ProductServiceProvider already registered, skipping")
		return nil
	}

	// Get container
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	// Register repository
	if err := container.Singleton("product.repository", func() interface{} {
		logger.Info("Creating ProductRepository instance")
		return repositories.NewProductRepository()
	}); err != nil {
		return fmt.Errorf("failed to register product repository: %w", err)
	}

	// Register service (deferred loading)
	if err := container.Bind("product.service", func() interface{} {
		logger.Info("Creating ProductService instance")
		// Create repository directly instead of resolving from container to avoid deadlock
		productRepo := repositories.NewProductRepository()
		return services.NewProductService(productRepo)
	}); err != nil {
		return fmt.Errorf("failed to register product service: %w", err)
	}

	logger.Info("ProductServiceProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization
func (p *ProductServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting ProductServiceProvider...")

	// Perform any boot-time initialization here
	// For example, we could pre-load some data or validate configuration

	// Verify our services are properly registered
	container := app.Container()
	if !container.IsBound("product.repository") {
		return fmt.Errorf("product repository not found in container")
	}
	if !container.IsBound("product.service") {
		return fmt.Errorf("product service not found in container")
	}

	logger.Info("ProductServiceProvider booted successfully")
	p.isBooted = true
	return nil
}

// GetProvides returns the services provided by this provider (for deferred loading)
func (p *ProductServiceProvider) GetProvides() []string {
	return []string{
		"product.repository",
		"product.service",
	}
}

// Terminate gracefully shuts down the product service provider
func (p *ProductServiceProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating ProductServiceProvider...")

	// Perform cleanup operations
	// For example, close database connections, save state, etc.

	// Simulate cleanup time
	select {
	case <-time.After(150 * time.Millisecond):
		logger.Info("ProductServiceProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("ProductServiceProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("ProductServiceProvider terminated successfully")
	return nil
}


// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*ProductServiceProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*ProductServiceProvider)(nil)
