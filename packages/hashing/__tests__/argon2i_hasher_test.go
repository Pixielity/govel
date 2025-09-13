package tests

import (
	"strings"
	"testing"

	containerMocks "govel/container/mocks"
	"govel/hashing"
	configMocks "govel/config/mocks"
)

func TestArgon2iHasher_BasicFunctionality(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	password := "argonTest123"

	// Test hashing
	hash, err := hasher.Make(password, nil)
	if err != nil {
		t.Fatalf("Make should not return error: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2i$") {
		t.Error("Argon2i hash should start with $argon2i$")
	}

	// Test verification
	if !hasher.Check(password, hash, nil) {
		t.Error("Check should return true for correct password")
	}

	if hasher.Check("wrongPassword", hash, nil) {
		t.Error("Check should return false for incorrect password")
	}
}

func TestArgon2iHasher_WithConfiguration(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set argon2i configuration
	config.Set("hashing.argon2i", map[string]interface{}{
		"memory":  32768, // 32 KB for fast testing
		"time":    1,
		"threads": uint8(1),
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	password := "configTest"
	hash, err := hasher.Make(password, nil)
	if err != nil {
		t.Fatalf("Make with config should not return error: %v", err)
	}

	// Verify hash info
	info := hasher.Info(hash)
	if info.AlgoName != "argon2i" {
		t.Errorf("Algorithm should be argon2i, got %s", info.AlgoName)
	}

	if memory, ok := info.Options["memory_cost"].(int); !ok || memory != 32768 {
		t.Errorf("Memory should be 32768, got %v", info.Options["memory_cost"])
	}
}

func TestArgon2iHasher_WithOverrideOptions(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set default argon2i configuration
	config.Set("hashing.argon2i", map[string]interface{}{
		"memory":  65536,
		"time":    1,
		"threads": uint8(1),
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	password := "overrideTest"

	// Hash with overridden memory
	overrideHash, err := hasher.Make(password, map[string]interface{}{
		"memory": 32768, // Override default 65536
	})
	if err != nil {
		t.Fatalf("Make with override should not return error: %v", err)
	}

	overrideInfo := hasher.Info(overrideHash)
	if memory, ok := overrideInfo.Options["memory_cost"].(int); !ok || memory != 32768 {
		t.Errorf("Override memory should be 32768, got %v", overrideInfo.Options["memory_cost"])
	}

	// Should verify correctly
	if !hasher.Check(password, overrideHash, nil) {
		t.Error("Check should work with override hash")
	}
}

func TestArgon2iHasher_NeedsRehash(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	password := "rehashTest"

	// Create hash with specific parameters
	hash, err := hasher.Make(password, map[string]interface{}{
		"memory":  32768,
		"time":    1,
		"threads": uint8(1),
		"keyLen":  32,
	})
	if err != nil {
		t.Fatalf("Make should not return error: %v", err)
	}

	// Should not need rehash with same parameters
	if hasher.NeedsRehash(hash, map[string]interface{}{
		"memory":  32768,
		"time":    1,
		"threads": uint8(1),
		"keyLen":  32,
	}) {
		t.Error("Should not need rehash with same parameters")
	}

	// Should need rehash with different memory
	if !hasher.NeedsRehash(hash, map[string]interface{}{
		"memory":  65536,
		"time":    1,
		"threads": uint8(1),
		"keyLen":  32,
	}) {
		t.Error("Should need rehash with different memory")
	}

	// Should need rehash for invalid hash
	if !hasher.NeedsRehash("invalid-hash", nil) {
		t.Error("Should need rehash for invalid hash")
	}
}

func TestArgon2iHasher_InvalidOptions(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	password := "testPassword"

	// Test invalid memory
	_, err = hasher.Make(password, map[string]interface{}{
		"memory": 0, // Invalid memory
	})
	if err == nil {
		t.Error("Expected error for invalid argon2i memory")
	}

	// Test invalid threads
	_, err = hasher.Make(password, map[string]interface{}{
		"threads": 0, // Invalid threads
	})
	if err == nil {
		t.Error("Expected error for invalid argon2i threads")
	}

	// Test invalid time
	_, err = hasher.Make(password, map[string]interface{}{
		"time": 0, // Invalid time
	})
	if err == nil {
		t.Error("Expected error for invalid argon2i time")
	}
}

func TestArgon2iHasher_InvalidHash(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	invalidHashes := []string{
		"invalid-argon2i-hash",
		"$argon2i$invalid",
		"",
		"plaintext",
	}

	for _, invalidHash := range invalidHashes {
		t.Run(invalidHash, func(t *testing.T) {
			// Check should return false for invalid hashes
			if hasher.Check("password", invalidHash, nil) {
				t.Errorf("Check should return false for invalid hash: %s", invalidHash)
			}

			// NeedsRehash should return true for invalid hashes
			if !hasher.NeedsRehash(invalidHash, nil) {
				t.Errorf("NeedsRehash should return true for invalid hash: %s", invalidHash)
			}
		})
	}
}

func TestArgon2iHasher_EmptyPassword(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	// Test empty password
	hash, err := hasher.Make("", nil)
	if err != nil {
		t.Fatalf("Make should not return error for empty password: %v", err)
	}

	// Should be able to verify empty password
	if !hasher.Check("", hash, nil) {
		t.Error("Check should return true for empty password")
	}

	// Should return false for non-empty password against empty password hash
	if hasher.Check("nonempty", hash, nil) {
		t.Error("Check should return false for non-empty password against empty password hash")
	}
}

func TestArgon2iHasher_LongPassword(t *testing.T) {
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
	hasher := manager.Hasher("argon2i")

	if hasher == nil {
		t.Fatal("Argon2i hasher should not be nil")
	}

	// Create a very long password
	longPassword := strings.Repeat("a", 1000)

	_, err = hasher.Make(longPassword, nil)
	if err != nil {
		// Check if it's a "too long" error
		if !strings.Contains(err.Error(), "too long") {
			t.Errorf("Expected 'too long' error for argon2i, got: %v", err)
		}
	}
	// Note: Argon2i might handle long passwords fine, so err == nil is also valid
}
