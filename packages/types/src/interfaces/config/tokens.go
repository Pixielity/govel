package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for config package
var (
	// CONFIG_TOKEN is the main service token for config
	CONFIG_TOKEN = symbol.For("govel.config")
	
	// CONFIG_FACTORY_TOKEN is the factory token for config
	CONFIG_FACTORY_TOKEN = symbol.For("govel.config.factory")
	
	// CONFIG_MANAGER_TOKEN is the manager token for config
	CONFIG_MANAGER_TOKEN = symbol.For("govel.config.manager")
	
	// CONFIG_INTERFACE_TOKEN is the interface token for config
	CONFIG_INTERFACE_TOKEN = symbol.For("govel.config.interface")
	
	// CONFIG_CONFIG_TOKEN is the config token for config
	CONFIG_CONFIG_TOKEN = symbol.For("govel.config.config")
	
	// CONFIG_ENVIRONMENT_TOKEN is the environment token for config
	CONFIG_ENVIRONMENT_TOKEN = symbol.For("govel.config.environment")
)
