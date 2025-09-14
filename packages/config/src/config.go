package config

import (
	"time"

	configInterfaces "govel/types/interfaces/config"

	"github.com/spf13/cast"
)

// Config is the main configuration struct that uses a driver for data operations.
// It provides type-safe methods for retrieving configuration values.
type Config struct {
	driver configInterfaces.DriverInterface
}

// NewConfig creates a new Config instance with the provided driver.
func NewConfig(driver configInterfaces.DriverInterface) *Config {
	return &Config{
		driver: driver,
	}
}

// GetString returns a string value for the given key with optional default.
func (c *Config) GetString(key string, defaultValue ...string) string {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return cast.ToString(value)
}

// GetInt returns an integer value for the given key with optional default.
func (c *Config) GetInt(key string, defaultValue ...int) int {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToInt(value)
}

// GetInt32 returns an int32 value for the given key with optional default.
func (c *Config) GetInt32(key string, defaultValue ...int32) int32 {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToInt32(value)
}

// GetInt64 returns an int64 value for the given key with optional default.
func (c *Config) GetInt64(key string, defaultValue ...int64) int64 {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToInt64(value)
}

// GetUint returns a uint value for the given key with optional default.
func (c *Config) GetUint(key string, defaultValue ...uint) uint {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToUint(value)
}

// GetUint32 returns a uint32 value for the given key with optional default.
func (c *Config) GetUint32(key string, defaultValue ...uint32) uint32 {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToUint32(value)
}

// GetUint64 returns a uint64 value for the given key with optional default.
func (c *Config) GetUint64(key string, defaultValue ...uint64) uint64 {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToUint64(value)
}

// GetFloat64 returns a float64 value for the given key with optional default.
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}
	return cast.ToFloat64(value)
}

// GetBool returns a boolean value for the given key with optional default.
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return cast.ToBool(value)
}

// GetTime returns a time.Time value for the given key with optional default.
func (c *Config) GetTime(key string, defaultValue ...time.Time) time.Time {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return time.Time{}
	}
	return cast.ToTime(value)
}

// GetDuration returns a time.Duration value for the given key with optional default.
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return cast.ToDuration(value)
}

// GetIntSlice returns a []int value for the given key with optional default.
func (c *Config) GetIntSlice(key string, defaultValue ...[]int) []int {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []int{}
	}
	return cast.ToIntSlice(value)
}

// GetStringSlice returns a []string value for the given key with optional default.
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return []string{}
	}
	return cast.ToStringSlice(value)
}

// GetStringMap returns a map[string]interface{} value for the given key with optional default.
func (c *Config) GetStringMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return map[string]interface{}{}
	}
	return cast.ToStringMap(value)
}

// GetStringMapString returns a map[string]string value for the given key with optional default.
func (c *Config) GetStringMapString(key string, defaultValue ...map[string]string) map[string]string {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return map[string]string{}
	}
	return cast.ToStringMapString(value)
}

// GetStringMapStringSlice returns a map[string][]string value for the given key with optional default.
func (c *Config) GetStringMapStringSlice(key string, defaultValue ...map[string][]string) map[string][]string {
	value, err := c.driver.Get(key)
	if err != nil || value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return map[string][]string{}
	}
	return cast.ToStringMapStringSlice(value)
}

// Get returns the raw value for the given key.
func (c *Config) Get(key string) interface{} {
	value, err := c.driver.Get(key)
	if err != nil {
		return nil
	}
	return value
}

// Set sets a configuration value.
func (c *Config) Set(key string, value interface{}) error {
	return c.driver.Set(key, value)
}

// Has checks if a configuration key exists.
func (c *Config) Has(key string) bool {
	return c.driver.Has(key)
}

// Delete removes a configuration key.
func (c *Config) Delete(key string) error {
	return c.driver.Delete(key)
}

// GetAll returns all configuration values.
func (c *Config) GetAll() (map[string]interface{}, error) {
	return c.driver.GetAll()
}

// Load loads configuration from the driver source.
func (c *Config) Load() error {
	return c.driver.Load()
}

// Watch starts watching for configuration changes.
func (c *Config) Watch(callback func()) error {
	return c.driver.Watch(callback)
}

// Unwatch stops watching for configuration changes.
func (c *Config) Unwatch() error {
	return c.driver.Unwatch()
}

// Invalidate invalidates the configuration cache.
func (c *Config) Invalidate() error {
	return c.driver.Invalidate()
}
