package providers

import (
	"path/filepath"

	"govel/application"
	"govel/package_manager/core"
	"govel/package_manager/interfaces"
	"govel/package_manager/parsers"
	"govel/package_manager/registry"
	"govel/package_manager/services"
)

// PackageManagerServiceProvider provides package manager services to the GoVel container
type PackageManagerServiceProvider struct{}

// Register registers the package manager services in the container
func (p *PackageManagerServiceProvider) Register(app application.ApplicationInterface) error {
	container := app.Container()

	// Register parser
	container.Singleton("package_manager.parser", func() interfaces.ParserInterface {
		return parsers.NewModuleParser()
	})

	// Register command executor
	container.Singleton("package_manager.executor", func() interfaces.ExecutorInterface {
		return services.NewCommandExecutor()
	})

	// Register state manager
	container.Singleton("package_manager.state_manager", func() interfaces.StateManagerInterface {
		stateDir := filepath.Join(app.BasePath(), ".govel")
		return services.NewStateManager(stateDir)
	})

	// Register package registry
	container.Singleton("package_manager.registry", func() interfaces.RegistryInterface {
		parser := container.Resolve("package_manager.parser").(interfaces.ParserInterface)
		stateManager := container.Resolve("package_manager.state_manager").(interfaces.StateManagerInterface)
		return registry.NewPackageRegistry(parser, stateManager)
	})

	// Register dependency resolver
	container.Singleton("package_manager.dependency_resolver", func() interfaces.DependencyResolverInterface {
		registry := container.Resolve("package_manager.registry").(interfaces.RegistryInterface)
		return services.NewDependencyResolver(registry)
	})

	// Register main package manager
	container.Singleton("package_manager", func() interfaces.PackageManagerInterface {
		rootPath := app.BasePath()
		parser := container.Resolve("package_manager.parser").(interfaces.ParserInterface)
		executor := container.Resolve("package_manager.executor").(interfaces.ExecutorInterface)
		stateManager := container.Resolve("package_manager.state_manager").(interfaces.StateManagerInterface)
		registry := container.Resolve("package_manager.registry").(interfaces.RegistryInterface)
		dependencyResolver := container.Resolve("package_manager.dependency_resolver").(interfaces.DependencyResolverInterface)

		return core.NewPackageManager(
			rootPath,
			registry,
			parser,
			executor,
			dependencyResolver,
			stateManager,
		)
	})

	return nil
}

// Boot initializes the package manager services
func (p *PackageManagerServiceProvider) Boot(app application.ApplicationInterface) error {
	// Initialize package registry by scanning for packages
	pm := app.Container().Resolve("package_manager").(interfaces.PackageManagerInterface)

	// Perform initial scan in the background
	go func() {
		if err := pm.Refresh(app.Context()); err != nil {
			app.Logger().Error("Failed to initialize package registry", "error", err)
		}
	}()

	return nil
}

// Provides returns the services provided by this service provider
func (p *PackageManagerServiceProvider) Provides() []string {
	return []string{
		"package_manager",
		"package_manager.parser",
		"package_manager.executor",
		"package_manager.state_manager",
		"package_manager.registry",
		"package_manager.dependency_resolver",
	}
}
