package models

import (
	"html/template"
)

// TemplateData holds all the data needed for the error page template
type TemplateData struct {
	Theme        string                `json:"theme"`
	ErrorType    string                `json:"error_type"`
	ErrorMessage string                `json:"error_message"`
	CSS          template.CSS          `json:"css"`
	JavaScript   template.JS           `json:"javascript"`
	RequireJS    template.JS           `json:"require_js"`
	ReportJSON   template.JS           `json:"report_json"`
	CustomHead   template.HTML         `json:"custom_head"`
	CustomBody   template.HTML         `json:"custom_body"`
	StackFrames  []StackFrameTemplate  `json:"stack_frames"`
}

// NewTemplateData creates a new template data instance
func NewTemplateData() *TemplateData {
	return &TemplateData{
		StackFrames: []StackFrameTemplate{},
	}
}

// GetTheme returns the theme
func (t *TemplateData) GetTheme() string {
	return t.Theme
}

// SetTheme sets the theme
func (t *TemplateData) SetTheme(theme string) {
	t.Theme = theme
}

// GetErrorType returns the error type
func (t *TemplateData) GetErrorType() string {
	return t.ErrorType
}

// SetErrorType sets the error type
func (t *TemplateData) SetErrorType(errorType string) {
	t.ErrorType = errorType
}

// GetErrorMessage returns the error message
func (t *TemplateData) GetErrorMessage() string {
	return t.ErrorMessage
}

// SetErrorMessage sets the error message
func (t *TemplateData) SetErrorMessage(errorMessage string) {
	t.ErrorMessage = errorMessage
}

// GetCSS returns the CSS content
func (t *TemplateData) GetCSS() template.CSS {
	return t.CSS
}

// SetCSS sets the CSS content
func (t *TemplateData) SetCSS(css template.CSS) {
	t.CSS = css
}

// SetCSSString sets the CSS content from a string
func (t *TemplateData) SetCSSString(css string) {
	t.CSS = template.CSS(css)
}

// GetJavaScript returns the JavaScript content
func (t *TemplateData) GetJavaScript() template.JS {
	return t.JavaScript
}

// SetJavaScript sets the JavaScript content
func (t *TemplateData) SetJavaScript(js template.JS) {
	t.JavaScript = js
}

// SetJavaScriptString sets the JavaScript content from a string
func (t *TemplateData) SetJavaScriptString(js string) {
	t.JavaScript = template.JS(js)
}

// GetRequireJS returns the RequireJS content
func (t *TemplateData) GetRequireJS() template.JS {
	return t.RequireJS
}

// SetRequireJS sets the RequireJS content
func (t *TemplateData) SetRequireJS(requireJS template.JS) {
	t.RequireJS = requireJS
}

// SetRequireJSString sets the RequireJS content from a string
func (t *TemplateData) SetRequireJSString(requireJS string) {
	t.RequireJS = template.JS(requireJS)
}

// GetReportJSON returns the report JSON
func (t *TemplateData) GetReportJSON() template.JS {
	return t.ReportJSON
}

// SetReportJSON sets the report JSON
func (t *TemplateData) SetReportJSON(reportJSON template.JS) {
	t.ReportJSON = reportJSON
}

// SetReportJSONBytes sets the report JSON from bytes
func (t *TemplateData) SetReportJSONBytes(reportJSON []byte) {
	t.ReportJSON = template.JS(reportJSON)
}

// GetCustomHead returns the custom head content
func (t *TemplateData) GetCustomHead() template.HTML {
	return t.CustomHead
}

// SetCustomHead sets the custom head content
func (t *TemplateData) SetCustomHead(customHead template.HTML) {
	t.CustomHead = customHead
}

// SetCustomHeadString sets the custom head content from a string
func (t *TemplateData) SetCustomHeadString(customHead string) {
	t.CustomHead = template.HTML(customHead)
}

// GetCustomBody returns the custom body content
func (t *TemplateData) GetCustomBody() template.HTML {
	return t.CustomBody
}

// SetCustomBody sets the custom body content
func (t *TemplateData) SetCustomBody(customBody template.HTML) {
	t.CustomBody = customBody
}

// SetCustomBodyString sets the custom body content from a string
func (t *TemplateData) SetCustomBodyString(customBody string) {
	t.CustomBody = template.HTML(customBody)
}

// GetStackFrames returns the stack frames
func (t *TemplateData) GetStackFrames() []StackFrameTemplate {
	return t.StackFrames
}

// SetStackFrames sets the stack frames
func (t *TemplateData) SetStackFrames(stackFrames []StackFrameTemplate) {
	t.StackFrames = stackFrames
}

// AddStackFrame adds a stack frame to the collection
func (t *TemplateData) AddStackFrame(frame StackFrameTemplate) {
	t.StackFrames = append(t.StackFrames, frame)
}

// GetStackFrameCount returns the number of stack frames
func (t *TemplateData) GetStackFrameCount() int {
	return len(t.StackFrames)
}

// HasStackFrames returns true if there are stack frames
func (t *TemplateData) HasStackFrames() bool {
	return len(t.StackFrames) > 0
}

// ClearStackFrames removes all stack frames
func (t *TemplateData) ClearStackFrames() {
	t.StackFrames = []StackFrameTemplate{}
}

// IsEmpty returns true if the template data is empty
func (t *TemplateData) IsEmpty() bool {
	return t.Theme == "" && t.ErrorType == "" && t.ErrorMessage == ""
}

// HasCustomContent returns true if custom head or body content is present
func (t *TemplateData) HasCustomContent() bool {
	return t.CustomHead != "" || t.CustomBody != ""
}

// CreateFromReport creates template data from an error report and configuration
func CreateFromReport(report *ErrorReport, theme, css, js, requireJS string, reportJSON []byte, customHead, customBody string) *TemplateData {
	// Convert stack frames for template
	stackFrames := make([]StackFrameTemplate, len(report.GetStack()))
	for i, frame := range report.GetStack() {
		stackFrames[i] = StackFrameTemplate{
			Function: frame.GetFunction(),
			File:     frame.GetFile(),
			Line:     frame.GetLine(),
		}
	}

	return &TemplateData{
		Theme:        theme,
		ErrorType:    report.GetType(),
		ErrorMessage: report.GetMessage(),
		CSS:          template.CSS(css),
		JavaScript:   template.JS(js),
		RequireJS:    template.JS(requireJS),
		ReportJSON:   template.JS(reportJSON),
		CustomHead:   template.HTML(customHead),
		CustomBody:   template.HTML(customBody),
		StackFrames:  stackFrames,
	}
}
