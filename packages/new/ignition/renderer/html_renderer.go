package renderer

import (
	"encoding/json"
	"html/template"
	"net/http"

	"govel/ignition/config"
	"govel/ignition/models"
	"govel/ignition/views"
)

// HTMLRenderer handles rendering error pages to HTML
type HTMLRenderer struct {
	tmpl *template.Template
}

// NewHTMLRenderer creates a new HTML renderer
func NewHTMLRenderer() *HTMLRenderer {
	templateContent := views.GetAssetString("templates/error_page.html")
	if templateContent == "" {
		panic("Failed to read error page template")
	}
	tmpl := template.Must(template.New("error_page").Parse(templateContent))
	return &HTMLRenderer{
		tmpl: tmpl,
	}
}

// RenderErrorPage renders the beautiful error page using the template system
func (r *HTMLRenderer) RenderErrorPage(report *models.ErrorReport, w http.ResponseWriter, req *http.Request, cfg *config.Config, applicationPath, customHead, customBody string) {
	// Create the complete Ignition data structure matching Laravel Ignition exactly
	ignitionData := models.CreateCompleteIgnitionData(
		report,
		req,
		applicationPath,
		cfg.Editor.String(),
		cfg.Theme.String(),
	)

	// Convert to JSON for JavaScript
	reportJSON, _ := json.Marshal(ignitionData)

	// Create template data
	templateData := models.CreateFromReport(
		report,
		cfg.Theme.String(),
		r.getAssetContent("css/ignition.css"),
		r.getAssetContent("js/ignition.js"),
		r.getAssetContent("js/require.js"),
		reportJSON,
		customHead,
		customBody,
	)

	// Set response headers
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)

	// Render template
	if err := r.tmpl.Execute(w, templateData); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Removed old convertStackToFrames function - now using CompleteStackFrame in models

// getAssetContent reads asset content from the embedded filesystem
func (r *HTMLRenderer) getAssetContent(filename string) string {
	var filepath string
	switch filename {
	case "css/ignition.css":
		filepath = "assets/css/ignition.css"
	case "js/ignition.js":
		filepath = "assets/js/ignition.js"
	case "js/require.js":
		filepath = "assets/js/require.js"
	default:
		return ""
	}

	return views.GetAssetString(filepath)
}
