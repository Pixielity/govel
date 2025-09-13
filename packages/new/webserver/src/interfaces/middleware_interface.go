// Package interfaces - Middleware interface definition
// This file defines the MiddlewareInterface contract for framework-agnostic middleware.
package interfaces

// MiddlewareInterface represents a middleware component that can process requests
// with separate before, main, and after phases.
//
// Middlewares are executed in three phases:
//  1. Before: Pre-processing logic (authentication, request validation, etc.)
//  2. Handle: Main middleware processing (can short-circuit the chain)
//  3. After: Post-processing logic (response modification, logging, etc.)
//
// Each phase is optional - middleware can embed BaseMiddleware and only override
// the methods they need to customize.
//
// Example:
//   type LoggingMiddleware struct {
//       BaseMiddleware
//   }
//   func (m *LoggingMiddleware) Before(req RequestInterface) error {
//       log.Printf("Request: %s %s", req.Method(), req.Path())
//       return nil
//   }
//   func (m *LoggingMiddleware) After(req RequestInterface, resp ResponseInterface) ResponseInterface {
//       log.Printf("Response: %d", resp.StatusCode())
//       return resp
//   }
type MiddlewareInterface interface {
	// Before is called before the main handler chain is executed.
	// This is where you can perform request preprocessing, validation, authentication, etc.
	//
	// Parameters:
	//   req: The incoming request
	//
	// Returns:
	//   error: If an error is returned, the middleware chain is short-circuited
	Before(req RequestInterface) error

	// Handle processes the incoming request and optionally calls the next handler.
	// This is the main middleware logic where you can modify the request,
	// short-circuit the chain, or pass control to the next handler.
	//
	// Parameters:
	//   req: The incoming request
	//   next: The next handler in the chain
	//
	// Returns:
	//   ResponseInterface: The response to send back to the client
	Handle(req RequestInterface, next HandlerInterface) ResponseInterface

	// After is called after the handler chain has been executed.
	// This is where you can perform response post-processing, logging, cleanup, etc.
	//
	// Parameters:
	//   req: The original request
	//   resp: The response from the handler chain
	//
	// Returns:
	//   ResponseInterface: The potentially modified response
	After(req RequestInterface, resp ResponseInterface) ResponseInterface

	// Priority returns the middleware execution priority.
	// Lower values execute earlier in the chain.
	// Default priority is 0.
	//
	// Returns:
	//   int: The priority value (lower = earlier execution)
	Priority() int
}
