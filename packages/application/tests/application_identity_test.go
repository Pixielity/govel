package tests

import (
	"testing"

	"govel/packages/application"
)

// TestApplicationIdentity tests application identity methods
func TestApplicationIdentity(t *testing.T) {
	app := application.New()

	// Test name methods
	testName := "Test Application"
	app.SetName(testName)

	if app.GetName() != testName {
		t.Errorf("Expected name %s, got %s", testName, app.GetName())
	}

	if app.GetName() != testName {
		t.Errorf("Expected Name() to return %s, got %s", testName, app.GetName())
	}

	// Test version methods
	testVersion := "2.1.0"
	app.SetVersion(testVersion)

	if app.GetVersion() != testVersion {
		t.Errorf("Expected version %s, got %s", testVersion, app.GetVersion())
	}

	if app.GetVersion() != testVersion {
		t.Errorf("Expected Version() to return %s, got %s", testVersion, app.GetVersion())
	}
}
