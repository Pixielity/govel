package security

import (
	"net/http"
	"strings"
)

// SameSitePolicy defines the SameSite cookie attribute policies.
// This enum provides Laravel-compatible SameSite attribute management
// for enhanced security against CSRF and other cross-site attacks.
//
// The SameSite attribute controls when cookies are sent in cross-site requests:
//   - Strict: Only sent with same-site requests (most secure)
//   - Lax: Sent with same-site requests and top-level navigations (default)
//   - None: Sent with all requests (requires Secure flag)
type SameSitePolicy int

const (
	// SameSiteDefault uses the default policy (Lax)
	SameSiteDefault SameSitePolicy = iota

	// SameSiteStrict only sends cookies with same-site requests.
	// This is the most secure option but may break some legitimate
	// cross-site functionality (like payment redirects).
	SameSiteStrict

	// SameSiteLax sends cookies with same-site requests and top-level navigations.
	// This is Laravel's default and provides good security while maintaining usability.
	SameSiteLax

	// SameSiteNone sends cookies with all requests (same-site and cross-site).
	// This requires the Secure flag and should only be used when necessary
	// for legitimate cross-site functionality.
	SameSiteNone
)

// String returns the string representation of the SameSite policy.
func (s SameSitePolicy) String() string {
	switch s {
	case SameSiteStrict:
		return "Strict"
	case SameSiteLax:
		return "Lax"
	case SameSiteNone:
		return "None"
	default:
		return "Lax" // Default to Lax
	}
}

// ToHTTP converts the policy to http.SameSite for use with Go's http.Cookie.
func (s SameSitePolicy) ToHTTP() http.SameSite {
	switch s {
	case SameSiteStrict:
		return http.SameSiteStrictMode
	case SameSiteLax:
		return http.SameSiteLaxMode
	case SameSiteNone:
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode // Default to Lax
	}
}

// FromHTTP converts http.SameSite to our SameSitePolicy.
func FromHTTP(s http.SameSite) SameSitePolicy {
	switch s {
	case http.SameSiteStrictMode:
		return SameSiteStrict
	case http.SameSiteLaxMode:
		return SameSiteLax
	case http.SameSiteNoneMode:
		return SameSiteNone
	default:
		return SameSiteDefault
	}
}

// SameSiteManager manages SameSite policies for different cookie types and contexts.
// This struct provides Laravel-compatible SameSite attribute management with
// support for different policies based on cookie names, paths, and request context.
//
// Features:
//   - Context-aware policy selection
//   - Cookie-specific policy overrides
//   - User agent compatibility checking
//   - Secure flag enforcement for SameSite=None
type SameSiteManager struct {
	// defaultPolicy is the default SameSite policy for all cookies
	defaultPolicy SameSitePolicy

	// cookiePolicies maps cookie names to specific SameSite policies
	cookiePolicies map[string]SameSitePolicy

	// pathPolicies maps cookie paths to specific SameSite policies
	pathPolicies map[string]SameSitePolicy

	// enforceSecure determines if Secure flag should be enforced for SameSite=None
	enforceSecure bool

	// checkUserAgent determines if user agent compatibility should be checked
	checkUserAgent bool
}

// NewSameSiteManager creates a new SameSite policy manager with Laravel defaults.
//
// Default configuration:
//   - Default policy: SameSiteLax (Laravel default)
//   - Enforce Secure: true (required for SameSite=None)
//   - Check User Agent: true (for compatibility)
//
// Parameters:
//   - options: Configuration options
func NewSameSiteManager(options ...SameSiteOption) *SameSiteManager {
	manager := &SameSiteManager{
		defaultPolicy:   SameSiteLax, // Laravel default
		cookiePolicies:  make(map[string]SameSitePolicy),
		pathPolicies:    make(map[string]SameSitePolicy),
		enforceSecure:   true,  // Required for SameSite=None
		checkUserAgent:  true,  // Check for compatibility
	}

	// Apply configuration options
	for _, option := range options {
		option(manager)
	}

	return manager
}

// GetPolicyForCookie determines the appropriate SameSite policy for a specific cookie.
// This method considers cookie name, path, and request context to select the best policy.
//
// Parameters:
//   - cookieName: Name of the cookie
//   - cookiePath: Path of the cookie
//   - r: HTTP request (for context)
//
// Returns:
//   - SameSitePolicy: The appropriate policy for this cookie
func (m *SameSiteManager) GetPolicyForCookie(cookieName, cookiePath string, r *http.Request) SameSitePolicy {
	// Check for cookie-specific policy first
	if policy, exists := m.cookiePolicies[cookieName]; exists {
		return policy
	}

	// Check for path-specific policy
	if policy, exists := m.pathPolicies[cookiePath]; exists {
		return policy
	}

	// Check for wildcard path policies
	for path, policy := range m.pathPolicies {
		if strings.HasSuffix(path, "*") {
			prefix := strings.TrimSuffix(path, "*")
			if strings.HasPrefix(cookiePath, prefix) {
				return policy
			}
		}
	}

	// Return default policy
	return m.defaultPolicy
}

// ApplySameSitePolicy applies the appropriate SameSite policy to a cookie.
// This method also enforces the Secure flag when required.
//
// Parameters:
//   - cookie: The cookie to modify
//   - r: HTTP request (for context)
//
// Returns:
//   - bool: True if the policy was applied successfully
func (m *SameSiteManager) ApplySameSitePolicy(cookie *http.Cookie, r *http.Request) bool {
	// Determine the appropriate policy
	policy := m.GetPolicyForCookie(cookie.Name, cookie.Path, r)

	// Check user agent compatibility if enabled
	if m.checkUserAgent && !m.isUserAgentCompatible(r, policy) {
		// Fall back to no SameSite attribute for incompatible browsers
		cookie.SameSite = http.SameSiteDefaultMode
		return true
	}

	// Apply the policy
	cookie.SameSite = policy.ToHTTP()

	// Enforce Secure flag for SameSite=None
	if policy == SameSiteNone && m.enforceSecure {
		cookie.Secure = true
	}

	return true
}

// isUserAgentCompatible checks if the user agent supports the given SameSite policy.
// This is important because some older browsers don't handle SameSite=None correctly.
func (m *SameSiteManager) isUserAgentCompatible(r *http.Request, policy SameSitePolicy) bool {
	// Only check compatibility for SameSite=None
	if policy != SameSiteNone {
		return true
	}

	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return true // Assume compatible if no user agent
	}

	// Check for known incompatible user agents
	// This is a simplified check - in production you might want more comprehensive detection

	// Chrome 51-66 doesn't handle SameSite=None correctly
	if strings.Contains(userAgent, "Chrome/5") ||
		strings.Contains(userAgent, "Chrome/6") {
		return false
	}

	// Safari on iOS 12 and macOS 10.14 don't handle SameSite=None correctly
	if strings.Contains(userAgent, "Safari") {
		if strings.Contains(userAgent, "Version/12") {
			return false
		}
	}

	// UC Browser before version 12.13 doesn't handle SameSite=None correctly
	if strings.Contains(userAgent, "UCBrowser") {
		// Simple check - you might want more sophisticated version parsing
		return false
	}

	return true
}

// Configuration options for SameSiteManager

// SameSiteOption defines a configuration function for the SameSite manager.
type SameSiteOption func(*SameSiteManager)

// WithDefaultPolicy sets the default SameSite policy for all cookies.
func WithDefaultPolicy(policy SameSitePolicy) SameSiteOption {
	return func(m *SameSiteManager) {
		m.defaultPolicy = policy
	}
}

// WithCookiePolicy sets a specific SameSite policy for a named cookie.
func WithCookiePolicy(cookieName string, policy SameSitePolicy) SameSiteOption {
	return func(m *SameSiteManager) {
		m.cookiePolicies[cookieName] = policy
	}
}

// WithPathPolicy sets a specific SameSite policy for a cookie path.
func WithPathPolicy(path string, policy SameSitePolicy) SameSiteOption {
	return func(m *SameSiteManager) {
		m.pathPolicies[path] = policy
	}
}

// WithEnforceSecure sets whether the Secure flag should be enforced for SameSite=None.
func WithEnforceSecure(enforce bool) SameSiteOption {
	return func(m *SameSiteManager) {
		m.enforceSecure = enforce
	}
}

// WithCheckUserAgent sets whether user agent compatibility should be checked.
func WithCheckUserAgent(check bool) SameSiteOption {
	return func(m *SameSiteManager) {
		m.checkUserAgent = check
	}
}

// Predefined policy configurations for common use cases

// StrictSameSiteConfig returns options for strict SameSite configuration.
// This configuration provides maximum security but may break some functionality.
func StrictSameSiteConfig() []SameSiteOption {
	return []SameSiteOption{
		WithDefaultPolicy(SameSiteStrict),
		WithCookiePolicy("csrf_token", SameSiteLax),     // CSRF tokens need Lax for forms
		WithCookiePolicy("XSRF-TOKEN", SameSiteLax),     // Laravel AJAX CSRF token
		WithCookiePolicy("language", SameSiteLax),       // Language preferences
		WithCookiePolicy("theme", SameSiteLax),          // Theme preferences
	}
}

// BalancedSameSiteConfig returns options for balanced SameSite configuration.
// This configuration balances security and usability (Laravel default).
func BalancedSameSiteConfig() []SameSiteOption {
	return []SameSiteOption{
		WithDefaultPolicy(SameSiteLax),
		WithCookiePolicy("api_token", SameSiteStrict),   // API tokens should be strict
		WithCookiePolicy("admin_session", SameSiteStrict), // Admin sessions should be strict
	}
}

// CompatibleSameSiteConfig returns options for maximum compatibility.
// This configuration prioritizes compatibility over security.
func CompatibleSameSiteConfig() []SameSiteOption {
	return []SameSiteOption{
		WithDefaultPolicy(SameSiteLax),
		WithCheckUserAgent(true),                       // Check for compatibility
		WithEnforceSecure(false),                       // Allow insecure for development
	}
}