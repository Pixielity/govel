package interfaces

import "govel/support/symbol"


// Standard tokens for queue package
var (
	// QUEUE_TOKEN is the main service token for queue
	QUEUE_TOKEN = symbol.For("govel.queue")

	// QUEUE_FACTORY_TOKEN is the factory token for queue
	QUEUE_FACTORY_TOKEN = symbol.For("govel.queue.factory")

	// QUEUE_MANAGER_TOKEN is the manager token for queue
	QUEUE_MANAGER_TOKEN = symbol.For("govel.queue.manager")

	// QUEUE_INTERFACE_TOKEN is the interface token for queue
	QUEUE_INTERFACE_TOKEN = symbol.For("govel.queue.interface")

	// QUEUE_CONFIG_TOKEN is the config token for queue
	QUEUE_CONFIG_TOKEN = symbol.For("govel.queue.config")
)
