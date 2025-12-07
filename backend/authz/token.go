package authz

import (
	"context"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
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

type auth0Token struct {
	*validator.ValidatedClaims
}

var _ Token = (*auth0Token)(nil)

func (t *auth0Token) Audience() ([]string, bool) {
	aud := t.RegisteredClaims.Audience
	return aud, len(aud) > 0
}

func (t *auth0Token) Expiration() (time.Time, bool) {
	epoch := t.RegisteredClaims.Expiry
	exp := time.Unix(epoch, 0)
	return exp, true
}

func (t *auth0Token) IssuedAt() (time.Time, bool) {
	epoch := t.RegisteredClaims.IssuedAt
	iat := time.Unix(epoch, 0)
	return iat, true
}

func (t *auth0Token) Issuer() (string, bool) {
	issuer := t.RegisteredClaims.Issuer
	return issuer, issuer != ""
}

func (t *auth0Token) JwtID() (string, bool) {
	id := t.RegisteredClaims.ID
	return id, id != ""
}

func (t *auth0Token) NotBefore() (time.Time, bool) {
	epoch := t.RegisteredClaims.NotBefore
	nb := time.Unix(epoch, 0)
	return nb, true
}

func (t *auth0Token) Subject() (string, bool) {
	sub := t.RegisteredClaims.Subject
	return sub, sub != ""
}
