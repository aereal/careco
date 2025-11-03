package providers

import (
	"careco/backend/config"
	"careco/backend/web"
)

func ProvidePort(e *config.Environment) web.Port {
	cast := config.Cast(config.StringAs[web.Port])
	retrieve := cast(config.EnvSource(e))
	return config.Yield(
		"WEB_PORT",
		config.WithDefaultValue[web.Port]("8080")(retrieve),
	)
}
