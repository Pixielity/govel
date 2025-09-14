package tests

import (
	"testing"

	containerMocks "govel/container/mocks"
	"govel/hashing"
	configMocks "govel/config/mocks"
)

func TestNewHashManager(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container (required by Manager)
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	// Create hash manager
	manager := src.NewHashManager(container)

	if manager == nil {
		t.Fatal("NewHashManager should return a valid manager")
	}
}

func TestHashManager_GetDefaultDriver(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)

	// Test getting default driver name
	defaultDriver := manager.GetDefaultDriver()
	if defaultDriver == "" {
		t.Error("Default driver name should not be empty")
	}
}

func TestHashManager_HasherWithDefaultDriver(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)

	// Test getting default driver
	defaultHasher := manager.Hasher()
	if defaultHasher == nil {
		t.Error("Default hasher should not be nil")
	}
}

func TestHashManager_HasherWithSpecificAlgorithm(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)

	algorithms := []string{"bcrypt", "argon2i", "argon2id"}

	for _, algorithm := range algorithms {
		t.Run(algorithm, func(t *testing.T) {
			hasher := manager.Hasher(algorithm)
			if hasher == nil {
				t.Errorf("%s hasher should not be nil", algorithm)
			}
		})
	}
}

func TestHashManager_HasherWithInvalidAlgorithm(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)

	invalidAlgorithms := []string{
		"md5",
		"sha1",
		"plaintext",
		"nonexistent",
		"scrypt",
	}

	for _, algorithm := range invalidAlgorithms {
		t.Run(algorithm, func(t *testing.T) {
			hasher := manager.Hasher(algorithm)
			if hasher != nil {
				t.Errorf("Invalid algorithm %s should return nil hasher", algorithm)
			}
		})
	}
}

func TestHashManager_HasherWithConfiguration(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set up configuration for different algorithms
	config.Set("hashing.bcrypt", map[string]interface{}{
		"cost": 8,
	})
	config.Set("hashing.argon2i", map[string]interface{}{
		"memory":  32768,
		"time":    1,
		"threads": 1,
		"keyLen":  32,
	})
	config.Set("hashing.argon2id", map[string]interface{}{
		"memory":  32768,
		"time":    1,
		"threads": 1,
		"keyLen":  32,
	})

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)

	algorithms := []string{"bcrypt", "argon2i", "argon2id"}

	for _, algorithm := range algorithms {
		t.Run(algorithm, func(t *testing.T) {
			hasher := manager.Hasher(algorithm)
			if hasher == nil {
				t.Errorf("%s hasher with configuration should not be nil", algorithm)
			}
		})
	}
}

func TestHashManager_DirectMethods(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)
	password := "testPassword"

	// Test direct Make method (uses default driver)
	hash, err := manager.Make(password, nil)
	if err != nil {
		t.Fatalf("Manager Make should not return error: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	// Test direct Check method
	if !manager.Check(password, hash, nil) {
		t.Error("Manager Check should return true for correct password")
	}

	if manager.Check("wrongPassword", hash, nil) {
		t.Error("Manager Check should return false for wrong password")
	}

	// Test direct Info method
	info := manager.Info(hash)
	if info.AlgoName == "" {
		t.Error("Manager Info should return algorithm name")
	}

	// Test direct IsHashed method
	if !manager.IsHashed(hash) {
		t.Error("Manager IsHashed should return true for hash")
	}

	if manager.IsHashed("plaintext") {
		t.Error("Manager IsHashed should return false for plaintext")
	}

	// Test direct NeedsRehash method
	needsRehash := manager.NeedsRehash(hash, nil)
	// Just check it doesn't panic - result depends on implementation
	_ = needsRehash
}
