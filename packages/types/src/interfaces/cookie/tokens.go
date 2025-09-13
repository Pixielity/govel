package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for cookie package
var (
	// COOKIE_TOKEN is the main service token for cookie
	COOKIE_TOKEN = symbol.For("govel.cookie")
	
	// COOKIE_FACTORY_TOKEN is the factory token for cookie
	COOKIE_FACTORY_TOKEN = symbol.For("govel.cookie.factory")
	
	// COOKIE_MANAGER_TOKEN is the manager token for cookie
	COOKIE_MANAGER_TOKEN = symbol.For("govel.cookie.manager")
	
	// COOKIE_INTERFACE_TOKEN is the interface token for cookie
	COOKIE_INTERFACE_TOKEN = symbol.For("govel.cookie.interface")
	
	// COOKIE_CONFIG_TOKEN is the config token for cookie
	COOKIE_CONFIG_TOKEN = symbol.For("govel.cookie.config")
	
	// COOKIE_JAR_TOKEN is the jar token for cookie
	COOKIE_JAR_TOKEN = symbol.For("govel.cookie.jar")
)
