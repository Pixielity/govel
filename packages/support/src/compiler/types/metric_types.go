// Package types provides supporting types for the GoVel compiler system.
//
// This file defines cache entries, metrics, and validation result types
// that support the main compiler functionality with comprehensive documentation
// and detailed field descriptions following Go best practices.
package types

import (
	"sync"
	"time"
)

// Metrics represents compilation statistics and performance metrics with thread safety.
// Provides comprehensive tracking of compiler performance and operational statistics.
type Metrics struct {
	// mu provides thread-safe access to metrics fields.
	mu sync.RWMutex `json:"-"`

	// TotalCompilations is the total number of compilation attempts.
	TotalCompilations int64 `json:"total_compilations"`

	// SuccessfulCompilations is the number of successful compilations.
	SuccessfulCompilations int64 `json:"successful_compilations"`

	// FailedCompilations is the number of failed compilations.
	FailedCompilations int64 `json:"failed_compilations"`

	// CacheHits is the number of times results were served from cache.
	CacheHits int64 `json:"cache_hits"`

	// CacheMisses is the number of times cache lookup failed.
	CacheMisses int64 `json:"cache_misses"`

	// AverageCompileTime is the average time taken for compilation.
	AverageCompileTime time.Duration `json:"average_compile_time"`

	// AverageExecutionTime is the average time taken for execution.
	AverageExecutionTime time.Duration `json:"average_execution_time"`

	// TotalCompileTime is the cumulative time spent on compilation.
	TotalCompileTime time.Duration `json:"total_compile_time"`

	// TotalExecutionTime is the cumulative time spent on execution.
	TotalExecutionTime time.Duration `json:"total_execution_time"`

	// PeakMemoryUsage is the highest memory usage recorded.
	PeakMemoryUsage int64 `json:"peak_memory_usage"`

	// AverageMemoryUsage is the average memory usage.
	AverageMemoryUsage int64 `json:"average_memory_usage"`

	// LastCompilationTime indicates when the last compilation occurred.
	LastCompilationTime time.Time `json:"last_compilation_time"`

	// StartTime indicates when metrics collection started.
	StartTime time.Time `json:"start_time"`
}

// NewMetrics creates a new Metrics instance with current start time.
// Initializes all counters to zero and sets start time to now.
//
// Returns:
//
//	*Metrics: A new metrics instance ready for use
func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

// RecordCompilation records metrics for a compilation operation with thread safety.
// Updates all relevant metrics including success/failure counts, timing, and memory usage.
//
// Parameters:
//
//	result: The compilation result to record metrics for
func (m *Metrics) RecordCompilation(result *Result) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalCompilations++
	m.LastCompilationTime = time.Now()

	// Update success/failure counts
	if result.Success {
		m.SuccessfulCompilations++
	} else {
		m.FailedCompilations++
	}

	// Update cache statistics
	if result.CacheHit {
		m.CacheHits++
	} else {
		m.CacheMisses++
	}

	// Update timing metrics
	m.TotalCompileTime += result.CompileTime
	m.TotalExecutionTime += result.ExecutionTime

	// Recalculate averages
	if m.TotalCompilations > 0 {
		m.AverageCompileTime = m.TotalCompileTime / time.Duration(m.TotalCompilations)
		m.AverageExecutionTime = m.TotalExecutionTime / time.Duration(m.TotalCompilations)
	}

	// Update memory metrics
	if result.MemoryUsed > m.PeakMemoryUsage {
		m.PeakMemoryUsage = result.MemoryUsed
	}

	// Update rolling average for memory usage
	if m.TotalCompilations > 0 {
		m.AverageMemoryUsage = (m.AverageMemoryUsage*(m.TotalCompilations-1) + result.MemoryUsed) / m.TotalCompilations
	}
}

// GetSuccessRate returns the success rate as a percentage (0-100).
// Calculates percentage of successful compilations out of total attempts.
//
// Returns:
//
//	float64: Success rate as a percentage (0.0 to 100.0)
func (m *Metrics) GetSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.TotalCompilations == 0 {
		return 0.0
	}
	return float64(m.SuccessfulCompilations) / float64(m.TotalCompilations) * 100.0
}

// GetCacheHitRate returns the cache hit rate as a percentage (0-100).
// Calculates percentage of cache hits out of all cache lookup attempts.
//
// Returns:
//
//	float64: Cache hit rate as a percentage (0.0 to 100.0)
func (m *Metrics) GetCacheHitRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := m.CacheHits + m.CacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(m.CacheHits) / float64(total) * 100.0
}

// GetUptime returns how long the compiler has been running.
// Calculates duration since metrics collection started.
//
// Returns:
//
//	time.Duration: Time since metrics collection started
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.StartTime)
}

// Reset resets all metrics to zero values with thread safety.
// Reinitializes all counters and timestamps for testing or period boundaries.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	*m = Metrics{StartTime: time.Now()}
}
