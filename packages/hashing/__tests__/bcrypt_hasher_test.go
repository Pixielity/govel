package tests

import (
	"strings"
	"testing"

	configMocks "govel/config/mocks"
	containerMocks "govel/container/mocks"
	"govel/hashing/src"
)

func TestBcryptHasher_BasicFunctionality(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	password := "testPassword123"

	// Test hashing
	hash, err := hasher.Make(password, nil)
	if err != nil {
		t.Fatalf("Make should not return error: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if !strings.HasPrefix(hash, "$2") {
		t.Error("Bcrypt hash should start with $2")
	}

	// Test verification
	if !hasher.Check(password, hash, nil) {
		t.Error("Check should return true for correct password")
	}

	if hasher.Check("wrongPassword", hash, nil) {
		t.Error("Check should return false for incorrect password")
	}
}

func TestBcryptHasher_WithConfiguration(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set bcrypt configuration
	config.Set("hashing.bcrypt", map[string]interface{}{
		"cost": 8, // Low cost for fast testing
	})

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	password := "configTest"
	hash, err := hasher.Make(password, nil)
	if err != nil {
		t.Fatalf("Make with config should not return error: %v", err)
	}

	// Verify hash info
	info := hasher.Info(hash)
	if info.AlgoName != "bcrypt" {
		t.Errorf("Algorithm should be bcrypt, got %s", info.AlgoName)
	}

	if cost, ok := info.Options["cost"].(int); !ok || cost != 8 {
		t.Errorf("Cost should be 8, got %v", info.Options["cost"])
	}
}

func TestBcryptHasher_WithOverrideOptions(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set default bcrypt configuration
	config.Set("hashing.bcrypt", map[string]interface{}{
		"cost": 10,
	})

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	password := "overrideTest"

	// Hash with overridden cost
	overrideHash, err := hasher.Make(password, map[string]interface{}{
		"cost": 8,
	})
	if err != nil {
		t.Fatalf("Make with override should not return error: %v", err)
	}

	overrideInfo := hasher.Info(overrideHash)
	if cost, ok := overrideInfo.Options["cost"].(int); !ok || cost != 8 {
		t.Errorf("Override cost should be 8, got %v", overrideInfo.Options["cost"])
	}

	// Should verify correctly
	if !hasher.Check(password, overrideHash, nil) {
		t.Error("Check should work with override hash")
	}
}

func TestBcryptHasher_NeedsRehash(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	password := "rehashTest"

	// Create hash with cost 8
	hash, err := hasher.Make(password, map[string]interface{}{
		"cost": 8,
	})
	if err != nil {
		t.Fatalf("Make should not return error: %v", err)
	}

	// Should not need rehash with same cost
	if hasher.NeedsRehash(hash, map[string]interface{}{"cost": 8}) {
		t.Error("Should not need rehash with same cost")
	}

	// Should need rehash with different cost
	if !hasher.NeedsRehash(hash, map[string]interface{}{"cost": 12}) {
		t.Error("Should need rehash with different cost")
	}

	// Should need rehash for invalid hash
	if !hasher.NeedsRehash("invalid-hash", nil) {
		t.Error("Should need rehash for invalid hash")
	}
}

func TestBcryptHasher_InvalidOptions(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	password := "testPassword"

	// Test invalid high cost
	_, err = hasher.Make(password, map[string]interface{}{
		"cost": 100, // Invalid high cost
	})
	if err == nil {
		t.Error("Expected error for invalid high bcrypt cost")
	}

	// Test invalid low cost
	_, err = hasher.Make(password, map[string]interface{}{
		"cost": 3, // Invalid low cost
	})
	if err == nil {
		t.Error("Expected error for invalid low bcrypt cost")
	}
}

func TestBcryptHasher_InvalidHash(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	invalidHashes := []string{
		"invalid-bcrypt-hash",
		"$2y$invalid",
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

func TestBcryptHasher_EmptyPassword(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
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

func TestBcryptHasher_LongPassword(t *testing.T) {
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
	hasher := manager.Hasher("bcrypt")

	if hasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	// Create a very long password (> 72 bytes for bcrypt)
	longPassword := strings.Repeat("a", 100)

	_, err = hasher.Make(longPassword, nil)
	if err != nil {
		// Check if it's a "too long" error
		if !strings.Contains(err.Error(), "too long") {
			t.Errorf("Expected 'too long' error for bcrypt, got: %v", err)
		}
	}
	// Note: bcrypt might handle long passwords by truncating, so err == nil is also valid
}
