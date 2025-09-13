package facades

import (
	seederInterfaces "govel/packages/types/src/interfaces/seeder"
	facade "govel/packages/support/src"
)

// Seeder provides a clean, static-like interface to the application's database seeding service.
//
// This facade implements the facade pattern, providing global access to the seeder
// service configured in the dependency injection container. It offers a Laravel-style
// API for database seeding with automatic service resolution, fixture management,
// test data generation, and seeding strategies.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved seeder service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent seeding operations across goroutines
//   - Supports multiple seeding strategies and data sources
//   - Built-in fixture management and relationship handling
//
// Behavior:
//   - First call: Resolves seeder service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if seeder service cannot be resolved (fail-fast behavior)
//   - Automatically handles data generation, relationship creation, and database operations
//
// Returns:
//   - SeederInterface: The application's seeder service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "seeder" service is not registered in the container
//   - If the resolved service doesn't implement SeederInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Seeder() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Seeding operations are thread-safe with proper database handling
//
// Usage Examples:
//
//	// Define a seeder struct
//	type UserSeeder struct{}
//
//	func (s *UserSeeder) Run() error {
//	    users := []User{
//	        {
//	            Name:     "John Doe",
//	            Email:    "john@example.com",
//	            Password: "password123",
//	            Active:   true,
//	        },
//	        {
//	            Name:     "Jane Smith",
//	            Email:    "jane@example.com",
//	            Password: "password456",
//	            Active:   true,
//	        },
//	    }
//
//	    for _, user := range users {
//	        if err := facades.ORM().Create(&user).Error; err != nil {
//	            return err
//	        }
//	    }
//
//	    return nil
//	}
//
//	// Register and run seeder
//	facades.Seeder().Register(&UserSeeder{})
//	err := facades.Seeder().Run("UserSeeder")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Run all registered seeders
//	err := facades.Seeder().RunAll()
//	if err != nil {
//	    log.Printf("Seeding failed: %v", err)
//	}
//
//	// Factory-based seeding
//	type UserFactory struct{}
//
//	func (f *UserFactory) Create(count int) []User {
//	    users := make([]User, count)
//
//	    for i := 0; i < count; i++ {
//	        users[i] = User{
//	            Name:     fmt.Sprintf("User %d", i+1),
//	            Email:    fmt.Sprintf("user%d@example.com", i+1),
//	            Password: "password123",
//	            Active:   true,
//	        }
//	    }
//
//	    return users
//	}
//
//	// Use factory to generate test data
//	factory := &UserFactory{}
//	users := factory.Create(50)
//
//	for _, user := range users {
//	    facades.ORM().Create(&user)
//	}
//
//	// Faker-powered seeding
//	type ProductSeeder struct{}
//
//	func (s *ProductSeeder) Run() error {
//	    faker := facades.Seeder().Faker()
//
//	    for i := 0; i < 100; i++ {
//	        product := Product{
//	            Name:        faker.Commerce().ProductName(),
//	            Description: faker.Lorem().Paragraph(3),
//	            Price:       faker.Commerce().Price(),
//	            SKU:         faker.Lorem().Text(10),
//	            InStock:     faker.Boolean(),
//	            CreatedAt:   faker.Date().Between(time.Now().AddDate(-1, 0, 0), time.Now()),
//	        }
//
//	        if err := facades.ORM().Create(&product).Error; err != nil {
//	            return err
//	        }
//	    }
//
//	    return nil
//	}
//
//	// Relationship seeding
//	type BlogSeeder struct{}
//
//	func (s *BlogSeeder) Run() error {
//	    // First, get existing users or create them
//	    var users []User
//	    facades.ORM().Find(&users)
//
//	    if len(users) == 0 {
//	        // Create users first
//	        userSeeder := &UserSeeder{}
//	        if err := userSeeder.Run(); err != nil {
//	            return err
//	        }
//	        facades.ORM().Find(&users)
//	    }
//
//	    faker := facades.Seeder().Faker()
//
//	    // Create blog posts for each user
//	    for _, user := range users {
//	        for i := 0; i < 5; i++ {
//	            post := Post{
//	                UserID:    user.ID,
//	                Title:     faker.Lorem().Sentence(6),
//	                Content:   faker.Lorem().Paragraphs(5),
//	                Published: faker.Boolean(),
//	                CreatedAt: faker.Date().Between(time.Now().AddDate(0, -6, 0), time.Now()),
//	            }
//
//	            if err := facades.ORM().Create(&post).Error; err != nil {
//	                return err
//	            }
//	        }
//	    }
//
//	    return nil
//	}
//
//	// Conditional seeding based on environment
//	type DevelopmentSeeder struct{}
//
//	func (s *DevelopmentSeeder) Run() error {
//	    if facades.App().Environment() != "development" {
//	        return nil // Skip in non-development environments
//	    }
//
//	    // Create development-specific data
//	    adminUser := User{
//	        Name:     "Admin User",
//	        Email:    "admin@localhost",
//	        Password: "admin123",
//	        IsAdmin:  true,
//	        Active:   true,
//	    }
//
//	    return facades.ORM().Create(&adminUser).Error
//	}
//
// Advanced Seeding Patterns:
//
//	// Batch seeding for performance
//	type BatchUserSeeder struct{}
//
//	func (s *BatchUserSeeder) Run() error {
//	    batchSize := 1000
//	    totalUsers := 10000
//
//	    faker := facades.Seeder().Faker()
//
//	    for batch := 0; batch < totalUsers/batchSize; batch++ {
//	        users := make([]User, batchSize)
//
//	        for i := 0; i < batchSize; i++ {
//	            users[i] = User{
//	                Name:      faker.Person().Name(),
//	                Email:     faker.Internet().Email(),
//	                Password:  "password123",
//	                Active:    true,
//	                CreatedAt: faker.Date().Between(time.Now().AddDate(-2, 0, 0), time.Now()),
//	            }
//	        }
//
//	        // Batch insert for better performance
//	        if err := facades.ORM().CreateInBatches(users, batchSize).Error; err != nil {
//	            return err
//	        }
//
//	        fmt.Printf("Seeded batch %d/%d\n", batch+1, totalUsers/batchSize)
//	    }
//
//	    return nil
//	}
//
//	// JSON fixture loading
//	type FixtureSeeder struct {
//	    FixturePath string
//	}
//
//	func (s *FixtureSeeder) Run() error {
//	    // Load fixture data from JSON file
//	    data, err := facades.Seeder().LoadFixture(s.FixturePath)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Parse JSON data into structs
//	    var users []User
//	    if err := json.Unmarshal(data, &users); err != nil {
//	        return err
//	    }
//
//	    // Insert into database
//	    for _, user := range users {
//	        if err := facades.ORM().Create(&user).Error; err != nil {
//	            return err
//	        }
//	    }
//
//	    return nil
//	}
//
//	// CSV data seeding
//	type CSVSeeder struct {
//	    CSVPath string
//	}
//
//	func (s *CSVSeeder) Run() error {
//	    records, err := facades.Seeder().LoadCSV(s.CSVPath)
//	    if err != nil {
//	        return err
//	    }
//
//	    for _, record := range records {
//	        user := User{
//	            Name:     record["name"],
//	            Email:    record["email"],
//	            Password: "password123",
//	            Active:   record["active"] == "true",
//	        }
//
//	        if err := facades.ORM().Create(&user).Error; err != nil {
//	            facades.Log().Warning("Failed to seed user from CSV", map[string]interface{}{
//	                "email": record["email"],
//	                "error": err.Error(),
//	            })
//	            continue
//	        }
//	    }
//
//	    return nil
//	}
//
//	// Seeder with dependencies
//	type DatabaseSeeder struct{}
//
//	func (s *DatabaseSeeder) Run() error {
//	    // Run seeders in dependency order
//	    seeders := []string{
//	        "UserSeeder",
//	        "CategorySeeder",
//	        "ProductSeeder",
//	        "OrderSeeder",
//	        "ReviewSeeder",
//	    }
//
//	    for _, seederName := range seeders {
//	        facades.Log().Info("Running seeder", map[string]interface{}{
//	            "seeder": seederName,
//	        })
//
//	        if err := facades.Seeder().Run(seederName); err != nil {
//	            return fmt.Errorf("seeder %s failed: %w", seederName, err)
//	        }
//	    }
//
//	    return nil
//	}
//
// Seeder Management:
//
//	// List available seeders
//	seeders := facades.Seeder().List()
//	for _, seeder := range seeders {
//	    fmt.Printf("Available seeder: %s\n", seeder)
//	}
//
//	// Check if seeder has been run
//	if facades.Seeder().HasRun("UserSeeder") {
//	    fmt.Println("UserSeeder has already been executed")
//	}
//
//	// Reset seeding history
//	facades.Seeder().Reset()
//
//	// Rollback specific seeder
//	err := facades.Seeder().Rollback("UserSeeder")
//	if err != nil {
//	    log.Printf("Rollback failed: %v", err)
//	}
//
// Testing Support:
//
//	// Database state management for tests
//	func TestUserService(t *testing.T) {
//	    // Seed test data
//	    testUsers := []User{
//	        {Name: "Test User 1", Email: "test1@example.com"},
//	        {Name: "Test User 2", Email: "test2@example.com"},
//	    }
//
//	    facades.Seeder().SeedTestData("users", testUsers)
//
//	    // Run test
//	    userService := NewUserService()
//	    users, err := userService.GetAllUsers()
//	    require.NoError(t, err)
//	    assert.Len(t, users, 2)
//
//	    // Cleanup after test
//	    facades.Seeder().CleanupTestData("users")
//	}
//
//	// Factory functions for testing
//	func TestFactories(t *testing.T) {
//	    // Create single instance
//	    user := facades.Seeder().Factory("User", map[string]interface{}{
//	        "email": "specific@example.com",
//	        "name":  "Specific Name",
//	    })
//
//	    // Create multiple instances
//	    users := facades.Seeder().FactoryMultiple("User", 5, map[string]interface{}{
//	        "active": true,
//	    })
//
//	    assert.Len(t, users, 5)
//	}
//
// Best Practices:
//   - Keep seeders idempotent (safe to run multiple times)
//   - Use transactions for complex seeding operations
//   - Handle seeding failures gracefully with proper error messages
//   - Use faker for realistic test data generation
//   - Organize seeders by logical groupings (users, products, etc.)
//   - Consider performance when seeding large datasets
//   - Use fixtures for consistent, reproducible data
//   - Implement proper cleanup for test environments
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume seeder service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	seeder, err := facade.TryResolve[SeederInterface]("seeder")
//	if err != nil {
//	    // Handle seeder service unavailability gracefully
//	    log.Printf("Seeder service unavailable: %v", err)
//	    return // Skip seeding
//	}
//	err = seeder.Run("UserSeeder")
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestSeeding(t *testing.T) {
//	    // Create a test seeder that tracks operations
//	    testSeeder := &TestSeeder{
//	        executedSeeders: []string{},
//	    }
//
//	    // Swap the real seeder with test seeder
//	    restore := support.SwapService("seeder", testSeeder)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Seeder() returns testSeeder
//	    err := facades.Seeder().Run("UserSeeder")
//	    require.NoError(t, err)
//
//	    // Verify seeding behavior
//	    executed := testSeeder.GetExecutedSeeders()
//	    assert.Contains(t, executed, "UserSeeder")
//	}
//
// Container Configuration:
// Ensure the seeder service is properly configured in your container:
//
//	// Example seeder registration
//	container.Singleton("seeder", func() interface{} {
//	    config := seeder.Config{
//	        // Database connection
//	        DB: facades.ORM(),
//
//	        // Faker configuration
//	        Faker: seeder.FakerConfig{
//	            Locale: "en_US",
//	            Seed:   12345, // For reproducible fake data
//	        },
//
//	        // Fixture paths
//	        FixturePaths: []string{
//	            "./database/fixtures",
//	            "./tests/fixtures",
//	        },
//
//	        // Seeder tracking
//	        TrackingEnabled: true,
//	        TrackingTable:   "seeder_runs",
//
//	        // Batch configuration
//	        DefaultBatchSize: 1000,
//
//	        // Performance options
//	        UseTransactions: true,
//	        LogProgress:     true,
//	    }
//
//	    seederService, err := seeder.NewSeederService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create seeder service: %v", err)
//	    }
//
//	    return seederService
//	})
func Seeder() seederInterfaces.SeederInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "seeder" service from the dependency injection container
	// - Performs type assertion to SeederInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[seederInterfaces.SeederInterface](seederInterfaces.SEEDER_TOKEN)
}

// SeederWithError provides error-safe access to the seeder service.
//
// This function offers the same functionality as Seeder() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle seeder service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Seeder() but with error handling.
//
// Returns:
//   - SeederInterface: The resolved seeder instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement SeederInterface
//
// Usage Examples:
//
//	// Basic error-safe seeding
//	seeder, err := facades.SeederWithError()
//	if err != nil {
//	    log.Printf("Seeder service unavailable: %v", err)
//	    return // Skip seeding operations
//	}
//	err = seeder.Run("UserSeeder")
//
//	// Conditional seeding
//	if seeder, err := facades.SeederWithError(); err == nil {
//	    // Run optional development seeders
//	    if facades.App().IsLocal() {
//	        seeder.Run("DevelopmentSeeder")
//	    }
//	}
func SeederWithError() (seederInterfaces.SeederInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "seeder" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[seederInterfaces.SeederInterface](seederInterfaces.SEEDER_TOKEN)
}
