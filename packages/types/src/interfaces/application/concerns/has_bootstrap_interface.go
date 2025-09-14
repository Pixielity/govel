package interfaces

// HasBootstrapInterface provides application bootstrap management functionality.
// This interface defines the contract for managing bootstrap classes and
// performing application bootstrapping operations.
type HasBootstrapInterface interface {
	// Bootstrap bootstraps the application with all registered bootstrappers.
	// This method checks if the application has been bootstrapped and performs
	// the complete bootstrap process including deferred provider loading.
	Bootstrap() error

	// BootstrapWith bootstraps the application with specific bootstrappers.
	// This method allows selective bootstrapping with a custom set of bootstrappers.
	//
	// Parameters:
	//   bootstrappers: A slice of bootstrapper instances to use for bootstrapping
	//
	// Returns:
	//   error: Any error that occurred during bootstrapping
	BootstrapWith(bootstrappers []interface{}) error

	// BootstrapWithoutProviders bootstraps the application without booting service providers.
	// This method performs bootstrapping but excludes provider-related bootstrappers.
	BootstrapWithoutProviders() error

	// HasBeenBootstrapped returns whether the application has been bootstrapped.
	//
	// Returns:
	//   bool: true if the application has been bootstrapped, false otherwise
	HasBeenBootstrapped() bool

	// SetBootstrapped marks the application as bootstrapped or not bootstrapped.
	//
	// Parameters:
	//   bootstrapped: Whether the application should be marked as bootstrapped
	SetBootstrapped(bootstrapped bool)

	// GetBootstrappers returns the list of registered bootstrap classes.
	//
	// Returns:
	//   []interface{}: List of bootstrap class instances
	GetBootstrappers() []interface{}

	// SetBootstrappers sets the bootstrap classes for the application.
	//
	// Parameters:
	//   bootstrappers: A slice of bootstrap class instances
	SetBootstrappers(bootstrappers []interface{})

	// AddBootstrapper adds a single bootstrap class to the application.
	//
	// Parameters:
	//   bootstrapper: The bootstrap class instance to add
	AddBootstrapper(bootstrapper interface{})

	// RemoveBootstrapper removes a bootstrap class from the application.
	//
	// Parameters:
	//   bootstrapper: The bootstrap class instance to remove
	//
	// Returns:
	//   bool: true if the bootstrapper was found and removed, false otherwise
	RemoveBootstrapper(bootstrapper interface{}) bool

	// ClearBootstrappers removes all bootstrap classes.
	ClearBootstrappers()

	// RegisterDefaultBootstrappers returns the default bootstrap classes for the application.
	// This method should be implemented to return the standard bootstrappers needed
	// for application initialization in the correct order.
	//
	// Returns:
	//   []interface{}: List of default bootstrap class instances
	RegisterDefaultBootstrappers() []interface{}
}
