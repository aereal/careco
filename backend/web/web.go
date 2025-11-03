package web

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Port string

func ProvideServer(port Port, tp trace.TracerProvider) *Server {
	return &Server{
		port: port,
		tp:   tp,
	}
}

type Server struct {
	port Port
	tp   trace.TracerProvider
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
	withOtel := otelhttp.NewMiddleware("",
		otelhttp.WithPropagators(propagation.TraceContext{}),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string { return r.Method + " " + r.URL.Path }),
		otelhttp.WithTracerProvider(s.tp),
	)
	return withOtel(mux)
}
