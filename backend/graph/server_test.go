package graph_test

import (
	"bytes"
	stdcmp "cmp"
	_ "embed"
	"encoding/json"
	"iter"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"

	"careco/backend/graph/test"

	"github.com/99designs/gqlgen/graphql"
	"gopkg.in/yaml.v3"
)

func TestServer(t *testing.T) {
	t.Parallel()

	root := new(casesRoot)
	if err := yaml.NewDecoder(bytes.NewReader(caseDefs)).Decode(root); err != nil {
		t.Fatal(err)
	}
	for tc := range root.All() {
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(tc.Request); err != nil {
				t.Fatal(err)
			}

			h := test.BuildHandler()
			srv := httptest.NewServer(h)
			t.Cleanup(srv.Close)

			req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, srv.URL, buf)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			gotResp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer gotResp.Body.Close()
			assertsResponse(t, tc.Response, gotResp)
		})
	}
}

//go:embed tests.yml
var caseDefs []byte

type gqlResponse struct {
	Errors []map[string]any `json:"errors,omitempty"`
	Data   map[string]any   `json:"data"`
}

type casesRoot struct {
	Cases map[string]*communication `json:"cases"`
}

type communication struct {
	Request  *graphql.RawParams `json:"request"`
	Response *gqlResponse       `json:"response"`
}

type testCase struct {
	*communication

	caseName string
}

func (r *casesRoot) All() iter.Seq[*testCase] {
	return func(yield func(*testCase) bool) {
		for _, caseName := range slices.SortedStableFunc(maps.Keys(r.Cases), stdcmp.Compare) {
			tc := &testCase{
				caseName:      caseName,
				communication: r.Cases[caseName],
			}
			if !yield(tc) {
				return
			}
		}
	}
}

func assertsResponse(t *testing.T, want *gqlResponse, got *http.Response) {
	t.Helper()
	jv, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	mv := make(map[string]any)
	if err := json.Unmarshal(jv, &mv); err != nil {
		t.Fatal(err)
	}
	exp := &responseExpectation{
		Status: http.StatusOK,
		Body:   mv,
	}
	if diff := cmp.Diff(exp, transformResponse(got)); diff != "" {
		t.Errorf("response (-want, +got):\n%s", diff)
	}
}

type responseExpectation struct {
	Status int
	Body   map[string]any
}

func transformResponse(hr *http.Response) *responseExpectation {
	defer hr.Body.Close()
	m := make(map[string]any)
	_ = json.NewDecoder(hr.Body).Decode(&m)
	return &responseExpectation{
		Status: hr.StatusCode,
		Body:   m,
	}
}
