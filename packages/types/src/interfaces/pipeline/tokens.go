package interfaces

import "govel/support/symbol"

// Standard tokens for pipeline package
var (
	// PIPELINE_TOKEN is the main service token for pipeline
	PIPELINE_TOKEN = symbol.For("govel.pipeline")

	// PIPELINE_FACTORY_TOKEN is the factory token for pipeline
	PIPELINE_FACTORY_TOKEN = symbol.For("govel.pipeline.factory")

	// PIPELINE_MANAGER_TOKEN is the manager token for pipeline
	PIPELINE_MANAGER_TOKEN = symbol.For("govel.pipeline.manager")

	// PIPELINE_INTERFACE_TOKEN is the interface token for pipeline
	PIPELINE_INTERFACE_TOKEN = symbol.For("govel.pipeline.interface")

	// PIPELINE_CONFIG_TOKEN is the config token for pipeline
	PIPELINE_CONFIG_TOKEN = symbol.For("govel.pipeline.config")

	// PIPELINE_HUB_TOKEN is the hub token for pipeline
	PIPELINE_HUB_TOKEN = symbol.For("govel.pipeline.hub")
)
