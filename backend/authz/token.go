package authz

import (
	"context"
	"iter"
	"time"

	"careco/backend/o11y/attr"

	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.opentelemetry.io/otel/attribute"
)

type tokenCtxKey struct{}

func contextWithToken(ctx context.Context, tok Token) context.Context {
	return context.WithValue(ctx, tokenCtxKey{}, tok)
}

func TokenFromContext(ctx context.Context) (Token, bool) {
	tok, ok := ctx.Value(tokenCtxKey{}).(Token)
	return tok, ok
}

type Token interface {
	Audience() ([]string, bool)
	Expiration() (time.Time, bool)
	IssuedAt() (time.Time, bool)
	Issuer() (string, bool)
	JwtID() (string, bool)
	NotBefore() (time.Time, bool)
	Subject() (string, bool)
}

var _ Token = jwt.Token(nil)

var (
	keyCredentialsId                  attribute.Key = "creds.id"
	keyCredentialsIssuer              attribute.Key = "creds.issuer"
	keyCredentialsSubject             attribute.Key = "creds.subject"
	keyCredentialsAudience            attribute.Key = "creds.audience"
	keyCredentialsExpirationTimestamp attribute.Key = "creds.expires_at.timestamp" //nolint:gosec // false-positive
	keyCredentialsExpirationRemaining attribute.Key = "creds.expires_at.remaining" //nolint:gosec // false-positive
	keyCredentialsIssuedAtTimestamp   attribute.Key = "creds.issued_at.timestamp"  //nolint:gosec // false-positive
	keyCredentialsIssuedAtElapsed     attribute.Key = "creds.issued_at.elapsed"    //nolint:gosec // false-positive
)

func iterateAttrs(token Token) iter.Seq[attribute.KeyValue] {
	return func(yield func(attribute.KeyValue) bool) {
		if id, ok := token.JwtID(); ok {
			if !yield(keyCredentialsId.String(id)) {
				return
			}
		}
		if issuer, ok := token.Issuer(); ok {
			if !yield(keyCredentialsIssuer.String(issuer)) {
				return
			}
		}
		if sub, ok := token.Subject(); ok {
			if !yield(keyCredentialsSubject.String(sub)) {
				return
			}
		}
		if aud, ok := token.Audience(); ok {
			if !yield(keyCredentialsAudience.StringSlice(aud)) {
				return
			}
		}
		if expiresAt, ok := token.Expiration(); ok {
			if !yield(attribute.KeyValue{Key: keyCredentialsExpirationTimestamp, Value: attr.ISO8601Value(expiresAt)}) {
				return
			}
			if !yield(attribute.KeyValue{Key: keyCredentialsExpirationRemaining, Value: attr.DurationSecondsValue(time.Until(expiresAt))}) {
				return
			}
		}
		if issuedAt, ok := token.IssuedAt(); ok {
			if !yield(attribute.KeyValue{Key: keyCredentialsIssuedAtTimestamp, Value: attr.ISO8601Value(issuedAt)}) {
				return
			}
			if !yield(attribute.KeyValue{Key: keyCredentialsIssuedAtElapsed, Value: attr.DurationSecondsValue(time.Since(issuedAt))}) {
				return
			}
		}
	}
}
