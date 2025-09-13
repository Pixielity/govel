package ignition

import (
	"net/http"

	"govel/packages/ignition/config"
	"govel/packages/ignition/enums"
	"govel/packages/ignition/handlers"
	"govel/packages/ignition/interfaces"
	"govel/packages/ignition/views"
)

// Ignition is the main facade that provides a clean API for error handling
type Ignition struct {
	handler *handlers.ErrorHandler
	config  *config.Config
}

// New creates a new Ignition instance with sensible defaults
func New() *Ignition {
	cfg := config.NewConfig()
	handler := handlers.NewErrorHandler(cfg)

	return &Ignition{
		handler: handler,
		config:  cfg,
	}
}

// Make creates a new Ignition instance (alias for New)
func Make() *Ignition {
	return New()
}

// SetTheme sets the error page theme
func (i *Ignition) SetTheme(theme enums.Theme) *Ignition {
	i.config.WithTheme(theme)
	return i
}

// SetEditor sets the preferred editor for opening files
func (i *Ignition) SetEditor(editor enums.Editor) *Ignition {
	i.config.WithEditor(editor)
	return i
}

// ShouldDisplayException sets whether to display exceptions
func (i *Ignition) ShouldDisplayException(should bool) *Ignition {
	i.handler.SetShouldDisplay(should)
	return i
}

// ApplicationPath sets the application root path
func (i *Ignition) ApplicationPath(path string) *Ignition {
	i.handler.SetApplicationPath(path)
	return i
}

// AddCustomHTMLToHead adds custom HTML to the error page head
func (i *Ignition) AddCustomHTMLToHead(html string) *Ignition {
	i.handler.AddCustomHTMLToHead(html)
	return i
}

// AddCustomHTMLToBody adds custom HTML to the error page body
func (i *Ignition) AddCustomHTMLToBody(html string) *Ignition {
	i.handler.AddCustomHTMLToBody(html)
	return i
}

// AddSolutionProviders adds solution providers
func (i *Ignition) AddSolutionProviders(providers []interfaces.SolutionProviderInterface) *Ignition {
	i.handler.AddSolutionProviders(providers)
	return i
}

// RegisterMiddleware registers middleware for processing reports
func (i *Ignition) RegisterMiddleware(middleware []interfaces.MiddlewareInterface) *Ignition {
	i.handler.RegisterMiddleware(middleware)
	return i
}

// HandleError handles an error and renders the error page
func (i *Ignition) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	i.handler.HandleError(err, w, r)
}

// Middleware creates an HTTP middleware for automatic error handling
func (i *Ignition) Middleware(next http.Handler) http.Handler {
	return i.handler.Middleware(next)
}

// GetAssetContent reads asset content from the embedded filesystem
func (i *Ignition) GetAssetContent(filename string) string {
	return views.GetAssetString(filename)
}

// Register sets up global error handlers (for backwards compatibility)
func (i *Ignition) Register() *Ignition {
	return i
}
