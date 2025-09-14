package parsers

import (
	"encoding/json"
	"fmt"
	"govel/package_manager/interfaces"
	"govel/package_manager/models"
	"os"
	"path/filepath"
	"strings"
)

// ModuleParser implements ParserInterface for parsing module.json files
type ModuleParser struct{}

// NewModuleParser creates a new module parser instance
func NewModuleParser() interfaces.ParserInterface {
	return &ModuleParser{}
}

// ParseFile parses a module.json file from the given file path
func (p *ModuleParser) ParseFile(filePath string) (*models.Package, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module.json file: %w", err)
	}

	pkg, err := p.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	// Set the package path based on the file location
	pkg.Path = filepath.Dir(filePath)

	return pkg, nil
}

// ParseBytes parses module.json data from byte slice
func (p *ModuleParser) ParseBytes(data []byte) (*models.Package, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	pkg := &models.Package{}

	// Parse basic fields
	if name, ok := rawData["name"].(string); ok {
		pkg.Name = name
	}

	if version, ok := rawData["version"].(string); ok {
		pkg.Version = version
	}

	if description, ok := rawData["description"].(string); ok {
		pkg.Description = description
	}

	if author, ok := rawData["author"].(string); ok {
		pkg.Author = author
	}

	if license, ok := rawData["license"].(string); ok {
		pkg.License = license
	}

	if homepage, ok := rawData["homepage"].(string); ok {
		pkg.Homepage = homepage
	}

	// Parse keywords
	if keywords, ok := rawData["keywords"].([]interface{}); ok {
		for _, keyword := range keywords {
			if kw, ok := keyword.(string); ok {
				pkg.Keywords = append(pkg.Keywords, kw)
			}
		}
	}

	// Parse repository
	if repo, ok := rawData["repository"].(map[string]interface{}); ok {
		pkg.Repository = models.Repository{
			Type:      getStringValue(repo, "type"),
			URL:       getStringValue(repo, "url"),
			Directory: getStringValue(repo, "directory"),
		}
	}

	// Parse bugs
	if bugs, ok := rawData["bugs"].(map[string]interface{}); ok {
		pkg.Bugs = models.Bugs{
			URL: getStringValue(bugs, "url"),
		}
	}

	// Parse dependencies
	if deps, ok := rawData["dependencies"].(map[string]interface{}); ok {
		pkg.Dependencies = make(map[string]string)
		for key, value := range deps {
			if v, ok := value.(string); ok {
				pkg.Dependencies[key] = v
			}
		}
	}

	// Parse scripts
	if scripts, ok := rawData["scripts"].(map[string]interface{}); ok {
		pkg.Scripts = make(map[string]string)
		for key, value := range scripts {
			if v, ok := value.(string); ok {
				pkg.Scripts[key] = v
			}
		}
	}

	// Parse hooks
	if hooks, ok := rawData["hooks"].(map[string]interface{}); ok {
		pkg.Hooks = models.Hooks{
			PreInstall:  parseStringArray(hooks, "pre-install"),
			PostInstall: parseStringArray(hooks, "post-install"),
			PreUpdate:   parseStringArray(hooks, "pre-update"),
			PostUpdate:  parseStringArray(hooks, "post-update"),
			PreBuild:    parseStringArray(hooks, "pre-build"),
			PostBuild:   parseStringArray(hooks, "post-build"),
			PreTest:     parseStringArray(hooks, "pre-test"),
			PostTest:    parseStringArray(hooks, "post-test"),
			PrePublish:  parseStringArray(hooks, "pre-publish"),
			PostPublish: parseStringArray(hooks, "post-publish"),
		}
	}

	// Parse engines
	if engines, ok := rawData["engines"].(map[string]interface{}); ok {
		pkg.Engines = make(map[string]string)
		for key, value := range engines {
			if v, ok := value.(string); ok {
				pkg.Engines[key] = v
			}
		}
	}

	// Parse files
	if files, ok := rawData["files"].([]interface{}); ok {
		for _, file := range files {
			if f, ok := file.(string); ok {
				pkg.Files = append(pkg.Files, f)
			}
		}
	}

	// Parse GoVel config
	if govel, ok := rawData["govel"].(map[string]interface{}); ok {
		pkg.GovelConfig = models.GovelConfig{
			Type:       getStringValue(govel, "type"),
			Category:   getStringValue(govel, "category"),
			Providers:  parseStringArray(govel, "providers"),
			Facades:    parseStringArray(govel, "facades"),
			Middleware: parseStringArray(govel, "middleware"),
			Commands:   parseStringArray(govel, "commands"),
			Provides:   parseStringArray(govel, "provides"),
		}
	}

	return pkg, nil
}

// ValidatePackage validates a parsed package for required fields and consistency
func (p *ModuleParser) ValidatePackage(pkg *models.Package) error {
	if pkg == nil {
		return fmt.Errorf("package cannot be nil")
	}

	// Validate required fields
	if strings.TrimSpace(pkg.Name) == "" {
		return fmt.Errorf("package name is required")
	}

	if strings.TrimSpace(pkg.Version) == "" {
		return fmt.Errorf("package version is required")
	}

	if strings.TrimSpace(pkg.Description) == "" {
		return fmt.Errorf("package description is required")
	}

	// Validate name format (should start with @govel/)
	if !strings.HasPrefix(pkg.Name, "@govel/") {
		return fmt.Errorf("package name must start with '@govel/'")
	}

	// Validate version format (basic semver check)
	versionParts := strings.Split(pkg.Version, ".")
	if len(versionParts) != 3 {
		return fmt.Errorf("version must be in semver format (x.y.z)")
	}

	// Validate GoVel config
	if pkg.GovelConfig.Type == "" {
		return fmt.Errorf("govel.type is required")
	}

	validTypes := []string{"package", "application", "plugin"}
	if !contains(validTypes, pkg.GovelConfig.Type) {
		return fmt.Errorf("govel.type must be one of: %s", strings.Join(validTypes, ", "))
	}

	if pkg.GovelConfig.Category == "" {
		return fmt.Errorf("govel.category is required")
	}

	return nil
}

// WriteFile writes a package to a module.json file
func (p *ModuleParser) WriteFile(pkg *models.Package, filePath string) error {
	if err := p.ValidatePackage(pkg); err != nil {
		return fmt.Errorf("package validation failed: %w", err)
	}

	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package to JSON: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write module.json file: %w", err)
	}

	return nil
}

// Helper functions

func getStringValue(m map[string]interface{}, key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func parseStringArray(m map[string]interface{}, key string) []string {
	var result []string
	if arr, ok := m[key].([]interface{}); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
