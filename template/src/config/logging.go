package config

// Logging returns the logging configuration map.
// This matches Laravel's logging.php configuration structure exactly.
// This configuration handles log channels, drivers, and formatting
// for application logging and debugging.
func Logging() map[string]any {
	return map[string]any{

		// Default Log Channel
		//
		//
		// This option defines the default log channel that is utilized to write
		// messages to your logs. The value provided here should match one of
		// the channels present in the list of "channels" configured below.
		//
		"default": Env("LOG_CHANNEL", "stack"),

		// Deprecations Log Channel
		//
		//
		// This option controls the log channel that should be used to log warnings
		// regarding deprecated language and library features. This allows you to get
		// your application ready for upcoming major versions of dependencies.
		//
		"deprecations": map[string]any{
			// Deprecation Log Channel
			//
			// The log channel to use for deprecation warnings. Set to "null"
			// to disable deprecation logging completely.
			"channel": Env("LOG_DEPRECATIONS_CHANNEL", "null"),

			// Include Stack Traces
			//
			// Whether to include stack traces in deprecation log entries.
			// This helps identify the source of deprecated code usage.
			"trace": Env("LOG_DEPRECATIONS_TRACE", false),
		},

		// Log Channels
		//
		//
		// Here you may configure the log channels for your application. The
		// framework utilizes various logging libraries, which includes a variety
		// of powerful log handlers and formatters that you're free to use.
		//
		// Available drivers: "single", "daily", "slack", "syslog",
		// "errorlog", "monolog", "custom", "stack"
		//
		"channels": map[string]any{

			// Stack Channel (Multi-Channel Logging)
			//
			// Allows logging to multiple channels simultaneously. Useful for
			// sending logs to both files and external services.
			"stack": map[string]any{
				// Log driver type
				"driver": "stack",

				// List of channels to write to simultaneously
				"channels": func() []string {
					channels := Env("LOG_STACK", "single").(string)
					if channels == "" {
						return []string{"single"}
					}
					return []string{channels} // In Go, we'll simplify this for now
				}(),

				// Whether to ignore exceptions from individual channels
				"ignore_exceptions": false,
			},

			// Single File Channel
			//
			// Writes all logs to a single file. Simple and straightforward
			// for small applications or development.
			"single": map[string]any{
				// Log driver type
				"driver": "single",

				// Path to the log file
				"path": StoragePath("logs/govel.log"),

				// Minimum log level to record
				"level": Env("LOG_LEVEL", "debug"),

				// Whether to replace placeholders in log messages
				"replace_placeholders": true,
			},

			// Daily Rotating File Channel
			//
			// Creates a new log file each day and automatically cleans up
			// old files. Good for production to manage log file sizes.
			"daily": map[string]any{
				// Log driver type
				"driver": "daily",

				// Base path for daily log files (date will be appended)
				"path": StoragePath("logs/govel.log"),

				// Minimum log level to record
				"level": Env("LOG_LEVEL", "debug"),

				// Number of days to retain log files
				"days": Env("LOG_DAILY_DAYS", 14),

				// Whether to replace placeholders in log messages
				"replace_placeholders": true,
			},

			// Slack Notification Channel
			//
			// Sends log messages to a Slack channel via webhook.
			// Typically used for critical errors and alerts.
			"slack": map[string]any{
				// Log driver type
				"driver": "slack",

				// Slack webhook URL for sending messages
				"url": Env("LOG_SLACK_WEBHOOK_URL", ""),

				// Username to display in Slack messages
				"username": Env("LOG_SLACK_USERNAME", "GoVel Log"),

				// Emoji icon for Slack messages
				"emoji": Env("LOG_SLACK_EMOJI", ":boom:"),

				// Minimum log level (usually critical/error for Slack)
				"level": Env("LOG_LEVEL", "critical"),

				// Whether to replace placeholders in log messages
				"replace_placeholders": true,
			},

			// Papertrail Log Management Service
			//
			// Sends logs to Papertrail for centralized log management.
			// Useful for cloud applications and distributed systems.
			"papertrail": map[string]any{
				// Log driver type (uses monolog with custom handler)
				"driver": "monolog",

				// Minimum log level to send
				"level": Env("LOG_LEVEL", "debug"),

				// Handler class for UDP syslog
				"handler": "SyslogUdpHandler", // Go equivalent class

				// Handler configuration
				"handler_with": map[string]any{
					// Papertrail hostname
					"host": Env("PAPERTRAIL_URL", ""),

					// Papertrail port number
					"port": Env("PAPERTRAIL_PORT", ""),

					// Full TLS connection string
					"connectionString": "tls://" + Env("PAPERTRAIL_URL", "").(string) + ":" + Env("PAPERTRAIL_PORT", "").(string),
				},

				// Message processors for formatting
				"processors": []string{"PsrLogMessageProcessor"},
			},

			// Standard Error Stream Channel
			//
			// Writes logs directly to stderr. Useful for containerized
			// applications where logs are captured by the container runtime.
			"stderr": map[string]any{
				// Log driver type (uses monolog)
				"driver": "monolog",

				// Minimum log level to output
				"level": Env("LOG_LEVEL", "debug"),

				// Handler class for stream output
				"handler": "StreamHandler",

				// Handler configuration
				"handler_with": map[string]any{
					// Output stream (stderr)
					"stream": "stderr",
				},

				// Custom formatter for stderr output
				"formatter": Env("LOG_STDERR_FORMATTER", ""),

				// Message processors for formatting
				"processors": []string{"PsrLogMessageProcessor"},
			},

			// System Log (syslog) Channel
			//
			// Sends logs to the system's syslog daemon. Good for
			// integration with system monitoring and log aggregation.
			"syslog": map[string]any{
				// Log driver type
				"driver": "syslog",

				// Minimum log level to record
				"level": Env("LOG_LEVEL", "debug"),

				// Syslog facility (LOG_USER, LOG_LOCAL0-7, etc.)
				"facility": Env("LOG_SYSLOG_FACILITY", "LOG_USER"),

				// Whether to replace placeholders in log messages
				"replace_placeholders": true,
			},

			// PHP Error Log Channel
			//
			// Writes to PHP's error log (adapted for Go runtime logs).
			// Useful for simple logging without file management.
			"errorlog": map[string]any{
				// Log driver type
				"driver": "errorlog",

				// Minimum log level to record
				"level": Env("LOG_LEVEL", "debug"),

				// Whether to replace placeholders in log messages
				"replace_placeholders": true,
			},

			// Null Channel (No Logging)
			//
			// Discards all log messages. Useful for disabling logging
			// in certain environments or for specific channels.
			"null": map[string]any{
				// Log driver type (uses monolog)
				"driver": "monolog",

				// Handler that discards all messages
				"handler": "NullHandler",
			},

			// Emergency Fallback Channel
			//
			// Simple fallback channel used when the primary logging
			// system fails. Always writes to a basic file.
			"emergency": map[string]any{
				// Path to emergency log file
				"path": StoragePath("logs/govel.log"),
			},
		},
	}
}
