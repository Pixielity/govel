package tests

import (
	"fmt"
	"testing"

	"govel/config"
	"govel/config/mocks"
)

// BenchmarkConfigSet benchmarks configuration setting operations
func BenchmarkConfigSet(b *testing.B) {
	cfg := config.New()
	
	for i := 0; i < b.N; i++ {
		cfg.Set("benchmark.key", i)
	}
}

// BenchmarkConfigGetString benchmarks string retrieval operations
func BenchmarkConfigGetString(b *testing.B) {
	cfg := config.New()
	cfg.Set("benchmark.key", "benchmark_value")
	
	for i := 0; i < b.N; i++ {
		_ = cfg.GetString("benchmark.key", "default")
	}
}

// BenchmarkConfigGetInt benchmarks integer retrieval operations
func BenchmarkConfigGetInt(b *testing.B) {
	cfg := config.New()
	cfg.Set("benchmark.number", 12345)
	
	for i := 0; i < b.N; i++ {
		_ = cfg.GetInt("benchmark.number", 0)
	}
}

// BenchmarkConfigGetBool benchmarks boolean retrieval operations
func BenchmarkConfigGetBool(b *testing.B) {
	cfg := config.New()
	cfg.Set("benchmark.flag", true)
	
	for i := 0; i < b.N; i++ {
		_ = cfg.GetBool("benchmark.flag", false)
	}
}

// BenchmarkConfigHasKey benchmarks key existence checking
func BenchmarkConfigHasKey(b *testing.B) {
	cfg := config.New()
	cfg.Set("benchmark.key", "value")
	
	for i := 0; i < b.N; i++ {
		_ = cfg.HasKey("benchmark.key")
	}
}

// BenchmarkConfigAllConfig benchmarks getting all configuration
func BenchmarkConfigAllConfig(b *testing.B) {
	cfg := config.New()
	
	// Setup some config values
	for i := 0; i < 100; i++ {
		cfg.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cfg.AllConfig()
	}
}

// BenchmarkMockConfigOperations benchmarks mock config operations
func BenchmarkMockConfigOperations(b *testing.B) {
	mockCfg := mocks.NewMockConfig()
	
	for i := 0; i < b.N; i++ {
		mockCfg.Set("benchmark.key", i)
		_ = mockCfg.GetString("benchmark.key", "default")
	}
}

// BenchmarkConfigMixedOperations benchmarks mixed config operations
func BenchmarkConfigMixedOperations(b *testing.B) {
	cfg := config.New()
	
	for i := 0; i < b.N; i++ {
		cfg.Set("mixed.string", "test")
		cfg.Set("mixed.int", i)
		cfg.Set("mixed.bool", i%2 == 0)
		
		_ = cfg.GetString("mixed.string", "")
		_ = cfg.GetInt("mixed.int", 0)
		_ = cfg.GetBool("mixed.bool", false)
		_ = cfg.HasKey("mixed.string")
	}
}
