package types

// ProviderCallback defines the signature for provider callback functions.
// These callbacks are used for pre-boot and post-boot operations in service providers,
// following Laravel's service provider callback pattern.
//
// Parameters:
//   app: The application instance for accessing services and configuration (interface{})
//
// Returns:
//   error: Any error that occurred during callback execution
type ProviderCallback func(app interface{}) error
