package config

// Jwt returns the JWT (JSON Web Token) configuration map.
// This configuration handles JWT authentication settings including
// secret keys, token expiration times, and JWT-specific options.
func Jwt() map[string]any {
	return map[string]any{

		// JWT Secret Key
		//
		// This key is used to sign and verify JWT tokens. It should be a long,
		// random string that is kept secure. If this key is compromised, all
		// existing tokens will be invalidated.
		"secret": Env("JWT_SECRET", ""),

		// JWT Algorithm
		//
		// The algorithm used to sign JWT tokens. Common algorithms include
		// HS256 (HMAC with SHA-256) and RS256 (RSA Signature with SHA-256).
		// Supported: "HS256", "HS384", "HS512", "RS256", "RS384", "RS512"
		"algorithm": Env("JWT_ALGORITHM", "HS256"),

		// Access Token TTL (Time To Live)
		//
		// How long access tokens remain valid (in minutes). Shorter lifetimes
		// are more secure but require more frequent token refreshes.
		"access_token_ttl": Env("JWT_ACCESS_TOKEN_TTL", 15), // minutes

		// Refresh Token TTL
		//
		// How long refresh tokens remain valid (in minutes). These can have
		// longer lifetimes since they're used less frequently.
		"refresh_token_ttl": Env("JWT_REFRESH_TOKEN_TTL", 20160), // minutes (2 weeks)

		// JWT Claims
		//
		//
		// Default claims that will be included in JWT tokens. These can be
		// overridden when generating specific tokens.
		//
		"claims": map[string]any{
			"issuer":   Env("JWT_ISSUER", Env("APP_NAME", "GoVel").(string)),
			"audience": Env("JWT_AUDIENCE", Env("APP_URL", "http://localhost").(string)),
			"subject":  Env("JWT_SUBJECT", "user"),
		},

		// RSA Keys (for RS256, RS384, RS512)
		//
		//
		// If using RSA algorithms, specify the paths to your RSA private and
		// public key files. The private key is used for signing tokens, and
		// the public key is used for verification.
		//
		"rsa": map[string]any{
			"private_key_path": Env("JWT_RSA_PRIVATE_KEY_PATH", ""),
			"public_key_path":  Env("JWT_RSA_PUBLIC_KEY_PATH", ""),
			"passphrase":       Env("JWT_RSA_PASSPHRASE", ""),
		},

		// Token Storage
		//
		//
		// Configuration for storing and managing JWT tokens, including blacklisting
		// and token invalidation strategies.
		//
		"storage": map[string]any{
			"driver": Env("JWT_STORAGE_DRIVER", "memory"), // memory, redis, database
			"prefix": Env("JWT_STORAGE_PREFIX", "jwt_tokens:"),
			"blacklist": map[string]any{
				"enabled":          Env("JWT_BLACKLIST_ENABLED", true),
				"grace_period":     Env("JWT_BLACKLIST_GRACE_PERIOD", 30),      // seconds
				"refresh_ttl":      Env("JWT_BLACKLIST_REFRESH_TTL", 20160),    // minutes
				"cleanup_interval": Env("JWT_BLACKLIST_CLEANUP_INTERVAL", 300), // seconds
			},
		},

		// Request Configuration
		//
		//
		// Settings for how JWT tokens are extracted from HTTP requests,
		// including headers, query parameters, and cookies.
		//
		"request": map[string]any{
			"header":        Env("JWT_REQUEST_HEADER", "Authorization"),
			"header_prefix": Env("JWT_REQUEST_HEADER_PREFIX", "Bearer "),
			"query_param":   Env("JWT_REQUEST_QUERY_PARAM", "token"),
			"cookie_name":   Env("JWT_REQUEST_COOKIE_NAME", "jwt_token"),
			"input_key":     Env("JWT_REQUEST_INPUT_KEY", "token"),
		},

		// Required Claims
		//
		//
		// List of claims that must be present in JWT tokens for them to be
		// considered valid. This adds an extra layer of security.
		//
		"required_claims": []string{
			"iss", // issuer
			"aud", // audience
			"exp", // expiration time
			"iat", // issued at
			"sub", // subject
		},

		// Leeway
		//
		//
		// Leeway time in seconds to account for clock skew between servers
		// when validating time-based claims (exp, nbf, iat).
		//
		"leeway": Env("JWT_LEEWAY", 0), // seconds

		// Providers
		//
		//
		// JWT can be used with different user providers. This allows you to
		// specify how user information should be retrieved when validating tokens.
		//
		"providers": map[string]any{
			"user": map[string]any{
				"model":       Env("JWT_USER_MODEL", "App\\Models\\User"),
				"table":       Env("JWT_USER_TABLE", "users"),
				"primary_key": Env("JWT_USER_PRIMARY_KEY", "id"),
			},
		},
	}
}
