// Package helpers provides utility functions for webserver configuration and validation.
// This file contains validation helpers extracted from the builder pattern.
package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"govel/packages/new/webserver/src/enums"
)

// ValidationHelper provides utilities for webserver configuration validation.
type ValidationHelper struct{}

// NewValidationHelper creates a new validation helper instance.
//
// Returns:
//
//	*ValidationHelper: A new validation helper
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{}
}

// ValidateConfiguration performs comprehensive validation on webserver configuration.
//
// Parameters:
//
//	config: The configuration map to validate
//	engine: The selected engine (optional, can be nil if custom adapter is used)
//
// Returns:
//
//	error: Validation error, or nil if configuration is valid
func (v *ValidationHelper) ValidateConfiguration(config map[string]interface{}, engine *enums.Engine) error {
	// Validate engine selection (if provided)
	if engine != nil && !engine.IsValid() {
		return fmt.Errorf("invalid engine: %v", engine)
	}

	// Basic network configuration validation
	if err := v.ValidatePort(config); err != nil {
		return err
	}
	if err := v.ValidateHost(config); err != nil {
		return err
	}
	if err := v.ValidateAddress(config); err != nil {
		return err
	}
	if err := v.ValidateTimeout(config); err != nil {
		return err
	}
	if err := v.ValidateBodySize(config); err != nil {
		return err
	}

	// Security and authentication validation
	if err := v.ValidateTLSConfig(config); err != nil {
		return err
	}
	if err := v.ValidateRateLimit(config); err != nil {
		return err
	}
	if err := v.ValidateAPIKeys(config); err != nil {
		return err
	}

	// Performance and optimization validation
	if err := v.ValidateConcurrency(config); err != nil {
		return err
	}

	// Static content and template validation
	if err := v.ValidateStaticFiles(config); err != nil {
		return err
	}
	if err := v.ValidateTemplateConfig(config); err != nil {
		return err
	}

	return nil
}

// ValidatePort validates port configuration.
//
// Parameters:
//
//	config: The configuration map containing port settings
//
// Returns:
//
//	error: Validation error, or nil if port is valid
func (v *ValidationHelper) ValidatePort(config map[string]interface{}) error {
	if portInterface, exists := config["port"]; exists {
		if port, ok := portInterface.(int); ok {
			if port < 1 || port > 65535 {
				return fmt.Errorf("invalid port number: %d (must be 1-65535)", port)
			}
		} else {
			return fmt.Errorf("port must be an integer, got: %T", portInterface)
		}
	}
	return nil
}

// ValidateHost validates host configuration.
//
// Parameters:
//
//	config: The configuration map containing host settings
//
// Returns:
//
//	error: Validation error, or nil if host is valid
func (v *ValidationHelper) ValidateHost(config map[string]interface{}) error {
	if hostInterface, exists := config["host"]; exists {
		if host, ok := hostInterface.(string); ok {
			if strings.TrimSpace(host) == "" {
				return fmt.Errorf("host cannot be empty")
			}
		} else {
			return fmt.Errorf("host must be a string, got: %T", hostInterface)
		}
	}
	return nil
}

// ValidateAddress validates address format configuration.
//
// Parameters:
//
//	config: The configuration map containing address settings
//
// Returns:
//
//	error: Validation error, or nil if address is valid
func (v *ValidationHelper) ValidateAddress(config map[string]interface{}) error {
	if addrInterface, exists := config["address"]; exists {
		if addr, ok := addrInterface.(string); ok {
			if err := v.ValidateAddressFormat(addr); err != nil {
				return fmt.Errorf("invalid address format: %w", err)
			}
		} else {
			return fmt.Errorf("address must be a string, got: %T", addrInterface)
		}
	}
	return nil
}

// ValidateAddressFormat validates the format of an address string.
//
// Parameters:
//
//	addr: The address string to validate
//
// Returns:
//
//	error: Validation error, or nil if address format is valid
func (v *ValidationHelper) ValidateAddressFormat(addr string) error {
	if addr == "" {
		return fmt.Errorf("address cannot be empty")
	}

	// Must contain exactly one colon for host:port format
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("address must be in format 'host:port', got: %s", addr)
	}

	// Validate port part if present
	portStr := strings.TrimSpace(parts[1])
	if portStr != "" {
		if port, err := strconv.Atoi(portStr); err != nil {
			return fmt.Errorf("invalid port in address: %s", portStr)
		} else if port < 1 || port > 65535 {
			return fmt.Errorf("port in address out of range: %d", port)
		}
	}

	return nil
}

// ValidateTimeout validates timeout configuration.
//
// Parameters:
//
//	config: The configuration map containing timeout settings
//
// Returns:
//
//	error: Validation error, or nil if timeout is valid
func (v *ValidationHelper) ValidateTimeout(config map[string]interface{}) error {
	if timeoutInterface, exists := config["timeout"]; exists {
		if timeout, ok := timeoutInterface.(int); ok {
			if timeout < 0 {
				return fmt.Errorf("timeout cannot be negative: %d", timeout)
			}
			if timeout > 3600 { // 1 hour maximum
				return fmt.Errorf("timeout too large: %d seconds (maximum 3600)", timeout)
			}
		} else {
			return fmt.Errorf("timeout must be an integer, got: %T", timeoutInterface)
		}
	}
	return nil
}

// ValidateBodySize validates max body size configuration.
//
// Parameters:
//
//	config: The configuration map containing body size settings
//
// Returns:
//
//	error: Validation error, or nil if body size is valid
func (v *ValidationHelper) ValidateBodySize(config map[string]interface{}) error {
	if bodySizeInterface, exists := config["max_body_size"]; exists {
		if bodySize, ok := bodySizeInterface.(int); ok {
			if bodySize < 0 {
				return fmt.Errorf("max_body_size cannot be negative: %d", bodySize)
			}
			// 100MB maximum
			if bodySize > 100*1024*1024 {
				return fmt.Errorf("max_body_size too large: %d bytes (maximum 100MB)", bodySize)
			}
		} else {
			return fmt.Errorf("max_body_size must be an integer, got: %T", bodySizeInterface)
		}
	}
	return nil
}

// ValidateEngine validates engine selection.
//
// Parameters:
//
//	engine: The engine to validate
//
// Returns:
//
//	error: Validation error, or nil if engine is valid
func (v *ValidationHelper) ValidateEngine(engine enums.Engine) error {
	if !engine.IsValid() {
		return fmt.Errorf("invalid engine: %v", engine)
	}
	return nil
}

// IsValidHTTPMethod checks if a method string represents a valid HTTP method.
//
// Parameters:
//
//	method: The HTTP method string to validate
//
// Returns:
//
//	bool: True if the method is valid, false otherwise
func (v *ValidationHelper) IsValidHTTPMethod(method string) bool {
	_, valid := enums.ParseHTTPMethod(method)
	return valid
}

// NormalizeHTTPMethod normalizes an HTTP method string to uppercase.
//
// Parameters:
//
//	method: The HTTP method string to normalize
//
// Returns:
//
//	string: The normalized method string (uppercase)
//	error: Error if the method is invalid
func (v *ValidationHelper) NormalizeHTTPMethod(method string) (string, error) {
	httpMethod, valid := enums.ParseHTTPMethod(method)
	if !valid {
		return "", fmt.Errorf("invalid HTTP method: %s", method)
	}
	return httpMethod.String(), nil
}

// ValidateMiddleware validates middleware configuration.
//
// Parameters:
//
//	middleware: The middleware slice to validate
//
// Returns:
//
//	error: Validation error, or nil if middleware is valid
func (v *ValidationHelper) ValidateMiddleware(middleware []interface{}) error {
	if middleware == nil {
		return nil
	}

	for i, mw := range middleware {
		if mw == nil {
			return fmt.Errorf("middleware at index %d is nil", i)
		}
	}

	return nil
}

// ValidateRequiredConfig validates that required configuration keys are present.
//
// Parameters:
//
//	config: The configuration map to validate
//	requiredKeys: Slice of required configuration keys
//
// Returns:
//
//	error: Validation error, or nil if all required keys are present
func (v *ValidationHelper) ValidateRequiredConfig(config map[string]interface{}, requiredKeys []string) error {
	for _, key := range requiredKeys {
		if _, exists := config[key]; !exists {
			return fmt.Errorf("required configuration key missing: %s", key)
		}
	}
	return nil
}

// ValidateTLSConfig validates TLS/SSL certificate configuration.
//
// Parameters:
//
//	config: The configuration map containing TLS settings
//
// Returns:
//
//	error: Validation error, or nil if TLS configuration is valid
func (v *ValidationHelper) ValidateTLSConfig(config map[string]interface{}) error {
	if enabled, exists := config["tls_enabled"]; exists && enabled == true {
		certFile, certExists := config["tls_cert_file"]
		keyFile, keyExists := config["tls_key_file"]
		
		if !certExists || certFile == "" {
			return fmt.Errorf("TLS certificate file is required when TLS is enabled")
		}
		if !keyExists || keyFile == "" {
			return fmt.Errorf("TLS key file is required when TLS is enabled")
		}
		
		// Validate that cert and key files are strings
		if _, ok := certFile.(string); !ok {
			return fmt.Errorf("TLS certificate file must be a string")
		}
		if _, ok := keyFile.(string); !ok {
			return fmt.Errorf("TLS key file must be a string")
		}
	}
	return nil
}

// ValidateRateLimit validates rate limiting configuration.
//
// Parameters:
//
//	config: The configuration map containing rate limit settings
//
// Returns:
//
//	error: Validation error, or nil if rate limit configuration is valid
func (v *ValidationHelper) ValidateRateLimit(config map[string]interface{}) error {
	if enabled, exists := config["rate_limit_enabled"]; exists && enabled == true {
		requests, requestsExists := config["rate_limit_requests"]
		window, windowExists := config["rate_limit_window"]
		
		if !requestsExists {
			return fmt.Errorf("rate_limit_requests is required when rate limiting is enabled")
		}
		if !windowExists {
			return fmt.Errorf("rate_limit_window is required when rate limiting is enabled")
		}
		
		// Validate requests is a positive integer
		if reqInt, ok := requests.(int); ok {
			if reqInt <= 0 {
				return fmt.Errorf("rate_limit_requests must be a positive integer, got: %d", reqInt)
			}
		} else {
			return fmt.Errorf("rate_limit_requests must be an integer, got: %T", requests)
		}
		
		// Validate window is a duration
		if _, ok := window.(time.Duration); !ok {
			return fmt.Errorf("rate_limit_window must be a time.Duration, got: %T", window)
		}
	}
	return nil
}

// ValidateConcurrency validates concurrency limit configuration.
//
// Parameters:
//
//	config: The configuration map containing concurrency settings
//
// Returns:
//
//	error: Validation error, or nil if concurrency configuration is valid
func (v *ValidationHelper) ValidateConcurrency(config map[string]interface{}) error {
	if limit, exists := config["max_concurrent_connections"]; exists {
		if limitInt, ok := limit.(int); ok {
			if limitInt <= 0 {
				return fmt.Errorf("max_concurrent_connections must be a positive integer, got: %d", limitInt)
			}
			// Reasonable upper limit to prevent resource exhaustion
			if limitInt > 1000000 {
				return fmt.Errorf("max_concurrent_connections too large: %d (maximum 1,000,000)", limitInt)
			}
		} else {
			return fmt.Errorf("max_concurrent_connections must be an integer, got: %T", limit)
		}
	}
	return nil
}

// ValidateAPIKeys validates API key configuration.
//
// Parameters:
//
//	config: The configuration map containing API key settings
//
// Returns:
//
//	error: Validation error, or nil if API key configuration is valid
func (v *ValidationHelper) ValidateAPIKeys(config map[string]interface{}) error {
	if enabled, exists := config["api_key_auth_enabled"]; exists && enabled == true {
		keys, keysExists := config["api_keys"]
		
		if !keysExists {
			return fmt.Errorf("api_keys is required when API key authentication is enabled")
		}
		
		// Validate that keys is a slice of strings
		if keySlice, ok := keys.([]string); ok {
			if len(keySlice) == 0 {
				return fmt.Errorf("api_keys cannot be empty when API key authentication is enabled")
			}
			
			// Validate each key is non-empty
			for i, key := range keySlice {
				if strings.TrimSpace(key) == "" {
					return fmt.Errorf("API key at index %d cannot be empty", i)
				}
			}
		} else {
			return fmt.Errorf("api_keys must be a slice of strings, got: %T", keys)
		}
	}
	return nil
}

// ValidateStaticFiles validates static file serving configuration.
//
// Parameters:
//
//	config: The configuration map containing static file settings
//
// Returns:
//
//	error: Validation error, or nil if static file configuration is valid
func (v *ValidationHelper) ValidateStaticFiles(config map[string]interface{}) error {
	if enabled, exists := config["static_files_enabled"]; exists && enabled == true {
		staticFiles, filesExists := config["static_files"]
		
		if filesExists {
			if filesMap, ok := staticFiles.(map[string]string); ok {
				for prefix, directory := range filesMap {
					if strings.TrimSpace(prefix) == "" {
						return fmt.Errorf("static file prefix cannot be empty")
					}
					if strings.TrimSpace(directory) == "" {
						return fmt.Errorf("static file directory cannot be empty for prefix: %s", prefix)
					}
				}
			} else {
				return fmt.Errorf("static_files must be a map[string]string, got: %T", staticFiles)
			}
		}
	}
	return nil
}

// ValidateTemplateConfig validates template engine configuration.
//
// Parameters:
//
//	config: The configuration map containing template settings
//
// Returns:
//
//	error: Validation error, or nil if template configuration is valid
func (v *ValidationHelper) ValidateTemplateConfig(config map[string]interface{}) error {
	if enabled, exists := config["template_enabled"]; exists && enabled == true {
		engine, engineExists := config["template_engine"]
		directory, dirExists := config["template_directory"]
		
		if !engineExists {
			return fmt.Errorf("template_engine is required when templates are enabled")
		}
		if !dirExists {
			return fmt.Errorf("template_directory is required when templates are enabled")
		}
		
		// Validate template engine
		if templateEngine, ok := engine.(enums.TemplateEngine); ok {
			if !templateEngine.IsValid() {
				return fmt.Errorf("invalid template engine: %v", templateEngine)
			}
		} else {
			return fmt.Errorf("template_engine must be a TemplateEngine enum, got: %T", engine)
		}
		
		// Validate directory is a non-empty string
		if dirStr, ok := directory.(string); ok {
			if strings.TrimSpace(dirStr) == "" {
				return fmt.Errorf("template_directory cannot be empty")
			}
		} else {
			return fmt.Errorf("template_directory must be a string, got: %T", directory)
		}
	}
	return nil
}
