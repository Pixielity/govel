package interfaces

import "govel/packages/support/src/symbol"


// Standard tokens for orm package
var (
	// ORM_TOKEN is the main service token for orm
	ORM_TOKEN = symbol.For("govel.orm")

	// ORM_FACTORY_TOKEN is the factory token for orm
	ORM_FACTORY_TOKEN = symbol.For("govel.orm.factory")

	// ORM_MANAGER_TOKEN is the manager token for orm
	ORM_MANAGER_TOKEN = symbol.For("govel.orm.manager")

	// ORM_INTERFACE_TOKEN is the interface token for orm
	ORM_INTERFACE_TOKEN = symbol.For("govel.orm.interface")

	// ORM_CONFIG_TOKEN is the config token for orm
	ORM_CONFIG_TOKEN = symbol.For("govel.orm.config")
)
