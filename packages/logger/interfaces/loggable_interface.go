package interfaces

/**
 * LoggableInterface defines the contract for components that provide
 * logger functionality. This interface follows the Interface Segregation
 * Principle by focusing solely on logger operations.
 *
 * By embedding LoggerInterface, all logger methods are automatically
 * available through this interface, providing transparent delegation.
 */
type LoggableInterface interface {
	/**
	 * Embed LoggerInterface to provide transparent access to all logger methods
	 */
	LoggerInterface

	/**
	 * GetLogger returns the logger instance.
	 *
	 * @return LoggerInterface The logger instance
	 */
	GetLogger() LoggerInterface

	/**
	 * SetLogger sets the logger instance.
	 *
	 * @param logger interface{} The logger instance to set (using interface{} to avoid circular import)
	 */
	SetLogger(logger interface{})

	/**
	 * HasLogger returns whether a logger instance is set.
	 *
	 * @return bool true if a logger is set
	 */
	HasLogger() bool

	/**
	 * GetLoggerInfo returns information about the logger.
	 *
	 * @return map[string]interface{} Logger information
	 */
	GetLoggerInfo() map[string]interface{}
}
