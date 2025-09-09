package traits

import (
	"govel/packages/application/core/hooks"
	traitInterfaces "govel/packages/application/interfaces/traits"
	"govel/packages/application/types"
)

/**
 * Hookable provides application hook management functionality by wrapping
 * the core hooks manager. This trait follows the self-contained pattern
 * and delegates all operations to the underlying manager.
 */
type Hookable struct {
	/**
	 * manager is the underlying hooks manager instance
	 */
	manager *hooks.Manager
}

/**
 * NewHookable creates a new Hookable instance with a hooks manager.
 *
 * @return *Hookable The newly created trait instance
 */
func NewHookable() *Hookable {
	return &Hookable{
		manager: hooks.NewManager(),
	}
}

/**
 * NewHookableWithManager creates a new Hookable with an existing manager.
 *
 * @param m *hooks.Manager The hooks manager to wrap
 * @return *Hookable The newly created trait instance
 */
func NewHookableWithManager(m *hooks.Manager) *Hookable {
	return &Hookable{
		manager: m,
	}
}

/**
 * RegisterHook registers a callback function for the specified hook name.
 *
 * @param name string The name of the hook
 * @param priority int The priority of the callback (lower numbers execute first)
 * @param callback HookCallback The callback function to register
 */
func (t *Hookable) RegisterHook(name string, priority int, callback types.HookCallback) {
	// Convert interface callback to hooks callback
	hookCallback := hooks.HookCallback(callback)
	t.manager.RegisterHook(name, priority, hookCallback)
}

/**
 * UnregisterHook removes all callbacks for the specified hook name.
 *
 * @param name string The name of the hook to unregister
 * @return bool true if any callbacks were removed
 */
func (t *Hookable) UnregisterHook(name string) bool {
	return t.manager.UnregisterHook(name)
}

/**
 * UnregisterHookCallback removes a specific callback from a hook.
 *
 * @param name string The name of the hook
 * @param callback HookCallback The specific callback to remove
 * @return bool true if the callback was found and removed
 */
func (t *Hookable) UnregisterHookCallback(name string, callback types.HookCallback) bool {
	// Convert interface callback to hooks callback
	hookCallback := hooks.HookCallback(callback)
	return t.manager.UnregisterHookCallback(name, hookCallback)
}

/**
 * HasHook returns whether any callbacks are registered for the specified hook.
 *
 * @param name string The name of the hook to check
 * @return bool true if the hook has registered callbacks
 */
func (t *Hookable) HasHook(name string) bool {
	return t.manager.HasHook(name)
}

/**
 * GetHooks returns all registered hooks.
 *
 * @return map[string][]HookCallback Map of hook names to their callbacks
 */
func (t *Hookable) GetHooks() map[string][]types.HookCallback {
	managerHooks := t.manager.GetHooks()
	result := make(map[string][]types.HookCallback)

	for name, callbacks := range managerHooks {
		interfaceCallbacks := make([]types.HookCallback, len(callbacks))
		for i, callback := range callbacks {
			interfaceCallbacks[i] = types.HookCallback(callback)
		}
		result[name] = interfaceCallbacks
	}

	return result
}

/**
 * GetHookCallbacks returns all callbacks for a specific hook.
 *
 * @param name string The name of the hook
 * @return []HookCallback Slice of callbacks for the hook
 */
func (t *Hookable) GetHookCallbacks(name string) []types.HookCallback {
	managerCallbacks := t.manager.GetHookCallbacks(name)
	callbacks := make([]types.HookCallback, len(managerCallbacks))

	for i, callback := range managerCallbacks {
		callbacks[i] = types.HookCallback(callback)
	}

	return callbacks
}

/**
 * CallHook executes all callbacks registered for the specified hook.
 *
 * @param name string The name of the hook to call
 * @param args ...interface{} Arguments to pass to the hook callbacks
 * @return []interface{} Results from all callback executions
 * @return error Any error that occurred during hook execution
 */
func (t *Hookable) CallHook(name string, args ...interface{}) ([]interface{}, error) {
	return t.manager.CallHook(name, args...)
}

/**
 * CallHookFirst executes callbacks until one returns a non-nil result.
 *
 * @param name string The name of the hook to call
 * @param args ...interface{} Arguments to pass to the hook callbacks
 * @return interface{} The first non-nil result from callbacks
 * @return error Any error that occurred during hook execution
 */
func (t *Hookable) CallHookFirst(name string, args ...interface{}) (interface{}, error) {
	return t.manager.CallHookFirst(name, args...)
}

/**
 * CallHookUntil executes callbacks until one returns true.
 *
 * @param name string The name of the hook to call
 * @param args ...interface{} Arguments to pass to the hook callbacks
 * @return bool true if any callback returned true
 * @return error Any error that occurred during hook execution
 */
func (t *Hookable) CallHookUntil(name string, args ...interface{}) (bool, error) {
	return t.manager.CallHookUntil(name, args...)
}

/**
 * GetHookCount returns the number of callbacks registered for a hook.
 *
 * @param name string The name of the hook
 * @return int The number of registered callbacks
 */
func (t *Hookable) GetHookCount(name string) int {
	return t.manager.GetHookCount(name)
}

/**
 * GetAllHookNames returns all registered hook names.
 *
 * @return []string Slice of all hook names
 */
func (t *Hookable) GetAllHookNames() []string {
	return t.manager.GetAllHookNames()
}

/**
 * ClearHooks removes all registered hooks.
 */
func (t *Hookable) ClearHooks() {
	t.manager.ClearHooks()
}

/**
 * GetHooksInfo returns comprehensive hooks information.
 *
 * @return map[string]interface{} Hooks details
 */
func (t *Hookable) GetHooksInfo() map[string]interface{} {
	return t.manager.GetHooksInfo()
}

/**
 * GetManager returns the underlying hooks manager instance.
 * This allows direct access to manager functionality when needed.
 *
 * @return *hooks.Manager The underlying manager
 */
func (t *Hookable) GetManager() *hooks.Manager {
	return t.manager
}

/**
 * SetManager sets a new underlying hooks manager instance.
 *
 * @param m *hooks.Manager The new manager to wrap
 */
func (t *Hookable) SetManager(m *hooks.Manager) {
	t.manager = m
}

// Compile-time interface compliance check
var _ traitInterfaces.HookableInterface = (*Hookable)(nil)
