package utils

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"govel/packages/ignition/config"
	"govel/packages/ignition/enums"
)

// SourceCodeExtractor handles extracting source code around error lines
type SourceCodeExtractor struct {
	contextLines int // Number of lines to show before and after the error line
}

// NewSourceCodeExtractor creates a new source code extractor
func NewSourceCodeExtractor(contextLines int) *SourceCodeExtractor {
	if contextLines < 1 {
		contextLines = 5 // Default to 5 lines of context
	}
	return &SourceCodeExtractor{
		contextLines: contextLines,
	}
}

// ExtractSourceCode extracts source code around the specified line
func (e *SourceCodeExtractor) ExtractSourceCode(filename string, lineNum int) map[string]string {
	code := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		// If we can't read the file, return a placeholder
		code[strconv.Itoa(lineNum)] = fmt.Sprintf("// Could not read source file: %s", filepath.Base(filename))
		return code
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	currentLine := 1

	// Read all lines
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		code[strconv.Itoa(lineNum)] = fmt.Sprintf("// Error reading file: %v", err)
		return code
	}

	// Calculate the range of lines to extract
	start := max(1, lineNum-e.contextLines)
	end := min(len(lines), lineNum+e.contextLines)

	// Extract the lines with line numbers
	for i := start; i <= end; i++ {
		if i <= len(lines) {
			lineContent := lines[i-1] // lines array is 0-indexed
			code[strconv.Itoa(i)] = lineContent
		}
	}

	return code
}

// ErrorAnalyzer provides advanced error analysis capabilities
type ErrorAnalyzer struct{}

// NewErrorAnalyzer creates a new error analyzer
func NewErrorAnalyzer() *ErrorAnalyzer {
	return &ErrorAnalyzer{}
}

// AnalyzeError performs detailed analysis of an error
func (a *ErrorAnalyzer) AnalyzeError(err error) map[string]interface{} {
	analysis := make(map[string]interface{})
	errMsg := err.Error()

	// Basic error categorization
	analysis["category"] = a.categorizeError(errMsg)
	analysis["severity"] = a.determineSeverity(errMsg)
	analysis["tags"] = a.extractTags(errMsg)

	return analysis
}

// categorizeError categorizes the error into common types
func (a *ErrorAnalyzer) categorizeError(errMsg string) string {
	errMsg = strings.ToLower(errMsg)

	categories := map[string][]string{
		"nil_pointer": {"nil pointer", "invalid memory address"},
		"type_error":  {"cannot use", "type assertion", "interface conversion"},
		"import":      {"cannot find package", "imported and not used"},
		"syntax":      {"syntax error", "unexpected"},
		"network":     {"connection refused", "timeout", "no such host"},
		"file":        {"no such file", "permission denied", "file not found"},
		"concurrency": {"concurrent map", "send on closed channel", "deadlock"},
		"database":    {"no rows", "connection refused", "sql"},
		"json":        {"json:", "unmarshal", "marshal"},
		"http":        {"http:", "status code", "request"},
	}

	for category, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(errMsg, keyword) {
				return category
			}
		}
	}

	return "unknown"
}

// determineSeverity determines the severity level of an error
func (a *ErrorAnalyzer) determineSeverity(errMsg string) string {
	errMsg = strings.ToLower(errMsg)

	critical := []string{"panic", "fatal", "segmentation fault", "stack overflow"}
	high := []string{"nil pointer", "index out of range", "deadlock"}
	medium := []string{"connection refused", "timeout", "file not found"}

	for _, keyword := range critical {
		if strings.Contains(errMsg, keyword) {
			return "critical"
		}
	}

	for _, keyword := range high {
		if strings.Contains(errMsg, keyword) {
			return "high"
		}
	}

	for _, keyword := range medium {
		if strings.Contains(errMsg, keyword) {
			return "medium"
		}
	}

	return "low"
}

// extractTags extracts relevant tags from the error message
func (a *ErrorAnalyzer) extractTags(errMsg string) []string {
	var tags []string
	errMsg = strings.ToLower(errMsg)

	tagMap := map[string]string{
		"nil":        "nil-pointer",
		"pointer":    "pointer",
		"json":       "json",
		"http":       "http",
		"sql":        "database",
		"timeout":    "timeout",
		"connection": "network",
		"concurrent": "concurrency",
		"goroutine":  "concurrency",
		"channel":    "channel",
		"mutex":      "synchronization",
		"interface":  "interface",
		"type":       "type-system",
		"import":     "import",
		"package":    "package",
		"syntax":     "syntax",
		"parse":      "parsing",
	}

	for keyword, tag := range tagMap {
		if strings.Contains(errMsg, keyword) {
			tags = append(tags, tag)
		}
	}

	return tags
}

// FilePathResolver resolves relative paths to absolute paths
type FilePathResolver struct {
	applicationPath string
}

// NewFilePathResolver creates a new file path resolver
func NewFilePathResolver(applicationPath string) *FilePathResolver {
	return &FilePathResolver{
		applicationPath: applicationPath,
	}
}

// ResolvePath resolves a path to an absolute path relative to the application
func (r *FilePathResolver) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if r.applicationPath != "" {
		resolved := filepath.Join(r.applicationPath, path)
		if _, err := os.Stat(resolved); err == nil {
			return resolved
		}
	}

	// If we can't resolve relative to application path, try current directory
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}

	return path
}

// RelativePath returns a path relative to the application root
func (r *FilePathResolver) RelativePath(path string) string {
	if r.applicationPath == "" {
		return path
	}

	if rel, err := filepath.Rel(r.applicationPath, path); err == nil {
		return rel
	}

	return path
}

// ThemeManager manages theme-related functionality
type ThemeManager struct {
	supportedThemes map[string]bool
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		supportedThemes: map[string]bool{
			"auto":  true,
			"light": true,
			"dark":  true,
		},
	}
}

// ValidateTheme validates if a theme is supported
func (t *ThemeManager) ValidateTheme(theme string) bool {
	return t.supportedThemes[theme]
}

// GetDefaultTheme returns the default theme
func (t *ThemeManager) GetDefaultTheme() string {
	return "auto"
}

// ConfigManager handles configuration-related operations
type ConfigManager struct{}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

// LoadConfig loads configuration from various sources
func (c *ConfigManager) LoadConfig() *config.Config {
	config := &config.Config{
		Theme:       enums.ThemeAuto,
		Editor:      enums.EditorVSCode,
		ShareReport: false,
	}

	// Try to load from environment variables
	if theme := os.Getenv("IGNITION_THEME"); theme != "" {
		config.Theme = enums.Theme(theme)
	}

	if editor := os.Getenv("IGNITION_EDITOR"); editor != "" {
		config.Editor = enums.Editor(editor)
	}

	// Add more configuration sources as needed (files, etc.)

	return config
}

// Utility functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// IsReadableFile checks if a file exists and is readable
func IsReadableFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// SanitizePath sanitizes a file path for display
func SanitizePath(path string) string {
	// Remove any potentially sensitive information
	path = filepath.Clean(path)

	// You might want to add more sanitization logic here
	// For example, replacing user home directory with ~
	if homeDir, err := os.UserHomeDir(); err == nil {
		if strings.HasPrefix(path, homeDir) {
			path = strings.Replace(path, homeDir, "~", 1)
		}
	}

	return path
}

// WalkFiles recursively walks through files in a directory
func WalkFiles(root string, walkFunc func(path string, info fs.FileInfo, err error) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return walkFunc(path, nil, err)
		}
		info, err := d.Info()
		if err != nil {
			return walkFunc(path, nil, err)
		}
		return walkFunc(path, info, nil)
	})
}
