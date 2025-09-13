package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for database package
var (
	// DATABASE_TOKEN is the main service token for database
	DATABASE_TOKEN = symbol.For("govel.database")
	
	// DATABASE_FACTORY_TOKEN is the factory token for database
	DATABASE_FACTORY_TOKEN = symbol.For("govel.database.factory")
	
	// DATABASE_MANAGER_TOKEN is the manager token for database
	DATABASE_MANAGER_TOKEN = symbol.For("govel.database.manager")
	
	// DATABASE_INTERFACE_TOKEN is the interface token for database
	DATABASE_INTERFACE_TOKEN = symbol.For("govel.database.interface")
	
	// DATABASE_CONFIG_TOKEN is the config token for database
	DATABASE_CONFIG_TOKEN = symbol.For("govel.database.config")
	
	// DB_TOKEN is the database token (alias)
	DB_TOKEN = symbol.For("govel.database")
)
