//go:build wireinject

package main

import (
	"context"

	"careco/backend/entrypoint/server"

	"github.com/google/wire"
)

func build(_ context.Context) (*server.Entrypoint, error) {
	wire.Build(server.DevProvider)
	return nil, nil
}
