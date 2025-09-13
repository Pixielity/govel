package cookie

import (
	types "govel/packages/types/src/types/cookie"
	"net/http"
	"time"
)

// WithExpiry sets the cookie expiration time.
func WithExpiry(expiry time.Time) types.CookieOption {
	return func(c *http.Cookie) {
		c.Expires = expiry
	}
}

// WithMaxAge sets the cookie max age in seconds.
func WithMaxAge(maxAge int) types.CookieOption {
	return func(c *http.Cookie) {
		c.MaxAge = maxAge
	}
}

// WithDomain sets the cookie domain.
func WithDomain(domain string) types.CookieOption {
	return func(c *http.Cookie) {
		c.Domain = domain
	}
}

// WithPath sets the cookie path.
func WithPath(path string) types.CookieOption {
	return func(c *http.Cookie) {
		c.Path = path
	}
}

// WithSecure sets whether the cookie should only be sent over HTTPS.
func WithSecure(secure bool) types.CookieOption {
	return func(c *http.Cookie) {
		c.Secure = secure
	}
}

// WithHttpOnly sets whether the cookie should be HTTP-only (not accessible via JavaScript).
func WithHttpOnly(httpOnly bool) types.CookieOption {
	return func(c *http.Cookie) {
		c.HttpOnly = httpOnly
	}
}

// WithSameSite sets the SameSite attribute of the cookie.
func WithSameSite(sameSite http.SameSite) types.CookieOption {
	return func(c *http.Cookie) {
		c.SameSite = sameSite
	}
}
