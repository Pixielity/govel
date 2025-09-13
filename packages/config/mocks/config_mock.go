package mocks

import (
	interfaces "govel/packages/types/src/interfaces/config"
	"time"
)

/**
 * MockConfig provides a mock implementation of ConfigInterface for testing.
 * This mock allows tests to verify configuration behavior without dependencies on actual file systems or environments.
 */
type MockConfig struct {
	// Configuration Data Storage
	Data map[string]interface{}

	// Environment Settings
	Environment       string
	EnvironmentPrefix string

	// Mock Control Flags
	ShouldFailLoadFile bool
	ShouldFailLoadEnv  bool

	// File Loading History
	LoadedFiles []string
	LoadedEnvs  []string
}

/**
 * NewMockConfig creates a new mock config with default values
 */
func NewMockConfig() *MockConfig {
	return &MockConfig{
		Data:              make(map[string]interface{}),
		Environment:       "testing",
		EnvironmentPrefix: "TEST_",
		LoadedFiles:       make([]string, 0),
		LoadedEnvs:        make([]string, 0),
	}
}

/**
 * NewMockConfigWithData creates a new mock config with predefined data
 */
func NewMockConfigWithData(data map[string]interface{}) *MockConfig {
	mock := NewMockConfig()
	mock.Data = data
	return mock
}

// ConfigInterface Implementation

func (m *MockConfig) GetString(key string, defaultValue ...string) string {
	default_ := ""
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return default_
}

func (m *MockConfig) GetInt(key string, defaultValue ...int) int {
	default_ := 0
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
		// Handle int64 to int conversion
		if i64, ok := val.(int64); ok {
			return int(i64)
		}
		// Handle float64 to int conversion (for JSON numbers)
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return default_
}

func (m *MockConfig) GetInt64(key string, defaultValue ...int64) int64 {
	default_ := int64(0)
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if i, ok := val.(int64); ok {
			return i
		}
		// Handle int to int64 conversion
		if i32, ok := val.(int); ok {
			return int64(i32)
		}
		// Handle float64 to int64 conversion
		if f, ok := val.(float64); ok {
			return int64(f)
		}
	}
	return default_
}

func (m *MockConfig) GetFloat64(key string, defaultValue ...float64) float64 {
	default_ := 0.0
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
		// Handle int to float64 conversion
		if i, ok := val.(int); ok {
			return float64(i)
		}
		// Handle int64 to float64 conversion
		if i64, ok := val.(int64); ok {
			return float64(i64)
		}
	}
	return default_
}

func (m *MockConfig) GetBool(key string, defaultValue ...bool) bool {
	default_ := false
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
		// Handle string to bool conversion
		if str, ok := val.(string); ok {
			switch str {
			case "true", "1", "yes", "on":
				return true
			case "false", "0", "no", "off":
				return false
			}
		}
	}
	return default_
}

func (m *MockConfig) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	default_ := time.Duration(0)
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if d, ok := val.(time.Duration); ok {
			return d
		}
		// Handle string to duration conversion
		if str, ok := val.(string); ok {
			if duration, err := time.ParseDuration(str); err == nil {
				return duration
			}
		}
		// Handle int/int64 as seconds
		if i, ok := val.(int); ok {
			return time.Duration(i) * time.Second
		}
		if i64, ok := val.(int64); ok {
			return time.Duration(i64) * time.Second
		}
	}
	return default_
}

func (m *MockConfig) GetStringSlice(key string, defaultValue ...[]string) []string {
	default_ := []string{}
	if len(defaultValue) > 0 {
		default_ = defaultValue[0]
	}

	if val, ok := m.Data[key]; ok {
		if slice, ok := val.([]string); ok {
			return slice
		}
		// Handle []interface{} conversion
		if interfaceSlice, ok := val.([]interface{}); ok {
			stringSlice := make([]string, len(interfaceSlice))
			for i, item := range interfaceSlice {
				if str, ok := item.(string); ok {
					stringSlice[i] = str
				} else {
					stringSlice[i] = ""
				}
			}
			return stringSlice
		}
	}
	return default_
}

func (m *MockConfig) Get(key string) (interface{}, bool) {
	val, exists := m.Data[key]
	return val, exists
}

func (m *MockConfig) Set(key string, value interface{}) {
	m.Data[key] = value
}

func (m *MockConfig) HasKey(key string) bool {
	_, exists := m.Data[key]
	return exists
}

func (m *MockConfig) AllConfig() map[string]interface{} {
	// Return a copy to prevent external modification
	result := make(map[string]interface{})
	for k, v := range m.Data {
		result[k] = v
	}
	return result
}

func (m *MockConfig) LoadFromFile(filePath string) error {
	if m.ShouldFailLoadFile {
		return &MockConfigError{Message: "mock file load failure", FilePath: filePath}
	}

	// Record the file loading attempt
	m.LoadedFiles = append(m.LoadedFiles, filePath)

	// Simulate loading some data based on file path
	switch filePath {
	case "config/app.json":
		m.Set("app.name", "Mock App")
		m.Set("app.version", "1.0.0")
	case "config/database.json":
		m.Set("database.host", "localhost")
		m.Set("database.port", 3306)
	case "config/logging.json":
		m.Set("logging.level", "debug")
		m.Set("logging.file", "/var/log/app.log")
	}

	return nil
}

func (m *MockConfig) LoadFromEnv(prefix string) error {
	if m.ShouldFailLoadEnv {
		return &MockConfigError{Message: "mock env load failure", Prefix: prefix}
	}

	// Record the env loading attempt
	m.LoadedEnvs = append(m.LoadedEnvs, prefix)

	// Simulate loading some environment variables
	envData := map[string]interface{}{
		"app.env":       m.Environment,
		"app.debug":     true,
		"server.port":   8080,
		"server.host":   "0.0.0.0",
		"cache.enabled": true,
	}

	for key, value := range envData {
		m.Set(key, value)
	}

	return nil
}

// Mock-specific helper methods

/**
 * GetEnvironment returns the current environment
 */
func (m *MockConfig) GetEnvironment() string {
	return m.Environment
}

/**
 * SetEnvironment sets the environment
 */
func (m *MockConfig) SetEnvironment(env string) {
	m.Environment = env
}

/**
 * GetLoadedFiles returns all files that were attempted to be loaded
 */
func (m *MockConfig) GetLoadedFiles() []string {
	return m.LoadedFiles
}

/**
 * GetLoadedEnvs returns all environment prefixes that were attempted to be loaded
 */
func (m *MockConfig) GetLoadedEnvs() []string {
	return m.LoadedEnvs
}

/**
 * SetFailureMode sets whether file/env loading should fail
 */
func (m *MockConfig) SetFailureMode(failFile, failEnv bool) {
	m.ShouldFailLoadFile = failFile
	m.ShouldFailLoadEnv = failEnv
}

/**
 * ClearData clears all configuration data
 */
func (m *MockConfig) ClearData() {
	m.Data = make(map[string]interface{})
}

/**
 * ClearHistory clears the loading history
 */
func (m *MockConfig) ClearHistory() {
	m.LoadedFiles = make([]string, 0)
	m.LoadedEnvs = make([]string, 0)
}

/**
 * SetBulkData sets multiple configuration values at once
 */
func (m *MockConfig) SetBulkData(data map[string]interface{}) {
	for key, value := range data {
		m.Data[key] = value
	}
}

/**
 * GetKeysWithPrefix returns all keys that start with the given prefix
 */
func (m *MockConfig) GetKeysWithPrefix(prefix string) []string {
	var keys []string
	for key := range m.Data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			keys = append(keys, key)
		}
	}
	return keys
}

/**
 * GetDataSize returns the number of configuration entries
 */
func (m *MockConfig) GetDataSize() int {
	return len(m.Data)
}

// Mock Error Type
type MockConfigError struct {
	Message  string
	FilePath string
	Prefix   string
}

func (e *MockConfigError) Error() string {
	if e.FilePath != "" {
		return "mock config error (file: " + e.FilePath + "): " + e.Message
	}
	if e.Prefix != "" {
		return "mock config error (prefix: " + e.Prefix + "): " + e.Message
	}
	return "mock config error: " + e.Message
}

// Compile-time interface compliance check removed

/**
 * MockConfigurable provides a mock implementation of ConfigurableInterface for testing.
 */
type MockConfigurable struct {
	*MockConfig

	ConfigInstance interfaces.ConfigInterface
	HasConfigValue bool
}

/**
 * NewMockConfigurable creates a new mock configurable with default values
 */
func NewMockConfigurable() *MockConfigurable {
	mockConfig := NewMockConfig()
	return &MockConfigurable{
		MockConfig:     mockConfig,
		ConfigInstance: mockConfig,
		HasConfigValue: true,
	}
}

// ConfigurableInterface Implementation

func (m *MockConfigurable) GetConfig() interfaces.ConfigInterface {
	return m.ConfigInstance
}

func (m *MockConfigurable) SetConfig(config interfaces.ConfigInterface) {
	if cfg, ok := config.(interfaces.ConfigInterface); ok {
		m.ConfigInstance = cfg
		m.HasConfigValue = true
	} else if cfg, ok := config.(*MockConfig); ok {
		m.ConfigInstance = cfg
		m.HasConfigValue = true
	}
}

func (m *MockConfigurable) HasConfig() bool {
	return m.HasConfigValue
}

func (m *MockConfigurable) GetConfigInfo() map[string]interface{} {
	info := map[string]interface{}{
		"has_config":  m.HasConfigValue,
		"config_type": "mock",
	}

	if m.ConfigInstance != nil {
		if mockConfig, ok := m.ConfigInstance.(*MockConfig); ok {
			info["config_keys_count"] = len(mockConfig.AllConfig())
			info["environment"] = mockConfig.GetEnvironment()
			info["loaded_files_count"] = len(mockConfig.GetLoadedFiles())
			info["loaded_envs_count"] = len(mockConfig.GetLoadedEnvs())
		} else {
			info["config_keys_count"] = 0
		}
	}

	return info
}

// Mock-specific helper methods for Configurable

/**
 * SetHasConfig controls whether the configurable reports having a config
 */
func (m *MockConfigurable) SetHasConfig(hasConfig bool) {
	m.HasConfigValue = hasConfig
}

/**
 * GetMockConfig returns the underlying MockConfig if available
 */
func (m *MockConfigurable) GetMockConfig() *MockConfig {
	if mockConfig, ok := m.ConfigInstance.(*MockConfig); ok {
		return mockConfig
	}
	return nil
}

// Compile-time interface compliance check removed
