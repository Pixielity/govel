package enums

import "fmt"

// Driver represents the available configuration drivers.
// Each driver provides different sources and methods for loading configuration data.
type Driver string

// Available configuration drivers
const (
	// FileDriver loads configuration from files (JSON, YAML, TOML, etc.)
	FileDriver Driver = "file"
	
	// EnvDriver loads configuration from environment variables
	EnvDriver Driver = "env"
	
	// RemoteDriver loads configuration from remote sources (HTTP, databases, etc.)
	RemoteDriver Driver = "remote"
	
	// MemoryDriver stores configuration in memory (for testing/caching)
	MemoryDriver Driver = "memory"
)

// DefaultDriver is the default configuration driver used when none is specified
const DefaultDriver = FileDriver

// String returns the string representation of the driver
func (d Driver) String() string {
	return string(d)
}

// IsValid checks if the driver is a valid configuration driver
func (d Driver) IsValid() bool {
	switch d {
	case FileDriver, EnvDriver, RemoteDriver, MemoryDriver:
		return true
	default:
		return false
	}
}

// GetDefaultDriver returns the default configuration driver
func GetDefaultDriver() string {
	return DefaultDriver.String()
}

// ValidateDriver validates if a driver is supported
func ValidateDriver(driver Driver) error {
	if !driver.IsValid() {
		return fmt.Errorf("unsupported configuration driver: %s", driver)
	}
	return nil
}

// AllDrivers returns a slice of all available drivers
func AllDrivers() []Driver {
	return []Driver{
		FileDriver,
		EnvDriver,
		RemoteDriver,
		MemoryDriver,
	}
}

// ParseDriver parses a string into a Driver type
func ParseDriver(s string) (Driver, error) {
	driver := Driver(s)
	if err := ValidateDriver(driver); err != nil {
		return "", err
	}
	return driver, nil
}

// MustParseDriver parses a string into a Driver type, panicking on error
func MustParseDriver(s string) Driver {
	driver, err := ParseDriver(s)
	if err != nil {
		panic(fmt.Sprintf("invalid driver: %s", s))
	}
	return driver
}