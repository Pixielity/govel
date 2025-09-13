package http

import (
	"strconv"

	"govel/exceptions/core/solution"
	solutionInterface "govel/exceptions/interfaces/solution"
)

// MethodNotAllowedSolution provides specific guidance for 405 Method Not Allowed errors
type MethodNotAllowedSolution struct {
	*solution.BaseSolution
	allowedMethods []string
	usedMethod     string
}

// NewMethodNotAllowedSolution creates a solution specifically for method not allowed errors
func NewMethodNotAllowedSolution(usedMethod string, allowedMethods []string) *MethodNotAllowedSolution {
	description := "The HTTP method you're using is not allowed for this endpoint.\n\n"
	
	if usedMethod != "" {
		description += "• You used: " + usedMethod + "\n"
	}
	
	if len(allowedMethods) > 0 {
		description += "• Allowed methods: "
		for i, method := range allowedMethods {
			if i > 0 {
				description += ", "
			}
			description += method
		}
		description += "\n"
	}
	
	description += "\nCheck your route definitions and HTTP client configuration."
	
	base := solution.NewBaseSolution("HTTP Method Not Allowed").
		SetSolutionDescription(description).
		AddDocumentationLink("HTTP 405 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405").
		AddDocumentationLink("GoVel Routing Methods", "https://govel.dev/docs/routing-methods").
		AddDocumentationLink("HTTP Methods Guide", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods")

	return &MethodNotAllowedSolution{
		BaseSolution:   base,
		allowedMethods: allowedMethods,
		usedMethod:     usedMethod,
	}
}

// GetAllowedMethods returns the allowed HTTP methods
func (s *MethodNotAllowedSolution) GetAllowedMethods() []string {
	return s.allowedMethods
}

// GetUsedMethod returns the HTTP method that was used
func (s *MethodNotAllowedSolution) GetUsedMethod() string {
	return s.usedMethod
}

// ValidationErrorSolution provides specific guidance for 422 Unprocessable Entity errors
type ValidationErrorSolution struct {
	*solution.BaseSolution
	validationErrors map[string][]string
}

// NewValidationErrorSolution creates a solution specifically for validation errors
func NewValidationErrorSolution(validationErrors map[string][]string) *ValidationErrorSolution {
	base := solution.NewBaseSolution("Request Validation Failed").
		SetSolutionDescription("The request data failed validation. Please check the following:\n\n• Ensure all required fields are provided\n• Verify field formats (email, date, etc.)\n• Check field length requirements\n• Validate data types and constraints").
		AddDocumentationLink("HTTP 422 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/422").
		AddDocumentationLink("GoVel Validation", "https://govel.dev/docs/validation").
		AddDocumentationLink("GoVel Form Requests", "https://govel.dev/docs/form-requests").
		AddDocumentationLink("Validation Rules Reference", "https://govel.dev/docs/validation-rules")

	return &ValidationErrorSolution{
		BaseSolution:     base,
		validationErrors: validationErrors,
	}
}

// GetValidationErrors returns the validation errors
func (s *ValidationErrorSolution) GetValidationErrors() map[string][]string {
	return s.validationErrors
}

// TooManyRequestsSolution provides specific guidance for 429 Too Many Requests errors
type TooManyRequestsSolution struct {
	*solution.BaseSolution
	retryAfter int
	rateLimit  int
}

// NewTooManyRequestsSolution creates a solution specifically for rate limiting errors
func NewTooManyRequestsSolution(retryAfter, rateLimit int) *TooManyRequestsSolution {
	description := "You've exceeded the rate limit. Solutions:\n\n• Wait before making more requests\n• Implement exponential backoff\n• Check if you can upgrade your rate limit\n• Cache responses to reduce API calls"
	
	if retryAfter > 0 {
		description += "\n• Retry after " + strconv.Itoa(retryAfter) + " seconds"
	}
	
	base := solution.NewBaseSolution("Rate Limit Exceeded").
		SetSolutionDescription(description).
		AddDocumentationLink("HTTP 429 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429").
		AddDocumentationLink("GoVel Rate Limiting", "https://govel.dev/docs/rate-limiting").
		AddDocumentationLink("API Rate Limits", "https://govel.dev/docs/api-limits").
		AddDocumentationLink("Caching Strategies", "https://govel.dev/docs/caching")

	return &TooManyRequestsSolution{
		BaseSolution: base,
		retryAfter:   retryAfter,
		rateLimit:    rateLimit,
	}
}

// GetRetryAfter returns the retry after time in seconds
func (s *TooManyRequestsSolution) GetRetryAfter() int {
	return s.retryAfter
}

// GetRateLimit returns the rate limit value
func (s *TooManyRequestsSolution) GetRateLimit() int {
	return s.rateLimit
}

// InternalServerErrorSolution provides specific guidance for 500 Internal Server errors
type InternalServerErrorSolution struct {
	*solution.BaseSolution
}

// NewInternalServerErrorSolution creates a solution specifically for 500 errors
func NewInternalServerErrorSolution() *InternalServerErrorSolution {
	base := solution.NewBaseSolution("Internal Server Error").
		SetSolutionDescription("An unexpected error occurred on the server. To troubleshoot:\n\n• Check the application logs for detailed error information\n• Verify database connections and external service availability\n• Review recent code changes that might have introduced bugs\n• Check server resources (memory, disk space)\n• Ensure all required environment variables are set").
		AddDocumentationLink("HTTP 500 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500").
		AddDocumentationLink("GoVel Error Handling", "https://govel.dev/docs/error-handling").
		AddDocumentationLink("GoVel Logging", "https://govel.dev/docs/logging").
		AddDocumentationLink("GoVel Debugging", "https://govel.dev/docs/debugging").
		AddDocumentationLink("Server Monitoring", "https://govel.dev/docs/monitoring")

	return &InternalServerErrorSolution{
		BaseSolution: base,
	}
}

// ServiceUnavailableSolution provides specific guidance for 503 Service Unavailable errors
type ServiceUnavailableSolution struct {
	*solution.BaseSolution
	retryAfter int
}

// NewServiceUnavailableSolution creates a solution specifically for 503 errors
func NewServiceUnavailableSolution(retryAfter int) *ServiceUnavailableSolution {
	description := "The service is temporarily unavailable. This might be due to:\n\n• Scheduled maintenance\n• Server overload\n• Database connectivity issues\n• External service dependencies\n• Resource exhaustion"
	
	if retryAfter > 0 {
		description += "\n\nTry again after " + strconv.Itoa(retryAfter) + " seconds."
	}
	
	base := solution.NewBaseSolution("Service Temporarily Unavailable").
		SetSolutionDescription(description).
		AddDocumentationLink("HTTP 503 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503").
		AddDocumentationLink("GoVel Health Checks", "https://govel.dev/docs/health-checks").
		AddDocumentationLink("GoVel Maintenance Mode", "https://govel.dev/docs/maintenance").
		AddDocumentationLink("Service Monitoring", "https://govel.dev/docs/monitoring")

	return &ServiceUnavailableSolution{
		BaseSolution: base,
		retryAfter:   retryAfter,
	}
}

// GetRetryAfter returns the retry after time in seconds
func (s *ServiceUnavailableSolution) GetRetryAfter() int {
	return s.retryAfter
}

// ConflictSolution provides specific guidance for 409 Conflict errors
type ConflictSolution struct {
	*solution.BaseSolution
}

// NewConflictSolution creates a solution specifically for 409 errors
func NewConflictSolution() *ConflictSolution {
	base := solution.NewBaseSolution("Resource Conflict").
		SetSolutionDescription("The request conflicts with the current state of the server. Common causes:\n\n• Resource already exists with the same identifier\n• Concurrent modification conflicts\n• Business rule violations\n• Version mismatch in optimistic locking\n• Duplicate operation attempts").
		AddDocumentationLink("HTTP 409 Reference", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409").
		AddDocumentationLink("GoVel Validation", "https://govel.dev/docs/validation").
		AddDocumentationLink("Concurrency Control", "https://govel.dev/docs/concurrency").
		AddDocumentationLink("Database Constraints", "https://govel.dev/docs/database-constraints")

	return &ConflictSolution{
		BaseSolution: base,
	}
}

// Ensure all solutions implement the Solution interface
var _ solutionInterface.Solution = (*MethodNotAllowedSolution)(nil)
var _ solutionInterface.Solution = (*ValidationErrorSolution)(nil)
var _ solutionInterface.Solution = (*TooManyRequestsSolution)(nil)
var _ solutionInterface.Solution = (*InternalServerErrorSolution)(nil)
var _ solutionInterface.Solution = (*ServiceUnavailableSolution)(nil)
var _ solutionInterface.Solution = (*ConflictSolution)(nil)
