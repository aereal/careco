//go:build wireinject

package test

import (
	"careco/backend/graph"
	"careco/backend/graph/resolver"

	"github.com/google/wire"
)

func BuildHandler() *TestHandler {
	wire.Build(
		graph.ProvideServer,
		provideTestServer,
		resolver.ProvideResolver,
	)
	return nil
}
