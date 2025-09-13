package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for auth package
var (
	// AUTH_TOKEN is the main service token for auth
	AUTH_TOKEN = symbol.For("govel.auth")

	// AUTH_FACTORY_TOKEN is the factory token for auth
	AUTH_FACTORY_TOKEN = symbol.For("govel.auth.factory")

	// AUTH_MANAGER_TOKEN is the manager token for auth
	AUTH_MANAGER_TOKEN = symbol.For("govel.auth.manager")

	// AUTH_INTERFACE_TOKEN is the interface token for auth
	AUTH_INTERFACE_TOKEN = symbol.For("govel.auth.interface")

	// AUTH_CONFIG_TOKEN is the config token for auth
	AUTH_CONFIG_TOKEN = symbol.For("govel.auth.config")
)
