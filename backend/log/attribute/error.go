package attribute

import "log/slog"

const KeyError = "error"

type ErrorValue struct{ Err error }

func Error(err error) slog.Attr { return slog.Any(KeyError, &ErrorValue{Err: err}) }
