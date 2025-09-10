package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	// Import our custom modules

	userProviders "service_provider_example/modules/user/providers"
	userServices "service_provider_example/modules/user/services"

	productProviders "service_provider_example/modules/product/providers"
	productServices "service_provider_example/modules/product/services"

	clientProviders "service_provider_example/modules/client/providers"
	clientServices "service_provider_example/modules/client/services"

	// Import GoVel framework components
	"govel/packages/application"
	"govel/packages/application/builders"
	configProviders "govel/packages/config/providers"
	containerProviders "govel/packages/container/providers"
	loggerProviders "govel/packages/logger/providers"
)

func main() {
	startTime := time.Now()
	fmt.Printf("🚀 Starting Service Provider Example at %s...\n", startTime.Format("15:04:05.000"))
	fmt.Println(strings.Repeat("=", 50))

	// Create service provider instances
	providers := []interface{}{
		// Core framework providers (essential services)
		configProviders.NewConfigServiceProvider(),
		loggerProviders.NewLoggerServiceProvider(),
		containerProviders.NewContainerServiceProvider(),

		// Custom module providers (ProductServiceProvider first to test ordering)
		productProviders.NewProductServiceProvider(),
		userProviders.NewUserServiceProvider(),

		// Client providers demonstrating different provider types
		clientProviders.NewClientServiceProvider(),                // 1. Standard/Eager provider
		clientProviders.NewClientDeferredServiceProvider(),       // 2. Deferred provider
		clientProviders.NewClientEventServiceProvider(),          // 3. Deferred + Event triggered provider

		// Test provider to isolate registration issues
		NewTestProvider(),
	}

	// Debug: Print provider count and types
	fmt.Printf("Registering %d providers:\n", len(providers))
	for i, provider := range providers {
		fmt.Printf("  %d: %T\n", i+1, provider)
	}

	// Build the application using the fluent builder with service providers
	app := builders.NewApp().
		WithName("Service Provider Example").
		WithVersion("1.0.0").
		WithEnvironment("development").
		WithDebug(true).
		WithServiceProviders(providers).
		Build()

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Boot the application and all service providers
	fmt.Println("\n📦 Booting application and service providers...")
	if err := app.BootProviders(ctx); err != nil {
		app.GetLogger().Error("Failed to boot service providers: %v", err)
		return
	}

	fmt.Println("✅ Application booted successfully!")
	fmt.Println("\n🔧 Testing Services...")
	fmt.Println(strings.Repeat("-", 30))

	// Test configuration service
	testConfigurationService(app)

	// Display container information
	displayContainerInfo(app)

	// Test test service
	testTestService(app)

	// Test client services (demonstrates different provider loading behaviors)
	testClientServices(app)

	// Test user service
	testUserService(app)

	// Test product service
	testProductService(app)

	// Demonstrate cross-module interaction
	fmt.Println("\n🔄 Cross-Module Operations...")
	fmt.Println(strings.Repeat("-", 30))
	testCrossModuleOperations(app)

	// Demonstrate graceful shutdown
	fmt.Println("\n🛑 Testing Graceful Shutdown...")
	fmt.Println(strings.Repeat("-", 30))
	testGracefulShutdown(app, ctx)

	fmt.Println("\n🎉 Service Provider Example completed successfully!")
}

// testConfigurationService demonstrates using the configuration service
func testConfigurationService(app *application.Application) {
	fmt.Println("📁 Configuration Service Test:")

	config := app.GetConfig()

	// Set some configuration values
	config.Set("app.name", "Service Provider Example App")
	config.Set("app.test_key", "test_value")
	config.Set("database.host", "localhost")
	config.Set("database.port", 5432)

	// Get configuration values with optional defaults
	appName := config.GetString("app.name")                           // No default - will return empty string if not found
	testKey := config.GetString("app.test_key", "default_test_value") // With default
	dbHost := config.GetString("database.host")                       // No default
	dbPort := config.GetInt("database.port")                          // No default - will return 0 if not found

	fmt.Printf("  • App Name: %s\n", appName)
	fmt.Printf("  • Test Key: %s\n", testKey)
	fmt.Printf("  • DB Host: %s\n", dbHost)
	fmt.Printf("  • DB Port: %d\n", dbPort)
	fmt.Printf("  • Has app.name: %v\n", config.HasKey("app.name"))
	fmt.Printf("  • Has non-existent: %v\n", config.HasKey("non.existent"))
}

// testTestService demonstrates using the test service
func testTestService(app *application.Application) {
	fmt.Println("\n🧪 Test Service Test:")

	container := app.GetContainer()

	// Get test service from container
	testServiceInterface, err := container.Make("test.service")
	if err != nil {
		fmt.Printf("  ❌ Failed to get test service: %v\n", err)
		return
	}

	fmt.Printf("  ✅ Retrieved test service: %v\n", testServiceInterface)
}

// testClientServices demonstrates the different client service provider types and their loading behaviors
func testClientServices(app *application.Application) {
	fmt.Println("\n👥 Client Services Test (Provider Type Demonstrations):")
	fmt.Println("This section demonstrates the different provider loading behaviors:")
	fmt.Println("  • Standard Provider: Loads immediately at startup")
	fmt.Println("  • Deferred Provider: Loads only when service is requested")
	fmt.Println("  • Event Provider: Loads when service is requested OR event is triggered")
	fmt.Println()

	container := app.GetContainer()

	// Test 1: Standard Client Service (should already be loaded)
	fmt.Println("📋 1. Testing Standard Client Service:")
	testStandardClientService(container)

	// Test 2: Deferred Client Service (will trigger loading)
	fmt.Println("\n⏰ 2. Testing Deferred Client Service:")
	fmt.Println("   [NOTE: This will trigger deferred loading]")
	testDeferredClientService(container)

	// Test 3: Event-Triggered Client Service (will trigger loading)
	fmt.Println("\n🎯 3. Testing Event-Triggered Client Service:")
	fmt.Println("   [NOTE: This will trigger event-based loading]")
	testEventClientService(container)
}

// testStandardClientService tests the standard (eager) client service
func testStandardClientService(container interface{}) {
	// Get the container interface - we'll need to cast it properly
	if containerService, ok := container.(interface {
		Make(string) (interface{}, error)
	}); ok {
		clientServiceInterface, err := containerService.Make("client.service")
		if err != nil {
			fmt.Printf("    ❌ Failed to get standard client service: %v\n", err)
			return
		}

		clientService, ok := clientServiceInterface.(clientServices.ClientServiceInterface)
		if !ok {
			fmt.Printf("    ❌ Invalid client service type: %T\n", clientServiceInterface)
			return
		}

		// Test basic operations
		clients, err := clientService.GetAllClients()
		if err != nil {
			fmt.Printf("    ❌ Failed to get clients: %v\n", err)
			return
		}

		fmt.Printf("    ✅ Standard service retrieved %d clients\n", len(clients))
		for _, client := range clients {
			fmt.Printf("      • %s (%s) - %s\n", client.GetDisplayName(), client.Email, client.Status)
		}

		// Test creating a new client
		newClient, err := clientService.CreateClient("Demo Corp", "demo@democorp.com", "Demo Corporation", "+1-555-DEMO")
		if err != nil {
			fmt.Printf("    ❌ Failed to create client: %v\n", err)
		} else {
			fmt.Printf("    ✅ Created client: %s\n", newClient.GetDisplayName())
		}
	} else {
		fmt.Printf("    ❌ Container does not support Make method\n")
	}
}

// testDeferredClientService tests the deferred client service
func testDeferredClientService(container interface{}) {
	if containerService, ok := container.(interface {
		Make(string) (interface{}, error)
	}); ok {
		clientServiceInterface, err := containerService.Make("client.deferred.service")
		if err != nil {
			fmt.Printf("    ❌ Failed to get deferred client service: %v\n", err)
			return
		}

		clientService, ok := clientServiceInterface.(clientServices.ClientServiceInterface)
		if !ok {
			fmt.Printf("    ❌ Invalid deferred client service type: %T\n", clientServiceInterface)
			return
		}

		// Test statistics operation
		stats, err := clientService.GetClientStatistics()
		if err != nil {
			fmt.Printf("    ❌ Failed to get statistics: %v\n", err)
			return
		}

		fmt.Printf("    ✅ Deferred service retrieved statistics:\n")
		for key, value := range stats {
			fmt.Printf("      • %s: %d\n", key, value)
		}

		// Test search functionality
		searchResults, err := clientService.SearchClients("corp")
		if err != nil {
			fmt.Printf("    ❌ Failed to search clients: %v\n", err)
		} else {
			fmt.Printf("    ✅ Found %d clients matching 'corp'\n", len(searchResults))
		}
	} else {
		fmt.Printf("    ❌ Container does not support Make method\n")
	}
}

// testEventClientService tests the event-triggered client service
func testEventClientService(container interface{}) {
	if containerService, ok := container.(interface {
		Make(string) (interface{}, error)
	}); ok {
		// Test the main event service
		clientServiceInterface, err := containerService.Make("client.event.service")
		if err != nil {
			fmt.Printf("    ❌ Failed to get event client service: %v\n", err)
			return
		}

		clientService, ok := clientServiceInterface.(clientServices.ClientServiceInterface)
		if !ok {
			fmt.Printf("    ❌ Invalid event client service type: %T\n", clientServiceInterface)
			return
		}

		// Test getting active clients count
		activeCount, err := clientService.GetActiveClientsCount()
		if err != nil {
			fmt.Printf("    ❌ Failed to get active clients count: %v\n", err)
		} else {
			fmt.Printf("    ✅ Event service found %d active clients\n", activeCount)
		}

		// Test the analytics service (also provided by event provider)
		analyticsServiceInterface, err := containerService.Make("client.analytics.service")
		if err != nil {
			fmt.Printf("    ❌ Failed to get analytics service: %v\n", err)
		} else {
			fmt.Printf("    ✅ Analytics service available: %T\n", analyticsServiceInterface)
		}
	} else {
		fmt.Printf("    ❌ Container does not support Make method\n")
	}
}

// testUserService demonstrates using the user service
func testUserService(app *application.Application) {
	fmt.Println("\n👤 User Service Test:")

	container := app.GetContainer()

	// Get user service from container
	userServiceInterface, err := container.Make("user.service")
	if err != nil {
		fmt.Printf("  ❌ Failed to get user service: %v\n", err)
		return
	}

	userService, ok := userServiceInterface.(userServices.UserServiceInterface)
	if !ok {
		fmt.Printf("  ❌ Invalid user service type: %T\n", userServiceInterface)
		return
	}

	// Test getting all users
	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Printf("  ❌ Failed to get users: %v\n", err)
		return
	}

	fmt.Printf("  📋 Total users: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("    • %d: %s (%s)\n", user.ID, user.Name, user.Email)
	}

	// Test creating a new user
	newUser, err := userService.CreateUser("Alice Johnson", "alice@example.com")
	if err != nil {
		fmt.Printf("  ❌ Failed to create user: %v\n", err)
	} else {
		fmt.Printf("  ✅ Created user: %d - %s (%s)\n", newUser.ID, newUser.Name, newUser.Email)
	}

	// Test searching users
	searchResults, err := userService.SearchUsers("john")
	if err != nil {
		fmt.Printf("  ❌ Failed to search users: %v\n", err)
	} else {
		fmt.Printf("  🔍 Search results for 'john': %d users\n", len(searchResults))
	}
}

// testProductService demonstrates using the product service
func testProductService(app *application.Application) {
	fmt.Println("\n🛍️  Product Service Test:")

	container := app.GetContainer()

	// Get product service from container
	productServiceInterface, err := container.Make("product.service")
	if err != nil {
		fmt.Printf("  ❌ Failed to get product service: %v\n", err)
		return
	}

	productService, ok := productServiceInterface.(productServices.ProductServiceInterface)
	if !ok {
		fmt.Printf("  ❌ Invalid product service type: %T\n", productServiceInterface)
		return
	}

	// Test getting all products
	products, err := productService.GetAllProducts()
	if err != nil {
		fmt.Printf("  ❌ Failed to get products: %v\n", err)
		return
	}

	fmt.Printf("  📋 Total products: %d\n", len(products))
	for _, product := range products {
		fmt.Printf("    • %d: %s ($%.2f) - %s\n", product.ID, product.Name, product.Price, product.Category)
	}

	// Test creating a new product
	newProduct, err := productService.CreateProduct("Wireless Mouse", "Ergonomic wireless mouse", "Electronics", 29.99)
	if err != nil {
		fmt.Printf("  ❌ Failed to create product: %v\n", err)
	} else {
		fmt.Printf("  ✅ Created product: %d - %s ($%.2f)\n", newProduct.ID, newProduct.Name, newProduct.Price)
	}

	// Test getting categories
	categories, err := productService.GetCategories()
	if err != nil {
		fmt.Printf("  ❌ Failed to get categories: %v\n", err)
	} else {
		fmt.Printf("  🏷️  Available categories: %v\n", categories)
	}

	// Test getting products by category
	electronicsProducts, err := productService.GetProductsByCategory("Electronics")
	if err != nil {
		fmt.Printf("  ❌ Failed to get electronics products: %v\n", err)
	} else {
		fmt.Printf("  🔌 Electronics products: %d\n", len(electronicsProducts))
	}
}

// testCrossModuleOperations demonstrates operations that span multiple modules
func testCrossModuleOperations(app *application.Application) {
	fmt.Println("🔗 Demonstrating cross-module operations:")

	container := app.GetContainer()
	logger := app.GetLogger()

	// Get both services
	userServiceInterface, _ := container.Make("user.service")
	productServiceInterface, _ := container.Make("product.service")

	if userServiceInterface == nil || productServiceInterface == nil {
		fmt.Println("  ❌ Could not get required services")
		return
	}

	userService := userServiceInterface.(userServices.UserServiceInterface)
	productService := productServiceInterface.(productServices.ProductServiceInterface)

	// Simulate a user purchasing products
	users, _ := userService.GetAllUsers()
	products, _ := productService.GetAllProducts()

	if len(users) > 0 && len(products) > 0 {
		user := users[0]
		product := products[0]

		logger.Info("User %s is interested in %s", user.Name, product.Name)
		fmt.Printf("  🛒 %s is viewing %s ($%.2f)\n", user.Name, product.Name, product.Price)

		// This would typically involve a separate order service
		fmt.Printf("  ✅ Cross-module operation completed successfully\n")
	}
}

// testGracefulShutdown demonstrates the graceful shutdown process
func testGracefulShutdown(app *application.Application, ctx context.Context) {
	fmt.Println("🔄 Initiating graceful shutdown...")

	// Create a shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Terminate all terminatable providers
	errors := app.TerminateProviders(shutdownCtx)

	if len(errors) > 0 {
		fmt.Printf("  ⚠️  Shutdown completed with %d errors:\n", len(errors))
		for i, err := range errors {
			fmt.Printf("    %d. %v\n", i+1, err)
		}
	} else {
		fmt.Printf("  ✅ All services terminated gracefully\n")
	}

	// Get final status
	loadedProviders := app.GetLoadedProviders()
	bootedProviders := app.GetBootedProviders()

	fmt.Printf("  📊 Final status:\n")
	fmt.Printf("    • Loaded providers: %d\n", len(loadedProviders))
	fmt.Printf("    • Booted providers: %d\n", len(bootedProviders))

	fmt.Printf("  🏁 Shutdown process completed\n")
}

// displayContainerInfo shows container bindings and statistics
func displayContainerInfo(app *application.Application) {
	fmt.Println("\n📦 Container Information:")
	fmt.Println(strings.Repeat("-", 30))

	container := app.GetContainer()
	
	// TODO: Check ContainerInterface for available methods
	// The GetBindings and GetStatistics methods don't exist on the interface
	fmt.Printf("📋 Container available (type: %T)\n", container)
	
	// TODO: Replace with actual available methods once we check the interface
	// bindings := container.GetBindings()
	// fmt.Printf("📋 Registered Services (%d total):\n", len(bindings))
	// for service := range bindings {
	//	fmt.Printf("  • %s\n", service)
	// }

	// stats := container.GetStatistics()
	// fmt.Printf("\n📊 Container Statistics (Raw):\n")
	// fmt.Printf("%+v\n", stats)
}
