// Package exceptions - Adapter not found exception
package exceptions

import "fmt"

// AdapterNotFoundException is thrown when an adapter can't be created/resolved.
type AdapterNotFoundException struct{ Engine string }

func (e *AdapterNotFoundException) Error() string {
	return fmt.Sprintf("adapter not found for engine: %s", e.Engine)
}
func NewAdapterNotFound(engine string) error { return &AdapterNotFoundException{Engine: engine} }
