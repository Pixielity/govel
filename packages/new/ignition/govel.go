package ignition

import (
	"net/http"
	"runtime"
)

// GoVelIgnition represents the Ignition integration for GoVel framework
type GoVelIgnition struct {
	*Ignition
	isProduction bool
}

// NewGoVelIgnition creates a new GoVel-specific Ignition instance
func NewGoVelIgnition() *GoVelIgnition {
	ignition := New()

	// Add default solution providers for Go errors
	ignition.AddSolutionProviders(GetDefaultSolutionProviders())

	return &GoVelIgnition{
		Ignition:     ignition,
		isProduction: false,
	}
}

// Production sets the environment to production mode
func (g *GoVelIgnition) Production(isProd bool) *GoVelIgnition {
	g.isProduction = isProd
	g.ShouldDisplayException(!isProd) // Don't show detailed errors in production
	return g
}

// WithApplication integrates with a GoVel application
func (g *GoVelIgnition) WithApplication(app interface{}) *GoVelIgnition {
	// This would integrate with the actual GoVel application type
	// For now, we'll add some application-specific configuration

	// Try to extract application path from the runtime caller
	if _, filename, _, ok := runtime.Caller(1); ok {
		g.ApplicationPath(filename)
	}

	return g
}

// ErrorHandler creates an HTTP handler for manual error handling
func (g *GoVelIgnition) ErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This would be called when you want to manually trigger an error page
		// Could be useful for testing or custom error scenarios

		// For demo purposes, let's create a sample error
		err := &GoVelError{
			Message: "Sample GoVel Error",
			Code:    "GOVEL_SAMPLE_ERROR",
			File:    "example.go",
			Line:    42,
		}

		g.HandleError(err, w, r)
	}
}

// GoVelError represents a GoVel-specific error type
type GoVelError struct {
	Message string
	Code    string
	File    string
	Line    int
	Data    map[string]interface{}
}

// Error implements the error interface
func (e *GoVelError) Error() string {
	return e.Message
}

// GoVelMiddleware creates GoVel-specific middleware
type GoVelMiddleware struct {
	ignition *GoVelIgnition
	logger   interface{} // This would be the actual GoVel logger interface
}

// NewGoVelMiddleware creates new GoVel middleware
func NewGoVelMiddleware(ignition *GoVelIgnition) *GoVelMiddleware {
	return &GoVelMiddleware{
		ignition: ignition,
	}
}

// Handle implements middleware interface for GoVel
func (m *GoVelMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// Log the error if logger is available
				if m.logger != nil {
					// m.logger.Error(recovered)
				}

				// Convert panic to error and handle with Ignition
				var err error
				if e, ok := recovered.(error); ok {
					err = e
				} else {
					err = &GoVelError{
						Message: "Panic occurred",
						Code:    "GOVEL_PANIC",
						Data:    map[string]interface{}{"panic": recovered},
					}
				}

				m.ignition.HandleError(err, w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Configuration helpers for GoVel integration

// ConfigureForDevelopment configures Ignition for development environment
func (g *GoVelIgnition) ConfigureForDevelopment() *GoVelIgnition {
	g.Production(false)
	g.SetTheme("auto")
	g.SetEditor("vscode")
	g.ShouldDisplayException(true)
	return g
}

// ConfigureForProduction configures Ignition for production environment
func (g *GoVelIgnition) ConfigureForProduction() *GoVelIgnition {
	g.Production(true)
	g.ShouldDisplayException(false)
	return g
}

// ConfigureForTesting configures Ignition for testing environment
func (g *GoVelIgnition) ConfigureForTesting() *GoVelIgnition {
	g.Production(false)
	g.ShouldDisplayException(true)
	g.SetTheme("light") // Light theme might be better for CI/testing screenshots
	return g
}

// Helper methods for common GoVel patterns

// HandleControllerError handles errors that occur in GoVel controllers
func (g *GoVelIgnition) HandleControllerError(err error, controller string, action string, w http.ResponseWriter, r *http.Request) {
	// Enhance the error with controller/action information
	enhancedErr := &GoVelError{
		Message: err.Error(),
		Code:    "GOVEL_CONTROLLER_ERROR",
		Data: map[string]interface{}{
			"controller": controller,
			"action":     action,
			"original":   err,
		},
	}

	g.HandleError(enhancedErr, w, r)
}

// HandleServiceError handles errors that occur in GoVel services
func (g *GoVelIgnition) HandleServiceError(err error, service string, method string, w http.ResponseWriter, r *http.Request) {
	enhancedErr := &GoVelError{
		Message: err.Error(),
		Code:    "GOVEL_SERVICE_ERROR",
		Data: map[string]interface{}{
			"service":  service,
			"method":   method,
			"original": err,
		},
	}

	g.HandleError(enhancedErr, w, r)
}

// HandleRepositoryError handles errors that occur in GoVel repositories
func (g *GoVelIgnition) HandleRepositoryError(err error, repository string, operation string, w http.ResponseWriter, r *http.Request) {
	enhancedErr := &GoVelError{
		Message: err.Error(),
		Code:    "GOVEL_REPOSITORY_ERROR",
		Data: map[string]interface{}{
			"repository": repository,
			"operation":  operation,
			"original":   err,
		},
	}

	g.HandleError(enhancedErr, w, r)
}

// Integration points with GoVel's existing packages

// WithLogger integrates with GoVel's logger package
func (g *GoVelIgnition) WithLogger(logger interface{}) *GoVelIgnition {
	// This would integrate with the actual GoVel logger interface
	// For now, just store it for potential future use
	return g
}

// WithContainer integrates with GoVel's container package
func (g *GoVelIgnition) WithContainer(container interface{}) *GoVelIgnition {
	// This would integrate with the actual GoVel container interface
	// Could be used to resolve dependencies needed for error handling
	return g
}

// WithConfig integrates with GoVel's config package
func (g *GoVelIgnition) WithConfig(config interface{}) *GoVelIgnition {
	// This would integrate with the actual GoVel config interface
	// Could be used to load Ignition configuration from GoVel's config
	return g
}
