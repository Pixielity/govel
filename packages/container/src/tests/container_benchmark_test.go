package tests

import (
	"testing"

	"govel/container"
	"govel/container/mocks"
)

// BenchmarkContainerBind benchmarks container binding operations
func BenchmarkContainerBind(b *testing.B) {
	c := container.New()
	
	for i := 0; i < b.N; i++ {
		c.Bind("benchmark.service", "benchmark_value")
	}
}

// BenchmarkContainerMake benchmarks container resolution operations
func BenchmarkContainerMake(b *testing.B) {
	c := container.New()
	c.Bind("benchmark.service", "benchmark_value")
	
	for i := 0; i < b.N; i++ {
		c.Make("benchmark.service")
	}
}

// BenchmarkContainerSingleton benchmarks singleton resolution
func BenchmarkContainerSingleton(b *testing.B) {
	c := container.New()
	c.Singleton("benchmark.singleton", func() interface{} {
		return &struct{ Value int }{Value: 42}
	})
	
	for i := 0; i < b.N; i++ {
		c.Make("benchmark.singleton")
	}
}

// BenchmarkMockContainerOperations benchmarks mock container operations
func BenchmarkMockContainerOperations(b *testing.B) {
	mockContainer := mocks.NewMockContainer()
	
	for i := 0; i < b.N; i++ {
		mockContainer.Bind("benchmark.service", i)
		mockContainer.Make("benchmark.service")
	}
}
