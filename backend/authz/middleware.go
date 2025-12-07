package authz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	validatorpkg "github.com/auth0/go-jwt-middleware/v2/validator"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Issuer url.URL

type Audience []string

type Middleware func(next http.Handler) http.Handler

func ProvideAuth0Middleware(issuer *Issuer, audience Audience, client *http.Client) (Middleware, error) {
	keyProvider := jwks.NewCachingProvider((*url.URL)(issuer), time.Minute*30, jwks.WithCustomClient(client))
	validator, err := validatorpkg.New(keyProvider.KeyFunc, validatorpkg.RS256, (*url.URL)(issuer).String(), audience)
	if err != nil {
		return nil, fmt.Errorf("validator.New: %w", err)
	}
	mw := jwtmiddleware.New(validator.ValidateToken)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				claims, ok := ctx.Value(jwtmiddleware.ContextKey{}).(*validatorpkg.ValidatedClaims)
				if !ok {
					next.ServeHTTP(w, r)
					return
				}

				span := trace.SpanFromContext(ctx)
				token := &auth0Token{claims}
				if id, ok := token.JwtID(); ok {
					span.SetAttributes(attribute.String("auth.credentials.id", id))
				}
				if issuer, ok := token.Issuer(); ok {
					span.SetAttributes(attribute.String("auth.credentials.issuer", issuer))
				}
				if sub, ok := token.Subject(); ok {
					span.SetAttributes(attribute.String("auth.credentials.subject", sub))
				}
				if aud, ok := token.Audience(); ok {
					span.SetAttributes(attribute.StringSlice("auth.credentials.audience", aud))
				}
				if expiresAt, ok := token.Expiration(); ok {
					span.SetAttributes(
						attribute.String("auth.credentials.expires_at.iso8601", expiresAt.Format(time.RFC3339)),
						attribute.Float64("auth.credentials.expires_at.remaining_seconds", time.Until(expiresAt).Seconds()),
					)
				}
				if issuedAt, ok := token.IssuedAt(); ok {
					span.SetAttributes(
						attribute.String("auth.credentials.issued_at.iso8601", issuedAt.Format(time.RFC3339)),
						attribute.Float64("auth.credentials.issued_at.elapsed_seconds", time.Since(issuedAt).Seconds()),
					)
				}

				ctx = context.WithValue(ctx, jwtmiddleware.ContextKey{}, nil)
				ctx = contextWithToken(ctx, token)
				next.ServeHTTP(w, r.WithContext(ctx))
			})).ServeHTTP(w, r)
		})
	}, nil
}
