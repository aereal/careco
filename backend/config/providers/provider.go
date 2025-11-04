package providers

import (
	"careco/backend/config"
	"careco/backend/web"
	"log/slog"
)

func ProvidePort(e *config.Environment) web.Port {
	cast := config.Cast(config.StringAs[web.Port])
	retrieve := cast(config.EnvSource(e))
	return config.Yield(
		"WEB_PORT",
		config.WithDefaultValue[web.Port]("8080")(retrieve),
	)
}

func ProvideLogLevel(e *config.Environment) slog.Level {
	cast := config.Cast(config.Unmarshal[slog.Level])
	retrieve := cast(config.EnvSource(e))
	withDefault := config.WithDefaultValue(slog.LevelInfo)
	return config.Yield(
		"LOG_LEVEL",
		withDefault(retrieve),
	)
}
