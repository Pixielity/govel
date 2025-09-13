package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for hashing package
var (
	// HASHING_TOKEN is the main service token for hashing
	HASHING_TOKEN = symbol.For("govel.hashing")
	
	// HASHING_FACTORY_TOKEN is the factory token for hashing
	HASHING_FACTORY_TOKEN = symbol.For("govel.hashing.factory")
	
	// HASHING_MANAGER_TOKEN is the manager token for hashing
	HASHING_MANAGER_TOKEN = symbol.For("govel.hashing.manager")
	
	// HASHING_INTERFACE_TOKEN is the interface token for hashing
	HASHING_INTERFACE_TOKEN = symbol.For("govel.hashing.interface")
	
	// HASHING_CONFIG_TOKEN is the config token for hashing
	HASHING_CONFIG_TOKEN = symbol.For("govel.hashing.config")
	
	// HASH_DRIVER_TOKEN is the driver token for hashing
	HASH_DRIVER_TOKEN = symbol.For("govel.hashing.driver")
	
	// HASH_FACTORY_TOKEN is the factory token for hashing
	HASH_FACTORY_TOKEN = symbol.For("govel.hashing.hash_factory")
)
