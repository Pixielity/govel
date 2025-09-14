package drivers

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"

	configInterfaces "govel/types/interfaces/config"
)

// RemoteDriver implements config.Driver interface for remote configuration sources.
// It uses Viper's remote capabilities to connect to etcd, Consul, etc.
type RemoteDriver struct {
	viper     *viper.Viper
	provider  string
	endpoint  string
	keyPath   string
	mu        sync.RWMutex
	watchFunc func()
}

// RemoteDriverOptions contains configuration options for RemoteDriver.
type RemoteDriverOptions struct {
	Provider string // Remote provider: "etcd", "consul", "firestore"
	Endpoint string // Remote endpoint URL
	KeyPath  string // Key path in the remote store
}

// NewRemoteDriver creates a new remote configuration driver.
// Note: This requires Viper remote dependencies to be installed.
func NewRemoteDriver(options *RemoteDriverOptions) *RemoteDriver {
	v := viper.New()

	// Set defaults
	if options == nil {
		options = &RemoteDriverOptions{
			Provider: "etcd",
			Endpoint: "http://localhost:2379",
			KeyPath:  "/config",
		}
	}

	return &RemoteDriver{
		viper:    v,
		provider: options.Provider,
		endpoint: options.Endpoint,
		keyPath:  options.KeyPath,
	}
}

// Get retrieves a configuration value by key.
func (r *RemoteDriver) Get(key string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.viper.IsSet(key) {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	return r.viper.Get(key), nil
}

// Set stores a configuration value by key.
// Note: This operation may not be supported by all remote providers.
func (r *RemoteDriver) Set(key string, value interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.viper.Set(key, value)
	return nil
}

// Load loads configuration from the remote source.
func (r *RemoteDriver) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Note: Viper remote functionality requires additional setup
	// This is a placeholder implementation
	return fmt.Errorf("remote driver load not implemented - requires viper remote dependencies")
}

// Watch starts watching for configuration changes from remote source.
func (r *RemoteDriver) Watch(callback func()) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.watchFunc != nil {
		return fmt.Errorf("already watching for config changes")
	}

	r.watchFunc = callback

	// Note: Remote watching requires additional implementation
	return fmt.Errorf("remote driver watch not implemented - requires viper remote dependencies")
}

// Unwatch stops watching for configuration changes.
func (r *RemoteDriver) Unwatch() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.watchFunc = nil
	return nil
}

// Invalidate clears the configuration cache.
func (r *RemoteDriver) Invalidate() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.viper = viper.New()
	return nil
}

// GetAll returns all configuration data from remote source.
func (r *RemoteDriver) GetAll() (map[string]interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.viper.AllSettings(), nil
}

// Has checks if a configuration key exists.
func (r *RemoteDriver) Has(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.viper.IsSet(key)
}

// Delete removes a configuration key.
func (r *RemoteDriver) Delete(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// This would require custom implementation for remote deletion
	return fmt.Errorf("remote driver delete not implemented")
}

// SetProvider sets the remote provider.
func (r *RemoteDriver) SetProvider(provider string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.provider = provider
}

// SetEndpoint sets the remote endpoint.
func (r *RemoteDriver) SetEndpoint(endpoint string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.endpoint = endpoint
}

// SetKeyPath sets the key path in remote store.
func (r *RemoteDriver) SetKeyPath(keyPath string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.keyPath = keyPath
}

// Ensure RemoteDriver implements the Driver interface
var _ configInterfaces.DriverInterface = (*RemoteDriver)(nil)
