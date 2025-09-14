package config

import (
	"os"
	"strconv"
	"strings"
)

// Env retrieves an environment variable with an optional default value.
// This is a helper function similar to Laravel's env() helper.
func Env(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Try to parse common types
	switch defaultValue.(type) {
	case bool:
		if lower := strings.ToLower(value); lower == "true" || lower == "1" || lower == "yes" || lower == "on" {
			return true
		}
		if lower := strings.ToLower(value); lower == "false" || lower == "0" || lower == "no" || lower == "off" {
			return false
		}
		return defaultValue
	case int:
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		return defaultValue
	case int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
		return defaultValue
	case float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
		return defaultValue
	default:
		return value
	}
}

// App returns the application configuration map.
// This configuration defines the core application settings
// and follows the Laravel/govel configuration pattern.
func App() map[string]any {
	return map[string]any{
		// Application Name
		//
		// This value is the name of your application. This value is used when the
		// framework needs to place the application's name in a notification or
		// any other location as required by the application or its packages.
		"name": Env("APP_NAME", "GoVel"),

		// Application Environment
		//
		// This value determines the "environment" your application is currently
		// running in. This may determine how you prefer to configure various
		// services the application utilizes. Set this in your ".env" file.
		"env": Env("APP_ENV", "production"),

		// Application Debug Mode
		//
		// When your application is in debug mode, detailed error messages with
		// stack traces will be shown on every error that occurs within your
		// application. If disabled, a simple generic error page is shown.
		"debug": Env("APP_DEBUG", false),

		// Application URL
		//
		// This URL is used by the console to properly generate URLs when using
		// the Artisan command line tool. You should set this to the root of
		// your application so that it is used when running Artisan tasks.
		"url": Env("APP_URL", "http://localhost:8080"),

		// Application Timezone
		//
		// Here you may specify the default timezone for your application, which
		// will be used by the Go time functions. You are free to set this
		// to any of the timezones which will be supported by the application.
		// Reference: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
		"timezone": Env("APP_TIMEZONE", "UTC"),

		// Application Locale Configuration
		//
		// The application locale determines the default locale that will be used
		// by the translation service provider. You are free to set this value
		// to any of the locales which will be supported by the application.
		"locale": Env("APP_LOCALE", "en"),

		// Application Fallback Locale
		//
		// The fallback locale determines the locale to use when the current one
		// is not available. You may change the value to correspond to any of
		// the language folders that are provided through your application.
		"fallback_locale": Env("APP_FALLBACK_LOCALE", "en"),

		// Application Lang Path
		//
		// The path to the language files for the application. You may change
		// the path to a different directory if you would like to customize it.
		"lang_path": Env("APP_LANG_PATH", "lang"),

		// Encryption Key
		//
		// This key is used by the encryption services and should be set
		// to a random, 32 character string, otherwise these encrypted strings
		// will not be safe. Please do this before deploying an application!
		"key": Env("APP_KEY", ""),

		// Application Version
		//
		// This value represents the version of your application. This version
		// number is used throughout the application for various purposes.
		"version": Env("APP_VERSION", "1.0.0"),

		// Application Cipher
		//
		// This cipher is used by the encryption services to encrypt data.
		// The cipher must be supported by the encryption library.
		"cipher": Env("APP_CIPHER", "AES-256-CBC"),

		// Maintenance Mode
		//
		// When maintenance mode is enabled, a maintenance page will be shown
		// for all requests to your application, except for administrators.
		"maintenance": map[string]any{
			"enabled": Env("MAINTENANCE_MODE", false),
			"secret":  Env("MAINTENANCE_SECRET", ""),
			"message": Env("MAINTENANCE_MESSAGE", "We're currently performing scheduled maintenance. Please check back soon."),
			"allowed": strings.Split(Env("MAINTENANCE_ALLOWED_IPS", "").(string), ","),
		},

		// Asset Configuration
		//
		// These settings control how assets are served and cached.
		"asset": map[string]any{
			"url":     Env("ASSET_URL", ""),
			"path":    Env("ASSET_PATH", "/assets"),
			"version": Env("ASSET_VERSION", ""),
		},

		// Application Features
		//
		// Here you can enable or disable various application features
		// based on environment variables or hardcoded values.
		"features": map[string]any{
			"registration_enabled": Env("FEATURE_REGISTRATION", true),
			"email_verification":   Env("FEATURE_EMAIL_VERIFICATION", false),
			"password_reset":       Env("FEATURE_PASSWORD_RESET", true),
			"api_enabled":          Env("FEATURE_API", true),
			"web_enabled":          Env("FEATURE_WEB", true),
			"admin_panel":          Env("FEATURE_ADMIN_PANEL", true),
			"user_profiles":        Env("FEATURE_USER_PROFILES", true),
			"notifications":        Env("FEATURE_NOTIFICATIONS", true),
			"real_time_updates":    Env("FEATURE_REAL_TIME", false),
			"file_uploads":         Env("FEATURE_FILE_UPLOADS", true),
			"social_login":         Env("FEATURE_SOCIAL_LOGIN", false),
			"two_factor_auth":      Env("FEATURE_2FA", false),
			"audit_logging":        Env("FEATURE_AUDIT_LOG", false),
			"rate_limiting":        Env("FEATURE_RATE_LIMITING", true),
			"caching":              Env("FEATURE_CACHING", true),
			"search":               Env("FEATURE_SEARCH", false),
			"analytics":            Env("FEATURE_ANALYTICS", false),
			"backups":              Env("FEATURE_BACKUPS", false),
			"health_checks":        Env("FEATURE_HEALTH_CHECKS", true),
		},

		// Application Limits
		//
		// These settings define various limits for the application.
		"limits": map[string]any{
			"max_upload_size":       Env("LIMIT_MAX_UPLOAD_SIZE", "10MB"),
			"max_request_size":      Env("LIMIT_MAX_REQUEST_SIZE", "50MB"),
			"max_users_per_account": Env("LIMIT_MAX_USERS", 100),
			"max_files_per_user":    Env("LIMIT_MAX_FILES", 1000),
			"session_lifetime":      Env("LIMIT_SESSION_LIFETIME", 120), // minutes
			"password_min_length":   Env("LIMIT_PASSWORD_MIN_LENGTH", 8),
			"username_min_length":   Env("LIMIT_USERNAME_MIN_LENGTH", 3),
			"max_login_attempts":    Env("LIMIT_LOGIN_ATTEMPTS", 5),
			"api_rate_limit":        Env("LIMIT_API_RATE", 1000), // requests per hour
		},

		// Security Configuration
		//
		// These settings control various security aspects of the application.
		"security": map[string]any{
			"password_timeout": Env("SECURITY_PASSWORD_TIMEOUT", 10800), // 3 hours in seconds
			"session_encrypt":  Env("SECURITY_SESSION_ENCRYPT", true),
			"csrf_protection":  Env("SECURITY_CSRF_PROTECTION", true),
			"force_https":      Env("SECURITY_FORCE_HTTPS", false),
			"hsts_enabled":     Env("SECURITY_HSTS", false),
			"content_sniffing": Env("SECURITY_CONTENT_SNIFFING", false),
			"frame_guard":      Env("SECURITY_FRAME_GUARD", true),
			"xss_protection":   Env("SECURITY_XSS_PROTECTION", true),
			"referrer_policy":  Env("SECURITY_REFERRER_POLICY", "strict-origin-when-cross-origin"),
		},

		// Performance Settings
		//
		// These settings control performance-related aspects of the application.
		"performance": map[string]any{
			"cache_config":        Env("PERFORMANCE_CACHE_CONFIG", true),
			"cache_routes":        Env("PERFORMANCE_CACHE_ROUTES", false),
			"cache_views":         Env("PERFORMANCE_CACHE_VIEWS", false),
			"optimize_autoloader": Env("PERFORMANCE_OPTIMIZE_AUTOLOADER", false),
			"preload_enabled":     Env("PERFORMANCE_PRELOAD", false),
			"opcache_enabled":     Env("PERFORMANCE_OPCACHE", true),
			"gzip_enabled":        Env("PERFORMANCE_GZIP", true),
			"minify_assets":       Env("PERFORMANCE_MINIFY_ASSETS", false),
		},

		// Third-party Service Configuration
		//
		// Configuration for external services and APIs.
		"services": map[string]any{
			"analytics": map[string]any{
				"google_analytics_id": Env("GOOGLE_ANALYTICS_ID", ""),
				"enabled":             Env("ANALYTICS_ENABLED", false),
			},
			"monitoring": map[string]any{
				"sentry_dsn": Env("SENTRY_DSN", ""),
				"enabled":    Env("MONITORING_ENABLED", false),
			},
			"payment": map[string]any{
				"stripe_key":    Env("STRIPE_PUBLIC_KEY", ""),
				"stripe_secret": Env("STRIPE_SECRET_KEY", ""),
				"paypal_client": Env("PAYPAL_CLIENT_ID", ""),
				"sandbox_mode":  Env("PAYMENT_SANDBOX", true),
			},
		},

		// Social Authentication Providers
		//
		// Configuration for OAuth providers.
		"auth": map[string]any{
			"providers": map[string]any{
				"google": map[string]any{
					"enabled":       Env("AUTH_GOOGLE_ENABLED", false),
					"client_id":     Env("GOOGLE_CLIENT_ID", ""),
					"client_secret": Env("GOOGLE_CLIENT_SECRET", ""),
					"redirect_url":  Env("GOOGLE_REDIRECT_URL", ""),
				},
				"github": map[string]any{
					"enabled":       Env("AUTH_GITHUB_ENABLED", false),
					"client_id":     Env("GITHUB_CLIENT_ID", ""),
					"client_secret": Env("GITHUB_CLIENT_SECRET", ""),
					"redirect_url":  Env("GITHUB_REDIRECT_URL", ""),
				},
				"facebook": map[string]any{
					"enabled":       Env("AUTH_FACEBOOK_ENABLED", false),
					"client_id":     Env("FACEBOOK_CLIENT_ID", ""),
					"client_secret": Env("FACEBOOK_CLIENT_SECRET", ""),
					"redirect_url":  Env("FACEBOOK_REDIRECT_URL", ""),
				},
				"twitter": map[string]any{
					"enabled":       Env("AUTH_TWITTER_ENABLED", false),
					"client_id":     Env("TWITTER_CLIENT_ID", ""),
					"client_secret": Env("TWITTER_CLIENT_SECRET", ""),
					"redirect_url":  Env("TWITTER_REDIRECT_URL", ""),
				},
			},
			"defaults": map[string]any{
				"guard":     Env("AUTH_GUARD", "web"),
				"passwords": Env("AUTH_PASSWORD_BROKER", "users"),
			},
		},

		// Developer Settings
		//
		// Settings useful during development.
		"dev": map[string]any{
			"show_errors":   Env("DEV_SHOW_ERRORS", false),
			"log_queries":   Env("DEV_LOG_QUERIES", false),
			"profiling":     Env("DEV_PROFILING", false),
			"hot_reload":    Env("DEV_HOT_RELOAD", false),
			"mock_external": Env("DEV_MOCK_EXTERNAL", false),
			"seed_database": Env("DEV_SEED_DATABASE", false),
			"fake_emails":   Env("DEV_FAKE_EMAILS", false),
		},
	}
}
