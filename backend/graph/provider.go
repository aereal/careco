package graph

import (
	"careco/backend/graph/exec"
	"careco/backend/graph/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/otelgqlgen"
	"go.opentelemetry.io/otel/trace"
)

func ProvideServer(tp trace.TracerProvider, r *resolver.Resolver) *handler.Server {
	es := exec.NewExecutableSchema(exec.Config{Resolvers: r})
	s := handler.New(es)
	s.AddTransport(transport.POST{})
	s.Use(otelgqlgen.New(otelgqlgen.WithTracerProvider(tp)))
	return s
}
