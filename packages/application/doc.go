// Package application provides the core application functionality for the GoVel framework.
//
// The application package is the central coordinator that brings together all the various
// interfaces and components that make up a GoVel application. It provides
// Laravel-inspired functionality while maintaining Go's design principles
// and performance characteristics.
//
// # Core Components
//
// The package is organized into several key areas:
//
//	┌─ application.go              # Main Application struct and core methods
//	├─ builders/           # Application builder patterns
//	├─ core/              # Core functionality (lifecycle, hooks, container)
//	├─ enums/             # Application enumerations
//	├─ interfaces/        # Interface definitions organized by domain
//	│  ├─ application/            # Main application interfaces
//	│  ├─ providers/      # Service provider interfaces
//	│  └─ traits/         # Trait interfaces (Has* pattern)
//	├─ maintenance/       # Maintenance mode system
//	├─ shutdown/          # Graceful shutdown system
//	├─ providers/         # Service provider infrastructure
//	└─ traits/           # Application trait implementations
//
// # Application Lifecycle
//
// The Application struct follows a predictable lifecycle:
//
//  1. Creation:     application := application.New()
//  2. Configuration: application.SetEnvironment(), application.UseConfigPath(), etc.
//  3. Registration: application.Register(serviceProvider)
//  4. Booting:      application.Boot(ctx)
//  5. Running:      // Application serves requests
//  6. Shutdown:     application.ShutdownWithSignals(ctx)
//
// # Service Providers
//
// GoVel applications are extensible through service providers that can:
//
//   - Register services in the container
//   - Boot services when the application starts
//   - Participate in graceful shutdown
//   - Provide deferred loading for performance
//
// # Maintenance Mode
//
// The package includes a comprehensive maintenance mode system:
//
//   - Put application into maintenance mode: application.Down(options)
//   - Check maintenance status: application.IsDown()
//   - Configure IP/path bypasses and custom messages
//   - Bring application back online: application.Up()
//
// # Graceful Shutdown
//
// Applications support graceful shutdown with:
//
//   - Signal handling (SIGINT, SIGTERM)
//   - Configurable timeout periods
//   - Service provider termination coordination
//   - Drain mode for completing active operations
//
// # Configuration and Environment
//
// The Application struct provides Laravel-style configuration management:
//
//   - Environment detection and configuration
//   - Multiple configuration sources (files, env vars, etc.)
//   - Directory path management and customization
//   - Locale and timezone configuration
//
// # Example Usage
//
// Basic application setup:
//
//	package main
//
//	import (
//	    "context"
//	    "log"
//	    "time"
//
//	    "govel/packages/application"
//	    "govel/packages/application/builders"
//	)
//
//	func main() {
//	    // Create and configure application
//	    application := builders.NewAppBuilder().
//	        SetEnvironment("production").
//	        SetName("My GoVel Application").
//	        SetVersion("1.0.0").
//	        Build()
//
//	    // Boot the application
//	    ctx := context.Background()
//	    if err := application.Boot(ctx); err != nil {
//	        log.Fatalf("Failed to boot application: %v", err)
//	    }
//
//	    // Setup graceful shutdown
//	    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	    defer cancel()
//
//	    // Your application logic here...
//
//	    // Graceful shutdown
//	    if err := application.ShutdownWithSignals(shutdownCtx); err != nil {
//	        log.Printf("Shutdown error: %v", err)
//	    }
//	}
//
// # Interface Segregation Principle
//
// The package follows ISP with granular interfaces:
//
//   - Directable: Directory path management
//   - Environmentable: Environment configuration
//   - Lifecycleable:   Lifecycle participation
//   - HasBootable:    Bootable components
//   - Shutdownable:    Shutdown participation
//   - Loggable:      Logging capabilities
//   - Configurable:      Configuration access
//
// This allows components to implement only the interfaces they need,
// promoting loose coupling and better testability.
package application
