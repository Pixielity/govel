package views

import "embed"

//go:embed assets templates
var FS embed.FS

// GetAsset reads an asset file from the embedded filesystem
func GetAsset(filename string) ([]byte, error) {
	return FS.ReadFile(filename)
}

// GetAssetString reads an asset file from the embedded filesystem and returns it as a string
func GetAssetString(filename string) string {
	content, err := FS.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(content)
}
