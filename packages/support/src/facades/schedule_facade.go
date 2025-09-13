package facades
import (
	scheduleInterfaces "govel/types/src/interfaces/schedule"
	facade "govel/support/src"
)
// Schedule provides a clean, static-like interface to the application's task scheduling service.
//
// This facade implements the facade pattern, providing global access to the scheduling
// service configured in the dependency injection container. It offers a Laravel-style
// API for job scheduling, cron management, task queuing, and background job processing
// with automatic service resolution and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved schedule service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent scheduling operations across goroutines
//   - Supports cron expressions, delayed tasks, and recurring jobs
//   - Built-in job queuing and background processing integration
//
// Behavior:
//   - First call: Resolves schedule service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if schedule service cannot be resolved (fail-fast behavior)
//   - Automatically handles job registration, execution, and monitoring
//
// Returns:
//   - ScheduleInterface: The application's scheduling service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "schedule" service is not registered in the container
//   - If the resolved service doesn't implement ScheduleInterface
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
//   - Multiple goroutines can call Schedule() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Job scheduling and execution are thread-safe and consistent
//
// Usage Examples:
//
//	// Basic job scheduling
//	facades.Schedule().Command("send:emails").Daily()
//	facades.Schedule().Command("cleanup:logs").Weekly()
//	facades.Schedule().Command("generate:reports").Monthly()
//
//	// Cron expression scheduling
//	facades.Schedule().Command("backup:database").Cron("0 2 * * *") // Every day at 2 AM
//	facades.Schedule().Command("update:cache").Cron("*/5 * * * *") // Every 5 minutes
//	facades.Schedule().Command("send:newsletter").Cron("0 9 * * MON") // Monday at 9 AM
//
//	// Time-based scheduling
//	facades.Schedule().Command("process:orders").EveryMinute()
//	facades.Schedule().Command("sync:data").EveryFiveMinutes()
//	facades.Schedule().Command("check:health").EveryTenMinutes()
//	facades.Schedule().Command("rotate:logs").EveryThirtyMinutes()
//	facades.Schedule().Command("update:stats").Hourly()
//	facades.Schedule().Command("send:reminders").DailyAt("14:30")
//
//	// Conditional scheduling
//	facades.Schedule().Command("send:report").Daily().When(func() bool {
//	    return time.Now().Weekday() != time.Saturday && time.Now().Weekday() != time.Sunday
//	})
//
//	facades.Schedule().Command("backup:files").Weekly().Unless(func() bool {
//	    return facades.App().IsLocal() // Skip backup in local environment
//	})
//
//	// Environment-specific scheduling
//	facades.Schedule().Command("debug:cleanup").Daily().Environments("local", "development")
//	facades.Schedule().Command("monitor:production").EveryMinute().Environments("production")
//
//	// Job callbacks and hooks
//	facades.Schedule().Command("process:payments").Hourly().
//	    Before(func() {
//	        log.Println("Starting payment processing...")
//	    }).
//	    After(func() {
//	        log.Println("Payment processing completed")
//	    }).
//	    OnSuccess(func() {
//	        facades.Log().Info("Payment processing succeeded")
//	    }).
//	    OnFailure(func() {
//	        facades.Log().Error("Payment processing failed")
//	        facades.Mail().Send("admin@example.com", "Payment Processing Failed", "...")
//	    })
//
//	// Background job scheduling
//	facades.Schedule().Job(NewEmailJob("user@example.com", "Welcome!")).Daily()
//	facades.Schedule().Job(NewReportGenerationJob(userID)).Weekly()
//	facades.Schedule().Job(NewDataSyncJob()).EveryFiveMinutes()
//
//	// Closure-based scheduling
//	facades.Schedule().Call(func() {
//	    // Clean up temporary files
//	    os.RemoveAll("/tmp/app-cache")
//	}).Daily()
//
//	facades.Schedule().Call(func() {
//	    // Update application metrics
//	    updateMetrics()
//	}).EveryTenMinutes()
//
//	// Queue job scheduling
//	facades.Schedule().QueueJob("send-email", map[string]interface{}{
//	    "to":      "user@example.com",
//	    "subject": "Daily Report",
//	    "template": "daily-report",
//	}).Daily()
//
//	// Delayed job scheduling
//	facades.Schedule().Delay(30*time.Minute, func() {
//	    // Execute after 30 minutes
//	    processDelayedTask()
//	})
//
//	facades.Schedule().DelayUntil(time.Now().Add(2*time.Hour), func() {
//	    // Execute at specific time
//	    processTimedTask()
//	})
//
//	// Recurring jobs with limits
//	facades.Schedule().Command("trial:cleanup").Daily().Times(30) // Run for 30 days
//	facades.Schedule().Command("migration:check").Hourly().Until(migrationComplete)
//
//	// Job monitoring and management
//	allJobs := facades.Schedule().GetJobs()
//	runningJobs := facades.Schedule().GetRunningJobs()
//	jobHistory := facades.Schedule().GetJobHistory("send:emails")
//
//	// Job control
//	facades.Schedule().PauseJob("backup:database")
//	facades.Schedule().ResumeJob("backup:database")
//	facades.Schedule().CancelJob("long:running:task")
//
// Advanced Scheduling Patterns:
//
//	// Job chains and dependencies
//	facades.Schedule().Chain([]ScheduledJob{
//	    facades.Schedule().Command("backup:database"),
//	    facades.Schedule().Command("compress:backup"),
//	    facades.Schedule().Command("upload:backup"),
//	    facades.Schedule().Command("cleanup:old:backups"),
//	}).Daily()
//
//	// Parallel job execution
//	facades.Schedule().Parallel([]ScheduledJob{
//	    facades.Schedule().Command("process:images"),
//	    facades.Schedule().Command("process:videos"),
//	    facades.Schedule().Command("process:documents"),
//	}).Hourly()
//
//	// Dynamic job registration
//	func RegisterDynamicJobs() {
//	    // Load jobs from configuration
//	    jobConfigs := facades.Config().Get("scheduler.jobs")
//	    for _, jobConfig := range jobConfigs {
//	        facades.Schedule().Command(jobConfig.Command).
//	            Cron(jobConfig.Schedule).
//	            Environments(jobConfig.Environments...)
//	    }
//
//	    // Register user-specific jobs
//	    activeUsers := getUsersWithScheduledTasks()
//	    for _, user := range activeUsers {
//	        facades.Schedule().Job(NewUserTaskJob(user.ID)).
//	            Daily().
//	            Name(fmt.Sprintf("user-task-%d", user.ID))
//	    }
//	}
//
//	// Job retry and failure handling
//	facades.Schedule().Command("critical:task").Hourly().
//	    Retry(3, 5*time.Minute). // Retry 3 times with 5-minute intervals
//	    OnFailure(func() {
//	        facades.Mail().Send("admin@example.com", "Critical Task Failed", "...")
//	        facades.Log().Error("Critical task failed after all retries")
//	    })
//
//	// Job timeouts and resource limits
//	facades.Schedule().Command("long:process").Daily().
//	    Timeout(30*time.Minute).
//	    MemoryLimit("512MB").
//	    OnTimeout(func() {
//	        facades.Log().Warning("Long process timed out")
//	    })
//
// Best Practices:
//   - Use descriptive names for scheduled jobs
//   - Set appropriate timeouts for long-running tasks
//   - Implement proper error handling and retry logic
//   - Use environment-specific scheduling when needed
//   - Monitor job performance and execution times
//   - Log job start, completion, and failure events
//   - Use job dependencies to ensure proper execution order
//   - Set resource limits for resource-intensive jobs
//
// Scheduling Patterns:
//  1. Use cron expressions for complex scheduling requirements
//  2. Implement job chains for dependent tasks
//  3. Use parallel execution for independent tasks
//  4. Set appropriate retry policies for critical jobs
//  5. Monitor job execution and handle failures gracefully
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume scheduling always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	scheduler, err := facade.TryResolve[ScheduleInterface]("schedule")
//	if err != nil {
//	    // Handle scheduler unavailability gracefully
//	    return fmt.Errorf("scheduler unavailable: %w", err)
//	}
//	scheduler.Command("backup").Daily()
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestJobScheduling(t *testing.T) {
//	    // Create a test scheduler
//	    testScheduler := &TestScheduler{
//	        jobs: make(map[string]ScheduledJob),
//	    }
//
//	    // Swap the real scheduler with test scheduler
//	    restore := support.SwapService("schedule", testScheduler)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Schedule() returns testScheduler
//	    facades.Schedule().Command("test:job").Daily()
//
//	    // Verify job registration
//	    assert.True(t, testScheduler.HasJob("test:job"))
//	    job := testScheduler.GetJob("test:job")
//	    assert.Equal(t, "daily", job.Frequency())
//	}
//
// Container Configuration:
// Ensure the schedule service is properly configured in your container:
//
//	// Example scheduler registration
//	container.Singleton("schedule", func() interface{} {
//	    config := scheduler.Config{
//	        // Scheduler configuration
//	        Timezone:        "UTC",
//	        MaxConcurrency: 10,
//	        JobTimeout:      30 * time.Minute,
//
//	        // Job storage
//	        JobStorage: "redis", // redis, database, memory
//	        StorageConfig: map[string]interface{}{
//	            "redis_url": "redis://localhost:6379/0",
//	        },
//
//	        // Job retry configuration
//	        DefaultRetries:      3,
//	        DefaultRetryDelay:   5 * time.Minute,
//	        MaxRetryDelay:       30 * time.Minute,
//
//	        // Monitoring and logging
//	        LogJobStart:    true,
//	        LogJobComplete: true,
//	        LogJobFailure:  true,
//	        MetricsEnabled: true,
//
//	        // Job discovery
//	        JobPaths: []string{"/app/jobs", "/app/commands"},
//	    }
//
//	    return scheduler.NewScheduler(config)
//	})
func Schedule() scheduleInterfaces.ScheduleInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "schedule" service from the dependency injection container
	// - Performs type assertion to ScheduleInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[scheduleInterfaces.ScheduleInterface](scheduleInterfaces.SCHEDULE_TOKEN)
}

// ScheduleWithError provides error-safe access to the scheduling service.
//
// This function offers the same functionality as Schedule() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle scheduling unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Schedule() but with error handling.
//
// Returns:
//   - ScheduleInterface: The resolved schedule instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ScheduleInterface
//
// Usage Examples:
//
//	// Basic error-safe scheduling
//	scheduler, err := facades.ScheduleWithError()
//	if err != nil {
//	    log.Printf("Scheduler unavailable: %v", err)
//	    return fmt.Errorf("scheduling service not available")
//	}
//	scheduler.Command("backup").Daily()
//
//	// Conditional job registration
//	if scheduler, err := facades.ScheduleWithError(); err == nil {
//	    scheduler.Command("optional:task").Weekly()
//	}
func ScheduleWithError() (scheduleInterfaces.ScheduleInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "schedule" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[scheduleInterfaces.ScheduleInterface](scheduleInterfaces.SCHEDULE_TOKEN)
}
