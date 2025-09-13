package pipeline_test

import (
	"testing"

	applicationMocks "govel/application/mocks"
	"govel/new/pipeline/src/providers"
)

// TestNewPipelineServiceProvider tests the service provider constructor
func TestNewPipelineServiceProvider(t *testing.T) {
	t.Run("NewPipelineServiceProvider creates valid instance", func(t *testing.T) {
		provider := providers.NewPipelineServiceProvider()

		if provider == nil {
			t.Error("NewPipelineServiceProvider should not return nil")
		}

		// Verify provider provides expected services
		services := provider.Provides()
		if len(services) == 0 {
			t.Error("Provider should provide at least one service")
		}
	})
}

// TestPipelineServiceProviderRegister tests the Register method
func TestPipelineServiceProviderRegister(t *testing.T) {
	t.Run("Register with full ApplicationInterface mock", func(t *testing.T) {
		app := applicationMocks.NewMockApplication()
		provider := providers.NewPipelineServiceProvider()

		// Test that Register works with the mock application
		err := provider.Register(app)
		if err != nil {
			t.Errorf("Register failed: %v", err)
		}

		// Verify that services were registered
		services := provider.Provides()
		for _, service := range services {
			if !app.IsBound(service) {
				t.Errorf("Service '%s' was not bound", service)
			}
		}

		t.Logf("Successfully registered %d services", len(services))
	})
}

// TestPipelineServiceProviderProvides tests the Provides method
func TestPipelineServiceProviderProvides(t *testing.T) {
	provider := providers.NewPipelineServiceProvider()

	services := provider.Provides()

	expectedServices := []string{"pipeline", "pipeline.hub", "pipeline.hub.contract"}

	if len(services) != len(expectedServices) {
		t.Errorf("Expected %d services, got %d", len(expectedServices), len(services))
	}

	// Check that all expected services are provided
	serviceMap := make(map[string]bool)
	for _, service := range services {
		serviceMap[service] = true
	}

	for _, expected := range expectedServices {
		if !serviceMap[expected] {
			t.Errorf("Service '%s' not in provided services", expected)
		}
	}
}
