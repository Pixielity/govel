package compiler

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	interfaces "govel/support/compiler/interfaces"
	types "govel/support/compiler/types"
)

// GoVelCompiler is the main implementation of the Compiler interface
// It provides a clean, efficient API for Go code compilation and execution
type GoVelCompiler struct {
	// config holds the compiler configuration
	config *types.Config

	// cache stores compilation results for faster subsequent executions
	cache   map[string]*types.CacheEntry
	cacheMu sync.RWMutex

	// metrics tracks compilation statistics
	metrics *types.Metrics

	// tempDirs holds temporary directories for cleanup
	tempDirs []string
	tempMu   sync.Mutex

	// closed indicates if the compiler has been closed
	closed bool
	mu     sync.RWMutex
}

// NewCompiler creates a new GoVel compiler instance with the specified configuration.
// This function initializes all internal components including cache, metrics tracking,
// and temporary directory management. It follows Go best practices for optional configuration.
//
// Parameters:
//
//	config: The compiler configuration to use. If nil, default configuration will be applied
//
// Returns:
//
//	interfaces.CompilerInterface: A fully initialized compiler instance ready for use
func NewCompiler(config *types.Config) interfaces.CompilerInterface {
	// Merge provided config with defaults
	defaultConfig := types.DefaultConfig()
	finalConfig := defaultConfig.Merge(config)

	// Validate the final configuration
	finalConfig.Validate()

	compiler := &GoVelCompiler{
		config:   finalConfig,
		cache:    make(map[string]*types.CacheEntry),
		metrics:  types.NewMetrics(),
		tempDirs: make([]string, 0),
		closed:   false,
	}

	return compiler
}

// Compile compiles and executes a Go file from the filesystem, returning the execution result.
// This is the main entry point for file-based compilation that provides a simple interface
// without context handling. For more control, use CompileWithContext instead.
//
// Parameters:
//
//	filePath: The absolute or relative path to the Go source file to compile and execute
//
// Returns:
//
//	*types.Result: The compilation and execution result containing output, errors, and metrics
func (c *GoVelCompiler) Compile(filePath string) *types.Result {
	return c.CompileWithContext(context.Background(), filePath)
}

// CompileWithContext compiles and executes a Go file with context support for cancellation and timeout control.
// This method provides fine-grained control over the compilation process, allowing for timeouts,
// cancellation, and deadline management. It reads the file from disk and delegates to CompileCodeWithContext.
//
// Parameters:
//
//	ctx: The context for controlling cancellation, timeouts, and deadlines
//	filePath: The absolute or relative path to the Go source file to compile and execute
//
// Returns:
//
//	*types.Result: The compilation and execution result with file path information included
func (c *GoVelCompiler) CompileWithContext(ctx context.Context, filePath string) *types.Result {
	startTime := time.Now()

	// Check if compiler is closed
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return types.NewErrorResult(types.StatusRuntimeError, fmt.Errorf("compiler is closed"))
	}
	c.mu.RUnlock()

	// Create context with timeout from config
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		result := types.NewErrorResult(types.StatusCompileError, fmt.Errorf("failed to read file %s: %w", filePath, err))
		c.recordResult(result, startTime)
		return result
	}

	// Use CompileCodeWithContext for actual compilation
	result := c.CompileCodeWithContext(ctx, string(content))
	result.FilePath = filePath

	return result
}

// CompileCode compiles and executes Go code provided as a string without context support.
// This method provides a simple interface for compiling code snippets or dynamically
// generated Go code. For timeout and cancellation support, use CompileCodeWithContext.
//
// Parameters:
//
//	code: The Go source code as a string to compile and execute
//
// Returns:
//
//	*types.Result: The compilation and execution result containing output, errors, and performance metrics
func (c *GoVelCompiler) CompileCode(code string) *types.Result {
	return c.CompileCodeWithContext(context.Background(), code)
}

// CompileCodeWithContext compiles and executes Go code with full context support and caching.
// This is the core compilation method that handles caching, validation, temporary file management,
// and execution. It provides comprehensive error handling and performance monitoring.
//
// Parameters:
//
//	ctx: The context for controlling cancellation, timeouts, and deadlines
//	code: The Go source code as a string to compile and execute
//
// Returns:
//
//	*types.Result: A comprehensive result object containing execution output, errors, timing, and cache information
func (c *GoVelCompiler) CompileCodeWithContext(ctx context.Context, code string) *types.Result {
	startTime := time.Now()

	// Check if compiler is closed
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return types.NewErrorResult(types.StatusRuntimeError, fmt.Errorf("compiler is closed"))
	}
	c.mu.RUnlock()

	// Generate cache key
	hash := c.generateHash(code)

	// Check cache first if enabled
	if c.config.EnableCache {
		if cachedResult := c.getCachedResult(hash); cachedResult != nil {
			cachedResult.CacheHit = true
			cachedResult.Status = types.StatusCacheHit
			c.recordResult(cachedResult, startTime)
			return cachedResult
		}
	}

	// Create result object
	result := types.NewResult()
	result.Hash = hash

	// Validate code if unsafe imports are not allowed
	if !c.config.AllowUnsafeImports {
		if validationResult := c.validateCode(code); !validationResult.Valid {
			result.Status = types.StatusValidationError
			result.Success = false
			result.CompilationError = strings.Join(validationResult.Errors, "; ")
			c.recordResult(result, startTime)
			return result
		}
	}

	// Create temporary directory
	tempDir, err := c.createTempDir()
	if err != nil {
		result.Status = types.StatusCompileError
		result.Success = false
		result.CompilationError = fmt.Sprintf("failed to create temp directory: %v", err)
		c.recordResult(result, startTime)
		return result
	}
	defer c.cleanupTempDir(tempDir)

	compileStart := time.Now()

	// Compile and execute the code
	if err := c.executeCode(ctx, code, tempDir, result); err != nil {
		result.Status = types.StatusRuntimeError
		result.Success = false
		result.RuntimeError = err.Error()
	}

	result.CompileTime = time.Since(compileStart)
	result.Duration = time.Since(startTime)
	result.ExecutionTime = result.Duration - result.CompileTime

	// Cache successful results if enabled
	if c.config.EnableCache && result.Success {
		c.cacheResult(hash, result)
	}

	// Record metrics
	c.recordResult(result, startTime)

	return result
}

// GetConfig returns the current compiler configuration settings.
// This method provides read-only access to the compiler's configuration,
// allowing inspection of current settings without modification.
//
// Returns:
//
//	*types.Config: The current configuration object containing all compiler settings
func (c *GoVelCompiler) GetConfig() *types.Config {
	return c.config
}

// UpdateConfig updates the compiler configuration with new settings.
// This method safely merges the provided configuration with the current one,
// validates the result, and applies it atomically. The compiler must not be closed.
//
// Parameters:
//
//	config: The new configuration settings to merge with the current configuration
//
// Returns:
//
//	error: An error if the compiler is closed or configuration validation fails
func (c *GoVelCompiler) UpdateConfig(config *types.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("cannot update config: compiler is closed")
	}

	// Merge with current config and validate
	newConfig := c.config.Merge(config)
	newConfig.Validate()

	c.config = newConfig
	return nil
}

// Close performs cleanup and resource deallocation for the compiler instance.
// This method safely shuts down the compiler, clears the cache, and optionally
// removes temporary directories based on configuration. It's safe to call multiple times.
//
// Returns:
//
//	error: Always returns nil, but implements the Closer interface for consistency
func (c *GoVelCompiler) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil // Already closed
	}

	c.closed = true

	// Clear cache
	c.cacheMu.Lock()
	c.cache = make(map[string]*types.CacheEntry)
	c.cacheMu.Unlock()

	// Clean up temp directories if configured
	if c.config.CleanupOnExit {
		c.tempMu.Lock()
		for _, tempDir := range c.tempDirs {
			os.RemoveAll(tempDir)
		}
		c.tempDirs = nil
		c.tempMu.Unlock()
	}

	return nil
}

// GetMetrics returns comprehensive compilation metrics and performance statistics.
// This method provides access to accumulated performance data including execution times,
// success rates, cache hit ratios, and resource usage across all compilation operations.
//
// Returns:
//
//	*types.Metrics: The metrics object containing performance statistics and counters
func (c *GoVelCompiler) GetMetrics() *types.Metrics {
	return c.metrics
}

// Private helper methods

// generateHash creates a SHA256 hash of the provided code for caching purposes.
// This method generates a deterministic hash key used for cache lookups and storage.
// Only the first 16 bytes of the hash are used to keep cache keys manageable.
//
// Parameters:
//
//	code: The Go source code to generate a hash for
//
// Returns:
//
//	string: A hexadecimal string representation of the truncated SHA256 hash
func (c *GoVelCompiler) generateHash(code string) string {
	hash := sha256.Sum256([]byte(code))
	return fmt.Sprintf("%x", hash[:16]) // Use first 16 bytes for shorter hash
}

// getCachedResult retrieves a compilation result from cache if available and not expired.
// This method performs cache lookup with TTL validation and automatic cleanup of expired entries.
// It also updates access statistics for the cache entry and returns a copy of the result.
//
// Parameters:
//
//	hash: The hash key generated from the source code to look up in cache
//
// Returns:
//
//	*types.Result: A copy of the cached result, or nil if not found or expired
func (c *GoVelCompiler) getCachedResult(hash string) *types.Result {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	entry, exists := c.cache[hash]
	if !exists {
		return nil
	}

	// Check if expired
	if entry.IsExpired(c.config.CacheTTL) {
		// Remove expired entry (cleanup happens in separate goroutine in production)
		go func() {
			c.cacheMu.Lock()
			delete(c.cache, hash)
			c.cacheMu.Unlock()
		}()
		return nil
	}

	// Update access statistics
	entry.Touch()

	// Return a copy of the result
	resultCopy := *entry.Result
	return &resultCopy
}

// cacheResult stores a compilation result in the cache with LRU eviction.
// This method handles cache size limits by evicting the oldest entry when necessary
// and creates a new cache entry with access tracking and metadata.
//
// Parameters:
//
//	hash: The hash key generated from the source code for cache storage
//	result: The compilation result to store in the cache
func (c *GoVelCompiler) cacheResult(hash string, result *types.Result) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	// Check cache size limit
	if len(c.cache) >= c.config.CacheSize {
		// Remove oldest entry (simple LRU)
		c.evictOldestCacheEntry()
	}

	// Create cache entry
	entry := &types.CacheEntry{
		Result:       result,
		Hash:         hash,
		Timestamp:    time.Now(),
		AccessCount:  1,
		LastAccessed: time.Now(),
		Size:         int64(len(result.Content) + len(result.ErrorOutput)),
	}

	c.cache[hash] = entry
}

// evictOldestCacheEntry removes the oldest cache entry to make space for new entries.
// This method implements a simple LRU (Least Recently Used) eviction policy by
// finding the entry with the oldest LastAccessed timestamp and removing it.
func (c *GoVelCompiler) evictOldestCacheEntry() {
	var oldestKey string
	var oldestTime time.Time = time.Now()

	for key, entry := range c.cache {
		if entry.LastAccessed.Before(oldestTime) {
			oldestTime = entry.LastAccessed
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

// createTempDir creates a temporary directory for compilation operations.
// This method creates a unique temporary directory using the configured base path
// and tracks it for cleanup. The directory is used for Go module initialization and compilation.
//
// Returns:
//
//	string: The absolute path to the created temporary directory
//	error: An error if directory creation fails
func (c *GoVelCompiler) createTempDir() (string, error) {
	baseDir := c.config.TempDir
	if baseDir == "" {
		baseDir = os.TempDir()
	}

	tempDir, err := os.MkdirTemp(baseDir, "govel-compiler-*")
	if err != nil {
		return "", err
	}

	// Track temp directory for cleanup
	c.tempMu.Lock()
	c.tempDirs = append(c.tempDirs, tempDir)
	c.tempMu.Unlock()

	return tempDir, nil
}

// cleanupTempDir removes a temporary directory and updates internal tracking.
// This method performs filesystem cleanup and removes the directory from the
// internal tracking list to prevent memory leaks.
//
// Parameters:
//
//	tempDir: The absolute path to the temporary directory to remove
func (c *GoVelCompiler) cleanupTempDir(tempDir string) {
	os.RemoveAll(tempDir)

	// Remove from tracking list
	c.tempMu.Lock()
	for i, dir := range c.tempDirs {
		if dir == tempDir {
			c.tempDirs = append(c.tempDirs[:i], c.tempDirs[i+1:]...)
			break
		}
	}
	c.tempMu.Unlock()
}

// validateCode performs basic validation and security checks on Go source code.
// This method checks for unsafe imports, calculates complexity scores, and
// generates warnings for potentially problematic code patterns.
//
// Parameters:
//
//	code: The Go source code to validate
//
// Returns:
//
//	*types.ValidationResult: A detailed validation result with errors, warnings, and metrics
func (c *GoVelCompiler) validateCode(code string) *types.ValidationResult {
	result := types.NewValidationResult()

	// Check for unsafe imports
	unsafeImports := []string{"unsafe", "syscall", "reflect"}
	for _, unsafeImport := range unsafeImports {
		if strings.Contains(code, fmt.Sprintf(`"%s"`, unsafeImport)) {
			result.AddUnsafeImport(unsafeImport)
			result.AddError(fmt.Sprintf("unsafe import '%s' not allowed", unsafeImport))
		}
	}

	// Basic complexity check
	lines := strings.Split(code, "\n")
	result.ComplexityScore = len(lines) / 10 // Very simple complexity measure
	if result.ComplexityScore > 100 {
		result.AddWarning("code complexity is high")
	}

	return result
}

// executeCode compiles and runs Go code in a temporary directory, populating the result object.
// This method handles the complete execution pipeline including file creation, module initialization,
// and program execution with proper context handling and error management.
//
// Parameters:
//
//	ctx: The context for controlling cancellation and timeouts
//	code: The Go source code to execute
//	tempDir: The temporary directory path for compilation artifacts
//	result: The result object to populate with execution data
//
// Returns:
//
//	error: An error if any step of the execution process fails
func (c *GoVelCompiler) executeCode(ctx context.Context, code, tempDir string, result *types.Result) error {
	// Write main.go file
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	// Initialize Go module if needed
	if err := c.initGoModule(tempDir); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Run the Go program
	return c.runGoProgram(ctx, tempDir, result)
}

// initGoModule initializes a Go module in the specified temporary directory.
// This method runs 'go mod init' with the configured module name and custom environment
// to prepare the temporary directory for compilation.
//
// Parameters:
//
//	tempDir: The temporary directory path where the Go module should be initialized
//
// Returns:
//
//	error: An error if module initialization fails
func (c *GoVelCompiler) initGoModule(tempDir string) error {
	cmd := exec.Command("go", "mod", "init", c.config.ModuleName)
	cmd.Dir = tempDir
	cmd.Env = c.buildEnvironment()

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init failed: %s", string(output))
	}

	return nil
}

// runGoProgram executes the Go program and captures comprehensive execution results.
// This method runs 'go run main.go' in the temporary directory, measures memory usage,
// captures output, and handles exit codes with proper context cancellation support.
//
// Parameters:
//
//	ctx: The context for controlling execution cancellation and timeouts
//	tempDir: The temporary directory containing the Go program to execute
//	result: The result object to populate with execution data and metrics
//
// Returns:
//
//	error: An error if program execution fails or is cancelled
func (c *GoVelCompiler) runGoProgram(ctx context.Context, tempDir string, result *types.Result) error {
	// Track memory usage
	var memStart, memEnd runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// Build command
	cmd := exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Dir = tempDir
	cmd.Env = c.buildEnvironment()

	// Execute and capture output
	output, err := cmd.CombinedOutput()

	// Measure memory usage
	runtime.ReadMemStats(&memEnd)
	result.MemoryUsed = int64(memEnd.Alloc - memStart.Alloc)

	// Set result content and error information
	result.Content = output
	if err != nil {
		result.ErrorOutput = output
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
		return err
	}

	result.ExitCode = 0
	return nil
}

// buildEnvironment builds the complete environment variable set for Go command execution.
// This method combines the system environment with custom variables and build tags
// from the compiler configuration to create the execution environment.
//
// Returns:
//
//	[]string: A slice of environment variables in KEY=VALUE format for command execution
func (c *GoVelCompiler) buildEnvironment() []string {
	env := os.Environ()

	// Add custom environment variables
	for key, value := range c.config.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add build tags if specified
	if len(c.config.BuildTags) > 0 {
		tags := strings.Join(c.config.BuildTags, ",")
		env = append(env, fmt.Sprintf("GOFLAGS=-tags=%s", tags))
	}

	return env
}

// recordResult records metrics and timing information for a compilation result.
// This method ensures proper duration calculation and delegates to the metrics
// system for performance tracking and statistical analysis.
//
// Parameters:
//
//	result: The compilation result to record metrics for
//	startTime: The timestamp when the compilation operation began
func (c *GoVelCompiler) recordResult(result *types.Result, startTime time.Time) {
	if result.Duration == 0 {
		result.Duration = time.Since(startTime)
	}
	c.metrics.RecordCompilation(result)
}
