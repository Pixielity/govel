package tests

import (
	"fmt"
	"testing"

	"govel/packages/logger"
	"govel/packages/logger/mocks"
)

// BenchmarkLoggerInfo benchmarks info level logging
func BenchmarkLoggerInfo(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.Info("Benchmark info message")
	}
}

// BenchmarkLoggerInfoFormatted benchmarks formatted info logging
func BenchmarkLoggerInfoFormatted(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.Info("Benchmark formatted message: %d", i)
	}
}

// BenchmarkLoggerDebug benchmarks debug level logging
func BenchmarkLoggerDebug(b *testing.B) {
	l := logger.New()
	l.SetLevel("debug") // Ensure debug messages are processed
	
	for i := 0; i < b.N; i++ {
		l.Debug("Benchmark debug message")
	}
}

// BenchmarkLoggerError benchmarks error level logging
func BenchmarkLoggerError(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.Error("Benchmark error message")
	}
}

// BenchmarkLoggerWarn benchmarks warning level logging
func BenchmarkLoggerWarn(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.Warn("Benchmark warning message")
	}
}

// BenchmarkLoggerWithField benchmarks single field logging
func BenchmarkLoggerWithField(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.WithField("request_id", i).Info("Benchmark single field message")
	}
}

// BenchmarkLoggerWithFields benchmarks multiple field logging
func BenchmarkLoggerWithFields(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.WithFields(map[string]interface{}{
			"request_id": i,
			"user_id":    fmt.Sprintf("user_%d", i),
			"action":     "benchmark",
		}).Info("Benchmark multiple fields message")
	}
}

// BenchmarkLoggerFieldReuse benchmarks reusing a field logger
func BenchmarkLoggerFieldReuse(b *testing.B) {
	l := logger.New()
	fieldLogger := l.WithFields(map[string]interface{}{
		"service": "api",
		"version": "1.0.0",
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fieldLogger.Info("Benchmark field reuse message %d", i)
	}
}

// BenchmarkLoggerLevelFiltering benchmarks level filtering performance
func BenchmarkLoggerLevelFiltering(b *testing.B) {
	l := logger.New()
	l.SetLevel("error") // Only error and above should be logged
	
	for i := 0; i < b.N; i++ {
		// These should be filtered out
		l.Debug("Debug message %d", i)
		l.Info("Info message %d", i)
		l.Warn("Warning message %d", i)
		
		// This should be logged
		l.Error("Error message %d", i)
	}
}

// BenchmarkMockLogger benchmarks mock logger operations
func BenchmarkMockLogger(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	
	for i := 0; i < b.N; i++ {
		mockLogger.Info("Mock benchmark message: %d", i)
	}
}

// BenchmarkMockLoggerWithFields benchmarks mock logger with fields
func BenchmarkMockLoggerWithFields(b *testing.B) {
	mockLogger := mocks.NewMockLogger()
	
	for i := 0; i < b.N; i++ {
		mockLogger.WithField("iteration", i).Info("Mock field benchmark: %d", i)
	}
}

// BenchmarkLoggerFieldChaining benchmarks field chaining performance
func BenchmarkLoggerFieldChaining(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.WithField("step1", i).
			WithField("step2", i*2).
			WithField("step3", i*3).
			Info("Chained fields benchmark")
	}
}

// BenchmarkLoggerComplexFormatting benchmarks complex formatted messages
func BenchmarkLoggerComplexFormatting(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		l.Info("Complex formatting: user=%s, id=%d, active=%t, score=%.2f, items=%v",
			fmt.Sprintf("user_%d", i), i, i%2 == 0, float64(i)*3.14159, []int{i, i + 1, i + 2})
	}
}

// BenchmarkLoggerMixedOperations benchmarks mixed logging operations
func BenchmarkLoggerMixedOperations(b *testing.B) {
	l := logger.New()
	
	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			l.Debug("Mixed operation debug %d", i)
		case 1:
			l.WithField("operation", "info").Info("Mixed operation info %d", i)
		case 2:
			l.Warn("Mixed operation warning %d", i)
		case 3:
			l.WithFields(map[string]interface{}{
				"operation": "error",
				"iteration": i,
			}).Error("Mixed operation error %d", i)
		}
	}
}

// BenchmarkLoggerSetLevel benchmarks level setting performance
func BenchmarkLoggerSetLevel(b *testing.B) {
	l := logger.New()
	levels := []string{"debug", "info", "warn", "error", "fatal"}
	
	for i := 0; i < b.N; i++ {
		l.SetLevel(levels[i%len(levels)])
	}
}

// BenchmarkLoggerGetLevel benchmarks level getting performance
func BenchmarkLoggerGetLevel(b *testing.B) {
	l := logger.New()
	l.SetLevel("info")
	
	for i := 0; i < b.N; i++ {
		_ = l.GetLevel()
	}
}
