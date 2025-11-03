package config

import (
	"encoding"
	"log/slog"
	"strconv"
)

func Unmarshal[T any, PT interface {
	encoding.TextUnmarshaler
	*T
}](s string) (T, error) {
	var t T
	if err := PT(&t).UnmarshalText([]byte(s)); err != nil {
		return *new(T), err
	}
	return t, nil
}

var _ CastFunc[string, slog.Level] = Unmarshal

func StringAs[T ~string](s string) (T, error) { return T(s), nil }

var _ CastFunc[string, string] = StringAs

func ParseBool[T ~bool](s string) (T, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return *new(T), err
	}
	return T(b), nil
}

type CastFunc[I any, O any] func(I) (O, error)

func Cast[I any, O any](convert CastFunc[I, O]) Middleware[I, O] {
	return func(r Retriever[I]) Retriever[O] {
		return RetrieverFunc[O](func(name string) (O, error) {
			incoming, err := r.Retrieve(name)
			if err != nil {
				return *new(O), err
			}
			return convert(incoming)
		})
	}
}

func WithDefaultValue[T any](dv T) func(Retriever[T]) Yielder[T] {
	return func(r Retriever[T]) Yielder[T] {
		return YielderFunc[T](func(name string) T {
			v, err := r.Retrieve(name)
			if _, ok := isMissingEnvError(err); ok {
				return dv
			}
			return v
		})
	}
}

func EnvSource(e *Environment) Retriever[string] {
	return RetrieverFunc[string](func(name string) (string, error) { return e.lookup(name) })
}

func Retrieve[T any](name string, r Retriever[T]) (T, error) {
	v, err := r.Retrieve(name)
	if err != nil {
		return *new(T), err
	}
	slog.Info("retrieve env", slog.String("name", name), slog.Any("value", v))
	return v, nil
}

func Yield[T any](name string, y Yielder[T]) T {
	v := y.Yield(name)
	slog.Info("yield env value", slog.String("name", name), slog.Any("value", v))
	return v
}

type Middleware[I any, O any] func(Retriever[I]) Retriever[O]

type Retriever[T any] interface {
	Retrieve(name string) (T, error)
}

type RetrieverFunc[T any] func(name string) (T, error)

var _ Retriever[any] = (RetrieverFunc[any])(nil)

func (f RetrieverFunc[T]) Retrieve(name string) (T, error) { return f(name) }

type Yielder[T any] interface {
	Retriever[T]
	Yield(name string) T
}

type YielderFunc[T any] func(name string) T

var _ Yielder[any] = (YielderFunc[any])(nil)

func (f YielderFunc[T]) Retrieve(name string) (T, error) { return f(name), nil }

func (f YielderFunc[T]) Yield(name string) T { return f(name) }
