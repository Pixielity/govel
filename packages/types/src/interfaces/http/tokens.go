package interfaces

import "govel/support/symbol"


// Standard tokens for http package
var (
	// HTTP_TOKEN is the main service token for http
	HTTP_TOKEN = symbol.For("govel.http")

	// HTTP_FACTORY_TOKEN is the factory token for http
	HTTP_FACTORY_TOKEN = symbol.For("govel.http.factory")

	// HTTP_MANAGER_TOKEN is the manager token for http
	HTTP_MANAGER_TOKEN = symbol.For("govel.http.manager")

	// HTTP_INTERFACE_TOKEN is the interface token for http
	HTTP_INTERFACE_TOKEN = symbol.For("govel.http.interface")

	// HTTP_CONFIG_TOKEN is the config token for http
	HTTP_CONFIG_TOKEN = symbol.For("govel.http.config")
)
