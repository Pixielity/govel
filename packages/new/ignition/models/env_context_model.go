package models

import (
	"runtime"
	"strings"

	"govel/packages/ignition/enums"
	"govel/packages/ignition/interfaces"
)

// EnvContext holds environment information
type EnvContext struct {
	GoVersion string         `json:"go_version"`
	OS        enums.OSType   `json:"os"`
	Arch      enums.ArchType `json:"arch"`
}

// NewEnvContext creates a new environment context with current runtime info
func NewEnvContext() *EnvContext {
	return &EnvContext{
		GoVersion: runtime.Version(),
		OS:        enums.ParseOSType(runtime.GOOS),
		Arch:      enums.ParseArchType(runtime.GOARCH),
	}
}

// GetGoVersion returns the Go version
func (e *EnvContext) GetGoVersion() string {
	return e.GoVersion
}

// SetGoVersion sets the Go version
func (e *EnvContext) SetGoVersion(goVersion string) {
	e.GoVersion = goVersion
}

// GetOS returns the operating system
func (e *EnvContext) GetOS() string {
	return e.OS.String()
}

// SetOS sets the operating system
func (e *EnvContext) SetOS(os string) {
	e.OS = enums.ParseOSType(os)
}

// GetOSType returns the operating system as enum
func (e *EnvContext) GetOSType() enums.OSType {
	return e.OS
}

// SetOSType sets the operating system from enum
func (e *EnvContext) SetOSType(osType enums.OSType) {
	e.OS = osType
}

// GetArch returns the architecture
func (e *EnvContext) GetArch() string {
	return e.Arch.String()
}

// SetArch sets the architecture
func (e *EnvContext) SetArch(arch string) {
	e.Arch = enums.ParseArchType(arch)
}

// GetArchType returns the architecture as enum
func (e *EnvContext) GetArchType() enums.ArchType {
	return e.Arch
}

// SetArchType sets the architecture from enum
func (e *EnvContext) SetArchType(archType enums.ArchType) {
	e.Arch = archType
}

// GetGoMajorVersion returns the major version of Go (e.g., "1.21")
func (e *EnvContext) GetGoMajorVersion() string {
	version := strings.TrimPrefix(e.GoVersion, "go")
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return version
}

// IsWindows returns true if running on Windows
func (e *EnvContext) IsWindows() bool {
	return e.OS == enums.OSWindows
}

// IsLinux returns true if running on Linux
func (e *EnvContext) IsLinux() bool {
	return e.OS == enums.OSLinux
}

// IsDarwin returns true if running on macOS
func (e *EnvContext) IsDarwin() bool {
	return e.OS == enums.OSDarwin
}

// IsAmd64 returns true if running on AMD64 architecture
func (e *EnvContext) IsAmd64() bool {
	return e.Arch == enums.ArchAMD64
}

// IsArm64 returns true if running on ARM64 architecture
func (e *EnvContext) IsArm64() bool {
	return e.Arch == enums.ArchARM64
}

// Is386 returns true if running on 386 architecture
func (e *EnvContext) Is386() bool {
	return e.Arch == enums.Arch386
}

// GetPlatformString returns a combined OS/Arch string
func (e *EnvContext) GetPlatformString() string {
	return e.OS.String() + "/" + e.Arch.String()
}

// GetDisplayName returns a human-readable platform name
func (e *EnvContext) GetDisplayName() string {
	return e.OS.DisplayName()
}

// GetArchDisplayName returns a human-readable architecture name
func (e *EnvContext) GetArchDisplayName() string {
	return e.Arch.DisplayName()
}

// Compile-time interface compliance check
var _ interfaces.EnvContextInterface = (*EnvContext)(nil)
