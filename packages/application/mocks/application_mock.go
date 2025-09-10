package mocks

import (
	"context"
	"os"
	"time"

	applicationInterfaces "govel/packages/application/interfaces/application"
	"govel/packages/application/types"
	configMocks "govel/packages/config/mocks"
	containerMocks "govel/packages/container/mocks"
	loggerMocks "govel/packages/logger/mocks"
)

/**
 * MockApplication provides a complete mock implementation of ApplicationInterface for testing.
 * This mock embeds dedicated mocks for config, container, and logger functionality.
 */
type MockApplication struct {
	// Application Identity
	NameValue    string
	VersionValue string

	// Runtime State
	RunningInConsoleValue bool
	RunningUnitTestsValue bool
	StartTimeValue        time.Time

	// Embedded Mocks for better reusability
	*configMocks.MockConfigurable
	*containerMocks.MockContainable
	*loggerMocks.MockLoggable

	// Mock Control Flags
	ShouldFailRegister bool
	ShouldFailBoot     bool
}

/**
 * NewMockApplication creates a new mock application with default values
 */
func NewMockApplication() *MockApplication {
	return &MockApplication{
		NameValue:             "MockApp",
		VersionValue:          "1.0.0-mock",
		RunningInConsoleValue: false,
		RunningUnitTestsValue: true,
		StartTimeValue:        time.Now(),
		MockConfigurable:      configMocks.NewMockConfigurable(),
		MockContainable:       containerMocks.NewMockContainable(),
		MockLoggable:          loggerMocks.NewMockLoggable(),
	}
}

// ApplicationInterface Implementation

// ApplicationIdentityInterface Methods
func (m *MockApplication) GetName() string     { return m.NameValue }
func (m *MockApplication) SetName(name string) { m.NameValue = name }
func (m *MockApplication) Name() string        { return m.NameValue }
func (m *MockApplication) GetVersion() string  { return m.VersionValue }
func (m *MockApplication) SetVersion(v string) { m.VersionValue = v }
func (m *MockApplication) Version() string     { return m.VersionValue }

// ApplicationRuntimeInterface Methods
func (m *MockApplication) IsRunningInConsole() bool   { return m.RunningInConsoleValue }
func (m *MockApplication) SetRunningInConsole(b bool) { m.RunningInConsoleValue = b }
func (m *MockApplication) IsRunningUnitTests() bool   { return m.RunningUnitTestsValue }
func (m *MockApplication) SetRunningUnitTests(b bool) { m.RunningUnitTestsValue = b }

// ApplicationTimingInterface Methods
func (m *MockApplication) GetStartTime() time.Time  { return m.StartTimeValue }
func (m *MockApplication) SetStartTime(t time.Time) { m.StartTimeValue = t }
func (m *MockApplication) GetUptime() time.Duration { return time.Since(m.StartTimeValue) }

// ApplicationInfoInterface Methods
func (m *MockApplication) GetApplicationInfo() map[string]interface{} {
	info := map[string]interface{}{
		"name":               m.NameValue,
		"version":            m.VersionValue,
		"running_in_console": m.RunningInConsoleValue,
		"running_unit_tests": m.RunningUnitTestsValue,
		"start_time":         m.StartTimeValue,
		"uptime":             m.GetUptime(),
		"type":               "mock",
	}

	// Add info from embedded mocks
	if m.MockConfigurable != nil {
		info["config"] = m.MockConfigurable.GetConfigInfo()
	}
	if m.MockContainable != nil {
		info["container"] = m.MockContainable.GetContainerInfo()
	}
	if m.MockLoggable != nil {
		info["logger"] = m.MockLoggable.GetLoggerInfo()
	}

	return info
}

// DirectableInterface Methods
func (m *MockApplication) GetBasePath() string             { return "/mock/path" }
func (m *MockApplication) SetBasePath(path string)         {}
func (m *MockApplication) PublicPath() string              { return "/mock/path/public" }
func (m *MockApplication) SetPublicPath(path string)       {}
func (m *MockApplication) StoragePath() string             { return "/mock/path/storage" }
func (m *MockApplication) SetStoragePath(path string)      {}
func (m *MockApplication) ConfigPath() string              { return "/mock/path/config" }
func (m *MockApplication) SetConfigPath(path string)       {}
func (m *MockApplication) LogPath() string                 { return "/mock/path/log" }
func (m *MockApplication) SetLogPath(path string)          {}
func (m *MockApplication) CachePath() string               { return "/mock/path/cache" }
func (m *MockApplication) SetCachePath(path string)        {}
func (m *MockApplication) ViewPath() string                { return "/mock/path/views" }
func (m *MockApplication) SetViewPath(path string)         {}
func (m *MockApplication) GetCustomPath(key string) string { return "/mock/path/" + key }
func (m *MockApplication) SetCustomPath(key, path string)  {}
func (m *MockApplication) GetAllCustomPaths() map[string]string {
	return map[string]string{"mock": "/mock/path"}
}
func (m *MockApplication) ClearCustomPaths()                       {}
func (m *MockApplication) EnsureDirectoryExists(path string) error { return nil }

// EnvironmentableInterface Methods
func (m *MockApplication) GetEnvironment() string             { return "testing" }
func (m *MockApplication) SetEnvironment(env string)          {}
func (m *MockApplication) IsProduction() bool                 { return false }
func (m *MockApplication) IsDevelopment() bool                { return false }
func (m *MockApplication) IsTesting() bool                    { return true }
func (m *MockApplication) IsStaging() bool                    { return false }
func (m *MockApplication) IsDebug() bool                      { return true }
func (m *MockApplication) SetDebug(debug bool)                {}
func (m *MockApplication) EnableDebug()                       {}
func (m *MockApplication) DisableDebug()                      {}
func (m *MockApplication) IsEnvironment(env string) bool      { return env == "testing" }
func (m *MockApplication) IsValidEnvironment(env string) bool { return true }
func (m *MockApplication) GetEnvironmentInfo() map[string]interface{} {
	return map[string]interface{}{
		"environment":   "testing",
		"is_production": false,
		"is_debug":      true,
	}
}

// HookableInterface Methods
func (m *MockApplication) RegisterHook(name string, priority int, callback types.HookCallback) {}
func (m *MockApplication) UnregisterHook(name string) bool                                     { return false }
func (m *MockApplication) UnregisterHookCallback(name string, callback types.HookCallback) bool {
	return false
}
func (m *MockApplication) HasHook(name string) bool { return false }
func (m *MockApplication) GetHooks() map[string][]types.HookCallback {
	return map[string][]types.HookCallback{}
}
func (m *MockApplication) GetHookCallbacks(name string) []types.HookCallback {
	return []types.HookCallback{}
}
func (m *MockApplication) CallHook(name string, args ...interface{}) ([]interface{}, error) {
	return []interface{}{}, nil
}
func (m *MockApplication) CallHookFirst(name string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
func (m *MockApplication) CallHookUntil(name string, args ...interface{}) (bool, error) {
	return false, nil
}
func (m *MockApplication) GetHookCount(name string) int { return 0 }
func (m *MockApplication) GetAllHookNames() []string    { return []string{} }
func (m *MockApplication) ClearHooks()                  {}
func (m *MockApplication) GetHooksInfo() map[string]interface{} {
	return map[string]interface{}{"total_hooks": 0}
}

// LifecycleableInterface Methods
func (m *MockApplication) Boot(ctx context.Context) error    { return nil }
func (m *MockApplication) IsBooted() bool                    { return true }
func (m *MockApplication) SetBooted(booted bool)             {}
func (m *MockApplication) Start(ctx context.Context) error   { return nil }
func (m *MockApplication) IsStarted() bool                   { return true }
func (m *MockApplication) SetStarted(started bool)           {}
func (m *MockApplication) Stop(ctx context.Context) error    { return nil }
func (m *MockApplication) IsStopped() bool                   { return false }
func (m *MockApplication) SetStopped(stopped bool)           {}
func (m *MockApplication) Restart(ctx context.Context) error { return nil }
func (m *MockApplication) GetState() string                  { return "running" }
func (m *MockApplication) IsState(state string) bool         { return state == "running" }
func (m *MockApplication) GetLifecycleInfo() map[string]interface{} {
	return map[string]interface{}{
		"is_booted":  true,
		"is_started": true,
		"is_stopped": false,
		"state":      "running",
	}
}

// LocalizableInterface Methods
func (m *MockApplication) GetLocale() string                       { return "en" }
func (m *MockApplication) SetLocale(locale string)                 {}
func (m *MockApplication) GetFallbackLocale() string               { return "en" }
func (m *MockApplication) SetFallbackLocale(fallbackLocale string) {}
func (m *MockApplication) GetTimezone() string                     { return "UTC" }
func (m *MockApplication) SetTimezone(timezone string)             {}
func (m *MockApplication) IsLocale(locale string) bool             { return locale == "en" }
func (m *MockApplication) IsTimezone(timezone string) bool         { return timezone == "UTC" }
func (m *MockApplication) GetLocaleWithFallback() string           { return "en" }
func (m *MockApplication) IsValidLocale(locale string) bool        { return true }
func (m *MockApplication) GetLanguageCode() string                 { return "en" }
func (m *MockApplication) GetCountryCode() string                  { return "" }
func (m *MockApplication) LocaleInfo() map[string]string {
	return map[string]string{
		"locale":          "en",
		"fallback_locale": "en",
		"timezone":        "UTC",
	}
}
func (m *MockApplication) SetLocaleInfo(locale, fallback, timezone string) {}
func (m *MockApplication) IsRTL() bool                                     { return false }
func (m *MockApplication) GetTextDirection() string                        { return "ltr" }

// MaintainableInterface Methods
func (m *MockApplication) IsDown() bool                                           { return false }
func (m *MockApplication) IsUp() bool                                             { return true }
func (m *MockApplication) Down(mode *types.MaintenanceMode) error                 { return nil }
func (m *MockApplication) Up() error                                              { return nil }
func (m *MockApplication) GetMaintenanceMode() *types.MaintenanceMode             { return nil }
func (m *MockApplication) CanBypass(clientIP, path, secret string) bool           { return true }
func (m *MockApplication) GetMaintenanceDuration() time.Duration                  { return 0 }
func (m *MockApplication) SetMaintenanceMessage(message string) error             { return nil }
func (m *MockApplication) SetRetryAfter(seconds int) error                        { return nil }
func (m *MockApplication) AddAllowedIP(ip string) error                           { return nil }
func (m *MockApplication) RemoveAllowedIP(ip string) error                        { return nil }
func (m *MockApplication) AddAllowedPath(path string) error                       { return nil }
func (m *MockApplication) RemoveAllowedPath(path string) error                    { return nil }
func (m *MockApplication) SetMaintenanceData(key string, value interface{}) error { return nil }
func (m *MockApplication) GetMaintenanceData(key string) interface{}              { return nil }
func (m *MockApplication) GetMaintenanceInfo() map[string]interface{} {
	return map[string]interface{}{
		"is_down":          false,
		"is_up":            true,
		"maintenance_mode": nil,
		"duration":         0,
	}
}

// ShutdownableInterface Methods
func (m *MockApplication) RegisterShutdownCallback(name string, callback types.ShutdownCallback) {}
func (m *MockApplication) UnregisterShutdownCallback(name string) bool                           { return false }
func (m *MockApplication) GetShutdownCallbacks() map[string]types.ShutdownCallback {
	return map[string]types.ShutdownCallback{}
}
func (m *MockApplication) Shutdown(ctx context.Context) error                       { return nil }
func (m *MockApplication) GracefulShutdown(timeout time.Duration) error             { return nil }
func (m *MockApplication) ForceShutdown()                                           {}
func (m *MockApplication) IsShuttingDown() bool                                     { return false }
func (m *MockApplication) SetShuttingDown(shutting bool)                            {}
func (m *MockApplication) IsShutdown() bool                                         { return false }
func (m *MockApplication) SetShutdown(shutdown bool)                                {}
func (m *MockApplication) HandleSignals(signals []os.Signal, timeout time.Duration) {}
func (m *MockApplication) GetShutdownTimeout() time.Duration                        { return 30 * time.Second }
func (m *MockApplication) SetShutdownTimeout(timeout time.Duration)                 {}
func (m *MockApplication) GetShutdownInfo() map[string]interface{} {
	return map[string]interface{}{
		"is_shutting_down": false,
		"is_shutdown":      false,
		"timeout":          30 * time.Second,
	}
}

// Trait capability methods
func (m *MockApplication) Directable() bool      { return true }
func (m *MockApplication) Environmentable() bool { return true }
func (m *MockApplication) Hookable() bool        { return true }
func (m *MockApplication) Lifecycleable() bool   { return true }
func (m *MockApplication) Localizable() bool     { return true }
func (m *MockApplication) Maintainable() bool    { return true }
func (m *MockApplication) Shutdownable() bool    { return true }

// Mock-specific helper methods

/**
 * GetMockConfig returns the embedded config mock for advanced testing
 */
func (m *MockApplication) GetMockConfig() *configMocks.MockConfig {
	if m.MockConfigurable != nil {
		return m.MockConfigurable.GetMockConfig()
	}
	return nil
}

/**
 * GetMockContainer returns the embedded container mock for advanced testing
 */
func (m *MockApplication) GetMockContainer() *containerMocks.MockContainer {
	if m.MockContainable != nil {
		return m.MockContainable.GetMockContainer()
	}
	return nil
}

/**
 * GetMockLogger returns the embedded logger mock for advanced testing
 */
func (m *MockApplication) GetMockLogger() *loggerMocks.MockLogger {
	if m.MockLoggable != nil {
		return m.MockLoggable.GetMockLogger()
	}
	return nil
}

// ApplicationProviderInterface Methods
func (m *MockApplication) RegisterProvider(provider interface{}) error {
	if m.ShouldFailRegister {
		return &MockError{Message: "register provider failed"}
	}
	return nil
}

func (m *MockApplication) RegisterProviders(providers []interface{}) error {
	if m.ShouldFailRegister {
		return &MockError{Message: "register providers failed"}
	}
	return nil
}

func (m *MockApplication) BootProviders(ctx context.Context) error {
	if m.ShouldFailBoot {
		return &MockError{Message: "boot providers failed"}
	}
	return nil
}

func (m *MockApplication) TerminateProviders(ctx context.Context) []error {
	return []error{}
}

func (m *MockApplication) LoadDeferredProvider(service string) error {
	return nil
}

func (m *MockApplication) GetProviderRepository() interface{} {
	return nil // Mock repository could be implemented if needed
}

func (m *MockApplication) GetRegisteredProviders() []interface{} {
	return []interface{}{}
}

func (m *MockApplication) GetLoadedProviders() []string {
	return []string{}
}

func (m *MockApplication) GetBootedProviders() []string {
	return []string{}
}

/**
 * SetFailureMode sets whether various operations should fail
 */
func (m *MockApplication) SetFailureMode(bind, make, register, boot bool) {
	m.ShouldFailRegister = register
	m.ShouldFailBoot = boot

	// Delegate to embedded mocks
	if mockContainer := m.GetMockContainer(); mockContainer != nil {
		mockContainer.SetFailureMode(bind, make, false)
	}
}

// Mock Error Type
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return "mock error: " + e.Message
}

// Compile-time interface compliance check
var _ applicationInterfaces.ApplicationInterface = (*MockApplication)(nil)
