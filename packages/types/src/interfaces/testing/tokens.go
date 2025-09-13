package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for testing package
var (
	// TESTING_TOKEN is the main service token for testing
	TESTING_TOKEN = symbol.For("govel.testing")

	// TESTING_FACTORY_TOKEN is the factory token for testing
	TESTING_FACTORY_TOKEN = symbol.For("govel.testing.factory")

	// TESTING_MANAGER_TOKEN is the manager token for testing
	TESTING_MANAGER_TOKEN = symbol.For("govel.testing.manager")

	// TESTING_INTERFACE_TOKEN is the interface token for testing
	TESTING_INTERFACE_TOKEN = symbol.For("govel.testing.interface")

	// TESTING_CONFIG_TOKEN is the config token for testing
	TESTING_CONFIG_TOKEN = symbol.For("govel.testing.config")
)
