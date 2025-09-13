package tests

import (
	"fmt"
	"strings"
	"testing"

	containerMocks "govel/packages/container/mocks"
	"govel/packages/hashing/src"
	configMocks "govel/packages/config/mocks"
	hashingInterfaces "govel/packages/types/src/interfaces/hashing"
)

func TestCompleteWorkflow(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set up comprehensive configuration
	config.Set("hashing.bcrypt", map[string]interface{}{
		"cost": 8, // Use lower cost for faster tests
	})
	config.Set("hashing.argon2i", map[string]interface{}{
		"memory":  32768, // 32 KB
		"time":    1,
		"threads": uint8(1),
		"keyLen":  32,
	})
	config.Set("hashing.argon2id", map[string]interface{}{
		"memory":  32768, // 32 KB
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

	passwords := []string{
		"simplePassword",
		"Complex@Password123!",
		"unicode-password-ñáéíóú-中文-🔐",
		"", // empty password
		"veryLongPasswordThatGoesOnAndOnAndOnToTestLengthHandling",
	}

	algorithms := []string{"bcrypt", "argon2i", "argon2id"}

	for _, password := range passwords {
		for _, algorithm := range algorithms {
			t.Run(fmt.Sprintf("%s_%s", algorithm, sanitizeTestName(password)), func(t *testing.T) {
				hasher := manager.Hasher(algorithm)
				if hasher == nil {
					t.Fatalf("Failed to get %s hasher", algorithm)
				}

				// Step 1: Hash the password
				hash, err := hasher.Make(password, nil)
				if err != nil {
					// For very long passwords, some algorithms may fail - that's acceptable
					if len(password) > 72 && algorithm == "bcrypt" {
						t.Logf("Expected error for long password in %s: %v", algorithm, err)
						return
					}
					t.Fatalf("Failed to hash password with %s: %v", algorithm, err)
				}

				// Step 2: Verify the hash format
				if !manager.IsHashed(hash) {
					t.Errorf("IsHashed should return true for generated hash")
				}

				// Step 3: Verify password matches
				if !hasher.Check(password, hash, nil) {
					t.Errorf("Check failed for correct password")
				}

				// Step 4: Verify wrong password doesn't match
				wrongPassword := password + "wrong"
				if hasher.Check(wrongPassword, hash, nil) {
					t.Errorf("Check should fail for incorrect password")
				}

				// Step 5: Extract hash info
				info := hasher.Info(hash)
				if info.AlgoName != algorithm {
					t.Errorf("Info should return algorithm name %s, got %s", algorithm, info.AlgoName)
				}

				if len(info.Options) == 0 {
					t.Errorf("Info should return options for %s", algorithm)
				}

				// Step 6: Test NeedsRehash with same options
				if hasher.NeedsRehash(hash, nil) {
					t.Errorf("NeedsRehash should return false for same configuration")
				}
			})
		}
	}
}

func TestCrossAlgorithmCompatibility(t *testing.T) {
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

	password := "crossTestPassword"

	// Generate hashes with different algorithms
	bcryptHasher := manager.Hasher("bcrypt")
	argon2iHasher := manager.Hasher("argon2i")
	argon2idHasher := manager.Hasher("argon2id")

	bcryptHash, _ := bcryptHasher.Make(password, map[string]interface{}{"cost": 8})
	argon2iHash, _ := argon2iHasher.Make(password, map[string]interface{}{
		"memory": 32768, "time": 1, "threads": uint8(1), "keyLen": 32,
	})
	argon2idHash, _ := argon2idHasher.Make(password, map[string]interface{}{
		"memory": 32768, "time": 1, "threads": uint8(1), "keyLen": 32,
	})

	testCases := []struct {
		hasher      hashingInterfaces.HasherInterface
		correctHash string
		wrongHashes []string
		name        string
	}{
		{
			hasher:      bcryptHasher,
			correctHash: bcryptHash,
			wrongHashes: []string{argon2iHash, argon2idHash},
			name:        "bcrypt",
		},
		{
			hasher:      argon2iHasher,
			correctHash: argon2iHash,
			wrongHashes: []string{bcryptHash, argon2idHash},
			name:        "argon2i",
		},
		{
			hasher:      argon2idHasher,
			correctHash: argon2idHash,
			wrongHashes: []string{bcryptHash, argon2iHash},
			name:        "argon2id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should verify correct hash
			if !tc.hasher.Check(password, tc.correctHash, nil) {
				t.Errorf("%s should verify its own hash", tc.name)
			}

			// Should not verify hashes from other algorithms
			for _, wrongHash := range tc.wrongHashes {
				if tc.hasher.Check(password, wrongHash, nil) {
					t.Errorf("%s should not verify hash from other algorithm", tc.name)
				}
			}
		})
	}
}

func TestConfigurationInheritanceAndOverrides(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set global configuration
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

	password := "configTest"

	// Test with global config
	globalHash, err := hasher.Make(password, nil)
	if err != nil {
		t.Fatalf("Failed with global config: %v", err)
	}

	globalInfo := hasher.Info(globalHash)
	if cost, ok := globalInfo.Options["cost"].(int); !ok || cost != 10 {
		t.Errorf("Global config cost should be 10, got %v", globalInfo.Options["cost"])
	}

	// Test with method-level override
	overrideHash, err := hasher.Make(password, map[string]interface{}{
		"cost": 8,
	})
	if err != nil {
		t.Fatalf("Failed with override config: %v", err)
	}

	overrideInfo := hasher.Info(overrideHash)
	if cost, ok := overrideInfo.Options["cost"].(int); !ok || cost != 8 {
		t.Errorf("Override config cost should be 8, got %v", overrideInfo.Options["cost"])
	}

	// Both should verify correctly
	if !hasher.Check(password, globalHash, nil) {
		t.Error("Global config hash should verify")
	}
	if !hasher.Check(password, overrideHash, nil) {
		t.Error("Override config hash should verify")
	}

	// Test NeedsRehash with different configs
	if !hasher.NeedsRehash(globalHash, map[string]interface{}{"cost": 8}) {
		t.Error("Should need rehash when cost changes")
	}
	if hasher.NeedsRehash(globalHash, map[string]interface{}{"cost": 10}) {
		t.Error("Should not need rehash with same cost")
	}
}

func TestManagerDirectMethods(t *testing.T) {
	// Create mock container with config
	container := containerMocks.NewMockContainer()
	config := configMocks.NewMockConfig()

	// Set default driver configuration
	config.Set("hashing.default", "bcrypt")
	config.Set("hashing.bcrypt", map[string]interface{}{
		"cost": 8,
	})

	// Bind config to container
	err := container.Singleton("config", func() interface{} {
		return config
	})
	if err != nil {
		t.Fatalf("Failed to bind config: %v", err)
	}

	manager := src.NewHashManager(container)
	password := "directMethodTest"

	// Test direct Make (should use default driver)
	hash, err := manager.Make(password, nil)
	if err != nil {
		t.Fatalf("Manager Make failed: %v", err)
	}

	// Test direct Check
	if !manager.Check(password, hash, nil) {
		t.Error("Manager Check failed for correct password")
	}

	if manager.Check("wrongPassword", hash, nil) {
		t.Error("Manager Check should fail for wrong password")
	}

	// Test direct IsHashed
	if !manager.IsHashed(hash) {
		t.Error("Manager IsHashed should return true for hash")
	}

	if manager.IsHashed("plaintext") {
		t.Error("Manager IsHashed should return false for plaintext")
	}

	// Test direct Info
	info := manager.Info(hash)
	if info.AlgoName == "" {
		t.Error("Manager Info should return algorithm name")
	}

	// Test direct NeedsRehash
	needsRehash := manager.NeedsRehash(hash, nil)
	// Just check it doesn't panic - result depends on implementation
	_ = needsRehash
}

func TestHashInfoExtraction(t *testing.T) {
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

	testCases := []struct {
		algorithm string
		prefix    string
	}{
		{"bcrypt", "$2"},
		{"argon2i", "$argon2i$"},
		{"argon2id", "$argon2id$"},
	}

	for _, tc := range testCases {
		t.Run(tc.algorithm, func(t *testing.T) {
			hasher := manager.Hasher(tc.algorithm)
			if hasher == nil {
				t.Fatalf("%s hasher should not be nil", tc.algorithm)
			}

			password := "infoTest"
			hash, err := hasher.Make(password, nil)
			if err != nil {
				t.Fatalf("Make should not return error: %v", err)
			}

			info := hasher.Info(hash)
			if info.AlgoName != tc.algorithm {
				t.Errorf("Algorithm name should be %s, got %s", tc.algorithm, info.AlgoName)
			}

			if len(info.Options) == 0 {
				t.Error("Options should not be empty")
			}
		})
	}
}

func TestIsHashedDetection(t *testing.T) {
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

	testCases := []struct {
		input    string
		expected bool
	}{
		{"plaintext", false},
		{"", false},
		{"$2y$12$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", true},
		{"$argon2i$v=19$m=65536,t=1,p=1$c2FsdA$hash", true},
		{"$argon2id$v=19$m=131072,t=2,p=1$c2FsdA$hash", true},
		{"invalid$hash$format", false},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			result := manager.IsHashed(tc.input)
			if result != tc.expected {
				t.Errorf("IsHashed(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestErrorPropagation(t *testing.T) {
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

	// Test unsupported algorithm error propagation
	unsupportedHasher := manager.Hasher("unsupported")
	if unsupportedHasher != nil {
		t.Error("Unsupported algorithm should return nil")
	}

	// Test invalid options error propagation
	bcryptHasher := manager.Hasher("bcrypt")
	if bcryptHasher == nil {
		t.Fatal("Bcrypt hasher should not be nil")
	}

	_, err = bcryptHasher.Make("password", map[string]interface{}{
		"cost": 100, // Invalid cost
	})
	if err == nil {
		t.Error("Invalid options should return error")
	}
}

// Helper function to sanitize test names
func sanitizeTestName(input string) string {
	if input == "" {
		return "empty"
	}

	// Replace problematic characters for test names
	result := strings.ReplaceAll(input, " ", "_")
	result = strings.ReplaceAll(result, "@", "at")
	result = strings.ReplaceAll(result, "!", "excl")
	result = strings.ReplaceAll(result, "#", "hash")
	result = strings.ReplaceAll(result, "$", "dollar")
	result = strings.ReplaceAll(result, "%", "percent")
	result = strings.ReplaceAll(result, "^", "caret")
	result = strings.ReplaceAll(result, "&", "amp")
	result = strings.ReplaceAll(result, "*", "star")
	result = strings.ReplaceAll(result, "(", "lparen")
	result = strings.ReplaceAll(result, ")", "rparen")
	result = strings.ReplaceAll(result, "+", "plus")
	result = strings.ReplaceAll(result, "=", "eq")
	result = strings.ReplaceAll(result, "{", "lbrace")
	result = strings.ReplaceAll(result, "}", "rbrace")
	result = strings.ReplaceAll(result, "[", "lbracket")
	result = strings.ReplaceAll(result, "]", "rbracket")
	result = strings.ReplaceAll(result, "|", "pipe")
	result = strings.ReplaceAll(result, "\\", "backslash")
	result = strings.ReplaceAll(result, ":", "colon")
	result = strings.ReplaceAll(result, ";", "semicolon")
	result = strings.ReplaceAll(result, "\"", "quote")
	result = strings.ReplaceAll(result, "'", "apostrophe")
	result = strings.ReplaceAll(result, "<", "lt")
	result = strings.ReplaceAll(result, ">", "gt")
	result = strings.ReplaceAll(result, ",", "comma")
	result = strings.ReplaceAll(result, ".", "dot")
	result = strings.ReplaceAll(result, "?", "question")
	result = strings.ReplaceAll(result, "/", "slash")
	result = strings.ReplaceAll(result, "~", "tilde")
	result = strings.ReplaceAll(result, "`", "backtick")

	// Truncate if too long
	if len(result) > 50 {
		result = result[:47] + "..."
	}

	return result
}
