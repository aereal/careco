package tests

import (
	"errors"
	"testing"
)

func AssertsErrors(t *testing.T, want, got error) bool {
	t.Helper()
	if !errors.Is(want, got) {
		t.Errorf("error:\n\twant: type=%T msg=%s\n\t got: type=%T msg=%s", want, want, got, got)
		return false
	}
	return true
}

type ErrorLiteral string //nolint:errname

func (err ErrorLiteral) Error() string { return string(err) }

func (err ErrorLiteral) Is(other error) bool {
	return string(err) == other.Error()
}
