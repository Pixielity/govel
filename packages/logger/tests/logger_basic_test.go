package tests

import (
	"testing"

	"govel/packages/logger"
)

// TestLoggerCreation tests basic logger creation
func TestLoggerCreation(t *testing.T) {
	l := logger.New()
	
	if l == nil {
		t.Fatal("Expected logger to be created, got nil")
	}
	
	// Test default log level
	if l.GetLevel() == "" {
		t.Error("Expected logger to have a default log level")
	}
}

// TestLoggerLevels tests different logging levels
func TestLoggerLevels(t *testing.T) {
	l := logger.New()
	
	// Test setting log level
	l.SetLevel("debug")
	if l.GetLevel() != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", l.GetLevel())
	}
	
	l.SetLevel("info")
	if l.GetLevel() != "info" {
		t.Errorf("Expected log level 'info', got '%s'", l.GetLevel())
	}
	
	l.SetLevel("warn")
	if l.GetLevel() != "warn" {
		t.Errorf("Expected log level 'warn', got '%s'", l.GetLevel())
	}
	
	l.SetLevel("error")
	if l.GetLevel() != "error" {
		t.Errorf("Expected log level 'error', got '%s'", l.GetLevel())
	}
	
	l.SetLevel("fatal")
	if l.GetLevel() != "fatal" {
		t.Errorf("Expected log level 'fatal', got '%s'", l.GetLevel())
	}
}

// TestLoggerBasicLogging tests basic logging methods
func TestLoggerBasicLogging(t *testing.T) {
	l := logger.New()
	
	// These tests mainly verify that methods don't panic
	// Actual output testing would depend on logger implementation
	
	l.Debug("Debug message")
	l.Info("Info message")
	l.Warn("Warning message")
	l.Error("Error message")
	
	// Test formatted logging
	l.Debug("Debug with args: %s, %d", "test", 123)
	l.Info("Info with args: %s, %d", "test", 456)
	l.Warn("Warning with args: %s, %d", "test", 789)
	l.Error("Error with args: %s, %d", "test", 101)
	
	// Note: We don't test Fatal() as it would terminate the test process
}
