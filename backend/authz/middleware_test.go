package authz_test

import (
	stdcmp "cmp"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"careco/backend/authz"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	expectedAudience := []string{"http://valid.rs.example"}
	signingPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: signingPrivateKey},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", "k1"),
	)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		wantSpans      tracetest.SpanStubs
		wantStatusCode int
		createToken    func(issuerURL *url.URL) (string, error)
	}{
		{
			name:           "ok",
			wantSpans:      expectedSpans.ok,
			wantStatusCode: http.StatusOK,
			createToken: func(issuerURL *url.URL) (string, error) {
				issuedAt := time.Now()
				return jwt.Signed(signer).
					Claims(validator.RegisteredClaims{
						Issuer:   issuerURL.String(),
						Subject:  "1234567890",
						Audience: expectedAudience,
						IssuedAt: issuedAt.Unix(),
						Expiry:   issuedAt.Add(time.Hour * 24).Unix(),
					}).
					CompactSerialize()
			},
		},
		{
			name:           "token from another issuer",
			wantSpans:      expectedSpans.tokenFromAnotherIssuer,
			wantStatusCode: http.StatusUnauthorized,
			createToken: func(_ *url.URL) (string, error) {
				issuedAt := time.Now()
				return jwt.Signed(signer).
					Claims(validator.RegisteredClaims{
						Issuer:   "http://another-issuer.test",
						Subject:  "1234567890",
						Audience: expectedAudience,
						IssuedAt: issuedAt.Unix(),
						Expiry:   issuedAt.Add(time.Hour * 24).Unix(),
					}).
					CompactSerialize()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issuerHandler := newIssuerServer(signingPrivateKey)
			issuerSrv := httptest.NewServer(issuerHandler)
			t.Cleanup(issuerSrv.Close)
			issuerURL, err := url.Parse(issuerSrv.URL)
			if err != nil {
				t.Fatal(err)
			}
			issuerHandler.origin = issuerURL.String()
			mw, err := authz.ProvideAuth0Middleware((*authz.Issuer)(issuerURL), expectedAudience, http.DefaultClient)
			if err != nil {
				t.Fatal(err)
			}
			exporter := tracetest.NewInMemoryExporter()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
			withOtel := otelhttp.NewMiddleware("default", otelhttp.WithTracerProvider(tp))
			appSrv := httptest.NewServer(withOtel(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
			t.Cleanup(appSrv.Close)

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, appSrv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tok, err := tc.createToken(issuerURL); err == nil {
				req.Header.Set("authorization", "Bearer "+tok)
			} else {
				t.Error(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if err := tp.ForceFlush(t.Context()); err != nil {
				t.Errorf("TracerProvider.ForceFlush: %s", err)
			}
			gotSpans := exporter.GetSpans()
			if diff := cmpSpans(tc.wantSpans, gotSpans); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}

			if resp.StatusCode != tc.wantStatusCode {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("response.status=%d body=%s", resp.StatusCode, string(body))
			}
		})
	}

}

var expectedSpans = struct {
	ok, tokenFromAnotherIssuer tracetest.SpanStubs
}{
	ok: tracetest.SpanStubs{
		{
			Name:     "TokenExtractor",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
		},
		{
			Name:     "KeyFunc",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
		},
		{
			Name:     "ValidateToken",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
		},
		{
			Name:     "Authenticate",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
			Attributes: []attribute.KeyValue{
				attribute.String("auth.credentials.issuer", valueReplacedStr),
				attribute.String("auth.credentials.subject", "1234567890"),
				attribute.StringSlice("auth.credentials.audience", []string{"http://valid.rs.example"}),
				attribute.String("auth.credentials.expires_at.iso8601", valueReplacedStr),
				attribute.Float64("auth.credentials.expires_at.remaining_seconds", valueReplacedFloat),
				attribute.String("auth.credentials.issued_at.iso8601", valueReplacedStr),
				attribute.Float64("auth.credentials.issued_at.elapsed_seconds", valueReplacedFloat),
			},
		},
		{
			Name:     "default",
			SpanKind: trace.SpanKindServer,
			Attributes: append(commonRequestSpanAttrs,
				attribute.Int64("http.response.status_code", 200),
			),
		},
	},
	tokenFromAnotherIssuer: tracetest.SpanStubs{
		{
			Name:     "TokenExtractor",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
		},
		{
			Name:     "KeyFunc",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Ok},
		},
		{
			Name:     "ValidateToken",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Error, Description: "expected claims not validated: go-jose/go-jose/jwt: validation failed, invalid issuer claim (iss)"},
			Events: []sdktrace.Event{
				{
					Name: "exception",
					Attributes: []attribute.KeyValue{
						attribute.String("exception.type", "*fmt.wrapError"),
						attribute.String("exception.message", "expected claims not validated: go-jose/go-jose/jwt: validation failed, invalid issuer claim (iss)"),
					},
				},
			},
		},
		{
			Name:     "Authenticate",
			SpanKind: trace.SpanKindServer,
			Status:   sdktrace.Status{Code: codes.Error, Description: "jwt invalid: expected claims not validated: go-jose/go-jose/jwt: validation failed, invalid issuer claim (iss)"},
			Events: []sdktrace.Event{
				{
					Name: "exception",
					Attributes: []attribute.KeyValue{
						attribute.String("exception.message", "jwt invalid: expected claims not validated: go-jose/go-jose/jwt: validation failed, invalid issuer claim (iss)"),
						attribute.String("exception.type", "*jwtmiddleware.invalidError"),
					},
				},
			},
		},
		{
			Name:     "default",
			SpanKind: trace.SpanKindServer,
			Attributes: append(commonRequestSpanAttrs,
				attribute.Int64("http.response.body.size", 29),
				attribute.Int64("http.response.status_code", 401),
			),
		},
	},
}

var commonRequestSpanAttrs = []attribute.KeyValue{
	attribute.String("client.address", "127.0.0.1"),
	attribute.String("http.request.method", "GET"),
	attribute.String("network.peer.address", "127.0.0.1"),
	attribute.String("network.peer.port", valueReplacedStr),
	attribute.String("network.protocol.version", "1.1"),
	attribute.String("server.address", "127.0.0.1"),
	attribute.String("server.port", valueReplacedStr),
	attribute.String("url.path", "/"),
	attribute.String("url.scheme", "http"),
	attribute.String("user_agent.original", "Go-http-client/1.1"),
}

func newIssuerServer(privKey *rsa.PrivateKey) *issuerServer {
	s := &issuerServer{mux: http.NewServeMux()}
	s.mux.Handle("GET /.well-known/openid-configuration", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"jwks_uri": s.origin + "/jwks.json"})
	}))
	s.mux.Handle("GET /jwks.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					Key:       &privKey.PublicKey,
					KeyID:     "k1",
					Algorithm: string(jose.RS256),
					Use:       "sig",
				},
			},
		})
	}))
	return s
}

type issuerServer struct {
	origin string
	mux    *http.ServeMux
}

func (s *issuerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func cmpSpans(want, got tracetest.SpanStubs) string {
	opts := []cmp.Option{
		cmp.Transformer("attribute.KeyValue", transformKeyValue),
		cmpopts.SortSlices(func(x, y tracetest.SpanStub) bool { return x.EndTime.Before(y.EndTime) }),
		cmpopts.SortSlices(func(x, y attribute.KeyValue) bool { return stdcmp.Less(x.Key, y.Key) }),
		cmpopts.IgnoreFields(sdktrace.Event{}, "Time"),
		cmpopts.IgnoreFields(
			tracetest.SpanStub{},
			"Parent", "SpanContext", "StartTime", "EndTime", "Links",
			"DroppedAttributes", "DroppedEvents", "DroppedLinks",
			"ChildSpanCount", "Resource", "InstrumentationLibrary", "InstrumentationScope",
		),
	}
	return cmp.Diff(want, got, opts...)
}

const (
	valueReplacedStr           = "<snip>"
	valueReplacedFloat float64 = 0.0
)

func transformKeyValue(kv attribute.KeyValue) map[attribute.Key]any {
	switch kv.Key {
	case "server.port", "network.peer.port",
		"auth.credentials.issuer",
		"auth.credentials.issued_at.iso8601", "auth.credentials.expires_at.iso8601":
		return map[attribute.Key]any{kv.Key: valueReplacedStr}
	case "auth.credentials.issued_at.elapsed_seconds", "auth.credentials.expires_at.remaining_seconds":
		return map[attribute.Key]any{kv.Key: valueReplacedFloat}
	}
	return map[attribute.Key]any{kv.Key: kv.Value.AsInterface()}
}
