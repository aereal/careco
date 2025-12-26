package authz

import "fmt"

type NoSignatureAlgorithmError struct{}

var _ error = NoSignatureAlgorithmError{}

func (NoSignatureAlgorithmError) Error() string { return "no signature algorithm" }

var ErrMissingToken MissingTokenError

type MissingTokenError struct{}

var _ error = MissingTokenError{}

func (MissingTokenError) Error() string { return "missing token" }

type UnexpectedAuthSchemeError struct {
	AuthScheme string
}

var _ error = (*UnexpectedAuthSchemeError)(nil)

func (e *UnexpectedAuthSchemeError) Error() string {
	return fmt.Sprintf("unexpected auth scheme: %q", e.AuthScheme)
}

func (e *UnexpectedAuthSchemeError) Is(other error) bool {
	rhs, ok := isKindOfUnexpectedAuthScheme(other)
	if !ok {
		return false
	}
	return e.AuthScheme == rhs.AuthScheme
}

func isKindOfUnexpectedAuthScheme(err error) (*UnexpectedAuthSchemeError, bool) {
	if err == nil {
		return nil, false
	}
	uasErr, ok := err.(*UnexpectedAuthSchemeError)
	return uasErr, ok
}

func isaUnexpectedAuthScheme(err error) bool {
	_, ok := isKindOfUnexpectedAuthScheme(err)
	return ok
}

type HTTPResponseError struct {
	Status int
}

var _ error = (*HTTPResponseError)(nil)

func (e *HTTPResponseError) Error() string { return fmt.Sprintf("http response: status: %d", e.Status) }

func (e *HTTPResponseError) Is(other error) bool {
	if other == nil {
		return false
	}
	httpErr, ok := other.(*HTTPResponseError)
	if !ok {
		return false
	}
	return e.Status == httpErr.Status
}
