package interfaces

import "govel/support/symbol"


// Standard tokens for encryption package
var (
	// ENCRYPTION_TOKEN is the main service token for encryption
	ENCRYPTION_TOKEN = symbol.For("govel.encryption")
	
	// ENCRYPTION_FACTORY_TOKEN is the factory token for encryption
	ENCRYPTION_FACTORY_TOKEN = symbol.For("govel.encryption.factory")
	
	// ENCRYPTION_MANAGER_TOKEN is the manager token for encryption
	ENCRYPTION_MANAGER_TOKEN = symbol.For("govel.encryption.manager")
	
	// ENCRYPTION_INTERFACE_TOKEN is the interface token for encryption
	ENCRYPTION_INTERFACE_TOKEN = symbol.For("govel.encryption.interface")
	
	// ENCRYPTION_CONFIG_TOKEN is the config token for encryption
	ENCRYPTION_CONFIG_TOKEN = symbol.For("govel.encryption.config")
	
	// ENCRYPTION_DRIVER_TOKEN is the driver token for encryption
	ENCRYPTION_DRIVER_TOKEN = symbol.For("govel.encryption.driver")
)
