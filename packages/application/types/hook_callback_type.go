package types

/**
 * HookCallback represents a function that can be registered as a hook.
 *
 * @param args ...interface{} Variable arguments passed to the hook
 * @return interface{} The result of the hook execution
 * @return error Any error that occurred during hook execution
 */
type HookCallback func(args ...interface{}) (interface{}, error)
