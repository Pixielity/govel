package cookie

import "net/http"

// QueuedCookie represents a cookie that has been queued for later processing.
// This struct provides additional metadata about when and why the cookie was queued.
type QueuedCookie struct {
	// Cookie is the actual HTTP cookie instance
	Cookie *http.Cookie

	// QueuedAt indicates when the cookie was added to the queue
	QueuedAt int64

	// Priority allows for ordering cookies when processing the queue
	// Lower numbers indicate higher priority (0 = highest priority)
	Priority int

	// Metadata can store additional information about the cookie
	// This might include the source of the cookie, processing instructions, etc.
	Metadata map[string]interface{}
}

// GetKey returns a unique key for the queued cookie based on name and path.
// This key is used for indexing cookies in the queue and checking for duplicates.
func (qc *QueuedCookie) GetKey() string {
	path := qc.Cookie.Path
	if path == "" {
		path = "/"
	}
	return qc.Cookie.Name + "@" + path
}

// IsExpired checks if the queued cookie has already expired.
// This can be used to clean up expired cookies from the queue.
func (qc *QueuedCookie) IsExpired() bool {
	if qc.Cookie.Expires.IsZero() {
		return false // Session cookies don't expire
	}
	return qc.Cookie.Expires.Before(qc.Cookie.Expires) // This should use current time
}
