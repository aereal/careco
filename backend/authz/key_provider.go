package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"sync"

	"careco/backend/o11y/traceutils"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func newOIDCKeyProvider(ksFetcher *oidcKeyFetcher) jws.KeyProvider {
	return &oidcKeyProvider{keySetFetcher: ksFetcher}
}

type oidcKeyProvider struct {
	keySetFetcher *oidcKeyFetcher
}

var _ jws.KeyProvider = (*oidcKeyProvider)(nil)

func (p *oidcKeyProvider) FetchKeys(ctx context.Context, sink jws.KeySink, sig *jws.Signature, _ *jws.Message) (err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("careco/backend/authz.oidcKeyProvider").
		Start(ctx, "FetchKeys")
	defer func() { traceutils.FinishSpan(span, err) }()
	algInSignature, ok := sig.ProtectedHeaders().Algorithm()
	if !ok {
		return NoSignatureAlgorithmError{}
	}

	ks, cacheHit, err := p.keySetFetcher.fetchKeySet(ctx)
	span.SetAttributes(attribute.Bool("cache.hit", cacheHit))
	if err != nil {
		return err
	}
	for alg, key := range iterateSignableKeys(ks) {
		if algInSignature != alg {
			continue
		}
		sink.Key(alg, key)
	}
	return nil
}

func iterateSignableKeys(keySet jwk.Set) iter.Seq2[jwa.SignatureAlgorithm, jwk.Key] {
	return func(yield func(jwa.SignatureAlgorithm, jwk.Key) bool) {
		for idx := range keySet.Len() {
			key, ok := keySet.Key(idx)
			if !ok {
				continue
			}
			alg, ok := key.Algorithm()
			if !ok {
				continue
			}
			sigAlg, ok := alg.(jwa.SignatureAlgorithm)
			if !ok {
				continue
			}
			if !yield(sigAlg, key) {
				return
			}
		}
	}
}

func wrapKeyProvider(ctx context.Context, kp jws.KeyProvider) jws.KeyProvider {
	return &contextualKeyProvider{ctx: ctx, kp: kp}
}

// contextualKeyProvider wraps a KeyProvider and replaces the context with the stored one.
// This is necessary because jwt.ParseString does not propagate the context to FetchKeys.
//
//nolint:containedctx // We need to store context to work around jwx library limitation
type contextualKeyProvider struct {
	kp  jws.KeyProvider
	ctx context.Context
}

func (c *contextualKeyProvider) FetchKeys(_ context.Context, sink jws.KeySink, sig *jws.Signature, msg *jws.Message) error {
	// Use the stored context instead of the one passed by jwx.
	// This is intentional to work around jwt.ParseString not propagating context.
	//nolint:contextcheck // We intentionally replace the context here
	return c.kp.FetchKeys(c.ctx, sink, sig, msg)
}

func newOIDCKeySetFetcher(client *http.Client, issuer *Issuer) *oidcKeyFetcher {
	return &oidcKeyFetcher{
		client: client,
		issuer: issuer,
	}
}

type oidcKeyFetcher struct {
	client *http.Client
	issuer *Issuer

	result jwk.Set
	err    error
	mux    sync.Mutex
	loaded bool
}

func (f *oidcKeyFetcher) fetchJwskURIFromOIDCConfiguration(ctx context.Context) (_ string, err error) {
	type wkEndpoints struct {
		JWKSURI string `json:"jwks_uri"`
	}
	body, err := request(ctx, f.client, f.issuer.oidcConfiguration().String(), func(r io.Reader) (*wkEndpoints, error) {
		var endpoints wkEndpoints
		if err := json.NewDecoder(r).Decode(&endpoints); err != nil {
			return nil, err
		}
		return &endpoints, nil
	})
	if err != nil {
		return "", err
	}
	return body.JWKSURI, nil
}

func (f *oidcKeyFetcher) fetchKeySet(ctx context.Context) (jwk.Set, bool, error) {
	f.mux.Lock()
	defer f.mux.Unlock()

	if f.loaded {
		return f.result, true, f.err
	}
	jwksURI, err := f.fetchJwskURIFromOIDCConfiguration(ctx)
	if err != nil {
		return nil, false, err
	}
	keySet, err := request(ctx, f.client, jwksURI, func(r io.Reader) (jwk.Set, error) { return jwk.ParseReader(r) })
	if err != nil {
		f.err = err
	} else {
		f.result = keySet
	}
	f.loaded = true
	return f.result, false, f.err
}

func request[T any](ctx context.Context, client *http.Client, uri string, parseBody func(r io.Reader) (T, error)) (T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return *new(T), fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return *new(T), fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return *new(T), &HTTPResponseError{Status: resp.StatusCode}
	}
	return parseBody(resp.Body)
}
