//go:build wireinject

package main

import (
	"context"

	"careco/backend/entrypoint/importdata"

	"github.com/google/wire"
)

func build(_ context.Context) (*importdata.Entrypoint, error) {
	wire.Build(
		importdata.Provider,
	)
	return nil, nil
}
