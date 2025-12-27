package authz

import (
	"errors"
	"iter"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"

	"careco/backend/log/attribute"

	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Issuer url.URL

func (i *Issuer) cloneURL() *url.URL {
	ret := &url.URL{
		Scheme:      i.Scheme,
		Opaque:      i.Opaque,
		Host:        i.Host,
		Path:        i.Path,
		RawPath:     i.RawPath,
		OmitHost:    i.OmitHost,
		ForceQuery:  i.ForceQuery,
		RawQuery:    i.RawQuery,
		Fragment:    i.Fragment,
		RawFragment: i.RawFragment,
	}
	if passwd, ok := i.User.Password(); ok {
		ret.User = url.UserPassword(i.User.Username(), passwd)
	} else {
		ret.User = url.User(i.User.Username())
	}
	return ret
}

func (i *Issuer) appendingPath(trail string) *url.URL {
	ret := i.cloneURL()
	ret.Path = path.Join(ret.Path, trail)
	return ret
}

func (i *Issuer) oidcConfiguration() *url.URL {
	return i.appendingPath(".well-known/openid-configuration")
}

type Audience []string

func (a Audience) jwtParseOptions() iter.Seq[jwt.ParseOption] {
	return func(yield func(jwt.ParseOption) bool) {
		for _, aud := range a {
			if !yield(jwt.WithAudience(aud)) {
				return
			}
		}
	}
}

type Middleware func(next http.Handler) http.Handler

func ProvideMiddleware(issuer *Issuer, audience Audience, client *http.Client) Middleware {
	kp := newOIDCKeyProvider(newOIDCKeySetFetcher(client, issuer))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			rawToken, err := getAuthzHeader(r.Header)
			if err != nil {
				handleError(w, r, err)
				return
			}

			parseOpts := []jwt.ParseOption{
				jwt.WithIssuer((*url.URL)(issuer).String()),
				jwt.WithKeyProvider(wrapKeyProvider(ctx, kp)),
			}
			parseOpts = slices.AppendSeq(parseOpts, audience.jwtParseOptions())
			tok, err := jwt.ParseString(rawToken, parseOpts...)
			if err != nil {
				handleError(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(contextWithToken(ctx, tok)))
		})
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	w.Header().Set("content-type", "application/json")

	switch {
	case errors.Is(err, ErrMissingToken):
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"no credentials provided"}`))
	case errors.Is(err, jwt.ParseError()) || isaUnexpectedAuthScheme(err):
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid credentials"}`))
	default:
		slog.WarnContext(r.Context(), "other kind of error caught", attribute.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"unexpected error happen"}`))
	}
}

func getAuthzHeader(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", ErrMissingToken
	}
	scheme, creds, ok := strings.Cut(val, " ")
	if !ok || strings.ToLower(scheme) != "bearer" {
		return "", &UnexpectedAuthSchemeError{AuthScheme: scheme}
	}
	return creds, nil
}
