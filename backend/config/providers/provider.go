package providers

import (
	"log/slog"

	"careco/backend/authz"
	"careco/backend/config"
	"careco/backend/infra/firestore"
	"careco/backend/usecases/interactions"
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

func ProvideLogLevel(e *config.Environment) slog.Level {
	cast := config.Cast(config.Unmarshal[slog.Level])
	retrieve := cast(config.EnvSource(e))
	withDefault := config.WithDefaultValue(slog.LevelInfo)
	return config.Yield(
		"LOG_LEVEL",
		withDefault(retrieve),
	)
}

func ProvideFirestoreEmulatorAddr(e *config.Environment) firestore.EmulatorAddr {
	cast := config.Cast(config.StringAs[firestore.EmulatorAddr])
	retrieve := cast(config.EnvSource(e))
	return config.Yield(
		"FIRESTORE_EMULATOR_ADDR",
		config.WithDefaultValue[firestore.EmulatorAddr]("localhost:8888")(retrieve),
	)
}

func ProvideExportFileName(e *config.Environment) (interactions.ExportFileName, error) {
	cast := config.Cast(config.StringAs[interactions.ExportFileName])
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"EXPORT_FILE_NAME",
		retrieve,
	)
}

func ProvideAudience(e *config.Environment) (authz.Audience, error) {
	cast := config.Cast(parseAudience)
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"AUTH0_AUDIENCE",
		retrieve,
	)
}

func ProvideIssuer(e *config.Environment) (*authz.Issuer, error) {
	cast := config.Cast(parseIssuer)
	retrieve := cast(config.EnvSource(e))
	return config.Retrieve(
		"AUTH0_ISSUER",
		retrieve,
	)
}
