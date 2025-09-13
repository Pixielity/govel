package constants

import "time"

/**
 * Shutdown constants define the standard shutdown configuration values.
 * These constants ensure consistency across different parts of the application.
 */

/**
 * MinShutdownTimeout is the minimum allowed shutdown timeout.
 */
const MinShutdownTimeout = 1 * time.Second

/**
 * MaxShutdownTimeout is the maximum allowed shutdown timeout.
 */
const MaxShutdownTimeout = 5 * time.Minute

/**
 * DefaultGracePeriod is the default grace period before force shutdown.
 */
const DefaultGracePeriod = 10 * time.Second

/**
 * ShutdownSignalBufferSize is the buffer size for the shutdown signal channel.
 */
const ShutdownSignalBufferSize = 1