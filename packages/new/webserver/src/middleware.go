// Package webserver - Core middleware helpers
// This file provides convenience functions to work with middleware chains from the root package.
package webserver

import (
	"govel/new/webserver/src/interfaces"
)

// BaseMiddleware provides default implementations for all middleware methods.
// Custom middleware can embed this struct and only override the methods they need.
//
// Example:
//
//	type AuthMiddleware struct {
//	    BaseMiddleware
//	}
//	func (m *AuthMiddleware) Before(req interfaces.RequestInterface) error {
//	    if !isAuthenticated(req) {
//	        return errors.New("unauthorized")
//	    }
//	    return nil
//	}
type BaseMiddleware struct{}

// Before provides a no-op implementation for pre-processing.
// Override this method to add custom before logic.
func (b *BaseMiddleware) Before(req interfaces.RequestInterface) error {
	return nil
}

// Handle provides a pass-through implementation.
// Override this method to add custom middleware logic.
func (b *BaseMiddleware) Handle(req interfaces.RequestInterface, next interfaces.HandlerInterface) interfaces.ResponseInterface {
	return next.Handle(req)
}

// After provides a pass-through implementation for post-processing.
// Override this method to add custom after logic.
func (b *BaseMiddleware) After(req interfaces.RequestInterface, resp interfaces.ResponseInterface) interfaces.ResponseInterface {
	return resp
}

// Priority returns the default priority of 0.
// Override this method to set custom middleware priority.
func (b *BaseMiddleware) Priority() int {
	return 0
}

// ApplyMiddleware executes the middleware chain with a final handler.
func ApplyMiddleware(req interfaces.RequestInterface, chain *MiddlewareChain, final interfaces.HandlerInterface) interfaces.ResponseInterface {
	if chain == nil || chain.IsEmpty() {
		return final.Handle(req)
	}
	return chain.Execute(req, final)
}

// Compile-time interface compliance check
var _ interfaces.MiddlewareInterface = (*BaseMiddleware)(nil)
