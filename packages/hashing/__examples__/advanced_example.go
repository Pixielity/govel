package main

import (
	"errors"
	"fmt"
	"log"

	"govel/hashing/src"
	"govel/hashing/src/exceptions"
)

// AdvancedContainer provides more sophisticated configuration
type AdvancedContainer struct{}

func (c *AdvancedContainer) Get(key string) (interface{}, bool) {
	return nil, false // No configuration - use defaults
}

func main() {
	fmt.Println("=== Advanced Hashing Example ===")

	// Create container and manager
	container := &AdvancedContainer{}
	manager := src.NewHashManager(container)

	// Example 1: Error handling with exceptions
	fmt.Println("1. Exception Handling:")
	demonstrateExceptionHandling(manager)
	fmt.Println()

	// Example 2: Runtime parameter overrides
	fmt.Println("2. Runtime Parameter Overrides:")
	demonstrateRuntimeOverrides(manager)
	fmt.Println()

	// Example 3: Hash migration and rehashing
	fmt.Println("3. Hash Migration and Rehashing:")
	demonstrateHashMigration(manager)
	fmt.Println()

	// Example 4: Hash analysis and inspection
	fmt.Println("4. Hash Analysis and Inspection:")
	demonstrateHashAnalysis(manager)
	fmt.Println()

	// Example 5: Invalid algorithm handling
	fmt.Println("5. Invalid Algorithm Handling:")
	demonstrateInvalidAlgorithm(manager)
}

func demonstrateExceptionHandling(manager *src.HashManager) {
	bcryptHasher := manager.Hasher("bcrypt")
	if bcryptHasher == nil {
		fmt.Println("❌ Failed to get bcrypt hasher")
		return
	}

	// Test ErrValueTooLong
	longPassword := make([]byte, 1000) // Very long password
	for i := range longPassword {
		longPassword[i] = 'a'
	}

	_, err := bcryptHasher.Make(string(longPassword), nil)
	if err != nil {
		if errors.Is(err, exceptions.ErrValueTooLong) {
			fmt.Printf("✅ ErrValueTooLong caught: %v\n", err)
		} else {
			fmt.Printf("❓ Unexpected error: %v\n", err)
		}
	}

	// Test ErrInvalidOptions
	_, err = bcryptHasher.Make("password", map[string]interface{}{
		"cost": 100, // Invalid cost (too high)
	})
	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidOptions) {
			fmt.Printf("✅ ErrInvalidOptions caught: %v\n", err)
		} else {
			fmt.Printf("❓ Unexpected error: %v\n", err)
		}
	}

	// Test ErrInvalidHash
	isValid := bcryptHasher.Check("password", "invalid-hash-format", nil)
	fmt.Printf("✅ Invalid hash verification (should be false): %t\n", isValid)
}

func demonstrateRuntimeOverrides(manager *src.HashManager) {
	password := "testPassword"

	// Bcrypt with different costs
	bcryptHasher := manager.Hasher("bcrypt")
	if bcryptHasher != nil {
		// Low cost (fast, less secure)
		hashLow, err := bcryptHasher.Make(password, map[string]interface{}{
			"cost": 8,
		})
		if err != nil {
			log.Printf("Error with low cost: %v", err)
		} else {
			info := bcryptHasher.Info(hashLow)
			fmt.Printf("Low cost hash (cost %v): %s\n", info.Options["cost"], hashLow[:20]+"...")
		}

		// High cost (slow, more secure)
		hashHigh, err := bcryptHasher.Make(password, map[string]interface{}{
			"cost": 14,
		})
		if err != nil {
			log.Printf("Error with high cost: %v", err)
		} else {
			info := bcryptHasher.Info(hashHigh)
			fmt.Printf("High cost hash (cost %v): %s\n", info.Options["cost"], hashHigh[:20]+"...")
		}
	}

	// Argon2id with different parameters
	argonHasher := manager.Hasher("argon2id")
	if argonHasher != nil {
		// High security configuration
		hashSecure, err := argonHasher.Make(password, map[string]interface{}{
			"memory":  uint32(256 * 1024), // 256MB
			"time":    uint32(4),          // 4 iterations
			"threads": uint8(2),           // 2 threads
			"keyLen":  uint32(64),         // 64-byte output
		})
		if err != nil {
			log.Printf("Error with secure config: %v", err)
		} else {
			info := argonHasher.Info(hashSecure)
			fmt.Printf("Secure Argon2id (mem: %v KB, time: %v): %s\n",
				info.Options["memory_cost"], info.Options["time_cost"], hashSecure[:30]+"...")
		}
	}
}

func demonstrateHashMigration(manager *src.HashManager) {
	password := "migrateMe"

	// Create hash with old parameters
	bcryptHasher := manager.Hasher("bcrypt")
	if bcryptHasher != nil {
		oldHash, err := bcryptHasher.Make(password, map[string]interface{}{
			"cost": 10, // Old cost
		})
		if err != nil {
			log.Printf("Error creating old hash: %v", err)
			return
		}

		fmt.Printf("Old hash (cost 10): %s\n", oldHash[:20]+"...")

		// Check if rehash needed with new parameters
		needsRehash := bcryptHasher.NeedsRehash(oldHash, map[string]interface{}{
			"cost": 13, // New higher cost
		})
		fmt.Printf("Needs rehash with cost 13: %t\n", needsRehash)

		if needsRehash {
			// Create new hash with updated parameters
			newHash, err := bcryptHasher.Make(password, map[string]interface{}{
				"cost": 13,
			})
			if err != nil {
				log.Printf("Error creating new hash: %v", err)
			} else {
				fmt.Printf("New hash (cost 13): %s\n", newHash[:20]+"...")

				// Verify both hashes work
				fmt.Printf("Old hash verifies: %t\n", bcryptHasher.Check(password, oldHash, nil))
				fmt.Printf("New hash verifies: %t\n", bcryptHasher.Check(password, newHash, nil))
			}
		}
	}
}

func demonstrateHashAnalysis(manager *src.HashManager) {
	testHashes := map[string]string{
		"bcrypt":   "$2y$12$abcdefghijklmnopqrstuvwxyz012345678901234567890123",
		"argon2i":  "$argon2i$v=19$m=65536,t=1,p=1$c2FsdA$hash",
		"argon2id": "$argon2id$v=19$m=131072,t=3,p=2$c2FsdA$hash",
	}

	for algoName, hash := range testHashes {
		fmt.Printf("\n--- Analyzing %s hash ---\n", algoName)

		// Use manager's Info method (detects algorithm automatically)
		info := manager.Info(hash)
		fmt.Printf("Detected algorithm: %s\n", info.AlgoName)

		if info.AlgoName != "" {
			fmt.Printf("Algorithm ID: %s\n", info.Algo)
			fmt.Printf("Options:\n")
			for key, value := range info.Options {
				fmt.Printf("  %s: %v\n", key, value)
			}

			// Check if it looks like a hash
			fmt.Printf("Is hashed: %t\n", manager.IsHashed(hash))
		} else {
			fmt.Println("❌ Unknown or invalid hash format")
		}
	}
}

func demonstrateInvalidAlgorithm(manager *src.HashManager) {
	// Try to get an invalid algorithm
	invalidHasher := manager.Hasher("md5") // Invalid algorithm
	if invalidHasher == nil {
		fmt.Println("✅ Correctly rejected invalid algorithm 'md5'")
	} else {
		fmt.Println("❌ Should have rejected invalid algorithm")
	}

	// Try empty algorithm name
	emptyHasher := manager.Hasher("")
	if emptyHasher != nil {
		fmt.Println("✅ Empty algorithm name uses default hasher")

		// Test it works
		hash, err := emptyHasher.Make("test", nil)
		if err != nil {
			log.Printf("Error with default hasher: %v", err)
		} else {
			fmt.Printf("Default hasher works: %s\n", hash[:20]+"...")
		}
	} else {
		fmt.Println("❌ Empty algorithm should use default")
	}
}
