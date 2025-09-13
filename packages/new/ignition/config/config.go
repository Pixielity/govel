package config

import (
	"govel/packages/ignition/enums"
)

// Config holds the Ignition configuration
type Config struct {
	Theme       enums.Theme  `json:"theme"`
	Editor      enums.Editor `json:"editor"`
	ShareReport bool         `json:"shareReport"`
}

// NewConfig creates a new configuration with sensible defaults
func NewConfig() *Config {
	return &Config{
		Theme:       enums.ThemeAuto,
		Editor:      enums.EditorVSCode,
		ShareReport: false,
	}
}

// WithTheme sets the theme
func (c *Config) WithTheme(theme enums.Theme) *Config {
	c.Theme = theme
	return c
}

// WithEditor sets the editor
func (c *Config) WithEditor(editor enums.Editor) *Config {
	c.Editor = editor
	return c
}

// WithShareReport enables or disables report sharing
func (c *Config) WithShareReport(share bool) *Config {
	c.ShareReport = share
	return c
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if !c.Theme.IsValid() {
		c.Theme = enums.ThemeAuto
	}

	if !c.Editor.IsValid() {
		c.Editor = enums.EditorVSCode
	}

	return nil
}
