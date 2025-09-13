package interfaces

import "govel/support/symbol"


// Standard tokens for cache package
var (
	// CACHE_TOKEN is the main service token for cache
	CACHE_TOKEN = symbol.For("govel.cache")

	// CACHE_FACTORY_TOKEN is the factory token for cache
	CACHE_FACTORY_TOKEN = symbol.For("govel.cache.factory")

	// CACHE_MANAGER_TOKEN is the manager token for cache
	CACHE_MANAGER_TOKEN = symbol.For("govel.cache.manager")

	// CACHE_INTERFACE_TOKEN is the interface token for cache
	CACHE_INTERFACE_TOKEN = symbol.For("govel.cache.interface")

	// CACHE_CONFIG_TOKEN is the config token for cache
	CACHE_CONFIG_TOKEN = symbol.For("govel.cache.config")
)
