package tests

import (
	"testing"

	"govel/logger"
)

// TestLoggerWithFields tests structured logging with fields
func TestLoggerWithFields(t *testing.T) {
	l := logger.New()
	
	// Test logging with context fields
	contextLogger1 := l.WithField("user_id", 12345)
	if contextLogger1 == nil {
		t.Error("Expected WithField to return a logger instance")
	} else {
		contextLogger1.Info("User logged in")
	}
	
	contextLogger2 := l.WithField("request_id", "req_123")
	if contextLogger2 == nil {
		t.Error("Expected WithField to return a logger instance")
	} else {
		contextLogger2.Debug("Processing request")
	}
	
	// Test multiple fields
	multiFieldLogger := l.WithFields(map[string]interface{}{
		"user_id":    12345,
		"session_id": "sess_456",
		"action":     "login",
	})
	if multiFieldLogger == nil {
		t.Error("Expected WithFields to return a logger instance")
	} else {
		multiFieldLogger.Info("User action performed")
	}
	
	// Test chaining with fields
	contextLogger := l.WithField("component", "auth")
	if contextLogger != nil {
		contextLogger.Info("Authentication started")
		
		// Chain another field
		methodLogger := contextLogger.WithField("method", "oauth")
		if methodLogger != nil {
			methodLogger.Info("Using OAuth method")
		}
	}
}

// TestLoggerFieldTypes tests different field value types
func TestLoggerFieldTypes(t *testing.T) {
	l := logger.New()
	
	// Test string field
	l.WithField("string_field", "test_value").Info("String field test")
	
	// Test integer field
	l.WithField("int_field", 42).Info("Integer field test")
	
	// Test boolean field
	l.WithField("bool_field", true).Info("Boolean field test")
	
	// Test float field
	l.WithField("float_field", 3.14159).Info("Float field test")
	
	// Test nil field
	l.WithField("nil_field", nil).Info("Nil field test")
	
	// Test struct field
	testStruct := struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	}
	l.WithField("struct_field", testStruct).Info("Struct field test")
	
	// Test slice field
	l.WithField("slice_field", []string{"a", "b", "c"}).Info("Slice field test")
	
	// Test map field
	l.WithField("map_field", map[string]int{"key1": 1, "key2": 2}).Info("Map field test")
}

// TestLoggerFieldChaining tests field chaining behavior
func TestLoggerFieldChaining(t *testing.T) {
	l := logger.New()
	
	// Create a base logger with fields
	baseLogger := l.WithFields(map[string]interface{}{
		"service": "api",
		"version": "1.0.0",
	})
	
	if baseLogger == nil {
		t.Error("Expected base logger to be created")
		return
	}
	
	// Add more fields to the base logger
	requestLogger := baseLogger.WithFields(map[string]interface{}{
		"request_id": "req_12345",
		"user_id":    67890,
	})
	
	if requestLogger != nil {
		requestLogger.Info("Processing API request")
		
		// Add even more context
		errorLogger := requestLogger.WithField("error_code", "AUTH_001")
		if errorLogger != nil {
			errorLogger.Error("Authentication failed")
		}
	}
}

// TestLoggerFieldOverride tests field override behavior
func TestLoggerFieldOverride(t *testing.T) {
	l := logger.New()
	
	// Create logger with initial field
	logger1 := l.WithField("environment", "development")
	if logger1 == nil {
		t.Error("Expected logger1 to be created")
		return
	}
	
	// Override the same field
	logger2 := logger1.WithField("environment", "production")
	if logger2 == nil {
		t.Error("Expected logger2 to be created")
		return
	}
	
	// Both loggers should work independently
	logger1.Info("Development message")
	logger2.Info("Production message")
	
	// Test overriding with WithFields
	logger3 := logger1.WithFields(map[string]interface{}{
		"environment": "staging",
		"region":      "us-east-1",
	})
	
	if logger3 != nil {
		logger3.Info("Staging message with region")
	}
}

// TestLoggerEmptyFields tests behavior with empty fields
func TestLoggerEmptyFields(t *testing.T) {
	l := logger.New()
	
	// Test WithFields with empty map
	emptyFieldsLogger := l.WithFields(map[string]interface{}{})
	if emptyFieldsLogger == nil {
		t.Error("Expected logger with empty fields to be created")
	} else {
		emptyFieldsLogger.Info("Message with empty fields")
	}
	
	// Test WithField with empty key
	emptyKeyLogger := l.WithField("", "empty_key_value")
	if emptyKeyLogger == nil {
		t.Error("Expected logger with empty key to be created")
	} else {
		emptyKeyLogger.Info("Message with empty key field")
	}
	
	// Test WithField with empty value
	emptyValueLogger := l.WithField("empty_value", "")
	if emptyValueLogger == nil {
		t.Error("Expected logger with empty value to be created")
	} else {
		emptyValueLogger.Info("Message with empty value field")
	}
}
