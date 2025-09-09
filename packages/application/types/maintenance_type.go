package types

import "time"

/**
 * MaintenanceMode represents the maintenance mode configuration.
 */
type MaintenanceMode struct {
	/**
	 * Active indicates whether maintenance mode is enabled
	 */
	Active bool `json:"active"`

	/**
	 * Message is the message displayed during maintenance
	 */
	Message string `json:"message"`

	/**
	 * RetryAfter is the number of seconds to retry after (for HTTP Retry-After header)
	 */
	RetryAfter int `json:"retry_after"`

	/**
	 * AllowedIPs contains IP addresses that can bypass maintenance mode
	 */
	AllowedIPs []string `json:"allowed_ips,omitempty"`

	/**
	 * AllowedPaths contains URL paths that can bypass maintenance mode
	 */
	AllowedPaths []string `json:"allowed_paths,omitempty"`

	/**
	 * Secret is a secret token that can be used to bypass maintenance mode
	 */
	Secret string `json:"secret,omitempty"`

	/**
	 * StartTime is when maintenance mode was activated
	 */
	StartTime time.Time `json:"start_time,omitempty"`

	/**
	 * EstimatedDuration is the estimated duration of maintenance
	 */
	EstimatedDuration time.Duration `json:"estimated_duration,omitempty"`

	/**
	 * MaintenanceType indicates the type of maintenance being performed
	 */
	MaintenanceType string `json:"maintenance_type,omitempty"`

	/**
	 * Data contains additional maintenance-related data
	 */
	Data map[string]interface{} `json:"data,omitempty"`
}
