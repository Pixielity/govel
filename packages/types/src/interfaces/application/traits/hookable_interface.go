package interfaces

// HookableInterface defines the contract for hooks and events management functionality.
type HookableInterface interface {
	// RegisterHook registers a hook callback for a specific event
	RegisterHook(event string, callback func(interface{}))
	
	// TriggerHook triggers all registered hooks for a specific event
	TriggerHook(event string, data interface{})
	
	// RemoveHook removes a hook callback for a specific event
	RemoveHook(event string, callback func(interface{}))
	
	// GetHooks returns all registered hooks for a specific event
	GetHooks(event string) []func(interface{})
	
	// ClearHooks removes all hooks for a specific event
	ClearHooks(event string)
	
	// ClearAllHooks removes all registered hooks
	ClearAllHooks()
}