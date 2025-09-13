package interfaces

import "govel/support/symbol"


// Standard tokens for validation package
var (
	// VALIDATION_TOKEN is the main service token for validation
	VALIDATION_TOKEN = symbol.For("govel.validation")

	// VALIDATION_FACTORY_TOKEN is the factory token for validation
	VALIDATION_FACTORY_TOKEN = symbol.For("govel.validation.factory")

	// VALIDATION_MANAGER_TOKEN is the manager token for validation
	VALIDATION_MANAGER_TOKEN = symbol.For("govel.validation.manager")

	// VALIDATION_INTERFACE_TOKEN is the interface token for validation
	VALIDATION_INTERFACE_TOKEN = symbol.For("govel.validation.interface")

	// VALIDATION_CONFIG_TOKEN is the config token for validation
	VALIDATION_CONFIG_TOKEN = symbol.For("govel.validation.config")
)
