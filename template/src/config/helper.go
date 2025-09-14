package config

import (
	"os"
	"strconv"
	"strings"
)

// Env retrieves an environment variable with an optional default value.
// This is a helper function similar to Laravel's env() helper.
// It provides automatic type conversion based on the default value type.
func Env(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Try to parse common types based on the default value type
	switch defaultValue.(type) {
	case bool:
		if lower := strings.ToLower(value); lower == "true" || lower == "1" || lower == "yes" || lower == "on" {
			return true
		}
		if lower := strings.ToLower(value); lower == "false" || lower == "0" || lower == "no" || lower == "off" {
			return false
		}
		return defaultValue
	case int:
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		return defaultValue
	case int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
		return defaultValue
	case float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
		return defaultValue
	default:
		return value
	}
}

// StoragePath returns the path to the storage directory
// Similar to Laravel's storage_path() helper function
func StoragePath(path string) string {
	baseStoragePath := Env("STORAGE_PATH", "storage").(string)
	if path == "" {
		return baseStoragePath
	}
	return baseStoragePath + "/" + path
}

// Slug converts a string to a URL-friendly slug
// Similar to Laravel's Str::slug() helper
func Slug(str string) string {
	// Basic slug implementation - convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(str)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

// PublicPath returns the path to the public directory
// Similar to Laravel's public_path() helper function
func PublicPath(path string) string {
	basePublicPath := Env("PUBLIC_PATH", "public").(string)
	if path == "" {
		return basePublicPath
	}
	return basePublicPath + "/" + path
}
