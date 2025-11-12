//go:build wireinject

package test

import (
	"careco/backend/graph"
	"careco/backend/graph/resolver"
	"careco/backend/o11y/test"

	"github.com/google/wire"
)

func BuildHandler() *TestHandler {
	wire.Build(
		graph.ProvideServer,
		provideTestServer,
		resolver.ProvideResolver,
		test.ProvideNoopTracerProvider,
	)
	return nil
}
