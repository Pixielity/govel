package interfaces

import "govel/support/symbol"


// Standard tokens for mail package
var (
	// MAIL_TOKEN is the main service token for mail
	MAIL_TOKEN = symbol.For("govel.mail")

	// MAIL_FACTORY_TOKEN is the factory token for mail
	MAIL_FACTORY_TOKEN = symbol.For("govel.mail.factory")

	// MAIL_MANAGER_TOKEN is the manager token for mail
	MAIL_MANAGER_TOKEN = symbol.For("govel.mail.manager")

	// MAIL_INTERFACE_TOKEN is the interface token for mail
	MAIL_INTERFACE_TOKEN = symbol.For("govel.mail.interface")

	// MAIL_CONFIG_TOKEN is the config token for mail
	MAIL_CONFIG_TOKEN = symbol.For("govel.mail.config")
)
