package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	config "govel/config"
	"govel/config/drivers"
)

func main() {
	fmt.Println("=== GoVel Configuration Package Demo ===\n")

	// // Set some environment variables for demonstration
	// setupEnvironmentVariables()

	// // Example 1: File-based configuration (YAML)
	// fmt.Println("1. File Driver Example (YAML)")
	// fileDriverExample()

	// // Example 2: File-based configuration (JSON)
	// fmt.Println("\n2. File Driver Example (JSON)")
	// jsonFileDriverExample()

	// // Example 3: File-based configuration (TOML)
	// fmt.Println("\n3. File Driver Example (TOML)")
	// tomlFileDriverExample()

	// // Example 4: Environment variable configuration
	// fmt.Println("\n4. Environment Driver Example")
	// envDriverExample()

	// // Example 5: Memory configuration
	// fmt.Println("\n5. Memory Driver Example")
	// memoryDriverExample()

	// // Example 6: Remote configuration (stub)
	// fmt.Println("\n6. Remote Driver Example")
	// remoteDriverExample()

	// // Example 7: Config watching
	// fmt.Println("\n7. Configuration Watching Example")
	// watchingExample()

	// // Example 8: Type-safe configuration access
	// fmt.Println("\n8. Type-safe Configuration Access")
	// typeSafeExample()

	// // Example 9: Driver comparison
	// fmt.Println("\n9. Driver Comparison")
	// driverComparisonExample()

	// Example 10: Go file driver configuration (app.go)
	fmt.Println("\n10. Go File Driver Configuration (app.go)")
	goFileConfigExample()

	fmt.Println("\n=== Demo Complete ===")
}

func setupEnvironmentVariables() {
	// Set up some environment variables for demo
	os.Setenv("APP_NAME", "GoVel from Environment")
	os.Setenv("APP_DEBUG", "false")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("APP_DATABASE_HOST", "env-database-host")
	os.Setenv("APP_DATABASE_PORT", "5432")
	os.Setenv("APP_CACHE_ENABLED", "true")
	os.Setenv("APP_CACHE_TTL", "7200")
}

func fileDriverExample() {
	// Create file driver with YAML configuration
	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "config",
		ConfigType:  "yaml",
	})

	config := config.NewConfig(driver)

	// Load the configuration
	if err := config.Load(); err != nil {
		log.Printf("Failed to load YAML config: %v", err)
		return
	}

	// Access configuration values
	appName := config.GetString("app.name", "Default App")
	debug := config.GetBool("app.debug", false)
	port := config.GetInt("server.port", 8080)
	timeout := config.GetDuration("server.timeout", 30*time.Second)

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Debug Mode: %t\n", debug)
	fmt.Printf("  Server Port: %d\n", port)
	fmt.Printf("  Timeout: %v\n", timeout)

	// Access nested configuration
	dbHost := config.GetString("database.host", "localhost")
	dbPort := config.GetInt("database.port", 5432)
	maxConnections := config.GetInt("database.max_connections", 10)

	fmt.Printf("  Database: %s:%d (max connections: %d)\n", dbHost, dbPort, maxConnections)

	// Access array values
	socialLogins := config.GetStringSlice("features.social_login", []string{})
	fmt.Printf("  Social Logins: %v\n", socialLogins)
}

func jsonFileDriverExample() {
	// Create file driver with JSON configuration
	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "config",
		ConfigType:  "json",
	})

	config := config.NewConfig(driver)

	// Load the configuration
	if err := config.Load(); err != nil {
		log.Printf("Failed to load JSON config: %v", err)
		return
	}

	// Access configuration values
	appName := config.GetString("app.name", "Default App")
	corsEnabled := config.GetBool("server.cors.enabled", false)
	corsOrigins := config.GetStringSlice("server.cors.origins", []string{})

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  CORS Enabled: %t\n", corsEnabled)
	fmt.Printf("  CORS Origins: %v\n", corsOrigins)

	// Access nested database configuration
	defaultDB := config.GetString("database.default", "mysql")
	mysqlHost := config.GetString("database.connections.mysql.host", "localhost")
	mysqlPort := config.GetInt("database.connections.mysql.port", 3306)

	fmt.Printf("  Default DB: %s\n", defaultDB)
	fmt.Printf("  MySQL: %s:%d\n", mysqlHost, mysqlPort)
}

func tomlFileDriverExample() {
	// Create file driver with TOML configuration
	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "config",
		ConfigType:  "toml",
	})

	config := config.NewConfig(driver)

	// Load the configuration
	if err := config.Load(); err != nil {
		log.Printf("Failed to load TOML config: %v", err)
		return
	}

	// Access configuration values
	appName := config.GetString("app.name", "Default App")
	cacheDefault := config.GetString("cache.default", "memory")
	redisHost := config.GetString("cache.stores.redis.host", "localhost")
	redisPort := config.GetInt("cache.stores.redis.port", 6379)

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Default Cache: %s\n", cacheDefault)
	fmt.Printf("  Redis: %s:%d\n", redisHost, redisPort)

	// Access logging configuration
	logLevel := config.GetString("logging.level", "info")
	logFormat := config.GetString("logging.format", "text")

	fmt.Printf("  Log Level: %s\n", logLevel)
	fmt.Printf("  Log Format: %s\n", logFormat)
}

func envDriverExample() {
	// Create environment driver
	driver := drivers.NewEnvDriver(&drivers.EnvDriverOptions{
		Prefix:       "APP_",
		AutomaticEnv: true,
		Replacer:     strings.NewReplacer(".", "_", "-", "_"),
	})

	config := config.NewConfig(driver)

	// Access environment variables
	appName := config.GetString("name", "Default App")
	debug := config.GetBool("debug", true)
	port := config.GetInt("port", 8080)
	dbHost := config.GetString("database.host", "localhost")
	dbPort := config.GetInt("database.port", 5432)
	cacheEnabled := config.GetBool("cache.enabled", false)
	cacheTTL := config.GetInt("cache.ttl", 3600)

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Debug Mode: %t\n", debug)
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  Database: %s:%d\n", dbHost, dbPort)
	fmt.Printf("  Cache Enabled: %t (TTL: %ds)\n", cacheEnabled, cacheTTL)

	// Demonstrate setting environment variables at runtime
	config.Set("runtime.test", "environment value")
	runtimeValue := config.GetString("runtime.test", "default")
	fmt.Printf("  Runtime Value: %s\n", runtimeValue)
}

func memoryDriverExample() {
	// Create memory driver with initial data
	initialData := map[string]interface{}{
		"app.name":        "Memory Config App",
		"app.debug":       true,
		"server.port":     9000,
		"database.host":   "memory-db-host",
		"database.port":   5432,
		"features.cache":  true,
		"features.queue":  false,
		"limits.requests": 1000,
		"limits.users":    []string{"admin", "user", "guest"},
	}

	driver := drivers.NewMemoryDriver(&drivers.MemoryDriverOptions{
		InitialData: initialData,
	})

	config := config.NewConfig(driver)

	// Access initial configuration
	appName := config.GetString("app.name")
	port := config.GetInt("server.port")
	dbHost := config.GetString("database.host")
	users := config.GetStringSlice("limits.users")

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  DB Host: %s\n", dbHost)
	fmt.Printf("  Users: %v\n", users)

	// Modify configuration at runtime
	config.Set("runtime.timestamp", time.Now().Format(time.RFC3339))
	config.Set("runtime.counter", 42)
	config.Set("features.analytics", true)

	timestamp := config.GetString("runtime.timestamp")
	counter := config.GetInt("runtime.counter")
	analytics := config.GetBool("features.analytics")

	fmt.Printf("  Runtime Timestamp: %s\n", timestamp)
	fmt.Printf("  Runtime Counter: %d\n", counter)
	fmt.Printf("  Analytics: %t\n", analytics)

	// Demonstrate additional memory operations
	memDriver := driver // Type assertion would be needed in real code

	// Load additional data
	additionalData := map[string]interface{}{
		"external.api.key":     "secret-api-key",
		"external.api.timeout": "30s",
		"external.retries":     3,
	}

	memDriver.LoadFromMap(additionalData)

	apiKey := config.GetString("external.api.key")
	apiTimeout := config.GetString("external.api.timeout")
	retries := config.GetInt("external.retries")

	fmt.Printf("  API Key: %s\n", apiKey)
	fmt.Printf("  API Timeout: %s\n", apiTimeout)
	fmt.Printf("  Retries: %d\n", retries)
}

func remoteDriverExample() {
	// Create remote driver (this is a stub implementation)
	driver := drivers.NewRemoteDriver(&drivers.RemoteDriverOptions{
		Provider: "etcd",
		Endpoint: "http://localhost:2379",
		KeyPath:  "/config",
	})

	config := config.NewConfig(driver)

	// Try to load remote configuration (will fail with stub)
	if err := config.Load(); err != nil {
		fmt.Printf("  Remote config load failed (expected): %v\n", err)
	}

	// Demonstrate setting values in remote driver
	config.Set("remote.demo", "remote value")
	config.Set("remote.timestamp", time.Now().Unix())

	demoValue := config.GetString("remote.demo", "default")
	timestamp := config.GetInt64("remote.timestamp", 0)

	fmt.Printf("  Demo Value: %s\n", demoValue)
	fmt.Printf("  Timestamp: %d\n", timestamp)

	fmt.Printf("  Note: This is a stub implementation. Real remote drivers would connect to etcd, Consul, etc.\n")
}

func watchingExample() {
	// Create file driver for watching
	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "config",
		ConfigType:  "yaml",
	})

	config := config.NewConfig(driver)

	// Load initial configuration
	if err := config.Load(); err != nil {
		log.Printf("Failed to load config for watching: %v", err)
		return
	}

	// Set up watching
	fmt.Printf("  Setting up file watching...\n")

	err := config.Watch(func() {
		fmt.Printf("  📝 Configuration file changed!\n")
		fmt.Printf("  📊 New app name: %s\n", config.GetString("app.name", "Unknown"))
	})

	if err != nil {
		fmt.Printf("  Failed to set up watching: %v\n", err)
		return
	}

	fmt.Printf("  ✅ Watching enabled. Try modifying the config/config.yaml file.\n")
	fmt.Printf("  ⏰ Waiting 3 seconds to demonstrate...\n")

	// Wait a bit to show watching is active
	time.Sleep(3 * time.Second)

	// Stop watching
	if err := config.Unwatch(); err != nil {
		fmt.Printf("  Failed to stop watching: %v\n", err)
	} else {
		fmt.Printf("  ⏹️  Watching stopped.\n")
	}
}

func typeSafeExample() {
	// Create memory driver for type safety demonstration
	driver := drivers.NewMemoryDriver(&drivers.MemoryDriverOptions{
		InitialData: map[string]interface{}{
			"string_value":   "hello world",
			"int_value":      42,
			"int64_value":    int64(9223372036854775807),
			"float_value":    3.14159,
			"bool_value":     true,
			"duration_value": "30s",
			"time_value":     "2023-01-01T00:00:00Z",
			"string_slice":   []string{"apple", "banana", "cherry"},
			"int_slice":      []int{1, 2, 3, 4, 5},
			"string_map": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			"string_map_string": map[string]string{
				"name":  "John",
				"email": "john@example.com",
			},
		},
	})

	config := config.NewConfig(driver)

	// Demonstrate type-safe access with spf13/cast
	fmt.Printf("  String Value: %s\n", config.GetString("string_value"))
	fmt.Printf("  Int Value: %d\n", config.GetInt("int_value"))
	fmt.Printf("  Int64 Value: %d\n", config.GetInt64("int64_value"))
	fmt.Printf("  Float Value: %.5f\n", config.GetFloat64("float_value"))
	fmt.Printf("  Bool Value: %t\n", config.GetBool("bool_value"))
	fmt.Printf("  Duration Value: %v\n", config.GetDuration("duration_value"))
	fmt.Printf("  Time Value: %v\n", config.GetTime("time_value"))
	fmt.Printf("  String Slice: %v\n", config.GetStringSlice("string_slice"))
	fmt.Printf("  Int Slice: %v\n", config.GetIntSlice("int_slice"))
	fmt.Printf("  String Map: %v\n", config.GetStringMap("string_map"))
	fmt.Printf("  String Map String: %v\n", config.GetStringMapString("string_map_string"))

	// Demonstrate type conversion with defaults
	fmt.Printf("\n  With defaults:\n")
	fmt.Printf("  Missing String: %s\n", config.GetString("missing_string", "default value"))
	fmt.Printf("  Missing Int: %d\n", config.GetInt("missing_int", 100))
	fmt.Printf("  Missing Bool: %t\n", config.GetBool("missing_bool", false))

	// Demonstrate type conversion safety
	fmt.Printf("\n  Type conversion safety:\n")
	fmt.Printf("  String as Int: %d\n", config.GetInt("string_value", 0)) // "hello world" -> 0
	fmt.Printf("  Int as String: %s\n", config.GetString("int_value"))    // 42 -> "42"
	fmt.Printf("  Float as Int: %d\n", config.GetInt("float_value"))      // 3.14159 -> 3
}

func driverComparisonExample() {
	fmt.Printf("  Comparing different drivers with the same interface:\n\n")

	// File driver
	fileDriver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "config",
		ConfigType:  "yaml",
	})
	fileConfig := config.NewConfig(fileDriver)
	if err := fileConfig.Load(); err == nil {
		fmt.Printf("  📁 File Driver - App Name: %s\n", fileConfig.GetString("app.name", "N/A"))
	} else {
		fmt.Printf("  📁 File Driver - Failed to load: %v\n", err)
	}

	// Environment driver
	envDriver := drivers.NewEnvDriver(&drivers.EnvDriverOptions{
		Prefix:       "APP_",
		AutomaticEnv: true,
	})
	envConfig := config.NewConfig(envDriver)
	fmt.Printf("  🌍 Env Driver - App Name: %s\n", envConfig.GetString("name", "N/A"))

	// Memory driver
	memoryDriver := drivers.NewMemoryDriver(&drivers.MemoryDriverOptions{
		InitialData: map[string]interface{}{
			"app.name": "Memory Config App",
		},
	})
	memoryConfig := config.NewConfig(memoryDriver)
	fmt.Printf("  🧠 Memory Driver - App Name: %s\n", memoryConfig.GetString("app.name", "N/A"))

	// Remote driver (stub)
	remoteDriver := drivers.NewRemoteDriver(&drivers.RemoteDriverOptions{
		Provider: "etcd",
		Endpoint: "http://localhost:2379",
	})
	remoteConfig := config.NewConfig(remoteDriver)
	remoteConfig.Set("app.name", "Remote Config App")
	fmt.Printf("  🌐 Remote Driver - App Name: %s\n", remoteConfig.GetString("app.name", "N/A"))

	fmt.Printf("\n  All drivers implement the same interface and provide consistent access patterns!\n")
}

func goFileConfigExample() {
	fmt.Printf("  Demonstrating Go file-based configuration with config.go:\n\n")

	// Create file driver specifically for Go files
	driver := drivers.NewFileDriver(&drivers.FileDriverOptions{
		ConfigPaths: []string{"./config"},
		ConfigName:  "app",
		ConfigType:  "go",
	})

	config := config.NewConfig(driver)

	// Load the Go configuration file
	if err := config.Load(); err != nil {
		log.Printf("Failed to load Go config: %v", err)
		return
	}

	// Access configuration values from the Go file
	appName := config.GetString("app.name")
	appVersion := config.GetString("app.version", "2.0.0")
	appEnv := config.GetString("app.env", "production")
	debug := config.GetBool("app.debug", false)
	appURL := config.GetString("app.url", "")
	timezone := config.GetString("app.timezone", "UTC")

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Version: %s\n", appVersion)
	fmt.Printf("  Environment: %s\n", appEnv)
	fmt.Printf("  Debug Mode: %t\n", debug)
	fmt.Printf("  URL: %s\n", appURL)
	fmt.Printf("  Timezone: %s\n", timezone)

	// Access database configuration
	dbDriver := config.GetString("database.driver", "sqlite")
	dbHost := config.GetString("database.host", "localhost")
	dbPort := config.GetInt("database.port", 5432)
	dbName := config.GetString("database.name", "govel")
	dbUser := config.GetString("database.user", "root")
	dbSSL := config.GetBool("database.ssl", false)

	fmt.Printf("\n  Database Configuration:\n")
	fmt.Printf("    Driver: %s\n", dbDriver)
	fmt.Printf("    Host: %s:%d\n", dbHost, dbPort)
	fmt.Printf("    Database: %s\n", dbName)
	fmt.Printf("    User: %s\n", dbUser)
	fmt.Printf("    SSL: %t\n", dbSSL)

	// Access cache configuration
	cacheDriver := config.GetString("cache.driver", "memory")
	cacheHost := config.GetString("cache.host", "localhost")
	cachePort := config.GetInt("cache.port", 6379)
	cacheTTL := config.GetInt("cache.ttl", 3600)

	fmt.Printf("\n  Cache Configuration:\n")
	fmt.Printf("    Driver: %s\n", cacheDriver)
	fmt.Printf("    Host: %s:%d\n", cacheHost, cachePort)
	fmt.Printf("    TTL: %ds\n", cacheTTL)

	// Access mail configuration
	mailDriver := config.GetString("mail.driver", "smtp")
	mailHost := config.GetString("mail.host", "localhost")
	mailPort := config.GetInt("mail.port", 587)
	mailFromAddress := config.GetString("mail.from.address", "")
	mailFromName := config.GetString("mail.from.name", "")

	fmt.Printf("\n  Mail Configuration:\n")
	fmt.Printf("    Driver: %s\n", mailDriver)
	fmt.Printf("    Host: %s:%d\n", mailHost, mailPort)
	fmt.Printf("    From: %s <%s>\n", mailFromName, mailFromAddress)

	// Access server configuration
	serverHost := config.GetString("server.host", "0.0.0.0")
	serverPort := config.GetInt("server.port", 8080)
	readTimeout := config.GetInt("server.read_timeout", 30)
	writeTimeout := config.GetInt("server.write_timeout", 30)

	fmt.Printf("\n  Server Configuration:\n")
	fmt.Printf("    Listen: %s:%d\n", serverHost, serverPort)
	fmt.Printf("    Read Timeout: %ds\n", readTimeout)
	fmt.Printf("    Write Timeout: %ds\n", writeTimeout)

	// Access logging configuration
	logLevel := config.GetString("logging.level", "info")
	logOutput := config.GetString("logging.output", "stdout")
	logFormat := config.GetString("logging.format", "json")

	fmt.Printf("\n  Logging Configuration:\n")
	fmt.Printf("    Level: %s\n", logLevel)
	fmt.Printf("    Output: %s\n", logOutput)
	fmt.Printf("    Format: %s\n", logFormat)

	fmt.Printf("\n  ✅ Go file configuration loaded successfully!\n")
	fmt.Printf("  📝 The config.go file can include dynamic logic and environment variable processing.\n")
}
