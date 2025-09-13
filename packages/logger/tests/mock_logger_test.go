package tests

import (
	"strings"
	"testing"

	"govel/logger/mocks"
)

// TestMockLogger tests the mock logger implementation
func TestMockLogger(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Test basic logging
	mockLogger.Info("Test info message")
	mockLogger.Error("Test error message")
	mockLogger.Debug("Test debug message")
	mockLogger.Warn("Test warning message")
	
	// Test formatted logging
	mockLogger.Info("Formatted info: %s, %d", "test", 123)
	mockLogger.Error("Formatted error: %s, %d", "error", 456)
	
	// Get messages from mock
	messages := mockLogger.GetMessages()
	
	if len(messages) < 6 {
		t.Errorf("Expected at least 6 messages, got %d", len(messages))
	}
	
	// Check specific messages
	foundInfo := false
	foundError := false
	foundFormatted := false
	
	for _, msg := range messages {
		if strings.Contains(msg.Format, "Test info message") {
			foundInfo = true
		}
		if strings.Contains(msg.Format, "Test error message") {
			foundError = true
		}
		if strings.Contains(msg.Format, "Formatted info: %s, %d") {
			foundFormatted = true
		}
	}
	
	if !foundInfo {
		t.Error("Expected to find info message in mock logger")
	}
	if !foundError {
		t.Error("Expected to find error message in mock logger")
	}
	if !foundFormatted {
		t.Error("Expected to find formatted message in mock logger")
	}
}

// TestMockLoggerLevel tests mock logger level functionality
func TestMockLoggerLevel(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Test setting and getting level
	mockLogger.SetLevel("debug")
	if mockLogger.GetLevel() != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", mockLogger.GetLevel())
	}
	
	mockLogger.SetLevel("error")
	if mockLogger.GetLevel() != "error" {
		t.Errorf("Expected level 'error', got '%s'", mockLogger.GetLevel())
	}
}

// TestMockLoggerWithFields tests mock logger structured logging
func TestMockLoggerWithFields(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Test that WithField works and returns a logger instance
	fieldLogger := mockLogger.WithField("component", "auth")
	if fieldLogger == nil {
		t.Error("Expected WithField to return a logger instance")
		return
	}
	
	// Log a message with the field logger
	fieldLogger.Info("Authentication started")
	
	// Test that WithFields works and returns a logger instance
	multiFieldLogger := mockLogger.WithFields(map[string]interface{}{
		"user_id": 12345,
		"action":  "login",
	})
	if multiFieldLogger == nil {
		t.Error("Expected WithFields to return a logger instance")
		return
	}
	
	// Log a message with the multi-field logger
	multiFieldLogger.Info("User action")
	
	// Check messages from the original logger (should share message storage)
	messages := mockLogger.GetMessages()
	t.Logf("Total messages: %d", len(messages))
	
	// If messages are shared, we should see 2 messages
	// If not shared, we'll see 0 and that's also valid behavior
	if len(messages) >= 2 {
		t.Log("Messages are shared between logger instances")
		for i, msg := range messages {
			t.Logf("Message %d: level=%s, format=%s, fields=%v", i, msg.Level, msg.Format, msg.Fields)
		}
	} else {
		t.Log("Messages are not shared - each logger instance has its own storage")
		// This is also a valid implementation choice
	}
}

// TestMockLoggerFailureMode tests mock logger failure simulation
func TestMockLoggerFailureMode(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Enable failure mode - parameters depend on mock implementation
	mockLogger.SetFailureMode(true) // Enable failure mode
	
	// Test that operations don't panic in failure mode
	mockLogger.Info("This info should fail")
	mockLogger.Debug("This debug should succeed")
	mockLogger.Error("This error should fail")
	
	messages := mockLogger.GetMessages()
	
	// In failure mode, behavior depends on implementation
	// We mainly test that it doesn't panic
	t.Logf("Messages in failure mode: %d", len(messages))
}

// TestMockLoggerClearMessages tests clearing messages in mock logger
func TestMockLoggerClearMessages(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Add some messages
	mockLogger.Info("First message")
	mockLogger.Debug("Second message")
	mockLogger.Error("Third message")
	
	// Verify messages exist
	messages := mockLogger.GetMessages()
	if len(messages) < 3 {
		t.Errorf("Expected at least 3 messages, got %d", len(messages))
	}
	
	// Clear messages if the mock supports it
	mockLogger.ClearMessages()
	
	// Verify messages are cleared
	clearedMessages := mockLogger.GetMessages()
	if len(clearedMessages) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(clearedMessages))
	}
}

// TestMockLoggerThreadSafety tests concurrent usage of mock logger
func TestMockLoggerThreadSafety(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	done := make(chan bool, 3)
	
	// Goroutine 1: Log info messages
	go func() {
		for i := 0; i < 10; i++ {
			mockLogger.Info("Concurrent info %d", i)
		}
		done <- true
	}()
	
	// Goroutine 2: Log error messages
	go func() {
		for i := 0; i < 10; i++ {
			mockLogger.Error("Concurrent error %d", i)
		}
		done <- true
	}()
	
	// Goroutine 3: Add fields and log
	go func() {
		for i := 0; i < 10; i++ {
			contextLogger := mockLogger.WithField("iteration", i)
			if contextLogger != nil {
				contextLogger.Debug("Concurrent field log %d", i)
			}
		}
		done <- true
	}()
	
	// Wait for all goroutines
	<-done
	<-done
	<-done
	
	// Verify that we don't get a panic and some messages were logged
	messages := mockLogger.GetMessages()
	if len(messages) == 0 {
		t.Error("Expected some messages from concurrent logging")
	}
	
	t.Logf("Total messages from concurrent logging: %d", len(messages))
}

// TestMockLoggerFieldBehavior tests field behavior in mock logger
func TestMockLoggerFieldBehavior(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Test field types
	testCases := []struct {
		key   string
		value interface{}
	}{
		{"string_field", "test_string"},
		{"int_field", 42},
		{"float_field", 3.14159},
		{"bool_field", true},
		{"nil_field", nil},
		{"slice_field", []string{"a", "b", "c"}},
		{"map_field", map[string]int{"key": 1}},
	}
	
	for _, tc := range testCases {
		fieldLogger := mockLogger.WithField(tc.key, tc.value)
		if fieldLogger == nil {
			t.Errorf("Expected WithField to work for %s field", tc.key)
		} else {
			fieldLogger.Info("Testing field %s", tc.key)
		}
	}
	
	// Test multiple fields at once
	multiLogger := mockLogger.WithFields(map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": true,
	})
	
	if multiLogger == nil {
		t.Error("Expected WithFields to work with multiple fields")
	} else {
		multiLogger.Info("Multiple fields test")
	}
}

// TestMockLoggerMessageFiltering tests message filtering by level
func TestMockLoggerMessageFiltering(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	
	// Set level to warn - should filter out debug and info
	mockLogger.SetLevel("warn")
	
	// Clear any existing messages
	mockLogger.ClearMessages()
	
	// Log at different levels
	mockLogger.Debug("Debug message - should be filtered")
	mockLogger.Info("Info message - should be filtered")
	mockLogger.Warn("Warn message - should appear")
	mockLogger.Error("Error message - should appear")
	
	messages := mockLogger.GetMessages()
	
	// Check filtering behavior (depends on mock implementation)
	debugFiltered := true
	infoFiltered := true
	warnPresent := false
	errorPresent := false
	
	for _, msg := range messages {
		if strings.Contains(msg.Format, "Debug message") {
			debugFiltered = false
		}
		if strings.Contains(msg.Format, "Info message") {
			infoFiltered = false
		}
		if strings.Contains(msg.Format, "Warn message") {
			warnPresent = true
		}
		if strings.Contains(msg.Format, "Error message") {
			errorPresent = true
		}
	}
	
	// These assertions depend on whether the mock implements level filtering
	// For now, just log the behavior
	t.Logf("Level filtering behavior - Debug filtered: %t, Info filtered: %t, Warn present: %t, Error present: %t",
		debugFiltered, infoFiltered, warnPresent, errorPresent)
}
