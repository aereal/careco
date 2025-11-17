package test

import (
	"careco/backend/config"
	"careco/backend/infra/firestore"
)

func provideEmulatorAddr(e *config.Environment) firestore.EmulatorAddr {
	cast := config.Cast(config.StringAs[firestore.EmulatorAddr])
	retrieve := cast(config.EnvSource(e))
	return config.Yield(
		"TEST_FIRESTORE_EMULATOR_ADDR",
		config.WithDefaultValue[firestore.EmulatorAddr]("localhost:8889")(retrieve),
	)
}
