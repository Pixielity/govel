// Package types - Middleware-related type definitions.
// This file defines types and utilities for working with middleware in the webserver package.
package webserver

import (
	"govel/new/webserver/src/interfaces"
)

// ConditionalMiddleware represents middleware that only executes under certain conditions.
// This is useful for middleware that should only apply to specific routes or requests.
type ConditionalMiddleware struct {
	// Middleware is the actual middleware implementation
	Middleware interfaces.MiddlewareInterface

	// Condition is a function that determines whether the middleware should execute
	Condition func(interfaces.RequestInterface) bool
}

// Before implements the MiddlewareInterface for ConditionalMiddleware.
// The Before method only executes if the condition returns true.
func (cm *ConditionalMiddleware) Before(req interfaces.RequestInterface) error {
	if cm.Condition(req) {
		return cm.Middleware.Before(req)
	}
	return nil
}

// Handle implements the MiddlewareInterface for ConditionalMiddleware.
// The middleware only executes if the condition returns true.
//
// Parameters:
//
//	req: The incoming HTTP request
//	next: The next handler in the middleware chain
//
// Returns:
//
//	interfaces.ResponseInterface: The response to send back to the client
func (cm *ConditionalMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	if cm.Condition(req) {
		return cm.Middleware.Handle(req, next)
	}
	return next.Handle(req)
}

// After implements the MiddlewareInterface for ConditionalMiddleware.
// The After method only executes if the condition returns true.
func (cm *ConditionalMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	if cm.Condition(req) {
		return cm.Middleware.After(req, resp)
	}
	return resp
}

// Priority implements the MiddlewareInterface for ConditionalMiddleware.
// Returns the priority of the wrapped middleware.
func (cm *ConditionalMiddleware) Priority() int {
	return cm.Middleware.Priority()
}

// NewConditionalMiddleware creates a new conditional middleware.
//
// Parameters:
//
//	middleware: The middleware to execute conditionally
//	condition: Function that determines when to execute the middleware
//
// Returns:
//
//	*ConditionalMiddleware: A new conditional middleware instance
//
// Example:
//
//	authMiddleware := NewConditionalMiddleware(
//	    myAuthMiddleware,
//	    func(req interfaces.RequestInterface) bool {
//	        return strings.HasPrefix(req.Path(), "/api/")
//	    },
//	)
func NewConditionalMiddleware(middleware interfaces.MiddlewareInterface, condition func(interfaces.RequestInterface) bool) *ConditionalMiddleware {
	return &ConditionalMiddleware{
		Middleware: middleware,
		Condition:  condition,
	}
}

// Ensure ConditionalMiddleware implements MiddlewareInterface at compile time.
var _ interfaces.MiddlewareInterface = (*ConditionalMiddleware)(nil)
