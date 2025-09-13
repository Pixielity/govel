package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	webserver "govel/new/webserver/src"
	"govel/new/webserver/src/interfaces"
)

// RequestIDMiddleware generates and manages unique request identifiers for
// distributed tracing, logging correlation, and request tracking across services.
//
// Features:
//   - Multiple ID generation strategies (UUID, Nanoid, Custom)
//   - Configurable header names and formats
//   - Request ID propagation to response headers
//   - Thread-safe ID generation
//   - Custom ID validation and sanitization
//   - Integration with logging and monitoring systems
//   - Support for existing request IDs from upstream services
//   - Flexible ID length and character sets
//
// Configuration:
//   - header_name: HTTP header name for request ID (default: "X-Request-ID")
//   - response_header: Include ID in response headers (default: true)
//   - generator: ID generation strategy ("uuid", "nanoid", "timestamp", "custom")
//   - id_length: Length of generated IDs (for nanoid/custom)
//   - charset: Character set for ID generation
//   - prefix: Prefix to add to generated IDs
//   - validate_existing: Validate existing request IDs
//   - custom_generator: Custom ID generation function
type RequestIDMiddleware struct {
	webserver.BaseMiddleware
	HeaderName       string
	ResponseHeader   bool
	Generator        string
	IDLength         int
	Charset          string
	Prefix           string
	ValidateExisting bool
	CustomGenerator  func() string
	CustomValidator  func(string) bool
}

// NewRequestIDMiddleware creates a new request ID middleware with sensible defaults
func NewRequestIDMiddleware(config map[string]interface{}) *RequestIDMiddleware {
	m := &RequestIDMiddleware{
		HeaderName:       "X-Request-ID",
		ResponseHeader:   true,
		Generator:        "nanoid",
		IDLength:         21, // Standard nanoid length
		Charset:          "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		Prefix:           "",
		ValidateExisting: true,
	}

	// Apply configuration
	if headerName, ok := config["header_name"].(string); ok {
		m.HeaderName = headerName
	}
	if responseHeader, ok := config["response_header"].(bool); ok {
		m.ResponseHeader = responseHeader
	}
	if generator, ok := config["generator"].(string); ok {
		m.Generator = generator
	}
	if idLength, ok := config["id_length"].(int); ok {
		m.IDLength = idLength
	}
	if charset, ok := config["charset"].(string); ok {
		m.Charset = charset
	}
	if prefix, ok := config["prefix"].(string); ok {
		m.Prefix = prefix
	}
	if validate, ok := config["validate_existing"].(bool); ok {
		m.ValidateExisting = validate
	}
	if customGen, ok := config["custom_generator"].(func() string); ok {
		m.CustomGenerator = customGen
	}
	if customVal, ok := config["custom_validator"].(func(string) bool); ok {
		m.CustomValidator = customVal
	}

	return m
}

// Before generates or validates the request ID and adds it to the request context
func (m *RequestIDMiddleware) Before(req interfaces.RequestInterface) error {
	// Check if request already has an ID
	existingID := req.Header(m.HeaderName)
	
	var requestID string
	
	if existingID != "" {
		// Validate existing ID if validation is enabled
		if m.ValidateExisting && !m.isValidRequestID(existingID) {
			// Generate new ID if existing one is invalid
			requestID = m.generateRequestID()
		} else {
			// Use existing ID
			requestID = existingID
		}
	} else {
		// Generate new ID
		requestID = m.generateRequestID()
	}
	
	// Store request ID in context for use by other middleware and handlers
	req.SetContext("__request_id", requestID)
	
	// Set the request ID header (in case it was generated or sanitized)
	req.SetHeader(m.HeaderName, requestID)
	
	return nil
}

// After adds the request ID to the response headers if configured
func (m *RequestIDMiddleware) After(req interfaces.RequestInterface, res interfaces.ResponseInterface) error {
	if !m.ResponseHeader {
		return nil
	}
	
	// Get request ID from context
	if requestID := req.GetContext("__request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			// Add request ID to response headers
			res.Header(m.HeaderName, id)
			
			// Also add it as X-Trace-ID for compatibility
			if m.HeaderName != "X-Trace-ID" {
				res.Header("X-Trace-ID", id)
			}
		}
	}
	
	return nil
}

// generateRequestID generates a new request ID based on the configured strategy
func (m *RequestIDMiddleware) generateRequestID() string {
	var id string
	
	switch m.Generator {
	case "uuid":
		id = m.generateUUID()
	case "nanoid":
		id = m.generateNanoid()
	case "timestamp":
		id = m.generateTimestampID()
	case "custom":
		if m.CustomGenerator != nil {
			id = m.CustomGenerator()
		} else {
			id = m.generateNanoid() // Fallback
		}
	default:
		id = m.generateNanoid()
	}
	
	// Add prefix if configured
	if m.Prefix != "" {
		id = m.Prefix + id
	}
	
	return id
}

// generateUUID generates a UUID-like identifier
func (m *RequestIDMiddleware) generateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	
	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// generateNanoid generates a nanoid-style identifier
func (m *RequestIDMiddleware) generateNanoid() string {
	if m.IDLength <= 0 {
		m.IDLength = 21
	}
	
	bytes := make([]byte, m.IDLength)
	charsetLen := len(m.Charset)
	
	for i := 0; i < m.IDLength; i++ {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		bytes[i] = m.Charset[int(randomBytes[0])%charsetLen]
	}
	
	return string(bytes)
}

// generateTimestampID generates a timestamp-based identifier
func (m *RequestIDMiddleware) generateTimestampID() string {
	now := time.Now()
	timestamp := now.UnixNano()
	
	// Add some randomness to avoid collisions
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	
	return fmt.Sprintf("%d-%s", timestamp, randomHex)
}

// isValidRequestID validates an existing request ID
func (m *RequestIDMiddleware) isValidRequestID(id string) bool {
	// Use custom validator if provided
	if m.CustomValidator != nil {
		return m.CustomValidator(id)
	}
	
	// Basic validation rules
	if len(id) == 0 {
		return false
	}
	
	// Check length limits (reasonable bounds)
	if len(id) < 8 || len(id) > 128 {
		return false
	}
	
	// Check for dangerous characters that could cause issues in logs or headers
	dangerousChars := []string{"\r", "\n", "\t", "\x00", "<", ">", "\"", "'"}
	for _, char := range dangerousChars {
		if strings.Contains(id, char) {
			return false
		}
	}
	
	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range id {
		if !((char >= '0' && char <= '9') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			char == '-' || char == '_' || char == '.') {
			return false
		}
	}
	
	return true
}

// GetRequestID retrieves the request ID from the request context
// This is a utility function for other middleware and handlers
func (m *RequestIDMiddleware) GetRequestID(req interfaces.RequestInterface) string {
	if id := req.GetContext("__request_id"); id != nil {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}

// Name returns the middleware name
func (m *RequestIDMiddleware) Name() string {
	return "request_id"
}

// Priority returns the middleware priority (should be very early)
func (m *RequestIDMiddleware) Priority() int {
	return 3 // Execute very early for tracing
}

// Helper functions for common request ID configurations

// RequestIDWithDefaults creates a request ID middleware with standard settings
func RequestIDWithDefaults() *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator": "nanoid",
		"id_length": 21,
	})
}

// RequestIDWithUUID creates a request ID middleware using UUID format
func RequestIDWithUUID() *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator": "uuid",
	})
}

// RequestIDWithTimestamp creates a request ID middleware using timestamp format
func RequestIDWithTimestamp() *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator": "timestamp",
	})
}

// RequestIDWithPrefix creates a request ID middleware with a custom prefix
func RequestIDWithPrefix(prefix string) *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator": "nanoid",
		"prefix":    prefix,
		"id_length": 16, // Shorter since we have prefix
	})
}

// RequestIDForMicroservices creates a request ID middleware optimized for microservices
func RequestIDForMicroservices(serviceName string) *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator":         "timestamp",
		"prefix":            serviceName + "-",
		"response_header":   true,
		"validate_existing": true,
	})
}

// RequestIDCustom creates a request ID middleware with a custom generator
func RequestIDCustom(generator func() string) *RequestIDMiddleware {
	return NewRequestIDMiddleware(map[string]interface{}{
		"generator":        "custom",
		"custom_generator": generator,
	})
}

// Global helper function to extract request ID from any request
// This can be used by other middleware and handlers
func GetRequestIDFromContext(req interfaces.RequestInterface) string {
	if id := req.GetContext("__request_id"); id != nil {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	// Fallback to header if context is empty
	return req.Header("X-Request-ID")
}
