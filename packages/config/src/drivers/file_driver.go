package drivers

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"govel/support/compiler"
)

// FileDriver implements config.Driver interface for file-based configuration.
// It uses Viper internally for file operations and supports multiple formats.
type FileDriver struct {
	viper       *viper.Viper
	configPaths []string
	configName  string
	configType  string
	mu          sync.RWMutex
	watcher     *fsnotify.Watcher
	watchFunc   func()
}

// FileDriverOptions contains configuration options for FileDriver.
type FileDriverOptions struct {
	ConfigPaths []string // List of paths to search for config files
	ConfigName  string   // Name of config file (without extension)
	ConfigType  string   // Type of config file (json, yaml, toml, etc.)
}

// NewFileDriver creates a new file-based configuration driver.
//
// Example:
//
//	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
//	    ConfigPaths: []string{".", "/etc/myapp", "$HOME/.myapp"},
//	    ConfigName:  "config",
//	    ConfigType:  "yaml",
//	})
func NewFileDriver(options *FileDriverOptions) *FileDriver {
	v := viper.New()

	// Set defaults
	if options == nil {
		options = &FileDriverOptions{
			ConfigPaths: []string{"."},
			ConfigName:  "config",
			ConfigType:  "yaml",
		}
	}

	// Configure viper
	v.SetConfigName(options.ConfigName)
	v.SetConfigType(options.ConfigType)

	// Add config paths
	for _, path := range options.ConfigPaths {
		v.AddConfigPath(path)
	}

	return &FileDriver{
		viper:       v,
		configPaths: options.ConfigPaths,
		configName:  options.ConfigName,
		configType:  options.ConfigType,
	}
}

// Get retrieves a configuration value by key.
func (f *FileDriver) Get(key string) (interface{}, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if !f.viper.IsSet(key) {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return f.viper.Get(key), nil
}

// Set stores a configuration value by key.
func (f *FileDriver) Set(key string, value interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.viper.Set(key, value)
	return nil
}

// Load loads the configuration from files.
func (f *FileDriver) Load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// First try to load Go config files
	if f.configType == "go" || f.configType == "" {
		if err := f.loadGoConfig(); err == nil {
			return nil
		} else if f.configType == "go" {
			// If explicitly looking for Go files and failed
			return err
		}
	}

	// Fallback to standard Viper loading
	if err := f.viper.ReadInConfig(); err != nil {
		// Check if it's a file not found error
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if you wish
			return fmt.Errorf("config file not found: %w", err)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

// Watch starts watching for configuration changes.
func (f *FileDriver) Watch(callback func()) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.watchFunc != nil {
		return fmt.Errorf("already watching for config changes")
	}

	f.watchFunc = callback

	// Use Viper's built-in watching functionality
	f.viper.WatchConfig()
	f.viper.OnConfigChange(func(e fsnotify.Event) {
		if f.watchFunc != nil {
			f.watchFunc()
		}
	})

	return nil
}

// Unwatch stops watching for configuration changes.
func (f *FileDriver) Unwatch() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.watchFunc = nil
	// Note: Viper doesn't provide a built-in way to stop watching
	// This is a limitation of the current Viper implementation
	return nil
}

// Invalidate clears the configuration cache and forces a reload on next access.
func (f *FileDriver) Invalidate() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Reset viper instance
	f.viper = viper.New()
	f.viper.SetConfigName(f.configName)
	f.viper.SetConfigType(f.configType)

	for _, path := range f.configPaths {
		f.viper.AddConfigPath(path)
	}

	return nil
}

// GetAll returns all configuration data as a map.
func (f *FileDriver) GetAll() (map[string]interface{}, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.viper.AllSettings(), nil
}

// Has checks if a configuration key exists.
func (f *FileDriver) Has(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.viper.IsSet(key)
}

// Delete removes a configuration key.
func (f *FileDriver) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Viper doesn't have a built-in delete method
	// We need to implement this by getting all settings and removing the key
	allSettings := f.viper.AllSettings()
	delete(allSettings, key)

	// Clear and reload settings
	f.viper = viper.New()
	f.viper.SetConfigName(f.configName)
	f.viper.SetConfigType(f.configType)

	for _, path := range f.configPaths {
		f.viper.AddConfigPath(path)
	}

	// Set all settings except the deleted key
	for k, v := range allSettings {
		f.viper.Set(k, v)
	}

	return nil
}

// SetConfigFile sets the exact config file path.
func (f *FileDriver) SetConfigFile(filepath string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.viper.SetConfigFile(filepath)
}

// GetConfigFile returns the config file path.
func (f *FileDriver) GetConfigFile() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.viper.ConfigFileUsed()
}

// WriteConfig writes the current configuration to file.
func (f *FileDriver) WriteConfig() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.viper.WriteConfig()
}

// WriteConfigAs writes the current configuration to a specific file.
func (f *FileDriver) WriteConfigAs(filename string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.viper.WriteConfigAs(filename)
}

// SafeWriteConfig writes the current configuration to file only if it doesn't exist.
func (f *FileDriver) SafeWriteConfig() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.viper.SafeWriteConfig()
}

// SafeWriteConfigAs writes the current configuration to a specific file only if it doesn't exist.
func (f *FileDriver) SafeWriteConfigAs(filename string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.viper.SafeWriteConfigAs(filename)
}

// loadGoConfig attempts to load configuration from Go files.
// Uses the compiler package to execute Go config files and extract configuration.
func (f *FileDriver) loadGoConfig() error {
	for _, configPath := range f.configPaths {
		goFile := filepath.Join(configPath, f.configName+".go")
		if _, err := os.Stat(goFile); err == nil {
			// Use compiler to execute Go config file
			configMap, err := compiler.NewCompiler().Compile(goFile).GetContent()
			if err != nil {
				return fmt.Errorf("failed to execute Go config file %s: %w", goFile, err)
			}

			// Use viper.Set to set the entire config under the config name
			f.viper.Set(f.configName, configMap)

			return nil
		}
	}
	return fmt.Errorf("Go config file '%s.go' not found in paths: %v", f.configName, f.configPaths)
}
