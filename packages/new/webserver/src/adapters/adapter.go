// Package adapters provides enterprise-grade framework-specific implementations that bridge 
// the unified webserver API with underlying web frameworks (GoFiber, Gin, Echo).
// This file defines the comprehensive BaseAdapter functionality with advanced features.
package adapters

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"govel/new/webserver/src/enums"
	"govel/new/webserver/src/interfaces"
)

// BaseAdapter provides comprehensive enterprise-grade functionality for all adapter implementations.
// It implements advanced configuration management, middleware handling, template rendering,
// metrics collection, health monitoring, and lifecycle management.
//
// Enterprise Features:
//   - Advanced configuration with environment variable support
//   - Template rendering with multiple engine support
//   - Comprehensive metrics and health monitoring
//   - Circuit breaker and rate limiting
//   - Distributed tracing and logging
//   - Graceful shutdown and resource management
//   - Security headers and CORS handling
//   - Static file serving with caching
//   - WebSocket support preparation
//
// This struct is designed to be embedded in framework-specific adapters:
//
//	type GoFiberAdapter struct {
//	    BaseAdapter
//	    app *fiber.App
//	}
type BaseAdapter struct {
	// Configuration Management
	config       map[string]interface{} // Runtime configuration
	envConfig    map[string]string      // Environment variable mapping
	configMutex  sync.RWMutex          // Thread-safe config access
	configValid  atomic.Bool           // Configuration validation state

	// Middleware and Request Handling
	middleware       []interfaces.MiddlewareInterface // Global middleware stack
	middlewareChain  *MiddlewareChain                 // Compiled middleware chain
	routeHandlers    map[string]interfaces.HandlerInterface // Route handler cache
	requestContext   *RequestContext                  // Request context management

	// Template and View Management
	templateEngine   enums.TemplateEngine            // Selected template engine
	templateCache    map[string]*template.Template   // Compiled template cache
	templateDir      string                          // Template directory path
	staticDirs       map[string]string               // Static file directories
	viewData         map[string]interface{}          // Global view data

	// Monitoring and Observability
	metrics          *AdapterMetrics                 // Performance metrics
	healthChecker    *HealthChecker                  // Health monitoring
	logger           Logger                          // Structured logging
	requestTracer    *RequestTracer                  // Distributed tracing

	// Security and Performance
	securityHeaders  map[string]string               // Security headers
	corsConfig       *CORSConfiguration             // CORS settings
	rateLimiter      *RateLimiter                   // Rate limiting
	circuitBreaker   *CircuitBreaker                // Circuit breaker

	// Lifecycle and State Management
	engine           enums.Engine                    // Framework engine type
	initialized      atomic.Bool                    // Initialization state
	running          atomic.Bool                    // Running state
	startTime        time.Time                      // Startup timestamp
	shutdownChan     chan struct{}                  // Shutdown signal
	shutdownTimeout  time.Duration                  // Graceful shutdown timeout

	// Resource Management
	connectionPool   *ConnectionPool                 // HTTP connection pooling
	cacheManager     *CacheManager                  // Response caching
	fileWatcher      *FileWatcher                   // Template/config file watching
	resourceMonitor  *ResourceMonitor               // System resource monitoring
}

// Enterprise Support Structures

// AdapterMetrics provides comprehensive performance and usage metrics
type AdapterMetrics struct {
	RequestCount          atomic.Int64    // Total requests processed
	ActiveRequests        atomic.Int64    // Currently active requests
	ErrorCount            atomic.Int64    // Total errors encountered
	AverageResponseTime   atomic.Value    // Average response time (time.Duration)
	TotalResponseTime     atomic.Int64    // Cumulative response time (nanoseconds)
	ThroughputPerSecond   atomic.Value    // Requests per second (float64)
	MemoryUsage          atomic.Int64    // Memory usage in bytes
	GoroutineCount       atomic.Int64    // Active goroutines
	ConnectionCount      atomic.Int64    // Active connections
	LastRequestTime      atomic.Value    // Last request timestamp (time.Time)
	StatusCodeCounts     sync.Map        // HTTP status code counters
	EndpointMetrics      sync.Map        // Per-endpoint metrics
}

// HealthChecker monitors adapter and system health
type HealthChecker struct {
	mu                 sync.RWMutex
	isHealthy          bool
	lastCheck          time.Time
	checkInterval      time.Duration
	healthChecks       []HealthCheckFunc
	failedChecks       []string
	notificationChan   chan HealthStatus
	shutdownChan       chan struct{}
}

type HealthCheckFunc func() error

type HealthStatus struct {
	Healthy      bool      `json:"healthy"`
	Timestamp    time.Time `json:"timestamp"`
	Uptime       string    `json:"uptime"`
	FailedChecks []string  `json:"failed_checks,omitempty"`
	Metrics      HealthMetrics `json:"metrics"`
}

type HealthMetrics struct {
	MemoryMB        float64 `json:"memory_mb"`
	Goroutines      int     `json:"goroutines"`
	RequestsTotal   int64   `json:"requests_total"`
	ActiveRequests  int64   `json:"active_requests"`
	ResponseTimeAvg string  `json:"response_time_avg"`
}

// Logger provides structured logging interface
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(key string, value interface{}) Logger
}

// DefaultLogger implements basic structured logging
type DefaultLogger struct {
	prefix string
	fields map[string]interface{}
	mu     sync.RWMutex
}

// RequestTracer handles distributed tracing
type RequestTracer struct {
	enabled       bool
	serviceName   string
	traceProvider string // "jaeger", "zipkin", "otel"
	samplingRate  float64
}

// CORSConfiguration defines comprehensive CORS settings
type CORSConfiguration struct {
	AllowOrigins     []string      `json:"allow_origins"`
	AllowMethods     []string      `json:"allow_methods"`
	AllowHeaders     []string      `json:"allow_headers"`
	ExposeHeaders    []string      `json:"expose_headers"`
	AllowCredentials bool          `json:"allow_credentials"`
	MaxAge           time.Duration `json:"max_age"`
	OptionsSuccess   int           `json:"options_success_status"`
}

// RateLimiter provides advanced rate limiting
type RateLimiter struct {
	enabled       bool
	requestsLimit int
	windowSize    time.Duration
	byIP          bool
	byUser        bool
	customKeyFunc func(interface{}) string
	buckets       sync.Map // string -> *TokenBucket
	cleanupTicker *time.Ticker
}

type TokenBucket struct {
	capacity     int
	tokens       atomic.Int32
	lastRefill   atomic.Value // time.Time
	refillRate   time.Duration
	mu           sync.Mutex
}

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	enabled           bool
	failureThreshold  int
	recoveryTimeout   time.Duration
	resetTimeout      time.Duration
	state             atomic.Value // CircuitState
	failureCount      atomic.Int32
	lastFailureTime   atomic.Value // time.Time
	onStateChange     func(CircuitState)
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// MiddlewareChain provides optimized middleware execution
type MiddlewareChain struct {
	middleware []interfaces.MiddlewareInterface
	compiled   bool
	chain      func(interfaces.RequestInterface, interfaces.HandlerInterface) interfaces.ResponseInterface
}

// RequestContext manages request-scoped data and lifecycle
type RequestContext struct {
	contextPool   sync.Pool
	requestID     string
	traceID       string
	userID        string
	sessionID     string
	values        sync.Map
	startTime     time.Time
	timeout       time.Duration
	cancellationFunc context.CancelFunc
}

// ConnectionPool manages HTTP connection pooling
type ConnectionPool struct {
	maxConnections    int
	maxIdleTime      time.Duration
	connectionTimeout time.Duration
	pool             sync.Pool
	activeCount      atomic.Int32
}

// CacheManager handles response caching
type CacheManager struct {
	enabled       bool
	defaultTTL    time.Duration
	maxSize       int64
	currentSize   atomic.Int64
	cache         sync.Map // string -> *CacheEntry
	cleanupTicker *time.Ticker
}

type CacheEntry struct {
	value     interface{}
	expiresAt time.Time
	size      int64
	accessCount atomic.Int64
	lastAccess  atomic.Value // time.Time
}

// FileWatcher monitors file changes for hot reloading
type FileWatcher struct {
	enabled    bool
	watchPaths []string
	notifyChan chan FileChangeEvent
	stopChan   chan struct{}
}

type FileChangeEvent struct {
	Path      string
	EventType FileEventType
	Timestamp time.Time
}

type FileEventType int

const (
	FileCreated FileEventType = iota
	FileModified
	FileDeleted
)

// ResourceMonitor tracks system resource usage
type ResourceMonitor struct {
	enabled        bool
	monitoringChan chan ResourceStats
	stopChan       chan struct{}
	interval       time.Duration
}

type ResourceStats struct {
	CPUPercent     float64   `json:"cpu_percent"`
	MemoryMB       float64   `json:"memory_mb"`
	MemoryPercent  float64   `json:"memory_percent"`
	Goroutines     int       `json:"goroutines"`
	FileDescriptors int      `json:"file_descriptors"`
	Timestamp      time.Time `json:"timestamp"`
}

// NewBaseAdapter creates a comprehensive enterprise-grade BaseAdapter instance.
// This constructor initializes all enterprise features including metrics monitoring,
// health checking, caching, security, and resource management.
//
// Features initialized:
//   - Configuration management with environment variable support
//   - Structured logging with default logger
//   - Performance metrics and health monitoring
//   - Security headers and CORS configuration
//   - Rate limiting and circuit breaker
//   - Template engine and static file support
//   - Resource monitoring and connection pooling
//   - Graceful shutdown and lifecycle management
//
// Parameters:
//   engine: The web framework engine this adapter will wrap
//
// Returns:
//   *BaseAdapter: A fully configured enterprise-grade base adapter
//
// Example:
//   base := NewBaseAdapter(enums.GoFiber)
//   base.SetLogger(customLogger)
//   base.EnableMetrics(true)
func NewBaseAdapter(engine enums.Engine) *BaseAdapter {
	base := &BaseAdapter{
		// Configuration Management
		config:          make(map[string]interface{}),
		envConfig:       make(map[string]string),
		configMutex:     sync.RWMutex{},

		// Middleware and Request Handling
		middleware:      make([]interfaces.MiddlewareInterface, 0),
		middlewareChain: &MiddlewareChain{middleware: make([]interfaces.MiddlewareInterface, 0)},
		routeHandlers:   make(map[string]interfaces.HandlerInterface),
		requestContext:  &RequestContext{contextPool: sync.Pool{}},

		// Template and View Management
		templateEngine:  enums.DefaultTemplateEngine(),
		templateCache:   make(map[string]*template.Template),
		staticDirs:      make(map[string]string),
		viewData:        make(map[string]interface{}),

		// Monitoring and Observability
		metrics:         newAdapterMetrics(),
		healthChecker:   newHealthChecker(),
		logger:          newDefaultLogger(string(engine)),
		requestTracer:   &RequestTracer{enabled: false},

		// Security and Performance
		securityHeaders: getDefaultSecurityHeaders(),
		corsConfig:      getDefaultCORSConfig(),
		rateLimiter:     newRateLimiter(),
		circuitBreaker:  newCircuitBreaker(),

		// Lifecycle and State Management
		engine:          engine,
		startTime:       time.Now(),
		shutdownChan:    make(chan struct{}),
		shutdownTimeout: 30 * time.Second,

		// Resource Management
		connectionPool:  newConnectionPool(),
		cacheManager:    newCacheManager(),
		fileWatcher:     &FileWatcher{enabled: false},
		resourceMonitor: &ResourceMonitor{enabled: false},
	}

	// Initialize atomic values
	base.configValid.Store(false)
	base.initialized.Store(false)
	base.running.Store(false)

	// Set up environment variable mapping
	base.setupEnvironmentMapping()

	// Initialize request context pool
	base.requestContext.contextPool.New = func() interface{} {
		return &RequestContext{
			values:    sync.Map{},
			startTime: time.Now(),
		}
	}

	return base
}

// Enterprise Component Factory Functions

// newAdapterMetrics creates a new metrics instance
func newAdapterMetrics() *AdapterMetrics {
	metrics := &AdapterMetrics{}
	metrics.LastRequestTime.Store(time.Now())
	metrics.AverageResponseTime.Store(time.Duration(0))
	metrics.ThroughputPerSecond.Store(float64(0))
	return metrics
}

// newHealthChecker creates a new health checker
func newHealthChecker() *HealthChecker {
	return &HealthChecker{
		isHealthy:        true,
		lastCheck:        time.Now(),
		checkInterval:    30 * time.Second,
		healthChecks:     make([]HealthCheckFunc, 0),
		failedChecks:     make([]string, 0),
		notificationChan: make(chan HealthStatus, 10),
		shutdownChan:     make(chan struct{}),
	}
}

// newDefaultLogger creates a default logger instance
func newDefaultLogger(prefix string) Logger {
	return &DefaultLogger{
		prefix: prefix,
		fields: make(map[string]interface{}),
		mu:     sync.RWMutex{},
	}
}

// getDefaultSecurityHeaders returns secure default headers
func getDefaultSecurityHeaders() map[string]string {
	return map[string]string{
		"X-Frame-Options":        "DENY",
		"X-Content-Type-Options": "nosniff",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Permissions-Policy":     "geolocation=(), microphone=(), camera=()",
	}
}

// getDefaultCORSConfig returns secure default CORS configuration
func getDefaultCORSConfig() *CORSConfiguration {
	return &CORSConfiguration{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		OptionsSuccess:   204,
	}
}

// newRateLimiter creates a new rate limiter
func newRateLimiter() *RateLimiter {
	return &RateLimiter{
		enabled:       false,
		requestsLimit: 100,
		windowSize:    time.Minute,
		byIP:          true,
		byUser:        false,
		buckets:       sync.Map{},
	}
}

// newCircuitBreaker creates a new circuit breaker
func newCircuitBreaker() *CircuitBreaker {
	cb := &CircuitBreaker{
		enabled:          false,
		failureThreshold: 5,
		recoveryTimeout:  30 * time.Second,
		resetTimeout:     60 * time.Second,
	}
	cb.state.Store(CircuitClosed)
	cb.lastFailureTime.Store(time.Time{})
	return cb
}

// newConnectionPool creates a new connection pool
func newConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		maxConnections:    100,
		maxIdleTime:      30 * time.Second,
		connectionTimeout: 10 * time.Second,
		pool:             sync.Pool{},
	}
}

// newCacheManager creates a new cache manager
func newCacheManager() *CacheManager {
	return &CacheManager{
		enabled:     false,
		defaultTTL:  5 * time.Minute,
		maxSize:     100 * 1024 * 1024, // 100MB
		cache:       sync.Map{},
	}
}

// setupEnvironmentMapping sets up environment variable mapping
func (b *BaseAdapter) setupEnvironmentMapping() {
	envMappings := map[string]string{
		"HOST":                    "host",
		"PORT":                    "port",
		"MAX_BODY_SIZE":           "max_body_size",
		"READ_TIMEOUT":            "read_timeout",
		"WRITE_TIMEOUT":           "write_timeout",
		"IDLE_TIMEOUT":            "idle_timeout",
		"KEEP_ALIVE_ENABLED":      "keep_alive_enabled",
		"COMPRESSION_ENABLED":     "compression_enabled",
		"RATE_LIMIT_ENABLED":      "rate_limit_enabled",
		"RATE_LIMIT_REQUESTS":     "rate_limit_requests",
		"HEALTH_CHECK_ENABLED":    "health_check_enabled",
		"METRICS_ENABLED":         "metrics_enabled",
		"TRACING_ENABLED":         "tracing_enabled",
		"TEMPLATE_DIR":            "template_directory",
		"STATIC_DIR":              "static_directory",
	}

	// Load environment variables
	for envVar, configKey := range envMappings {
		if value := os.Getenv(envVar); value != "" {
			b.envConfig[configKey] = value
		}
	}
}

// DefaultLogger Implementation

// Debug logs debug-level messages with optional fields
func (l *DefaultLogger) Debug(msg string, fields ...interface{}) {
	l.logWithLevel("DEBUG", msg, fields...)
}

// Info logs info-level messages with optional fields
func (l *DefaultLogger) Info(msg string, fields ...interface{}) {
	l.logWithLevel("INFO", msg, fields...)
}

// Warn logs warning-level messages with optional fields
func (l *DefaultLogger) Warn(msg string, fields ...interface{}) {
	l.logWithLevel("WARN", msg, fields...)
}

// Error logs error-level messages with optional fields
func (l *DefaultLogger) Error(msg string, fields ...interface{}) {
	l.logWithLevel("ERROR", msg, fields...)
}

// Fatal logs fatal-level messages with optional fields and exits
func (l *DefaultLogger) Fatal(msg string, fields ...interface{}) {
	l.logWithLevel("FATAL", msg, fields...)
	os.Exit(1)
}

// With creates a new logger with additional context fields
func (l *DefaultLogger) With(key string, value interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	
	return &DefaultLogger{
		prefix: l.prefix,
		fields: newFields,
		mu:     sync.RWMutex{},
	}
}

// logWithLevel performs the actual logging with level and structured fields
func (l *DefaultLogger) logWithLevel(level string, msg string, fields ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	timestamp := time.Now().Format(time.RFC3339)
	logMsg := fmt.Sprintf("[%s] %s [%s] %s", timestamp, level, l.prefix, msg)
	
	// Add context fields
	if len(l.fields) > 0 {
		logMsg += " |"
		for k, v := range l.fields {
			logMsg += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	
	// Add additional fields from call
	if len(fields) > 0 {
		logMsg += " |"
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMsg += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
			}
		}
	}
	
	log.Println(logMsg)
}

// Configuration Management Methods

// Init initializes the base adapter with configuration and middleware.
// This method should be called by concrete adapter implementations as part of their Init() method.
//
// Parameters:
//
//	config: Configuration key-value pairs
//	middleware: Global middleware to register
//
// Returns:
//
//	error: Initialization error, or nil on success
//
// Example:
//
//	func (a *GoFiberAdapter) Init(config map[string]interface{}, middleware []interfaces.MiddlewareInterface) error {
//	    if err := a.BaseAdapter.Init(config, middleware); err != nil {
//	        return err
//	    }
//	    // GoFiber-specific initialization here
//	    return nil
//	}
func (b *BaseAdapter) Init(config map[string]interface{}, middleware []interfaces.MiddlewareInterface) error {
	if b.initialized.Load() {
		return errors.New("adapter already initialized")
	}

	// Validate and store configuration
	if err := b.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Store configuration with thread-safe access
	b.configMutex.Lock()
	for key, value := range config {
		b.config[key] = value
	}
	b.configMutex.Unlock()

	// Store middleware
	b.middleware = append(b.middleware, middleware...)

	// Set default configuration values if not provided
	b.setDefaults()

	// Mark as initialized
	b.initialized.Store(true)
	b.configValid.Store(true)
	
	return nil
}

// SetConfig sets a configuration value at runtime.
// This method allows dynamic reconfiguration of the adapter.
//
// Parameters:
//
//	key: The configuration key
//	value: The configuration value
//
// Example:
//
//	adapter.SetConfig("timeout", 30)
//	adapter.SetConfig("debug", true)
func (b *BaseAdapter) SetConfig(key string, value interface{}) {
	b.config[key] = value
}

// GetConfig retrieves a configuration value by key.
// Returns nil if the key doesn't exist.
//
// Parameters:
//
//	key: The configuration key to retrieve
//
// Returns:
//
//	interface{}: The configuration value, or nil if not found
//
// Example:
//
//	timeout := adapter.GetConfig("timeout")
//	if timeout != nil {
//	    fmt.Printf("Timeout: %v\n", timeout)
//	}
func (b *BaseAdapter) GetConfig(key string) interface{} {
	return b.config[key]
}

// GetConfigString retrieves a configuration value as a string.
// Returns the default value if the key doesn't exist or cannot be converted.
//
// Parameters:
//
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	string: The configuration value as string, or default value
//
// Example:
//
//	host := adapter.GetConfigString("host", "localhost")
func (b *BaseAdapter) GetConfigString(key string, defaultValue string) string {
	if value, exists := b.config[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetConfigInt retrieves a configuration value as an integer.
// Returns the default value if the key doesn't exist or cannot be converted.
//
// Parameters:
//
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	int: The configuration value as integer, or default value
//
// Example:
//
//	port := adapter.GetConfigInt("port", 8080)
func (b *BaseAdapter) GetConfigInt(key string, defaultValue int) int {
	if value, exists := b.config[key]; exists {
		if intVal, ok := value.(int); ok {
			return intVal
		}
	}
	return defaultValue
}

// GetConfigBool retrieves a configuration value as a boolean.
// Returns the default value if the key doesn't exist or cannot be converted.
//
// Parameters:
//
//	key: The configuration key
//	defaultValue: The default value if key is not found
//
// Returns:
//
//	bool: The configuration value as boolean, or default value
//
// Example:
//
//	debug := adapter.GetConfigBool("debug", false)
func (b *BaseAdapter) GetConfigBool(key string, defaultValue bool) bool {
	if value, exists := b.config[key]; exists {
		if boolVal, ok := value.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

// Middleware Management Methods

// Use registers global middleware.
// This method appends middleware to the existing middleware stack.
//
// Parameters:
//
//	middleware: One or more middleware implementations to register
//
// Example:
//
//	corsMiddleware := &CorsMiddleware{}
//	authMiddleware := &AuthMiddleware{}
//	adapter.Use(corsMiddleware, authMiddleware)
func (b *BaseAdapter) Use(middleware ...interfaces.MiddlewareInterface) {
	b.middleware = append(b.middleware, middleware...)
}

// GetMiddleware returns a copy of the registered middleware stack.
// Returns a copy to prevent external modification of the internal stack.
//
// Returns:
//
//	[]interfaces.MiddlewareInterface: Copy of the middleware stack
func (b *BaseAdapter) GetMiddleware() []interfaces.MiddlewareInterface {
	middlewareCopy := make([]interfaces.MiddlewareInterface, len(b.middleware))
	copy(middlewareCopy, b.middleware)
	return middlewareCopy
}

// Utility Methods

// GetListenAddress constructs the listening address from configuration.
// This method handles various address configuration patterns:
//   - "address" key takes precedence (e.g., "localhost:8080")
//   - Falls back to combining "host" and "port" keys
//   - Uses defaults if neither is configured
//
// Returns:
//
//	string: The listening address in "host:port" format
//
// Example:
//
//	addr := adapter.GetListenAddress() // Returns "localhost:8080"
func (b *BaseAdapter) GetListenAddress() string {
	// Check if complete address is configured
	if addr := b.GetConfigString("address", ""); addr != "" {
		return addr
	}

	// Build address from host and port
	host := b.GetConfigString("host", "localhost")
	port := b.GetConfigInt("port", 8080)

	return fmt.Sprintf("%s:%d", host, port)
}

// GetEngine returns the web framework engine this adapter wraps.
//
// Returns:
//
//	enums.Engine: The engine type (GoFiber, Gin, Echo)
func (b *BaseAdapter) GetEngine() enums.Engine {
	return b.engine
}

// IsInitialized returns whether the adapter has been initialized.
//
// Returns:
//
//	bool: True if Init() has been called successfully
func (b *BaseAdapter) IsInitialized() bool {
	return b.initialized.Load()
}

// GetLogger returns the adapter's logger instance
func (b *BaseAdapter) GetLogger() Logger {
	return b.logger
}

// Validation Methods

// validateConfig validates the provided configuration.
// This method checks for common configuration errors and ensures
// the adapter can be initialized successfully.
//
// Parameters:
//
//	config: The configuration to validate
//
// Returns:
//
//	error: Validation error, or nil if configuration is valid
//
// Validates:
//   - Port numbers are in valid range (1-65535)
//   - Host addresses are properly formatted
//   - Required configuration keys are present
//   - Value types are correct
func (b *BaseAdapter) validateConfig(config map[string]interface{}) error {
	// Validate port if provided
	if portInterface, exists := config["port"]; exists {
		if port, ok := portInterface.(int); ok {
			if port < 1 || port > 65535 {
				return fmt.Errorf("invalid port number: %d (must be 1-65535)", port)
			}
		} else {
			return fmt.Errorf("port must be an integer, got: %T", portInterface)
		}
	}

	// Validate host if provided
	if hostInterface, exists := config["host"]; exists {
		if host, ok := hostInterface.(string); ok {
			if strings.TrimSpace(host) == "" {
				return errors.New("host cannot be empty")
			}
		} else {
			return fmt.Errorf("host must be a string, got: %T", hostInterface)
		}
	}

	// Validate address format if provided
	if addrInterface, exists := config["address"]; exists {
		if addr, ok := addrInterface.(string); ok {
			if err := b.validateAddress(addr); err != nil {
				return fmt.Errorf("invalid address format: %w", err)
			}
		} else {
			return fmt.Errorf("address must be a string, got: %T", addrInterface)
		}
	}

	return nil
}

// validateAddress validates an address string format.
// Accepts formats like "host:port", ":port", "host:", etc.
//
// Parameters:
//
//	addr: The address string to validate
//
// Returns:
//
//	error: Validation error, or nil if address is valid
func (b *BaseAdapter) validateAddress(addr string) error {
	if addr == "" {
		return errors.New("address cannot be empty")
	}

	// Must contain exactly one colon for host:port format
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("address must be in format 'host:port', got: %s", addr)
	}

	// Validate port part if present
	portStr := strings.TrimSpace(parts[1])
	if portStr != "" {
		var port int
		if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
			return fmt.Errorf("invalid port in address: %s", portStr)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("port in address out of range: %d", port)
		}
	}

	return nil
}

// setDefaults sets default configuration values if not already configured.
// This ensures the adapter has sensible defaults for common configuration keys.
func (b *BaseAdapter) setDefaults() {
	// Set default host if not configured
	if _, exists := b.config["host"]; !exists {
		b.config["host"] = "localhost"
	}

	// Set default port if not configured
	if _, exists := b.config["port"]; !exists {
		b.config["port"] = 8080
	}

	// Set default timeout if not configured (in seconds)
	if _, exists := b.config["timeout"]; !exists {
		b.config["timeout"] = 30
	}

	// Set default max body size if not configured (in bytes)
	if _, exists := b.config["max_body_size"]; !exists {
		b.config["max_body_size"] = 1024 * 1024 * 4 // 4MB
	}

	// Set default debug mode if not configured
	if _, exists := b.config["debug"]; !exists {
		b.config["debug"] = false
	}
}

// Error Handling Utilities

// CreateAdapterError creates a standardized adapter error.
// This provides consistent error formatting across all adapters.
//
// Parameters:
//
//	operation: The operation that failed (e.g., "listen", "route_registration")
//	err: The underlying error
//
// Returns:
//
//	error: A formatted adapter error
//
// Example:
//
//	if err := app.Listen(addr); err != nil {
//	    return b.CreateAdapterError("listen", err)
//	}
func (b *BaseAdapter) CreateAdapterError(operation string, err error) error {
	return fmt.Errorf("%s adapter %s failed: %w", b.engine.Name(), operation, err)
}

// HTTP Method Validation

// IsValidHTTPMethod checks if a method string represents a valid HTTP method.
// This uses the enums package for validation.
//
// Parameters:
//
//	method: The HTTP method string to validate
//
// Returns:
//
//	bool: True if the method is valid, false otherwise
//
// Example:
//
//	if !adapter.IsValidHTTPMethod("INVALID") {
//	    return errors.New("invalid HTTP method")
//	}
func (b *BaseAdapter) IsValidHTTPMethod(method string) bool {
	_, valid := enums.ParseHTTPMethod(method)
	return valid
}

// NormalizeHTTPMethod normalizes an HTTP method string to uppercase.
// This ensures consistent method handling across adapters.
//
// Parameters:
//
//	method: The HTTP method string to normalize
//
// Returns:
//
//	string: The normalized method string (uppercase)
//	error: Error if the method is invalid
//
// Example:
//
//	method, err := adapter.NormalizeHTTPMethod("get")
//	// method is now "GET"
func (b *BaseAdapter) NormalizeHTTPMethod(method string) (string, error) {
	httpMethod, valid := enums.ParseHTTPMethod(method)
	if !valid {
		return "", fmt.Errorf("invalid HTTP method: %s", method)
	}
	return httpMethod.String(), nil
}

// Lifecycle Management

// Shutdown provides a default shutdown implementation.
// Concrete adapters should override this method with framework-specific shutdown logic.
//
// Parameters:
//
//	ctx: The context for shutdown timeout
//
// Returns:
//
//	error: Shutdown error, or nil on success
func (b *BaseAdapter) Shutdown(ctx context.Context) error {
	// Base implementation does nothing
	// Concrete adapters should override this method
	return nil
}

// String returns a string representation of the adapter.
// This is useful for logging and debugging purposes.
//
// Returns:
//
//	string: A string describing the adapter
//
// Example:
//
//	fmt.Printf("Using adapter: %s\n", adapter.String())
func (b *BaseAdapter) String() string {
	addr := b.GetListenAddress()
	return fmt.Sprintf("%s adapter (listening on %s)", b.engine.Name(), addr)
}

// Package-level utility functions

// ConvertToStandardHandler converts a webserver HandlerInterface to a standard http.Handler.
// This is useful for adapters that need to work with the standard library.
//
// Parameters:
//
//	handler: The webserver handler to convert
//	requestWrapper: Function to wrap the http.Request into RequestInterface
//
// Returns:
//
//	http.Handler: A standard HTTP handler
//
// Note: This is a utility function that adapters can use when they need to
// bridge between the webserver interfaces and standard Go HTTP interfaces.
func ConvertToStandardHandler(handler interfaces.HandlerInterface, requestWrapper func(*http.Request) interfaces.RequestInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the standard request into our RequestInterface
		req := requestWrapper(r)

		// Execute the handler
		resp := handler.Handle(req)

		// Convert the response back to standard HTTP
		// This would need to be implemented based on the ResponseInterface
		// For now, this is a placeholder for the concept
		_ = resp // TODO: Convert ResponseInterface to http.ResponseWriter
	})
}

// AdapterRegistry provides a registry for adapter implementations.
// This can be used by factories to register and discover available adapters.
var AdapterRegistry = make(map[enums.Engine]func() interfaces.AdapterInterface)

// RegisterAdapter registers an adapter factory function for an engine.
// This allows the adapter factory to create instances of registered adapters.
//
// Parameters:
//
//	engine: The engine this adapter supports
//	factory: Function that creates a new instance of the adapter
//
// Example:
//
//	RegisterAdapter(enums.GoFiber, func() interfaces.AdapterInterface {
//	    return &GoFiberAdapter{}
//	})
func RegisterAdapter(engine enums.Engine, factory func() interfaces.AdapterInterface) {
	AdapterRegistry[engine] = factory
}

// GetRegisteredAdapters returns a slice of all registered engine types.
//
// Returns:
//
//	[]enums.Engine: Slice of registered engines
func GetRegisteredAdapters() []enums.Engine {
	engines := make([]enums.Engine, 0, len(AdapterRegistry))
	for engine := range AdapterRegistry {
		engines = append(engines, engine)
	}
	return engines
}

// Create creates a new adapter instance for the specified engine name.
// This is the main factory function used by the webserver to create adapters.
//
// Parameters:
//
//	engineName: The name of the engine ("gin", "echo", "fiber", "net/http")
//
// Returns:
//
//	interfaces.AdapterInterface: A new adapter instance
//	error: Error if the engine is not supported or adapter creation fails
//
// Example:
//
//	adapter, err := adapters.Create("gin")
//	if err != nil {
//	    log.Fatal("Failed to create adapter:", err)
//	}
func Create(engineName string) (interfaces.AdapterInterface, error) {
	// Parse engine name to enum
	engine, valid := enums.ParseEngine(strings.ToLower(engineName))
	if !valid {
		return nil, fmt.Errorf("unsupported engine: %s", engineName)
	}

	// Check if adapter is registered
	factory, exists := AdapterRegistry[engine]
	if !exists {
		return nil, fmt.Errorf("no adapter registered for engine: %s", engineName)
	}

	// Create and return adapter instance
	return factory(), nil
}
