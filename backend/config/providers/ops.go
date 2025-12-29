package providers

import (
	"fmt"
	"iter"
	"net/url"
	"slices"
	"strings"

	"careco/backend/authz"

	"github.com/aereal/coll"
)

func parseAudience(s string) (authz.Audience, error) {
	return slices.Collect(parseCommaSeparatedList[string](s)), nil
}

func parseAllowedSubjects(s string) (*coll.Set[authz.AllowedSubject], error) {
	return coll.NewSet(slices.Collect(parseCommaSeparatedList[authz.AllowedSubject](s))...), nil
}

func parseCommaSeparatedList[T ~string](s string) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x := range strings.SplitSeq(s, ",") {
			if !yield(T(x)) {
				return
			}
		}
	}
}

func parseIssuer(s string) (*authz.Issuer, error) {
	parsed, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("net/url.Parse: %w", err)
	}
	return (*authz.Issuer)(parsed), nil
}
