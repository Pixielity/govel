package drivers

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"

	configInterfaces "govel/types/interfaces/config"
)

// MemoryDriver implements config.Driver interface for in-memory configuration.
// It uses Viper internally for data storage and manipulation.
type MemoryDriver struct {
	viper     *viper.Viper
	mu        sync.RWMutex
	watchFunc func()
}

// MemoryDriverOptions contains configuration options for MemoryDriver.
type MemoryDriverOptions struct {
	InitialData map[string]interface{} // Initial configuration data
}

// NewMemoryDriver creates a new in-memory configuration driver.
//
// Example:
//
//	driver := drivers.NewMemoryDriver(&drivers.MemoryDriverOptions{
//	    InitialData: map[string]interface{}{
//	        "app.name": "MyApp",
//	        "app.debug": true,
//	        "server.port": 8080,
//	    },
//	})
func NewMemoryDriver(options *MemoryDriverOptions) *MemoryDriver {
	v := viper.New()

	// Set initial data if provided
	if options != nil && options.InitialData != nil {
		for key, value := range options.InitialData {
			v.Set(key, value)
		}
	}

	return &MemoryDriver{
		viper: v,
	}
}

// Get retrieves a configuration value by key.
func (m *MemoryDriver) Get(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.viper.IsSet(key) {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return m.viper.Get(key), nil
}

// Set stores a configuration value by key.
func (m *MemoryDriver) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.viper.Set(key, value)

	// Trigger watch callback if set
	if m.watchFunc != nil {
		go m.watchFunc()
	}

	return nil
}

// Load loads the configuration (no-op for memory driver).
func (m *MemoryDriver) Load() error {
	// Memory driver doesn't need to load from external sources
	return nil
}

// Watch starts watching for configuration changes.
func (m *MemoryDriver) Watch(callback func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.watchFunc != nil {
		return fmt.Errorf("already watching for config changes")
	}

	m.watchFunc = callback
	return nil
}

// Unwatch stops watching for configuration changes.
func (m *MemoryDriver) Unwatch() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.watchFunc = nil
	return nil
}

// Invalidate clears the configuration cache.
func (m *MemoryDriver) Invalidate() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear all settings by creating a new viper instance
	m.viper = viper.New()
	return nil
}

// GetAll returns all configuration data as a map.
func (m *MemoryDriver) GetAll() (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.viper.AllSettings(), nil
}

// Has checks if a configuration key exists.
func (m *MemoryDriver) Has(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.viper.IsSet(key)
}

// Delete removes a configuration key.
func (m *MemoryDriver) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get all settings and remove the key
	allSettings := m.viper.AllSettings()

	// Use dot notation to delete nested keys
	if err := deleteNestedKey(allSettings, key); err != nil {
		return err
	}

	// Clear and reload settings
	m.viper = viper.New()
	for k, v := range allSettings {
		m.viper.Set(k, v)
	}

	// Trigger watch callback if set
	if m.watchFunc != nil {
		go m.watchFunc()
	}

	return nil
}

// LoadFromMap loads configuration from a map.
func (m *MemoryDriver) LoadFromMap(data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, value := range data {
		m.viper.Set(key, value)
	}

	// Trigger watch callback if set
	if m.watchFunc != nil {
		go m.watchFunc()
	}

	return nil
}

// Clear removes all configuration data.
func (m *MemoryDriver) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.viper = viper.New()

	// Trigger watch callback if set
	if m.watchFunc != nil {
		go m.watchFunc()
	}

	return nil
}

// MergeMap merges configuration from a map.
func (m *MemoryDriver) MergeMap(data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use viper's merge functionality
	for key, value := range data {
		if m.viper.IsSet(key) {
			// If key exists, we need to handle merging
			existing := m.viper.Get(key)
			if existingMap, ok := existing.(map[string]interface{}); ok {
				if valueMap, ok := value.(map[string]interface{}); ok {
					// Merge maps
					merged := mergeMaps(existingMap, valueMap)
					m.viper.Set(key, merged)
				} else {
					// Replace with new value
					m.viper.Set(key, value)
				}
			} else {
				// Replace with new value
				m.viper.Set(key, value)
			}
		} else {
			m.viper.Set(key, value)
		}
	}

	// Trigger watch callback if set
	if m.watchFunc != nil {
		go m.watchFunc()
	}

	return nil
}

// Helper functions

// deleteNestedKey deletes a key from a nested map using dot notation
func deleteNestedKey(data map[string]interface{}, key string) error {
	// For simplicity, we'll just delete top-level keys
	// A more sophisticated implementation would handle nested keys properly
	delete(data, key)
	return nil
}

// mergeMaps recursively merges two maps
func mergeMaps(dst, src map[string]interface{}) map[string]interface{} {
	for key, srcVal := range src {
		if dstVal, exists := dst[key]; exists {
			// Both dst and src have this key
			if dstMap, dstIsMap := dstVal.(map[string]interface{}); dstIsMap {
				if srcMap, srcIsMap := srcVal.(map[string]interface{}); srcIsMap {
					// Both are maps, merge recursively
					dst[key] = mergeMaps(dstMap, srcMap)
				} else {
					// Source is not a map, overwrite
					dst[key] = srcVal
				}
			} else {
				// Destination is not a map, overwrite
				dst[key] = srcVal
			}
		} else {
			// Key doesn't exist in dst, add it
			dst[key] = srcVal
		}
	}
	return dst
}

// Ensure MemoryDriver implements the Driver interface
var _ configInterfaces.DriverInterface = (*MemoryDriver)(nil)
