package authz_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"careco/backend/authz"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	signingPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: signingPrivateKey},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", "k1"),
	)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		audience []string
	}{
		{
			name:     "ok",
			audience: []string{"http://rs.example"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issuerHandler := newIssuerServer(signingPrivateKey)
			issuerSrv := httptest.NewServer(issuerHandler)
			t.Cleanup(issuerSrv.Close)
			issuerURL, err := url.Parse(issuerSrv.URL)
			if err != nil {
				t.Fatal(err)
			}
			issuerHandler.origin = issuerURL.String()
			mw, err := authz.ProvideAuth0Middleware((*authz.Issuer)(issuerURL), tc.audience, http.DefaultClient)
			if err != nil {
				t.Fatal(err)
			}
			appSrv := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
			t.Cleanup(appSrv.Close)

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, appSrv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			issuedAt := time.Now()
			tok, err := jwt.Signed(signer).
				Claims(validator.RegisteredClaims{
					Issuer:   issuerURL.String(),
					Subject:  "1234567890",
					Audience: tc.audience,
					IssuedAt: issuedAt.Unix(),
					Expiry:   issuedAt.Add(time.Hour * 24).Unix(),
				}).
				CompactSerialize()
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("authorization", "Bearer "+tok)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("response.status=%d body=%s", resp.StatusCode, string(body))
			}
		})
	}

}

func newIssuerServer(privKey *rsa.PrivateKey) *issuerServer {
	s := &issuerServer{mux: http.NewServeMux()}
	s.mux.Handle("GET /.well-known/openid-configuration", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"jwks_uri": s.origin + "/jwks.json"})
	}))
	s.mux.Handle("GET /jwks.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{
				{
					Key:       &privKey.PublicKey,
					KeyID:     "k1",
					Algorithm: string(jose.RS256),
					Use:       "sig",
				},
			},
		})
	}))
	return s
}

type issuerServer struct {
	origin string
	mux    *http.ServeMux
}

func (s *issuerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
