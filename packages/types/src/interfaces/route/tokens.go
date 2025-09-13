package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for route package
var (
	// ROUTE_TOKEN is the main service token for route
	ROUTE_TOKEN = symbol.For("govel.route")

	// ROUTE_FACTORY_TOKEN is the factory token for route
	ROUTE_FACTORY_TOKEN = symbol.For("govel.route.factory")

	// ROUTE_MANAGER_TOKEN is the manager token for route
	ROUTE_MANAGER_TOKEN = symbol.For("govel.route.manager")

	// ROUTE_INTERFACE_TOKEN is the interface token for route
	ROUTE_INTERFACE_TOKEN = symbol.For("govel.route.interface")

	// ROUTE_CONFIG_TOKEN is the config token for route
	ROUTE_CONFIG_TOKEN = symbol.For("govel.route.config")
)
