package interfaces

import "govel/packages/application/types"

/**
 * HookableInterface defines the contract for components that provide
 * application hook management functionality. This interface follows the
 * Interface Segregation Principle by focusing solely on hook operations.
 */
type HookableInterface interface {
	/**
	 * RegisterHook registers a callback function for the specified hook name.
	 *
	 * @param name string The name of the hook
	 * @param priority int The priority of the callback (lower numbers execute first)
	 * @param callback HookCallback The callback function to register
	 */
	RegisterHook(name string, priority int, callback types.HookCallback)

	/**
	 * UnregisterHook removes all callbacks for the specified hook name.
	 *
	 * @param name string The name of the hook to unregister
	 * @return bool true if any callbacks were removed
	 */
	UnregisterHook(name string) bool

	/**
	 * UnregisterHookCallback removes a specific callback from a hook.
	 *
	 * @param name string The name of the hook
	 * @param callback HookCallback The specific callback to remove
	 * @return bool true if the callback was found and removed
	 */
	UnregisterHookCallback(name string, callback types.HookCallback) bool

	/**
	 * HasHook returns whether any callbacks are registered for the specified hook.
	 *
	 * @param name string The name of the hook to check
	 * @return bool true if the hook has registered callbacks
	 */
	HasHook(name string) bool

	/**
	 * GetHooks returns all registered hooks.
	 *
	 * @return map[string][]HookCallback Map of hook names to their callbacks
	 */
	GetHooks() map[string][]types.HookCallback

	/**
	 * GetHookCallbacks returns all callbacks for a specific hook.
	 *
	 * @param name string The name of the hook
	 * @return []HookCallback Slice of callbacks for the hook
	 */
	GetHookCallbacks(name string) []types.HookCallback

	/**
	 * CallHook executes all callbacks registered for the specified hook.
	 *
	 * @param name string The name of the hook to call
	 * @param args ...interface{} Arguments to pass to the hook callbacks
	 * @return []interface{} Results from all callback executions
	 * @return error Any error that occurred during hook execution
	 */
	CallHook(name string, args ...interface{}) ([]interface{}, error)

	/**
	 * CallHookFirst executes callbacks until one returns a non-nil result.
	 *
	 * @param name string The name of the hook to call
	 * @param args ...interface{} Arguments to pass to the hook callbacks
	 * @return interface{} The first non-nil result from callbacks
	 * @return error Any error that occurred during hook execution
	 */
	CallHookFirst(name string, args ...interface{}) (interface{}, error)

	/**
	 * CallHookUntil executes callbacks until one returns true.
	 *
	 * @param name string The name of the hook to call
	 * @param args ...interface{} Arguments to pass to the hook callbacks
	 * @return bool true if any callback returned true
	 * @return error Any error that occurred during hook execution
	 */
	CallHookUntil(name string, args ...interface{}) (bool, error)

	/**
	 * GetHookCount returns the number of callbacks registered for a hook.
	 *
	 * @param name string The name of the hook
	 * @return int The number of registered callbacks
	 */
	GetHookCount(name string) int

	/**
	 * GetAllHookNames returns all registered hook names.
	 *
	 * @return []string Slice of all hook names
	 */
	GetAllHookNames() []string

	/**
	 * ClearHooks removes all registered hooks.
	 */
	ClearHooks()

	/**
	 * GetHooksInfo returns comprehensive hooks information.
	 *
	 * @return map[string]interface{} Hooks details
	 */
	GetHooksInfo() map[string]interface{}
}
