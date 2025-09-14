package config

// Filesystems returns the filesystem configuration map.
// This matches Laravel's filesystems.php configuration structure exactly.
// This configuration handles file storage drivers and disk configurations
// for managing file uploads, storage, and serving static assets.
func Filesystems() map[string]any {
	return map[string]any{

		// Default Filesystem Disk
		//
		//
		// Here you may specify the default filesystem disk that should be used
		// by the framework. The "local" disk, as well as a variety of cloud
		// based disks are available to your application for file storage.
		//
		"default": Env("FILESYSTEM_DISK", "local"),

		// Filesystem Disks
		//
		//
		// Below you may configure as many filesystem disks as necessary, and you
		// may even configure multiple disks for the same driver. Examples for
		// most supported storage drivers are configured here for reference.
		//
		// Supported drivers: "local", "ftp", "sftp", "s3"
		//
		"disks": map[string]any{

			// Local Private Storage Disk
			//
			// This disk stores files locally on the server's filesystem.
			// Files are private by default and not directly web-accessible.
			"local": map[string]any{
				// Storage driver type
				"driver": "local",

				// Root directory path for file storage
				"root": StoragePath("app/private"),

				// Whether to serve files directly (for private files)
				"serve": true,

				// Whether to throw exceptions on errors
				"throw": false,

				// Whether to report errors to logging system
				"report": false,
			},

			// Public Local Storage Disk
			//
			// This disk stores files locally with public web access.
			// Files can be accessed directly via web URLs.
			"public": map[string]any{
				// Storage driver type
				"driver": "local",

				// Root directory path for public file storage
				"root": StoragePath("app/public"),

				// Public URL base for accessing files
				"url": Env("APP_URL", "").(string) + "/storage",

				// File visibility (public means web-accessible)
				"visibility": "public",

				// Whether to throw exceptions on errors
				"throw": false,

				// Whether to report errors to logging system
				"report": false,
			},

			// Amazon S3 Cloud Storage Disk
			//
			// This disk stores files in Amazon S3 cloud storage.
			// Provides scalable, reliable cloud-based file storage.
			"s3": map[string]any{
				// Storage driver type
				"driver": "s3",

				// AWS Access Key ID for authentication
				"key": Env("AWS_ACCESS_KEY_ID", ""),

				// AWS Secret Access Key for authentication
				"secret": Env("AWS_SECRET_ACCESS_KEY", ""),

				// AWS region where the S3 bucket is located
				"region": Env("AWS_DEFAULT_REGION", ""),

				// S3 bucket name for file storage
				"bucket": Env("AWS_BUCKET", ""),

				// Custom URL for accessing files (optional)
				"url": Env("AWS_URL", ""),

				// Custom S3 endpoint (for S3-compatible services)
				"endpoint": Env("AWS_ENDPOINT", ""),

				// Use path-style URLs instead of virtual-hosted-style
				"use_path_style_endpoint": Env("AWS_USE_PATH_STYLE_ENDPOINT", false),

				// Whether to throw exceptions on errors
				"throw": false,

				// Whether to report errors to logging system
				"report": false,
			},
		},

		// Symbolic Links
		//
		//
		// Here you may configure the symbolic links that will be created when the
		// `storage:link` command is executed. The array keys should be
		// the locations of the links and the values should be their targets.
		//
		"links": map[string]string{
			PublicPath("storage"): StoragePath("app/public"),
		},
	}
}
