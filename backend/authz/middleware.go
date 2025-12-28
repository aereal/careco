package authz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	validatorpkg "github.com/auth0/go-jwt-middleware/v2/validator"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Issuer url.URL

type Audience []string

type Middleware func(next http.Handler) http.Handler

func ProvideAuth0Middleware(issuer *Issuer, audience Audience, client *http.Client) (Middleware, error) {
	keyProvider := jwks.NewCachingProvider((*url.URL)(issuer), time.Minute*30, jwks.WithCustomClient(client))
	validator, err := validatorpkg.New(keyFuncWithTracing(keyProvider.KeyFunc), validatorpkg.RS256, (*url.URL)(issuer).String(), audience, validatorpkg.WithAllowedClockSkew(time.Minute))
	if err != nil {
		return nil, fmt.Errorf("validator.New: %w", err)
	}
	handleError := errorHandlerWithTracing(jwtmiddleware.DefaultErrorHandler)
	mw := jwtmiddleware.New(
		validateTokenWithTracing(validator.ValidateToken),
		jwtmiddleware.WithErrorHandler(handleError),
		jwtmiddleware.WithTokenExtractor(tokenExtractorWithTracing(jwtmiddleware.AuthHeaderTokenExtractor)),
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nested := mw.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				rawToken := ctx.Value(jwtmiddleware.ContextKey{})
				claims, ok := rawToken.(*validatorpkg.ValidatedClaims)
				if !ok {
					handleError(w, r, unexpectedClaimsType(rawToken))
					return
				}
				token := &auth0Token{claims}
				trace.SpanFromContext(ctx).SetAttributes(slices.Collect(token.Attrs())...)
				next.ServeHTTP(w, r.Clone(contextWithToken(ctx, token)))
			}))
			nested.ServeHTTP(w, r)
		})
	}, nil
}

func unexpectedClaimsType(token any) error {
	return fmt.Errorf("unexpected token type: %T", token) //nolint:err113
}

func tokenExtractorWithTracing(f jwtmiddleware.TokenExtractor) jwtmiddleware.TokenExtractor {
	return func(r *http.Request) (string, error) {
		ctx, span := getTracer(r.Context()).Start(r.Context(), "TokenExtractor", trace.WithSpanKind(trace.SpanKindServer))
		token, err := f(r.WithContext(ctx))
		closeSpan(span, err)
		return token, err
	}
}

func errorHandlerWithTracing(h jwtmiddleware.ErrorHandler) jwtmiddleware.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		h(w, r, err)
	}
}

func validateTokenWithTracing(validate jwtmiddleware.ValidateToken) jwtmiddleware.ValidateToken {
	return func(ctx context.Context, s string) (any, error) {
		ctx, span := getTracer(ctx).Start(ctx, "ValidateToken", trace.WithSpanKind(trace.SpanKindServer))
		validated, err := validate(ctx, s)
		closeSpan(span, err)
		return validated, err
	}
}

func keyFuncWithTracing(keyFunc func(context.Context) (any, error)) func(context.Context) (any, error) {
	return func(ctx context.Context) (any, error) {
		childCtx, span := getTracer(ctx).Start(ctx, "KeyFunc", trace.WithSpanKind(trace.SpanKindServer))
		key, err := keyFunc(childCtx)
		closeSpan(span, err)
		return key, err
	}
}

func closeSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End()
}

const tracerName = "careco/backend/authz.Middleware"

func getTracer(ctx context.Context) trace.Tracer {
	return trace.SpanFromContext(ctx).TracerProvider().Tracer(tracerName)
}
