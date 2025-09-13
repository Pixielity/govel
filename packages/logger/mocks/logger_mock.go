package mocks

import (
	"time"

	loggerInterfaces "govel/packages/types/src/interfaces/logger"
)

/**
 * MockLogger provides a mock implementation of LoggerInterface for testing.
 * This mock allows tests to verify logging behavior and capture log messages for assertion.
 */
type MockLogger struct {
	// Log Messages Storage
	Messages []MockLogMessage

	// Logger Configuration
	Level  string
	Fields map[string]interface{}

	// Mock Control Flags
	ShouldFailLog bool
}

/**
 * MockLogMessage represents a logged message for testing verification
 */
type MockLogMessage struct {
	Level  string
	Format string
	Args   []interface{}
	Time   time.Time
	Fields map[string]interface{} // Contextual fields
}

/**
 * NewMockLogger creates a new mock logger with default values
 */
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Messages: make([]MockLogMessage, 0),
		Level:    "debug",
		Fields:   make(map[string]interface{}),
	}
}

/**
 * NewMockLoggerWithLevel creates a new mock logger with a specific level
 */
func NewMockLoggerWithLevel(level string) *MockLogger {
	logger := NewMockLogger()
	logger.Level = level
	return logger
}

// LoggerInterface Implementation

func (m *MockLogger) Debug(format string, args ...interface{}) {
	if m.ShouldFailLog {
		return
	}

	m.addMessage("debug", format, args)
}

func (m *MockLogger) Info(format string, args ...interface{}) {
	if m.ShouldFailLog {
		return
	}

	m.addMessage("info", format, args)
}

func (m *MockLogger) Warn(format string, args ...interface{}) {
	if m.ShouldFailLog {
		return
	}

	m.addMessage("warn", format, args)
}

func (m *MockLogger) Error(format string, args ...interface{}) {
	if m.ShouldFailLog {
		return
	}

	m.addMessage("error", format, args)
}

func (m *MockLogger) Fatal(format string, args ...interface{}) {
	if m.ShouldFailLog {
		return
	}

	m.addMessage("fatal", format, args)
	// Note: Real Fatal would call os.Exit(1), but we don't want that in tests
}

func (m *MockLogger) WithField(key string, value interface{}) loggerInterfaces.LoggerInterface {
	// Create a new logger with the added field
	newLogger := &MockLogger{
		Messages:      m.Messages, // Share the same message slice
		Level:         m.Level,
		ShouldFailLog: m.ShouldFailLog,
		Fields:        make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range m.Fields {
		newLogger.Fields[k] = v
	}

	// Add the new field
	newLogger.Fields[key] = value

	return newLogger
}

func (m *MockLogger) WithFields(fields map[string]interface{}) loggerInterfaces.LoggerInterface {
	// Create a new logger with the added fields
	newLogger := &MockLogger{
		Messages:      m.Messages, // Share the same message slice
		Level:         m.Level,
		ShouldFailLog: m.ShouldFailLog,
		Fields:        make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range m.Fields {
		newLogger.Fields[k] = v
	}

	// Add the new fields
	for k, v := range fields {
		newLogger.Fields[k] = v
	}

	return newLogger
}

func (m *MockLogger) SetLevel(level string) {
	m.Level = level
}

func (m *MockLogger) GetLevel() string {
	return m.Level
}

// Internal helper method
func (m *MockLogger) addMessage(level, format string, args []interface{}) {
	// Copy fields to avoid reference issues
	fieldsCopy := make(map[string]interface{})
	for k, v := range m.Fields {
		fieldsCopy[k] = v
	}

	message := MockLogMessage{
		Level:  level,
		Format: format,
		Args:   args,
		Time:   time.Now(),
		Fields: fieldsCopy,
	}

	m.Messages = append(m.Messages, message)
}

// Mock-specific helper methods

/**
 * GetMessages returns all logged messages
 */
func (m *MockLogger) GetMessages() []MockLogMessage {
	return m.Messages
}

/**
 * GetMessagesOfLevel returns all log messages of a specific level
 */
func (m *MockLogger) GetMessagesOfLevel(level string) []MockLogMessage {
	var messages []MockLogMessage
	for _, msg := range m.Messages {
		if msg.Level == level {
			messages = append(messages, msg)
		}
	}
	return messages
}

/**
 * GetMessageCount returns the total number of logged messages
 */
func (m *MockLogger) GetMessageCount() int {
	return len(m.Messages)
}

/**
 * GetMessageCountOfLevel returns the number of messages of a specific level
 */
func (m *MockLogger) GetMessageCountOfLevel(level string) int {
	count := 0
	for _, msg := range m.Messages {
		if msg.Level == level {
			count++
		}
	}
	return count
}

/**
 * ClearMessages clears all logged messages
 */
func (m *MockLogger) ClearMessages() {
	m.Messages = make([]MockLogMessage, 0)
}

/**
 * SetFailureMode sets whether logging should fail
 */
func (m *MockLogger) SetFailureMode(shouldFail bool) {
	m.ShouldFailLog = shouldFail
}

/**
 * GetFields returns the current contextual fields
 */
func (m *MockLogger) GetFields() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m.Fields {
		result[k] = v
	}
	return result
}

/**
 * HasMessage checks if any message contains the specified text
 */
func (m *MockLogger) HasMessage(text string) bool {
	for _, msg := range m.Messages {
		if msg.Format == text {
			return true
		}
		// Could also check formatted message, but that's more complex
	}
	return false
}

/**
 * HasMessageWithLevel checks if any message of a specific level contains the specified text
 */
func (m *MockLogger) HasMessageWithLevel(level, text string) bool {
	for _, msg := range m.Messages {
		if msg.Level == level && msg.Format == text {
			return true
		}
	}
	return false
}

/**
 * GetLastMessage returns the most recently logged message
 */
func (m *MockLogger) GetLastMessage() *MockLogMessage {
	if len(m.Messages) == 0 {
		return nil
	}
	return &m.Messages[len(m.Messages)-1]
}

/**
 * GetLastMessageOfLevel returns the most recent message of a specific level
 */
func (m *MockLogger) GetLastMessageOfLevel(level string) *MockLogMessage {
	for i := len(m.Messages) - 1; i >= 0; i-- {
		if m.Messages[i].Level == level {
			return &m.Messages[i]
		}
	}
	return nil
}

/**
 * GetMessagesWithField returns all messages that have a specific field
 */
func (m *MockLogger) GetMessagesWithField(key string) []MockLogMessage {
	var messages []MockLogMessage
	for _, msg := range m.Messages {
		if _, exists := msg.Fields[key]; exists {
			messages = append(messages, msg)
		}
	}
	return messages
}

/**
 * GetMessagesWithFieldValue returns all messages that have a field with a specific value
 */
func (m *MockLogger) GetMessagesWithFieldValue(key string, value interface{}) []MockLogMessage {
	var messages []MockLogMessage
	for _, msg := range m.Messages {
		if fieldValue, exists := msg.Fields[key]; exists && fieldValue == value {
			messages = append(messages, msg)
		}
	}
	return messages
}

// Compile-time interface compliance check
var _ loggerInterfaces.LoggerInterface = (*MockLogger)(nil)

/**
 * MockLoggable provides a mock implementation of LoggableInterface for testing.
 */
type MockLoggable struct {
	*MockLogger

	LoggerInstance loggerInterfaces.LoggerInterface
	HasLoggerValue bool
}

/**
 * NewMockLoggable creates a new mock loggable with default values
 */
func NewMockLoggable() *MockLoggable {
	mockLogger := NewMockLogger()
	return &MockLoggable{
		MockLogger:     mockLogger,
		LoggerInstance: mockLogger,
		HasLoggerValue: true,
	}
}

// LoggableInterface Implementation

func (m *MockLoggable) GetLogger() loggerInterfaces.LoggerInterface {
	return m.LoggerInstance
}

func (m *MockLoggable) SetLogger(logger interface{}) {
	if log, ok := logger.(loggerInterfaces.LoggerInterface); ok {
		m.LoggerInstance = log
		m.HasLoggerValue = true
	} else if log, ok := logger.(*MockLogger); ok {
		m.LoggerInstance = log
		m.HasLoggerValue = true
	}
}

func (m *MockLoggable) HasLogger() bool {
	return m.HasLoggerValue
}

func (m *MockLoggable) GetLoggerInfo() map[string]interface{} {
	info := map[string]interface{}{
		"has_logger":  m.HasLoggerValue,
		"logger_type": "mock",
	}

	if m.LoggerInstance != nil {
		info["current_level"] = m.LoggerInstance.GetLevel()

		if mockLogger, ok := m.LoggerInstance.(*MockLogger); ok {
			info["messages_count"] = mockLogger.GetMessageCount()
			info["debug_count"] = mockLogger.GetMessageCountOfLevel("debug")
			info["info_count"] = mockLogger.GetMessageCountOfLevel("info")
			info["warn_count"] = mockLogger.GetMessageCountOfLevel("warn")
			info["error_count"] = mockLogger.GetMessageCountOfLevel("error")
			info["fatal_count"] = mockLogger.GetMessageCountOfLevel("fatal")
			info["fields_count"] = len(mockLogger.GetFields())
		}
	}

	return info
}

// Mock-specific helper methods for Loggable

/**
 * SetHasLogger controls whether the loggable reports having a logger
 */
func (m *MockLoggable) SetHasLogger(hasLogger bool) {
	m.HasLoggerValue = hasLogger
}

/**
 * GetMockLogger returns the underlying MockLogger if available
 */
func (m *MockLoggable) GetMockLogger() *MockLogger {
	if mockLogger, ok := m.LoggerInstance.(*MockLogger); ok {
		return mockLogger
	}
	return nil
}

// Compile-time interface compliance check
var _ loggerInterfaces.LoggableInterface = (*MockLoggable)(nil)
