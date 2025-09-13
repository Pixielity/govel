package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"govel/packages/ignition/config"
	"govel/packages/ignition/interfaces"
	"govel/packages/ignition/models"
	"govel/packages/ignition/renderer"
)

// ErrorHandler handles errors and creates reports
type ErrorHandler struct {
	config            *config.Config
	shouldDisplay     bool
	applicationPath   string
	middleware        []interfaces.MiddlewareInterface
	solutionProviders []interfaces.SolutionProviderInterface
	renderer          *renderer.HTMLRenderer
	customHTMLHead    string
	customHTMLBody    string
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(cfg *config.Config) *ErrorHandler {
	return &ErrorHandler{
		config:            cfg,
		shouldDisplay:     true,
		middleware:        []interfaces.MiddlewareInterface{},
		solutionProviders: []interfaces.SolutionProviderInterface{},
		renderer:          renderer.NewHTMLRenderer(),
	}
}

// SetShouldDisplay sets whether to display exceptions
func (h *ErrorHandler) SetShouldDisplay(should bool) *ErrorHandler {
	h.shouldDisplay = should
	return h
}

// SetApplicationPath sets the application root path
func (h *ErrorHandler) SetApplicationPath(path string) *ErrorHandler {
	h.applicationPath = path
	return h
}

// AddCustomHTMLToHead adds custom HTML to the error page head
func (h *ErrorHandler) AddCustomHTMLToHead(html string) *ErrorHandler {
	h.customHTMLHead += html
	return h
}

// AddCustomHTMLToBody adds custom HTML to the error page body
func (h *ErrorHandler) AddCustomHTMLToBody(html string) *ErrorHandler {
	h.customHTMLBody += html
	return h
}

// AddSolutionProviders adds solution providers
func (h *ErrorHandler) AddSolutionProviders(providers []interfaces.SolutionProviderInterface) *ErrorHandler {
	h.solutionProviders = append(h.solutionProviders, providers...)
	return h
}

// RegisterMiddleware registers middleware for processing reports
func (h *ErrorHandler) RegisterMiddleware(middleware []interfaces.MiddlewareInterface) *ErrorHandler {
	h.middleware = append(h.middleware, middleware...)
	return h
}

// HandleError handles an error and renders the error page
func (h *ErrorHandler) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	if !h.shouldDisplay {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	report := h.CreateReport(err, r)
	h.renderer.RenderErrorPage(report, w, r, h.config, h.applicationPath, h.customHTMLHead, h.customHTMLBody)
}

// CreateReport creates a structured error report
func (h *ErrorHandler) CreateReport(err error, r *http.Request) *models.ErrorReport {
	stack := h.buildStackTrace(err)

	report := models.NewErrorReport()
	report.SetMessage(err.Error())
	report.SetStack(stack)
	report.SetType(h.classifyErrorType(err, stack))
	report.SetContext(h.buildContext(r))
	report.SetSolutions(h.getSolutions(err))
	report.SetTimestamp(time.Now())

	// Set file and line from first stack frame
	if len(stack) > 0 {
		report.SetFile(stack[0].GetFile())
		report.SetLine(stack[0].GetLine())
	}

	// Apply middleware (temporarily disabled due to interface compatibility issues)
	_ = h.middleware
	// for _, middleware := range h.middleware {
	// 	middleware.Process(report)
	// }

	return report
}

// buildStackTrace builds a detailed stack trace
func (h *ErrorHandler) buildStackTrace(err error) []models.StackFrame {
	var frames []models.StackFrame

	// Get the current stack
	pcs := make([]uintptr, 50)
	n := runtime.Callers(0, pcs)

	for idx := 0; idx < n; idx++ {
		pc := pcs[idx]
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		file, line := fn.FileLine(pc)

		// Skip internal Go files and our own stack frames
		if strings.Contains(file, "runtime/") ||
			strings.Contains(file, "error_handler.go") {
			continue
		}

		frame := models.NewStackFrame()
		frame.SetFunction(fn.Name())
		frame.SetFile(file)
		frame.SetLine(line)
		frame.SetCode(h.getSourceCode(file, line))

		frames = append(frames, *frame)
	}

	return frames
}

// getSourceCode retrieves source code around the error line
func (h *ErrorHandler) getSourceCode(file string, line int) map[string]string {
	// Use the enhanced source code extractor
	extractor := NewSourceCodeExtractor(5)
	return extractor.ExtractSourceCode(file, line)
}

// buildContext builds the error context
func (h *ErrorHandler) buildContext(r *http.Request) models.ErrorContext {
	// For now, use the basic context. The complete context is built in the renderer
	context := models.NewErrorContext()
	envContext := models.NewEnvContext()
	context.SetEnvironment(envContext)
	return *context
}

// getSolutions gets solutions for the error
func (h *ErrorHandler) getSolutions(err error) []models.Solution {
	var solutions []models.Solution

	for _, provider := range h.solutionProviders {
		providerSolutions := provider.GetSolutions(err)
		for _, solution := range providerSolutions {
			// Convert interface to concrete type
			if concreteSolution, ok := solution.(*models.Solution); ok {
				solutions = append(solutions, *concreteSolution)
			}
		}
	}

	return solutions
}

// classifyErrorType returns the filename where the error occurred (for exception_class)
// It takes the stack trace to get the actual error location
func (h *ErrorHandler) classifyErrorType(err error, stack []models.StackFrame) string {
	// Use the first application frame from the stack trace
	for _, frame := range stack {
		if frame.IsApplicationFrame() {
			file := frame.GetFile()
			// Return the relative path from application root if possible
			if h.applicationPath != "" && strings.HasPrefix(file, h.applicationPath) {
				return strings.TrimPrefix(file, h.applicationPath+"/")
			}
			return file
		}
	}

	// If no application frame found, use the first frame
	if len(stack) > 0 {
		file := stack[0].GetFile()
		if h.applicationPath != "" && strings.HasPrefix(file, h.applicationPath) {
			return strings.TrimPrefix(file, h.applicationPath+"/")
		}
		return file
	}

	// Fallback to generic error if we can't determine the file
	return "UnknownError"
}

// Middleware creates an HTTP middleware for automatic error handling
func (h *ErrorHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				var err error
				if e, ok := recovered.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", recovered)
				}
				h.HandleError(err, w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
