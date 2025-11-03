//go:build wireinject

package main

import (
	"careco/backend/cmd/config"
	"careco/backend/cmd/config/providers"
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
