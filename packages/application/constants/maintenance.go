package constants

// Maintenance Mode Constants
// These constants define values used for maintenance mode functionality
// throughout the GoVel application.

const (
	// DefaultMaintenanceMessage is the default message shown during maintenance mode
	DefaultMaintenanceMessage = "Application is currently undergoing maintenance. Please try again later."

	// DefaultMaintenanceRetryAfter is the default retry-after value in seconds
	DefaultMaintenanceRetryAfter = 60

	// MaintenanceFileActiveValue is the value indicating maintenance mode is active
	MaintenanceFileActiveValue = true

	// MaintenanceFileInactiveValue is the value indicating maintenance mode is inactive
	MaintenanceFileInactiveValue = false
)

// Maintenance Mode Bypass Types
const (
	// BypassTypeIP indicates bypass by IP address
	BypassTypeIP = "ip"

	// BypassTypePath indicates bypass by URL path
	BypassTypePath = "path"

	// BypassTypeSecret indicates bypass by secret token
	BypassTypeSecret = "secret"

	// BypassTypeUser indicates bypass by user authentication
	BypassTypeUser = "user"

	// BypassTypeRole indicates bypass by user role
	BypassTypeRole = "role"
)

// Common Maintenance Messages
const (
	// MaintenanceMessageScheduled is for scheduled maintenance
	MaintenanceMessageScheduled = "We're performing scheduled maintenance and will be back shortly."

	// MaintenanceMessageUpgrade is for system upgrades
	MaintenanceMessageUpgrade = "We're upgrading our system to serve you better. Please check back soon."

	// MaintenanceMessageEmergency is for emergency maintenance
	MaintenanceMessageEmergency = "We're experiencing technical difficulties and are working to resolve them quickly."

	// MaintenanceMessageDatabase is for database maintenance
	MaintenanceMessageDatabase = "Database maintenance is in progress. Service will be restored shortly."

	// MaintenanceMessageDeploy is for deployment maintenance
	MaintenanceMessageDeploy = "New features are being deployed. We'll be back online soon!"
)

// Maintenance Retry-After Values (in seconds)
const (
	// RetryAfterShort is for short maintenance periods (1 minute)
	RetryAfterShort = 60

	// RetryAfterMedium is for medium maintenance periods (15 minutes)
	RetryAfterMedium = 900

	// RetryAfterLong is for long maintenance periods (1 hour)
	RetryAfterLong = 3600

	// RetryAfterExtended is for extended maintenance periods (4 hours)
	RetryAfterExtended = 14400
)

// Default Allowed IPs for bypass
var (
	// DefaultAllowedIPs contains common localhost IP addresses
	DefaultAllowedIPs = []string{
		"127.0.0.1", // IPv4 localhost
		"::1",       // IPv6 localhost
	}
)

// Default Allowed Paths for bypass
var (
	// DefaultAllowedPaths contains paths that should remain accessible during maintenance
	DefaultAllowedPaths = []string{
		"/health",      // Health check endpoint
		"/status",      // Status endpoint
		"/metrics",     // Metrics endpoint
		"/admin",       // Admin panel
		"/maintenance", // Maintenance status page
	}
)

// Maintenance Status Codes
const (
	// StatusCodeMaintenance is the HTTP status code returned during maintenance
	StatusCodeMaintenance = 503 // Service Unavailable

	// StatusCodeBypass is the status code when maintenance is bypassed
	StatusCodeBypass = 200 // OK

	// StatusCodeUnauthorizedBypass is the status code for unauthorized bypass attempts
	StatusCodeUnauthorizedBypass = 403 // Forbidden
)

// Maintenance Data Keys
// These constants are used as keys in the maintenance mode data map
const (
	// DataKeyEstimatedDuration is the key for estimated maintenance duration
	DataKeyEstimatedDuration = "estimated_duration"

	// DataKeyContactInfo is the key for contact information during maintenance
	DataKeyContactInfo = "contact_info"

	// DataKeyProgressPercentage is the key for maintenance progress percentage
	DataKeyProgressPercentage = "progress_percentage"

	// DataKeyLastUpdate is the key for last update timestamp
	DataKeyLastUpdate = "last_update"

	// DataKeyMaintenanceType is the key for maintenance type
	DataKeyMaintenanceType = "maintenance_type"

	// DataKeyAffectedServices is the key for list of affected services
	DataKeyAffectedServices = "affected_services"
)

// Maintenance Types
const (
	// MaintenanceTypeScheduled indicates scheduled maintenance
	MaintenanceTypeScheduled = "scheduled"

	// MaintenanceTypeEmergency indicates emergency maintenance
	MaintenanceTypeEmergency = "emergency"

	// MaintenanceTypeDatabase indicates database maintenance
	MaintenanceTypeDatabase = "database"

	// MaintenanceTypeDeployment indicates deployment maintenance
	MaintenanceTypeDeployment = "deployment"

	// MaintenanceTypeUpgrade indicates system upgrade maintenance
	MaintenanceTypeUpgrade = "upgrade"

	// MaintenanceTypeSecurity indicates security-related maintenance
	MaintenanceTypeSecurity = "security"
)
