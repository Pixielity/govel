package tests

import (
	"testing"
	"time"

	"govel/application/builders"
	"govel/container"
)

// BenchmarkAppBuilderSimple benchmarks simple builder usage
func BenchmarkAppBuilderSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := builders.NewApp().
			WithName("Benchmark App").
			WithVersion("1.0.0").
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderComplex benchmarks complex builder usage
func BenchmarkAppBuilderComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := builders.NewApp().
			WithName("Complex Benchmark App").
			WithVersion("2.0.0").
			WithEnvironment("production").
			WithLocale("en").
			WithFallbackLocale("en").
			WithTimezone("UTC").
			WithBasePath("/opt/app").
			WithShutdownTimeout(60 * time.Second).
			WithDebug(false).
			InConsole().
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderConvenience benchmarks convenience methods
func BenchmarkAppBuilderConvenience(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := builders.NewApp().
			ForProduction().
			WithName("Convenience App").
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderWithContainer benchmarks builder with container
func BenchmarkAppBuilderWithContainer(b *testing.B) {
	customContainer := container.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := builders.NewApp().
			WithContainer(customContainer).
			WithName("Container App").
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderForEnvironments benchmarks different environment configurations
func BenchmarkAppBuilderForEnvironments(b *testing.B) {
	environments := []func(*builders.AppBuilder) *builders.AppBuilder{
		func(builder *builders.AppBuilder) *builders.AppBuilder { return builder.ForDevelopment() },
		func(builder *builders.AppBuilder) *builders.AppBuilder { return builder.ForTesting() },
		func(builder *builders.AppBuilder) *builders.AppBuilder { return builder.ForProduction() },
	}

	for i := 0; i < b.N; i++ {
		envMethod := environments[i%len(environments)]
		app := envMethod(builders.NewApp()).
			WithName("Environment App").
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderChaining benchmarks long method chains
func BenchmarkAppBuilderChaining(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := builders.NewApp().
			WithName("Chaining App").
			WithVersion("1.0.0").
			WithEnvironment("development").
			WithDebug(true).
			WithLocale("en").
			WithFallbackLocale("en").
			WithTimezone("UTC").
			WithBasePath(".").
			WithShutdownTimeout(30 * time.Second).
			InConsole().
			InTesting().
			Build()
		_ = app // Use the app to prevent optimization
	}
}

// BenchmarkAppBuilderCreation benchmarks just the builder creation
func BenchmarkAppBuilderCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := builders.NewApp()
		_ = builder // Use the builder to prevent optimization
	}
}

// BenchmarkAppBuilderBuild benchmarks just the build process
func BenchmarkAppBuilderBuild(b *testing.B) {
	// Setup builder outside of timing
	builder := builders.NewApp().
		WithName("Build Test").
		WithVersion("1.0.0").
		ForProduction()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := builder.Build()
		_ = app // Use the app to prevent optimization
	}
}
