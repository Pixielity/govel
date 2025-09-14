package models

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// CompleteStackFrame represents a complete stack frame with all required fields
type CompleteStackFrame struct {
	File             string            `json:"file"`
	LineNumber       int               `json:"line_number"`
	Method           string            `json:"method"`
	Class            interface{}       `json:"class"`
	CodeSnippet      map[string]string `json:"code_snippet"`
	Arguments        []interface{}     `json:"arguments"`
	ApplicationFrame bool              `json:"application_frame"`
}

// CompleteReport represents the full error report structure
type CompleteReport struct {
	Notifier           string               `json:"notifier"`
	Language           string               `json:"language"`
	FrameworkVersion   string               `json:"framework_version"`
	LanguageVersion    string               `json:"language_version"`
	ExceptionClass     string               `json:"exception_class"`
	SeenAt             int64                `json:"seen_at"`
	Message            string               `json:"message"`
	Glows              []interface{}        `json:"glows"`
	Solutions          []Solution           `json:"solutions"`
	DocumentationLinks []string             `json:"documentation_links"`
	Stacktrace         []CompleteStackFrame `json:"stacktrace"`
	Context            CompleteErrorContext `json:"context"`
	Stage              string               `json:"stage"`
	MessageLevel       interface{}          `json:"message_level"`
	OpenFrameIndex     interface{}          `json:"open_frame_index"`
	ApplicationPath    string               `json:"application_path"`
	ApplicationVersion interface{}          `json:"application_version"`
	TrackingUuid       string               `json:"tracking_uuid"`
	Handled            interface{}          `json:"handled"`
	OverriddenGrouping interface{}          `json:"overridden_grouping"`
}

// ShareableReport represents the shareable version of the report (may have filtered data)
type ShareableReport struct {
	Notifier           string               `json:"notifier"`
	Language           string               `json:"language"`
	FrameworkVersion   string               `json:"framework_version"`
	LanguageVersion    string               `json:"language_version"`
	ExceptionClass     string               `json:"exception_class"`
	SeenAt             int64                `json:"seen_at"`
	Message            string               `json:"message"`
	Glows              []interface{}        `json:"glows"`
	Solutions          []Solution           `json:"solutions"`
	DocumentationLinks []string             `json:"documentation_links"`
	Stacktrace         []CompleteStackFrame `json:"stacktrace"`
	Context            CompleteErrorContext `json:"context"`
	Stage              string               `json:"stage"`
	MessageLevel       interface{}          `json:"message_level"`
	OpenFrameIndex     interface{}          `json:"open_frame_index"`
	ApplicationPath    string               `json:"application_path"`
	ApplicationVersion interface{}          `json:"application_version"`
	TrackingUuid       string               `json:"tracking_uuid"`
	Handled            interface{}          `json:"handled"`
	OverriddenGrouping interface{}          `json:"overridden_grouping"`
}

// CompleteIgnitionData represents the full Ignition data structure
type CompleteIgnitionData struct {
	Report               CompleteReport  `json:"report"`
	ShareableReport      ShareableReport `json:"shareableReport"`
	Config               ConfigData      `json:"config"`
	Solutions            []Solution      `json:"solutions"`
	UpdateConfigEndpoint string          `json:"updateConfigEndpoint"`
}

// ConfigData represents the configuration data for Ignition
type ConfigData struct {
	Editor                  string                  `json:"editor"`
	Theme                   string                  `json:"theme"`
	HideSolutions           bool                    `json:"hideSolutions"`
	RemoteSitesPath         string                  `json:"remoteSitesPath"`
	LocalSitesPath          string                  `json:"localSitesPath"`
	EnableShareButton       bool                    `json:"enableShareButton"`
	EnableRunnableSolutions bool                    `json:"enableRunnableSolutions"`
	DirectorySeparator      string                  `json:"directorySeparator"`
	EditorOptions           map[string]EditorOption `json:"editorOptions"`
	ShareEndpoint           string                  `json:"shareEndpoint"`
}

// EditorOption represents editor configuration options
type EditorOption struct {
	Label     string `json:"label"`
	URL       string `json:"url"`
	Clipboard bool   `json:"clipboard,omitempty"`
}

// CreateCompleteReport creates a complete report from ErrorReport and HTTP request
func CreateCompleteReport(errorReport *ErrorReport, req *http.Request, applicationPath string) *CompleteReport {
	context := NewCompleteErrorContext(req)

	// Convert stack frames
	stackFrames := make([]CompleteStackFrame, 0, len(errorReport.GetStack()))
	for _, frame := range errorReport.GetStack() {
		stackFrames = append(stackFrames, CompleteStackFrame{
			File:             ensureString(frame.GetFile()),
			LineNumber:       frame.GetLine(),
			Method:           ensureString(frame.GetFunction()),
			Class:            nil,
			CodeSnippet:      frame.GetCode(),
			Arguments:        []interface{}{},
			ApplicationFrame: isApplicationFrame(frame.GetFile()),
		})
	}

	now := time.Now()
	return &CompleteReport{
		Notifier:         "GoVel Client",
		Language:         "Go",
		FrameworkVersion: "1.0.0",
		LanguageVersion:  runtime.Version(),
		ExceptionClass:   errorReport.GetType(),
		SeenAt:           now.Unix(),
		Message:          errorReport.GetMessage(),
		Glows:            []interface{}{},
		Solutions:        ensureSolutionsArray(errorReport.GetSolutions()),
		DocumentationLinks: []string{
			"https://golang.org/doc/",
			"https://golang.org/pkg/",
			"https://gobyexample.com/",
		},
		Stacktrace:         stackFrames,
		Context:            *context,
		Stage:              "local",
		MessageLevel:       nil,
		OpenFrameIndex:     nil,
		ApplicationPath:    applicationPath,
		ApplicationVersion: nil,
		TrackingUuid:       generateUUID(),
		Handled:            nil,
		OverriddenGrouping: nil,
	}
}

// CreateShareableReport creates a shareable report (potentially with filtered sensitive data)
func CreateShareableReport(report *CompleteReport) *ShareableReport {
	// For now, shareable report is identical to the main report
	// In production, you might want to filter sensitive information
	return &ShareableReport{
		Notifier:           report.Notifier,
		Language:           report.Language,
		FrameworkVersion:   report.FrameworkVersion,
		LanguageVersion:    report.LanguageVersion,
		ExceptionClass:     report.ExceptionClass,
		SeenAt:             report.SeenAt,
		Message:            report.Message,
		Glows:              report.Glows,
		Solutions:          ensureSolutionsArray(report.Solutions),
		DocumentationLinks: report.DocumentationLinks,
		Stacktrace:         report.Stacktrace,
		Context:            report.Context,
		Stage:              report.Stage,
		MessageLevel:       report.MessageLevel,
		OpenFrameIndex:     report.OpenFrameIndex,
		ApplicationPath:    report.ApplicationPath,
		ApplicationVersion: report.ApplicationVersion,
		TrackingUuid:       report.TrackingUuid,
		Handled:            report.Handled,
		OverriddenGrouping: report.OverriddenGrouping,
	}
}

// CreateCompleteIgnitionData creates the complete Ignition data structure
func CreateCompleteIgnitionData(errorReport *ErrorReport, req *http.Request, applicationPath, editor, theme string) *CompleteIgnitionData {
	report := CreateCompleteReport(errorReport, req, applicationPath)
	shareableReport := CreateShareableReport(report)

	config := ConfigData{
		Editor:                  editor,
		Theme:                   theme,
		HideSolutions:           false,
		RemoteSitesPath:         applicationPath,
		LocalSitesPath:          "",
		EnableShareButton:       true,
		EnableRunnableSolutions: false,
		DirectorySeparator:      "/",
		EditorOptions:           getEditorOptions(),
		ShareEndpoint:           "https://flareapp.io/api/public-reports",
	}

	return &CompleteIgnitionData{
		Report:               *report,
		ShareableReport:      *shareableReport,
		Config:               config,
		Solutions:            ensureSolutionsArray(errorReport.GetSolutions()),
		UpdateConfigEndpoint: "/_ignition/update-config",
	}
}

// Helper functions

func ensureString(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}

func ensureSolutionsArray(solutions []Solution) []Solution {
	if solutions == nil {
		return []Solution{}
	}
	return solutions
}

func isApplicationFrame(file string) bool {
	// Consider a frame as application frame if it's not in Go runtime or standard library
	return !strings.Contains(file, "runtime/") &&
		!strings.Contains(file, "/go") &&
		!strings.Contains(file, "vendor/")
}

func generateUUID() string {
	// Simple UUID generation - in production you'd use a proper UUID library
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

func getEditorOptions() map[string]EditorOption {
	return map[string]EditorOption{
		"clipboard": {
			Label:     "Clipboard",
			URL:       "%path:%line",
			Clipboard: true,
		},
		"sublime": {
			Label: "Sublime",
			URL:   "subl://open?url=file://%path&line=%line",
		},
		"textmate": {
			Label: "TextMate",
			URL:   "txmt://open?url=file://%path&line=%line",
		},
		"emacs": {
			Label: "Emacs",
			URL:   "emacs://open?url=file://%path&line=%line",
		},
		"macvim": {
			Label: "MacVim",
			URL:   "mvim://open/?url=file://%path&line=%line",
		},
		"phpstorm": {
			Label: "PhpStorm",
			URL:   "phpstorm://open?file=%path&line=%line",
		},
		"phpstorm-remote": {
			Label: "PHPStorm Remote",
			URL:   "javascript:r = new XMLHttpRequest;r.open(\"get\", \"http://localhost:63342/api/file/%path:%line\");r.send()",
		},
		"idea": {
			Label: "Idea",
			URL:   "idea://open?file=%path&line=%line",
		},
		"vscode": {
			Label: "VS Code",
			URL:   "vscode://file/%path:%line",
		},
		"vscode-insiders": {
			Label: "VS Code Insiders",
			URL:   "vscode-insiders://file/%path:%line",
		},
		"vscode-remote": {
			Label: "VS Code Remote",
			URL:   "vscode://vscode-remote/%path:%line",
		},
		"vscode-insiders-remote": {
			Label: "VS Code Insiders Remote",
			URL:   "vscode-insiders://vscode-remote/%path:%line",
		},
		"vscodium": {
			Label: "VS Codium",
			URL:   "vscodium://file/%path:%line",
		},
		"cursor": {
			Label: "Cursor",
			URL:   "cursor://file/%path:%line",
		},
		"atom": {
			Label: "Atom",
			URL:   "atom://core/open/file?filename=%path&line=%line",
		},
		"nova": {
			Label: "Nova",
			URL:   "nova://open?path=%path&line=%line",
		},
		"netbeans": {
			Label: "NetBeans",
			URL:   "netbeans://open/?f=%path:%line",
		},
		"xdebug": {
			Label: "Xdebug",
			URL:   "xdebug://%path@%line",
		},
	}
}
