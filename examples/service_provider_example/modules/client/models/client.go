package models

import "time"

// Client represents a client entity in the system.
// This model defines the structure for client data and business logic.
//
// Each client has basic information like name, email, and company details,
// plus tracking information for registration and last activity.
type Client struct {
	// ID is the unique identifier for the client
	ID int `json:"id"`

	// Name is the client's full name
	Name string `json:"name"`

	// Email is the client's email address
	Email string `json:"email"`

	// Company is the client's company name
	Company string `json:"company"`

	// Phone is the client's phone number
	Phone string `json:"phone"`

	// Status represents the client's current status (active, inactive, suspended)
	Status string `json:"status"`

	// RegisteredAt is the timestamp when the client was first registered
	RegisteredAt time.Time `json:"registered_at"`

	// LastActivityAt is the timestamp of the client's last activity
	LastActivityAt time.Time `json:"last_activity_at"`
}

// ClientStatus defines the possible client statuses
type ClientStatus struct {
	Active    string
	Inactive  string
	Suspended string
}

// GetClientStatuses returns available client statuses
func GetClientStatuses() ClientStatus {
	return ClientStatus{
		Active:    "active",
		Inactive:  "inactive",
		Suspended: "suspended",
	}
}

// IsActive checks if the client is active
func (c *Client) IsActive() bool {
	return c.Status == GetClientStatuses().Active
}

// GetDisplayName returns a formatted display name for the client
func (c *Client) GetDisplayName() string {
	if c.Company != "" {
		return c.Name + " (" + c.Company + ")"
	}
	return c.Name
}
