package tests

import (
	"testing"
	"time"

	"govel/packages/application"
)

// TestApplicationTiming tests application timing functionality
func TestApplicationTiming(t *testing.T) {
	app := application.New()
	
	startTime := time.Now()
	app.SetStartTime(startTime)
	
	retrievedStartTime := app.GetStartTime()
	if !retrievedStartTime.Equal(startTime) {
		t.Errorf("Expected start time %v, got %v", startTime, retrievedStartTime)
	}
	
	// Test uptime calculation
	time.Sleep(10 * time.Millisecond) // Small delay to test uptime
	uptime := app.GetUptime()
	
	if uptime <= 0 {
		t.Error("Expected uptime to be positive")
	}
	
	if uptime < 10*time.Millisecond {
		t.Error("Expected uptime to be at least 10ms")
	}
}
