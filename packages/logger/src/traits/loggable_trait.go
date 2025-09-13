package traits

import (
	"sync"

	"govel/logger"
	"govel/types/interfaces/logger"
)

/**
 * Loggable provides logger management functionality in a thread-safe manner.
 * This trait follows the embedding pattern where it embeds the logger instance directly,
 * providing both trait-specific management methods and transparent access to all
 * logger interface methods through embedding.
 */
type Loggable struct {
	/**
	 * mutex provides thread-safe access to logger properties
	 */
	mutex sync.RWMutex

	/**
	 * Logger instance embedded directly to provide transparent delegation
	 * All LoggerInterface methods are automatically available
	 */
	*logger.Logger
}

/**
 * NewLoggable creates a new Loggable instance with a logger.
 *
 * @param loggerInstance *logger.Logger The logger instance to use
 * @return *Loggable The newly created trait instance
 */
func NewLoggable(loggerInstance *logger.Logger) *Loggable {
	if loggerInstance == nil {
		loggerInstance = logger.New()
	}

	return &Loggable{
		Logger: loggerInstance,
	}
}

/**
 * NewLoggableDefault creates a new Loggable instance with a default logger.
 * This is a convenience constructor that creates a trait with a default logger instance.
 *
 * @return *Loggable The newly created trait instance with default logger
 */
func NewLoggableDefault() *Loggable {
	return &Loggable{
		Logger: logger.New(),
	}
}

/**
 * GetLogger returns the logger instance.
 *
 * @return logger_interfaces.LoggerInterface The logger instance
 */
func (t *Loggable) GetLogger() interfaces.LoggerInterface {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Logger
}

/**
 * SetLogger sets the logger instance.
 *
 * @param loggerInstance *logger.Logger The logger instance to set
 */
func (t *Loggable) SetLogger(loggerInstance interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if loggerPtr, ok := loggerInstance.(*logger.Logger); ok {
		if loggerPtr == nil {
			loggerPtr = logger.New()
		}
		t.Logger = loggerPtr
	} else {
		// Fallback to default logger if invalid type
		t.Logger = logger.New()
	}
}

/**
 * HasLogger returns whether a logger instance is set.
 *
 * @return bool true if a logger is set
 */
func (t *Loggable) HasLogger() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.Logger != nil
}

/**
 * GetLoggerInfo returns information about the logger.
 *
 * @return map[string]interface{} Logger information
 */
func (t *Loggable) GetLoggerInfo() map[string]interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	info := map[string]interface{}{
		"has_logger": t.Logger != nil,
	}

	if t.Logger != nil {
		// Add logger-specific information if available
		// This would depend on the actual logger implementation
		info["logger_type"] = "default"
	}

	return info
}

// Compile-time interface compliance check
var _ interfaces.LoggableInterface = (*Loggable)(nil)
