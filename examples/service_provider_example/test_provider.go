package main

import (
	"context"
	"fmt"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	providerInterfaces "govel/packages/application/interfaces/providers"
)

// TestProvider is a minimal test provider to isolate registration issues
type TestProvider struct {
	isRegistered bool
	isBooted     bool
	logger       interface{}
}

// NewTestProvider creates a new test provider
func NewTestProvider() *TestProvider {
	fmt.Println("[DEBUG] Creating new TestProvider")
	return &TestProvider{
		isRegistered: false,
		isBooted:     false,
	}
}

// Register registers test services with the application container
func (p *TestProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	fmt.Println("[DEBUG] TestProvider.Register() called")
	if p.isRegistered {
		fmt.Println("[DEBUG] TestProvider already registered, skipping")
		return nil
	}

	// Get container
	container := app.Container()
	logger := app.GetLogger()
	p.logger = logger

	// Register a simple test service
	if err := container.Bind("test.service", func() interface{} {
		logger.Info("Creating TestService instance")
		return "TestService"
	}); err != nil {
		return fmt.Errorf("failed to register test service: %w", err)
	}

	logger.Info("TestProvider registered successfully")
	p.isRegistered = true
	return nil
}

// Boot performs post-registration initialization
func (p *TestProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	logger := app.GetLogger()
	logger.Info("Booting TestProvider...")

	p.isBooted = true
	return nil
}

// GetProvides returns the services provided by this provider
func (p *TestProvider) GetProvides() []string {
	return []string{"test.service"}
}

// Terminate gracefully shuts down the test provider
func (p *TestProvider) Terminate(ctx context.Context, app applicationInterfaces.ApplicationInterface) error {
	logger := app.GetLogger()
	logger.Info("Terminating TestProvider...")

	select {
	case <-time.After(50 * time.Millisecond):
		logger.Info("TestProvider cleanup completed")
	case <-ctx.Done():
		logger.Warn("TestProvider cleanup cancelled due to timeout")
		return ctx.Err()
	}

	p.isBooted = false
	p.isRegistered = false

	logger.Info("TestProvider terminated successfully")
	return nil
}


// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*TestProvider)(nil)
var _ providerInterfaces.TerminatableProvider = (*TestProvider)(nil)
