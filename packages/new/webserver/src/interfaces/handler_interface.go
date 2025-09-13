// Package interfaces - Handler interface definition
// This file defines the HandlerInterface contract representing a route handler.
package interfaces

// HandlerInterface represents an executable handler for a route.
// Implementations typically wrap a function that receives a request
// and returns a response.
//
// Example:
//   type MyHandler struct {}
//   func (h *MyHandler) Handle(req RequestInterface) ResponseInterface {
//       return NewResponse().Json(map[string]string{"message": "ok"})
//   }
type HandlerInterface interface {
	// Handle executes the handler logic for a given request.
	//
	// Parameters:
	//   req: The incoming request
	//
	// Returns:
	//   ResponseInterface: The response to send back to the client
	Handle(req RequestInterface) ResponseInterface
}
