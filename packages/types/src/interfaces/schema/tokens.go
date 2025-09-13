package interfaces

import "govel/support/symbol"


// Standard tokens for schema package
var (
	// SCHEMA_TOKEN is the main service token for schema
	SCHEMA_TOKEN = symbol.For("govel.schema")

	// SCHEMA_FACTORY_TOKEN is the factory token for schema
	SCHEMA_FACTORY_TOKEN = symbol.For("govel.schema.factory")

	// SCHEMA_MANAGER_TOKEN is the manager token for schema
	SCHEMA_MANAGER_TOKEN = symbol.For("govel.schema.manager")

	// SCHEMA_INTERFACE_TOKEN is the interface token for schema
	SCHEMA_INTERFACE_TOKEN = symbol.For("govel.schema.interface")

	// SCHEMA_CONFIG_TOKEN is the config token for schema
	SCHEMA_CONFIG_TOKEN = symbol.For("govel.schema.config")
)
