package config

// Mail returns the mail configuration map.
// This matches Laravel's mail.php configuration structure exactly.
// This configuration handles email sending drivers, SMTP settings,
// and global email settings for the application.
func Mail() map[string]any {
	return map[string]any{

		// Default Mailer
		//
		//
		// This option controls the default mailer that is used to send all email
		// messages unless another mailer is explicitly specified when sending
		// the message. All additional mailers can be configured within the
		// "mailers" array. Examples of each type of mailer are provided.
		//
		"default": Env("MAIL_MAILER", "log"),

		// Mailer Configurations
		//
		//
		// Here you may configure all of the mailers used by your application plus
		// their respective settings. Several examples have been configured for
		// you and you are free to add your own as your application requires.
		//
		// The framework supports a variety of mail "transport" drivers that can be used
		// when delivering an email. You may specify which one you're using for
		// your mailers below. You may also add additional mailers if needed.
		//
		// Supported: "smtp", "sendmail", "mailgun", "ses", "ses-v2",
		// "postmark", "resend", "log", "array",
		// "failover", "roundrobin"
		//
		"mailers": map[string]any{

			// SMTP Mail Configuration
			//
			// Standard SMTP configuration for sending emails through any
			// SMTP server. This is the most common email sending method.
			"smtp": map[string]any{
				// Mail transport driver
				"transport": "smtp",

				// Connection scheme (usually "smtp" or "smtps")
				"scheme": Env("MAIL_SCHEME", ""),

				// Full SMTP URL (optional, overrides individual settings)
				"url": Env("MAIL_URL", ""),

				// SMTP server hostname
				"host": Env("MAIL_HOST", "127.0.0.1"),

				// SMTP server port (25, 587, or 465 typically)
				"port": Env("MAIL_PORT", 2525),

				// SMTP authentication username
				"username": Env("MAIL_USERNAME", ""),

				// SMTP authentication password
				"password": Env("MAIL_PASSWORD", ""),

				// Connection timeout in seconds
				"timeout": nil,

				// Local domain for EHLO command
				"local_domain": Env("MAIL_EHLO_DOMAIN", Env("APP_URL", "http://localhost").(string)),
			},

			// Amazon SES (Simple Email Service)
			//
			// Amazon's cloud-based email service. Requires AWS credentials
			// to be configured in your environment or AWS config files.
			"ses": map[string]any{
				// Mail transport driver
				"transport": "ses",

				// AWS Region (optional, defaults to us-east-1)
				// "region": Env("AWS_DEFAULT_REGION", "us-east-1"),

				// AWS Access Key ID (optional if using IAM roles)
				// "key": Env("AWS_ACCESS_KEY_ID", ""),

				// AWS Secret Access Key (optional if using IAM roles)
				// "secret": Env("AWS_SECRET_ACCESS_KEY", ""),

				// AWS Session Token (for temporary credentials)
				// "token": Env("AWS_SESSION_TOKEN", ""),
			},

			// Postmark Email Service
			//
			// Postmark is a transactional email service with high delivery rates.
			// Requires a Postmark server token to be set in environment variables.
			"postmark": map[string]any{
				// Mail transport driver
				"transport": "postmark",

				// Postmark Server Token (required)
				// "token": Env("POSTMARK_TOKEN", ""),

				// Message Stream ID (optional for organizing messages)
				// "message_stream_id": Env("POSTMARK_MESSAGE_STREAM_ID", ""),

				// HTTP Client Configuration
				// "client": map[string]any{
				//     // Request timeout in seconds
				//     "timeout": Env("POSTMARK_TIMEOUT", 5),
				// },
			},

			// Resend Email Service
			//
			// Modern email API with good deliverability and developer experience.
			// Requires a Resend API key to be configured.
			"resend": map[string]any{
				// Mail transport driver
				"transport": "resend",

				// Resend API Key (required)
				// "key": Env("RESEND_KEY", ""),
			},

			// Sendmail Local Transport
			//
			// Uses the local sendmail binary to send emails. Good for servers
			// with properly configured local mail systems.
			"sendmail": map[string]any{
				// Mail transport driver
				"transport": "sendmail",

				// Path to sendmail binary with flags
				// -bs: Use SMTP mode, -i: Don't treat single dots as message end
				"path": Env("MAIL_SENDMAIL_PATH", "/usr/sbin/sendmail -bs -i"),
			},

			// Log Transport (Development)
			//
			// Logs emails instead of sending them. Useful for development
			// and testing where you want to see email content without delivery.
			"log": map[string]any{
				// Mail transport driver
				"transport": "log",

				// Log channel to write emails to (empty = default channel)
				"channel": Env("MAIL_LOG_CHANNEL", ""),
			},

			// Array Transport (Testing)
			//
			// Stores emails in memory array instead of sending them.
			// Perfect for unit testing email functionality.
			"array": map[string]any{
				// Mail transport driver
				"transport": "array",
			},

			// Failover Transport
			//
			// Tries multiple mailers in sequence until one succeeds.
			// Provides redundancy and fault tolerance for email delivery.
			"failover": map[string]any{
				// Mail transport driver
				"transport": "failover",

				// List of mailers to try in order
				"mailers": []string{
					"smtp",
					"log",
				},

				// Seconds to wait before retrying a failed mailer
				"retry_after": 60,
			},

			// Round Robin Transport
			//
			// Distributes emails across multiple mailers in rotation.
			// Helps balance load and avoid rate limits.
			"roundrobin": map[string]any{
				// Mail transport driver
				"transport": "roundrobin",

				// List of mailers to rotate through
				"mailers": []string{
					"ses",
					"postmark",
				},

				// Seconds to wait before retrying a failed mailer
				"retry_after": 60,
			},
		},

		// Global "From" Address
		//
		//
		// You may wish for all emails sent by your application to be sent from
		// the same address. Here you may specify a name and address that is
		// used globally for all emails that are sent by your application.
		//
		"from": map[string]any{
			// Default From Email Address
			//
			// The email address that will be used as the "from" address
			// for all outgoing emails unless explicitly overridden.
			"address": Env("MAIL_FROM_ADDRESS", "hello@example.com"),

			// Default From Name
			//
			// The name that will be displayed as the sender name
			// for all outgoing emails unless explicitly overridden.
			"name": Env("MAIL_FROM_NAME", "Example"),
		},
	}
}
