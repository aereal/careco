package config

import (
	"os"
	"strings"
	"sync"
)

func ProvideEnvironment() *Environment {
	e := &Environment{dirty: map[string]string{}}
	for _, pair := range os.Environ() {
		k, v, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		e.dirty[k] = v
	}
	return e
}

type Environment struct {
	mux   sync.Mutex
	dirty map[string]string
}

func (e *Environment) lookup(name string) (string, error) {
	e.mux.Lock()
	defer e.mux.Unlock()
	v, ok := e.dirty[name]
	if !ok {
		return "", &MissingEnvError{Name: name}
	}
	return v, nil
}
