package attribute

import (
	"log/slog"
	"reflect"
)

const KeyError = "error"

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	typErr := reflect.TypeOf(err)
	attrs := make([]slog.Attr, 0, 3)
	attrs = append(attrs,
		slog.String("type", typErr.String()),
		slog.String("msg", err.Error()),
	)
	if pkg := typErr.PkgPath(); pkg != "" {
		attrs = append(attrs, slog.String("pkg", pkg))
	}
	return slog.GroupAttrs(KeyError, attrs...)
}
