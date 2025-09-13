package interfaces

import "govel/support/symbol"


// Standard tokens for event package
var (
	// EVENT_TOKEN is the main service token for event
	EVENT_TOKEN = symbol.For("govel.event")

	// EVENT_FACTORY_TOKEN is the factory token for event
	EVENT_FACTORY_TOKEN = symbol.For("govel.event.factory")

	// EVENT_MANAGER_TOKEN is the manager token for event
	EVENT_MANAGER_TOKEN = symbol.For("govel.event.manager")

	// EVENT_INTERFACE_TOKEN is the interface token for event
	EVENT_INTERFACE_TOKEN = symbol.For("govel.event.interface")

	// EVENT_CONFIG_TOKEN is the config token for event
	EVENT_CONFIG_TOKEN = symbol.For("govel.event.config")
)
