package config

// Auth returns the authentication configuration map.
// This matches Laravel's auth.php configuration structure exactly.
func Auth() map[string]any {
	return map[string]any{

		// Authentication Defaults
		//
		// This option defines the default authentication "guard" and password
		// reset "broker" for your application. You may change these values
		// as required, but they're a perfect start for most applications.
		"defaults": map[string]any{
			// Default Authentication Guard
			//
			// The default guard to use for authentication. This guard will be
			// used when no specific guard is specified in authentication calls.
			"guard": Env("AUTH_GUARD", "web"),

			// Default Password Broker
			//
			// The default password reset broker to use when no specific broker
			// is specified for password reset operations.
			"passwords": Env("AUTH_PASSWORD_BROKER", "users"),
		},

		// Authentication Guards
		//
		// Next, you may define every authentication guard for your application.
		// Of course, a great default configuration has been defined for you
		// which utilizes session storage plus the Eloquent user provider.
		// Supported: "session"
		"guards": map[string]any{
			// Web Guard Configuration
			//
			// The default web guard using session-based authentication.
			// This guard is typically used for web-based user interfaces.
			"web": map[string]any{
				// Authentication Driver
				// The driver determines how authentication state is maintained
				"driver": "session",

				// User Provider
				// The user provider defines how users are retrieved from storage
				"provider": "users",
			},
		},

		// User Providers
		//
		//
		// All authentication guards have a user provider, which defines how the
		// users are actually retrieved out of your database or other storage
		// system used by the application. Typically, Eloquent is utilized.
		//
		// If you have multiple user tables or models you may configure multiple
		// providers to represent the model / table. These providers may then
		// be assigned to any extra authentication guards you have defined.
		//
		// Supported: "database", "eloquent"
		//
		"providers": map[string]any{
			// Users Provider Configuration
			//
			// Configuration for the primary user provider. This defines how
			// users are retrieved from storage for authentication purposes.
			"users": map[string]any{
				// Provider Driver Type
				//
				// The driver determines how users are retrieved from storage.
				// "eloquent" uses ORM models, "database" uses direct queries.
				"driver": "eloquent",

				// User Model Class
				//
				// The fully qualified class name of the User model when using
				// the eloquent driver. This model represents user entities.
				"model": Env("AUTH_MODEL", "App\\Models\\User"),
			},
		},

		// Resetting Passwords
		//
		//
		// These configuration options specify the behavior of the password
		// reset functionality, including the table utilized for token storage
		// and the user provider that is invoked to actually retrieve users.
		//
		// The expiry time is the number of minutes that each reset token will be
		// considered valid. This security feature keeps tokens short-lived so
		// they have less time to be guessed. You may change this as needed.
		//
		// The throttle setting is the number of seconds a user must wait before
		// generating more password reset tokens. This prevents the user from
		// quickly generating a very large amount of password reset tokens.
		//
		"passwords": map[string]any{
			// User Password Reset Configuration
			//
			// Configuration for password reset functionality for the users table.
			// This defines how password reset tokens are managed and validated.
			"users": map[string]any{
				// Password Reset Provider
				//
				// The user provider to use for password resets. This should match
				// one of the providers defined in the providers section above.
				"provider": "users",

				// Password Reset Token Table
				//
				// The database table where password reset tokens will be stored.
				// This table stores temporary tokens used for password reset links.
				"table": Env("AUTH_PASSWORD_RESET_TOKEN_TABLE", "password_reset_tokens"),

				// Token Expiration Time
				//
				// The number of minutes that each reset token will be considered valid.
				// After this time, the token expires and cannot be used for reset.
				"expire": 60,

				// Reset Request Throttle
				//
				// The number of seconds a user must wait before generating more
				// password reset tokens. This prevents spam and abuse.
				"throttle": 60,
			},
		},

		// Password Confirmation Timeout
		//
		//
		// Here you may define the number of seconds before a password confirmation
		// window expires and users are asked to re-enter their password via the
		// confirmation screen. By default, the timeout lasts for three hours.
		//
		"password_timeout": Env("AUTH_PASSWORD_TIMEOUT", 10800),
	}
}
