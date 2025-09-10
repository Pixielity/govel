package main

import (
	"fmt"

	applicationInterfaces "govel/packages/application/interfaces"
	providerInterfaces "govel/packages/application/interfaces/providers"
)

// SimpleProvider demonstrates a provider that implements only DeferrableProvider
// without EventTriggeredProvider, showing that the When() method is now optional
type SimpleProvider struct {
	isRegistered bool
	isBooted     bool
}

// NewSimpleProvider creates a new SimpleProvider instance
func NewSimpleProvider() *SimpleProvider {
	return &SimpleProvider{}
}

// Register implements ServiceProviderInterface
func (p *SimpleProvider) Register(app applicationInterfaces.ApplicationInterface) error {
	if p.isRegistered {
		return nil
	}

	fmt.Println("SimpleProvider: Registering services...")
	// Register services here

	p.isRegistered = true
	return nil
}

// Boot implements ServiceProviderInterface
func (p *SimpleProvider) Boot(app applicationInterfaces.ApplicationInterface) error {
	if p.isBooted {
		return nil
	}

	fmt.Println("SimpleProvider: Booting...")
	// Boot services here

	p.isBooted = true
	return nil
}

// GetProvides implements ServiceProviderInterface
func (p *SimpleProvider) GetProvides() []string {
	return []string{
		"simple.service",
		"simple.helper",
	}
}

// IsDeferred implements ServiceProviderInterface
func (p *SimpleProvider) IsDeferred() bool {
	return true // This provider supports deferred loading
}

// Provides implements DeferrableProvider
func (p *SimpleProvider) Provides() []string {
	return p.GetProvides()
}

// NOTE: We intentionally DO NOT implement When() method to show that it's now optional

// Compile-time interface compliance checks
var _ providerInterfaces.ServiceProviderInterface = (*SimpleProvider)(nil)
var _ providerInterfaces.DeferrableProvider = (*SimpleProvider)(nil)

// This would cause a compilation error if we tried to assign it to EventTriggeredProvider:
// var _ providerInterfaces.EventTriggeredProvider = (*SimpleProvider)(nil) // This would fail

func main() {
	provider := NewSimpleProvider()

	fmt.Println("SimpleProvider created successfully!")
	fmt.Printf("Services provided: %v\n", provider.Provides())
	fmt.Printf("Is deferred: %v\n", provider.IsDeferred())

	// Demonstrate interface segregation - we can use it as DeferrableProvider
	var deferrableProvider providerInterfaces.DeferrableProvider = provider
	fmt.Printf("As DeferrableProvider - Services: %v\n", deferrableProvider.Provides())

	// But we cannot use it as EventTriggeredProvider without implementing When()
	// This would compile but we'd need to check at runtime:
	var serviceProvider providerInterfaces.ServiceProviderInterface = provider
	if eventProvider, ok := serviceProvider.(providerInterfaces.EventTriggeredProvider); ok {
		fmt.Printf("Events: %v\n", eventProvider.When())
	} else {
		fmt.Println("Provider does not implement EventTriggeredProvider (as expected)")
	}
}
