// Package controllers provides HTTP controllers for health check endpoints.
// These controllers handle HTTP requests and provide health status responses.
package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"govel/packages/healthcheck/src/interfaces"
)

// HealthController handles HTTP requests for health check endpoints.
// It provides multiple endpoints for different health check scenarios.
type HealthController struct {
	// registry is the health check registry
	registry interfaces.HealthRegistryInterface

	// resultStore is the result storage backend
	resultStore interfaces.ResultStoreInterface
}

// NewHealthController creates a new health controller instance.
//
// Parameters:
//   registry: The health check registry to use
//
// Returns:
//   *HealthController: A new controller instance
func NewHealthController(registry interfaces.HealthRegistryInterface) *HealthController {
	return &HealthController{
		registry: registry,
	}
}

// HandleHealthCheck handles the main health check endpoint.
// Responds with JSON by default, HTML if Accept header includes text/html.
//
// Endpoints:
//   GET /health - Main health check endpoint
func (hc *HealthController) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Convert standard HTTP interfaces to our interfaces
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Execute health checks
	results := hc.registry.RunChecks(ctx)

	// Determine response format based on Accept header
	acceptHeader := r.Header.Get("Accept")
	
	if strings.Contains(acceptHeader, "text/html") {
		hc.handleHTMLResponse(w, r, results)
	} else {
		hc.handleJSONResponse(w, r, results)
	}
}

// HandleHealthCheckJSON handles JSON health check endpoint.
// Always responds with JSON format.
//
// Endpoints:
//   GET /health.json - JSON health check endpoint
func (hc *HealthController) HandleHealthCheckJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Check for fresh results query parameter
	if r.URL.Query().Get("fresh") == "true" {
		results := hc.registry.RunChecks(ctx)
		hc.handleJSONResponse(w, r, results)
		return
	}

	// Try to get cached results if result store is available
	if hc.resultStore != nil {
		if storedResults, err := hc.resultStore.Get(); err == nil && storedResults != nil {
			hc.handleJSONResponse(w, r, storedResults)
			return
		}
	}

	// Fall back to running fresh checks
	results := hc.registry.RunChecks(ctx)
	hc.handleJSONResponse(w, r, results)
}

// HandleSimpleHealthCheck handles simple text health check endpoint.
// Responds with simple "OK" or "FAILED" text.
//
// Endpoints:
//   GET /health/simple - Simple text health check
func (hc *HealthController) HandleSimpleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	results := hc.registry.RunChecks(ctx)
	
	// Set content type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

	// Determine overall status
	if results.ContainsFailingCheck() {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("FAILED"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// HandleReadinessCheck handles Kubernetes-style readiness probe.
// This endpoint should return 200 when the application is ready to serve traffic.
//
// Endpoints:
//   GET /health/ready - Readiness probe
func (hc *HealthController) HandleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// For readiness, we might want to run only critical checks
	// For now, we'll run all checks but you could filter by tags or names
	results := hc.registry.RunChecks(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	// Readiness fails if any critical checks fail
	if results.ContainsFailingCheck() {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Return minimal JSON response
	if jsonResponse, err := results.ToJSON(); err == nil {
		w.Write([]byte(jsonResponse))
	} else {
		w.Write([]byte(`{"status": "error", "message": "Failed to generate response"}`))
	}
}

// HandleLivenessCheck handles Kubernetes-style liveness probe.
// This endpoint should return 200 when the application is alive and not deadlocked.
//
// Endpoints:
//   GET /health/live - Liveness probe
func (hc *HealthController) HandleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	// For liveness, we typically want very basic checks that just verify
	// the application is responsive and not deadlocked
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	// Simple alive check - if we can respond, we're alive
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{
		"status": "alive",
		"timestamp": "%s",
		"uptime": "unknown"
	}`, time.Now().Format(time.RFC3339))
	
	w.Write([]byte(response))
}

// SetRegistry sets the health registry for the controller.
//
// Parameters:
//   registry: The health registry to use
//
// Returns:
//   *HealthController: Self for method chaining (interface compatible)
func (hc *HealthController) SetRegistry(registry interfaces.HealthRegistryInterface) *HealthController {
	hc.registry = registry
	return hc
}

// GetRegistry returns the configured health registry.
//
// Returns:
//   interfaces.HealthRegistryInterface: The configured registry
func (hc *HealthController) GetRegistry() interfaces.HealthRegistryInterface {
	return hc.registry
}

// SetResultStore sets the result store for the controller.
//
// Parameters:
//   store: The result store to use
//
// Returns:
//   *HealthController: Self for method chaining (interface compatible)
func (hc *HealthController) SetResultStore(store interfaces.ResultStoreInterface) *HealthController {
	hc.resultStore = store
	return hc
}

// GetResultStore returns the configured result store.
//
// Returns:
//   interfaces.ResultStoreInterface: The configured result store
func (hc *HealthController) GetResultStore() interfaces.ResultStoreInterface {
	return hc.resultStore
}

// handleJSONResponse handles JSON response formatting.
func (hc *HealthController) handleJSONResponse(w http.ResponseWriter, r *http.Request, results interfaces.CheckResultsInterface) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

	// Set HTTP status based on health check results
	httpStatus := http.StatusOK
	if results.ContainsFailingCheck() {
		httpStatus = http.StatusServiceUnavailable
	}

	w.WriteHeader(httpStatus)

	// Generate JSON response
	if jsonResponse, err := results.ToJSON(); err == nil {
		w.Write([]byte(jsonResponse))
	} else {
		// Fallback error response
		errorResponse := fmt.Sprintf(`{
			"status": "error",
			"message": "Failed to generate health check response: %s",
			"timestamp": "%s"
		}`, err.Error(), time.Now().Format(time.RFC3339))
		w.Write([]byte(errorResponse))
	}
}

// handleHTMLResponse handles HTML response formatting.
func (hc *HealthController) handleHTMLResponse(w http.ResponseWriter, r *http.Request, results interfaces.CheckResultsInterface) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")

	// Set HTTP status based on health check results
	httpStatus := http.StatusOK
	if results.ContainsFailingCheck() {
		httpStatus = http.StatusServiceUnavailable
	}

	w.WriteHeader(httpStatus)

	// Generate HTML response
	html := hc.generateHTMLResponse(results)
	w.Write([]byte(html))
}

// generateHTMLResponse creates a basic HTML dashboard for health checks.
func (hc *HealthController) generateHTMLResponse(results interfaces.CheckResultsInterface) string {
	summary := results.GetHealthSummary()
	overallStatus := "healthy"
	overallColor := "#28a745"
	
	if results.ContainsFailingCheck() {
		overallStatus = "unhealthy"
		overallColor = "#dc3545"
	} else if results.ContainsWarningCheck() {
		overallStatus = "degraded"
		overallColor = "#ffc107"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Health Check Dashboard</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6; 
            margin: 0; 
            padding: 20px; 
            background-color: #f5f5f5;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: %s; 
            color: white; 
            padding: 20px; 
            text-align: center;
        }
        .summary { 
            padding: 20px; 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); 
            gap: 20px;
            background: #f8f9fa;
        }
        .stat { 
            text-align: center; 
            padding: 15px;
            background: white;
            border-radius: 6px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .stat-number { 
            font-size: 2em; 
            font-weight: bold; 
            margin-bottom: 5px;
        }
        .checks { 
            padding: 20px; 
        }
        .check { 
            display: flex; 
            justify-content: space-between; 
            align-items: center;
            padding: 15px;
            margin-bottom: 10px;
            border-left: 4px solid;
            background: #f8f9fa;
            border-radius: 0 4px 4px 0;
        }
        .check.ok { border-color: #28a745; }
        .check.warning { border-color: #ffc107; }
        .check.failed { border-color: #dc3545; }
        .check-name { font-weight: 600; }
        .check-status { padding: 4px 12px; border-radius: 20px; font-size: 0.85em; font-weight: 500; }
        .status-ok { background: #d4edda; color: #155724; }
        .status-warning { background: #fff3cd; color: #856404; }
        .status-failed { background: #f8d7da; color: #721c24; }
        .timestamp { 
            text-align: center; 
            padding: 20px; 
            color: #6c757d; 
            border-top: 1px solid #dee2e6;
            font-size: 0.9em;
        }
        .duration { font-size: 0.85em; color: #6c757d; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏥 Health Check Dashboard</h1>
            <p>System Status: <strong>%s</strong></p>
        </div>
        
        <div class="summary">
            <div class="stat">
                <div class="stat-number">%d</div>
                <div>Total Checks</div>
            </div>
            <div class="stat">
                <div class="stat-number" style="color: #28a745;">%d</div>
                <div>Healthy</div>
            </div>
            <div class="stat">
                <div class="stat-number" style="color: #ffc107;">%d</div>
                <div>Warnings</div>
            </div>
            <div class="stat">
                <div class="stat-number" style="color: #dc3545;">%d</div>
                <div>Failed</div>
            </div>
            <div class="stat">
                <div class="stat-number">%dms</div>
                <div>Total Duration</div>
            </div>
        </div>
        
        <div class="checks">
            <h2>Check Results</h2>`,
		overallColor,
		strings.Title(overallStatus),
		int(summary["total"].(int)),
		int(summary["healthy"].(int)),
		int(summary["warnings"].(int)),
		int(summary["failed"].(int)),
		int(summary["total_duration_ms"].(int64)),
	)

	// Add individual check results
	for _, result := range results.GetResults() {
		status := result.GetStatus()
		if status == nil {
			continue
		}

		checkName := "Unknown"
		if result.GetCheck() != nil {
			checkName = result.GetCheck().GetName()
		}

		statusClass := "ok"
		statusText := "OK"
		statusColor := "status-ok"

		switch status.String() {
		case "warning":
			statusClass = "warning"
			statusText = "WARNING"
			statusColor = "status-warning"
		case "failed", "crashed":
			statusClass = "failed"
			statusText = "FAILED"
			statusColor = "status-failed"
		}

		duration := result.GetDuration().Milliseconds()
		
		html += fmt.Sprintf(`
            <div class="check %s">
                <div>
                    <div class="check-name">%s</div>
                    <div class="duration">%dms</div>
                </div>
                <div class="check-status %s">%s</div>
            </div>`,
			statusClass,
			checkName,
			duration,
			statusColor,
			statusText,
		)
	}

	// Add footer
	executedAt := time.Now()
	if results.GetExecutedAt() != nil {
		executedAt = *results.GetExecutedAt()
	}

	html += fmt.Sprintf(`
        </div>
        
        <div class="timestamp">
            Last updated: %s
        </div>
    </div>
    
    <script>
        // Auto-refresh every 30 seconds
        setTimeout(function() {
            window.location.reload();
        }, 30000);
    </script>
</body>
</html>`, executedAt.Format("2006-01-02 15:04:05 MST"))

	return html
}

// RegisterRoutes is a helper method to register routes with common HTTP frameworks.
// This method can be used with frameworks like Gin, Echo, or standard http.ServeMux.
//
// Parameters:
//   mux: HTTP request multiplexer (e.g., http.ServeMux, gin.Engine)
//
// Example usage:
//   mux := http.NewServeMux()
//   controller.RegisterRoutes(mux)
func (hc *HealthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", hc.HandleHealthCheck)
	mux.HandleFunc("/health.json", hc.HandleHealthCheckJSON)
	mux.HandleFunc("/health/simple", hc.HandleSimpleHealthCheck)
	mux.HandleFunc("/health/ready", hc.HandleReadinessCheck)
	mux.HandleFunc("/health/live", hc.HandleLivenessCheck)
}

// WithTimeout returns an HTTP handler that wraps the health check with a timeout.
// This is useful for preventing health checks from hanging indefinitely.
//
// Parameters:
//   handler: The handler function to wrap
//   timeout: Maximum time to wait for the handler to complete
//
// Returns:
//   http.HandlerFunc: Wrapped handler with timeout
func WithTimeout(handler http.HandlerFunc, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)
		handler(w, r)
	}
}

// WithCORS returns an HTTP handler that adds CORS headers.
// This is useful for allowing browser-based health check monitoring tools.
//
// Parameters:
//   handler: The handler function to wrap
//
// Returns:
//   http.HandlerFunc: Wrapped handler with CORS headers
func WithCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}
