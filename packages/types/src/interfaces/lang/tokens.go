package interfaces

import "govel/support/symbol"


// Standard tokens for lang package
var (
	// LANG_TOKEN is the main service token for lang
	LANG_TOKEN = symbol.For("govel.lang")

	// LANG_FACTORY_TOKEN is the factory token for lang
	LANG_FACTORY_TOKEN = symbol.For("govel.lang.factory")

	// LANG_MANAGER_TOKEN is the manager token for lang
	LANG_MANAGER_TOKEN = symbol.For("govel.lang.manager")

	// LANG_INTERFACE_TOKEN is the interface token for lang
	LANG_INTERFACE_TOKEN = symbol.For("govel.lang.interface")

	// LANG_CONFIG_TOKEN is the config token for lang
	LANG_CONFIG_TOKEN = symbol.For("govel.lang.config")
)
