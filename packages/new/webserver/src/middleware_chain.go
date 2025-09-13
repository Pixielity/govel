// Package types - Middleware-related type definitions.
// This file defines types and utilities for working with middleware in the webserver package.
package webserver

import (
	"govel/new/webserver/interfaces"
	"govel/new/webserver/types"
)

// MiddlewareChain represents a chain of middleware that can be executed in sequence.
// This type provides utilities for building and executing middleware chains.
type MiddlewareChain struct {
	// middleware stores the middleware stack in execution order
	middleware []interfaces.MiddlewareInterface
}

// NewMiddlewareChain creates a new middleware chain with the provided middleware.
//
// Parameters:
//
//	middleware: Zero or more middleware implementations to include in the chain
//
// Returns:
//
//	*MiddlewareChain: A new middleware chain instance
//
// Example:
//
//	chain := NewMiddlewareChain(corsMiddleware, authMiddleware, loggerMiddleware)
func NewMiddlewareChain(middleware ...interfaces.MiddlewareInterface) *MiddlewareChain {
	return &MiddlewareChain{
		middleware: append([]interfaces.MiddlewareInterface{}, middleware...),
	}
}

// Add appends middleware to the chain.
//
// Parameters:
//
//	middleware: One or more middleware implementations to add
//
// Returns:
//
//	*MiddlewareChain: The middleware chain instance for method chaining
func (mc *MiddlewareChain) Add(middleware ...interfaces.MiddlewareInterface) *MiddlewareChain {
	mc.middleware = append(mc.middleware, middleware...)
	return mc
}

// Prepend adds middleware to the beginning of the chain.
//
// Parameters:
//
//	middleware: One or more middleware implementations to prepend
//
// Returns:
//
//	*MiddlewareChain: The middleware chain instance for method chaining
func (mc *MiddlewareChain) Prepend(middleware ...interfaces.MiddlewareInterface) *MiddlewareChain {
	mc.middleware = append(middleware, mc.middleware...)
	return mc
}

// Execute executes the middleware chain with the provided handler as the final handler.
//
// Parameters:
//
//	req: The HTTP request to process
//	finalHandler: The final handler to execute after all middleware
//
// Returns:
//
//	interfaces.ResponseInterface: The response from the middleware chain
func (mc *MiddlewareChain) Execute(req interfaces.RequestInterface, finalHandler interfaces.HandlerInterface) interfaces.ResponseInterface {
	if len(mc.middleware) == 0 {
		return finalHandler.Handle(req)
	}

	// Build the chain from the end backwards
	handler := finalHandler

	for i := len(mc.middleware) - 1; i >= 0; i-- {
		middleware := mc.middleware[i]
		currentHandler := handler

		handler = types.HandlerFunc(func(req interfaces.RequestInterface) interfaces.ResponseInterface {
			return middleware.Handle(req, currentHandler)
		})
	}

	return handler.Handle(req)
}

// Length returns the number of middleware in the chain.
//
// Returns:
//
//	int: The number of middleware in the chain
func (mc *MiddlewareChain) Length() int {
	return len(mc.middleware)
}

// IsEmpty returns true if the middleware chain is empty.
//
// Returns:
//
//	bool: True if the chain contains no middleware, false otherwise
func (mc *MiddlewareChain) IsEmpty() bool {
	return len(mc.middleware) == 0
}

// ToSlice returns a copy of the middleware slice.
//
// Returns:
//
//	[]interfaces.MiddlewareInterface: A copy of the middleware stack
func (mc *MiddlewareChain) ToSlice() []interfaces.MiddlewareInterface {
	middlewareCopy := make([]interfaces.MiddlewareInterface, len(mc.middleware))
	copy(middlewareCopy, mc.middleware)
	return middlewareCopy
}

// Clone creates a deep copy of the middleware chain.
//
// Returns:
//
//	*MiddlewareChain: A new middleware chain with the same middleware
func (mc *MiddlewareChain) Clone() *MiddlewareChain {
	return NewMiddlewareChain(mc.ToSlice()...)
}
