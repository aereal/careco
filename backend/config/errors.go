package config

import (
	"errors"
	"fmt"
)

func isMissingEnvError(err error) (*MissingEnvError, bool) {
	meErr := new(MissingEnvError)
	return meErr, errors.As(err, &meErr)
}

type MissingEnvError struct {
	Name string
}

var _ error = (*MissingEnvError)(nil)

func (e *MissingEnvError) Error() string {
	return fmt.Sprintf("missing environment variable: %q", e.Name)
}

func (e *MissingEnvError) Is(target error) bool {
	meErr, ok := isMissingEnvError(target)
	if !ok {
		return false
	}
	return e.Name == meErr.Name
}
