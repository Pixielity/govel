package types

import "time"

// Config represents comprehensive configuration options for the GoVel compiler.
// Provides fine-grained control over compiler behavior, performance, and security settings.
type Config struct {
	// TempDir specifies the base directory for temporary files during compilation.
	// Empty string uses the system's default temporary directory.
	TempDir string `json:"temp_dir,omitempty"`

	// ModuleName sets the default Go module name for compiled code.
	// Used when initializing Go modules for code without existing module definition.
	ModuleName string `json:"module_name,omitempty"`

	// Timeout specifies the maximum time allowed for compilation and execution operations.
	// Applies to the entire compilation pipeline with graceful cancellation on timeout.
	Timeout time.Duration `json:"timeout,omitempty"`

	// MaxMemory sets the maximum memory limit for compilation processes in bytes.
	// Prevents memory exhaustion attacks and ensures system stability.
	MaxMemory int64 `json:"max_memory,omitempty"`

	// EnableCache determines whether to cache compilation results for improved performance.
	// Uses SHA-256 hashes as keys with LRU eviction and TTL-based expiration.
	EnableCache bool `json:"enable_cache"`

	// CacheSize sets the maximum number of entries in the compilation cache.
	// Prevents unbounded memory growth with LRU eviction when limit is reached.
	CacheSize int `json:"cache_size,omitempty"`

	// CacheTTL sets how long cache entries remain valid before expiration.
	// Provides automatic cleanup with lazy expiration on cache access.
	CacheTTL time.Duration `json:"cache_ttl,omitempty"`

	// Debug enables verbose debug logging throughout the compilation process.
	// Includes detailed information about pipeline steps, cache operations, and timing.
	Debug bool `json:"debug"`

	// Verbose enables verbose output during compilation operations.
	// Focuses on user-facing operational information rather than internal debugging details.
	Verbose bool `json:"verbose"`

	// CleanupOnExit determines whether to clean temporary directories on compiler close.
	// Automatically removes temporary directories and files to prevent resource leaks.
	CleanupOnExit bool `json:"cleanup_on_exit"`

	// MaxConcurrentJobs sets the maximum number of concurrent compilation jobs.
	// Limits simultaneous operations to prevent resource exhaustion and maintain system stability.
	MaxConcurrentJobs int `json:"max_concurrent_jobs,omitempty"`

	// AllowUnsafeImports permits potentially dangerous imports in compiled code.
	// When disabled (default), blocks unsafe, syscall, and reflect imports for security.
	AllowUnsafeImports bool `json:"allow_unsafe_imports"`

	// Environment contains environment variables to set during compilation.
	// Variables are passed to the Go toolchain for compilation behavior control.
	Environment map[string]string `json:"environment,omitempty"`

	// BuildTags specifies build tags to use during compilation.
	// Enables conditional compilation for platform-specific code and feature toggles.
	BuildTags []string `json:"build_tags,omitempty"`
}

// DefaultConfig returns a Config struct with all default values set.
// Provides baseline configuration with sensible defaults for production use.
//
// Returns:
//
//	*Config: A fully configured Config struct with production-ready defaults
func DefaultConfig() *Config {
	return &Config{
		TempDir:            "", // Use system temp directory
		ModuleName:         "temp-module",
		Timeout:            30 * time.Second,
		MaxMemory:          256 * 1024 * 1024, // 256MB
		EnableCache:        true,
		CacheSize:          100,
		CacheTTL:           10 * time.Minute,
		Debug:              false,
		Verbose:            false,
		CleanupOnExit:      true,
		MaxConcurrentJobs:  4,
		AllowUnsafeImports: false,
		Environment:        make(map[string]string),
		BuildTags:          make([]string, 0),
	}
}

// Merge merges the provided config with defaults, returning a new Config.
// Non-zero/non-nil values from the provided config take precedence over current values.
//
// Parameters:
//
//	other: The configuration to merge with current config
//
// Returns:
//
//	*Config: A new Config struct with merged values
func (c *Config) Merge(other *Config) *Config {
	// Handle nil cases
	if c == nil {
		c = DefaultConfig()
	}
	if other == nil {
		return c
	}

	// Create a copy of the current config
	merged := *c

	// Override with non-zero values from other config
	if other.TempDir != "" {
		merged.TempDir = other.TempDir
	}
	if other.ModuleName != "" {
		merged.ModuleName = other.ModuleName
	}
	if other.Timeout != 0 {
		merged.Timeout = other.Timeout
	}
	if other.MaxMemory != 0 {
		merged.MaxMemory = other.MaxMemory
	}

	// Boolean fields always override (including false values)
	merged.EnableCache = other.EnableCache
	merged.Debug = other.Debug
	merged.Verbose = other.Verbose
	merged.CleanupOnExit = other.CleanupOnExit
	merged.AllowUnsafeImports = other.AllowUnsafeImports

	if other.CacheSize != 0 {
		merged.CacheSize = other.CacheSize
	}
	if other.CacheTTL != 0 {
		merged.CacheTTL = other.CacheTTL
	}
	if other.MaxConcurrentJobs != 0 {
		merged.MaxConcurrentJobs = other.MaxConcurrentJobs
	}

	// Map and slice fields replace entirely if non-empty
	if other.Environment != nil && len(other.Environment) > 0 {
		merged.Environment = other.Environment
	}
	if other.BuildTags != nil && len(other.BuildTags) > 0 {
		merged.BuildTags = other.BuildTags
	}

	return &merged
}

// Validate checks if the configuration values are valid and reasonable.
// Performs validation and auto-corrects invalid values to safe minimums.
//
// Returns:
//
//	error: Always nil in current implementation
func (c *Config) Validate() error {
	// Enforce minimum timeout to prevent immediate failures
	if c.Timeout < time.Second {
		c.Timeout = time.Second
	}

	// Enforce minimum memory to allow basic compilation
	if c.MaxMemory < 1024*1024 { // 1MB minimum
		c.MaxMemory = 1024 * 1024
	}

	// Enforce minimum cache size to maintain functionality
	if c.CacheSize < 1 {
		c.CacheSize = 1
	}

	// Enforce minimum TTL for reasonable cache behavior
	if c.CacheTTL < time.Minute {
		c.CacheTTL = time.Minute
	}

	// Enforce minimum concurrent jobs
	if c.MaxConcurrentJobs < 1 {
		c.MaxConcurrentJobs = 1
	}

	return nil
}
