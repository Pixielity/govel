package main

import (
	"fmt"
	"log"

	"govel/packages/hashing/src"
)

// MockContainer implements a simple container for testing
type MockContainer struct{}

func (c *MockContainer) Get(key string) (interface{}, bool) {
	// Return default configurations for different algorithms
	switch key {
	case "hashing.bcrypt":
		return map[string]interface{}{
			"cost": 12, // Higher security
		}, true
	case "hashing.argon2i":
		return map[string]interface{}{
			"memory":  uint32(64 * 1024), // 64MB
			"time":    uint32(1),         // 1 iteration
			"threads": uint8(1),          // 1 thread
			"keyLen":  uint32(32),        // 32 bytes
		}, true
	case "hashing.argon2id":
		return map[string]interface{}{
			"memory":  uint32(128 * 1024), // 128MB for higher security
			"time":    uint32(3),          // 3 iterations
			"threads": uint8(2),           // 2 threads
			"keyLen":  uint32(32),         // 32 bytes
		}, true
	}
	return nil, false
}

func main() {
	fmt.Println("=== Basic Hashing Example ===")

	// Create a mock container with configuration
	container := &MockContainer{}

	// Create hash manager
	manager := src.NewHashManager(container)

	// Password to hash
	password := "mySecurePassword123"

	fmt.Println("Original password:", password)
	fmt.Println()

	// Example 1: Using default hasher
	fmt.Println("1. Using default hasher:")
	defaultHasher := manager.Hasher()
	if defaultHasher != nil {
		hash, err := defaultHasher.Make(password, nil)
		if err != nil {
			log.Printf("Error hashing with default: %v", err)
		} else {
			fmt.Printf("Default hash: %s\n", hash)
			fmt.Printf("Verification: %t\n", defaultHasher.Check(password, hash, nil))
		}
	}
	fmt.Println()

	// Example 2: Using bcrypt hasher
	fmt.Println("2. Using bcrypt hasher:")
	bcryptHasher := manager.Hasher("bcrypt")
	if bcryptHasher != nil {
		hash, err := bcryptHasher.Make(password, nil)
		if err != nil {
			log.Printf("Error hashing with bcrypt: %v", err)
		} else {
			fmt.Printf("Bcrypt hash: %s\n", hash)
			fmt.Printf("Verification: %t\n", bcryptHasher.Check(password, hash, nil))

			// Get hash info
			info := bcryptHasher.Info(hash)
			fmt.Printf("Algorithm: %s, Cost: %v\n", info.AlgoName, info.Options["cost"])
		}
	}
	fmt.Println()

	// Example 3: Using Argon2i hasher
	fmt.Println("3. Using Argon2i hasher:")
	argonHasher := manager.Hasher("argon2i")
	if argonHasher != nil {
		hash, err := argonHasher.Make(password, nil)
		if err != nil {
			log.Printf("Error hashing with Argon2i: %v", err)
		} else {
			fmt.Printf("Argon2i hash: %s\n", hash)
			fmt.Printf("Verification: %t\n", argonHasher.Check(password, hash, nil))

			// Get hash info
			info := argonHasher.Info(hash)
			fmt.Printf("Algorithm: %s\n", info.AlgoName)
			fmt.Printf("Memory: %v KB, Time: %v, Threads: %v\n",
				info.Options["memory_cost"], info.Options["time_cost"], info.Options["threads"])
		}
	}
	fmt.Println()

	// Example 4: Using Argon2id hasher (recommended)
	fmt.Println("4. Using Argon2id hasher (recommended):")
	argon2idHasher := manager.Hasher("argon2id")
	if argon2idHasher != nil {
		hash, err := argon2idHasher.Make(password, nil)
		if err != nil {
			log.Printf("Error hashing with Argon2id: %v", err)
		} else {
			fmt.Printf("Argon2id hash: %s\n", hash)
			fmt.Printf("Verification: %t\n", argon2idHasher.Check(password, hash, nil))

			// Get hash info
			info := argon2idHasher.Info(hash)
			fmt.Printf("Algorithm: %s\n", info.AlgoName)
			fmt.Printf("Memory: %v KB, Time: %v, Threads: %v\n",
				info.Options["memory_cost"], info.Options["time_cost"], info.Options["threads"])
		}
	}
	fmt.Println()

	// Example 5: Using manager directly (delegates to default)
	fmt.Println("5. Using manager directly:")
	hash, err := manager.Make(password, nil)
	if err != nil {
		log.Printf("Error hashing with manager: %v", err)
	} else {
		fmt.Printf("Manager hash: %s\n", hash)
		fmt.Printf("Verification: %t\n", manager.Check(password, hash, nil))
		fmt.Printf("Is hashed: %t\n", manager.IsHashed(hash))
		fmt.Printf("Needs rehash: %t\n", manager.NeedsRehash(hash, nil))
	}
}
