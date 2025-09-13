package interfaces

import "govel/support/symbol"


// Standard tokens for session package
var (
	// SESSION_TOKEN is the main service token for session
	SESSION_TOKEN = symbol.For("govel.session")

	// SESSION_FACTORY_TOKEN is the factory token for session
	SESSION_FACTORY_TOKEN = symbol.For("govel.session.factory")

	// SESSION_MANAGER_TOKEN is the manager token for session
	SESSION_MANAGER_TOKEN = symbol.For("govel.session.manager")

	// SESSION_INTERFACE_TOKEN is the interface token for session
	SESSION_INTERFACE_TOKEN = symbol.For("govel.session.interface")

	// SESSION_CONFIG_TOKEN is the config token for session
	SESSION_CONFIG_TOKEN = symbol.For("govel.session.config")
)
