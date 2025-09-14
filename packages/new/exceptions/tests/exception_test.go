package tests

import (
	"net/http"
	"testing"
	"time"

	"govel/exceptions"
)

func TestNewException(t *testing.T) {
	message := "Test error"
	statusCode := http.StatusBadRequest

	exc := exceptions.NewException(message, statusCode)

	if exc.GetMessage() != message {
		t.Errorf("Expected message '%s', got '%s'", message, exc.GetMessage())
	}

	if exc.GetStatusCode() != statusCode {
		t.Errorf("Expected status code %d, got %d", statusCode, exc.GetStatusCode())
	}

	if exc.Error() != message {
		t.Errorf("Expected Error() to return '%s', got '%s'", message, exc.Error())
	}
}

func TestExceptionDefaultStatusCode(t *testing.T) {
	exc := exceptions.NewException("Test error")

	if exc.GetStatusCode() != http.StatusInternalServerError {
		t.Errorf("Expected default status code %d, got %d", http.StatusInternalServerError, exc.GetStatusCode())
	}
}

func TestExceptionWithHeader(t *testing.T) {
	exc := exceptions.NewException("Test error", http.StatusUnauthorized)
	exc.WithHeader("WWW-Authenticate", "Bearer")

	headers := exc.GetHeaders()
	if headers["WWW-Authenticate"] != "Bearer" {
		t.Errorf("Expected header 'WWW-Authenticate: Bearer', got '%s'", headers["WWW-Authenticate"])
	}
}

func TestExceptionWithContext(t *testing.T) {
	exc := exceptions.NewException("Test error")
	exc.WithContext("field", "email")
	exc.WithContext("validation", "required")

	context := exc.GetContext()
	if context["field"] != "email" {
		t.Errorf("Expected context field 'email', got '%v'", context["field"])
	}
	if context["validation"] != "required" {
		t.Errorf("Expected context validation 'required', got '%v'", context["validation"])
	}
}

func TestNotFoundException(t *testing.T) {
	exc := exceptions.NewNotFoundException("Resource not found")

	if exc.GetStatusCode() != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, exc.GetStatusCode())
	}

	if exc.GetMessage() != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", exc.GetMessage())
	}
}

func TestUnauthorizedException(t *testing.T) {
	exc := exceptions.NewUnauthorizedException("Authentication required")

	if exc.GetStatusCode() != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, exc.GetStatusCode())
	}

	// Should have WWW-Authenticate header
	headers := exc.GetHeaders()
	if headers["WWW-Authenticate"] != "Bearer" {
		t.Errorf("Expected WWW-Authenticate header, got '%v'", headers["WWW-Authenticate"])
	}
}

func TestAbortFunction(t *testing.T) {
	exc := exceptions.Abort(404, "Not found")

	if exc.GetStatusCode() != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, exc.GetStatusCode())
	}

	if exc.GetMessage() != "Not found" {
		t.Errorf("Expected message 'Not found', got '%s'", exc.GetMessage())
	}
}

func TestAbortIf(t *testing.T) {
	// Test condition true - should return exception
	exc := exceptions.AbortIf(true, 400, "Bad request")
	if exc == nil {
		t.Error("Expected exception when condition is true, got nil")
	}
	if exc.GetStatusCode() != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, exc.GetStatusCode())
	}

	// Test condition false - should return nil
	exc = exceptions.AbortIf(false, 400, "Bad request")
	if exc != nil {
		t.Error("Expected nil when condition is false, got exception")
	}
}

func TestAbortUnless(t *testing.T) {
	// Test condition false - should return exception
	exc := exceptions.AbortUnless(false, 401, "Unauthorized")
	if exc == nil {
		t.Error("Expected exception when condition is false, got nil")
	}
	if exc.GetStatusCode() != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, exc.GetStatusCode())
	}

	// Test condition true - should return nil
	exc = exceptions.AbortUnless(true, 401, "Unauthorized")
	if exc != nil {
		t.Error("Expected nil when condition is true, got exception")
	}
}

func TestShorthandAbortFunctions(t *testing.T) {
	testCases := []struct {
		function     func(...string) exceptions.ExceptionInterface
		expectedCode int
		name         string
	}{
		{exceptions.Abort400, http.StatusBadRequest, "Abort400"},
		{exceptions.Abort401, http.StatusUnauthorized, "Abort401"},
		{exceptions.Abort403, http.StatusForbidden, "Abort403"},
		{exceptions.Abort404, http.StatusNotFound, "Abort404"},
		{exceptions.Abort500, http.StatusInternalServerError, "Abort500"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exc := tc.function("Test message")
			if exc.GetStatusCode() != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, exc.GetStatusCode())
			}
		})
	}
}

func TestExceptionRender(t *testing.T) {
	exc := exceptions.NewBadRequestException("Invalid input").
		WithHeader("Content-Type", "application/json").
		WithContext("field", "email")

	rendered := exc.Render()

	// Check required fields
	if rendered["error"] != true {
		t.Error("Expected error field to be true")
	}

	if rendered["status_code"] != http.StatusBadRequest {
		t.Errorf("Expected status_code %d, got %v", http.StatusBadRequest, rendered["status_code"])
	}

	if rendered["message"] != "Invalid input" {
		t.Errorf("Expected message 'Invalid input', got %v", rendered["message"])
	}

	// Check timestamp format
	if timestamp, ok := rendered["timestamp"].(string); ok {
		if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
			t.Errorf("Timestamp is not in RFC3339 format: %s", timestamp)
		}
	} else {
		t.Error("Expected timestamp to be a string")
	}

	// Check context and headers
	if rendered["context"] == nil {
		t.Error("Expected context to be present")
	}

	if rendered["headers"] == nil {
		t.Error("Expected headers to be present")
	}
}

func TestMethodNotAllowedException(t *testing.T) {
	exc := exceptions.NewMethodNotAllowedException("Method not allowed", "GET", "POST")

	if exc.GetStatusCode() != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, exc.GetStatusCode())
	}

	headers := exc.GetHeaders()
	if headers["Allow"] != "GET, POST" {
		t.Errorf("Expected Allow header 'GET, POST', got '%s'", headers["Allow"])
	}
}

func TestTooManyRequestsException(t *testing.T) {
	exc := exceptions.NewTooManyRequestsException("Rate limited", 60)

	if exc.GetStatusCode() != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, exc.GetStatusCode())
	}

	headers := exc.GetHeaders()
	if headers["Retry-After"] != "60" {
		t.Errorf("Expected Retry-After header '60', got '%s'", headers["Retry-After"])
	}
}

func TestExceptionStackTrace(t *testing.T) {
	exc := exceptions.NewException("Test error")

	stackTrace := exc.GetStackTrace()
	if len(stackTrace) == 0 {
		t.Error("Expected stack trace to be captured")
	}
}

func TestExceptionIsHTTP(t *testing.T) {
	// HTTP status code
	exc := exceptions.NewException("Test", 404)
	if !exc.IsHTTPException() {
		t.Error("Expected exception with HTTP status to be HTTP exception")
	}

	// Non-HTTP status code
	exc = exceptions.NewException("Test", 999)
	if exc.IsHTTPException() {
		t.Error("Expected exception with non-HTTP status to not be HTTP exception")
	}
}
