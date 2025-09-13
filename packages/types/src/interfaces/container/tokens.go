package interfaces

import "govel/support/symbol"


// Standard tokens for container package
var (
	// CONTAINER_TOKEN is the main service token for container
	CONTAINER_TOKEN = symbol.For("govel.container")

	// CONTAINER_FACTORY_TOKEN is the factory token for container
	CONTAINER_FACTORY_TOKEN = symbol.For("govel.container.factory")

	// CONTAINER_MANAGER_TOKEN is the manager token for container
	CONTAINER_MANAGER_TOKEN = symbol.For("govel.container.manager")

	// CONTAINER_INTERFACE_TOKEN is the interface token for container
	CONTAINER_INTERFACE_TOKEN = symbol.For("govel.container.interface")

	// CONTAINER_CONFIG_TOKEN is the config token for container
	CONTAINER_CONFIG_TOKEN = symbol.For("govel.container.config")
	
	// CONTAINER_BINDINGS_TOKEN is the token for container bindings introspection
	CONTAINER_BINDINGS_TOKEN = symbol.For("govel.container.bindings")
	
	// CONTAINER_STATS_TOKEN is the token for container statistics
	CONTAINER_STATS_TOKEN = symbol.For("govel.container.stats")
)

// Additional package-specific tokens can be added below
