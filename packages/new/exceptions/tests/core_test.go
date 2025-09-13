package tests

import (
	"net/http"
	"testing"

	"govel/exceptions/core"
	"govel/exceptions/helpers"
	httpExceptions "govel/exceptions/http"
	"govel/exceptions/interfaces"
)

// TestCoreExceptionCreation tests basic exception creation
func TestCoreExceptionCreation(t *testing.T) {
	exc := core.NewException("Test error", 400)

	if exc.GetMessage() != "Test error" {
		t.Errorf("Expected message 'Test error', got '%s'", exc.GetMessage())
	}

	if exc.GetStatusCode() != 400 {
		t.Errorf("Expected status code 400, got %d", exc.GetStatusCode())
	}

	if exc.Error() != "Test error" {
		t.Errorf("Expected Error() to return 'Test error', got '%s'", exc.Error())
	}
}

// TestISPInterfaces tests that ISP interfaces work correctly
func TestISPInterfaces(t *testing.T) {
	exc := core.NewException("Test error", 500)

	// Test HTTPable interface
	var httpable interfaces.HTTPable = exc
	if httpable.GetStatusCode() != 500 {
		t.Error("HTTPable interface not working correctly")
	}

	// Test Contextable interface
	var contextable interfaces.Contextable = exc
	contextable.WithContext("test", "value")
	context := contextable.GetContext()
	if context["test"] != "value" {
		t.Error("Contextable interface not working correctly")
	}

	// Test Renderable interface
	var renderable interfaces.Renderable = exc
	rendered := renderable.Render()
	if rendered["status_code"] != 500 {
		t.Error("Renderable interface not working correctly")
	}

	// Test Stackable interface
	var stackable interfaces.Stackable = exc
	stack := stackable.GetStackTrace()
	if len(stack) == 0 {
		t.Error("Stackable interface not working correctly")
	}

	// Test Solutionable interface
	var solutionable interfaces.Solutionable = exc
	if solutionable.HasSolution() {
		t.Error("Should not have solution initially")
	}
}

// TestHTTPException tests HTTP exception types
func TestHTTPException(t *testing.T) {
	exc := httpExceptions.NewNotFoundException("Resource not found")

	if exc.GetStatusCode() != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, exc.GetStatusCode())
	}

	if exc.GetMessage() != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", exc.GetMessage())
	}

	if !exc.HasSolution() {
		t.Error("NotFoundException should have a solution")
	}

	solution := exc.GetSolution()
	if solution.GetSolutionTitle() != "Resource Not Found" {
		t.Errorf("Expected solution title 'Resource Not Found', got '%s'", solution.GetSolutionTitle())
	}
}

// TestHelperFunctions tests helper functions
func TestHelperFunctions(t *testing.T) {
	// Test Abort
	exc := helpers.Abort(404, "Not found")
	if exc.GetStatusCode() != 404 {
		t.Error("Abort function not working correctly")
	}

	// Test AbortIf
	exc = helpers.AbortIf(true, 400, "Bad request")
	if exc == nil {
		t.Error("AbortIf should return exception when condition is true")
	}

	exc = helpers.AbortIf(false, 400, "Bad request")
	if exc != nil {
		t.Error("AbortIf should return nil when condition is false")
	}

	// Test shortcuts
	exc = helpers.Abort404("Not found")
	if exc.GetStatusCode() != 404 {
		t.Error("Abort404 shortcut not working correctly")
	}
}
