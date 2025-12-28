package authz

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
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
