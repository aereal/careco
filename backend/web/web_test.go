package web_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"careco/backend/web"
)

func TestWithCORS(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		method      string
		origin      string
		headers     []string
		wantSuccess bool
	}{
		{
			name:        "ok",
			method:      http.MethodPost,
			origin:      "http://localhost:1234",
			headers:     []string{"origin"},
			wantSuccess: true,
		},
		{
			name:        "DELETE is not allowed",
			method:      http.MethodDelete,
			origin:      "http://localhost:1234",
			headers:     []string{"origin"},
			wantSuccess: false,
		},
		{
			name:        "Vercel production",
			method:      http.MethodPost,
			origin:      "https://careco-nine.vercel.app",
			headers:     []string{"origin"},
			wantSuccess: true,
		},
		{
			name:        "Vercel preview production",
			method:      http.MethodPost,
			origin:      "https://careco-0xdeadbeaf-aereals-projects.vercel.app",
			headers:     []string{"origin"},
			wantSuccess: true,
		},
		{
			name:        "other Vercel app",
			method:      http.MethodPost,
			origin:      "https://hoge.vercel.app",
			headers:     []string{"origin"},
			wantSuccess: false,
		},
		{
			name:        "same suffix but other project",
			method:      http.MethodPost,
			origin:      "https://other-0xdeadbeaf-aereals-projects.vercel.app",
			headers:     []string{"origin"},
			wantSuccess: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			srv :=
				httptest.NewServer(
					web.WithCors(slog.New(slog.DiscardHandler))(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }),
					))
			t.Cleanup(srv.Close)

			req, err := newPreflightRequest(t.Context(), srv.URL, tc.origin, tc.method, tc.headers)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			preflightErr := isPreflightSuccessful(resp)
			if tc.wantSuccess != (preflightErr == nil) {
				t.Errorf("wantSuccess=%v response.status=%d error=%s", tc.wantSuccess, resp.StatusCode, preflightErr)
			}
		})
	}
}

func isPreflightSuccessful(resp *http.Response) error {
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode) //nolint:err113
	}
	if resp.Header.Get("access-control-allow-origin") == "" {
		return errors.New("the origin is not allowed") //nolint:err113
	}
	return nil
}

func newPreflightRequest(ctx context.Context, url string, origin string, method string, headers []string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodOptions, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("origin", origin)
	req.Header.Set("access-control-request-method", method)
	req.Header.Set("access-control-request-headers", strings.Join(headers, ","))
	return req, nil
}
