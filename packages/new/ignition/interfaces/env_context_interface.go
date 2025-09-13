package interfaces

// EnvContextInterface interface for environment information
type EnvContextInterface interface {
	GetGoVersion() string
	SetGoVersion(string)
	GetOS() string
	SetOS(string)
	GetArch() string
	SetArch(string)
	GetGoMajorVersion() string
	IsWindows() bool
	IsLinux() bool
	IsDarwin() bool
	IsAmd64() bool
	IsArm64() bool
	Is386() bool
	GetPlatformString() string
	GetDisplayName() string
	GetArchDisplayName() string
}
