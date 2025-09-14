package concerns

import (
	"runtime"
	"sync"
	"time"

	"govel/application/helpers"
	concernsInterfaces "govel/types/interfaces/application/concerns"
)

/**
 * HasInfo provides comprehensive application information functionality.
 * This trait implements the HasInfoInterface and manages collection and
 * presentation of detailed application metadata, system information, and
 * runtime statistics.
 *
 * Features:
 * - Comprehensive application information collection
 * - Runtime system information (Go version, OS, architecture)
 * - Performance metrics and memory statistics
 * - Thread-safe access to information data
 * - Configurable information gathering scope
 */
type HasInfo struct {
	/**
	 * additionalInfo holds custom application information
	 */
	additionalInfo map[string]interface{}

	/**
	 * mutex provides thread-safe access to info fields
	 */
	mutex sync.RWMutex

	/**
	 * envHelper provides environment variable access
	 */
	envHelper *helpers.EnvHelper
}

// NewInfo creates a new info trait with optional additional information.
//
// Parameters:
//
//	additionalInfo: Optional map of additional information to include
//
// Returns:
//
//	*HasInfo: A new info trait instance
//
// Example:
//
//	// Basic info
//	info := NewInfo()
//	// With additional information
//	info := NewInfo(map[string]interface{}{
//	    "build_number": "12345",
//	    "commit_hash": "abc123def",
//	})
func NewInfo(additionalInfo ...map[string]interface{}) *HasInfo {
	envHelper := helpers.NewEnvHelper()

	// Initialize additional info
	additional := make(map[string]interface{})
	if len(additionalInfo) > 0 && additionalInfo[0] != nil {
		additional = additionalInfo[0]
	}

	return &HasInfo{
		additionalInfo: additional,
		envHelper:      envHelper,
	}
}

// GetApplicationInfo returns comprehensive application information.
// This method collects and returns a detailed map of application metadata,
// system information, runtime statistics, and custom information.
//
// Returns:
//
//	map[string]interface{}: Comprehensive application information
//
// Example:
//
//	info := app.GetApplicationInfo()
//	fmt.Printf("App: %s v%s\n", info["name"], info["version"])
//	fmt.Printf("Environment: %s\n", info["environment"])
//	fmt.Printf("Go Version: %s\n", info["go_version"])
func (i *HasInfo) GetApplicationInfo() map[string]interface{} {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// Collect runtime memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Build comprehensive information map
	info := map[string]interface{}{
		// Basic application information
		"name":        i.envHelper.GetAppName(),
		"version":     i.envHelper.GetAppVersion(),
		"environment": i.envHelper.GetAppEnvironment(),
		"debug":       i.envHelper.GetAppDebug(),

		// Runtime information
		"go_version":     runtime.Version(),
		"go_os":          runtime.GOOS,
		"go_arch":        runtime.GOARCH,
		"num_cpu":        runtime.NumCPU(),
		"num_goroutines": runtime.NumGoroutine(),

		// Memory statistics (in bytes)
		"memory": map[string]interface{}{
			"allocated":   memStats.Alloc,
			"total_alloc": memStats.TotalAlloc,
			"sys":         memStats.Sys,
			"heap_alloc":  memStats.HeapAlloc,
			"heap_sys":    memStats.HeapSys,
			"heap_inuse":  memStats.HeapInuse,
			"heap_idle":   memStats.HeapIdle,
			"stack_inuse": memStats.StackInuse,
			"stack_sys":   memStats.StackSys,
			"num_gc":      memStats.NumGC,
		},

		// Runtime state
		"runtime": map[string]interface{}{
			"running_in_console": i.envHelper.GetRunningInConsole(),
			"running_unit_tests": i.envHelper.GetRunningUnitTests(),
		},

		// Timing information (will be set if timing concern is available)
		"timing": map[string]interface{}{
			"current_time": time.Now(),
		},

		// Build and deployment information
		"build": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		},

		// System paths (will be set if directable concern is available)
		"paths": map[string]interface{}{
			"working_directory": i.getWorkingDirectory(),
		},

		// Configuration information
		"config": map[string]interface{}{
			"timezone": i.envHelper.GetAppTimezone(),
			"locale":   i.envHelper.GetAppLocale(),
		},
	}

	// Add any additional custom information
	for key, value := range i.additionalInfo {
		info[key] = value
	}

	return info
}

// SetAdditionalInfo sets custom additional information.
//
// Parameters:
//
//	key: The information key
//	value: The information value
//
// Example:
//
//	app.SetAdditionalInfo("deployment_id", "deploy-123")
//	app.SetAdditionalInfo("cluster_name", "production-east")
func (i *HasInfo) SetAdditionalInfo(key string, value interface{}) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if i.additionalInfo == nil {
		i.additionalInfo = make(map[string]interface{})
	}

	i.additionalInfo[key] = value
}

// GetAdditionalInfo returns a specific piece of additional information.
//
// Parameters:
//
//	key: The information key to retrieve
//
// Returns:
//
//	interface{}: The information value, or nil if not found
//
// Example:
//
//	deploymentId := app.GetAdditionalInfo("deployment_id")
func (i *HasInfo) GetAdditionalInfo(key string) interface{} {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if i.additionalInfo == nil {
		return nil
	}

	return i.additionalInfo[key]
}

// ClearAdditionalInfo removes all additional information.
//
// Example:
//
//	app.ClearAdditionalInfo() // Reset custom info
func (i *HasInfo) ClearAdditionalInfo() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.additionalInfo = make(map[string]interface{})
}

// getWorkingDirectory safely gets the current working directory
func (i *HasInfo) getWorkingDirectory() string {
	// This is a helper method that can be overridden if directable concern is available
	return "." // Default fallback
}

// Compile-time interface compliance check
var _ concernsInterfaces.HasInfoInterface = (*HasInfo)(nil)
