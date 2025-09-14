package tests

import (
	"testing"

	"govel/exceptions"
)

// TestBackwardCompatibility tests that the old API still works
func TestBackwardCompatibility(t *testing.T) {
	// Test main exception creation (should work like before)
	exc := exceptions.NewException("Test error", 400)

	if exc.GetMessage() != "Test error" {
		t.Errorf("Expected message 'Test error', got '%s'", exc.GetMessage())
	}

	if exc.GetStatusCode() != 400 {
		t.Errorf("Expected status code 400, got %d", exc.GetStatusCode())
	}

	// Test HTTP exceptions (should work like before)
	notFoundExc := exceptions.NewNotFoundException("Resource not found")

	if notFoundExc.GetStatusCode() != 404 {
		t.Errorf("Expected status code 404, got %d", notFoundExc.GetStatusCode())
	}

	if notFoundExc.GetMessage() != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", notFoundExc.GetMessage())
	}

	// Test helpers (should work like before)
	abortExc := exceptions.Abort(500, "Server error")

	if abortExc.GetStatusCode() != 500 {
		t.Error("Abort helper not working")
	}

	// Test shortcuts (should work like before)
	notFoundShortcut := exceptions.Abort404("Not found")

	if notFoundShortcut.GetStatusCode() != 404 {
		t.Error("Abort404 shortcut not working")
	}

	// Test solutions (new functionality should work)
	if !notFoundExc.HasSolution() {
		t.Error("NotFoundException should have a solution")
	}

	solution := notFoundExc.GetSolution()
	if solution == nil {
		t.Error("Solution should not be nil")
	}
}
