package interfaces

import "govel/support/symbol"


// Standard tokens for logger package
var (
	// LOGGER_TOKEN is the main service token for logger
	LOGGER_TOKEN = symbol.For("govel.logger")
	
	// LOGGER_FACTORY_TOKEN is the factory token for logger
	LOGGER_FACTORY_TOKEN = symbol.For("govel.logger.factory")
	
	// LOGGER_MANAGER_TOKEN is the manager token for logger
	LOGGER_MANAGER_TOKEN = symbol.For("govel.logger.manager")
	
	// LOGGER_INTERFACE_TOKEN is the interface token for logger
	LOGGER_INTERFACE_TOKEN = symbol.For("govel.logger.interface")
	
	// LOGGER_CONFIG_TOKEN is the config token for logger
	LOGGER_CONFIG_TOKEN = symbol.For("govel.logger.config")
	
	// LOGGER_LEVEL_TOKEN is the level token for logger
	LOGGER_LEVEL_TOKEN = symbol.For("govel.logger.level")
)
