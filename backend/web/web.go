package web

import (
	"context"
	"net"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Port string

func ProvideServer(port Port, tp trace.TracerProvider, gh *handler.Server) *Server {
	return &Server{
		port: port,
		tp:   tp,
		gh:   gh,
	}
}

type Server struct {
	port Port
	tp   trace.TracerProvider
	gh   *handler.Server
}

func (s *Server) Start(ctx context.Context) error {
	hs := &http.Server{
		Addr:              net.JoinHostPort("", string(s.port)),
		Handler:           s.handler(),
		ReadHeaderTimeout: time.Second * 3,
	}
	return start(ctx, hs, time.Second*5)
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /graphql", s.gh)
	withOtel := otelhttp.NewMiddleware("",
		otelhttp.WithPropagators(propagation.TraceContext{}),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string { return r.Method + " " + r.URL.Path }),
		otelhttp.WithTracerProvider(s.tp),
	)
	return withOtel(withCors(mux))
}

func withCors(next http.Handler) http.Handler {
	opts := cors.Options{
		AllowOriginFunc: func(origin string) bool {
			url, err := neturl.Parse(origin)
			if err != nil {
				return false
			}
			return url.Scheme == "http" && url.Hostname() == "localhost"
		},
		AllowedMethods: []string{
			http.MethodPost, http.MethodGet, http.MethodHead,
		},
		AllowedHeaders: []string{"*"},
	}
	return cors.New(opts).Handler(next)
}
