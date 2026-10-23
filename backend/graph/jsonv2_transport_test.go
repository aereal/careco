package graph_test

import (
	"context"
	jsonv1 "encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"careco/backend/graph"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/google/go-cmp/cmp"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestTransport(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		accept                  []string
		wantResponseContentType string
	}{
		{
			accept:                  []string{"application/json"},
			wantResponseContentType: "application/json",
		},
		{
			accept:                  []string{"application/graphql-response+json"},
			wantResponseContentType: "application/graphql-response+json",
		},
		{
			accept:                  []string{"application/json", "application/graphql-response+json"},
			wantResponseContentType: "application/json",
		},
		{
			accept:                  []string{"*/*"},
			wantResponseContentType: "application/graphql-response+json",
		},
	}
	for _, tc := range testCases {
		t.Run(strings.Join(tc.accept, ","), func(t *testing.T) {
			t.Parallel()

			h := newServer()
			h.AddTransport(graph.POSTTransport())

			resp := doRequest(t.Context(), h, requestParams{
				method:      http.MethodPost,
				target:      "/graphql",
				body:        `{"query":"query getUsers($num: Int!){users(num:$num){name age}}","variables":{"num":1}}`,
				accept:      tc.accept,
				contentType: "application/json",
			})

			if resp.Code != http.StatusOK {
				t.Errorf("Code: %d", resp.Code)
			}
			if got := resp.Header().Get("Content-Type"); got != tc.wantResponseContentType {
				t.Errorf("Content-Type: %s", got)
			}
		})
	}
}

func TestTransport_decodeCompatibility(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "integer",
			body: `{"query":"query getUsers($num: Int!){users(num:$num){name age}}","variables":{"num":1}}`,
		},
		{
			name: "integer beyond float64 precision",
			body: `{"query":"query getUsers($num: Int!){users(num:$num){name age}}","variables":{"num":9007199254740993}}`,
		},
		{
			name: "float and exponent",
			body: `{"query":"query { users(num: 1) { name } }","variables":{"f":1.5,"e":1e3,"neg":-0}}`,
		},
		{
			name: "nested numbers",
			body: `{"query":"query { users(num: 1) { name } }","variables":{"obj":{"n":1,"arr":[1,2.5,{"deep":3}]}}}`,
		},
		{
			name: "numbers in extensions",
			body: `{"query":"query { users(num: 1) { name } }","extensions":{"persistedQuery":{"version":1}}}`,
		},
		{
			name: "non-number values",
			body: `{"query":"query { users(num: 1) { name } }","variables":{"s":"1","b":true,"n":null}}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			want := captureRawParams(t, transport.POST{}, tc.body)
			got := captureRawParams(t, graph.POSTTransport(), tc.body)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("decoded parameters mismatch (-transport.POST +PostJSONV2Transport):\n%s", diff)
			}
		})
	}
}

func TestTransport_decodeNumberAsJSONNumber(t *testing.T) {
	t.Parallel()

	got := captureRawParams(t, graph.POSTTransport(),
		`{"query":"query getUsers($num: Int!){users(num:$num){name age}}","variables":{"num":9007199254740993}}`)
	num, ok := got.Variables["num"].(jsonv1.Number)
	if !ok {
		t.Fatalf("variables.num: want json.Number, got %T", got.Variables["num"])
	}
	if num.String() != "9007199254740993" {
		t.Errorf("variables.num: %s", num)
	}
}

type requestParams struct {
	method      string
	target      string
	body        string
	accept      []string
	contentType string
}

func doRequest(
	ctx context.Context,
	handler http.Handler,
	params requestParams,
) *httptest.ResponseRecorder {
	r := httptest.NewRequestWithContext(ctx, params.method, params.target, strings.NewReader(params.body))
	for _, v := range params.accept {
		r.Header.Add("Accept", v)
	}
	r.Header.Set("Content-Type", params.contentType)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
	return w
}

func newServer() *handler.Server {
	schema := gqlparser.MustLoadSchema(&ast.Source{
		Input: `
		type User {
			name: String!
			age: Int!
		}

		type Query {
			users(num: Int!): [User!]!
		}
		`})
	return handler.New(&graphql.ExecutableSchemaMock{
		ExecFunc: func(ctx context.Context) graphql.ResponseHandler {
			opCtx := graphql.GetOperationContext(ctx)
			if opCtx.Operation.Operation != ast.Query {
				return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
			}
			return func(ctx context.Context) *graphql.Response {
				resp, err := graphql.
					GetOperationContext(ctx).
					ResolverMiddleware(ctx, graphql.Resolver(func(ctx context.Context) (res any, err error) {
						return &graphql.Response{
							Data: []byte(`{"users":[{"name":"Yuno","age":17}]}`),
						}, nil
					}))
				if err != nil {
					return graphql.ErrorResponse(ctx, "resolver error: %+v", err)
				}
				return resp.(*graphql.Response) //nolint:forcetypeassert,errcheck
			}
		},
		SchemaFunc: func() *ast.Schema { return schema },
		ComplexityFunc: func(_ context.Context, _, _ string, _ int, _ map[string]any) (int, bool) {
			return 0, false
		},
	})
}

type capturedParams struct {
	Variables  map[string]any
	Extensions map[string]any
}

type rawParamsCapturer struct {
	captured *capturedParams
}

var (
	_ graphql.HandlerExtension          = (*rawParamsCapturer)(nil)
	_ graphql.OperationParameterMutator = (*rawParamsCapturer)(nil)
)

func (*rawParamsCapturer) ExtensionName() string { return "rawParamsCapturer" }

func (*rawParamsCapturer) Validate(graphql.ExecutableSchema) error { return nil }

func (c *rawParamsCapturer) MutateOperationParameters(_ context.Context, params *graphql.RawParams) *gqlerror.Error {
	// RawParams is pooled and reset by the transport, so keep the maps only.
	c.captured = &capturedParams{Variables: params.Variables, Extensions: params.Extensions}
	return nil
}

func captureRawParams(t *testing.T, tr graphql.Transport, body string) *capturedParams {
	t.Helper()

	capturer := &rawParamsCapturer{}
	h := newServer()
	h.AddTransport(tr)
	h.Use(capturer)
	resp := doRequest(t.Context(), h, requestParams{
		method:      http.MethodPost,
		target:      "/graphql",
		body:        body,
		accept:      []string{"application/json"},
		contentType: "application/json",
	})
	if capturer.captured == nil {
		t.Fatalf("parameters are not captured: code=%d body=%s", resp.Code, resp.Body)
	}
	return capturer.captured
}
