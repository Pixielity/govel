package facades

import (
	eventInterfaces "govel/types/src/interfaces/event"
	facade "govel/support/src"
)

// Event provides a clean, static-like interface to the application's event dispatching service.
//
// This facade implements the facade pattern, providing global access to the event
// service configured in the dependency injection container. It offers a Laravel-style
// API for event-driven architecture with automatic service resolution, listener management,
// asynchronous event processing, and robust error handling.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved event service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent event operations across goroutines
//   - Supports synchronous and asynchronous event dispatching
//   - Built-in event listener registration and subscription management
//
// Behavior:
//   - First call: Resolves event service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if event service cannot be resolved (fail-fast behavior)
//   - Automatically handles event routing, listener execution, and error propagation
//
// Returns:
//   - EventInterface: The application's event service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "event" service is not registered in the container
//   - If the resolved service doesn't implement EventInterface
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
//   - Multiple goroutines can call Event() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Event dispatching and listener execution are thread-safe
//
// Usage Examples:
//
//	// Basic event dispatching
//	event := UserRegisteredEvent{
//	    UserID:    123,
//	    Email:     "user@example.com",
//	    Name:      "John Doe",
//	    Timestamp: time.Now(),
//	}
//
//	facades.Event().Dispatch("user.registered", event)
//
//	// Dispatch with context
//	ctx := context.WithTimeout(context.Background(), 5*time.Second)
//	err := facades.Event().DispatchWithContext(ctx, "user.login", UserLoginEvent{
//	    UserID:    123,
//	    IPAddress: "192.168.1.100",
//	    UserAgent: "Mozilla/5.0...",
//	    LoginTime: time.Now(),
//	})
//	if err != nil {
//	    log.Printf("Event dispatch failed: %v", err)
//	}
//
//	// Asynchronous event dispatching
//	facades.Event().DispatchAsync("order.placed", OrderPlacedEvent{
//	    OrderID:    "ORD-123",
//	    CustomerID: 456,
//	    Amount:     99.99,
//	    Items:      []string{"item1", "item2"},
//	})
//
//	// Register event listeners
//	facades.Event().Listen("user.registered", func(event interface{}) error {
//	    userEvent := event.(UserRegisteredEvent)
//
//	    // Send welcome email
//	    return facades.Mail().Send("welcome", userEvent.Email, map[string]interface{}{
//	        "name": userEvent.Name,
//	    })
//	})
//
//	// Register multiple listeners for same event
//	facades.Event().Listen("user.registered", func(event interface{}) error {
//	    userEvent := event.(UserRegisteredEvent)
//
//	    // Create user profile
//	    return facades.DB().Exec(
//	        "INSERT INTO user_profiles (user_id, created_at) VALUES (?, ?)",
//	        userEvent.UserID, time.Now(),
//	    )
//	})
//
//	facades.Event().Listen("user.registered", func(event interface{}) error {
//	    userEvent := event.(UserRegisteredEvent)
//
//	    // Log registration
//	    facades.Log().Info("New user registered", map[string]interface{}{
//	        "user_id": userEvent.UserID,
//	        "email":   userEvent.Email,
//	    })
//	    return nil
//	})
//
//	// Subscribe to events with priority
//	facades.Event().Subscribe(&UserEventSubscriber{})
//
//	// Event subscriber example
//	type UserEventSubscriber struct{}
//
//	func (s *UserEventSubscriber) Subscribe(events EventInterface) {
//	    // High priority listener (executes first)
//	    events.ListenWithPriority("user.registered", s.HandleUserRegistration, 100)
//
//	    // Normal priority listeners
//	    events.Listen("user.login", s.HandleUserLogin)
//	    events.Listen("user.logout", s.HandleUserLogout)
//	    events.Listen("user.password_changed", s.HandlePasswordChange)
//	}
//
//	func (s *UserEventSubscriber) HandleUserRegistration(event interface{}) error {
//	    userEvent := event.(UserRegisteredEvent)
//
//	    // Critical registration processing
//	    return s.setupUserAccount(userEvent.UserID)
//	}
//
//	// Conditional event dispatching
//	if facades.Event().HasListeners("payment.processed") {
//	    facades.Event().Dispatch("payment.processed", PaymentProcessedEvent{
//	        PaymentID: "PAY-123",
//	        Amount:    49.99,
//	        Status:    "completed",
//	    })
//	}
//
//	// Event forwarder pattern
//	facades.Event().Listen("order.placed", func(event interface{}) error {
//	    orderEvent := event.(OrderPlacedEvent)
//
//	    // Forward to inventory system
//	    return facades.Event().DispatchAsync("inventory.reserve", InventoryReserveEvent{
//	        OrderID: orderEvent.OrderID,
//	        Items:   orderEvent.Items,
//	    })
//	})
//
//	// Event with timeout handling
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	err := facades.Event().DispatchWithContext(ctx, "email.send", EmailEvent{
//	    To:      "user@example.com",
//	    Subject: "Important Notification",
//	    Body:    "This is an important message",
//	})
//	if err != nil {
//	    if errors.Is(err, context.DeadlineExceeded) {
//	        log.Printf("Email event timed out")
//	    } else {
//	        log.Printf("Email event failed: %v", err)
//	    }
//	}
//
// Advanced Event Patterns:
//
//	// Domain event pattern
//	type User struct {
//	    ID     int
//	    Email  string
//	    events []interface{}
//	}
//
//	func (u *User) Register(email string) {
//	    u.Email = email
//	    u.AddEvent(UserRegisteredEvent{
//	        UserID:    u.ID,
//	        Email:     email,
//	        Timestamp: time.Now(),
//	    })
//	}
//
//	func (u *User) AddEvent(event interface{}) {
//	    u.events = append(u.events, event)
//	}
//
//	func (u *User) FlushEvents() {
//	    for _, event := range u.events {
//	        facades.Event().Dispatch(getEventName(event), event)
//	    }
//	    u.events = nil
//	}
//
//	// Event sourcing pattern
//	func ReplayUserEvents(userID int) (*User, error) {
//	    events, err := facades.Event().GetHistory("user", userID)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    user := &User{ID: userID}
//	    for _, event := range events {
//	        user.ApplyEvent(event)
//	    }
//
//	    return user, nil
//	}
//
//	// Event middleware pattern
//	facades.Event().Use(func(eventName string, event interface{}, next func() error) error {
//	    // Log all events
//	    start := time.Now()
//	    facades.Log().Debug("Event dispatching", map[string]interface{}{
//	        "event": eventName,
//	    })
//
//	    err := next() // Execute listeners
//
//	    duration := time.Since(start)
//	    if err != nil {
//	        facades.Log().Error("Event failed", map[string]interface{}{
//	            "event":    eventName,
//	            "duration": duration,
//	            "error":    err.Error(),
//	        })
//	    } else {
//	        facades.Log().Debug("Event completed", map[string]interface{}{
//	            "event":    eventName,
//	            "duration": duration,
//	        })
//	    }
//
//	    return err
//	})
//
//	// Event batching for performance
//	batch := facades.Event().NewBatch()
//	batch.Add("user.login", UserLoginEvent{UserID: 1})
//	batch.Add("user.login", UserLoginEvent{UserID: 2})
//	batch.Add("user.login", UserLoginEvent{UserID: 3})
//
//	// Dispatch all events in batch
//	err := batch.Dispatch()
//	if err != nil {
//	    log.Printf("Batch dispatch failed: %v", err)
//	}
//
// Best Practices:
//   - Use descriptive event names with dot notation ("user.registered", "order.placed")
//   - Design events to be immutable and self-contained
//   - Keep event listeners lightweight and focused on single responsibilities
//   - Use asynchronous dispatch for non-critical operations
//   - Implement proper error handling and retry logic for failed listeners
//   - Use event subscribers to organize related listeners
//   - Consider event versioning for backward compatibility
//   - Monitor event processing performance and failure rates
//
// Error Handling Strategies:
//   - Graceful degradation: Continue processing other listeners if one fails
//   - Retry mechanisms: Automatically retry failed listeners
//   - Dead letter queues: Store failed events for later processing
//   - Circuit breakers: Temporarily disable failing listeners
//   - Monitoring: Track event processing metrics and alerts
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume event service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	eventService, err := facade.Resolve[EventInterface]("event")
//	if err != nil {
//	    // Handle event service unavailability gracefully
//	    log.Printf("Event service unavailable, skipping event: %v", err)
//	    return
//	}
//	eventService.Dispatch("user.action", event)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserRegistration(t *testing.T) {
//	    // Create a test event service that captures dispatched events
//	    testEvents := &TestEventService{
//	        dispatchedEvents: make(map[string][]interface{}),
//	    }
//
//	    // Swap the real event service with test service
//	    restore := support.SwapService("event", testEvents)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Event() returns testEvents
//	    userService := NewUserService()
//	    user, err := userService.Register("test@example.com", "password")
//	    require.NoError(t, err)
//
//	    // Verify events were dispatched
//	    events := testEvents.GetEvents("user.registered")
//	    assert.Len(t, events, 1)
//
//	    registeredEvent := events[0].(UserRegisteredEvent)
//	    assert.Equal(t, user.ID, registeredEvent.UserID)
//	    assert.Equal(t, "test@example.com", registeredEvent.Email)
//	}
//
// Container Configuration:
// Ensure the event service is properly configured in your container:
//
//	// Example event registration
//	container.Singleton("event", func() interface{} {
//	    config := event.Config{
//	        // Event processing configuration
//	        AsyncWorkers:    10,                    // Number of background workers
//	        QueueSize:       1000,                 // Event queue buffer size
//	        ProcessTimeout:  time.Second * 30,     // Max processing time per event
//	        RetryAttempts:   3,                    // Retry failed events
//	        RetryDelay:      time.Second * 2,      // Delay between retries
//
//	        // Event storage (for event sourcing)
//	        StorageEnabled:  true,
//	        StorageDriver:   "database",            // or "redis", "memory"
//	        StorageTable:    "events",
//
//	        // Middleware configuration
//	        EnableLogging:   true,
//	        EnableMetrics:   true,
//	        EnableTracing:   true,
//
//	        // Error handling
//	        StopOnError:     false,                // Continue processing other listeners
//	        DeadLetterQueue: true,                 // Store failed events
//
//	        // Performance tuning
//	        BatchSize:       100,                 // Events per batch
//	        FlushInterval:   time.Second * 5,     // Batch flush interval
//	    }
//
//	    eventService, err := event.NewEventDispatcher(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create event service: %v", err)
//	    }
//
//	    return eventService
//	})
func Event() eventInterfaces.EventInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "event" service from the dependency injection container
	// - Performs type assertion to EventInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[eventInterfaces.EventInterface](eventInterfaces.EVENT_TOKEN)
}

// EventWithError provides error-safe access to the event service.
//
// This function offers the same functionality as Event() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle event service unavailability gracefully.
//
// This is a convenience wrapper around facade.Resolve() that provides
// the same caching and performance benefits as Event() but with error handling.
//
// Returns:
//   - EventInterface: The resolved event instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement EventInterface
//
// Usage Examples:
//
//	// Basic error-safe event dispatching
//	eventService, err := facades.EventWithError()
//	if err != nil {
//	    log.Printf("Event service unavailable: %v", err)
//	    // Continue without events or use fallback
//	    return
//	}
//	eventService.Dispatch("user.action", userEvent)
//
//	// Conditional event dispatching
//	if eventService, err := facades.EventWithError(); err == nil {
//	    // Dispatch optional events
//	    eventService.DispatchAsync("analytics.track", trackingEvent)
//	}
//
//	// Health check pattern
//	func CheckEventHealth() error {
//	    eventService, err := facades.EventWithError()
//	    if err != nil {
//	        return fmt.Errorf("event service unavailable: %w", err)
//	    }
//
//	    // Test basic event functionality
//	    testEvent := map[string]interface{}{"test": "health-check"}
//	    err = eventService.Dispatch("system.health_check", testEvent)
//	    if err != nil {
//	        return fmt.Errorf("event dispatch failed: %w", err)
//	    }
//
//	    return nil
//	}
func EventWithError() (eventInterfaces.EventInterface, error) {
	// Use facade.Resolve() for error-return behavior:
	// - Resolves "event" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[eventInterfaces.EventInterface](eventInterfaces.EVENT_TOKEN)
}
