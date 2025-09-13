package interfaces

// PipeHandlerInterface defines the contract for pipeline pipe handlers.
// Objects implementing this interface can be used as pipes in the pipeline system.
// The Handle method is called when the pipe is executed during pipeline processing.
//
// Key features:
//   - Standard pipe execution contract
//   - Support for passable object transformation
//   - Integration with pipeline "next" function for chaining
//   - Error handling and propagation support
//   - Compatible with pipeline "Russian Doll" execution pattern
type PipeHandlerInterface interface {
	// Handle processes the passable object and continues the pipeline chain.
	// This method is called when the pipe is executed during pipeline processing.
	//
	// The pipe can:
	//   - Modify the passable object before passing to next
	//   - Perform side effects (logging, validation, etc.)
	//   - Conditionally call next based on logic
	//   - Handle errors from downstream pipes
	//   - Return early without calling next to short-circuit the pipeline
	//
	// Parameters:
	//   - passable: The object being passed through the pipeline
	//   - next: Function to call the next pipe in the chain
	//
	// Returns:
	//   - interface{}: The processed result (may be modified passable or completely new value)
	//   - error: Any error that occurred during pipe processing
	//
	// Example implementation:
	//
	//   func (p *MyPipe) Handle(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error) {
	//       // Pre-processing logic
	//       modifiedPassable := preprocessPassable(passable)
	//       
	//       // Call next pipe in chain
	//       result, err := next(modifiedPassable)
	//       if err != nil {
	//           return nil, err
	//       }
	//       
	//       // Post-processing logic
	//       return postprocessResult(result), nil
	//   }
	Handle(passable interface{}, next func(interface{}) (interface{}, error)) (interface{}, error)
}
