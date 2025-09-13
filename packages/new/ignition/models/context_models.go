package models

import (
	"net/http"
	"os"
	"runtime"
	"time"
)

// RequestContext represents HTTP request information
type RequestContext struct {
	URL       string `json:"url"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	UserAgent string `json:"useragent"`
}

// RequestData represents request payload data
type RequestData struct {
	QueryString []interface{} `json:"queryString"`
	Body        []interface{} `json:"body"`
	Files       []interface{} `json:"files"`
}

// RouteContext represents routing information
type RouteContext struct {
	Route            interface{}   `json:"route"`
	RouteParameters  []interface{} `json:"routeParameters"`
	ControllerAction string        `json:"controllerAction"`
	Middleware       []string      `json:"middleware"`
}

// EnvironmentContext represents runtime environment
type EnvironmentContext struct {
	GoVersion       string `json:"go_version"`
	GovelVersion    string `json:"govel_version"`
	GovelLocale     string `json:"govel_locale"`
	GovelConfigCached bool `json:"govel_config_cached"`
	AppDebug        bool   `json:"app_debug"`
	AppEnv          string `json:"app_env"`
}

// QueryInfo represents database query information
type QueryInfo struct {
	SQL            string        `json:"sql"`
	Time           float64       `json:"time"`
	ConnectionName string        `json:"connection_name"`
	Bindings       []interface{} `json:"bindings"`
	Microtime      float64       `json:"microtime"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Context   map[string]interface{} `json:"context"`
	Timestamp time.Time `json:"timestamp"`
}

// DumpEntry represents a dump/debug output
type DumpEntry struct {
	Content   string `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// CompleteErrorContext represents the full error context
type CompleteErrorContext struct {
	Request     RequestContext     `json:"request"`
	RequestData RequestData        `json:"request_data"`
	Headers     map[string]string  `json:"headers"`
	Cookies     map[string]string  `json:"cookies"`
	Session     map[string]interface{} `json:"session"`
	Route       RouteContext       `json:"route"`
	Env         EnvironmentContext `json:"env"`
	Dumps       []DumpEntry        `json:"dumps"`
	Logs        []LogEntry         `json:"logs"`
	Queries     []QueryInfo        `json:"queries"`
}

// NewCompleteErrorContext creates a complete error context from HTTP request
func NewCompleteErrorContext(r *http.Request) *CompleteErrorContext {
	ctx := &CompleteErrorContext{
		Request: RequestContext{
			URL:       getRequestURL(r),
			IP:        getClientIP(r),
			Method:    getMethod(r),
			UserAgent: getUserAgent(r),
		},
		RequestData: RequestData{
			QueryString: []interface{}{},
			Body:        []interface{}{},
			Files:       []interface{}{},
		},
		Headers:     getHeaders(r),
		Cookies:     getCookies(r),
		Session:     getSession(r),
		Route:       getRoute(r),
		Env:         getEnvironment(),
		Dumps:       []DumpEntry{},
		Logs:        []LogEntry{},
		Queries:     []QueryInfo{},
	}

	return ctx
}

// Helper functions for context data extraction

func getRequestURL(r *http.Request) string {
	if r == nil {
		return ""
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + r.RequestURI
}

func getClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	return r.RemoteAddr
}

func getMethod(r *http.Request) string {
	if r == nil {
		return "UNKNOWN"
	}
	return r.Method
}

func getUserAgent(r *http.Request) string {
	if r == nil {
		return ""
	}
	return r.Header.Get("User-Agent")
}

// Removed getQueryString and getRequestBody - now using arrays directly

func getHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	if r == nil {
		return headers
	}
	
	for name, values := range r.Header {
		if len(values) > 0 {
			// Censor sensitive headers
			if isSensitiveHeader(name) {
				headers[name] = "<CENSORED>"
			} else {
				headers[name] = values[0]
			}
		}
	}
	
	return headers
}

func getCookies(r *http.Request) map[string]string {
	cookies := make(map[string]string)
	if r == nil {
		return cookies
	}
	
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	
	return cookies
}

func getSession(r *http.Request) map[string]interface{} {
	// For now, return empty session. In a real implementation,
	// you'd integrate with your session management system
	return map[string]interface{}{
		"_token": "go-session-token",
		"_flash": map[string]interface{}{
			"old": []interface{}{},
			"new": []interface{}{},
		},
	}
}

func getRoute(r *http.Request) RouteContext {
	return RouteContext{
		Route:            nil,
		RouteParameters:  []interface{}{},
		ControllerAction: "Go Handler",
		Middleware:       []string{"web"},
	}
}

func getEnvironment() EnvironmentContext {
	return EnvironmentContext{
		GoVersion:         runtime.Version(),
		GovelVersion:      "1.0.0",
		GovelLocale:       "en",
		GovelConfigCached: false,
		AppDebug:          true,
		AppEnv:            getAppEnv(),
	}
}

func getAppEnv() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	if env := os.Getenv("GO_ENV"); env != "" {
		return env
	}
	return "development"
}

func isSensitiveHeader(name string) bool {
	sensitive := []string{"Cookie", "Authorization", "X-Api-Key", "X-Auth-Token"}
	for _, s := range sensitive {
		if name == s {
			return true
		}
	}
	return false
}
