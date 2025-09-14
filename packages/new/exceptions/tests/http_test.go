package tests

import (
	"net/http"
	"testing"

	"govel/exceptions"
)

func TestAllHTTPExceptions(t *testing.T) {
	testCases := []struct {
		name           string
		constructor    func(...string) exceptions.ExceptionInterface
		expectedStatus int
		defaultMessage string
	}{
		{"BadRequest", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewBadRequestException(msgs...)
		}, http.StatusBadRequest, "Bad Request"},

		{"Unauthorized", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewUnauthorizedException(msgs...)
		}, http.StatusUnauthorized, "Unauthorized"},

		{"Forbidden", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewForbiddenException(msgs...)
		}, http.StatusForbidden, "Forbidden"},

		{"NotFound", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewNotFoundException(msgs...)
		}, http.StatusNotFound, "Not Found"},

		{"UnprocessableEntity", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewUnprocessableEntityException(msgs...)
		}, http.StatusUnprocessableEntity, "Unprocessable Entity"},

		{"InternalServerError", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewInternalServerErrorException(msgs...)
		}, http.StatusInternalServerError, "Internal Server Error"},

		{"ServiceUnavailable", func(msgs ...string) exceptions.ExceptionInterface {
			if len(msgs) > 0 {
				return exceptions.NewServiceUnavailableException(msgs[0])
			}
			return exceptions.NewServiceUnavailableException("")
		}, http.StatusServiceUnavailable, "Service Unavailable"},

		{"Conflict", func(msgs ...string) exceptions.ExceptionInterface {
			return exceptions.NewConflictException(msgs...)
		}, http.StatusConflict, "Conflict"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with default message
			exc := tc.constructor()

			if exc.GetStatusCode() != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, exc.GetStatusCode())
			}

			if exc.GetMessage() != tc.defaultMessage {
				t.Errorf("Expected message '%s', got '%s'", tc.defaultMessage, exc.GetMessage())
			}

			// Test with custom message
			customMessage := "Custom " + tc.name + " message"
			excCustom := tc.constructor(customMessage)

			if excCustom.GetMessage() != customMessage {
				t.Errorf("Expected custom message '%s', got '%s'", customMessage, excCustom.GetMessage())
			}

			// Test that exception has solution
			if !exc.HasSolution() {
				t.Errorf("%s should have a solution", tc.name)
			}

			solution := exc.GetSolution()
			if solution == nil {
				t.Errorf("%s solution should not be nil", tc.name)
			}

			if solution.GetSolutionTitle() == "" {
				t.Errorf("%s solution should have a title", tc.name)
			}

			if solution.GetSolutionDescription() == "" {
				t.Errorf("%s solution should have a description", tc.name)
			}

			// Test exception implements all interfaces
			testInterfaces(t, exc, tc.name)
		})
	}
}

func TestSpecialHTTPExceptions(t *testing.T) {
	// Test MethodNotAllowed with allowed methods
	methodNotAllowed := exceptions.NewMethodNotAllowedException("Method not allowed", "GET", "POST")

	if methodNotAllowed.GetStatusCode() != http.StatusMethodNotAllowed {
		t.Error("MethodNotAllowed should have 405 status")
	}

	headers := methodNotAllowed.GetHeaders()
	if headers["Allow"] != "GET, POST" {
		t.Errorf("Expected Allow header 'GET, POST', got '%s'", headers["Allow"])
	}

	// Test TooManyRequests with retry after
	tooManyRequests := exceptions.NewTooManyRequestsException("Rate limited", 60)

	if tooManyRequests.GetStatusCode() != http.StatusTooManyRequests {
		t.Error("TooManyRequests should have 429 status")
	}

	headers = tooManyRequests.GetHeaders()
	if headers["Retry-After"] != "60" {
		t.Errorf("Expected Retry-After header '60', got '%s'", headers["Retry-After"])
	}

	// Test ServiceUnavailable with retry after
	serviceUnavailable := exceptions.NewServiceUnavailableException("Under maintenance", 3600)

	if serviceUnavailable.GetStatusCode() != http.StatusServiceUnavailable {
		t.Error("ServiceUnavailable should have 503 status")
	}

	headers = serviceUnavailable.GetHeaders()
	if headers["Retry-After"] != "3600" {
		t.Errorf("Expected Retry-After header '3600', got '%s'", headers["Retry-After"])
	}
}

func testInterfaces(t *testing.T, exc exceptions.ExceptionInterface, name string) {
	// Test HTTPable
	if exc.GetStatusCode() == 0 {
		t.Errorf("%s HTTPable interface: status code should not be 0", name)
	}

	// Test Contextable
	exc.WithContext("test", "value")
	context := exc.GetContext()
	if context["test"] != "value" {
		t.Errorf("%s Contextable interface: context not working", name)
	}

	// Test Renderable
	rendered := exc.Render()
	if rendered == nil {
		t.Errorf("%s Renderable interface: render should not return nil", name)
	}

	if rendered["error"] != true {
		t.Errorf("%s Renderable interface: error field should be true", name)
	}

	// Test Stackable
	stack := exc.GetStackTrace()
	if len(stack) == 0 {
		t.Errorf("%s Stackable interface: should have stack trace", name)
	}

	// Test Solutionable
	if !exc.HasSolution() {
		t.Errorf("%s Solutionable interface: should have solution", name)
	}
}
