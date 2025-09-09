package tests

import (
	"testing"

	"govel/packages/application"
	"govel/packages/application/mocks"
)

// BenchmarkApplicationCreation benchmarks application creation
func BenchmarkApplicationCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := application.New()
		_ = app.GetName() // Use the app to prevent optimization
	}
}

// BenchmarkApplicationIdentity benchmarks identity operations
func BenchmarkApplicationIdentity(b *testing.B) {
	app := application.New()
	
	for i := 0; i < b.N; i++ {
		app.SetName("Benchmark App")
		_ = app.GetName()
		app.SetVersion("1.0.0")
		_ = app.GetVersion()
	}
}

// BenchmarkApplicationConfiguration benchmarks configuration operations
func BenchmarkApplicationConfiguration(b *testing.B) {
	app := application.New()
	
	for i := 0; i < b.N; i++ {
		app.Set("benchmark.key", i)
		_ = app.GetString("benchmark.key", "default")
		_ = app.GetInt("benchmark.key", 0)
		_ = app.HasKey("benchmark.key")
	}
}

// BenchmarkApplicationContainer benchmarks container operations
func BenchmarkApplicationContainer(b *testing.B) {
	app := application.New()
	
	// Setup
	app.Bind("benchmark.service", "benchmark_value")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = app.IsBound("benchmark.service")
		app.Make("benchmark.service")
	}
}

// BenchmarkApplicationInfo benchmarks GetApplicationInfo
func BenchmarkApplicationInfo(b *testing.B) {
	app := application.New()
	app.SetName("Benchmark App")
	app.SetVersion("1.0.0")
	
	for i := 0; i < b.N; i++ {
		_ = app.GetApplicationInfo()
	}
}

// BenchmarkMockApplicationCreation benchmarks mock application creation
func BenchmarkMockApplicationCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mockApp := mocks.NewMockApplication()
		_ = mockApp.GetName() // Use the mock to prevent optimization
	}
}

// BenchmarkMockApplicationOperations benchmarks mock application operations
func BenchmarkMockApplicationOperations(b *testing.B) {
	mockApp := mocks.NewMockApplication()
	
	for i := 0; i < b.N; i++ {
		mockApp.Set("benchmark.key", i)
		_ = mockApp.GetString("benchmark.key", "default")
		mockApp.Bind("benchmark.service", i)
		mockApp.Make("benchmark.service")
	}
}
