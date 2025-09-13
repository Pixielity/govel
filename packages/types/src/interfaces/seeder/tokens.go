package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for seeder package
var (
	// SEEDER_TOKEN is the main service token for seeder
	SEEDER_TOKEN = symbol.For("govel.seeder")

	// SEEDER_FACTORY_TOKEN is the factory token for seeder
	SEEDER_FACTORY_TOKEN = symbol.For("govel.seeder.factory")

	// SEEDER_MANAGER_TOKEN is the manager token for seeder
	SEEDER_MANAGER_TOKEN = symbol.For("govel.seeder.manager")

	// SEEDER_INTERFACE_TOKEN is the interface token for seeder
	SEEDER_INTERFACE_TOKEN = symbol.For("govel.seeder.interface")

	// SEEDER_CONFIG_TOKEN is the config token for seeder
	SEEDER_CONFIG_TOKEN = symbol.For("govel.seeder.config")
)
