package facades

import (
	testingInterfaces "govel/types/src/interfaces/testing"
	facade "govel/support/src"
)

// Testing provides a clean, static-like interface to the application's testing utilities service.
//
// This facade implements the facade pattern, providing global access to the testing
// service configured in the dependency injection container. It offers a Laravel-style
// API for testing utilities with automatic service resolution, mock management,
// test data factories, and comprehensive testing patterns.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved testing service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent testing operations across goroutines
//   - Supports mock services, test data factories, and assertion helpers
//   - Built-in database state management and service isolation
//
// Behavior:
//   - First call: Resolves testing service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if testing service cannot be resolved (fail-fast behavior)
//   - Automatically handles test isolation, cleanup, and state management
//
// Returns:
//   - TestingInterface: The application's testing service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "testing" service is not registered in the container
//   - If the resolved service doesn't implement TestingInterface
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
//   - Multiple goroutines can call Testing() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Testing operations are thread-safe with proper isolation
//
// Usage Examples:
//
//	// Basic service mocking
//	func TestUserService(t *testing.T) {
//	    // Create a mock user repository
//	    mockRepo := &MockUserRepository{
//	        users: make(map[int]*User),
//	    }
//
//	    // Swap the real repository with mock
//	    facades.Testing().Mock("user.repository", mockRepo)
//	    defer facades.Testing().RestoreMocks()
//
//	    // Test the user service
//	    userService := NewUserService()
//	    user := &User{Name: "John Doe", Email: "john@example.com"}
//
//	    err := userService.CreateUser(user)
//	    require.NoError(t, err)
//
//	    // Verify mock interactions
//	    assert.Equal(t, 1, mockRepo.CreateCallCount)
//	    assert.Equal(t, user, mockRepo.LastCreatedUser)
//	}
//
//	// Database testing with transactions
//	func TestUserCRUD(t *testing.T) {
//	    // Start a database transaction for test isolation
//	    facades.Testing().BeginTransaction()
//	    defer facades.Testing().RollbackTransaction()
//
//	    // Create test user
//	    user := &User{
//	        Name:  "Test User",
//	        Email: "test@example.com",
//	    }
//
//	    // Test creation
//	    err := facades.ORM().Create(user).Error
//	    require.NoError(t, err)
//	    assert.NotZero(t, user.ID)
//
//	    // Test retrieval
//	    var foundUser User
//	    err = facades.ORM().First(&foundUser, user.ID).Error
//	    require.NoError(t, err)
//	    assert.Equal(t, user.Email, foundUser.Email)
//
//	    // Test update
//	    foundUser.Name = "Updated Name"
//	    err = facades.ORM().Save(&foundUser).Error
//	    require.NoError(t, err)
//
//	    // Test deletion
//	    err = facades.ORM().Delete(&foundUser).Error
//	    require.NoError(t, err)
//
//	    // Transaction will be rolled back automatically
//	}
//
//	// Test data factories
//	type UserFactory struct{}
//
//	func (f *UserFactory) Make(overrides ...map[string]interface{}) *User {
//	    faker := facades.Testing().Faker()
//
//	    user := &User{
//	        Name:      faker.Person().Name(),
//	        Email:     faker.Internet().Email(),
//	        Password:  "password123",
//	        Active:    true,
//	        CreatedAt: time.Now(),
//	    }
//
//	    // Apply overrides
//	    for _, override := range overrides {
//	        if name, ok := override["name"].(string); ok {
//	            user.Name = name
//	        }
//	        if email, ok := override["email"].(string); ok {
//	            user.Email = email
//	        }
//	        if active, ok := override["active"].(bool); ok {
//	            user.Active = active
//	        }
//	    }
//
//	    return user
//	}
//
//	func (f *UserFactory) Create(overrides ...map[string]interface{}) *User {
//	    user := f.Make(overrides...)
//	    facades.ORM().Create(user)
//	    return user
//	}
//
//	// Using factories in tests
//	func TestUserFactory(t *testing.T) {
//	    facades.Testing().BeginTransaction()
//	    defer facades.Testing().RollbackTransaction()
//
//	    factory := &UserFactory{}
//
//	    // Create user with default values
//	    user1 := factory.Create()
//	    assert.NotEmpty(t, user1.Name)
//	    assert.NotEmpty(t, user1.Email)
//	    assert.True(t, user1.Active)
//
//	    // Create user with overrides
//	    user2 := factory.Create(map[string]interface{}{
//	        "name":   "Specific Name",
//	        "email":  "specific@example.com",
//	        "active": false,
//	    })
//	    assert.Equal(t, "Specific Name", user2.Name)
//	    assert.Equal(t, "specific@example.com", user2.Email)
//	    assert.False(t, user2.Active)
//	}
//
//	// HTTP testing
//	func TestUserAPI(t *testing.T) {
//	    // Create test server
//	    server := facades.Testing().CreateTestServer()
//	    defer server.Close()
//
//	    // Mock dependencies
//	    mockUserService := &MockUserService{}
//	    facades.Testing().Mock("user.service", mockUserService)
//	    defer facades.Testing().RestoreMocks()
//
//	    // Test GET request
//	    response := facades.Testing().Get(server.URL + "/api/users/1")
//	    assert.Equal(t, http.StatusOK, response.StatusCode)
//
//	    var user User
//	    err := response.JSON(&user)
//	    require.NoError(t, err)
//	    assert.Equal(t, "John Doe", user.Name)
//
//	    // Test POST request
//	    newUser := map[string]interface{}{
//	        "name":  "Jane Doe",
//	        "email": "jane@example.com",
//	    }
//
//	    response = facades.Testing().Post(server.URL+"/api/users", newUser)
//	    assert.Equal(t, http.StatusCreated, response.StatusCode)
//
//	    // Verify mock was called
//	    assert.Equal(t, 1, mockUserService.CreateCallCount)
//	}
//
// Advanced Testing Patterns:
//
//	// Custom assertion helpers
//	func TestCustomAssertions(t *testing.T) {
//	    user := &User{
//	        Name:  "John Doe",
//	        Email: "john@example.com",
//	    }
//
//	    // Custom assertion for user validation
//	    facades.Testing().AssertValidUser(t, user)
//
//	    // Custom assertion for database state
//	    facades.Testing().AssertUserExistsInDatabase(t, user.ID)
//
//	    // Custom assertion for API responses
//	    response := facades.Testing().Get("/api/users/1")
//	    facades.Testing().AssertSuccessfulResponse(t, response)
//	    facades.Testing().AssertJSONStructure(t, response, map[string]interface{}{
//	        "id":    "integer",
//	        "name":  "string",
//	        "email": "string",
//	    })
//	}
//
//	// Parallel testing with isolation
//	func TestParallelOperations(t *testing.T) {
//	    t.Parallel()
//
//	    // Each parallel test gets isolated environment
//	    testEnv := facades.Testing().CreateIsolatedEnvironment()
//	    defer testEnv.Cleanup()
//
//	    // Use isolated database connection
//	    db := testEnv.Database()
//
//	    // Run test with isolation
//	    user := &User{Name: "Test User", Email: "test@example.com"}
//	    err := db.Create(user).Error
//	    require.NoError(t, err)
//
//	    var count int64
//	    db.Model(&User{}).Count(&count)
//	    assert.Equal(t, int64(1), count)
//	}
//
//	// Testing with time manipulation
//	func TestTimeDependent(t *testing.T) {
//	    // Freeze time for consistent testing
//	    fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
//	    facades.Testing().FreezeTime(fixedTime)
//	    defer facades.Testing().UnfreezeTime()
//
//	    // Test time-dependent functionality
//	    user := &User{Name: "Test User", Email: "test@example.com"}
//	    err := facades.ORM().Create(user).Error
//	    require.NoError(t, err)
//
//	    // Verify timestamp
//	    assert.Equal(t, fixedTime.Unix(), user.CreatedAt.Unix())
//
//	    // Travel forward in time
//	    facades.Testing().TravelTo(fixedTime.Add(24 * time.Hour))
//
//	    // Test after time travel
//	    user.Name = "Updated Name"
//	    err = facades.ORM().Save(user).Error
//	    require.NoError(t, err)
//
//	    expectedUpdateTime := fixedTime.Add(24 * time.Hour)
//	    assert.Equal(t, expectedUpdateTime.Unix(), user.UpdatedAt.Unix())
//	}
//
// Mock Service Patterns:
//
//	// Comprehensive mock service
//	type MockEmailService struct {
//	    SentEmails    []Email
//	    SendCallCount int
//	    SendError     error
//	    mutex         sync.RWMutex
//	}
//
//	func (m *MockEmailService) Send(email Email) error {
//	    m.mutex.Lock()
//	    defer m.mutex.Unlock()
//
//	    m.SendCallCount++
//
//	    if m.SendError != nil {
//	        return m.SendError
//	    }
//
//	    m.SentEmails = append(m.SentEmails, email)
//	    return nil
//	}
//
//	func (m *MockEmailService) GetSentEmails() []Email {
//	    m.mutex.RLock()
//	    defer m.mutex.RUnlock()
//
//	    emails := make([]Email, len(m.SentEmails))
//	    copy(emails, m.SentEmails)
//	    return emails
//	}
//
//	func (m *MockEmailService) Reset() {
//	    m.mutex.Lock()
//	    defer m.mutex.Unlock()
//
//	    m.SentEmails = nil
//	    m.SendCallCount = 0
//	    m.SendError = nil
//	}
//
//	// Using mock in tests
//	func TestEmailNotification(t *testing.T) {
//	    mockEmail := &MockEmailService{}
//	    facades.Testing().Mock("email", mockEmail)
//	    defer facades.Testing().RestoreMocks()
//
//	    // Test email sending
//	    notificationService := NewNotificationService()
//	    err := notificationService.SendWelcomeEmail("user@example.com")
//	    require.NoError(t, err)
//
//	    // Verify email was sent
//	    sentEmails := mockEmail.GetSentEmails()
//	    assert.Len(t, sentEmails, 1)
//	    assert.Equal(t, "user@example.com", sentEmails[0].To)
//	    assert.Contains(t, sentEmails[0].Subject, "Welcome")
//
//	    // Test error handling
//	    mockEmail.Reset()
//	    mockEmail.SendError = errors.New("SMTP error")
//
//	    err = notificationService.SendWelcomeEmail("user2@example.com")
//	    assert.Error(t, err)
//	    assert.Contains(t, err.Error(), "SMTP error")
//	}
//
// File System Testing:
//
//	// Test file operations with temporary directories
//	func TestFileOperations(t *testing.T) {
//	    // Create temporary directory for testing
//	    tempDir := facades.Testing().CreateTempDir()
//	    defer facades.Testing().CleanupTempDir(tempDir)
//
//	    // Mock storage to use temp directory
//	    mockStorage := facades.Testing().MockFileSystem(tempDir)
//	    facades.Testing().Mock("storage", mockStorage)
//	    defer facades.Testing().RestoreMocks()
//
//	    // Test file operations
//	    fileService := NewFileService()
//
//	    content := []byte("test file content")
//	    err := fileService.SaveFile("test.txt", content)
//	    require.NoError(t, err)
//
//	    // Verify file exists
//	    assert.True(t, facades.Storage().Exists("test.txt"))
//
//	    // Verify content
//	    retrievedContent, err := fileService.ReadFile("test.txt")
//	    require.NoError(t, err)
//	    assert.Equal(t, content, retrievedContent)
//	}
//
// Cache Testing:
//
//	// Test cache operations with mock cache
//	func TestCacheOperations(t *testing.T) {
//	    mockCache := facades.Testing().MockCache()
//	    facades.Testing().Mock("cache", mockCache)
//	    defer facades.Testing().RestoreMocks()
//
//	    cacheService := NewCacheService()
//
//	    // Test cache set and get
//	    key := "test_key"
//	    value := "test_value"
//
//	    err := cacheService.Set(key, value, time.Hour)
//	    require.NoError(t, err)
//
//	    retrievedValue, err := cacheService.Get(key)
//	    require.NoError(t, err)
//	    assert.Equal(t, value, retrievedValue)
//
//	    // Test cache expiration
//	    facades.Testing().TravelTo(time.Now().Add(2 * time.Hour))
//
//	    _, err = cacheService.Get(key)
//	    assert.Error(t, err) // Should be expired
//
//	    // Verify mock interactions
//	    interactions := mockCache.GetInteractions()
//	    assert.Len(t, interactions, 3) // Set, Get, Get (expired)
//	}
//
// Queue Testing:
//
//	// Test job queue with sync processing
//	func TestJobQueue(t *testing.T) {
//	    // Use synchronous queue for testing
//	    syncQueue := facades.Testing().MockSyncQueue()
//	    facades.Testing().Mock("queue", syncQueue)
//	    defer facades.Testing().RestoreMocks()
//
//	    // Test job dispatch
//	    job := &SendEmailJob{
//	        To:      "user@example.com",
//	        Subject: "Test Email",
//	        Body:    "Test content",
//	    }
//
//	    err := facades.Queue().Dispatch(job)
//	    require.NoError(t, err)
//
//	    // Verify job was processed immediately
//	    processedJobs := syncQueue.GetProcessedJobs()
//	    assert.Len(t, processedJobs, 1)
//	    assert.Equal(t, job, processedJobs[0])
//
//	    // Test job failure
//	    failingJob := &FailingJob{}
//	    err = facades.Queue().Dispatch(failingJob)
//	    assert.Error(t, err)
//
//	    failedJobs := syncQueue.GetFailedJobs()
//	    assert.Len(t, failedJobs, 1)
//	}
//
// API Testing Helpers:
//
//	// Comprehensive API testing
//	func TestUserAPIEndpoints(t *testing.T) {
//	    server := facades.Testing().CreateTestServer()
//	    defer server.Close()
//
//	    client := facades.Testing().CreateAPIClient(server.URL)
//
//	    // Test authentication
//	    authResponse := client.Post("/auth/login", map[string]string{
//	        "email":    "user@example.com",
//	        "password": "password123",
//	    })
//
//	    facades.Testing().AssertStatus(t, authResponse, http.StatusOK)
//
//	    var authData map[string]interface{}
//	    err := authResponse.JSON(&authData)
//	    require.NoError(t, err)
//
//	    token, ok := authData["token"].(string)
//	    require.True(t, ok)
//
//	    // Set authorization header for subsequent requests
//	    client.SetHeader("Authorization", "Bearer "+token)
//
//	    // Test protected endpoint
//	    userResponse := client.Get("/api/user")
//	    facades.Testing().AssertStatus(t, userResponse, http.StatusOK)
//	    facades.Testing().AssertJSONStructure(t, userResponse, map[string]interface{}{
//	        "id":    "integer",
//	        "name":  "string",
//	        "email": "string",
//	    })
//
//	    // Test data creation
//	    createData := map[string]interface{}{
//	        "name":        "New Item",
//	        "description": "Test description",
//	        "price":       29.99,
//	    }
//
//	    createResponse := client.Post("/api/items", createData)
//	    facades.Testing().AssertStatus(t, createResponse, http.StatusCreated)
//
//	    var createdItem map[string]interface{}
//	    err = createResponse.JSON(&createdItem)
//	    require.NoError(t, err)
//
//	    itemID := int(createdItem["id"].(float64))
//
//	    // Test data retrieval
//	    getResponse := client.Get(fmt.Sprintf("/api/items/%d", itemID))
//	    facades.Testing().AssertStatus(t, getResponse, http.StatusOK)
//
//	    // Test data update
//	    updateData := map[string]interface{}{
//	        "name":  "Updated Item",
//	        "price": 39.99,
//	    }
//
//	    updateResponse := client.Put(fmt.Sprintf("/api/items/%d", itemID), updateData)
//	    facades.Testing().AssertStatus(t, updateResponse, http.StatusOK)
//
//	    // Test data deletion
//	    deleteResponse := client.Delete(fmt.Sprintf("/api/items/%d", itemID))
//	    facades.Testing().AssertStatus(t, deleteResponse, http.StatusNoContent)
//	}
//
// Integration Testing:
//
//	// Full application integration test
//	func TestUserRegistrationFlow(t *testing.T) {
//	    // Setup test environment
//	    testEnv := facades.Testing().CreateIntegrationEnvironment()
//	    defer testEnv.Cleanup()
//
//	    // Mock external services
//	    mockEmail := &MockEmailService{}
//	    facades.Testing().Mock("email", mockEmail)
//	    defer facades.Testing().RestoreMocks()
//
//	    client := testEnv.APIClient()
//
//	    // Test user registration
//	    userData := map[string]interface{}{
//	        "name":                  "John Doe",
//	        "email":                 "john@example.com",
//	        "password":              "password123",
//	        "password_confirmation": "password123",
//	    }
//
//	    response := client.Post("/auth/register", userData)
//	    facades.Testing().AssertStatus(t, response, http.StatusCreated)
//
//	    // Verify welcome email was sent
//	    sentEmails := mockEmail.GetSentEmails()
//	    assert.Len(t, sentEmails, 1)
//	    assert.Equal(t, "john@example.com", sentEmails[0].To)
//	    assert.Contains(t, sentEmails[0].Subject, "Welcome")
//
//	    // Verify user was created in database
//	    var user User
//	    err := testEnv.Database().Where("email = ?", "john@example.com").First(&user).Error
//	    require.NoError(t, err)
//	    assert.Equal(t, "John Doe", user.Name)
//
//	    // Test email verification
//	    verificationToken := user.VerificationToken
//	    verifyResponse := client.Get("/auth/verify/" + verificationToken)
//	    facades.Testing().AssertStatus(t, verifyResponse, http.StatusOK)
//
//	    // Verify user is now verified
//	    err = testEnv.Database().First(&user, user.ID).Error
//	    require.NoError(t, err)
//	    assert.True(t, user.EmailVerifiedAt != nil)
//
//	    // Test login with verified user
//	    loginResponse := client.Post("/auth/login", map[string]string{
//	        "email":    "john@example.com",
//	        "password": "password123",
//	    })
//	    facades.Testing().AssertStatus(t, loginResponse, http.StatusOK)
//	}
//
// Best Practices:
//   - Always clean up resources (transactions, mocks, temp files)
//   - Use test isolation to prevent test interference
//   - Mock external dependencies for reliable testing
//   - Use factories for consistent test data generation
//   - Test both success and failure scenarios
//   - Use descriptive test names and clear assertions
//   - Group related tests in test suites
//   - Use parallel testing when appropriate
//   - Test time-dependent functionality with time manipulation
//   - Verify mock interactions and call counts
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume testing service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	testing, err := facade.TryResolve[TestingInterface]("testing")
//	if err != nil {
//	    // Handle testing service unavailability gracefully
//	    log.Printf("Testing service unavailable: %v", err)
//	    return // Skip testing operations
//	}
//	testing.Mock("service", mockService)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestTestingBehavior(t *testing.T) {
//	    // Create a test testing service (meta-testing!)
//	    testTesting := &TestTesting{
//	        mocks: make(map[string]interface{}),
//	    }
//
//	    // Swap the real testing with test testing
//	    restore := support.SwapService("testing", testTesting)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Testing() returns testTesting
//	    facades.Testing().Mock("service", &MockService{})
//
//	    // Verify testing behavior
//	    assert.Len(t, testTesting.mocks, 1)
//	    assert.Contains(t, testTesting.mocks, "service")
//	}
//
// Container Configuration:
// Ensure the testing service is properly configured in your container:
//
//	// Example testing registration
//	container.Singleton("testing", func() interface{} {
//	    config := testing.Config{
//	        // Test database configuration
//	        TestDatabase: testing.DatabaseConfig{
//	            Driver:     "sqlite",
//	            Connection: ":memory:",
//	            Migrations: "./database/migrations",
//	        },
//
//	        // Mock service configuration
//	        MockConfig: testing.MockConfig{
//	            AutoRestore:    true,
//	            StrictMatching: true,
//	        },
//
//	        // Faker configuration
//	        Faker: testing.FakerConfig{
//	            Locale: "en_US",
//	            Seed:   12345, // For reproducible fake data
//	        },
//
//	        // HTTP testing configuration
//	        HTTPTesting: testing.HTTPConfig{
//	            BaseURL:        "http://localhost:8080",
//	            DefaultHeaders: map[string]string{
//	                "Content-Type": "application/json",
//	                "Accept":       "application/json",
//	            },
//	            Timeout: 30 * time.Second,
//	        },
//
//	        // Time manipulation
//	        TimeTesting: testing.TimeConfig{
//	            AllowTimeTravel: true,
//	            DefaultTZ:       "UTC",
//	        },
//
//	        // Cleanup configuration
//	        Cleanup: testing.CleanupConfig{
//	            AutoCleanup:     true,
//	            CleanupTimeout:  10 * time.Second,
//	            TempDirPrefix:   "govel_test_",
//	        },
//	    }
//
//	    testingService, err := testing.NewTestingService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create testing service: %v", err)
//	    }
//
//	    return testingService
//	})
func Testing() testingInterfaces.TestingInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "testing" service from the dependency injection container
	// - Performs type assertion to TestingInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[testingInterfaces.TestingInterface](testingInterfaces.TESTING_TOKEN)
}

// TestingWithError provides error-safe access to the testing service.
//
// This function offers the same functionality as Testing() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle testing service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Testing() but with error handling.
//
// Returns:
//   - TestingInterface: The resolved testing instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement TestingInterface
//
// Usage Examples:
//
//	// Basic error-safe testing operations
//	testing, err := facades.TestingWithError()
//	if err != nil {
//	    log.Printf("Testing service unavailable: %v", err)
//	    return // Skip testing operations
//	}
//	testing.Mock("service", mockService)
//
//	// Conditional testing setup
//	if testing, err := facades.TestingWithError(); err == nil {
//	    // Setup optional test mocks
//	    testing.Mock("optional_service", mockService)
//	}
func TestingWithError() (testingInterfaces.TestingInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "testing" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[testingInterfaces.TestingInterface](testingInterfaces.TESTING_TOKEN)
}
