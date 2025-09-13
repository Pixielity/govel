package facades

import (
	queueInterfaces "govel/packages/types/src/interfaces/queue"
	facade "govel/packages/support/src"
)

// Queue provides a clean, static-like interface to the application's job queue service.
//
// This facade implements the facade pattern, providing global access to the queue
// service configured in the dependency injection container. It offers a Laravel-style
// API for background job processing with automatic service resolution, job dispatching,
// worker management, retry handling, and multiple queue driver support.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved queue service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent job operations across goroutines
//   - Supports multiple queue drivers (Redis, Database, Memory, SQS, etc.)
//   - Built-in job serialization, retry logic, and failure handling
//
// Behavior:
//   - First call: Resolves queue service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if queue service cannot be resolved (fail-fast behavior)
//   - Automatically handles job serialization, queuing, and worker coordination
//
// Returns:
//   - QueueInterface: The application's queue service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "queue" service is not registered in the container
//   - If the resolved service doesn't implement QueueInterface
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
//   - Multiple goroutines can call Queue() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Job dispatching and processing are thread-safe
//
// Usage Examples:
//
//	// Basic job dispatching
//	type SendEmailJob struct {
//	    To      string
//	    Subject string
//	    Message string
//	}
//
//	func (j *SendEmailJob) Handle() error {
//	    return facades.Mail().Send(j.To, j.Subject, j.Message)
//	}
//
//	// Dispatch job to default queue
//	job := &SendEmailJob{
//	    To:      "user@example.com",
//	    Subject: "Welcome!",
//	    Message: "Welcome to our platform!",
//	}
//
//	err := facades.Queue().Dispatch(job)
//	if err != nil {
//	    log.Printf("Failed to dispatch job: %v", err)
//	}
//
//	// Dispatch job to specific queue
//	err := facades.Queue().DispatchTo("emails", job)
//
//	// Delayed job execution
//	delay := 5 * time.Minute
//	err := facades.Queue().DispatchAfter(delay, job)
//
//	// Dispatch job at specific time
//	runAt := time.Now().Add(1 * time.Hour)
//	err := facades.Queue().DispatchAt(runAt, job)
//
//	// Job with custom configuration
//	err := facades.Queue().DispatchWithConfig(job, QueueConfig{
//	    Queue:     "high-priority",
//	    Delay:     0,
//	    MaxTries:  5,
//	    Timeout:   time.Minute * 10,
//	    Backoff:   "exponential",
//	})
//
//	// Batch job processing
//	users := []User{
//	    {Email: "user1@example.com", Name: "User 1"},
//	    {Email: "user2@example.com", Name: "User 2"},
//	    {Email: "user3@example.com", Name: "User 3"},
//	}
//
//	batch := facades.Queue().NewBatch()
//	for _, user := range users {
//	    job := &SendEmailJob{
//	        To:      user.Email,
//	        Subject: "Newsletter",
//	        Message: fmt.Sprintf("Hello %s!", user.Name),
//	    }
//	    batch.Add(job)
//	}
//
//	batchResult, err := batch.Dispatch()
//	if err != nil {
//	    log.Printf("Batch dispatch failed: %v", err)
//	}
//
//	// Monitor batch progress
//	for !batchResult.Finished() {
//	    progress := batchResult.Progress()
//	    fmt.Printf("Batch progress: %d/%d completed\n", progress.Completed, progress.Total)
//	    time.Sleep(1 * time.Second)
//	}
//
//	// Job chaining
//	chain := facades.Queue().NewChain()
//	chain.Add(&ProcessOrderJob{OrderID: 123})
//	chain.Add(&SendInvoiceJob{OrderID: 123})
//	chain.Add(&UpdateInventoryJob{OrderID: 123})
//
//	err := chain.Dispatch()
//
// Advanced Job Types:
//
//	// Job with retry logic
//	type ProcessPaymentJob struct {
//	    PaymentID string
//	    Amount    float64
//	    attempts  int
//	}
//
//	func (j *ProcessPaymentJob) Handle() error {
//	    j.attempts++
//
//	    err := processPayment(j.PaymentID, j.Amount)
//	    if err != nil {
//	        if j.attempts < 3 {
//	            // Exponential backoff
//	            delay := time.Duration(math.Pow(2, float64(j.attempts))) * time.Second
//	            return facades.Queue().RetryAfter(delay, err)
//	        }
//	        return err // Give up after 3 attempts
//	    }
//
//	    return nil
//	}
//
//	func (j *ProcessPaymentJob) Failed(err error) {
//	    // Handle permanent failure
//	    facades.Log().Error("Payment processing failed permanently", map[string]interface{}{
//	        "payment_id": j.PaymentID,
//	        "error":      err.Error(),
//	        "attempts":   j.attempts,
//	    })
//
//	    // Notify customer
//	    facades.Mail().Send(
//	        customer.Email,
//	        "Payment Failed",
//	        "We were unable to process your payment. Please try again.",
//	    )
//	}
//
//	// Job with progress tracking
//	type ImportDataJob struct {
//	    FilePath string
//	    BatchID  string
//	}
//
//	func (j *ImportDataJob) Handle() error {
//	    records, err := readCSVFile(j.FilePath)
//	    if err != nil {
//	        return err
//	    }
//
//	    total := len(records)
//	    for i, record := range records {
//	        // Process each record
//	        err := processRecord(record)
//	        if err != nil {
//	            facades.Log().Warning("Failed to process record", map[string]interface{}{
//	                "record": record,
//	                "error":  err.Error(),
//	            })
//	            continue
//	        }
//
//	        // Update progress
//	        progress := float64(i+1) / float64(total) * 100
//	        facades.Queue().SetProgress(j.BatchID, progress)
//	    }
//
//	    return nil
//	}
//
//	// Recurring job
//	type CleanupJob struct {
//	    OlderThan time.Duration
//	}
//
//	func (j *CleanupJob) Handle() error {
//	    cutoff := time.Now().Add(-j.OlderThan)
//
//	    // Clean up old log files
//	    err := cleanupLogs(cutoff)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Clean up temporary files
//	    err = cleanupTempFiles(cutoff)
//	    if err != nil {
//	        return err
//	    }
//
//	    // Schedule next cleanup
//	    nextRun := time.Now().Add(24 * time.Hour)
//	    return facades.Queue().DispatchAt(nextRun, &CleanupJob{
//	        OlderThan: 7 * 24 * time.Hour, // 7 days
//	    })
//	}
//
// Worker Management:
//
//	// Start queue workers
//	func StartWorkers() {
//	    // Start workers for default queue
//	    go facades.Queue().StartWorker("default", 3) // 3 concurrent workers
//
//	    // Start workers for specific queues
//	    go facades.Queue().StartWorker("emails", 5)
//	    go facades.Queue().StartWorker("high-priority", 2)
//	    go facades.Queue().StartWorker("background", 1)
//	}
//
//	// Worker with custom configuration
//	workerConfig := WorkerConfig{
//	    Queue:       "processing",
//	    Workers:     4,
//	    Timeout:     time.Minute * 30,
//	    Sleep:       time.Second * 3,
//	    MaxTries:    5,
//	    BackoffType: "exponential",
//	}
//
//	go facades.Queue().StartWorkerWithConfig(workerConfig)
//
//	// Graceful worker shutdown
//	func GracefulShutdown() {
//	    facades.Queue().StopWorkers()
//
//	    // Wait for current jobs to finish
//	    timeout := time.After(30 * time.Second)
//	    ticker := time.Tick(1 * time.Second)
//
//	loop:
//	    for {
//	        select {
//	        case <-timeout:
//	            fmt.Println("Timeout reached, forcing shutdown")
//	            break loop
//	        case <-ticker:
//	            if facades.Queue().ActiveJobs() == 0 {
//	                fmt.Println("All jobs completed, shutting down")
//	                break loop
//	            }
//	        }
//	    }
//	}
//
// Queue Management:
//
//	// Queue statistics
//	stats := facades.Queue().Stats("default")
//	fmt.Printf("Queue stats: %d waiting, %d processing, %d failed\n",
//	    stats.Waiting, stats.Processing, stats.Failed)
//
//	// Pause and resume queues
//	facades.Queue().Pause("maintenance")
//	performMaintenance()
//	facades.Queue().Resume("maintenance")
//
//	// Clear queue
//	facades.Queue().Clear("failed-jobs")
//
//	// Retry failed jobs
//	failureCount := facades.Queue().RetryFailed("default")
//	fmt.Printf("Retried %d failed jobs\n", failureCount)
//
//	// Get failed jobs
//	failedJobs := facades.Queue().GetFailed(10, 0) // 10 jobs, offset 0
//	for _, job := range failedJobs {
//	    fmt.Printf("Failed job: %s, Error: %s\n", job.Type, job.Error)
//	}
//
// Monitoring and Observability:
//
//	// Job middleware for monitoring
//	type MonitoringMiddleware struct{}
//
//	func (m *MonitoringMiddleware) Handle(job Job, next func() error) error {
//	    start := time.Now()
//
//	    facades.Log().Info("Job started", map[string]interface{}{
//	        "job_type": job.Type(),
//	        "job_id":   job.ID(),
//	        "queue":    job.Queue(),
//	    })
//
//	    err := next()
//	    duration := time.Since(start)
//
//	    if err != nil {
//	        facades.Log().Error("Job failed", map[string]interface{}{
//	            "job_type": job.Type(),
//	            "job_id":   job.ID(),
//	            "duration": duration,
//	            "error":    err.Error(),
//	        })
//	    } else {
//	        facades.Log().Info("Job completed", map[string]interface{}{
//	            "job_type": job.Type(),
//	            "job_id":   job.ID(),
//	            "duration": duration,
//	        })
//	    }
//
//	    return err
//	}
//
//	// Register middleware
//	facades.Queue().Use(&MonitoringMiddleware{})
//
// Best Practices:
//   - Design jobs to be idempotent (safe to run multiple times)
//   - Use appropriate queue names for different job types
//   - Implement proper error handling and retry logic
//   - Monitor queue health and performance metrics
//   - Use batching for related operations
//   - Implement job timeouts to prevent hanging
//   - Use priority queues for time-sensitive jobs
//   - Keep job payloads small and efficient to serialize
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume queue service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	queue, err := facade.TryResolve[QueueInterface]("queue")
//	if err != nil {
//	    // Handle queue service unavailability gracefully
//	    log.Printf("Queue service unavailable: %v", err)
//	    // Maybe execute job synchronously or skip
//	    return job.Handle() // Execute immediately
//	}
//	err = queue.Dispatch(job)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestJobDispatch(t *testing.T) {
//	    // Create a test queue that captures dispatched jobs
//	    testQueue := &TestQueue{
//	        dispatchedJobs: []Job{},
//	    }
//
//	    // Swap the real queue with test queue
//	    restore := support.SwapService("queue", testQueue)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Queue() returns testQueue
//	    emailService := NewEmailService()
//
//	    err := emailService.SendWelcomeEmail("user@example.com")
//	    require.NoError(t, err)
//
//	    // Verify job was dispatched
//	    jobs := testQueue.GetDispatchedJobs()
//	    assert.Len(t, jobs, 1)
//
//	    emailJob, ok := jobs[0].(*SendEmailJob)
//	    require.True(t, ok)
//	    assert.Equal(t, "user@example.com", emailJob.To)
//	    assert.Contains(t, emailJob.Subject, "Welcome")
//	}
//
// Container Configuration:
// Ensure the queue service is properly configured in your container:
//
//	// Example queue registration
//	container.Singleton("queue", func() interface{} {
//	    config := queue.Config{
//	        // Default connection
//	        Default: "redis",
//
//	        // Queue connections
//	        Connections: map[string]queue.ConnectionConfig{
//	            "redis": {
//	                Driver: "redis",
//	                Host:   "localhost:6379",
//	                DB:     0,
//	                Prefix: "queues:",
//	            },
//	            "database": {
//	                Driver:     "database",
//	                Connection: "default",
//	                Table:      "jobs",
//	            },
//	            "sqs": {
//	                Driver:    "sqs",
//	                Region:    "us-east-1",
//	                AccessKey: "your-access-key",
//	                SecretKey: "your-secret-key",
//	                QueueURL:  "https://sqs.us-east-1.amazonaws.com/123456789/my-queue",
//	            },
//	        },
//
//	        // Failed job configuration
//	        Failed: queue.FailedConfig{
//	            Driver:     "database",
//	            Connection: "default",
//	            Table:      "failed_jobs",
//	        },
//
//	        // Worker configuration
//	        Worker: queue.WorkerConfig{
//	            Sleep:       time.Second * 3,
//	            MaxTries:    3,
//	            Timeout:     time.Minute * 60,
//	            BackoffType: "exponential",
//	        },
//
//	        // Serialization
//	        Serializer: "json", // json, gob, msgpack
//	    }
//
//	    queueService, err := queue.NewQueueManager(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create queue service: %v", err)
//	    }
//
//	    return queueService
//	})
func Queue() queueInterfaces.QueueInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "queue" service from the dependency injection container
	// - Performs type assertion to QueueInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[queueInterfaces.QueueInterface](queueInterfaces.QUEUE_TOKEN)
}

// QueueWithError provides error-safe access to the queue service.
//
// This function offers the same functionality as Queue() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle queue service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Queue() but with error handling.
//
// Returns:
//   - QueueInterface: The resolved queue instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement QueueInterface
//
// Usage Examples:
//
//	// Basic error-safe job dispatching
//	queue, err := facades.QueueWithError()
//	if err != nil {
//	    log.Printf("Queue service unavailable: %v", err)
//	    // Execute job synchronously as fallback
//	    return job.Handle()
//	}
//	err = queue.Dispatch(job)
//
//	// Conditional job queuing
//	if queue, err := facades.QueueWithError(); err == nil {
//	    // Queue optional background tasks
//	    queue.DispatchAfter(time.Hour, &CleanupJob{})
//	} else {
//	    // Log that background task couldn't be queued
//	    log.Printf("Background cleanup skipped: %v", err)
//	}
//
//	// Health check pattern
//	func CheckQueueHealth() error {
//	    queue, err := facades.QueueWithError()
//	    if err != nil {
//	        return fmt.Errorf("queue service unavailable: %w", err)
//	    }
//
//	    // Test basic queue functionality
//	    stats := queue.Stats("default")
//	    if stats.Workers == 0 {
//	        return fmt.Errorf("no workers running for default queue")
//	    }
//
//	    return nil
//	}
func QueueWithError() (queueInterfaces.QueueInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "queue" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[queueInterfaces.QueueInterface](queueInterfaces.QUEUE_TOKEN)
}
