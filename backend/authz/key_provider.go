package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"

	"careco/backend/o11y/traceutils"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"go.opentelemetry.io/otel/trace"
)

func newOIDCKeyProvider(issuer *Issuer, client *http.Client) jws.KeyProvider {
	return &oidcKeyProvider{issuer: issuer, client: client}
}

type oidcKeyProvider struct {
	issuer *Issuer
	client *http.Client
}

var _ jws.KeyProvider = (*oidcKeyProvider)(nil)

func (p *oidcKeyProvider) FetchKeys(ctx context.Context, sink jws.KeySink, sig *jws.Signature, _ *jws.Message) (err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("careco/backend/authz.oidcKeyProvider").
		Start(ctx, "FetchKeys", trace.WithSpanKind(trace.SpanKindServer))
	defer func() { traceutils.FinishSpan(span, err) }()
	algInSignature, ok := sig.ProtectedHeaders().Algorithm()
	if !ok {
		return NoSignatureAlgorithmError{}
	}

	jwksURI, err := p.fetchJwskURIFromOIDCConfiguration(ctx)
	if err != nil {
		return fmt.Errorf("fetchJwskURIFromOIDCConfiguration: %w", err)
	}
	ks, err := jwk.Fetch(ctx, jwksURI, jwk.WithHTTPClient(p.client))
	if err != nil {
		return fmt.Errorf("jwk.Fetch: %w", err)
	}
	for alg, key := range iterateSignableKeys(ks) {
		if algInSignature != alg {
			continue
		}
		sink.Key(alg, key)
	}
	return nil
}

func (p *oidcKeyProvider) fetchJwskURIFromOIDCConfiguration(ctx context.Context) (_ string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.issuer.oidcConfiguration().String(), nil)
	if err != nil {
		return "", fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http.Client.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", &HTTPResponseError{Status: resp.StatusCode}
	}
	var wellKnownEndpoints struct {
		JWKSURI string `json:"jwks_uri"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wellKnownEndpoints); err != nil {
		return "", fmt.Errorf("json.Decoder.Decode: %w", err)
	}
	return wellKnownEndpoints.JWKSURI, nil
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
