package config

// Hashing returns the hashing configuration map.
// This configuration handles password hashing algorithms and their settings
// for secure password storage and verification in the application.
func Hashing() map[string]any {
	return map[string]any{

		// Default Hash Driver
		//
		// This option controls the default hash driver that will be used to hash
		// passwords for your application. By default, the bcrypt algorithm is
		// used; however, you remain free to modify this option if you wish.
		// Supported: "bcrypt", "argon", "argon2id"
		"driver": Env("HASH_DRIVER", "bcrypt"),

		// Bcrypt Options
		//
		// Here you may specify the configuration options that should be used when
		// passwords are hashed using the Bcrypt algorithm. This will allow you
		// to control the amount of time it takes to hash the given password.
		"bcrypt": map[string]any{
			// Cost Factor - Number of iterations (4-31, higher = more secure but slower)
			"rounds": Env("BCRYPT_ROUNDS", 12),

			// Enable hash verification
			"verify": true,
		},

		// Argon2 Options
		//
		// Here you may specify the configuration options that should be used when
		// passwords are hashed using the Argon2 algorithm. These will allow you
		// to control the amount of time it takes to hash the given password.
		"argon": map[string]any{
			// Memory usage in KB
			"memory": Env("ARGON_MEMORY", 65536),

			// Number of iterations
			"time": Env("ARGON_TIME", 4),

			// Number of parallel threads
			"threads": Env("ARGON_THREADS", 3),

			// Enable hash verification
			"verify": true,
		},

		// Scrypt Options
		//
		// Here you may specify the configuration options that should be used when
		// passwords are hashed using the Scrypt algorithm. These will allow you
		// to control the amount of time it takes to hash the given password.
		"scrypt": map[string]any{
			// CPU/memory cost parameter (must be power of 2)
			"n": Env("SCRYPT_N", 16384),

			// Block size parameter
			"r": Env("SCRYPT_R", 8),

			// Parallelization parameter
			"p": Env("SCRYPT_P", 1),

			// Length of derived key in bytes
			"length": Env("SCRYPT_LENGTH", 32),

			// Enable hash verification
			"verify": true,
		},
	}
}
