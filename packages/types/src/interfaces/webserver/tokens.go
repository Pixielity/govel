package interfaces

import "govel/support/symbol"


// Standard tokens for webserver package
var (
	// WEBSERVER_TOKEN is the main service token for webserver
	WEBSERVER_TOKEN = symbol.For("govel.webserver")
	
	// WEBSERVER_FACTORY_TOKEN is the factory token for webserver
	WEBSERVER_FACTORY_TOKEN = symbol.For("govel.webserver.factory")
	
	// WEBSERVER_MANAGER_TOKEN is the manager token for webserver
	WEBSERVER_MANAGER_TOKEN = symbol.For("govel.webserver.manager")
	
	// WEBSERVER_INTERFACE_TOKEN is the interface token for webserver
	WEBSERVER_INTERFACE_TOKEN = symbol.For("govel.webserver.interface")
	
	// WEBSERVER_CONFIG_TOKEN is the config token for webserver
	WEBSERVER_CONFIG_TOKEN = symbol.For("govel.webserver.config")
)

// Specific webserver service tokens
var (
	// WEBSERVER_MAIN_FACTORY_TOKEN is the main webserver factory
	WEBSERVER_MAIN_FACTORY_TOKEN = "webserver.factory"
	
	// WEBSERVER_ADAPTER_FACTORY_TOKEN is the adapter factory
	WEBSERVER_ADAPTER_FACTORY_TOKEN = "webserver.adapter.factory"
	
	// WEBSERVER_MIDDLEWARE_FACTORY_TOKEN is the middleware factory
	WEBSERVER_MIDDLEWARE_FACTORY_TOKEN = "webserver.middleware.factory"
	
	// WEBSERVER_CREATE_TOKEN is the webserver creation helper
	WEBSERVER_CREATE_TOKEN = "webserver.create"
	
	// WEBSERVER_DEFAULT_TOKEN is the default webserver helper
	WEBSERVER_DEFAULT_TOKEN = "webserver.default"
)

// Additional webserver-specific tokens can be added below