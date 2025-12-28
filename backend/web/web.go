package web

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"careco/backend/authz"
	"careco/backend/log/attribute"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Port string

func ProvideServer(port Port, tp trace.TracerProvider, gh *handler.Server, authMiddleware authz.Middleware, lh slog.Handler) *Server {
	return &Server{
		port:   port,
		tp:     tp,
		gh:     gh,
		authMW: authMiddleware,
		lh:     lh,
	}
}

type Server struct {
	port   Port
	tp     trace.TracerProvider
	gh     *handler.Server
	authMW authz.Middleware
	lh     slog.Handler
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
	withOtel := otelhttp.NewMiddleware("",
		otelhttp.WithPropagators(propagation.TraceContext{}),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string { return r.Method + " " + r.URL.Path }),
		otelhttp.WithTracerProvider(s.tp),
		otelhttp.WithFilter(func(r *http.Request) bool { return r.Method != http.MethodOptions }),
	)
	mux := http.NewServeMux()
	withCors := WithCors(slog.New(s.lh))
	mux.Handle("/graphql", withOtel(withCors(s.authMW(s.gh))))
	mux.Handle("GET /allow-cors", withOtel(withCors(s.authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"allow-cors":true,"auth":true}`))
	})))))
	mux.Handle("GET /no-cors", withOtel(s.authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"allow-cors":false,"auth":true}`))
	}))))
	mux.Handle("GET /no-auth", withOtel(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"allow-cors":false,"auth":false}`))
	})))
	return mux
}

func WithCors(l *slog.Logger) func(http.Handler) http.Handler {
	opts := cors.Options{
		AllowOriginVaryRequestFunc: func(r *http.Request, origin string) (bool, []string) {
			url, err := neturl.Parse(origin)
			if err != nil {
				slog.WarnContext(r.Context(), "failed to parse origin", slog.String("origin", origin), attribute.Error(err))
				return false, nil
			}
			return isLocalhost(url) || isOwnVercelURL(url), []string{"origin"}
		},
		AllowedMethods: []string{
			http.MethodPost, http.MethodGet, http.MethodHead,
		},
		AllowedHeaders: []string{"*"},
	}
	return cors.New(opts).Handler
}

func isLocalhost(url *neturl.URL) bool {
	return url.Scheme == "http" && url.Hostname() == "localhost"
}

func isOwnVercelURL(url *neturl.URL) bool {
	return url.Scheme == "https" && (isProductionVercelHost(url.Hostname()) || isPreviewVercelHost(url.Hostname()))
}

func isProductionVercelHost(host string) bool {
	return host == "careco-nine.vercel.app"
}

func isPreviewVercelHost(host string) bool {
	parts := strings.SplitN(host, ".", 2)
	if len(parts) < 2 {
		return false
	}
	if parts[1] != "vercel.app" {
		return false
	}
	return strings.HasPrefix(parts[0], "careco-") && strings.HasSuffix(parts[0], "-aereals-projects")
}
