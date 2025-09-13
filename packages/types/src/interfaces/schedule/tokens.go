package interfaces

import "govel/support/symbol"


// Standard tokens for schedule package
var (
	// SCHEDULE_TOKEN is the main service token for schedule
	SCHEDULE_TOKEN = symbol.For("govel.schedule")

	// SCHEDULE_FACTORY_TOKEN is the factory token for schedule
	SCHEDULE_FACTORY_TOKEN = symbol.For("govel.schedule.factory")

	// SCHEDULE_MANAGER_TOKEN is the manager token for schedule
	SCHEDULE_MANAGER_TOKEN = symbol.For("govel.schedule.manager")

	// SCHEDULE_INTERFACE_TOKEN is the interface token for schedule
	SCHEDULE_INTERFACE_TOKEN = symbol.For("govel.schedule.interface")

	// SCHEDULE_CONFIG_TOKEN is the config token for schedule
	SCHEDULE_CONFIG_TOKEN = symbol.For("govel.schedule.config")
)
