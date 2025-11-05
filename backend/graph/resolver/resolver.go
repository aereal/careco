//go:generate ./generate.bash

package resolver

import "errors"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func ProvideResolver() *Resolver {
	return &Resolver{}
}

type Resolver struct{}

var errNotImplemented = errors.New("not implemented yet")
