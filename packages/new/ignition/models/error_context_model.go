package models

// ErrorContext holds contextual information about the error
type ErrorContext struct {
	Request     *RequestContext `json:"request,omitempty"`
	Environment *EnvContext     `json:"environment,omitempty"`
	User        interface{}     `json:"user,omitempty"`
}

// NewErrorContext creates a new error context
func NewErrorContext() *ErrorContext {
	return &ErrorContext{}
}

// GetRequest returns the request context
func (e *ErrorContext) GetRequest() *RequestContext {
	return e.Request
}

// SetRequest sets the request context
func (e *ErrorContext) SetRequest(request *RequestContext) {
	e.Request = request
}

// GetEnvironment returns the environment context
func (e *ErrorContext) GetEnvironment() *EnvContext {
	return e.Environment
}

// SetEnvironment sets the environment context
func (e *ErrorContext) SetEnvironment(environment *EnvContext) {
	e.Environment = environment
}

// GetUser returns the user context
func (e *ErrorContext) GetUser() interface{} {
	return e.User
}

// SetUser sets the user context
func (e *ErrorContext) SetUser(user interface{}) {
	e.User = user
}

// HasRequest returns true if request context is present
func (e *ErrorContext) HasRequest() bool {
	return e.Request != nil
}

// HasEnvironment returns true if environment context is present
func (e *ErrorContext) HasEnvironment() bool {
	return e.Environment != nil
}

// HasUser returns true if user context is present
func (e *ErrorContext) HasUser() bool {
	return e.User != nil
}

// IsEmpty returns true if the context is empty
func (e *ErrorContext) IsEmpty() bool {
	return e.Request == nil && e.Environment == nil && e.User == nil
}

// Note: Compile-time interface compliance check removed due to circular dependencies
// The ErrorContext interface references concrete types instead of interfaces
