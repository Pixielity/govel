// Package checks provides built-in health check implementations.
package checks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"govel/healthcheck/src/checks"
	"govel/healthcheck/src/enums"
	"govel/healthcheck/src/interfaces"
)

// PingCheck performs HTTP connectivity checks to verify endpoint availability.
// It closely mirrors the Laravel health PingCheck pattern.
type PingCheck struct {
	*checks.BaseCheck

	// url is the HTTP(S) endpoint to check
	url string

	// failureMessage is the custom message to use when the check fails
	failureMessage string

	// timeout is the maximum time to wait for the request to complete
	timeout time.Duration

	// retryTimes is the number of times to retry a failed request
	retryTimes int

	// method is the HTTP method to use (default: GET)
	method string

	// headers are additional HTTP headers to send with the request
	headers map[string]string
}

// NewPingCheck creates a new PingCheck instance with default settings.
//
// Returns:
//
//	*PingCheck: A new ping check instance
func NewPingCheck() *PingCheck {
	return &PingCheck{
		BaseCheck:  checks.NewBaseCheck(),
		timeout:    1 * time.Second,
		retryTimes: 1,
		method:     "GET",
		headers:    make(map[string]string),
	}
}

// URL sets the endpoint URL to check.
//
// Parameters:
//
//	url: The URL to check
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) URL(url string) *PingCheck {
	pc.url = url
	return pc
}

// Timeout sets the timeout for the HTTP request.
//
// Parameters:
//
//	seconds: Timeout in seconds
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) Timeout(seconds int) *PingCheck {
	pc.timeout = time.Duration(seconds) * time.Second
	return pc
}

// Method sets the HTTP method to use.
//
// Parameters:
//
//	method: The HTTP method (GET, POST, etc.)
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) Method(method string) *PingCheck {
	pc.method = method
	return pc
}

// RetryTimes sets the number of times to retry a failed request.
//
// Parameters:
//
//	times: Number of retry attempts
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) RetryTimes(times int) *PingCheck {
	pc.retryTimes = times
	return pc
}

// Headers sets HTTP headers to include with the request.
//
// Parameters:
//
//	headers: Map of header name to header value
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) Headers(headers map[string]string) *PingCheck {
	pc.headers = headers
	return pc
}

// FailureMessage sets a custom message to display when the check fails.
//
// Parameters:
//
//	failureMessage: Custom failure message
//
// Returns:
//
//	*PingCheck: Self for method chaining
func (pc *PingCheck) FailureMessage(failureMessage string) *PingCheck {
	pc.failureMessage = failureMessage
	return pc
}

// Run performs the ping health check.
//
// Returns:
//
//	interfaces.ResultInterface: The health check result
func (pc *PingCheck) Run() interfaces.ResultInterface {
	// If URL is not set, return a failed result
	if pc.url == "" {
		return checks.NewResult().
			SetStatus(enums.StatusFailed).
			SetShortSummary("URL not set").
			SetNotificationMessage("URL for ping check is not set")
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: pc.timeout,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), pc.timeout)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, pc.method, pc.url, nil)
	if err != nil {
		return pc.failedResult()
	}

	// Add headers
	for key, value := range pc.headers {
		req.Header.Set(key, value)
	}

	// Perform request with retries
	var response *http.Response
	var lastErr error

	for attempt := 0; attempt <= pc.retryTimes; attempt++ {
		response, lastErr = client.Do(req)
		if lastErr == nil && response != nil && response.StatusCode >= 200 && response.StatusCode < 400 {
			// Success
			if response.Body != nil {
				defer response.Body.Close()
			}
			return checks.NewResult().
				SetStatus(enums.StatusOK).
				SetShortSummary("Reachable")
		}

		// Don't retry after the last attempt
		if attempt < pc.retryTimes {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Clean up response if it exists
	if response != nil && response.Body != nil {
		response.Body.Close()
	}

	// Request failed
	return pc.failedResult()
}

// failedResult creates a standardized result for failed checks
func (pc *PingCheck) failedResult() interfaces.ResultInterface {
	notificationMessage := pc.failureMessage
	if notificationMessage == "" {
		notificationMessage = fmt.Sprintf("Pinging %s failed.", pc.GetName())
	}

	return checks.NewResult().
		SetStatus(enums.StatusFailed).
		SetShortSummary("Unreachable").
		SetNotificationMessage(notificationMessage)
}

// Compile-time interface compliance check
var _ interfaces.CheckInterface = (*PingCheck)(nil)
