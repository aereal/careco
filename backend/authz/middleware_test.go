package authz_test

import (
	stdcmp "cmp"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"careco/backend/authz"

	"github.com/aereal/coll"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

var (
	expectedAudience = []string{
		"http://valid.rs.example",
	}
)

type addAuthorizationHeaderFunc func(issuerURL *url.URL, add func(value string)) error

func TestMiddleware(t *testing.T) {
	t.Parallel()

	signingPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		wantStatus     int
		wantSpans      tracetest.SpanStubs
		addAuthzHeader addAuthorizationHeaderFunc
	}{
		{
			name:       "ok",
			wantStatus: http.StatusOK,
			wantSpans:  expectedSpans.ok,
			addAuthzHeader: func(issuerURL *url.URL, add func(value string)) error {
				tok, err := createToken(signingPrivateKey, issuerURL.String())
				if err != nil {
					return err
				}
				add("Bearer " + tok)
				return nil
			},
		},
		{
			name:       "malformed token",
			wantStatus: http.StatusUnauthorized,
			wantSpans:  expectedSpans.malformedToken,
			addAuthzHeader: func(_ *url.URL, add func(value string)) error {
				add("Bearer " + "0xdeadbeaf")
				return nil
			},
		},
		{
			name:       "valid token but unexpected auth-scheme",
			wantStatus: http.StatusUnauthorized,
			wantSpans:  expectedSpans.unexpectedAuthScheme,
			addAuthzHeader: func(issuerURL *url.URL, add func(value string)) error {
				tok, err := createToken(signingPrivateKey, issuerURL.String())
				if err != nil {
					return err
				}
				add("unknown " + tok)
				return nil
			},
		},
		{
			name:           "no authz header",
			wantStatus:     http.StatusBadRequest,
			wantSpans:      expectedSpans.noHeader,
			addAuthzHeader: func(_ *url.URL, _ func(value string)) error { return nil },
		},
		{
			name:       "token from another issuer",
			wantStatus: http.StatusUnauthorized,
			wantSpans:  expectedSpans.tokenFromAnotherIssuer,
			addAuthzHeader: func(_ *url.URL, add func(value string)) error {
				tok, err := createToken(signingPrivateKey, "http://another-issuer.example/")
				if err != nil {
					return err
				}
				add("Bearer " + tok)
				return nil
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issuerHandler, err := newIssuerServer(signingPrivateKey)
			if err != nil {
				t.Fatal(err)
			}
			issuerSrv := httptest.NewServer(issuerHandler)
			t.Cleanup(issuerSrv.Close)
			issuerURL, err := url.Parse(issuerSrv.URL)
			if err != nil {
				t.Fatal(err)
			}
			issuerHandler.origin = issuerURL.String()
			mw := authz.ProvideMiddleware((*authz.Issuer)(issuerURL), expectedAudience, http.DefaultClient, coll.NewSet[authz.AllowedSubject]())
			exporter := tracetest.NewInMemoryExporter()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
			withOtel := otelhttp.NewMiddleware("default", otelhttp.WithTracerProvider(tp))
			appSrv := httptest.NewServer(withOtel(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
			t.Cleanup(appSrv.Close)

			if err := sendRequest(t.Context(), appSrv.URL, issuerURL, tc.wantStatus, tc.addAuthzHeader); err != nil {
				t.Error(err)
			}

			if err := tp.ForceFlush(t.Context()); err != nil {
				t.Error(err)
			}
			gotSpans := exporter.GetSpans()
			if diff := cmpSpans(tc.wantSpans, gotSpans); diff != "" {
				t.Errorf("(-want, +got):\n%s", diff)
			}
		})
	}
}

func TestMiddleware_cache(t *testing.T) {
	t.Parallel()

	signingPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	issuerHandler, err := newIssuerServer(signingPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	issuerSrv := httptest.NewServer(issuerHandler)
	t.Cleanup(issuerSrv.Close)
	issuerURL, err := url.Parse(issuerSrv.URL)
	if err != nil {
		t.Fatal(err)
	}
	issuerHandler.origin = issuerURL.String()
	mw := authz.ProvideMiddleware((*authz.Issuer)(issuerURL), expectedAudience, http.DefaultClient, coll.NewSet[authz.AllowedSubject]())
	appSrv := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	t.Cleanup(appSrv.Close)

	addAuthz := addAuthorizationHeaderFunc(func(issuerURL *url.URL, add func(value string)) error {
		tok, err := createToken(signingPrivateKey, issuerURL.String())
		if err != nil {
			return err
		}
		add("Bearer " + tok)
		return nil
	})
	for range 3 {
		if err := sendRequest(t.Context(), appSrv.URL, issuerURL, http.StatusOK, addAuthz); err != nil {
			t.Error(err)
		}
	}
	if actual := issuerHandler.callCount.wellKnownOIDCConfig.Load(); actual != 1 {
		t.Errorf("request to /.well-known/openid-configuration must be cached but not: count=%d", actual)
	}
	if actual := issuerHandler.callCount.jwks.Load(); actual != 1 {
		t.Errorf("request to JWKs endpoint must be cached but not: count=%d", actual)
	}
}

func TestMiddleware_specify_subject(t *testing.T) {
	t.Parallel()

	signingPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	issuerHandler, err := newIssuerServer(signingPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	issuerSrv := httptest.NewServer(issuerHandler)
	t.Cleanup(issuerSrv.Close)
	issuerURL, err := url.Parse(issuerSrv.URL)
	if err != nil {
		t.Fatal(err)
	}
	issuerHandler.origin = issuerURL.String()
	mw := authz.ProvideMiddleware(
		(*authz.Issuer)(issuerURL),
		expectedAudience,
		http.DefaultClient,
		coll.NewSet[authz.AllowedSubject]("sub1", "sub2"),
	)
	appSrv := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	t.Cleanup(appSrv.Close)

	testCases := []struct {
		sub        string
		wantStatus int
	}{
		{sub: "sub1", wantStatus: http.StatusOK},
		{sub: "sub2", wantStatus: http.StatusOK},
		{sub: "other_sub", wantStatus: http.StatusUnauthorized},
	}
	for _, tc := range testCases {
		t.Run(tc.sub, func(t *testing.T) {
			t.Parallel()

			addAuthz := addAuthorizationHeaderFunc(func(issuerURL *url.URL, add func(value string)) error {
				tok, err := createToken(signingPrivateKey, issuerURL.String(), withSubject(tc.sub))
				if err != nil {
					return err
				}
				add("Bearer " + tok)
				return nil
			})
			if err := sendRequest(t.Context(), appSrv.URL, issuerURL, tc.wantStatus, addAuthz); err != nil {
				t.Error(err)
			}
		})
	}
}

func sendRequest(ctx context.Context, appSrvURL string, issuerURL *url.URL, wantStatus int, addFunc addAuthorizationHeaderFunc) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appSrvURL, nil)
	if err != nil {
		return err
	}
	if err := addFunc(issuerURL, func(value string) { req.Header.Set("Authorization", value) }); err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != wantStatus {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("response.status=%d body=%s", resp.StatusCode, string(body)) //nolint:err113 // just a test
	}
	return nil
}

func newIssuerServer(privKey *rsa.PrivateKey) (*issuerServer, error) {
	key, err := jwk.Import(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("jwk.Import: %w", err)
	}
	if err := key.Set(jwk.AlgorithmKey, jwa.RS256()); err != nil {
		return nil, fmt.Errorf(`jwk.Key.Set("alg"): %w`, err)
	}
	if err := key.Set(jwk.KeyIDKey, "k1"); err != nil {
		return nil, fmt.Errorf(`jwk.Key.Set("kid"): %w`, err)
	}
	if err := key.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return nil, fmt.Errorf(`jwk.Key.Set("use"): %w`, err)
	}
	set := jwk.NewSet()
	if err := set.AddKey(key); err != nil {
		return nil, fmt.Errorf("jwk.Set.AddKey: %w", err)
	}

	s := &issuerServer{mux: http.NewServeMux()}
	s.mux.Handle("GET /.well-known/openid-configuration", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.callCount.wellKnownOIDCConfig.Add(1)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"jwks_uri": s.origin + "/jwks.json"})
	}))
	s.mux.Handle("GET /jwks.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.callCount.jwks.Add(1)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(set)
	}))
	return s, nil
}

type issuerServer struct {
	origin    string
	mux       *http.ServeMux
	callCount struct {
		wellKnownOIDCConfig atomic.Int32
		jwks                atomic.Int32
	}
}

func (s *issuerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type createTokenOption func(b *jwt.Builder) *jwt.Builder

func withSubject(sub string) createTokenOption {
	return func(b *jwt.Builder) *jwt.Builder { return b.Subject(sub) }
}

func createToken(signingPrivateKey *rsa.PrivateKey, issuerURL string, opts ...createTokenOption) (string, error) {
	issuedAt := time.Now()
	b := jwt.NewBuilder().
		Issuer(issuerURL).
		Subject("1234567890").
		Audience(expectedAudience).
		IssuedAt(issuedAt).
		Expiration(issuedAt.Add(time.Hour * 24))
	for _, o := range opts {
		b = o(b)
	}
	tok, err := b.Build()
	if err != nil {
		return "", fmt.Errorf("jwt.Builder.Build: %w", err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256(), signingPrivateKey))
	if err != nil {
		return "", fmt.Errorf("jwt.Sign: %w", err)
	}
	return string(signed), nil
}

var (
	expectedSpans = struct {
		ok, malformedToken, unexpectedAuthScheme, tokenFromAnotherIssuer, noHeader tracetest.SpanStubs
	}{
		ok: tracetest.SpanStubs{
			{
				Name:     "FetchKeys",
				SpanKind: trace.SpanKindInternal,
				Status:   sdktrace.Status{Code: codes.Ok},
				Attributes: []attribute.KeyValue{
					attribute.Bool("cache.hit", false),
				},
			},
			{
				Name:     "GET",
				SpanKind: trace.SpanKindServer,
				Attributes: append(commonRequestSpanAttrs,
					attribute.Int64("http.response.status_code", 200),
					attribute.StringSlice("creds.audience", []string{"http://valid.rs.example"}),
					attribute.Float64("creds.expires_at.remaining", valueReplacedFloat),
					attribute.String("creds.expires_at.timestamp", valueReplacedStr),
					attribute.Float64("creds.issued_at.elapsed", valueReplacedFloat),
					attribute.String("creds.issued_at.timestamp", valueReplacedStr),
					attribute.String("creds.issuer", valueReplacedStr),
					attribute.String("creds.subject", "1234567890"),
				),
			},
		},
		malformedToken: tracetest.SpanStubs{
			{
				Name:     "GET",
				SpanKind: trace.SpanKindServer,
				Attributes: append(commonRequestSpanAttrs,
					attribute.Int64("http.response.body.size", 31),
					attribute.Int64("http.response.status_code", 401),
				),
				Status: sdktrace.Status{Code: codes.Error, Description: `jwt.ParseString: failed to parse string: unknown payload type (payload is not JWT?)`},
				Events: []sdktrace.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.message", `jwt.ParseString: failed to parse string: unknown payload type (payload is not JWT?)`),
							attribute.String("exception.type", "github.com/lestrrat-go/jwx/v3/jwt/internal/errors.ParseError"),
						},
					},
				},
			},
		},
		unexpectedAuthScheme: tracetest.SpanStubs{
			{
				Name:     "GET",
				SpanKind: trace.SpanKindServer,
				Attributes: append(commonRequestSpanAttrs,
					attribute.Int64("http.response.body.size", 31),
					attribute.Int64("http.response.status_code", 401),
				),
				Status: sdktrace.Status{Code: codes.Error, Description: `unexpected auth scheme: "unknown"`},
				Events: []sdktrace.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.message", `unexpected auth scheme: "unknown"`),
							attribute.String("exception.type", "*authz.UnexpectedAuthSchemeError"),
						},
					},
				},
			},
		},
		tokenFromAnotherIssuer: tracetest.SpanStubs{
			{
				Name:     "FetchKeys",
				SpanKind: trace.SpanKindInternal,
				Status:   sdktrace.Status{Code: codes.Ok},
				Attributes: []attribute.KeyValue{
					attribute.Bool("cache.hit", false),
				},
			},
			{
				Name:     "GET",
				SpanKind: trace.SpanKindServer,
				Attributes: append(commonRequestSpanAttrs,
					attribute.Int64("http.response.body.size", 31),
					attribute.Int64("http.response.status_code", 401),
				),
				Status: sdktrace.Status{Code: codes.Error, Description: `jwt.ParseString: failed to parse string: jwt.Validate: validation failed: "iss" not satisfied: claim "iss" does not have the expected value`},
				Events: []sdktrace.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.message", `jwt.ParseString: failed to parse string: jwt.Validate: validation failed: "iss" not satisfied: claim "iss" does not have the expected value`),
							attribute.String("exception.type", "github.com/lestrrat-go/jwx/v3/jwt/internal/errors.ParseError"),
						},
					},
				},
			},
		},
		noHeader: tracetest.SpanStubs{
			{
				Name:     "GET",
				SpanKind: trace.SpanKindServer,
				Attributes: append(commonRequestSpanAttrs,
					attribute.Int64("http.response.body.size", 35),
					attribute.Int64("http.response.status_code", 400),
				),
				Status: sdktrace.Status{Code: codes.Error, Description: `missing token`},
				Events: []sdktrace.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.message", `missing token`),
							attribute.String("exception.type", "careco/backend/authz.MissingTokenError"),
						},
					},
				},
			},
		},
	}
	commonRequestSpanAttrs = []attribute.KeyValue{
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
)

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
		"creds.issuer",
		"creds.issued_at.timestamp", "creds.expires_at.timestamp":
		return map[attribute.Key]any{kv.Key: valueReplacedStr}
	case "creds.issued_at.elapsed", "creds.expires_at.remaining":
		return map[attribute.Key]any{kv.Key: valueReplacedFloat}
	}
	return map[attribute.Key]any{kv.Key: kv.Value.AsInterface()}
}
