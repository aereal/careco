//go:build wireinject

package main

import (
	"careco/backend/config"
	"careco/backend/config/providers"
	"careco/backend/web"

	"github.com/google/wire"
)

func build() *web.Server {
	wire.Build(
		config.ProvideEnvironment,
		providers.ProvidePort,
		web.ProvideServer,
	)
	return nil
}
