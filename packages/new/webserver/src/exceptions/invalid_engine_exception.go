// Package exceptions - Invalid engine exception
package exceptions

import "fmt"

// InvalidEngineException is thrown when an unsupported engine is selected.
type InvalidEngineException struct{ Engine string }

func (e *InvalidEngineException) Error() string {
	return fmt.Sprintf("invalid webserver engine: %s", e.Engine)
}
func NewInvalidEngine(engine string) error { return &InvalidEngineException{Engine: engine} }
