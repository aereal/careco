package providers

import (
	"fmt"
	"net/url"
	"strings"

	"careco/backend/authz"
)

func parseAudience(s string) (authz.Audience, error) {
	return strings.Split(s, ","), nil
}

func parseIssuer(s string) (*authz.Issuer, error) {
	parsed, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("net/url.Parse: %w", err)
	}
	return (*authz.Issuer)(parsed), nil
}
