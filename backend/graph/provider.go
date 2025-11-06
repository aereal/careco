package graph

import (
	"careco/backend/graph/exec"
	"careco/backend/graph/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

func ProvideServer(r *resolver.Resolver) *handler.Server {
	es := exec.NewExecutableSchema(exec.Config{Resolvers: r})
	s := handler.New(es)
	s.AddTransport(transport.POST{})
	return s
}
