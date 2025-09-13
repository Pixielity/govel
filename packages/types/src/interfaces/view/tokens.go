package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for view package
var (
	// VIEW_TOKEN is the main service token for view
	VIEW_TOKEN = symbol.For("govel.view")

	// VIEW_FACTORY_TOKEN is the factory token for view
	VIEW_FACTORY_TOKEN = symbol.For("govel.view.factory")

	// VIEW_MANAGER_TOKEN is the manager token for view
	VIEW_MANAGER_TOKEN = symbol.For("govel.view.manager")

	// VIEW_INTERFACE_TOKEN is the interface token for view
	VIEW_INTERFACE_TOKEN = symbol.For("govel.view.interface")

	// VIEW_CONFIG_TOKEN is the config token for view
	VIEW_CONFIG_TOKEN = symbol.For("govel.view.config")
)
