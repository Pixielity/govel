package providers

import (
	"context"
	"fmt"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
	"service_provider_example/modules/user/repositories"
	"service_provider_example/modules/user/services"
)

func init() {
	fmt.Println("[DEBUG] UserServiceProvider module loaded")
}

// UserServiceProvider registers user-related services
type UserServiceProvider struct {
	isRegistered bool
	isBooted     bool
	logger       interface{} // Will hold logger instance
}

// NewUserServiceProvider creates a new user service provider
func NewUserServiceProvider() *UserServiceProvider {
	fmt.Println("[DEBUG] Creating new UserServiceProvider")
	return &UserServiceProvider{
		isRegistered: false,
		isBooted:     false,
	}
}

// Register registers user services with the application container
func (p *UserServiceProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] UserServiceProvider.Register() called")
	if p.isRegistered {
		fmt.Println("[DEBUG] UserServiceProvider already registered, skipping")
		return nil
	}

	// Get container
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	fmt.Println("[DEBUG] UserServiceProvider: Registering services")

	// Register repository
	if err := container.Singleton("user.repository", func() interface{} {
		logger.Info("Creating UserRepository instance")
		return repositories.NewUserRepository()
	}); err != nil {
		return fmt.Errorf("failed to register user repository: %w", err)
	}

	// Register service
	if err := container.Bind("user.service", func() interface{} {
		logger.Info("Creating UserService instance")
		// Create repository directly instead of resolving from container to avoid deadlock
		userRepo := repositories.NewUserRepository()
		return services.NewUserService(userRepo)
	}); err != nil {
		return fmt.Errorf("failed to register user service: %w", err)
	}

	logger.Info("UserServiceProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization
func (p *UserServiceProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting UserServiceProvider...")

	// Perform any boot-time initialization here
	// For example, we could pre-load some data or validate configuration

	// Verify our services are properly registered
	container := app.Container()
	if !container.IsBound("user.repository") {
		return fmt.Errorf("user repository not found in container")
	}
	if !container.IsBound("user.service") {
		return fmt.Errorf("user service not found in container")
	}

	logger.Info("UserServiceProvider booted successfully")
	p.isBooted = true
	return nil
}

// GetProvides returns the services provided by this provider (for deferred loading)
func (p *UserServiceProvider) GetProvides() []string {
	return []string{
		"user.repository",
		"user.service",
	}
}

// Terminate gracefully shuts down the user service provider
func (p *UserServiceProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating UserServiceProvider...")

	// Perform cleanup operations
	// For example, close database connections, save state, etc.

	// Simulate cleanup time
	select {
	case <-time.After(100 * time.Millisecond):
		logger.Info("UserServiceProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("UserServiceProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("UserServiceProvider terminated successfully")
	return nil
}

// TEST 3: Remove both deferred interfaces
// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*UserServiceProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*UserServiceProvider)(nil)
