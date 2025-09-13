package interfaces

import "govel/support/symbol"


// Standard tokens for storage package
var (
	// STORAGE_TOKEN is the main service token for storage
	STORAGE_TOKEN = symbol.For("govel.storage")

	// STORAGE_FACTORY_TOKEN is the factory token for storage
	STORAGE_FACTORY_TOKEN = symbol.For("govel.storage.factory")

	// STORAGE_MANAGER_TOKEN is the manager token for storage
	STORAGE_MANAGER_TOKEN = symbol.For("govel.storage.manager")

	// STORAGE_INTERFACE_TOKEN is the interface token for storage
	STORAGE_INTERFACE_TOKEN = symbol.For("govel.storage.interface")

	// STORAGE_CONFIG_TOKEN is the config token for storage
	STORAGE_CONFIG_TOKEN = symbol.For("govel.storage.config")
)
