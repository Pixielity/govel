package application

import (
	"fmt"

	appBootstrappers "govel/application/bootstrappers"
	configBootstrappers "govel/config/bootstrappers"
)

// BootstrapHelper provides utility functions for application bootstrapping.
type BootstrapHelper struct {
	app *Application
}

// NewBootstrapHelper creates a new bootstrap helper for the given application.
//
// Parameters:
//   app: The application instance
//
// Returns:
//   *BootstrapHelper: A new bootstrap helper instance
//
// Example:
//   helper := NewBootstrapHelper(app)
func NewBootstrapHelper(app *Application) *BootstrapHelper {
	return &BootstrapHelper{
		app: app,
	}
}

// SetupDefaultBootstrappers configures the application with the standard GoVel bootstrappers.
// This method provides a convenient way to set up the typical bootstrap sequence for
// a GoVel application.
//
// Returns:
//   error: Any error that occurred during setup
//
// Example:
//   helper := NewBootstrapHelper(app)
//   err := helper.SetupDefaultBootstrappers()
//   if err != nil {
//       log.Fatalf("Failed to setup bootstrappers: %v", err)
//   }
func (h *BootstrapHelper) SetupDefaultBootstrappers() error {
	// Get the default bootstrappers from the application
	// This will return an empty slice by default, but can be overridden
	defaultBootstrappers := h.app.RegisterDefaultBootstrappers()

	// If no default bootstrappers are defined, set up the standard ones
	if len(defaultBootstrappers) == 0 {
		standardBootstrappers := []interface{}{
			configBootstrappers.NewLoadEnvironmentVariables(),
			configBootstrappers.NewLoadConfiguration(),
			appBootstrappers.NewRegisterFacades(),
		}
		h.app.SetBootstrappers(standardBootstrappers)
	} else {
		// Use the application-specific bootstrappers
		h.app.SetBootstrappers(defaultBootstrappers)
	}

	return nil
}

// BootstrapApplication performs the complete bootstrap process for the application.
// This includes setting up default bootstrappers and running the bootstrap process.
//
// Returns:
//   error: Any error that occurred during bootstrapping
//
// Example:
//   helper := NewBootstrapHelper(app)
//   err := helper.BootstrapApplication()
//   if err != nil {
//       log.Fatalf("Application bootstrap failed: %v", err)
//   }
func (h *BootstrapHelper) BootstrapApplication() error {
	// Setup default bootstrappers if none are configured
	if len(h.app.GetBootstrappers()) == 0 {
		if err := h.SetupDefaultBootstrappers(); err != nil {
			return fmt.Errorf("failed to setup default bootstrappers: %w", err)
		}
	}

	// Run the bootstrap process
	if err := h.app.Bootstrap(); err != nil {
		return fmt.Errorf("application bootstrap failed: %w", err)
	}

	return nil
}

// AddCustomBootstrapper adds a custom bootstrapper to the application.
// This allows for extending the bootstrap process with application-specific logic.
//
// Parameters:
//   bootstrapper: The custom bootstrapper instance to add
//
// Example:
//   helper := NewBootstrapHelper(app)
//   helper.AddCustomBootstrapper(&MyCustomBootstrapper{})
func (h *BootstrapHelper) AddCustomBootstrapper(bootstrapper interface{}) {
	h.app.AddBootstrapper(bootstrapper)
}

// GetBootstrapperCount returns the number of registered bootstrappers.
//
// Returns:
//   int: The number of bootstrappers
//
// Example:
//   helper := NewBootstrapHelper(app)
//   count := helper.GetBootstrapperCount()
//   fmt.Printf("Registered %d bootstrappers\n", count)
func (h *BootstrapHelper) GetBootstrapperCount() int {
	return len(h.app.GetBootstrappers())
}

// ValidateBootstrappers checks if the application is ready for bootstrapping.
//
// Returns:
//   error: Any validation error
//
// Example:
//   helper := NewBootstrapHelper(app)
//   if err := helper.ValidateBootstrappers(); err != nil {
//       log.Printf("Bootstrap validation failed: %v", err)
//   }
func (h *BootstrapHelper) ValidateBootstrappers() error {
	bootstrappers := h.app.GetBootstrappers()
	
	if len(bootstrappers) == 0 {
		return fmt.Errorf("no bootstrappers configured - call SetupDefaultBootstrappers() first")
	}

	// Additional validation logic can be added here
	// For example, checking for required bootstrapper types

	return nil
}

// CreateCustomApplicationWithBootstrappers creates a new application with custom bootstrappers.
// This is a factory function for creating applications with specific bootstrap configurations.
//
// Parameters:
//   customBootstrappers: List of custom bootstrapper instances
//
// Returns:
//   *Application: A new application with the specified bootstrappers
//   error: Any creation error
//
// Example:
//   customBootstrappers := []interface{}{
//       &MyEnvironmentBootstrapper{},
//       &MyConfigBootstrapper{},
//       &MyFacadeBootstrapper{},
//   }
//   app, err := CreateCustomApplicationWithBootstrappers(customBootstrappers)
func CreateCustomApplicationWithBootstrappers(customBootstrappers []interface{}) (*Application, error) {
	app := New()
	
	// Set the custom bootstrappers
	app.SetBootstrappers(customBootstrappers)
	
	// Create bootstrap helper and validate
	helper := NewBootstrapHelper(app)
	if err := helper.ValidateBootstrappers(); err != nil {
		return nil, fmt.Errorf("bootstrap validation failed: %w", err)
	}

	return app, nil
}

// BootstrapWithDefaults is a convenience function that creates a new application
// and bootstraps it with the default bootstrappers.
//
// Returns:
//   *Application: A bootstrapped application instance
//   error: Any error that occurred
//
// Example:
//   app, err := BootstrapWithDefaults()
//   if err != nil {
//       log.Fatalf("Failed to bootstrap application: %v", err)
//   }
//   // Application is now ready to use
func BootstrapWithDefaults() (*Application, error) {
	app := New()
	helper := NewBootstrapHelper(app)
	
	if err := helper.BootstrapApplication(); err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return app, nil
}