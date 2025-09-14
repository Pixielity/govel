package config

// App returns the application configuration map.
// This matches Laravel's app.php configuration structure exactly.
func App() map[string]any {
	return map[string]any{

		// Application Name
		//
		// This value is the name of your application, which will be used when the
		// framework needs to place the application's name in a notification or
		// other UI elements where an application name needs to be displayed.
		"name": Env("APP_NAME", "GoVel"),

		// Application Environment
		//
		// This value determines the "environment" your application is currently
		// running in. This may determine how you prefer to configure various
		// services the application utilizes. Set this in your ".env" file.
		"env": Env("APP_ENV", "production"),

		// Application Debug Mode
		//
		//
		// When your application is in debug mode, detailed error messages with
		// stack traces will be shown on every error that occurs within your
		// application. If disabled, a simple generic error page is shown.
		//
		"debug": Env("APP_DEBUG", false),

		// Application URL
		//
		//
		// This URL is used by the console to properly generate URLs when using
		// the command line tool. You should set this to the root of
		// the application so that it's available within CLI commands.
		//
		"url": Env("APP_URL", "http://localhost"),

		// Application Timezone
		//
		//
		// Here you may specify the default timezone for your application, which
		// will be used by the time functions. The timezone
		// is set to "UTC" by default as it is suitable for most use cases.
		//
		"timezone": "UTC",

		// Application Locale Configuration
		//
		//
		// The application locale determines the default locale that will be used
		// by the translation / localization methods. This option can be
		// set to any locale for which you plan to have translation strings.
		//
		"locale": Env("APP_LOCALE", "en"),

		// Application Fallback Locale
		//
		// The fallback locale determines the locale to use when the current one
		// is not available. You may change the value to correspond to any of
		// the language folders that are provided through your application.
		"fallback_locale": Env("APP_FALLBACK_LOCALE", "en"),

		// Application Cipher
		//
		// This cipher is used by the encryption services to encrypt data.
		// The cipher must be supported by the encryption library.
		"cipher": "AES-256-CBC",

		// Encryption Key
		//
		// This key is utilized by the encryption services and should be set
		// to a random, 32 character string to ensure that all encrypted values
		// are secure. You should do this prior to deploying the application.
		"key": Env("APP_KEY", ""),

		// Maintenance Mode Driver
		//
		//
		// These configuration options determine the driver used to determine and
		// manage the application's "maintenance mode" status. The "cache" driver will
		// allow maintenance mode to be controlled across multiple machines.
		//
		// Supported drivers: "file", "cache"
		//
		"maintenance": map[string]any{
			// Maintenance Mode Driver
			//
			// The driver used to determine maintenance mode status.
			// "file" stores status in a file, "cache" uses cache store.
			"driver": Env("APP_MAINTENANCE_DRIVER", "file"),

			// Maintenance Mode Store
			//
			// The cache store to use when using "cache" driver for maintenance mode.
			// This must correspond to a cache store defined in your cache config.
			"store": Env("APP_MAINTENANCE_STORE", "database"),
		},
	}
}
